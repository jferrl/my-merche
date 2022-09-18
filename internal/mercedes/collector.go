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

type bundle struct {
	LockStatus          []*merche.VehicleLockStatus
	Status              []*merche.VehicleStatus
	PayAsYouDriveStatus []*merche.PayAsYouDriveStatus
	FuelStatus          []*merche.FuelStatus
}

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

	var b bundle

	opts := &merche.Options{
		VehicleID: string(c.vehicle),
	}

	vls, _, err := c.m.VehicleLockStatus.GetVehicleLockStatus(ctx, opts)
	if err != nil {
		return "", err
	}
	b.LockStatus = vls

	vs, _, err := c.m.VehicleStatus.GetVehicleStatus(ctx, opts)
	if err != nil {
		return "", err
	}
	b.Status = vs

	psd, _, err := c.m.PayAsYouDrive.GetPayAsYouDriveStatus(ctx, opts)
	if err != nil {
		return "", err
	}
	b.PayAsYouDriveStatus = psd

	fs, _, err := c.m.FuelStatus.GetFuelStatus(ctx, opts)
	if err != nil {
		return "", err
	}
	b.FuelStatus = fs

	e, err := json.MarshalIndent(&b, "", "  ")
	if err != nil {
		return "", err
	}

	return string(e), nil
}
