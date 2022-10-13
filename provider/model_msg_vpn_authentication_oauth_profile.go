package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// Terraform schema for MsgVpnAuthenticationOauthProfile
func MsgVpnAuthenticationOauthProfileSchema(requiredAttributes ...string) tfsdk.Schema {
	schema := tfsdk.Schema{
		Description: "MsgVpnAuthenticationOauthProfile",
		Attributes: map[string]tfsdk.Attribute{
			"authorization_groups_claim_name": {
				Type:        types.StringType,
				Description: "The name of the groups claim. If non-empty, the specified claim will be used to determine groups for authorization. If empty, the authorizationType attribute of the Message VPN will be used to determine authorization. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"groups\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"client_id": {
				Type:        types.StringType,
				Description: "The OAuth client id. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"client_required_type": {
				Type:        types.StringType,
				Description: "The required value for the TYP field in the ID token header. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"JWT\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"client_secret": {
				Type:        types.StringType,
				Description: "The OAuth client secret. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"client_validate_type_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable verification of the TYP field in the ID token header. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"disconnect_on_token_expiration_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable the disconnection of clients when their tokens expire. Changing this value does not affect existing clients, only new client connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable the OAuth profile. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"endpoint_discovery": {
				Type:        types.StringType,
				Description: "The OpenID Connect discovery endpoint or OAuth Authorization Server Metadata endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"endpoint_discovery_refresh_interval": {
				Type:        types.Int64Type,
				Description: "The number of seconds between discovery endpoint requests. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `86400`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"endpoint_introspection": {
				Type:        types.StringType,
				Description: "The OAuth introspection endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"endpoint_introspection_timeout": {
				Type:        types.Int64Type,
				Description: "The maximum time in seconds a token introspection request is allowed to take. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"endpoint_jwks": {
				Type:        types.StringType,
				Description: "The OAuth JWKS endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"endpoint_jwks_refresh_interval": {
				Type:        types.Int64Type,
				Description: "The number of seconds between JWKS endpoint requests. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `86400`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"endpoint_userinfo": {
				Type:        types.StringType,
				Description: "The OpenID Connect Userinfo endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"endpoint_userinfo_timeout": {
				Type:        types.Int64Type,
				Description: "The maximum time in seconds a userinfo request is allowed to take. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"issuer": {
				Type:        types.StringType,
				Description: "The Issuer Identifier for the OAuth provider. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"mqtt_username_validate_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable whether the API provided MQTT client username will be validated against the username calculated from the token(s). When enabled, connection attempts by MQTT clients are rejected if they differ. Note that this value only applies to MQTT clients; SMF client usernames will not be validated. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"msg_vpn_name": {
				Type:        types.StringType,
				Description: "The name of the Message VPN.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"oauth_profile_name": {
				Type:        types.StringType,
				Description: "The name of the OAuth profile.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"oauth_role": {
				Type:        types.StringType,
				Description: "The OAuth role of the broker. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"client\"`. The allowed values and their meaning are:  <pre> \"client\" - The broker is in the OAuth client role. \"resource-server\" - The broker is in the OAuth resource server role. </pre> ",
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf("client", "resource-server"),
				},
			},
			"resource_server_parse_access_token_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable parsing of the access token as a JWT. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"resource_server_required_audience": {
				Type:        types.StringType,
				Description: "The required audience value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"resource_server_required_issuer": {
				Type:        types.StringType,
				Description: "The required issuer value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"resource_server_required_scope": {
				Type:        types.StringType,
				Description: "A space-separated list of scopes that must be present in the scope claim. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"resource_server_required_type": {
				Type:        types.StringType,
				Description: "The required TYP value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"at+jwt\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"resource_server_validate_audience_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable verification of the audience claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"resource_server_validate_issuer_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable verification of the issuer claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"resource_server_validate_scope_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable verification of the scope claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"resource_server_validate_type_enabled": {
				Type:        types.BoolType,
				Description: "Enable or disable verification of the TYP field in the access token header. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
			"username_claim_name": {
				Type:        types.StringType,
				Description: "The name of the username claim. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"sub\"`.",
				Optional:    true,
				Validators:  []tfsdk.AttributeValidator{},
			},
		},
	}

	return WithRequiredAttributes(schema, requiredAttributes)
}
