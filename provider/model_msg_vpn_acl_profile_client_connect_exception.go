package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"telusag/terraform-provider-solace/sempv2"
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

// Terraform DataSource schema for MsgVpnAclProfileClientConnectException
func MsgVpnAclProfileClientConnectExceptionDataSourceSchema(requiredAttributes ...string) dschema.Schema {
	schema := dschema.Schema{
		Description: "MsgVpnAclProfileClientConnectException",
		Attributes: map[string]dschema.Attribute{
			"acl_profile_name": dschema.StringAttribute{
				Description: "The name of the ACL Profile.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
				},
			},
			"client_connect_exception_address": dschema.StringAttribute{
				Description: "The IP address/netmask of the client connect exception in CIDR form.",
				Optional:    true,
			},
			"msg_vpn_name": dschema.StringAttribute{
				Description: "The name of the Message VPN.",
				Optional:    true,
			},
		},
	}

	return schema
}

// Terraform Resource schema for MsgVpnAclProfileClientConnectException
func MsgVpnAclProfileClientConnectExceptionResourceSchema(requiredAttributes ...string) rschema.Schema {
	schema := rschema.Schema{
		Description: "MsgVpnAclProfileClientConnectException",
		Attributes: map[string]rschema.Attribute{
			"acl_profile_name": rschema.StringAttribute{
				Description: "The name of the ACL Profile.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(32),
				},
			},
			"client_connect_exception_address": rschema.StringAttribute{
				Description: "The IP address/netmask of the client connect exception in CIDR form.",
				Optional:    true,
			},
			"msg_vpn_name": rschema.StringAttribute{
				Description: "The name of the Message VPN.",
				Optional:    true,
			},
		},
	}

	return schema
}
