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

// MsgVpnQueueTemplate struct for MsgVpnQueueTemplate
type MsgVpnQueueTemplate struct {
	// The access type for delivering messages to consumer flows. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"exclusive\"`. The allowed values and their meaning are:  <pre> \"exclusive\" - Exclusive delivery of messages to the first bound consumer flow. \"non-exclusive\" - Non-exclusive delivery of messages to all bound consumer flows in a round-robin fashion. </pre>
	AccessType *string `json:"accessType,omitempty"`
	// Enable or disable the propagation of consumer acknowledgements (ACKs) received on the active replication Message VPN to the standby replication Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	ConsumerAckPropagationEnabled *bool `json:"consumerAckPropagationEnabled,omitempty"`
	// The name of the Dead Message Queue (DMQ). Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"#DEAD_MSG_QUEUE\"`.
	DeadMsgQueue *string `json:"deadMsgQueue,omitempty"`
	// The delay, in seconds, to apply to messages arriving on the Queue before the messages are eligible for delivery. This attribute does not apply to MQTT queues created from this template, but it may apply in future releases. Therefore, to maintain forward compatibility, do not set this value on templates that might be used for MQTT queues. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.22.
	DeliveryDelay *int64 `json:"deliveryDelay,omitempty"`
	// Controls the durability of queues created from this template. If non-durable, the created queue will be non-durable, regardless of the specified durability. If none, the created queue will have the requested durability. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"none\"`. The allowed values and their meaning are:  <pre> \"none\" - The durability of the endpoint will be as requested on create. \"non-durable\" - The durability of the created queue will be non-durable, regardless of what was requested. </pre>
	DurabilityOverride                      *string         `json:"durabilityOverride,omitempty"`
	EventBindCountThreshold                 *EventThreshold `json:"eventBindCountThreshold,omitempty"`
	EventMsgSpoolUsageThreshold             *EventThreshold `json:"eventMsgSpoolUsageThreshold,omitempty"`
	EventRejectLowPriorityMsgLimitThreshold *EventThreshold `json:"eventRejectLowPriorityMsgLimitThreshold,omitempty"`
	// The maximum number of consumer flows that can bind. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.
	MaxBindCount *int64 `json:"maxBindCount,omitempty"`
	// The maximum number of messages delivered but not acknowledged per flow. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `10000`.
	MaxDeliveredUnackedMsgsPerFlow *int64 `json:"maxDeliveredUnackedMsgsPerFlow,omitempty"`
	// The maximum message size allowed, in bytes (B). Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `10000000`.
	MaxMsgSize *int32 `json:"maxMsgSize,omitempty"`
	// The maximum message spool usage allowed, in megabytes (MB). A value of 0 only allows spooling of the last message received and disables quota checking. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `5000`.
	MaxMsgSpoolUsage *int64 `json:"maxMsgSpoolUsage,omitempty"`
	// The maximum number of message redelivery attempts that will occur prior to the message being discarded or moved to the DMQ. A value of 0 means to retry forever. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.
	MaxRedeliveryCount *int64 `json:"maxRedeliveryCount,omitempty"`
	// The maximum time in seconds a message can stay in a Queue when `respectTtlEnabled` is `\"true\"`. A message expires when the lesser of the sender assigned time-to-live (TTL) in the message and the `maxTtl` configured for the Queue, is exceeded. A value of 0 disables expiry. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.
	MaxTtl *int64 `json:"maxTtl,omitempty"`
	// The name of the Message VPN.
	MsgVpnName *string `json:"msgVpnName,omitempty"`
	// The permission level for all consumers, excluding the owner. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"no-access\"`. The allowed values and their meaning are:  <pre> \"no-access\" - Disallows all access. \"read-only\" - Read-only access to the messages. \"consume\" - Consume (read and remove) messages. \"modify-topic\" - Consume messages or modify the topic/selector. \"delete\" - Consume messages, modify the topic/selector or delete the Client created endpoint altogether. </pre>
	Permission *string `json:"permission,omitempty"`
	// A wildcardable pattern used to determine which Queues use settings from this Template. Two different wildcards are supported: * and >. Similar to topic filters or subscription patterns, a > matches anything (but only when used at the end), and a * matches zero or more characters but never a slash (/). A > is only a wildcard when used at the end, after a /. A * is only allowed at the end, after a slash (/). Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.
	QueueNameFilter *string `json:"queueNameFilter,omitempty"`
	// The name of the Queue Template.
	QueueTemplateName *string `json:"queueTemplateName,omitempty"`
	// Enable or disable message redelivery. When enabled, the number of redelivery attempts is controlled by maxRedeliveryCount. When disabled, the message will never be delivered from the queue more than once. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`. Available since 2.18.
	RedeliveryEnabled *bool `json:"redeliveryEnabled,omitempty"`
	// Enable or disable the checking of low priority messages against the `rejectLowPriorityMsgLimit`. This may only be enabled if `rejectMsgToSenderOnDiscardBehavior` does not have a value of `\"never\"`. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	RejectLowPriorityMsgEnabled *bool `json:"rejectLowPriorityMsgEnabled,omitempty"`
	// The number of messages of any priority above which low priority messages are not admitted but higher priority messages are allowed. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.
	RejectLowPriorityMsgLimit *int64 `json:"rejectLowPriorityMsgLimit,omitempty"`
	// Determines when to return negative acknowledgements (NACKs) to sending clients on message discards. Note that NACKs prevent the message from being delivered to any destination and Transacted Session commits to fail. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"when-queue-enabled\"`. The allowed values and their meaning are:  <pre> \"always\" - Always return a negative acknowledgment (NACK) to the sending client on message discard. \"when-queue-enabled\" - Only return a negative acknowledgment (NACK) to the sending client on message discard when the Queue is enabled. \"never\" - Never return a negative acknowledgment (NACK) to the sending client on message discard. </pre>
	RejectMsgToSenderOnDiscardBehavior *string `json:"rejectMsgToSenderOnDiscardBehavior,omitempty"`
	// Enable or disable the respecting of message priority. When enabled, messages are delivered in priority order, from 9 (highest) to 0 (lowest). Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	RespectMsgPriorityEnabled *bool `json:"respectMsgPriorityEnabled,omitempty"`
	// Enable or disable the respecting of the time-to-live (TTL) for messages. When enabled, expired messages are discarded or moved to the DMQ. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	RespectTtlEnabled *bool `json:"respectTtlEnabled,omitempty"`
}

// NewMsgVpnQueueTemplate instantiates a new MsgVpnQueueTemplate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMsgVpnQueueTemplate() *MsgVpnQueueTemplate {
	this := MsgVpnQueueTemplate{}
	return &this
}

// NewMsgVpnQueueTemplateWithDefaults instantiates a new MsgVpnQueueTemplate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMsgVpnQueueTemplateWithDefaults() *MsgVpnQueueTemplate {
	this := MsgVpnQueueTemplate{}
	return &this
}

// GetAccessType returns the AccessType field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetAccessType() string {
	if o == nil || o.AccessType == nil {
		var ret string
		return ret
	}
	return *o.AccessType
}

// GetAccessTypeOk returns a tuple with the AccessType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetAccessTypeOk() (*string, bool) {
	if o == nil || o.AccessType == nil {
		return nil, false
	}
	return o.AccessType, true
}

// HasAccessType returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasAccessType() bool {
	if o != nil && o.AccessType != nil {
		return true
	}

	return false
}

// SetAccessType gets a reference to the given string and assigns it to the AccessType field.
func (o *MsgVpnQueueTemplate) SetAccessType(v string) {
	o.AccessType = &v
}

// GetConsumerAckPropagationEnabled returns the ConsumerAckPropagationEnabled field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetConsumerAckPropagationEnabled() bool {
	if o == nil || o.ConsumerAckPropagationEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ConsumerAckPropagationEnabled
}

// GetConsumerAckPropagationEnabledOk returns a tuple with the ConsumerAckPropagationEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetConsumerAckPropagationEnabledOk() (*bool, bool) {
	if o == nil || o.ConsumerAckPropagationEnabled == nil {
		return nil, false
	}
	return o.ConsumerAckPropagationEnabled, true
}

// HasConsumerAckPropagationEnabled returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasConsumerAckPropagationEnabled() bool {
	if o != nil && o.ConsumerAckPropagationEnabled != nil {
		return true
	}

	return false
}

// SetConsumerAckPropagationEnabled gets a reference to the given bool and assigns it to the ConsumerAckPropagationEnabled field.
func (o *MsgVpnQueueTemplate) SetConsumerAckPropagationEnabled(v bool) {
	o.ConsumerAckPropagationEnabled = &v
}

// GetDeadMsgQueue returns the DeadMsgQueue field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetDeadMsgQueue() string {
	if o == nil || o.DeadMsgQueue == nil {
		var ret string
		return ret
	}
	return *o.DeadMsgQueue
}

// GetDeadMsgQueueOk returns a tuple with the DeadMsgQueue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetDeadMsgQueueOk() (*string, bool) {
	if o == nil || o.DeadMsgQueue == nil {
		return nil, false
	}
	return o.DeadMsgQueue, true
}

// HasDeadMsgQueue returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasDeadMsgQueue() bool {
	if o != nil && o.DeadMsgQueue != nil {
		return true
	}

	return false
}

// SetDeadMsgQueue gets a reference to the given string and assigns it to the DeadMsgQueue field.
func (o *MsgVpnQueueTemplate) SetDeadMsgQueue(v string) {
	o.DeadMsgQueue = &v
}

// GetDeliveryDelay returns the DeliveryDelay field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetDeliveryDelay() int64 {
	if o == nil || o.DeliveryDelay == nil {
		var ret int64
		return ret
	}
	return *o.DeliveryDelay
}

// GetDeliveryDelayOk returns a tuple with the DeliveryDelay field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetDeliveryDelayOk() (*int64, bool) {
	if o == nil || o.DeliveryDelay == nil {
		return nil, false
	}
	return o.DeliveryDelay, true
}

// HasDeliveryDelay returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasDeliveryDelay() bool {
	if o != nil && o.DeliveryDelay != nil {
		return true
	}

	return false
}

// SetDeliveryDelay gets a reference to the given int64 and assigns it to the DeliveryDelay field.
func (o *MsgVpnQueueTemplate) SetDeliveryDelay(v int64) {
	o.DeliveryDelay = &v
}

// GetDurabilityOverride returns the DurabilityOverride field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetDurabilityOverride() string {
	if o == nil || o.DurabilityOverride == nil {
		var ret string
		return ret
	}
	return *o.DurabilityOverride
}

// GetDurabilityOverrideOk returns a tuple with the DurabilityOverride field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetDurabilityOverrideOk() (*string, bool) {
	if o == nil || o.DurabilityOverride == nil {
		return nil, false
	}
	return o.DurabilityOverride, true
}

// HasDurabilityOverride returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasDurabilityOverride() bool {
	if o != nil && o.DurabilityOverride != nil {
		return true
	}

	return false
}

// SetDurabilityOverride gets a reference to the given string and assigns it to the DurabilityOverride field.
func (o *MsgVpnQueueTemplate) SetDurabilityOverride(v string) {
	o.DurabilityOverride = &v
}

// GetEventBindCountThreshold returns the EventBindCountThreshold field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetEventBindCountThreshold() EventThreshold {
	if o == nil || o.EventBindCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventBindCountThreshold
}

// GetEventBindCountThresholdOk returns a tuple with the EventBindCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetEventBindCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventBindCountThreshold == nil {
		return nil, false
	}
	return o.EventBindCountThreshold, true
}

// HasEventBindCountThreshold returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasEventBindCountThreshold() bool {
	if o != nil && o.EventBindCountThreshold != nil {
		return true
	}

	return false
}

// SetEventBindCountThreshold gets a reference to the given EventThreshold and assigns it to the EventBindCountThreshold field.
func (o *MsgVpnQueueTemplate) SetEventBindCountThreshold(v EventThreshold) {
	o.EventBindCountThreshold = &v
}

// GetEventMsgSpoolUsageThreshold returns the EventMsgSpoolUsageThreshold field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetEventMsgSpoolUsageThreshold() EventThreshold {
	if o == nil || o.EventMsgSpoolUsageThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventMsgSpoolUsageThreshold
}

// GetEventMsgSpoolUsageThresholdOk returns a tuple with the EventMsgSpoolUsageThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetEventMsgSpoolUsageThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventMsgSpoolUsageThreshold == nil {
		return nil, false
	}
	return o.EventMsgSpoolUsageThreshold, true
}

// HasEventMsgSpoolUsageThreshold returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasEventMsgSpoolUsageThreshold() bool {
	if o != nil && o.EventMsgSpoolUsageThreshold != nil {
		return true
	}

	return false
}

// SetEventMsgSpoolUsageThreshold gets a reference to the given EventThreshold and assigns it to the EventMsgSpoolUsageThreshold field.
func (o *MsgVpnQueueTemplate) SetEventMsgSpoolUsageThreshold(v EventThreshold) {
	o.EventMsgSpoolUsageThreshold = &v
}

// GetEventRejectLowPriorityMsgLimitThreshold returns the EventRejectLowPriorityMsgLimitThreshold field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetEventRejectLowPriorityMsgLimitThreshold() EventThreshold {
	if o == nil || o.EventRejectLowPriorityMsgLimitThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventRejectLowPriorityMsgLimitThreshold
}

// GetEventRejectLowPriorityMsgLimitThresholdOk returns a tuple with the EventRejectLowPriorityMsgLimitThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetEventRejectLowPriorityMsgLimitThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventRejectLowPriorityMsgLimitThreshold == nil {
		return nil, false
	}
	return o.EventRejectLowPriorityMsgLimitThreshold, true
}

// HasEventRejectLowPriorityMsgLimitThreshold returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasEventRejectLowPriorityMsgLimitThreshold() bool {
	if o != nil && o.EventRejectLowPriorityMsgLimitThreshold != nil {
		return true
	}

	return false
}

// SetEventRejectLowPriorityMsgLimitThreshold gets a reference to the given EventThreshold and assigns it to the EventRejectLowPriorityMsgLimitThreshold field.
func (o *MsgVpnQueueTemplate) SetEventRejectLowPriorityMsgLimitThreshold(v EventThreshold) {
	o.EventRejectLowPriorityMsgLimitThreshold = &v
}

// GetMaxBindCount returns the MaxBindCount field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetMaxBindCount() int64 {
	if o == nil || o.MaxBindCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxBindCount
}

// GetMaxBindCountOk returns a tuple with the MaxBindCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetMaxBindCountOk() (*int64, bool) {
	if o == nil || o.MaxBindCount == nil {
		return nil, false
	}
	return o.MaxBindCount, true
}

// HasMaxBindCount returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasMaxBindCount() bool {
	if o != nil && o.MaxBindCount != nil {
		return true
	}

	return false
}

// SetMaxBindCount gets a reference to the given int64 and assigns it to the MaxBindCount field.
func (o *MsgVpnQueueTemplate) SetMaxBindCount(v int64) {
	o.MaxBindCount = &v
}

// GetMaxDeliveredUnackedMsgsPerFlow returns the MaxDeliveredUnackedMsgsPerFlow field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetMaxDeliveredUnackedMsgsPerFlow() int64 {
	if o == nil || o.MaxDeliveredUnackedMsgsPerFlow == nil {
		var ret int64
		return ret
	}
	return *o.MaxDeliveredUnackedMsgsPerFlow
}

// GetMaxDeliveredUnackedMsgsPerFlowOk returns a tuple with the MaxDeliveredUnackedMsgsPerFlow field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetMaxDeliveredUnackedMsgsPerFlowOk() (*int64, bool) {
	if o == nil || o.MaxDeliveredUnackedMsgsPerFlow == nil {
		return nil, false
	}
	return o.MaxDeliveredUnackedMsgsPerFlow, true
}

// HasMaxDeliveredUnackedMsgsPerFlow returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasMaxDeliveredUnackedMsgsPerFlow() bool {
	if o != nil && o.MaxDeliveredUnackedMsgsPerFlow != nil {
		return true
	}

	return false
}

// SetMaxDeliveredUnackedMsgsPerFlow gets a reference to the given int64 and assigns it to the MaxDeliveredUnackedMsgsPerFlow field.
func (o *MsgVpnQueueTemplate) SetMaxDeliveredUnackedMsgsPerFlow(v int64) {
	o.MaxDeliveredUnackedMsgsPerFlow = &v
}

// GetMaxMsgSize returns the MaxMsgSize field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetMaxMsgSize() int32 {
	if o == nil || o.MaxMsgSize == nil {
		var ret int32
		return ret
	}
	return *o.MaxMsgSize
}

// GetMaxMsgSizeOk returns a tuple with the MaxMsgSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetMaxMsgSizeOk() (*int32, bool) {
	if o == nil || o.MaxMsgSize == nil {
		return nil, false
	}
	return o.MaxMsgSize, true
}

// HasMaxMsgSize returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasMaxMsgSize() bool {
	if o != nil && o.MaxMsgSize != nil {
		return true
	}

	return false
}

// SetMaxMsgSize gets a reference to the given int32 and assigns it to the MaxMsgSize field.
func (o *MsgVpnQueueTemplate) SetMaxMsgSize(v int32) {
	o.MaxMsgSize = &v
}

// GetMaxMsgSpoolUsage returns the MaxMsgSpoolUsage field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetMaxMsgSpoolUsage() int64 {
	if o == nil || o.MaxMsgSpoolUsage == nil {
		var ret int64
		return ret
	}
	return *o.MaxMsgSpoolUsage
}

// GetMaxMsgSpoolUsageOk returns a tuple with the MaxMsgSpoolUsage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetMaxMsgSpoolUsageOk() (*int64, bool) {
	if o == nil || o.MaxMsgSpoolUsage == nil {
		return nil, false
	}
	return o.MaxMsgSpoolUsage, true
}

// HasMaxMsgSpoolUsage returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasMaxMsgSpoolUsage() bool {
	if o != nil && o.MaxMsgSpoolUsage != nil {
		return true
	}

	return false
}

// SetMaxMsgSpoolUsage gets a reference to the given int64 and assigns it to the MaxMsgSpoolUsage field.
func (o *MsgVpnQueueTemplate) SetMaxMsgSpoolUsage(v int64) {
	o.MaxMsgSpoolUsage = &v
}

// GetMaxRedeliveryCount returns the MaxRedeliveryCount field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetMaxRedeliveryCount() int64 {
	if o == nil || o.MaxRedeliveryCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxRedeliveryCount
}

// GetMaxRedeliveryCountOk returns a tuple with the MaxRedeliveryCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetMaxRedeliveryCountOk() (*int64, bool) {
	if o == nil || o.MaxRedeliveryCount == nil {
		return nil, false
	}
	return o.MaxRedeliveryCount, true
}

// HasMaxRedeliveryCount returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasMaxRedeliveryCount() bool {
	if o != nil && o.MaxRedeliveryCount != nil {
		return true
	}

	return false
}

// SetMaxRedeliveryCount gets a reference to the given int64 and assigns it to the MaxRedeliveryCount field.
func (o *MsgVpnQueueTemplate) SetMaxRedeliveryCount(v int64) {
	o.MaxRedeliveryCount = &v
}

// GetMaxTtl returns the MaxTtl field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetMaxTtl() int64 {
	if o == nil || o.MaxTtl == nil {
		var ret int64
		return ret
	}
	return *o.MaxTtl
}

// GetMaxTtlOk returns a tuple with the MaxTtl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetMaxTtlOk() (*int64, bool) {
	if o == nil || o.MaxTtl == nil {
		return nil, false
	}
	return o.MaxTtl, true
}

// HasMaxTtl returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasMaxTtl() bool {
	if o != nil && o.MaxTtl != nil {
		return true
	}

	return false
}

// SetMaxTtl gets a reference to the given int64 and assigns it to the MaxTtl field.
func (o *MsgVpnQueueTemplate) SetMaxTtl(v int64) {
	o.MaxTtl = &v
}

// GetMsgVpnName returns the MsgVpnName field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetMsgVpnName() string {
	if o == nil || o.MsgVpnName == nil {
		var ret string
		return ret
	}
	return *o.MsgVpnName
}

// GetMsgVpnNameOk returns a tuple with the MsgVpnName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetMsgVpnNameOk() (*string, bool) {
	if o == nil || o.MsgVpnName == nil {
		return nil, false
	}
	return o.MsgVpnName, true
}

// HasMsgVpnName returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasMsgVpnName() bool {
	if o != nil && o.MsgVpnName != nil {
		return true
	}

	return false
}

// SetMsgVpnName gets a reference to the given string and assigns it to the MsgVpnName field.
func (o *MsgVpnQueueTemplate) SetMsgVpnName(v string) {
	o.MsgVpnName = &v
}

// GetPermission returns the Permission field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetPermission() string {
	if o == nil || o.Permission == nil {
		var ret string
		return ret
	}
	return *o.Permission
}

// GetPermissionOk returns a tuple with the Permission field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetPermissionOk() (*string, bool) {
	if o == nil || o.Permission == nil {
		return nil, false
	}
	return o.Permission, true
}

// HasPermission returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasPermission() bool {
	if o != nil && o.Permission != nil {
		return true
	}

	return false
}

// SetPermission gets a reference to the given string and assigns it to the Permission field.
func (o *MsgVpnQueueTemplate) SetPermission(v string) {
	o.Permission = &v
}

// GetQueueNameFilter returns the QueueNameFilter field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetQueueNameFilter() string {
	if o == nil || o.QueueNameFilter == nil {
		var ret string
		return ret
	}
	return *o.QueueNameFilter
}

// GetQueueNameFilterOk returns a tuple with the QueueNameFilter field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetQueueNameFilterOk() (*string, bool) {
	if o == nil || o.QueueNameFilter == nil {
		return nil, false
	}
	return o.QueueNameFilter, true
}

// HasQueueNameFilter returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasQueueNameFilter() bool {
	if o != nil && o.QueueNameFilter != nil {
		return true
	}

	return false
}

// SetQueueNameFilter gets a reference to the given string and assigns it to the QueueNameFilter field.
func (o *MsgVpnQueueTemplate) SetQueueNameFilter(v string) {
	o.QueueNameFilter = &v
}

// GetQueueTemplateName returns the QueueTemplateName field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetQueueTemplateName() string {
	if o == nil || o.QueueTemplateName == nil {
		var ret string
		return ret
	}
	return *o.QueueTemplateName
}

// GetQueueTemplateNameOk returns a tuple with the QueueTemplateName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetQueueTemplateNameOk() (*string, bool) {
	if o == nil || o.QueueTemplateName == nil {
		return nil, false
	}
	return o.QueueTemplateName, true
}

// HasQueueTemplateName returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasQueueTemplateName() bool {
	if o != nil && o.QueueTemplateName != nil {
		return true
	}

	return false
}

// SetQueueTemplateName gets a reference to the given string and assigns it to the QueueTemplateName field.
func (o *MsgVpnQueueTemplate) SetQueueTemplateName(v string) {
	o.QueueTemplateName = &v
}

// GetRedeliveryEnabled returns the RedeliveryEnabled field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetRedeliveryEnabled() bool {
	if o == nil || o.RedeliveryEnabled == nil {
		var ret bool
		return ret
	}
	return *o.RedeliveryEnabled
}

// GetRedeliveryEnabledOk returns a tuple with the RedeliveryEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetRedeliveryEnabledOk() (*bool, bool) {
	if o == nil || o.RedeliveryEnabled == nil {
		return nil, false
	}
	return o.RedeliveryEnabled, true
}

// HasRedeliveryEnabled returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasRedeliveryEnabled() bool {
	if o != nil && o.RedeliveryEnabled != nil {
		return true
	}

	return false
}

// SetRedeliveryEnabled gets a reference to the given bool and assigns it to the RedeliveryEnabled field.
func (o *MsgVpnQueueTemplate) SetRedeliveryEnabled(v bool) {
	o.RedeliveryEnabled = &v
}

// GetRejectLowPriorityMsgEnabled returns the RejectLowPriorityMsgEnabled field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetRejectLowPriorityMsgEnabled() bool {
	if o == nil || o.RejectLowPriorityMsgEnabled == nil {
		var ret bool
		return ret
	}
	return *o.RejectLowPriorityMsgEnabled
}

// GetRejectLowPriorityMsgEnabledOk returns a tuple with the RejectLowPriorityMsgEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetRejectLowPriorityMsgEnabledOk() (*bool, bool) {
	if o == nil || o.RejectLowPriorityMsgEnabled == nil {
		return nil, false
	}
	return o.RejectLowPriorityMsgEnabled, true
}

// HasRejectLowPriorityMsgEnabled returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasRejectLowPriorityMsgEnabled() bool {
	if o != nil && o.RejectLowPriorityMsgEnabled != nil {
		return true
	}

	return false
}

// SetRejectLowPriorityMsgEnabled gets a reference to the given bool and assigns it to the RejectLowPriorityMsgEnabled field.
func (o *MsgVpnQueueTemplate) SetRejectLowPriorityMsgEnabled(v bool) {
	o.RejectLowPriorityMsgEnabled = &v
}

// GetRejectLowPriorityMsgLimit returns the RejectLowPriorityMsgLimit field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetRejectLowPriorityMsgLimit() int64 {
	if o == nil || o.RejectLowPriorityMsgLimit == nil {
		var ret int64
		return ret
	}
	return *o.RejectLowPriorityMsgLimit
}

// GetRejectLowPriorityMsgLimitOk returns a tuple with the RejectLowPriorityMsgLimit field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetRejectLowPriorityMsgLimitOk() (*int64, bool) {
	if o == nil || o.RejectLowPriorityMsgLimit == nil {
		return nil, false
	}
	return o.RejectLowPriorityMsgLimit, true
}

// HasRejectLowPriorityMsgLimit returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasRejectLowPriorityMsgLimit() bool {
	if o != nil && o.RejectLowPriorityMsgLimit != nil {
		return true
	}

	return false
}

// SetRejectLowPriorityMsgLimit gets a reference to the given int64 and assigns it to the RejectLowPriorityMsgLimit field.
func (o *MsgVpnQueueTemplate) SetRejectLowPriorityMsgLimit(v int64) {
	o.RejectLowPriorityMsgLimit = &v
}

// GetRejectMsgToSenderOnDiscardBehavior returns the RejectMsgToSenderOnDiscardBehavior field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetRejectMsgToSenderOnDiscardBehavior() string {
	if o == nil || o.RejectMsgToSenderOnDiscardBehavior == nil {
		var ret string
		return ret
	}
	return *o.RejectMsgToSenderOnDiscardBehavior
}

// GetRejectMsgToSenderOnDiscardBehaviorOk returns a tuple with the RejectMsgToSenderOnDiscardBehavior field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetRejectMsgToSenderOnDiscardBehaviorOk() (*string, bool) {
	if o == nil || o.RejectMsgToSenderOnDiscardBehavior == nil {
		return nil, false
	}
	return o.RejectMsgToSenderOnDiscardBehavior, true
}

// HasRejectMsgToSenderOnDiscardBehavior returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasRejectMsgToSenderOnDiscardBehavior() bool {
	if o != nil && o.RejectMsgToSenderOnDiscardBehavior != nil {
		return true
	}

	return false
}

// SetRejectMsgToSenderOnDiscardBehavior gets a reference to the given string and assigns it to the RejectMsgToSenderOnDiscardBehavior field.
func (o *MsgVpnQueueTemplate) SetRejectMsgToSenderOnDiscardBehavior(v string) {
	o.RejectMsgToSenderOnDiscardBehavior = &v
}

// GetRespectMsgPriorityEnabled returns the RespectMsgPriorityEnabled field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetRespectMsgPriorityEnabled() bool {
	if o == nil || o.RespectMsgPriorityEnabled == nil {
		var ret bool
		return ret
	}
	return *o.RespectMsgPriorityEnabled
}

// GetRespectMsgPriorityEnabledOk returns a tuple with the RespectMsgPriorityEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetRespectMsgPriorityEnabledOk() (*bool, bool) {
	if o == nil || o.RespectMsgPriorityEnabled == nil {
		return nil, false
	}
	return o.RespectMsgPriorityEnabled, true
}

// HasRespectMsgPriorityEnabled returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasRespectMsgPriorityEnabled() bool {
	if o != nil && o.RespectMsgPriorityEnabled != nil {
		return true
	}

	return false
}

// SetRespectMsgPriorityEnabled gets a reference to the given bool and assigns it to the RespectMsgPriorityEnabled field.
func (o *MsgVpnQueueTemplate) SetRespectMsgPriorityEnabled(v bool) {
	o.RespectMsgPriorityEnabled = &v
}

// GetRespectTtlEnabled returns the RespectTtlEnabled field value if set, zero value otherwise.
func (o *MsgVpnQueueTemplate) GetRespectTtlEnabled() bool {
	if o == nil || o.RespectTtlEnabled == nil {
		var ret bool
		return ret
	}
	return *o.RespectTtlEnabled
}

// GetRespectTtlEnabledOk returns a tuple with the RespectTtlEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnQueueTemplate) GetRespectTtlEnabledOk() (*bool, bool) {
	if o == nil || o.RespectTtlEnabled == nil {
		return nil, false
	}
	return o.RespectTtlEnabled, true
}

// HasRespectTtlEnabled returns a boolean if a field has been set.
func (o *MsgVpnQueueTemplate) HasRespectTtlEnabled() bool {
	if o != nil && o.RespectTtlEnabled != nil {
		return true
	}

	return false
}

// SetRespectTtlEnabled gets a reference to the given bool and assigns it to the RespectTtlEnabled field.
func (o *MsgVpnQueueTemplate) SetRespectTtlEnabled(v bool) {
	o.RespectTtlEnabled = &v
}

func (o MsgVpnQueueTemplate) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AccessType != nil {
		toSerialize["accessType"] = o.AccessType
	}
	if o.ConsumerAckPropagationEnabled != nil {
		toSerialize["consumerAckPropagationEnabled"] = o.ConsumerAckPropagationEnabled
	}
	if o.DeadMsgQueue != nil {
		toSerialize["deadMsgQueue"] = o.DeadMsgQueue
	}
	if o.DeliveryDelay != nil {
		toSerialize["deliveryDelay"] = o.DeliveryDelay
	}
	if o.DurabilityOverride != nil {
		toSerialize["durabilityOverride"] = o.DurabilityOverride
	}
	if o.EventBindCountThreshold != nil {
		toSerialize["eventBindCountThreshold"] = o.EventBindCountThreshold
	}
	if o.EventMsgSpoolUsageThreshold != nil {
		toSerialize["eventMsgSpoolUsageThreshold"] = o.EventMsgSpoolUsageThreshold
	}
	if o.EventRejectLowPriorityMsgLimitThreshold != nil {
		toSerialize["eventRejectLowPriorityMsgLimitThreshold"] = o.EventRejectLowPriorityMsgLimitThreshold
	}
	if o.MaxBindCount != nil {
		toSerialize["maxBindCount"] = o.MaxBindCount
	}
	if o.MaxDeliveredUnackedMsgsPerFlow != nil {
		toSerialize["maxDeliveredUnackedMsgsPerFlow"] = o.MaxDeliveredUnackedMsgsPerFlow
	}
	if o.MaxMsgSize != nil {
		toSerialize["maxMsgSize"] = o.MaxMsgSize
	}
	if o.MaxMsgSpoolUsage != nil {
		toSerialize["maxMsgSpoolUsage"] = o.MaxMsgSpoolUsage
	}
	if o.MaxRedeliveryCount != nil {
		toSerialize["maxRedeliveryCount"] = o.MaxRedeliveryCount
	}
	if o.MaxTtl != nil {
		toSerialize["maxTtl"] = o.MaxTtl
	}
	if o.MsgVpnName != nil {
		toSerialize["msgVpnName"] = o.MsgVpnName
	}
	if o.Permission != nil {
		toSerialize["permission"] = o.Permission
	}
	if o.QueueNameFilter != nil {
		toSerialize["queueNameFilter"] = o.QueueNameFilter
	}
	if o.QueueTemplateName != nil {
		toSerialize["queueTemplateName"] = o.QueueTemplateName
	}
	if o.RedeliveryEnabled != nil {
		toSerialize["redeliveryEnabled"] = o.RedeliveryEnabled
	}
	if o.RejectLowPriorityMsgEnabled != nil {
		toSerialize["rejectLowPriorityMsgEnabled"] = o.RejectLowPriorityMsgEnabled
	}
	if o.RejectLowPriorityMsgLimit != nil {
		toSerialize["rejectLowPriorityMsgLimit"] = o.RejectLowPriorityMsgLimit
	}
	if o.RejectMsgToSenderOnDiscardBehavior != nil {
		toSerialize["rejectMsgToSenderOnDiscardBehavior"] = o.RejectMsgToSenderOnDiscardBehavior
	}
	if o.RespectMsgPriorityEnabled != nil {
		toSerialize["respectMsgPriorityEnabled"] = o.RespectMsgPriorityEnabled
	}
	if o.RespectTtlEnabled != nil {
		toSerialize["respectTtlEnabled"] = o.RespectTtlEnabled
	}
	return json.Marshal(toSerialize)
}

type NullableMsgVpnQueueTemplate struct {
	value *MsgVpnQueueTemplate
	isSet bool
}

func (v NullableMsgVpnQueueTemplate) Get() *MsgVpnQueueTemplate {
	return v.value
}

func (v *NullableMsgVpnQueueTemplate) Set(val *MsgVpnQueueTemplate) {
	v.value = val
	v.isSet = true
}

func (v NullableMsgVpnQueueTemplate) IsSet() bool {
	return v.isSet
}

func (v *NullableMsgVpnQueueTemplate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMsgVpnQueueTemplate(val *MsgVpnQueueTemplate) *NullableMsgVpnQueueTemplate {
	return &NullableMsgVpnQueueTemplate{value: val, isSet: true}
}

func (v NullableMsgVpnQueueTemplate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMsgVpnQueueTemplate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
