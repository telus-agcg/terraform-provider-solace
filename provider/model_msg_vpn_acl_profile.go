package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"telusag/terraform-provider-solace/sempv2"
)

// MsgVpnAclProfile struct for MsgVpnAclProfile
type MsgVpnAclProfile struct {
	AclProfileName                  *string `tfsdk:"acl_profile_name"`
	ClientConnectDefaultAction      *string `tfsdk:"client_connect_default_action"`
	MsgVpnName                      *string `tfsdk:"msg_vpn_name"`
	PublishTopicDefaultAction       *string `tfsdk:"publish_topic_default_action"`
	SubscribeShareNameDefaultAction *string `tfsdk:"subscribe_share_name_default_action"`
	SubscribeTopicDefaultAction     *string `tfsdk:"subscribe_topic_default_action"`
}

func (tfData *MsgVpnAclProfile) ToTF(apiData *sempv2.MsgVpnAclProfile) {
	AssignIfDstNotNil(&tfData.AclProfileName, apiData.AclProfileName)
	AssignIfDstNotNil(&tfData.ClientConnectDefaultAction, apiData.ClientConnectDefaultAction)
	AssignIfDstNotNil(&tfData.MsgVpnName, apiData.MsgVpnName)
	AssignIfDstNotNil(&tfData.PublishTopicDefaultAction, apiData.PublishTopicDefaultAction)
	AssignIfDstNotNil(&tfData.SubscribeShareNameDefaultAction, apiData.SubscribeShareNameDefaultAction)
	AssignIfDstNotNil(&tfData.SubscribeTopicDefaultAction, apiData.SubscribeTopicDefaultAction)
}

func (tfData *MsgVpnAclProfile) ToApi() *sempv2.MsgVpnAclProfile {
	return &sempv2.MsgVpnAclProfile{
		AclProfileName:                  tfData.AclProfileName,
		ClientConnectDefaultAction:      tfData.ClientConnectDefaultAction,
		MsgVpnName:                      tfData.MsgVpnName,
		PublishTopicDefaultAction:       tfData.PublishTopicDefaultAction,
		SubscribeShareNameDefaultAction: tfData.SubscribeShareNameDefaultAction,
		SubscribeTopicDefaultAction:     tfData.SubscribeTopicDefaultAction,
	}
}

// Terraform DataSource schema for MsgVpnAclProfile
func MsgVpnAclProfileDataSourceSchema(requiredAttributes ...string) dschema.Schema {
	schema := dschema.Schema{
		Description: "MsgVpnAclProfile",
		Attributes: map[string]dschema.Attribute{
			"acl_profile_name": dschema.StringAttribute{
				Description: "The name of the ACL Profile.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
				},
			},
			"client_connect_default_action": dschema.StringAttribute{
				Description: "The default action to take when a client using the ACL Profile connects to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"disallow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow client connection unless an exception is found for it. \"disallow\" - Disallow client connection unless an exception is found for it. </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("allow", "disallow"),
				},
			},
			"msg_vpn_name": dschema.StringAttribute{
				Description: "The name of the Message VPN.",
				Optional:    true,
			},
			"publish_topic_default_action": dschema.StringAttribute{
				Description: "The default action to take when a client using the ACL Profile publishes to a topic in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"disallow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow topic unless an exception is found for it. \"disallow\" - Disallow topic unless an exception is found for it. </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("allow", "disallow"),
				},
			},
			"subscribe_share_name_default_action": dschema.StringAttribute{
				Description: "The default action to take when a client using the ACL Profile subscribes to a share-name subscription in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"allow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow topic unless an exception is found for it. \"disallow\" - Disallow topic unless an exception is found for it. </pre>  Available since 2.14.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("allow", "disallow"),
				},
			},
			"subscribe_topic_default_action": dschema.StringAttribute{
				Description: "The default action to take when a client using the ACL Profile subscribes to a topic in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"disallow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow topic unless an exception is found for it. \"disallow\" - Disallow topic unless an exception is found for it. </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("allow", "disallow"),
				},
			},
		},
	}

	return schema
}

// Terraform Resource schema for MsgVpnAclProfile
func MsgVpnAclProfileResourceSchema(requiredAttributes ...string) rschema.Schema {
	schema := rschema.Schema{
		Description: "MsgVpnAclProfile",
		Attributes: map[string]rschema.Attribute{
			"acl_profile_name": rschema.StringAttribute{
				Description: "The name of the ACL Profile.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
				},
			},
			"client_connect_default_action": rschema.StringAttribute{
				Description: "The default action to take when a client using the ACL Profile connects to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"disallow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow client connection unless an exception is found for it. \"disallow\" - Disallow client connection unless an exception is found for it. </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("allow", "disallow"),
				},
			},
			"msg_vpn_name": rschema.StringAttribute{
				Description: "The name of the Message VPN.",
				Optional:    true,
			},
			"publish_topic_default_action": rschema.StringAttribute{
				Description: "The default action to take when a client using the ACL Profile publishes to a topic in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"disallow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow topic unless an exception is found for it. \"disallow\" - Disallow topic unless an exception is found for it. </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("allow", "disallow"),
				},
			},
			"subscribe_share_name_default_action": rschema.StringAttribute{
				Description: "The default action to take when a client using the ACL Profile subscribes to a share-name subscription in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"allow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow topic unless an exception is found for it. \"disallow\" - Disallow topic unless an exception is found for it. </pre>  Available since 2.14.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("allow", "disallow"),
				},
			},
			"subscribe_topic_default_action": rschema.StringAttribute{
				Description: "The default action to take when a client using the ACL Profile subscribes to a topic in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"disallow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow topic unless an exception is found for it. \"disallow\" - Disallow topic unless an exception is found for it. </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("allow", "disallow"),
				},
			},
		},
	}

	return schema
}
