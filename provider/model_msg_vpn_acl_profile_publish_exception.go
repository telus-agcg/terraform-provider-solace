package provider

import (
	"telusag/terraform-provider-solace/sempv2"
	"telusag/terraform-provider-solace/util"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// Terraform schema for MsgVpnAclProfilePublishException
func MsgVpnAclProfilePublishExceptionSchema(requiredAttributes ...string) tfsdk.Schema {
	schema := tfsdk.Schema{
		Description: "MsgVpnAclProfilePublishException",
		Attributes: map[string]tfsdk.Attribute{
			"acl_profile_name": {
				Type:        types.StringType,
				Description: "The name of the ACL Profile. Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Optional:    true,
			},
			"msg_vpn_name": {
				Type:        types.StringType,
				Description: "The name of the Message VPN. Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Optional:    true,
			},
			"publish_exception_topic": {
				Type:        types.StringType,
				Description: "The topic for the exception to the default action taken. May include wildcard characters. Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Optional:    true,
			},
			"topic_syntax": {
				Type:        types.StringType,
				Description: "The syntax of the topic for the exception to the default action taken. The allowed values and their meaning are:  <pre> \"smf\" - Topic uses SMF syntax. \"mqtt\" - Topic uses MQTT syntax. </pre>  Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					util.StringOneOfValidator("smf", "mqtt"),
				},
			},
		},
	}

	return WithRequiredAttributes(schema, requiredAttributes)
}
