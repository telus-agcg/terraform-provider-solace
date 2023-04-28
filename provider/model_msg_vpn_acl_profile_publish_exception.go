package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"telusag/terraform-provider-solace/sempv2"
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

// Terraform DataSource schema for MsgVpnAclProfilePublishException
func MsgVpnAclProfilePublishExceptionDataSourceSchema(requiredAttributes ...string) dschema.Schema {
	schema := dschema.Schema{
		Description: "MsgVpnAclProfilePublishException",
		Attributes: map[string]dschema.Attribute{
			"acl_profile_name": dschema.StringAttribute{
				Description: "The name of the ACL Profile. Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
				},
			},
			"msg_vpn_name": dschema.StringAttribute{
				Description: "The name of the Message VPN. Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Optional:    true,
			},
			"publish_exception_topic": dschema.StringAttribute{
				Description: "The topic for the exception to the default action taken. May include wildcard characters. Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Optional:    true,
			},
			"topic_syntax": dschema.StringAttribute{
				Description: "The syntax of the topic for the exception to the default action taken. The allowed values and their meaning are:  <pre> \"smf\" - Topic uses SMF syntax. \"mqtt\" - Topic uses MQTT syntax. </pre>  Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("smf", "mqtt"),
				},
			},
		},
	}

	return schema
}

// Terraform Resource schema for MsgVpnAclProfilePublishException
func MsgVpnAclProfilePublishExceptionResourceSchema(requiredAttributes ...string) rschema.Schema {
	schema := rschema.Schema{
		Description: "MsgVpnAclProfilePublishException",
		Attributes: map[string]rschema.Attribute{
			"acl_profile_name": rschema.StringAttribute{
				Description: "The name of the ACL Profile. Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
				},
			},
			"msg_vpn_name": rschema.StringAttribute{
				Description: "The name of the Message VPN. Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Optional:    true,
			},
			"publish_exception_topic": rschema.StringAttribute{
				Description: "The topic for the exception to the default action taken. May include wildcard characters. Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Optional:    true,
			},
			"topic_syntax": rschema.StringAttribute{
				Description: "The syntax of the topic for the exception to the default action taken. The allowed values and their meaning are:  <pre> \"smf\" - Topic uses SMF syntax. \"mqtt\" - Topic uses MQTT syntax. </pre>  Deprecated since 2.14. Replaced by publishTopicExceptions.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("smf", "mqtt"),
				},
			},
		},
	}

	return schema
}
