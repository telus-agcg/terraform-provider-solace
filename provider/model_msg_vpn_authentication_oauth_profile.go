package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"telusag/terraform-provider-solace/sempv2"
)

// MsgVpnAuthenticationOauthProfile struct for MsgVpnAuthenticationOauthProfile
type MsgVpnAuthenticationOauthProfile struct {
	AuthorizationGroupsClaimName          *string `tfsdk:"authorization_groups_claim_name"`
	ClientId                              *string `tfsdk:"client_id"`
	ClientRequiredType                    *string `tfsdk:"client_required_type"`
	ClientSecret                          *string `tfsdk:"client_secret"`
	ClientValidateTypeEnabled             *bool   `tfsdk:"client_validate_type_enabled"`
	DisconnectOnTokenExpirationEnabled    *bool   `tfsdk:"disconnect_on_token_expiration_enabled"`
	Enabled                               *bool   `tfsdk:"enabled"`
	EndpointDiscovery                     *string `tfsdk:"endpoint_discovery"`
	EndpointDiscoveryRefreshInterval      *int32  `tfsdk:"endpoint_discovery_refresh_interval"`
	EndpointIntrospection                 *string `tfsdk:"endpoint_introspection"`
	EndpointIntrospectionTimeout          *int32  `tfsdk:"endpoint_introspection_timeout"`
	EndpointJwks                          *string `tfsdk:"endpoint_jwks"`
	EndpointJwksRefreshInterval           *int32  `tfsdk:"endpoint_jwks_refresh_interval"`
	EndpointUserinfo                      *string `tfsdk:"endpoint_userinfo"`
	EndpointUserinfoTimeout               *int32  `tfsdk:"endpoint_userinfo_timeout"`
	Issuer                                *string `tfsdk:"issuer"`
	MqttUsernameValidateEnabled           *bool   `tfsdk:"mqtt_username_validate_enabled"`
	MsgVpnName                            *string `tfsdk:"msg_vpn_name"`
	OauthProfileName                      *string `tfsdk:"oauth_profile_name"`
	OauthRole                             *string `tfsdk:"oauth_role"`
	ResourceServerParseAccessTokenEnabled *bool   `tfsdk:"resource_server_parse_access_token_enabled"`
	ResourceServerRequiredAudience        *string `tfsdk:"resource_server_required_audience"`
	ResourceServerRequiredIssuer          *string `tfsdk:"resource_server_required_issuer"`
	ResourceServerRequiredScope           *string `tfsdk:"resource_server_required_scope"`
	ResourceServerRequiredType            *string `tfsdk:"resource_server_required_type"`
	ResourceServerValidateAudienceEnabled *bool   `tfsdk:"resource_server_validate_audience_enabled"`
	ResourceServerValidateIssuerEnabled   *bool   `tfsdk:"resource_server_validate_issuer_enabled"`
	ResourceServerValidateScopeEnabled    *bool   `tfsdk:"resource_server_validate_scope_enabled"`
	ResourceServerValidateTypeEnabled     *bool   `tfsdk:"resource_server_validate_type_enabled"`
	UsernameClaimName                     *string `tfsdk:"username_claim_name"`
}

func (tfData *MsgVpnAuthenticationOauthProfile) ToTF(apiData *sempv2.MsgVpnAuthenticationOauthProfile) {
	AssignIfDstNotNil(&tfData.AuthorizationGroupsClaimName, apiData.AuthorizationGroupsClaimName)
	AssignIfDstNotNil(&tfData.ClientId, apiData.ClientId)
	AssignIfDstNotNil(&tfData.ClientRequiredType, apiData.ClientRequiredType)
	AssignIfDstNotNil(&tfData.ClientSecret, apiData.ClientSecret)
	AssignIfDstNotNil(&tfData.ClientValidateTypeEnabled, apiData.ClientValidateTypeEnabled)
	AssignIfDstNotNil(&tfData.DisconnectOnTokenExpirationEnabled, apiData.DisconnectOnTokenExpirationEnabled)
	AssignIfDstNotNil(&tfData.Enabled, apiData.Enabled)
	AssignIfDstNotNil(&tfData.EndpointDiscovery, apiData.EndpointDiscovery)
	AssignIfDstNotNil(&tfData.EndpointDiscoveryRefreshInterval, apiData.EndpointDiscoveryRefreshInterval)
	AssignIfDstNotNil(&tfData.EndpointIntrospection, apiData.EndpointIntrospection)
	AssignIfDstNotNil(&tfData.EndpointIntrospectionTimeout, apiData.EndpointIntrospectionTimeout)
	AssignIfDstNotNil(&tfData.EndpointJwks, apiData.EndpointJwks)
	AssignIfDstNotNil(&tfData.EndpointJwksRefreshInterval, apiData.EndpointJwksRefreshInterval)
	AssignIfDstNotNil(&tfData.EndpointUserinfo, apiData.EndpointUserinfo)
	AssignIfDstNotNil(&tfData.EndpointUserinfoTimeout, apiData.EndpointUserinfoTimeout)
	AssignIfDstNotNil(&tfData.Issuer, apiData.Issuer)
	AssignIfDstNotNil(&tfData.MqttUsernameValidateEnabled, apiData.MqttUsernameValidateEnabled)
	AssignIfDstNotNil(&tfData.MsgVpnName, apiData.MsgVpnName)
	AssignIfDstNotNil(&tfData.OauthProfileName, apiData.OauthProfileName)
	AssignIfDstNotNil(&tfData.OauthRole, apiData.OauthRole)
	AssignIfDstNotNil(&tfData.ResourceServerParseAccessTokenEnabled, apiData.ResourceServerParseAccessTokenEnabled)
	AssignIfDstNotNil(&tfData.ResourceServerRequiredAudience, apiData.ResourceServerRequiredAudience)
	AssignIfDstNotNil(&tfData.ResourceServerRequiredIssuer, apiData.ResourceServerRequiredIssuer)
	AssignIfDstNotNil(&tfData.ResourceServerRequiredScope, apiData.ResourceServerRequiredScope)
	AssignIfDstNotNil(&tfData.ResourceServerRequiredType, apiData.ResourceServerRequiredType)
	AssignIfDstNotNil(&tfData.ResourceServerValidateAudienceEnabled, apiData.ResourceServerValidateAudienceEnabled)
	AssignIfDstNotNil(&tfData.ResourceServerValidateIssuerEnabled, apiData.ResourceServerValidateIssuerEnabled)
	AssignIfDstNotNil(&tfData.ResourceServerValidateScopeEnabled, apiData.ResourceServerValidateScopeEnabled)
	AssignIfDstNotNil(&tfData.ResourceServerValidateTypeEnabled, apiData.ResourceServerValidateTypeEnabled)
	AssignIfDstNotNil(&tfData.UsernameClaimName, apiData.UsernameClaimName)
}

func (tfData *MsgVpnAuthenticationOauthProfile) ToApi() *sempv2.MsgVpnAuthenticationOauthProfile {
	return &sempv2.MsgVpnAuthenticationOauthProfile{
		AuthorizationGroupsClaimName:          tfData.AuthorizationGroupsClaimName,
		ClientId:                              tfData.ClientId,
		ClientRequiredType:                    tfData.ClientRequiredType,
		ClientSecret:                          tfData.ClientSecret,
		ClientValidateTypeEnabled:             tfData.ClientValidateTypeEnabled,
		DisconnectOnTokenExpirationEnabled:    tfData.DisconnectOnTokenExpirationEnabled,
		Enabled:                               tfData.Enabled,
		EndpointDiscovery:                     tfData.EndpointDiscovery,
		EndpointDiscoveryRefreshInterval:      tfData.EndpointDiscoveryRefreshInterval,
		EndpointIntrospection:                 tfData.EndpointIntrospection,
		EndpointIntrospectionTimeout:          tfData.EndpointIntrospectionTimeout,
		EndpointJwks:                          tfData.EndpointJwks,
		EndpointJwksRefreshInterval:           tfData.EndpointJwksRefreshInterval,
		EndpointUserinfo:                      tfData.EndpointUserinfo,
		EndpointUserinfoTimeout:               tfData.EndpointUserinfoTimeout,
		Issuer:                                tfData.Issuer,
		MqttUsernameValidateEnabled:           tfData.MqttUsernameValidateEnabled,
		MsgVpnName:                            tfData.MsgVpnName,
		OauthProfileName:                      tfData.OauthProfileName,
		OauthRole:                             tfData.OauthRole,
		ResourceServerParseAccessTokenEnabled: tfData.ResourceServerParseAccessTokenEnabled,
		ResourceServerRequiredAudience:        tfData.ResourceServerRequiredAudience,
		ResourceServerRequiredIssuer:          tfData.ResourceServerRequiredIssuer,
		ResourceServerRequiredScope:           tfData.ResourceServerRequiredScope,
		ResourceServerRequiredType:            tfData.ResourceServerRequiredType,
		ResourceServerValidateAudienceEnabled: tfData.ResourceServerValidateAudienceEnabled,
		ResourceServerValidateIssuerEnabled:   tfData.ResourceServerValidateIssuerEnabled,
		ResourceServerValidateScopeEnabled:    tfData.ResourceServerValidateScopeEnabled,
		ResourceServerValidateTypeEnabled:     tfData.ResourceServerValidateTypeEnabled,
		UsernameClaimName:                     tfData.UsernameClaimName,
	}
}

// Terraform DataSource schema for MsgVpnAuthenticationOauthProfile
func MsgVpnAuthenticationOauthProfileDataSourceSchema(requiredAttributes ...string) dschema.Schema {
	schema := dschema.Schema{
		Description: "MsgVpnAuthenticationOauthProfile",
		Attributes: map[string]dschema.Attribute{
			"authorization_groups_claim_name": dschema.StringAttribute{
				Description: "The name of the groups claim. If non-empty, the specified claim will be used to determine groups for authorization. If empty, the authorizationType attribute of the Message VPN will be used to determine authorization. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"groups\"`.",
				Optional:    true,
			},
			"client_id": dschema.StringAttribute{
				Description: "The OAuth client id. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"client_required_type": dschema.StringAttribute{
				Description: "The required value for the TYP field in the ID token header. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"JWT\"`.",
				Optional:    true,
			},
			"client_secret": dschema.StringAttribute{
				Description: "The OAuth client secret. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"client_validate_type_enabled": dschema.BoolAttribute{
				Description: "Enable or disable verification of the TYP field in the ID token header. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"disconnect_on_token_expiration_enabled": dschema.BoolAttribute{
				Description: "Enable or disable the disconnection of clients when their tokens expire. Changing this value does not affect existing clients, only new client connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"enabled": dschema.BoolAttribute{
				Description: "Enable or disable the OAuth profile. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"endpoint_discovery": dschema.StringAttribute{
				Description: "The OpenID Connect discovery endpoint or OAuth Authorization Server Metadata endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"endpoint_discovery_refresh_interval": dschema.Int64Attribute{
				Description: "The number of seconds between discovery endpoint requests. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `86400`.",
				Optional:    true,
			},
			"endpoint_introspection": dschema.StringAttribute{
				Description: "The OAuth introspection endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"endpoint_introspection_timeout": dschema.Int64Attribute{
				Description: "The maximum time in seconds a token introspection request is allowed to take. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1`.",
				Optional:    true,
			},
			"endpoint_jwks": dschema.StringAttribute{
				Description: "The OAuth JWKS endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"endpoint_jwks_refresh_interval": dschema.Int64Attribute{
				Description: "The number of seconds between JWKS endpoint requests. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `86400`.",
				Optional:    true,
			},
			"endpoint_userinfo": dschema.StringAttribute{
				Description: "The OpenID Connect Userinfo endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"endpoint_userinfo_timeout": dschema.Int64Attribute{
				Description: "The maximum time in seconds a userinfo request is allowed to take. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1`.",
				Optional:    true,
			},
			"issuer": dschema.StringAttribute{
				Description: "The Issuer Identifier for the OAuth provider. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"mqtt_username_validate_enabled": dschema.BoolAttribute{
				Description: "Enable or disable whether the API provided MQTT client username will be validated against the username calculated from the token(s). When enabled, connection attempts by MQTT clients are rejected if they differ. Note that this value only applies to MQTT clients; SMF client usernames will not be validated. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"msg_vpn_name": dschema.StringAttribute{
				Description: "The name of the Message VPN.",
				Optional:    true,
			},
			"oauth_profile_name": dschema.StringAttribute{
				Description: "The name of the OAuth profile.",
				Optional:    true,
			},
			"oauth_role": dschema.StringAttribute{
				Description: "The OAuth role of the broker. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"client\"`. The allowed values and their meaning are:  <pre> \"client\" - The broker is in the OAuth client role. \"resource-server\" - The broker is in the OAuth resource server role. </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("client", "resource-server"),
				},
			},
			"resource_server_parse_access_token_enabled": dschema.BoolAttribute{
				Description: "Enable or disable parsing of the access token as a JWT. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"resource_server_required_audience": dschema.StringAttribute{
				Description: "The required audience value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"resource_server_required_issuer": dschema.StringAttribute{
				Description: "The required issuer value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"resource_server_required_scope": dschema.StringAttribute{
				Description: "A space-separated list of scopes that must be present in the scope claim. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"resource_server_required_type": dschema.StringAttribute{
				Description: "The required TYP value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"at+jwt\"`.",
				Optional:    true,
			},
			"resource_server_validate_audience_enabled": dschema.BoolAttribute{
				Description: "Enable or disable verification of the audience claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"resource_server_validate_issuer_enabled": dschema.BoolAttribute{
				Description: "Enable or disable verification of the issuer claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"resource_server_validate_scope_enabled": dschema.BoolAttribute{
				Description: "Enable or disable verification of the scope claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"resource_server_validate_type_enabled": dschema.BoolAttribute{
				Description: "Enable or disable verification of the TYP field in the access token header. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"username_claim_name": dschema.StringAttribute{
				Description: "The name of the username claim. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"sub\"`.",
				Optional:    true,
			},
		},
	}

	return schema
}

// Terraform Resource schema for MsgVpnAuthenticationOauthProfile
func MsgVpnAuthenticationOauthProfileResourceSchema(requiredAttributes ...string) rschema.Schema {
	schema := rschema.Schema{
		Description: "MsgVpnAuthenticationOauthProfile",
		Attributes: map[string]rschema.Attribute{
			"authorization_groups_claim_name": rschema.StringAttribute{
				Description: "The name of the groups claim. If non-empty, the specified claim will be used to determine groups for authorization. If empty, the authorizationType attribute of the Message VPN will be used to determine authorization. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"groups\"`.",
				Optional:    true,
			},
			"client_id": rschema.StringAttribute{
				Description: "The OAuth client id. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"client_required_type": rschema.StringAttribute{
				Description: "The required value for the TYP field in the ID token header. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"JWT\"`.",
				Optional:    true,
			},
			"client_secret": rschema.StringAttribute{
				Description: "The OAuth client secret. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"client_validate_type_enabled": rschema.BoolAttribute{
				Description: "Enable or disable verification of the TYP field in the ID token header. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"disconnect_on_token_expiration_enabled": rschema.BoolAttribute{
				Description: "Enable or disable the disconnection of clients when their tokens expire. Changing this value does not affect existing clients, only new client connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"enabled": rschema.BoolAttribute{
				Description: "Enable or disable the OAuth profile. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"endpoint_discovery": rschema.StringAttribute{
				Description: "The OpenID Connect discovery endpoint or OAuth Authorization Server Metadata endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"endpoint_discovery_refresh_interval": rschema.Int64Attribute{
				Description: "The number of seconds between discovery endpoint requests. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `86400`.",
				Optional:    true,
			},
			"endpoint_introspection": rschema.StringAttribute{
				Description: "The OAuth introspection endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"endpoint_introspection_timeout": rschema.Int64Attribute{
				Description: "The maximum time in seconds a token introspection request is allowed to take. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1`.",
				Optional:    true,
			},
			"endpoint_jwks": rschema.StringAttribute{
				Description: "The OAuth JWKS endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"endpoint_jwks_refresh_interval": rschema.Int64Attribute{
				Description: "The number of seconds between JWKS endpoint requests. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `86400`.",
				Optional:    true,
			},
			"endpoint_userinfo": rschema.StringAttribute{
				Description: "The OpenID Connect Userinfo endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"endpoint_userinfo_timeout": rschema.Int64Attribute{
				Description: "The maximum time in seconds a userinfo request is allowed to take. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1`.",
				Optional:    true,
			},
			"issuer": rschema.StringAttribute{
				Description: "The Issuer Identifier for the OAuth provider. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"mqtt_username_validate_enabled": rschema.BoolAttribute{
				Description: "Enable or disable whether the API provided MQTT client username will be validated against the username calculated from the token(s). When enabled, connection attempts by MQTT clients are rejected if they differ. Note that this value only applies to MQTT clients; SMF client usernames will not be validated. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"msg_vpn_name": rschema.StringAttribute{
				Description: "The name of the Message VPN.",
				Optional:    true,
			},
			"oauth_profile_name": rschema.StringAttribute{
				Description: "The name of the OAuth profile.",
				Optional:    true,
			},
			"oauth_role": rschema.StringAttribute{
				Description: "The OAuth role of the broker. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"client\"`. The allowed values and their meaning are:  <pre> \"client\" - The broker is in the OAuth client role. \"resource-server\" - The broker is in the OAuth resource server role. </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("client", "resource-server"),
				},
			},
			"resource_server_parse_access_token_enabled": rschema.BoolAttribute{
				Description: "Enable or disable parsing of the access token as a JWT. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"resource_server_required_audience": rschema.StringAttribute{
				Description: "The required audience value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"resource_server_required_issuer": rschema.StringAttribute{
				Description: "The required issuer value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"resource_server_required_scope": rschema.StringAttribute{
				Description: "A space-separated list of scopes that must be present in the scope claim. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"resource_server_required_type": rschema.StringAttribute{
				Description: "The required TYP value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"at+jwt\"`.",
				Optional:    true,
			},
			"resource_server_validate_audience_enabled": rschema.BoolAttribute{
				Description: "Enable or disable verification of the audience claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"resource_server_validate_issuer_enabled": rschema.BoolAttribute{
				Description: "Enable or disable verification of the issuer claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"resource_server_validate_scope_enabled": rschema.BoolAttribute{
				Description: "Enable or disable verification of the scope claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"resource_server_validate_type_enabled": rschema.BoolAttribute{
				Description: "Enable or disable verification of the TYP field in the access token header. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"username_claim_name": rschema.StringAttribute{
				Description: "The name of the username claim. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"sub\"`.",
				Optional:    true,
			},
		},
	}

	return schema
}
