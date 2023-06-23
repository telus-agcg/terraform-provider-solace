package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var EventThresholdByValueResourceAttributes map[string]rschema.Attribute = map[string]rschema.Attribute{
	"clear_value": rschema.Int64Attribute{
		Required: true,
		Optional: false,
	},
	"set_value": rschema.Int64Attribute{
		Required: true,
		Optional: false,
	},
}

var EventThresholdByValueDatasourceAttributes map[string]dschema.Attribute = map[string]dschema.Attribute{
	"clear_value": dschema.Int64Attribute{
		Required: true,
		Optional: false,
	},
	"set_value": dschema.Int64Attribute{
		Required: true,
		Optional: false,
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
