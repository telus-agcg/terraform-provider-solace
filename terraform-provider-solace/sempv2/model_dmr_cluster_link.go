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

// DmrClusterLink struct for DmrClusterLink
type DmrClusterLink struct {
	// The password used to authenticate with the remote node when using basic internal authentication. If this per-Link password is not configured, the Cluster's password is used instead. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	AuthenticationBasicPassword *string `json:"authenticationBasicPassword,omitempty"`
	// The authentication scheme to be used by the Link which initiates connections to the remote node. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"basic\"`. The allowed values and their meaning are:  <pre> \"basic\" - Basic Authentication Scheme (via username and password). \"client-certificate\" - Client Certificate Authentication Scheme (via certificate file or content). </pre>
	AuthenticationScheme *string `json:"authenticationScheme,omitempty"`
	// The maximum depth of the \"Control 1\" (C-1) priority queue, in work units. Each work unit is 2048 bytes of message data. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `20000`.
	ClientProfileQueueControl1MaxDepth *int32 `json:"clientProfileQueueControl1MaxDepth,omitempty"`
	// The number of messages that are always allowed entry into the \"Control 1\" (C-1) priority queue, regardless of the `clientProfileQueueControl1MaxDepth` value. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `4`.
	ClientProfileQueueControl1MinMsgBurst *int32 `json:"clientProfileQueueControl1MinMsgBurst,omitempty"`
	// The maximum depth of the \"Direct 1\" (D-1) priority queue, in work units. Each work unit is 2048 bytes of message data. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `20000`.
	ClientProfileQueueDirect1MaxDepth *int32 `json:"clientProfileQueueDirect1MaxDepth,omitempty"`
	// The number of messages that are always allowed entry into the \"Direct 1\" (D-1) priority queue, regardless of the `clientProfileQueueDirect1MaxDepth` value. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `4`.
	ClientProfileQueueDirect1MinMsgBurst *int32 `json:"clientProfileQueueDirect1MinMsgBurst,omitempty"`
	// The maximum depth of the \"Direct 2\" (D-2) priority queue, in work units. Each work unit is 2048 bytes of message data. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `20000`.
	ClientProfileQueueDirect2MaxDepth *int32 `json:"clientProfileQueueDirect2MaxDepth,omitempty"`
	// The number of messages that are always allowed entry into the \"Direct 2\" (D-2) priority queue, regardless of the `clientProfileQueueDirect2MaxDepth` value. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `4`.
	ClientProfileQueueDirect2MinMsgBurst *int32 `json:"clientProfileQueueDirect2MinMsgBurst,omitempty"`
	// The maximum depth of the \"Direct 3\" (D-3) priority queue, in work units. Each work unit is 2048 bytes of message data. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `20000`.
	ClientProfileQueueDirect3MaxDepth *int32 `json:"clientProfileQueueDirect3MaxDepth,omitempty"`
	// The number of messages that are always allowed entry into the \"Direct 3\" (D-3) priority queue, regardless of the `clientProfileQueueDirect3MaxDepth` value. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `4`.
	ClientProfileQueueDirect3MinMsgBurst *int32 `json:"clientProfileQueueDirect3MinMsgBurst,omitempty"`
	// The maximum depth of the \"Guaranteed 1\" (G-1) priority queue, in work units. Each work unit is 2048 bytes of message data. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `20000`.
	ClientProfileQueueGuaranteed1MaxDepth *int32 `json:"clientProfileQueueGuaranteed1MaxDepth,omitempty"`
	// The number of messages that are always allowed entry into the \"Guaranteed 1\" (G-3) priority queue, regardless of the `clientProfileQueueGuaranteed1MaxDepth` value. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `255`.
	ClientProfileQueueGuaranteed1MinMsgBurst *int32 `json:"clientProfileQueueGuaranteed1MinMsgBurst,omitempty"`
	// The TCP initial congestion window size, in multiples of the TCP Maximum Segment Size (MSS). Changing the value from its default of 2 results in non-compliance with RFC 2581. Contact support before changing this value. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `2`.
	ClientProfileTcpCongestionWindowSize *int64 `json:"clientProfileTcpCongestionWindowSize,omitempty"`
	// The number of TCP keepalive retransmissions to be carried out before declaring that the remote end is not available. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `5`.
	ClientProfileTcpKeepaliveCount *int64 `json:"clientProfileTcpKeepaliveCount,omitempty"`
	// The amount of time a connection must remain idle before TCP begins sending keepalive probes, in seconds. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `3`.
	ClientProfileTcpKeepaliveIdleTime *int64 `json:"clientProfileTcpKeepaliveIdleTime,omitempty"`
	// The amount of time between TCP keepalive retransmissions when no acknowledgement is received, in seconds. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `1`.
	ClientProfileTcpKeepaliveInterval *int64 `json:"clientProfileTcpKeepaliveInterval,omitempty"`
	// The TCP maximum segment size, in bytes. Changes are applied to all existing connections. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `1460`.
	ClientProfileTcpMaxSegmentSize *int64 `json:"clientProfileTcpMaxSegmentSize,omitempty"`
	// The TCP maximum window size, in kilobytes. Changes are applied to all existing connections. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `256`.
	ClientProfileTcpMaxWindowSize *int64 `json:"clientProfileTcpMaxWindowSize,omitempty"`
	// The name of the Cluster.
	DmrClusterName *string `json:"dmrClusterName,omitempty"`
	// The number of outstanding guaranteed messages that can be sent over the Link before acknowledgement is received by the sender. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `255`.
	EgressFlowWindowSize *int64 `json:"egressFlowWindowSize,omitempty"`
	// Enable or disable the Link. When disabled, subscription sets of this and the remote node are not kept up-to-date, and messages are not exchanged with the remote node. Published guaranteed messages will be queued up for future delivery based on current subscription sets. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.
	Enabled *bool `json:"enabled,omitempty"`
	// The initiator of the Link's TCP connections. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"lexical\"`. The allowed values and their meaning are:  <pre> \"lexical\" - The \"higher\" node-name initiates. \"local\" - The local node initiates. \"remote\" - The remote node initiates. </pre>
	Initiator *string `json:"initiator,omitempty"`
	// The name of the Dead Message Queue (DMQ) used by the Queue for discarded messages. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"#DEAD_MSG_QUEUE\"`.
	QueueDeadMsgQueue             *string         `json:"queueDeadMsgQueue,omitempty"`
	QueueEventSpoolUsageThreshold *EventThreshold `json:"queueEventSpoolUsageThreshold,omitempty"`
	// The maximum number of messages delivered but not acknowledged per flow for the Queue. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `1000000`.
	QueueMaxDeliveredUnackedMsgsPerFlow *int64 `json:"queueMaxDeliveredUnackedMsgsPerFlow,omitempty"`
	// The maximum message spool usage by the Queue (quota), in megabytes (MB). Changes to this attribute are synchronized to HA mates via config-sync. The default value is `800000`.
	QueueMaxMsgSpoolUsage *int64 `json:"queueMaxMsgSpoolUsage,omitempty"`
	// The maximum number of times the Queue will attempt redelivery of a message prior to it being discarded or moved to the DMQ. A value of 0 means to retry forever. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `0`.
	QueueMaxRedeliveryCount *int64 `json:"queueMaxRedeliveryCount,omitempty"`
	// The maximum time in seconds a message can stay in the Queue when `queueRespectTtlEnabled` is `true`. A message expires when the lesser of the sender assigned time-to-live (TTL) in the message and the `queueMaxTtl` configured for the Queue, is exceeded. A value of 0 disables expiry. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `0`.
	QueueMaxTtl *int64 `json:"queueMaxTtl,omitempty"`
	// Determines when to return negative acknowledgements (NACKs) to sending clients on message discards. Note that NACKs cause the message to not be delivered to any destination and Transacted Session commits to fail. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"always\"`. The allowed values and their meaning are:  <pre> \"always\" - Always return a negative acknowledgment (NACK) to the sending client on message discard. \"when-queue-enabled\" - Only return a negative acknowledgment (NACK) to the sending client on message discard when the Queue is enabled. \"never\" - Never return a negative acknowledgment (NACK) to the sending client on message discard. </pre>
	QueueRejectMsgToSenderOnDiscardBehavior *string `json:"queueRejectMsgToSenderOnDiscardBehavior,omitempty"`
	// Enable or disable the respecting of the time-to-live (TTL) for messages in the Queue. When enabled, expired messages are discarded or moved to the DMQ. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.
	QueueRespectTtlEnabled *bool `json:"queueRespectTtlEnabled,omitempty"`
	// The name of the node at the remote end of the Link.
	RemoteNodeName *string `json:"remoteNodeName,omitempty"`
	// The span of the Link, either internal or external. Internal Links connect nodes within the same Cluster. External Links connect nodes within different Clusters. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"external\"`. The allowed values and their meaning are:  <pre> \"internal\" - Link to same cluster. \"external\" - Link to other cluster. </pre>
	Span *string `json:"span,omitempty"`
	// Enable or disable compression on the Link. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.
	TransportCompressedEnabled *bool `json:"transportCompressedEnabled,omitempty"`
	// Enable or disable encryption (TLS) on the Link. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.
	TransportTlsEnabled *bool `json:"transportTlsEnabled,omitempty"`
}

// NewDmrClusterLink instantiates a new DmrClusterLink object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDmrClusterLink() *DmrClusterLink {
	this := DmrClusterLink{}
	return &this
}

// NewDmrClusterLinkWithDefaults instantiates a new DmrClusterLink object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDmrClusterLinkWithDefaults() *DmrClusterLink {
	this := DmrClusterLink{}
	return &this
}

// GetAuthenticationBasicPassword returns the AuthenticationBasicPassword field value if set, zero value otherwise.
func (o *DmrClusterLink) GetAuthenticationBasicPassword() string {
	if o == nil || o.AuthenticationBasicPassword == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationBasicPassword
}

// GetAuthenticationBasicPasswordOk returns a tuple with the AuthenticationBasicPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetAuthenticationBasicPasswordOk() (*string, bool) {
	if o == nil || o.AuthenticationBasicPassword == nil {
		return nil, false
	}
	return o.AuthenticationBasicPassword, true
}

// HasAuthenticationBasicPassword returns a boolean if a field has been set.
func (o *DmrClusterLink) HasAuthenticationBasicPassword() bool {
	if o != nil && o.AuthenticationBasicPassword != nil {
		return true
	}

	return false
}

// SetAuthenticationBasicPassword gets a reference to the given string and assigns it to the AuthenticationBasicPassword field.
func (o *DmrClusterLink) SetAuthenticationBasicPassword(v string) {
	o.AuthenticationBasicPassword = &v
}

// GetAuthenticationScheme returns the AuthenticationScheme field value if set, zero value otherwise.
func (o *DmrClusterLink) GetAuthenticationScheme() string {
	if o == nil || o.AuthenticationScheme == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationScheme
}

// GetAuthenticationSchemeOk returns a tuple with the AuthenticationScheme field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetAuthenticationSchemeOk() (*string, bool) {
	if o == nil || o.AuthenticationScheme == nil {
		return nil, false
	}
	return o.AuthenticationScheme, true
}

// HasAuthenticationScheme returns a boolean if a field has been set.
func (o *DmrClusterLink) HasAuthenticationScheme() bool {
	if o != nil && o.AuthenticationScheme != nil {
		return true
	}

	return false
}

// SetAuthenticationScheme gets a reference to the given string and assigns it to the AuthenticationScheme field.
func (o *DmrClusterLink) SetAuthenticationScheme(v string) {
	o.AuthenticationScheme = &v
}

// GetClientProfileQueueControl1MaxDepth returns the ClientProfileQueueControl1MaxDepth field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileQueueControl1MaxDepth() int32 {
	if o == nil || o.ClientProfileQueueControl1MaxDepth == nil {
		var ret int32
		return ret
	}
	return *o.ClientProfileQueueControl1MaxDepth
}

// GetClientProfileQueueControl1MaxDepthOk returns a tuple with the ClientProfileQueueControl1MaxDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileQueueControl1MaxDepthOk() (*int32, bool) {
	if o == nil || o.ClientProfileQueueControl1MaxDepth == nil {
		return nil, false
	}
	return o.ClientProfileQueueControl1MaxDepth, true
}

// HasClientProfileQueueControl1MaxDepth returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileQueueControl1MaxDepth() bool {
	if o != nil && o.ClientProfileQueueControl1MaxDepth != nil {
		return true
	}

	return false
}

// SetClientProfileQueueControl1MaxDepth gets a reference to the given int32 and assigns it to the ClientProfileQueueControl1MaxDepth field.
func (o *DmrClusterLink) SetClientProfileQueueControl1MaxDepth(v int32) {
	o.ClientProfileQueueControl1MaxDepth = &v
}

// GetClientProfileQueueControl1MinMsgBurst returns the ClientProfileQueueControl1MinMsgBurst field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileQueueControl1MinMsgBurst() int32 {
	if o == nil || o.ClientProfileQueueControl1MinMsgBurst == nil {
		var ret int32
		return ret
	}
	return *o.ClientProfileQueueControl1MinMsgBurst
}

// GetClientProfileQueueControl1MinMsgBurstOk returns a tuple with the ClientProfileQueueControl1MinMsgBurst field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileQueueControl1MinMsgBurstOk() (*int32, bool) {
	if o == nil || o.ClientProfileQueueControl1MinMsgBurst == nil {
		return nil, false
	}
	return o.ClientProfileQueueControl1MinMsgBurst, true
}

// HasClientProfileQueueControl1MinMsgBurst returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileQueueControl1MinMsgBurst() bool {
	if o != nil && o.ClientProfileQueueControl1MinMsgBurst != nil {
		return true
	}

	return false
}

// SetClientProfileQueueControl1MinMsgBurst gets a reference to the given int32 and assigns it to the ClientProfileQueueControl1MinMsgBurst field.
func (o *DmrClusterLink) SetClientProfileQueueControl1MinMsgBurst(v int32) {
	o.ClientProfileQueueControl1MinMsgBurst = &v
}

// GetClientProfileQueueDirect1MaxDepth returns the ClientProfileQueueDirect1MaxDepth field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileQueueDirect1MaxDepth() int32 {
	if o == nil || o.ClientProfileQueueDirect1MaxDepth == nil {
		var ret int32
		return ret
	}
	return *o.ClientProfileQueueDirect1MaxDepth
}

// GetClientProfileQueueDirect1MaxDepthOk returns a tuple with the ClientProfileQueueDirect1MaxDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileQueueDirect1MaxDepthOk() (*int32, bool) {
	if o == nil || o.ClientProfileQueueDirect1MaxDepth == nil {
		return nil, false
	}
	return o.ClientProfileQueueDirect1MaxDepth, true
}

// HasClientProfileQueueDirect1MaxDepth returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileQueueDirect1MaxDepth() bool {
	if o != nil && o.ClientProfileQueueDirect1MaxDepth != nil {
		return true
	}

	return false
}

// SetClientProfileQueueDirect1MaxDepth gets a reference to the given int32 and assigns it to the ClientProfileQueueDirect1MaxDepth field.
func (o *DmrClusterLink) SetClientProfileQueueDirect1MaxDepth(v int32) {
	o.ClientProfileQueueDirect1MaxDepth = &v
}

// GetClientProfileQueueDirect1MinMsgBurst returns the ClientProfileQueueDirect1MinMsgBurst field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileQueueDirect1MinMsgBurst() int32 {
	if o == nil || o.ClientProfileQueueDirect1MinMsgBurst == nil {
		var ret int32
		return ret
	}
	return *o.ClientProfileQueueDirect1MinMsgBurst
}

// GetClientProfileQueueDirect1MinMsgBurstOk returns a tuple with the ClientProfileQueueDirect1MinMsgBurst field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileQueueDirect1MinMsgBurstOk() (*int32, bool) {
	if o == nil || o.ClientProfileQueueDirect1MinMsgBurst == nil {
		return nil, false
	}
	return o.ClientProfileQueueDirect1MinMsgBurst, true
}

// HasClientProfileQueueDirect1MinMsgBurst returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileQueueDirect1MinMsgBurst() bool {
	if o != nil && o.ClientProfileQueueDirect1MinMsgBurst != nil {
		return true
	}

	return false
}

// SetClientProfileQueueDirect1MinMsgBurst gets a reference to the given int32 and assigns it to the ClientProfileQueueDirect1MinMsgBurst field.
func (o *DmrClusterLink) SetClientProfileQueueDirect1MinMsgBurst(v int32) {
	o.ClientProfileQueueDirect1MinMsgBurst = &v
}

// GetClientProfileQueueDirect2MaxDepth returns the ClientProfileQueueDirect2MaxDepth field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileQueueDirect2MaxDepth() int32 {
	if o == nil || o.ClientProfileQueueDirect2MaxDepth == nil {
		var ret int32
		return ret
	}
	return *o.ClientProfileQueueDirect2MaxDepth
}

// GetClientProfileQueueDirect2MaxDepthOk returns a tuple with the ClientProfileQueueDirect2MaxDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileQueueDirect2MaxDepthOk() (*int32, bool) {
	if o == nil || o.ClientProfileQueueDirect2MaxDepth == nil {
		return nil, false
	}
	return o.ClientProfileQueueDirect2MaxDepth, true
}

// HasClientProfileQueueDirect2MaxDepth returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileQueueDirect2MaxDepth() bool {
	if o != nil && o.ClientProfileQueueDirect2MaxDepth != nil {
		return true
	}

	return false
}

// SetClientProfileQueueDirect2MaxDepth gets a reference to the given int32 and assigns it to the ClientProfileQueueDirect2MaxDepth field.
func (o *DmrClusterLink) SetClientProfileQueueDirect2MaxDepth(v int32) {
	o.ClientProfileQueueDirect2MaxDepth = &v
}

// GetClientProfileQueueDirect2MinMsgBurst returns the ClientProfileQueueDirect2MinMsgBurst field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileQueueDirect2MinMsgBurst() int32 {
	if o == nil || o.ClientProfileQueueDirect2MinMsgBurst == nil {
		var ret int32
		return ret
	}
	return *o.ClientProfileQueueDirect2MinMsgBurst
}

// GetClientProfileQueueDirect2MinMsgBurstOk returns a tuple with the ClientProfileQueueDirect2MinMsgBurst field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileQueueDirect2MinMsgBurstOk() (*int32, bool) {
	if o == nil || o.ClientProfileQueueDirect2MinMsgBurst == nil {
		return nil, false
	}
	return o.ClientProfileQueueDirect2MinMsgBurst, true
}

// HasClientProfileQueueDirect2MinMsgBurst returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileQueueDirect2MinMsgBurst() bool {
	if o != nil && o.ClientProfileQueueDirect2MinMsgBurst != nil {
		return true
	}

	return false
}

// SetClientProfileQueueDirect2MinMsgBurst gets a reference to the given int32 and assigns it to the ClientProfileQueueDirect2MinMsgBurst field.
func (o *DmrClusterLink) SetClientProfileQueueDirect2MinMsgBurst(v int32) {
	o.ClientProfileQueueDirect2MinMsgBurst = &v
}

// GetClientProfileQueueDirect3MaxDepth returns the ClientProfileQueueDirect3MaxDepth field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileQueueDirect3MaxDepth() int32 {
	if o == nil || o.ClientProfileQueueDirect3MaxDepth == nil {
		var ret int32
		return ret
	}
	return *o.ClientProfileQueueDirect3MaxDepth
}

// GetClientProfileQueueDirect3MaxDepthOk returns a tuple with the ClientProfileQueueDirect3MaxDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileQueueDirect3MaxDepthOk() (*int32, bool) {
	if o == nil || o.ClientProfileQueueDirect3MaxDepth == nil {
		return nil, false
	}
	return o.ClientProfileQueueDirect3MaxDepth, true
}

// HasClientProfileQueueDirect3MaxDepth returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileQueueDirect3MaxDepth() bool {
	if o != nil && o.ClientProfileQueueDirect3MaxDepth != nil {
		return true
	}

	return false
}

// SetClientProfileQueueDirect3MaxDepth gets a reference to the given int32 and assigns it to the ClientProfileQueueDirect3MaxDepth field.
func (o *DmrClusterLink) SetClientProfileQueueDirect3MaxDepth(v int32) {
	o.ClientProfileQueueDirect3MaxDepth = &v
}

// GetClientProfileQueueDirect3MinMsgBurst returns the ClientProfileQueueDirect3MinMsgBurst field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileQueueDirect3MinMsgBurst() int32 {
	if o == nil || o.ClientProfileQueueDirect3MinMsgBurst == nil {
		var ret int32
		return ret
	}
	return *o.ClientProfileQueueDirect3MinMsgBurst
}

// GetClientProfileQueueDirect3MinMsgBurstOk returns a tuple with the ClientProfileQueueDirect3MinMsgBurst field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileQueueDirect3MinMsgBurstOk() (*int32, bool) {
	if o == nil || o.ClientProfileQueueDirect3MinMsgBurst == nil {
		return nil, false
	}
	return o.ClientProfileQueueDirect3MinMsgBurst, true
}

// HasClientProfileQueueDirect3MinMsgBurst returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileQueueDirect3MinMsgBurst() bool {
	if o != nil && o.ClientProfileQueueDirect3MinMsgBurst != nil {
		return true
	}

	return false
}

// SetClientProfileQueueDirect3MinMsgBurst gets a reference to the given int32 and assigns it to the ClientProfileQueueDirect3MinMsgBurst field.
func (o *DmrClusterLink) SetClientProfileQueueDirect3MinMsgBurst(v int32) {
	o.ClientProfileQueueDirect3MinMsgBurst = &v
}

// GetClientProfileQueueGuaranteed1MaxDepth returns the ClientProfileQueueGuaranteed1MaxDepth field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileQueueGuaranteed1MaxDepth() int32 {
	if o == nil || o.ClientProfileQueueGuaranteed1MaxDepth == nil {
		var ret int32
		return ret
	}
	return *o.ClientProfileQueueGuaranteed1MaxDepth
}

// GetClientProfileQueueGuaranteed1MaxDepthOk returns a tuple with the ClientProfileQueueGuaranteed1MaxDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileQueueGuaranteed1MaxDepthOk() (*int32, bool) {
	if o == nil || o.ClientProfileQueueGuaranteed1MaxDepth == nil {
		return nil, false
	}
	return o.ClientProfileQueueGuaranteed1MaxDepth, true
}

// HasClientProfileQueueGuaranteed1MaxDepth returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileQueueGuaranteed1MaxDepth() bool {
	if o != nil && o.ClientProfileQueueGuaranteed1MaxDepth != nil {
		return true
	}

	return false
}

// SetClientProfileQueueGuaranteed1MaxDepth gets a reference to the given int32 and assigns it to the ClientProfileQueueGuaranteed1MaxDepth field.
func (o *DmrClusterLink) SetClientProfileQueueGuaranteed1MaxDepth(v int32) {
	o.ClientProfileQueueGuaranteed1MaxDepth = &v
}

// GetClientProfileQueueGuaranteed1MinMsgBurst returns the ClientProfileQueueGuaranteed1MinMsgBurst field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileQueueGuaranteed1MinMsgBurst() int32 {
	if o == nil || o.ClientProfileQueueGuaranteed1MinMsgBurst == nil {
		var ret int32
		return ret
	}
	return *o.ClientProfileQueueGuaranteed1MinMsgBurst
}

// GetClientProfileQueueGuaranteed1MinMsgBurstOk returns a tuple with the ClientProfileQueueGuaranteed1MinMsgBurst field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileQueueGuaranteed1MinMsgBurstOk() (*int32, bool) {
	if o == nil || o.ClientProfileQueueGuaranteed1MinMsgBurst == nil {
		return nil, false
	}
	return o.ClientProfileQueueGuaranteed1MinMsgBurst, true
}

// HasClientProfileQueueGuaranteed1MinMsgBurst returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileQueueGuaranteed1MinMsgBurst() bool {
	if o != nil && o.ClientProfileQueueGuaranteed1MinMsgBurst != nil {
		return true
	}

	return false
}

// SetClientProfileQueueGuaranteed1MinMsgBurst gets a reference to the given int32 and assigns it to the ClientProfileQueueGuaranteed1MinMsgBurst field.
func (o *DmrClusterLink) SetClientProfileQueueGuaranteed1MinMsgBurst(v int32) {
	o.ClientProfileQueueGuaranteed1MinMsgBurst = &v
}

// GetClientProfileTcpCongestionWindowSize returns the ClientProfileTcpCongestionWindowSize field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileTcpCongestionWindowSize() int64 {
	if o == nil || o.ClientProfileTcpCongestionWindowSize == nil {
		var ret int64
		return ret
	}
	return *o.ClientProfileTcpCongestionWindowSize
}

// GetClientProfileTcpCongestionWindowSizeOk returns a tuple with the ClientProfileTcpCongestionWindowSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileTcpCongestionWindowSizeOk() (*int64, bool) {
	if o == nil || o.ClientProfileTcpCongestionWindowSize == nil {
		return nil, false
	}
	return o.ClientProfileTcpCongestionWindowSize, true
}

// HasClientProfileTcpCongestionWindowSize returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileTcpCongestionWindowSize() bool {
	if o != nil && o.ClientProfileTcpCongestionWindowSize != nil {
		return true
	}

	return false
}

// SetClientProfileTcpCongestionWindowSize gets a reference to the given int64 and assigns it to the ClientProfileTcpCongestionWindowSize field.
func (o *DmrClusterLink) SetClientProfileTcpCongestionWindowSize(v int64) {
	o.ClientProfileTcpCongestionWindowSize = &v
}

// GetClientProfileTcpKeepaliveCount returns the ClientProfileTcpKeepaliveCount field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileTcpKeepaliveCount() int64 {
	if o == nil || o.ClientProfileTcpKeepaliveCount == nil {
		var ret int64
		return ret
	}
	return *o.ClientProfileTcpKeepaliveCount
}

// GetClientProfileTcpKeepaliveCountOk returns a tuple with the ClientProfileTcpKeepaliveCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileTcpKeepaliveCountOk() (*int64, bool) {
	if o == nil || o.ClientProfileTcpKeepaliveCount == nil {
		return nil, false
	}
	return o.ClientProfileTcpKeepaliveCount, true
}

// HasClientProfileTcpKeepaliveCount returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileTcpKeepaliveCount() bool {
	if o != nil && o.ClientProfileTcpKeepaliveCount != nil {
		return true
	}

	return false
}

// SetClientProfileTcpKeepaliveCount gets a reference to the given int64 and assigns it to the ClientProfileTcpKeepaliveCount field.
func (o *DmrClusterLink) SetClientProfileTcpKeepaliveCount(v int64) {
	o.ClientProfileTcpKeepaliveCount = &v
}

// GetClientProfileTcpKeepaliveIdleTime returns the ClientProfileTcpKeepaliveIdleTime field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileTcpKeepaliveIdleTime() int64 {
	if o == nil || o.ClientProfileTcpKeepaliveIdleTime == nil {
		var ret int64
		return ret
	}
	return *o.ClientProfileTcpKeepaliveIdleTime
}

// GetClientProfileTcpKeepaliveIdleTimeOk returns a tuple with the ClientProfileTcpKeepaliveIdleTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileTcpKeepaliveIdleTimeOk() (*int64, bool) {
	if o == nil || o.ClientProfileTcpKeepaliveIdleTime == nil {
		return nil, false
	}
	return o.ClientProfileTcpKeepaliveIdleTime, true
}

// HasClientProfileTcpKeepaliveIdleTime returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileTcpKeepaliveIdleTime() bool {
	if o != nil && o.ClientProfileTcpKeepaliveIdleTime != nil {
		return true
	}

	return false
}

// SetClientProfileTcpKeepaliveIdleTime gets a reference to the given int64 and assigns it to the ClientProfileTcpKeepaliveIdleTime field.
func (o *DmrClusterLink) SetClientProfileTcpKeepaliveIdleTime(v int64) {
	o.ClientProfileTcpKeepaliveIdleTime = &v
}

// GetClientProfileTcpKeepaliveInterval returns the ClientProfileTcpKeepaliveInterval field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileTcpKeepaliveInterval() int64 {
	if o == nil || o.ClientProfileTcpKeepaliveInterval == nil {
		var ret int64
		return ret
	}
	return *o.ClientProfileTcpKeepaliveInterval
}

// GetClientProfileTcpKeepaliveIntervalOk returns a tuple with the ClientProfileTcpKeepaliveInterval field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileTcpKeepaliveIntervalOk() (*int64, bool) {
	if o == nil || o.ClientProfileTcpKeepaliveInterval == nil {
		return nil, false
	}
	return o.ClientProfileTcpKeepaliveInterval, true
}

// HasClientProfileTcpKeepaliveInterval returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileTcpKeepaliveInterval() bool {
	if o != nil && o.ClientProfileTcpKeepaliveInterval != nil {
		return true
	}

	return false
}

// SetClientProfileTcpKeepaliveInterval gets a reference to the given int64 and assigns it to the ClientProfileTcpKeepaliveInterval field.
func (o *DmrClusterLink) SetClientProfileTcpKeepaliveInterval(v int64) {
	o.ClientProfileTcpKeepaliveInterval = &v
}

// GetClientProfileTcpMaxSegmentSize returns the ClientProfileTcpMaxSegmentSize field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileTcpMaxSegmentSize() int64 {
	if o == nil || o.ClientProfileTcpMaxSegmentSize == nil {
		var ret int64
		return ret
	}
	return *o.ClientProfileTcpMaxSegmentSize
}

// GetClientProfileTcpMaxSegmentSizeOk returns a tuple with the ClientProfileTcpMaxSegmentSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileTcpMaxSegmentSizeOk() (*int64, bool) {
	if o == nil || o.ClientProfileTcpMaxSegmentSize == nil {
		return nil, false
	}
	return o.ClientProfileTcpMaxSegmentSize, true
}

// HasClientProfileTcpMaxSegmentSize returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileTcpMaxSegmentSize() bool {
	if o != nil && o.ClientProfileTcpMaxSegmentSize != nil {
		return true
	}

	return false
}

// SetClientProfileTcpMaxSegmentSize gets a reference to the given int64 and assigns it to the ClientProfileTcpMaxSegmentSize field.
func (o *DmrClusterLink) SetClientProfileTcpMaxSegmentSize(v int64) {
	o.ClientProfileTcpMaxSegmentSize = &v
}

// GetClientProfileTcpMaxWindowSize returns the ClientProfileTcpMaxWindowSize field value if set, zero value otherwise.
func (o *DmrClusterLink) GetClientProfileTcpMaxWindowSize() int64 {
	if o == nil || o.ClientProfileTcpMaxWindowSize == nil {
		var ret int64
		return ret
	}
	return *o.ClientProfileTcpMaxWindowSize
}

// GetClientProfileTcpMaxWindowSizeOk returns a tuple with the ClientProfileTcpMaxWindowSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetClientProfileTcpMaxWindowSizeOk() (*int64, bool) {
	if o == nil || o.ClientProfileTcpMaxWindowSize == nil {
		return nil, false
	}
	return o.ClientProfileTcpMaxWindowSize, true
}

// HasClientProfileTcpMaxWindowSize returns a boolean if a field has been set.
func (o *DmrClusterLink) HasClientProfileTcpMaxWindowSize() bool {
	if o != nil && o.ClientProfileTcpMaxWindowSize != nil {
		return true
	}

	return false
}

// SetClientProfileTcpMaxWindowSize gets a reference to the given int64 and assigns it to the ClientProfileTcpMaxWindowSize field.
func (o *DmrClusterLink) SetClientProfileTcpMaxWindowSize(v int64) {
	o.ClientProfileTcpMaxWindowSize = &v
}

// GetDmrClusterName returns the DmrClusterName field value if set, zero value otherwise.
func (o *DmrClusterLink) GetDmrClusterName() string {
	if o == nil || o.DmrClusterName == nil {
		var ret string
		return ret
	}
	return *o.DmrClusterName
}

// GetDmrClusterNameOk returns a tuple with the DmrClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetDmrClusterNameOk() (*string, bool) {
	if o == nil || o.DmrClusterName == nil {
		return nil, false
	}
	return o.DmrClusterName, true
}

// HasDmrClusterName returns a boolean if a field has been set.
func (o *DmrClusterLink) HasDmrClusterName() bool {
	if o != nil && o.DmrClusterName != nil {
		return true
	}

	return false
}

// SetDmrClusterName gets a reference to the given string and assigns it to the DmrClusterName field.
func (o *DmrClusterLink) SetDmrClusterName(v string) {
	o.DmrClusterName = &v
}

// GetEgressFlowWindowSize returns the EgressFlowWindowSize field value if set, zero value otherwise.
func (o *DmrClusterLink) GetEgressFlowWindowSize() int64 {
	if o == nil || o.EgressFlowWindowSize == nil {
		var ret int64
		return ret
	}
	return *o.EgressFlowWindowSize
}

// GetEgressFlowWindowSizeOk returns a tuple with the EgressFlowWindowSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetEgressFlowWindowSizeOk() (*int64, bool) {
	if o == nil || o.EgressFlowWindowSize == nil {
		return nil, false
	}
	return o.EgressFlowWindowSize, true
}

// HasEgressFlowWindowSize returns a boolean if a field has been set.
func (o *DmrClusterLink) HasEgressFlowWindowSize() bool {
	if o != nil && o.EgressFlowWindowSize != nil {
		return true
	}

	return false
}

// SetEgressFlowWindowSize gets a reference to the given int64 and assigns it to the EgressFlowWindowSize field.
func (o *DmrClusterLink) SetEgressFlowWindowSize(v int64) {
	o.EgressFlowWindowSize = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *DmrClusterLink) GetEnabled() bool {
	if o == nil || o.Enabled == nil {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetEnabledOk() (*bool, bool) {
	if o == nil || o.Enabled == nil {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *DmrClusterLink) HasEnabled() bool {
	if o != nil && o.Enabled != nil {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *DmrClusterLink) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetInitiator returns the Initiator field value if set, zero value otherwise.
func (o *DmrClusterLink) GetInitiator() string {
	if o == nil || o.Initiator == nil {
		var ret string
		return ret
	}
	return *o.Initiator
}

// GetInitiatorOk returns a tuple with the Initiator field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetInitiatorOk() (*string, bool) {
	if o == nil || o.Initiator == nil {
		return nil, false
	}
	return o.Initiator, true
}

// HasInitiator returns a boolean if a field has been set.
func (o *DmrClusterLink) HasInitiator() bool {
	if o != nil && o.Initiator != nil {
		return true
	}

	return false
}

// SetInitiator gets a reference to the given string and assigns it to the Initiator field.
func (o *DmrClusterLink) SetInitiator(v string) {
	o.Initiator = &v
}

// GetQueueDeadMsgQueue returns the QueueDeadMsgQueue field value if set, zero value otherwise.
func (o *DmrClusterLink) GetQueueDeadMsgQueue() string {
	if o == nil || o.QueueDeadMsgQueue == nil {
		var ret string
		return ret
	}
	return *o.QueueDeadMsgQueue
}

// GetQueueDeadMsgQueueOk returns a tuple with the QueueDeadMsgQueue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetQueueDeadMsgQueueOk() (*string, bool) {
	if o == nil || o.QueueDeadMsgQueue == nil {
		return nil, false
	}
	return o.QueueDeadMsgQueue, true
}

// HasQueueDeadMsgQueue returns a boolean if a field has been set.
func (o *DmrClusterLink) HasQueueDeadMsgQueue() bool {
	if o != nil && o.QueueDeadMsgQueue != nil {
		return true
	}

	return false
}

// SetQueueDeadMsgQueue gets a reference to the given string and assigns it to the QueueDeadMsgQueue field.
func (o *DmrClusterLink) SetQueueDeadMsgQueue(v string) {
	o.QueueDeadMsgQueue = &v
}

// GetQueueEventSpoolUsageThreshold returns the QueueEventSpoolUsageThreshold field value if set, zero value otherwise.
func (o *DmrClusterLink) GetQueueEventSpoolUsageThreshold() EventThreshold {
	if o == nil || o.QueueEventSpoolUsageThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.QueueEventSpoolUsageThreshold
}

// GetQueueEventSpoolUsageThresholdOk returns a tuple with the QueueEventSpoolUsageThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetQueueEventSpoolUsageThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.QueueEventSpoolUsageThreshold == nil {
		return nil, false
	}
	return o.QueueEventSpoolUsageThreshold, true
}

// HasQueueEventSpoolUsageThreshold returns a boolean if a field has been set.
func (o *DmrClusterLink) HasQueueEventSpoolUsageThreshold() bool {
	if o != nil && o.QueueEventSpoolUsageThreshold != nil {
		return true
	}

	return false
}

// SetQueueEventSpoolUsageThreshold gets a reference to the given EventThreshold and assigns it to the QueueEventSpoolUsageThreshold field.
func (o *DmrClusterLink) SetQueueEventSpoolUsageThreshold(v EventThreshold) {
	o.QueueEventSpoolUsageThreshold = &v
}

// GetQueueMaxDeliveredUnackedMsgsPerFlow returns the QueueMaxDeliveredUnackedMsgsPerFlow field value if set, zero value otherwise.
func (o *DmrClusterLink) GetQueueMaxDeliveredUnackedMsgsPerFlow() int64 {
	if o == nil || o.QueueMaxDeliveredUnackedMsgsPerFlow == nil {
		var ret int64
		return ret
	}
	return *o.QueueMaxDeliveredUnackedMsgsPerFlow
}

// GetQueueMaxDeliveredUnackedMsgsPerFlowOk returns a tuple with the QueueMaxDeliveredUnackedMsgsPerFlow field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetQueueMaxDeliveredUnackedMsgsPerFlowOk() (*int64, bool) {
	if o == nil || o.QueueMaxDeliveredUnackedMsgsPerFlow == nil {
		return nil, false
	}
	return o.QueueMaxDeliveredUnackedMsgsPerFlow, true
}

// HasQueueMaxDeliveredUnackedMsgsPerFlow returns a boolean if a field has been set.
func (o *DmrClusterLink) HasQueueMaxDeliveredUnackedMsgsPerFlow() bool {
	if o != nil && o.QueueMaxDeliveredUnackedMsgsPerFlow != nil {
		return true
	}

	return false
}

// SetQueueMaxDeliveredUnackedMsgsPerFlow gets a reference to the given int64 and assigns it to the QueueMaxDeliveredUnackedMsgsPerFlow field.
func (o *DmrClusterLink) SetQueueMaxDeliveredUnackedMsgsPerFlow(v int64) {
	o.QueueMaxDeliveredUnackedMsgsPerFlow = &v
}

// GetQueueMaxMsgSpoolUsage returns the QueueMaxMsgSpoolUsage field value if set, zero value otherwise.
func (o *DmrClusterLink) GetQueueMaxMsgSpoolUsage() int64 {
	if o == nil || o.QueueMaxMsgSpoolUsage == nil {
		var ret int64
		return ret
	}
	return *o.QueueMaxMsgSpoolUsage
}

// GetQueueMaxMsgSpoolUsageOk returns a tuple with the QueueMaxMsgSpoolUsage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetQueueMaxMsgSpoolUsageOk() (*int64, bool) {
	if o == nil || o.QueueMaxMsgSpoolUsage == nil {
		return nil, false
	}
	return o.QueueMaxMsgSpoolUsage, true
}

// HasQueueMaxMsgSpoolUsage returns a boolean if a field has been set.
func (o *DmrClusterLink) HasQueueMaxMsgSpoolUsage() bool {
	if o != nil && o.QueueMaxMsgSpoolUsage != nil {
		return true
	}

	return false
}

// SetQueueMaxMsgSpoolUsage gets a reference to the given int64 and assigns it to the QueueMaxMsgSpoolUsage field.
func (o *DmrClusterLink) SetQueueMaxMsgSpoolUsage(v int64) {
	o.QueueMaxMsgSpoolUsage = &v
}

// GetQueueMaxRedeliveryCount returns the QueueMaxRedeliveryCount field value if set, zero value otherwise.
func (o *DmrClusterLink) GetQueueMaxRedeliveryCount() int64 {
	if o == nil || o.QueueMaxRedeliveryCount == nil {
		var ret int64
		return ret
	}
	return *o.QueueMaxRedeliveryCount
}

// GetQueueMaxRedeliveryCountOk returns a tuple with the QueueMaxRedeliveryCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetQueueMaxRedeliveryCountOk() (*int64, bool) {
	if o == nil || o.QueueMaxRedeliveryCount == nil {
		return nil, false
	}
	return o.QueueMaxRedeliveryCount, true
}

// HasQueueMaxRedeliveryCount returns a boolean if a field has been set.
func (o *DmrClusterLink) HasQueueMaxRedeliveryCount() bool {
	if o != nil && o.QueueMaxRedeliveryCount != nil {
		return true
	}

	return false
}

// SetQueueMaxRedeliveryCount gets a reference to the given int64 and assigns it to the QueueMaxRedeliveryCount field.
func (o *DmrClusterLink) SetQueueMaxRedeliveryCount(v int64) {
	o.QueueMaxRedeliveryCount = &v
}

// GetQueueMaxTtl returns the QueueMaxTtl field value if set, zero value otherwise.
func (o *DmrClusterLink) GetQueueMaxTtl() int64 {
	if o == nil || o.QueueMaxTtl == nil {
		var ret int64
		return ret
	}
	return *o.QueueMaxTtl
}

// GetQueueMaxTtlOk returns a tuple with the QueueMaxTtl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetQueueMaxTtlOk() (*int64, bool) {
	if o == nil || o.QueueMaxTtl == nil {
		return nil, false
	}
	return o.QueueMaxTtl, true
}

// HasQueueMaxTtl returns a boolean if a field has been set.
func (o *DmrClusterLink) HasQueueMaxTtl() bool {
	if o != nil && o.QueueMaxTtl != nil {
		return true
	}

	return false
}

// SetQueueMaxTtl gets a reference to the given int64 and assigns it to the QueueMaxTtl field.
func (o *DmrClusterLink) SetQueueMaxTtl(v int64) {
	o.QueueMaxTtl = &v
}

// GetQueueRejectMsgToSenderOnDiscardBehavior returns the QueueRejectMsgToSenderOnDiscardBehavior field value if set, zero value otherwise.
func (o *DmrClusterLink) GetQueueRejectMsgToSenderOnDiscardBehavior() string {
	if o == nil || o.QueueRejectMsgToSenderOnDiscardBehavior == nil {
		var ret string
		return ret
	}
	return *o.QueueRejectMsgToSenderOnDiscardBehavior
}

// GetQueueRejectMsgToSenderOnDiscardBehaviorOk returns a tuple with the QueueRejectMsgToSenderOnDiscardBehavior field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetQueueRejectMsgToSenderOnDiscardBehaviorOk() (*string, bool) {
	if o == nil || o.QueueRejectMsgToSenderOnDiscardBehavior == nil {
		return nil, false
	}
	return o.QueueRejectMsgToSenderOnDiscardBehavior, true
}

// HasQueueRejectMsgToSenderOnDiscardBehavior returns a boolean if a field has been set.
func (o *DmrClusterLink) HasQueueRejectMsgToSenderOnDiscardBehavior() bool {
	if o != nil && o.QueueRejectMsgToSenderOnDiscardBehavior != nil {
		return true
	}

	return false
}

// SetQueueRejectMsgToSenderOnDiscardBehavior gets a reference to the given string and assigns it to the QueueRejectMsgToSenderOnDiscardBehavior field.
func (o *DmrClusterLink) SetQueueRejectMsgToSenderOnDiscardBehavior(v string) {
	o.QueueRejectMsgToSenderOnDiscardBehavior = &v
}

// GetQueueRespectTtlEnabled returns the QueueRespectTtlEnabled field value if set, zero value otherwise.
func (o *DmrClusterLink) GetQueueRespectTtlEnabled() bool {
	if o == nil || o.QueueRespectTtlEnabled == nil {
		var ret bool
		return ret
	}
	return *o.QueueRespectTtlEnabled
}

// GetQueueRespectTtlEnabledOk returns a tuple with the QueueRespectTtlEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetQueueRespectTtlEnabledOk() (*bool, bool) {
	if o == nil || o.QueueRespectTtlEnabled == nil {
		return nil, false
	}
	return o.QueueRespectTtlEnabled, true
}

// HasQueueRespectTtlEnabled returns a boolean if a field has been set.
func (o *DmrClusterLink) HasQueueRespectTtlEnabled() bool {
	if o != nil && o.QueueRespectTtlEnabled != nil {
		return true
	}

	return false
}

// SetQueueRespectTtlEnabled gets a reference to the given bool and assigns it to the QueueRespectTtlEnabled field.
func (o *DmrClusterLink) SetQueueRespectTtlEnabled(v bool) {
	o.QueueRespectTtlEnabled = &v
}

// GetRemoteNodeName returns the RemoteNodeName field value if set, zero value otherwise.
func (o *DmrClusterLink) GetRemoteNodeName() string {
	if o == nil || o.RemoteNodeName == nil {
		var ret string
		return ret
	}
	return *o.RemoteNodeName
}

// GetRemoteNodeNameOk returns a tuple with the RemoteNodeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetRemoteNodeNameOk() (*string, bool) {
	if o == nil || o.RemoteNodeName == nil {
		return nil, false
	}
	return o.RemoteNodeName, true
}

// HasRemoteNodeName returns a boolean if a field has been set.
func (o *DmrClusterLink) HasRemoteNodeName() bool {
	if o != nil && o.RemoteNodeName != nil {
		return true
	}

	return false
}

// SetRemoteNodeName gets a reference to the given string and assigns it to the RemoteNodeName field.
func (o *DmrClusterLink) SetRemoteNodeName(v string) {
	o.RemoteNodeName = &v
}

// GetSpan returns the Span field value if set, zero value otherwise.
func (o *DmrClusterLink) GetSpan() string {
	if o == nil || o.Span == nil {
		var ret string
		return ret
	}
	return *o.Span
}

// GetSpanOk returns a tuple with the Span field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetSpanOk() (*string, bool) {
	if o == nil || o.Span == nil {
		return nil, false
	}
	return o.Span, true
}

// HasSpan returns a boolean if a field has been set.
func (o *DmrClusterLink) HasSpan() bool {
	if o != nil && o.Span != nil {
		return true
	}

	return false
}

// SetSpan gets a reference to the given string and assigns it to the Span field.
func (o *DmrClusterLink) SetSpan(v string) {
	o.Span = &v
}

// GetTransportCompressedEnabled returns the TransportCompressedEnabled field value if set, zero value otherwise.
func (o *DmrClusterLink) GetTransportCompressedEnabled() bool {
	if o == nil || o.TransportCompressedEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TransportCompressedEnabled
}

// GetTransportCompressedEnabledOk returns a tuple with the TransportCompressedEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetTransportCompressedEnabledOk() (*bool, bool) {
	if o == nil || o.TransportCompressedEnabled == nil {
		return nil, false
	}
	return o.TransportCompressedEnabled, true
}

// HasTransportCompressedEnabled returns a boolean if a field has been set.
func (o *DmrClusterLink) HasTransportCompressedEnabled() bool {
	if o != nil && o.TransportCompressedEnabled != nil {
		return true
	}

	return false
}

// SetTransportCompressedEnabled gets a reference to the given bool and assigns it to the TransportCompressedEnabled field.
func (o *DmrClusterLink) SetTransportCompressedEnabled(v bool) {
	o.TransportCompressedEnabled = &v
}

// GetTransportTlsEnabled returns the TransportTlsEnabled field value if set, zero value otherwise.
func (o *DmrClusterLink) GetTransportTlsEnabled() bool {
	if o == nil || o.TransportTlsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TransportTlsEnabled
}

// GetTransportTlsEnabledOk returns a tuple with the TransportTlsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrClusterLink) GetTransportTlsEnabledOk() (*bool, bool) {
	if o == nil || o.TransportTlsEnabled == nil {
		return nil, false
	}
	return o.TransportTlsEnabled, true
}

// HasTransportTlsEnabled returns a boolean if a field has been set.
func (o *DmrClusterLink) HasTransportTlsEnabled() bool {
	if o != nil && o.TransportTlsEnabled != nil {
		return true
	}

	return false
}

// SetTransportTlsEnabled gets a reference to the given bool and assigns it to the TransportTlsEnabled field.
func (o *DmrClusterLink) SetTransportTlsEnabled(v bool) {
	o.TransportTlsEnabled = &v
}

func (o DmrClusterLink) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AuthenticationBasicPassword != nil {
		toSerialize["authenticationBasicPassword"] = o.AuthenticationBasicPassword
	}
	if o.AuthenticationScheme != nil {
		toSerialize["authenticationScheme"] = o.AuthenticationScheme
	}
	if o.ClientProfileQueueControl1MaxDepth != nil {
		toSerialize["clientProfileQueueControl1MaxDepth"] = o.ClientProfileQueueControl1MaxDepth
	}
	if o.ClientProfileQueueControl1MinMsgBurst != nil {
		toSerialize["clientProfileQueueControl1MinMsgBurst"] = o.ClientProfileQueueControl1MinMsgBurst
	}
	if o.ClientProfileQueueDirect1MaxDepth != nil {
		toSerialize["clientProfileQueueDirect1MaxDepth"] = o.ClientProfileQueueDirect1MaxDepth
	}
	if o.ClientProfileQueueDirect1MinMsgBurst != nil {
		toSerialize["clientProfileQueueDirect1MinMsgBurst"] = o.ClientProfileQueueDirect1MinMsgBurst
	}
	if o.ClientProfileQueueDirect2MaxDepth != nil {
		toSerialize["clientProfileQueueDirect2MaxDepth"] = o.ClientProfileQueueDirect2MaxDepth
	}
	if o.ClientProfileQueueDirect2MinMsgBurst != nil {
		toSerialize["clientProfileQueueDirect2MinMsgBurst"] = o.ClientProfileQueueDirect2MinMsgBurst
	}
	if o.ClientProfileQueueDirect3MaxDepth != nil {
		toSerialize["clientProfileQueueDirect3MaxDepth"] = o.ClientProfileQueueDirect3MaxDepth
	}
	if o.ClientProfileQueueDirect3MinMsgBurst != nil {
		toSerialize["clientProfileQueueDirect3MinMsgBurst"] = o.ClientProfileQueueDirect3MinMsgBurst
	}
	if o.ClientProfileQueueGuaranteed1MaxDepth != nil {
		toSerialize["clientProfileQueueGuaranteed1MaxDepth"] = o.ClientProfileQueueGuaranteed1MaxDepth
	}
	if o.ClientProfileQueueGuaranteed1MinMsgBurst != nil {
		toSerialize["clientProfileQueueGuaranteed1MinMsgBurst"] = o.ClientProfileQueueGuaranteed1MinMsgBurst
	}
	if o.ClientProfileTcpCongestionWindowSize != nil {
		toSerialize["clientProfileTcpCongestionWindowSize"] = o.ClientProfileTcpCongestionWindowSize
	}
	if o.ClientProfileTcpKeepaliveCount != nil {
		toSerialize["clientProfileTcpKeepaliveCount"] = o.ClientProfileTcpKeepaliveCount
	}
	if o.ClientProfileTcpKeepaliveIdleTime != nil {
		toSerialize["clientProfileTcpKeepaliveIdleTime"] = o.ClientProfileTcpKeepaliveIdleTime
	}
	if o.ClientProfileTcpKeepaliveInterval != nil {
		toSerialize["clientProfileTcpKeepaliveInterval"] = o.ClientProfileTcpKeepaliveInterval
	}
	if o.ClientProfileTcpMaxSegmentSize != nil {
		toSerialize["clientProfileTcpMaxSegmentSize"] = o.ClientProfileTcpMaxSegmentSize
	}
	if o.ClientProfileTcpMaxWindowSize != nil {
		toSerialize["clientProfileTcpMaxWindowSize"] = o.ClientProfileTcpMaxWindowSize
	}
	if o.DmrClusterName != nil {
		toSerialize["dmrClusterName"] = o.DmrClusterName
	}
	if o.EgressFlowWindowSize != nil {
		toSerialize["egressFlowWindowSize"] = o.EgressFlowWindowSize
	}
	if o.Enabled != nil {
		toSerialize["enabled"] = o.Enabled
	}
	if o.Initiator != nil {
		toSerialize["initiator"] = o.Initiator
	}
	if o.QueueDeadMsgQueue != nil {
		toSerialize["queueDeadMsgQueue"] = o.QueueDeadMsgQueue
	}
	if o.QueueEventSpoolUsageThreshold != nil {
		toSerialize["queueEventSpoolUsageThreshold"] = o.QueueEventSpoolUsageThreshold
	}
	if o.QueueMaxDeliveredUnackedMsgsPerFlow != nil {
		toSerialize["queueMaxDeliveredUnackedMsgsPerFlow"] = o.QueueMaxDeliveredUnackedMsgsPerFlow
	}
	if o.QueueMaxMsgSpoolUsage != nil {
		toSerialize["queueMaxMsgSpoolUsage"] = o.QueueMaxMsgSpoolUsage
	}
	if o.QueueMaxRedeliveryCount != nil {
		toSerialize["queueMaxRedeliveryCount"] = o.QueueMaxRedeliveryCount
	}
	if o.QueueMaxTtl != nil {
		toSerialize["queueMaxTtl"] = o.QueueMaxTtl
	}
	if o.QueueRejectMsgToSenderOnDiscardBehavior != nil {
		toSerialize["queueRejectMsgToSenderOnDiscardBehavior"] = o.QueueRejectMsgToSenderOnDiscardBehavior
	}
	if o.QueueRespectTtlEnabled != nil {
		toSerialize["queueRespectTtlEnabled"] = o.QueueRespectTtlEnabled
	}
	if o.RemoteNodeName != nil {
		toSerialize["remoteNodeName"] = o.RemoteNodeName
	}
	if o.Span != nil {
		toSerialize["span"] = o.Span
	}
	if o.TransportCompressedEnabled != nil {
		toSerialize["transportCompressedEnabled"] = o.TransportCompressedEnabled
	}
	if o.TransportTlsEnabled != nil {
		toSerialize["transportTlsEnabled"] = o.TransportTlsEnabled
	}
	return json.Marshal(toSerialize)
}

type NullableDmrClusterLink struct {
	value *DmrClusterLink
	isSet bool
}

func (v NullableDmrClusterLink) Get() *DmrClusterLink {
	return v.value
}

func (v *NullableDmrClusterLink) Set(val *DmrClusterLink) {
	v.value = val
	v.isSet = true
}

func (v NullableDmrClusterLink) IsSet() bool {
	return v.isSet
}

func (v *NullableDmrClusterLink) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDmrClusterLink(val *DmrClusterLink) *NullableDmrClusterLink {
	return &NullableDmrClusterLink{value: val, isSet: true}
}

func (v NullableDmrClusterLink) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDmrClusterLink) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
