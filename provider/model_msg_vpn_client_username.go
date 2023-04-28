package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
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

// Terraform Resource schema for MsgVpnClientUsername
func MsgVpnClientUsernameResourceSchema(requiredAttributes ...string) schema.Schema {
	schema := schema.Schema{
		Description: "MsgVpnClientUsername",
		Attributes: map[string]schema.Attribute{
			"acl_profile_name": schema.StringAttribute{
				Description: "The ACL Profile of the Client Username. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"default\"`.",
				Required:    contains(requiredAttributes, "acl_profile_name"),
				Optional:    !contains(requiredAttributes, "acl_profile_name"),

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
				},
				PlanModifiers: StringPlanModifiersFor("acl_profile_name", requiredAttributes),
			},
			"client_profile_name": schema.StringAttribute{
				Description: "The Client Profile of the Client Username. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"default\"`.",
				Required:    contains(requiredAttributes, "client_profile_name"),
				Optional:    !contains(requiredAttributes, "client_profile_name"),

				PlanModifiers: StringPlanModifiersFor("client_profile_name", requiredAttributes),
			},
			"client_username": schema.StringAttribute{
				Description: "The name of the Client Username.",
				Required:    contains(requiredAttributes, "client_username"),
				Optional:    !contains(requiredAttributes, "client_username"),

				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[[:print:]]{1,189}$"), "Does not match pattern '^[[:print:]]{1,189}$'"),
				},
				PlanModifiers: StringPlanModifiersFor("client_username", requiredAttributes),
			},
			"enabled": schema.BoolAttribute{
				Description: "Enable or disable the Client Username. When disabled, all clients currently connected as the Client Username are disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "enabled"),
				Optional:    !contains(requiredAttributes, "enabled"),
			},
			"guaranteed_endpoint_permission_override_enabled": schema.BoolAttribute{
				Description: "Enable or disable guaranteed endpoint permission override for the Client Username. When enabled all guaranteed endpoints may be accessed, modified or deleted with the same permission as the owner. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "guaranteed_endpoint_permission_override_enabled"),
				Optional:    !contains(requiredAttributes, "guaranteed_endpoint_permission_override_enabled"),
			},
			"msg_vpn_name": schema.StringAttribute{
				Description: "The name of the Message VPN.",
				Required:    contains(requiredAttributes, "msg_vpn_name"),
				Optional:    !contains(requiredAttributes, "msg_vpn_name"),

				PlanModifiers: StringPlanModifiersFor("msg_vpn_name", requiredAttributes),
			},
			"password": schema.StringAttribute{
				Description: "The password for the Client Username. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "password"),
				Optional:    !contains(requiredAttributes, "password"),

				PlanModifiers: StringPlanModifiersFor("password", requiredAttributes),
			},
			"subscription_manager_enabled": schema.BoolAttribute{
				Description: "Enable or disable the subscription management capability of the Client Username. This is the ability to manage subscriptions on behalf of other Client Usernames. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "subscription_manager_enabled"),
				Optional:    !contains(requiredAttributes, "subscription_manager_enabled"),
			},
		},
	}

	return schema
}
