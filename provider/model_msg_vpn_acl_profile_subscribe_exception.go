package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// MsgVpnAclProfileSubscribeException struct for MsgVpnAclProfileSubscribeException
type MsgVpnAclProfileSubscribeException struct {
	AclProfileName          *string `tfsdk:"acl_profile_name"`
	MsgVpnName              *string `tfsdk:"msg_vpn_name"`
	SubscribeExceptionTopic *string `tfsdk:"subscribe_exception_topic"`
	TopicSyntax             *string `tfsdk:"topic_syntax"`
}

func (tfData *MsgVpnAclProfileSubscribeException) ToTF(apiData *sempv2.MsgVpnAclProfileSubscribeException) {
	AssignIfDstNotNil(&tfData.AclProfileName, apiData.AclProfileName)
	AssignIfDstNotNil(&tfData.MsgVpnName, apiData.MsgVpnName)
	AssignIfDstNotNil(&tfData.SubscribeExceptionTopic, apiData.SubscribeExceptionTopic)
	AssignIfDstNotNil(&tfData.TopicSyntax, apiData.TopicSyntax)
}

func (tfData *MsgVpnAclProfileSubscribeException) ToApi() *sempv2.MsgVpnAclProfileSubscribeException {
	return &sempv2.MsgVpnAclProfileSubscribeException{
		AclProfileName:          tfData.AclProfileName,
		MsgVpnName:              tfData.MsgVpnName,
		SubscribeExceptionTopic: tfData.SubscribeExceptionTopic,
		TopicSyntax:             tfData.TopicSyntax,
	}
}

// Terraform Resource schema for MsgVpnAclProfileSubscribeException
func MsgVpnAclProfileSubscribeExceptionResourceSchema(requiredAttributes ...string) schema.Schema {
	schema := schema.Schema{
		Description: "MsgVpnAclProfileSubscribeException",
		Attributes: map[string]schema.Attribute{
			"acl_profile_name": schema.StringAttribute{
				Description: "The name of the ACL Profile. Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Required:    contains(requiredAttributes, "acl_profile_name"),
				Optional:    !contains(requiredAttributes, "acl_profile_name"),

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
				},
				PlanModifiers: StringPlanModifiersFor("acl_profile_name", requiredAttributes),
			},
			"msg_vpn_name": schema.StringAttribute{
				Description: "The name of the Message VPN. Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Required:    contains(requiredAttributes, "msg_vpn_name"),
				Optional:    !contains(requiredAttributes, "msg_vpn_name"),

				PlanModifiers: StringPlanModifiersFor("msg_vpn_name", requiredAttributes),
			},
			"subscribe_exception_topic": schema.StringAttribute{
				Description: "The topic for the exception to the default action taken. May include wildcard characters. Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Required:    contains(requiredAttributes, "subscribe_exception_topic"),
				Optional:    !contains(requiredAttributes, "subscribe_exception_topic"),

				PlanModifiers: StringPlanModifiersFor("subscribe_exception_topic", requiredAttributes),
			},
			"topic_syntax": schema.StringAttribute{
				Description: "The syntax of the topic for the exception to the default action taken. The allowed values and their meaning are:  <pre> \"smf\" - Topic uses SMF syntax. \"mqtt\" - Topic uses MQTT syntax. </pre>  Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Required:    contains(requiredAttributes, "topic_syntax"),
				Optional:    !contains(requiredAttributes, "topic_syntax"),

				Validators: []validator.String{
					stringvalidator.OneOf("smf", "mqtt"),
				},
				PlanModifiers: StringPlanModifiersFor("topic_syntax", requiredAttributes),
			},
		},
	}

	return schema
}
