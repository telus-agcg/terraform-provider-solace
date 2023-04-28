package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// MsgVpnAclProfileClientConnectException struct for MsgVpnAclProfileClientConnectException
type MsgVpnAclProfileClientConnectException struct {
	AclProfileName                *string `tfsdk:"acl_profile_name"`
	ClientConnectExceptionAddress *string `tfsdk:"client_connect_exception_address"`
	MsgVpnName                    *string `tfsdk:"msg_vpn_name"`
}

func (tfData *MsgVpnAclProfileClientConnectException) ToTF(apiData *sempv2.MsgVpnAclProfileClientConnectException) {
	AssignIfDstNotNil(&tfData.AclProfileName, apiData.AclProfileName)
	AssignIfDstNotNil(&tfData.ClientConnectExceptionAddress, apiData.ClientConnectExceptionAddress)
	AssignIfDstNotNil(&tfData.MsgVpnName, apiData.MsgVpnName)
}

func (tfData *MsgVpnAclProfileClientConnectException) ToApi() *sempv2.MsgVpnAclProfileClientConnectException {
	return &sempv2.MsgVpnAclProfileClientConnectException{
		AclProfileName:                tfData.AclProfileName,
		ClientConnectExceptionAddress: tfData.ClientConnectExceptionAddress,
		MsgVpnName:                    tfData.MsgVpnName,
	}
}

// Terraform Resource schema for MsgVpnAclProfileClientConnectException
func MsgVpnAclProfileClientConnectExceptionResourceSchema(requiredAttributes ...string) schema.Schema {
	schema := schema.Schema{
		Description: "MsgVpnAclProfileClientConnectException",
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
			"client_connect_exception_address": schema.StringAttribute{
				Description: "The IP address/netmask of the client connect exception in CIDR form.",
				Required:    contains(requiredAttributes, "client_connect_exception_address"),
				Optional:    !contains(requiredAttributes, "client_connect_exception_address"),

				PlanModifiers: StringPlanModifiersFor("client_connect_exception_address", requiredAttributes),
			},
			"msg_vpn_name": schema.StringAttribute{
				Description: "The name of the Message VPN.",
				Required:    contains(requiredAttributes, "msg_vpn_name"),
				Optional:    !contains(requiredAttributes, "msg_vpn_name"),

				PlanModifiers: StringPlanModifiersFor("msg_vpn_name", requiredAttributes),
			},
		},
	}

	return schema
}
