package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var EventThresholdByPercentResourceAttributes map[string]schema.Attribute = map[string]schema.Attribute{
	"clear_percent": schema.Int64Attribute{
		Required: true,
		Optional: false,
	},
	"set_percent": schema.Int64Attribute{
		Required: true,
		Optional: false,
	},
}

type EventThresholdByPercent struct {
	ClearPercent *int64 `tfsdk:"clear_percent"`
	SetPercent   *int64 `tfsdk:"set_percent"`
}

func (tfData *EventThresholdByPercent) ToApi() *sempv2.EventThresholdByPercent {
	if tfData == nil {
		return nil
	}

	return &sempv2.EventThresholdByPercent{
		ClearPercent: tfData.ClearPercent,
		SetPercent:   tfData.SetPercent,
	}
}

func EventThresholdByPercentToTF(api *sempv2.EventThresholdByPercent) *EventThresholdByPercent {
	if api == nil {
		return nil
	}

	return &EventThresholdByPercent{
		ClearPercent: api.ClearPercent,
		SetPercent:   api.SetPercent,
	}
}
