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

// MsgVpnLinks struct for MsgVpnLinks
type MsgVpnLinks struct {
	// The URI of this Message VPN's collection of ACL Profile objects.
	AclProfilesUri *string `json:"aclProfilesUri,omitempty"`
	// The URI of this Message VPN's collection of OAuth Profile objects. Available since 2.25.
	AuthenticationOauthProfilesUri *string `json:"authenticationOauthProfilesUri,omitempty"`
	// The URI of this Message VPN's collection of OAuth Provider objects. Deprecated since 2.25. Replaced by authenticationOauthProfiles.
	AuthenticationOauthProvidersUri *string `json:"authenticationOauthProvidersUri,omitempty"`
	// The URI of this Message VPN's collection of Authorization Group objects.
	AuthorizationGroupsUri *string `json:"authorizationGroupsUri,omitempty"`
	// The URI of this Message VPN's collection of Bridge objects.
	BridgesUri *string `json:"bridgesUri,omitempty"`
	// The URI of this Message VPN's collection of Client Profile objects.
	ClientProfilesUri *string `json:"clientProfilesUri,omitempty"`
	// The URI of this Message VPN's collection of Client Username objects.
	ClientUsernamesUri *string `json:"clientUsernamesUri,omitempty"`
	// The URI of this Message VPN's collection of Distributed Cache objects. Available since 2.11.
	DistributedCachesUri *string `json:"distributedCachesUri,omitempty"`
	// The URI of this Message VPN's collection of DMR Bridge objects. Available since 2.11.
	DmrBridgesUri *string `json:"dmrBridgesUri,omitempty"`
	// The URI of this Message VPN's collection of JNDI Connection Factory objects. Available since 2.2.
	JndiConnectionFactoriesUri *string `json:"jndiConnectionFactoriesUri,omitempty"`
	// The URI of this Message VPN's collection of JNDI Queue objects. Available since 2.2.
	JndiQueuesUri *string `json:"jndiQueuesUri,omitempty"`
	// The URI of this Message VPN's collection of JNDI Topic objects. Available since 2.2.
	JndiTopicsUri *string `json:"jndiTopicsUri,omitempty"`
	// The URI of this Message VPN's collection of MQTT Retain Cache objects. Available since 2.11.
	MqttRetainCachesUri *string `json:"mqttRetainCachesUri,omitempty"`
	// The URI of this Message VPN's collection of MQTT Session objects. Available since 2.1.
	MqttSessionsUri *string `json:"mqttSessionsUri,omitempty"`
	// The URI of this Message VPN's collection of Queue Template objects. Available since 2.14.
	QueueTemplatesUri *string `json:"queueTemplatesUri,omitempty"`
	// The URI of this Message VPN's collection of Queue objects.
	QueuesUri *string `json:"queuesUri,omitempty"`
	// The URI of this Message VPN's collection of Replay Log objects. Available since 2.10.
	ReplayLogsUri *string `json:"replayLogsUri,omitempty"`
	// The URI of this Message VPN's collection of Replicated Topic objects. Available since 2.1.
	ReplicatedTopicsUri *string `json:"replicatedTopicsUri,omitempty"`
	// The URI of this Message VPN's collection of REST Delivery Point objects.
	RestDeliveryPointsUri *string `json:"restDeliveryPointsUri,omitempty"`
	// The URI of this Message VPN's collection of Sequenced Topic objects.
	SequencedTopicsUri *string `json:"sequencedTopicsUri,omitempty"`
	// The URI of this Message VPN's collection of Topic Endpoint Template objects. Available since 2.14.
	TopicEndpointTemplatesUri *string `json:"topicEndpointTemplatesUri,omitempty"`
	// The URI of this Message VPN's collection of Topic Endpoint objects. Available since 2.1.
	TopicEndpointsUri *string `json:"topicEndpointsUri,omitempty"`
	// The URI of this Message VPN object.
	Uri *string `json:"uri,omitempty"`
}

// NewMsgVpnLinks instantiates a new MsgVpnLinks object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMsgVpnLinks() *MsgVpnLinks {
	this := MsgVpnLinks{}
	return &this
}

// NewMsgVpnLinksWithDefaults instantiates a new MsgVpnLinks object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMsgVpnLinksWithDefaults() *MsgVpnLinks {
	this := MsgVpnLinks{}
	return &this
}

// GetAclProfilesUri returns the AclProfilesUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetAclProfilesUri() string {
	if o == nil || o.AclProfilesUri == nil {
		var ret string
		return ret
	}
	return *o.AclProfilesUri
}

// GetAclProfilesUriOk returns a tuple with the AclProfilesUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetAclProfilesUriOk() (*string, bool) {
	if o == nil || o.AclProfilesUri == nil {
		return nil, false
	}
	return o.AclProfilesUri, true
}

// HasAclProfilesUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasAclProfilesUri() bool {
	if o != nil && o.AclProfilesUri != nil {
		return true
	}

	return false
}

// SetAclProfilesUri gets a reference to the given string and assigns it to the AclProfilesUri field.
func (o *MsgVpnLinks) SetAclProfilesUri(v string) {
	o.AclProfilesUri = &v
}

// GetAuthenticationOauthProfilesUri returns the AuthenticationOauthProfilesUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetAuthenticationOauthProfilesUri() string {
	if o == nil || o.AuthenticationOauthProfilesUri == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationOauthProfilesUri
}

// GetAuthenticationOauthProfilesUriOk returns a tuple with the AuthenticationOauthProfilesUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetAuthenticationOauthProfilesUriOk() (*string, bool) {
	if o == nil || o.AuthenticationOauthProfilesUri == nil {
		return nil, false
	}
	return o.AuthenticationOauthProfilesUri, true
}

// HasAuthenticationOauthProfilesUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasAuthenticationOauthProfilesUri() bool {
	if o != nil && o.AuthenticationOauthProfilesUri != nil {
		return true
	}

	return false
}

// SetAuthenticationOauthProfilesUri gets a reference to the given string and assigns it to the AuthenticationOauthProfilesUri field.
func (o *MsgVpnLinks) SetAuthenticationOauthProfilesUri(v string) {
	o.AuthenticationOauthProfilesUri = &v
}

// GetAuthenticationOauthProvidersUri returns the AuthenticationOauthProvidersUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetAuthenticationOauthProvidersUri() string {
	if o == nil || o.AuthenticationOauthProvidersUri == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationOauthProvidersUri
}

// GetAuthenticationOauthProvidersUriOk returns a tuple with the AuthenticationOauthProvidersUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetAuthenticationOauthProvidersUriOk() (*string, bool) {
	if o == nil || o.AuthenticationOauthProvidersUri == nil {
		return nil, false
	}
	return o.AuthenticationOauthProvidersUri, true
}

// HasAuthenticationOauthProvidersUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasAuthenticationOauthProvidersUri() bool {
	if o != nil && o.AuthenticationOauthProvidersUri != nil {
		return true
	}

	return false
}

// SetAuthenticationOauthProvidersUri gets a reference to the given string and assigns it to the AuthenticationOauthProvidersUri field.
func (o *MsgVpnLinks) SetAuthenticationOauthProvidersUri(v string) {
	o.AuthenticationOauthProvidersUri = &v
}

// GetAuthorizationGroupsUri returns the AuthorizationGroupsUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetAuthorizationGroupsUri() string {
	if o == nil || o.AuthorizationGroupsUri == nil {
		var ret string
		return ret
	}
	return *o.AuthorizationGroupsUri
}

// GetAuthorizationGroupsUriOk returns a tuple with the AuthorizationGroupsUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetAuthorizationGroupsUriOk() (*string, bool) {
	if o == nil || o.AuthorizationGroupsUri == nil {
		return nil, false
	}
	return o.AuthorizationGroupsUri, true
}

// HasAuthorizationGroupsUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasAuthorizationGroupsUri() bool {
	if o != nil && o.AuthorizationGroupsUri != nil {
		return true
	}

	return false
}

// SetAuthorizationGroupsUri gets a reference to the given string and assigns it to the AuthorizationGroupsUri field.
func (o *MsgVpnLinks) SetAuthorizationGroupsUri(v string) {
	o.AuthorizationGroupsUri = &v
}

// GetBridgesUri returns the BridgesUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetBridgesUri() string {
	if o == nil || o.BridgesUri == nil {
		var ret string
		return ret
	}
	return *o.BridgesUri
}

// GetBridgesUriOk returns a tuple with the BridgesUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetBridgesUriOk() (*string, bool) {
	if o == nil || o.BridgesUri == nil {
		return nil, false
	}
	return o.BridgesUri, true
}

// HasBridgesUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasBridgesUri() bool {
	if o != nil && o.BridgesUri != nil {
		return true
	}

	return false
}

// SetBridgesUri gets a reference to the given string and assigns it to the BridgesUri field.
func (o *MsgVpnLinks) SetBridgesUri(v string) {
	o.BridgesUri = &v
}

// GetClientProfilesUri returns the ClientProfilesUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetClientProfilesUri() string {
	if o == nil || o.ClientProfilesUri == nil {
		var ret string
		return ret
	}
	return *o.ClientProfilesUri
}

// GetClientProfilesUriOk returns a tuple with the ClientProfilesUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetClientProfilesUriOk() (*string, bool) {
	if o == nil || o.ClientProfilesUri == nil {
		return nil, false
	}
	return o.ClientProfilesUri, true
}

// HasClientProfilesUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasClientProfilesUri() bool {
	if o != nil && o.ClientProfilesUri != nil {
		return true
	}

	return false
}

// SetClientProfilesUri gets a reference to the given string and assigns it to the ClientProfilesUri field.
func (o *MsgVpnLinks) SetClientProfilesUri(v string) {
	o.ClientProfilesUri = &v
}

// GetClientUsernamesUri returns the ClientUsernamesUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetClientUsernamesUri() string {
	if o == nil || o.ClientUsernamesUri == nil {
		var ret string
		return ret
	}
	return *o.ClientUsernamesUri
}

// GetClientUsernamesUriOk returns a tuple with the ClientUsernamesUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetClientUsernamesUriOk() (*string, bool) {
	if o == nil || o.ClientUsernamesUri == nil {
		return nil, false
	}
	return o.ClientUsernamesUri, true
}

// HasClientUsernamesUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasClientUsernamesUri() bool {
	if o != nil && o.ClientUsernamesUri != nil {
		return true
	}

	return false
}

// SetClientUsernamesUri gets a reference to the given string and assigns it to the ClientUsernamesUri field.
func (o *MsgVpnLinks) SetClientUsernamesUri(v string) {
	o.ClientUsernamesUri = &v
}

// GetDistributedCachesUri returns the DistributedCachesUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetDistributedCachesUri() string {
	if o == nil || o.DistributedCachesUri == nil {
		var ret string
		return ret
	}
	return *o.DistributedCachesUri
}

// GetDistributedCachesUriOk returns a tuple with the DistributedCachesUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetDistributedCachesUriOk() (*string, bool) {
	if o == nil || o.DistributedCachesUri == nil {
		return nil, false
	}
	return o.DistributedCachesUri, true
}

// HasDistributedCachesUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasDistributedCachesUri() bool {
	if o != nil && o.DistributedCachesUri != nil {
		return true
	}

	return false
}

// SetDistributedCachesUri gets a reference to the given string and assigns it to the DistributedCachesUri field.
func (o *MsgVpnLinks) SetDistributedCachesUri(v string) {
	o.DistributedCachesUri = &v
}

// GetDmrBridgesUri returns the DmrBridgesUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetDmrBridgesUri() string {
	if o == nil || o.DmrBridgesUri == nil {
		var ret string
		return ret
	}
	return *o.DmrBridgesUri
}

// GetDmrBridgesUriOk returns a tuple with the DmrBridgesUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetDmrBridgesUriOk() (*string, bool) {
	if o == nil || o.DmrBridgesUri == nil {
		return nil, false
	}
	return o.DmrBridgesUri, true
}

// HasDmrBridgesUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasDmrBridgesUri() bool {
	if o != nil && o.DmrBridgesUri != nil {
		return true
	}

	return false
}

// SetDmrBridgesUri gets a reference to the given string and assigns it to the DmrBridgesUri field.
func (o *MsgVpnLinks) SetDmrBridgesUri(v string) {
	o.DmrBridgesUri = &v
}

// GetJndiConnectionFactoriesUri returns the JndiConnectionFactoriesUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetJndiConnectionFactoriesUri() string {
	if o == nil || o.JndiConnectionFactoriesUri == nil {
		var ret string
		return ret
	}
	return *o.JndiConnectionFactoriesUri
}

// GetJndiConnectionFactoriesUriOk returns a tuple with the JndiConnectionFactoriesUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetJndiConnectionFactoriesUriOk() (*string, bool) {
	if o == nil || o.JndiConnectionFactoriesUri == nil {
		return nil, false
	}
	return o.JndiConnectionFactoriesUri, true
}

// HasJndiConnectionFactoriesUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasJndiConnectionFactoriesUri() bool {
	if o != nil && o.JndiConnectionFactoriesUri != nil {
		return true
	}

	return false
}

// SetJndiConnectionFactoriesUri gets a reference to the given string and assigns it to the JndiConnectionFactoriesUri field.
func (o *MsgVpnLinks) SetJndiConnectionFactoriesUri(v string) {
	o.JndiConnectionFactoriesUri = &v
}

// GetJndiQueuesUri returns the JndiQueuesUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetJndiQueuesUri() string {
	if o == nil || o.JndiQueuesUri == nil {
		var ret string
		return ret
	}
	return *o.JndiQueuesUri
}

// GetJndiQueuesUriOk returns a tuple with the JndiQueuesUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetJndiQueuesUriOk() (*string, bool) {
	if o == nil || o.JndiQueuesUri == nil {
		return nil, false
	}
	return o.JndiQueuesUri, true
}

// HasJndiQueuesUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasJndiQueuesUri() bool {
	if o != nil && o.JndiQueuesUri != nil {
		return true
	}

	return false
}

// SetJndiQueuesUri gets a reference to the given string and assigns it to the JndiQueuesUri field.
func (o *MsgVpnLinks) SetJndiQueuesUri(v string) {
	o.JndiQueuesUri = &v
}

// GetJndiTopicsUri returns the JndiTopicsUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetJndiTopicsUri() string {
	if o == nil || o.JndiTopicsUri == nil {
		var ret string
		return ret
	}
	return *o.JndiTopicsUri
}

// GetJndiTopicsUriOk returns a tuple with the JndiTopicsUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetJndiTopicsUriOk() (*string, bool) {
	if o == nil || o.JndiTopicsUri == nil {
		return nil, false
	}
	return o.JndiTopicsUri, true
}

// HasJndiTopicsUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasJndiTopicsUri() bool {
	if o != nil && o.JndiTopicsUri != nil {
		return true
	}

	return false
}

// SetJndiTopicsUri gets a reference to the given string and assigns it to the JndiTopicsUri field.
func (o *MsgVpnLinks) SetJndiTopicsUri(v string) {
	o.JndiTopicsUri = &v
}

// GetMqttRetainCachesUri returns the MqttRetainCachesUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetMqttRetainCachesUri() string {
	if o == nil || o.MqttRetainCachesUri == nil {
		var ret string
		return ret
	}
	return *o.MqttRetainCachesUri
}

// GetMqttRetainCachesUriOk returns a tuple with the MqttRetainCachesUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetMqttRetainCachesUriOk() (*string, bool) {
	if o == nil || o.MqttRetainCachesUri == nil {
		return nil, false
	}
	return o.MqttRetainCachesUri, true
}

// HasMqttRetainCachesUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasMqttRetainCachesUri() bool {
	if o != nil && o.MqttRetainCachesUri != nil {
		return true
	}

	return false
}

// SetMqttRetainCachesUri gets a reference to the given string and assigns it to the MqttRetainCachesUri field.
func (o *MsgVpnLinks) SetMqttRetainCachesUri(v string) {
	o.MqttRetainCachesUri = &v
}

// GetMqttSessionsUri returns the MqttSessionsUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetMqttSessionsUri() string {
	if o == nil || o.MqttSessionsUri == nil {
		var ret string
		return ret
	}
	return *o.MqttSessionsUri
}

// GetMqttSessionsUriOk returns a tuple with the MqttSessionsUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetMqttSessionsUriOk() (*string, bool) {
	if o == nil || o.MqttSessionsUri == nil {
		return nil, false
	}
	return o.MqttSessionsUri, true
}

// HasMqttSessionsUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasMqttSessionsUri() bool {
	if o != nil && o.MqttSessionsUri != nil {
		return true
	}

	return false
}

// SetMqttSessionsUri gets a reference to the given string and assigns it to the MqttSessionsUri field.
func (o *MsgVpnLinks) SetMqttSessionsUri(v string) {
	o.MqttSessionsUri = &v
}

// GetQueueTemplatesUri returns the QueueTemplatesUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetQueueTemplatesUri() string {
	if o == nil || o.QueueTemplatesUri == nil {
		var ret string
		return ret
	}
	return *o.QueueTemplatesUri
}

// GetQueueTemplatesUriOk returns a tuple with the QueueTemplatesUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetQueueTemplatesUriOk() (*string, bool) {
	if o == nil || o.QueueTemplatesUri == nil {
		return nil, false
	}
	return o.QueueTemplatesUri, true
}

// HasQueueTemplatesUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasQueueTemplatesUri() bool {
	if o != nil && o.QueueTemplatesUri != nil {
		return true
	}

	return false
}

// SetQueueTemplatesUri gets a reference to the given string and assigns it to the QueueTemplatesUri field.
func (o *MsgVpnLinks) SetQueueTemplatesUri(v string) {
	o.QueueTemplatesUri = &v
}

// GetQueuesUri returns the QueuesUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetQueuesUri() string {
	if o == nil || o.QueuesUri == nil {
		var ret string
		return ret
	}
	return *o.QueuesUri
}

// GetQueuesUriOk returns a tuple with the QueuesUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetQueuesUriOk() (*string, bool) {
	if o == nil || o.QueuesUri == nil {
		return nil, false
	}
	return o.QueuesUri, true
}

// HasQueuesUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasQueuesUri() bool {
	if o != nil && o.QueuesUri != nil {
		return true
	}

	return false
}

// SetQueuesUri gets a reference to the given string and assigns it to the QueuesUri field.
func (o *MsgVpnLinks) SetQueuesUri(v string) {
	o.QueuesUri = &v
}

// GetReplayLogsUri returns the ReplayLogsUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetReplayLogsUri() string {
	if o == nil || o.ReplayLogsUri == nil {
		var ret string
		return ret
	}
	return *o.ReplayLogsUri
}

// GetReplayLogsUriOk returns a tuple with the ReplayLogsUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetReplayLogsUriOk() (*string, bool) {
	if o == nil || o.ReplayLogsUri == nil {
		return nil, false
	}
	return o.ReplayLogsUri, true
}

// HasReplayLogsUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasReplayLogsUri() bool {
	if o != nil && o.ReplayLogsUri != nil {
		return true
	}

	return false
}

// SetReplayLogsUri gets a reference to the given string and assigns it to the ReplayLogsUri field.
func (o *MsgVpnLinks) SetReplayLogsUri(v string) {
	o.ReplayLogsUri = &v
}

// GetReplicatedTopicsUri returns the ReplicatedTopicsUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetReplicatedTopicsUri() string {
	if o == nil || o.ReplicatedTopicsUri == nil {
		var ret string
		return ret
	}
	return *o.ReplicatedTopicsUri
}

// GetReplicatedTopicsUriOk returns a tuple with the ReplicatedTopicsUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetReplicatedTopicsUriOk() (*string, bool) {
	if o == nil || o.ReplicatedTopicsUri == nil {
		return nil, false
	}
	return o.ReplicatedTopicsUri, true
}

// HasReplicatedTopicsUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasReplicatedTopicsUri() bool {
	if o != nil && o.ReplicatedTopicsUri != nil {
		return true
	}

	return false
}

// SetReplicatedTopicsUri gets a reference to the given string and assigns it to the ReplicatedTopicsUri field.
func (o *MsgVpnLinks) SetReplicatedTopicsUri(v string) {
	o.ReplicatedTopicsUri = &v
}

// GetRestDeliveryPointsUri returns the RestDeliveryPointsUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetRestDeliveryPointsUri() string {
	if o == nil || o.RestDeliveryPointsUri == nil {
		var ret string
		return ret
	}
	return *o.RestDeliveryPointsUri
}

// GetRestDeliveryPointsUriOk returns a tuple with the RestDeliveryPointsUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetRestDeliveryPointsUriOk() (*string, bool) {
	if o == nil || o.RestDeliveryPointsUri == nil {
		return nil, false
	}
	return o.RestDeliveryPointsUri, true
}

// HasRestDeliveryPointsUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasRestDeliveryPointsUri() bool {
	if o != nil && o.RestDeliveryPointsUri != nil {
		return true
	}

	return false
}

// SetRestDeliveryPointsUri gets a reference to the given string and assigns it to the RestDeliveryPointsUri field.
func (o *MsgVpnLinks) SetRestDeliveryPointsUri(v string) {
	o.RestDeliveryPointsUri = &v
}

// GetSequencedTopicsUri returns the SequencedTopicsUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetSequencedTopicsUri() string {
	if o == nil || o.SequencedTopicsUri == nil {
		var ret string
		return ret
	}
	return *o.SequencedTopicsUri
}

// GetSequencedTopicsUriOk returns a tuple with the SequencedTopicsUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetSequencedTopicsUriOk() (*string, bool) {
	if o == nil || o.SequencedTopicsUri == nil {
		return nil, false
	}
	return o.SequencedTopicsUri, true
}

// HasSequencedTopicsUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasSequencedTopicsUri() bool {
	if o != nil && o.SequencedTopicsUri != nil {
		return true
	}

	return false
}

// SetSequencedTopicsUri gets a reference to the given string and assigns it to the SequencedTopicsUri field.
func (o *MsgVpnLinks) SetSequencedTopicsUri(v string) {
	o.SequencedTopicsUri = &v
}

// GetTopicEndpointTemplatesUri returns the TopicEndpointTemplatesUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetTopicEndpointTemplatesUri() string {
	if o == nil || o.TopicEndpointTemplatesUri == nil {
		var ret string
		return ret
	}
	return *o.TopicEndpointTemplatesUri
}

// GetTopicEndpointTemplatesUriOk returns a tuple with the TopicEndpointTemplatesUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetTopicEndpointTemplatesUriOk() (*string, bool) {
	if o == nil || o.TopicEndpointTemplatesUri == nil {
		return nil, false
	}
	return o.TopicEndpointTemplatesUri, true
}

// HasTopicEndpointTemplatesUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasTopicEndpointTemplatesUri() bool {
	if o != nil && o.TopicEndpointTemplatesUri != nil {
		return true
	}

	return false
}

// SetTopicEndpointTemplatesUri gets a reference to the given string and assigns it to the TopicEndpointTemplatesUri field.
func (o *MsgVpnLinks) SetTopicEndpointTemplatesUri(v string) {
	o.TopicEndpointTemplatesUri = &v
}

// GetTopicEndpointsUri returns the TopicEndpointsUri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetTopicEndpointsUri() string {
	if o == nil || o.TopicEndpointsUri == nil {
		var ret string
		return ret
	}
	return *o.TopicEndpointsUri
}

// GetTopicEndpointsUriOk returns a tuple with the TopicEndpointsUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetTopicEndpointsUriOk() (*string, bool) {
	if o == nil || o.TopicEndpointsUri == nil {
		return nil, false
	}
	return o.TopicEndpointsUri, true
}

// HasTopicEndpointsUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasTopicEndpointsUri() bool {
	if o != nil && o.TopicEndpointsUri != nil {
		return true
	}

	return false
}

// SetTopicEndpointsUri gets a reference to the given string and assigns it to the TopicEndpointsUri field.
func (o *MsgVpnLinks) SetTopicEndpointsUri(v string) {
	o.TopicEndpointsUri = &v
}

// GetUri returns the Uri field value if set, zero value otherwise.
func (o *MsgVpnLinks) GetUri() string {
	if o == nil || o.Uri == nil {
		var ret string
		return ret
	}
	return *o.Uri
}

// GetUriOk returns a tuple with the Uri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnLinks) GetUriOk() (*string, bool) {
	if o == nil || o.Uri == nil {
		return nil, false
	}
	return o.Uri, true
}

// HasUri returns a boolean if a field has been set.
func (o *MsgVpnLinks) HasUri() bool {
	if o != nil && o.Uri != nil {
		return true
	}

	return false
}

// SetUri gets a reference to the given string and assigns it to the Uri field.
func (o *MsgVpnLinks) SetUri(v string) {
	o.Uri = &v
}

func (o MsgVpnLinks) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AclProfilesUri != nil {
		toSerialize["aclProfilesUri"] = o.AclProfilesUri
	}
	if o.AuthenticationOauthProfilesUri != nil {
		toSerialize["authenticationOauthProfilesUri"] = o.AuthenticationOauthProfilesUri
	}
	if o.AuthenticationOauthProvidersUri != nil {
		toSerialize["authenticationOauthProvidersUri"] = o.AuthenticationOauthProvidersUri
	}
	if o.AuthorizationGroupsUri != nil {
		toSerialize["authorizationGroupsUri"] = o.AuthorizationGroupsUri
	}
	if o.BridgesUri != nil {
		toSerialize["bridgesUri"] = o.BridgesUri
	}
	if o.ClientProfilesUri != nil {
		toSerialize["clientProfilesUri"] = o.ClientProfilesUri
	}
	if o.ClientUsernamesUri != nil {
		toSerialize["clientUsernamesUri"] = o.ClientUsernamesUri
	}
	if o.DistributedCachesUri != nil {
		toSerialize["distributedCachesUri"] = o.DistributedCachesUri
	}
	if o.DmrBridgesUri != nil {
		toSerialize["dmrBridgesUri"] = o.DmrBridgesUri
	}
	if o.JndiConnectionFactoriesUri != nil {
		toSerialize["jndiConnectionFactoriesUri"] = o.JndiConnectionFactoriesUri
	}
	if o.JndiQueuesUri != nil {
		toSerialize["jndiQueuesUri"] = o.JndiQueuesUri
	}
	if o.JndiTopicsUri != nil {
		toSerialize["jndiTopicsUri"] = o.JndiTopicsUri
	}
	if o.MqttRetainCachesUri != nil {
		toSerialize["mqttRetainCachesUri"] = o.MqttRetainCachesUri
	}
	if o.MqttSessionsUri != nil {
		toSerialize["mqttSessionsUri"] = o.MqttSessionsUri
	}
	if o.QueueTemplatesUri != nil {
		toSerialize["queueTemplatesUri"] = o.QueueTemplatesUri
	}
	if o.QueuesUri != nil {
		toSerialize["queuesUri"] = o.QueuesUri
	}
	if o.ReplayLogsUri != nil {
		toSerialize["replayLogsUri"] = o.ReplayLogsUri
	}
	if o.ReplicatedTopicsUri != nil {
		toSerialize["replicatedTopicsUri"] = o.ReplicatedTopicsUri
	}
	if o.RestDeliveryPointsUri != nil {
		toSerialize["restDeliveryPointsUri"] = o.RestDeliveryPointsUri
	}
	if o.SequencedTopicsUri != nil {
		toSerialize["sequencedTopicsUri"] = o.SequencedTopicsUri
	}
	if o.TopicEndpointTemplatesUri != nil {
		toSerialize["topicEndpointTemplatesUri"] = o.TopicEndpointTemplatesUri
	}
	if o.TopicEndpointsUri != nil {
		toSerialize["topicEndpointsUri"] = o.TopicEndpointsUri
	}
	if o.Uri != nil {
		toSerialize["uri"] = o.Uri
	}
	return json.Marshal(toSerialize)
}

type NullableMsgVpnLinks struct {
	value *MsgVpnLinks
	isSet bool
}

func (v NullableMsgVpnLinks) Get() *MsgVpnLinks {
	return v.value
}

func (v *NullableMsgVpnLinks) Set(val *MsgVpnLinks) {
	v.value = val
	v.isSet = true
}

func (v NullableMsgVpnLinks) IsSet() bool {
	return v.isSet
}

func (v *NullableMsgVpnLinks) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMsgVpnLinks(val *MsgVpnLinks) *NullableMsgVpnLinks {
	return &NullableMsgVpnLinks{value: val, isSet: true}
}

func (v NullableMsgVpnLinks) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMsgVpnLinks) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


