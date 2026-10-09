package repository

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// * Directorio de migraciones goose, es la unica fuente de verdad del esquema:
// * asi el test no se desincroniza cuando cambia una migracion
const migrationsDir = "../db/migrations"

// * Pool compartido por todos los tests del paquete.
// * Es nil cuando se corre con -short (no hay contenedor)
var testPool *pgxpool.Pool

// * Levanta UN contenedor de Postgres para todo el paquete y aplica las
// * migraciones. Cada test limpia las tablas antes de correr
func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		os.Exit(m.Run())
	}

	ctx := context.Background()

	container, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("failed_to_start_postgres_container: %v", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		testcontainers.TerminateContainer(container) //nolint:errcheck // best-effort cleanup in TestMain
		log.Fatalf("failed_to_build_dsn: %v", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		testcontainers.TerminateContainer(container) //nolint:errcheck // best-effort cleanup in TestMain
		log.Fatalf("failed_to_create_pool: %v", err)
	}

	if err := applyMigrations(ctx, pool); err != nil {
		pool.Close()
		testcontainers.TerminateContainer(container) //nolint:errcheck // best-effort cleanup in TestMain
		log.Fatalf("failed_to_apply_migrations: %v", err)
	}

	testPool = pool

	code := m.Run()

	pool.Close()
	testcontainers.TerminateContainer(container) //nolint:errcheck // best-effort cleanup in TestMain

	os.Exit(code)
}

// * Aplica la seccion Up de cada migracion en orden de version
func applyMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read_migrations_dir: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			names = append(names, entry.Name())
		}
	}
	slices.Sort(names)

	if len(names) == 0 {
		return fmt.Errorf("no_migrations_found_in %q", migrationsDir)
	}

	for _, name := range names {
		raw, err := os.ReadFile(filepath.Join(migrationsDir, name))
		if err != nil {
			return fmt.Errorf("read_migration %s: %w", name, err)
		}

		up, err := gooseUpSection(string(raw))
		if err != nil {
			return fmt.Errorf("parse_migration %s: %w", name, err)
		}

		// * pgx usa el protocolo simple cuando no hay argumentos,
		// * por eso admite varios statements en un solo Exec
		if _, err := pool.Exec(ctx, up); err != nil {
			return fmt.Errorf("exec_migration %s: %w", name, err)
		}
	}

	return nil
}

// * Extrae lo que hay entre "-- +goose Up" y "-- +goose Down"
// * y descarta el resto de anotaciones (StatementBegin/End, NO TRANSACTION)
func gooseUpSection(content string) (string, error) {
	_, afterUp, found := strings.Cut(content, "-- +goose Up")
	if !found {
		return "", errors.New("missing_goose_up_marker")
	}

	up, _, _ := strings.Cut(afterUp, "-- +goose Down")

	lines := strings.Split(up, "\n")
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "-- +goose") {
			continue
		}
		kept = append(kept, line)
	}

	statements := strings.TrimSpace(strings.Join(kept, "\n"))
	if statements == "" {
		return "", errors.New("empty_goose_up_section")
	}

	return statements, nil
}

// * Deja la base limpia y devuelve el repositorio bajo prueba
func newTestRepo(t *testing.T) (*QuoteRepository, context.Context) {
	t.Helper()

	if testPool == nil {
		t.Skip("integration_test_skipped: se corrio con -short")
	}

	ctx := context.Background()
	if _, err := testPool.Exec(ctx, "TRUNCATE quote_items, quotes, leads CASCADE"); err != nil {
		t.Fatalf("Failed_to_truncate: %v", err)
	}

	return NewQuotesRepository(testPool), ctx
}

func testCurrency(t *testing.T) domain.Currency {
	t.Helper()

	c, err := domain.ParseCurrency("MXN")
	if err != nil {
		t.Fatalf("ParseCurrency_failed: %v", err)
	}

	return c
}

// * Inserta un lead valido para poder probar el FK quotes.lead_id
func seedLead(t *testing.T, ctx context.Context) uuid.UUID {
	t.Helper()

	id := uuid.New()
	const insert = `
		INSERT INTO leads (id, name, email, message, source, phone_number, business_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := testPool.Exec(ctx, insert,
		id, "Erick Gutierrez", "lead@test.com", "Necesito un backend", "landing_page", "5512345678", "Greco Labs",
	); err != nil {
		t.Fatalf("Failed_to_seed_lead: %v", err)
	}

	return id
}

/*
Cotizacion valida contra el esquema real:

  - subtotal/total en centavos, con total >= subtotal (quotes_total_consistency)
  - cada item cumple subtotal = quantity * unit_price (quote_items_subtotal_consistency)
  - la suma de los items coincide con el subtotal de la cotizacion
*/
func newTestQuote(t *testing.T) *domain.QuotesModel {
	t.Helper()

	id := uuid.New()
	now := time.Now().UTC().Truncate(time.Microsecond)

	items := []domain.QuoteItem{
		{
			ID:          uuid.New(),
			QuoteID:     id,
			Description: "Desarrollo backend",
			Quantity:    decimal.NewFromInt(2),
			Unit:        "hora",
			UnitPrice:   domain.FromCents(50000), // $500.00
			SubTotal:    domain.FromCents(100000),
			SortOrder:   0,
			CreatedAt:   now,
		},
		{
			ID:          uuid.New(),
			QuoteID:     id,
			Description: "Diseno UI",
			Quantity:    decimal.NewFromInt(1),
			Unit:        "mockup",
			UnitPrice:   domain.FromCents(25000), // $250.00
			SubTotal:    domain.FromCents(25000),
			SortOrder:   1,
			CreatedAt:   now,
		},
	}

	return &domain.QuotesModel{
		ID:                           id,
		LeadID:                       nil, // ** cotizacion manual
		QuoteNumber:                  "Q-2026-" + id.String()[:8],
		ClientName:                   "Erick Gutierrez",
		ClientBusiness:               new("Greco Labs"),
		ClientEmail:                  "test@test.com",
		ClientPhone:                  "5512345678",
		ProjectTitle:                 "Test Project",
		ProjectDescription:           "Test description",
		ProjectResume:                new("Resumen del proyecto"),
		ProjectEstimatedDurationDays: 30,
		FunctionalityExample:         new("https://example.com/demo"),
		Items:                        items,
		PaymentMethod:                domain.PaymentUpfront,
		Status:                       domain.QuoteStatusDraft,
		CreatedAt:                    now,
		Subtotal:                     domain.FromCents(125000),
		Total:                        domain.FromCents(125000),
		Currency:                     testCurrency(t),
		Notes:                        new("Sin notas"),
	}
}

func assertCents(t *testing.T, field string, got domain.Money, wantCents int64) {
	t.Helper()

	if got.ToCents() != wantCents {
		t.Errorf("Expected_%s: %v Got: %v", field, wantCents, got.ToCents())
	}
}

func TestCreateQuote(t *testing.T) {
	repo, ctx := newTestRepo(t)

	quote := newTestQuote(t)

	result, err := repo.CreateQuote(ctx, quote)
	if err != nil {
		t.Fatalf("CreateQuote_failed: %v", err)
	}

	if result.ID != quote.ID {
		t.Errorf("Expected_ID: %v Got: %v", quote.ID, result.ID)
	}

	if result.QuoteNumber != quote.QuoteNumber {
		t.Errorf("Expected_quote_number: %v Got: %v", quote.QuoteNumber, result.QuoteNumber)
	}

	if result.Status != domain.QuoteStatusDraft {
		t.Errorf("Expected_status: %v Got: %v", domain.QuoteStatusDraft, result.Status)
	}

	if result.LeadID != nil {
		t.Errorf("Expected_lead_id: nil Got: %v", *result.LeadID)
	}

	if result.Currency.String() != "MXN" {
		t.Errorf("Expected_currency: MXN Got: %v", result.Currency.String())
	}

	if result.ClientBusiness == nil || *result.ClientBusiness != "Greco Labs" {
		t.Errorf("Expected_client_business: Greco Labs Got: %v", result.ClientBusiness)
	}

	assertCents(t, "subtotal", result.Subtotal, 125000)
	assertCents(t, "total", result.Total, 125000)

	if !result.CreatedAt.Equal(quote.CreatedAt) {
		t.Errorf("Expected_created_at: %v Got: %v", quote.CreatedAt, result.CreatedAt)
	}

	if len(result.Items) != len(quote.Items) {
		t.Fatalf("Expected_items: %v Got: %v", len(quote.Items), len(result.Items))
	}

	for i, item := range result.Items {
		want := quote.Items[i]

		if item.ID != want.ID {
			t.Errorf("Expected_item_id: %v Got: %v", want.ID, item.ID)
		}
		if item.QuoteID != quote.ID {
			t.Errorf("Expected_item_quote_id: %v Got: %v", quote.ID, item.QuoteID)
		}
		if item.Description != want.Description {
			t.Errorf("Expected_item_description: %v Got: %v", want.Description, item.Description)
		}
		if !item.Quantity.Equal(want.Quantity) {
			t.Errorf("Expected_item_quantity: %v Got: %v", want.Quantity, item.Quantity)
		}
		if item.SortOrder != i {
			t.Errorf("Expected_item_sort_order: %v Got: %v", i, item.SortOrder)
		}
		assertCents(t, "item_unit_price", item.UnitPrice, want.UnitPrice.ToCents())
		assertCents(t, "item_subtotal", item.SubTotal, want.SubTotal.ToCents())
	}
}

func TestCreateQuote_FromLead(t *testing.T) {
	repo, ctx := newTestRepo(t)

	leadID := seedLead(t, ctx)

	quote := newTestQuote(t)
	quote.LeadID = &leadID

	result, err := repo.CreateQuote(ctx, quote)
	if err != nil {
		t.Fatalf("CreateQuote_failed: %v", err)
	}

	if result.LeadID == nil {
		t.Fatalf("Expected_lead_id: %v Got: nil", leadID)
	}

	if *result.LeadID != leadID {
		t.Errorf("Expected_lead_id: %v Got: %v", leadID, *result.LeadID)
	}
}

func TestCreateQuote_WithoutItems(t *testing.T) {
	repo, ctx := newTestRepo(t)

	quote := newTestQuote(t)
	quote.Items = nil
	quote.Subtotal = domain.FromCents(0)
	quote.Total = domain.FromCents(0)

	result, err := repo.CreateQuote(ctx, quote)
	if err != nil {
		t.Fatalf("CreateQuote_failed: %v", err)
	}

	if len(result.Items) != 0 {
		t.Errorf("Expected_items: 0 Got: %v", len(result.Items))
	}
}

// * Cada caso parte de una cotizacion valida y rompe UNA sola invariante
func TestCreateQuote_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T, ctx context.Context, repo *QuoteRepository)
		mutate  func(t *testing.T, ctx context.Context, quote *domain.QuotesModel)
		wantErr error
	}{
		{
			name:    "valid quote creation",
			wantErr: nil,
		},
		{
			name: "duplicate quote number should fail",
			setup: func(t *testing.T, ctx context.Context, repo *QuoteRepository) {
				existing := newTestQuote(t)
				existing.QuoteNumber = "Q-2026-DUP"
				if _, err := repo.CreateQuote(ctx, existing); err != nil {
					t.Fatalf("Setup_CreateQuote_failed: %v", err)
				}
			},
			mutate: func(_ *testing.T, _ context.Context, quote *domain.QuotesModel) {
				quote.QuoteNumber = "Q-2026-DUP"
			},
			wantErr: domain.ErrAlreadyExists,
		},
		{
			name: "unknown lead id should fail",
			mutate: func(_ *testing.T, _ context.Context, quote *domain.QuotesModel) {
				quote.LeadID = new(uuid.New())
			},
			wantErr: domain.ErrInvalidData,
		},
		{
			name: "invalid payment method should fail",
			mutate: func(_ *testing.T, _ context.Context, quote *domain.QuotesModel) {
				quote.PaymentMethod = domain.PayMethod("bitcoin")
			},
			wantErr: domain.ErrInvalidData,
		},
		{
			name: "invalid status should fail",
			mutate: func(_ *testing.T, _ context.Context, quote *domain.QuotesModel) {
				quote.Status = domain.QuoteStatus("archived")
			},
			wantErr: domain.ErrInvalidData,
		},
		{
			name: "negative subtotal should fail",
			mutate: func(_ *testing.T, _ context.Context, quote *domain.QuotesModel) {
				quote.Subtotal = domain.FromCents(-1)
				quote.Total = domain.FromCents(-1)
			},
			wantErr: domain.ErrInvalidData,
		},
		{
			name: "total lower than subtotal should fail",
			mutate: func(_ *testing.T, _ context.Context, quote *domain.QuotesModel) {
				quote.Total = domain.FromCents(quote.Subtotal.ToCents() - 1)
			},
			wantErr: domain.ErrInvalidData,
		},
		{
			name: "item subtotal mismatch should fail",
			mutate: func(_ *testing.T, _ context.Context, quote *domain.QuotesModel) {
				quote.Items[0].SubTotal = domain.FromCents(1)
			},
			wantErr: domain.ErrInvalidData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, ctx := newTestRepo(t)

			if tt.setup != nil {
				tt.setup(t, ctx, repo)
			}

			quote := newTestQuote(t)
			if tt.mutate != nil {
				tt.mutate(t, ctx, quote)
			}

			created, err := repo.CreateQuote(ctx, quote)

			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("Expected_no_error Got: %v", err)
				}
				if created == nil {
					t.Fatal("Expected_created_quote Got: nil")
				}
				return
			}

			if err == nil {
				t.Fatalf("Expected_error: %v Got: nil", tt.wantErr)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Expected_error: %v Got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestGetQuoteByID(t *testing.T) {
	repo, ctx := newTestRepo(t)

	quote := newTestQuote(t)
	if _, err := repo.CreateQuote(ctx, quote); err != nil {
		t.Fatalf("CreateQuote_failed: %v", err)
	}

	found, err := repo.GetQuoteByID(ctx, quote.ID)
	if err != nil {
		t.Fatalf("GetQuoteByID_failed: %v", err)
	}

	if found.ID != quote.ID {
		t.Errorf("Expected_ID: %v Got: %v", quote.ID, found.ID)
	}

	if found.QuoteNumber != quote.QuoteNumber {
		t.Errorf("Expected_quote_number: %v Got: %v", quote.QuoteNumber, found.QuoteNumber)
	}

	if len(found.Items) != len(quote.Items) {
		t.Fatalf("Expected_items: %v Got: %v", len(quote.Items), len(found.Items))
	}

	// * ListQuoteItemsByQuoteID ordena por sort_order
	for i, item := range found.Items {
		if item.SortOrder != i {
			t.Errorf("Expected_item_sort_order: %v Got: %v", i, item.SortOrder)
		}
		if item.Description != quote.Items[i].Description {
			t.Errorf("Expected_item_description: %v Got: %v", quote.Items[i].Description, item.Description)
		}
	}

	assertCents(t, "total", found.Total, quote.Total.ToCents())
}

func TestGetQuoteByID_NotFound(t *testing.T) {
	repo, ctx := newTestRepo(t)

	_, err := repo.GetQuoteByID(ctx, uuid.New())
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Expected_error: %v Got: %v", domain.ErrNotFound, err)
	}
}

func TestUpdateQuote(t *testing.T) {
	repo, ctx := newTestRepo(t)

	quote := newTestQuote(t)
	if _, err := repo.CreateQuote(ctx, quote); err != nil {
		t.Fatalf("CreateQuote_failed: %v", err)
	}

	// * UpdateQuote reemplaza por completo los items
	updated := newTestQuote(t)
	updated.ID = quote.ID
	updated.QuoteNumber = quote.QuoteNumber
	updated.ClientName = "Cliente Actualizado"
	updated.ProjectTitle = "Proyecto Actualizado"
	updated.PaymentMethod = domain.PaymentMilestone
	updated.Notes = new("Notas actualizadas")
	updated.Subtotal = domain.FromCents(300000)
	updated.Total = domain.FromCents(300000)
	updated.Items = []domain.QuoteItem{
		{
			ID:          uuid.New(),
			QuoteID:     quote.ID,
			Description: "Consultoria estrategica",
			Quantity:    decimal.NewFromInt(3),
			Unit:        "sesion",
			UnitPrice:   domain.FromCents(100000),
			SubTotal:    domain.FromCents(300000),
			CreatedAt:   time.Now().UTC().Truncate(time.Microsecond),
		},
	}

	result, err := repo.UpdateQuote(ctx, updated)
	if err != nil {
		t.Fatalf("UpdateQuote_failed: %v", err)
	}

	if result.ClientName != "Cliente Actualizado" {
		t.Errorf("Expected_client_name: Cliente Actualizado Got: %v", result.ClientName)
	}

	if result.ProjectTitle != "Proyecto Actualizado" {
		t.Errorf("Expected_project_title: Proyecto Actualizado Got: %v", result.ProjectTitle)
	}

	if result.PaymentMethod != domain.PaymentMilestone {
		t.Errorf("Expected_payment_method: %v Got: %v", domain.PaymentMilestone, result.PaymentMethod)
	}

	assertCents(t, "total", result.Total, 300000)

	if len(result.Items) != 1 {
		t.Fatalf("Expected_items: 1 Got: %v", len(result.Items))
	}

	if result.Items[0].Description != "Consultoria estrategica" {
		t.Errorf("Expected_item_description: Consultoria estrategica Got: %v", result.Items[0].Description)
	}

	// * los items viejos deben haber desaparecido
	persisted, err := repo.GetQuoteByID(ctx, quote.ID)
	if err != nil {
		t.Fatalf("GetQuoteByID_failed: %v", err)
	}

	if len(persisted.Items) != 1 {
		t.Fatalf("Expected_persisted_items: 1 Got: %v", len(persisted.Items))
	}

	if persisted.Items[0].ID != updated.Items[0].ID {
		t.Errorf("Expected_item_id: %v Got: %v", updated.Items[0].ID, persisted.Items[0].ID)
	}
}

func TestUpdateQuote_NotFound(t *testing.T) {
	repo, ctx := newTestRepo(t)

	quote := newTestQuote(t)

	_, err := repo.UpdateQuote(ctx, quote)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Expected_error: %v Got: %v", domain.ErrNotFound, err)
	}
}

// * UpdateQuote solo aplica sobre drafts (WHERE status = 'draft')
func TestUpdateQuote_NonDraftStatus(t *testing.T) {
	repo, ctx := newTestRepo(t)

	quote := newTestQuote(t)
	quote.Status = domain.QuoteStatusSent
	if _, err := repo.CreateQuote(ctx, quote); err != nil {
		t.Fatalf("CreateQuote_failed: %v", err)
	}

	quote.ClientName = "No deberia actualizarse"

	_, err := repo.UpdateQuote(ctx, quote)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Expected_error: %v Got: %v", domain.ErrNotFound, err)
	}
}
