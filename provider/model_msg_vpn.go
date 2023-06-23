package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// MsgVpn struct for MsgVpn
type MsgVpn struct {
	Alias                                                   *string                `tfsdk:"alias"`
	AuthenticationBasicEnabled                              *bool                  `tfsdk:"authentication_basic_enabled"`
	AuthenticationBasicProfileName                          *string                `tfsdk:"authentication_basic_profile_name"`
	AuthenticationBasicRadiusDomain                         *string                `tfsdk:"authentication_basic_radius_domain"`
	AuthenticationBasicType                                 *string                `tfsdk:"authentication_basic_type"`
	AuthenticationClientCertAllowApiProvidedUsernameEnabled *bool                  `tfsdk:"authentication_client_cert_allow_api_provided_username_enabled"`
	AuthenticationClientCertEnabled                         *bool                  `tfsdk:"authentication_client_cert_enabled"`
	AuthenticationClientCertMaxChainDepth                   *int64                 `tfsdk:"authentication_client_cert_max_chain_depth"`
	AuthenticationClientCertRevocationCheckMode             *string                `tfsdk:"authentication_client_cert_revocation_check_mode"`
	AuthenticationClientCertUsernameSource                  *string                `tfsdk:"authentication_client_cert_username_source"`
	AuthenticationClientCertValidateDateEnabled             *bool                  `tfsdk:"authentication_client_cert_validate_date_enabled"`
	AuthenticationKerberosAllowApiProvidedUsernameEnabled   *bool                  `tfsdk:"authentication_kerberos_allow_api_provided_username_enabled"`
	AuthenticationKerberosEnabled                           *bool                  `tfsdk:"authentication_kerberos_enabled"`
	AuthenticationOauthDefaultProfileName                   *string                `tfsdk:"authentication_oauth_default_profile_name"`
	AuthenticationOauthDefaultProviderName                  *string                `tfsdk:"authentication_oauth_default_provider_name"`
	AuthenticationOauthEnabled                              *bool                  `tfsdk:"authentication_oauth_enabled"`
	AuthorizationLdapGroupMembershipAttributeName           *string                `tfsdk:"authorization_ldap_group_membership_attribute_name"`
	AuthorizationLdapTrimClientUsernameDomainEnabled        *bool                  `tfsdk:"authorization_ldap_trim_client_username_domain_enabled"`
	AuthorizationProfileName                                *string                `tfsdk:"authorization_profile_name"`
	AuthorizationType                                       *string                `tfsdk:"authorization_type"`
	BridgingTlsServerCertEnforceTrustedCommonNameEnabled    *bool                  `tfsdk:"bridging_tls_server_cert_enforce_trusted_common_name_enabled"`
	BridgingTlsServerCertMaxChainDepth                      *int64                 `tfsdk:"bridging_tls_server_cert_max_chain_depth"`
	BridgingTlsServerCertValidateDateEnabled                *bool                  `tfsdk:"bridging_tls_server_cert_validate_date_enabled"`
	BridgingTlsServerCertValidateNameEnabled                *bool                  `tfsdk:"bridging_tls_server_cert_validate_name_enabled"`
	DistributedCacheManagementEnabled                       *bool                  `tfsdk:"distributed_cache_management_enabled"`
	DmrEnabled                                              *bool                  `tfsdk:"dmr_enabled"`
	Enabled                                                 *bool                  `tfsdk:"enabled"`
	EventConnectionCountThreshold                           *EventThreshold        `tfsdk:"event_connection_count_threshold"`
	EventEgressFlowCountThreshold                           *EventThreshold        `tfsdk:"event_egress_flow_count_threshold"`
	EventEgressMsgRateThreshold                             *EventThresholdByValue `tfsdk:"event_egress_msg_rate_threshold"`
	EventEndpointCountThreshold                             *EventThreshold        `tfsdk:"event_endpoint_count_threshold"`
	EventIngressFlowCountThreshold                          *EventThreshold        `tfsdk:"event_ingress_flow_count_threshold"`
	EventIngressMsgRateThreshold                            *EventThresholdByValue `tfsdk:"event_ingress_msg_rate_threshold"`
	EventLargeMsgThreshold                                  *int64                 `tfsdk:"event_large_msg_threshold"`
	EventLogTag                                             *string                `tfsdk:"event_log_tag"`
	EventMsgSpoolUsageThreshold                             *EventThreshold        `tfsdk:"event_msg_spool_usage_threshold"`
	EventPublishClientEnabled                               *bool                  `tfsdk:"event_publish_client_enabled"`
	EventPublishMsgVpnEnabled                               *bool                  `tfsdk:"event_publish_msg_vpn_enabled"`
	EventPublishSubscriptionMode                            *string                `tfsdk:"event_publish_subscription_mode"`
	EventPublishTopicFormatMqttEnabled                      *bool                  `tfsdk:"event_publish_topic_format_mqtt_enabled"`
	EventPublishTopicFormatSmfEnabled                       *bool                  `tfsdk:"event_publish_topic_format_smf_enabled"`
	EventServiceAmqpConnectionCountThreshold                *EventThreshold        `tfsdk:"event_service_amqp_connection_count_threshold"`
	EventServiceMqttConnectionCountThreshold                *EventThreshold        `tfsdk:"event_service_mqtt_connection_count_threshold"`
	EventServiceRestIncomingConnectionCountThreshold        *EventThreshold        `tfsdk:"event_service_rest_incoming_connection_count_threshold"`
	EventServiceSmfConnectionCountThreshold                 *EventThreshold        `tfsdk:"event_service_smf_connection_count_threshold"`
	EventServiceWebConnectionCountThreshold                 *EventThreshold        `tfsdk:"event_service_web_connection_count_threshold"`
	EventSubscriptionCountThreshold                         *EventThreshold        `tfsdk:"event_subscription_count_threshold"`
	EventTransactedSessionCountThreshold                    *EventThreshold        `tfsdk:"event_transacted_session_count_threshold"`
	EventTransactionCountThreshold                          *EventThreshold        `tfsdk:"event_transaction_count_threshold"`
	ExportSubscriptionsEnabled                              *bool                  `tfsdk:"export_subscriptions_enabled"`
	JndiEnabled                                             *bool                  `tfsdk:"jndi_enabled"`
	MaxConnectionCount                                      *int64                 `tfsdk:"max_connection_count"`
	MaxEgressFlowCount                                      *int64                 `tfsdk:"max_egress_flow_count"`
	MaxEndpointCount                                        *int64                 `tfsdk:"max_endpoint_count"`
	MaxIngressFlowCount                                     *int64                 `tfsdk:"max_ingress_flow_count"`
	MaxMsgSpoolUsage                                        *int64                 `tfsdk:"max_msg_spool_usage"`
	MaxSubscriptionCount                                    *int64                 `tfsdk:"max_subscription_count"`
	MaxTransactedSessionCount                               *int64                 `tfsdk:"max_transacted_session_count"`
	MaxTransactionCount                                     *int64                 `tfsdk:"max_transaction_count"`
	MqttRetainMaxMemory                                     *int32                 `tfsdk:"mqtt_retain_max_memory"`
	MsgVpnName                                              *string                `tfsdk:"msg_vpn_name"`
	ReplicationAckPropagationIntervalMsgCount               *int64                 `tfsdk:"replication_ack_propagation_interval_msg_count"`
	ReplicationBridgeAuthenticationBasicClientUsername      *string                `tfsdk:"replication_bridge_authentication_basic_client_username"`
	ReplicationBridgeAuthenticationBasicPassword            *string                `tfsdk:"replication_bridge_authentication_basic_password"`
	ReplicationBridgeAuthenticationClientCertContent        *string                `tfsdk:"replication_bridge_authentication_client_cert_content"`
	ReplicationBridgeAuthenticationClientCertPassword       *string                `tfsdk:"replication_bridge_authentication_client_cert_password"`
	ReplicationBridgeAuthenticationScheme                   *string                `tfsdk:"replication_bridge_authentication_scheme"`
	ReplicationBridgeCompressedDataEnabled                  *bool                  `tfsdk:"replication_bridge_compressed_data_enabled"`
	ReplicationBridgeEgressFlowWindowSize                   *int64                 `tfsdk:"replication_bridge_egress_flow_window_size"`
	ReplicationBridgeRetryDelay                             *int64                 `tfsdk:"replication_bridge_retry_delay"`
	ReplicationBridgeTlsEnabled                             *bool                  `tfsdk:"replication_bridge_tls_enabled"`
	ReplicationBridgeUnidirectionalClientProfileName        *string                `tfsdk:"replication_bridge_unidirectional_client_profile_name"`
	ReplicationEnabled                                      *bool                  `tfsdk:"replication_enabled"`
	ReplicationEnabledQueueBehavior                         *string                `tfsdk:"replication_enabled_queue_behavior"`
	ReplicationQueueMaxMsgSpoolUsage                        *int64                 `tfsdk:"replication_queue_max_msg_spool_usage"`
	ReplicationQueueRejectMsgToSenderOnDiscardEnabled       *bool                  `tfsdk:"replication_queue_reject_msg_to_sender_on_discard_enabled"`
	ReplicationRejectMsgWhenSyncIneligibleEnabled           *bool                  `tfsdk:"replication_reject_msg_when_sync_ineligible_enabled"`
	ReplicationRole                                         *string                `tfsdk:"replication_role"`
	ReplicationTransactionMode                              *string                `tfsdk:"replication_transaction_mode"`
	RestTlsServerCertEnforceTrustedCommonNameEnabled        *bool                  `tfsdk:"rest_tls_server_cert_enforce_trusted_common_name_enabled"`
	RestTlsServerCertMaxChainDepth                          *int64                 `tfsdk:"rest_tls_server_cert_max_chain_depth"`
	RestTlsServerCertValidateDateEnabled                    *bool                  `tfsdk:"rest_tls_server_cert_validate_date_enabled"`
	RestTlsServerCertValidateNameEnabled                    *bool                  `tfsdk:"rest_tls_server_cert_validate_name_enabled"`
	SempOverMsgBusAdminClientEnabled                        *bool                  `tfsdk:"semp_over_msg_bus_admin_client_enabled"`
	SempOverMsgBusAdminDistributedCacheEnabled              *bool                  `tfsdk:"semp_over_msg_bus_admin_distributed_cache_enabled"`
	SempOverMsgBusAdminEnabled                              *bool                  `tfsdk:"semp_over_msg_bus_admin_enabled"`
	SempOverMsgBusEnabled                                   *bool                  `tfsdk:"semp_over_msg_bus_enabled"`
	SempOverMsgBusShowEnabled                               *bool                  `tfsdk:"semp_over_msg_bus_show_enabled"`
	ServiceAmqpMaxConnectionCount                           *int64                 `tfsdk:"service_amqp_max_connection_count"`
	ServiceAmqpPlainTextEnabled                             *bool                  `tfsdk:"service_amqp_plain_text_enabled"`
	ServiceAmqpPlainTextListenPort                          *int64                 `tfsdk:"service_amqp_plain_text_listen_port"`
	ServiceAmqpTlsEnabled                                   *bool                  `tfsdk:"service_amqp_tls_enabled"`
	ServiceAmqpTlsListenPort                                *int64                 `tfsdk:"service_amqp_tls_listen_port"`
	ServiceMqttAuthenticationClientCertRequest              *string                `tfsdk:"service_mqtt_authentication_client_cert_request"`
	ServiceMqttMaxConnectionCount                           *int64                 `tfsdk:"service_mqtt_max_connection_count"`
	ServiceMqttPlainTextEnabled                             *bool                  `tfsdk:"service_mqtt_plain_text_enabled"`
	ServiceMqttPlainTextListenPort                          *int64                 `tfsdk:"service_mqtt_plain_text_listen_port"`
	ServiceMqttTlsEnabled                                   *bool                  `tfsdk:"service_mqtt_tls_enabled"`
	ServiceMqttTlsListenPort                                *int64                 `tfsdk:"service_mqtt_tls_listen_port"`
	ServiceMqttTlsWebSocketEnabled                          *bool                  `tfsdk:"service_mqtt_tls_web_socket_enabled"`
	ServiceMqttTlsWebSocketListenPort                       *int64                 `tfsdk:"service_mqtt_tls_web_socket_listen_port"`
	ServiceMqttWebSocketEnabled                             *bool                  `tfsdk:"service_mqtt_web_socket_enabled"`
	ServiceMqttWebSocketListenPort                          *int64                 `tfsdk:"service_mqtt_web_socket_listen_port"`
	ServiceRestIncomingAuthenticationClientCertRequest      *string                `tfsdk:"service_rest_incoming_authentication_client_cert_request"`
	ServiceRestIncomingAuthorizationHeaderHandling          *string                `tfsdk:"service_rest_incoming_authorization_header_handling"`
	ServiceRestIncomingMaxConnectionCount                   *int64                 `tfsdk:"service_rest_incoming_max_connection_count"`
	ServiceRestIncomingPlainTextEnabled                     *bool                  `tfsdk:"service_rest_incoming_plain_text_enabled"`
	ServiceRestIncomingPlainTextListenPort                  *int64                 `tfsdk:"service_rest_incoming_plain_text_listen_port"`
	ServiceRestIncomingTlsEnabled                           *bool                  `tfsdk:"service_rest_incoming_tls_enabled"`
	ServiceRestIncomingTlsListenPort                        *int64                 `tfsdk:"service_rest_incoming_tls_listen_port"`
	ServiceRestMode                                         *string                `tfsdk:"service_rest_mode"`
	ServiceRestOutgoingMaxConnectionCount                   *int64                 `tfsdk:"service_rest_outgoing_max_connection_count"`
	ServiceSmfMaxConnectionCount                            *int64                 `tfsdk:"service_smf_max_connection_count"`
	ServiceSmfPlainTextEnabled                              *bool                  `tfsdk:"service_smf_plain_text_enabled"`
	ServiceSmfTlsEnabled                                    *bool                  `tfsdk:"service_smf_tls_enabled"`
	ServiceWebAuthenticationClientCertRequest               *string                `tfsdk:"service_web_authentication_client_cert_request"`
	ServiceWebMaxConnectionCount                            *int64                 `tfsdk:"service_web_max_connection_count"`
	ServiceWebPlainTextEnabled                              *bool                  `tfsdk:"service_web_plain_text_enabled"`
	ServiceWebTlsEnabled                                    *bool                  `tfsdk:"service_web_tls_enabled"`
	TlsAllowDowngradeToPlainTextEnabled                     *bool                  `tfsdk:"tls_allow_downgrade_to_plain_text_enabled"`
}

func (tfData *MsgVpn) ToTF(apiData *sempv2.MsgVpn) {
	AssignIfDstNotNil(&tfData.Alias, apiData.Alias)
	AssignIfDstNotNil(&tfData.AuthenticationBasicEnabled, apiData.AuthenticationBasicEnabled)
	AssignIfDstNotNil(&tfData.AuthenticationBasicProfileName, apiData.AuthenticationBasicProfileName)
	AssignIfDstNotNil(&tfData.AuthenticationBasicRadiusDomain, apiData.AuthenticationBasicRadiusDomain)
	AssignIfDstNotNil(&tfData.AuthenticationBasicType, apiData.AuthenticationBasicType)
	AssignIfDstNotNil(&tfData.AuthenticationClientCertAllowApiProvidedUsernameEnabled, apiData.AuthenticationClientCertAllowApiProvidedUsernameEnabled)
	AssignIfDstNotNil(&tfData.AuthenticationClientCertEnabled, apiData.AuthenticationClientCertEnabled)
	AssignIfDstNotNil(&tfData.AuthenticationClientCertMaxChainDepth, apiData.AuthenticationClientCertMaxChainDepth)
	AssignIfDstNotNil(&tfData.AuthenticationClientCertRevocationCheckMode, apiData.AuthenticationClientCertRevocationCheckMode)
	AssignIfDstNotNil(&tfData.AuthenticationClientCertUsernameSource, apiData.AuthenticationClientCertUsernameSource)
	AssignIfDstNotNil(&tfData.AuthenticationClientCertValidateDateEnabled, apiData.AuthenticationClientCertValidateDateEnabled)
	AssignIfDstNotNil(&tfData.AuthenticationKerberosAllowApiProvidedUsernameEnabled, apiData.AuthenticationKerberosAllowApiProvidedUsernameEnabled)
	AssignIfDstNotNil(&tfData.AuthenticationKerberosEnabled, apiData.AuthenticationKerberosEnabled)
	AssignIfDstNotNil(&tfData.AuthenticationOauthDefaultProfileName, apiData.AuthenticationOauthDefaultProfileName)
	AssignIfDstNotNil(&tfData.AuthenticationOauthDefaultProviderName, apiData.AuthenticationOauthDefaultProviderName)
	AssignIfDstNotNil(&tfData.AuthenticationOauthEnabled, apiData.AuthenticationOauthEnabled)
	AssignIfDstNotNil(&tfData.AuthorizationLdapGroupMembershipAttributeName, apiData.AuthorizationLdapGroupMembershipAttributeName)
	AssignIfDstNotNil(&tfData.AuthorizationLdapTrimClientUsernameDomainEnabled, apiData.AuthorizationLdapTrimClientUsernameDomainEnabled)
	AssignIfDstNotNil(&tfData.AuthorizationProfileName, apiData.AuthorizationProfileName)
	AssignIfDstNotNil(&tfData.AuthorizationType, apiData.AuthorizationType)
	AssignIfDstNotNil(&tfData.BridgingTlsServerCertEnforceTrustedCommonNameEnabled, apiData.BridgingTlsServerCertEnforceTrustedCommonNameEnabled)
	AssignIfDstNotNil(&tfData.BridgingTlsServerCertMaxChainDepth, apiData.BridgingTlsServerCertMaxChainDepth)
	AssignIfDstNotNil(&tfData.BridgingTlsServerCertValidateDateEnabled, apiData.BridgingTlsServerCertValidateDateEnabled)
	AssignIfDstNotNil(&tfData.BridgingTlsServerCertValidateNameEnabled, apiData.BridgingTlsServerCertValidateNameEnabled)
	AssignIfDstNotNil(&tfData.DistributedCacheManagementEnabled, apiData.DistributedCacheManagementEnabled)
	AssignIfDstNotNil(&tfData.DmrEnabled, apiData.DmrEnabled)
	AssignIfDstNotNil(&tfData.Enabled, apiData.Enabled)
	AssignIfDstNotNil(&tfData.EventConnectionCountThreshold, EventThresholdToTF(apiData.EventConnectionCountThreshold))
	AssignIfDstNotNil(&tfData.EventEgressFlowCountThreshold, EventThresholdToTF(apiData.EventEgressFlowCountThreshold))
	AssignIfDstNotNil(&tfData.EventEgressMsgRateThreshold, EventThresholdByValueToTF(apiData.EventEgressMsgRateThreshold))
	AssignIfDstNotNil(&tfData.EventEndpointCountThreshold, EventThresholdToTF(apiData.EventEndpointCountThreshold))
	AssignIfDstNotNil(&tfData.EventIngressFlowCountThreshold, EventThresholdToTF(apiData.EventIngressFlowCountThreshold))
	AssignIfDstNotNil(&tfData.EventIngressMsgRateThreshold, EventThresholdByValueToTF(apiData.EventIngressMsgRateThreshold))
	AssignIfDstNotNil(&tfData.EventLargeMsgThreshold, apiData.EventLargeMsgThreshold)
	AssignIfDstNotNil(&tfData.EventLogTag, apiData.EventLogTag)
	AssignIfDstNotNil(&tfData.EventMsgSpoolUsageThreshold, EventThresholdToTF(apiData.EventMsgSpoolUsageThreshold))
	AssignIfDstNotNil(&tfData.EventPublishClientEnabled, apiData.EventPublishClientEnabled)
	AssignIfDstNotNil(&tfData.EventPublishMsgVpnEnabled, apiData.EventPublishMsgVpnEnabled)
	AssignIfDstNotNil(&tfData.EventPublishSubscriptionMode, apiData.EventPublishSubscriptionMode)
	AssignIfDstNotNil(&tfData.EventPublishTopicFormatMqttEnabled, apiData.EventPublishTopicFormatMqttEnabled)
	AssignIfDstNotNil(&tfData.EventPublishTopicFormatSmfEnabled, apiData.EventPublishTopicFormatSmfEnabled)
	AssignIfDstNotNil(&tfData.EventServiceAmqpConnectionCountThreshold, EventThresholdToTF(apiData.EventServiceAmqpConnectionCountThreshold))
	AssignIfDstNotNil(&tfData.EventServiceMqttConnectionCountThreshold, EventThresholdToTF(apiData.EventServiceMqttConnectionCountThreshold))
	AssignIfDstNotNil(&tfData.EventServiceRestIncomingConnectionCountThreshold, EventThresholdToTF(apiData.EventServiceRestIncomingConnectionCountThreshold))
	AssignIfDstNotNil(&tfData.EventServiceSmfConnectionCountThreshold, EventThresholdToTF(apiData.EventServiceSmfConnectionCountThreshold))
	AssignIfDstNotNil(&tfData.EventServiceWebConnectionCountThreshold, EventThresholdToTF(apiData.EventServiceWebConnectionCountThreshold))
	AssignIfDstNotNil(&tfData.EventSubscriptionCountThreshold, EventThresholdToTF(apiData.EventSubscriptionCountThreshold))
	AssignIfDstNotNil(&tfData.EventTransactedSessionCountThreshold, EventThresholdToTF(apiData.EventTransactedSessionCountThreshold))
	AssignIfDstNotNil(&tfData.EventTransactionCountThreshold, EventThresholdToTF(apiData.EventTransactionCountThreshold))
	AssignIfDstNotNil(&tfData.ExportSubscriptionsEnabled, apiData.ExportSubscriptionsEnabled)
	AssignIfDstNotNil(&tfData.JndiEnabled, apiData.JndiEnabled)
	AssignIfDstNotNil(&tfData.MaxConnectionCount, apiData.MaxConnectionCount)
	AssignIfDstNotNil(&tfData.MaxEgressFlowCount, apiData.MaxEgressFlowCount)
	AssignIfDstNotNil(&tfData.MaxEndpointCount, apiData.MaxEndpointCount)
	AssignIfDstNotNil(&tfData.MaxIngressFlowCount, apiData.MaxIngressFlowCount)
	AssignIfDstNotNil(&tfData.MaxMsgSpoolUsage, apiData.MaxMsgSpoolUsage)
	AssignIfDstNotNil(&tfData.MaxSubscriptionCount, apiData.MaxSubscriptionCount)
	AssignIfDstNotNil(&tfData.MaxTransactedSessionCount, apiData.MaxTransactedSessionCount)
	AssignIfDstNotNil(&tfData.MaxTransactionCount, apiData.MaxTransactionCount)
	AssignIfDstNotNil(&tfData.MqttRetainMaxMemory, apiData.MqttRetainMaxMemory)
	AssignIfDstNotNil(&tfData.MsgVpnName, apiData.MsgVpnName)
	AssignIfDstNotNil(&tfData.ReplicationAckPropagationIntervalMsgCount, apiData.ReplicationAckPropagationIntervalMsgCount)
	AssignIfDstNotNil(&tfData.ReplicationBridgeAuthenticationBasicClientUsername, apiData.ReplicationBridgeAuthenticationBasicClientUsername)
	AssignIfDstNotNil(&tfData.ReplicationBridgeAuthenticationBasicPassword, apiData.ReplicationBridgeAuthenticationBasicPassword)
	AssignIfDstNotNil(&tfData.ReplicationBridgeAuthenticationClientCertContent, apiData.ReplicationBridgeAuthenticationClientCertContent)
	AssignIfDstNotNil(&tfData.ReplicationBridgeAuthenticationClientCertPassword, apiData.ReplicationBridgeAuthenticationClientCertPassword)
	AssignIfDstNotNil(&tfData.ReplicationBridgeAuthenticationScheme, apiData.ReplicationBridgeAuthenticationScheme)
	AssignIfDstNotNil(&tfData.ReplicationBridgeCompressedDataEnabled, apiData.ReplicationBridgeCompressedDataEnabled)
	AssignIfDstNotNil(&tfData.ReplicationBridgeEgressFlowWindowSize, apiData.ReplicationBridgeEgressFlowWindowSize)
	AssignIfDstNotNil(&tfData.ReplicationBridgeRetryDelay, apiData.ReplicationBridgeRetryDelay)
	AssignIfDstNotNil(&tfData.ReplicationBridgeTlsEnabled, apiData.ReplicationBridgeTlsEnabled)
	AssignIfDstNotNil(&tfData.ReplicationBridgeUnidirectionalClientProfileName, apiData.ReplicationBridgeUnidirectionalClientProfileName)
	AssignIfDstNotNil(&tfData.ReplicationEnabled, apiData.ReplicationEnabled)
	AssignIfDstNotNil(&tfData.ReplicationEnabledQueueBehavior, apiData.ReplicationEnabledQueueBehavior)
	AssignIfDstNotNil(&tfData.ReplicationQueueMaxMsgSpoolUsage, apiData.ReplicationQueueMaxMsgSpoolUsage)
	AssignIfDstNotNil(&tfData.ReplicationQueueRejectMsgToSenderOnDiscardEnabled, apiData.ReplicationQueueRejectMsgToSenderOnDiscardEnabled)
	AssignIfDstNotNil(&tfData.ReplicationRejectMsgWhenSyncIneligibleEnabled, apiData.ReplicationRejectMsgWhenSyncIneligibleEnabled)
	AssignIfDstNotNil(&tfData.ReplicationRole, apiData.ReplicationRole)
	AssignIfDstNotNil(&tfData.ReplicationTransactionMode, apiData.ReplicationTransactionMode)
	AssignIfDstNotNil(&tfData.RestTlsServerCertEnforceTrustedCommonNameEnabled, apiData.RestTlsServerCertEnforceTrustedCommonNameEnabled)
	AssignIfDstNotNil(&tfData.RestTlsServerCertMaxChainDepth, apiData.RestTlsServerCertMaxChainDepth)
	AssignIfDstNotNil(&tfData.RestTlsServerCertValidateDateEnabled, apiData.RestTlsServerCertValidateDateEnabled)
	AssignIfDstNotNil(&tfData.RestTlsServerCertValidateNameEnabled, apiData.RestTlsServerCertValidateNameEnabled)
	AssignIfDstNotNil(&tfData.SempOverMsgBusAdminClientEnabled, apiData.SempOverMsgBusAdminClientEnabled)
	AssignIfDstNotNil(&tfData.SempOverMsgBusAdminDistributedCacheEnabled, apiData.SempOverMsgBusAdminDistributedCacheEnabled)
	AssignIfDstNotNil(&tfData.SempOverMsgBusAdminEnabled, apiData.SempOverMsgBusAdminEnabled)
	AssignIfDstNotNil(&tfData.SempOverMsgBusEnabled, apiData.SempOverMsgBusEnabled)
	AssignIfDstNotNil(&tfData.SempOverMsgBusShowEnabled, apiData.SempOverMsgBusShowEnabled)
	AssignIfDstNotNil(&tfData.ServiceAmqpMaxConnectionCount, apiData.ServiceAmqpMaxConnectionCount)
	AssignIfDstNotNil(&tfData.ServiceAmqpPlainTextEnabled, apiData.ServiceAmqpPlainTextEnabled)
	AssignIfDstNotNil(&tfData.ServiceAmqpPlainTextListenPort, apiData.ServiceAmqpPlainTextListenPort)
	AssignIfDstNotNil(&tfData.ServiceAmqpTlsEnabled, apiData.ServiceAmqpTlsEnabled)
	AssignIfDstNotNil(&tfData.ServiceAmqpTlsListenPort, apiData.ServiceAmqpTlsListenPort)
	AssignIfDstNotNil(&tfData.ServiceMqttAuthenticationClientCertRequest, apiData.ServiceMqttAuthenticationClientCertRequest)
	AssignIfDstNotNil(&tfData.ServiceMqttMaxConnectionCount, apiData.ServiceMqttMaxConnectionCount)
	AssignIfDstNotNil(&tfData.ServiceMqttPlainTextEnabled, apiData.ServiceMqttPlainTextEnabled)
	AssignIfDstNotNil(&tfData.ServiceMqttPlainTextListenPort, apiData.ServiceMqttPlainTextListenPort)
	AssignIfDstNotNil(&tfData.ServiceMqttTlsEnabled, apiData.ServiceMqttTlsEnabled)
	AssignIfDstNotNil(&tfData.ServiceMqttTlsListenPort, apiData.ServiceMqttTlsListenPort)
	AssignIfDstNotNil(&tfData.ServiceMqttTlsWebSocketEnabled, apiData.ServiceMqttTlsWebSocketEnabled)
	AssignIfDstNotNil(&tfData.ServiceMqttTlsWebSocketListenPort, apiData.ServiceMqttTlsWebSocketListenPort)
	AssignIfDstNotNil(&tfData.ServiceMqttWebSocketEnabled, apiData.ServiceMqttWebSocketEnabled)
	AssignIfDstNotNil(&tfData.ServiceMqttWebSocketListenPort, apiData.ServiceMqttWebSocketListenPort)
	AssignIfDstNotNil(&tfData.ServiceRestIncomingAuthenticationClientCertRequest, apiData.ServiceRestIncomingAuthenticationClientCertRequest)
	AssignIfDstNotNil(&tfData.ServiceRestIncomingAuthorizationHeaderHandling, apiData.ServiceRestIncomingAuthorizationHeaderHandling)
	AssignIfDstNotNil(&tfData.ServiceRestIncomingMaxConnectionCount, apiData.ServiceRestIncomingMaxConnectionCount)
	AssignIfDstNotNil(&tfData.ServiceRestIncomingPlainTextEnabled, apiData.ServiceRestIncomingPlainTextEnabled)
	AssignIfDstNotNil(&tfData.ServiceRestIncomingPlainTextListenPort, apiData.ServiceRestIncomingPlainTextListenPort)
	AssignIfDstNotNil(&tfData.ServiceRestIncomingTlsEnabled, apiData.ServiceRestIncomingTlsEnabled)
	AssignIfDstNotNil(&tfData.ServiceRestIncomingTlsListenPort, apiData.ServiceRestIncomingTlsListenPort)
	AssignIfDstNotNil(&tfData.ServiceRestMode, apiData.ServiceRestMode)
	AssignIfDstNotNil(&tfData.ServiceRestOutgoingMaxConnectionCount, apiData.ServiceRestOutgoingMaxConnectionCount)
	AssignIfDstNotNil(&tfData.ServiceSmfMaxConnectionCount, apiData.ServiceSmfMaxConnectionCount)
	AssignIfDstNotNil(&tfData.ServiceSmfPlainTextEnabled, apiData.ServiceSmfPlainTextEnabled)
	AssignIfDstNotNil(&tfData.ServiceSmfTlsEnabled, apiData.ServiceSmfTlsEnabled)
	AssignIfDstNotNil(&tfData.ServiceWebAuthenticationClientCertRequest, apiData.ServiceWebAuthenticationClientCertRequest)
	AssignIfDstNotNil(&tfData.ServiceWebMaxConnectionCount, apiData.ServiceWebMaxConnectionCount)
	AssignIfDstNotNil(&tfData.ServiceWebPlainTextEnabled, apiData.ServiceWebPlainTextEnabled)
	AssignIfDstNotNil(&tfData.ServiceWebTlsEnabled, apiData.ServiceWebTlsEnabled)
	AssignIfDstNotNil(&tfData.TlsAllowDowngradeToPlainTextEnabled, apiData.TlsAllowDowngradeToPlainTextEnabled)
}

func (tfData *MsgVpn) ToApi() *sempv2.MsgVpn {
	return &sempv2.MsgVpn{
		Alias:                           tfData.Alias,
		AuthenticationBasicEnabled:      tfData.AuthenticationBasicEnabled,
		AuthenticationBasicProfileName:  tfData.AuthenticationBasicProfileName,
		AuthenticationBasicRadiusDomain: tfData.AuthenticationBasicRadiusDomain,
		AuthenticationBasicType:         tfData.AuthenticationBasicType,
		AuthenticationClientCertAllowApiProvidedUsernameEnabled: tfData.AuthenticationClientCertAllowApiProvidedUsernameEnabled,
		AuthenticationClientCertEnabled:                         tfData.AuthenticationClientCertEnabled,
		AuthenticationClientCertMaxChainDepth:                   tfData.AuthenticationClientCertMaxChainDepth,
		AuthenticationClientCertRevocationCheckMode:             tfData.AuthenticationClientCertRevocationCheckMode,
		AuthenticationClientCertUsernameSource:                  tfData.AuthenticationClientCertUsernameSource,
		AuthenticationClientCertValidateDateEnabled:             tfData.AuthenticationClientCertValidateDateEnabled,
		AuthenticationKerberosAllowApiProvidedUsernameEnabled:   tfData.AuthenticationKerberosAllowApiProvidedUsernameEnabled,
		AuthenticationKerberosEnabled:                           tfData.AuthenticationKerberosEnabled,
		AuthenticationOauthDefaultProfileName:                   tfData.AuthenticationOauthDefaultProfileName,
		AuthenticationOauthDefaultProviderName:                  tfData.AuthenticationOauthDefaultProviderName,
		AuthenticationOauthEnabled:                              tfData.AuthenticationOauthEnabled,
		AuthorizationLdapGroupMembershipAttributeName:           tfData.AuthorizationLdapGroupMembershipAttributeName,
		AuthorizationLdapTrimClientUsernameDomainEnabled:        tfData.AuthorizationLdapTrimClientUsernameDomainEnabled,
		AuthorizationProfileName:                                tfData.AuthorizationProfileName,
		AuthorizationType:                                       tfData.AuthorizationType,
		BridgingTlsServerCertEnforceTrustedCommonNameEnabled:    tfData.BridgingTlsServerCertEnforceTrustedCommonNameEnabled,
		BridgingTlsServerCertMaxChainDepth:                      tfData.BridgingTlsServerCertMaxChainDepth,
		BridgingTlsServerCertValidateDateEnabled:                tfData.BridgingTlsServerCertValidateDateEnabled,
		BridgingTlsServerCertValidateNameEnabled:                tfData.BridgingTlsServerCertValidateNameEnabled,
		DistributedCacheManagementEnabled:                       tfData.DistributedCacheManagementEnabled,
		DmrEnabled:                                              tfData.DmrEnabled,
		Enabled:                                                 tfData.Enabled,
		EventConnectionCountThreshold:                           tfData.EventConnectionCountThreshold.ToApi(),
		EventEgressFlowCountThreshold:                           tfData.EventEgressFlowCountThreshold.ToApi(),
		EventEgressMsgRateThreshold:                             tfData.EventEgressMsgRateThreshold.ToApi(),
		EventEndpointCountThreshold:                             tfData.EventEndpointCountThreshold.ToApi(),
		EventIngressFlowCountThreshold:                          tfData.EventIngressFlowCountThreshold.ToApi(),
		EventIngressMsgRateThreshold:                            tfData.EventIngressMsgRateThreshold.ToApi(),
		EventLargeMsgThreshold:                                  tfData.EventLargeMsgThreshold,
		EventLogTag:                                             tfData.EventLogTag,
		EventMsgSpoolUsageThreshold:                             tfData.EventMsgSpoolUsageThreshold.ToApi(),
		EventPublishClientEnabled:                               tfData.EventPublishClientEnabled,
		EventPublishMsgVpnEnabled:                               tfData.EventPublishMsgVpnEnabled,
		EventPublishSubscriptionMode:                            tfData.EventPublishSubscriptionMode,
		EventPublishTopicFormatMqttEnabled:                      tfData.EventPublishTopicFormatMqttEnabled,
		EventPublishTopicFormatSmfEnabled:                       tfData.EventPublishTopicFormatSmfEnabled,
		EventServiceAmqpConnectionCountThreshold:                tfData.EventServiceAmqpConnectionCountThreshold.ToApi(),
		EventServiceMqttConnectionCountThreshold:                tfData.EventServiceMqttConnectionCountThreshold.ToApi(),
		EventServiceRestIncomingConnectionCountThreshold:        tfData.EventServiceRestIncomingConnectionCountThreshold.ToApi(),
		EventServiceSmfConnectionCountThreshold:                 tfData.EventServiceSmfConnectionCountThreshold.ToApi(),
		EventServiceWebConnectionCountThreshold:                 tfData.EventServiceWebConnectionCountThreshold.ToApi(),
		EventSubscriptionCountThreshold:                         tfData.EventSubscriptionCountThreshold.ToApi(),
		EventTransactedSessionCountThreshold:                    tfData.EventTransactedSessionCountThreshold.ToApi(),
		EventTransactionCountThreshold:                          tfData.EventTransactionCountThreshold.ToApi(),
		ExportSubscriptionsEnabled:                              tfData.ExportSubscriptionsEnabled,
		JndiEnabled:                                             tfData.JndiEnabled,
		MaxConnectionCount:                                      tfData.MaxConnectionCount,
		MaxEgressFlowCount:                                      tfData.MaxEgressFlowCount,
		MaxEndpointCount:                                        tfData.MaxEndpointCount,
		MaxIngressFlowCount:                                     tfData.MaxIngressFlowCount,
		MaxMsgSpoolUsage:                                        tfData.MaxMsgSpoolUsage,
		MaxSubscriptionCount:                                    tfData.MaxSubscriptionCount,
		MaxTransactedSessionCount:                               tfData.MaxTransactedSessionCount,
		MaxTransactionCount:                                     tfData.MaxTransactionCount,
		MqttRetainMaxMemory:                                     tfData.MqttRetainMaxMemory,
		MsgVpnName:                                              tfData.MsgVpnName,
		ReplicationAckPropagationIntervalMsgCount:               tfData.ReplicationAckPropagationIntervalMsgCount,
		ReplicationBridgeAuthenticationBasicClientUsername:      tfData.ReplicationBridgeAuthenticationBasicClientUsername,
		ReplicationBridgeAuthenticationBasicPassword:            tfData.ReplicationBridgeAuthenticationBasicPassword,
		ReplicationBridgeAuthenticationClientCertContent:        tfData.ReplicationBridgeAuthenticationClientCertContent,
		ReplicationBridgeAuthenticationClientCertPassword:       tfData.ReplicationBridgeAuthenticationClientCertPassword,
		ReplicationBridgeAuthenticationScheme:                   tfData.ReplicationBridgeAuthenticationScheme,
		ReplicationBridgeCompressedDataEnabled:                  tfData.ReplicationBridgeCompressedDataEnabled,
		ReplicationBridgeEgressFlowWindowSize:                   tfData.ReplicationBridgeEgressFlowWindowSize,
		ReplicationBridgeRetryDelay:                             tfData.ReplicationBridgeRetryDelay,
		ReplicationBridgeTlsEnabled:                             tfData.ReplicationBridgeTlsEnabled,
		ReplicationBridgeUnidirectionalClientProfileName:        tfData.ReplicationBridgeUnidirectionalClientProfileName,
		ReplicationEnabled:                                      tfData.ReplicationEnabled,
		ReplicationEnabledQueueBehavior:                         tfData.ReplicationEnabledQueueBehavior,
		ReplicationQueueMaxMsgSpoolUsage:                        tfData.ReplicationQueueMaxMsgSpoolUsage,
		ReplicationQueueRejectMsgToSenderOnDiscardEnabled:       tfData.ReplicationQueueRejectMsgToSenderOnDiscardEnabled,
		ReplicationRejectMsgWhenSyncIneligibleEnabled:           tfData.ReplicationRejectMsgWhenSyncIneligibleEnabled,
		ReplicationRole:                                         tfData.ReplicationRole,
		ReplicationTransactionMode:                              tfData.ReplicationTransactionMode,
		RestTlsServerCertEnforceTrustedCommonNameEnabled:        tfData.RestTlsServerCertEnforceTrustedCommonNameEnabled,
		RestTlsServerCertMaxChainDepth:                          tfData.RestTlsServerCertMaxChainDepth,
		RestTlsServerCertValidateDateEnabled:                    tfData.RestTlsServerCertValidateDateEnabled,
		RestTlsServerCertValidateNameEnabled:                    tfData.RestTlsServerCertValidateNameEnabled,
		SempOverMsgBusAdminClientEnabled:                        tfData.SempOverMsgBusAdminClientEnabled,
		SempOverMsgBusAdminDistributedCacheEnabled:              tfData.SempOverMsgBusAdminDistributedCacheEnabled,
		SempOverMsgBusAdminEnabled:                              tfData.SempOverMsgBusAdminEnabled,
		SempOverMsgBusEnabled:                                   tfData.SempOverMsgBusEnabled,
		SempOverMsgBusShowEnabled:                               tfData.SempOverMsgBusShowEnabled,
		ServiceAmqpMaxConnectionCount:                           tfData.ServiceAmqpMaxConnectionCount,
		ServiceAmqpPlainTextEnabled:                             tfData.ServiceAmqpPlainTextEnabled,
		ServiceAmqpPlainTextListenPort:                          tfData.ServiceAmqpPlainTextListenPort,
		ServiceAmqpTlsEnabled:                                   tfData.ServiceAmqpTlsEnabled,
		ServiceAmqpTlsListenPort:                                tfData.ServiceAmqpTlsListenPort,
		ServiceMqttAuthenticationClientCertRequest:              tfData.ServiceMqttAuthenticationClientCertRequest,
		ServiceMqttMaxConnectionCount:                           tfData.ServiceMqttMaxConnectionCount,
		ServiceMqttPlainTextEnabled:                             tfData.ServiceMqttPlainTextEnabled,
		ServiceMqttPlainTextListenPort:                          tfData.ServiceMqttPlainTextListenPort,
		ServiceMqttTlsEnabled:                                   tfData.ServiceMqttTlsEnabled,
		ServiceMqttTlsListenPort:                                tfData.ServiceMqttTlsListenPort,
		ServiceMqttTlsWebSocketEnabled:                          tfData.ServiceMqttTlsWebSocketEnabled,
		ServiceMqttTlsWebSocketListenPort:                       tfData.ServiceMqttTlsWebSocketListenPort,
		ServiceMqttWebSocketEnabled:                             tfData.ServiceMqttWebSocketEnabled,
		ServiceMqttWebSocketListenPort:                          tfData.ServiceMqttWebSocketListenPort,
		ServiceRestIncomingAuthenticationClientCertRequest:      tfData.ServiceRestIncomingAuthenticationClientCertRequest,
		ServiceRestIncomingAuthorizationHeaderHandling:          tfData.ServiceRestIncomingAuthorizationHeaderHandling,
		ServiceRestIncomingMaxConnectionCount:                   tfData.ServiceRestIncomingMaxConnectionCount,
		ServiceRestIncomingPlainTextEnabled:                     tfData.ServiceRestIncomingPlainTextEnabled,
		ServiceRestIncomingPlainTextListenPort:                  tfData.ServiceRestIncomingPlainTextListenPort,
		ServiceRestIncomingTlsEnabled:                           tfData.ServiceRestIncomingTlsEnabled,
		ServiceRestIncomingTlsListenPort:                        tfData.ServiceRestIncomingTlsListenPort,
		ServiceRestMode:                                         tfData.ServiceRestMode,
		ServiceRestOutgoingMaxConnectionCount:                   tfData.ServiceRestOutgoingMaxConnectionCount,
		ServiceSmfMaxConnectionCount:                            tfData.ServiceSmfMaxConnectionCount,
		ServiceSmfPlainTextEnabled:                              tfData.ServiceSmfPlainTextEnabled,
		ServiceSmfTlsEnabled:                                    tfData.ServiceSmfTlsEnabled,
		ServiceWebAuthenticationClientCertRequest:               tfData.ServiceWebAuthenticationClientCertRequest,
		ServiceWebMaxConnectionCount:                            tfData.ServiceWebMaxConnectionCount,
		ServiceWebPlainTextEnabled:                              tfData.ServiceWebPlainTextEnabled,
		ServiceWebTlsEnabled:                                    tfData.ServiceWebTlsEnabled,
		TlsAllowDowngradeToPlainTextEnabled:                     tfData.TlsAllowDowngradeToPlainTextEnabled,
	}
}

// Terraform Resource schema for MsgVpn
func MsgVpnResourceSchema(requiredAttributes ...string) schema.Schema {
	schema := schema.Schema{
		Description: "MsgVpn",
		Attributes: map[string]schema.Attribute{
			"alias": schema.StringAttribute{
				Description: "The name of another Message VPN which this Message VPN is an alias for. When this Message VPN is enabled, the alias has no effect. When this Message VPN is disabled, Clients (but not Bridges and routing Links) logging into this Message VPN are automatically logged in to the other Message VPN, and authentication and authorization take place in the context of the other Message VPN.  Aliases may form a non-circular chain, cascading one to the next. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.14.",
				Required:    contains(requiredAttributes, "alias"),
				Optional:    !contains(requiredAttributes, "alias"),

				PlanModifiers: StringPlanModifiersFor("alias", requiredAttributes),
			},
			"authentication_basic_enabled": schema.BoolAttribute{
				Description: "Enable or disable basic authentication for clients connecting to the Message VPN. Basic authentication is authentication that involves the use of a username and password to prove identity. If a user provides credentials for a different authentication scheme, this setting is not applicable. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "authentication_basic_enabled"),
				Optional:    !contains(requiredAttributes, "authentication_basic_enabled"),
			},
			"authentication_basic_profile_name": schema.StringAttribute{
				Description: "The name of the RADIUS or LDAP Profile to use for basic authentication. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"default\"`.",
				Required:    contains(requiredAttributes, "authentication_basic_profile_name"),
				Optional:    !contains(requiredAttributes, "authentication_basic_profile_name"),

				PlanModifiers: StringPlanModifiersFor("authentication_basic_profile_name", requiredAttributes),
			},
			"authentication_basic_radius_domain": schema.StringAttribute{
				Description: "The RADIUS domain to use for basic authentication. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "authentication_basic_radius_domain"),
				Optional:    !contains(requiredAttributes, "authentication_basic_radius_domain"),

				PlanModifiers: StringPlanModifiersFor("authentication_basic_radius_domain", requiredAttributes),
			},
			"authentication_basic_type": schema.StringAttribute{
				Description: "The type of basic authentication to use for clients connecting to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"radius\"`. The allowed values and their meaning are:  <pre> \"internal\" - Internal database. Authentication is against Client Usernames. \"ldap\" - LDAP authentication. An LDAP profile name must be provided. \"radius\" - RADIUS authentication. A RADIUS profile name must be provided. \"none\" - No authentication. Anonymous login allowed. </pre> ",
				Required:    contains(requiredAttributes, "authentication_basic_type"),
				Optional:    !contains(requiredAttributes, "authentication_basic_type"),

				Validators: []validator.String{
					stringvalidator.OneOf("internal", "ldap", "radius", "none"),
				},
				PlanModifiers: StringPlanModifiersFor("authentication_basic_type", requiredAttributes),
			},
			"authentication_client_cert_allow_api_provided_username_enabled": schema.BoolAttribute{
				Description: "Enable or disable allowing a client to specify a Client Username via the API connect method. When disabled, the certificate CN (Common Name) is always used. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "authentication_client_cert_allow_api_provided_username_enabled"),
				Optional:    !contains(requiredAttributes, "authentication_client_cert_allow_api_provided_username_enabled"),
			},
			"authentication_client_cert_enabled": schema.BoolAttribute{
				Description: "Enable or disable client certificate authentication in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "authentication_client_cert_enabled"),
				Optional:    !contains(requiredAttributes, "authentication_client_cert_enabled"),
			},
			"authentication_client_cert_max_chain_depth": schema.Int64Attribute{
				Description: "The maximum depth for a client certificate chain. The depth of a chain is defined as the number of signing CA certificates that are present in the chain back to a trusted self-signed root CA certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.",
				Required:    contains(requiredAttributes, "authentication_client_cert_max_chain_depth"),
				Optional:    !contains(requiredAttributes, "authentication_client_cert_max_chain_depth"),
			},
			"authentication_client_cert_revocation_check_mode": schema.StringAttribute{
				Description: "The desired behavior for client certificate revocation checking. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"allow-valid\"`. The allowed values and their meaning are:  <pre> \"allow-all\" - Allow the client to authenticate, the result of client certificate revocation check is ignored. \"allow-unknown\" - Allow the client to authenticate even if the revocation status of his certificate cannot be determined. \"allow-valid\" - Allow the client to authenticate only when the revocation check returned an explicit positive response. </pre>  Available since 2.6.",
				Required:    contains(requiredAttributes, "authentication_client_cert_revocation_check_mode"),
				Optional:    !contains(requiredAttributes, "authentication_client_cert_revocation_check_mode"),

				Validators: []validator.String{
					stringvalidator.OneOf("allow-all", "allow-unknown", "allow-valid"),
				},
				PlanModifiers: StringPlanModifiersFor("authentication_client_cert_revocation_check_mode", requiredAttributes),
			},
			"authentication_client_cert_username_source": schema.StringAttribute{
				Description: "The field from the client certificate to use as the client username. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"common-name\"`. The allowed values and their meaning are:  <pre> \"certificate-thumbprint\" - The username is computed as the SHA-1 hash over the entire DER-encoded contents of the client certificate. \"common-name\" - The username is extracted from the certificate's first instance of the Common Name attribute in the Subject DN. \"common-name-last\" - The username is extracted from the certificate's last instance of the Common Name attribute in the Subject DN. \"subject-alternate-name-msupn\" - The username is extracted from the certificate's Other Name type of the Subject Alternative Name and must have the msUPN signature. \"uid\" - The username is extracted from the certificate's first instance of the User Identifier attribute in the Subject DN. \"uid-last\" - The username is extracted from the certificate's last instance of the User Identifier attribute in the Subject DN. </pre>  Available since 2.6.",
				Required:    contains(requiredAttributes, "authentication_client_cert_username_source"),
				Optional:    !contains(requiredAttributes, "authentication_client_cert_username_source"),

				Validators: []validator.String{
					stringvalidator.OneOf("certificate-thumbprint", "common-name", "common-name-last", "subject-alternate-name-msupn", "uid", "uid-last"),
				},
				PlanModifiers: StringPlanModifiersFor("authentication_client_cert_username_source", requiredAttributes),
			},
			"authentication_client_cert_validate_date_enabled": schema.BoolAttribute{
				Description: "Enable or disable validation of the \"Not Before\" and \"Not After\" validity dates in the client certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "authentication_client_cert_validate_date_enabled"),
				Optional:    !contains(requiredAttributes, "authentication_client_cert_validate_date_enabled"),
			},
			"authentication_kerberos_allow_api_provided_username_enabled": schema.BoolAttribute{
				Description: "Enable or disable allowing a client to specify a Client Username via the API connect method. When disabled, the Kerberos Principal name is always used. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "authentication_kerberos_allow_api_provided_username_enabled"),
				Optional:    !contains(requiredAttributes, "authentication_kerberos_allow_api_provided_username_enabled"),
			},
			"authentication_kerberos_enabled": schema.BoolAttribute{
				Description: "Enable or disable Kerberos authentication in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "authentication_kerberos_enabled"),
				Optional:    !contains(requiredAttributes, "authentication_kerberos_enabled"),
			},
			"authentication_oauth_default_profile_name": schema.StringAttribute{
				Description: "The name of the profile to use when the client does not supply a profile name. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.25.",
				Required:    contains(requiredAttributes, "authentication_oauth_default_profile_name"),
				Optional:    !contains(requiredAttributes, "authentication_oauth_default_profile_name"),

				PlanModifiers: StringPlanModifiersFor("authentication_oauth_default_profile_name", requiredAttributes),
			},
			"authentication_oauth_default_provider_name": schema.StringAttribute{
				Description: "The name of the provider to use when the client does not supply a provider name. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Deprecated since 2.25. authenticationOauthDefaultProviderName and authenticationOauthProviders replaced by authenticationOauthDefaultProfileName and authenticationOauthProfiles.",
				Required:    contains(requiredAttributes, "authentication_oauth_default_provider_name"),
				Optional:    !contains(requiredAttributes, "authentication_oauth_default_provider_name"),

				PlanModifiers: StringPlanModifiersFor("authentication_oauth_default_provider_name", requiredAttributes),
			},
			"authentication_oauth_enabled": schema.BoolAttribute{
				Description: "Enable or disable OAuth authentication. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.13.",
				Required:    contains(requiredAttributes, "authentication_oauth_enabled"),
				Optional:    !contains(requiredAttributes, "authentication_oauth_enabled"),
			},
			"authorization_ldap_group_membership_attribute_name": schema.StringAttribute{
				Description: "The name of the attribute that is retrieved from the LDAP server as part of the LDAP search when authorizing a client connecting to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"memberOf\"`.",
				Required:    contains(requiredAttributes, "authorization_ldap_group_membership_attribute_name"),
				Optional:    !contains(requiredAttributes, "authorization_ldap_group_membership_attribute_name"),

				PlanModifiers: StringPlanModifiersFor("authorization_ldap_group_membership_attribute_name", requiredAttributes),
			},
			"authorization_ldap_trim_client_username_domain_enabled": schema.BoolAttribute{
				Description: "Enable or disable client-username domain trimming for LDAP lookups of client connections. When enabled, the value of $CLIENT_USERNAME (when used for searching) will be truncated at the first occurance of the @ character. For example, if the client-username is in the form of an email address, then the domain portion will be removed. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.13.",
				Required:    contains(requiredAttributes, "authorization_ldap_trim_client_username_domain_enabled"),
				Optional:    !contains(requiredAttributes, "authorization_ldap_trim_client_username_domain_enabled"),
			},
			"authorization_profile_name": schema.StringAttribute{
				Description: "The name of the LDAP Profile to use for client authorization. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "authorization_profile_name"),
				Optional:    !contains(requiredAttributes, "authorization_profile_name"),

				PlanModifiers: StringPlanModifiersFor("authorization_profile_name", requiredAttributes),
			},
			"authorization_type": schema.StringAttribute{
				Description: "The type of authorization to use for clients connecting to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"internal\"`. The allowed values and their meaning are:  <pre> \"ldap\" - LDAP authorization. \"internal\" - Internal authorization. </pre> ",
				Required:    contains(requiredAttributes, "authorization_type"),
				Optional:    !contains(requiredAttributes, "authorization_type"),

				Validators: []validator.String{
					stringvalidator.OneOf("ldap", "internal"),
				},
				PlanModifiers: StringPlanModifiersFor("authorization_type", requiredAttributes),
			},
			"bridging_tls_server_cert_enforce_trusted_common_name_enabled": schema.BoolAttribute{
				Description: "Enable or disable validation of the Common Name (CN) in the server certificate from the remote broker. If enabled, the Common Name is checked against the list of Trusted Common Names configured for the Bridge. Common Name validation is not performed if Server Certificate Name Validation is enabled, even if Common Name validation is enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Deprecated since 2.18. Common Name validation has been replaced by Server Certificate Name validation.",
				Required:    contains(requiredAttributes, "bridging_tls_server_cert_enforce_trusted_common_name_enabled"),
				Optional:    !contains(requiredAttributes, "bridging_tls_server_cert_enforce_trusted_common_name_enabled"),
			},
			"bridging_tls_server_cert_max_chain_depth": schema.Int64Attribute{
				Description: "The maximum depth for a server certificate chain. The depth of a chain is defined as the number of signing CA certificates that are present in the chain back to a trusted self-signed root CA certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.",
				Required:    contains(requiredAttributes, "bridging_tls_server_cert_max_chain_depth"),
				Optional:    !contains(requiredAttributes, "bridging_tls_server_cert_max_chain_depth"),
			},
			"bridging_tls_server_cert_validate_date_enabled": schema.BoolAttribute{
				Description: "Enable or disable validation of the \"Not Before\" and \"Not After\" validity dates in the server certificate. When disabled, a certificate will be accepted even if the certificate is not valid based on these dates. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "bridging_tls_server_cert_validate_date_enabled"),
				Optional:    !contains(requiredAttributes, "bridging_tls_server_cert_validate_date_enabled"),
			},
			"bridging_tls_server_cert_validate_name_enabled": schema.BoolAttribute{
				Description: "Enable or disable the standard TLS authentication mechanism of verifying the name used to connect to the bridge. If enabled, the name used to connect to the bridge is checked against the names specified in the certificate returned by the remote router. Legacy Common Name validation is not performed if Server Certificate Name Validation is enabled, even if Common Name validation is also enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`. Available since 2.18.",
				Required:    contains(requiredAttributes, "bridging_tls_server_cert_validate_name_enabled"),
				Optional:    !contains(requiredAttributes, "bridging_tls_server_cert_validate_name_enabled"),
			},
			"distributed_cache_management_enabled": schema.BoolAttribute{
				Description: "Enable or disable managing of cache instances over the message bus. The default value is `true`.",
				Required:    contains(requiredAttributes, "distributed_cache_management_enabled"),
				Optional:    !contains(requiredAttributes, "distributed_cache_management_enabled"),
			},
			"dmr_enabled": schema.BoolAttribute{
				Description: "Enable or disable Dynamic Message Routing (DMR) for the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.11.",
				Required:    contains(requiredAttributes, "dmr_enabled"),
				Optional:    !contains(requiredAttributes, "dmr_enabled"),
			},
			"enabled": schema.BoolAttribute{
				Description: "Enable or disable the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "enabled"),
				Optional:    !contains(requiredAttributes, "enabled"),
			},
			"event_connection_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_connection_count_threshold"),
				Optional:    !contains(requiredAttributes, "event_connection_count_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"event_egress_flow_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_egress_flow_count_threshold"),
				Optional:    !contains(requiredAttributes, "event_egress_flow_count_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"event_egress_msg_rate_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_egress_msg_rate_threshold"),
				Optional:    !contains(requiredAttributes, "event_egress_msg_rate_threshold"),

				Attributes: EventThresholdByValueResourceAttributes,
			},
			"event_endpoint_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_endpoint_count_threshold"),
				Optional:    !contains(requiredAttributes, "event_endpoint_count_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"event_ingress_flow_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_ingress_flow_count_threshold"),
				Optional:    !contains(requiredAttributes, "event_ingress_flow_count_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"event_ingress_msg_rate_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_ingress_msg_rate_threshold"),
				Optional:    !contains(requiredAttributes, "event_ingress_msg_rate_threshold"),

				Attributes: EventThresholdByValueResourceAttributes,
			},
			"event_large_msg_threshold": schema.Int64Attribute{
				Description: "The threshold, in kilobytes, after which a message is considered to be large for the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1024`.",
				Required:    contains(requiredAttributes, "event_large_msg_threshold"),
				Optional:    !contains(requiredAttributes, "event_large_msg_threshold"),
			},
			"event_log_tag": schema.StringAttribute{
				Description: "A prefix applied to all published Events in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "event_log_tag"),
				Optional:    !contains(requiredAttributes, "event_log_tag"),

				PlanModifiers: StringPlanModifiersFor("event_log_tag", requiredAttributes),
			},
			"event_msg_spool_usage_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_msg_spool_usage_threshold"),
				Optional:    !contains(requiredAttributes, "event_msg_spool_usage_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"event_publish_client_enabled": schema.BoolAttribute{
				Description: "Enable or disable Client level Event message publishing. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "event_publish_client_enabled"),
				Optional:    !contains(requiredAttributes, "event_publish_client_enabled"),
			},
			"event_publish_msg_vpn_enabled": schema.BoolAttribute{
				Description: "Enable or disable Message VPN level Event message publishing. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "event_publish_msg_vpn_enabled"),
				Optional:    !contains(requiredAttributes, "event_publish_msg_vpn_enabled"),
			},
			"event_publish_subscription_mode": schema.StringAttribute{
				Description: "Subscription level Event message publishing mode. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"off\"`. The allowed values and their meaning are:  <pre> \"off\" - Disable client level event message publishing. \"on-with-format-v1\" - Enable client level event message publishing with format v1. \"on-with-no-unsubscribe-events-on-disconnect-format-v1\" - As \"on-with-format-v1\", but unsubscribe events are not generated when a client disconnects. Unsubscribe events are still raised when a client explicitly unsubscribes from its subscriptions. \"on-with-format-v2\" - Enable client level event message publishing with format v2. \"on-with-no-unsubscribe-events-on-disconnect-format-v2\" - As \"on-with-format-v2\", but unsubscribe events are not generated when a client disconnects. Unsubscribe events are still raised when a client explicitly unsubscribes from its subscriptions. </pre> ",
				Required:    contains(requiredAttributes, "event_publish_subscription_mode"),
				Optional:    !contains(requiredAttributes, "event_publish_subscription_mode"),

				Validators: []validator.String{
					stringvalidator.OneOf("off", "on-with-format-v1", "on-with-no-unsubscribe-events-on-disconnect-format-v1", "on-with-format-v2", "on-with-no-unsubscribe-events-on-disconnect-format-v2"),
				},
				PlanModifiers: StringPlanModifiersFor("event_publish_subscription_mode", requiredAttributes),
			},
			"event_publish_topic_format_mqtt_enabled": schema.BoolAttribute{
				Description: "Enable or disable Event publish topics in MQTT format. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "event_publish_topic_format_mqtt_enabled"),
				Optional:    !contains(requiredAttributes, "event_publish_topic_format_mqtt_enabled"),
			},
			"event_publish_topic_format_smf_enabled": schema.BoolAttribute{
				Description: "Enable or disable Event publish topics in SMF format. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "event_publish_topic_format_smf_enabled"),
				Optional:    !contains(requiredAttributes, "event_publish_topic_format_smf_enabled"),
			},
			"event_service_amqp_connection_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_service_amqp_connection_count_threshold"),
				Optional:    !contains(requiredAttributes, "event_service_amqp_connection_count_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"event_service_mqtt_connection_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_service_mqtt_connection_count_threshold"),
				Optional:    !contains(requiredAttributes, "event_service_mqtt_connection_count_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"event_service_rest_incoming_connection_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_service_rest_incoming_connection_count_threshold"),
				Optional:    !contains(requiredAttributes, "event_service_rest_incoming_connection_count_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"event_service_smf_connection_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_service_smf_connection_count_threshold"),
				Optional:    !contains(requiredAttributes, "event_service_smf_connection_count_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"event_service_web_connection_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_service_web_connection_count_threshold"),
				Optional:    !contains(requiredAttributes, "event_service_web_connection_count_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"event_subscription_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_subscription_count_threshold"),
				Optional:    !contains(requiredAttributes, "event_subscription_count_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"event_transacted_session_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_transacted_session_count_threshold"),
				Optional:    !contains(requiredAttributes, "event_transacted_session_count_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"event_transaction_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_transaction_count_threshold"),
				Optional:    !contains(requiredAttributes, "event_transaction_count_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"export_subscriptions_enabled": schema.BoolAttribute{
				Description: "Enable or disable the export of subscriptions in the Message VPN to other routers in the network over Neighbor links. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "export_subscriptions_enabled"),
				Optional:    !contains(requiredAttributes, "export_subscriptions_enabled"),
			},
			"jndi_enabled": schema.BoolAttribute{
				Description: "Enable or disable JNDI access for clients in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.2.",
				Required:    contains(requiredAttributes, "jndi_enabled"),
				Optional:    !contains(requiredAttributes, "jndi_enabled"),
			},
			"max_connection_count": schema.Int64Attribute{
				Description: "The maximum number of client connections to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform.",
				Required:    contains(requiredAttributes, "max_connection_count"),
				Optional:    !contains(requiredAttributes, "max_connection_count"),
			},
			"max_egress_flow_count": schema.Int64Attribute{
				Description: "The maximum number of transmit flows that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.",
				Required:    contains(requiredAttributes, "max_egress_flow_count"),
				Optional:    !contains(requiredAttributes, "max_egress_flow_count"),
			},
			"max_endpoint_count": schema.Int64Attribute{
				Description: "The maximum number of Queues and Topic Endpoints that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.",
				Required:    contains(requiredAttributes, "max_endpoint_count"),
				Optional:    !contains(requiredAttributes, "max_endpoint_count"),
			},
			"max_ingress_flow_count": schema.Int64Attribute{
				Description: "The maximum number of receive flows that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.",
				Required:    contains(requiredAttributes, "max_ingress_flow_count"),
				Optional:    !contains(requiredAttributes, "max_ingress_flow_count"),
			},
			"max_msg_spool_usage": schema.Int64Attribute{
				Description: "The maximum message spool usage by the Message VPN, in megabytes. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.",
				Required:    contains(requiredAttributes, "max_msg_spool_usage"),
				Optional:    !contains(requiredAttributes, "max_msg_spool_usage"),
			},
			"max_subscription_count": schema.Int64Attribute{
				Description: "The maximum number of local client subscriptions that can be added to the Message VPN. This limit is not enforced when a subscription is added using a management interface, such as CLI or SEMP. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.",
				Required:    contains(requiredAttributes, "max_subscription_count"),
				Optional:    !contains(requiredAttributes, "max_subscription_count"),
			},
			"max_transacted_session_count": schema.Int64Attribute{
				Description: "The maximum number of transacted sessions that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.",
				Required:    contains(requiredAttributes, "max_transacted_session_count"),
				Optional:    !contains(requiredAttributes, "max_transacted_session_count"),
			},
			"max_transaction_count": schema.Int64Attribute{
				Description: "The maximum number of transactions that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.",
				Required:    contains(requiredAttributes, "max_transaction_count"),
				Optional:    !contains(requiredAttributes, "max_transaction_count"),
			},
			"mqtt_retain_max_memory": schema.Int64Attribute{
				Description: "The maximum total memory usage of the MQTT Retain feature for this Message VPN, in MB. If the maximum memory is reached, any arriving retain messages that require more memory are discarded. A value of -1 indicates that the memory is bounded only by the global max memory limit. A value of 0 prevents MQTT Retain from becoming operational. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `-1`. Available since 2.11.",
				Required:    contains(requiredAttributes, "mqtt_retain_max_memory"),
				Optional:    !contains(requiredAttributes, "mqtt_retain_max_memory"),
			},
			"msg_vpn_name": schema.StringAttribute{
				Description: "The name of the Message VPN.",
				Required:    contains(requiredAttributes, "msg_vpn_name"),
				Optional:    !contains(requiredAttributes, "msg_vpn_name"),

				PlanModifiers: StringPlanModifiersFor("msg_vpn_name", requiredAttributes),
			},
			"replication_ack_propagation_interval_msg_count": schema.Int64Attribute{
				Description: "The acknowledgement (ACK) propagation interval for the replication Bridge, in number of replicated messages. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `20`.",
				Required:    contains(requiredAttributes, "replication_ack_propagation_interval_msg_count"),
				Optional:    !contains(requiredAttributes, "replication_ack_propagation_interval_msg_count"),
			},
			"replication_bridge_authentication_basic_client_username": schema.StringAttribute{
				Description: "The Client Username the replication Bridge uses to login to the remote Message VPN. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "replication_bridge_authentication_basic_client_username"),
				Optional:    !contains(requiredAttributes, "replication_bridge_authentication_basic_client_username"),

				PlanModifiers: StringPlanModifiersFor("replication_bridge_authentication_basic_client_username", requiredAttributes),
			},
			"replication_bridge_authentication_basic_password": schema.StringAttribute{
				Description: "The password for the Client Username. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "replication_bridge_authentication_basic_password"),
				Optional:    !contains(requiredAttributes, "replication_bridge_authentication_basic_password"),

				PlanModifiers: StringPlanModifiersFor("replication_bridge_authentication_basic_password", requiredAttributes),
			},
			"replication_bridge_authentication_client_cert_content": schema.StringAttribute{
				Description: "The PEM formatted content for the client certificate used by this bridge to login to the Remote Message VPN. It must consist of a private key and between one and three certificates comprising the certificate trust chain. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changing this attribute requires an HTTPS connection. The default value is `\"\"`. Available since 2.9.",
				Required:    contains(requiredAttributes, "replication_bridge_authentication_client_cert_content"),
				Optional:    !contains(requiredAttributes, "replication_bridge_authentication_client_cert_content"),

				PlanModifiers: StringPlanModifiersFor("replication_bridge_authentication_client_cert_content", requiredAttributes),
			},
			"replication_bridge_authentication_client_cert_password": schema.StringAttribute{
				Description: "The password for the client certificate. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changing this attribute requires an HTTPS connection. The default value is `\"\"`. Available since 2.9.",
				Required:    contains(requiredAttributes, "replication_bridge_authentication_client_cert_password"),
				Optional:    !contains(requiredAttributes, "replication_bridge_authentication_client_cert_password"),

				PlanModifiers: StringPlanModifiersFor("replication_bridge_authentication_client_cert_password", requiredAttributes),
			},
			"replication_bridge_authentication_scheme": schema.StringAttribute{
				Description: "The authentication scheme for the replication Bridge in the Message VPN. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"basic\"`. The allowed values and their meaning are:  <pre> \"basic\" - Basic Authentication Scheme (via username and password). \"client-certificate\" - Client Certificate Authentication Scheme (via certificate file or content). </pre> ",
				Required:    contains(requiredAttributes, "replication_bridge_authentication_scheme"),
				Optional:    !contains(requiredAttributes, "replication_bridge_authentication_scheme"),

				Validators: []validator.String{
					stringvalidator.OneOf("basic", "client-certificate"),
				},
				PlanModifiers: StringPlanModifiersFor("replication_bridge_authentication_scheme", requiredAttributes),
			},
			"replication_bridge_compressed_data_enabled": schema.BoolAttribute{
				Description: "Enable or disable use of compression for the replication Bridge. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "replication_bridge_compressed_data_enabled"),
				Optional:    !contains(requiredAttributes, "replication_bridge_compressed_data_enabled"),
			},
			"replication_bridge_egress_flow_window_size": schema.Int64Attribute{
				Description: "The size of the window used for guaranteed messages published to the replication Bridge, in messages. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `255`.",
				Required:    contains(requiredAttributes, "replication_bridge_egress_flow_window_size"),
				Optional:    !contains(requiredAttributes, "replication_bridge_egress_flow_window_size"),
			},
			"replication_bridge_retry_delay": schema.Int64Attribute{
				Description: "The number of seconds that must pass before retrying the replication Bridge connection. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `3`.",
				Required:    contains(requiredAttributes, "replication_bridge_retry_delay"),
				Optional:    !contains(requiredAttributes, "replication_bridge_retry_delay"),
			},
			"replication_bridge_tls_enabled": schema.BoolAttribute{
				Description: "Enable or disable use of encryption (TLS) for the replication Bridge connection. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "replication_bridge_tls_enabled"),
				Optional:    !contains(requiredAttributes, "replication_bridge_tls_enabled"),
			},
			"replication_bridge_unidirectional_client_profile_name": schema.StringAttribute{
				Description: "The Client Profile for the unidirectional replication Bridge in the Message VPN. It is used only for the TCP parameters. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"#client-profile\"`.",
				Required:    contains(requiredAttributes, "replication_bridge_unidirectional_client_profile_name"),
				Optional:    !contains(requiredAttributes, "replication_bridge_unidirectional_client_profile_name"),

				PlanModifiers: StringPlanModifiersFor("replication_bridge_unidirectional_client_profile_name", requiredAttributes),
			},
			"replication_enabled": schema.BoolAttribute{
				Description: "Enable or disable replication for the Message VPN. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "replication_enabled"),
				Optional:    !contains(requiredAttributes, "replication_enabled"),
			},
			"replication_enabled_queue_behavior": schema.StringAttribute{
				Description: "The behavior to take when enabling replication for the Message VPN, depending on the existence of the replication Queue. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"fail-on-existing-queue\"`. The allowed values and their meaning are:  <pre> \"fail-on-existing-queue\" - The data replication queue must not already exist. \"force-use-existing-queue\" - The data replication queue must already exist. Any data messages on the Queue will be forwarded to interested applications. IMPORTANT: Before using this mode be certain that the messages are not stale or otherwise unsuitable to be forwarded. This mode can only be specified when the existing queue is configured the same as is currently specified under replication configuration otherwise the enabling of replication will fail. \"force-recreate-queue\" - The data replication queue must already exist. Any data messages on the Queue will be discarded. IMPORTANT: Before using this mode be certain that the messages on the existing data replication queue are not needed by interested applications. </pre> ",
				Required:    contains(requiredAttributes, "replication_enabled_queue_behavior"),
				Optional:    !contains(requiredAttributes, "replication_enabled_queue_behavior"),

				Validators: []validator.String{
					stringvalidator.OneOf("fail-on-existing-queue", "force-use-existing-queue", "force-recreate-queue"),
				},
				PlanModifiers: StringPlanModifiersFor("replication_enabled_queue_behavior", requiredAttributes),
			},
			"replication_queue_max_msg_spool_usage": schema.Int64Attribute{
				Description: "The maximum message spool usage by the replication Bridge local Queue (quota), in megabytes. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `60000`.",
				Required:    contains(requiredAttributes, "replication_queue_max_msg_spool_usage"),
				Optional:    !contains(requiredAttributes, "replication_queue_max_msg_spool_usage"),
			},
			"replication_queue_reject_msg_to_sender_on_discard_enabled": schema.BoolAttribute{
				Description: "Enable or disable whether messages discarded on the replication Bridge local Queue are rejected back to the sender. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "replication_queue_reject_msg_to_sender_on_discard_enabled"),
				Optional:    !contains(requiredAttributes, "replication_queue_reject_msg_to_sender_on_discard_enabled"),
			},
			"replication_reject_msg_when_sync_ineligible_enabled": schema.BoolAttribute{
				Description: "Enable or disable whether guaranteed messages published to synchronously replicated Topics are rejected back to the sender when synchronous replication becomes ineligible. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "replication_reject_msg_when_sync_ineligible_enabled"),
				Optional:    !contains(requiredAttributes, "replication_reject_msg_when_sync_ineligible_enabled"),
			},
			"replication_role": schema.StringAttribute{
				Description: "The replication role for the Message VPN. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"standby\"`. The allowed values and their meaning are:  <pre> \"active\" - Assume the Active role in replication for the Message VPN. \"standby\" - Assume the Standby role in replication for the Message VPN. </pre> ",
				Required:    contains(requiredAttributes, "replication_role"),
				Optional:    !contains(requiredAttributes, "replication_role"),

				Validators: []validator.String{
					stringvalidator.OneOf("active", "standby"),
				},
				PlanModifiers: StringPlanModifiersFor("replication_role", requiredAttributes),
			},
			"replication_transaction_mode": schema.StringAttribute{
				Description: "The transaction replication mode for all transactions within the Message VPN. Changing this value during operation will not affect existing transactions; it is only used upon starting a transaction. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"async\"`. The allowed values and their meaning are:  <pre> \"sync\" - Messages are acknowledged when replicated (spooled remotely). \"async\" - Messages are acknowledged when pending replication (spooled locally). </pre> ",
				Required:    contains(requiredAttributes, "replication_transaction_mode"),
				Optional:    !contains(requiredAttributes, "replication_transaction_mode"),

				Validators: []validator.String{
					stringvalidator.OneOf("sync", "async"),
				},
				PlanModifiers: StringPlanModifiersFor("replication_transaction_mode", requiredAttributes),
			},
			"rest_tls_server_cert_enforce_trusted_common_name_enabled": schema.BoolAttribute{
				Description: "Enable or disable validation of the Common Name (CN) in the server certificate from the remote REST Consumer. If enabled, the Common Name is checked against the list of Trusted Common Names configured for the REST Consumer. Common Name validation is not performed if Server Certificate Name Validation is enabled, even if Common Name validation is enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Deprecated since 2.17. Common Name validation has been replaced by Server Certificate Name validation.",
				Required:    contains(requiredAttributes, "rest_tls_server_cert_enforce_trusted_common_name_enabled"),
				Optional:    !contains(requiredAttributes, "rest_tls_server_cert_enforce_trusted_common_name_enabled"),
			},
			"rest_tls_server_cert_max_chain_depth": schema.Int64Attribute{
				Description: "The maximum depth for a REST Consumer server certificate chain. The depth of a chain is defined as the number of signing CA certificates that are present in the chain back to a trusted self-signed root CA certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.",
				Required:    contains(requiredAttributes, "rest_tls_server_cert_max_chain_depth"),
				Optional:    !contains(requiredAttributes, "rest_tls_server_cert_max_chain_depth"),
			},
			"rest_tls_server_cert_validate_date_enabled": schema.BoolAttribute{
				Description: "Enable or disable validation of the \"Not Before\" and \"Not After\" validity dates in the REST Consumer server certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "rest_tls_server_cert_validate_date_enabled"),
				Optional:    !contains(requiredAttributes, "rest_tls_server_cert_validate_date_enabled"),
			},
			"rest_tls_server_cert_validate_name_enabled": schema.BoolAttribute{
				Description: "Enable or disable the standard TLS authentication mechanism of verifying the name used to connect to the remote REST Consumer. If enabled, the name used to connect to the remote REST Consumer is checked against the names specified in the certificate returned by the remote router. Legacy Common Name validation is not performed if Server Certificate Name Validation is enabled, even if Common Name validation is also enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`. Available since 2.17.",
				Required:    contains(requiredAttributes, "rest_tls_server_cert_validate_name_enabled"),
				Optional:    !contains(requiredAttributes, "rest_tls_server_cert_validate_name_enabled"),
			},
			"semp_over_msg_bus_admin_client_enabled": schema.BoolAttribute{
				Description: "Enable or disable \"admin client\" SEMP over the message bus commands for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "semp_over_msg_bus_admin_client_enabled"),
				Optional:    !contains(requiredAttributes, "semp_over_msg_bus_admin_client_enabled"),
			},
			"semp_over_msg_bus_admin_distributed_cache_enabled": schema.BoolAttribute{
				Description: "Enable or disable \"admin distributed-cache\" SEMP over the message bus commands for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "semp_over_msg_bus_admin_distributed_cache_enabled"),
				Optional:    !contains(requiredAttributes, "semp_over_msg_bus_admin_distributed_cache_enabled"),
			},
			"semp_over_msg_bus_admin_enabled": schema.BoolAttribute{
				Description: "Enable or disable \"admin\" SEMP over the message bus commands for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "semp_over_msg_bus_admin_enabled"),
				Optional:    !contains(requiredAttributes, "semp_over_msg_bus_admin_enabled"),
			},
			"semp_over_msg_bus_enabled": schema.BoolAttribute{
				Description: "Enable or disable SEMP over the message bus for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "semp_over_msg_bus_enabled"),
				Optional:    !contains(requiredAttributes, "semp_over_msg_bus_enabled"),
			},
			"semp_over_msg_bus_show_enabled": schema.BoolAttribute{
				Description: "Enable or disable \"show\" SEMP over the message bus commands for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "semp_over_msg_bus_show_enabled"),
				Optional:    !contains(requiredAttributes, "semp_over_msg_bus_show_enabled"),
			},
			"service_amqp_max_connection_count": schema.Int64Attribute{
				Description: "The maximum number of AMQP client connections that can be simultaneously connected to the Message VPN. This value may be higher than supported by the platform. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform. Available since 2.7.",
				Required:    contains(requiredAttributes, "service_amqp_max_connection_count"),
				Optional:    !contains(requiredAttributes, "service_amqp_max_connection_count"),
			},
			"service_amqp_plain_text_enabled": schema.BoolAttribute{
				Description: "Enable or disable the plain-text AMQP service in the Message VPN. Disabling causes clients connected to the corresponding listen-port to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.7.",
				Required:    contains(requiredAttributes, "service_amqp_plain_text_enabled"),
				Optional:    !contains(requiredAttributes, "service_amqp_plain_text_enabled"),
			},
			"service_amqp_plain_text_listen_port": schema.Int64Attribute{
				Description: "The port number for plain-text AMQP clients that connect to the Message VPN. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.7.",
				Required:    contains(requiredAttributes, "service_amqp_plain_text_listen_port"),
				Optional:    !contains(requiredAttributes, "service_amqp_plain_text_listen_port"),
			},
			"service_amqp_tls_enabled": schema.BoolAttribute{
				Description: "Enable or disable the use of encryption (TLS) for the AMQP service in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.7.",
				Required:    contains(requiredAttributes, "service_amqp_tls_enabled"),
				Optional:    !contains(requiredAttributes, "service_amqp_tls_enabled"),
			},
			"service_amqp_tls_listen_port": schema.Int64Attribute{
				Description: "The port number for AMQP clients that connect to the Message VPN over TLS. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.7.",
				Required:    contains(requiredAttributes, "service_amqp_tls_listen_port"),
				Optional:    !contains(requiredAttributes, "service_amqp_tls_listen_port"),
			},
			"service_mqtt_authentication_client_cert_request": schema.StringAttribute{
				Description: "Determines when to request a client certificate from an incoming MQTT client connecting via a TLS port. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"when-enabled-in-message-vpn\"`. The allowed values and their meaning are:  <pre> \"always\" - Always ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"never\" - Never ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"when-enabled-in-message-vpn\" - Only ask for a client-certificate if client certificate authentication is enabled under \"message-vpn >  authentication > client-certificate > shutdown\". </pre>  Available since 2.21.",
				Required:    contains(requiredAttributes, "service_mqtt_authentication_client_cert_request"),
				Optional:    !contains(requiredAttributes, "service_mqtt_authentication_client_cert_request"),

				Validators: []validator.String{
					stringvalidator.OneOf("always", "never", "when-enabled-in-message-vpn"),
				},
				PlanModifiers: StringPlanModifiersFor("service_mqtt_authentication_client_cert_request", requiredAttributes),
			},
			"service_mqtt_max_connection_count": schema.Int64Attribute{
				Description: "The maximum number of MQTT client connections that can be simultaneously connected to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform. Available since 2.1.",
				Required:    contains(requiredAttributes, "service_mqtt_max_connection_count"),
				Optional:    !contains(requiredAttributes, "service_mqtt_max_connection_count"),
			},
			"service_mqtt_plain_text_enabled": schema.BoolAttribute{
				Description: "Enable or disable the plain-text MQTT service in the Message VPN. Disabling causes clients currently connected to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.1.",
				Required:    contains(requiredAttributes, "service_mqtt_plain_text_enabled"),
				Optional:    !contains(requiredAttributes, "service_mqtt_plain_text_enabled"),
			},
			"service_mqtt_plain_text_listen_port": schema.Int64Attribute{
				Description: "The port number for plain-text MQTT clients that connect to the Message VPN. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.1.",
				Required:    contains(requiredAttributes, "service_mqtt_plain_text_listen_port"),
				Optional:    !contains(requiredAttributes, "service_mqtt_plain_text_listen_port"),
			},
			"service_mqtt_tls_enabled": schema.BoolAttribute{
				Description: "Enable or disable the use of encryption (TLS) for the MQTT service in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.1.",
				Required:    contains(requiredAttributes, "service_mqtt_tls_enabled"),
				Optional:    !contains(requiredAttributes, "service_mqtt_tls_enabled"),
			},
			"service_mqtt_tls_listen_port": schema.Int64Attribute{
				Description: "The port number for MQTT clients that connect to the Message VPN over TLS. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.1.",
				Required:    contains(requiredAttributes, "service_mqtt_tls_listen_port"),
				Optional:    !contains(requiredAttributes, "service_mqtt_tls_listen_port"),
			},
			"service_mqtt_tls_web_socket_enabled": schema.BoolAttribute{
				Description: "Enable or disable the use of encrypted WebSocket (WebSocket over TLS) for the MQTT service in the Message VPN. Disabling causes clients currently connected by encrypted WebSocket to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.1.",
				Required:    contains(requiredAttributes, "service_mqtt_tls_web_socket_enabled"),
				Optional:    !contains(requiredAttributes, "service_mqtt_tls_web_socket_enabled"),
			},
			"service_mqtt_tls_web_socket_listen_port": schema.Int64Attribute{
				Description: "The port number for MQTT clients that connect to the Message VPN using WebSocket over TLS. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.1.",
				Required:    contains(requiredAttributes, "service_mqtt_tls_web_socket_listen_port"),
				Optional:    !contains(requiredAttributes, "service_mqtt_tls_web_socket_listen_port"),
			},
			"service_mqtt_web_socket_enabled": schema.BoolAttribute{
				Description: "Enable or disable the use of WebSocket for the MQTT service in the Message VPN. Disabling causes clients currently connected by WebSocket to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.1.",
				Required:    contains(requiredAttributes, "service_mqtt_web_socket_enabled"),
				Optional:    !contains(requiredAttributes, "service_mqtt_web_socket_enabled"),
			},
			"service_mqtt_web_socket_listen_port": schema.Int64Attribute{
				Description: "The port number for plain-text MQTT clients that connect to the Message VPN using WebSocket. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.1.",
				Required:    contains(requiredAttributes, "service_mqtt_web_socket_listen_port"),
				Optional:    !contains(requiredAttributes, "service_mqtt_web_socket_listen_port"),
			},
			"service_rest_incoming_authentication_client_cert_request": schema.StringAttribute{
				Description: "Determines when to request a client certificate from an incoming REST Producer connecting via a TLS port. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"when-enabled-in-message-vpn\"`. The allowed values and their meaning are:  <pre> \"always\" - Always ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"never\" - Never ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"when-enabled-in-message-vpn\" - Only ask for a client-certificate if client certificate authentication is enabled under \"message-vpn >  authentication > client-certificate > shutdown\". </pre>  Available since 2.21.",
				Required:    contains(requiredAttributes, "service_rest_incoming_authentication_client_cert_request"),
				Optional:    !contains(requiredAttributes, "service_rest_incoming_authentication_client_cert_request"),

				Validators: []validator.String{
					stringvalidator.OneOf("always", "never", "when-enabled-in-message-vpn"),
				},
				PlanModifiers: StringPlanModifiersFor("service_rest_incoming_authentication_client_cert_request", requiredAttributes),
			},
			"service_rest_incoming_authorization_header_handling": schema.StringAttribute{
				Description: "The handling of Authorization headers for incoming REST connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"drop\"`. The allowed values and their meaning are:  <pre> \"drop\" - Do not attach the Authorization header to the message as a user property. This configuration is most secure. \"forward\" - Forward the Authorization header, attaching it to the message as a user property in the same way as other headers. For best security, use the drop setting. \"legacy\" - If the Authorization header was used for authentication to the broker, do not attach it to the message. If the Authorization header was not used for authentication to the broker, attach it to the message as a user property in the same way as other headers. For best security, use the drop setting. </pre>  Available since 2.19.",
				Required:    contains(requiredAttributes, "service_rest_incoming_authorization_header_handling"),
				Optional:    !contains(requiredAttributes, "service_rest_incoming_authorization_header_handling"),

				Validators: []validator.String{
					stringvalidator.OneOf("drop", "forward", "legacy"),
				},
				PlanModifiers: StringPlanModifiersFor("service_rest_incoming_authorization_header_handling", requiredAttributes),
			},
			"service_rest_incoming_max_connection_count": schema.Int64Attribute{
				Description: "The maximum number of REST incoming client connections that can be simultaneously connected to the Message VPN. This value may be higher than supported by the platform. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform.",
				Required:    contains(requiredAttributes, "service_rest_incoming_max_connection_count"),
				Optional:    !contains(requiredAttributes, "service_rest_incoming_max_connection_count"),
			},
			"service_rest_incoming_plain_text_enabled": schema.BoolAttribute{
				Description: "Enable or disable the plain-text REST service for incoming clients in the Message VPN. Disabling causes clients currently connected to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "service_rest_incoming_plain_text_enabled"),
				Optional:    !contains(requiredAttributes, "service_rest_incoming_plain_text_enabled"),
			},
			"service_rest_incoming_plain_text_listen_port": schema.Int64Attribute{
				Description: "The port number for incoming plain-text REST clients that connect to the Message VPN. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.",
				Required:    contains(requiredAttributes, "service_rest_incoming_plain_text_listen_port"),
				Optional:    !contains(requiredAttributes, "service_rest_incoming_plain_text_listen_port"),
			},
			"service_rest_incoming_tls_enabled": schema.BoolAttribute{
				Description: "Enable or disable the use of encryption (TLS) for the REST service for incoming clients in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "service_rest_incoming_tls_enabled"),
				Optional:    !contains(requiredAttributes, "service_rest_incoming_tls_enabled"),
			},
			"service_rest_incoming_tls_listen_port": schema.Int64Attribute{
				Description: "The port number for incoming REST clients that connect to the Message VPN over TLS. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.",
				Required:    contains(requiredAttributes, "service_rest_incoming_tls_listen_port"),
				Optional:    !contains(requiredAttributes, "service_rest_incoming_tls_listen_port"),
			},
			"service_rest_mode": schema.StringAttribute{
				Description: "The REST service mode for incoming REST clients that connect to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"messaging\"`. The allowed values and their meaning are:  <pre> \"gateway\" - Act as a message gateway through which REST messages are propagated. \"messaging\" - Act as a message broker on which REST messages are queued. </pre>  Available since 2.6.",
				Required:    contains(requiredAttributes, "service_rest_mode"),
				Optional:    !contains(requiredAttributes, "service_rest_mode"),

				Validators: []validator.String{
					stringvalidator.OneOf("gateway", "messaging"),
				},
				PlanModifiers: StringPlanModifiersFor("service_rest_mode", requiredAttributes),
			},
			"service_rest_outgoing_max_connection_count": schema.Int64Attribute{
				Description: "The maximum number of REST Consumer (outgoing) client connections that can be simultaneously connected to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.",
				Required:    contains(requiredAttributes, "service_rest_outgoing_max_connection_count"),
				Optional:    !contains(requiredAttributes, "service_rest_outgoing_max_connection_count"),
			},
			"service_smf_max_connection_count": schema.Int64Attribute{
				Description: "The maximum number of SMF client connections that can be simultaneously connected to the Message VPN. This value may be higher than supported by the platform. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.",
				Required:    contains(requiredAttributes, "service_smf_max_connection_count"),
				Optional:    !contains(requiredAttributes, "service_smf_max_connection_count"),
			},
			"service_smf_plain_text_enabled": schema.BoolAttribute{
				Description: "Enable or disable the plain-text SMF service in the Message VPN. Disabling causes clients currently connected to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "service_smf_plain_text_enabled"),
				Optional:    !contains(requiredAttributes, "service_smf_plain_text_enabled"),
			},
			"service_smf_tls_enabled": schema.BoolAttribute{
				Description: "Enable or disable the use of encryption (TLS) for the SMF service in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "service_smf_tls_enabled"),
				Optional:    !contains(requiredAttributes, "service_smf_tls_enabled"),
			},
			"service_web_authentication_client_cert_request": schema.StringAttribute{
				Description: "Determines when to request a client certificate from a Web Transport client connecting via a TLS port. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"when-enabled-in-message-vpn\"`. The allowed values and their meaning are:  <pre> \"always\" - Always ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"never\" - Never ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"when-enabled-in-message-vpn\" - Only ask for a client-certificate if client certificate authentication is enabled under \"message-vpn >  authentication > client-certificate > shutdown\". </pre>  Available since 2.21.",
				Required:    contains(requiredAttributes, "service_web_authentication_client_cert_request"),
				Optional:    !contains(requiredAttributes, "service_web_authentication_client_cert_request"),

				Validators: []validator.String{
					stringvalidator.OneOf("always", "never", "when-enabled-in-message-vpn"),
				},
				PlanModifiers: StringPlanModifiersFor("service_web_authentication_client_cert_request", requiredAttributes),
			},
			"service_web_max_connection_count": schema.Int64Attribute{
				Description: "The maximum number of Web Transport client connections that can be simultaneously connected to the Message VPN. This value may be higher than supported by the platform. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform.",
				Required:    contains(requiredAttributes, "service_web_max_connection_count"),
				Optional:    !contains(requiredAttributes, "service_web_max_connection_count"),
			},
			"service_web_plain_text_enabled": schema.BoolAttribute{
				Description: "Enable or disable the plain-text Web Transport service in the Message VPN. Disabling causes clients currently connected to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "service_web_plain_text_enabled"),
				Optional:    !contains(requiredAttributes, "service_web_plain_text_enabled"),
			},
			"service_web_tls_enabled": schema.BoolAttribute{
				Description: "Enable or disable the use of TLS for the Web Transport service in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "service_web_tls_enabled"),
				Optional:    !contains(requiredAttributes, "service_web_tls_enabled"),
			},
			"tls_allow_downgrade_to_plain_text_enabled": schema.BoolAttribute{
				Description: "Enable or disable the allowing of TLS SMF clients to downgrade their connections to plain-text connections. Changing this will not affect existing connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "tls_allow_downgrade_to_plain_text_enabled"),
				Optional:    !contains(requiredAttributes, "tls_allow_downgrade_to_plain_text_enabled"),
			},
		},
	}

	return schema
}
