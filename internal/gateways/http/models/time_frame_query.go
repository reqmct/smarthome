package models

import (
	"context"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

type TimeFraneQuery struct {
	Start *time.Time `form:"start_date"`
	End   *time.Time `form:"end_date"`
}

func (m *TimeFraneQuery) Validate(_ strfmt.Registry) error {
	if err := validate.Required("start_date", "query", m.Start); err != nil {
		return err
	}
	if err := validate.Required("end_date", "query", m.End); err != nil {
		return err
	}
	return nil
}

func (m *TimeFraneQuery) ContextValidate(_ context.Context, _ strfmt.Registry) error {
	return nil
}
