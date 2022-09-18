package mercedes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jferrl/go-merche"
)

type Resource struct {
	Timestamp int64
	Value     string
}

type ResouceID string

type Resouces map[ResouceID]Resource

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

func (c *Collector) Collect(ctx context.Context) (string, error) {
	if c.m == nil {
		return "", errors.New("mercedes api client must be defined")
	}

	vls, _, err := c.m.VehicleLockStatus.GetVehicleLockStatus(ctx, &merche.Options{VehicleID: string(c.vehicle)})
	if err != nil {
		return "", err
	}

	e, err := json.Marshal(&vls)
	if err != nil {
		return "", err
	}

	return string(e), nil
}
