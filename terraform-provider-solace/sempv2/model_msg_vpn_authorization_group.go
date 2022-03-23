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

// MsgVpnAuthorizationGroup struct for MsgVpnAuthorizationGroup
type MsgVpnAuthorizationGroup struct {
	// The ACL Profile of the Authorization Group. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"default\"`.
	AclProfileName *string `json:"aclProfileName,omitempty"`
	// The name of the Authorization Group. For LDAP groups, special care is needed if the group name contains special characters such as '#', '+', ';', '=' as the value of the group name returned from the LDAP server might prepend those characters with '\\'. For example a group name called 'test#,lab,com' will be returned from the LDAP server as 'test\\#,lab,com'.
	AuthorizationGroupName *string `json:"authorizationGroupName,omitempty"`
	// The Client Profile of the Authorization Group. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"default\"`.
	ClientProfileName *string `json:"clientProfileName,omitempty"`
	// Enable or disable the Authorization Group in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	Enabled *bool `json:"enabled,omitempty"`
	// The name of the Message VPN.
	MsgVpnName *string `json:"msgVpnName,omitempty"`
	// Lower the priority to be less than this group. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is not applicable.
	OrderAfterAuthorizationGroupName *string `json:"orderAfterAuthorizationGroupName,omitempty"`
	// Raise the priority to be greater than this group. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is not applicable.
	OrderBeforeAuthorizationGroupName *string `json:"orderBeforeAuthorizationGroupName,omitempty"`
}

// NewMsgVpnAuthorizationGroup instantiates a new MsgVpnAuthorizationGroup object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMsgVpnAuthorizationGroup() *MsgVpnAuthorizationGroup {
	this := MsgVpnAuthorizationGroup{}
	return &this
}

// NewMsgVpnAuthorizationGroupWithDefaults instantiates a new MsgVpnAuthorizationGroup object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMsgVpnAuthorizationGroupWithDefaults() *MsgVpnAuthorizationGroup {
	this := MsgVpnAuthorizationGroup{}
	return &this
}

// GetAclProfileName returns the AclProfileName field value if set, zero value otherwise.
func (o *MsgVpnAuthorizationGroup) GetAclProfileName() string {
	if o == nil || o.AclProfileName == nil {
		var ret string
		return ret
	}
	return *o.AclProfileName
}

// GetAclProfileNameOk returns a tuple with the AclProfileName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthorizationGroup) GetAclProfileNameOk() (*string, bool) {
	if o == nil || o.AclProfileName == nil {
		return nil, false
	}
	return o.AclProfileName, true
}

// HasAclProfileName returns a boolean if a field has been set.
func (o *MsgVpnAuthorizationGroup) HasAclProfileName() bool {
	if o != nil && o.AclProfileName != nil {
		return true
	}

	return false
}

// SetAclProfileName gets a reference to the given string and assigns it to the AclProfileName field.
func (o *MsgVpnAuthorizationGroup) SetAclProfileName(v string) {
	o.AclProfileName = &v
}

// GetAuthorizationGroupName returns the AuthorizationGroupName field value if set, zero value otherwise.
func (o *MsgVpnAuthorizationGroup) GetAuthorizationGroupName() string {
	if o == nil || o.AuthorizationGroupName == nil {
		var ret string
		return ret
	}
	return *o.AuthorizationGroupName
}

// GetAuthorizationGroupNameOk returns a tuple with the AuthorizationGroupName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthorizationGroup) GetAuthorizationGroupNameOk() (*string, bool) {
	if o == nil || o.AuthorizationGroupName == nil {
		return nil, false
	}
	return o.AuthorizationGroupName, true
}

// HasAuthorizationGroupName returns a boolean if a field has been set.
func (o *MsgVpnAuthorizationGroup) HasAuthorizationGroupName() bool {
	if o != nil && o.AuthorizationGroupName != nil {
		return true
	}

	return false
}

// SetAuthorizationGroupName gets a reference to the given string and assigns it to the AuthorizationGroupName field.
func (o *MsgVpnAuthorizationGroup) SetAuthorizationGroupName(v string) {
	o.AuthorizationGroupName = &v
}

// GetClientProfileName returns the ClientProfileName field value if set, zero value otherwise.
func (o *MsgVpnAuthorizationGroup) GetClientProfileName() string {
	if o == nil || o.ClientProfileName == nil {
		var ret string
		return ret
	}
	return *o.ClientProfileName
}

// GetClientProfileNameOk returns a tuple with the ClientProfileName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthorizationGroup) GetClientProfileNameOk() (*string, bool) {
	if o == nil || o.ClientProfileName == nil {
		return nil, false
	}
	return o.ClientProfileName, true
}

// HasClientProfileName returns a boolean if a field has been set.
func (o *MsgVpnAuthorizationGroup) HasClientProfileName() bool {
	if o != nil && o.ClientProfileName != nil {
		return true
	}

	return false
}

// SetClientProfileName gets a reference to the given string and assigns it to the ClientProfileName field.
func (o *MsgVpnAuthorizationGroup) SetClientProfileName(v string) {
	o.ClientProfileName = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *MsgVpnAuthorizationGroup) GetEnabled() bool {
	if o == nil || o.Enabled == nil {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthorizationGroup) GetEnabledOk() (*bool, bool) {
	if o == nil || o.Enabled == nil {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *MsgVpnAuthorizationGroup) HasEnabled() bool {
	if o != nil && o.Enabled != nil {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *MsgVpnAuthorizationGroup) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetMsgVpnName returns the MsgVpnName field value if set, zero value otherwise.
func (o *MsgVpnAuthorizationGroup) GetMsgVpnName() string {
	if o == nil || o.MsgVpnName == nil {
		var ret string
		return ret
	}
	return *o.MsgVpnName
}

// GetMsgVpnNameOk returns a tuple with the MsgVpnName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthorizationGroup) GetMsgVpnNameOk() (*string, bool) {
	if o == nil || o.MsgVpnName == nil {
		return nil, false
	}
	return o.MsgVpnName, true
}

// HasMsgVpnName returns a boolean if a field has been set.
func (o *MsgVpnAuthorizationGroup) HasMsgVpnName() bool {
	if o != nil && o.MsgVpnName != nil {
		return true
	}

	return false
}

// SetMsgVpnName gets a reference to the given string and assigns it to the MsgVpnName field.
func (o *MsgVpnAuthorizationGroup) SetMsgVpnName(v string) {
	o.MsgVpnName = &v
}

// GetOrderAfterAuthorizationGroupName returns the OrderAfterAuthorizationGroupName field value if set, zero value otherwise.
func (o *MsgVpnAuthorizationGroup) GetOrderAfterAuthorizationGroupName() string {
	if o == nil || o.OrderAfterAuthorizationGroupName == nil {
		var ret string
		return ret
	}
	return *o.OrderAfterAuthorizationGroupName
}

// GetOrderAfterAuthorizationGroupNameOk returns a tuple with the OrderAfterAuthorizationGroupName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthorizationGroup) GetOrderAfterAuthorizationGroupNameOk() (*string, bool) {
	if o == nil || o.OrderAfterAuthorizationGroupName == nil {
		return nil, false
	}
	return o.OrderAfterAuthorizationGroupName, true
}

// HasOrderAfterAuthorizationGroupName returns a boolean if a field has been set.
func (o *MsgVpnAuthorizationGroup) HasOrderAfterAuthorizationGroupName() bool {
	if o != nil && o.OrderAfterAuthorizationGroupName != nil {
		return true
	}

	return false
}

// SetOrderAfterAuthorizationGroupName gets a reference to the given string and assigns it to the OrderAfterAuthorizationGroupName field.
func (o *MsgVpnAuthorizationGroup) SetOrderAfterAuthorizationGroupName(v string) {
	o.OrderAfterAuthorizationGroupName = &v
}

// GetOrderBeforeAuthorizationGroupName returns the OrderBeforeAuthorizationGroupName field value if set, zero value otherwise.
func (o *MsgVpnAuthorizationGroup) GetOrderBeforeAuthorizationGroupName() string {
	if o == nil || o.OrderBeforeAuthorizationGroupName == nil {
		var ret string
		return ret
	}
	return *o.OrderBeforeAuthorizationGroupName
}

// GetOrderBeforeAuthorizationGroupNameOk returns a tuple with the OrderBeforeAuthorizationGroupName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthorizationGroup) GetOrderBeforeAuthorizationGroupNameOk() (*string, bool) {
	if o == nil || o.OrderBeforeAuthorizationGroupName == nil {
		return nil, false
	}
	return o.OrderBeforeAuthorizationGroupName, true
}

// HasOrderBeforeAuthorizationGroupName returns a boolean if a field has been set.
func (o *MsgVpnAuthorizationGroup) HasOrderBeforeAuthorizationGroupName() bool {
	if o != nil && o.OrderBeforeAuthorizationGroupName != nil {
		return true
	}

	return false
}

// SetOrderBeforeAuthorizationGroupName gets a reference to the given string and assigns it to the OrderBeforeAuthorizationGroupName field.
func (o *MsgVpnAuthorizationGroup) SetOrderBeforeAuthorizationGroupName(v string) {
	o.OrderBeforeAuthorizationGroupName = &v
}

func (o MsgVpnAuthorizationGroup) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AclProfileName != nil {
		toSerialize["aclProfileName"] = o.AclProfileName
	}
	if o.AuthorizationGroupName != nil {
		toSerialize["authorizationGroupName"] = o.AuthorizationGroupName
	}
	if o.ClientProfileName != nil {
		toSerialize["clientProfileName"] = o.ClientProfileName
	}
	if o.Enabled != nil {
		toSerialize["enabled"] = o.Enabled
	}
	if o.MsgVpnName != nil {
		toSerialize["msgVpnName"] = o.MsgVpnName
	}
	if o.OrderAfterAuthorizationGroupName != nil {
		toSerialize["orderAfterAuthorizationGroupName"] = o.OrderAfterAuthorizationGroupName
	}
	if o.OrderBeforeAuthorizationGroupName != nil {
		toSerialize["orderBeforeAuthorizationGroupName"] = o.OrderBeforeAuthorizationGroupName
	}
	return json.Marshal(toSerialize)
}

type NullableMsgVpnAuthorizationGroup struct {
	value *MsgVpnAuthorizationGroup
	isSet bool
}

func (v NullableMsgVpnAuthorizationGroup) Get() *MsgVpnAuthorizationGroup {
	return v.value
}

func (v *NullableMsgVpnAuthorizationGroup) Set(val *MsgVpnAuthorizationGroup) {
	v.value = val
	v.isSet = true
}

func (v NullableMsgVpnAuthorizationGroup) IsSet() bool {
	return v.isSet
}

func (v *NullableMsgVpnAuthorizationGroup) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMsgVpnAuthorizationGroup(val *MsgVpnAuthorizationGroup) *NullableMsgVpnAuthorizationGroup {
	return &NullableMsgVpnAuthorizationGroup{value: val, isSet: true}
}

func (v NullableMsgVpnAuthorizationGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMsgVpnAuthorizationGroup) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
