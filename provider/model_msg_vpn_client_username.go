package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MsgVpnClientUsername struct for MsgVpnClientUsername
type MsgVpnClientUsername struct {
	AclProfileName                              *string `tfsdk:"acl_profile_name"`
	ClientProfileName                           *string `tfsdk:"client_profile_name"`
	ClientUsername                              *string `tfsdk:"client_username"`
	Enabled                                     *bool   `tfsdk:"enabled"`
	GuaranteedEndpointPermissionOverrideEnabled *bool   `tfsdk:"guaranteed_endpoint_permission_override_enabled"`
	MsgVpnName                                  *string `tfsdk:"msg_vpn_name"`
	Password                                    *string `tfsdk:"password"`
	SubscriptionManagerEnabled                  *bool   `tfsdk:"subscription_manager_enabled"`
}

func (tfData *MsgVpnClientUsername) ToTF(apiData *sempv2.MsgVpnClientUsername) {
	AssignIfDstNotNil(&tfData.AclProfileName, apiData.AclProfileName)
	AssignIfDstNotNil(&tfData.ClientProfileName, apiData.ClientProfileName)
	AssignIfDstNotNil(&tfData.ClientUsername, apiData.ClientUsername)
	AssignIfDstNotNil(&tfData.Enabled, apiData.Enabled)
	AssignIfDstNotNil(&tfData.GuaranteedEndpointPermissionOverrideEnabled, apiData.GuaranteedEndpointPermissionOverrideEnabled)
	AssignIfDstNotNil(&tfData.MsgVpnName, apiData.MsgVpnName)
	AssignIfDstNotNil(&tfData.Password, apiData.Password)
	AssignIfDstNotNil(&tfData.SubscriptionManagerEnabled, apiData.SubscriptionManagerEnabled)
}

func (tfData *MsgVpnClientUsername) ToApi() *sempv2.MsgVpnClientUsername {
	return &sempv2.MsgVpnClientUsername{
		AclProfileName:    tfData.AclProfileName,
		ClientProfileName: tfData.ClientProfileName,
		ClientUsername:    tfData.ClientUsername,
		Enabled:           tfData.Enabled,
		GuaranteedEndpointPermissionOverrideEnabled: tfData.GuaranteedEndpointPermissionOverrideEnabled,
		MsgVpnName:                 tfData.MsgVpnName,
		Password:                   tfData.Password,
		SubscriptionManagerEnabled: tfData.SubscriptionManagerEnabled,
	}
}

// Terraform schema for MsgVpnClientUsername
func MsgVpnClientUsernameSchema(requiredAttributes ...string) tfsdk.Schema {
	schema := tfsdk.Schema{
		Description: "MsgVpnClientUsername",
		Attributes: map[string]tfsdk.Attribute{
			"acl_profile_name": {
				Type:        types.StringType,
				Description: "The ACL Profile of the Client Username. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"default\"`.",
				Optional:    true,
			},
			"client_profile_name": {
				Type:        types.StringType,
				Description: "The Client Profile of the Client Username. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"default\"`.",
				Optional:    true,
			},
			"client_username": {
				Type:        types.StringType,
				Description: "The name of the Client Username.",
				Optional:    true,
			},
			"enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable the Client Username. When disabled, all clients currently connected as the Client Username are disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"guaranteed_endpoint_permission_override_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable guaranteed endpoint permission override for the Client Username. When enabled all guaranteed endpoints may be accessed, modified or deleted with the same permission as the owner. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"msg_vpn_name": {
				Type:        types.StringType,
				Description: "The name of the Message VPN.",
				Optional:    true,
			},
			"password": {
				Type:        types.StringType,
				Description: "The password for the Client Username. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"subscription_manager_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable the subscription management capability of the Client Username. This is the ability to manage subscriptions on behalf of other Client Usernames. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
		},
	}

	return WithRequiredAttributes(schema, requiredAttributes)
}
