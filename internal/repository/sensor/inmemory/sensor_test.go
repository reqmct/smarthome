package inmemory

import (
	"context"
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecase"
	"math/rand/v2"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSensorRepository_SaveSensor(t *testing.T) {
	t.Run("err, sensor is nil", func(t *testing.T) {
		sr := NewSensorRepository()
		err := sr.SaveSensor(context.Background(), nil)
		assert.Error(t, err)
	})

	t.Run("fail, ctx cancelled", func(t *testing.T) {
		sr := NewSensorRepository()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := sr.SaveSensor(ctx, &domain.Sensor{})
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("fail, ctx deadline exceeded", func(t *testing.T) {
		sr := NewSensorRepository()
		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()

		err := sr.SaveSensor(ctx, &domain.Sensor{})
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("ok, save and get one", func(t *testing.T) {
		sr := NewSensorRepository()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sensor := &domain.Sensor{
			SerialNumber: "12345678",
			Type:         domain.SensorTypeContactClosure,
			CurrentState: 0,
			Description:  "sensor description",
			IsActive:     true,
		}

		err := sr.SaveSensor(ctx, sensor)
		assert.NoError(t, err)

		actualSensor, err := sr.GetSensorByID(ctx, sensor.ID)
		assert.NoError(t, err)
		assert.NotNil(t, actualSensor)
		assert.Equal(t, sensor.ID, actualSensor.ID)
		assert.Equal(t, sensor.SerialNumber, actualSensor.SerialNumber)
		assert.Equal(t, sensor.Description, actualSensor.Description)
		assert.Equal(t, sensor.IsActive, actualSensor.IsActive)
		assert.NotEmpty(t, actualSensor.RegisteredAt)
		assert.Empty(t, actualSensor.LastActivity)
	})

	t.Run("ok, collision test", func(t *testing.T) {
		sr := NewSensorRepository()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg := sync.WaitGroup{}
		r := rand.New(rand.NewPCG(42, 1024))
		for i := 0; i < 1000; i++ {
			sensor := &domain.Sensor{
				SerialNumber: strconv.Itoa(r.IntN(1000000000)),
				Type:         domain.SensorTypeADC,
				CurrentState: 0,
				Description:  fmt.Sprintf("some description %d", i),
				IsActive:     false,
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				assert.NoError(t, sr.SaveSensor(ctx, sensor))
			}()
		}

		wg.Wait()

		sensors, err := sr.GetSensors(ctx)
		assert.NoError(t, err)
		assert.Len(t, sensors, 1000)
	})
}

func TestSensorRepository_GetSensors(t *testing.T) {
	t.Run("fail, ctx cancelled", func(t *testing.T) {
		sr := NewSensorRepository()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := sr.GetSensors(ctx)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("fail, ctx deadline exceeded", func(t *testing.T) {
		sr := NewSensorRepository()
		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()

		_, err := sr.GetSensors(ctx)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("ok, get empty list", func(t *testing.T) {
		sr := NewSensorRepository()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sensors, err := sr.GetSensors(ctx)
		assert.NoError(t, err)
		assert.Len(t, sensors, 0)
	})

	t.Run("ok, get list", func(t *testing.T) {
		sr := NewSensorRepository()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		r := rand.New(rand.NewPCG(42, 1024))
		for i := 0; i < 10; i++ {
			sensor := &domain.Sensor{
				SerialNumber: strconv.Itoa(r.IntN(1000000000)),
				Type:         domain.SensorTypeADC,
				CurrentState: 0,
				Description:  fmt.Sprintf("some description %d", i),
				IsActive:     false,
			}
			assert.NoError(t, sr.SaveSensor(ctx, sensor))
		}

		sensors, err := sr.GetSensors(ctx)
		assert.NoError(t, err)
		assert.Len(t, sensors, 10)
	})
}

func TestSensorRepository_GetSensorByID(t *testing.T) {
	t.Run("fail, ctx cancelled", func(t *testing.T) {
		sr := NewSensorRepository()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := sr.GetSensorByID(ctx, 0)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("fail, ctx deadline exceeded", func(t *testing.T) {
		sr := NewSensorRepository()
		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()

		_, err := sr.GetSensorByID(ctx, 0)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("fail, not found", func(t *testing.T) {
		sr := NewSensorRepository()

		_, err := sr.GetSensorByID(context.Background(), 123)
		assert.ErrorIs(t, err, usecase.ErrSensorNotFound)
	})
}

func TestSensorRepository_GetSensorBySerialNumber(t *testing.T) {
	t.Run("fail, ctx cancelled", func(t *testing.T) {
		sr := NewSensorRepository()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := sr.GetSensorBySerialNumber(ctx, "1234")
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("fail, ctx deadline exceeded", func(t *testing.T) {
		sr := NewSensorRepository()
		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()

		_, err := sr.GetSensorBySerialNumber(ctx, "12345")
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("fail, not found", func(t *testing.T) {
		sr := NewSensorRepository()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		_, err := sr.GetSensorBySerialNumber(ctx, "12345")
		assert.ErrorIs(t, err, usecase.ErrSensorNotFound)
	})

	t.Run("ok, save and get one", func(t *testing.T) {
		sr := NewSensorRepository()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sensor := &domain.Sensor{
			SerialNumber: "12345678",
			Type:         domain.SensorTypeContactClosure,
			CurrentState: 0,
			Description:  "sensor description",
			IsActive:     true,
		}

		err := sr.SaveSensor(ctx, sensor)
		assert.NoError(t, err)

		actualSensor, err := sr.GetSensorBySerialNumber(ctx, sensor.SerialNumber)
		assert.NoError(t, err)
		assert.NotNil(t, actualSensor)
		assert.Equal(t, sensor.ID, actualSensor.ID)
		assert.Equal(t, sensor.SerialNumber, actualSensor.SerialNumber)
		assert.Equal(t, sensor.Description, actualSensor.Description)
		assert.Equal(t, sensor.IsActive, actualSensor.IsActive)
		assert.NotEmpty(t, actualSensor.RegisteredAt)
		assert.Empty(t, actualSensor.LastActivity)
	})
}
