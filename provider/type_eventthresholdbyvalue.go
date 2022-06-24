package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var EventThresholdByValueType attr.Type = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"clear_value": types.Int64Type,
		"set_value":   types.Int64Type,
	},
}

type EventThresholdByValue struct {
	ClearValue *int64 `tfsdk:"clear_value"`
	SetValue   *int64 `tfsdk:"set_value"`
}

func (tfData *EventThresholdByValue) ToApi() *sempv2.EventThresholdByValue {
	if tfData == nil {
		return nil
	}

	return &sempv2.EventThresholdByValue{
		ClearValue: tfData.ClearValue,
		SetValue:   tfData.SetValue,
	}
}

func EventThresholdByValueToTF(api *sempv2.EventThresholdByValue) *EventThresholdByValue {
	if api == nil {
		return nil
	}

	return &EventThresholdByValue{
		ClearValue: api.ClearValue,
		SetValue:   api.SetValue,
	}
}
