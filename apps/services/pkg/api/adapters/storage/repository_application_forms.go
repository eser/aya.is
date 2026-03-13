package storage

import (
	"context"
	"database/sql"

	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

func (r *Repository) GetActiveApplicationForm(
	ctx context.Context,
	profileID string,
) (*profiles.ApplicationForm, error) {
	row, err := r.queries.GetActiveApplicationForm(ctx, GetActiveApplicationFormParams{
		ProfileID: profileID,
	})
	if err != nil {
		return nil, err
	}

	fields, err := r.ListApplicationFormFields(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	return &profiles.ApplicationForm{
		ID:                  row.ID,
		ProfileID:           row.ProfileID,
		PresetKey:           vars.ToStringPtr(row.PresetKey),
		IsActive:            row.IsActive,
		ResponsesVisibility: row.ResponsesVisibility,
		Fields:              fields,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           vars.ToTimePtr(row.UpdatedAt),
	}, nil
}

func (r *Repository) GetApplicationFormByProfileID(
	ctx context.Context,
	profileID string,
) (*profiles.ApplicationForm, string, error) {
	row, err := r.queries.GetApplicationFormByProfileID(ctx, GetApplicationFormByProfileIDParams{
		ProfileID: profileID,
	})
	if err != nil {
		return nil, "", err
	}

	fields, err := r.ListApplicationFormFields(ctx, row.ID)
	if err != nil {
		return nil, "", err
	}

	form := &profiles.ApplicationForm{
		ID:                  row.ID,
		ProfileID:           row.ProfileID,
		PresetKey:           vars.ToStringPtr(row.PresetKey),
		IsActive:            row.IsActive,
		ResponsesVisibility: row.ResponsesVisibility,
		Fields:              fields,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           vars.ToTimePtr(row.UpdatedAt),
	}

	return form, row.FeatureApplications, nil
}

func (r *Repository) CreateApplicationForm(
	ctx context.Context,
	formID string,
	profileID string,
	presetKey *string,
	responsesVisibility string,
) (*profiles.ApplicationForm, error) {
	var presetKeyParam sql.NullString
	if presetKey != nil {
		presetKeyParam = sql.NullString{String: *presetKey, Valid: true}
	}

	row, err := r.queries.CreateApplicationForm(ctx, CreateApplicationFormParams{
		ID:                  formID,
		ProfileID:           profileID,
		PresetKey:           presetKeyParam,
		ResponsesVisibility: responsesVisibility,
	})
	if err != nil {
		return nil, err
	}

	return &profiles.ApplicationForm{
		ID:                  row.ID,
		ProfileID:           row.ProfileID,
		PresetKey:           vars.ToStringPtr(row.PresetKey),
		IsActive:            row.IsActive,
		ResponsesVisibility: row.ResponsesVisibility,
		Fields:              []*profiles.ApplicationFormField{},
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           vars.ToTimePtr(row.UpdatedAt),
	}, nil
}

func (r *Repository) UpdateApplicationForm(
	ctx context.Context,
	formID string,
	presetKey *string,
	responsesVisibility string,
) error {
	var presetKeyParam sql.NullString
	if presetKey != nil {
		presetKeyParam = sql.NullString{String: *presetKey, Valid: true}
	}

	return r.queries.UpdateApplicationForm(ctx, UpdateApplicationFormParams{
		ID:                  formID,
		PresetKey:           presetKeyParam,
		ResponsesVisibility: responsesVisibility,
	})
}

func (r *Repository) DeactivateApplicationForms(
	ctx context.Context,
	profileID string,
) error {
	return r.queries.DeactivateApplicationForms(ctx, DeactivateApplicationFormsParams{
		ProfileID: profileID,
	})
}

func (r *Repository) ListApplicationFormFields(
	ctx context.Context,
	formID string,
) ([]*profiles.ApplicationFormField, error) {
	rows, err := r.queries.ListApplicationFormFields(ctx, ListApplicationFormFieldsParams{
		FormID: formID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*profiles.ApplicationFormField, 0, len(rows))
	for _, row := range rows {
		result = append(result, &profiles.ApplicationFormField{
			ID:          row.ID,
			FormID:      row.FormID,
			Label:       row.Label,
			FieldType:   row.FieldType,
			IsRequired:  row.IsRequired,
			SortOrder:   int(row.SortOrder),
			Placeholder: vars.ToStringPtr(row.Placeholder),
			CreatedAt:   row.CreatedAt,
		})
	}

	return result, nil
}

func (r *Repository) CreateApplicationFormField(
	ctx context.Context,
	fieldID string,
	formID string,
	label string,
	fieldType string,
	isRequired bool,
	sortOrder int,
	placeholder *string,
) (*profiles.ApplicationFormField, error) {
	var placeholderParam sql.NullString
	if placeholder != nil {
		placeholderParam = sql.NullString{String: *placeholder, Valid: true}
	}

	row, err := r.queries.CreateApplicationFormField(ctx, CreateApplicationFormFieldParams{
		ID:          fieldID,
		FormID:      formID,
		Label:       label,
		FieldType:   fieldType,
		IsRequired:  isRequired,
		SortOrder:   int32(sortOrder),
		Placeholder: placeholderParam,
	})
	if err != nil {
		return nil, err
	}

	return &profiles.ApplicationFormField{
		ID:          row.ID,
		FormID:      row.FormID,
		Label:       row.Label,
		FieldType:   row.FieldType,
		IsRequired:  row.IsRequired,
		SortOrder:   int(row.SortOrder),
		Placeholder: vars.ToStringPtr(row.Placeholder),
		CreatedAt:   row.CreatedAt,
	}, nil
}

func (r *Repository) DeleteApplicationFormFields(
	ctx context.Context,
	formID string,
) error {
	return r.queries.DeleteApplicationFormFields(ctx, DeleteApplicationFormFieldsParams{
		FormID: formID,
	})
}

func (r *Repository) SoftDeleteCandidate(
	ctx context.Context,
	candidateID string,
) error {
	_, err := r.queries.SoftDeleteCandidate(ctx, SoftDeleteCandidateParams{
		ID: candidateID,
	})

	return err
}

func (r *Repository) CreateCandidateResponse(
	ctx context.Context,
	responseID string,
	candidateID string,
	formFieldID string,
	value string,
) error {
	_, err := r.queries.CreateCandidateResponse(ctx, CreateCandidateResponseParams{
		ID:          responseID,
		CandidateID: candidateID,
		FormFieldID: formFieldID,
		Value:       value,
	})

	return err
}

func (r *Repository) ListCandidateResponses(
	ctx context.Context,
	candidateID string,
) ([]*profiles.CandidateFormResponse, error) {
	rows, err := r.queries.ListCandidateResponses(ctx, ListCandidateResponsesParams{
		CandidateID: candidateID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*profiles.CandidateFormResponse, 0, len(rows))
	for _, row := range rows {
		result = append(result, &profiles.CandidateFormResponse{
			ID:          row.ID,
			CandidateID: row.CandidateID,
			FormFieldID: row.FormFieldID,
			FieldLabel:  row.FieldLabel,
			FieldType:   row.FieldType,
			Value:       row.Value,
			SortOrder:   int(row.SortOrder),
			IsRequired:  row.IsRequired,
		})
	}

	return result, nil
}
