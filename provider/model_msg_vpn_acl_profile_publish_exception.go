package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// MsgVpnAclProfilePublishException struct for MsgVpnAclProfilePublishException
type MsgVpnAclProfilePublishException struct {
	AclProfileName        *string `tfsdk:"acl_profile_name"`
	MsgVpnName            *string `tfsdk:"msg_vpn_name"`
	PublishExceptionTopic *string `tfsdk:"publish_exception_topic"`
	TopicSyntax           *string `tfsdk:"topic_syntax"`
}

func (tfData *MsgVpnAclProfilePublishException) ToTF(apiData *sempv2.MsgVpnAclProfilePublishException) {
	AssignIfDstNotNil(&tfData.AclProfileName, apiData.AclProfileName)
	AssignIfDstNotNil(&tfData.MsgVpnName, apiData.MsgVpnName)
	AssignIfDstNotNil(&tfData.PublishExceptionTopic, apiData.PublishExceptionTopic)
	AssignIfDstNotNil(&tfData.TopicSyntax, apiData.TopicSyntax)
}

func (tfData *MsgVpnAclProfilePublishException) ToApi() *sempv2.MsgVpnAclProfilePublishException {
	return &sempv2.MsgVpnAclProfilePublishException{
		AclProfileName:        tfData.AclProfileName,
		MsgVpnName:            tfData.MsgVpnName,
		PublishExceptionTopic: tfData.PublishExceptionTopic,
		TopicSyntax:           tfData.TopicSyntax,
	}
}

// Terraform Resource schema for MsgVpnAclProfilePublishException
func MsgVpnAclProfilePublishExceptionResourceSchema(requiredAttributes ...string) schema.Schema {
	schema := schema.Schema{
		Description: "MsgVpnAclProfilePublishException",
		Attributes: map[string]schema.Attribute{
			"acl_profile_name": schema.StringAttribute{
				Description: "The name of the ACL Profile. Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Required:    contains(requiredAttributes, "acl_profile_name"),
				Optional:    !contains(requiredAttributes, "acl_profile_name"),

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
				},
				PlanModifiers: StringPlanModifiersFor("acl_profile_name", requiredAttributes),
			},
			"msg_vpn_name": schema.StringAttribute{
				Description: "The name of the Message VPN. Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Required:    contains(requiredAttributes, "msg_vpn_name"),
				Optional:    !contains(requiredAttributes, "msg_vpn_name"),

				PlanModifiers: StringPlanModifiersFor("msg_vpn_name", requiredAttributes),
			},
			"publish_exception_topic": schema.StringAttribute{
				Description: "The topic for the exception to the default action taken. May include wildcard characters. Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Required:    contains(requiredAttributes, "publish_exception_topic"),
				Optional:    !contains(requiredAttributes, "publish_exception_topic"),

				PlanModifiers: StringPlanModifiersFor("publish_exception_topic", requiredAttributes),
			},
			"topic_syntax": schema.StringAttribute{
				Description: "The syntax of the topic for the exception to the default action taken. The allowed values and their meaning are:  <pre> \"smf\" - Topic uses SMF syntax. \"mqtt\" - Topic uses MQTT syntax. </pre>  Deprecated since 2.14. Replaced by publishTopicExceptions.",
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
