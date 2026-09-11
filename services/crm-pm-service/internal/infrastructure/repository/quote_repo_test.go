package repository

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/go-playground/locales/currency"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/moby/moby/api/types/network"
	"github.com/shopspring/decimal"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Inicia un contenedor PostgreSQL para tests
func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test_user",
			"POSTGRES_PASSWORD": "test_password",
			"POSTGRES_DB":       "test_db",
		},
		WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(host string, port network.Port) string {
			return "postgres://test_user:test_password@" + host + ":" + port.Port() + "/test_db?sslmode=disable"
		}).WithStartupTimeout(60 * time.Second).WithPollInterval(100 * time.Millisecond),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed_to_get_port: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed_to_get_host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed_to_get_port: %v", err)
	}

	dsn := "postgres://test_user:test_password@" + host + ":" + port.Port() + "/test_db?sslmode=disable"

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Failed_to_create_pool: %v", err)
	}

	createTableSQL := `
		CREATE TABLE quotes(
    		id                          UUID PRIMARY KEY,
    		lead_id                     UUID NOT NULL,
    		quote_number                VARCHAR(50) NOT NULL UNIQUE,

    		client_name                 VARCHAR(255) NOT NULL,
    		client_business             VARCHAR(255),
    		client_email                VARCHAR(255) NOT NULL,
    		client_phone                VARCHAR(20) NOT NULL,

    		-- Informacion del proyecto
    		project_title               VARCHAR(500) NOT NULL,
    		project_description         TEXT NOT NULL,
    		project_resume              TEXT,
    		project_estimated_duration  INTERVAL,

    		-- Detalles de la oferta
    		functionality_example       TEXT,
    		quote_items                 JSONB, -- quote_items
    		payment_method              VARCHAR(50) NOT NULL
    		                            CHECK (payment_method IN ('upfront', 'milestone', 'upon_delivery', 'split')),

    		status                      VARCHAR(30) NOT NULL DEFAULT 'draft' 
    		                            CHECK (status IN ('draft', 'sent', 'accepted', 'rejected', 'expired', 'converted')),
    		created_at                   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- se usa tambien para auditoria
    		valid_until                 TIMESTAMP WITH TIME ZONE,
    		sent_at                     TIMESTAMP WITH TIME ZONE,
    		accepted_at                 TIMESTAMP WITH TIME ZONE,
    		rejected_at                 TIMESTAMP WITH TIME ZONE,

    		subtotal_cents              BIGINT NOT NULL CHECK (subtotal_cents >= 0),
    		total_cents                 BIGINT NOT NULL CHECK (total_cents >= 0),

    		currency                    SMALLINT NOT NULL,
    		notes                       TEXT,

    		updated_at                  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    		CONSTRAINT quotes_valid_dates CHECK (valid_until > created_at),
    		CONSTRAINT quotes_total_consistency CHECK (
    		    total_cents = subtotal_cents
    		)
		);
	`

	if _, err := pool.Exec(ctx, createTableSQL); err != nil {
		t.Fatalf("Failed_to_create_table: %v", err)
	}

	cleanUp := func() {
		pool.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		container.Terminate(ctx)
	}

	return pool, cleanUp
}

func TestCreateQuote(t *testing.T) {
	pool, cleanUp := setupTestDB(t)
	defer cleanUp()

	repo := NewQuotesRepository(pool)
	ctx := context.Background()

	quote := &domain.QuotesModel{
		ID:                 uuid.New(),
		LeadID:             uuid.New(),
		QuoteNumber:        "Q-2026-001-test",
		ClientName:         "Erick Gutierrez",
		ClientEmail:        "test@test.com",
		ClientPhone:        "5512345678",
		ProjectTitle:       "Test Project",
		ProjectDescription: "Test description",
		PaymentMethod:      "upfront",
		Status:             "draft",
		CreatedAt:          time.Now(),
		Subtotal:           decimal.NewFromInt(100000),
		Total:              decimal.NewFromInt(100000),
		Currency:           currency.MXN,
	}

	result, err := repo.CreateQuote(ctx, quote)
	if err != nil {
		t.Fatalf("CreateQuote_failed: %v", err)
	}

	if result.ID != quote.ID {
		t.Errorf("Expected_ID: %v, Got: %v", quote.ID, result.ID)
	}

	if result.QuoteNumber != "Q-2026-001-test" {
		t.Errorf("Expected_quote_number: %v, Got: %v", quote.QuoteNumber, result.QuoteNumber)
	}
}

func TestCreateQuote_TableDriven(t *testing.T) {
	pool, cleanUp := setupTestDB(t)
	defer cleanUp()

	repo := NewQuotesRepository(pool)
	ctx := context.Background()

	tests := []struct {
		name    string
		quote   *domain.QuotesModel
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid quote creation",
			quote: &domain.QuotesModel{
				ID:                 uuid.New(),
				LeadID:             uuid.New(),
				QuoteNumber:        "Q-2026-001",
				ClientName:         "ACME Corp",
				ClientEmail:        "contact@acme.com",
				ProjectTitle:       "Backend API",
				ProjectDescription: "REST API",
				PaymentMethod:      "upfront",
				Status:             "draft",
				CreatedAt:          time.Now(),
				ValidUntil:         time.Now().AddDate(0, 0, 30),
				Subtotal:           decimal.NewFromInt(500000),
				Total:              decimal.NewFromInt(500000),
				Currency:           currency.MXN,
			},
			wantErr: false,
		},
		{
			name: "duplicate quote number should fail",
			quote: &domain.QuotesModel{
				ID:                 uuid.New(),
				LeadID:             uuid.New(),
				QuoteNumber:        "Q-2026-001", // ← Mismo número
				ClientName:         "Another Client",
				ClientEmail:        "another@example.com",
				ProjectTitle:       "Another Project",
				ProjectDescription: "Another description",
				PaymentMethod:      "split",
				Status:             "draft",
				CreatedAt:          time.Now(),
				ValidUntil:         time.Now().AddDate(0, 0, 30),
				Subtotal:           decimal.NewFromInt(200000),
				Total:              decimal.NewFromInt(200000),
				Currency:           currency.MXN,
			},
			wantErr: true,
			errMsg:  "duplicate",
		},
		{
			name: "valid until before issued at should fail",
			quote: &domain.QuotesModel{
				ID:                 uuid.New(),
				LeadID:             uuid.New(),
				QuoteNumber:        "Q-2026-002",
				ClientName:         "Test Client",
				ClientEmail:        "test@example.com",
				ProjectTitle:       "Test",
				ProjectDescription: "Test",
				PaymentMethod:      "upfront",
				Status:             "draft",
				CreatedAt:          time.Now(),
				ValidUntil:         time.Now().Add(-24 * time.Hour), // ← Antes de issued_at
				Subtotal:           decimal.NewFromInt(100000),
				Total:              decimal.NewFromInt(100000),
				Currency:           currency.MXN,
			},
			wantErr: true,
			errMsg:  "valid_dates",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.CreateQuote(ctx, tt.quote)
			if err != nil {
				t.Errorf("wantErr: %v, gotErr: %v", tt.wantErr, err)
			}

			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected_error_containing: %q Got: %v", tt.errMsg, err)
				}
			}
		})
	}
}
