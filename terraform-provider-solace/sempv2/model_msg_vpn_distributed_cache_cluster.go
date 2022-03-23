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

// MsgVpnDistributedCacheCluster struct for MsgVpnDistributedCacheCluster
type MsgVpnDistributedCacheCluster struct {
	// The name of the Distributed Cache.
	CacheName *string `json:"cacheName,omitempty"`
	// The name of the Cache Cluster.
	ClusterName *string `json:"clusterName,omitempty"`
	// Enable or disable deliver-to-one override for the Cache Cluster. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	DeliverToOneOverrideEnabled *bool `json:"deliverToOneOverrideEnabled,omitempty"`
	// Enable or disable the Cache Cluster. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	Enabled                         *bool                    `json:"enabled,omitempty"`
	EventDataByteRateThreshold      *EventThresholdByValue   `json:"eventDataByteRateThreshold,omitempty"`
	EventDataMsgRateThreshold       *EventThresholdByValue   `json:"eventDataMsgRateThreshold,omitempty"`
	EventMaxMemoryThreshold         *EventThresholdByPercent `json:"eventMaxMemoryThreshold,omitempty"`
	EventMaxTopicsThreshold         *EventThresholdByPercent `json:"eventMaxTopicsThreshold,omitempty"`
	EventRequestQueueDepthThreshold *EventThresholdByPercent `json:"eventRequestQueueDepthThreshold,omitempty"`
	EventRequestRateThreshold       *EventThresholdByValue   `json:"eventRequestRateThreshold,omitempty"`
	EventResponseRateThreshold      *EventThresholdByValue   `json:"eventResponseRateThreshold,omitempty"`
	// Enable or disable global caching for the Cache Cluster. When enabled, the Cache Instances will fetch topics from remote Home Cache Clusters when requested, and subscribe to those topics to cache them locally. When disabled, the Cache Instances will remove all subscriptions and cached messages for topics from remote Home Cache Clusters. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	GlobalCachingEnabled *bool `json:"globalCachingEnabled,omitempty"`
	// The heartbeat interval, in seconds, used by the Cache Instances to monitor connectivity with the remote Home Cache Clusters. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.
	GlobalCachingHeartbeat *int64 `json:"globalCachingHeartbeat,omitempty"`
	// The topic lifetime, in seconds. If no client requests are received for a given global topic over the duration of the topic lifetime, then the Cache Instance will remove the subscription and cached messages for that topic. A value of 0 disables aging. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3600`.
	GlobalCachingTopicLifetime *int64 `json:"globalCachingTopicLifetime,omitempty"`
	// The maximum memory usage, in megabytes (MB), for each Cache Instance in the Cache Cluster. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `2048`.
	MaxMemory *int64 `json:"maxMemory,omitempty"`
	// The maximum number of messages per topic for each Cache Instance in the Cache Cluster. When at the maximum, old messages are removed as new messages arrive. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1`.
	MaxMsgsPerTopic *int64 `json:"maxMsgsPerTopic,omitempty"`
	// The maximum queue depth for cache requests received by the Cache Cluster. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `100000`.
	MaxRequestQueueDepth *int64 `json:"maxRequestQueueDepth,omitempty"`
	// The maximum number of topics for each Cache Instance in the Cache Cluster. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `2000000`.
	MaxTopicCount *int64 `json:"maxTopicCount,omitempty"`
	// The message lifetime, in seconds. If a message remains cached for the duration of its lifetime, the Cache Instance will remove the message. A lifetime of 0 results in the message being retained indefinitely. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.
	MsgLifetime *int64 `json:"msgLifetime,omitempty"`
	// The name of the Message VPN.
	MsgVpnName *string `json:"msgVpnName,omitempty"`
	// Enable or disable the advertising, onto the message bus, of new topics learned by each Cache Instance in the Cache Cluster. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	NewTopicAdvertisementEnabled *bool `json:"newTopicAdvertisementEnabled,omitempty"`
}

// NewMsgVpnDistributedCacheCluster instantiates a new MsgVpnDistributedCacheCluster object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMsgVpnDistributedCacheCluster() *MsgVpnDistributedCacheCluster {
	this := MsgVpnDistributedCacheCluster{}
	return &this
}

// NewMsgVpnDistributedCacheClusterWithDefaults instantiates a new MsgVpnDistributedCacheCluster object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMsgVpnDistributedCacheClusterWithDefaults() *MsgVpnDistributedCacheCluster {
	this := MsgVpnDistributedCacheCluster{}
	return &this
}

// GetCacheName returns the CacheName field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetCacheName() string {
	if o == nil || o.CacheName == nil {
		var ret string
		return ret
	}
	return *o.CacheName
}

// GetCacheNameOk returns a tuple with the CacheName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetCacheNameOk() (*string, bool) {
	if o == nil || o.CacheName == nil {
		return nil, false
	}
	return o.CacheName, true
}

// HasCacheName returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasCacheName() bool {
	if o != nil && o.CacheName != nil {
		return true
	}

	return false
}

// SetCacheName gets a reference to the given string and assigns it to the CacheName field.
func (o *MsgVpnDistributedCacheCluster) SetCacheName(v string) {
	o.CacheName = &v
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetClusterName() string {
	if o == nil || o.ClusterName == nil {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetClusterNameOk() (*string, bool) {
	if o == nil || o.ClusterName == nil {
		return nil, false
	}
	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasClusterName() bool {
	if o != nil && o.ClusterName != nil {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *MsgVpnDistributedCacheCluster) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetDeliverToOneOverrideEnabled returns the DeliverToOneOverrideEnabled field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetDeliverToOneOverrideEnabled() bool {
	if o == nil || o.DeliverToOneOverrideEnabled == nil {
		var ret bool
		return ret
	}
	return *o.DeliverToOneOverrideEnabled
}

// GetDeliverToOneOverrideEnabledOk returns a tuple with the DeliverToOneOverrideEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetDeliverToOneOverrideEnabledOk() (*bool, bool) {
	if o == nil || o.DeliverToOneOverrideEnabled == nil {
		return nil, false
	}
	return o.DeliverToOneOverrideEnabled, true
}

// HasDeliverToOneOverrideEnabled returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasDeliverToOneOverrideEnabled() bool {
	if o != nil && o.DeliverToOneOverrideEnabled != nil {
		return true
	}

	return false
}

// SetDeliverToOneOverrideEnabled gets a reference to the given bool and assigns it to the DeliverToOneOverrideEnabled field.
func (o *MsgVpnDistributedCacheCluster) SetDeliverToOneOverrideEnabled(v bool) {
	o.DeliverToOneOverrideEnabled = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetEnabled() bool {
	if o == nil || o.Enabled == nil {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetEnabledOk() (*bool, bool) {
	if o == nil || o.Enabled == nil {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasEnabled() bool {
	if o != nil && o.Enabled != nil {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *MsgVpnDistributedCacheCluster) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetEventDataByteRateThreshold returns the EventDataByteRateThreshold field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetEventDataByteRateThreshold() EventThresholdByValue {
	if o == nil || o.EventDataByteRateThreshold == nil {
		var ret EventThresholdByValue
		return ret
	}
	return *o.EventDataByteRateThreshold
}

// GetEventDataByteRateThresholdOk returns a tuple with the EventDataByteRateThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetEventDataByteRateThresholdOk() (*EventThresholdByValue, bool) {
	if o == nil || o.EventDataByteRateThreshold == nil {
		return nil, false
	}
	return o.EventDataByteRateThreshold, true
}

// HasEventDataByteRateThreshold returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasEventDataByteRateThreshold() bool {
	if o != nil && o.EventDataByteRateThreshold != nil {
		return true
	}

	return false
}

// SetEventDataByteRateThreshold gets a reference to the given EventThresholdByValue and assigns it to the EventDataByteRateThreshold field.
func (o *MsgVpnDistributedCacheCluster) SetEventDataByteRateThreshold(v EventThresholdByValue) {
	o.EventDataByteRateThreshold = &v
}

// GetEventDataMsgRateThreshold returns the EventDataMsgRateThreshold field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetEventDataMsgRateThreshold() EventThresholdByValue {
	if o == nil || o.EventDataMsgRateThreshold == nil {
		var ret EventThresholdByValue
		return ret
	}
	return *o.EventDataMsgRateThreshold
}

// GetEventDataMsgRateThresholdOk returns a tuple with the EventDataMsgRateThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetEventDataMsgRateThresholdOk() (*EventThresholdByValue, bool) {
	if o == nil || o.EventDataMsgRateThreshold == nil {
		return nil, false
	}
	return o.EventDataMsgRateThreshold, true
}

// HasEventDataMsgRateThreshold returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasEventDataMsgRateThreshold() bool {
	if o != nil && o.EventDataMsgRateThreshold != nil {
		return true
	}

	return false
}

// SetEventDataMsgRateThreshold gets a reference to the given EventThresholdByValue and assigns it to the EventDataMsgRateThreshold field.
func (o *MsgVpnDistributedCacheCluster) SetEventDataMsgRateThreshold(v EventThresholdByValue) {
	o.EventDataMsgRateThreshold = &v
}

// GetEventMaxMemoryThreshold returns the EventMaxMemoryThreshold field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetEventMaxMemoryThreshold() EventThresholdByPercent {
	if o == nil || o.EventMaxMemoryThreshold == nil {
		var ret EventThresholdByPercent
		return ret
	}
	return *o.EventMaxMemoryThreshold
}

// GetEventMaxMemoryThresholdOk returns a tuple with the EventMaxMemoryThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetEventMaxMemoryThresholdOk() (*EventThresholdByPercent, bool) {
	if o == nil || o.EventMaxMemoryThreshold == nil {
		return nil, false
	}
	return o.EventMaxMemoryThreshold, true
}

// HasEventMaxMemoryThreshold returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasEventMaxMemoryThreshold() bool {
	if o != nil && o.EventMaxMemoryThreshold != nil {
		return true
	}

	return false
}

// SetEventMaxMemoryThreshold gets a reference to the given EventThresholdByPercent and assigns it to the EventMaxMemoryThreshold field.
func (o *MsgVpnDistributedCacheCluster) SetEventMaxMemoryThreshold(v EventThresholdByPercent) {
	o.EventMaxMemoryThreshold = &v
}

// GetEventMaxTopicsThreshold returns the EventMaxTopicsThreshold field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetEventMaxTopicsThreshold() EventThresholdByPercent {
	if o == nil || o.EventMaxTopicsThreshold == nil {
		var ret EventThresholdByPercent
		return ret
	}
	return *o.EventMaxTopicsThreshold
}

// GetEventMaxTopicsThresholdOk returns a tuple with the EventMaxTopicsThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetEventMaxTopicsThresholdOk() (*EventThresholdByPercent, bool) {
	if o == nil || o.EventMaxTopicsThreshold == nil {
		return nil, false
	}
	return o.EventMaxTopicsThreshold, true
}

// HasEventMaxTopicsThreshold returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasEventMaxTopicsThreshold() bool {
	if o != nil && o.EventMaxTopicsThreshold != nil {
		return true
	}

	return false
}

// SetEventMaxTopicsThreshold gets a reference to the given EventThresholdByPercent and assigns it to the EventMaxTopicsThreshold field.
func (o *MsgVpnDistributedCacheCluster) SetEventMaxTopicsThreshold(v EventThresholdByPercent) {
	o.EventMaxTopicsThreshold = &v
}

// GetEventRequestQueueDepthThreshold returns the EventRequestQueueDepthThreshold field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetEventRequestQueueDepthThreshold() EventThresholdByPercent {
	if o == nil || o.EventRequestQueueDepthThreshold == nil {
		var ret EventThresholdByPercent
		return ret
	}
	return *o.EventRequestQueueDepthThreshold
}

// GetEventRequestQueueDepthThresholdOk returns a tuple with the EventRequestQueueDepthThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetEventRequestQueueDepthThresholdOk() (*EventThresholdByPercent, bool) {
	if o == nil || o.EventRequestQueueDepthThreshold == nil {
		return nil, false
	}
	return o.EventRequestQueueDepthThreshold, true
}

// HasEventRequestQueueDepthThreshold returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasEventRequestQueueDepthThreshold() bool {
	if o != nil && o.EventRequestQueueDepthThreshold != nil {
		return true
	}

	return false
}

// SetEventRequestQueueDepthThreshold gets a reference to the given EventThresholdByPercent and assigns it to the EventRequestQueueDepthThreshold field.
func (o *MsgVpnDistributedCacheCluster) SetEventRequestQueueDepthThreshold(v EventThresholdByPercent) {
	o.EventRequestQueueDepthThreshold = &v
}

// GetEventRequestRateThreshold returns the EventRequestRateThreshold field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetEventRequestRateThreshold() EventThresholdByValue {
	if o == nil || o.EventRequestRateThreshold == nil {
		var ret EventThresholdByValue
		return ret
	}
	return *o.EventRequestRateThreshold
}

// GetEventRequestRateThresholdOk returns a tuple with the EventRequestRateThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetEventRequestRateThresholdOk() (*EventThresholdByValue, bool) {
	if o == nil || o.EventRequestRateThreshold == nil {
		return nil, false
	}
	return o.EventRequestRateThreshold, true
}

// HasEventRequestRateThreshold returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasEventRequestRateThreshold() bool {
	if o != nil && o.EventRequestRateThreshold != nil {
		return true
	}

	return false
}

// SetEventRequestRateThreshold gets a reference to the given EventThresholdByValue and assigns it to the EventRequestRateThreshold field.
func (o *MsgVpnDistributedCacheCluster) SetEventRequestRateThreshold(v EventThresholdByValue) {
	o.EventRequestRateThreshold = &v
}

// GetEventResponseRateThreshold returns the EventResponseRateThreshold field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetEventResponseRateThreshold() EventThresholdByValue {
	if o == nil || o.EventResponseRateThreshold == nil {
		var ret EventThresholdByValue
		return ret
	}
	return *o.EventResponseRateThreshold
}

// GetEventResponseRateThresholdOk returns a tuple with the EventResponseRateThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetEventResponseRateThresholdOk() (*EventThresholdByValue, bool) {
	if o == nil || o.EventResponseRateThreshold == nil {
		return nil, false
	}
	return o.EventResponseRateThreshold, true
}

// HasEventResponseRateThreshold returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasEventResponseRateThreshold() bool {
	if o != nil && o.EventResponseRateThreshold != nil {
		return true
	}

	return false
}

// SetEventResponseRateThreshold gets a reference to the given EventThresholdByValue and assigns it to the EventResponseRateThreshold field.
func (o *MsgVpnDistributedCacheCluster) SetEventResponseRateThreshold(v EventThresholdByValue) {
	o.EventResponseRateThreshold = &v
}

// GetGlobalCachingEnabled returns the GlobalCachingEnabled field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetGlobalCachingEnabled() bool {
	if o == nil || o.GlobalCachingEnabled == nil {
		var ret bool
		return ret
	}
	return *o.GlobalCachingEnabled
}

// GetGlobalCachingEnabledOk returns a tuple with the GlobalCachingEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetGlobalCachingEnabledOk() (*bool, bool) {
	if o == nil || o.GlobalCachingEnabled == nil {
		return nil, false
	}
	return o.GlobalCachingEnabled, true
}

// HasGlobalCachingEnabled returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasGlobalCachingEnabled() bool {
	if o != nil && o.GlobalCachingEnabled != nil {
		return true
	}

	return false
}

// SetGlobalCachingEnabled gets a reference to the given bool and assigns it to the GlobalCachingEnabled field.
func (o *MsgVpnDistributedCacheCluster) SetGlobalCachingEnabled(v bool) {
	o.GlobalCachingEnabled = &v
}

// GetGlobalCachingHeartbeat returns the GlobalCachingHeartbeat field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetGlobalCachingHeartbeat() int64 {
	if o == nil || o.GlobalCachingHeartbeat == nil {
		var ret int64
		return ret
	}
	return *o.GlobalCachingHeartbeat
}

// GetGlobalCachingHeartbeatOk returns a tuple with the GlobalCachingHeartbeat field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetGlobalCachingHeartbeatOk() (*int64, bool) {
	if o == nil || o.GlobalCachingHeartbeat == nil {
		return nil, false
	}
	return o.GlobalCachingHeartbeat, true
}

// HasGlobalCachingHeartbeat returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasGlobalCachingHeartbeat() bool {
	if o != nil && o.GlobalCachingHeartbeat != nil {
		return true
	}

	return false
}

// SetGlobalCachingHeartbeat gets a reference to the given int64 and assigns it to the GlobalCachingHeartbeat field.
func (o *MsgVpnDistributedCacheCluster) SetGlobalCachingHeartbeat(v int64) {
	o.GlobalCachingHeartbeat = &v
}

// GetGlobalCachingTopicLifetime returns the GlobalCachingTopicLifetime field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetGlobalCachingTopicLifetime() int64 {
	if o == nil || o.GlobalCachingTopicLifetime == nil {
		var ret int64
		return ret
	}
	return *o.GlobalCachingTopicLifetime
}

// GetGlobalCachingTopicLifetimeOk returns a tuple with the GlobalCachingTopicLifetime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetGlobalCachingTopicLifetimeOk() (*int64, bool) {
	if o == nil || o.GlobalCachingTopicLifetime == nil {
		return nil, false
	}
	return o.GlobalCachingTopicLifetime, true
}

// HasGlobalCachingTopicLifetime returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasGlobalCachingTopicLifetime() bool {
	if o != nil && o.GlobalCachingTopicLifetime != nil {
		return true
	}

	return false
}

// SetGlobalCachingTopicLifetime gets a reference to the given int64 and assigns it to the GlobalCachingTopicLifetime field.
func (o *MsgVpnDistributedCacheCluster) SetGlobalCachingTopicLifetime(v int64) {
	o.GlobalCachingTopicLifetime = &v
}

// GetMaxMemory returns the MaxMemory field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetMaxMemory() int64 {
	if o == nil || o.MaxMemory == nil {
		var ret int64
		return ret
	}
	return *o.MaxMemory
}

// GetMaxMemoryOk returns a tuple with the MaxMemory field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetMaxMemoryOk() (*int64, bool) {
	if o == nil || o.MaxMemory == nil {
		return nil, false
	}
	return o.MaxMemory, true
}

// HasMaxMemory returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasMaxMemory() bool {
	if o != nil && o.MaxMemory != nil {
		return true
	}

	return false
}

// SetMaxMemory gets a reference to the given int64 and assigns it to the MaxMemory field.
func (o *MsgVpnDistributedCacheCluster) SetMaxMemory(v int64) {
	o.MaxMemory = &v
}

// GetMaxMsgsPerTopic returns the MaxMsgsPerTopic field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetMaxMsgsPerTopic() int64 {
	if o == nil || o.MaxMsgsPerTopic == nil {
		var ret int64
		return ret
	}
	return *o.MaxMsgsPerTopic
}

// GetMaxMsgsPerTopicOk returns a tuple with the MaxMsgsPerTopic field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetMaxMsgsPerTopicOk() (*int64, bool) {
	if o == nil || o.MaxMsgsPerTopic == nil {
		return nil, false
	}
	return o.MaxMsgsPerTopic, true
}

// HasMaxMsgsPerTopic returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasMaxMsgsPerTopic() bool {
	if o != nil && o.MaxMsgsPerTopic != nil {
		return true
	}

	return false
}

// SetMaxMsgsPerTopic gets a reference to the given int64 and assigns it to the MaxMsgsPerTopic field.
func (o *MsgVpnDistributedCacheCluster) SetMaxMsgsPerTopic(v int64) {
	o.MaxMsgsPerTopic = &v
}

// GetMaxRequestQueueDepth returns the MaxRequestQueueDepth field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetMaxRequestQueueDepth() int64 {
	if o == nil || o.MaxRequestQueueDepth == nil {
		var ret int64
		return ret
	}
	return *o.MaxRequestQueueDepth
}

// GetMaxRequestQueueDepthOk returns a tuple with the MaxRequestQueueDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetMaxRequestQueueDepthOk() (*int64, bool) {
	if o == nil || o.MaxRequestQueueDepth == nil {
		return nil, false
	}
	return o.MaxRequestQueueDepth, true
}

// HasMaxRequestQueueDepth returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasMaxRequestQueueDepth() bool {
	if o != nil && o.MaxRequestQueueDepth != nil {
		return true
	}

	return false
}

// SetMaxRequestQueueDepth gets a reference to the given int64 and assigns it to the MaxRequestQueueDepth field.
func (o *MsgVpnDistributedCacheCluster) SetMaxRequestQueueDepth(v int64) {
	o.MaxRequestQueueDepth = &v
}

// GetMaxTopicCount returns the MaxTopicCount field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetMaxTopicCount() int64 {
	if o == nil || o.MaxTopicCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxTopicCount
}

// GetMaxTopicCountOk returns a tuple with the MaxTopicCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetMaxTopicCountOk() (*int64, bool) {
	if o == nil || o.MaxTopicCount == nil {
		return nil, false
	}
	return o.MaxTopicCount, true
}

// HasMaxTopicCount returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasMaxTopicCount() bool {
	if o != nil && o.MaxTopicCount != nil {
		return true
	}

	return false
}

// SetMaxTopicCount gets a reference to the given int64 and assigns it to the MaxTopicCount field.
func (o *MsgVpnDistributedCacheCluster) SetMaxTopicCount(v int64) {
	o.MaxTopicCount = &v
}

// GetMsgLifetime returns the MsgLifetime field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetMsgLifetime() int64 {
	if o == nil || o.MsgLifetime == nil {
		var ret int64
		return ret
	}
	return *o.MsgLifetime
}

// GetMsgLifetimeOk returns a tuple with the MsgLifetime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetMsgLifetimeOk() (*int64, bool) {
	if o == nil || o.MsgLifetime == nil {
		return nil, false
	}
	return o.MsgLifetime, true
}

// HasMsgLifetime returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasMsgLifetime() bool {
	if o != nil && o.MsgLifetime != nil {
		return true
	}

	return false
}

// SetMsgLifetime gets a reference to the given int64 and assigns it to the MsgLifetime field.
func (o *MsgVpnDistributedCacheCluster) SetMsgLifetime(v int64) {
	o.MsgLifetime = &v
}

// GetMsgVpnName returns the MsgVpnName field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetMsgVpnName() string {
	if o == nil || o.MsgVpnName == nil {
		var ret string
		return ret
	}
	return *o.MsgVpnName
}

// GetMsgVpnNameOk returns a tuple with the MsgVpnName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetMsgVpnNameOk() (*string, bool) {
	if o == nil || o.MsgVpnName == nil {
		return nil, false
	}
	return o.MsgVpnName, true
}

// HasMsgVpnName returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasMsgVpnName() bool {
	if o != nil && o.MsgVpnName != nil {
		return true
	}

	return false
}

// SetMsgVpnName gets a reference to the given string and assigns it to the MsgVpnName field.
func (o *MsgVpnDistributedCacheCluster) SetMsgVpnName(v string) {
	o.MsgVpnName = &v
}

// GetNewTopicAdvertisementEnabled returns the NewTopicAdvertisementEnabled field value if set, zero value otherwise.
func (o *MsgVpnDistributedCacheCluster) GetNewTopicAdvertisementEnabled() bool {
	if o == nil || o.NewTopicAdvertisementEnabled == nil {
		var ret bool
		return ret
	}
	return *o.NewTopicAdvertisementEnabled
}

// GetNewTopicAdvertisementEnabledOk returns a tuple with the NewTopicAdvertisementEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnDistributedCacheCluster) GetNewTopicAdvertisementEnabledOk() (*bool, bool) {
	if o == nil || o.NewTopicAdvertisementEnabled == nil {
		return nil, false
	}
	return o.NewTopicAdvertisementEnabled, true
}

// HasNewTopicAdvertisementEnabled returns a boolean if a field has been set.
func (o *MsgVpnDistributedCacheCluster) HasNewTopicAdvertisementEnabled() bool {
	if o != nil && o.NewTopicAdvertisementEnabled != nil {
		return true
	}

	return false
}

// SetNewTopicAdvertisementEnabled gets a reference to the given bool and assigns it to the NewTopicAdvertisementEnabled field.
func (o *MsgVpnDistributedCacheCluster) SetNewTopicAdvertisementEnabled(v bool) {
	o.NewTopicAdvertisementEnabled = &v
}

func (o MsgVpnDistributedCacheCluster) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.CacheName != nil {
		toSerialize["cacheName"] = o.CacheName
	}
	if o.ClusterName != nil {
		toSerialize["clusterName"] = o.ClusterName
	}
	if o.DeliverToOneOverrideEnabled != nil {
		toSerialize["deliverToOneOverrideEnabled"] = o.DeliverToOneOverrideEnabled
	}
	if o.Enabled != nil {
		toSerialize["enabled"] = o.Enabled
	}
	if o.EventDataByteRateThreshold != nil {
		toSerialize["eventDataByteRateThreshold"] = o.EventDataByteRateThreshold
	}
	if o.EventDataMsgRateThreshold != nil {
		toSerialize["eventDataMsgRateThreshold"] = o.EventDataMsgRateThreshold
	}
	if o.EventMaxMemoryThreshold != nil {
		toSerialize["eventMaxMemoryThreshold"] = o.EventMaxMemoryThreshold
	}
	if o.EventMaxTopicsThreshold != nil {
		toSerialize["eventMaxTopicsThreshold"] = o.EventMaxTopicsThreshold
	}
	if o.EventRequestQueueDepthThreshold != nil {
		toSerialize["eventRequestQueueDepthThreshold"] = o.EventRequestQueueDepthThreshold
	}
	if o.EventRequestRateThreshold != nil {
		toSerialize["eventRequestRateThreshold"] = o.EventRequestRateThreshold
	}
	if o.EventResponseRateThreshold != nil {
		toSerialize["eventResponseRateThreshold"] = o.EventResponseRateThreshold
	}
	if o.GlobalCachingEnabled != nil {
		toSerialize["globalCachingEnabled"] = o.GlobalCachingEnabled
	}
	if o.GlobalCachingHeartbeat != nil {
		toSerialize["globalCachingHeartbeat"] = o.GlobalCachingHeartbeat
	}
	if o.GlobalCachingTopicLifetime != nil {
		toSerialize["globalCachingTopicLifetime"] = o.GlobalCachingTopicLifetime
	}
	if o.MaxMemory != nil {
		toSerialize["maxMemory"] = o.MaxMemory
	}
	if o.MaxMsgsPerTopic != nil {
		toSerialize["maxMsgsPerTopic"] = o.MaxMsgsPerTopic
	}
	if o.MaxRequestQueueDepth != nil {
		toSerialize["maxRequestQueueDepth"] = o.MaxRequestQueueDepth
	}
	if o.MaxTopicCount != nil {
		toSerialize["maxTopicCount"] = o.MaxTopicCount
	}
	if o.MsgLifetime != nil {
		toSerialize["msgLifetime"] = o.MsgLifetime
	}
	if o.MsgVpnName != nil {
		toSerialize["msgVpnName"] = o.MsgVpnName
	}
	if o.NewTopicAdvertisementEnabled != nil {
		toSerialize["newTopicAdvertisementEnabled"] = o.NewTopicAdvertisementEnabled
	}
	return json.Marshal(toSerialize)
}

type NullableMsgVpnDistributedCacheCluster struct {
	value *MsgVpnDistributedCacheCluster
	isSet bool
}

func (v NullableMsgVpnDistributedCacheCluster) Get() *MsgVpnDistributedCacheCluster {
	return v.value
}

func (v *NullableMsgVpnDistributedCacheCluster) Set(val *MsgVpnDistributedCacheCluster) {
	v.value = val
	v.isSet = true
}

func (v NullableMsgVpnDistributedCacheCluster) IsSet() bool {
	return v.isSet
}

func (v *NullableMsgVpnDistributedCacheCluster) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMsgVpnDistributedCacheCluster(val *MsgVpnDistributedCacheCluster) *NullableMsgVpnDistributedCacheCluster {
	return &NullableMsgVpnDistributedCacheCluster{value: val, isSet: true}
}

func (v NullableMsgVpnDistributedCacheCluster) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMsgVpnDistributedCacheCluster) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
