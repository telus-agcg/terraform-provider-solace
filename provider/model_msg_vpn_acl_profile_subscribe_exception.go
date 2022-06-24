package provider

import (
	"telusag/terraform-provider-solace/sempv2"
	"telusag/terraform-provider-solace/util"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// Terraform schema for MsgVpnAclProfileSubscribeException
func MsgVpnAclProfileSubscribeExceptionSchema(requiredAttributes ...string) tfsdk.Schema {
	schema := tfsdk.Schema{
		Description: "MsgVpnAclProfileSubscribeException",
		Attributes: map[string]tfsdk.Attribute{
			"acl_profile_name": {
				Type:        types.StringType,
				Description: "The name of the ACL Profile. Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Optional:    true,
				// PlanModifiers: []tfsdk.AttributePlanModifier{
				// 	tfsdk.RequiresReplace(),
				// },
			},
			"msg_vpn_name": {
				Type:        types.StringType,
				Description: "The name of the Message VPN. Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Optional:    true,
				// PlanModifiers: []tfsdk.AttributePlanModifier{
				// 	tfsdk.RequiresReplace(),
				// },
			},
			"subscribe_exception_topic": {
				Type:        types.StringType,
				Description: "The topic for the exception to the default action taken. May include wildcard characters. Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Optional:    true,
				// PlanModifiers: []tfsdk.AttributePlanModifier{
				// 	tfsdk.RequiresReplace(),
				// },
			},
			"topic_syntax": {
				Type:        types.StringType,
				Description: "The syntax of the topic for the exception to the default action taken. The allowed values and their meaning are:  <pre> \"smf\" - Topic uses SMF syntax. \"mqtt\" - Topic uses MQTT syntax. </pre>  Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Optional:    true,
				// PlanModifiers: []tfsdk.AttributePlanModifier{
				// 	tfsdk.RequiresReplace(),
				// },
				Validators: []tfsdk.AttributeValidator{
					util.StringOneOfValidator("smf", "mqtt"),
				},
			},
		},
	}

	return WithRequiredAttributes(schema, requiredAttributes)
}
