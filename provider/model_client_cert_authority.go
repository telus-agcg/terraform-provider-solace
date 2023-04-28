package provider

import (
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"telusag/terraform-provider-solace/sempv2"
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

// Terraform DataSource schema for ClientCertAuthority
func ClientCertAuthorityDataSourceSchema(requiredAttributes ...string) dschema.Schema {
	schema := dschema.Schema{
		Description: "ClientCertAuthority",
		Attributes: map[string]dschema.Attribute{
			"cert_authority_name": dschema.StringAttribute{
				Description: "The name of the Certificate Authority.",
				Optional:    true,
			},
			"cert_content": dschema.StringAttribute{
				Description: "The PEM formatted content for the trusted root certificate of a client Certificate Authority. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"crl_day_list": dschema.StringAttribute{
				Description: "The scheduled CRL refresh day(s), specified as \"daily\" or a comma-separated list of days. Days must be specified as \"Sun\", \"Mon\", \"Tue\", \"Wed\", \"Thu\", \"Fri\", or \"Sat\", with no spaces, and in sorted order from Sunday to Saturday. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"daily\"`.",
				Optional:    true,
			},
			"crl_time_list": dschema.StringAttribute{
				Description: "The scheduled CRL refresh time(s), specified as \"hourly\" or a comma-separated list of 24-hour times in the form hh:mm, or h:mm. There must be no spaces, and times (up to 4) must be in sorted order from 0:00 to 23:59. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"3:00\"`.",
				Optional:    true,
			},
			"crl_url": dschema.StringAttribute{
				Description: "The URL for the CRL source. This is a required attribute for CRL to be operational and the URL must be complete with http:// included. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"ocsp_non_responder_cert_enabled": dschema.BoolAttribute{
				Description: "Enable or disable allowing a non-responder certificate to sign an OCSP response. Typically used with an OCSP override URL in cases where a single certificate is used to sign client certificates and OCSP responses. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"ocsp_override_url": dschema.StringAttribute{
				Description: "The OCSP responder URL to use for overriding the one supplied in the client certificate. The URL must be complete with http:// included. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"ocsp_timeout": dschema.Int64Attribute{
				Description: "The timeout in seconds to receive a response from the OCSP responder after sending a request or making the initial connection attempt. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `5`.",
				Optional:    true,
			},
			"revocation_check_enabled": dschema.BoolAttribute{
				Description: "Enable or disable Certificate Authority revocation checking. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.",
				Optional:    true,
			},
		},
	}

	return schema
}

// Terraform Resource schema for ClientCertAuthority
func ClientCertAuthorityResourceSchema(requiredAttributes ...string) rschema.Schema {
	schema := rschema.Schema{
		Description: "ClientCertAuthority",
		Attributes: map[string]rschema.Attribute{
			"cert_authority_name": rschema.StringAttribute{
				Description: "The name of the Certificate Authority.",
				Optional:    true,
			},
			"cert_content": rschema.StringAttribute{
				Description: "The PEM formatted content for the trusted root certificate of a client Certificate Authority. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"crl_day_list": rschema.StringAttribute{
				Description: "The scheduled CRL refresh day(s), specified as \"daily\" or a comma-separated list of days. Days must be specified as \"Sun\", \"Mon\", \"Tue\", \"Wed\", \"Thu\", \"Fri\", or \"Sat\", with no spaces, and in sorted order from Sunday to Saturday. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"daily\"`.",
				Optional:    true,
			},
			"crl_time_list": rschema.StringAttribute{
				Description: "The scheduled CRL refresh time(s), specified as \"hourly\" or a comma-separated list of 24-hour times in the form hh:mm, or h:mm. There must be no spaces, and times (up to 4) must be in sorted order from 0:00 to 23:59. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"3:00\"`.",
				Optional:    true,
			},
			"crl_url": rschema.StringAttribute{
				Description: "The URL for the CRL source. This is a required attribute for CRL to be operational and the URL must be complete with http:// included. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"ocsp_non_responder_cert_enabled": rschema.BoolAttribute{
				Description: "Enable or disable allowing a non-responder certificate to sign an OCSP response. Typically used with an OCSP override URL in cases where a single certificate is used to sign client certificates and OCSP responses. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"ocsp_override_url": rschema.StringAttribute{
				Description: "The OCSP responder URL to use for overriding the one supplied in the client certificate. The URL must be complete with http:// included. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"ocsp_timeout": rschema.Int64Attribute{
				Description: "The timeout in seconds to receive a response from the OCSP responder after sending a request or making the initial connection attempt. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `5`.",
				Optional:    true,
			},
			"revocation_check_enabled": rschema.BoolAttribute{
				Description: "Enable or disable Certificate Authority revocation checking. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.",
				Optional:    true,
			},
		},
	}

	return schema
}
