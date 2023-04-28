package provider

import (
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"telusag/terraform-provider-solace/sempv2"
)

// MsgVpnQueueSubscription struct for MsgVpnQueueSubscription
type MsgVpnQueueSubscription struct {
	MsgVpnName        *string `tfsdk:"msg_vpn_name"`
	QueueName         *string `tfsdk:"queue_name"`
	SubscriptionTopic *string `tfsdk:"subscription_topic"`
}

func (tfData *MsgVpnQueueSubscription) ToTF(apiData *sempv2.MsgVpnQueueSubscription) {
	AssignIfDstNotNil(&tfData.MsgVpnName, apiData.MsgVpnName)
	AssignIfDstNotNil(&tfData.QueueName, apiData.QueueName)
	AssignIfDstNotNil(&tfData.SubscriptionTopic, apiData.SubscriptionTopic)
}

func (tfData *MsgVpnQueueSubscription) ToApi() *sempv2.MsgVpnQueueSubscription {
	return &sempv2.MsgVpnQueueSubscription{
		MsgVpnName:        tfData.MsgVpnName,
		QueueName:         tfData.QueueName,
		SubscriptionTopic: tfData.SubscriptionTopic,
	}
}

// Terraform DataSource schema for MsgVpnQueueSubscription
func MsgVpnQueueSubscriptionDataSourceSchema(requiredAttributes ...string) dschema.Schema {
	schema := dschema.Schema{
		Description: "MsgVpnQueueSubscription",
		Attributes: map[string]dschema.Attribute{
			"msg_vpn_name": dschema.StringAttribute{
				Description: "The name of the Message VPN.",
				Optional:    true,
			},
			"queue_name": dschema.StringAttribute{
				Description: "The name of the Queue.",
				Optional:    true,
			},
			"subscription_topic": dschema.StringAttribute{
				Description: "The topic of the Subscription.",
				Optional:    true,
			},
		},
	}

	return schema
}

// Terraform Resource schema for MsgVpnQueueSubscription
func MsgVpnQueueSubscriptionResourceSchema(requiredAttributes ...string) rschema.Schema {
	schema := rschema.Schema{
		Description: "MsgVpnQueueSubscription",
		Attributes: map[string]rschema.Attribute{
			"msg_vpn_name": rschema.StringAttribute{
				Description: "The name of the Message VPN.",
				Optional:    true,
			},
			"queue_name": rschema.StringAttribute{
				Description: "The name of the Queue.",
				Optional:    true,
			},
			"subscription_topic": rschema.StringAttribute{
				Description: "The topic of the Subscription.",
				Optional:    true,
			},
		},
	}

	return schema
}
