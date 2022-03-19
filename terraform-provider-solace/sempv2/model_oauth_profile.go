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

// OauthProfile struct for OauthProfile
type OauthProfile struct {
	// The name of the groups claim. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"groups\"`.
	AccessLevelGroupsClaimName *string `json:"accessLevelGroupsClaimName,omitempty"`
	// The OAuth client id. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	ClientId *string `json:"clientId,omitempty"`
	// The OAuth redirect URI. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	ClientRedirectUri *string `json:"clientRedirectUri,omitempty"`
	// The required value for the TYP field in the ID token header. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"JWT\"`.
	ClientRequiredType *string `json:"clientRequiredType,omitempty"`
	// The OAuth scope. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"openid email\"`.
	ClientScope *string `json:"clientScope,omitempty"`
	// The OAuth client secret. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	ClientSecret *string `json:"clientSecret,omitempty"`
	// Enable or disable verification of the TYP field in the ID token header. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`.
	ClientValidateTypeEnabled *bool `json:"clientValidateTypeEnabled,omitempty"`
	// The default global access level for this OAuth profile. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"none\"`. The allowed values and their meaning are:  <pre> \"none\" - User has no access to global data. \"read-only\" - User has read-only access to global data. \"read-write\" - User has read-write access to most global data. \"admin\" - User has read-write access to all global data. </pre> 
	DefaultGlobalAccessLevel *string `json:"defaultGlobalAccessLevel,omitempty"`
	// The default message VPN access level for the OAuth profile. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"none\"`. The allowed values and their meaning are:  <pre> \"none\" - User has no access to a Message VPN. \"read-only\" - User has read-only access to a Message VPN. \"read-write\" - User has read-write access to most Message VPN settings. </pre> 
	DefaultMsgVpnAccessLevel *string `json:"defaultMsgVpnAccessLevel,omitempty"`
	// The user friendly name for the OAuth profile. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	DisplayName *string `json:"displayName,omitempty"`
	// Enable or disable the OAuth profile. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.
	Enabled *bool `json:"enabled,omitempty"`
	// The OAuth authorization endpoint. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	EndpointAuthorization *string `json:"endpointAuthorization,omitempty"`
	// The OpenID Connect discovery endpoint or OAuth Authorization Server Metadata endpoint. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	EndpointDiscovery *string `json:"endpointDiscovery,omitempty"`
	// The number of seconds between discovery endpoint requests. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `86400`.
	EndpointDiscoveryRefreshInterval *int32 `json:"endpointDiscoveryRefreshInterval,omitempty"`
	// The OAuth introspection endpoint. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	EndpointIntrospection *string `json:"endpointIntrospection,omitempty"`
	// The maximum time in seconds a token introspection request is allowed to take. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `1`.
	EndpointIntrospectionTimeout *int32 `json:"endpointIntrospectionTimeout,omitempty"`
	// The OAuth JWKS endpoint. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	EndpointJwks *string `json:"endpointJwks,omitempty"`
	// The number of seconds between JWKS endpoint requests. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `86400`.
	EndpointJwksRefreshInterval *int32 `json:"endpointJwksRefreshInterval,omitempty"`
	// The OAuth token endpoint. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	EndpointToken *string `json:"endpointToken,omitempty"`
	// The maximum time in seconds a token request is allowed to take. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `1`.
	EndpointTokenTimeout *int32 `json:"endpointTokenTimeout,omitempty"`
	// The OpenID Connect Userinfo endpoint. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	EndpointUserinfo *string `json:"endpointUserinfo,omitempty"`
	// The maximum time in seconds a userinfo request is allowed to take. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `1`.
	EndpointUserinfoTimeout *int32 `json:"endpointUserinfoTimeout,omitempty"`
	// Enable or disable interactive logins via this OAuth provider. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`.
	InteractiveEnabled *bool `json:"interactiveEnabled,omitempty"`
	// The value of the prompt parameter provided to the OAuth authorization server for login requests where the session has expired. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	InteractivePromptForExpiredSession *string `json:"interactivePromptForExpiredSession,omitempty"`
	// The value of the prompt parameter provided to the OAuth authorization server for login requests where the session is new or the user has explicitly logged out. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"select_account\"`.
	InteractivePromptForNewSession *string `json:"interactivePromptForNewSession,omitempty"`
	// The Issuer Identifier for the OAuth provider. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	Issuer *string `json:"issuer,omitempty"`
	// The name of the OAuth profile.
	OauthProfileName *string `json:"oauthProfileName,omitempty"`
	// The OAuth role of the broker. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"client\"`. The allowed values and their meaning are:  <pre> \"client\" - The broker is in the OAuth client role. \"resource-server\" - The broker is in the OAuth resource server role. </pre> 
	OauthRole *string `json:"oauthRole,omitempty"`
	// Enable or disable parsing of the access token as a JWT. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`.
	ResourceServerParseAccessTokenEnabled *bool `json:"resourceServerParseAccessTokenEnabled,omitempty"`
	// The required audience value. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	ResourceServerRequiredAudience *string `json:"resourceServerRequiredAudience,omitempty"`
	// The required issuer value. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	ResourceServerRequiredIssuer *string `json:"resourceServerRequiredIssuer,omitempty"`
	// A space-separated list of scopes that must be present in the scope claim. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	ResourceServerRequiredScope *string `json:"resourceServerRequiredScope,omitempty"`
	// The required TYP value. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"at+jwt\"`.
	ResourceServerRequiredType *string `json:"resourceServerRequiredType,omitempty"`
	// Enable or disable verification of the audience claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`.
	ResourceServerValidateAudienceEnabled *bool `json:"resourceServerValidateAudienceEnabled,omitempty"`
	// Enable or disable verification of the issuer claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`.
	ResourceServerValidateIssuerEnabled *bool `json:"resourceServerValidateIssuerEnabled,omitempty"`
	// Enable or disable verification of the scope claim in the access token or introspection response. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`.
	ResourceServerValidateScopeEnabled *bool `json:"resourceServerValidateScopeEnabled,omitempty"`
	// Enable or disable verification of the TYP field in the access token header. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`.
	ResourceServerValidateTypeEnabled *bool `json:"resourceServerValidateTypeEnabled,omitempty"`
	// Enable or disable authentication of SEMP requests with OAuth tokens. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`.
	SempEnabled *bool `json:"sempEnabled,omitempty"`
	// The name of the username claim. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"sub\"`.
	UsernameClaimName *string `json:"usernameClaimName,omitempty"`
}

// NewOauthProfile instantiates a new OauthProfile object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOauthProfile() *OauthProfile {
	this := OauthProfile{}
	return &this
}

// NewOauthProfileWithDefaults instantiates a new OauthProfile object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOauthProfileWithDefaults() *OauthProfile {
	this := OauthProfile{}
	return &this
}

// GetAccessLevelGroupsClaimName returns the AccessLevelGroupsClaimName field value if set, zero value otherwise.
func (o *OauthProfile) GetAccessLevelGroupsClaimName() string {
	if o == nil || o.AccessLevelGroupsClaimName == nil {
		var ret string
		return ret
	}
	return *o.AccessLevelGroupsClaimName
}

// GetAccessLevelGroupsClaimNameOk returns a tuple with the AccessLevelGroupsClaimName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetAccessLevelGroupsClaimNameOk() (*string, bool) {
	if o == nil || o.AccessLevelGroupsClaimName == nil {
		return nil, false
	}
	return o.AccessLevelGroupsClaimName, true
}

// HasAccessLevelGroupsClaimName returns a boolean if a field has been set.
func (o *OauthProfile) HasAccessLevelGroupsClaimName() bool {
	if o != nil && o.AccessLevelGroupsClaimName != nil {
		return true
	}

	return false
}

// SetAccessLevelGroupsClaimName gets a reference to the given string and assigns it to the AccessLevelGroupsClaimName field.
func (o *OauthProfile) SetAccessLevelGroupsClaimName(v string) {
	o.AccessLevelGroupsClaimName = &v
}

// GetClientId returns the ClientId field value if set, zero value otherwise.
func (o *OauthProfile) GetClientId() string {
	if o == nil || o.ClientId == nil {
		var ret string
		return ret
	}
	return *o.ClientId
}

// GetClientIdOk returns a tuple with the ClientId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetClientIdOk() (*string, bool) {
	if o == nil || o.ClientId == nil {
		return nil, false
	}
	return o.ClientId, true
}

// HasClientId returns a boolean if a field has been set.
func (o *OauthProfile) HasClientId() bool {
	if o != nil && o.ClientId != nil {
		return true
	}

	return false
}

// SetClientId gets a reference to the given string and assigns it to the ClientId field.
func (o *OauthProfile) SetClientId(v string) {
	o.ClientId = &v
}

// GetClientRedirectUri returns the ClientRedirectUri field value if set, zero value otherwise.
func (o *OauthProfile) GetClientRedirectUri() string {
	if o == nil || o.ClientRedirectUri == nil {
		var ret string
		return ret
	}
	return *o.ClientRedirectUri
}

// GetClientRedirectUriOk returns a tuple with the ClientRedirectUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetClientRedirectUriOk() (*string, bool) {
	if o == nil || o.ClientRedirectUri == nil {
		return nil, false
	}
	return o.ClientRedirectUri, true
}

// HasClientRedirectUri returns a boolean if a field has been set.
func (o *OauthProfile) HasClientRedirectUri() bool {
	if o != nil && o.ClientRedirectUri != nil {
		return true
	}

	return false
}

// SetClientRedirectUri gets a reference to the given string and assigns it to the ClientRedirectUri field.
func (o *OauthProfile) SetClientRedirectUri(v string) {
	o.ClientRedirectUri = &v
}

// GetClientRequiredType returns the ClientRequiredType field value if set, zero value otherwise.
func (o *OauthProfile) GetClientRequiredType() string {
	if o == nil || o.ClientRequiredType == nil {
		var ret string
		return ret
	}
	return *o.ClientRequiredType
}

// GetClientRequiredTypeOk returns a tuple with the ClientRequiredType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetClientRequiredTypeOk() (*string, bool) {
	if o == nil || o.ClientRequiredType == nil {
		return nil, false
	}
	return o.ClientRequiredType, true
}

// HasClientRequiredType returns a boolean if a field has been set.
func (o *OauthProfile) HasClientRequiredType() bool {
	if o != nil && o.ClientRequiredType != nil {
		return true
	}

	return false
}

// SetClientRequiredType gets a reference to the given string and assigns it to the ClientRequiredType field.
func (o *OauthProfile) SetClientRequiredType(v string) {
	o.ClientRequiredType = &v
}

// GetClientScope returns the ClientScope field value if set, zero value otherwise.
func (o *OauthProfile) GetClientScope() string {
	if o == nil || o.ClientScope == nil {
		var ret string
		return ret
	}
	return *o.ClientScope
}

// GetClientScopeOk returns a tuple with the ClientScope field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetClientScopeOk() (*string, bool) {
	if o == nil || o.ClientScope == nil {
		return nil, false
	}
	return o.ClientScope, true
}

// HasClientScope returns a boolean if a field has been set.
func (o *OauthProfile) HasClientScope() bool {
	if o != nil && o.ClientScope != nil {
		return true
	}

	return false
}

// SetClientScope gets a reference to the given string and assigns it to the ClientScope field.
func (o *OauthProfile) SetClientScope(v string) {
	o.ClientScope = &v
}

// GetClientSecret returns the ClientSecret field value if set, zero value otherwise.
func (o *OauthProfile) GetClientSecret() string {
	if o == nil || o.ClientSecret == nil {
		var ret string
		return ret
	}
	return *o.ClientSecret
}

// GetClientSecretOk returns a tuple with the ClientSecret field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetClientSecretOk() (*string, bool) {
	if o == nil || o.ClientSecret == nil {
		return nil, false
	}
	return o.ClientSecret, true
}

// HasClientSecret returns a boolean if a field has been set.
func (o *OauthProfile) HasClientSecret() bool {
	if o != nil && o.ClientSecret != nil {
		return true
	}

	return false
}

// SetClientSecret gets a reference to the given string and assigns it to the ClientSecret field.
func (o *OauthProfile) SetClientSecret(v string) {
	o.ClientSecret = &v
}

// GetClientValidateTypeEnabled returns the ClientValidateTypeEnabled field value if set, zero value otherwise.
func (o *OauthProfile) GetClientValidateTypeEnabled() bool {
	if o == nil || o.ClientValidateTypeEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ClientValidateTypeEnabled
}

// GetClientValidateTypeEnabledOk returns a tuple with the ClientValidateTypeEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetClientValidateTypeEnabledOk() (*bool, bool) {
	if o == nil || o.ClientValidateTypeEnabled == nil {
		return nil, false
	}
	return o.ClientValidateTypeEnabled, true
}

// HasClientValidateTypeEnabled returns a boolean if a field has been set.
func (o *OauthProfile) HasClientValidateTypeEnabled() bool {
	if o != nil && o.ClientValidateTypeEnabled != nil {
		return true
	}

	return false
}

// SetClientValidateTypeEnabled gets a reference to the given bool and assigns it to the ClientValidateTypeEnabled field.
func (o *OauthProfile) SetClientValidateTypeEnabled(v bool) {
	o.ClientValidateTypeEnabled = &v
}

// GetDefaultGlobalAccessLevel returns the DefaultGlobalAccessLevel field value if set, zero value otherwise.
func (o *OauthProfile) GetDefaultGlobalAccessLevel() string {
	if o == nil || o.DefaultGlobalAccessLevel == nil {
		var ret string
		return ret
	}
	return *o.DefaultGlobalAccessLevel
}

// GetDefaultGlobalAccessLevelOk returns a tuple with the DefaultGlobalAccessLevel field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetDefaultGlobalAccessLevelOk() (*string, bool) {
	if o == nil || o.DefaultGlobalAccessLevel == nil {
		return nil, false
	}
	return o.DefaultGlobalAccessLevel, true
}

// HasDefaultGlobalAccessLevel returns a boolean if a field has been set.
func (o *OauthProfile) HasDefaultGlobalAccessLevel() bool {
	if o != nil && o.DefaultGlobalAccessLevel != nil {
		return true
	}

	return false
}

// SetDefaultGlobalAccessLevel gets a reference to the given string and assigns it to the DefaultGlobalAccessLevel field.
func (o *OauthProfile) SetDefaultGlobalAccessLevel(v string) {
	o.DefaultGlobalAccessLevel = &v
}

// GetDefaultMsgVpnAccessLevel returns the DefaultMsgVpnAccessLevel field value if set, zero value otherwise.
func (o *OauthProfile) GetDefaultMsgVpnAccessLevel() string {
	if o == nil || o.DefaultMsgVpnAccessLevel == nil {
		var ret string
		return ret
	}
	return *o.DefaultMsgVpnAccessLevel
}

// GetDefaultMsgVpnAccessLevelOk returns a tuple with the DefaultMsgVpnAccessLevel field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetDefaultMsgVpnAccessLevelOk() (*string, bool) {
	if o == nil || o.DefaultMsgVpnAccessLevel == nil {
		return nil, false
	}
	return o.DefaultMsgVpnAccessLevel, true
}

// HasDefaultMsgVpnAccessLevel returns a boolean if a field has been set.
func (o *OauthProfile) HasDefaultMsgVpnAccessLevel() bool {
	if o != nil && o.DefaultMsgVpnAccessLevel != nil {
		return true
	}

	return false
}

// SetDefaultMsgVpnAccessLevel gets a reference to the given string and assigns it to the DefaultMsgVpnAccessLevel field.
func (o *OauthProfile) SetDefaultMsgVpnAccessLevel(v string) {
	o.DefaultMsgVpnAccessLevel = &v
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise.
func (o *OauthProfile) GetDisplayName() string {
	if o == nil || o.DisplayName == nil {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetDisplayNameOk() (*string, bool) {
	if o == nil || o.DisplayName == nil {
		return nil, false
	}
	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *OauthProfile) HasDisplayName() bool {
	if o != nil && o.DisplayName != nil {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *OauthProfile) SetDisplayName(v string) {
	o.DisplayName = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *OauthProfile) GetEnabled() bool {
	if o == nil || o.Enabled == nil {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetEnabledOk() (*bool, bool) {
	if o == nil || o.Enabled == nil {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *OauthProfile) HasEnabled() bool {
	if o != nil && o.Enabled != nil {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *OauthProfile) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetEndpointAuthorization returns the EndpointAuthorization field value if set, zero value otherwise.
func (o *OauthProfile) GetEndpointAuthorization() string {
	if o == nil || o.EndpointAuthorization == nil {
		var ret string
		return ret
	}
	return *o.EndpointAuthorization
}

// GetEndpointAuthorizationOk returns a tuple with the EndpointAuthorization field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetEndpointAuthorizationOk() (*string, bool) {
	if o == nil || o.EndpointAuthorization == nil {
		return nil, false
	}
	return o.EndpointAuthorization, true
}

// HasEndpointAuthorization returns a boolean if a field has been set.
func (o *OauthProfile) HasEndpointAuthorization() bool {
	if o != nil && o.EndpointAuthorization != nil {
		return true
	}

	return false
}

// SetEndpointAuthorization gets a reference to the given string and assigns it to the EndpointAuthorization field.
func (o *OauthProfile) SetEndpointAuthorization(v string) {
	o.EndpointAuthorization = &v
}

// GetEndpointDiscovery returns the EndpointDiscovery field value if set, zero value otherwise.
func (o *OauthProfile) GetEndpointDiscovery() string {
	if o == nil || o.EndpointDiscovery == nil {
		var ret string
		return ret
	}
	return *o.EndpointDiscovery
}

// GetEndpointDiscoveryOk returns a tuple with the EndpointDiscovery field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetEndpointDiscoveryOk() (*string, bool) {
	if o == nil || o.EndpointDiscovery == nil {
		return nil, false
	}
	return o.EndpointDiscovery, true
}

// HasEndpointDiscovery returns a boolean if a field has been set.
func (o *OauthProfile) HasEndpointDiscovery() bool {
	if o != nil && o.EndpointDiscovery != nil {
		return true
	}

	return false
}

// SetEndpointDiscovery gets a reference to the given string and assigns it to the EndpointDiscovery field.
func (o *OauthProfile) SetEndpointDiscovery(v string) {
	o.EndpointDiscovery = &v
}

// GetEndpointDiscoveryRefreshInterval returns the EndpointDiscoveryRefreshInterval field value if set, zero value otherwise.
func (o *OauthProfile) GetEndpointDiscoveryRefreshInterval() int32 {
	if o == nil || o.EndpointDiscoveryRefreshInterval == nil {
		var ret int32
		return ret
	}
	return *o.EndpointDiscoveryRefreshInterval
}

// GetEndpointDiscoveryRefreshIntervalOk returns a tuple with the EndpointDiscoveryRefreshInterval field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetEndpointDiscoveryRefreshIntervalOk() (*int32, bool) {
	if o == nil || o.EndpointDiscoveryRefreshInterval == nil {
		return nil, false
	}
	return o.EndpointDiscoveryRefreshInterval, true
}

// HasEndpointDiscoveryRefreshInterval returns a boolean if a field has been set.
func (o *OauthProfile) HasEndpointDiscoveryRefreshInterval() bool {
	if o != nil && o.EndpointDiscoveryRefreshInterval != nil {
		return true
	}

	return false
}

// SetEndpointDiscoveryRefreshInterval gets a reference to the given int32 and assigns it to the EndpointDiscoveryRefreshInterval field.
func (o *OauthProfile) SetEndpointDiscoveryRefreshInterval(v int32) {
	o.EndpointDiscoveryRefreshInterval = &v
}

// GetEndpointIntrospection returns the EndpointIntrospection field value if set, zero value otherwise.
func (o *OauthProfile) GetEndpointIntrospection() string {
	if o == nil || o.EndpointIntrospection == nil {
		var ret string
		return ret
	}
	return *o.EndpointIntrospection
}

// GetEndpointIntrospectionOk returns a tuple with the EndpointIntrospection field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetEndpointIntrospectionOk() (*string, bool) {
	if o == nil || o.EndpointIntrospection == nil {
		return nil, false
	}
	return o.EndpointIntrospection, true
}

// HasEndpointIntrospection returns a boolean if a field has been set.
func (o *OauthProfile) HasEndpointIntrospection() bool {
	if o != nil && o.EndpointIntrospection != nil {
		return true
	}

	return false
}

// SetEndpointIntrospection gets a reference to the given string and assigns it to the EndpointIntrospection field.
func (o *OauthProfile) SetEndpointIntrospection(v string) {
	o.EndpointIntrospection = &v
}

// GetEndpointIntrospectionTimeout returns the EndpointIntrospectionTimeout field value if set, zero value otherwise.
func (o *OauthProfile) GetEndpointIntrospectionTimeout() int32 {
	if o == nil || o.EndpointIntrospectionTimeout == nil {
		var ret int32
		return ret
	}
	return *o.EndpointIntrospectionTimeout
}

// GetEndpointIntrospectionTimeoutOk returns a tuple with the EndpointIntrospectionTimeout field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetEndpointIntrospectionTimeoutOk() (*int32, bool) {
	if o == nil || o.EndpointIntrospectionTimeout == nil {
		return nil, false
	}
	return o.EndpointIntrospectionTimeout, true
}

// HasEndpointIntrospectionTimeout returns a boolean if a field has been set.
func (o *OauthProfile) HasEndpointIntrospectionTimeout() bool {
	if o != nil && o.EndpointIntrospectionTimeout != nil {
		return true
	}

	return false
}

// SetEndpointIntrospectionTimeout gets a reference to the given int32 and assigns it to the EndpointIntrospectionTimeout field.
func (o *OauthProfile) SetEndpointIntrospectionTimeout(v int32) {
	o.EndpointIntrospectionTimeout = &v
}

// GetEndpointJwks returns the EndpointJwks field value if set, zero value otherwise.
func (o *OauthProfile) GetEndpointJwks() string {
	if o == nil || o.EndpointJwks == nil {
		var ret string
		return ret
	}
	return *o.EndpointJwks
}

// GetEndpointJwksOk returns a tuple with the EndpointJwks field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetEndpointJwksOk() (*string, bool) {
	if o == nil || o.EndpointJwks == nil {
		return nil, false
	}
	return o.EndpointJwks, true
}

// HasEndpointJwks returns a boolean if a field has been set.
func (o *OauthProfile) HasEndpointJwks() bool {
	if o != nil && o.EndpointJwks != nil {
		return true
	}

	return false
}

// SetEndpointJwks gets a reference to the given string and assigns it to the EndpointJwks field.
func (o *OauthProfile) SetEndpointJwks(v string) {
	o.EndpointJwks = &v
}

// GetEndpointJwksRefreshInterval returns the EndpointJwksRefreshInterval field value if set, zero value otherwise.
func (o *OauthProfile) GetEndpointJwksRefreshInterval() int32 {
	if o == nil || o.EndpointJwksRefreshInterval == nil {
		var ret int32
		return ret
	}
	return *o.EndpointJwksRefreshInterval
}

// GetEndpointJwksRefreshIntervalOk returns a tuple with the EndpointJwksRefreshInterval field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetEndpointJwksRefreshIntervalOk() (*int32, bool) {
	if o == nil || o.EndpointJwksRefreshInterval == nil {
		return nil, false
	}
	return o.EndpointJwksRefreshInterval, true
}

// HasEndpointJwksRefreshInterval returns a boolean if a field has been set.
func (o *OauthProfile) HasEndpointJwksRefreshInterval() bool {
	if o != nil && o.EndpointJwksRefreshInterval != nil {
		return true
	}

	return false
}

// SetEndpointJwksRefreshInterval gets a reference to the given int32 and assigns it to the EndpointJwksRefreshInterval field.
func (o *OauthProfile) SetEndpointJwksRefreshInterval(v int32) {
	o.EndpointJwksRefreshInterval = &v
}

// GetEndpointToken returns the EndpointToken field value if set, zero value otherwise.
func (o *OauthProfile) GetEndpointToken() string {
	if o == nil || o.EndpointToken == nil {
		var ret string
		return ret
	}
	return *o.EndpointToken
}

// GetEndpointTokenOk returns a tuple with the EndpointToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetEndpointTokenOk() (*string, bool) {
	if o == nil || o.EndpointToken == nil {
		return nil, false
	}
	return o.EndpointToken, true
}

// HasEndpointToken returns a boolean if a field has been set.
func (o *OauthProfile) HasEndpointToken() bool {
	if o != nil && o.EndpointToken != nil {
		return true
	}

	return false
}

// SetEndpointToken gets a reference to the given string and assigns it to the EndpointToken field.
func (o *OauthProfile) SetEndpointToken(v string) {
	o.EndpointToken = &v
}

// GetEndpointTokenTimeout returns the EndpointTokenTimeout field value if set, zero value otherwise.
func (o *OauthProfile) GetEndpointTokenTimeout() int32 {
	if o == nil || o.EndpointTokenTimeout == nil {
		var ret int32
		return ret
	}
	return *o.EndpointTokenTimeout
}

// GetEndpointTokenTimeoutOk returns a tuple with the EndpointTokenTimeout field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetEndpointTokenTimeoutOk() (*int32, bool) {
	if o == nil || o.EndpointTokenTimeout == nil {
		return nil, false
	}
	return o.EndpointTokenTimeout, true
}

// HasEndpointTokenTimeout returns a boolean if a field has been set.
func (o *OauthProfile) HasEndpointTokenTimeout() bool {
	if o != nil && o.EndpointTokenTimeout != nil {
		return true
	}

	return false
}

// SetEndpointTokenTimeout gets a reference to the given int32 and assigns it to the EndpointTokenTimeout field.
func (o *OauthProfile) SetEndpointTokenTimeout(v int32) {
	o.EndpointTokenTimeout = &v
}

// GetEndpointUserinfo returns the EndpointUserinfo field value if set, zero value otherwise.
func (o *OauthProfile) GetEndpointUserinfo() string {
	if o == nil || o.EndpointUserinfo == nil {
		var ret string
		return ret
	}
	return *o.EndpointUserinfo
}

// GetEndpointUserinfoOk returns a tuple with the EndpointUserinfo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetEndpointUserinfoOk() (*string, bool) {
	if o == nil || o.EndpointUserinfo == nil {
		return nil, false
	}
	return o.EndpointUserinfo, true
}

// HasEndpointUserinfo returns a boolean if a field has been set.
func (o *OauthProfile) HasEndpointUserinfo() bool {
	if o != nil && o.EndpointUserinfo != nil {
		return true
	}

	return false
}

// SetEndpointUserinfo gets a reference to the given string and assigns it to the EndpointUserinfo field.
func (o *OauthProfile) SetEndpointUserinfo(v string) {
	o.EndpointUserinfo = &v
}

// GetEndpointUserinfoTimeout returns the EndpointUserinfoTimeout field value if set, zero value otherwise.
func (o *OauthProfile) GetEndpointUserinfoTimeout() int32 {
	if o == nil || o.EndpointUserinfoTimeout == nil {
		var ret int32
		return ret
	}
	return *o.EndpointUserinfoTimeout
}

// GetEndpointUserinfoTimeoutOk returns a tuple with the EndpointUserinfoTimeout field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetEndpointUserinfoTimeoutOk() (*int32, bool) {
	if o == nil || o.EndpointUserinfoTimeout == nil {
		return nil, false
	}
	return o.EndpointUserinfoTimeout, true
}

// HasEndpointUserinfoTimeout returns a boolean if a field has been set.
func (o *OauthProfile) HasEndpointUserinfoTimeout() bool {
	if o != nil && o.EndpointUserinfoTimeout != nil {
		return true
	}

	return false
}

// SetEndpointUserinfoTimeout gets a reference to the given int32 and assigns it to the EndpointUserinfoTimeout field.
func (o *OauthProfile) SetEndpointUserinfoTimeout(v int32) {
	o.EndpointUserinfoTimeout = &v
}

// GetInteractiveEnabled returns the InteractiveEnabled field value if set, zero value otherwise.
func (o *OauthProfile) GetInteractiveEnabled() bool {
	if o == nil || o.InteractiveEnabled == nil {
		var ret bool
		return ret
	}
	return *o.InteractiveEnabled
}

// GetInteractiveEnabledOk returns a tuple with the InteractiveEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetInteractiveEnabledOk() (*bool, bool) {
	if o == nil || o.InteractiveEnabled == nil {
		return nil, false
	}
	return o.InteractiveEnabled, true
}

// HasInteractiveEnabled returns a boolean if a field has been set.
func (o *OauthProfile) HasInteractiveEnabled() bool {
	if o != nil && o.InteractiveEnabled != nil {
		return true
	}

	return false
}

// SetInteractiveEnabled gets a reference to the given bool and assigns it to the InteractiveEnabled field.
func (o *OauthProfile) SetInteractiveEnabled(v bool) {
	o.InteractiveEnabled = &v
}

// GetInteractivePromptForExpiredSession returns the InteractivePromptForExpiredSession field value if set, zero value otherwise.
func (o *OauthProfile) GetInteractivePromptForExpiredSession() string {
	if o == nil || o.InteractivePromptForExpiredSession == nil {
		var ret string
		return ret
	}
	return *o.InteractivePromptForExpiredSession
}

// GetInteractivePromptForExpiredSessionOk returns a tuple with the InteractivePromptForExpiredSession field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetInteractivePromptForExpiredSessionOk() (*string, bool) {
	if o == nil || o.InteractivePromptForExpiredSession == nil {
		return nil, false
	}
	return o.InteractivePromptForExpiredSession, true
}

// HasInteractivePromptForExpiredSession returns a boolean if a field has been set.
func (o *OauthProfile) HasInteractivePromptForExpiredSession() bool {
	if o != nil && o.InteractivePromptForExpiredSession != nil {
		return true
	}

	return false
}

// SetInteractivePromptForExpiredSession gets a reference to the given string and assigns it to the InteractivePromptForExpiredSession field.
func (o *OauthProfile) SetInteractivePromptForExpiredSession(v string) {
	o.InteractivePromptForExpiredSession = &v
}

// GetInteractivePromptForNewSession returns the InteractivePromptForNewSession field value if set, zero value otherwise.
func (o *OauthProfile) GetInteractivePromptForNewSession() string {
	if o == nil || o.InteractivePromptForNewSession == nil {
		var ret string
		return ret
	}
	return *o.InteractivePromptForNewSession
}

// GetInteractivePromptForNewSessionOk returns a tuple with the InteractivePromptForNewSession field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetInteractivePromptForNewSessionOk() (*string, bool) {
	if o == nil || o.InteractivePromptForNewSession == nil {
		return nil, false
	}
	return o.InteractivePromptForNewSession, true
}

// HasInteractivePromptForNewSession returns a boolean if a field has been set.
func (o *OauthProfile) HasInteractivePromptForNewSession() bool {
	if o != nil && o.InteractivePromptForNewSession != nil {
		return true
	}

	return false
}

// SetInteractivePromptForNewSession gets a reference to the given string and assigns it to the InteractivePromptForNewSession field.
func (o *OauthProfile) SetInteractivePromptForNewSession(v string) {
	o.InteractivePromptForNewSession = &v
}

// GetIssuer returns the Issuer field value if set, zero value otherwise.
func (o *OauthProfile) GetIssuer() string {
	if o == nil || o.Issuer == nil {
		var ret string
		return ret
	}
	return *o.Issuer
}

// GetIssuerOk returns a tuple with the Issuer field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetIssuerOk() (*string, bool) {
	if o == nil || o.Issuer == nil {
		return nil, false
	}
	return o.Issuer, true
}

// HasIssuer returns a boolean if a field has been set.
func (o *OauthProfile) HasIssuer() bool {
	if o != nil && o.Issuer != nil {
		return true
	}

	return false
}

// SetIssuer gets a reference to the given string and assigns it to the Issuer field.
func (o *OauthProfile) SetIssuer(v string) {
	o.Issuer = &v
}

// GetOauthProfileName returns the OauthProfileName field value if set, zero value otherwise.
func (o *OauthProfile) GetOauthProfileName() string {
	if o == nil || o.OauthProfileName == nil {
		var ret string
		return ret
	}
	return *o.OauthProfileName
}

// GetOauthProfileNameOk returns a tuple with the OauthProfileName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetOauthProfileNameOk() (*string, bool) {
	if o == nil || o.OauthProfileName == nil {
		return nil, false
	}
	return o.OauthProfileName, true
}

// HasOauthProfileName returns a boolean if a field has been set.
func (o *OauthProfile) HasOauthProfileName() bool {
	if o != nil && o.OauthProfileName != nil {
		return true
	}

	return false
}

// SetOauthProfileName gets a reference to the given string and assigns it to the OauthProfileName field.
func (o *OauthProfile) SetOauthProfileName(v string) {
	o.OauthProfileName = &v
}

// GetOauthRole returns the OauthRole field value if set, zero value otherwise.
func (o *OauthProfile) GetOauthRole() string {
	if o == nil || o.OauthRole == nil {
		var ret string
		return ret
	}
	return *o.OauthRole
}

// GetOauthRoleOk returns a tuple with the OauthRole field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetOauthRoleOk() (*string, bool) {
	if o == nil || o.OauthRole == nil {
		return nil, false
	}
	return o.OauthRole, true
}

// HasOauthRole returns a boolean if a field has been set.
func (o *OauthProfile) HasOauthRole() bool {
	if o != nil && o.OauthRole != nil {
		return true
	}

	return false
}

// SetOauthRole gets a reference to the given string and assigns it to the OauthRole field.
func (o *OauthProfile) SetOauthRole(v string) {
	o.OauthRole = &v
}

// GetResourceServerParseAccessTokenEnabled returns the ResourceServerParseAccessTokenEnabled field value if set, zero value otherwise.
func (o *OauthProfile) GetResourceServerParseAccessTokenEnabled() bool {
	if o == nil || o.ResourceServerParseAccessTokenEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ResourceServerParseAccessTokenEnabled
}

// GetResourceServerParseAccessTokenEnabledOk returns a tuple with the ResourceServerParseAccessTokenEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetResourceServerParseAccessTokenEnabledOk() (*bool, bool) {
	if o == nil || o.ResourceServerParseAccessTokenEnabled == nil {
		return nil, false
	}
	return o.ResourceServerParseAccessTokenEnabled, true
}

// HasResourceServerParseAccessTokenEnabled returns a boolean if a field has been set.
func (o *OauthProfile) HasResourceServerParseAccessTokenEnabled() bool {
	if o != nil && o.ResourceServerParseAccessTokenEnabled != nil {
		return true
	}

	return false
}

// SetResourceServerParseAccessTokenEnabled gets a reference to the given bool and assigns it to the ResourceServerParseAccessTokenEnabled field.
func (o *OauthProfile) SetResourceServerParseAccessTokenEnabled(v bool) {
	o.ResourceServerParseAccessTokenEnabled = &v
}

// GetResourceServerRequiredAudience returns the ResourceServerRequiredAudience field value if set, zero value otherwise.
func (o *OauthProfile) GetResourceServerRequiredAudience() string {
	if o == nil || o.ResourceServerRequiredAudience == nil {
		var ret string
		return ret
	}
	return *o.ResourceServerRequiredAudience
}

// GetResourceServerRequiredAudienceOk returns a tuple with the ResourceServerRequiredAudience field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetResourceServerRequiredAudienceOk() (*string, bool) {
	if o == nil || o.ResourceServerRequiredAudience == nil {
		return nil, false
	}
	return o.ResourceServerRequiredAudience, true
}

// HasResourceServerRequiredAudience returns a boolean if a field has been set.
func (o *OauthProfile) HasResourceServerRequiredAudience() bool {
	if o != nil && o.ResourceServerRequiredAudience != nil {
		return true
	}

	return false
}

// SetResourceServerRequiredAudience gets a reference to the given string and assigns it to the ResourceServerRequiredAudience field.
func (o *OauthProfile) SetResourceServerRequiredAudience(v string) {
	o.ResourceServerRequiredAudience = &v
}

// GetResourceServerRequiredIssuer returns the ResourceServerRequiredIssuer field value if set, zero value otherwise.
func (o *OauthProfile) GetResourceServerRequiredIssuer() string {
	if o == nil || o.ResourceServerRequiredIssuer == nil {
		var ret string
		return ret
	}
	return *o.ResourceServerRequiredIssuer
}

// GetResourceServerRequiredIssuerOk returns a tuple with the ResourceServerRequiredIssuer field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetResourceServerRequiredIssuerOk() (*string, bool) {
	if o == nil || o.ResourceServerRequiredIssuer == nil {
		return nil, false
	}
	return o.ResourceServerRequiredIssuer, true
}

// HasResourceServerRequiredIssuer returns a boolean if a field has been set.
func (o *OauthProfile) HasResourceServerRequiredIssuer() bool {
	if o != nil && o.ResourceServerRequiredIssuer != nil {
		return true
	}

	return false
}

// SetResourceServerRequiredIssuer gets a reference to the given string and assigns it to the ResourceServerRequiredIssuer field.
func (o *OauthProfile) SetResourceServerRequiredIssuer(v string) {
	o.ResourceServerRequiredIssuer = &v
}

// GetResourceServerRequiredScope returns the ResourceServerRequiredScope field value if set, zero value otherwise.
func (o *OauthProfile) GetResourceServerRequiredScope() string {
	if o == nil || o.ResourceServerRequiredScope == nil {
		var ret string
		return ret
	}
	return *o.ResourceServerRequiredScope
}

// GetResourceServerRequiredScopeOk returns a tuple with the ResourceServerRequiredScope field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetResourceServerRequiredScopeOk() (*string, bool) {
	if o == nil || o.ResourceServerRequiredScope == nil {
		return nil, false
	}
	return o.ResourceServerRequiredScope, true
}

// HasResourceServerRequiredScope returns a boolean if a field has been set.
func (o *OauthProfile) HasResourceServerRequiredScope() bool {
	if o != nil && o.ResourceServerRequiredScope != nil {
		return true
	}

	return false
}

// SetResourceServerRequiredScope gets a reference to the given string and assigns it to the ResourceServerRequiredScope field.
func (o *OauthProfile) SetResourceServerRequiredScope(v string) {
	o.ResourceServerRequiredScope = &v
}

// GetResourceServerRequiredType returns the ResourceServerRequiredType field value if set, zero value otherwise.
func (o *OauthProfile) GetResourceServerRequiredType() string {
	if o == nil || o.ResourceServerRequiredType == nil {
		var ret string
		return ret
	}
	return *o.ResourceServerRequiredType
}

// GetResourceServerRequiredTypeOk returns a tuple with the ResourceServerRequiredType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetResourceServerRequiredTypeOk() (*string, bool) {
	if o == nil || o.ResourceServerRequiredType == nil {
		return nil, false
	}
	return o.ResourceServerRequiredType, true
}

// HasResourceServerRequiredType returns a boolean if a field has been set.
func (o *OauthProfile) HasResourceServerRequiredType() bool {
	if o != nil && o.ResourceServerRequiredType != nil {
		return true
	}

	return false
}

// SetResourceServerRequiredType gets a reference to the given string and assigns it to the ResourceServerRequiredType field.
func (o *OauthProfile) SetResourceServerRequiredType(v string) {
	o.ResourceServerRequiredType = &v
}

// GetResourceServerValidateAudienceEnabled returns the ResourceServerValidateAudienceEnabled field value if set, zero value otherwise.
func (o *OauthProfile) GetResourceServerValidateAudienceEnabled() bool {
	if o == nil || o.ResourceServerValidateAudienceEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ResourceServerValidateAudienceEnabled
}

// GetResourceServerValidateAudienceEnabledOk returns a tuple with the ResourceServerValidateAudienceEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetResourceServerValidateAudienceEnabledOk() (*bool, bool) {
	if o == nil || o.ResourceServerValidateAudienceEnabled == nil {
		return nil, false
	}
	return o.ResourceServerValidateAudienceEnabled, true
}

// HasResourceServerValidateAudienceEnabled returns a boolean if a field has been set.
func (o *OauthProfile) HasResourceServerValidateAudienceEnabled() bool {
	if o != nil && o.ResourceServerValidateAudienceEnabled != nil {
		return true
	}

	return false
}

// SetResourceServerValidateAudienceEnabled gets a reference to the given bool and assigns it to the ResourceServerValidateAudienceEnabled field.
func (o *OauthProfile) SetResourceServerValidateAudienceEnabled(v bool) {
	o.ResourceServerValidateAudienceEnabled = &v
}

// GetResourceServerValidateIssuerEnabled returns the ResourceServerValidateIssuerEnabled field value if set, zero value otherwise.
func (o *OauthProfile) GetResourceServerValidateIssuerEnabled() bool {
	if o == nil || o.ResourceServerValidateIssuerEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ResourceServerValidateIssuerEnabled
}

// GetResourceServerValidateIssuerEnabledOk returns a tuple with the ResourceServerValidateIssuerEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetResourceServerValidateIssuerEnabledOk() (*bool, bool) {
	if o == nil || o.ResourceServerValidateIssuerEnabled == nil {
		return nil, false
	}
	return o.ResourceServerValidateIssuerEnabled, true
}

// HasResourceServerValidateIssuerEnabled returns a boolean if a field has been set.
func (o *OauthProfile) HasResourceServerValidateIssuerEnabled() bool {
	if o != nil && o.ResourceServerValidateIssuerEnabled != nil {
		return true
	}

	return false
}

// SetResourceServerValidateIssuerEnabled gets a reference to the given bool and assigns it to the ResourceServerValidateIssuerEnabled field.
func (o *OauthProfile) SetResourceServerValidateIssuerEnabled(v bool) {
	o.ResourceServerValidateIssuerEnabled = &v
}

// GetResourceServerValidateScopeEnabled returns the ResourceServerValidateScopeEnabled field value if set, zero value otherwise.
func (o *OauthProfile) GetResourceServerValidateScopeEnabled() bool {
	if o == nil || o.ResourceServerValidateScopeEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ResourceServerValidateScopeEnabled
}

// GetResourceServerValidateScopeEnabledOk returns a tuple with the ResourceServerValidateScopeEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetResourceServerValidateScopeEnabledOk() (*bool, bool) {
	if o == nil || o.ResourceServerValidateScopeEnabled == nil {
		return nil, false
	}
	return o.ResourceServerValidateScopeEnabled, true
}

// HasResourceServerValidateScopeEnabled returns a boolean if a field has been set.
func (o *OauthProfile) HasResourceServerValidateScopeEnabled() bool {
	if o != nil && o.ResourceServerValidateScopeEnabled != nil {
		return true
	}

	return false
}

// SetResourceServerValidateScopeEnabled gets a reference to the given bool and assigns it to the ResourceServerValidateScopeEnabled field.
func (o *OauthProfile) SetResourceServerValidateScopeEnabled(v bool) {
	o.ResourceServerValidateScopeEnabled = &v
}

// GetResourceServerValidateTypeEnabled returns the ResourceServerValidateTypeEnabled field value if set, zero value otherwise.
func (o *OauthProfile) GetResourceServerValidateTypeEnabled() bool {
	if o == nil || o.ResourceServerValidateTypeEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ResourceServerValidateTypeEnabled
}

// GetResourceServerValidateTypeEnabledOk returns a tuple with the ResourceServerValidateTypeEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetResourceServerValidateTypeEnabledOk() (*bool, bool) {
	if o == nil || o.ResourceServerValidateTypeEnabled == nil {
		return nil, false
	}
	return o.ResourceServerValidateTypeEnabled, true
}

// HasResourceServerValidateTypeEnabled returns a boolean if a field has been set.
func (o *OauthProfile) HasResourceServerValidateTypeEnabled() bool {
	if o != nil && o.ResourceServerValidateTypeEnabled != nil {
		return true
	}

	return false
}

// SetResourceServerValidateTypeEnabled gets a reference to the given bool and assigns it to the ResourceServerValidateTypeEnabled field.
func (o *OauthProfile) SetResourceServerValidateTypeEnabled(v bool) {
	o.ResourceServerValidateTypeEnabled = &v
}

// GetSempEnabled returns the SempEnabled field value if set, zero value otherwise.
func (o *OauthProfile) GetSempEnabled() bool {
	if o == nil || o.SempEnabled == nil {
		var ret bool
		return ret
	}
	return *o.SempEnabled
}

// GetSempEnabledOk returns a tuple with the SempEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetSempEnabledOk() (*bool, bool) {
	if o == nil || o.SempEnabled == nil {
		return nil, false
	}
	return o.SempEnabled, true
}

// HasSempEnabled returns a boolean if a field has been set.
func (o *OauthProfile) HasSempEnabled() bool {
	if o != nil && o.SempEnabled != nil {
		return true
	}

	return false
}

// SetSempEnabled gets a reference to the given bool and assigns it to the SempEnabled field.
func (o *OauthProfile) SetSempEnabled(v bool) {
	o.SempEnabled = &v
}

// GetUsernameClaimName returns the UsernameClaimName field value if set, zero value otherwise.
func (o *OauthProfile) GetUsernameClaimName() string {
	if o == nil || o.UsernameClaimName == nil {
		var ret string
		return ret
	}
	return *o.UsernameClaimName
}

// GetUsernameClaimNameOk returns a tuple with the UsernameClaimName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OauthProfile) GetUsernameClaimNameOk() (*string, bool) {
	if o == nil || o.UsernameClaimName == nil {
		return nil, false
	}
	return o.UsernameClaimName, true
}

// HasUsernameClaimName returns a boolean if a field has been set.
func (o *OauthProfile) HasUsernameClaimName() bool {
	if o != nil && o.UsernameClaimName != nil {
		return true
	}

	return false
}

// SetUsernameClaimName gets a reference to the given string and assigns it to the UsernameClaimName field.
func (o *OauthProfile) SetUsernameClaimName(v string) {
	o.UsernameClaimName = &v
}

func (o OauthProfile) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AccessLevelGroupsClaimName != nil {
		toSerialize["accessLevelGroupsClaimName"] = o.AccessLevelGroupsClaimName
	}
	if o.ClientId != nil {
		toSerialize["clientId"] = o.ClientId
	}
	if o.ClientRedirectUri != nil {
		toSerialize["clientRedirectUri"] = o.ClientRedirectUri
	}
	if o.ClientRequiredType != nil {
		toSerialize["clientRequiredType"] = o.ClientRequiredType
	}
	if o.ClientScope != nil {
		toSerialize["clientScope"] = o.ClientScope
	}
	if o.ClientSecret != nil {
		toSerialize["clientSecret"] = o.ClientSecret
	}
	if o.ClientValidateTypeEnabled != nil {
		toSerialize["clientValidateTypeEnabled"] = o.ClientValidateTypeEnabled
	}
	if o.DefaultGlobalAccessLevel != nil {
		toSerialize["defaultGlobalAccessLevel"] = o.DefaultGlobalAccessLevel
	}
	if o.DefaultMsgVpnAccessLevel != nil {
		toSerialize["defaultMsgVpnAccessLevel"] = o.DefaultMsgVpnAccessLevel
	}
	if o.DisplayName != nil {
		toSerialize["displayName"] = o.DisplayName
	}
	if o.Enabled != nil {
		toSerialize["enabled"] = o.Enabled
	}
	if o.EndpointAuthorization != nil {
		toSerialize["endpointAuthorization"] = o.EndpointAuthorization
	}
	if o.EndpointDiscovery != nil {
		toSerialize["endpointDiscovery"] = o.EndpointDiscovery
	}
	if o.EndpointDiscoveryRefreshInterval != nil {
		toSerialize["endpointDiscoveryRefreshInterval"] = o.EndpointDiscoveryRefreshInterval
	}
	if o.EndpointIntrospection != nil {
		toSerialize["endpointIntrospection"] = o.EndpointIntrospection
	}
	if o.EndpointIntrospectionTimeout != nil {
		toSerialize["endpointIntrospectionTimeout"] = o.EndpointIntrospectionTimeout
	}
	if o.EndpointJwks != nil {
		toSerialize["endpointJwks"] = o.EndpointJwks
	}
	if o.EndpointJwksRefreshInterval != nil {
		toSerialize["endpointJwksRefreshInterval"] = o.EndpointJwksRefreshInterval
	}
	if o.EndpointToken != nil {
		toSerialize["endpointToken"] = o.EndpointToken
	}
	if o.EndpointTokenTimeout != nil {
		toSerialize["endpointTokenTimeout"] = o.EndpointTokenTimeout
	}
	if o.EndpointUserinfo != nil {
		toSerialize["endpointUserinfo"] = o.EndpointUserinfo
	}
	if o.EndpointUserinfoTimeout != nil {
		toSerialize["endpointUserinfoTimeout"] = o.EndpointUserinfoTimeout
	}
	if o.InteractiveEnabled != nil {
		toSerialize["interactiveEnabled"] = o.InteractiveEnabled
	}
	if o.InteractivePromptForExpiredSession != nil {
		toSerialize["interactivePromptForExpiredSession"] = o.InteractivePromptForExpiredSession
	}
	if o.InteractivePromptForNewSession != nil {
		toSerialize["interactivePromptForNewSession"] = o.InteractivePromptForNewSession
	}
	if o.Issuer != nil {
		toSerialize["issuer"] = o.Issuer
	}
	if o.OauthProfileName != nil {
		toSerialize["oauthProfileName"] = o.OauthProfileName
	}
	if o.OauthRole != nil {
		toSerialize["oauthRole"] = o.OauthRole
	}
	if o.ResourceServerParseAccessTokenEnabled != nil {
		toSerialize["resourceServerParseAccessTokenEnabled"] = o.ResourceServerParseAccessTokenEnabled
	}
	if o.ResourceServerRequiredAudience != nil {
		toSerialize["resourceServerRequiredAudience"] = o.ResourceServerRequiredAudience
	}
	if o.ResourceServerRequiredIssuer != nil {
		toSerialize["resourceServerRequiredIssuer"] = o.ResourceServerRequiredIssuer
	}
	if o.ResourceServerRequiredScope != nil {
		toSerialize["resourceServerRequiredScope"] = o.ResourceServerRequiredScope
	}
	if o.ResourceServerRequiredType != nil {
		toSerialize["resourceServerRequiredType"] = o.ResourceServerRequiredType
	}
	if o.ResourceServerValidateAudienceEnabled != nil {
		toSerialize["resourceServerValidateAudienceEnabled"] = o.ResourceServerValidateAudienceEnabled
	}
	if o.ResourceServerValidateIssuerEnabled != nil {
		toSerialize["resourceServerValidateIssuerEnabled"] = o.ResourceServerValidateIssuerEnabled
	}
	if o.ResourceServerValidateScopeEnabled != nil {
		toSerialize["resourceServerValidateScopeEnabled"] = o.ResourceServerValidateScopeEnabled
	}
	if o.ResourceServerValidateTypeEnabled != nil {
		toSerialize["resourceServerValidateTypeEnabled"] = o.ResourceServerValidateTypeEnabled
	}
	if o.SempEnabled != nil {
		toSerialize["sempEnabled"] = o.SempEnabled
	}
	if o.UsernameClaimName != nil {
		toSerialize["usernameClaimName"] = o.UsernameClaimName
	}
	return json.Marshal(toSerialize)
}

type NullableOauthProfile struct {
	value *OauthProfile
	isSet bool
}

func (v NullableOauthProfile) Get() *OauthProfile {
	return v.value
}

func (v *NullableOauthProfile) Set(val *OauthProfile) {
	v.value = val
	v.isSet = true
}

func (v NullableOauthProfile) IsSet() bool {
	return v.isSet
}

func (v *NullableOauthProfile) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableOauthProfile(val *OauthProfile) *NullableOauthProfile {
	return &NullableOauthProfile{value: val, isSet: true}
}

func (v NullableOauthProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableOauthProfile) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


