package repository

import (
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/repository/sqlgen"
)

// * Convierte el modelo de dominio a los parametros esperados
// * por la query generada
func leadDomainToGen(lead *domain.LeadsModel) sqlgen.CreateLeadParams {
	return sqlgen.CreateLeadParams{
		ID:           lead.ID,
		Name:         lead.Name,
		BusinessName: lead.BusinessName,
		PhoneNumber:  lead.PhoneNumber,
		Email:        lead.Email,
		Message:      lead.Message,
		DataQuality:  string(lead.DataQuality),
		QualityScore: lead.QualityScore,
		Status:       string(lead.Status),
		Source:       string(lead.Source),
	}
}

// * Convierte el modelo generado por sqlc al modelo de dominio
// * Es el unico lugar donde sqlgen.Lead cruza hacia domain
func genToLeadDomain(row sqlgen.Lead) *domain.LeadsModel {
	return &domain.LeadsModel{
		ID:           row.ID,
		Name:         row.Name,
		BusinessName: row.BusinessName,
		PhoneNumber:  row.PhoneNumber,
		Email:        row.Email,
		Message:      row.Message,
		DataQuality:  domain.Quality(row.DataQuality),
		QualityScore: row.QualityScore,
		Status:       domain.LeadStatus(row.Status),
		Source:       domain.LeadSource(row.Source),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
