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

// MsgVpnRestDeliveryPointRestConsumer struct for MsgVpnRestDeliveryPointRestConsumer
type MsgVpnRestDeliveryPointRestConsumer struct {
	// The AWS access key id. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.26.
	AuthenticationAwsAccessKeyId *string `json:"authenticationAwsAccessKeyId,omitempty"`
	// The AWS region id. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.26.
	AuthenticationAwsRegion *string `json:"authenticationAwsRegion,omitempty"`
	// The AWS secret access key. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.26.
	AuthenticationAwsSecretAccessKey *string `json:"authenticationAwsSecretAccessKey,omitempty"`
	// The AWS service id. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.26.
	AuthenticationAwsService *string `json:"authenticationAwsService,omitempty"`
	// The PEM formatted content for the client certificate that the REST Consumer will present to the REST host. It must consist of a private key and between one and three certificates comprising the certificate trust chain. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changing this attribute requires an HTTPS connection. The default value is `\"\"`. Available since 2.9.
	AuthenticationClientCertContent *string `json:"authenticationClientCertContent,omitempty"`
	// The password for the client certificate. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changing this attribute requires an HTTPS connection. The default value is `\"\"`. Available since 2.9.
	AuthenticationClientCertPassword *string `json:"authenticationClientCertPassword,omitempty"`
	// The password for the username. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.
	AuthenticationHttpBasicPassword *string `json:"authenticationHttpBasicPassword,omitempty"`
	// The username that the REST Consumer will use to login to the REST host. Normally a username is only configured when basic authentication is selected for the REST Consumer. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.
	AuthenticationHttpBasicUsername *string `json:"authenticationHttpBasicUsername,omitempty"`
	// The authentication header name. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.15.
	AuthenticationHttpHeaderName *string `json:"authenticationHttpHeaderName,omitempty"`
	// The authentication header value. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.15.
	AuthenticationHttpHeaderValue *string `json:"authenticationHttpHeaderValue,omitempty"`
	// The OAuth client ID. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.19.
	AuthenticationOauthClientId *string `json:"authenticationOauthClientId,omitempty"`
	// The OAuth scope. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.19.
	AuthenticationOauthClientScope *string `json:"authenticationOauthClientScope,omitempty"`
	// The OAuth client secret. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.19.
	AuthenticationOauthClientSecret *string `json:"authenticationOauthClientSecret,omitempty"`
	// The OAuth token endpoint URL that the REST Consumer will use to request a token for login to the REST host. Must begin with \"https\". Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.19.
	AuthenticationOauthClientTokenEndpoint *string `json:"authenticationOauthClientTokenEndpoint,omitempty"`
	// The OAuth secret key used to sign the token request JWT. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.21.
	AuthenticationOauthJwtSecretKey *string `json:"authenticationOauthJwtSecretKey,omitempty"`
	// The OAuth token endpoint URL that the REST Consumer will use to request a token for login to the REST host. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.21.
	AuthenticationOauthJwtTokenEndpoint *string `json:"authenticationOauthJwtTokenEndpoint,omitempty"`
	// The authentication scheme used by the REST Consumer to login to the REST host. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"none\"`. The allowed values and their meaning are:  <pre> \"none\" - Login with no authentication. This may be useful for anonymous connections or when a REST Consumer does not require authentication. \"http-basic\" - Login with a username and optional password according to HTTP Basic authentication as per RFC2616. \"client-certificate\" - Login with a client TLS certificate as per RFC5246. Client certificate authentication is only available on TLS connections. \"http-header\" - Login with a specified HTTP header. \"oauth-client\" - Login with OAuth 2.0 client credentials. \"oauth-jwt\" - Login with OAuth (RFC 7523 JWT Profile). \"transparent\" - Login using the Authorization header from the message properties, if present. Transparent authentication passes along existing Authorization header metadata instead of discarding it. Note that if the message is coming from a REST producer, the REST service must be configured to forward the Authorization header. \"aws\" - Login using AWS Signature Version 4 authentication (AWS4-HMAC-SHA256). </pre> 
	AuthenticationScheme *string `json:"authenticationScheme,omitempty"`
	// Enable or disable the REST Consumer. When disabled, no connections are initiated or messages delivered to this particular REST Consumer. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	Enabled *bool `json:"enabled,omitempty"`
	// The HTTP method to use (POST or PUT). This is used only when operating in the REST service \"messaging\" mode and is ignored in \"gateway\" mode. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"post\"`. The allowed values and their meaning are:  <pre> \"post\" - Use the POST HTTP method. \"put\" - Use the PUT HTTP method. </pre>  Available since 2.17.
	HttpMethod *string `json:"httpMethod,omitempty"`
	// The interface that will be used for all outgoing connections associated with the REST Consumer. When unspecified, an interface is automatically chosen. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.
	LocalInterface *string `json:"localInterface,omitempty"`
	// The maximum amount of time (in seconds) to wait for an HTTP POST response from the REST Consumer. Once this time is exceeded, the TCP connection is reset. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `30`.
	MaxPostWaitTime *int32 `json:"maxPostWaitTime,omitempty"`
	// The name of the Message VPN.
	MsgVpnName *string `json:"msgVpnName,omitempty"`
	// The number of concurrent TCP connections open to the REST Consumer. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.
	OutgoingConnectionCount *int32 `json:"outgoingConnectionCount,omitempty"`
	// The IP address or DNS name to which the broker is to connect to deliver messages for the REST Consumer. A host value must be configured for the REST Consumer to be operationally up. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.
	RemoteHost *string `json:"remoteHost,omitempty"`
	// The port associated with the host of the REST Consumer. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `8080`.
	RemotePort *int64 `json:"remotePort,omitempty"`
	// The name of the REST Consumer.
	RestConsumerName *string `json:"restConsumerName,omitempty"`
	// The name of the REST Delivery Point.
	RestDeliveryPointName *string `json:"restDeliveryPointName,omitempty"`
	// The number of seconds that must pass before retrying the remote REST Consumer connection. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.
	RetryDelay *int32 `json:"retryDelay,omitempty"`
	// The colon-separated list of cipher suites the REST Consumer uses in its encrypted connection. The value `\"default\"` implies all supported suites ordered from most secure to least secure. The list of default cipher suites is available in the `tlsCipherSuiteMsgBackboneDefaultList` attribute of the Broker object in the Monitoring API. The REST Consumer should choose the first suite from this list that it supports. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"default\"`.
	TlsCipherSuiteList *string `json:"tlsCipherSuiteList,omitempty"`
	// Enable or disable encryption (TLS) for the REST Consumer. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	TlsEnabled *bool `json:"tlsEnabled,omitempty"`
}

// NewMsgVpnRestDeliveryPointRestConsumer instantiates a new MsgVpnRestDeliveryPointRestConsumer object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMsgVpnRestDeliveryPointRestConsumer() *MsgVpnRestDeliveryPointRestConsumer {
	this := MsgVpnRestDeliveryPointRestConsumer{}
	return &this
}

// NewMsgVpnRestDeliveryPointRestConsumerWithDefaults instantiates a new MsgVpnRestDeliveryPointRestConsumer object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMsgVpnRestDeliveryPointRestConsumerWithDefaults() *MsgVpnRestDeliveryPointRestConsumer {
	this := MsgVpnRestDeliveryPointRestConsumer{}
	return &this
}

// GetAuthenticationAwsAccessKeyId returns the AuthenticationAwsAccessKeyId field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationAwsAccessKeyId() string {
	if o == nil || o.AuthenticationAwsAccessKeyId == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationAwsAccessKeyId
}

// GetAuthenticationAwsAccessKeyIdOk returns a tuple with the AuthenticationAwsAccessKeyId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationAwsAccessKeyIdOk() (*string, bool) {
	if o == nil || o.AuthenticationAwsAccessKeyId == nil {
		return nil, false
	}
	return o.AuthenticationAwsAccessKeyId, true
}

// HasAuthenticationAwsAccessKeyId returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationAwsAccessKeyId() bool {
	if o != nil && o.AuthenticationAwsAccessKeyId != nil {
		return true
	}

	return false
}

// SetAuthenticationAwsAccessKeyId gets a reference to the given string and assigns it to the AuthenticationAwsAccessKeyId field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationAwsAccessKeyId(v string) {
	o.AuthenticationAwsAccessKeyId = &v
}

// GetAuthenticationAwsRegion returns the AuthenticationAwsRegion field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationAwsRegion() string {
	if o == nil || o.AuthenticationAwsRegion == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationAwsRegion
}

// GetAuthenticationAwsRegionOk returns a tuple with the AuthenticationAwsRegion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationAwsRegionOk() (*string, bool) {
	if o == nil || o.AuthenticationAwsRegion == nil {
		return nil, false
	}
	return o.AuthenticationAwsRegion, true
}

// HasAuthenticationAwsRegion returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationAwsRegion() bool {
	if o != nil && o.AuthenticationAwsRegion != nil {
		return true
	}

	return false
}

// SetAuthenticationAwsRegion gets a reference to the given string and assigns it to the AuthenticationAwsRegion field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationAwsRegion(v string) {
	o.AuthenticationAwsRegion = &v
}

// GetAuthenticationAwsSecretAccessKey returns the AuthenticationAwsSecretAccessKey field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationAwsSecretAccessKey() string {
	if o == nil || o.AuthenticationAwsSecretAccessKey == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationAwsSecretAccessKey
}

// GetAuthenticationAwsSecretAccessKeyOk returns a tuple with the AuthenticationAwsSecretAccessKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationAwsSecretAccessKeyOk() (*string, bool) {
	if o == nil || o.AuthenticationAwsSecretAccessKey == nil {
		return nil, false
	}
	return o.AuthenticationAwsSecretAccessKey, true
}

// HasAuthenticationAwsSecretAccessKey returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationAwsSecretAccessKey() bool {
	if o != nil && o.AuthenticationAwsSecretAccessKey != nil {
		return true
	}

	return false
}

// SetAuthenticationAwsSecretAccessKey gets a reference to the given string and assigns it to the AuthenticationAwsSecretAccessKey field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationAwsSecretAccessKey(v string) {
	o.AuthenticationAwsSecretAccessKey = &v
}

// GetAuthenticationAwsService returns the AuthenticationAwsService field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationAwsService() string {
	if o == nil || o.AuthenticationAwsService == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationAwsService
}

// GetAuthenticationAwsServiceOk returns a tuple with the AuthenticationAwsService field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationAwsServiceOk() (*string, bool) {
	if o == nil || o.AuthenticationAwsService == nil {
		return nil, false
	}
	return o.AuthenticationAwsService, true
}

// HasAuthenticationAwsService returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationAwsService() bool {
	if o != nil && o.AuthenticationAwsService != nil {
		return true
	}

	return false
}

// SetAuthenticationAwsService gets a reference to the given string and assigns it to the AuthenticationAwsService field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationAwsService(v string) {
	o.AuthenticationAwsService = &v
}

// GetAuthenticationClientCertContent returns the AuthenticationClientCertContent field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationClientCertContent() string {
	if o == nil || o.AuthenticationClientCertContent == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationClientCertContent
}

// GetAuthenticationClientCertContentOk returns a tuple with the AuthenticationClientCertContent field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationClientCertContentOk() (*string, bool) {
	if o == nil || o.AuthenticationClientCertContent == nil {
		return nil, false
	}
	return o.AuthenticationClientCertContent, true
}

// HasAuthenticationClientCertContent returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationClientCertContent() bool {
	if o != nil && o.AuthenticationClientCertContent != nil {
		return true
	}

	return false
}

// SetAuthenticationClientCertContent gets a reference to the given string and assigns it to the AuthenticationClientCertContent field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationClientCertContent(v string) {
	o.AuthenticationClientCertContent = &v
}

// GetAuthenticationClientCertPassword returns the AuthenticationClientCertPassword field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationClientCertPassword() string {
	if o == nil || o.AuthenticationClientCertPassword == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationClientCertPassword
}

// GetAuthenticationClientCertPasswordOk returns a tuple with the AuthenticationClientCertPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationClientCertPasswordOk() (*string, bool) {
	if o == nil || o.AuthenticationClientCertPassword == nil {
		return nil, false
	}
	return o.AuthenticationClientCertPassword, true
}

// HasAuthenticationClientCertPassword returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationClientCertPassword() bool {
	if o != nil && o.AuthenticationClientCertPassword != nil {
		return true
	}

	return false
}

// SetAuthenticationClientCertPassword gets a reference to the given string and assigns it to the AuthenticationClientCertPassword field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationClientCertPassword(v string) {
	o.AuthenticationClientCertPassword = &v
}

// GetAuthenticationHttpBasicPassword returns the AuthenticationHttpBasicPassword field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationHttpBasicPassword() string {
	if o == nil || o.AuthenticationHttpBasicPassword == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationHttpBasicPassword
}

// GetAuthenticationHttpBasicPasswordOk returns a tuple with the AuthenticationHttpBasicPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationHttpBasicPasswordOk() (*string, bool) {
	if o == nil || o.AuthenticationHttpBasicPassword == nil {
		return nil, false
	}
	return o.AuthenticationHttpBasicPassword, true
}

// HasAuthenticationHttpBasicPassword returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationHttpBasicPassword() bool {
	if o != nil && o.AuthenticationHttpBasicPassword != nil {
		return true
	}

	return false
}

// SetAuthenticationHttpBasicPassword gets a reference to the given string and assigns it to the AuthenticationHttpBasicPassword field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationHttpBasicPassword(v string) {
	o.AuthenticationHttpBasicPassword = &v
}

// GetAuthenticationHttpBasicUsername returns the AuthenticationHttpBasicUsername field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationHttpBasicUsername() string {
	if o == nil || o.AuthenticationHttpBasicUsername == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationHttpBasicUsername
}

// GetAuthenticationHttpBasicUsernameOk returns a tuple with the AuthenticationHttpBasicUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationHttpBasicUsernameOk() (*string, bool) {
	if o == nil || o.AuthenticationHttpBasicUsername == nil {
		return nil, false
	}
	return o.AuthenticationHttpBasicUsername, true
}

// HasAuthenticationHttpBasicUsername returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationHttpBasicUsername() bool {
	if o != nil && o.AuthenticationHttpBasicUsername != nil {
		return true
	}

	return false
}

// SetAuthenticationHttpBasicUsername gets a reference to the given string and assigns it to the AuthenticationHttpBasicUsername field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationHttpBasicUsername(v string) {
	o.AuthenticationHttpBasicUsername = &v
}

// GetAuthenticationHttpHeaderName returns the AuthenticationHttpHeaderName field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationHttpHeaderName() string {
	if o == nil || o.AuthenticationHttpHeaderName == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationHttpHeaderName
}

// GetAuthenticationHttpHeaderNameOk returns a tuple with the AuthenticationHttpHeaderName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationHttpHeaderNameOk() (*string, bool) {
	if o == nil || o.AuthenticationHttpHeaderName == nil {
		return nil, false
	}
	return o.AuthenticationHttpHeaderName, true
}

// HasAuthenticationHttpHeaderName returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationHttpHeaderName() bool {
	if o != nil && o.AuthenticationHttpHeaderName != nil {
		return true
	}

	return false
}

// SetAuthenticationHttpHeaderName gets a reference to the given string and assigns it to the AuthenticationHttpHeaderName field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationHttpHeaderName(v string) {
	o.AuthenticationHttpHeaderName = &v
}

// GetAuthenticationHttpHeaderValue returns the AuthenticationHttpHeaderValue field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationHttpHeaderValue() string {
	if o == nil || o.AuthenticationHttpHeaderValue == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationHttpHeaderValue
}

// GetAuthenticationHttpHeaderValueOk returns a tuple with the AuthenticationHttpHeaderValue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationHttpHeaderValueOk() (*string, bool) {
	if o == nil || o.AuthenticationHttpHeaderValue == nil {
		return nil, false
	}
	return o.AuthenticationHttpHeaderValue, true
}

// HasAuthenticationHttpHeaderValue returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationHttpHeaderValue() bool {
	if o != nil && o.AuthenticationHttpHeaderValue != nil {
		return true
	}

	return false
}

// SetAuthenticationHttpHeaderValue gets a reference to the given string and assigns it to the AuthenticationHttpHeaderValue field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationHttpHeaderValue(v string) {
	o.AuthenticationHttpHeaderValue = &v
}

// GetAuthenticationOauthClientId returns the AuthenticationOauthClientId field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationOauthClientId() string {
	if o == nil || o.AuthenticationOauthClientId == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationOauthClientId
}

// GetAuthenticationOauthClientIdOk returns a tuple with the AuthenticationOauthClientId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationOauthClientIdOk() (*string, bool) {
	if o == nil || o.AuthenticationOauthClientId == nil {
		return nil, false
	}
	return o.AuthenticationOauthClientId, true
}

// HasAuthenticationOauthClientId returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationOauthClientId() bool {
	if o != nil && o.AuthenticationOauthClientId != nil {
		return true
	}

	return false
}

// SetAuthenticationOauthClientId gets a reference to the given string and assigns it to the AuthenticationOauthClientId field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationOauthClientId(v string) {
	o.AuthenticationOauthClientId = &v
}

// GetAuthenticationOauthClientScope returns the AuthenticationOauthClientScope field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationOauthClientScope() string {
	if o == nil || o.AuthenticationOauthClientScope == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationOauthClientScope
}

// GetAuthenticationOauthClientScopeOk returns a tuple with the AuthenticationOauthClientScope field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationOauthClientScopeOk() (*string, bool) {
	if o == nil || o.AuthenticationOauthClientScope == nil {
		return nil, false
	}
	return o.AuthenticationOauthClientScope, true
}

// HasAuthenticationOauthClientScope returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationOauthClientScope() bool {
	if o != nil && o.AuthenticationOauthClientScope != nil {
		return true
	}

	return false
}

// SetAuthenticationOauthClientScope gets a reference to the given string and assigns it to the AuthenticationOauthClientScope field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationOauthClientScope(v string) {
	o.AuthenticationOauthClientScope = &v
}

// GetAuthenticationOauthClientSecret returns the AuthenticationOauthClientSecret field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationOauthClientSecret() string {
	if o == nil || o.AuthenticationOauthClientSecret == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationOauthClientSecret
}

// GetAuthenticationOauthClientSecretOk returns a tuple with the AuthenticationOauthClientSecret field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationOauthClientSecretOk() (*string, bool) {
	if o == nil || o.AuthenticationOauthClientSecret == nil {
		return nil, false
	}
	return o.AuthenticationOauthClientSecret, true
}

// HasAuthenticationOauthClientSecret returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationOauthClientSecret() bool {
	if o != nil && o.AuthenticationOauthClientSecret != nil {
		return true
	}

	return false
}

// SetAuthenticationOauthClientSecret gets a reference to the given string and assigns it to the AuthenticationOauthClientSecret field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationOauthClientSecret(v string) {
	o.AuthenticationOauthClientSecret = &v
}

// GetAuthenticationOauthClientTokenEndpoint returns the AuthenticationOauthClientTokenEndpoint field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationOauthClientTokenEndpoint() string {
	if o == nil || o.AuthenticationOauthClientTokenEndpoint == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationOauthClientTokenEndpoint
}

// GetAuthenticationOauthClientTokenEndpointOk returns a tuple with the AuthenticationOauthClientTokenEndpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationOauthClientTokenEndpointOk() (*string, bool) {
	if o == nil || o.AuthenticationOauthClientTokenEndpoint == nil {
		return nil, false
	}
	return o.AuthenticationOauthClientTokenEndpoint, true
}

// HasAuthenticationOauthClientTokenEndpoint returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationOauthClientTokenEndpoint() bool {
	if o != nil && o.AuthenticationOauthClientTokenEndpoint != nil {
		return true
	}

	return false
}

// SetAuthenticationOauthClientTokenEndpoint gets a reference to the given string and assigns it to the AuthenticationOauthClientTokenEndpoint field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationOauthClientTokenEndpoint(v string) {
	o.AuthenticationOauthClientTokenEndpoint = &v
}

// GetAuthenticationOauthJwtSecretKey returns the AuthenticationOauthJwtSecretKey field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationOauthJwtSecretKey() string {
	if o == nil || o.AuthenticationOauthJwtSecretKey == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationOauthJwtSecretKey
}

// GetAuthenticationOauthJwtSecretKeyOk returns a tuple with the AuthenticationOauthJwtSecretKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationOauthJwtSecretKeyOk() (*string, bool) {
	if o == nil || o.AuthenticationOauthJwtSecretKey == nil {
		return nil, false
	}
	return o.AuthenticationOauthJwtSecretKey, true
}

// HasAuthenticationOauthJwtSecretKey returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationOauthJwtSecretKey() bool {
	if o != nil && o.AuthenticationOauthJwtSecretKey != nil {
		return true
	}

	return false
}

// SetAuthenticationOauthJwtSecretKey gets a reference to the given string and assigns it to the AuthenticationOauthJwtSecretKey field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationOauthJwtSecretKey(v string) {
	o.AuthenticationOauthJwtSecretKey = &v
}

// GetAuthenticationOauthJwtTokenEndpoint returns the AuthenticationOauthJwtTokenEndpoint field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationOauthJwtTokenEndpoint() string {
	if o == nil || o.AuthenticationOauthJwtTokenEndpoint == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationOauthJwtTokenEndpoint
}

// GetAuthenticationOauthJwtTokenEndpointOk returns a tuple with the AuthenticationOauthJwtTokenEndpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationOauthJwtTokenEndpointOk() (*string, bool) {
	if o == nil || o.AuthenticationOauthJwtTokenEndpoint == nil {
		return nil, false
	}
	return o.AuthenticationOauthJwtTokenEndpoint, true
}

// HasAuthenticationOauthJwtTokenEndpoint returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationOauthJwtTokenEndpoint() bool {
	if o != nil && o.AuthenticationOauthJwtTokenEndpoint != nil {
		return true
	}

	return false
}

// SetAuthenticationOauthJwtTokenEndpoint gets a reference to the given string and assigns it to the AuthenticationOauthJwtTokenEndpoint field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationOauthJwtTokenEndpoint(v string) {
	o.AuthenticationOauthJwtTokenEndpoint = &v
}

// GetAuthenticationScheme returns the AuthenticationScheme field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationScheme() string {
	if o == nil || o.AuthenticationScheme == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationScheme
}

// GetAuthenticationSchemeOk returns a tuple with the AuthenticationScheme field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetAuthenticationSchemeOk() (*string, bool) {
	if o == nil || o.AuthenticationScheme == nil {
		return nil, false
	}
	return o.AuthenticationScheme, true
}

// HasAuthenticationScheme returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasAuthenticationScheme() bool {
	if o != nil && o.AuthenticationScheme != nil {
		return true
	}

	return false
}

// SetAuthenticationScheme gets a reference to the given string and assigns it to the AuthenticationScheme field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetAuthenticationScheme(v string) {
	o.AuthenticationScheme = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetEnabled() bool {
	if o == nil || o.Enabled == nil {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetEnabledOk() (*bool, bool) {
	if o == nil || o.Enabled == nil {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasEnabled() bool {
	if o != nil && o.Enabled != nil {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetHttpMethod returns the HttpMethod field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetHttpMethod() string {
	if o == nil || o.HttpMethod == nil {
		var ret string
		return ret
	}
	return *o.HttpMethod
}

// GetHttpMethodOk returns a tuple with the HttpMethod field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetHttpMethodOk() (*string, bool) {
	if o == nil || o.HttpMethod == nil {
		return nil, false
	}
	return o.HttpMethod, true
}

// HasHttpMethod returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasHttpMethod() bool {
	if o != nil && o.HttpMethod != nil {
		return true
	}

	return false
}

// SetHttpMethod gets a reference to the given string and assigns it to the HttpMethod field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetHttpMethod(v string) {
	o.HttpMethod = &v
}

// GetLocalInterface returns the LocalInterface field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetLocalInterface() string {
	if o == nil || o.LocalInterface == nil {
		var ret string
		return ret
	}
	return *o.LocalInterface
}

// GetLocalInterfaceOk returns a tuple with the LocalInterface field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetLocalInterfaceOk() (*string, bool) {
	if o == nil || o.LocalInterface == nil {
		return nil, false
	}
	return o.LocalInterface, true
}

// HasLocalInterface returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasLocalInterface() bool {
	if o != nil && o.LocalInterface != nil {
		return true
	}

	return false
}

// SetLocalInterface gets a reference to the given string and assigns it to the LocalInterface field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetLocalInterface(v string) {
	o.LocalInterface = &v
}

// GetMaxPostWaitTime returns the MaxPostWaitTime field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetMaxPostWaitTime() int32 {
	if o == nil || o.MaxPostWaitTime == nil {
		var ret int32
		return ret
	}
	return *o.MaxPostWaitTime
}

// GetMaxPostWaitTimeOk returns a tuple with the MaxPostWaitTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetMaxPostWaitTimeOk() (*int32, bool) {
	if o == nil || o.MaxPostWaitTime == nil {
		return nil, false
	}
	return o.MaxPostWaitTime, true
}

// HasMaxPostWaitTime returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasMaxPostWaitTime() bool {
	if o != nil && o.MaxPostWaitTime != nil {
		return true
	}

	return false
}

// SetMaxPostWaitTime gets a reference to the given int32 and assigns it to the MaxPostWaitTime field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetMaxPostWaitTime(v int32) {
	o.MaxPostWaitTime = &v
}

// GetMsgVpnName returns the MsgVpnName field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetMsgVpnName() string {
	if o == nil || o.MsgVpnName == nil {
		var ret string
		return ret
	}
	return *o.MsgVpnName
}

// GetMsgVpnNameOk returns a tuple with the MsgVpnName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetMsgVpnNameOk() (*string, bool) {
	if o == nil || o.MsgVpnName == nil {
		return nil, false
	}
	return o.MsgVpnName, true
}

// HasMsgVpnName returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasMsgVpnName() bool {
	if o != nil && o.MsgVpnName != nil {
		return true
	}

	return false
}

// SetMsgVpnName gets a reference to the given string and assigns it to the MsgVpnName field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetMsgVpnName(v string) {
	o.MsgVpnName = &v
}

// GetOutgoingConnectionCount returns the OutgoingConnectionCount field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetOutgoingConnectionCount() int32 {
	if o == nil || o.OutgoingConnectionCount == nil {
		var ret int32
		return ret
	}
	return *o.OutgoingConnectionCount
}

// GetOutgoingConnectionCountOk returns a tuple with the OutgoingConnectionCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetOutgoingConnectionCountOk() (*int32, bool) {
	if o == nil || o.OutgoingConnectionCount == nil {
		return nil, false
	}
	return o.OutgoingConnectionCount, true
}

// HasOutgoingConnectionCount returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasOutgoingConnectionCount() bool {
	if o != nil && o.OutgoingConnectionCount != nil {
		return true
	}

	return false
}

// SetOutgoingConnectionCount gets a reference to the given int32 and assigns it to the OutgoingConnectionCount field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetOutgoingConnectionCount(v int32) {
	o.OutgoingConnectionCount = &v
}

// GetRemoteHost returns the RemoteHost field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetRemoteHost() string {
	if o == nil || o.RemoteHost == nil {
		var ret string
		return ret
	}
	return *o.RemoteHost
}

// GetRemoteHostOk returns a tuple with the RemoteHost field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetRemoteHostOk() (*string, bool) {
	if o == nil || o.RemoteHost == nil {
		return nil, false
	}
	return o.RemoteHost, true
}

// HasRemoteHost returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasRemoteHost() bool {
	if o != nil && o.RemoteHost != nil {
		return true
	}

	return false
}

// SetRemoteHost gets a reference to the given string and assigns it to the RemoteHost field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetRemoteHost(v string) {
	o.RemoteHost = &v
}

// GetRemotePort returns the RemotePort field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetRemotePort() int64 {
	if o == nil || o.RemotePort == nil {
		var ret int64
		return ret
	}
	return *o.RemotePort
}

// GetRemotePortOk returns a tuple with the RemotePort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetRemotePortOk() (*int64, bool) {
	if o == nil || o.RemotePort == nil {
		return nil, false
	}
	return o.RemotePort, true
}

// HasRemotePort returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasRemotePort() bool {
	if o != nil && o.RemotePort != nil {
		return true
	}

	return false
}

// SetRemotePort gets a reference to the given int64 and assigns it to the RemotePort field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetRemotePort(v int64) {
	o.RemotePort = &v
}

// GetRestConsumerName returns the RestConsumerName field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetRestConsumerName() string {
	if o == nil || o.RestConsumerName == nil {
		var ret string
		return ret
	}
	return *o.RestConsumerName
}

// GetRestConsumerNameOk returns a tuple with the RestConsumerName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetRestConsumerNameOk() (*string, bool) {
	if o == nil || o.RestConsumerName == nil {
		return nil, false
	}
	return o.RestConsumerName, true
}

// HasRestConsumerName returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasRestConsumerName() bool {
	if o != nil && o.RestConsumerName != nil {
		return true
	}

	return false
}

// SetRestConsumerName gets a reference to the given string and assigns it to the RestConsumerName field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetRestConsumerName(v string) {
	o.RestConsumerName = &v
}

// GetRestDeliveryPointName returns the RestDeliveryPointName field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetRestDeliveryPointName() string {
	if o == nil || o.RestDeliveryPointName == nil {
		var ret string
		return ret
	}
	return *o.RestDeliveryPointName
}

// GetRestDeliveryPointNameOk returns a tuple with the RestDeliveryPointName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetRestDeliveryPointNameOk() (*string, bool) {
	if o == nil || o.RestDeliveryPointName == nil {
		return nil, false
	}
	return o.RestDeliveryPointName, true
}

// HasRestDeliveryPointName returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasRestDeliveryPointName() bool {
	if o != nil && o.RestDeliveryPointName != nil {
		return true
	}

	return false
}

// SetRestDeliveryPointName gets a reference to the given string and assigns it to the RestDeliveryPointName field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetRestDeliveryPointName(v string) {
	o.RestDeliveryPointName = &v
}

// GetRetryDelay returns the RetryDelay field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetRetryDelay() int32 {
	if o == nil || o.RetryDelay == nil {
		var ret int32
		return ret
	}
	return *o.RetryDelay
}

// GetRetryDelayOk returns a tuple with the RetryDelay field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetRetryDelayOk() (*int32, bool) {
	if o == nil || o.RetryDelay == nil {
		return nil, false
	}
	return o.RetryDelay, true
}

// HasRetryDelay returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasRetryDelay() bool {
	if o != nil && o.RetryDelay != nil {
		return true
	}

	return false
}

// SetRetryDelay gets a reference to the given int32 and assigns it to the RetryDelay field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetRetryDelay(v int32) {
	o.RetryDelay = &v
}

// GetTlsCipherSuiteList returns the TlsCipherSuiteList field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetTlsCipherSuiteList() string {
	if o == nil || o.TlsCipherSuiteList == nil {
		var ret string
		return ret
	}
	return *o.TlsCipherSuiteList
}

// GetTlsCipherSuiteListOk returns a tuple with the TlsCipherSuiteList field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetTlsCipherSuiteListOk() (*string, bool) {
	if o == nil || o.TlsCipherSuiteList == nil {
		return nil, false
	}
	return o.TlsCipherSuiteList, true
}

// HasTlsCipherSuiteList returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasTlsCipherSuiteList() bool {
	if o != nil && o.TlsCipherSuiteList != nil {
		return true
	}

	return false
}

// SetTlsCipherSuiteList gets a reference to the given string and assigns it to the TlsCipherSuiteList field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetTlsCipherSuiteList(v string) {
	o.TlsCipherSuiteList = &v
}

// GetTlsEnabled returns the TlsEnabled field value if set, zero value otherwise.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetTlsEnabled() bool {
	if o == nil || o.TlsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TlsEnabled
}

// GetTlsEnabledOk returns a tuple with the TlsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) GetTlsEnabledOk() (*bool, bool) {
	if o == nil || o.TlsEnabled == nil {
		return nil, false
	}
	return o.TlsEnabled, true
}

// HasTlsEnabled returns a boolean if a field has been set.
func (o *MsgVpnRestDeliveryPointRestConsumer) HasTlsEnabled() bool {
	if o != nil && o.TlsEnabled != nil {
		return true
	}

	return false
}

// SetTlsEnabled gets a reference to the given bool and assigns it to the TlsEnabled field.
func (o *MsgVpnRestDeliveryPointRestConsumer) SetTlsEnabled(v bool) {
	o.TlsEnabled = &v
}

func (o MsgVpnRestDeliveryPointRestConsumer) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AuthenticationAwsAccessKeyId != nil {
		toSerialize["authenticationAwsAccessKeyId"] = o.AuthenticationAwsAccessKeyId
	}
	if o.AuthenticationAwsRegion != nil {
		toSerialize["authenticationAwsRegion"] = o.AuthenticationAwsRegion
	}
	if o.AuthenticationAwsSecretAccessKey != nil {
		toSerialize["authenticationAwsSecretAccessKey"] = o.AuthenticationAwsSecretAccessKey
	}
	if o.AuthenticationAwsService != nil {
		toSerialize["authenticationAwsService"] = o.AuthenticationAwsService
	}
	if o.AuthenticationClientCertContent != nil {
		toSerialize["authenticationClientCertContent"] = o.AuthenticationClientCertContent
	}
	if o.AuthenticationClientCertPassword != nil {
		toSerialize["authenticationClientCertPassword"] = o.AuthenticationClientCertPassword
	}
	if o.AuthenticationHttpBasicPassword != nil {
		toSerialize["authenticationHttpBasicPassword"] = o.AuthenticationHttpBasicPassword
	}
	if o.AuthenticationHttpBasicUsername != nil {
		toSerialize["authenticationHttpBasicUsername"] = o.AuthenticationHttpBasicUsername
	}
	if o.AuthenticationHttpHeaderName != nil {
		toSerialize["authenticationHttpHeaderName"] = o.AuthenticationHttpHeaderName
	}
	if o.AuthenticationHttpHeaderValue != nil {
		toSerialize["authenticationHttpHeaderValue"] = o.AuthenticationHttpHeaderValue
	}
	if o.AuthenticationOauthClientId != nil {
		toSerialize["authenticationOauthClientId"] = o.AuthenticationOauthClientId
	}
	if o.AuthenticationOauthClientScope != nil {
		toSerialize["authenticationOauthClientScope"] = o.AuthenticationOauthClientScope
	}
	if o.AuthenticationOauthClientSecret != nil {
		toSerialize["authenticationOauthClientSecret"] = o.AuthenticationOauthClientSecret
	}
	if o.AuthenticationOauthClientTokenEndpoint != nil {
		toSerialize["authenticationOauthClientTokenEndpoint"] = o.AuthenticationOauthClientTokenEndpoint
	}
	if o.AuthenticationOauthJwtSecretKey != nil {
		toSerialize["authenticationOauthJwtSecretKey"] = o.AuthenticationOauthJwtSecretKey
	}
	if o.AuthenticationOauthJwtTokenEndpoint != nil {
		toSerialize["authenticationOauthJwtTokenEndpoint"] = o.AuthenticationOauthJwtTokenEndpoint
	}
	if o.AuthenticationScheme != nil {
		toSerialize["authenticationScheme"] = o.AuthenticationScheme
	}
	if o.Enabled != nil {
		toSerialize["enabled"] = o.Enabled
	}
	if o.HttpMethod != nil {
		toSerialize["httpMethod"] = o.HttpMethod
	}
	if o.LocalInterface != nil {
		toSerialize["localInterface"] = o.LocalInterface
	}
	if o.MaxPostWaitTime != nil {
		toSerialize["maxPostWaitTime"] = o.MaxPostWaitTime
	}
	if o.MsgVpnName != nil {
		toSerialize["msgVpnName"] = o.MsgVpnName
	}
	if o.OutgoingConnectionCount != nil {
		toSerialize["outgoingConnectionCount"] = o.OutgoingConnectionCount
	}
	if o.RemoteHost != nil {
		toSerialize["remoteHost"] = o.RemoteHost
	}
	if o.RemotePort != nil {
		toSerialize["remotePort"] = o.RemotePort
	}
	if o.RestConsumerName != nil {
		toSerialize["restConsumerName"] = o.RestConsumerName
	}
	if o.RestDeliveryPointName != nil {
		toSerialize["restDeliveryPointName"] = o.RestDeliveryPointName
	}
	if o.RetryDelay != nil {
		toSerialize["retryDelay"] = o.RetryDelay
	}
	if o.TlsCipherSuiteList != nil {
		toSerialize["tlsCipherSuiteList"] = o.TlsCipherSuiteList
	}
	if o.TlsEnabled != nil {
		toSerialize["tlsEnabled"] = o.TlsEnabled
	}
	return json.Marshal(toSerialize)
}

type NullableMsgVpnRestDeliveryPointRestConsumer struct {
	value *MsgVpnRestDeliveryPointRestConsumer
	isSet bool
}

func (v NullableMsgVpnRestDeliveryPointRestConsumer) Get() *MsgVpnRestDeliveryPointRestConsumer {
	return v.value
}

func (v *NullableMsgVpnRestDeliveryPointRestConsumer) Set(val *MsgVpnRestDeliveryPointRestConsumer) {
	v.value = val
	v.isSet = true
}

func (v NullableMsgVpnRestDeliveryPointRestConsumer) IsSet() bool {
	return v.isSet
}

func (v *NullableMsgVpnRestDeliveryPointRestConsumer) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMsgVpnRestDeliveryPointRestConsumer(val *MsgVpnRestDeliveryPointRestConsumer) *NullableMsgVpnRestDeliveryPointRestConsumer {
	return &NullableMsgVpnRestDeliveryPointRestConsumer{value: val, isSet: true}
}

func (v NullableMsgVpnRestDeliveryPointRestConsumer) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMsgVpnRestDeliveryPointRestConsumer) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


