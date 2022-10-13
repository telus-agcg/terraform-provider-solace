package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ClientCertAuthority struct for ClientCertAuthority
type ClientCertAuthority struct {
	CertAuthorityName           *string `tfsdk:"cert_authority_name"`
	CertContent                 *string `tfsdk:"cert_content"`
	CrlDayList                  *string `tfsdk:"crl_day_list"`
	CrlTimeList                 *string `tfsdk:"crl_time_list"`
	CrlUrl                      *string `tfsdk:"crl_url"`
	OcspNonResponderCertEnabled *bool   `tfsdk:"ocsp_non_responder_cert_enabled"`
	OcspOverrideUrl             *string `tfsdk:"ocsp_override_url"`
	OcspTimeout                 *int64  `tfsdk:"ocsp_timeout"`
	RevocationCheckEnabled      *bool   `tfsdk:"revocation_check_enabled"`
}

func (tfData *ClientCertAuthority) ToTF(apiData *sempv2.ClientCertAuthority) {
	AssignIfDstNotNil(&tfData.CertAuthorityName, apiData.CertAuthorityName)
	AssignIfDstNotNil(&tfData.CertContent, apiData.CertContent)
	AssignIfDstNotNil(&tfData.CrlDayList, apiData.CrlDayList)
	AssignIfDstNotNil(&tfData.CrlTimeList, apiData.CrlTimeList)
	AssignIfDstNotNil(&tfData.CrlUrl, apiData.CrlUrl)
	AssignIfDstNotNil(&tfData.OcspNonResponderCertEnabled, apiData.OcspNonResponderCertEnabled)
	AssignIfDstNotNil(&tfData.OcspOverrideUrl, apiData.OcspOverrideUrl)
	AssignIfDstNotNil(&tfData.OcspTimeout, apiData.OcspTimeout)
	AssignIfDstNotNil(&tfData.RevocationCheckEnabled, apiData.RevocationCheckEnabled)
}

func (tfData *ClientCertAuthority) ToApi() *sempv2.ClientCertAuthority {
	return &sempv2.ClientCertAuthority{
		CertAuthorityName:           tfData.CertAuthorityName,
		CertContent:                 tfData.CertContent,
		CrlDayList:                  tfData.CrlDayList,
		CrlTimeList:                 tfData.CrlTimeList,
		CrlUrl:                      tfData.CrlUrl,
		OcspNonResponderCertEnabled: tfData.OcspNonResponderCertEnabled,
		OcspOverrideUrl:             tfData.OcspOverrideUrl,
		OcspTimeout:                 tfData.OcspTimeout,
		RevocationCheckEnabled:      tfData.RevocationCheckEnabled,
	}
}

// Terraform schema for ClientCertAuthority
func ClientCertAuthoritySchema(requiredAttributes ...string) tfsdk.Schema {
	schema := tfsdk.Schema{
		Description: "ClientCertAuthority",
		Attributes: map[string]tfsdk.Attribute{
			"cert_authority_name": {
				Type:        types.StringType,
				Description: "The name of the Certificate Authority.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"cert_content": {
				Type:        types.StringType,
				Description: "The PEM formatted content for the trusted root certificate of a client Certificate Authority. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"crl_day_list": {
				Type:        types.StringType,
				Description: "The scheduled CRL refresh day(s), specified as \"daily\" or a comma-separated list of days. Days must be specified as \"Sun\", \"Mon\", \"Tue\", \"Wed\", \"Thu\", \"Fri\", or \"Sat\", with no spaces, and in sorted order from Sunday to Saturday. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"daily\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"crl_time_list": {
				Type:        types.StringType,
				Description: "The scheduled CRL refresh time(s), specified as \"hourly\" or a comma-separated list of 24-hour times in the form hh:mm, or h:mm. There must be no spaces, and times (up to 4) must be in sorted order from 0:00 to 23:59. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"3:00\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"crl_url": {
				Type:        types.StringType,
				Description: "The URL for the CRL source. This is a required attribute for CRL to be operational and the URL must be complete with http:// included. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"ocsp_non_responder_cert_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable allowing a non-responder certificate to sign an OCSP response. Typically used with an OCSP override URL in cases where a single certificate is used to sign client certificates and OCSP responses. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"ocsp_override_url": {
				Type:        types.StringType,
				Description: "The OCSP responder URL to use for overriding the one supplied in the client certificate. The URL must be complete with http:// included. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"ocsp_timeout": {
				Type:        types.Int64Type,
				Description: "The timeout in seconds to receive a response from the OCSP responder after sending a request or making the initial connection attempt. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `5`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"revocation_check_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable Certificate Authority revocation checking. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
		},
	}

	return WithRequiredAttributes(schema, requiredAttributes)
}
