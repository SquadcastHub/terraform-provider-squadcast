package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/squadcast/terraform-provider-squadcast/internal/tf"
)

type ServiceMaintenanceWindow struct {
	From              string `json:"maintenance_from" tf:"from"`
	Till              string `json:"maintenance_till" tf:"till"`
	RepeatTill        string `json:"repeat_till" tf:"repeat_till"`
	RepeatFrequency   string `json:"-" tf:"repeat_frequency"`
	RepeatDaily       bool   `json:"repetition_daily" tf:"-"`
	RepeatWeekly      bool   `json:"repetition_weekly" tf:"-"`
	RepeatTwoWeekly   bool   `json:"repetition_two_weekly" tf:"-"`
	RepeatThreeWeekly bool   `json:"repetition_three_weekly" tf:"-"`
	RepeatMonthly     bool   `json:"repetition_monthly" tf:"-"`
}

func (s *ServiceMaintenanceWindow) Encode() (tf.M, error) {
	if s.RepeatDaily {
		s.RepeatFrequency = "day"
	} else if s.RepeatWeekly {
		s.RepeatFrequency = "week"
	} else if s.RepeatTwoWeekly {
		s.RepeatFrequency = "2 weeks"
	} else if s.RepeatThreeWeekly {
		s.RepeatFrequency = "3 weeks"
	} else if s.RepeatMonthly {
		s.RepeatFrequency = "month"
	}

	if s.RepeatFrequency == "" {
		s.RepeatTill = ""
	}

	m, err := tf.Encode(s)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (client *Client) GetServiceMaintenanceWindows(ctx context.Context, serviceID string) ([]*ServiceMaintenanceWindow, error) {
	url := fmt.Sprintf("%s/services/%s/maintenance", client.BaseURLV3, serviceID)

	return RequestSlice[any, ServiceMaintenanceWindow](http.MethodGet, url, client, ctx, nil)
}

type UpdateServiceMaintenanceWindowsWindow struct {
	From        string `json:"maintenanceStartDate"`
	Till        string `json:"maintenanceEndDate"`
	RepeatTill  string `json:"repeatTill"`
	Daily       bool   `json:"daily"`
	Weekly      bool   `json:"weekly"`
	TwoWeekly   bool   `json:"twoWeekly"`
	ThreeWeekly bool   `json:"threeWeekly"`
	Monthly     bool   `json:"monthly"`
}

type UpdateServiceMaintenanceWindowsData struct {
	ServiceMaintenanceWindows []UpdateServiceMaintenanceWindowsWindow `json:"serviceMaintenance"`
}

type UpdateServiceMaintenanceWindows struct {
	Data           UpdateServiceMaintenanceWindowsData `json:"data"`
	OrganizationID string                              `json:"organizationId"`
	ServiceID      string                              `json:"serviceId"`
}

func (client *Client) UpdateServiceMaintenance(ctx context.Context, serviceID string, req *UpdateServiceMaintenanceWindows) (*any, error) {
	url := fmt.Sprintf("%s/services/%s/maintenance", client.BaseURLV3, serviceID)
	return Request[UpdateServiceMaintenanceWindows, any](http.MethodPost, url, client, ctx, req)
}
