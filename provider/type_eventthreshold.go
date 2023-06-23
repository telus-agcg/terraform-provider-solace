package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var EventThresholdResourceAttributes map[string]rschema.Attribute = map[string]rschema.Attribute{
	"clear_value": rschema.Int64Attribute{
		Required: false,
		Optional: true,
	},
	"clear_percent": rschema.Int64Attribute{
		Required: false,
		Optional: true,
	},
	"set_value": rschema.Int64Attribute{
		Required: false,
		Optional: true,
	},
	"set_percent": rschema.Int64Attribute{
		Required: false,
		Optional: true,
	},
}

var EventThresholdDatasourceAttributes map[string]dschema.Attribute = map[string]dschema.Attribute{
	"clear_value": dschema.Int64Attribute{
		Required: false,
		Optional: true,
	},
	"clear_percent": dschema.Int64Attribute{
		Required: false,
		Optional: true,
	},
	"set_value": dschema.Int64Attribute{
		Required: false,
		Optional: true,
	},
	"set_percent": dschema.Int64Attribute{
		Required: false,
		Optional: true,
	},
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
