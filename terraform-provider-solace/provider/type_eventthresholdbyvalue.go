package provider

import (
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
