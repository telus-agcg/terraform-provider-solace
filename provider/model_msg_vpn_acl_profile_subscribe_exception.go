package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"telusag/terraform-provider-solace/sempv2"
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

// Terraform DataSource schema for MsgVpnAclProfileSubscribeException
func MsgVpnAclProfileSubscribeExceptionDataSourceSchema(requiredAttributes ...string) dschema.Schema {
	schema := dschema.Schema{
		Description: "MsgVpnAclProfileSubscribeException",
		Attributes: map[string]dschema.Attribute{
			"acl_profile_name": dschema.StringAttribute{
				Description: "The name of the ACL Profile. Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
				},
			},
			"msg_vpn_name": dschema.StringAttribute{
				Description: "The name of the Message VPN. Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Optional:    true,
			},
			"subscribe_exception_topic": dschema.StringAttribute{
				Description: "The topic for the exception to the default action taken. May include wildcard characters. Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Optional:    true,
			},
			"topic_syntax": dschema.StringAttribute{
				Description: "The syntax of the topic for the exception to the default action taken. The allowed values and their meaning are:  <pre> \"smf\" - Topic uses SMF syntax. \"mqtt\" - Topic uses MQTT syntax. </pre>  Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("smf", "mqtt"),
				},
			},
		},
	}

	return schema
}

// Terraform Resource schema for MsgVpnAclProfileSubscribeException
func MsgVpnAclProfileSubscribeExceptionResourceSchema(requiredAttributes ...string) rschema.Schema {
	schema := rschema.Schema{
		Description: "MsgVpnAclProfileSubscribeException",
		Attributes: map[string]rschema.Attribute{
			"acl_profile_name": rschema.StringAttribute{
				Description: "The name of the ACL Profile. Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
				},
			},
			"msg_vpn_name": rschema.StringAttribute{
				Description: "The name of the Message VPN. Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Optional:    true,
			},
			"subscribe_exception_topic": rschema.StringAttribute{
				Description: "The topic for the exception to the default action taken. May include wildcard characters. Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Optional:    true,
			},
			"topic_syntax": rschema.StringAttribute{
				Description: "The syntax of the topic for the exception to the default action taken. The allowed values and their meaning are:  <pre> \"smf\" - Topic uses SMF syntax. \"mqtt\" - Topic uses MQTT syntax. </pre>  Deprecated since 2.14. Replaced by subscribeTopicExceptions.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("smf", "mqtt"),
				},
			},
		},
	}

	return schema
}
