package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var EventThresholdByPercentType attr.Type = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"clear_percent": types.Int64Type,
		"set_percent":   types.Int64Type,
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
