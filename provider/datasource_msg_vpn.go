package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func NewMsgVpnDataSource() datasource.DataSource {
	return &dataSource[MsgVpn]{spds: &msgVpnDataSource{}}
}

var _ solaceProviderDataSource[MsgVpn] = &msgVpnDataSource{}

type msgVpnDataSource struct {
	*solaceProvider
}

func (r msgVpnDataSource) Name() string {
	return "msgvpn"
}

func (r msgVpnDataSource) Schema() schema.Schema {
	return schema.Schema{
		Description: "MsgVpn",
		Attributes: map[string]schema.Attribute{
			"alias": schema.StringAttribute{
				Description: "The name of another Message VPN which this Message VPN is an alias for. When this Message VPN is enabled, the alias has no effect. When this Message VPN is disabled, Clients (but not Bridges and routing Links) logging into this Message VPN are automatically logged in to the other Message VPN, and authentication and authorization take place in the context of the other Message VPN.  Aliases may form a non-circular chain, cascading one to the next. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.14.",
				Optional:    true,
			},
			"authentication_basic_enabled": schema.BoolAttribute{
				Description: "Enable or disable basic authentication for clients connecting to the Message VPN. Basic authentication is authentication that involves the use of a username and password to prove identity. If a user provides credentials for a different authentication scheme, this setting is not applicable. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"authentication_basic_profile_name": schema.StringAttribute{
				Description: "The name of the RADIUS or LDAP Profile to use for basic authentication. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"default\"`.",
				Optional:    true,
			},
			"authentication_basic_radius_domain": schema.StringAttribute{
				Description: "The RADIUS domain to use for basic authentication. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"authentication_basic_type": schema.StringAttribute{
				Description: "The type of basic authentication to use for clients connecting to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"radius\"`. The allowed values and their meaning are:  <pre> \"internal\" - Internal database. Authentication is against Client Usernames. \"ldap\" - LDAP authentication. An LDAP profile name must be provided. \"radius\" - RADIUS authentication. A RADIUS profile name must be provided. \"none\" - No authentication. Anonymous login allowed. </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("internal", "ldap", "radius", "none"),
				},
			},
			"authentication_client_cert_allow_api_provided_username_enabled": schema.BoolAttribute{
				Description: "Enable or disable allowing a client to specify a Client Username via the API connect method. When disabled, the certificate CN (Common Name) is always used. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"authentication_client_cert_enabled": schema.BoolAttribute{
				Description: "Enable or disable client certificate authentication in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"authentication_client_cert_max_chain_depth": schema.Int64Attribute{
				Description: "The maximum depth for a client certificate chain. The depth of a chain is defined as the number of signing CA certificates that are present in the chain back to a trusted self-signed root CA certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.",
				Optional:    true,
			},
			"authentication_client_cert_revocation_check_mode": schema.StringAttribute{
				Description: "The desired behavior for client certificate revocation checking. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"allow-valid\"`. The allowed values and their meaning are:  <pre> \"allow-all\" - Allow the client to authenticate, the result of client certificate revocation check is ignored. \"allow-unknown\" - Allow the client to authenticate even if the revocation status of his certificate cannot be determined. \"allow-valid\" - Allow the client to authenticate only when the revocation check returned an explicit positive response. </pre>  Available since 2.6.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("allow-all", "allow-unknown", "allow-valid"),
				},
			},
			"authentication_client_cert_username_source": schema.StringAttribute{
				Description: "The field from the client certificate to use as the client username. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"common-name\"`. The allowed values and their meaning are:  <pre> \"certificate-thumbprint\" - The username is computed as the SHA-1 hash over the entire DER-encoded contents of the client certificate. \"common-name\" - The username is extracted from the certificate's first instance of the Common Name attribute in the Subject DN. \"common-name-last\" - The username is extracted from the certificate's last instance of the Common Name attribute in the Subject DN. \"subject-alternate-name-msupn\" - The username is extracted from the certificate's Other Name type of the Subject Alternative Name and must have the msUPN signature. \"uid\" - The username is extracted from the certificate's first instance of the User Identifier attribute in the Subject DN. \"uid-last\" - The username is extracted from the certificate's last instance of the User Identifier attribute in the Subject DN. </pre>  Available since 2.6.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("certificate-thumbprint", "common-name", "common-name-last", "subject-alternate-name-msupn", "uid", "uid-last"),
				},
			},
			"authentication_client_cert_validate_date_enabled": schema.BoolAttribute{
				Description: "Enable or disable validation of the \"Not Before\" and \"Not After\" validity dates in the client certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"authentication_kerberos_allow_api_provided_username_enabled": schema.BoolAttribute{
				Description: "Enable or disable allowing a client to specify a Client Username via the API connect method. When disabled, the Kerberos Principal name is always used. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"authentication_kerberos_enabled": schema.BoolAttribute{
				Description: "Enable or disable Kerberos authentication in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"authentication_oauth_default_profile_name": schema.StringAttribute{
				Description: "The name of the profile to use when the client does not supply a profile name. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.25.",
				Optional:    true,
			},
			"authentication_oauth_default_provider_name": schema.StringAttribute{
				Description: "The name of the provider to use when the client does not supply a provider name. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Deprecated since 2.25. authenticationOauthDefaultProviderName and authenticationOauthProviders replaced by authenticationOauthDefaultProfileName and authenticationOauthProfiles.",
				Optional:    true,
			},
			"authentication_oauth_enabled": schema.BoolAttribute{
				Description: "Enable or disable OAuth authentication. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.13.",
				Optional:    true,
			},
			"authorization_ldap_group_membership_attribute_name": schema.StringAttribute{
				Description: "The name of the attribute that is retrieved from the LDAP server as part of the LDAP search when authorizing a client connecting to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"memberOf\"`.",
				Optional:    true,
			},
			"authorization_ldap_trim_client_username_domain_enabled": schema.BoolAttribute{
				Description: "Enable or disable client-username domain trimming for LDAP lookups of client connections. When enabled, the value of $CLIENT_USERNAME (when used for searching) will be truncated at the first occurance of the @ character. For example, if the client-username is in the form of an email address, then the domain portion will be removed. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.13.",
				Optional:    true,
			},
			"authorization_profile_name": schema.StringAttribute{
				Description: "The name of the LDAP Profile to use for client authorization. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"authorization_type": schema.StringAttribute{
				Description: "The type of authorization to use for clients connecting to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"internal\"`. The allowed values and their meaning are:  <pre> \"ldap\" - LDAP authorization. \"internal\" - Internal authorization. </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("ldap", "internal"),
				},
			},
			"bridging_tls_server_cert_enforce_trusted_common_name_enabled": schema.BoolAttribute{
				Description: "Enable or disable validation of the Common Name (CN) in the server certificate from the remote broker. If enabled, the Common Name is checked against the list of Trusted Common Names configured for the Bridge. Common Name validation is not performed if Server Certificate Name Validation is enabled, even if Common Name validation is enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Deprecated since 2.18. Common Name validation has been replaced by Server Certificate Name validation.",
				Optional:    true,
			},
			"bridging_tls_server_cert_max_chain_depth": schema.Int64Attribute{
				Description: "The maximum depth for a server certificate chain. The depth of a chain is defined as the number of signing CA certificates that are present in the chain back to a trusted self-signed root CA certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.",
				Optional:    true,
			},
			"bridging_tls_server_cert_validate_date_enabled": schema.BoolAttribute{
				Description: "Enable or disable validation of the \"Not Before\" and \"Not After\" validity dates in the server certificate. When disabled, a certificate will be accepted even if the certificate is not valid based on these dates. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"bridging_tls_server_cert_validate_name_enabled": schema.BoolAttribute{
				Description: "Enable or disable the standard TLS authentication mechanism of verifying the name used to connect to the bridge. If enabled, the name used to connect to the bridge is checked against the names specified in the certificate returned by the remote router. Legacy Common Name validation is not performed if Server Certificate Name Validation is enabled, even if Common Name validation is also enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`. Available since 2.18.",
				Optional:    true,
			},
			"distributed_cache_management_enabled": schema.BoolAttribute{
				Description: "Enable or disable managing of cache instances over the message bus. The default value is `true`.",
				Optional:    true,
			},
			"dmr_enabled": schema.BoolAttribute{
				Description: "Enable or disable Dynamic Message Routing (DMR) for the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.11.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Enable or disable the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"event_connection_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: EventThresholdDatasourceAttributes,
			},
			"event_egress_flow_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: EventThresholdDatasourceAttributes,
			},
			"event_egress_msg_rate_threshold": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: EventThresholdDatasourceAttributes,
			},
			"event_endpoint_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: EventThresholdDatasourceAttributes,
			},
			"event_ingress_flow_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: EventThresholdDatasourceAttributes,
			},
			"event_ingress_msg_rate_threshold": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: EventThresholdByValueDatasourceAttributes,
			},
			"event_large_msg_threshold": schema.Int64Attribute{
				Description: "The threshold, in kilobytes, after which a message is considered to be large for the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1024`.",
				Optional:    true,
			},
			"event_log_tag": schema.StringAttribute{
				Description: "A prefix applied to all published Events in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"event_msg_spool_usage_threshold": schema.ObjectAttribute{
				Description: "",
				Optional:    true,
			},
			"event_publish_client_enabled": schema.BoolAttribute{
				Description: "Enable or disable Client level Event message publishing. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"event_publish_msg_vpn_enabled": schema.BoolAttribute{
				Description: "Enable or disable Message VPN level Event message publishing. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"event_publish_subscription_mode": schema.StringAttribute{
				Description: "Subscription level Event message publishing mode. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"off\"`. The allowed values and their meaning are:  <pre> \"off\" - Disable client level event message publishing. \"on-with-format-v1\" - Enable client level event message publishing with format v1. \"on-with-no-unsubscribe-events-on-disconnect-format-v1\" - As \"on-with-format-v1\", but unsubscribe events are not generated when a client disconnects. Unsubscribe events are still raised when a client explicitly unsubscribes from its subscriptions. \"on-with-format-v2\" - Enable client level event message publishing with format v2. \"on-with-no-unsubscribe-events-on-disconnect-format-v2\" - As \"on-with-format-v2\", but unsubscribe events are not generated when a client disconnects. Unsubscribe events are still raised when a client explicitly unsubscribes from its subscriptions. </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("off", "on-with-format-v1", "on-with-no-unsubscribe-events-on-disconnect-format-v1", "on-with-format-v2", "on-with-no-unsubscribe-events-on-disconnect-format-v2"),
				},
			},
			"event_publish_topic_format_mqtt_enabled": schema.BoolAttribute{
				Description: "Enable or disable Event publish topics in MQTT format. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"event_publish_topic_format_smf_enabled": schema.BoolAttribute{
				Description: "Enable or disable Event publish topics in SMF format. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"event_service_amqp_connection_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: EventThresholdDatasourceAttributes,
			},
			"event_service_mqtt_connection_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: EventThresholdDatasourceAttributes,
			},
			"event_service_rest_incoming_connection_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: EventThresholdDatasourceAttributes,
			},
			"event_service_smf_connection_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: EventThresholdDatasourceAttributes,
			},
			"event_service_web_connection_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: EventThresholdDatasourceAttributes,
			},
			"event_subscription_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: EventThresholdDatasourceAttributes,
			},
			"event_transacted_session_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: EventThresholdDatasourceAttributes,
			},
			"event_transaction_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,

				Attributes: EventThresholdDatasourceAttributes,
			},
			"export_subscriptions_enabled": schema.BoolAttribute{
				Description: "Enable or disable the export of subscriptions in the Message VPN to other routers in the network over Neighbor links. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"jndi_enabled": schema.BoolAttribute{
				Description: "Enable or disable JNDI access for clients in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.2.",
				Optional:    true,
			},
			"max_connection_count": schema.Int64Attribute{
				Description: "The maximum number of client connections to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform.",
				Optional:    true,
			},
			"max_egress_flow_count": schema.Int64Attribute{
				Description: "The maximum number of transmit flows that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.",
				Optional:    true,
			},
			"max_endpoint_count": schema.Int64Attribute{
				Description: "The maximum number of Queues and Topic Endpoints that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.",
				Optional:    true,
			},
			"max_ingress_flow_count": schema.Int64Attribute{
				Description: "The maximum number of receive flows that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.",
				Optional:    true,
			},
			"max_msg_spool_usage": schema.Int64Attribute{
				Description: "The maximum message spool usage by the Message VPN, in megabytes. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.",
				Optional:    true,
			},
			"max_subscription_count": schema.Int64Attribute{
				Description: "The maximum number of local client subscriptions that can be added to the Message VPN. This limit is not enforced when a subscription is added using a management interface, such as CLI or SEMP. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.",
				Optional:    true,
			},
			"max_transacted_session_count": schema.Int64Attribute{
				Description: "The maximum number of transacted sessions that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.",
				Optional:    true,
			},
			"max_transaction_count": schema.Int64Attribute{
				Description: "The maximum number of transactions that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.",
				Optional:    true,
			},
			"mqtt_retain_max_memory": schema.Int64Attribute{
				Description: "The maximum total memory usage of the MQTT Retain feature for this Message VPN, in MB. If the maximum memory is reached, any arriving retain messages that require more memory are discarded. A value of -1 indicates that the memory is bounded only by the global max memory limit. A value of 0 prevents MQTT Retain from becoming operational. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `-1`. Available since 2.11.",
				Optional:    true,
			},
			"msg_vpn_name": schema.StringAttribute{
				Description: "The name of the Message VPN.",
				Required:    true,
			},
			"replication_ack_propagation_interval_msg_count": schema.Int64Attribute{
				Description: "The acknowledgement (ACK) propagation interval for the replication Bridge, in number of replicated messages. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `20`.",
				Optional:    true,
			},
			"replication_bridge_authentication_basic_client_username": schema.StringAttribute{
				Description: "The Client Username the replication Bridge uses to login to the remote Message VPN. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"replication_bridge_authentication_basic_password": schema.StringAttribute{
				Description: "The password for the Client Username. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.",
				Optional:    true,
			},
			"replication_bridge_authentication_client_cert_content": schema.StringAttribute{
				Description: "The PEM formatted content for the client certificate used by this bridge to login to the Remote Message VPN. It must consist of a private key and between one and three certificates comprising the certificate trust chain. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changing this attribute requires an HTTPS connection. The default value is `\"\"`. Available since 2.9.",
				Optional:    true,
			},
			"replication_bridge_authentication_client_cert_password": schema.StringAttribute{
				Description: "The password for the client certificate. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changing this attribute requires an HTTPS connection. The default value is `\"\"`. Available since 2.9.",
				Optional:    true,
			},
			"replication_bridge_authentication_scheme": schema.StringAttribute{
				Description: "The authentication scheme for the replication Bridge in the Message VPN. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"basic\"`. The allowed values and their meaning are:  <pre> \"basic\" - Basic Authentication Scheme (via username and password). \"client-certificate\" - Client Certificate Authentication Scheme (via certificate file or content). </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("basic", "client-certificate"),
				},
			},
			"replication_bridge_compressed_data_enabled": schema.BoolAttribute{
				Description: "Enable or disable use of compression for the replication Bridge. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"replication_bridge_egress_flow_window_size": schema.Int64Attribute{
				Description: "The size of the window used for guaranteed messages published to the replication Bridge, in messages. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `255`.",
				Optional:    true,
			},
			"replication_bridge_retry_delay": schema.Int64Attribute{
				Description: "The number of seconds that must pass before retrying the replication Bridge connection. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `3`.",
				Optional:    true,
			},
			"replication_bridge_tls_enabled": schema.BoolAttribute{
				Description: "Enable or disable use of encryption (TLS) for the replication Bridge connection. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"replication_bridge_unidirectional_client_profile_name": schema.StringAttribute{
				Description: "The Client Profile for the unidirectional replication Bridge in the Message VPN. It is used only for the TCP parameters. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"#client-profile\"`.",
				Optional:    true,
			},
			"replication_enabled": schema.BoolAttribute{
				Description: "Enable or disable replication for the Message VPN. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"replication_enabled_queue_behavior": schema.StringAttribute{
				Description: "The behavior to take when enabling replication for the Message VPN, depending on the existence of the replication Queue. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"fail-on-existing-queue\"`. The allowed values and their meaning are:  <pre> \"fail-on-existing-queue\" - The data replication queue must not already exist. \"force-use-existing-queue\" - The data replication queue must already exist. Any data messages on the Queue will be forwarded to interested applications. IMPORTANT: Before using this mode be certain that the messages are not stale or otherwise unsuitable to be forwarded. This mode can only be specified when the existing queue is configured the same as is currently specified under replication configuration otherwise the enabling of replication will fail. \"force-recreate-queue\" - The data replication queue must already exist. Any data messages on the Queue will be discarded. IMPORTANT: Before using this mode be certain that the messages on the existing data replication queue are not needed by interested applications. </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("fail-on-existing-queue", "force-use-existing-queue", "force-recreate-queue"),
				},
			},
			"replication_queue_max_msg_spool_usage": schema.Int64Attribute{
				Description: "The maximum message spool usage by the replication Bridge local Queue (quota), in megabytes. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `60000`.",
				Optional:    true,
			},
			"replication_queue_reject_msg_to_sender_on_discard_enabled": schema.BoolAttribute{
				Description: "Enable or disable whether messages discarded on the replication Bridge local Queue are rejected back to the sender. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"replication_reject_msg_when_sync_ineligible_enabled": schema.BoolAttribute{
				Description: "Enable or disable whether guaranteed messages published to synchronously replicated Topics are rejected back to the sender when synchronous replication becomes ineligible. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"replication_role": schema.StringAttribute{
				Description: "The replication role for the Message VPN. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"standby\"`. The allowed values and their meaning are:  <pre> \"active\" - Assume the Active role in replication for the Message VPN. \"standby\" - Assume the Standby role in replication for the Message VPN. </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("active", "standby"),
				},
			},
			"replication_transaction_mode": schema.StringAttribute{
				Description: "The transaction replication mode for all transactions within the Message VPN. Changing this value during operation will not affect existing transactions; it is only used upon starting a transaction. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"async\"`. The allowed values and their meaning are:  <pre> \"sync\" - Messages are acknowledged when replicated (spooled remotely). \"async\" - Messages are acknowledged when pending replication (spooled locally). </pre> ",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("sync", "async"),
				},
			},
			"rest_tls_server_cert_enforce_trusted_common_name_enabled": schema.BoolAttribute{
				Description: "Enable or disable validation of the Common Name (CN) in the server certificate from the remote REST Consumer. If enabled, the Common Name is checked against the list of Trusted Common Names configured for the REST Consumer. Common Name validation is not performed if Server Certificate Name Validation is enabled, even if Common Name validation is enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Deprecated since 2.17. Common Name validation has been replaced by Server Certificate Name validation.",
				Optional:    true,
			},
			"rest_tls_server_cert_max_chain_depth": schema.Int64Attribute{
				Description: "The maximum depth for a REST Consumer server certificate chain. The depth of a chain is defined as the number of signing CA certificates that are present in the chain back to a trusted self-signed root CA certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.",
				Optional:    true,
			},
			"rest_tls_server_cert_validate_date_enabled": schema.BoolAttribute{
				Description: "Enable or disable validation of the \"Not Before\" and \"Not After\" validity dates in the REST Consumer server certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"rest_tls_server_cert_validate_name_enabled": schema.BoolAttribute{
				Description: "Enable or disable the standard TLS authentication mechanism of verifying the name used to connect to the remote REST Consumer. If enabled, the name used to connect to the remote REST Consumer is checked against the names specified in the certificate returned by the remote router. Legacy Common Name validation is not performed if Server Certificate Name Validation is enabled, even if Common Name validation is also enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`. Available since 2.17.",
				Optional:    true,
			},
			"semp_over_msg_bus_admin_client_enabled": schema.BoolAttribute{
				Description: "Enable or disable \"admin client\" SEMP over the message bus commands for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"semp_over_msg_bus_admin_distributed_cache_enabled": schema.BoolAttribute{
				Description: "Enable or disable \"admin distributed-cache\" SEMP over the message bus commands for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"semp_over_msg_bus_admin_enabled": schema.BoolAttribute{
				Description: "Enable or disable \"admin\" SEMP over the message bus commands for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"semp_over_msg_bus_enabled": schema.BoolAttribute{
				Description: "Enable or disable SEMP over the message bus for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"semp_over_msg_bus_show_enabled": schema.BoolAttribute{
				Description: "Enable or disable \"show\" SEMP over the message bus commands for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"service_amqp_max_connection_count": schema.Int64Attribute{
				Description: "The maximum number of AMQP client connections that can be simultaneously connected to the Message VPN. This value may be higher than supported by the platform. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform. Available since 2.7.",
				Optional:    true,
			},
			"service_amqp_plain_text_enabled": schema.BoolAttribute{
				Description: "Enable or disable the plain-text AMQP service in the Message VPN. Disabling causes clients connected to the corresponding listen-port to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.7.",
				Optional:    true,
			},
			"service_amqp_plain_text_listen_port": schema.Int64Attribute{
				Description: "The port number for plain-text AMQP clients that connect to the Message VPN. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.7.",
				Optional:    true,
			},
			"service_amqp_tls_enabled": schema.BoolAttribute{
				Description: "Enable or disable the use of encryption (TLS) for the AMQP service in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.7.",
				Optional:    true,
			},
			"service_amqp_tls_listen_port": schema.Int64Attribute{
				Description: "The port number for AMQP clients that connect to the Message VPN over TLS. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.7.",
				Optional:    true,
			},
			"service_mqtt_authentication_client_cert_request": schema.StringAttribute{
				Description: "Determines when to request a client certificate from an incoming MQTT client connecting via a TLS port. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"when-enabled-in-message-vpn\"`. The allowed values and their meaning are:  <pre> \"always\" - Always ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"never\" - Never ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"when-enabled-in-message-vpn\" - Only ask for a client-certificate if client certificate authentication is enabled under \"message-vpn >  authentication > client-certificate > shutdown\". </pre>  Available since 2.21.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("always", "never", "when-enabled-in-message-vpn"),
				},
			},
			"service_mqtt_max_connection_count": schema.Int64Attribute{
				Description: "The maximum number of MQTT client connections that can be simultaneously connected to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform. Available since 2.1.",
				Optional:    true,
			},
			"service_mqtt_plain_text_enabled": schema.BoolAttribute{
				Description: "Enable or disable the plain-text MQTT service in the Message VPN. Disabling causes clients currently connected to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.1.",
				Optional:    true,
			},
			"service_mqtt_plain_text_listen_port": schema.Int64Attribute{
				Description: "The port number for plain-text MQTT clients that connect to the Message VPN. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.1.",
				Optional:    true,
			},
			"service_mqtt_tls_enabled": schema.BoolAttribute{
				Description: "Enable or disable the use of encryption (TLS) for the MQTT service in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.1.",
				Optional:    true,
			},
			"service_mqtt_tls_listen_port": schema.Int64Attribute{
				Description: "The port number for MQTT clients that connect to the Message VPN over TLS. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.1.",
				Optional:    true,
			},
			"service_mqtt_tls_web_socket_enabled": schema.BoolAttribute{
				Description: "Enable or disable the use of encrypted WebSocket (WebSocket over TLS) for the MQTT service in the Message VPN. Disabling causes clients currently connected by encrypted WebSocket to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.1.",
				Optional:    true,
			},
			"service_mqtt_tls_web_socket_listen_port": schema.Int64Attribute{
				Description: "The port number for MQTT clients that connect to the Message VPN using WebSocket over TLS. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.1.",
				Optional:    true,
			},
			"service_mqtt_web_socket_enabled": schema.BoolAttribute{
				Description: "Enable or disable the use of WebSocket for the MQTT service in the Message VPN. Disabling causes clients currently connected by WebSocket to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.1.",
				Optional:    true,
			},
			"service_mqtt_web_socket_listen_port": schema.Int64Attribute{
				Description: "The port number for plain-text MQTT clients that connect to the Message VPN using WebSocket. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.1.",
				Optional:    true,
			},
			"service_rest_incoming_authentication_client_cert_request": schema.StringAttribute{
				Description: "Determines when to request a client certificate from an incoming REST Producer connecting via a TLS port. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"when-enabled-in-message-vpn\"`. The allowed values and their meaning are:  <pre> \"always\" - Always ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"never\" - Never ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"when-enabled-in-message-vpn\" - Only ask for a client-certificate if client certificate authentication is enabled under \"message-vpn >  authentication > client-certificate > shutdown\". </pre>  Available since 2.21.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("always", "never", "when-enabled-in-message-vpn"),
				},
			},
			"service_rest_incoming_authorization_header_handling": schema.StringAttribute{
				Description: "The handling of Authorization headers for incoming REST connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"drop\"`. The allowed values and their meaning are:  <pre> \"drop\" - Do not attach the Authorization header to the message as a user property. This configuration is most secure. \"forward\" - Forward the Authorization header, attaching it to the message as a user property in the same way as other headers. For best security, use the drop setting. \"legacy\" - If the Authorization header was used for authentication to the broker, do not attach it to the message. If the Authorization header was not used for authentication to the broker, attach it to the message as a user property in the same way as other headers. For best security, use the drop setting. </pre>  Available since 2.19.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("drop", "forward", "legacy"),
				},
			},
			"service_rest_incoming_max_connection_count": schema.Int64Attribute{
				Description: "The maximum number of REST incoming client connections that can be simultaneously connected to the Message VPN. This value may be higher than supported by the platform. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform.",
				Optional:    true,
			},
			"service_rest_incoming_plain_text_enabled": schema.BoolAttribute{
				Description: "Enable or disable the plain-text REST service for incoming clients in the Message VPN. Disabling causes clients currently connected to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"service_rest_incoming_plain_text_listen_port": schema.Int64Attribute{
				Description: "The port number for incoming plain-text REST clients that connect to the Message VPN. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.",
				Optional:    true,
			},
			"service_rest_incoming_tls_enabled": schema.BoolAttribute{
				Description: "Enable or disable the use of encryption (TLS) for the REST service for incoming clients in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
			"service_rest_incoming_tls_listen_port": schema.Int64Attribute{
				Description: "The port number for incoming REST clients that connect to the Message VPN over TLS. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.",
				Optional:    true,
			},
			"service_rest_mode": schema.StringAttribute{
				Description: "The REST service mode for incoming REST clients that connect to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"messaging\"`. The allowed values and their meaning are:  <pre> \"gateway\" - Act as a message gateway through which REST messages are propagated. \"messaging\" - Act as a message broker on which REST messages are queued. </pre>  Available since 2.6.",
				Optional:    true,

				Validators: []validator.String{
					stringvalidator.OneOf("gateway", "messaging"),
				},
			},
			"service_rest_outgoing_max_connection_count": schema.Int64Attribute{
				Description: "The maximum number of REST Consumer (outgoing) client connections that can be simultaneously connected to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.",
				Optional:    true,
			},
			"service_smf_max_connection_count": schema.Int64Attribute{
				Description: "The maximum number of SMF client connections that can be simultaneously connected to the Message VPN. This value may be higher than supported by the platform. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.",
				Optional:    true,
			},
			"service_smf_plain_text_enabled": schema.BoolAttribute{
				Description: "Enable or disable the plain-text SMF service in the Message VPN. Disabling causes clients currently connected to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"service_smf_tls_enabled": schema.BoolAttribute{
				Description: "Enable or disable the use of encryption (TLS) for the SMF service in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"service_web_authentication_client_cert_request": schema.StringAttribute{
				Description: "Determines when to request a client certificate from a Web Transport client connecting via a TLS port. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"when-enabled-in-message-vpn\"`. The allowed values and their meaning are:  <pre> \"always\" - Always ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"never\" - Never ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"when-enabled-in-message-vpn\" - Only ask for a client-certificate if client certificate authentication is enabled under \"message-vpn >  authentication > client-certificate > shutdown\". </pre>  Available since 2.21.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("always", "never", "when-enabled-in-message-vpn"),
				},
			},
			"service_web_max_connection_count": schema.Int64Attribute{
				Description: "The maximum number of Web Transport client connections that can be simultaneously connected to the Message VPN. This value may be higher than supported by the platform. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform.",
				Optional:    true,
			},
			"service_web_plain_text_enabled": schema.BoolAttribute{
				Description: "Enable or disable the plain-text Web Transport service in the Message VPN. Disabling causes clients currently connected to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"service_web_tls_enabled": schema.BoolAttribute{
				Description: "Enable or disable the use of TLS for the Web Transport service in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Optional:    true,
			},
			"tls_allow_downgrade_to_plain_text_enabled": schema.BoolAttribute{
				Description: "Enable or disable the allowing of TLS SMF clients to downgrade their connections to plain-text connections. Changing this will not affect existing connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Optional:    true,
			},
		},
	}
}

func (r *msgVpnDataSource) SetProvider(provider *solaceProvider) {
	r.solaceProvider = provider
}

func (r msgVpnDataSource) NewData() *MsgVpn {
	return &MsgVpn{}
}

func (r msgVpnDataSource) Read(data *MsgVpn, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.MsgVpnApi.GetMsgVpn(r.Context, *data.MsgVpnName)
	apiResponse, httpResponse, err := apiReq.Execute()
	if err == nil && apiResponse != nil && apiResponse.Data != nil {
		data.ToTF(apiResponse.Data)
	}
	return httpResponse, err
}
