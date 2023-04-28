package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

// Terraform Resource schema for MsgVpnQueueSubscription
func MsgVpnQueueSubscriptionResourceSchema(requiredAttributes ...string) schema.Schema {
	schema := schema.Schema{
		Description: "MsgVpnQueueSubscription",
		Attributes: map[string]schema.Attribute{
			"msg_vpn_name": schema.StringAttribute{
				Description: "The name of the Message VPN.",
				Required:    contains(requiredAttributes, "msg_vpn_name"),
				Optional:    !contains(requiredAttributes, "msg_vpn_name"),

				PlanModifiers: StringPlanModifiersFor("msg_vpn_name", requiredAttributes),
			},
			"queue_name": schema.StringAttribute{
				Description: "The name of the Queue.",
				Required:    contains(requiredAttributes, "queue_name"),
				Optional:    !contains(requiredAttributes, "queue_name"),

				PlanModifiers: StringPlanModifiersFor("queue_name", requiredAttributes),
			},
			"subscription_topic": schema.StringAttribute{
				Description: "The topic of the Subscription.",
				Required:    contains(requiredAttributes, "subscription_topic"),
				Optional:    !contains(requiredAttributes, "subscription_topic"),

				PlanModifiers: StringPlanModifiersFor("subscription_topic", requiredAttributes),
			},
		},
	}

	return schema
}
