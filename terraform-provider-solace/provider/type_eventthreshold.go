package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var EventThresholdType attr.Type = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"clear_value":   types.Int64Type,
		"clear_percent": types.Int64Type,
		"set_value":     types.Int64Type,
		"set_percent":   types.Int64Type,
	},
}

type EventThreshold struct {
	ClearPercent *int64 `tfsdk:"clear_percent"`
	ClearValue   *int64 `tfsdk:"clear_value"`
	SetPercent   *int64 `tfsdk:"set_percent"`
	SetValue     *int64 `tfsdk:"set_value"`
}
