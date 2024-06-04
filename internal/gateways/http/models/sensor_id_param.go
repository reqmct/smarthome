package models

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

type SensorIDParam struct {
	// Идентификатор датчика
	// Required: true
	// Minimum: 1
	SensorID *int64 `uri:"sensor_id"`
}

// Validate validates this sensor
func (m *SensorIDParam) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSensorID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}

	return nil
}

func (m *SensorIDParam) validateSensorID(_ strfmt.Registry) error {
	if err := validate.Required("sensor_id", "uri", m.SensorID); err != nil {
		return err
	}

	if err := validate.MinimumInt("sensor_id", "uri", *m.SensorID, 1, false); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this sensor
func (m *SensorIDParam) ContextValidate(_ context.Context, _ strfmt.Registry) error {
	return nil
}
