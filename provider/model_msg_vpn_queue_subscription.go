package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// Terraform schema for MsgVpnQueueSubscription
func MsgVpnQueueSubscriptionSchema(requiredAttributes ...string) tfsdk.Schema {
	schema := tfsdk.Schema{
		Description: "MsgVpnQueueSubscription",
		Attributes: map[string]tfsdk.Attribute{
			"msg_vpn_name": {
				Type:        types.StringType,
				Description: "The name of the Message VPN.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"queue_name": {
				Type:        types.StringType,
				Description: "The name of the Queue.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"subscription_topic": {
				Type:        types.StringType,
				Description: "The topic of the Subscription.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
		},
	}

	return WithRequiredAttributes(schema, requiredAttributes)
}
