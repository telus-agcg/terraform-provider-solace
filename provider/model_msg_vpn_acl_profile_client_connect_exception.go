package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// Terraform schema for MsgVpnAclProfileClientConnectException
func MsgVpnAclProfileClientConnectExceptionSchema(requiredAttributes ...string) tfsdk.Schema {
	schema := tfsdk.Schema{
		Description: "MsgVpnAclProfileClientConnectException",
		Attributes: map[string]tfsdk.Attribute{
			"acl_profile_name": {
				Type:        types.StringType,
				Description: "The name of the ACL Profile.",
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
				},
			},
			"client_connect_exception_address": {
				Type:        types.StringType,
				Description: "The IP address/netmask of the client connect exception in CIDR form.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"msg_vpn_name": {
				Type:        types.StringType,
				Description: "The name of the Message VPN.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
		},
	}

	return WithRequiredAttributes(schema, requiredAttributes)
}
