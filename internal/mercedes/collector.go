package mercedes

import (
	"context"
	"errors"
	"net/http"
	"reflect"

	"github.com/jferrl/go-merche"
)

type Resource struct {
	Timestamp int64
	Value     string
}

type ResouceID string

type Resouces map[ResouceID]any

type VehicleID string

type Collector struct {
	m       *merche.Client
	vehicle VehicleID
}

func NewCollector(vID VehicleID) *Collector {
	return &Collector{
		vehicle: vID,
	}
}

func (c *Collector) Bootstrap(httpClient *http.Client) {
	c.m = merche.NewClient(httpClient)
}

func (c *Collector) Collect(ctx context.Context) (Resouces, error) {
	if c.m == nil {
		return nil, errors.New("mercedes api client must be defined")
	}

	vls, _, err := c.m.VehicleLockStatus.GetVehicleLockStatus(ctx, &merche.Options{VehicleID: string(c.vehicle)})
	if err != nil {
		return nil, err
	}

	var resources Resouces
	for _, ls := range vls {
		v := reflect.Indirect(reflect.ValueOf(ls))
		for i := 0; i < v.NumField(); i++ {
			resources[ResouceID(v.Field(i).Type().Name())] = v.Field(i).Interface()
		}
	}

	return resources, nil
}
