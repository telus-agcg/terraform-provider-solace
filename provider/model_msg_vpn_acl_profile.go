package provider

import (
	"telusag/terraform-provider-solace/sempv2"
	"telusag/terraform-provider-solace/util"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func (tfData *MsgVpnAclProfile) ToApi() sempv2.MsgVpnAclProfile {
	return sempv2.MsgVpnAclProfile{
		AclProfileName:                  tfData.AclProfileName,
		ClientConnectDefaultAction:      tfData.ClientConnectDefaultAction,
		MsgVpnName:                      tfData.MsgVpnName,
		PublishTopicDefaultAction:       tfData.PublishTopicDefaultAction,
		SubscribeShareNameDefaultAction: tfData.SubscribeShareNameDefaultAction,
		SubscribeTopicDefaultAction:     tfData.SubscribeTopicDefaultAction,
	}
}

// Terraform schema for MsgVpnAclProfile
func MsgVpnAclProfileSchema(requiredAttributes ...string) tfsdk.Schema {
	schema := tfsdk.Schema{
		Description: "MsgVpnAclProfile",
		Attributes: map[string]tfsdk.Attribute{
			"acl_profile_name": {
				Type:        types.StringType,
				Description: "The name of the ACL Profile.",
				Optional:    true,
				// PlanModifiers: []tfsdk.AttributePlanModifier{
				// 	tfsdk.RequiresReplace(),
				// },
			},
			"client_connect_default_action": {
				Type:        types.StringType,
				Description: "The default action to take when a client using the ACL Profile connects to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"disallow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow client connection unless an exception is found for it. \"disallow\" - Disallow client connection unless an exception is found for it. </pre> ",
				Optional:    true,
				// PlanModifiers: []tfsdk.AttributePlanModifier{
				// 	tfsdk.RequiresReplace(),
				// },
				Validators: []tfsdk.AttributeValidator{
					util.StringOneOfValidator("allow", "disallow"),
				},
			},
			"msg_vpn_name": {
				Type:        types.StringType,
				Description: "The name of the Message VPN.",
				Optional:    true,
				// PlanModifiers: []tfsdk.AttributePlanModifier{
				// 	tfsdk.RequiresReplace(),
				// },
			},
			"publish_topic_default_action": {
				Type:        types.StringType,
				Description: "The default action to take when a client using the ACL Profile publishes to a topic in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"disallow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow topic unless an exception is found for it. \"disallow\" - Disallow topic unless an exception is found for it. </pre> ",
				Optional:    true,
				// PlanModifiers: []tfsdk.AttributePlanModifier{
				// 	tfsdk.RequiresReplace(),
				// },
				Validators: []tfsdk.AttributeValidator{
					util.StringOneOfValidator("allow", "disallow"),
				},
			},
			"subscribe_share_name_default_action": {
				Type:        types.StringType,
				Description: "The default action to take when a client using the ACL Profile subscribes to a share-name subscription in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"allow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow topic unless an exception is found for it. \"disallow\" - Disallow topic unless an exception is found for it. </pre>  Available since 2.14.",
				Optional:    true,
				// PlanModifiers: []tfsdk.AttributePlanModifier{
				// 	tfsdk.RequiresReplace(),
				// },
				Validators: []tfsdk.AttributeValidator{
					util.StringOneOfValidator("allow", "disallow"),
				},
			},
			"subscribe_topic_default_action": {
				Type:        types.StringType,
				Description: "The default action to take when a client using the ACL Profile subscribes to a topic in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"disallow\"`. The allowed values and their meaning are:  <pre> \"allow\" - Allow topic unless an exception is found for it. \"disallow\" - Disallow topic unless an exception is found for it. </pre> ",
				Optional:    true,
				// PlanModifiers: []tfsdk.AttributePlanModifier{
				// 	tfsdk.RequiresReplace(),
				// },
				Validators: []tfsdk.AttributeValidator{
					util.StringOneOfValidator("allow", "disallow"),
				},
			},
		},
	}

	return WithRequiredAttributes(schema, requiredAttributes)
}
