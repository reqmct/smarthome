package models

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

type UserIDParam struct {
	// Идентификатор
	// Required: true
	// Minimum: 1
	UserID *int64 `uri:"user_id"`
}

// Validate validates this user
func (m *UserIDParam) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateUserID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *UserIDParam) validateUserID(_ strfmt.Registry) error {
	if err := validate.Required("user_id", "body", m.UserID); err != nil {
		return err
	}

	if err := validate.MinimumInt("user_id", "body", *m.UserID, 1, false); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this user
func (m *UserIDParam) ContextValidate(_ context.Context, _ strfmt.Registry) error {
	return nil
}
