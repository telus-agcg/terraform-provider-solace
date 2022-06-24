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

// MsgVpnClientProfile struct for MsgVpnClientProfile
type MsgVpnClientProfile struct {
	// Enable or disable allowing Bridge clients using the Client Profile to connect. Changing this setting does not affect existing Bridge client connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	AllowBridgeConnectionsEnabled *bool `json:"allowBridgeConnectionsEnabled,omitempty"`
	// Enable or disable allowing clients using the Client Profile to bind to endpoints with the cut-through forwarding delivery mode. Changing this value does not affect existing client connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Deprecated since 2.22. This attribute has been deprecated. Please visit the Solace Product Lifecycle Policy web page for details on deprecated features.
	AllowCutThroughForwardingEnabled *bool `json:"allowCutThroughForwardingEnabled,omitempty"`
	// The types of Queues and Topic Endpoints that clients using the client-profile can create. Changing this value does not affect existing client connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"all\"`. The allowed values and their meaning are:  <pre> \"all\" - Client can create any type of endpoint. \"durable\" - Client can create only durable endpoints. \"non-durable\" - Client can create only non-durable endpoints. </pre>  Available since 2.14.
	AllowGuaranteedEndpointCreateDurability *string `json:"allowGuaranteedEndpointCreateDurability,omitempty"`
	// Enable or disable allowing clients using the Client Profile to create topic endponts or queues. Changing this value does not affect existing client connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	AllowGuaranteedEndpointCreateEnabled *bool `json:"allowGuaranteedEndpointCreateEnabled,omitempty"`
	// Enable or disable allowing clients using the Client Profile to receive guaranteed messages. Changing this setting does not affect existing client connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	AllowGuaranteedMsgReceiveEnabled *bool `json:"allowGuaranteedMsgReceiveEnabled,omitempty"`
	// Enable or disable allowing clients using the Client Profile to send guaranteed messages. Changing this setting does not affect existing client connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	AllowGuaranteedMsgSendEnabled *bool `json:"allowGuaranteedMsgSendEnabled,omitempty"`
	// Enable or disable allowing shared subscriptions. Changing this setting does not affect existing subscriptions. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.11.
	AllowSharedSubscriptionsEnabled *bool `json:"allowSharedSubscriptionsEnabled,omitempty"`
	// Enable or disable allowing clients using the Client Profile to establish transacted sessions. Changing this setting does not affect existing client connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	AllowTransactedSessionsEnabled *bool `json:"allowTransactedSessionsEnabled,omitempty"`
	// The name of a queue to copy settings from when a new queue is created by a client using the Client Profile. The referenced queue must exist in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Deprecated since 2.14. This attribute has been replaced with `apiQueueManagementCopyFromOnCreateTemplateName`.
	ApiQueueManagementCopyFromOnCreateName *string `json:"apiQueueManagementCopyFromOnCreateName,omitempty"`
	// The name of a queue template to copy settings from when a new queue is created by a client using the Client Profile. If the referenced queue template does not exist, queue creation will fail when it tries to resolve this template. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.14.
	ApiQueueManagementCopyFromOnCreateTemplateName *string `json:"apiQueueManagementCopyFromOnCreateTemplateName,omitempty"`
	// The name of a topic endpoint to copy settings from when a new topic endpoint is created by a client using the Client Profile. The referenced topic endpoint must exist in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Deprecated since 2.14. This attribute has been replaced with `apiTopicEndpointManagementCopyFromOnCreateTemplateName`.
	ApiTopicEndpointManagementCopyFromOnCreateName *string `json:"apiTopicEndpointManagementCopyFromOnCreateName,omitempty"`
	// The name of a topic endpoint template to copy settings from when a new topic endpoint is created by a client using the Client Profile. If the referenced topic endpoint template does not exist, topic endpoint creation will fail when it tries to resolve this template. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.14.
	ApiTopicEndpointManagementCopyFromOnCreateTemplateName *string `json:"apiTopicEndpointManagementCopyFromOnCreateTemplateName,omitempty"`
	// The name of the Client Profile.
	ClientProfileName *string `json:"clientProfileName,omitempty"`
	// Enable or disable allowing clients using the Client Profile to use compression. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`. Available since 2.10.
	CompressionEnabled *bool `json:"compressionEnabled,omitempty"`
	// The amount of time to delay the delivery of messages to clients using the Client Profile after the initial message has been delivered (the eliding delay interval), in milliseconds. A value of 0 means there is no delay in delivering messages to clients. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.
	ElidingDelay *int64 `json:"elidingDelay,omitempty"`
	// Enable or disable message eliding for clients using the Client Profile. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	ElidingEnabled *bool `json:"elidingEnabled,omitempty"`
	// The maximum number of topics tracked for message eliding per client connection using the Client Profile. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `256`.
	ElidingMaxTopicCount                                     *int64                   `json:"elidingMaxTopicCount,omitempty"`
	EventClientProvisionedEndpointSpoolUsageThreshold        *EventThresholdByPercent `json:"eventClientProvisionedEndpointSpoolUsageThreshold,omitempty"`
	EventConnectionCountPerClientUsernameThreshold           *EventThreshold          `json:"eventConnectionCountPerClientUsernameThreshold,omitempty"`
	EventEgressFlowCountThreshold                            *EventThreshold          `json:"eventEgressFlowCountThreshold,omitempty"`
	EventEndpointCountPerClientUsernameThreshold             *EventThreshold          `json:"eventEndpointCountPerClientUsernameThreshold,omitempty"`
	EventIngressFlowCountThreshold                           *EventThreshold          `json:"eventIngressFlowCountThreshold,omitempty"`
	EventServiceSmfConnectionCountPerClientUsernameThreshold *EventThreshold          `json:"eventServiceSmfConnectionCountPerClientUsernameThreshold,omitempty"`
	EventServiceWebConnectionCountPerClientUsernameThreshold *EventThreshold          `json:"eventServiceWebConnectionCountPerClientUsernameThreshold,omitempty"`
	EventSubscriptionCountThreshold                          *EventThreshold          `json:"eventSubscriptionCountThreshold,omitempty"`
	EventTransactedSessionCountThreshold                     *EventThreshold          `json:"eventTransactedSessionCountThreshold,omitempty"`
	EventTransactionCountThreshold                           *EventThreshold          `json:"eventTransactionCountThreshold,omitempty"`
	// The maximum number of client connections per Client Username using the Client Profile. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform.
	MaxConnectionCountPerClientUsername *int64 `json:"maxConnectionCountPerClientUsername,omitempty"`
	// The maximum number of transmit flows that can be created by one client using the Client Profile. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.
	MaxEgressFlowCount *int64 `json:"maxEgressFlowCount,omitempty"`
	// The maximum number of queues and topic endpoints that can be created by clients with the same Client Username using the Client Profile. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.
	MaxEndpointCountPerClientUsername *int64 `json:"maxEndpointCountPerClientUsername,omitempty"`
	// The maximum number of receive flows that can be created by one client using the Client Profile. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.
	MaxIngressFlowCount *int64 `json:"maxIngressFlowCount,omitempty"`
	// The maximum number of publisher and consumer messages combined that is allowed within a transaction for each client associated with this client-profile. Exceeding this limit will result in a transaction prepare or commit failure. Changing this value during operation will not affect existing sessions. It is only validated at transaction creation time. Large transactions consume more resources and are more likely to require retrieving messages from the ADB or from disk to process the transaction prepare or commit requests. The transaction processing rate may diminish if a large number of messages must be retrieved from the ADB or from disk. Care should be taken to not use excessively large transactions needlessly to avoid exceeding resource limits and to avoid reducing the overall broker performance. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `256`. Available since 2.20.
	MaxMsgsPerTransaction *int32 `json:"maxMsgsPerTransaction,omitempty"`
	// The maximum number of subscriptions per client using the Client Profile. This limit is not enforced when a client adds a subscription to an endpoint, except for MQTT QoS 1 subscriptions. In addition, this limit is not enforced when a subscription is added using a management interface, such as CLI or SEMP. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.
	MaxSubscriptionCount *int64 `json:"maxSubscriptionCount,omitempty"`
	// The maximum number of transacted sessions that can be created by one client using the Client Profile. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `10`.
	MaxTransactedSessionCount *int64 `json:"maxTransactedSessionCount,omitempty"`
	// The maximum number of transactions that can be created by one client using the Client Profile. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.
	MaxTransactionCount *int64 `json:"maxTransactionCount,omitempty"`
	// The name of the Message VPN.
	MsgVpnName *string `json:"msgVpnName,omitempty"`
	// The maximum depth of the \"Control 1\" (C-1) priority queue, in work units. Each work unit is 2048 bytes of message data. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `20000`.
	QueueControl1MaxDepth *int32 `json:"queueControl1MaxDepth,omitempty"`
	// The number of messages that are always allowed entry into the \"Control 1\" (C-1) priority queue, regardless of the `queueControl1MaxDepth` value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `4`.
	QueueControl1MinMsgBurst *int32 `json:"queueControl1MinMsgBurst,omitempty"`
	// The maximum depth of the \"Direct 1\" (D-1) priority queue, in work units. Each work unit is 2048 bytes of message data. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `20000`.
	QueueDirect1MaxDepth *int32 `json:"queueDirect1MaxDepth,omitempty"`
	// The number of messages that are always allowed entry into the \"Direct 1\" (D-1) priority queue, regardless of the `queueDirect1MaxDepth` value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `4`.
	QueueDirect1MinMsgBurst *int32 `json:"queueDirect1MinMsgBurst,omitempty"`
	// The maximum depth of the \"Direct 2\" (D-2) priority queue, in work units. Each work unit is 2048 bytes of message data. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `20000`.
	QueueDirect2MaxDepth *int32 `json:"queueDirect2MaxDepth,omitempty"`
	// The number of messages that are always allowed entry into the \"Direct 2\" (D-2) priority queue, regardless of the `queueDirect2MaxDepth` value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `4`.
	QueueDirect2MinMsgBurst *int32 `json:"queueDirect2MinMsgBurst,omitempty"`
	// The maximum depth of the \"Direct 3\" (D-3) priority queue, in work units. Each work unit is 2048 bytes of message data. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `20000`.
	QueueDirect3MaxDepth *int32 `json:"queueDirect3MaxDepth,omitempty"`
	// The number of messages that are always allowed entry into the \"Direct 3\" (D-3) priority queue, regardless of the `queueDirect3MaxDepth` value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `4`.
	QueueDirect3MinMsgBurst *int32 `json:"queueDirect3MinMsgBurst,omitempty"`
	// The maximum depth of the \"Guaranteed 1\" (G-1) priority queue, in work units. Each work unit is 2048 bytes of message data. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `20000`.
	QueueGuaranteed1MaxDepth *int32 `json:"queueGuaranteed1MaxDepth,omitempty"`
	// The number of messages that are always allowed entry into the \"Guaranteed 1\" (G-3) priority queue, regardless of the `queueGuaranteed1MaxDepth` value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `255`.
	QueueGuaranteed1MinMsgBurst *int32 `json:"queueGuaranteed1MinMsgBurst,omitempty"`
	// Enable or disable the sending of a negative acknowledgement (NACK) to a client using the Client Profile when discarding a guaranteed message due to no matching subscription found. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.2.
	RejectMsgToSenderOnNoSubscriptionMatchEnabled *bool `json:"rejectMsgToSenderOnNoSubscriptionMatchEnabled,omitempty"`
	// Enable or disable allowing clients using the Client Profile to connect to the Message VPN when its replication state is standby. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	ReplicationAllowClientConnectWhenStandbyEnabled *bool `json:"replicationAllowClientConnectWhenStandbyEnabled,omitempty"`
	// The minimum client keepalive timeout which will be enforced for client connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `30`. Available since 2.19.
	ServiceMinKeepaliveTimeout *int32 `json:"serviceMinKeepaliveTimeout,omitempty"`
	// The maximum number of SMF client connections per Client Username using the Client Profile. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform.
	ServiceSmfMaxConnectionCountPerClientUsername *int64 `json:"serviceSmfMaxConnectionCountPerClientUsername,omitempty"`
	// Enable or disable the enforcement of a minimum keepalive timeout for SMF clients. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.19.
	ServiceSmfMinKeepaliveEnabled *bool `json:"serviceSmfMinKeepaliveEnabled,omitempty"`
	// The timeout for inactive Web Transport client sessions using the Client Profile, in seconds. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `30`.
	ServiceWebInactiveTimeout *int64 `json:"serviceWebInactiveTimeout,omitempty"`
	// The maximum number of Web Transport client connections per Client Username using the Client Profile. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform.
	ServiceWebMaxConnectionCountPerClientUsername *int64 `json:"serviceWebMaxConnectionCountPerClientUsername,omitempty"`
	// The maximum Web Transport payload size before fragmentation occurs for clients using the Client Profile, in bytes. The size of the header is not included. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000000`.
	ServiceWebMaxPayload *int64 `json:"serviceWebMaxPayload,omitempty"`
	// The TCP initial congestion window size for clients using the Client Profile, in multiples of the TCP Maximum Segment Size (MSS). Changing the value from its default of 2 results in non-compliance with RFC 2581. Contact support before changing this value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `2`.
	TcpCongestionWindowSize *int64 `json:"tcpCongestionWindowSize,omitempty"`
	// The number of TCP keepalive retransmissions to a client using the Client Profile before declaring that it is not available. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `5`.
	TcpKeepaliveCount *int64 `json:"tcpKeepaliveCount,omitempty"`
	// The amount of time a client connection using the Client Profile must remain idle before TCP begins sending keepalive probes, in seconds. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.
	TcpKeepaliveIdleTime *int64 `json:"tcpKeepaliveIdleTime,omitempty"`
	// The amount of time between TCP keepalive retransmissions to a client using the Client Profile when no acknowledgement is received, in seconds. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1`.
	TcpKeepaliveInterval *int64 `json:"tcpKeepaliveInterval,omitempty"`
	// The TCP maximum segment size for clients using the Client Profile, in bytes. Changes are applied to all existing connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1460`.
	TcpMaxSegmentSize *int64 `json:"tcpMaxSegmentSize,omitempty"`
	// The TCP maximum window size for clients using the Client Profile, in kilobytes. Changes are applied to all existing connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `256`.
	TcpMaxWindowSize *int64 `json:"tcpMaxWindowSize,omitempty"`
	// Enable or disable allowing a client using the Client Profile to downgrade an encrypted connection to plain text. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`. Available since 2.8.
	TlsAllowDowngradeToPlainTextEnabled *bool `json:"tlsAllowDowngradeToPlainTextEnabled,omitempty"`
}

// NewMsgVpnClientProfile instantiates a new MsgVpnClientProfile object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMsgVpnClientProfile() *MsgVpnClientProfile {
	this := MsgVpnClientProfile{}
	return &this
}

// NewMsgVpnClientProfileWithDefaults instantiates a new MsgVpnClientProfile object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMsgVpnClientProfileWithDefaults() *MsgVpnClientProfile {
	this := MsgVpnClientProfile{}
	return &this
}

// GetAllowBridgeConnectionsEnabled returns the AllowBridgeConnectionsEnabled field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetAllowBridgeConnectionsEnabled() bool {
	if o == nil || o.AllowBridgeConnectionsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AllowBridgeConnectionsEnabled
}

// GetAllowBridgeConnectionsEnabledOk returns a tuple with the AllowBridgeConnectionsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetAllowBridgeConnectionsEnabledOk() (*bool, bool) {
	if o == nil || o.AllowBridgeConnectionsEnabled == nil {
		return nil, false
	}
	return o.AllowBridgeConnectionsEnabled, true
}

// HasAllowBridgeConnectionsEnabled returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasAllowBridgeConnectionsEnabled() bool {
	if o != nil && o.AllowBridgeConnectionsEnabled != nil {
		return true
	}

	return false
}

// SetAllowBridgeConnectionsEnabled gets a reference to the given bool and assigns it to the AllowBridgeConnectionsEnabled field.
func (o *MsgVpnClientProfile) SetAllowBridgeConnectionsEnabled(v bool) {
	o.AllowBridgeConnectionsEnabled = &v
}

// GetAllowCutThroughForwardingEnabled returns the AllowCutThroughForwardingEnabled field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetAllowCutThroughForwardingEnabled() bool {
	if o == nil || o.AllowCutThroughForwardingEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AllowCutThroughForwardingEnabled
}

// GetAllowCutThroughForwardingEnabledOk returns a tuple with the AllowCutThroughForwardingEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetAllowCutThroughForwardingEnabledOk() (*bool, bool) {
	if o == nil || o.AllowCutThroughForwardingEnabled == nil {
		return nil, false
	}
	return o.AllowCutThroughForwardingEnabled, true
}

// HasAllowCutThroughForwardingEnabled returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasAllowCutThroughForwardingEnabled() bool {
	if o != nil && o.AllowCutThroughForwardingEnabled != nil {
		return true
	}

	return false
}

// SetAllowCutThroughForwardingEnabled gets a reference to the given bool and assigns it to the AllowCutThroughForwardingEnabled field.
func (o *MsgVpnClientProfile) SetAllowCutThroughForwardingEnabled(v bool) {
	o.AllowCutThroughForwardingEnabled = &v
}

// GetAllowGuaranteedEndpointCreateDurability returns the AllowGuaranteedEndpointCreateDurability field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetAllowGuaranteedEndpointCreateDurability() string {
	if o == nil || o.AllowGuaranteedEndpointCreateDurability == nil {
		var ret string
		return ret
	}
	return *o.AllowGuaranteedEndpointCreateDurability
}

// GetAllowGuaranteedEndpointCreateDurabilityOk returns a tuple with the AllowGuaranteedEndpointCreateDurability field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetAllowGuaranteedEndpointCreateDurabilityOk() (*string, bool) {
	if o == nil || o.AllowGuaranteedEndpointCreateDurability == nil {
		return nil, false
	}
	return o.AllowGuaranteedEndpointCreateDurability, true
}

// HasAllowGuaranteedEndpointCreateDurability returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasAllowGuaranteedEndpointCreateDurability() bool {
	if o != nil && o.AllowGuaranteedEndpointCreateDurability != nil {
		return true
	}

	return false
}

// SetAllowGuaranteedEndpointCreateDurability gets a reference to the given string and assigns it to the AllowGuaranteedEndpointCreateDurability field.
func (o *MsgVpnClientProfile) SetAllowGuaranteedEndpointCreateDurability(v string) {
	o.AllowGuaranteedEndpointCreateDurability = &v
}

// GetAllowGuaranteedEndpointCreateEnabled returns the AllowGuaranteedEndpointCreateEnabled field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetAllowGuaranteedEndpointCreateEnabled() bool {
	if o == nil || o.AllowGuaranteedEndpointCreateEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AllowGuaranteedEndpointCreateEnabled
}

// GetAllowGuaranteedEndpointCreateEnabledOk returns a tuple with the AllowGuaranteedEndpointCreateEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetAllowGuaranteedEndpointCreateEnabledOk() (*bool, bool) {
	if o == nil || o.AllowGuaranteedEndpointCreateEnabled == nil {
		return nil, false
	}
	return o.AllowGuaranteedEndpointCreateEnabled, true
}

// HasAllowGuaranteedEndpointCreateEnabled returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasAllowGuaranteedEndpointCreateEnabled() bool {
	if o != nil && o.AllowGuaranteedEndpointCreateEnabled != nil {
		return true
	}

	return false
}

// SetAllowGuaranteedEndpointCreateEnabled gets a reference to the given bool and assigns it to the AllowGuaranteedEndpointCreateEnabled field.
func (o *MsgVpnClientProfile) SetAllowGuaranteedEndpointCreateEnabled(v bool) {
	o.AllowGuaranteedEndpointCreateEnabled = &v
}

// GetAllowGuaranteedMsgReceiveEnabled returns the AllowGuaranteedMsgReceiveEnabled field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetAllowGuaranteedMsgReceiveEnabled() bool {
	if o == nil || o.AllowGuaranteedMsgReceiveEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AllowGuaranteedMsgReceiveEnabled
}

// GetAllowGuaranteedMsgReceiveEnabledOk returns a tuple with the AllowGuaranteedMsgReceiveEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetAllowGuaranteedMsgReceiveEnabledOk() (*bool, bool) {
	if o == nil || o.AllowGuaranteedMsgReceiveEnabled == nil {
		return nil, false
	}
	return o.AllowGuaranteedMsgReceiveEnabled, true
}

// HasAllowGuaranteedMsgReceiveEnabled returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasAllowGuaranteedMsgReceiveEnabled() bool {
	if o != nil && o.AllowGuaranteedMsgReceiveEnabled != nil {
		return true
	}

	return false
}

// SetAllowGuaranteedMsgReceiveEnabled gets a reference to the given bool and assigns it to the AllowGuaranteedMsgReceiveEnabled field.
func (o *MsgVpnClientProfile) SetAllowGuaranteedMsgReceiveEnabled(v bool) {
	o.AllowGuaranteedMsgReceiveEnabled = &v
}

// GetAllowGuaranteedMsgSendEnabled returns the AllowGuaranteedMsgSendEnabled field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetAllowGuaranteedMsgSendEnabled() bool {
	if o == nil || o.AllowGuaranteedMsgSendEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AllowGuaranteedMsgSendEnabled
}

// GetAllowGuaranteedMsgSendEnabledOk returns a tuple with the AllowGuaranteedMsgSendEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetAllowGuaranteedMsgSendEnabledOk() (*bool, bool) {
	if o == nil || o.AllowGuaranteedMsgSendEnabled == nil {
		return nil, false
	}
	return o.AllowGuaranteedMsgSendEnabled, true
}

// HasAllowGuaranteedMsgSendEnabled returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasAllowGuaranteedMsgSendEnabled() bool {
	if o != nil && o.AllowGuaranteedMsgSendEnabled != nil {
		return true
	}

	return false
}

// SetAllowGuaranteedMsgSendEnabled gets a reference to the given bool and assigns it to the AllowGuaranteedMsgSendEnabled field.
func (o *MsgVpnClientProfile) SetAllowGuaranteedMsgSendEnabled(v bool) {
	o.AllowGuaranteedMsgSendEnabled = &v
}

// GetAllowSharedSubscriptionsEnabled returns the AllowSharedSubscriptionsEnabled field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetAllowSharedSubscriptionsEnabled() bool {
	if o == nil || o.AllowSharedSubscriptionsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AllowSharedSubscriptionsEnabled
}

// GetAllowSharedSubscriptionsEnabledOk returns a tuple with the AllowSharedSubscriptionsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetAllowSharedSubscriptionsEnabledOk() (*bool, bool) {
	if o == nil || o.AllowSharedSubscriptionsEnabled == nil {
		return nil, false
	}
	return o.AllowSharedSubscriptionsEnabled, true
}

// HasAllowSharedSubscriptionsEnabled returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasAllowSharedSubscriptionsEnabled() bool {
	if o != nil && o.AllowSharedSubscriptionsEnabled != nil {
		return true
	}

	return false
}

// SetAllowSharedSubscriptionsEnabled gets a reference to the given bool and assigns it to the AllowSharedSubscriptionsEnabled field.
func (o *MsgVpnClientProfile) SetAllowSharedSubscriptionsEnabled(v bool) {
	o.AllowSharedSubscriptionsEnabled = &v
}

// GetAllowTransactedSessionsEnabled returns the AllowTransactedSessionsEnabled field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetAllowTransactedSessionsEnabled() bool {
	if o == nil || o.AllowTransactedSessionsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AllowTransactedSessionsEnabled
}

// GetAllowTransactedSessionsEnabledOk returns a tuple with the AllowTransactedSessionsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetAllowTransactedSessionsEnabledOk() (*bool, bool) {
	if o == nil || o.AllowTransactedSessionsEnabled == nil {
		return nil, false
	}
	return o.AllowTransactedSessionsEnabled, true
}

// HasAllowTransactedSessionsEnabled returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasAllowTransactedSessionsEnabled() bool {
	if o != nil && o.AllowTransactedSessionsEnabled != nil {
		return true
	}

	return false
}

// SetAllowTransactedSessionsEnabled gets a reference to the given bool and assigns it to the AllowTransactedSessionsEnabled field.
func (o *MsgVpnClientProfile) SetAllowTransactedSessionsEnabled(v bool) {
	o.AllowTransactedSessionsEnabled = &v
}

// GetApiQueueManagementCopyFromOnCreateName returns the ApiQueueManagementCopyFromOnCreateName field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetApiQueueManagementCopyFromOnCreateName() string {
	if o == nil || o.ApiQueueManagementCopyFromOnCreateName == nil {
		var ret string
		return ret
	}
	return *o.ApiQueueManagementCopyFromOnCreateName
}

// GetApiQueueManagementCopyFromOnCreateNameOk returns a tuple with the ApiQueueManagementCopyFromOnCreateName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetApiQueueManagementCopyFromOnCreateNameOk() (*string, bool) {
	if o == nil || o.ApiQueueManagementCopyFromOnCreateName == nil {
		return nil, false
	}
	return o.ApiQueueManagementCopyFromOnCreateName, true
}

// HasApiQueueManagementCopyFromOnCreateName returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasApiQueueManagementCopyFromOnCreateName() bool {
	if o != nil && o.ApiQueueManagementCopyFromOnCreateName != nil {
		return true
	}

	return false
}

// SetApiQueueManagementCopyFromOnCreateName gets a reference to the given string and assigns it to the ApiQueueManagementCopyFromOnCreateName field.
func (o *MsgVpnClientProfile) SetApiQueueManagementCopyFromOnCreateName(v string) {
	o.ApiQueueManagementCopyFromOnCreateName = &v
}

// GetApiQueueManagementCopyFromOnCreateTemplateName returns the ApiQueueManagementCopyFromOnCreateTemplateName field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetApiQueueManagementCopyFromOnCreateTemplateName() string {
	if o == nil || o.ApiQueueManagementCopyFromOnCreateTemplateName == nil {
		var ret string
		return ret
	}
	return *o.ApiQueueManagementCopyFromOnCreateTemplateName
}

// GetApiQueueManagementCopyFromOnCreateTemplateNameOk returns a tuple with the ApiQueueManagementCopyFromOnCreateTemplateName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetApiQueueManagementCopyFromOnCreateTemplateNameOk() (*string, bool) {
	if o == nil || o.ApiQueueManagementCopyFromOnCreateTemplateName == nil {
		return nil, false
	}
	return o.ApiQueueManagementCopyFromOnCreateTemplateName, true
}

// HasApiQueueManagementCopyFromOnCreateTemplateName returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasApiQueueManagementCopyFromOnCreateTemplateName() bool {
	if o != nil && o.ApiQueueManagementCopyFromOnCreateTemplateName != nil {
		return true
	}

	return false
}

// SetApiQueueManagementCopyFromOnCreateTemplateName gets a reference to the given string and assigns it to the ApiQueueManagementCopyFromOnCreateTemplateName field.
func (o *MsgVpnClientProfile) SetApiQueueManagementCopyFromOnCreateTemplateName(v string) {
	o.ApiQueueManagementCopyFromOnCreateTemplateName = &v
}

// GetApiTopicEndpointManagementCopyFromOnCreateName returns the ApiTopicEndpointManagementCopyFromOnCreateName field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetApiTopicEndpointManagementCopyFromOnCreateName() string {
	if o == nil || o.ApiTopicEndpointManagementCopyFromOnCreateName == nil {
		var ret string
		return ret
	}
	return *o.ApiTopicEndpointManagementCopyFromOnCreateName
}

// GetApiTopicEndpointManagementCopyFromOnCreateNameOk returns a tuple with the ApiTopicEndpointManagementCopyFromOnCreateName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetApiTopicEndpointManagementCopyFromOnCreateNameOk() (*string, bool) {
	if o == nil || o.ApiTopicEndpointManagementCopyFromOnCreateName == nil {
		return nil, false
	}
	return o.ApiTopicEndpointManagementCopyFromOnCreateName, true
}

// HasApiTopicEndpointManagementCopyFromOnCreateName returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasApiTopicEndpointManagementCopyFromOnCreateName() bool {
	if o != nil && o.ApiTopicEndpointManagementCopyFromOnCreateName != nil {
		return true
	}

	return false
}

// SetApiTopicEndpointManagementCopyFromOnCreateName gets a reference to the given string and assigns it to the ApiTopicEndpointManagementCopyFromOnCreateName field.
func (o *MsgVpnClientProfile) SetApiTopicEndpointManagementCopyFromOnCreateName(v string) {
	o.ApiTopicEndpointManagementCopyFromOnCreateName = &v
}

// GetApiTopicEndpointManagementCopyFromOnCreateTemplateName returns the ApiTopicEndpointManagementCopyFromOnCreateTemplateName field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetApiTopicEndpointManagementCopyFromOnCreateTemplateName() string {
	if o == nil || o.ApiTopicEndpointManagementCopyFromOnCreateTemplateName == nil {
		var ret string
		return ret
	}
	return *o.ApiTopicEndpointManagementCopyFromOnCreateTemplateName
}

// GetApiTopicEndpointManagementCopyFromOnCreateTemplateNameOk returns a tuple with the ApiTopicEndpointManagementCopyFromOnCreateTemplateName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetApiTopicEndpointManagementCopyFromOnCreateTemplateNameOk() (*string, bool) {
	if o == nil || o.ApiTopicEndpointManagementCopyFromOnCreateTemplateName == nil {
		return nil, false
	}
	return o.ApiTopicEndpointManagementCopyFromOnCreateTemplateName, true
}

// HasApiTopicEndpointManagementCopyFromOnCreateTemplateName returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasApiTopicEndpointManagementCopyFromOnCreateTemplateName() bool {
	if o != nil && o.ApiTopicEndpointManagementCopyFromOnCreateTemplateName != nil {
		return true
	}

	return false
}

// SetApiTopicEndpointManagementCopyFromOnCreateTemplateName gets a reference to the given string and assigns it to the ApiTopicEndpointManagementCopyFromOnCreateTemplateName field.
func (o *MsgVpnClientProfile) SetApiTopicEndpointManagementCopyFromOnCreateTemplateName(v string) {
	o.ApiTopicEndpointManagementCopyFromOnCreateTemplateName = &v
}

// GetClientProfileName returns the ClientProfileName field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetClientProfileName() string {
	if o == nil || o.ClientProfileName == nil {
		var ret string
		return ret
	}
	return *o.ClientProfileName
}

// GetClientProfileNameOk returns a tuple with the ClientProfileName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetClientProfileNameOk() (*string, bool) {
	if o == nil || o.ClientProfileName == nil {
		return nil, false
	}
	return o.ClientProfileName, true
}

// HasClientProfileName returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasClientProfileName() bool {
	if o != nil && o.ClientProfileName != nil {
		return true
	}

	return false
}

// SetClientProfileName gets a reference to the given string and assigns it to the ClientProfileName field.
func (o *MsgVpnClientProfile) SetClientProfileName(v string) {
	o.ClientProfileName = &v
}

// GetCompressionEnabled returns the CompressionEnabled field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetCompressionEnabled() bool {
	if o == nil || o.CompressionEnabled == nil {
		var ret bool
		return ret
	}
	return *o.CompressionEnabled
}

// GetCompressionEnabledOk returns a tuple with the CompressionEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetCompressionEnabledOk() (*bool, bool) {
	if o == nil || o.CompressionEnabled == nil {
		return nil, false
	}
	return o.CompressionEnabled, true
}

// HasCompressionEnabled returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasCompressionEnabled() bool {
	if o != nil && o.CompressionEnabled != nil {
		return true
	}

	return false
}

// SetCompressionEnabled gets a reference to the given bool and assigns it to the CompressionEnabled field.
func (o *MsgVpnClientProfile) SetCompressionEnabled(v bool) {
	o.CompressionEnabled = &v
}

// GetElidingDelay returns the ElidingDelay field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetElidingDelay() int64 {
	if o == nil || o.ElidingDelay == nil {
		var ret int64
		return ret
	}
	return *o.ElidingDelay
}

// GetElidingDelayOk returns a tuple with the ElidingDelay field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetElidingDelayOk() (*int64, bool) {
	if o == nil || o.ElidingDelay == nil {
		return nil, false
	}
	return o.ElidingDelay, true
}

// HasElidingDelay returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasElidingDelay() bool {
	if o != nil && o.ElidingDelay != nil {
		return true
	}

	return false
}

// SetElidingDelay gets a reference to the given int64 and assigns it to the ElidingDelay field.
func (o *MsgVpnClientProfile) SetElidingDelay(v int64) {
	o.ElidingDelay = &v
}

// GetElidingEnabled returns the ElidingEnabled field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetElidingEnabled() bool {
	if o == nil || o.ElidingEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ElidingEnabled
}

// GetElidingEnabledOk returns a tuple with the ElidingEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetElidingEnabledOk() (*bool, bool) {
	if o == nil || o.ElidingEnabled == nil {
		return nil, false
	}
	return o.ElidingEnabled, true
}

// HasElidingEnabled returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasElidingEnabled() bool {
	if o != nil && o.ElidingEnabled != nil {
		return true
	}

	return false
}

// SetElidingEnabled gets a reference to the given bool and assigns it to the ElidingEnabled field.
func (o *MsgVpnClientProfile) SetElidingEnabled(v bool) {
	o.ElidingEnabled = &v
}

// GetElidingMaxTopicCount returns the ElidingMaxTopicCount field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetElidingMaxTopicCount() int64 {
	if o == nil || o.ElidingMaxTopicCount == nil {
		var ret int64
		return ret
	}
	return *o.ElidingMaxTopicCount
}

// GetElidingMaxTopicCountOk returns a tuple with the ElidingMaxTopicCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetElidingMaxTopicCountOk() (*int64, bool) {
	if o == nil || o.ElidingMaxTopicCount == nil {
		return nil, false
	}
	return o.ElidingMaxTopicCount, true
}

// HasElidingMaxTopicCount returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasElidingMaxTopicCount() bool {
	if o != nil && o.ElidingMaxTopicCount != nil {
		return true
	}

	return false
}

// SetElidingMaxTopicCount gets a reference to the given int64 and assigns it to the ElidingMaxTopicCount field.
func (o *MsgVpnClientProfile) SetElidingMaxTopicCount(v int64) {
	o.ElidingMaxTopicCount = &v
}

// GetEventClientProvisionedEndpointSpoolUsageThreshold returns the EventClientProvisionedEndpointSpoolUsageThreshold field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetEventClientProvisionedEndpointSpoolUsageThreshold() EventThresholdByPercent {
	if o == nil || o.EventClientProvisionedEndpointSpoolUsageThreshold == nil {
		var ret EventThresholdByPercent
		return ret
	}
	return *o.EventClientProvisionedEndpointSpoolUsageThreshold
}

// GetEventClientProvisionedEndpointSpoolUsageThresholdOk returns a tuple with the EventClientProvisionedEndpointSpoolUsageThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetEventClientProvisionedEndpointSpoolUsageThresholdOk() (*EventThresholdByPercent, bool) {
	if o == nil || o.EventClientProvisionedEndpointSpoolUsageThreshold == nil {
		return nil, false
	}
	return o.EventClientProvisionedEndpointSpoolUsageThreshold, true
}

// HasEventClientProvisionedEndpointSpoolUsageThreshold returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasEventClientProvisionedEndpointSpoolUsageThreshold() bool {
	if o != nil && o.EventClientProvisionedEndpointSpoolUsageThreshold != nil {
		return true
	}

	return false
}

// SetEventClientProvisionedEndpointSpoolUsageThreshold gets a reference to the given EventThresholdByPercent and assigns it to the EventClientProvisionedEndpointSpoolUsageThreshold field.
func (o *MsgVpnClientProfile) SetEventClientProvisionedEndpointSpoolUsageThreshold(v EventThresholdByPercent) {
	o.EventClientProvisionedEndpointSpoolUsageThreshold = &v
}

// GetEventConnectionCountPerClientUsernameThreshold returns the EventConnectionCountPerClientUsernameThreshold field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetEventConnectionCountPerClientUsernameThreshold() EventThreshold {
	if o == nil || o.EventConnectionCountPerClientUsernameThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventConnectionCountPerClientUsernameThreshold
}

// GetEventConnectionCountPerClientUsernameThresholdOk returns a tuple with the EventConnectionCountPerClientUsernameThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetEventConnectionCountPerClientUsernameThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventConnectionCountPerClientUsernameThreshold == nil {
		return nil, false
	}
	return o.EventConnectionCountPerClientUsernameThreshold, true
}

// HasEventConnectionCountPerClientUsernameThreshold returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasEventConnectionCountPerClientUsernameThreshold() bool {
	if o != nil && o.EventConnectionCountPerClientUsernameThreshold != nil {
		return true
	}

	return false
}

// SetEventConnectionCountPerClientUsernameThreshold gets a reference to the given EventThreshold and assigns it to the EventConnectionCountPerClientUsernameThreshold field.
func (o *MsgVpnClientProfile) SetEventConnectionCountPerClientUsernameThreshold(v EventThreshold) {
	o.EventConnectionCountPerClientUsernameThreshold = &v
}

// GetEventEgressFlowCountThreshold returns the EventEgressFlowCountThreshold field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetEventEgressFlowCountThreshold() EventThreshold {
	if o == nil || o.EventEgressFlowCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventEgressFlowCountThreshold
}

// GetEventEgressFlowCountThresholdOk returns a tuple with the EventEgressFlowCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetEventEgressFlowCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventEgressFlowCountThreshold == nil {
		return nil, false
	}
	return o.EventEgressFlowCountThreshold, true
}

// HasEventEgressFlowCountThreshold returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasEventEgressFlowCountThreshold() bool {
	if o != nil && o.EventEgressFlowCountThreshold != nil {
		return true
	}

	return false
}

// SetEventEgressFlowCountThreshold gets a reference to the given EventThreshold and assigns it to the EventEgressFlowCountThreshold field.
func (o *MsgVpnClientProfile) SetEventEgressFlowCountThreshold(v EventThreshold) {
	o.EventEgressFlowCountThreshold = &v
}

// GetEventEndpointCountPerClientUsernameThreshold returns the EventEndpointCountPerClientUsernameThreshold field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetEventEndpointCountPerClientUsernameThreshold() EventThreshold {
	if o == nil || o.EventEndpointCountPerClientUsernameThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventEndpointCountPerClientUsernameThreshold
}

// GetEventEndpointCountPerClientUsernameThresholdOk returns a tuple with the EventEndpointCountPerClientUsernameThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetEventEndpointCountPerClientUsernameThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventEndpointCountPerClientUsernameThreshold == nil {
		return nil, false
	}
	return o.EventEndpointCountPerClientUsernameThreshold, true
}

// HasEventEndpointCountPerClientUsernameThreshold returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasEventEndpointCountPerClientUsernameThreshold() bool {
	if o != nil && o.EventEndpointCountPerClientUsernameThreshold != nil {
		return true
	}

	return false
}

// SetEventEndpointCountPerClientUsernameThreshold gets a reference to the given EventThreshold and assigns it to the EventEndpointCountPerClientUsernameThreshold field.
func (o *MsgVpnClientProfile) SetEventEndpointCountPerClientUsernameThreshold(v EventThreshold) {
	o.EventEndpointCountPerClientUsernameThreshold = &v
}

// GetEventIngressFlowCountThreshold returns the EventIngressFlowCountThreshold field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetEventIngressFlowCountThreshold() EventThreshold {
	if o == nil || o.EventIngressFlowCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventIngressFlowCountThreshold
}

// GetEventIngressFlowCountThresholdOk returns a tuple with the EventIngressFlowCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetEventIngressFlowCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventIngressFlowCountThreshold == nil {
		return nil, false
	}
	return o.EventIngressFlowCountThreshold, true
}

// HasEventIngressFlowCountThreshold returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasEventIngressFlowCountThreshold() bool {
	if o != nil && o.EventIngressFlowCountThreshold != nil {
		return true
	}

	return false
}

// SetEventIngressFlowCountThreshold gets a reference to the given EventThreshold and assigns it to the EventIngressFlowCountThreshold field.
func (o *MsgVpnClientProfile) SetEventIngressFlowCountThreshold(v EventThreshold) {
	o.EventIngressFlowCountThreshold = &v
}

// GetEventServiceSmfConnectionCountPerClientUsernameThreshold returns the EventServiceSmfConnectionCountPerClientUsernameThreshold field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetEventServiceSmfConnectionCountPerClientUsernameThreshold() EventThreshold {
	if o == nil || o.EventServiceSmfConnectionCountPerClientUsernameThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventServiceSmfConnectionCountPerClientUsernameThreshold
}

// GetEventServiceSmfConnectionCountPerClientUsernameThresholdOk returns a tuple with the EventServiceSmfConnectionCountPerClientUsernameThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetEventServiceSmfConnectionCountPerClientUsernameThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventServiceSmfConnectionCountPerClientUsernameThreshold == nil {
		return nil, false
	}
	return o.EventServiceSmfConnectionCountPerClientUsernameThreshold, true
}

// HasEventServiceSmfConnectionCountPerClientUsernameThreshold returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasEventServiceSmfConnectionCountPerClientUsernameThreshold() bool {
	if o != nil && o.EventServiceSmfConnectionCountPerClientUsernameThreshold != nil {
		return true
	}

	return false
}

// SetEventServiceSmfConnectionCountPerClientUsernameThreshold gets a reference to the given EventThreshold and assigns it to the EventServiceSmfConnectionCountPerClientUsernameThreshold field.
func (o *MsgVpnClientProfile) SetEventServiceSmfConnectionCountPerClientUsernameThreshold(v EventThreshold) {
	o.EventServiceSmfConnectionCountPerClientUsernameThreshold = &v
}

// GetEventServiceWebConnectionCountPerClientUsernameThreshold returns the EventServiceWebConnectionCountPerClientUsernameThreshold field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetEventServiceWebConnectionCountPerClientUsernameThreshold() EventThreshold {
	if o == nil || o.EventServiceWebConnectionCountPerClientUsernameThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventServiceWebConnectionCountPerClientUsernameThreshold
}

// GetEventServiceWebConnectionCountPerClientUsernameThresholdOk returns a tuple with the EventServiceWebConnectionCountPerClientUsernameThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetEventServiceWebConnectionCountPerClientUsernameThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventServiceWebConnectionCountPerClientUsernameThreshold == nil {
		return nil, false
	}
	return o.EventServiceWebConnectionCountPerClientUsernameThreshold, true
}

// HasEventServiceWebConnectionCountPerClientUsernameThreshold returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasEventServiceWebConnectionCountPerClientUsernameThreshold() bool {
	if o != nil && o.EventServiceWebConnectionCountPerClientUsernameThreshold != nil {
		return true
	}

	return false
}

// SetEventServiceWebConnectionCountPerClientUsernameThreshold gets a reference to the given EventThreshold and assigns it to the EventServiceWebConnectionCountPerClientUsernameThreshold field.
func (o *MsgVpnClientProfile) SetEventServiceWebConnectionCountPerClientUsernameThreshold(v EventThreshold) {
	o.EventServiceWebConnectionCountPerClientUsernameThreshold = &v
}

// GetEventSubscriptionCountThreshold returns the EventSubscriptionCountThreshold field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetEventSubscriptionCountThreshold() EventThreshold {
	if o == nil || o.EventSubscriptionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventSubscriptionCountThreshold
}

// GetEventSubscriptionCountThresholdOk returns a tuple with the EventSubscriptionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetEventSubscriptionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventSubscriptionCountThreshold == nil {
		return nil, false
	}
	return o.EventSubscriptionCountThreshold, true
}

// HasEventSubscriptionCountThreshold returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasEventSubscriptionCountThreshold() bool {
	if o != nil && o.EventSubscriptionCountThreshold != nil {
		return true
	}

	return false
}

// SetEventSubscriptionCountThreshold gets a reference to the given EventThreshold and assigns it to the EventSubscriptionCountThreshold field.
func (o *MsgVpnClientProfile) SetEventSubscriptionCountThreshold(v EventThreshold) {
	o.EventSubscriptionCountThreshold = &v
}

// GetEventTransactedSessionCountThreshold returns the EventTransactedSessionCountThreshold field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetEventTransactedSessionCountThreshold() EventThreshold {
	if o == nil || o.EventTransactedSessionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventTransactedSessionCountThreshold
}

// GetEventTransactedSessionCountThresholdOk returns a tuple with the EventTransactedSessionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetEventTransactedSessionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventTransactedSessionCountThreshold == nil {
		return nil, false
	}
	return o.EventTransactedSessionCountThreshold, true
}

// HasEventTransactedSessionCountThreshold returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasEventTransactedSessionCountThreshold() bool {
	if o != nil && o.EventTransactedSessionCountThreshold != nil {
		return true
	}

	return false
}

// SetEventTransactedSessionCountThreshold gets a reference to the given EventThreshold and assigns it to the EventTransactedSessionCountThreshold field.
func (o *MsgVpnClientProfile) SetEventTransactedSessionCountThreshold(v EventThreshold) {
	o.EventTransactedSessionCountThreshold = &v
}

// GetEventTransactionCountThreshold returns the EventTransactionCountThreshold field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetEventTransactionCountThreshold() EventThreshold {
	if o == nil || o.EventTransactionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventTransactionCountThreshold
}

// GetEventTransactionCountThresholdOk returns a tuple with the EventTransactionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetEventTransactionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventTransactionCountThreshold == nil {
		return nil, false
	}
	return o.EventTransactionCountThreshold, true
}

// HasEventTransactionCountThreshold returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasEventTransactionCountThreshold() bool {
	if o != nil && o.EventTransactionCountThreshold != nil {
		return true
	}

	return false
}

// SetEventTransactionCountThreshold gets a reference to the given EventThreshold and assigns it to the EventTransactionCountThreshold field.
func (o *MsgVpnClientProfile) SetEventTransactionCountThreshold(v EventThreshold) {
	o.EventTransactionCountThreshold = &v
}

// GetMaxConnectionCountPerClientUsername returns the MaxConnectionCountPerClientUsername field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetMaxConnectionCountPerClientUsername() int64 {
	if o == nil || o.MaxConnectionCountPerClientUsername == nil {
		var ret int64
		return ret
	}
	return *o.MaxConnectionCountPerClientUsername
}

// GetMaxConnectionCountPerClientUsernameOk returns a tuple with the MaxConnectionCountPerClientUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetMaxConnectionCountPerClientUsernameOk() (*int64, bool) {
	if o == nil || o.MaxConnectionCountPerClientUsername == nil {
		return nil, false
	}
	return o.MaxConnectionCountPerClientUsername, true
}

// HasMaxConnectionCountPerClientUsername returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasMaxConnectionCountPerClientUsername() bool {
	if o != nil && o.MaxConnectionCountPerClientUsername != nil {
		return true
	}

	return false
}

// SetMaxConnectionCountPerClientUsername gets a reference to the given int64 and assigns it to the MaxConnectionCountPerClientUsername field.
func (o *MsgVpnClientProfile) SetMaxConnectionCountPerClientUsername(v int64) {
	o.MaxConnectionCountPerClientUsername = &v
}

// GetMaxEgressFlowCount returns the MaxEgressFlowCount field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetMaxEgressFlowCount() int64 {
	if o == nil || o.MaxEgressFlowCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxEgressFlowCount
}

// GetMaxEgressFlowCountOk returns a tuple with the MaxEgressFlowCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetMaxEgressFlowCountOk() (*int64, bool) {
	if o == nil || o.MaxEgressFlowCount == nil {
		return nil, false
	}
	return o.MaxEgressFlowCount, true
}

// HasMaxEgressFlowCount returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasMaxEgressFlowCount() bool {
	if o != nil && o.MaxEgressFlowCount != nil {
		return true
	}

	return false
}

// SetMaxEgressFlowCount gets a reference to the given int64 and assigns it to the MaxEgressFlowCount field.
func (o *MsgVpnClientProfile) SetMaxEgressFlowCount(v int64) {
	o.MaxEgressFlowCount = &v
}

// GetMaxEndpointCountPerClientUsername returns the MaxEndpointCountPerClientUsername field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetMaxEndpointCountPerClientUsername() int64 {
	if o == nil || o.MaxEndpointCountPerClientUsername == nil {
		var ret int64
		return ret
	}
	return *o.MaxEndpointCountPerClientUsername
}

// GetMaxEndpointCountPerClientUsernameOk returns a tuple with the MaxEndpointCountPerClientUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetMaxEndpointCountPerClientUsernameOk() (*int64, bool) {
	if o == nil || o.MaxEndpointCountPerClientUsername == nil {
		return nil, false
	}
	return o.MaxEndpointCountPerClientUsername, true
}

// HasMaxEndpointCountPerClientUsername returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasMaxEndpointCountPerClientUsername() bool {
	if o != nil && o.MaxEndpointCountPerClientUsername != nil {
		return true
	}

	return false
}

// SetMaxEndpointCountPerClientUsername gets a reference to the given int64 and assigns it to the MaxEndpointCountPerClientUsername field.
func (o *MsgVpnClientProfile) SetMaxEndpointCountPerClientUsername(v int64) {
	o.MaxEndpointCountPerClientUsername = &v
}

// GetMaxIngressFlowCount returns the MaxIngressFlowCount field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetMaxIngressFlowCount() int64 {
	if o == nil || o.MaxIngressFlowCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxIngressFlowCount
}

// GetMaxIngressFlowCountOk returns a tuple with the MaxIngressFlowCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetMaxIngressFlowCountOk() (*int64, bool) {
	if o == nil || o.MaxIngressFlowCount == nil {
		return nil, false
	}
	return o.MaxIngressFlowCount, true
}

// HasMaxIngressFlowCount returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasMaxIngressFlowCount() bool {
	if o != nil && o.MaxIngressFlowCount != nil {
		return true
	}

	return false
}

// SetMaxIngressFlowCount gets a reference to the given int64 and assigns it to the MaxIngressFlowCount field.
func (o *MsgVpnClientProfile) SetMaxIngressFlowCount(v int64) {
	o.MaxIngressFlowCount = &v
}

// GetMaxMsgsPerTransaction returns the MaxMsgsPerTransaction field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetMaxMsgsPerTransaction() int32 {
	if o == nil || o.MaxMsgsPerTransaction == nil {
		var ret int32
		return ret
	}
	return *o.MaxMsgsPerTransaction
}

// GetMaxMsgsPerTransactionOk returns a tuple with the MaxMsgsPerTransaction field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetMaxMsgsPerTransactionOk() (*int32, bool) {
	if o == nil || o.MaxMsgsPerTransaction == nil {
		return nil, false
	}
	return o.MaxMsgsPerTransaction, true
}

// HasMaxMsgsPerTransaction returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasMaxMsgsPerTransaction() bool {
	if o != nil && o.MaxMsgsPerTransaction != nil {
		return true
	}

	return false
}

// SetMaxMsgsPerTransaction gets a reference to the given int32 and assigns it to the MaxMsgsPerTransaction field.
func (o *MsgVpnClientProfile) SetMaxMsgsPerTransaction(v int32) {
	o.MaxMsgsPerTransaction = &v
}

// GetMaxSubscriptionCount returns the MaxSubscriptionCount field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetMaxSubscriptionCount() int64 {
	if o == nil || o.MaxSubscriptionCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxSubscriptionCount
}

// GetMaxSubscriptionCountOk returns a tuple with the MaxSubscriptionCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetMaxSubscriptionCountOk() (*int64, bool) {
	if o == nil || o.MaxSubscriptionCount == nil {
		return nil, false
	}
	return o.MaxSubscriptionCount, true
}

// HasMaxSubscriptionCount returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasMaxSubscriptionCount() bool {
	if o != nil && o.MaxSubscriptionCount != nil {
		return true
	}

	return false
}

// SetMaxSubscriptionCount gets a reference to the given int64 and assigns it to the MaxSubscriptionCount field.
func (o *MsgVpnClientProfile) SetMaxSubscriptionCount(v int64) {
	o.MaxSubscriptionCount = &v
}

// GetMaxTransactedSessionCount returns the MaxTransactedSessionCount field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetMaxTransactedSessionCount() int64 {
	if o == nil || o.MaxTransactedSessionCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxTransactedSessionCount
}

// GetMaxTransactedSessionCountOk returns a tuple with the MaxTransactedSessionCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetMaxTransactedSessionCountOk() (*int64, bool) {
	if o == nil || o.MaxTransactedSessionCount == nil {
		return nil, false
	}
	return o.MaxTransactedSessionCount, true
}

// HasMaxTransactedSessionCount returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasMaxTransactedSessionCount() bool {
	if o != nil && o.MaxTransactedSessionCount != nil {
		return true
	}

	return false
}

// SetMaxTransactedSessionCount gets a reference to the given int64 and assigns it to the MaxTransactedSessionCount field.
func (o *MsgVpnClientProfile) SetMaxTransactedSessionCount(v int64) {
	o.MaxTransactedSessionCount = &v
}

// GetMaxTransactionCount returns the MaxTransactionCount field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetMaxTransactionCount() int64 {
	if o == nil || o.MaxTransactionCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxTransactionCount
}

// GetMaxTransactionCountOk returns a tuple with the MaxTransactionCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetMaxTransactionCountOk() (*int64, bool) {
	if o == nil || o.MaxTransactionCount == nil {
		return nil, false
	}
	return o.MaxTransactionCount, true
}

// HasMaxTransactionCount returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasMaxTransactionCount() bool {
	if o != nil && o.MaxTransactionCount != nil {
		return true
	}

	return false
}

// SetMaxTransactionCount gets a reference to the given int64 and assigns it to the MaxTransactionCount field.
func (o *MsgVpnClientProfile) SetMaxTransactionCount(v int64) {
	o.MaxTransactionCount = &v
}

// GetMsgVpnName returns the MsgVpnName field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetMsgVpnName() string {
	if o == nil || o.MsgVpnName == nil {
		var ret string
		return ret
	}
	return *o.MsgVpnName
}

// GetMsgVpnNameOk returns a tuple with the MsgVpnName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetMsgVpnNameOk() (*string, bool) {
	if o == nil || o.MsgVpnName == nil {
		return nil, false
	}
	return o.MsgVpnName, true
}

// HasMsgVpnName returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasMsgVpnName() bool {
	if o != nil && o.MsgVpnName != nil {
		return true
	}

	return false
}

// SetMsgVpnName gets a reference to the given string and assigns it to the MsgVpnName field.
func (o *MsgVpnClientProfile) SetMsgVpnName(v string) {
	o.MsgVpnName = &v
}

// GetQueueControl1MaxDepth returns the QueueControl1MaxDepth field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetQueueControl1MaxDepth() int32 {
	if o == nil || o.QueueControl1MaxDepth == nil {
		var ret int32
		return ret
	}
	return *o.QueueControl1MaxDepth
}

// GetQueueControl1MaxDepthOk returns a tuple with the QueueControl1MaxDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetQueueControl1MaxDepthOk() (*int32, bool) {
	if o == nil || o.QueueControl1MaxDepth == nil {
		return nil, false
	}
	return o.QueueControl1MaxDepth, true
}

// HasQueueControl1MaxDepth returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasQueueControl1MaxDepth() bool {
	if o != nil && o.QueueControl1MaxDepth != nil {
		return true
	}

	return false
}

// SetQueueControl1MaxDepth gets a reference to the given int32 and assigns it to the QueueControl1MaxDepth field.
func (o *MsgVpnClientProfile) SetQueueControl1MaxDepth(v int32) {
	o.QueueControl1MaxDepth = &v
}

// GetQueueControl1MinMsgBurst returns the QueueControl1MinMsgBurst field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetQueueControl1MinMsgBurst() int32 {
	if o == nil || o.QueueControl1MinMsgBurst == nil {
		var ret int32
		return ret
	}
	return *o.QueueControl1MinMsgBurst
}

// GetQueueControl1MinMsgBurstOk returns a tuple with the QueueControl1MinMsgBurst field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetQueueControl1MinMsgBurstOk() (*int32, bool) {
	if o == nil || o.QueueControl1MinMsgBurst == nil {
		return nil, false
	}
	return o.QueueControl1MinMsgBurst, true
}

// HasQueueControl1MinMsgBurst returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasQueueControl1MinMsgBurst() bool {
	if o != nil && o.QueueControl1MinMsgBurst != nil {
		return true
	}

	return false
}

// SetQueueControl1MinMsgBurst gets a reference to the given int32 and assigns it to the QueueControl1MinMsgBurst field.
func (o *MsgVpnClientProfile) SetQueueControl1MinMsgBurst(v int32) {
	o.QueueControl1MinMsgBurst = &v
}

// GetQueueDirect1MaxDepth returns the QueueDirect1MaxDepth field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetQueueDirect1MaxDepth() int32 {
	if o == nil || o.QueueDirect1MaxDepth == nil {
		var ret int32
		return ret
	}
	return *o.QueueDirect1MaxDepth
}

// GetQueueDirect1MaxDepthOk returns a tuple with the QueueDirect1MaxDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetQueueDirect1MaxDepthOk() (*int32, bool) {
	if o == nil || o.QueueDirect1MaxDepth == nil {
		return nil, false
	}
	return o.QueueDirect1MaxDepth, true
}

// HasQueueDirect1MaxDepth returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasQueueDirect1MaxDepth() bool {
	if o != nil && o.QueueDirect1MaxDepth != nil {
		return true
	}

	return false
}

// SetQueueDirect1MaxDepth gets a reference to the given int32 and assigns it to the QueueDirect1MaxDepth field.
func (o *MsgVpnClientProfile) SetQueueDirect1MaxDepth(v int32) {
	o.QueueDirect1MaxDepth = &v
}

// GetQueueDirect1MinMsgBurst returns the QueueDirect1MinMsgBurst field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetQueueDirect1MinMsgBurst() int32 {
	if o == nil || o.QueueDirect1MinMsgBurst == nil {
		var ret int32
		return ret
	}
	return *o.QueueDirect1MinMsgBurst
}

// GetQueueDirect1MinMsgBurstOk returns a tuple with the QueueDirect1MinMsgBurst field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetQueueDirect1MinMsgBurstOk() (*int32, bool) {
	if o == nil || o.QueueDirect1MinMsgBurst == nil {
		return nil, false
	}
	return o.QueueDirect1MinMsgBurst, true
}

// HasQueueDirect1MinMsgBurst returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasQueueDirect1MinMsgBurst() bool {
	if o != nil && o.QueueDirect1MinMsgBurst != nil {
		return true
	}

	return false
}

// SetQueueDirect1MinMsgBurst gets a reference to the given int32 and assigns it to the QueueDirect1MinMsgBurst field.
func (o *MsgVpnClientProfile) SetQueueDirect1MinMsgBurst(v int32) {
	o.QueueDirect1MinMsgBurst = &v
}

// GetQueueDirect2MaxDepth returns the QueueDirect2MaxDepth field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetQueueDirect2MaxDepth() int32 {
	if o == nil || o.QueueDirect2MaxDepth == nil {
		var ret int32
		return ret
	}
	return *o.QueueDirect2MaxDepth
}

// GetQueueDirect2MaxDepthOk returns a tuple with the QueueDirect2MaxDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetQueueDirect2MaxDepthOk() (*int32, bool) {
	if o == nil || o.QueueDirect2MaxDepth == nil {
		return nil, false
	}
	return o.QueueDirect2MaxDepth, true
}

// HasQueueDirect2MaxDepth returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasQueueDirect2MaxDepth() bool {
	if o != nil && o.QueueDirect2MaxDepth != nil {
		return true
	}

	return false
}

// SetQueueDirect2MaxDepth gets a reference to the given int32 and assigns it to the QueueDirect2MaxDepth field.
func (o *MsgVpnClientProfile) SetQueueDirect2MaxDepth(v int32) {
	o.QueueDirect2MaxDepth = &v
}

// GetQueueDirect2MinMsgBurst returns the QueueDirect2MinMsgBurst field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetQueueDirect2MinMsgBurst() int32 {
	if o == nil || o.QueueDirect2MinMsgBurst == nil {
		var ret int32
		return ret
	}
	return *o.QueueDirect2MinMsgBurst
}

// GetQueueDirect2MinMsgBurstOk returns a tuple with the QueueDirect2MinMsgBurst field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetQueueDirect2MinMsgBurstOk() (*int32, bool) {
	if o == nil || o.QueueDirect2MinMsgBurst == nil {
		return nil, false
	}
	return o.QueueDirect2MinMsgBurst, true
}

// HasQueueDirect2MinMsgBurst returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasQueueDirect2MinMsgBurst() bool {
	if o != nil && o.QueueDirect2MinMsgBurst != nil {
		return true
	}

	return false
}

// SetQueueDirect2MinMsgBurst gets a reference to the given int32 and assigns it to the QueueDirect2MinMsgBurst field.
func (o *MsgVpnClientProfile) SetQueueDirect2MinMsgBurst(v int32) {
	o.QueueDirect2MinMsgBurst = &v
}

// GetQueueDirect3MaxDepth returns the QueueDirect3MaxDepth field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetQueueDirect3MaxDepth() int32 {
	if o == nil || o.QueueDirect3MaxDepth == nil {
		var ret int32
		return ret
	}
	return *o.QueueDirect3MaxDepth
}

// GetQueueDirect3MaxDepthOk returns a tuple with the QueueDirect3MaxDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetQueueDirect3MaxDepthOk() (*int32, bool) {
	if o == nil || o.QueueDirect3MaxDepth == nil {
		return nil, false
	}
	return o.QueueDirect3MaxDepth, true
}

// HasQueueDirect3MaxDepth returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasQueueDirect3MaxDepth() bool {
	if o != nil && o.QueueDirect3MaxDepth != nil {
		return true
	}

	return false
}

// SetQueueDirect3MaxDepth gets a reference to the given int32 and assigns it to the QueueDirect3MaxDepth field.
func (o *MsgVpnClientProfile) SetQueueDirect3MaxDepth(v int32) {
	o.QueueDirect3MaxDepth = &v
}

// GetQueueDirect3MinMsgBurst returns the QueueDirect3MinMsgBurst field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetQueueDirect3MinMsgBurst() int32 {
	if o == nil || o.QueueDirect3MinMsgBurst == nil {
		var ret int32
		return ret
	}
	return *o.QueueDirect3MinMsgBurst
}

// GetQueueDirect3MinMsgBurstOk returns a tuple with the QueueDirect3MinMsgBurst field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetQueueDirect3MinMsgBurstOk() (*int32, bool) {
	if o == nil || o.QueueDirect3MinMsgBurst == nil {
		return nil, false
	}
	return o.QueueDirect3MinMsgBurst, true
}

// HasQueueDirect3MinMsgBurst returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasQueueDirect3MinMsgBurst() bool {
	if o != nil && o.QueueDirect3MinMsgBurst != nil {
		return true
	}

	return false
}

// SetQueueDirect3MinMsgBurst gets a reference to the given int32 and assigns it to the QueueDirect3MinMsgBurst field.
func (o *MsgVpnClientProfile) SetQueueDirect3MinMsgBurst(v int32) {
	o.QueueDirect3MinMsgBurst = &v
}

// GetQueueGuaranteed1MaxDepth returns the QueueGuaranteed1MaxDepth field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetQueueGuaranteed1MaxDepth() int32 {
	if o == nil || o.QueueGuaranteed1MaxDepth == nil {
		var ret int32
		return ret
	}
	return *o.QueueGuaranteed1MaxDepth
}

// GetQueueGuaranteed1MaxDepthOk returns a tuple with the QueueGuaranteed1MaxDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetQueueGuaranteed1MaxDepthOk() (*int32, bool) {
	if o == nil || o.QueueGuaranteed1MaxDepth == nil {
		return nil, false
	}
	return o.QueueGuaranteed1MaxDepth, true
}

// HasQueueGuaranteed1MaxDepth returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasQueueGuaranteed1MaxDepth() bool {
	if o != nil && o.QueueGuaranteed1MaxDepth != nil {
		return true
	}

	return false
}

// SetQueueGuaranteed1MaxDepth gets a reference to the given int32 and assigns it to the QueueGuaranteed1MaxDepth field.
func (o *MsgVpnClientProfile) SetQueueGuaranteed1MaxDepth(v int32) {
	o.QueueGuaranteed1MaxDepth = &v
}

// GetQueueGuaranteed1MinMsgBurst returns the QueueGuaranteed1MinMsgBurst field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetQueueGuaranteed1MinMsgBurst() int32 {
	if o == nil || o.QueueGuaranteed1MinMsgBurst == nil {
		var ret int32
		return ret
	}
	return *o.QueueGuaranteed1MinMsgBurst
}

// GetQueueGuaranteed1MinMsgBurstOk returns a tuple with the QueueGuaranteed1MinMsgBurst field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetQueueGuaranteed1MinMsgBurstOk() (*int32, bool) {
	if o == nil || o.QueueGuaranteed1MinMsgBurst == nil {
		return nil, false
	}
	return o.QueueGuaranteed1MinMsgBurst, true
}

// HasQueueGuaranteed1MinMsgBurst returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasQueueGuaranteed1MinMsgBurst() bool {
	if o != nil && o.QueueGuaranteed1MinMsgBurst != nil {
		return true
	}

	return false
}

// SetQueueGuaranteed1MinMsgBurst gets a reference to the given int32 and assigns it to the QueueGuaranteed1MinMsgBurst field.
func (o *MsgVpnClientProfile) SetQueueGuaranteed1MinMsgBurst(v int32) {
	o.QueueGuaranteed1MinMsgBurst = &v
}

// GetRejectMsgToSenderOnNoSubscriptionMatchEnabled returns the RejectMsgToSenderOnNoSubscriptionMatchEnabled field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetRejectMsgToSenderOnNoSubscriptionMatchEnabled() bool {
	if o == nil || o.RejectMsgToSenderOnNoSubscriptionMatchEnabled == nil {
		var ret bool
		return ret
	}
	return *o.RejectMsgToSenderOnNoSubscriptionMatchEnabled
}

// GetRejectMsgToSenderOnNoSubscriptionMatchEnabledOk returns a tuple with the RejectMsgToSenderOnNoSubscriptionMatchEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetRejectMsgToSenderOnNoSubscriptionMatchEnabledOk() (*bool, bool) {
	if o == nil || o.RejectMsgToSenderOnNoSubscriptionMatchEnabled == nil {
		return nil, false
	}
	return o.RejectMsgToSenderOnNoSubscriptionMatchEnabled, true
}

// HasRejectMsgToSenderOnNoSubscriptionMatchEnabled returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasRejectMsgToSenderOnNoSubscriptionMatchEnabled() bool {
	if o != nil && o.RejectMsgToSenderOnNoSubscriptionMatchEnabled != nil {
		return true
	}

	return false
}

// SetRejectMsgToSenderOnNoSubscriptionMatchEnabled gets a reference to the given bool and assigns it to the RejectMsgToSenderOnNoSubscriptionMatchEnabled field.
func (o *MsgVpnClientProfile) SetRejectMsgToSenderOnNoSubscriptionMatchEnabled(v bool) {
	o.RejectMsgToSenderOnNoSubscriptionMatchEnabled = &v
}

// GetReplicationAllowClientConnectWhenStandbyEnabled returns the ReplicationAllowClientConnectWhenStandbyEnabled field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetReplicationAllowClientConnectWhenStandbyEnabled() bool {
	if o == nil || o.ReplicationAllowClientConnectWhenStandbyEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ReplicationAllowClientConnectWhenStandbyEnabled
}

// GetReplicationAllowClientConnectWhenStandbyEnabledOk returns a tuple with the ReplicationAllowClientConnectWhenStandbyEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetReplicationAllowClientConnectWhenStandbyEnabledOk() (*bool, bool) {
	if o == nil || o.ReplicationAllowClientConnectWhenStandbyEnabled == nil {
		return nil, false
	}
	return o.ReplicationAllowClientConnectWhenStandbyEnabled, true
}

// HasReplicationAllowClientConnectWhenStandbyEnabled returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasReplicationAllowClientConnectWhenStandbyEnabled() bool {
	if o != nil && o.ReplicationAllowClientConnectWhenStandbyEnabled != nil {
		return true
	}

	return false
}

// SetReplicationAllowClientConnectWhenStandbyEnabled gets a reference to the given bool and assigns it to the ReplicationAllowClientConnectWhenStandbyEnabled field.
func (o *MsgVpnClientProfile) SetReplicationAllowClientConnectWhenStandbyEnabled(v bool) {
	o.ReplicationAllowClientConnectWhenStandbyEnabled = &v
}

// GetServiceMinKeepaliveTimeout returns the ServiceMinKeepaliveTimeout field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetServiceMinKeepaliveTimeout() int32 {
	if o == nil || o.ServiceMinKeepaliveTimeout == nil {
		var ret int32
		return ret
	}
	return *o.ServiceMinKeepaliveTimeout
}

// GetServiceMinKeepaliveTimeoutOk returns a tuple with the ServiceMinKeepaliveTimeout field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetServiceMinKeepaliveTimeoutOk() (*int32, bool) {
	if o == nil || o.ServiceMinKeepaliveTimeout == nil {
		return nil, false
	}
	return o.ServiceMinKeepaliveTimeout, true
}

// HasServiceMinKeepaliveTimeout returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasServiceMinKeepaliveTimeout() bool {
	if o != nil && o.ServiceMinKeepaliveTimeout != nil {
		return true
	}

	return false
}

// SetServiceMinKeepaliveTimeout gets a reference to the given int32 and assigns it to the ServiceMinKeepaliveTimeout field.
func (o *MsgVpnClientProfile) SetServiceMinKeepaliveTimeout(v int32) {
	o.ServiceMinKeepaliveTimeout = &v
}

// GetServiceSmfMaxConnectionCountPerClientUsername returns the ServiceSmfMaxConnectionCountPerClientUsername field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetServiceSmfMaxConnectionCountPerClientUsername() int64 {
	if o == nil || o.ServiceSmfMaxConnectionCountPerClientUsername == nil {
		var ret int64
		return ret
	}
	return *o.ServiceSmfMaxConnectionCountPerClientUsername
}

// GetServiceSmfMaxConnectionCountPerClientUsernameOk returns a tuple with the ServiceSmfMaxConnectionCountPerClientUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetServiceSmfMaxConnectionCountPerClientUsernameOk() (*int64, bool) {
	if o == nil || o.ServiceSmfMaxConnectionCountPerClientUsername == nil {
		return nil, false
	}
	return o.ServiceSmfMaxConnectionCountPerClientUsername, true
}

// HasServiceSmfMaxConnectionCountPerClientUsername returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasServiceSmfMaxConnectionCountPerClientUsername() bool {
	if o != nil && o.ServiceSmfMaxConnectionCountPerClientUsername != nil {
		return true
	}

	return false
}

// SetServiceSmfMaxConnectionCountPerClientUsername gets a reference to the given int64 and assigns it to the ServiceSmfMaxConnectionCountPerClientUsername field.
func (o *MsgVpnClientProfile) SetServiceSmfMaxConnectionCountPerClientUsername(v int64) {
	o.ServiceSmfMaxConnectionCountPerClientUsername = &v
}

// GetServiceSmfMinKeepaliveEnabled returns the ServiceSmfMinKeepaliveEnabled field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetServiceSmfMinKeepaliveEnabled() bool {
	if o == nil || o.ServiceSmfMinKeepaliveEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceSmfMinKeepaliveEnabled
}

// GetServiceSmfMinKeepaliveEnabledOk returns a tuple with the ServiceSmfMinKeepaliveEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetServiceSmfMinKeepaliveEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceSmfMinKeepaliveEnabled == nil {
		return nil, false
	}
	return o.ServiceSmfMinKeepaliveEnabled, true
}

// HasServiceSmfMinKeepaliveEnabled returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasServiceSmfMinKeepaliveEnabled() bool {
	if o != nil && o.ServiceSmfMinKeepaliveEnabled != nil {
		return true
	}

	return false
}

// SetServiceSmfMinKeepaliveEnabled gets a reference to the given bool and assigns it to the ServiceSmfMinKeepaliveEnabled field.
func (o *MsgVpnClientProfile) SetServiceSmfMinKeepaliveEnabled(v bool) {
	o.ServiceSmfMinKeepaliveEnabled = &v
}

// GetServiceWebInactiveTimeout returns the ServiceWebInactiveTimeout field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetServiceWebInactiveTimeout() int64 {
	if o == nil || o.ServiceWebInactiveTimeout == nil {
		var ret int64
		return ret
	}
	return *o.ServiceWebInactiveTimeout
}

// GetServiceWebInactiveTimeoutOk returns a tuple with the ServiceWebInactiveTimeout field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetServiceWebInactiveTimeoutOk() (*int64, bool) {
	if o == nil || o.ServiceWebInactiveTimeout == nil {
		return nil, false
	}
	return o.ServiceWebInactiveTimeout, true
}

// HasServiceWebInactiveTimeout returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasServiceWebInactiveTimeout() bool {
	if o != nil && o.ServiceWebInactiveTimeout != nil {
		return true
	}

	return false
}

// SetServiceWebInactiveTimeout gets a reference to the given int64 and assigns it to the ServiceWebInactiveTimeout field.
func (o *MsgVpnClientProfile) SetServiceWebInactiveTimeout(v int64) {
	o.ServiceWebInactiveTimeout = &v
}

// GetServiceWebMaxConnectionCountPerClientUsername returns the ServiceWebMaxConnectionCountPerClientUsername field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetServiceWebMaxConnectionCountPerClientUsername() int64 {
	if o == nil || o.ServiceWebMaxConnectionCountPerClientUsername == nil {
		var ret int64
		return ret
	}
	return *o.ServiceWebMaxConnectionCountPerClientUsername
}

// GetServiceWebMaxConnectionCountPerClientUsernameOk returns a tuple with the ServiceWebMaxConnectionCountPerClientUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetServiceWebMaxConnectionCountPerClientUsernameOk() (*int64, bool) {
	if o == nil || o.ServiceWebMaxConnectionCountPerClientUsername == nil {
		return nil, false
	}
	return o.ServiceWebMaxConnectionCountPerClientUsername, true
}

// HasServiceWebMaxConnectionCountPerClientUsername returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasServiceWebMaxConnectionCountPerClientUsername() bool {
	if o != nil && o.ServiceWebMaxConnectionCountPerClientUsername != nil {
		return true
	}

	return false
}

// SetServiceWebMaxConnectionCountPerClientUsername gets a reference to the given int64 and assigns it to the ServiceWebMaxConnectionCountPerClientUsername field.
func (o *MsgVpnClientProfile) SetServiceWebMaxConnectionCountPerClientUsername(v int64) {
	o.ServiceWebMaxConnectionCountPerClientUsername = &v
}

// GetServiceWebMaxPayload returns the ServiceWebMaxPayload field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetServiceWebMaxPayload() int64 {
	if o == nil || o.ServiceWebMaxPayload == nil {
		var ret int64
		return ret
	}
	return *o.ServiceWebMaxPayload
}

// GetServiceWebMaxPayloadOk returns a tuple with the ServiceWebMaxPayload field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetServiceWebMaxPayloadOk() (*int64, bool) {
	if o == nil || o.ServiceWebMaxPayload == nil {
		return nil, false
	}
	return o.ServiceWebMaxPayload, true
}

// HasServiceWebMaxPayload returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasServiceWebMaxPayload() bool {
	if o != nil && o.ServiceWebMaxPayload != nil {
		return true
	}

	return false
}

// SetServiceWebMaxPayload gets a reference to the given int64 and assigns it to the ServiceWebMaxPayload field.
func (o *MsgVpnClientProfile) SetServiceWebMaxPayload(v int64) {
	o.ServiceWebMaxPayload = &v
}

// GetTcpCongestionWindowSize returns the TcpCongestionWindowSize field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetTcpCongestionWindowSize() int64 {
	if o == nil || o.TcpCongestionWindowSize == nil {
		var ret int64
		return ret
	}
	return *o.TcpCongestionWindowSize
}

// GetTcpCongestionWindowSizeOk returns a tuple with the TcpCongestionWindowSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetTcpCongestionWindowSizeOk() (*int64, bool) {
	if o == nil || o.TcpCongestionWindowSize == nil {
		return nil, false
	}
	return o.TcpCongestionWindowSize, true
}

// HasTcpCongestionWindowSize returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasTcpCongestionWindowSize() bool {
	if o != nil && o.TcpCongestionWindowSize != nil {
		return true
	}

	return false
}

// SetTcpCongestionWindowSize gets a reference to the given int64 and assigns it to the TcpCongestionWindowSize field.
func (o *MsgVpnClientProfile) SetTcpCongestionWindowSize(v int64) {
	o.TcpCongestionWindowSize = &v
}

// GetTcpKeepaliveCount returns the TcpKeepaliveCount field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetTcpKeepaliveCount() int64 {
	if o == nil || o.TcpKeepaliveCount == nil {
		var ret int64
		return ret
	}
	return *o.TcpKeepaliveCount
}

// GetTcpKeepaliveCountOk returns a tuple with the TcpKeepaliveCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetTcpKeepaliveCountOk() (*int64, bool) {
	if o == nil || o.TcpKeepaliveCount == nil {
		return nil, false
	}
	return o.TcpKeepaliveCount, true
}

// HasTcpKeepaliveCount returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasTcpKeepaliveCount() bool {
	if o != nil && o.TcpKeepaliveCount != nil {
		return true
	}

	return false
}

// SetTcpKeepaliveCount gets a reference to the given int64 and assigns it to the TcpKeepaliveCount field.
func (o *MsgVpnClientProfile) SetTcpKeepaliveCount(v int64) {
	o.TcpKeepaliveCount = &v
}

// GetTcpKeepaliveIdleTime returns the TcpKeepaliveIdleTime field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetTcpKeepaliveIdleTime() int64 {
	if o == nil || o.TcpKeepaliveIdleTime == nil {
		var ret int64
		return ret
	}
	return *o.TcpKeepaliveIdleTime
}

// GetTcpKeepaliveIdleTimeOk returns a tuple with the TcpKeepaliveIdleTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetTcpKeepaliveIdleTimeOk() (*int64, bool) {
	if o == nil || o.TcpKeepaliveIdleTime == nil {
		return nil, false
	}
	return o.TcpKeepaliveIdleTime, true
}

// HasTcpKeepaliveIdleTime returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasTcpKeepaliveIdleTime() bool {
	if o != nil && o.TcpKeepaliveIdleTime != nil {
		return true
	}

	return false
}

// SetTcpKeepaliveIdleTime gets a reference to the given int64 and assigns it to the TcpKeepaliveIdleTime field.
func (o *MsgVpnClientProfile) SetTcpKeepaliveIdleTime(v int64) {
	o.TcpKeepaliveIdleTime = &v
}

// GetTcpKeepaliveInterval returns the TcpKeepaliveInterval field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetTcpKeepaliveInterval() int64 {
	if o == nil || o.TcpKeepaliveInterval == nil {
		var ret int64
		return ret
	}
	return *o.TcpKeepaliveInterval
}

// GetTcpKeepaliveIntervalOk returns a tuple with the TcpKeepaliveInterval field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetTcpKeepaliveIntervalOk() (*int64, bool) {
	if o == nil || o.TcpKeepaliveInterval == nil {
		return nil, false
	}
	return o.TcpKeepaliveInterval, true
}

// HasTcpKeepaliveInterval returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasTcpKeepaliveInterval() bool {
	if o != nil && o.TcpKeepaliveInterval != nil {
		return true
	}

	return false
}

// SetTcpKeepaliveInterval gets a reference to the given int64 and assigns it to the TcpKeepaliveInterval field.
func (o *MsgVpnClientProfile) SetTcpKeepaliveInterval(v int64) {
	o.TcpKeepaliveInterval = &v
}

// GetTcpMaxSegmentSize returns the TcpMaxSegmentSize field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetTcpMaxSegmentSize() int64 {
	if o == nil || o.TcpMaxSegmentSize == nil {
		var ret int64
		return ret
	}
	return *o.TcpMaxSegmentSize
}

// GetTcpMaxSegmentSizeOk returns a tuple with the TcpMaxSegmentSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetTcpMaxSegmentSizeOk() (*int64, bool) {
	if o == nil || o.TcpMaxSegmentSize == nil {
		return nil, false
	}
	return o.TcpMaxSegmentSize, true
}

// HasTcpMaxSegmentSize returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasTcpMaxSegmentSize() bool {
	if o != nil && o.TcpMaxSegmentSize != nil {
		return true
	}

	return false
}

// SetTcpMaxSegmentSize gets a reference to the given int64 and assigns it to the TcpMaxSegmentSize field.
func (o *MsgVpnClientProfile) SetTcpMaxSegmentSize(v int64) {
	o.TcpMaxSegmentSize = &v
}

// GetTcpMaxWindowSize returns the TcpMaxWindowSize field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetTcpMaxWindowSize() int64 {
	if o == nil || o.TcpMaxWindowSize == nil {
		var ret int64
		return ret
	}
	return *o.TcpMaxWindowSize
}

// GetTcpMaxWindowSizeOk returns a tuple with the TcpMaxWindowSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetTcpMaxWindowSizeOk() (*int64, bool) {
	if o == nil || o.TcpMaxWindowSize == nil {
		return nil, false
	}
	return o.TcpMaxWindowSize, true
}

// HasTcpMaxWindowSize returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasTcpMaxWindowSize() bool {
	if o != nil && o.TcpMaxWindowSize != nil {
		return true
	}

	return false
}

// SetTcpMaxWindowSize gets a reference to the given int64 and assigns it to the TcpMaxWindowSize field.
func (o *MsgVpnClientProfile) SetTcpMaxWindowSize(v int64) {
	o.TcpMaxWindowSize = &v
}

// GetTlsAllowDowngradeToPlainTextEnabled returns the TlsAllowDowngradeToPlainTextEnabled field value if set, zero value otherwise.
func (o *MsgVpnClientProfile) GetTlsAllowDowngradeToPlainTextEnabled() bool {
	if o == nil || o.TlsAllowDowngradeToPlainTextEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TlsAllowDowngradeToPlainTextEnabled
}

// GetTlsAllowDowngradeToPlainTextEnabledOk returns a tuple with the TlsAllowDowngradeToPlainTextEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnClientProfile) GetTlsAllowDowngradeToPlainTextEnabledOk() (*bool, bool) {
	if o == nil || o.TlsAllowDowngradeToPlainTextEnabled == nil {
		return nil, false
	}
	return o.TlsAllowDowngradeToPlainTextEnabled, true
}

// HasTlsAllowDowngradeToPlainTextEnabled returns a boolean if a field has been set.
func (o *MsgVpnClientProfile) HasTlsAllowDowngradeToPlainTextEnabled() bool {
	if o != nil && o.TlsAllowDowngradeToPlainTextEnabled != nil {
		return true
	}

	return false
}

// SetTlsAllowDowngradeToPlainTextEnabled gets a reference to the given bool and assigns it to the TlsAllowDowngradeToPlainTextEnabled field.
func (o *MsgVpnClientProfile) SetTlsAllowDowngradeToPlainTextEnabled(v bool) {
	o.TlsAllowDowngradeToPlainTextEnabled = &v
}

func (o MsgVpnClientProfile) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AllowBridgeConnectionsEnabled != nil {
		toSerialize["allowBridgeConnectionsEnabled"] = o.AllowBridgeConnectionsEnabled
	}
	if o.AllowCutThroughForwardingEnabled != nil {
		toSerialize["allowCutThroughForwardingEnabled"] = o.AllowCutThroughForwardingEnabled
	}
	if o.AllowGuaranteedEndpointCreateDurability != nil {
		toSerialize["allowGuaranteedEndpointCreateDurability"] = o.AllowGuaranteedEndpointCreateDurability
	}
	if o.AllowGuaranteedEndpointCreateEnabled != nil {
		toSerialize["allowGuaranteedEndpointCreateEnabled"] = o.AllowGuaranteedEndpointCreateEnabled
	}
	if o.AllowGuaranteedMsgReceiveEnabled != nil {
		toSerialize["allowGuaranteedMsgReceiveEnabled"] = o.AllowGuaranteedMsgReceiveEnabled
	}
	if o.AllowGuaranteedMsgSendEnabled != nil {
		toSerialize["allowGuaranteedMsgSendEnabled"] = o.AllowGuaranteedMsgSendEnabled
	}
	if o.AllowSharedSubscriptionsEnabled != nil {
		toSerialize["allowSharedSubscriptionsEnabled"] = o.AllowSharedSubscriptionsEnabled
	}
	if o.AllowTransactedSessionsEnabled != nil {
		toSerialize["allowTransactedSessionsEnabled"] = o.AllowTransactedSessionsEnabled
	}
	if o.ApiQueueManagementCopyFromOnCreateName != nil {
		toSerialize["apiQueueManagementCopyFromOnCreateName"] = o.ApiQueueManagementCopyFromOnCreateName
	}
	if o.ApiQueueManagementCopyFromOnCreateTemplateName != nil {
		toSerialize["apiQueueManagementCopyFromOnCreateTemplateName"] = o.ApiQueueManagementCopyFromOnCreateTemplateName
	}
	if o.ApiTopicEndpointManagementCopyFromOnCreateName != nil {
		toSerialize["apiTopicEndpointManagementCopyFromOnCreateName"] = o.ApiTopicEndpointManagementCopyFromOnCreateName
	}
	if o.ApiTopicEndpointManagementCopyFromOnCreateTemplateName != nil {
		toSerialize["apiTopicEndpointManagementCopyFromOnCreateTemplateName"] = o.ApiTopicEndpointManagementCopyFromOnCreateTemplateName
	}
	if o.ClientProfileName != nil {
		toSerialize["clientProfileName"] = o.ClientProfileName
	}
	if o.CompressionEnabled != nil {
		toSerialize["compressionEnabled"] = o.CompressionEnabled
	}
	if o.ElidingDelay != nil {
		toSerialize["elidingDelay"] = o.ElidingDelay
	}
	if o.ElidingEnabled != nil {
		toSerialize["elidingEnabled"] = o.ElidingEnabled
	}
	if o.ElidingMaxTopicCount != nil {
		toSerialize["elidingMaxTopicCount"] = o.ElidingMaxTopicCount
	}
	if o.EventClientProvisionedEndpointSpoolUsageThreshold != nil {
		toSerialize["eventClientProvisionedEndpointSpoolUsageThreshold"] = o.EventClientProvisionedEndpointSpoolUsageThreshold
	}
	if o.EventConnectionCountPerClientUsernameThreshold != nil {
		toSerialize["eventConnectionCountPerClientUsernameThreshold"] = o.EventConnectionCountPerClientUsernameThreshold
	}
	if o.EventEgressFlowCountThreshold != nil {
		toSerialize["eventEgressFlowCountThreshold"] = o.EventEgressFlowCountThreshold
	}
	if o.EventEndpointCountPerClientUsernameThreshold != nil {
		toSerialize["eventEndpointCountPerClientUsernameThreshold"] = o.EventEndpointCountPerClientUsernameThreshold
	}
	if o.EventIngressFlowCountThreshold != nil {
		toSerialize["eventIngressFlowCountThreshold"] = o.EventIngressFlowCountThreshold
	}
	if o.EventServiceSmfConnectionCountPerClientUsernameThreshold != nil {
		toSerialize["eventServiceSmfConnectionCountPerClientUsernameThreshold"] = o.EventServiceSmfConnectionCountPerClientUsernameThreshold
	}
	if o.EventServiceWebConnectionCountPerClientUsernameThreshold != nil {
		toSerialize["eventServiceWebConnectionCountPerClientUsernameThreshold"] = o.EventServiceWebConnectionCountPerClientUsernameThreshold
	}
	if o.EventSubscriptionCountThreshold != nil {
		toSerialize["eventSubscriptionCountThreshold"] = o.EventSubscriptionCountThreshold
	}
	if o.EventTransactedSessionCountThreshold != nil {
		toSerialize["eventTransactedSessionCountThreshold"] = o.EventTransactedSessionCountThreshold
	}
	if o.EventTransactionCountThreshold != nil {
		toSerialize["eventTransactionCountThreshold"] = o.EventTransactionCountThreshold
	}
	if o.MaxConnectionCountPerClientUsername != nil {
		toSerialize["maxConnectionCountPerClientUsername"] = o.MaxConnectionCountPerClientUsername
	}
	if o.MaxEgressFlowCount != nil {
		toSerialize["maxEgressFlowCount"] = o.MaxEgressFlowCount
	}
	if o.MaxEndpointCountPerClientUsername != nil {
		toSerialize["maxEndpointCountPerClientUsername"] = o.MaxEndpointCountPerClientUsername
	}
	if o.MaxIngressFlowCount != nil {
		toSerialize["maxIngressFlowCount"] = o.MaxIngressFlowCount
	}
	if o.MaxMsgsPerTransaction != nil {
		toSerialize["maxMsgsPerTransaction"] = o.MaxMsgsPerTransaction
	}
	if o.MaxSubscriptionCount != nil {
		toSerialize["maxSubscriptionCount"] = o.MaxSubscriptionCount
	}
	if o.MaxTransactedSessionCount != nil {
		toSerialize["maxTransactedSessionCount"] = o.MaxTransactedSessionCount
	}
	if o.MaxTransactionCount != nil {
		toSerialize["maxTransactionCount"] = o.MaxTransactionCount
	}
	if o.MsgVpnName != nil {
		toSerialize["msgVpnName"] = o.MsgVpnName
	}
	if o.QueueControl1MaxDepth != nil {
		toSerialize["queueControl1MaxDepth"] = o.QueueControl1MaxDepth
	}
	if o.QueueControl1MinMsgBurst != nil {
		toSerialize["queueControl1MinMsgBurst"] = o.QueueControl1MinMsgBurst
	}
	if o.QueueDirect1MaxDepth != nil {
		toSerialize["queueDirect1MaxDepth"] = o.QueueDirect1MaxDepth
	}
	if o.QueueDirect1MinMsgBurst != nil {
		toSerialize["queueDirect1MinMsgBurst"] = o.QueueDirect1MinMsgBurst
	}
	if o.QueueDirect2MaxDepth != nil {
		toSerialize["queueDirect2MaxDepth"] = o.QueueDirect2MaxDepth
	}
	if o.QueueDirect2MinMsgBurst != nil {
		toSerialize["queueDirect2MinMsgBurst"] = o.QueueDirect2MinMsgBurst
	}
	if o.QueueDirect3MaxDepth != nil {
		toSerialize["queueDirect3MaxDepth"] = o.QueueDirect3MaxDepth
	}
	if o.QueueDirect3MinMsgBurst != nil {
		toSerialize["queueDirect3MinMsgBurst"] = o.QueueDirect3MinMsgBurst
	}
	if o.QueueGuaranteed1MaxDepth != nil {
		toSerialize["queueGuaranteed1MaxDepth"] = o.QueueGuaranteed1MaxDepth
	}
	if o.QueueGuaranteed1MinMsgBurst != nil {
		toSerialize["queueGuaranteed1MinMsgBurst"] = o.QueueGuaranteed1MinMsgBurst
	}
	if o.RejectMsgToSenderOnNoSubscriptionMatchEnabled != nil {
		toSerialize["rejectMsgToSenderOnNoSubscriptionMatchEnabled"] = o.RejectMsgToSenderOnNoSubscriptionMatchEnabled
	}
	if o.ReplicationAllowClientConnectWhenStandbyEnabled != nil {
		toSerialize["replicationAllowClientConnectWhenStandbyEnabled"] = o.ReplicationAllowClientConnectWhenStandbyEnabled
	}
	if o.ServiceMinKeepaliveTimeout != nil {
		toSerialize["serviceMinKeepaliveTimeout"] = o.ServiceMinKeepaliveTimeout
	}
	if o.ServiceSmfMaxConnectionCountPerClientUsername != nil {
		toSerialize["serviceSmfMaxConnectionCountPerClientUsername"] = o.ServiceSmfMaxConnectionCountPerClientUsername
	}
	if o.ServiceSmfMinKeepaliveEnabled != nil {
		toSerialize["serviceSmfMinKeepaliveEnabled"] = o.ServiceSmfMinKeepaliveEnabled
	}
	if o.ServiceWebInactiveTimeout != nil {
		toSerialize["serviceWebInactiveTimeout"] = o.ServiceWebInactiveTimeout
	}
	if o.ServiceWebMaxConnectionCountPerClientUsername != nil {
		toSerialize["serviceWebMaxConnectionCountPerClientUsername"] = o.ServiceWebMaxConnectionCountPerClientUsername
	}
	if o.ServiceWebMaxPayload != nil {
		toSerialize["serviceWebMaxPayload"] = o.ServiceWebMaxPayload
	}
	if o.TcpCongestionWindowSize != nil {
		toSerialize["tcpCongestionWindowSize"] = o.TcpCongestionWindowSize
	}
	if o.TcpKeepaliveCount != nil {
		toSerialize["tcpKeepaliveCount"] = o.TcpKeepaliveCount
	}
	if o.TcpKeepaliveIdleTime != nil {
		toSerialize["tcpKeepaliveIdleTime"] = o.TcpKeepaliveIdleTime
	}
	if o.TcpKeepaliveInterval != nil {
		toSerialize["tcpKeepaliveInterval"] = o.TcpKeepaliveInterval
	}
	if o.TcpMaxSegmentSize != nil {
		toSerialize["tcpMaxSegmentSize"] = o.TcpMaxSegmentSize
	}
	if o.TcpMaxWindowSize != nil {
		toSerialize["tcpMaxWindowSize"] = o.TcpMaxWindowSize
	}
	if o.TlsAllowDowngradeToPlainTextEnabled != nil {
		toSerialize["tlsAllowDowngradeToPlainTextEnabled"] = o.TlsAllowDowngradeToPlainTextEnabled
	}
	return json.Marshal(toSerialize)
}

type NullableMsgVpnClientProfile struct {
	value *MsgVpnClientProfile
	isSet bool
}

func (v NullableMsgVpnClientProfile) Get() *MsgVpnClientProfile {
	return v.value
}

func (v *NullableMsgVpnClientProfile) Set(val *MsgVpnClientProfile) {
	v.value = val
	v.isSet = true
}

func (v NullableMsgVpnClientProfile) IsSet() bool {
	return v.isSet
}

func (v *NullableMsgVpnClientProfile) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMsgVpnClientProfile(val *MsgVpnClientProfile) *NullableMsgVpnClientProfile {
	return &NullableMsgVpnClientProfile{value: val, isSet: true}
}

func (v NullableMsgVpnClientProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMsgVpnClientProfile) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
