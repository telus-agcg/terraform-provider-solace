package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

// Terraform Resource schema for MsgVpnAclProfile
func MsgVpnAclProfileResourceSchema(requiredAttributes ...string) schema.Schema {
	schema := schema.Schema{
		Description: "MsgVpnAclProfile",
		Attributes: map[string]schema.Attribute{
			"acl_profile_name": schema.StringAttribute{
				Description: "The name of the ACL Profile.",
				Required:    contains(requiredAttributes, "acl_profile_name"),
				Optional:    !contains(requiredAttributes, "acl_profile_name"),

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
				},
				PlanModifiers: StringPlanModifiersFor("acl_profile_name", requiredAttributes),
			},
			"client_connect_default_action": schema.StringAttribute{
				Description: "The default action to take when a client using the ACL Profile connects to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"disallow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow client connection unless an exception is found for it. \"disallow\" - Disallow client connection unless an exception is found for it. </pre> ",
				Required:    contains(requiredAttributes, "client_connect_default_action"),
				Optional:    !contains(requiredAttributes, "client_connect_default_action"),

				Validators: []validator.String{
					stringvalidator.OneOf("allow", "disallow"),
				},
				PlanModifiers: StringPlanModifiersFor("client_connect_default_action", requiredAttributes),
			},
			"msg_vpn_name": schema.StringAttribute{
				Description: "The name of the Message VPN.",
				Required:    contains(requiredAttributes, "msg_vpn_name"),
				Optional:    !contains(requiredAttributes, "msg_vpn_name"),

				PlanModifiers: StringPlanModifiersFor("msg_vpn_name", requiredAttributes),
			},
			"publish_topic_default_action": schema.StringAttribute{
				Description: "The default action to take when a client using the ACL Profile publishes to a topic in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"disallow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow topic unless an exception is found for it. \"disallow\" - Disallow topic unless an exception is found for it. </pre> ",
				Required:    contains(requiredAttributes, "publish_topic_default_action"),
				Optional:    !contains(requiredAttributes, "publish_topic_default_action"),

				Validators: []validator.String{
					stringvalidator.OneOf("allow", "disallow"),
				},
				PlanModifiers: StringPlanModifiersFor("publish_topic_default_action", requiredAttributes),
			},
			"subscribe_share_name_default_action": schema.StringAttribute{
				Description: "The default action to take when a client using the ACL Profile subscribes to a share-name subscription in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"allow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow topic unless an exception is found for it. \"disallow\" - Disallow topic unless an exception is found for it. </pre>  Available since 2.14.",
				Required:    contains(requiredAttributes, "subscribe_share_name_default_action"),
				Optional:    !contains(requiredAttributes, "subscribe_share_name_default_action"),

				Validators: []validator.String{
					stringvalidator.OneOf("allow", "disallow"),
				},
				PlanModifiers: StringPlanModifiersFor("subscribe_share_name_default_action", requiredAttributes),
			},
			"subscribe_topic_default_action": schema.StringAttribute{
				Description: "The default action to take when a client using the ACL Profile subscribes to a topic in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"disallow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow topic unless an exception is found for it. \"disallow\" - Disallow topic unless an exception is found for it. </pre> ",
				Required:    contains(requiredAttributes, "subscribe_topic_default_action"),
				Optional:    !contains(requiredAttributes, "subscribe_topic_default_action"),

				Validators: []validator.String{
					stringvalidator.OneOf("allow", "disallow"),
				},
				PlanModifiers: StringPlanModifiersFor("subscribe_topic_default_action", requiredAttributes),
			},
		},
	}

	return schema
}
