package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var EventThresholdAttributeTypes map[string]attr.Type = map[string]attr.Type{
	"clear_value":   types.Int64Type,
	"clear_percent": types.Int64Type,
	"set_value":     types.Int64Type,
	"set_percent":   types.Int64Type,
}

type EventThreshold struct {
	ClearPercent *int64 `tfsdk:"clear_percent"`
	ClearValue   *int64 `tfsdk:"clear_value"`
	SetPercent   *int64 `tfsdk:"set_percent"`
	SetValue     *int64 `tfsdk:"set_value"`
}

func (tfData *EventThreshold) ToApi() *sempv2.EventThreshold {
	if tfData == nil {
		return nil
	}

	return &sempv2.EventThreshold{
		ClearPercent: tfData.ClearPercent,
		ClearValue:   tfData.ClearValue,
		SetPercent:   tfData.SetPercent,
		SetValue:     tfData.SetValue,
	}
}

func EventThresholdToTF(api *sempv2.EventThreshold) *EventThreshold {
	if api == nil {
		return nil
	}

	return &EventThreshold{
		ClearPercent: api.ClearPercent,
		ClearValue:   api.ClearValue,
		SetPercent:   api.SetPercent,
		SetValue:     api.SetValue,
	}
}
