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

// MsgVpnMqttSession struct for MsgVpnMqttSession
type MsgVpnMqttSession struct {
	// Enable or disable the MQTT Session. When disabled, the client is disconnected, new messages matching QoS 0 subscriptions are discarded, and new messages matching QoS 1 subscriptions are stored for future delivery. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	Enabled *bool `json:"enabled,omitempty"`
	// The Client ID of the MQTT Session, which corresponds to the ClientId provided in the MQTT CONNECT packet.
	MqttSessionClientId *string `json:"mqttSessionClientId,omitempty"`
	// The virtual router of the MQTT Session. The allowed values and their meaning are:  <pre> \"primary\" - The MQTT Session belongs to the primary virtual router. \"backup\" - The MQTT Session belongs to the backup virtual router. \"auto\" - The MQTT Session is automatically assigned a virtual router at creation, depending on the broker's active-standby role. </pre>
	MqttSessionVirtualRouter *string `json:"mqttSessionVirtualRouter,omitempty"`
	// The name of the Message VPN.
	MsgVpnName *string `json:"msgVpnName,omitempty"`
	// The owner of the MQTT Session. For externally-created sessions this defaults to the Client Username of the connecting client. For management-created sessions this defaults to empty. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.
	Owner *string `json:"owner,omitempty"`
	// Enable or disable the propagation of consumer acknowledgements (ACKs) received on the active replication Message VPN to the standby replication Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`. Available since 2.14.
	QueueConsumerAckPropagationEnabled *bool `json:"queueConsumerAckPropagationEnabled,omitempty"`
	// The name of the Dead Message Queue (DMQ) used by the MQTT Session Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"#DEAD_MSG_QUEUE\"`. Available since 2.14.
	QueueDeadMsgQueue                            *string         `json:"queueDeadMsgQueue,omitempty"`
	QueueEventBindCountThreshold                 *EventThreshold `json:"queueEventBindCountThreshold,omitempty"`
	QueueEventMsgSpoolUsageThreshold             *EventThreshold `json:"queueEventMsgSpoolUsageThreshold,omitempty"`
	QueueEventRejectLowPriorityMsgLimitThreshold *EventThreshold `json:"queueEventRejectLowPriorityMsgLimitThreshold,omitempty"`
	// The maximum number of consumer flows that can bind to the MQTT Session Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`. Available since 2.14.
	QueueMaxBindCount *int64 `json:"queueMaxBindCount,omitempty"`
	// The maximum number of messages delivered but not acknowledged per flow for the MQTT Session Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `10000`. Available since 2.14.
	QueueMaxDeliveredUnackedMsgsPerFlow *int64 `json:"queueMaxDeliveredUnackedMsgsPerFlow,omitempty"`
	// The maximum message size allowed in the MQTT Session Queue, in bytes (B). Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `10000000`. Available since 2.14.
	QueueMaxMsgSize *int32 `json:"queueMaxMsgSize,omitempty"`
	// The maximum message spool usage allowed by the MQTT Session Queue, in megabytes (MB). A value of 0 only allows spooling of the last message received and disables quota checking. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `5000`. Available since 2.14.
	QueueMaxMsgSpoolUsage *int64 `json:"queueMaxMsgSpoolUsage,omitempty"`
	// The maximum number of times the MQTT Session Queue will attempt redelivery of a message prior to it being discarded or moved to the DMQ. A value of 0 means to retry forever. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.14.
	QueueMaxRedeliveryCount *int64 `json:"queueMaxRedeliveryCount,omitempty"`
	// The maximum time in seconds a message can stay in the MQTT Session Queue when `queueRespectTtlEnabled` is `\"true\"`. A message expires when the lesser of the sender assigned time-to-live (TTL) in the message and the `queueMaxTtl` configured for the MQTT Session Queue, is exceeded. A value of 0 disables expiry. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.14.
	QueueMaxTtl *int64 `json:"queueMaxTtl,omitempty"`
	// Enable or disable the checking of low priority messages against the `queueRejectLowPriorityMsgLimit`. This may only be enabled if `queueRejectMsgToSenderOnDiscardBehavior` does not have a value of `\"never\"`. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.14.
	QueueRejectLowPriorityMsgEnabled *bool `json:"queueRejectLowPriorityMsgEnabled,omitempty"`
	// The number of messages of any priority in the MQTT Session Queue above which low priority messages are not admitted but higher priority messages are allowed. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.14.
	QueueRejectLowPriorityMsgLimit *int64 `json:"queueRejectLowPriorityMsgLimit,omitempty"`
	// Determines when to return negative acknowledgements (NACKs) to sending clients on message discards. Note that NACKs cause the message to not be delivered to any destination and Transacted Session commits to fail. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"when-queue-enabled\"`. The allowed values and their meaning are:  <pre> \"always\" - Always return a negative acknowledgment (NACK) to the sending client on message discard. \"when-queue-enabled\" - Only return a negative acknowledgment (NACK) to the sending client on message discard when the Queue is enabled. \"never\" - Never return a negative acknowledgment (NACK) to the sending client on message discard. </pre>  Available since 2.14.
	QueueRejectMsgToSenderOnDiscardBehavior *string `json:"queueRejectMsgToSenderOnDiscardBehavior,omitempty"`
	// Enable or disable the respecting of the time-to-live (TTL) for messages in the MQTT Session Queue. When enabled, expired messages are discarded or moved to the DMQ. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.14.
	QueueRespectTtlEnabled *bool `json:"queueRespectTtlEnabled,omitempty"`
}

// NewMsgVpnMqttSession instantiates a new MsgVpnMqttSession object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMsgVpnMqttSession() *MsgVpnMqttSession {
	this := MsgVpnMqttSession{}
	return &this
}

// NewMsgVpnMqttSessionWithDefaults instantiates a new MsgVpnMqttSession object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMsgVpnMqttSessionWithDefaults() *MsgVpnMqttSession {
	this := MsgVpnMqttSession{}
	return &this
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetEnabled() bool {
	if o == nil || o.Enabled == nil {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetEnabledOk() (*bool, bool) {
	if o == nil || o.Enabled == nil {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasEnabled() bool {
	if o != nil && o.Enabled != nil {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *MsgVpnMqttSession) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetMqttSessionClientId returns the MqttSessionClientId field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetMqttSessionClientId() string {
	if o == nil || o.MqttSessionClientId == nil {
		var ret string
		return ret
	}
	return *o.MqttSessionClientId
}

// GetMqttSessionClientIdOk returns a tuple with the MqttSessionClientId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetMqttSessionClientIdOk() (*string, bool) {
	if o == nil || o.MqttSessionClientId == nil {
		return nil, false
	}
	return o.MqttSessionClientId, true
}

// HasMqttSessionClientId returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasMqttSessionClientId() bool {
	if o != nil && o.MqttSessionClientId != nil {
		return true
	}

	return false
}

// SetMqttSessionClientId gets a reference to the given string and assigns it to the MqttSessionClientId field.
func (o *MsgVpnMqttSession) SetMqttSessionClientId(v string) {
	o.MqttSessionClientId = &v
}

// GetMqttSessionVirtualRouter returns the MqttSessionVirtualRouter field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetMqttSessionVirtualRouter() string {
	if o == nil || o.MqttSessionVirtualRouter == nil {
		var ret string
		return ret
	}
	return *o.MqttSessionVirtualRouter
}

// GetMqttSessionVirtualRouterOk returns a tuple with the MqttSessionVirtualRouter field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetMqttSessionVirtualRouterOk() (*string, bool) {
	if o == nil || o.MqttSessionVirtualRouter == nil {
		return nil, false
	}
	return o.MqttSessionVirtualRouter, true
}

// HasMqttSessionVirtualRouter returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasMqttSessionVirtualRouter() bool {
	if o != nil && o.MqttSessionVirtualRouter != nil {
		return true
	}

	return false
}

// SetMqttSessionVirtualRouter gets a reference to the given string and assigns it to the MqttSessionVirtualRouter field.
func (o *MsgVpnMqttSession) SetMqttSessionVirtualRouter(v string) {
	o.MqttSessionVirtualRouter = &v
}

// GetMsgVpnName returns the MsgVpnName field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetMsgVpnName() string {
	if o == nil || o.MsgVpnName == nil {
		var ret string
		return ret
	}
	return *o.MsgVpnName
}

// GetMsgVpnNameOk returns a tuple with the MsgVpnName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetMsgVpnNameOk() (*string, bool) {
	if o == nil || o.MsgVpnName == nil {
		return nil, false
	}
	return o.MsgVpnName, true
}

// HasMsgVpnName returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasMsgVpnName() bool {
	if o != nil && o.MsgVpnName != nil {
		return true
	}

	return false
}

// SetMsgVpnName gets a reference to the given string and assigns it to the MsgVpnName field.
func (o *MsgVpnMqttSession) SetMsgVpnName(v string) {
	o.MsgVpnName = &v
}

// GetOwner returns the Owner field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetOwner() string {
	if o == nil || o.Owner == nil {
		var ret string
		return ret
	}
	return *o.Owner
}

// GetOwnerOk returns a tuple with the Owner field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetOwnerOk() (*string, bool) {
	if o == nil || o.Owner == nil {
		return nil, false
	}
	return o.Owner, true
}

// HasOwner returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasOwner() bool {
	if o != nil && o.Owner != nil {
		return true
	}

	return false
}

// SetOwner gets a reference to the given string and assigns it to the Owner field.
func (o *MsgVpnMqttSession) SetOwner(v string) {
	o.Owner = &v
}

// GetQueueConsumerAckPropagationEnabled returns the QueueConsumerAckPropagationEnabled field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueConsumerAckPropagationEnabled() bool {
	if o == nil || o.QueueConsumerAckPropagationEnabled == nil {
		var ret bool
		return ret
	}
	return *o.QueueConsumerAckPropagationEnabled
}

// GetQueueConsumerAckPropagationEnabledOk returns a tuple with the QueueConsumerAckPropagationEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueConsumerAckPropagationEnabledOk() (*bool, bool) {
	if o == nil || o.QueueConsumerAckPropagationEnabled == nil {
		return nil, false
	}
	return o.QueueConsumerAckPropagationEnabled, true
}

// HasQueueConsumerAckPropagationEnabled returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueConsumerAckPropagationEnabled() bool {
	if o != nil && o.QueueConsumerAckPropagationEnabled != nil {
		return true
	}

	return false
}

// SetQueueConsumerAckPropagationEnabled gets a reference to the given bool and assigns it to the QueueConsumerAckPropagationEnabled field.
func (o *MsgVpnMqttSession) SetQueueConsumerAckPropagationEnabled(v bool) {
	o.QueueConsumerAckPropagationEnabled = &v
}

// GetQueueDeadMsgQueue returns the QueueDeadMsgQueue field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueDeadMsgQueue() string {
	if o == nil || o.QueueDeadMsgQueue == nil {
		var ret string
		return ret
	}
	return *o.QueueDeadMsgQueue
}

// GetQueueDeadMsgQueueOk returns a tuple with the QueueDeadMsgQueue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueDeadMsgQueueOk() (*string, bool) {
	if o == nil || o.QueueDeadMsgQueue == nil {
		return nil, false
	}
	return o.QueueDeadMsgQueue, true
}

// HasQueueDeadMsgQueue returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueDeadMsgQueue() bool {
	if o != nil && o.QueueDeadMsgQueue != nil {
		return true
	}

	return false
}

// SetQueueDeadMsgQueue gets a reference to the given string and assigns it to the QueueDeadMsgQueue field.
func (o *MsgVpnMqttSession) SetQueueDeadMsgQueue(v string) {
	o.QueueDeadMsgQueue = &v
}

// GetQueueEventBindCountThreshold returns the QueueEventBindCountThreshold field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueEventBindCountThreshold() EventThreshold {
	if o == nil || o.QueueEventBindCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.QueueEventBindCountThreshold
}

// GetQueueEventBindCountThresholdOk returns a tuple with the QueueEventBindCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueEventBindCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.QueueEventBindCountThreshold == nil {
		return nil, false
	}
	return o.QueueEventBindCountThreshold, true
}

// HasQueueEventBindCountThreshold returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueEventBindCountThreshold() bool {
	if o != nil && o.QueueEventBindCountThreshold != nil {
		return true
	}

	return false
}

// SetQueueEventBindCountThreshold gets a reference to the given EventThreshold and assigns it to the QueueEventBindCountThreshold field.
func (o *MsgVpnMqttSession) SetQueueEventBindCountThreshold(v EventThreshold) {
	o.QueueEventBindCountThreshold = &v
}

// GetQueueEventMsgSpoolUsageThreshold returns the QueueEventMsgSpoolUsageThreshold field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueEventMsgSpoolUsageThreshold() EventThreshold {
	if o == nil || o.QueueEventMsgSpoolUsageThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.QueueEventMsgSpoolUsageThreshold
}

// GetQueueEventMsgSpoolUsageThresholdOk returns a tuple with the QueueEventMsgSpoolUsageThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueEventMsgSpoolUsageThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.QueueEventMsgSpoolUsageThreshold == nil {
		return nil, false
	}
	return o.QueueEventMsgSpoolUsageThreshold, true
}

// HasQueueEventMsgSpoolUsageThreshold returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueEventMsgSpoolUsageThreshold() bool {
	if o != nil && o.QueueEventMsgSpoolUsageThreshold != nil {
		return true
	}

	return false
}

// SetQueueEventMsgSpoolUsageThreshold gets a reference to the given EventThreshold and assigns it to the QueueEventMsgSpoolUsageThreshold field.
func (o *MsgVpnMqttSession) SetQueueEventMsgSpoolUsageThreshold(v EventThreshold) {
	o.QueueEventMsgSpoolUsageThreshold = &v
}

// GetQueueEventRejectLowPriorityMsgLimitThreshold returns the QueueEventRejectLowPriorityMsgLimitThreshold field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueEventRejectLowPriorityMsgLimitThreshold() EventThreshold {
	if o == nil || o.QueueEventRejectLowPriorityMsgLimitThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.QueueEventRejectLowPriorityMsgLimitThreshold
}

// GetQueueEventRejectLowPriorityMsgLimitThresholdOk returns a tuple with the QueueEventRejectLowPriorityMsgLimitThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueEventRejectLowPriorityMsgLimitThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.QueueEventRejectLowPriorityMsgLimitThreshold == nil {
		return nil, false
	}
	return o.QueueEventRejectLowPriorityMsgLimitThreshold, true
}

// HasQueueEventRejectLowPriorityMsgLimitThreshold returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueEventRejectLowPriorityMsgLimitThreshold() bool {
	if o != nil && o.QueueEventRejectLowPriorityMsgLimitThreshold != nil {
		return true
	}

	return false
}

// SetQueueEventRejectLowPriorityMsgLimitThreshold gets a reference to the given EventThreshold and assigns it to the QueueEventRejectLowPriorityMsgLimitThreshold field.
func (o *MsgVpnMqttSession) SetQueueEventRejectLowPriorityMsgLimitThreshold(v EventThreshold) {
	o.QueueEventRejectLowPriorityMsgLimitThreshold = &v
}

// GetQueueMaxBindCount returns the QueueMaxBindCount field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueMaxBindCount() int64 {
	if o == nil || o.QueueMaxBindCount == nil {
		var ret int64
		return ret
	}
	return *o.QueueMaxBindCount
}

// GetQueueMaxBindCountOk returns a tuple with the QueueMaxBindCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueMaxBindCountOk() (*int64, bool) {
	if o == nil || o.QueueMaxBindCount == nil {
		return nil, false
	}
	return o.QueueMaxBindCount, true
}

// HasQueueMaxBindCount returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueMaxBindCount() bool {
	if o != nil && o.QueueMaxBindCount != nil {
		return true
	}

	return false
}

// SetQueueMaxBindCount gets a reference to the given int64 and assigns it to the QueueMaxBindCount field.
func (o *MsgVpnMqttSession) SetQueueMaxBindCount(v int64) {
	o.QueueMaxBindCount = &v
}

// GetQueueMaxDeliveredUnackedMsgsPerFlow returns the QueueMaxDeliveredUnackedMsgsPerFlow field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueMaxDeliveredUnackedMsgsPerFlow() int64 {
	if o == nil || o.QueueMaxDeliveredUnackedMsgsPerFlow == nil {
		var ret int64
		return ret
	}
	return *o.QueueMaxDeliveredUnackedMsgsPerFlow
}

// GetQueueMaxDeliveredUnackedMsgsPerFlowOk returns a tuple with the QueueMaxDeliveredUnackedMsgsPerFlow field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueMaxDeliveredUnackedMsgsPerFlowOk() (*int64, bool) {
	if o == nil || o.QueueMaxDeliveredUnackedMsgsPerFlow == nil {
		return nil, false
	}
	return o.QueueMaxDeliveredUnackedMsgsPerFlow, true
}

// HasQueueMaxDeliveredUnackedMsgsPerFlow returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueMaxDeliveredUnackedMsgsPerFlow() bool {
	if o != nil && o.QueueMaxDeliveredUnackedMsgsPerFlow != nil {
		return true
	}

	return false
}

// SetQueueMaxDeliveredUnackedMsgsPerFlow gets a reference to the given int64 and assigns it to the QueueMaxDeliveredUnackedMsgsPerFlow field.
func (o *MsgVpnMqttSession) SetQueueMaxDeliveredUnackedMsgsPerFlow(v int64) {
	o.QueueMaxDeliveredUnackedMsgsPerFlow = &v
}

// GetQueueMaxMsgSize returns the QueueMaxMsgSize field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueMaxMsgSize() int32 {
	if o == nil || o.QueueMaxMsgSize == nil {
		var ret int32
		return ret
	}
	return *o.QueueMaxMsgSize
}

// GetQueueMaxMsgSizeOk returns a tuple with the QueueMaxMsgSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueMaxMsgSizeOk() (*int32, bool) {
	if o == nil || o.QueueMaxMsgSize == nil {
		return nil, false
	}
	return o.QueueMaxMsgSize, true
}

// HasQueueMaxMsgSize returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueMaxMsgSize() bool {
	if o != nil && o.QueueMaxMsgSize != nil {
		return true
	}

	return false
}

// SetQueueMaxMsgSize gets a reference to the given int32 and assigns it to the QueueMaxMsgSize field.
func (o *MsgVpnMqttSession) SetQueueMaxMsgSize(v int32) {
	o.QueueMaxMsgSize = &v
}

// GetQueueMaxMsgSpoolUsage returns the QueueMaxMsgSpoolUsage field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueMaxMsgSpoolUsage() int64 {
	if o == nil || o.QueueMaxMsgSpoolUsage == nil {
		var ret int64
		return ret
	}
	return *o.QueueMaxMsgSpoolUsage
}

// GetQueueMaxMsgSpoolUsageOk returns a tuple with the QueueMaxMsgSpoolUsage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueMaxMsgSpoolUsageOk() (*int64, bool) {
	if o == nil || o.QueueMaxMsgSpoolUsage == nil {
		return nil, false
	}
	return o.QueueMaxMsgSpoolUsage, true
}

// HasQueueMaxMsgSpoolUsage returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueMaxMsgSpoolUsage() bool {
	if o != nil && o.QueueMaxMsgSpoolUsage != nil {
		return true
	}

	return false
}

// SetQueueMaxMsgSpoolUsage gets a reference to the given int64 and assigns it to the QueueMaxMsgSpoolUsage field.
func (o *MsgVpnMqttSession) SetQueueMaxMsgSpoolUsage(v int64) {
	o.QueueMaxMsgSpoolUsage = &v
}

// GetQueueMaxRedeliveryCount returns the QueueMaxRedeliveryCount field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueMaxRedeliveryCount() int64 {
	if o == nil || o.QueueMaxRedeliveryCount == nil {
		var ret int64
		return ret
	}
	return *o.QueueMaxRedeliveryCount
}

// GetQueueMaxRedeliveryCountOk returns a tuple with the QueueMaxRedeliveryCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueMaxRedeliveryCountOk() (*int64, bool) {
	if o == nil || o.QueueMaxRedeliveryCount == nil {
		return nil, false
	}
	return o.QueueMaxRedeliveryCount, true
}

// HasQueueMaxRedeliveryCount returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueMaxRedeliveryCount() bool {
	if o != nil && o.QueueMaxRedeliveryCount != nil {
		return true
	}

	return false
}

// SetQueueMaxRedeliveryCount gets a reference to the given int64 and assigns it to the QueueMaxRedeliveryCount field.
func (o *MsgVpnMqttSession) SetQueueMaxRedeliveryCount(v int64) {
	o.QueueMaxRedeliveryCount = &v
}

// GetQueueMaxTtl returns the QueueMaxTtl field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueMaxTtl() int64 {
	if o == nil || o.QueueMaxTtl == nil {
		var ret int64
		return ret
	}
	return *o.QueueMaxTtl
}

// GetQueueMaxTtlOk returns a tuple with the QueueMaxTtl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueMaxTtlOk() (*int64, bool) {
	if o == nil || o.QueueMaxTtl == nil {
		return nil, false
	}
	return o.QueueMaxTtl, true
}

// HasQueueMaxTtl returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueMaxTtl() bool {
	if o != nil && o.QueueMaxTtl != nil {
		return true
	}

	return false
}

// SetQueueMaxTtl gets a reference to the given int64 and assigns it to the QueueMaxTtl field.
func (o *MsgVpnMqttSession) SetQueueMaxTtl(v int64) {
	o.QueueMaxTtl = &v
}

// GetQueueRejectLowPriorityMsgEnabled returns the QueueRejectLowPriorityMsgEnabled field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueRejectLowPriorityMsgEnabled() bool {
	if o == nil || o.QueueRejectLowPriorityMsgEnabled == nil {
		var ret bool
		return ret
	}
	return *o.QueueRejectLowPriorityMsgEnabled
}

// GetQueueRejectLowPriorityMsgEnabledOk returns a tuple with the QueueRejectLowPriorityMsgEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueRejectLowPriorityMsgEnabledOk() (*bool, bool) {
	if o == nil || o.QueueRejectLowPriorityMsgEnabled == nil {
		return nil, false
	}
	return o.QueueRejectLowPriorityMsgEnabled, true
}

// HasQueueRejectLowPriorityMsgEnabled returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueRejectLowPriorityMsgEnabled() bool {
	if o != nil && o.QueueRejectLowPriorityMsgEnabled != nil {
		return true
	}

	return false
}

// SetQueueRejectLowPriorityMsgEnabled gets a reference to the given bool and assigns it to the QueueRejectLowPriorityMsgEnabled field.
func (o *MsgVpnMqttSession) SetQueueRejectLowPriorityMsgEnabled(v bool) {
	o.QueueRejectLowPriorityMsgEnabled = &v
}

// GetQueueRejectLowPriorityMsgLimit returns the QueueRejectLowPriorityMsgLimit field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueRejectLowPriorityMsgLimit() int64 {
	if o == nil || o.QueueRejectLowPriorityMsgLimit == nil {
		var ret int64
		return ret
	}
	return *o.QueueRejectLowPriorityMsgLimit
}

// GetQueueRejectLowPriorityMsgLimitOk returns a tuple with the QueueRejectLowPriorityMsgLimit field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueRejectLowPriorityMsgLimitOk() (*int64, bool) {
	if o == nil || o.QueueRejectLowPriorityMsgLimit == nil {
		return nil, false
	}
	return o.QueueRejectLowPriorityMsgLimit, true
}

// HasQueueRejectLowPriorityMsgLimit returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueRejectLowPriorityMsgLimit() bool {
	if o != nil && o.QueueRejectLowPriorityMsgLimit != nil {
		return true
	}

	return false
}

// SetQueueRejectLowPriorityMsgLimit gets a reference to the given int64 and assigns it to the QueueRejectLowPriorityMsgLimit field.
func (o *MsgVpnMqttSession) SetQueueRejectLowPriorityMsgLimit(v int64) {
	o.QueueRejectLowPriorityMsgLimit = &v
}

// GetQueueRejectMsgToSenderOnDiscardBehavior returns the QueueRejectMsgToSenderOnDiscardBehavior field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueRejectMsgToSenderOnDiscardBehavior() string {
	if o == nil || o.QueueRejectMsgToSenderOnDiscardBehavior == nil {
		var ret string
		return ret
	}
	return *o.QueueRejectMsgToSenderOnDiscardBehavior
}

// GetQueueRejectMsgToSenderOnDiscardBehaviorOk returns a tuple with the QueueRejectMsgToSenderOnDiscardBehavior field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueRejectMsgToSenderOnDiscardBehaviorOk() (*string, bool) {
	if o == nil || o.QueueRejectMsgToSenderOnDiscardBehavior == nil {
		return nil, false
	}
	return o.QueueRejectMsgToSenderOnDiscardBehavior, true
}

// HasQueueRejectMsgToSenderOnDiscardBehavior returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueRejectMsgToSenderOnDiscardBehavior() bool {
	if o != nil && o.QueueRejectMsgToSenderOnDiscardBehavior != nil {
		return true
	}

	return false
}

// SetQueueRejectMsgToSenderOnDiscardBehavior gets a reference to the given string and assigns it to the QueueRejectMsgToSenderOnDiscardBehavior field.
func (o *MsgVpnMqttSession) SetQueueRejectMsgToSenderOnDiscardBehavior(v string) {
	o.QueueRejectMsgToSenderOnDiscardBehavior = &v
}

// GetQueueRespectTtlEnabled returns the QueueRespectTtlEnabled field value if set, zero value otherwise.
func (o *MsgVpnMqttSession) GetQueueRespectTtlEnabled() bool {
	if o == nil || o.QueueRespectTtlEnabled == nil {
		var ret bool
		return ret
	}
	return *o.QueueRespectTtlEnabled
}

// GetQueueRespectTtlEnabledOk returns a tuple with the QueueRespectTtlEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnMqttSession) GetQueueRespectTtlEnabledOk() (*bool, bool) {
	if o == nil || o.QueueRespectTtlEnabled == nil {
		return nil, false
	}
	return o.QueueRespectTtlEnabled, true
}

// HasQueueRespectTtlEnabled returns a boolean if a field has been set.
func (o *MsgVpnMqttSession) HasQueueRespectTtlEnabled() bool {
	if o != nil && o.QueueRespectTtlEnabled != nil {
		return true
	}

	return false
}

// SetQueueRespectTtlEnabled gets a reference to the given bool and assigns it to the QueueRespectTtlEnabled field.
func (o *MsgVpnMqttSession) SetQueueRespectTtlEnabled(v bool) {
	o.QueueRespectTtlEnabled = &v
}

func (o MsgVpnMqttSession) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Enabled != nil {
		toSerialize["enabled"] = o.Enabled
	}
	if o.MqttSessionClientId != nil {
		toSerialize["mqttSessionClientId"] = o.MqttSessionClientId
	}
	if o.MqttSessionVirtualRouter != nil {
		toSerialize["mqttSessionVirtualRouter"] = o.MqttSessionVirtualRouter
	}
	if o.MsgVpnName != nil {
		toSerialize["msgVpnName"] = o.MsgVpnName
	}
	if o.Owner != nil {
		toSerialize["owner"] = o.Owner
	}
	if o.QueueConsumerAckPropagationEnabled != nil {
		toSerialize["queueConsumerAckPropagationEnabled"] = o.QueueConsumerAckPropagationEnabled
	}
	if o.QueueDeadMsgQueue != nil {
		toSerialize["queueDeadMsgQueue"] = o.QueueDeadMsgQueue
	}
	if o.QueueEventBindCountThreshold != nil {
		toSerialize["queueEventBindCountThreshold"] = o.QueueEventBindCountThreshold
	}
	if o.QueueEventMsgSpoolUsageThreshold != nil {
		toSerialize["queueEventMsgSpoolUsageThreshold"] = o.QueueEventMsgSpoolUsageThreshold
	}
	if o.QueueEventRejectLowPriorityMsgLimitThreshold != nil {
		toSerialize["queueEventRejectLowPriorityMsgLimitThreshold"] = o.QueueEventRejectLowPriorityMsgLimitThreshold
	}
	if o.QueueMaxBindCount != nil {
		toSerialize["queueMaxBindCount"] = o.QueueMaxBindCount
	}
	if o.QueueMaxDeliveredUnackedMsgsPerFlow != nil {
		toSerialize["queueMaxDeliveredUnackedMsgsPerFlow"] = o.QueueMaxDeliveredUnackedMsgsPerFlow
	}
	if o.QueueMaxMsgSize != nil {
		toSerialize["queueMaxMsgSize"] = o.QueueMaxMsgSize
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
	if o.QueueRejectLowPriorityMsgEnabled != nil {
		toSerialize["queueRejectLowPriorityMsgEnabled"] = o.QueueRejectLowPriorityMsgEnabled
	}
	if o.QueueRejectLowPriorityMsgLimit != nil {
		toSerialize["queueRejectLowPriorityMsgLimit"] = o.QueueRejectLowPriorityMsgLimit
	}
	if o.QueueRejectMsgToSenderOnDiscardBehavior != nil {
		toSerialize["queueRejectMsgToSenderOnDiscardBehavior"] = o.QueueRejectMsgToSenderOnDiscardBehavior
	}
	if o.QueueRespectTtlEnabled != nil {
		toSerialize["queueRespectTtlEnabled"] = o.QueueRespectTtlEnabled
	}
	return json.Marshal(toSerialize)
}

type NullableMsgVpnMqttSession struct {
	value *MsgVpnMqttSession
	isSet bool
}

func (v NullableMsgVpnMqttSession) Get() *MsgVpnMqttSession {
	return v.value
}

func (v *NullableMsgVpnMqttSession) Set(val *MsgVpnMqttSession) {
	v.value = val
	v.isSet = true
}

func (v NullableMsgVpnMqttSession) IsSet() bool {
	return v.isSet
}

func (v *NullableMsgVpnMqttSession) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMsgVpnMqttSession(val *MsgVpnMqttSession) *NullableMsgVpnMqttSession {
	return &NullableMsgVpnMqttSession{value: val, isSet: true}
}

func (v NullableMsgVpnMqttSession) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMsgVpnMqttSession) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
