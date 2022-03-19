package provider

import (
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
