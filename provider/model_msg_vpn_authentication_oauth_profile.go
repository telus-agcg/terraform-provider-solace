package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

// Terraform Resource schema for MsgVpnAuthenticationOauthProfile
func MsgVpnAuthenticationOauthProfileResourceSchema(requiredAttributes ...string) schema.Schema {
	schema := schema.Schema{
		Description: "MsgVpnAuthenticationOauthProfile",
		Attributes: map[string]schema.Attribute{
			"authorization_groups_claim_name": schema.StringAttribute{
				Description: "The name of the groups claim. If non-empty, the specified claim will be used to determine groups for authorization. If empty, the authorizationType attribute of the Message VPN will be used to determine authorization. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"groups\"`.",
				Required:    contains(requiredAttributes, "authorization_groups_claim_name"),
				Optional:    !contains(requiredAttributes, "authorization_groups_claim_name"),

				PlanModifiers: StringPlanModifiersFor("authorization_groups_claim_name", requiredAttributes),
			},
			"client_id": schema.StringAttribute{
				Description: "The OAuth client id. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "client_id"),
				Optional:    !contains(requiredAttributes, "client_id"),

				PlanModifiers: StringPlanModifiersFor("client_id", requiredAttributes),
			},
			"client_required_type": schema.StringAttribute{
				Description: "The required value for the TYP field in the ID token header. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"JWT\"`.",
				Required:    contains(requiredAttributes, "client_required_type"),
				Optional:    !contains(requiredAttributes, "client_required_type"),

				PlanModifiers: StringPlanModifiersFor("client_required_type", requiredAttributes),
			},
			"client_secret": schema.StringAttribute{
				Description: "The OAuth client secret. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "client_secret"),
				Optional:    !contains(requiredAttributes, "client_secret"),

				PlanModifiers: StringPlanModifiersFor("client_secret", requiredAttributes),
			},
			"client_validate_type_enabled": schema.BoolAttribute{
				Description: "Enable or disable verification of the TYP field in the ID token header. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "client_validate_type_enabled"),
				Optional:    !contains(requiredAttributes, "client_validate_type_enabled"),
			},
			"disconnect_on_token_expiration_enabled": schema.BoolAttribute{
				Description: "Enable or disable the disconnection of clients when their tokens expire. Changing this value does not affect existing clients, only new client connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "disconnect_on_token_expiration_enabled"),
				Optional:    !contains(requiredAttributes, "disconnect_on_token_expiration_enabled"),
			},
			"enabled": schema.BoolAttribute{
				Description: "Enable or disable the OAuth profile. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "enabled"),
				Optional:    !contains(requiredAttributes, "enabled"),
			},
			"endpoint_discovery": schema.StringAttribute{
				Description: "The OpenID Connect discovery endpoint or OAuth Authorization Server Metadata endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "endpoint_discovery"),
				Optional:    !contains(requiredAttributes, "endpoint_discovery"),

				PlanModifiers: StringPlanModifiersFor("endpoint_discovery", requiredAttributes),
			},
			"endpoint_discovery_refresh_interval": schema.Int64Attribute{
				Description: "The number of seconds between discovery endpoint requests. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `86400`.",
				Required:    contains(requiredAttributes, "endpoint_discovery_refresh_interval"),
				Optional:    !contains(requiredAttributes, "endpoint_discovery_refresh_interval"),
			},
			"endpoint_introspection": schema.StringAttribute{
				Description: "The OAuth introspection endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "endpoint_introspection"),
				Optional:    !contains(requiredAttributes, "endpoint_introspection"),

				PlanModifiers: StringPlanModifiersFor("endpoint_introspection", requiredAttributes),
			},
			"endpoint_introspection_timeout": schema.Int64Attribute{
				Description: "The maximum time in seconds a token introspection request is allowed to take. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1`.",
				Required:    contains(requiredAttributes, "endpoint_introspection_timeout"),
				Optional:    !contains(requiredAttributes, "endpoint_introspection_timeout"),
			},
			"endpoint_jwks": schema.StringAttribute{
				Description: "The OAuth JWKS endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "endpoint_jwks"),
				Optional:    !contains(requiredAttributes, "endpoint_jwks"),

				PlanModifiers: StringPlanModifiersFor("endpoint_jwks", requiredAttributes),
			},
			"endpoint_jwks_refresh_interval": schema.Int64Attribute{
				Description: "The number of seconds between JWKS endpoint requests. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `86400`.",
				Required:    contains(requiredAttributes, "endpoint_jwks_refresh_interval"),
				Optional:    !contains(requiredAttributes, "endpoint_jwks_refresh_interval"),
			},
			"endpoint_userinfo": schema.StringAttribute{
				Description: "The OpenID Connect Userinfo endpoint. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "endpoint_userinfo"),
				Optional:    !contains(requiredAttributes, "endpoint_userinfo"),

				PlanModifiers: StringPlanModifiersFor("endpoint_userinfo", requiredAttributes),
			},
			"endpoint_userinfo_timeout": schema.Int64Attribute{
				Description: "The maximum time in seconds a userinfo request is allowed to take. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1`.",
				Required:    contains(requiredAttributes, "endpoint_userinfo_timeout"),
				Optional:    !contains(requiredAttributes, "endpoint_userinfo_timeout"),
			},
			"issuer": schema.StringAttribute{
				Description: "The Issuer Identifier for the OAuth provider. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "issuer"),
				Optional:    !contains(requiredAttributes, "issuer"),

				PlanModifiers: StringPlanModifiersFor("issuer", requiredAttributes),
			},
			"mqtt_username_validate_enabled": schema.BoolAttribute{
				Description: "Enable or disable whether the API provided MQTT client username will be validated against the username calculated from the token(s). When enabled, connection attempts by MQTT clients are rejected if they differ. Note that this value only applies to MQTT clients; SMF client usernames will not be validated. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "mqtt_username_validate_enabled"),
				Optional:    !contains(requiredAttributes, "mqtt_username_validate_enabled"),
			},
			"msg_vpn_name": schema.StringAttribute{
				Description: "The name of the Message VPN.",
				Required:    contains(requiredAttributes, "msg_vpn_name"),
				Optional:    !contains(requiredAttributes, "msg_vpn_name"),

				PlanModifiers: StringPlanModifiersFor("msg_vpn_name", requiredAttributes),
			},
			"oauth_profile_name": schema.StringAttribute{
				Description: "The name of the OAuth profile.",
				Required:    contains(requiredAttributes, "oauth_profile_name"),
				Optional:    !contains(requiredAttributes, "oauth_profile_name"),

				PlanModifiers: StringPlanModifiersFor("oauth_profile_name", requiredAttributes),
			},
			"oauth_role": schema.StringAttribute{
				Description: "The OAuth role of the broker. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"client\"`. The allowed values and their meaning are:  <pre> \"client\" - The broker is in the OAuth client role. \"resource-server\" - The broker is in the OAuth resource server role. </pre> ",
				Required:    contains(requiredAttributes, "oauth_role"),
				Optional:    !contains(requiredAttributes, "oauth_role"),

				Validators: []validator.String{
					stringvalidator.OneOf("client", "resource-server"),
				},
				PlanModifiers: StringPlanModifiersFor("oauth_role", requiredAttributes),
			},
			"resource_server_parse_access_token_enabled": schema.BoolAttribute{
				Description: "Enable or disable parsing of the access token as a JWT. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "resource_server_parse_access_token_enabled"),
				Optional:    !contains(requiredAttributes, "resource_server_parse_access_token_enabled"),
			},
			"resource_server_required_audience": schema.StringAttribute{
				Description: "The required audience value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "resource_server_required_audience"),
				Optional:    !contains(requiredAttributes, "resource_server_required_audience"),

				PlanModifiers: StringPlanModifiersFor("resource_server_required_audience", requiredAttributes),
			},
			"resource_server_required_issuer": schema.StringAttribute{
				Description: "The required issuer value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "resource_server_required_issuer"),
				Optional:    !contains(requiredAttributes, "resource_server_required_issuer"),

				PlanModifiers: StringPlanModifiersFor("resource_server_required_issuer", requiredAttributes),
			},
			"resource_server_required_scope": schema.StringAttribute{
				Description: "A space-separated list of scopes that must be present in the scope claim. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "resource_server_required_scope"),
				Optional:    !contains(requiredAttributes, "resource_server_required_scope"),

				PlanModifiers: StringPlanModifiersFor("resource_server_required_scope", requiredAttributes),
			},
			"resource_server_required_type": schema.StringAttribute{
				Description: "The required TYP value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"at+jwt\"`.",
				Required:    contains(requiredAttributes, "resource_server_required_type"),
				Optional:    !contains(requiredAttributes, "resource_server_required_type"),

				PlanModifiers: StringPlanModifiersFor("resource_server_required_type", requiredAttributes),
			},
			"resource_server_validate_audience_enabled": schema.BoolAttribute{
				Description: "Enable or disable verification of the audience claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "resource_server_validate_audience_enabled"),
				Optional:    !contains(requiredAttributes, "resource_server_validate_audience_enabled"),
			},
			"resource_server_validate_issuer_enabled": schema.BoolAttribute{
				Description: "Enable or disable verification of the issuer claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "resource_server_validate_issuer_enabled"),
				Optional:    !contains(requiredAttributes, "resource_server_validate_issuer_enabled"),
			},
			"resource_server_validate_scope_enabled": schema.BoolAttribute{
				Description: "Enable or disable verification of the scope claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "resource_server_validate_scope_enabled"),
				Optional:    !contains(requiredAttributes, "resource_server_validate_scope_enabled"),
			},
			"resource_server_validate_type_enabled": schema.BoolAttribute{
				Description: "Enable or disable verification of the TYP field in the access token header. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "resource_server_validate_type_enabled"),
				Optional:    !contains(requiredAttributes, "resource_server_validate_type_enabled"),
			},
			"username_claim_name": schema.StringAttribute{
				Description: "The name of the username claim. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"sub\"`.",
				Required:    contains(requiredAttributes, "username_claim_name"),
				Optional:    !contains(requiredAttributes, "username_claim_name"),

				PlanModifiers: StringPlanModifiersFor("username_claim_name", requiredAttributes),
			},
		},
	}

	return schema
}
