/*
SEMP (Solace Element Management Protocol)

SEMP (starting in `v2`, see note 1) is a RESTful API for configuring, monitoring, and administering a Solace PubSub+ broker.  SEMP uses URIs to address manageable **resources** of the Solace PubSub+ broker. Resources are individual **objects**, **collections** of objects, or (exclusively in the action API) **actions**. This document applies to the following API:   API|Base Path|Purpose|Comments :---|:---|:---|:--- Configuration|/SEMP/v2/config|Reading and writing config state|See note 2    The following APIs are also available:   API|Base Path|Purpose|Comments :---|:---|:---|:--- Action|/SEMP/v2/action|Performing actions|See note 2 Monitoring|/SEMP/v2/monitor|Querying operational parameters|See note 2    Resources are always nouns, with individual objects being singular and collections being plural.  Objects within a collection are identified by an `obj-id`, which follows the collection name with the form `collection-name/obj-id`.  Actions within an object are identified by an `action-id`, which follows the object name with the form `obj-id/action-id`.  Some examples:  ``` /SEMP/v2/config/msgVpns                        ; MsgVpn collection /SEMP/v2/config/msgVpns/a                      ; MsgVpn object named \"a\" /SEMP/v2/config/msgVpns/a/queues               ; Queue collection in MsgVpn \"a\" /SEMP/v2/config/msgVpns/a/queues/b             ; Queue object named \"b\" in MsgVpn \"a\" /SEMP/v2/action/msgVpns/a/queues/b/startReplay ; Action that starts a replay on Queue \"b\" in MsgVpn \"a\" /SEMP/v2/monitor/msgVpns/a/clients             ; Client collection in MsgVpn \"a\" /SEMP/v2/monitor/msgVpns/a/clients/c           ; Client object named \"c\" in MsgVpn \"a\" ```  ## Collection Resources  Collections are unordered lists of objects (unless described as otherwise), and are described by JSON arrays. Each item in the array represents an object in the same manner as the individual object would normally be represented. In the configuration API, the creation of a new object is done through its collection resource.  ## Object and Action Resources  Objects are composed of attributes, actions, collections, and other objects. They are described by JSON objects as name/value pairs. The collections and actions of an object are not contained directly in the object's JSON content; rather the content includes an attribute containing a URI which points to the collections and actions. These contained resources must be managed through this URI. At a minimum, every object has one or more identifying attributes, and its own `uri` attribute which contains the URI pointing to itself.  Actions are also composed of attributes, and are described by JSON objects as name/value pairs. Unlike objects, however, they are not members of a collection and cannot be retrieved, only performed. Actions only exist in the action API.  Attributes in an object or action may have any combination of the following properties:   Property|Meaning|Comments :---|:---|:--- Identifying|Attribute is involved in unique identification of the object, and appears in its URI| Required|Attribute must be provided in the request| Read-Only|Attribute can only be read, not written.|See note 3 Write-Only|Attribute can only be written, not read, unless the attribute is also opaque|See the documentation for the opaque property Requires-Disable|Attribute can only be changed when object is disabled| Deprecated|Attribute is deprecated, and will disappear in the next SEMP version| Opaque|Attribute can be set or retrieved in opaque form when the `opaquePassword` query parameter is present|See the `opaquePassword` query parameter documentation    In some requests, certain attributes may only be provided in certain combinations with other attributes:   Relationship|Meaning :---|:--- Requires|Attribute may only be changed by a request if a particular attribute or combination of attributes is also provided in the request Conflicts|Attribute may only be provided in a request if a particular attribute or combination of attributes is not also provided in the request    In the monitoring API, any non-identifying attribute may not be returned in a GET.  ## HTTP Methods  The following HTTP methods manipulate resources in accordance with these general principles. Note that some methods are only used in certain APIs:   Method|Resource|Meaning|Request Body|Response Body|Missing Request Attributes :---|:---|:---|:---|:---|:--- POST|Collection|Create object|Initial attribute values|Object attributes and metadata|Set to default PUT|Object|Create or replace object (see note 5)|New attribute values|Object attributes and metadata|Set to default, with certain exceptions (see note 4) PUT|Action|Performs action|Action arguments|Action metadata|N/A PATCH|Object|Update object|New attribute values|Object attributes and metadata|unchanged DELETE|Object|Delete object|Empty|Object metadata|N/A GET|Object|Get object|Empty|Object attributes and metadata|N/A GET|Collection|Get collection|Empty|Object attributes and collection metadata|N/A    ## Common Query Parameters  The following are some common query parameters that are supported by many method/URI combinations. Individual URIs may document additional parameters. Note that multiple query parameters can be used together in a single URI, separated by the ampersand character. For example:  ``` ; Request for the MsgVpns collection using two hypothetical query parameters ; \"q1\" and \"q2\" with values \"val1\" and \"val2\" respectively /SEMP/v2/config/msgVpns?q1=val1&q2=val2 ```  ### select  Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. Use this query parameter to limit the size of the returned data for each returned object, return only those fields that are desired, or exclude fields that are not desired.  The value of `select` is a comma-separated list of attribute names. If the list contains attribute names that are not prefaced by `-`, only those attributes are included in the response. If the list contains attribute names that are prefaced by `-`, those attributes are excluded from the response. If the list contains both types, then the difference of the first set of attributes and the second set of attributes is returned. If the list is empty (i.e. `select=`), no attributes are returned.  All attributes that are prefaced by `-` must follow all attributes that are not prefaced by `-`. In addition, each attribute name in the list must match at least one attribute in the object.  Names may include the `*` wildcard (zero or more characters). Nested attribute names are supported using periods (e.g. `parentName.childName`).  Some examples:  ``` ; List of all MsgVpn names /SEMP/v2/config/msgVpns?select=msgVpnName ; List of all MsgVpn and their attributes except for their names /SEMP/v2/config/msgVpns?select=-msgVpnName ; Authentication attributes of MsgVpn \"finance\" /SEMP/v2/config/msgVpns/finance?select=authentication* ; All attributes of MsgVpn \"finance\" except for authentication attributes /SEMP/v2/config/msgVpns/finance?select=-authentication* ; Access related attributes of Queue \"orderQ\" of MsgVpn \"finance\" /SEMP/v2/config/msgVpns/finance/queues/orderQ?select=owner,permission ```  ### where  Include in the response only objects where certain conditions are true. Use this query parameter to limit which objects are returned to those whose attribute values meet the given conditions.  The value of `where` is a comma-separated list of expressions. All expressions must be true for the object to be included in the response. Each expression takes the form:  ``` expression  = attribute-name OP value OP          = '==' | '!=' | '&lt;' | '&gt;' | '&lt;=' | '&gt;=' ```  `value` may be a number, string, `true`, or `false`, as appropriate for the type of `attribute-name`. Greater-than and less-than comparisons only work for numbers. A `*` in a string `value` is interpreted as a wildcard (zero or more characters). Some examples:  ``` ; Only enabled MsgVpns /SEMP/v2/config/msgVpns?where=enabled==true ; Only MsgVpns using basic non-LDAP authentication /SEMP/v2/config/msgVpns?where=authenticationBasicEnabled==true,authenticationBasicType!=ldap ; Only MsgVpns that allow more than 100 client connections /SEMP/v2/config/msgVpns?where=maxConnectionCount>100 ; Only MsgVpns with msgVpnName starting with \"B\": /SEMP/v2/config/msgVpns?where=msgVpnName==B* ```  ### count  Limit the count of objects in the response. This can be useful to limit the size of the response for large collections. The minimum value for `count` is `1` and the default is `10`. There is also a per-collection maximum value to limit request handling time.  `count` does not guarantee that a minimum number of objects will be returned. A page may contain fewer than `count` objects or even be empty. Additional objects may nonetheless be available for retrieval on subsequent pages. See the `cursor` query parameter documentation for more information on paging.  For example: ``` ; Up to 25 MsgVpns /SEMP/v2/config/msgVpns?count=25 ```  ### cursor  The cursor, or position, for the next page of objects. Cursors are opaque data that should not be created or interpreted by SEMP clients, and should only be used as described below.  When a request is made for a collection and there may be additional objects available for retrieval that are not included in the initial response, the response will include a `cursorQuery` field containing a cursor. The value of this field can be specified in the `cursor` query parameter of a subsequent request to retrieve the next page of objects. For convenience, an appropriate URI is constructed automatically by the broker and included in the `nextPageUri` field of the response. This URI can be used directly to retrieve the next page of objects.  Applications must continue to follow the `nextPageUri` if one is provided in order to retrieve the full set of objects associated with the request, even if a page contains fewer than the requested number of objects (see the `count` query parameter documentation) or is empty.  ### opaquePassword  Attributes with the opaque property are also write-only and so cannot normally be retrieved in a GET. However, when a password is provided in the `opaquePassword` query parameter, attributes with the opaque property are retrieved in a GET in opaque form, encrypted with this password. The query parameter can also be used on a POST, PATCH, or PUT to set opaque attributes using opaque attribute values retrieved in a GET, so long as:  1. the same password that was used to retrieve the opaque attribute values is provided; and  2. the broker to which the request is being sent has the same major and minor SEMP version as the broker that produced the opaque attribute values.  The password provided in the query parameter must be a minimum of 8 characters and a maximum of 128 characters.  The query parameter can only be used in the configuration API, and only over HTTPS.  ## Authentication  When a client makes its first SEMPv2 request, it must supply a username and password using HTTP Basic authentication, or an OAuth token or tokens using HTTP Bearer authentication.  When HTTP Basic authentication is used, the broker returns a cookie containing a session key. The client can omit the username and password from subsequent requests, because the broker can use the session cookie for authentication instead. When the session expires or is deleted, the client must provide the username and password again, and the broker creates a new session.  There are a limited number of session slots available on the broker. The broker returns 529 No SEMP Session Available if it is not able to allocate a session.  If certain attributes—such as a user's password—are changed, the broker automatically deletes the affected sessions. These attributes are documented below. However, changes in external user configuration data stored on a RADIUS or LDAP server do not trigger the broker to delete the associated session(s), therefore you must do this manually, if required.  A client can retrieve its current session information using the /about/user endpoint and delete its own session using the /about/user/logout endpoint. A client with appropriate permissions can also manage all sessions using the /sessions endpoint.  Sessions are not created when authenticating with an OAuth token or tokens using HTTP Bearer authentication. If a session cookie is provided, it is ignored.  ## Help  Visit [our website](https://solace.com) to learn more about Solace.  You can also download the SEMP API specifications by clicking [here](https://solace.com/downloads/).  If you need additional support, please contact us at [support@solace.com](mailto:support@solace.com).  ## Notes  Note|Description :---:|:--- 1|This specification defines SEMP starting in \"v2\", and not the original SEMP \"v1\" interface. Request and response formats between \"v1\" and \"v2\" are entirely incompatible, although both protocols share a common port configuration on the Solace PubSub+ broker. They are differentiated by the initial portion of the URI path, one of either \"/SEMP/\" or \"/SEMP/v2/\" 2|This API is partially implemented. Only a subset of all objects are available. 3|Read-only attributes may appear in POST and PUT/PATCH requests. However, if a read-only attribute is not marked as identifying, it will be ignored during a PUT/PATCH. 4|On a PUT, if the SEMP user is not authorized to modify the attribute, its value is left unchanged rather than set to default. In addition, the values of write-only attributes are not set to their defaults on a PUT, except in the following two cases: there is a mutual requires relationship with another non-write-only attribute, both attributes are absent from the request, and the non-write-only attribute is not currently set to its default value; or the attribute is also opaque and the `opaquePassword` query parameter is provided in the request. 5|On a PUT, if the object does not exist, it is created first.  

API version: 2.26
Contact: support@solace.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package sempv2

import (
	"encoding/json"
)

// Broker struct for Broker
type Broker struct {
	// The client certificate revocation checking mode used when a client authenticates with a client certificate. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"none\"`. The allowed values and their meaning are:  <pre> \"none\" - Do not perform any certificate revocation checking. \"ocsp\" - Use the Open Certificate Status Protcol (OCSP) for certificate revocation checking. \"crl\" - Use Certificate Revocation Lists (CRL) for certificate revocation checking. \"ocsp-crl\" - Use OCSP first, but if OCSP fails to return an unambiguous result, then check via CRL. </pre> 
	AuthClientCertRevocationCheckMode *string `json:"authClientCertRevocationCheckMode,omitempty"`
	// The maximum depth for a client certificate chain. The depth of a chain is defined as the number of signing CA certificates that are present in the chain back to a trusted self-signed root CA certificate. The default value is `3`. Available since 2.22.
	ConfigSyncAuthenticationClientCertMaxChainDepth *int64 `json:"configSyncAuthenticationClientCertMaxChainDepth,omitempty"`
	// Enable or disable validation of the \"Not Before\" and \"Not After\" validity dates in the authentication certificate(s). The default value is `true`. Available since 2.22.
	ConfigSyncAuthenticationClientCertValidateDateEnabled *bool `json:"configSyncAuthenticationClientCertValidateDateEnabled,omitempty"`
	// The TCP initial congestion window size for Config Sync clients, in multiples of the TCP Maximum Segment Size (MSS). Changing the value from its default of 2 results in non-compliance with RFC 2581. Contact support before changing this value. The default value is `2`. Available since 2.22.
	ConfigSyncClientProfileTcpInitialCongestionWindow *int64 `json:"configSyncClientProfileTcpInitialCongestionWindow,omitempty"`
	// The number of TCP keepalive retransmissions to a client using the Client Profile before declaring that it is not available. The default value is `5`. Available since 2.22.
	ConfigSyncClientProfileTcpKeepaliveCount *int64 `json:"configSyncClientProfileTcpKeepaliveCount,omitempty"`
	// The amount of time a client connection using the Client Profile must remain idle before TCP begins sending keepalive probes, in seconds. The default value is `3`. Available since 2.22.
	ConfigSyncClientProfileTcpKeepaliveIdle *int64 `json:"configSyncClientProfileTcpKeepaliveIdle,omitempty"`
	// The amount of time between TCP keepalive retransmissions to a client using the Client Profile when no acknowledgement is received, in seconds. The default value is `1`. Available since 2.22.
	ConfigSyncClientProfileTcpKeepaliveInterval *int64 `json:"configSyncClientProfileTcpKeepaliveInterval,omitempty"`
	// The TCP maximum window size for clients using the Client Profile, in kilobytes. Changes are applied to all existing connections. The default value is `256`. Available since 2.22.
	ConfigSyncClientProfileTcpMaxWindow *int64 `json:"configSyncClientProfileTcpMaxWindow,omitempty"`
	// The TCP maximum segment size for clients using the Client Profile, in bytes. Changes are applied to all existing connections. The default value is `1460`. Available since 2.22.
	ConfigSyncClientProfileTcpMss *int64 `json:"configSyncClientProfileTcpMss,omitempty"`
	// Enable or disable configuration synchronization for High Availability or Disaster Recovery. The default value is `false`. Available since 2.22.
	ConfigSyncEnabled *bool `json:"configSyncEnabled,omitempty"`
	// Enable or disable the synchronizing of usernames within High Availability groups. The transition from not synchronizing to synchronizing will cause the High Availability mate to fall out of sync. Recommendation: leave this as enabled. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`. Available since 2.22.
	ConfigSyncSynchronizeUsernameEnabled *bool `json:"configSyncSynchronizeUsernameEnabled,omitempty"`
	// Enable or disable the use of TLS encryption of the configuration synchronization communications between brokers in High Availability groups and/or Disaster Recovery sites. The default value is `false`. Available since 2.22.
	ConfigSyncTlsEnabled *bool `json:"configSyncTlsEnabled,omitempty"`
	// The days of the week to schedule defragmentation runs, specified as \"daily\" or as a comma-separated list of days. Days must be specified as \"Sun\", \"Mon\", \"Tue\", \"Wed\", \"Thu\", \"Fri, or \"Sat\", with no spaces, and in sorted order from Sunday to Saturday. Please note \"Sun,Mon,Tue,Wed,Thu,Fri,Sat\" is not allowed, use \"daily\" instead. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"daily\"`. Available since 2.25.
	GuaranteedMsgingDefragmentationScheduleDayList *string `json:"guaranteedMsgingDefragmentationScheduleDayList,omitempty"`
	// Enable or disable schedule-based defragmentation of Guaranteed Messaging spool files. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`. Available since 2.25.
	GuaranteedMsgingDefragmentationScheduleEnabled *bool `json:"guaranteedMsgingDefragmentationScheduleEnabled,omitempty"`
	// The times of the day to schedule defragmentation runs, specified as \"hourly\" or as a comma-separated list of 24-hour times in the form hh:mm, or h:mm. There must be no spaces, and times (up to 4) must be in sorted order from 0:00 to 23:59. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"0:00\"`. Available since 2.25.
	GuaranteedMsgingDefragmentationScheduleTimeList *string `json:"guaranteedMsgingDefragmentationScheduleTimeList,omitempty"`
	// Enable or disable threshold-based defragmentation of Guaranteed Messaging spool files. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`. Available since 2.25.
	GuaranteedMsgingDefragmentationThresholdEnabled *bool `json:"guaranteedMsgingDefragmentationThresholdEnabled,omitempty"`
	// Percentage of spool fragmentation needed to trigger defragmentation run. The minimum value allowed is 30%. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `50`. Available since 2.25.
	GuaranteedMsgingDefragmentationThresholdFragmentationPercentage *int64 `json:"guaranteedMsgingDefragmentationThresholdFragmentationPercentage,omitempty"`
	// Minimum interval of time (in minutes) between defragmentation runs triggered by thresholds. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `15`. Available since 2.25.
	GuaranteedMsgingDefragmentationThresholdMinInterval *int64 `json:"guaranteedMsgingDefragmentationThresholdMinInterval,omitempty"`
	// Percentage of spool usage needed to trigger defragmentation run. The minimum value allowed is 30%. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `50`. Available since 2.25.
	GuaranteedMsgingDefragmentationThresholdUsagePercentage *int64 `json:"guaranteedMsgingDefragmentationThresholdUsagePercentage,omitempty"`
	// Enable or disable Guaranteed Messaging. The default value is `false`. Available since 2.18.
	GuaranteedMsgingEnabled *bool `json:"guaranteedMsgingEnabled,omitempty"`
	GuaranteedMsgingEventCacheUsageThreshold *EventThreshold `json:"guaranteedMsgingEventCacheUsageThreshold,omitempty"`
	GuaranteedMsgingEventDeliveredUnackedThreshold *EventThresholdByPercent `json:"guaranteedMsgingEventDeliveredUnackedThreshold,omitempty"`
	GuaranteedMsgingEventDiskUsageThreshold *EventThresholdByPercent `json:"guaranteedMsgingEventDiskUsageThreshold,omitempty"`
	GuaranteedMsgingEventEgressFlowCountThreshold *EventThreshold `json:"guaranteedMsgingEventEgressFlowCountThreshold,omitempty"`
	GuaranteedMsgingEventEndpointCountThreshold *EventThreshold `json:"guaranteedMsgingEventEndpointCountThreshold,omitempty"`
	GuaranteedMsgingEventIngressFlowCountThreshold *EventThreshold `json:"guaranteedMsgingEventIngressFlowCountThreshold,omitempty"`
	GuaranteedMsgingEventMsgCountThreshold *EventThresholdByPercent `json:"guaranteedMsgingEventMsgCountThreshold,omitempty"`
	GuaranteedMsgingEventMsgSpoolFileCountThreshold *EventThresholdByPercent `json:"guaranteedMsgingEventMsgSpoolFileCountThreshold,omitempty"`
	GuaranteedMsgingEventMsgSpoolUsageThreshold *EventThreshold `json:"guaranteedMsgingEventMsgSpoolUsageThreshold,omitempty"`
	GuaranteedMsgingEventTransactedSessionCountThreshold *EventThreshold `json:"guaranteedMsgingEventTransactedSessionCountThreshold,omitempty"`
	GuaranteedMsgingEventTransactedSessionResourceCountThreshold *EventThresholdByPercent `json:"guaranteedMsgingEventTransactedSessionResourceCountThreshold,omitempty"`
	GuaranteedMsgingEventTransactionCountThreshold *EventThreshold `json:"guaranteedMsgingEventTransactionCountThreshold,omitempty"`
	// Guaranteed messaging cache usage limit. Expressed as a maximum percentage of the NAB's egress queueing. resources that the guaranteed message cache is allowed to use. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `10`. Available since 2.18.
	GuaranteedMsgingMaxCacheUsage *int32 `json:"guaranteedMsgingMaxCacheUsage,omitempty"`
	// The maximum total message spool usage allowed across all VPNs on this broker, in megabytes. Recommendation: the maximum value should be less than 90% of the disk space allocated for the guaranteed message spool. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `1500`. Available since 2.18.
	GuaranteedMsgingMaxMsgSpoolUsage *int64 `json:"guaranteedMsgingMaxMsgSpoolUsage,omitempty"`
	// The maximum time, in milliseconds, that can be tolerated for remote acknowledgement of synchronization messages before which the remote system will be considered out of sync. The default value is `10000`. Available since 2.18.
	GuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout *int64 `json:"guaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout,omitempty"`
	// The maximum time, in milliseconds, that can be tolerated for remote disk writes before which the remote system will be considered out of sync. The default value is `10000`. Available since 2.18.
	GuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout *int64 `json:"guaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout,omitempty"`
	// The replication compatibility mode for the router. The default value is `\"legacy\"`. The allowed values and their meaning are:\"legacy\" - All transactions originated by clients are replicated to the standby site without using transactions.\"transacted\" - All transactions originated by clients are replicated to the standby site using transactions. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"legacy\"`. The allowed values and their meaning are:  <pre> \"legacy\" - All transactions originated by clients are replicated to the standby site without using transactions. \"transacted\" - All transactions originated by clients are replicated to the standby site using transactions. </pre>  Available since 2.18.
	GuaranteedMsgingTransactionReplicationCompatibilityMode *string `json:"guaranteedMsgingTransactionReplicationCompatibilityMode,omitempty"`
	// The default OAuth profile for OAuth authenticated SEMP requests. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`. Available since 2.24.
	OauthProfileDefault *string `json:"oauthProfileDefault,omitempty"`
	// Enable or disable the AMQP service. When disabled new AMQP Clients may not connect through the global or per-VPN AMQP listen-ports, and all currently connected AMQP Clients are immediately disconnected. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`. Available since 2.17.
	ServiceAmqpEnabled *bool `json:"serviceAmqpEnabled,omitempty"`
	// TCP port number that AMQP clients can use to connect to the broker using raw TCP over TLS. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `0`. Available since 2.17.
	ServiceAmqpTlsListenPort *int64 `json:"serviceAmqpTlsListenPort,omitempty"`
	ServiceEventConnectionCountThreshold *EventThreshold `json:"serviceEventConnectionCountThreshold,omitempty"`
	// Enable or disable the health-check service. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`. Available since 2.17.
	ServiceHealthCheckEnabled *bool `json:"serviceHealthCheckEnabled,omitempty"`
	// The port number for the health-check service. The port must be unique across the message backbone. The health-check service must be disabled to change the port. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `5550`. Available since 2.17.
	ServiceHealthCheckListenPort *int64 `json:"serviceHealthCheckListenPort,omitempty"`
	// Enable or disable the mate-link service. The default value is `true`. Available since 2.17.
	ServiceMateLinkEnabled *bool `json:"serviceMateLinkEnabled,omitempty"`
	// The port number for the mate-link service. The port must be unique across the message backbone. The mate-link service must be disabled to change the port. The default value is `8741`. Available since 2.17.
	ServiceMateLinkListenPort *int64 `json:"serviceMateLinkListenPort,omitempty"`
	// Enable or disable the MQTT service. When disabled new MQTT Clients may not connect through the per-VPN MQTT listen-ports, and all currently connected MQTT Clients are immediately disconnected. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`. Available since 2.17.
	ServiceMqttEnabled *bool `json:"serviceMqttEnabled,omitempty"`
	// Enable or disable the msg-backbone service. When disabled new Clients may not connect through global or per-VPN listen-ports, and all currently connected Clients are immediately disconnected. The default value is `true`. Available since 2.17.
	ServiceMsgBackboneEnabled *bool `json:"serviceMsgBackboneEnabled,omitempty"`
	// Enable or disable the redundancy service. The default value is `true`. Available since 2.17.
	ServiceRedundancyEnabled *bool `json:"serviceRedundancyEnabled,omitempty"`
	// The first listen-port used for the redundancy service. Redundancy uses this port and the subsequent 2 ports. These port must be unique across the message backbone. The redundancy service must be disabled to change this port. The default value is `8300`. Available since 2.17.
	ServiceRedundancyFirstListenPort *int64 `json:"serviceRedundancyFirstListenPort,omitempty"`
	ServiceRestEventOutgoingConnectionCountThreshold *EventThreshold `json:"serviceRestEventOutgoingConnectionCountThreshold,omitempty"`
	// Enable or disable the REST service incoming connections on the router. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`. Available since 2.17.
	ServiceRestIncomingEnabled *bool `json:"serviceRestIncomingEnabled,omitempty"`
	// Enable or disable the REST service outgoing connections on the router. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`. Available since 2.17.
	ServiceRestOutgoingEnabled *bool `json:"serviceRestOutgoingEnabled,omitempty"`
	// Enable or disable cross origin resource requests for the SEMP service. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`. Available since 2.24.
	ServiceSempCorsAllowAnyHostEnabled *bool `json:"serviceSempCorsAllowAnyHostEnabled,omitempty"`
	// Enable or disable extended SEMP timeouts for paged GETs. When a request times out, it returns the current page of content, even if the page is not full.  When enabled, the timeout is 60 seconds. When disabled, the timeout is 5 seconds.  The recommended setting is disabled (no legacy-timeout).  This parameter is intended as a temporary workaround to be used until SEMP clients can handle short pages.  This setting will be removed in a future release. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`. Available since 2.18.
	ServiceSempLegacyTimeoutEnabled *bool `json:"serviceSempLegacyTimeoutEnabled,omitempty"`
	// Enable or disable plain-text SEMP service. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`. Available since 2.17.
	ServiceSempPlainTextEnabled *bool `json:"serviceSempPlainTextEnabled,omitempty"`
	// The TCP port for plain-text SEMP client connections. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `80`. Available since 2.17.
	ServiceSempPlainTextListenPort *int64 `json:"serviceSempPlainTextListenPort,omitempty"`
	// The session idle timeout, in minutes. Sessions will be invalidated if there is no activity in this period of time. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `15`. Available since 2.21.
	ServiceSempSessionIdleTimeout *int32 `json:"serviceSempSessionIdleTimeout,omitempty"`
	// The maximum lifetime of a session, in minutes. Sessions will be invalidated after this period of time, regardless of activity. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `43200`. Available since 2.21.
	ServiceSempSessionMaxLifetime *int32 `json:"serviceSempSessionMaxLifetime,omitempty"`
	// Enable or disable TLS SEMP service. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`. Available since 2.17.
	ServiceSempTlsEnabled *bool `json:"serviceSempTlsEnabled,omitempty"`
	// The TCP port for TLS SEMP client connections. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `1943`. Available since 2.17.
	ServiceSempTlsListenPort *int64 `json:"serviceSempTlsListenPort,omitempty"`
	// TCP port number that SMF clients can use to connect to the broker using raw compression TCP. The default value is `55003`. Available since 2.17.
	ServiceSmfCompressionListenPort *int64 `json:"serviceSmfCompressionListenPort,omitempty"`
	// Enable or disable the SMF service. When disabled new SMF Clients may not connect through the global listen-ports, and all currently connected SMF Clients are immediately disconnected. The default value is `true`. Available since 2.17.
	ServiceSmfEnabled *bool `json:"serviceSmfEnabled,omitempty"`
	ServiceSmfEventConnectionCountThreshold *EventThreshold `json:"serviceSmfEventConnectionCountThreshold,omitempty"`
	// TCP port number that SMF clients can use to connect to the broker using raw TCP. The default value is `55555`. Available since 2.17.
	ServiceSmfPlainTextListenPort *int64 `json:"serviceSmfPlainTextListenPort,omitempty"`
	// TCP port number that SMF clients can use to connect to the broker using raw routing control TCP. The default value is `55556`. Available since 2.17.
	ServiceSmfRoutingControlListenPort *int64 `json:"serviceSmfRoutingControlListenPort,omitempty"`
	// TCP port number that SMF clients can use to connect to the broker using raw TCP over TLS. The default value is `55443`. Available since 2.17.
	ServiceSmfTlsListenPort *int64 `json:"serviceSmfTlsListenPort,omitempty"`
	ServiceTlsEventConnectionCountThreshold *EventThreshold `json:"serviceTlsEventConnectionCountThreshold,omitempty"`
	// Enable or disable the web-transport service. When disabled new web-transport Clients may not connect through the global listen-ports, and all currently connected web-transport Clients are immediately disconnected. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`. Available since 2.17.
	ServiceWebTransportEnabled *bool `json:"serviceWebTransportEnabled,omitempty"`
	// The TCP port for plain-text WEB client connections. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `8008`. Available since 2.17.
	ServiceWebTransportPlainTextListenPort *int64 `json:"serviceWebTransportPlainTextListenPort,omitempty"`
	// The TCP port for TLS WEB client connections. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `1443`. Available since 2.17.
	ServiceWebTransportTlsListenPort *int64 `json:"serviceWebTransportTlsListenPort,omitempty"`
	// Used to specify the Web URL suffix that will be used by Web clients when communicating with the broker. The default value is `\"\"`. Available since 2.17.
	ServiceWebTransportWebUrlSuffix *string `json:"serviceWebTransportWebUrlSuffix,omitempty"`
	// Enable or disable the blocking of TLS version 1.1 connections. When blocked, all existing incoming and outgoing TLS 1.1 connections with Clients, SEMP users, and LDAP servers remain connected while new connections are blocked. Note that support for TLS 1.1 will eventually be discontinued, at which time TLS 1.1 connections will be blocked regardless of this setting. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.
	TlsBlockVersion11Enabled *bool `json:"tlsBlockVersion11Enabled,omitempty"`
	// The colon-separated list of cipher suites used for TLS management connections (e.g. SEMP, LDAP). The value \"default\" implies all supported suites ordered from most secure to least secure. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"default\"`.
	TlsCipherSuiteManagementList *string `json:"tlsCipherSuiteManagementList,omitempty"`
	// The colon-separated list of cipher suites used for TLS data connections (e.g. client pub/sub). The value \"default\" implies all supported suites ordered from most secure to least secure. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"default\"`.
	TlsCipherSuiteMsgBackboneList *string `json:"tlsCipherSuiteMsgBackboneList,omitempty"`
	// The colon-separated list of cipher suites used for TLS secure shell connections (e.g. SSH, SFTP, SCP). The value \"default\" implies all supported suites ordered from most secure to least secure. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"default\"`.
	TlsCipherSuiteSecureShellList *string `json:"tlsCipherSuiteSecureShellList,omitempty"`
	// Enable or disable protection against the CRIME exploit. When enabled, TLS+compressed messaging performance is degraded. This protection should only be disabled if sufficient ACL and authentication features are being employed such that a potential attacker does not have sufficient access to trigger the exploit. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`.
	TlsCrimeExploitProtectionEnabled *bool `json:"tlsCrimeExploitProtectionEnabled,omitempty"`
	// The PEM formatted content for the server certificate used for TLS connections. It must consist of a private key and between one and three certificates comprising the certificate trust chain. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changing this attribute requires an HTTPS connection. The default value is `\"\"`.
	TlsServerCertContent *string `json:"tlsServerCertContent,omitempty"`
	// The password for the server certificate used for TLS connections. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changing this attribute requires an HTTPS connection. The default value is `\"\"`.
	TlsServerCertPassword *string `json:"tlsServerCertPassword,omitempty"`
	// Enable or disable the standard domain certificate authority list. The default value is `true`. Available since 2.19.
	TlsStandardDomainCertificateAuthoritiesEnabled *bool `json:"tlsStandardDomainCertificateAuthoritiesEnabled,omitempty"`
	// The TLS ticket lifetime in seconds. When a client connects with TLS, a session with a session ticket is created using the TLS ticket lifetime which determines how long the client has to resume the session. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `86400`.
	TlsTicketLifetime *int32 `json:"tlsTicketLifetime,omitempty"`
	// Reserved for internal use by Solace. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`. Available since 2.25.
	WebManagerCustomization *string `json:"webManagerCustomization,omitempty"`
	// Enable or disable redirection of HTTP requests for the broker manager to HTTPS. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`. Available since 2.24.
	WebManagerRedirectHttpEnabled *bool `json:"webManagerRedirectHttpEnabled,omitempty"`
	// The HTTPS port that HTTP requests will be redirected towards in a HTTP 301 redirect response. Zero is a special value that means use the value specified for the SEMP TLS port value. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `0`. Available since 2.24.
	WebManagerRedirectHttpOverrideTlsPort *int32 `json:"webManagerRedirectHttpOverrideTlsPort,omitempty"`
}

// NewBroker instantiates a new Broker object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBroker() *Broker {
	this := Broker{}
	return &this
}

// NewBrokerWithDefaults instantiates a new Broker object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBrokerWithDefaults() *Broker {
	this := Broker{}
	return &this
}

// GetAuthClientCertRevocationCheckMode returns the AuthClientCertRevocationCheckMode field value if set, zero value otherwise.
func (o *Broker) GetAuthClientCertRevocationCheckMode() string {
	if o == nil || o.AuthClientCertRevocationCheckMode == nil {
		var ret string
		return ret
	}
	return *o.AuthClientCertRevocationCheckMode
}

// GetAuthClientCertRevocationCheckModeOk returns a tuple with the AuthClientCertRevocationCheckMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetAuthClientCertRevocationCheckModeOk() (*string, bool) {
	if o == nil || o.AuthClientCertRevocationCheckMode == nil {
		return nil, false
	}
	return o.AuthClientCertRevocationCheckMode, true
}

// HasAuthClientCertRevocationCheckMode returns a boolean if a field has been set.
func (o *Broker) HasAuthClientCertRevocationCheckMode() bool {
	if o != nil && o.AuthClientCertRevocationCheckMode != nil {
		return true
	}

	return false
}

// SetAuthClientCertRevocationCheckMode gets a reference to the given string and assigns it to the AuthClientCertRevocationCheckMode field.
func (o *Broker) SetAuthClientCertRevocationCheckMode(v string) {
	o.AuthClientCertRevocationCheckMode = &v
}

// GetConfigSyncAuthenticationClientCertMaxChainDepth returns the ConfigSyncAuthenticationClientCertMaxChainDepth field value if set, zero value otherwise.
func (o *Broker) GetConfigSyncAuthenticationClientCertMaxChainDepth() int64 {
	if o == nil || o.ConfigSyncAuthenticationClientCertMaxChainDepth == nil {
		var ret int64
		return ret
	}
	return *o.ConfigSyncAuthenticationClientCertMaxChainDepth
}

// GetConfigSyncAuthenticationClientCertMaxChainDepthOk returns a tuple with the ConfigSyncAuthenticationClientCertMaxChainDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetConfigSyncAuthenticationClientCertMaxChainDepthOk() (*int64, bool) {
	if o == nil || o.ConfigSyncAuthenticationClientCertMaxChainDepth == nil {
		return nil, false
	}
	return o.ConfigSyncAuthenticationClientCertMaxChainDepth, true
}

// HasConfigSyncAuthenticationClientCertMaxChainDepth returns a boolean if a field has been set.
func (o *Broker) HasConfigSyncAuthenticationClientCertMaxChainDepth() bool {
	if o != nil && o.ConfigSyncAuthenticationClientCertMaxChainDepth != nil {
		return true
	}

	return false
}

// SetConfigSyncAuthenticationClientCertMaxChainDepth gets a reference to the given int64 and assigns it to the ConfigSyncAuthenticationClientCertMaxChainDepth field.
func (o *Broker) SetConfigSyncAuthenticationClientCertMaxChainDepth(v int64) {
	o.ConfigSyncAuthenticationClientCertMaxChainDepth = &v
}

// GetConfigSyncAuthenticationClientCertValidateDateEnabled returns the ConfigSyncAuthenticationClientCertValidateDateEnabled field value if set, zero value otherwise.
func (o *Broker) GetConfigSyncAuthenticationClientCertValidateDateEnabled() bool {
	if o == nil || o.ConfigSyncAuthenticationClientCertValidateDateEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ConfigSyncAuthenticationClientCertValidateDateEnabled
}

// GetConfigSyncAuthenticationClientCertValidateDateEnabledOk returns a tuple with the ConfigSyncAuthenticationClientCertValidateDateEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetConfigSyncAuthenticationClientCertValidateDateEnabledOk() (*bool, bool) {
	if o == nil || o.ConfigSyncAuthenticationClientCertValidateDateEnabled == nil {
		return nil, false
	}
	return o.ConfigSyncAuthenticationClientCertValidateDateEnabled, true
}

// HasConfigSyncAuthenticationClientCertValidateDateEnabled returns a boolean if a field has been set.
func (o *Broker) HasConfigSyncAuthenticationClientCertValidateDateEnabled() bool {
	if o != nil && o.ConfigSyncAuthenticationClientCertValidateDateEnabled != nil {
		return true
	}

	return false
}

// SetConfigSyncAuthenticationClientCertValidateDateEnabled gets a reference to the given bool and assigns it to the ConfigSyncAuthenticationClientCertValidateDateEnabled field.
func (o *Broker) SetConfigSyncAuthenticationClientCertValidateDateEnabled(v bool) {
	o.ConfigSyncAuthenticationClientCertValidateDateEnabled = &v
}

// GetConfigSyncClientProfileTcpInitialCongestionWindow returns the ConfigSyncClientProfileTcpInitialCongestionWindow field value if set, zero value otherwise.
func (o *Broker) GetConfigSyncClientProfileTcpInitialCongestionWindow() int64 {
	if o == nil || o.ConfigSyncClientProfileTcpInitialCongestionWindow == nil {
		var ret int64
		return ret
	}
	return *o.ConfigSyncClientProfileTcpInitialCongestionWindow
}

// GetConfigSyncClientProfileTcpInitialCongestionWindowOk returns a tuple with the ConfigSyncClientProfileTcpInitialCongestionWindow field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetConfigSyncClientProfileTcpInitialCongestionWindowOk() (*int64, bool) {
	if o == nil || o.ConfigSyncClientProfileTcpInitialCongestionWindow == nil {
		return nil, false
	}
	return o.ConfigSyncClientProfileTcpInitialCongestionWindow, true
}

// HasConfigSyncClientProfileTcpInitialCongestionWindow returns a boolean if a field has been set.
func (o *Broker) HasConfigSyncClientProfileTcpInitialCongestionWindow() bool {
	if o != nil && o.ConfigSyncClientProfileTcpInitialCongestionWindow != nil {
		return true
	}

	return false
}

// SetConfigSyncClientProfileTcpInitialCongestionWindow gets a reference to the given int64 and assigns it to the ConfigSyncClientProfileTcpInitialCongestionWindow field.
func (o *Broker) SetConfigSyncClientProfileTcpInitialCongestionWindow(v int64) {
	o.ConfigSyncClientProfileTcpInitialCongestionWindow = &v
}

// GetConfigSyncClientProfileTcpKeepaliveCount returns the ConfigSyncClientProfileTcpKeepaliveCount field value if set, zero value otherwise.
func (o *Broker) GetConfigSyncClientProfileTcpKeepaliveCount() int64 {
	if o == nil || o.ConfigSyncClientProfileTcpKeepaliveCount == nil {
		var ret int64
		return ret
	}
	return *o.ConfigSyncClientProfileTcpKeepaliveCount
}

// GetConfigSyncClientProfileTcpKeepaliveCountOk returns a tuple with the ConfigSyncClientProfileTcpKeepaliveCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetConfigSyncClientProfileTcpKeepaliveCountOk() (*int64, bool) {
	if o == nil || o.ConfigSyncClientProfileTcpKeepaliveCount == nil {
		return nil, false
	}
	return o.ConfigSyncClientProfileTcpKeepaliveCount, true
}

// HasConfigSyncClientProfileTcpKeepaliveCount returns a boolean if a field has been set.
func (o *Broker) HasConfigSyncClientProfileTcpKeepaliveCount() bool {
	if o != nil && o.ConfigSyncClientProfileTcpKeepaliveCount != nil {
		return true
	}

	return false
}

// SetConfigSyncClientProfileTcpKeepaliveCount gets a reference to the given int64 and assigns it to the ConfigSyncClientProfileTcpKeepaliveCount field.
func (o *Broker) SetConfigSyncClientProfileTcpKeepaliveCount(v int64) {
	o.ConfigSyncClientProfileTcpKeepaliveCount = &v
}

// GetConfigSyncClientProfileTcpKeepaliveIdle returns the ConfigSyncClientProfileTcpKeepaliveIdle field value if set, zero value otherwise.
func (o *Broker) GetConfigSyncClientProfileTcpKeepaliveIdle() int64 {
	if o == nil || o.ConfigSyncClientProfileTcpKeepaliveIdle == nil {
		var ret int64
		return ret
	}
	return *o.ConfigSyncClientProfileTcpKeepaliveIdle
}

// GetConfigSyncClientProfileTcpKeepaliveIdleOk returns a tuple with the ConfigSyncClientProfileTcpKeepaliveIdle field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetConfigSyncClientProfileTcpKeepaliveIdleOk() (*int64, bool) {
	if o == nil || o.ConfigSyncClientProfileTcpKeepaliveIdle == nil {
		return nil, false
	}
	return o.ConfigSyncClientProfileTcpKeepaliveIdle, true
}

// HasConfigSyncClientProfileTcpKeepaliveIdle returns a boolean if a field has been set.
func (o *Broker) HasConfigSyncClientProfileTcpKeepaliveIdle() bool {
	if o != nil && o.ConfigSyncClientProfileTcpKeepaliveIdle != nil {
		return true
	}

	return false
}

// SetConfigSyncClientProfileTcpKeepaliveIdle gets a reference to the given int64 and assigns it to the ConfigSyncClientProfileTcpKeepaliveIdle field.
func (o *Broker) SetConfigSyncClientProfileTcpKeepaliveIdle(v int64) {
	o.ConfigSyncClientProfileTcpKeepaliveIdle = &v
}

// GetConfigSyncClientProfileTcpKeepaliveInterval returns the ConfigSyncClientProfileTcpKeepaliveInterval field value if set, zero value otherwise.
func (o *Broker) GetConfigSyncClientProfileTcpKeepaliveInterval() int64 {
	if o == nil || o.ConfigSyncClientProfileTcpKeepaliveInterval == nil {
		var ret int64
		return ret
	}
	return *o.ConfigSyncClientProfileTcpKeepaliveInterval
}

// GetConfigSyncClientProfileTcpKeepaliveIntervalOk returns a tuple with the ConfigSyncClientProfileTcpKeepaliveInterval field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetConfigSyncClientProfileTcpKeepaliveIntervalOk() (*int64, bool) {
	if o == nil || o.ConfigSyncClientProfileTcpKeepaliveInterval == nil {
		return nil, false
	}
	return o.ConfigSyncClientProfileTcpKeepaliveInterval, true
}

// HasConfigSyncClientProfileTcpKeepaliveInterval returns a boolean if a field has been set.
func (o *Broker) HasConfigSyncClientProfileTcpKeepaliveInterval() bool {
	if o != nil && o.ConfigSyncClientProfileTcpKeepaliveInterval != nil {
		return true
	}

	return false
}

// SetConfigSyncClientProfileTcpKeepaliveInterval gets a reference to the given int64 and assigns it to the ConfigSyncClientProfileTcpKeepaliveInterval field.
func (o *Broker) SetConfigSyncClientProfileTcpKeepaliveInterval(v int64) {
	o.ConfigSyncClientProfileTcpKeepaliveInterval = &v
}

// GetConfigSyncClientProfileTcpMaxWindow returns the ConfigSyncClientProfileTcpMaxWindow field value if set, zero value otherwise.
func (o *Broker) GetConfigSyncClientProfileTcpMaxWindow() int64 {
	if o == nil || o.ConfigSyncClientProfileTcpMaxWindow == nil {
		var ret int64
		return ret
	}
	return *o.ConfigSyncClientProfileTcpMaxWindow
}

// GetConfigSyncClientProfileTcpMaxWindowOk returns a tuple with the ConfigSyncClientProfileTcpMaxWindow field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetConfigSyncClientProfileTcpMaxWindowOk() (*int64, bool) {
	if o == nil || o.ConfigSyncClientProfileTcpMaxWindow == nil {
		return nil, false
	}
	return o.ConfigSyncClientProfileTcpMaxWindow, true
}

// HasConfigSyncClientProfileTcpMaxWindow returns a boolean if a field has been set.
func (o *Broker) HasConfigSyncClientProfileTcpMaxWindow() bool {
	if o != nil && o.ConfigSyncClientProfileTcpMaxWindow != nil {
		return true
	}

	return false
}

// SetConfigSyncClientProfileTcpMaxWindow gets a reference to the given int64 and assigns it to the ConfigSyncClientProfileTcpMaxWindow field.
func (o *Broker) SetConfigSyncClientProfileTcpMaxWindow(v int64) {
	o.ConfigSyncClientProfileTcpMaxWindow = &v
}

// GetConfigSyncClientProfileTcpMss returns the ConfigSyncClientProfileTcpMss field value if set, zero value otherwise.
func (o *Broker) GetConfigSyncClientProfileTcpMss() int64 {
	if o == nil || o.ConfigSyncClientProfileTcpMss == nil {
		var ret int64
		return ret
	}
	return *o.ConfigSyncClientProfileTcpMss
}

// GetConfigSyncClientProfileTcpMssOk returns a tuple with the ConfigSyncClientProfileTcpMss field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetConfigSyncClientProfileTcpMssOk() (*int64, bool) {
	if o == nil || o.ConfigSyncClientProfileTcpMss == nil {
		return nil, false
	}
	return o.ConfigSyncClientProfileTcpMss, true
}

// HasConfigSyncClientProfileTcpMss returns a boolean if a field has been set.
func (o *Broker) HasConfigSyncClientProfileTcpMss() bool {
	if o != nil && o.ConfigSyncClientProfileTcpMss != nil {
		return true
	}

	return false
}

// SetConfigSyncClientProfileTcpMss gets a reference to the given int64 and assigns it to the ConfigSyncClientProfileTcpMss field.
func (o *Broker) SetConfigSyncClientProfileTcpMss(v int64) {
	o.ConfigSyncClientProfileTcpMss = &v
}

// GetConfigSyncEnabled returns the ConfigSyncEnabled field value if set, zero value otherwise.
func (o *Broker) GetConfigSyncEnabled() bool {
	if o == nil || o.ConfigSyncEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ConfigSyncEnabled
}

// GetConfigSyncEnabledOk returns a tuple with the ConfigSyncEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetConfigSyncEnabledOk() (*bool, bool) {
	if o == nil || o.ConfigSyncEnabled == nil {
		return nil, false
	}
	return o.ConfigSyncEnabled, true
}

// HasConfigSyncEnabled returns a boolean if a field has been set.
func (o *Broker) HasConfigSyncEnabled() bool {
	if o != nil && o.ConfigSyncEnabled != nil {
		return true
	}

	return false
}

// SetConfigSyncEnabled gets a reference to the given bool and assigns it to the ConfigSyncEnabled field.
func (o *Broker) SetConfigSyncEnabled(v bool) {
	o.ConfigSyncEnabled = &v
}

// GetConfigSyncSynchronizeUsernameEnabled returns the ConfigSyncSynchronizeUsernameEnabled field value if set, zero value otherwise.
func (o *Broker) GetConfigSyncSynchronizeUsernameEnabled() bool {
	if o == nil || o.ConfigSyncSynchronizeUsernameEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ConfigSyncSynchronizeUsernameEnabled
}

// GetConfigSyncSynchronizeUsernameEnabledOk returns a tuple with the ConfigSyncSynchronizeUsernameEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetConfigSyncSynchronizeUsernameEnabledOk() (*bool, bool) {
	if o == nil || o.ConfigSyncSynchronizeUsernameEnabled == nil {
		return nil, false
	}
	return o.ConfigSyncSynchronizeUsernameEnabled, true
}

// HasConfigSyncSynchronizeUsernameEnabled returns a boolean if a field has been set.
func (o *Broker) HasConfigSyncSynchronizeUsernameEnabled() bool {
	if o != nil && o.ConfigSyncSynchronizeUsernameEnabled != nil {
		return true
	}

	return false
}

// SetConfigSyncSynchronizeUsernameEnabled gets a reference to the given bool and assigns it to the ConfigSyncSynchronizeUsernameEnabled field.
func (o *Broker) SetConfigSyncSynchronizeUsernameEnabled(v bool) {
	o.ConfigSyncSynchronizeUsernameEnabled = &v
}

// GetConfigSyncTlsEnabled returns the ConfigSyncTlsEnabled field value if set, zero value otherwise.
func (o *Broker) GetConfigSyncTlsEnabled() bool {
	if o == nil || o.ConfigSyncTlsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ConfigSyncTlsEnabled
}

// GetConfigSyncTlsEnabledOk returns a tuple with the ConfigSyncTlsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetConfigSyncTlsEnabledOk() (*bool, bool) {
	if o == nil || o.ConfigSyncTlsEnabled == nil {
		return nil, false
	}
	return o.ConfigSyncTlsEnabled, true
}

// HasConfigSyncTlsEnabled returns a boolean if a field has been set.
func (o *Broker) HasConfigSyncTlsEnabled() bool {
	if o != nil && o.ConfigSyncTlsEnabled != nil {
		return true
	}

	return false
}

// SetConfigSyncTlsEnabled gets a reference to the given bool and assigns it to the ConfigSyncTlsEnabled field.
func (o *Broker) SetConfigSyncTlsEnabled(v bool) {
	o.ConfigSyncTlsEnabled = &v
}

// GetGuaranteedMsgingDefragmentationScheduleDayList returns the GuaranteedMsgingDefragmentationScheduleDayList field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingDefragmentationScheduleDayList() string {
	if o == nil || o.GuaranteedMsgingDefragmentationScheduleDayList == nil {
		var ret string
		return ret
	}
	return *o.GuaranteedMsgingDefragmentationScheduleDayList
}

// GetGuaranteedMsgingDefragmentationScheduleDayListOk returns a tuple with the GuaranteedMsgingDefragmentationScheduleDayList field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingDefragmentationScheduleDayListOk() (*string, bool) {
	if o == nil || o.GuaranteedMsgingDefragmentationScheduleDayList == nil {
		return nil, false
	}
	return o.GuaranteedMsgingDefragmentationScheduleDayList, true
}

// HasGuaranteedMsgingDefragmentationScheduleDayList returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingDefragmentationScheduleDayList() bool {
	if o != nil && o.GuaranteedMsgingDefragmentationScheduleDayList != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingDefragmentationScheduleDayList gets a reference to the given string and assigns it to the GuaranteedMsgingDefragmentationScheduleDayList field.
func (o *Broker) SetGuaranteedMsgingDefragmentationScheduleDayList(v string) {
	o.GuaranteedMsgingDefragmentationScheduleDayList = &v
}

// GetGuaranteedMsgingDefragmentationScheduleEnabled returns the GuaranteedMsgingDefragmentationScheduleEnabled field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingDefragmentationScheduleEnabled() bool {
	if o == nil || o.GuaranteedMsgingDefragmentationScheduleEnabled == nil {
		var ret bool
		return ret
	}
	return *o.GuaranteedMsgingDefragmentationScheduleEnabled
}

// GetGuaranteedMsgingDefragmentationScheduleEnabledOk returns a tuple with the GuaranteedMsgingDefragmentationScheduleEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingDefragmentationScheduleEnabledOk() (*bool, bool) {
	if o == nil || o.GuaranteedMsgingDefragmentationScheduleEnabled == nil {
		return nil, false
	}
	return o.GuaranteedMsgingDefragmentationScheduleEnabled, true
}

// HasGuaranteedMsgingDefragmentationScheduleEnabled returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingDefragmentationScheduleEnabled() bool {
	if o != nil && o.GuaranteedMsgingDefragmentationScheduleEnabled != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingDefragmentationScheduleEnabled gets a reference to the given bool and assigns it to the GuaranteedMsgingDefragmentationScheduleEnabled field.
func (o *Broker) SetGuaranteedMsgingDefragmentationScheduleEnabled(v bool) {
	o.GuaranteedMsgingDefragmentationScheduleEnabled = &v
}

// GetGuaranteedMsgingDefragmentationScheduleTimeList returns the GuaranteedMsgingDefragmentationScheduleTimeList field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingDefragmentationScheduleTimeList() string {
	if o == nil || o.GuaranteedMsgingDefragmentationScheduleTimeList == nil {
		var ret string
		return ret
	}
	return *o.GuaranteedMsgingDefragmentationScheduleTimeList
}

// GetGuaranteedMsgingDefragmentationScheduleTimeListOk returns a tuple with the GuaranteedMsgingDefragmentationScheduleTimeList field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingDefragmentationScheduleTimeListOk() (*string, bool) {
	if o == nil || o.GuaranteedMsgingDefragmentationScheduleTimeList == nil {
		return nil, false
	}
	return o.GuaranteedMsgingDefragmentationScheduleTimeList, true
}

// HasGuaranteedMsgingDefragmentationScheduleTimeList returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingDefragmentationScheduleTimeList() bool {
	if o != nil && o.GuaranteedMsgingDefragmentationScheduleTimeList != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingDefragmentationScheduleTimeList gets a reference to the given string and assigns it to the GuaranteedMsgingDefragmentationScheduleTimeList field.
func (o *Broker) SetGuaranteedMsgingDefragmentationScheduleTimeList(v string) {
	o.GuaranteedMsgingDefragmentationScheduleTimeList = &v
}

// GetGuaranteedMsgingDefragmentationThresholdEnabled returns the GuaranteedMsgingDefragmentationThresholdEnabled field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingDefragmentationThresholdEnabled() bool {
	if o == nil || o.GuaranteedMsgingDefragmentationThresholdEnabled == nil {
		var ret bool
		return ret
	}
	return *o.GuaranteedMsgingDefragmentationThresholdEnabled
}

// GetGuaranteedMsgingDefragmentationThresholdEnabledOk returns a tuple with the GuaranteedMsgingDefragmentationThresholdEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingDefragmentationThresholdEnabledOk() (*bool, bool) {
	if o == nil || o.GuaranteedMsgingDefragmentationThresholdEnabled == nil {
		return nil, false
	}
	return o.GuaranteedMsgingDefragmentationThresholdEnabled, true
}

// HasGuaranteedMsgingDefragmentationThresholdEnabled returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingDefragmentationThresholdEnabled() bool {
	if o != nil && o.GuaranteedMsgingDefragmentationThresholdEnabled != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingDefragmentationThresholdEnabled gets a reference to the given bool and assigns it to the GuaranteedMsgingDefragmentationThresholdEnabled field.
func (o *Broker) SetGuaranteedMsgingDefragmentationThresholdEnabled(v bool) {
	o.GuaranteedMsgingDefragmentationThresholdEnabled = &v
}

// GetGuaranteedMsgingDefragmentationThresholdFragmentationPercentage returns the GuaranteedMsgingDefragmentationThresholdFragmentationPercentage field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingDefragmentationThresholdFragmentationPercentage() int64 {
	if o == nil || o.GuaranteedMsgingDefragmentationThresholdFragmentationPercentage == nil {
		var ret int64
		return ret
	}
	return *o.GuaranteedMsgingDefragmentationThresholdFragmentationPercentage
}

// GetGuaranteedMsgingDefragmentationThresholdFragmentationPercentageOk returns a tuple with the GuaranteedMsgingDefragmentationThresholdFragmentationPercentage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingDefragmentationThresholdFragmentationPercentageOk() (*int64, bool) {
	if o == nil || o.GuaranteedMsgingDefragmentationThresholdFragmentationPercentage == nil {
		return nil, false
	}
	return o.GuaranteedMsgingDefragmentationThresholdFragmentationPercentage, true
}

// HasGuaranteedMsgingDefragmentationThresholdFragmentationPercentage returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingDefragmentationThresholdFragmentationPercentage() bool {
	if o != nil && o.GuaranteedMsgingDefragmentationThresholdFragmentationPercentage != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingDefragmentationThresholdFragmentationPercentage gets a reference to the given int64 and assigns it to the GuaranteedMsgingDefragmentationThresholdFragmentationPercentage field.
func (o *Broker) SetGuaranteedMsgingDefragmentationThresholdFragmentationPercentage(v int64) {
	o.GuaranteedMsgingDefragmentationThresholdFragmentationPercentage = &v
}

// GetGuaranteedMsgingDefragmentationThresholdMinInterval returns the GuaranteedMsgingDefragmentationThresholdMinInterval field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingDefragmentationThresholdMinInterval() int64 {
	if o == nil || o.GuaranteedMsgingDefragmentationThresholdMinInterval == nil {
		var ret int64
		return ret
	}
	return *o.GuaranteedMsgingDefragmentationThresholdMinInterval
}

// GetGuaranteedMsgingDefragmentationThresholdMinIntervalOk returns a tuple with the GuaranteedMsgingDefragmentationThresholdMinInterval field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingDefragmentationThresholdMinIntervalOk() (*int64, bool) {
	if o == nil || o.GuaranteedMsgingDefragmentationThresholdMinInterval == nil {
		return nil, false
	}
	return o.GuaranteedMsgingDefragmentationThresholdMinInterval, true
}

// HasGuaranteedMsgingDefragmentationThresholdMinInterval returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingDefragmentationThresholdMinInterval() bool {
	if o != nil && o.GuaranteedMsgingDefragmentationThresholdMinInterval != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingDefragmentationThresholdMinInterval gets a reference to the given int64 and assigns it to the GuaranteedMsgingDefragmentationThresholdMinInterval field.
func (o *Broker) SetGuaranteedMsgingDefragmentationThresholdMinInterval(v int64) {
	o.GuaranteedMsgingDefragmentationThresholdMinInterval = &v
}

// GetGuaranteedMsgingDefragmentationThresholdUsagePercentage returns the GuaranteedMsgingDefragmentationThresholdUsagePercentage field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingDefragmentationThresholdUsagePercentage() int64 {
	if o == nil || o.GuaranteedMsgingDefragmentationThresholdUsagePercentage == nil {
		var ret int64
		return ret
	}
	return *o.GuaranteedMsgingDefragmentationThresholdUsagePercentage
}

// GetGuaranteedMsgingDefragmentationThresholdUsagePercentageOk returns a tuple with the GuaranteedMsgingDefragmentationThresholdUsagePercentage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingDefragmentationThresholdUsagePercentageOk() (*int64, bool) {
	if o == nil || o.GuaranteedMsgingDefragmentationThresholdUsagePercentage == nil {
		return nil, false
	}
	return o.GuaranteedMsgingDefragmentationThresholdUsagePercentage, true
}

// HasGuaranteedMsgingDefragmentationThresholdUsagePercentage returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingDefragmentationThresholdUsagePercentage() bool {
	if o != nil && o.GuaranteedMsgingDefragmentationThresholdUsagePercentage != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingDefragmentationThresholdUsagePercentage gets a reference to the given int64 and assigns it to the GuaranteedMsgingDefragmentationThresholdUsagePercentage field.
func (o *Broker) SetGuaranteedMsgingDefragmentationThresholdUsagePercentage(v int64) {
	o.GuaranteedMsgingDefragmentationThresholdUsagePercentage = &v
}

// GetGuaranteedMsgingEnabled returns the GuaranteedMsgingEnabled field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingEnabled() bool {
	if o == nil || o.GuaranteedMsgingEnabled == nil {
		var ret bool
		return ret
	}
	return *o.GuaranteedMsgingEnabled
}

// GetGuaranteedMsgingEnabledOk returns a tuple with the GuaranteedMsgingEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingEnabledOk() (*bool, bool) {
	if o == nil || o.GuaranteedMsgingEnabled == nil {
		return nil, false
	}
	return o.GuaranteedMsgingEnabled, true
}

// HasGuaranteedMsgingEnabled returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingEnabled() bool {
	if o != nil && o.GuaranteedMsgingEnabled != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingEnabled gets a reference to the given bool and assigns it to the GuaranteedMsgingEnabled field.
func (o *Broker) SetGuaranteedMsgingEnabled(v bool) {
	o.GuaranteedMsgingEnabled = &v
}

// GetGuaranteedMsgingEventCacheUsageThreshold returns the GuaranteedMsgingEventCacheUsageThreshold field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingEventCacheUsageThreshold() EventThreshold {
	if o == nil || o.GuaranteedMsgingEventCacheUsageThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.GuaranteedMsgingEventCacheUsageThreshold
}

// GetGuaranteedMsgingEventCacheUsageThresholdOk returns a tuple with the GuaranteedMsgingEventCacheUsageThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingEventCacheUsageThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.GuaranteedMsgingEventCacheUsageThreshold == nil {
		return nil, false
	}
	return o.GuaranteedMsgingEventCacheUsageThreshold, true
}

// HasGuaranteedMsgingEventCacheUsageThreshold returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingEventCacheUsageThreshold() bool {
	if o != nil && o.GuaranteedMsgingEventCacheUsageThreshold != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingEventCacheUsageThreshold gets a reference to the given EventThreshold and assigns it to the GuaranteedMsgingEventCacheUsageThreshold field.
func (o *Broker) SetGuaranteedMsgingEventCacheUsageThreshold(v EventThreshold) {
	o.GuaranteedMsgingEventCacheUsageThreshold = &v
}

// GetGuaranteedMsgingEventDeliveredUnackedThreshold returns the GuaranteedMsgingEventDeliveredUnackedThreshold field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingEventDeliveredUnackedThreshold() EventThresholdByPercent {
	if o == nil || o.GuaranteedMsgingEventDeliveredUnackedThreshold == nil {
		var ret EventThresholdByPercent
		return ret
	}
	return *o.GuaranteedMsgingEventDeliveredUnackedThreshold
}

// GetGuaranteedMsgingEventDeliveredUnackedThresholdOk returns a tuple with the GuaranteedMsgingEventDeliveredUnackedThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingEventDeliveredUnackedThresholdOk() (*EventThresholdByPercent, bool) {
	if o == nil || o.GuaranteedMsgingEventDeliveredUnackedThreshold == nil {
		return nil, false
	}
	return o.GuaranteedMsgingEventDeliveredUnackedThreshold, true
}

// HasGuaranteedMsgingEventDeliveredUnackedThreshold returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingEventDeliveredUnackedThreshold() bool {
	if o != nil && o.GuaranteedMsgingEventDeliveredUnackedThreshold != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingEventDeliveredUnackedThreshold gets a reference to the given EventThresholdByPercent and assigns it to the GuaranteedMsgingEventDeliveredUnackedThreshold field.
func (o *Broker) SetGuaranteedMsgingEventDeliveredUnackedThreshold(v EventThresholdByPercent) {
	o.GuaranteedMsgingEventDeliveredUnackedThreshold = &v
}

// GetGuaranteedMsgingEventDiskUsageThreshold returns the GuaranteedMsgingEventDiskUsageThreshold field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingEventDiskUsageThreshold() EventThresholdByPercent {
	if o == nil || o.GuaranteedMsgingEventDiskUsageThreshold == nil {
		var ret EventThresholdByPercent
		return ret
	}
	return *o.GuaranteedMsgingEventDiskUsageThreshold
}

// GetGuaranteedMsgingEventDiskUsageThresholdOk returns a tuple with the GuaranteedMsgingEventDiskUsageThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingEventDiskUsageThresholdOk() (*EventThresholdByPercent, bool) {
	if o == nil || o.GuaranteedMsgingEventDiskUsageThreshold == nil {
		return nil, false
	}
	return o.GuaranteedMsgingEventDiskUsageThreshold, true
}

// HasGuaranteedMsgingEventDiskUsageThreshold returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingEventDiskUsageThreshold() bool {
	if o != nil && o.GuaranteedMsgingEventDiskUsageThreshold != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingEventDiskUsageThreshold gets a reference to the given EventThresholdByPercent and assigns it to the GuaranteedMsgingEventDiskUsageThreshold field.
func (o *Broker) SetGuaranteedMsgingEventDiskUsageThreshold(v EventThresholdByPercent) {
	o.GuaranteedMsgingEventDiskUsageThreshold = &v
}

// GetGuaranteedMsgingEventEgressFlowCountThreshold returns the GuaranteedMsgingEventEgressFlowCountThreshold field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingEventEgressFlowCountThreshold() EventThreshold {
	if o == nil || o.GuaranteedMsgingEventEgressFlowCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.GuaranteedMsgingEventEgressFlowCountThreshold
}

// GetGuaranteedMsgingEventEgressFlowCountThresholdOk returns a tuple with the GuaranteedMsgingEventEgressFlowCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingEventEgressFlowCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.GuaranteedMsgingEventEgressFlowCountThreshold == nil {
		return nil, false
	}
	return o.GuaranteedMsgingEventEgressFlowCountThreshold, true
}

// HasGuaranteedMsgingEventEgressFlowCountThreshold returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingEventEgressFlowCountThreshold() bool {
	if o != nil && o.GuaranteedMsgingEventEgressFlowCountThreshold != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingEventEgressFlowCountThreshold gets a reference to the given EventThreshold and assigns it to the GuaranteedMsgingEventEgressFlowCountThreshold field.
func (o *Broker) SetGuaranteedMsgingEventEgressFlowCountThreshold(v EventThreshold) {
	o.GuaranteedMsgingEventEgressFlowCountThreshold = &v
}

// GetGuaranteedMsgingEventEndpointCountThreshold returns the GuaranteedMsgingEventEndpointCountThreshold field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingEventEndpointCountThreshold() EventThreshold {
	if o == nil || o.GuaranteedMsgingEventEndpointCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.GuaranteedMsgingEventEndpointCountThreshold
}

// GetGuaranteedMsgingEventEndpointCountThresholdOk returns a tuple with the GuaranteedMsgingEventEndpointCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingEventEndpointCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.GuaranteedMsgingEventEndpointCountThreshold == nil {
		return nil, false
	}
	return o.GuaranteedMsgingEventEndpointCountThreshold, true
}

// HasGuaranteedMsgingEventEndpointCountThreshold returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingEventEndpointCountThreshold() bool {
	if o != nil && o.GuaranteedMsgingEventEndpointCountThreshold != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingEventEndpointCountThreshold gets a reference to the given EventThreshold and assigns it to the GuaranteedMsgingEventEndpointCountThreshold field.
func (o *Broker) SetGuaranteedMsgingEventEndpointCountThreshold(v EventThreshold) {
	o.GuaranteedMsgingEventEndpointCountThreshold = &v
}

// GetGuaranteedMsgingEventIngressFlowCountThreshold returns the GuaranteedMsgingEventIngressFlowCountThreshold field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingEventIngressFlowCountThreshold() EventThreshold {
	if o == nil || o.GuaranteedMsgingEventIngressFlowCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.GuaranteedMsgingEventIngressFlowCountThreshold
}

// GetGuaranteedMsgingEventIngressFlowCountThresholdOk returns a tuple with the GuaranteedMsgingEventIngressFlowCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingEventIngressFlowCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.GuaranteedMsgingEventIngressFlowCountThreshold == nil {
		return nil, false
	}
	return o.GuaranteedMsgingEventIngressFlowCountThreshold, true
}

// HasGuaranteedMsgingEventIngressFlowCountThreshold returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingEventIngressFlowCountThreshold() bool {
	if o != nil && o.GuaranteedMsgingEventIngressFlowCountThreshold != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingEventIngressFlowCountThreshold gets a reference to the given EventThreshold and assigns it to the GuaranteedMsgingEventIngressFlowCountThreshold field.
func (o *Broker) SetGuaranteedMsgingEventIngressFlowCountThreshold(v EventThreshold) {
	o.GuaranteedMsgingEventIngressFlowCountThreshold = &v
}

// GetGuaranteedMsgingEventMsgCountThreshold returns the GuaranteedMsgingEventMsgCountThreshold field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingEventMsgCountThreshold() EventThresholdByPercent {
	if o == nil || o.GuaranteedMsgingEventMsgCountThreshold == nil {
		var ret EventThresholdByPercent
		return ret
	}
	return *o.GuaranteedMsgingEventMsgCountThreshold
}

// GetGuaranteedMsgingEventMsgCountThresholdOk returns a tuple with the GuaranteedMsgingEventMsgCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingEventMsgCountThresholdOk() (*EventThresholdByPercent, bool) {
	if o == nil || o.GuaranteedMsgingEventMsgCountThreshold == nil {
		return nil, false
	}
	return o.GuaranteedMsgingEventMsgCountThreshold, true
}

// HasGuaranteedMsgingEventMsgCountThreshold returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingEventMsgCountThreshold() bool {
	if o != nil && o.GuaranteedMsgingEventMsgCountThreshold != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingEventMsgCountThreshold gets a reference to the given EventThresholdByPercent and assigns it to the GuaranteedMsgingEventMsgCountThreshold field.
func (o *Broker) SetGuaranteedMsgingEventMsgCountThreshold(v EventThresholdByPercent) {
	o.GuaranteedMsgingEventMsgCountThreshold = &v
}

// GetGuaranteedMsgingEventMsgSpoolFileCountThreshold returns the GuaranteedMsgingEventMsgSpoolFileCountThreshold field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingEventMsgSpoolFileCountThreshold() EventThresholdByPercent {
	if o == nil || o.GuaranteedMsgingEventMsgSpoolFileCountThreshold == nil {
		var ret EventThresholdByPercent
		return ret
	}
	return *o.GuaranteedMsgingEventMsgSpoolFileCountThreshold
}

// GetGuaranteedMsgingEventMsgSpoolFileCountThresholdOk returns a tuple with the GuaranteedMsgingEventMsgSpoolFileCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingEventMsgSpoolFileCountThresholdOk() (*EventThresholdByPercent, bool) {
	if o == nil || o.GuaranteedMsgingEventMsgSpoolFileCountThreshold == nil {
		return nil, false
	}
	return o.GuaranteedMsgingEventMsgSpoolFileCountThreshold, true
}

// HasGuaranteedMsgingEventMsgSpoolFileCountThreshold returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingEventMsgSpoolFileCountThreshold() bool {
	if o != nil && o.GuaranteedMsgingEventMsgSpoolFileCountThreshold != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingEventMsgSpoolFileCountThreshold gets a reference to the given EventThresholdByPercent and assigns it to the GuaranteedMsgingEventMsgSpoolFileCountThreshold field.
func (o *Broker) SetGuaranteedMsgingEventMsgSpoolFileCountThreshold(v EventThresholdByPercent) {
	o.GuaranteedMsgingEventMsgSpoolFileCountThreshold = &v
}

// GetGuaranteedMsgingEventMsgSpoolUsageThreshold returns the GuaranteedMsgingEventMsgSpoolUsageThreshold field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingEventMsgSpoolUsageThreshold() EventThreshold {
	if o == nil || o.GuaranteedMsgingEventMsgSpoolUsageThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.GuaranteedMsgingEventMsgSpoolUsageThreshold
}

// GetGuaranteedMsgingEventMsgSpoolUsageThresholdOk returns a tuple with the GuaranteedMsgingEventMsgSpoolUsageThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingEventMsgSpoolUsageThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.GuaranteedMsgingEventMsgSpoolUsageThreshold == nil {
		return nil, false
	}
	return o.GuaranteedMsgingEventMsgSpoolUsageThreshold, true
}

// HasGuaranteedMsgingEventMsgSpoolUsageThreshold returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingEventMsgSpoolUsageThreshold() bool {
	if o != nil && o.GuaranteedMsgingEventMsgSpoolUsageThreshold != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingEventMsgSpoolUsageThreshold gets a reference to the given EventThreshold and assigns it to the GuaranteedMsgingEventMsgSpoolUsageThreshold field.
func (o *Broker) SetGuaranteedMsgingEventMsgSpoolUsageThreshold(v EventThreshold) {
	o.GuaranteedMsgingEventMsgSpoolUsageThreshold = &v
}

// GetGuaranteedMsgingEventTransactedSessionCountThreshold returns the GuaranteedMsgingEventTransactedSessionCountThreshold field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingEventTransactedSessionCountThreshold() EventThreshold {
	if o == nil || o.GuaranteedMsgingEventTransactedSessionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.GuaranteedMsgingEventTransactedSessionCountThreshold
}

// GetGuaranteedMsgingEventTransactedSessionCountThresholdOk returns a tuple with the GuaranteedMsgingEventTransactedSessionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingEventTransactedSessionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.GuaranteedMsgingEventTransactedSessionCountThreshold == nil {
		return nil, false
	}
	return o.GuaranteedMsgingEventTransactedSessionCountThreshold, true
}

// HasGuaranteedMsgingEventTransactedSessionCountThreshold returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingEventTransactedSessionCountThreshold() bool {
	if o != nil && o.GuaranteedMsgingEventTransactedSessionCountThreshold != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingEventTransactedSessionCountThreshold gets a reference to the given EventThreshold and assigns it to the GuaranteedMsgingEventTransactedSessionCountThreshold field.
func (o *Broker) SetGuaranteedMsgingEventTransactedSessionCountThreshold(v EventThreshold) {
	o.GuaranteedMsgingEventTransactedSessionCountThreshold = &v
}

// GetGuaranteedMsgingEventTransactedSessionResourceCountThreshold returns the GuaranteedMsgingEventTransactedSessionResourceCountThreshold field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingEventTransactedSessionResourceCountThreshold() EventThresholdByPercent {
	if o == nil || o.GuaranteedMsgingEventTransactedSessionResourceCountThreshold == nil {
		var ret EventThresholdByPercent
		return ret
	}
	return *o.GuaranteedMsgingEventTransactedSessionResourceCountThreshold
}

// GetGuaranteedMsgingEventTransactedSessionResourceCountThresholdOk returns a tuple with the GuaranteedMsgingEventTransactedSessionResourceCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingEventTransactedSessionResourceCountThresholdOk() (*EventThresholdByPercent, bool) {
	if o == nil || o.GuaranteedMsgingEventTransactedSessionResourceCountThreshold == nil {
		return nil, false
	}
	return o.GuaranteedMsgingEventTransactedSessionResourceCountThreshold, true
}

// HasGuaranteedMsgingEventTransactedSessionResourceCountThreshold returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingEventTransactedSessionResourceCountThreshold() bool {
	if o != nil && o.GuaranteedMsgingEventTransactedSessionResourceCountThreshold != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingEventTransactedSessionResourceCountThreshold gets a reference to the given EventThresholdByPercent and assigns it to the GuaranteedMsgingEventTransactedSessionResourceCountThreshold field.
func (o *Broker) SetGuaranteedMsgingEventTransactedSessionResourceCountThreshold(v EventThresholdByPercent) {
	o.GuaranteedMsgingEventTransactedSessionResourceCountThreshold = &v
}

// GetGuaranteedMsgingEventTransactionCountThreshold returns the GuaranteedMsgingEventTransactionCountThreshold field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingEventTransactionCountThreshold() EventThreshold {
	if o == nil || o.GuaranteedMsgingEventTransactionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.GuaranteedMsgingEventTransactionCountThreshold
}

// GetGuaranteedMsgingEventTransactionCountThresholdOk returns a tuple with the GuaranteedMsgingEventTransactionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingEventTransactionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.GuaranteedMsgingEventTransactionCountThreshold == nil {
		return nil, false
	}
	return o.GuaranteedMsgingEventTransactionCountThreshold, true
}

// HasGuaranteedMsgingEventTransactionCountThreshold returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingEventTransactionCountThreshold() bool {
	if o != nil && o.GuaranteedMsgingEventTransactionCountThreshold != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingEventTransactionCountThreshold gets a reference to the given EventThreshold and assigns it to the GuaranteedMsgingEventTransactionCountThreshold field.
func (o *Broker) SetGuaranteedMsgingEventTransactionCountThreshold(v EventThreshold) {
	o.GuaranteedMsgingEventTransactionCountThreshold = &v
}

// GetGuaranteedMsgingMaxCacheUsage returns the GuaranteedMsgingMaxCacheUsage field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingMaxCacheUsage() int32 {
	if o == nil || o.GuaranteedMsgingMaxCacheUsage == nil {
		var ret int32
		return ret
	}
	return *o.GuaranteedMsgingMaxCacheUsage
}

// GetGuaranteedMsgingMaxCacheUsageOk returns a tuple with the GuaranteedMsgingMaxCacheUsage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingMaxCacheUsageOk() (*int32, bool) {
	if o == nil || o.GuaranteedMsgingMaxCacheUsage == nil {
		return nil, false
	}
	return o.GuaranteedMsgingMaxCacheUsage, true
}

// HasGuaranteedMsgingMaxCacheUsage returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingMaxCacheUsage() bool {
	if o != nil && o.GuaranteedMsgingMaxCacheUsage != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingMaxCacheUsage gets a reference to the given int32 and assigns it to the GuaranteedMsgingMaxCacheUsage field.
func (o *Broker) SetGuaranteedMsgingMaxCacheUsage(v int32) {
	o.GuaranteedMsgingMaxCacheUsage = &v
}

// GetGuaranteedMsgingMaxMsgSpoolUsage returns the GuaranteedMsgingMaxMsgSpoolUsage field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingMaxMsgSpoolUsage() int64 {
	if o == nil || o.GuaranteedMsgingMaxMsgSpoolUsage == nil {
		var ret int64
		return ret
	}
	return *o.GuaranteedMsgingMaxMsgSpoolUsage
}

// GetGuaranteedMsgingMaxMsgSpoolUsageOk returns a tuple with the GuaranteedMsgingMaxMsgSpoolUsage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingMaxMsgSpoolUsageOk() (*int64, bool) {
	if o == nil || o.GuaranteedMsgingMaxMsgSpoolUsage == nil {
		return nil, false
	}
	return o.GuaranteedMsgingMaxMsgSpoolUsage, true
}

// HasGuaranteedMsgingMaxMsgSpoolUsage returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingMaxMsgSpoolUsage() bool {
	if o != nil && o.GuaranteedMsgingMaxMsgSpoolUsage != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingMaxMsgSpoolUsage gets a reference to the given int64 and assigns it to the GuaranteedMsgingMaxMsgSpoolUsage field.
func (o *Broker) SetGuaranteedMsgingMaxMsgSpoolUsage(v int64) {
	o.GuaranteedMsgingMaxMsgSpoolUsage = &v
}

// GetGuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout returns the GuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout() int64 {
	if o == nil || o.GuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout == nil {
		var ret int64
		return ret
	}
	return *o.GuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout
}

// GetGuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeoutOk returns a tuple with the GuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeoutOk() (*int64, bool) {
	if o == nil || o.GuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout == nil {
		return nil, false
	}
	return o.GuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout, true
}

// HasGuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout() bool {
	if o != nil && o.GuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout gets a reference to the given int64 and assigns it to the GuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout field.
func (o *Broker) SetGuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout(v int64) {
	o.GuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout = &v
}

// GetGuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout returns the GuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout() int64 {
	if o == nil || o.GuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout == nil {
		var ret int64
		return ret
	}
	return *o.GuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout
}

// GetGuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeoutOk returns a tuple with the GuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeoutOk() (*int64, bool) {
	if o == nil || o.GuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout == nil {
		return nil, false
	}
	return o.GuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout, true
}

// HasGuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout() bool {
	if o != nil && o.GuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout gets a reference to the given int64 and assigns it to the GuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout field.
func (o *Broker) SetGuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout(v int64) {
	o.GuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout = &v
}

// GetGuaranteedMsgingTransactionReplicationCompatibilityMode returns the GuaranteedMsgingTransactionReplicationCompatibilityMode field value if set, zero value otherwise.
func (o *Broker) GetGuaranteedMsgingTransactionReplicationCompatibilityMode() string {
	if o == nil || o.GuaranteedMsgingTransactionReplicationCompatibilityMode == nil {
		var ret string
		return ret
	}
	return *o.GuaranteedMsgingTransactionReplicationCompatibilityMode
}

// GetGuaranteedMsgingTransactionReplicationCompatibilityModeOk returns a tuple with the GuaranteedMsgingTransactionReplicationCompatibilityMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetGuaranteedMsgingTransactionReplicationCompatibilityModeOk() (*string, bool) {
	if o == nil || o.GuaranteedMsgingTransactionReplicationCompatibilityMode == nil {
		return nil, false
	}
	return o.GuaranteedMsgingTransactionReplicationCompatibilityMode, true
}

// HasGuaranteedMsgingTransactionReplicationCompatibilityMode returns a boolean if a field has been set.
func (o *Broker) HasGuaranteedMsgingTransactionReplicationCompatibilityMode() bool {
	if o != nil && o.GuaranteedMsgingTransactionReplicationCompatibilityMode != nil {
		return true
	}

	return false
}

// SetGuaranteedMsgingTransactionReplicationCompatibilityMode gets a reference to the given string and assigns it to the GuaranteedMsgingTransactionReplicationCompatibilityMode field.
func (o *Broker) SetGuaranteedMsgingTransactionReplicationCompatibilityMode(v string) {
	o.GuaranteedMsgingTransactionReplicationCompatibilityMode = &v
}

// GetOauthProfileDefault returns the OauthProfileDefault field value if set, zero value otherwise.
func (o *Broker) GetOauthProfileDefault() string {
	if o == nil || o.OauthProfileDefault == nil {
		var ret string
		return ret
	}
	return *o.OauthProfileDefault
}

// GetOauthProfileDefaultOk returns a tuple with the OauthProfileDefault field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetOauthProfileDefaultOk() (*string, bool) {
	if o == nil || o.OauthProfileDefault == nil {
		return nil, false
	}
	return o.OauthProfileDefault, true
}

// HasOauthProfileDefault returns a boolean if a field has been set.
func (o *Broker) HasOauthProfileDefault() bool {
	if o != nil && o.OauthProfileDefault != nil {
		return true
	}

	return false
}

// SetOauthProfileDefault gets a reference to the given string and assigns it to the OauthProfileDefault field.
func (o *Broker) SetOauthProfileDefault(v string) {
	o.OauthProfileDefault = &v
}

// GetServiceAmqpEnabled returns the ServiceAmqpEnabled field value if set, zero value otherwise.
func (o *Broker) GetServiceAmqpEnabled() bool {
	if o == nil || o.ServiceAmqpEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceAmqpEnabled
}

// GetServiceAmqpEnabledOk returns a tuple with the ServiceAmqpEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceAmqpEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceAmqpEnabled == nil {
		return nil, false
	}
	return o.ServiceAmqpEnabled, true
}

// HasServiceAmqpEnabled returns a boolean if a field has been set.
func (o *Broker) HasServiceAmqpEnabled() bool {
	if o != nil && o.ServiceAmqpEnabled != nil {
		return true
	}

	return false
}

// SetServiceAmqpEnabled gets a reference to the given bool and assigns it to the ServiceAmqpEnabled field.
func (o *Broker) SetServiceAmqpEnabled(v bool) {
	o.ServiceAmqpEnabled = &v
}

// GetServiceAmqpTlsListenPort returns the ServiceAmqpTlsListenPort field value if set, zero value otherwise.
func (o *Broker) GetServiceAmqpTlsListenPort() int64 {
	if o == nil || o.ServiceAmqpTlsListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceAmqpTlsListenPort
}

// GetServiceAmqpTlsListenPortOk returns a tuple with the ServiceAmqpTlsListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceAmqpTlsListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceAmqpTlsListenPort == nil {
		return nil, false
	}
	return o.ServiceAmqpTlsListenPort, true
}

// HasServiceAmqpTlsListenPort returns a boolean if a field has been set.
func (o *Broker) HasServiceAmqpTlsListenPort() bool {
	if o != nil && o.ServiceAmqpTlsListenPort != nil {
		return true
	}

	return false
}

// SetServiceAmqpTlsListenPort gets a reference to the given int64 and assigns it to the ServiceAmqpTlsListenPort field.
func (o *Broker) SetServiceAmqpTlsListenPort(v int64) {
	o.ServiceAmqpTlsListenPort = &v
}

// GetServiceEventConnectionCountThreshold returns the ServiceEventConnectionCountThreshold field value if set, zero value otherwise.
func (o *Broker) GetServiceEventConnectionCountThreshold() EventThreshold {
	if o == nil || o.ServiceEventConnectionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.ServiceEventConnectionCountThreshold
}

// GetServiceEventConnectionCountThresholdOk returns a tuple with the ServiceEventConnectionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceEventConnectionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.ServiceEventConnectionCountThreshold == nil {
		return nil, false
	}
	return o.ServiceEventConnectionCountThreshold, true
}

// HasServiceEventConnectionCountThreshold returns a boolean if a field has been set.
func (o *Broker) HasServiceEventConnectionCountThreshold() bool {
	if o != nil && o.ServiceEventConnectionCountThreshold != nil {
		return true
	}

	return false
}

// SetServiceEventConnectionCountThreshold gets a reference to the given EventThreshold and assigns it to the ServiceEventConnectionCountThreshold field.
func (o *Broker) SetServiceEventConnectionCountThreshold(v EventThreshold) {
	o.ServiceEventConnectionCountThreshold = &v
}

// GetServiceHealthCheckEnabled returns the ServiceHealthCheckEnabled field value if set, zero value otherwise.
func (o *Broker) GetServiceHealthCheckEnabled() bool {
	if o == nil || o.ServiceHealthCheckEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceHealthCheckEnabled
}

// GetServiceHealthCheckEnabledOk returns a tuple with the ServiceHealthCheckEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceHealthCheckEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceHealthCheckEnabled == nil {
		return nil, false
	}
	return o.ServiceHealthCheckEnabled, true
}

// HasServiceHealthCheckEnabled returns a boolean if a field has been set.
func (o *Broker) HasServiceHealthCheckEnabled() bool {
	if o != nil && o.ServiceHealthCheckEnabled != nil {
		return true
	}

	return false
}

// SetServiceHealthCheckEnabled gets a reference to the given bool and assigns it to the ServiceHealthCheckEnabled field.
func (o *Broker) SetServiceHealthCheckEnabled(v bool) {
	o.ServiceHealthCheckEnabled = &v
}

// GetServiceHealthCheckListenPort returns the ServiceHealthCheckListenPort field value if set, zero value otherwise.
func (o *Broker) GetServiceHealthCheckListenPort() int64 {
	if o == nil || o.ServiceHealthCheckListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceHealthCheckListenPort
}

// GetServiceHealthCheckListenPortOk returns a tuple with the ServiceHealthCheckListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceHealthCheckListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceHealthCheckListenPort == nil {
		return nil, false
	}
	return o.ServiceHealthCheckListenPort, true
}

// HasServiceHealthCheckListenPort returns a boolean if a field has been set.
func (o *Broker) HasServiceHealthCheckListenPort() bool {
	if o != nil && o.ServiceHealthCheckListenPort != nil {
		return true
	}

	return false
}

// SetServiceHealthCheckListenPort gets a reference to the given int64 and assigns it to the ServiceHealthCheckListenPort field.
func (o *Broker) SetServiceHealthCheckListenPort(v int64) {
	o.ServiceHealthCheckListenPort = &v
}

// GetServiceMateLinkEnabled returns the ServiceMateLinkEnabled field value if set, zero value otherwise.
func (o *Broker) GetServiceMateLinkEnabled() bool {
	if o == nil || o.ServiceMateLinkEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceMateLinkEnabled
}

// GetServiceMateLinkEnabledOk returns a tuple with the ServiceMateLinkEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceMateLinkEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceMateLinkEnabled == nil {
		return nil, false
	}
	return o.ServiceMateLinkEnabled, true
}

// HasServiceMateLinkEnabled returns a boolean if a field has been set.
func (o *Broker) HasServiceMateLinkEnabled() bool {
	if o != nil && o.ServiceMateLinkEnabled != nil {
		return true
	}

	return false
}

// SetServiceMateLinkEnabled gets a reference to the given bool and assigns it to the ServiceMateLinkEnabled field.
func (o *Broker) SetServiceMateLinkEnabled(v bool) {
	o.ServiceMateLinkEnabled = &v
}

// GetServiceMateLinkListenPort returns the ServiceMateLinkListenPort field value if set, zero value otherwise.
func (o *Broker) GetServiceMateLinkListenPort() int64 {
	if o == nil || o.ServiceMateLinkListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceMateLinkListenPort
}

// GetServiceMateLinkListenPortOk returns a tuple with the ServiceMateLinkListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceMateLinkListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceMateLinkListenPort == nil {
		return nil, false
	}
	return o.ServiceMateLinkListenPort, true
}

// HasServiceMateLinkListenPort returns a boolean if a field has been set.
func (o *Broker) HasServiceMateLinkListenPort() bool {
	if o != nil && o.ServiceMateLinkListenPort != nil {
		return true
	}

	return false
}

// SetServiceMateLinkListenPort gets a reference to the given int64 and assigns it to the ServiceMateLinkListenPort field.
func (o *Broker) SetServiceMateLinkListenPort(v int64) {
	o.ServiceMateLinkListenPort = &v
}

// GetServiceMqttEnabled returns the ServiceMqttEnabled field value if set, zero value otherwise.
func (o *Broker) GetServiceMqttEnabled() bool {
	if o == nil || o.ServiceMqttEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceMqttEnabled
}

// GetServiceMqttEnabledOk returns a tuple with the ServiceMqttEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceMqttEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceMqttEnabled == nil {
		return nil, false
	}
	return o.ServiceMqttEnabled, true
}

// HasServiceMqttEnabled returns a boolean if a field has been set.
func (o *Broker) HasServiceMqttEnabled() bool {
	if o != nil && o.ServiceMqttEnabled != nil {
		return true
	}

	return false
}

// SetServiceMqttEnabled gets a reference to the given bool and assigns it to the ServiceMqttEnabled field.
func (o *Broker) SetServiceMqttEnabled(v bool) {
	o.ServiceMqttEnabled = &v
}

// GetServiceMsgBackboneEnabled returns the ServiceMsgBackboneEnabled field value if set, zero value otherwise.
func (o *Broker) GetServiceMsgBackboneEnabled() bool {
	if o == nil || o.ServiceMsgBackboneEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceMsgBackboneEnabled
}

// GetServiceMsgBackboneEnabledOk returns a tuple with the ServiceMsgBackboneEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceMsgBackboneEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceMsgBackboneEnabled == nil {
		return nil, false
	}
	return o.ServiceMsgBackboneEnabled, true
}

// HasServiceMsgBackboneEnabled returns a boolean if a field has been set.
func (o *Broker) HasServiceMsgBackboneEnabled() bool {
	if o != nil && o.ServiceMsgBackboneEnabled != nil {
		return true
	}

	return false
}

// SetServiceMsgBackboneEnabled gets a reference to the given bool and assigns it to the ServiceMsgBackboneEnabled field.
func (o *Broker) SetServiceMsgBackboneEnabled(v bool) {
	o.ServiceMsgBackboneEnabled = &v
}

// GetServiceRedundancyEnabled returns the ServiceRedundancyEnabled field value if set, zero value otherwise.
func (o *Broker) GetServiceRedundancyEnabled() bool {
	if o == nil || o.ServiceRedundancyEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceRedundancyEnabled
}

// GetServiceRedundancyEnabledOk returns a tuple with the ServiceRedundancyEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceRedundancyEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceRedundancyEnabled == nil {
		return nil, false
	}
	return o.ServiceRedundancyEnabled, true
}

// HasServiceRedundancyEnabled returns a boolean if a field has been set.
func (o *Broker) HasServiceRedundancyEnabled() bool {
	if o != nil && o.ServiceRedundancyEnabled != nil {
		return true
	}

	return false
}

// SetServiceRedundancyEnabled gets a reference to the given bool and assigns it to the ServiceRedundancyEnabled field.
func (o *Broker) SetServiceRedundancyEnabled(v bool) {
	o.ServiceRedundancyEnabled = &v
}

// GetServiceRedundancyFirstListenPort returns the ServiceRedundancyFirstListenPort field value if set, zero value otherwise.
func (o *Broker) GetServiceRedundancyFirstListenPort() int64 {
	if o == nil || o.ServiceRedundancyFirstListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceRedundancyFirstListenPort
}

// GetServiceRedundancyFirstListenPortOk returns a tuple with the ServiceRedundancyFirstListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceRedundancyFirstListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceRedundancyFirstListenPort == nil {
		return nil, false
	}
	return o.ServiceRedundancyFirstListenPort, true
}

// HasServiceRedundancyFirstListenPort returns a boolean if a field has been set.
func (o *Broker) HasServiceRedundancyFirstListenPort() bool {
	if o != nil && o.ServiceRedundancyFirstListenPort != nil {
		return true
	}

	return false
}

// SetServiceRedundancyFirstListenPort gets a reference to the given int64 and assigns it to the ServiceRedundancyFirstListenPort field.
func (o *Broker) SetServiceRedundancyFirstListenPort(v int64) {
	o.ServiceRedundancyFirstListenPort = &v
}

// GetServiceRestEventOutgoingConnectionCountThreshold returns the ServiceRestEventOutgoingConnectionCountThreshold field value if set, zero value otherwise.
func (o *Broker) GetServiceRestEventOutgoingConnectionCountThreshold() EventThreshold {
	if o == nil || o.ServiceRestEventOutgoingConnectionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.ServiceRestEventOutgoingConnectionCountThreshold
}

// GetServiceRestEventOutgoingConnectionCountThresholdOk returns a tuple with the ServiceRestEventOutgoingConnectionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceRestEventOutgoingConnectionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.ServiceRestEventOutgoingConnectionCountThreshold == nil {
		return nil, false
	}
	return o.ServiceRestEventOutgoingConnectionCountThreshold, true
}

// HasServiceRestEventOutgoingConnectionCountThreshold returns a boolean if a field has been set.
func (o *Broker) HasServiceRestEventOutgoingConnectionCountThreshold() bool {
	if o != nil && o.ServiceRestEventOutgoingConnectionCountThreshold != nil {
		return true
	}

	return false
}

// SetServiceRestEventOutgoingConnectionCountThreshold gets a reference to the given EventThreshold and assigns it to the ServiceRestEventOutgoingConnectionCountThreshold field.
func (o *Broker) SetServiceRestEventOutgoingConnectionCountThreshold(v EventThreshold) {
	o.ServiceRestEventOutgoingConnectionCountThreshold = &v
}

// GetServiceRestIncomingEnabled returns the ServiceRestIncomingEnabled field value if set, zero value otherwise.
func (o *Broker) GetServiceRestIncomingEnabled() bool {
	if o == nil || o.ServiceRestIncomingEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceRestIncomingEnabled
}

// GetServiceRestIncomingEnabledOk returns a tuple with the ServiceRestIncomingEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceRestIncomingEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceRestIncomingEnabled == nil {
		return nil, false
	}
	return o.ServiceRestIncomingEnabled, true
}

// HasServiceRestIncomingEnabled returns a boolean if a field has been set.
func (o *Broker) HasServiceRestIncomingEnabled() bool {
	if o != nil && o.ServiceRestIncomingEnabled != nil {
		return true
	}

	return false
}

// SetServiceRestIncomingEnabled gets a reference to the given bool and assigns it to the ServiceRestIncomingEnabled field.
func (o *Broker) SetServiceRestIncomingEnabled(v bool) {
	o.ServiceRestIncomingEnabled = &v
}

// GetServiceRestOutgoingEnabled returns the ServiceRestOutgoingEnabled field value if set, zero value otherwise.
func (o *Broker) GetServiceRestOutgoingEnabled() bool {
	if o == nil || o.ServiceRestOutgoingEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceRestOutgoingEnabled
}

// GetServiceRestOutgoingEnabledOk returns a tuple with the ServiceRestOutgoingEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceRestOutgoingEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceRestOutgoingEnabled == nil {
		return nil, false
	}
	return o.ServiceRestOutgoingEnabled, true
}

// HasServiceRestOutgoingEnabled returns a boolean if a field has been set.
func (o *Broker) HasServiceRestOutgoingEnabled() bool {
	if o != nil && o.ServiceRestOutgoingEnabled != nil {
		return true
	}

	return false
}

// SetServiceRestOutgoingEnabled gets a reference to the given bool and assigns it to the ServiceRestOutgoingEnabled field.
func (o *Broker) SetServiceRestOutgoingEnabled(v bool) {
	o.ServiceRestOutgoingEnabled = &v
}

// GetServiceSempCorsAllowAnyHostEnabled returns the ServiceSempCorsAllowAnyHostEnabled field value if set, zero value otherwise.
func (o *Broker) GetServiceSempCorsAllowAnyHostEnabled() bool {
	if o == nil || o.ServiceSempCorsAllowAnyHostEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceSempCorsAllowAnyHostEnabled
}

// GetServiceSempCorsAllowAnyHostEnabledOk returns a tuple with the ServiceSempCorsAllowAnyHostEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceSempCorsAllowAnyHostEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceSempCorsAllowAnyHostEnabled == nil {
		return nil, false
	}
	return o.ServiceSempCorsAllowAnyHostEnabled, true
}

// HasServiceSempCorsAllowAnyHostEnabled returns a boolean if a field has been set.
func (o *Broker) HasServiceSempCorsAllowAnyHostEnabled() bool {
	if o != nil && o.ServiceSempCorsAllowAnyHostEnabled != nil {
		return true
	}

	return false
}

// SetServiceSempCorsAllowAnyHostEnabled gets a reference to the given bool and assigns it to the ServiceSempCorsAllowAnyHostEnabled field.
func (o *Broker) SetServiceSempCorsAllowAnyHostEnabled(v bool) {
	o.ServiceSempCorsAllowAnyHostEnabled = &v
}

// GetServiceSempLegacyTimeoutEnabled returns the ServiceSempLegacyTimeoutEnabled field value if set, zero value otherwise.
func (o *Broker) GetServiceSempLegacyTimeoutEnabled() bool {
	if o == nil || o.ServiceSempLegacyTimeoutEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceSempLegacyTimeoutEnabled
}

// GetServiceSempLegacyTimeoutEnabledOk returns a tuple with the ServiceSempLegacyTimeoutEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceSempLegacyTimeoutEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceSempLegacyTimeoutEnabled == nil {
		return nil, false
	}
	return o.ServiceSempLegacyTimeoutEnabled, true
}

// HasServiceSempLegacyTimeoutEnabled returns a boolean if a field has been set.
func (o *Broker) HasServiceSempLegacyTimeoutEnabled() bool {
	if o != nil && o.ServiceSempLegacyTimeoutEnabled != nil {
		return true
	}

	return false
}

// SetServiceSempLegacyTimeoutEnabled gets a reference to the given bool and assigns it to the ServiceSempLegacyTimeoutEnabled field.
func (o *Broker) SetServiceSempLegacyTimeoutEnabled(v bool) {
	o.ServiceSempLegacyTimeoutEnabled = &v
}

// GetServiceSempPlainTextEnabled returns the ServiceSempPlainTextEnabled field value if set, zero value otherwise.
func (o *Broker) GetServiceSempPlainTextEnabled() bool {
	if o == nil || o.ServiceSempPlainTextEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceSempPlainTextEnabled
}

// GetServiceSempPlainTextEnabledOk returns a tuple with the ServiceSempPlainTextEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceSempPlainTextEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceSempPlainTextEnabled == nil {
		return nil, false
	}
	return o.ServiceSempPlainTextEnabled, true
}

// HasServiceSempPlainTextEnabled returns a boolean if a field has been set.
func (o *Broker) HasServiceSempPlainTextEnabled() bool {
	if o != nil && o.ServiceSempPlainTextEnabled != nil {
		return true
	}

	return false
}

// SetServiceSempPlainTextEnabled gets a reference to the given bool and assigns it to the ServiceSempPlainTextEnabled field.
func (o *Broker) SetServiceSempPlainTextEnabled(v bool) {
	o.ServiceSempPlainTextEnabled = &v
}

// GetServiceSempPlainTextListenPort returns the ServiceSempPlainTextListenPort field value if set, zero value otherwise.
func (o *Broker) GetServiceSempPlainTextListenPort() int64 {
	if o == nil || o.ServiceSempPlainTextListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceSempPlainTextListenPort
}

// GetServiceSempPlainTextListenPortOk returns a tuple with the ServiceSempPlainTextListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceSempPlainTextListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceSempPlainTextListenPort == nil {
		return nil, false
	}
	return o.ServiceSempPlainTextListenPort, true
}

// HasServiceSempPlainTextListenPort returns a boolean if a field has been set.
func (o *Broker) HasServiceSempPlainTextListenPort() bool {
	if o != nil && o.ServiceSempPlainTextListenPort != nil {
		return true
	}

	return false
}

// SetServiceSempPlainTextListenPort gets a reference to the given int64 and assigns it to the ServiceSempPlainTextListenPort field.
func (o *Broker) SetServiceSempPlainTextListenPort(v int64) {
	o.ServiceSempPlainTextListenPort = &v
}

// GetServiceSempSessionIdleTimeout returns the ServiceSempSessionIdleTimeout field value if set, zero value otherwise.
func (o *Broker) GetServiceSempSessionIdleTimeout() int32 {
	if o == nil || o.ServiceSempSessionIdleTimeout == nil {
		var ret int32
		return ret
	}
	return *o.ServiceSempSessionIdleTimeout
}

// GetServiceSempSessionIdleTimeoutOk returns a tuple with the ServiceSempSessionIdleTimeout field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceSempSessionIdleTimeoutOk() (*int32, bool) {
	if o == nil || o.ServiceSempSessionIdleTimeout == nil {
		return nil, false
	}
	return o.ServiceSempSessionIdleTimeout, true
}

// HasServiceSempSessionIdleTimeout returns a boolean if a field has been set.
func (o *Broker) HasServiceSempSessionIdleTimeout() bool {
	if o != nil && o.ServiceSempSessionIdleTimeout != nil {
		return true
	}

	return false
}

// SetServiceSempSessionIdleTimeout gets a reference to the given int32 and assigns it to the ServiceSempSessionIdleTimeout field.
func (o *Broker) SetServiceSempSessionIdleTimeout(v int32) {
	o.ServiceSempSessionIdleTimeout = &v
}

// GetServiceSempSessionMaxLifetime returns the ServiceSempSessionMaxLifetime field value if set, zero value otherwise.
func (o *Broker) GetServiceSempSessionMaxLifetime() int32 {
	if o == nil || o.ServiceSempSessionMaxLifetime == nil {
		var ret int32
		return ret
	}
	return *o.ServiceSempSessionMaxLifetime
}

// GetServiceSempSessionMaxLifetimeOk returns a tuple with the ServiceSempSessionMaxLifetime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceSempSessionMaxLifetimeOk() (*int32, bool) {
	if o == nil || o.ServiceSempSessionMaxLifetime == nil {
		return nil, false
	}
	return o.ServiceSempSessionMaxLifetime, true
}

// HasServiceSempSessionMaxLifetime returns a boolean if a field has been set.
func (o *Broker) HasServiceSempSessionMaxLifetime() bool {
	if o != nil && o.ServiceSempSessionMaxLifetime != nil {
		return true
	}

	return false
}

// SetServiceSempSessionMaxLifetime gets a reference to the given int32 and assigns it to the ServiceSempSessionMaxLifetime field.
func (o *Broker) SetServiceSempSessionMaxLifetime(v int32) {
	o.ServiceSempSessionMaxLifetime = &v
}

// GetServiceSempTlsEnabled returns the ServiceSempTlsEnabled field value if set, zero value otherwise.
func (o *Broker) GetServiceSempTlsEnabled() bool {
	if o == nil || o.ServiceSempTlsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceSempTlsEnabled
}

// GetServiceSempTlsEnabledOk returns a tuple with the ServiceSempTlsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceSempTlsEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceSempTlsEnabled == nil {
		return nil, false
	}
	return o.ServiceSempTlsEnabled, true
}

// HasServiceSempTlsEnabled returns a boolean if a field has been set.
func (o *Broker) HasServiceSempTlsEnabled() bool {
	if o != nil && o.ServiceSempTlsEnabled != nil {
		return true
	}

	return false
}

// SetServiceSempTlsEnabled gets a reference to the given bool and assigns it to the ServiceSempTlsEnabled field.
func (o *Broker) SetServiceSempTlsEnabled(v bool) {
	o.ServiceSempTlsEnabled = &v
}

// GetServiceSempTlsListenPort returns the ServiceSempTlsListenPort field value if set, zero value otherwise.
func (o *Broker) GetServiceSempTlsListenPort() int64 {
	if o == nil || o.ServiceSempTlsListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceSempTlsListenPort
}

// GetServiceSempTlsListenPortOk returns a tuple with the ServiceSempTlsListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceSempTlsListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceSempTlsListenPort == nil {
		return nil, false
	}
	return o.ServiceSempTlsListenPort, true
}

// HasServiceSempTlsListenPort returns a boolean if a field has been set.
func (o *Broker) HasServiceSempTlsListenPort() bool {
	if o != nil && o.ServiceSempTlsListenPort != nil {
		return true
	}

	return false
}

// SetServiceSempTlsListenPort gets a reference to the given int64 and assigns it to the ServiceSempTlsListenPort field.
func (o *Broker) SetServiceSempTlsListenPort(v int64) {
	o.ServiceSempTlsListenPort = &v
}

// GetServiceSmfCompressionListenPort returns the ServiceSmfCompressionListenPort field value if set, zero value otherwise.
func (o *Broker) GetServiceSmfCompressionListenPort() int64 {
	if o == nil || o.ServiceSmfCompressionListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceSmfCompressionListenPort
}

// GetServiceSmfCompressionListenPortOk returns a tuple with the ServiceSmfCompressionListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceSmfCompressionListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceSmfCompressionListenPort == nil {
		return nil, false
	}
	return o.ServiceSmfCompressionListenPort, true
}

// HasServiceSmfCompressionListenPort returns a boolean if a field has been set.
func (o *Broker) HasServiceSmfCompressionListenPort() bool {
	if o != nil && o.ServiceSmfCompressionListenPort != nil {
		return true
	}

	return false
}

// SetServiceSmfCompressionListenPort gets a reference to the given int64 and assigns it to the ServiceSmfCompressionListenPort field.
func (o *Broker) SetServiceSmfCompressionListenPort(v int64) {
	o.ServiceSmfCompressionListenPort = &v
}

// GetServiceSmfEnabled returns the ServiceSmfEnabled field value if set, zero value otherwise.
func (o *Broker) GetServiceSmfEnabled() bool {
	if o == nil || o.ServiceSmfEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceSmfEnabled
}

// GetServiceSmfEnabledOk returns a tuple with the ServiceSmfEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceSmfEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceSmfEnabled == nil {
		return nil, false
	}
	return o.ServiceSmfEnabled, true
}

// HasServiceSmfEnabled returns a boolean if a field has been set.
func (o *Broker) HasServiceSmfEnabled() bool {
	if o != nil && o.ServiceSmfEnabled != nil {
		return true
	}

	return false
}

// SetServiceSmfEnabled gets a reference to the given bool and assigns it to the ServiceSmfEnabled field.
func (o *Broker) SetServiceSmfEnabled(v bool) {
	o.ServiceSmfEnabled = &v
}

// GetServiceSmfEventConnectionCountThreshold returns the ServiceSmfEventConnectionCountThreshold field value if set, zero value otherwise.
func (o *Broker) GetServiceSmfEventConnectionCountThreshold() EventThreshold {
	if o == nil || o.ServiceSmfEventConnectionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.ServiceSmfEventConnectionCountThreshold
}

// GetServiceSmfEventConnectionCountThresholdOk returns a tuple with the ServiceSmfEventConnectionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceSmfEventConnectionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.ServiceSmfEventConnectionCountThreshold == nil {
		return nil, false
	}
	return o.ServiceSmfEventConnectionCountThreshold, true
}

// HasServiceSmfEventConnectionCountThreshold returns a boolean if a field has been set.
func (o *Broker) HasServiceSmfEventConnectionCountThreshold() bool {
	if o != nil && o.ServiceSmfEventConnectionCountThreshold != nil {
		return true
	}

	return false
}

// SetServiceSmfEventConnectionCountThreshold gets a reference to the given EventThreshold and assigns it to the ServiceSmfEventConnectionCountThreshold field.
func (o *Broker) SetServiceSmfEventConnectionCountThreshold(v EventThreshold) {
	o.ServiceSmfEventConnectionCountThreshold = &v
}

// GetServiceSmfPlainTextListenPort returns the ServiceSmfPlainTextListenPort field value if set, zero value otherwise.
func (o *Broker) GetServiceSmfPlainTextListenPort() int64 {
	if o == nil || o.ServiceSmfPlainTextListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceSmfPlainTextListenPort
}

// GetServiceSmfPlainTextListenPortOk returns a tuple with the ServiceSmfPlainTextListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceSmfPlainTextListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceSmfPlainTextListenPort == nil {
		return nil, false
	}
	return o.ServiceSmfPlainTextListenPort, true
}

// HasServiceSmfPlainTextListenPort returns a boolean if a field has been set.
func (o *Broker) HasServiceSmfPlainTextListenPort() bool {
	if o != nil && o.ServiceSmfPlainTextListenPort != nil {
		return true
	}

	return false
}

// SetServiceSmfPlainTextListenPort gets a reference to the given int64 and assigns it to the ServiceSmfPlainTextListenPort field.
func (o *Broker) SetServiceSmfPlainTextListenPort(v int64) {
	o.ServiceSmfPlainTextListenPort = &v
}

// GetServiceSmfRoutingControlListenPort returns the ServiceSmfRoutingControlListenPort field value if set, zero value otherwise.
func (o *Broker) GetServiceSmfRoutingControlListenPort() int64 {
	if o == nil || o.ServiceSmfRoutingControlListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceSmfRoutingControlListenPort
}

// GetServiceSmfRoutingControlListenPortOk returns a tuple with the ServiceSmfRoutingControlListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceSmfRoutingControlListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceSmfRoutingControlListenPort == nil {
		return nil, false
	}
	return o.ServiceSmfRoutingControlListenPort, true
}

// HasServiceSmfRoutingControlListenPort returns a boolean if a field has been set.
func (o *Broker) HasServiceSmfRoutingControlListenPort() bool {
	if o != nil && o.ServiceSmfRoutingControlListenPort != nil {
		return true
	}

	return false
}

// SetServiceSmfRoutingControlListenPort gets a reference to the given int64 and assigns it to the ServiceSmfRoutingControlListenPort field.
func (o *Broker) SetServiceSmfRoutingControlListenPort(v int64) {
	o.ServiceSmfRoutingControlListenPort = &v
}

// GetServiceSmfTlsListenPort returns the ServiceSmfTlsListenPort field value if set, zero value otherwise.
func (o *Broker) GetServiceSmfTlsListenPort() int64 {
	if o == nil || o.ServiceSmfTlsListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceSmfTlsListenPort
}

// GetServiceSmfTlsListenPortOk returns a tuple with the ServiceSmfTlsListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceSmfTlsListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceSmfTlsListenPort == nil {
		return nil, false
	}
	return o.ServiceSmfTlsListenPort, true
}

// HasServiceSmfTlsListenPort returns a boolean if a field has been set.
func (o *Broker) HasServiceSmfTlsListenPort() bool {
	if o != nil && o.ServiceSmfTlsListenPort != nil {
		return true
	}

	return false
}

// SetServiceSmfTlsListenPort gets a reference to the given int64 and assigns it to the ServiceSmfTlsListenPort field.
func (o *Broker) SetServiceSmfTlsListenPort(v int64) {
	o.ServiceSmfTlsListenPort = &v
}

// GetServiceTlsEventConnectionCountThreshold returns the ServiceTlsEventConnectionCountThreshold field value if set, zero value otherwise.
func (o *Broker) GetServiceTlsEventConnectionCountThreshold() EventThreshold {
	if o == nil || o.ServiceTlsEventConnectionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.ServiceTlsEventConnectionCountThreshold
}

// GetServiceTlsEventConnectionCountThresholdOk returns a tuple with the ServiceTlsEventConnectionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceTlsEventConnectionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.ServiceTlsEventConnectionCountThreshold == nil {
		return nil, false
	}
	return o.ServiceTlsEventConnectionCountThreshold, true
}

// HasServiceTlsEventConnectionCountThreshold returns a boolean if a field has been set.
func (o *Broker) HasServiceTlsEventConnectionCountThreshold() bool {
	if o != nil && o.ServiceTlsEventConnectionCountThreshold != nil {
		return true
	}

	return false
}

// SetServiceTlsEventConnectionCountThreshold gets a reference to the given EventThreshold and assigns it to the ServiceTlsEventConnectionCountThreshold field.
func (o *Broker) SetServiceTlsEventConnectionCountThreshold(v EventThreshold) {
	o.ServiceTlsEventConnectionCountThreshold = &v
}

// GetServiceWebTransportEnabled returns the ServiceWebTransportEnabled field value if set, zero value otherwise.
func (o *Broker) GetServiceWebTransportEnabled() bool {
	if o == nil || o.ServiceWebTransportEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceWebTransportEnabled
}

// GetServiceWebTransportEnabledOk returns a tuple with the ServiceWebTransportEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceWebTransportEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceWebTransportEnabled == nil {
		return nil, false
	}
	return o.ServiceWebTransportEnabled, true
}

// HasServiceWebTransportEnabled returns a boolean if a field has been set.
func (o *Broker) HasServiceWebTransportEnabled() bool {
	if o != nil && o.ServiceWebTransportEnabled != nil {
		return true
	}

	return false
}

// SetServiceWebTransportEnabled gets a reference to the given bool and assigns it to the ServiceWebTransportEnabled field.
func (o *Broker) SetServiceWebTransportEnabled(v bool) {
	o.ServiceWebTransportEnabled = &v
}

// GetServiceWebTransportPlainTextListenPort returns the ServiceWebTransportPlainTextListenPort field value if set, zero value otherwise.
func (o *Broker) GetServiceWebTransportPlainTextListenPort() int64 {
	if o == nil || o.ServiceWebTransportPlainTextListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceWebTransportPlainTextListenPort
}

// GetServiceWebTransportPlainTextListenPortOk returns a tuple with the ServiceWebTransportPlainTextListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceWebTransportPlainTextListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceWebTransportPlainTextListenPort == nil {
		return nil, false
	}
	return o.ServiceWebTransportPlainTextListenPort, true
}

// HasServiceWebTransportPlainTextListenPort returns a boolean if a field has been set.
func (o *Broker) HasServiceWebTransportPlainTextListenPort() bool {
	if o != nil && o.ServiceWebTransportPlainTextListenPort != nil {
		return true
	}

	return false
}

// SetServiceWebTransportPlainTextListenPort gets a reference to the given int64 and assigns it to the ServiceWebTransportPlainTextListenPort field.
func (o *Broker) SetServiceWebTransportPlainTextListenPort(v int64) {
	o.ServiceWebTransportPlainTextListenPort = &v
}

// GetServiceWebTransportTlsListenPort returns the ServiceWebTransportTlsListenPort field value if set, zero value otherwise.
func (o *Broker) GetServiceWebTransportTlsListenPort() int64 {
	if o == nil || o.ServiceWebTransportTlsListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceWebTransportTlsListenPort
}

// GetServiceWebTransportTlsListenPortOk returns a tuple with the ServiceWebTransportTlsListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceWebTransportTlsListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceWebTransportTlsListenPort == nil {
		return nil, false
	}
	return o.ServiceWebTransportTlsListenPort, true
}

// HasServiceWebTransportTlsListenPort returns a boolean if a field has been set.
func (o *Broker) HasServiceWebTransportTlsListenPort() bool {
	if o != nil && o.ServiceWebTransportTlsListenPort != nil {
		return true
	}

	return false
}

// SetServiceWebTransportTlsListenPort gets a reference to the given int64 and assigns it to the ServiceWebTransportTlsListenPort field.
func (o *Broker) SetServiceWebTransportTlsListenPort(v int64) {
	o.ServiceWebTransportTlsListenPort = &v
}

// GetServiceWebTransportWebUrlSuffix returns the ServiceWebTransportWebUrlSuffix field value if set, zero value otherwise.
func (o *Broker) GetServiceWebTransportWebUrlSuffix() string {
	if o == nil || o.ServiceWebTransportWebUrlSuffix == nil {
		var ret string
		return ret
	}
	return *o.ServiceWebTransportWebUrlSuffix
}

// GetServiceWebTransportWebUrlSuffixOk returns a tuple with the ServiceWebTransportWebUrlSuffix field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetServiceWebTransportWebUrlSuffixOk() (*string, bool) {
	if o == nil || o.ServiceWebTransportWebUrlSuffix == nil {
		return nil, false
	}
	return o.ServiceWebTransportWebUrlSuffix, true
}

// HasServiceWebTransportWebUrlSuffix returns a boolean if a field has been set.
func (o *Broker) HasServiceWebTransportWebUrlSuffix() bool {
	if o != nil && o.ServiceWebTransportWebUrlSuffix != nil {
		return true
	}

	return false
}

// SetServiceWebTransportWebUrlSuffix gets a reference to the given string and assigns it to the ServiceWebTransportWebUrlSuffix field.
func (o *Broker) SetServiceWebTransportWebUrlSuffix(v string) {
	o.ServiceWebTransportWebUrlSuffix = &v
}

// GetTlsBlockVersion11Enabled returns the TlsBlockVersion11Enabled field value if set, zero value otherwise.
func (o *Broker) GetTlsBlockVersion11Enabled() bool {
	if o == nil || o.TlsBlockVersion11Enabled == nil {
		var ret bool
		return ret
	}
	return *o.TlsBlockVersion11Enabled
}

// GetTlsBlockVersion11EnabledOk returns a tuple with the TlsBlockVersion11Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetTlsBlockVersion11EnabledOk() (*bool, bool) {
	if o == nil || o.TlsBlockVersion11Enabled == nil {
		return nil, false
	}
	return o.TlsBlockVersion11Enabled, true
}

// HasTlsBlockVersion11Enabled returns a boolean if a field has been set.
func (o *Broker) HasTlsBlockVersion11Enabled() bool {
	if o != nil && o.TlsBlockVersion11Enabled != nil {
		return true
	}

	return false
}

// SetTlsBlockVersion11Enabled gets a reference to the given bool and assigns it to the TlsBlockVersion11Enabled field.
func (o *Broker) SetTlsBlockVersion11Enabled(v bool) {
	o.TlsBlockVersion11Enabled = &v
}

// GetTlsCipherSuiteManagementList returns the TlsCipherSuiteManagementList field value if set, zero value otherwise.
func (o *Broker) GetTlsCipherSuiteManagementList() string {
	if o == nil || o.TlsCipherSuiteManagementList == nil {
		var ret string
		return ret
	}
	return *o.TlsCipherSuiteManagementList
}

// GetTlsCipherSuiteManagementListOk returns a tuple with the TlsCipherSuiteManagementList field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetTlsCipherSuiteManagementListOk() (*string, bool) {
	if o == nil || o.TlsCipherSuiteManagementList == nil {
		return nil, false
	}
	return o.TlsCipherSuiteManagementList, true
}

// HasTlsCipherSuiteManagementList returns a boolean if a field has been set.
func (o *Broker) HasTlsCipherSuiteManagementList() bool {
	if o != nil && o.TlsCipherSuiteManagementList != nil {
		return true
	}

	return false
}

// SetTlsCipherSuiteManagementList gets a reference to the given string and assigns it to the TlsCipherSuiteManagementList field.
func (o *Broker) SetTlsCipherSuiteManagementList(v string) {
	o.TlsCipherSuiteManagementList = &v
}

// GetTlsCipherSuiteMsgBackboneList returns the TlsCipherSuiteMsgBackboneList field value if set, zero value otherwise.
func (o *Broker) GetTlsCipherSuiteMsgBackboneList() string {
	if o == nil || o.TlsCipherSuiteMsgBackboneList == nil {
		var ret string
		return ret
	}
	return *o.TlsCipherSuiteMsgBackboneList
}

// GetTlsCipherSuiteMsgBackboneListOk returns a tuple with the TlsCipherSuiteMsgBackboneList field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetTlsCipherSuiteMsgBackboneListOk() (*string, bool) {
	if o == nil || o.TlsCipherSuiteMsgBackboneList == nil {
		return nil, false
	}
	return o.TlsCipherSuiteMsgBackboneList, true
}

// HasTlsCipherSuiteMsgBackboneList returns a boolean if a field has been set.
func (o *Broker) HasTlsCipherSuiteMsgBackboneList() bool {
	if o != nil && o.TlsCipherSuiteMsgBackboneList != nil {
		return true
	}

	return false
}

// SetTlsCipherSuiteMsgBackboneList gets a reference to the given string and assigns it to the TlsCipherSuiteMsgBackboneList field.
func (o *Broker) SetTlsCipherSuiteMsgBackboneList(v string) {
	o.TlsCipherSuiteMsgBackboneList = &v
}

// GetTlsCipherSuiteSecureShellList returns the TlsCipherSuiteSecureShellList field value if set, zero value otherwise.
func (o *Broker) GetTlsCipherSuiteSecureShellList() string {
	if o == nil || o.TlsCipherSuiteSecureShellList == nil {
		var ret string
		return ret
	}
	return *o.TlsCipherSuiteSecureShellList
}

// GetTlsCipherSuiteSecureShellListOk returns a tuple with the TlsCipherSuiteSecureShellList field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetTlsCipherSuiteSecureShellListOk() (*string, bool) {
	if o == nil || o.TlsCipherSuiteSecureShellList == nil {
		return nil, false
	}
	return o.TlsCipherSuiteSecureShellList, true
}

// HasTlsCipherSuiteSecureShellList returns a boolean if a field has been set.
func (o *Broker) HasTlsCipherSuiteSecureShellList() bool {
	if o != nil && o.TlsCipherSuiteSecureShellList != nil {
		return true
	}

	return false
}

// SetTlsCipherSuiteSecureShellList gets a reference to the given string and assigns it to the TlsCipherSuiteSecureShellList field.
func (o *Broker) SetTlsCipherSuiteSecureShellList(v string) {
	o.TlsCipherSuiteSecureShellList = &v
}

// GetTlsCrimeExploitProtectionEnabled returns the TlsCrimeExploitProtectionEnabled field value if set, zero value otherwise.
func (o *Broker) GetTlsCrimeExploitProtectionEnabled() bool {
	if o == nil || o.TlsCrimeExploitProtectionEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TlsCrimeExploitProtectionEnabled
}

// GetTlsCrimeExploitProtectionEnabledOk returns a tuple with the TlsCrimeExploitProtectionEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetTlsCrimeExploitProtectionEnabledOk() (*bool, bool) {
	if o == nil || o.TlsCrimeExploitProtectionEnabled == nil {
		return nil, false
	}
	return o.TlsCrimeExploitProtectionEnabled, true
}

// HasTlsCrimeExploitProtectionEnabled returns a boolean if a field has been set.
func (o *Broker) HasTlsCrimeExploitProtectionEnabled() bool {
	if o != nil && o.TlsCrimeExploitProtectionEnabled != nil {
		return true
	}

	return false
}

// SetTlsCrimeExploitProtectionEnabled gets a reference to the given bool and assigns it to the TlsCrimeExploitProtectionEnabled field.
func (o *Broker) SetTlsCrimeExploitProtectionEnabled(v bool) {
	o.TlsCrimeExploitProtectionEnabled = &v
}

// GetTlsServerCertContent returns the TlsServerCertContent field value if set, zero value otherwise.
func (o *Broker) GetTlsServerCertContent() string {
	if o == nil || o.TlsServerCertContent == nil {
		var ret string
		return ret
	}
	return *o.TlsServerCertContent
}

// GetTlsServerCertContentOk returns a tuple with the TlsServerCertContent field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetTlsServerCertContentOk() (*string, bool) {
	if o == nil || o.TlsServerCertContent == nil {
		return nil, false
	}
	return o.TlsServerCertContent, true
}

// HasTlsServerCertContent returns a boolean if a field has been set.
func (o *Broker) HasTlsServerCertContent() bool {
	if o != nil && o.TlsServerCertContent != nil {
		return true
	}

	return false
}

// SetTlsServerCertContent gets a reference to the given string and assigns it to the TlsServerCertContent field.
func (o *Broker) SetTlsServerCertContent(v string) {
	o.TlsServerCertContent = &v
}

// GetTlsServerCertPassword returns the TlsServerCertPassword field value if set, zero value otherwise.
func (o *Broker) GetTlsServerCertPassword() string {
	if o == nil || o.TlsServerCertPassword == nil {
		var ret string
		return ret
	}
	return *o.TlsServerCertPassword
}

// GetTlsServerCertPasswordOk returns a tuple with the TlsServerCertPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetTlsServerCertPasswordOk() (*string, bool) {
	if o == nil || o.TlsServerCertPassword == nil {
		return nil, false
	}
	return o.TlsServerCertPassword, true
}

// HasTlsServerCertPassword returns a boolean if a field has been set.
func (o *Broker) HasTlsServerCertPassword() bool {
	if o != nil && o.TlsServerCertPassword != nil {
		return true
	}

	return false
}

// SetTlsServerCertPassword gets a reference to the given string and assigns it to the TlsServerCertPassword field.
func (o *Broker) SetTlsServerCertPassword(v string) {
	o.TlsServerCertPassword = &v
}

// GetTlsStandardDomainCertificateAuthoritiesEnabled returns the TlsStandardDomainCertificateAuthoritiesEnabled field value if set, zero value otherwise.
func (o *Broker) GetTlsStandardDomainCertificateAuthoritiesEnabled() bool {
	if o == nil || o.TlsStandardDomainCertificateAuthoritiesEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TlsStandardDomainCertificateAuthoritiesEnabled
}

// GetTlsStandardDomainCertificateAuthoritiesEnabledOk returns a tuple with the TlsStandardDomainCertificateAuthoritiesEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetTlsStandardDomainCertificateAuthoritiesEnabledOk() (*bool, bool) {
	if o == nil || o.TlsStandardDomainCertificateAuthoritiesEnabled == nil {
		return nil, false
	}
	return o.TlsStandardDomainCertificateAuthoritiesEnabled, true
}

// HasTlsStandardDomainCertificateAuthoritiesEnabled returns a boolean if a field has been set.
func (o *Broker) HasTlsStandardDomainCertificateAuthoritiesEnabled() bool {
	if o != nil && o.TlsStandardDomainCertificateAuthoritiesEnabled != nil {
		return true
	}

	return false
}

// SetTlsStandardDomainCertificateAuthoritiesEnabled gets a reference to the given bool and assigns it to the TlsStandardDomainCertificateAuthoritiesEnabled field.
func (o *Broker) SetTlsStandardDomainCertificateAuthoritiesEnabled(v bool) {
	o.TlsStandardDomainCertificateAuthoritiesEnabled = &v
}

// GetTlsTicketLifetime returns the TlsTicketLifetime field value if set, zero value otherwise.
func (o *Broker) GetTlsTicketLifetime() int32 {
	if o == nil || o.TlsTicketLifetime == nil {
		var ret int32
		return ret
	}
	return *o.TlsTicketLifetime
}

// GetTlsTicketLifetimeOk returns a tuple with the TlsTicketLifetime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetTlsTicketLifetimeOk() (*int32, bool) {
	if o == nil || o.TlsTicketLifetime == nil {
		return nil, false
	}
	return o.TlsTicketLifetime, true
}

// HasTlsTicketLifetime returns a boolean if a field has been set.
func (o *Broker) HasTlsTicketLifetime() bool {
	if o != nil && o.TlsTicketLifetime != nil {
		return true
	}

	return false
}

// SetTlsTicketLifetime gets a reference to the given int32 and assigns it to the TlsTicketLifetime field.
func (o *Broker) SetTlsTicketLifetime(v int32) {
	o.TlsTicketLifetime = &v
}

// GetWebManagerCustomization returns the WebManagerCustomization field value if set, zero value otherwise.
func (o *Broker) GetWebManagerCustomization() string {
	if o == nil || o.WebManagerCustomization == nil {
		var ret string
		return ret
	}
	return *o.WebManagerCustomization
}

// GetWebManagerCustomizationOk returns a tuple with the WebManagerCustomization field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetWebManagerCustomizationOk() (*string, bool) {
	if o == nil || o.WebManagerCustomization == nil {
		return nil, false
	}
	return o.WebManagerCustomization, true
}

// HasWebManagerCustomization returns a boolean if a field has been set.
func (o *Broker) HasWebManagerCustomization() bool {
	if o != nil && o.WebManagerCustomization != nil {
		return true
	}

	return false
}

// SetWebManagerCustomization gets a reference to the given string and assigns it to the WebManagerCustomization field.
func (o *Broker) SetWebManagerCustomization(v string) {
	o.WebManagerCustomization = &v
}

// GetWebManagerRedirectHttpEnabled returns the WebManagerRedirectHttpEnabled field value if set, zero value otherwise.
func (o *Broker) GetWebManagerRedirectHttpEnabled() bool {
	if o == nil || o.WebManagerRedirectHttpEnabled == nil {
		var ret bool
		return ret
	}
	return *o.WebManagerRedirectHttpEnabled
}

// GetWebManagerRedirectHttpEnabledOk returns a tuple with the WebManagerRedirectHttpEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetWebManagerRedirectHttpEnabledOk() (*bool, bool) {
	if o == nil || o.WebManagerRedirectHttpEnabled == nil {
		return nil, false
	}
	return o.WebManagerRedirectHttpEnabled, true
}

// HasWebManagerRedirectHttpEnabled returns a boolean if a field has been set.
func (o *Broker) HasWebManagerRedirectHttpEnabled() bool {
	if o != nil && o.WebManagerRedirectHttpEnabled != nil {
		return true
	}

	return false
}

// SetWebManagerRedirectHttpEnabled gets a reference to the given bool and assigns it to the WebManagerRedirectHttpEnabled field.
func (o *Broker) SetWebManagerRedirectHttpEnabled(v bool) {
	o.WebManagerRedirectHttpEnabled = &v
}

// GetWebManagerRedirectHttpOverrideTlsPort returns the WebManagerRedirectHttpOverrideTlsPort field value if set, zero value otherwise.
func (o *Broker) GetWebManagerRedirectHttpOverrideTlsPort() int32 {
	if o == nil || o.WebManagerRedirectHttpOverrideTlsPort == nil {
		var ret int32
		return ret
	}
	return *o.WebManagerRedirectHttpOverrideTlsPort
}

// GetWebManagerRedirectHttpOverrideTlsPortOk returns a tuple with the WebManagerRedirectHttpOverrideTlsPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Broker) GetWebManagerRedirectHttpOverrideTlsPortOk() (*int32, bool) {
	if o == nil || o.WebManagerRedirectHttpOverrideTlsPort == nil {
		return nil, false
	}
	return o.WebManagerRedirectHttpOverrideTlsPort, true
}

// HasWebManagerRedirectHttpOverrideTlsPort returns a boolean if a field has been set.
func (o *Broker) HasWebManagerRedirectHttpOverrideTlsPort() bool {
	if o != nil && o.WebManagerRedirectHttpOverrideTlsPort != nil {
		return true
	}

	return false
}

// SetWebManagerRedirectHttpOverrideTlsPort gets a reference to the given int32 and assigns it to the WebManagerRedirectHttpOverrideTlsPort field.
func (o *Broker) SetWebManagerRedirectHttpOverrideTlsPort(v int32) {
	o.WebManagerRedirectHttpOverrideTlsPort = &v
}

func (o Broker) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AuthClientCertRevocationCheckMode != nil {
		toSerialize["authClientCertRevocationCheckMode"] = o.AuthClientCertRevocationCheckMode
	}
	if o.ConfigSyncAuthenticationClientCertMaxChainDepth != nil {
		toSerialize["configSyncAuthenticationClientCertMaxChainDepth"] = o.ConfigSyncAuthenticationClientCertMaxChainDepth
	}
	if o.ConfigSyncAuthenticationClientCertValidateDateEnabled != nil {
		toSerialize["configSyncAuthenticationClientCertValidateDateEnabled"] = o.ConfigSyncAuthenticationClientCertValidateDateEnabled
	}
	if o.ConfigSyncClientProfileTcpInitialCongestionWindow != nil {
		toSerialize["configSyncClientProfileTcpInitialCongestionWindow"] = o.ConfigSyncClientProfileTcpInitialCongestionWindow
	}
	if o.ConfigSyncClientProfileTcpKeepaliveCount != nil {
		toSerialize["configSyncClientProfileTcpKeepaliveCount"] = o.ConfigSyncClientProfileTcpKeepaliveCount
	}
	if o.ConfigSyncClientProfileTcpKeepaliveIdle != nil {
		toSerialize["configSyncClientProfileTcpKeepaliveIdle"] = o.ConfigSyncClientProfileTcpKeepaliveIdle
	}
	if o.ConfigSyncClientProfileTcpKeepaliveInterval != nil {
		toSerialize["configSyncClientProfileTcpKeepaliveInterval"] = o.ConfigSyncClientProfileTcpKeepaliveInterval
	}
	if o.ConfigSyncClientProfileTcpMaxWindow != nil {
		toSerialize["configSyncClientProfileTcpMaxWindow"] = o.ConfigSyncClientProfileTcpMaxWindow
	}
	if o.ConfigSyncClientProfileTcpMss != nil {
		toSerialize["configSyncClientProfileTcpMss"] = o.ConfigSyncClientProfileTcpMss
	}
	if o.ConfigSyncEnabled != nil {
		toSerialize["configSyncEnabled"] = o.ConfigSyncEnabled
	}
	if o.ConfigSyncSynchronizeUsernameEnabled != nil {
		toSerialize["configSyncSynchronizeUsernameEnabled"] = o.ConfigSyncSynchronizeUsernameEnabled
	}
	if o.ConfigSyncTlsEnabled != nil {
		toSerialize["configSyncTlsEnabled"] = o.ConfigSyncTlsEnabled
	}
	if o.GuaranteedMsgingDefragmentationScheduleDayList != nil {
		toSerialize["guaranteedMsgingDefragmentationScheduleDayList"] = o.GuaranteedMsgingDefragmentationScheduleDayList
	}
	if o.GuaranteedMsgingDefragmentationScheduleEnabled != nil {
		toSerialize["guaranteedMsgingDefragmentationScheduleEnabled"] = o.GuaranteedMsgingDefragmentationScheduleEnabled
	}
	if o.GuaranteedMsgingDefragmentationScheduleTimeList != nil {
		toSerialize["guaranteedMsgingDefragmentationScheduleTimeList"] = o.GuaranteedMsgingDefragmentationScheduleTimeList
	}
	if o.GuaranteedMsgingDefragmentationThresholdEnabled != nil {
		toSerialize["guaranteedMsgingDefragmentationThresholdEnabled"] = o.GuaranteedMsgingDefragmentationThresholdEnabled
	}
	if o.GuaranteedMsgingDefragmentationThresholdFragmentationPercentage != nil {
		toSerialize["guaranteedMsgingDefragmentationThresholdFragmentationPercentage"] = o.GuaranteedMsgingDefragmentationThresholdFragmentationPercentage
	}
	if o.GuaranteedMsgingDefragmentationThresholdMinInterval != nil {
		toSerialize["guaranteedMsgingDefragmentationThresholdMinInterval"] = o.GuaranteedMsgingDefragmentationThresholdMinInterval
	}
	if o.GuaranteedMsgingDefragmentationThresholdUsagePercentage != nil {
		toSerialize["guaranteedMsgingDefragmentationThresholdUsagePercentage"] = o.GuaranteedMsgingDefragmentationThresholdUsagePercentage
	}
	if o.GuaranteedMsgingEnabled != nil {
		toSerialize["guaranteedMsgingEnabled"] = o.GuaranteedMsgingEnabled
	}
	if o.GuaranteedMsgingEventCacheUsageThreshold != nil {
		toSerialize["guaranteedMsgingEventCacheUsageThreshold"] = o.GuaranteedMsgingEventCacheUsageThreshold
	}
	if o.GuaranteedMsgingEventDeliveredUnackedThreshold != nil {
		toSerialize["guaranteedMsgingEventDeliveredUnackedThreshold"] = o.GuaranteedMsgingEventDeliveredUnackedThreshold
	}
	if o.GuaranteedMsgingEventDiskUsageThreshold != nil {
		toSerialize["guaranteedMsgingEventDiskUsageThreshold"] = o.GuaranteedMsgingEventDiskUsageThreshold
	}
	if o.GuaranteedMsgingEventEgressFlowCountThreshold != nil {
		toSerialize["guaranteedMsgingEventEgressFlowCountThreshold"] = o.GuaranteedMsgingEventEgressFlowCountThreshold
	}
	if o.GuaranteedMsgingEventEndpointCountThreshold != nil {
		toSerialize["guaranteedMsgingEventEndpointCountThreshold"] = o.GuaranteedMsgingEventEndpointCountThreshold
	}
	if o.GuaranteedMsgingEventIngressFlowCountThreshold != nil {
		toSerialize["guaranteedMsgingEventIngressFlowCountThreshold"] = o.GuaranteedMsgingEventIngressFlowCountThreshold
	}
	if o.GuaranteedMsgingEventMsgCountThreshold != nil {
		toSerialize["guaranteedMsgingEventMsgCountThreshold"] = o.GuaranteedMsgingEventMsgCountThreshold
	}
	if o.GuaranteedMsgingEventMsgSpoolFileCountThreshold != nil {
		toSerialize["guaranteedMsgingEventMsgSpoolFileCountThreshold"] = o.GuaranteedMsgingEventMsgSpoolFileCountThreshold
	}
	if o.GuaranteedMsgingEventMsgSpoolUsageThreshold != nil {
		toSerialize["guaranteedMsgingEventMsgSpoolUsageThreshold"] = o.GuaranteedMsgingEventMsgSpoolUsageThreshold
	}
	if o.GuaranteedMsgingEventTransactedSessionCountThreshold != nil {
		toSerialize["guaranteedMsgingEventTransactedSessionCountThreshold"] = o.GuaranteedMsgingEventTransactedSessionCountThreshold
	}
	if o.GuaranteedMsgingEventTransactedSessionResourceCountThreshold != nil {
		toSerialize["guaranteedMsgingEventTransactedSessionResourceCountThreshold"] = o.GuaranteedMsgingEventTransactedSessionResourceCountThreshold
	}
	if o.GuaranteedMsgingEventTransactionCountThreshold != nil {
		toSerialize["guaranteedMsgingEventTransactionCountThreshold"] = o.GuaranteedMsgingEventTransactionCountThreshold
	}
	if o.GuaranteedMsgingMaxCacheUsage != nil {
		toSerialize["guaranteedMsgingMaxCacheUsage"] = o.GuaranteedMsgingMaxCacheUsage
	}
	if o.GuaranteedMsgingMaxMsgSpoolUsage != nil {
		toSerialize["guaranteedMsgingMaxMsgSpoolUsage"] = o.GuaranteedMsgingMaxMsgSpoolUsage
	}
	if o.GuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout != nil {
		toSerialize["guaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout"] = o.GuaranteedMsgingMsgSpoolSyncMirroredMsgAckTimeout
	}
	if o.GuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout != nil {
		toSerialize["guaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout"] = o.GuaranteedMsgingMsgSpoolSyncMirroredSpoolFileAckTimeout
	}
	if o.GuaranteedMsgingTransactionReplicationCompatibilityMode != nil {
		toSerialize["guaranteedMsgingTransactionReplicationCompatibilityMode"] = o.GuaranteedMsgingTransactionReplicationCompatibilityMode
	}
	if o.OauthProfileDefault != nil {
		toSerialize["oauthProfileDefault"] = o.OauthProfileDefault
	}
	if o.ServiceAmqpEnabled != nil {
		toSerialize["serviceAmqpEnabled"] = o.ServiceAmqpEnabled
	}
	if o.ServiceAmqpTlsListenPort != nil {
		toSerialize["serviceAmqpTlsListenPort"] = o.ServiceAmqpTlsListenPort
	}
	if o.ServiceEventConnectionCountThreshold != nil {
		toSerialize["serviceEventConnectionCountThreshold"] = o.ServiceEventConnectionCountThreshold
	}
	if o.ServiceHealthCheckEnabled != nil {
		toSerialize["serviceHealthCheckEnabled"] = o.ServiceHealthCheckEnabled
	}
	if o.ServiceHealthCheckListenPort != nil {
		toSerialize["serviceHealthCheckListenPort"] = o.ServiceHealthCheckListenPort
	}
	if o.ServiceMateLinkEnabled != nil {
		toSerialize["serviceMateLinkEnabled"] = o.ServiceMateLinkEnabled
	}
	if o.ServiceMateLinkListenPort != nil {
		toSerialize["serviceMateLinkListenPort"] = o.ServiceMateLinkListenPort
	}
	if o.ServiceMqttEnabled != nil {
		toSerialize["serviceMqttEnabled"] = o.ServiceMqttEnabled
	}
	if o.ServiceMsgBackboneEnabled != nil {
		toSerialize["serviceMsgBackboneEnabled"] = o.ServiceMsgBackboneEnabled
	}
	if o.ServiceRedundancyEnabled != nil {
		toSerialize["serviceRedundancyEnabled"] = o.ServiceRedundancyEnabled
	}
	if o.ServiceRedundancyFirstListenPort != nil {
		toSerialize["serviceRedundancyFirstListenPort"] = o.ServiceRedundancyFirstListenPort
	}
	if o.ServiceRestEventOutgoingConnectionCountThreshold != nil {
		toSerialize["serviceRestEventOutgoingConnectionCountThreshold"] = o.ServiceRestEventOutgoingConnectionCountThreshold
	}
	if o.ServiceRestIncomingEnabled != nil {
		toSerialize["serviceRestIncomingEnabled"] = o.ServiceRestIncomingEnabled
	}
	if o.ServiceRestOutgoingEnabled != nil {
		toSerialize["serviceRestOutgoingEnabled"] = o.ServiceRestOutgoingEnabled
	}
	if o.ServiceSempCorsAllowAnyHostEnabled != nil {
		toSerialize["serviceSempCorsAllowAnyHostEnabled"] = o.ServiceSempCorsAllowAnyHostEnabled
	}
	if o.ServiceSempLegacyTimeoutEnabled != nil {
		toSerialize["serviceSempLegacyTimeoutEnabled"] = o.ServiceSempLegacyTimeoutEnabled
	}
	if o.ServiceSempPlainTextEnabled != nil {
		toSerialize["serviceSempPlainTextEnabled"] = o.ServiceSempPlainTextEnabled
	}
	if o.ServiceSempPlainTextListenPort != nil {
		toSerialize["serviceSempPlainTextListenPort"] = o.ServiceSempPlainTextListenPort
	}
	if o.ServiceSempSessionIdleTimeout != nil {
		toSerialize["serviceSempSessionIdleTimeout"] = o.ServiceSempSessionIdleTimeout
	}
	if o.ServiceSempSessionMaxLifetime != nil {
		toSerialize["serviceSempSessionMaxLifetime"] = o.ServiceSempSessionMaxLifetime
	}
	if o.ServiceSempTlsEnabled != nil {
		toSerialize["serviceSempTlsEnabled"] = o.ServiceSempTlsEnabled
	}
	if o.ServiceSempTlsListenPort != nil {
		toSerialize["serviceSempTlsListenPort"] = o.ServiceSempTlsListenPort
	}
	if o.ServiceSmfCompressionListenPort != nil {
		toSerialize["serviceSmfCompressionListenPort"] = o.ServiceSmfCompressionListenPort
	}
	if o.ServiceSmfEnabled != nil {
		toSerialize["serviceSmfEnabled"] = o.ServiceSmfEnabled
	}
	if o.ServiceSmfEventConnectionCountThreshold != nil {
		toSerialize["serviceSmfEventConnectionCountThreshold"] = o.ServiceSmfEventConnectionCountThreshold
	}
	if o.ServiceSmfPlainTextListenPort != nil {
		toSerialize["serviceSmfPlainTextListenPort"] = o.ServiceSmfPlainTextListenPort
	}
	if o.ServiceSmfRoutingControlListenPort != nil {
		toSerialize["serviceSmfRoutingControlListenPort"] = o.ServiceSmfRoutingControlListenPort
	}
	if o.ServiceSmfTlsListenPort != nil {
		toSerialize["serviceSmfTlsListenPort"] = o.ServiceSmfTlsListenPort
	}
	if o.ServiceTlsEventConnectionCountThreshold != nil {
		toSerialize["serviceTlsEventConnectionCountThreshold"] = o.ServiceTlsEventConnectionCountThreshold
	}
	if o.ServiceWebTransportEnabled != nil {
		toSerialize["serviceWebTransportEnabled"] = o.ServiceWebTransportEnabled
	}
	if o.ServiceWebTransportPlainTextListenPort != nil {
		toSerialize["serviceWebTransportPlainTextListenPort"] = o.ServiceWebTransportPlainTextListenPort
	}
	if o.ServiceWebTransportTlsListenPort != nil {
		toSerialize["serviceWebTransportTlsListenPort"] = o.ServiceWebTransportTlsListenPort
	}
	if o.ServiceWebTransportWebUrlSuffix != nil {
		toSerialize["serviceWebTransportWebUrlSuffix"] = o.ServiceWebTransportWebUrlSuffix
	}
	if o.TlsBlockVersion11Enabled != nil {
		toSerialize["tlsBlockVersion11Enabled"] = o.TlsBlockVersion11Enabled
	}
	if o.TlsCipherSuiteManagementList != nil {
		toSerialize["tlsCipherSuiteManagementList"] = o.TlsCipherSuiteManagementList
	}
	if o.TlsCipherSuiteMsgBackboneList != nil {
		toSerialize["tlsCipherSuiteMsgBackboneList"] = o.TlsCipherSuiteMsgBackboneList
	}
	if o.TlsCipherSuiteSecureShellList != nil {
		toSerialize["tlsCipherSuiteSecureShellList"] = o.TlsCipherSuiteSecureShellList
	}
	if o.TlsCrimeExploitProtectionEnabled != nil {
		toSerialize["tlsCrimeExploitProtectionEnabled"] = o.TlsCrimeExploitProtectionEnabled
	}
	if o.TlsServerCertContent != nil {
		toSerialize["tlsServerCertContent"] = o.TlsServerCertContent
	}
	if o.TlsServerCertPassword != nil {
		toSerialize["tlsServerCertPassword"] = o.TlsServerCertPassword
	}
	if o.TlsStandardDomainCertificateAuthoritiesEnabled != nil {
		toSerialize["tlsStandardDomainCertificateAuthoritiesEnabled"] = o.TlsStandardDomainCertificateAuthoritiesEnabled
	}
	if o.TlsTicketLifetime != nil {
		toSerialize["tlsTicketLifetime"] = o.TlsTicketLifetime
	}
	if o.WebManagerCustomization != nil {
		toSerialize["webManagerCustomization"] = o.WebManagerCustomization
	}
	if o.WebManagerRedirectHttpEnabled != nil {
		toSerialize["webManagerRedirectHttpEnabled"] = o.WebManagerRedirectHttpEnabled
	}
	if o.WebManagerRedirectHttpOverrideTlsPort != nil {
		toSerialize["webManagerRedirectHttpOverrideTlsPort"] = o.WebManagerRedirectHttpOverrideTlsPort
	}
	return json.Marshal(toSerialize)
}

type NullableBroker struct {
	value *Broker
	isSet bool
}

func (v NullableBroker) Get() *Broker {
	return v.value
}

func (v *NullableBroker) Set(val *Broker) {
	v.value = val
	v.isSet = true
}

func (v NullableBroker) IsSet() bool {
	return v.isSet
}

func (v *NullableBroker) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBroker(val *Broker) *NullableBroker {
	return &NullableBroker{value: val, isSet: true}
}

func (v NullableBroker) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBroker) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


