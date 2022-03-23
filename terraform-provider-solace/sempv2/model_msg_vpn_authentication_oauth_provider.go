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

// MsgVpnAuthenticationOauthProvider struct for MsgVpnAuthenticationOauthProvider
type MsgVpnAuthenticationOauthProvider struct {
	// The audience claim name, indicating which part of the object to use for determining the audience. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"aud\"`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	AudienceClaimName *string `json:"audienceClaimName,omitempty"`
	// The audience claim source, indicating where to search for the audience value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"id-token\"`. The allowed values and their meaning are:  <pre> \"access-token\" - The OAuth v2 access_token. \"id-token\" - The OpenID Connect id_token. \"introspection\" - The result of introspecting the OAuth v2 access_token. </pre>  Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	AudienceClaimSource *string `json:"audienceClaimSource,omitempty"`
	// The required audience value for a token to be considered valid. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	AudienceClaimValue *string `json:"audienceClaimValue,omitempty"`
	// Enable or disable audience validation. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	AudienceValidationEnabled *bool `json:"audienceValidationEnabled,omitempty"`
	// The authorization group claim name, indicating which part of the object to use for determining the authorization group. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"scope\"`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	AuthorizationGroupClaimName *string `json:"authorizationGroupClaimName,omitempty"`
	// The authorization group claim source, indicating where to search for the authorization group name. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"id-token\"`. The allowed values and their meaning are:  <pre> \"access-token\" - The OAuth v2 access_token. \"id-token\" - The OpenID Connect id_token. \"introspection\" - The result of introspecting the OAuth v2 access_token. </pre>  Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	AuthorizationGroupClaimSource *string `json:"authorizationGroupClaimSource,omitempty"`
	// Enable or disable OAuth based authorization. When enabled, the configured authorization type for OAuth clients is overridden. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	AuthorizationGroupEnabled *bool `json:"authorizationGroupEnabled,omitempty"`
	// Enable or disable the disconnection of clients when their tokens expire. Changing this value does not affect existing clients, only new client connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	DisconnectOnTokenExpirationEnabled *bool `json:"disconnectOnTokenExpirationEnabled,omitempty"`
	// Enable or disable OAuth Provider client authentication. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	Enabled *bool `json:"enabled,omitempty"`
	// The number of seconds between forced JWKS public key refreshing. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `86400`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	JwksRefreshInterval *int32 `json:"jwksRefreshInterval,omitempty"`
	// The URI where the OAuth provider publishes its JWKS public keys. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	JwksUri *string `json:"jwksUri,omitempty"`
	// The name of the Message VPN. Deprecated since 2.25. Replaced by authenticationOauthProfiles.
	MsgVpnName *string `json:"msgVpnName,omitempty"`
	// The name of the OAuth Provider. Deprecated since 2.25. Replaced by authenticationOauthProfiles.
	OauthProviderName *string `json:"oauthProviderName,omitempty"`
	// Enable or disable whether to ignore time limits and accept tokens that are not yet valid or are no longer valid. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	TokenIgnoreTimeLimitsEnabled *bool `json:"tokenIgnoreTimeLimitsEnabled,omitempty"`
	// The parameter name used to identify the token during access token introspection. A standards compliant OAuth introspection server expects \"token\". Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"token\"`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	TokenIntrospectionParameterName *string `json:"tokenIntrospectionParameterName,omitempty"`
	// The password to use when logging into the token introspection URI. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	TokenIntrospectionPassword *string `json:"tokenIntrospectionPassword,omitempty"`
	// The maximum time in seconds a token introspection is allowed to take. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	TokenIntrospectionTimeout *int32 `json:"tokenIntrospectionTimeout,omitempty"`
	// The token introspection URI of the OAuth authentication server. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	TokenIntrospectionUri *string `json:"tokenIntrospectionUri,omitempty"`
	// The username to use when logging into the token introspection URI. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	TokenIntrospectionUsername *string `json:"tokenIntrospectionUsername,omitempty"`
	// The username claim name, indicating which part of the object to use for determining the username. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"sub\"`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	UsernameClaimName *string `json:"usernameClaimName,omitempty"`
	// The username claim source, indicating where to search for the username value. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"id-token\"`. The allowed values and their meaning are:  <pre> \"access-token\" - The OAuth v2 access_token. \"id-token\" - The OpenID Connect id_token. \"introspection\" - The result of introspecting the OAuth v2 access_token. </pre>  Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	UsernameClaimSource *string `json:"usernameClaimSource,omitempty"`
	// Enable or disable whether the API provided username will be validated against the username calculated from the token(s); the connection attempt is rejected if they differ. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Deprecated since 2.25. authenticationOauthProviders replaced by authenticationOauthProfiles.
	UsernameValidateEnabled *bool `json:"usernameValidateEnabled,omitempty"`
}

// NewMsgVpnAuthenticationOauthProvider instantiates a new MsgVpnAuthenticationOauthProvider object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMsgVpnAuthenticationOauthProvider() *MsgVpnAuthenticationOauthProvider {
	this := MsgVpnAuthenticationOauthProvider{}
	return &this
}

// NewMsgVpnAuthenticationOauthProviderWithDefaults instantiates a new MsgVpnAuthenticationOauthProvider object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMsgVpnAuthenticationOauthProviderWithDefaults() *MsgVpnAuthenticationOauthProvider {
	this := MsgVpnAuthenticationOauthProvider{}
	return &this
}

// GetAudienceClaimName returns the AudienceClaimName field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetAudienceClaimName() string {
	if o == nil || o.AudienceClaimName == nil {
		var ret string
		return ret
	}
	return *o.AudienceClaimName
}

// GetAudienceClaimNameOk returns a tuple with the AudienceClaimName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetAudienceClaimNameOk() (*string, bool) {
	if o == nil || o.AudienceClaimName == nil {
		return nil, false
	}
	return o.AudienceClaimName, true
}

// HasAudienceClaimName returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasAudienceClaimName() bool {
	if o != nil && o.AudienceClaimName != nil {
		return true
	}

	return false
}

// SetAudienceClaimName gets a reference to the given string and assigns it to the AudienceClaimName field.
func (o *MsgVpnAuthenticationOauthProvider) SetAudienceClaimName(v string) {
	o.AudienceClaimName = &v
}

// GetAudienceClaimSource returns the AudienceClaimSource field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetAudienceClaimSource() string {
	if o == nil || o.AudienceClaimSource == nil {
		var ret string
		return ret
	}
	return *o.AudienceClaimSource
}

// GetAudienceClaimSourceOk returns a tuple with the AudienceClaimSource field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetAudienceClaimSourceOk() (*string, bool) {
	if o == nil || o.AudienceClaimSource == nil {
		return nil, false
	}
	return o.AudienceClaimSource, true
}

// HasAudienceClaimSource returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasAudienceClaimSource() bool {
	if o != nil && o.AudienceClaimSource != nil {
		return true
	}

	return false
}

// SetAudienceClaimSource gets a reference to the given string and assigns it to the AudienceClaimSource field.
func (o *MsgVpnAuthenticationOauthProvider) SetAudienceClaimSource(v string) {
	o.AudienceClaimSource = &v
}

// GetAudienceClaimValue returns the AudienceClaimValue field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetAudienceClaimValue() string {
	if o == nil || o.AudienceClaimValue == nil {
		var ret string
		return ret
	}
	return *o.AudienceClaimValue
}

// GetAudienceClaimValueOk returns a tuple with the AudienceClaimValue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetAudienceClaimValueOk() (*string, bool) {
	if o == nil || o.AudienceClaimValue == nil {
		return nil, false
	}
	return o.AudienceClaimValue, true
}

// HasAudienceClaimValue returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasAudienceClaimValue() bool {
	if o != nil && o.AudienceClaimValue != nil {
		return true
	}

	return false
}

// SetAudienceClaimValue gets a reference to the given string and assigns it to the AudienceClaimValue field.
func (o *MsgVpnAuthenticationOauthProvider) SetAudienceClaimValue(v string) {
	o.AudienceClaimValue = &v
}

// GetAudienceValidationEnabled returns the AudienceValidationEnabled field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetAudienceValidationEnabled() bool {
	if o == nil || o.AudienceValidationEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AudienceValidationEnabled
}

// GetAudienceValidationEnabledOk returns a tuple with the AudienceValidationEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetAudienceValidationEnabledOk() (*bool, bool) {
	if o == nil || o.AudienceValidationEnabled == nil {
		return nil, false
	}
	return o.AudienceValidationEnabled, true
}

// HasAudienceValidationEnabled returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasAudienceValidationEnabled() bool {
	if o != nil && o.AudienceValidationEnabled != nil {
		return true
	}

	return false
}

// SetAudienceValidationEnabled gets a reference to the given bool and assigns it to the AudienceValidationEnabled field.
func (o *MsgVpnAuthenticationOauthProvider) SetAudienceValidationEnabled(v bool) {
	o.AudienceValidationEnabled = &v
}

// GetAuthorizationGroupClaimName returns the AuthorizationGroupClaimName field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetAuthorizationGroupClaimName() string {
	if o == nil || o.AuthorizationGroupClaimName == nil {
		var ret string
		return ret
	}
	return *o.AuthorizationGroupClaimName
}

// GetAuthorizationGroupClaimNameOk returns a tuple with the AuthorizationGroupClaimName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetAuthorizationGroupClaimNameOk() (*string, bool) {
	if o == nil || o.AuthorizationGroupClaimName == nil {
		return nil, false
	}
	return o.AuthorizationGroupClaimName, true
}

// HasAuthorizationGroupClaimName returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasAuthorizationGroupClaimName() bool {
	if o != nil && o.AuthorizationGroupClaimName != nil {
		return true
	}

	return false
}

// SetAuthorizationGroupClaimName gets a reference to the given string and assigns it to the AuthorizationGroupClaimName field.
func (o *MsgVpnAuthenticationOauthProvider) SetAuthorizationGroupClaimName(v string) {
	o.AuthorizationGroupClaimName = &v
}

// GetAuthorizationGroupClaimSource returns the AuthorizationGroupClaimSource field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetAuthorizationGroupClaimSource() string {
	if o == nil || o.AuthorizationGroupClaimSource == nil {
		var ret string
		return ret
	}
	return *o.AuthorizationGroupClaimSource
}

// GetAuthorizationGroupClaimSourceOk returns a tuple with the AuthorizationGroupClaimSource field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetAuthorizationGroupClaimSourceOk() (*string, bool) {
	if o == nil || o.AuthorizationGroupClaimSource == nil {
		return nil, false
	}
	return o.AuthorizationGroupClaimSource, true
}

// HasAuthorizationGroupClaimSource returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasAuthorizationGroupClaimSource() bool {
	if o != nil && o.AuthorizationGroupClaimSource != nil {
		return true
	}

	return false
}

// SetAuthorizationGroupClaimSource gets a reference to the given string and assigns it to the AuthorizationGroupClaimSource field.
func (o *MsgVpnAuthenticationOauthProvider) SetAuthorizationGroupClaimSource(v string) {
	o.AuthorizationGroupClaimSource = &v
}

// GetAuthorizationGroupEnabled returns the AuthorizationGroupEnabled field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetAuthorizationGroupEnabled() bool {
	if o == nil || o.AuthorizationGroupEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AuthorizationGroupEnabled
}

// GetAuthorizationGroupEnabledOk returns a tuple with the AuthorizationGroupEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetAuthorizationGroupEnabledOk() (*bool, bool) {
	if o == nil || o.AuthorizationGroupEnabled == nil {
		return nil, false
	}
	return o.AuthorizationGroupEnabled, true
}

// HasAuthorizationGroupEnabled returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasAuthorizationGroupEnabled() bool {
	if o != nil && o.AuthorizationGroupEnabled != nil {
		return true
	}

	return false
}

// SetAuthorizationGroupEnabled gets a reference to the given bool and assigns it to the AuthorizationGroupEnabled field.
func (o *MsgVpnAuthenticationOauthProvider) SetAuthorizationGroupEnabled(v bool) {
	o.AuthorizationGroupEnabled = &v
}

// GetDisconnectOnTokenExpirationEnabled returns the DisconnectOnTokenExpirationEnabled field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetDisconnectOnTokenExpirationEnabled() bool {
	if o == nil || o.DisconnectOnTokenExpirationEnabled == nil {
		var ret bool
		return ret
	}
	return *o.DisconnectOnTokenExpirationEnabled
}

// GetDisconnectOnTokenExpirationEnabledOk returns a tuple with the DisconnectOnTokenExpirationEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetDisconnectOnTokenExpirationEnabledOk() (*bool, bool) {
	if o == nil || o.DisconnectOnTokenExpirationEnabled == nil {
		return nil, false
	}
	return o.DisconnectOnTokenExpirationEnabled, true
}

// HasDisconnectOnTokenExpirationEnabled returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasDisconnectOnTokenExpirationEnabled() bool {
	if o != nil && o.DisconnectOnTokenExpirationEnabled != nil {
		return true
	}

	return false
}

// SetDisconnectOnTokenExpirationEnabled gets a reference to the given bool and assigns it to the DisconnectOnTokenExpirationEnabled field.
func (o *MsgVpnAuthenticationOauthProvider) SetDisconnectOnTokenExpirationEnabled(v bool) {
	o.DisconnectOnTokenExpirationEnabled = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetEnabled() bool {
	if o == nil || o.Enabled == nil {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetEnabledOk() (*bool, bool) {
	if o == nil || o.Enabled == nil {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasEnabled() bool {
	if o != nil && o.Enabled != nil {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *MsgVpnAuthenticationOauthProvider) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetJwksRefreshInterval returns the JwksRefreshInterval field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetJwksRefreshInterval() int32 {
	if o == nil || o.JwksRefreshInterval == nil {
		var ret int32
		return ret
	}
	return *o.JwksRefreshInterval
}

// GetJwksRefreshIntervalOk returns a tuple with the JwksRefreshInterval field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetJwksRefreshIntervalOk() (*int32, bool) {
	if o == nil || o.JwksRefreshInterval == nil {
		return nil, false
	}
	return o.JwksRefreshInterval, true
}

// HasJwksRefreshInterval returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasJwksRefreshInterval() bool {
	if o != nil && o.JwksRefreshInterval != nil {
		return true
	}

	return false
}

// SetJwksRefreshInterval gets a reference to the given int32 and assigns it to the JwksRefreshInterval field.
func (o *MsgVpnAuthenticationOauthProvider) SetJwksRefreshInterval(v int32) {
	o.JwksRefreshInterval = &v
}

// GetJwksUri returns the JwksUri field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetJwksUri() string {
	if o == nil || o.JwksUri == nil {
		var ret string
		return ret
	}
	return *o.JwksUri
}

// GetJwksUriOk returns a tuple with the JwksUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetJwksUriOk() (*string, bool) {
	if o == nil || o.JwksUri == nil {
		return nil, false
	}
	return o.JwksUri, true
}

// HasJwksUri returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasJwksUri() bool {
	if o != nil && o.JwksUri != nil {
		return true
	}

	return false
}

// SetJwksUri gets a reference to the given string and assigns it to the JwksUri field.
func (o *MsgVpnAuthenticationOauthProvider) SetJwksUri(v string) {
	o.JwksUri = &v
}

// GetMsgVpnName returns the MsgVpnName field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetMsgVpnName() string {
	if o == nil || o.MsgVpnName == nil {
		var ret string
		return ret
	}
	return *o.MsgVpnName
}

// GetMsgVpnNameOk returns a tuple with the MsgVpnName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetMsgVpnNameOk() (*string, bool) {
	if o == nil || o.MsgVpnName == nil {
		return nil, false
	}
	return o.MsgVpnName, true
}

// HasMsgVpnName returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasMsgVpnName() bool {
	if o != nil && o.MsgVpnName != nil {
		return true
	}

	return false
}

// SetMsgVpnName gets a reference to the given string and assigns it to the MsgVpnName field.
func (o *MsgVpnAuthenticationOauthProvider) SetMsgVpnName(v string) {
	o.MsgVpnName = &v
}

// GetOauthProviderName returns the OauthProviderName field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetOauthProviderName() string {
	if o == nil || o.OauthProviderName == nil {
		var ret string
		return ret
	}
	return *o.OauthProviderName
}

// GetOauthProviderNameOk returns a tuple with the OauthProviderName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetOauthProviderNameOk() (*string, bool) {
	if o == nil || o.OauthProviderName == nil {
		return nil, false
	}
	return o.OauthProviderName, true
}

// HasOauthProviderName returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasOauthProviderName() bool {
	if o != nil && o.OauthProviderName != nil {
		return true
	}

	return false
}

// SetOauthProviderName gets a reference to the given string and assigns it to the OauthProviderName field.
func (o *MsgVpnAuthenticationOauthProvider) SetOauthProviderName(v string) {
	o.OauthProviderName = &v
}

// GetTokenIgnoreTimeLimitsEnabled returns the TokenIgnoreTimeLimitsEnabled field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetTokenIgnoreTimeLimitsEnabled() bool {
	if o == nil || o.TokenIgnoreTimeLimitsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TokenIgnoreTimeLimitsEnabled
}

// GetTokenIgnoreTimeLimitsEnabledOk returns a tuple with the TokenIgnoreTimeLimitsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetTokenIgnoreTimeLimitsEnabledOk() (*bool, bool) {
	if o == nil || o.TokenIgnoreTimeLimitsEnabled == nil {
		return nil, false
	}
	return o.TokenIgnoreTimeLimitsEnabled, true
}

// HasTokenIgnoreTimeLimitsEnabled returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasTokenIgnoreTimeLimitsEnabled() bool {
	if o != nil && o.TokenIgnoreTimeLimitsEnabled != nil {
		return true
	}

	return false
}

// SetTokenIgnoreTimeLimitsEnabled gets a reference to the given bool and assigns it to the TokenIgnoreTimeLimitsEnabled field.
func (o *MsgVpnAuthenticationOauthProvider) SetTokenIgnoreTimeLimitsEnabled(v bool) {
	o.TokenIgnoreTimeLimitsEnabled = &v
}

// GetTokenIntrospectionParameterName returns the TokenIntrospectionParameterName field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetTokenIntrospectionParameterName() string {
	if o == nil || o.TokenIntrospectionParameterName == nil {
		var ret string
		return ret
	}
	return *o.TokenIntrospectionParameterName
}

// GetTokenIntrospectionParameterNameOk returns a tuple with the TokenIntrospectionParameterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetTokenIntrospectionParameterNameOk() (*string, bool) {
	if o == nil || o.TokenIntrospectionParameterName == nil {
		return nil, false
	}
	return o.TokenIntrospectionParameterName, true
}

// HasTokenIntrospectionParameterName returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasTokenIntrospectionParameterName() bool {
	if o != nil && o.TokenIntrospectionParameterName != nil {
		return true
	}

	return false
}

// SetTokenIntrospectionParameterName gets a reference to the given string and assigns it to the TokenIntrospectionParameterName field.
func (o *MsgVpnAuthenticationOauthProvider) SetTokenIntrospectionParameterName(v string) {
	o.TokenIntrospectionParameterName = &v
}

// GetTokenIntrospectionPassword returns the TokenIntrospectionPassword field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetTokenIntrospectionPassword() string {
	if o == nil || o.TokenIntrospectionPassword == nil {
		var ret string
		return ret
	}
	return *o.TokenIntrospectionPassword
}

// GetTokenIntrospectionPasswordOk returns a tuple with the TokenIntrospectionPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetTokenIntrospectionPasswordOk() (*string, bool) {
	if o == nil || o.TokenIntrospectionPassword == nil {
		return nil, false
	}
	return o.TokenIntrospectionPassword, true
}

// HasTokenIntrospectionPassword returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasTokenIntrospectionPassword() bool {
	if o != nil && o.TokenIntrospectionPassword != nil {
		return true
	}

	return false
}

// SetTokenIntrospectionPassword gets a reference to the given string and assigns it to the TokenIntrospectionPassword field.
func (o *MsgVpnAuthenticationOauthProvider) SetTokenIntrospectionPassword(v string) {
	o.TokenIntrospectionPassword = &v
}

// GetTokenIntrospectionTimeout returns the TokenIntrospectionTimeout field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetTokenIntrospectionTimeout() int32 {
	if o == nil || o.TokenIntrospectionTimeout == nil {
		var ret int32
		return ret
	}
	return *o.TokenIntrospectionTimeout
}

// GetTokenIntrospectionTimeoutOk returns a tuple with the TokenIntrospectionTimeout field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetTokenIntrospectionTimeoutOk() (*int32, bool) {
	if o == nil || o.TokenIntrospectionTimeout == nil {
		return nil, false
	}
	return o.TokenIntrospectionTimeout, true
}

// HasTokenIntrospectionTimeout returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasTokenIntrospectionTimeout() bool {
	if o != nil && o.TokenIntrospectionTimeout != nil {
		return true
	}

	return false
}

// SetTokenIntrospectionTimeout gets a reference to the given int32 and assigns it to the TokenIntrospectionTimeout field.
func (o *MsgVpnAuthenticationOauthProvider) SetTokenIntrospectionTimeout(v int32) {
	o.TokenIntrospectionTimeout = &v
}

// GetTokenIntrospectionUri returns the TokenIntrospectionUri field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetTokenIntrospectionUri() string {
	if o == nil || o.TokenIntrospectionUri == nil {
		var ret string
		return ret
	}
	return *o.TokenIntrospectionUri
}

// GetTokenIntrospectionUriOk returns a tuple with the TokenIntrospectionUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetTokenIntrospectionUriOk() (*string, bool) {
	if o == nil || o.TokenIntrospectionUri == nil {
		return nil, false
	}
	return o.TokenIntrospectionUri, true
}

// HasTokenIntrospectionUri returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasTokenIntrospectionUri() bool {
	if o != nil && o.TokenIntrospectionUri != nil {
		return true
	}

	return false
}

// SetTokenIntrospectionUri gets a reference to the given string and assigns it to the TokenIntrospectionUri field.
func (o *MsgVpnAuthenticationOauthProvider) SetTokenIntrospectionUri(v string) {
	o.TokenIntrospectionUri = &v
}

// GetTokenIntrospectionUsername returns the TokenIntrospectionUsername field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetTokenIntrospectionUsername() string {
	if o == nil || o.TokenIntrospectionUsername == nil {
		var ret string
		return ret
	}
	return *o.TokenIntrospectionUsername
}

// GetTokenIntrospectionUsernameOk returns a tuple with the TokenIntrospectionUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetTokenIntrospectionUsernameOk() (*string, bool) {
	if o == nil || o.TokenIntrospectionUsername == nil {
		return nil, false
	}
	return o.TokenIntrospectionUsername, true
}

// HasTokenIntrospectionUsername returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasTokenIntrospectionUsername() bool {
	if o != nil && o.TokenIntrospectionUsername != nil {
		return true
	}

	return false
}

// SetTokenIntrospectionUsername gets a reference to the given string and assigns it to the TokenIntrospectionUsername field.
func (o *MsgVpnAuthenticationOauthProvider) SetTokenIntrospectionUsername(v string) {
	o.TokenIntrospectionUsername = &v
}

// GetUsernameClaimName returns the UsernameClaimName field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetUsernameClaimName() string {
	if o == nil || o.UsernameClaimName == nil {
		var ret string
		return ret
	}
	return *o.UsernameClaimName
}

// GetUsernameClaimNameOk returns a tuple with the UsernameClaimName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetUsernameClaimNameOk() (*string, bool) {
	if o == nil || o.UsernameClaimName == nil {
		return nil, false
	}
	return o.UsernameClaimName, true
}

// HasUsernameClaimName returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasUsernameClaimName() bool {
	if o != nil && o.UsernameClaimName != nil {
		return true
	}

	return false
}

// SetUsernameClaimName gets a reference to the given string and assigns it to the UsernameClaimName field.
func (o *MsgVpnAuthenticationOauthProvider) SetUsernameClaimName(v string) {
	o.UsernameClaimName = &v
}

// GetUsernameClaimSource returns the UsernameClaimSource field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetUsernameClaimSource() string {
	if o == nil || o.UsernameClaimSource == nil {
		var ret string
		return ret
	}
	return *o.UsernameClaimSource
}

// GetUsernameClaimSourceOk returns a tuple with the UsernameClaimSource field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetUsernameClaimSourceOk() (*string, bool) {
	if o == nil || o.UsernameClaimSource == nil {
		return nil, false
	}
	return o.UsernameClaimSource, true
}

// HasUsernameClaimSource returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasUsernameClaimSource() bool {
	if o != nil && o.UsernameClaimSource != nil {
		return true
	}

	return false
}

// SetUsernameClaimSource gets a reference to the given string and assigns it to the UsernameClaimSource field.
func (o *MsgVpnAuthenticationOauthProvider) SetUsernameClaimSource(v string) {
	o.UsernameClaimSource = &v
}

// GetUsernameValidateEnabled returns the UsernameValidateEnabled field value if set, zero value otherwise.
func (o *MsgVpnAuthenticationOauthProvider) GetUsernameValidateEnabled() bool {
	if o == nil || o.UsernameValidateEnabled == nil {
		var ret bool
		return ret
	}
	return *o.UsernameValidateEnabled
}

// GetUsernameValidateEnabledOk returns a tuple with the UsernameValidateEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnAuthenticationOauthProvider) GetUsernameValidateEnabledOk() (*bool, bool) {
	if o == nil || o.UsernameValidateEnabled == nil {
		return nil, false
	}
	return o.UsernameValidateEnabled, true
}

// HasUsernameValidateEnabled returns a boolean if a field has been set.
func (o *MsgVpnAuthenticationOauthProvider) HasUsernameValidateEnabled() bool {
	if o != nil && o.UsernameValidateEnabled != nil {
		return true
	}

	return false
}

// SetUsernameValidateEnabled gets a reference to the given bool and assigns it to the UsernameValidateEnabled field.
func (o *MsgVpnAuthenticationOauthProvider) SetUsernameValidateEnabled(v bool) {
	o.UsernameValidateEnabled = &v
}

func (o MsgVpnAuthenticationOauthProvider) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AudienceClaimName != nil {
		toSerialize["audienceClaimName"] = o.AudienceClaimName
	}
	if o.AudienceClaimSource != nil {
		toSerialize["audienceClaimSource"] = o.AudienceClaimSource
	}
	if o.AudienceClaimValue != nil {
		toSerialize["audienceClaimValue"] = o.AudienceClaimValue
	}
	if o.AudienceValidationEnabled != nil {
		toSerialize["audienceValidationEnabled"] = o.AudienceValidationEnabled
	}
	if o.AuthorizationGroupClaimName != nil {
		toSerialize["authorizationGroupClaimName"] = o.AuthorizationGroupClaimName
	}
	if o.AuthorizationGroupClaimSource != nil {
		toSerialize["authorizationGroupClaimSource"] = o.AuthorizationGroupClaimSource
	}
	if o.AuthorizationGroupEnabled != nil {
		toSerialize["authorizationGroupEnabled"] = o.AuthorizationGroupEnabled
	}
	if o.DisconnectOnTokenExpirationEnabled != nil {
		toSerialize["disconnectOnTokenExpirationEnabled"] = o.DisconnectOnTokenExpirationEnabled
	}
	if o.Enabled != nil {
		toSerialize["enabled"] = o.Enabled
	}
	if o.JwksRefreshInterval != nil {
		toSerialize["jwksRefreshInterval"] = o.JwksRefreshInterval
	}
	if o.JwksUri != nil {
		toSerialize["jwksUri"] = o.JwksUri
	}
	if o.MsgVpnName != nil {
		toSerialize["msgVpnName"] = o.MsgVpnName
	}
	if o.OauthProviderName != nil {
		toSerialize["oauthProviderName"] = o.OauthProviderName
	}
	if o.TokenIgnoreTimeLimitsEnabled != nil {
		toSerialize["tokenIgnoreTimeLimitsEnabled"] = o.TokenIgnoreTimeLimitsEnabled
	}
	if o.TokenIntrospectionParameterName != nil {
		toSerialize["tokenIntrospectionParameterName"] = o.TokenIntrospectionParameterName
	}
	if o.TokenIntrospectionPassword != nil {
		toSerialize["tokenIntrospectionPassword"] = o.TokenIntrospectionPassword
	}
	if o.TokenIntrospectionTimeout != nil {
		toSerialize["tokenIntrospectionTimeout"] = o.TokenIntrospectionTimeout
	}
	if o.TokenIntrospectionUri != nil {
		toSerialize["tokenIntrospectionUri"] = o.TokenIntrospectionUri
	}
	if o.TokenIntrospectionUsername != nil {
		toSerialize["tokenIntrospectionUsername"] = o.TokenIntrospectionUsername
	}
	if o.UsernameClaimName != nil {
		toSerialize["usernameClaimName"] = o.UsernameClaimName
	}
	if o.UsernameClaimSource != nil {
		toSerialize["usernameClaimSource"] = o.UsernameClaimSource
	}
	if o.UsernameValidateEnabled != nil {
		toSerialize["usernameValidateEnabled"] = o.UsernameValidateEnabled
	}
	return json.Marshal(toSerialize)
}

type NullableMsgVpnAuthenticationOauthProvider struct {
	value *MsgVpnAuthenticationOauthProvider
	isSet bool
}

func (v NullableMsgVpnAuthenticationOauthProvider) Get() *MsgVpnAuthenticationOauthProvider {
	return v.value
}

func (v *NullableMsgVpnAuthenticationOauthProvider) Set(val *MsgVpnAuthenticationOauthProvider) {
	v.value = val
	v.isSet = true
}

func (v NullableMsgVpnAuthenticationOauthProvider) IsSet() bool {
	return v.isSet
}

func (v *NullableMsgVpnAuthenticationOauthProvider) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMsgVpnAuthenticationOauthProvider(val *MsgVpnAuthenticationOauthProvider) *NullableMsgVpnAuthenticationOauthProvider {
	return &NullableMsgVpnAuthenticationOauthProvider{value: val, isSet: true}
}

func (v NullableMsgVpnAuthenticationOauthProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMsgVpnAuthenticationOauthProvider) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
