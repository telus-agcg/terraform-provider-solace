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

// DmrCluster struct for DmrCluster
type DmrCluster struct {
	// Enable or disable basic authentication for Cluster Links. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`.
	AuthenticationBasicEnabled *bool `json:"authenticationBasicEnabled,omitempty"`
	// The password used to authenticate incoming Cluster Links when using basic internal authentication. The same password is also used by outgoing Cluster Links if a per-Link password is not configured. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	AuthenticationBasicPassword *string `json:"authenticationBasicPassword,omitempty"`
	// The type of basic authentication to use for Cluster Links. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"internal\"`. The allowed values and their meaning are:  <pre> \"internal\" - Use locally configured password. \"none\" - No authentication. </pre> 
	AuthenticationBasicType *string `json:"authenticationBasicType,omitempty"`
	// The PEM formatted content for the client certificate used to login to the remote node. It must consist of a private key and between one and three certificates comprising the certificate trust chain. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changing this attribute requires an HTTPS connection. The default value is `\"\"`.
	AuthenticationClientCertContent *string `json:"authenticationClientCertContent,omitempty"`
	// Enable or disable client certificate authentication for Cluster Links. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`.
	AuthenticationClientCertEnabled *bool `json:"authenticationClientCertEnabled,omitempty"`
	// The password for the client certificate. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changing this attribute requires an HTTPS connection. The default value is `\"\"`.
	AuthenticationClientCertPassword *string `json:"authenticationClientCertPassword,omitempty"`
	// Enable or disable direct messaging only. Guaranteed messages will not be transmitted through the cluster. The default value is `false`.
	DirectOnlyEnabled *bool `json:"directOnlyEnabled,omitempty"`
	// The name of the Cluster.
	DmrClusterName *string `json:"dmrClusterName,omitempty"`
	// Enable or disable the Cluster. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.
	Enabled *bool `json:"enabled,omitempty"`
	// The name of this node in the Cluster. This is the name that this broker (or redundant group of brokers) is know by to other nodes in the Cluster. The name is chosen automatically to be either this broker's Router Name or Mate Router Name, depending on which Active Standby Role (primary or backup) this broker plays in its redundancy group.
	NodeName *string `json:"nodeName,omitempty"`
	// Enable or disable the enforcing of the common name provided by the remote broker against the list of trusted common names configured for the Link. If enabled, the certificate's common name must match one of the trusted common names for the Link to be accepted. Common Name validation is not performed if Server Certificate Name Validation is enabled, even if Common Name validation is enabled. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`. Deprecated since 2.18. Common Name validation has been replaced by Server Certificate Name validation.
	TlsServerCertEnforceTrustedCommonNameEnabled *bool `json:"tlsServerCertEnforceTrustedCommonNameEnabled,omitempty"`
	// The maximum allowed depth of a certificate chain. The depth of a chain is defined as the number of signing CA certificates that are present in the chain back to a trusted self-signed root CA certificate. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `3`.
	TlsServerCertMaxChainDepth *int64 `json:"tlsServerCertMaxChainDepth,omitempty"`
	// Enable or disable the validation of the \"Not Before\" and \"Not After\" validity dates in the certificate. When disabled, the certificate is accepted even if the certificate is not valid based on these dates. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`.
	TlsServerCertValidateDateEnabled *bool `json:"tlsServerCertValidateDateEnabled,omitempty"`
	// Enable or disable the standard TLS authentication mechanism of verifying the name used to connect to the bridge. If enabled, the name used to connect to the bridge is checked against the names specified in the certificate returned by the remote router. Legacy Common Name validation is not performed if Server Certificate Name Validation is enabled, even if Common Name validation is also enabled. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `true`. Available since 2.18.
	TlsServerCertValidateNameEnabled *bool `json:"tlsServerCertValidateNameEnabled,omitempty"`
}

// NewDmrCluster instantiates a new DmrCluster object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDmrCluster() *DmrCluster {
	this := DmrCluster{}
	return &this
}

// NewDmrClusterWithDefaults instantiates a new DmrCluster object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDmrClusterWithDefaults() *DmrCluster {
	this := DmrCluster{}
	return &this
}

// GetAuthenticationBasicEnabled returns the AuthenticationBasicEnabled field value if set, zero value otherwise.
func (o *DmrCluster) GetAuthenticationBasicEnabled() bool {
	if o == nil || o.AuthenticationBasicEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AuthenticationBasicEnabled
}

// GetAuthenticationBasicEnabledOk returns a tuple with the AuthenticationBasicEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrCluster) GetAuthenticationBasicEnabledOk() (*bool, bool) {
	if o == nil || o.AuthenticationBasicEnabled == nil {
		return nil, false
	}
	return o.AuthenticationBasicEnabled, true
}

// HasAuthenticationBasicEnabled returns a boolean if a field has been set.
func (o *DmrCluster) HasAuthenticationBasicEnabled() bool {
	if o != nil && o.AuthenticationBasicEnabled != nil {
		return true
	}

	return false
}

// SetAuthenticationBasicEnabled gets a reference to the given bool and assigns it to the AuthenticationBasicEnabled field.
func (o *DmrCluster) SetAuthenticationBasicEnabled(v bool) {
	o.AuthenticationBasicEnabled = &v
}

// GetAuthenticationBasicPassword returns the AuthenticationBasicPassword field value if set, zero value otherwise.
func (o *DmrCluster) GetAuthenticationBasicPassword() string {
	if o == nil || o.AuthenticationBasicPassword == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationBasicPassword
}

// GetAuthenticationBasicPasswordOk returns a tuple with the AuthenticationBasicPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrCluster) GetAuthenticationBasicPasswordOk() (*string, bool) {
	if o == nil || o.AuthenticationBasicPassword == nil {
		return nil, false
	}
	return o.AuthenticationBasicPassword, true
}

// HasAuthenticationBasicPassword returns a boolean if a field has been set.
func (o *DmrCluster) HasAuthenticationBasicPassword() bool {
	if o != nil && o.AuthenticationBasicPassword != nil {
		return true
	}

	return false
}

// SetAuthenticationBasicPassword gets a reference to the given string and assigns it to the AuthenticationBasicPassword field.
func (o *DmrCluster) SetAuthenticationBasicPassword(v string) {
	o.AuthenticationBasicPassword = &v
}

// GetAuthenticationBasicType returns the AuthenticationBasicType field value if set, zero value otherwise.
func (o *DmrCluster) GetAuthenticationBasicType() string {
	if o == nil || o.AuthenticationBasicType == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationBasicType
}

// GetAuthenticationBasicTypeOk returns a tuple with the AuthenticationBasicType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrCluster) GetAuthenticationBasicTypeOk() (*string, bool) {
	if o == nil || o.AuthenticationBasicType == nil {
		return nil, false
	}
	return o.AuthenticationBasicType, true
}

// HasAuthenticationBasicType returns a boolean if a field has been set.
func (o *DmrCluster) HasAuthenticationBasicType() bool {
	if o != nil && o.AuthenticationBasicType != nil {
		return true
	}

	return false
}

// SetAuthenticationBasicType gets a reference to the given string and assigns it to the AuthenticationBasicType field.
func (o *DmrCluster) SetAuthenticationBasicType(v string) {
	o.AuthenticationBasicType = &v
}

// GetAuthenticationClientCertContent returns the AuthenticationClientCertContent field value if set, zero value otherwise.
func (o *DmrCluster) GetAuthenticationClientCertContent() string {
	if o == nil || o.AuthenticationClientCertContent == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationClientCertContent
}

// GetAuthenticationClientCertContentOk returns a tuple with the AuthenticationClientCertContent field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrCluster) GetAuthenticationClientCertContentOk() (*string, bool) {
	if o == nil || o.AuthenticationClientCertContent == nil {
		return nil, false
	}
	return o.AuthenticationClientCertContent, true
}

// HasAuthenticationClientCertContent returns a boolean if a field has been set.
func (o *DmrCluster) HasAuthenticationClientCertContent() bool {
	if o != nil && o.AuthenticationClientCertContent != nil {
		return true
	}

	return false
}

// SetAuthenticationClientCertContent gets a reference to the given string and assigns it to the AuthenticationClientCertContent field.
func (o *DmrCluster) SetAuthenticationClientCertContent(v string) {
	o.AuthenticationClientCertContent = &v
}

// GetAuthenticationClientCertEnabled returns the AuthenticationClientCertEnabled field value if set, zero value otherwise.
func (o *DmrCluster) GetAuthenticationClientCertEnabled() bool {
	if o == nil || o.AuthenticationClientCertEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AuthenticationClientCertEnabled
}

// GetAuthenticationClientCertEnabledOk returns a tuple with the AuthenticationClientCertEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrCluster) GetAuthenticationClientCertEnabledOk() (*bool, bool) {
	if o == nil || o.AuthenticationClientCertEnabled == nil {
		return nil, false
	}
	return o.AuthenticationClientCertEnabled, true
}

// HasAuthenticationClientCertEnabled returns a boolean if a field has been set.
func (o *DmrCluster) HasAuthenticationClientCertEnabled() bool {
	if o != nil && o.AuthenticationClientCertEnabled != nil {
		return true
	}

	return false
}

// SetAuthenticationClientCertEnabled gets a reference to the given bool and assigns it to the AuthenticationClientCertEnabled field.
func (o *DmrCluster) SetAuthenticationClientCertEnabled(v bool) {
	o.AuthenticationClientCertEnabled = &v
}

// GetAuthenticationClientCertPassword returns the AuthenticationClientCertPassword field value if set, zero value otherwise.
func (o *DmrCluster) GetAuthenticationClientCertPassword() string {
	if o == nil || o.AuthenticationClientCertPassword == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationClientCertPassword
}

// GetAuthenticationClientCertPasswordOk returns a tuple with the AuthenticationClientCertPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrCluster) GetAuthenticationClientCertPasswordOk() (*string, bool) {
	if o == nil || o.AuthenticationClientCertPassword == nil {
		return nil, false
	}
	return o.AuthenticationClientCertPassword, true
}

// HasAuthenticationClientCertPassword returns a boolean if a field has been set.
func (o *DmrCluster) HasAuthenticationClientCertPassword() bool {
	if o != nil && o.AuthenticationClientCertPassword != nil {
		return true
	}

	return false
}

// SetAuthenticationClientCertPassword gets a reference to the given string and assigns it to the AuthenticationClientCertPassword field.
func (o *DmrCluster) SetAuthenticationClientCertPassword(v string) {
	o.AuthenticationClientCertPassword = &v
}

// GetDirectOnlyEnabled returns the DirectOnlyEnabled field value if set, zero value otherwise.
func (o *DmrCluster) GetDirectOnlyEnabled() bool {
	if o == nil || o.DirectOnlyEnabled == nil {
		var ret bool
		return ret
	}
	return *o.DirectOnlyEnabled
}

// GetDirectOnlyEnabledOk returns a tuple with the DirectOnlyEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrCluster) GetDirectOnlyEnabledOk() (*bool, bool) {
	if o == nil || o.DirectOnlyEnabled == nil {
		return nil, false
	}
	return o.DirectOnlyEnabled, true
}

// HasDirectOnlyEnabled returns a boolean if a field has been set.
func (o *DmrCluster) HasDirectOnlyEnabled() bool {
	if o != nil && o.DirectOnlyEnabled != nil {
		return true
	}

	return false
}

// SetDirectOnlyEnabled gets a reference to the given bool and assigns it to the DirectOnlyEnabled field.
func (o *DmrCluster) SetDirectOnlyEnabled(v bool) {
	o.DirectOnlyEnabled = &v
}

// GetDmrClusterName returns the DmrClusterName field value if set, zero value otherwise.
func (o *DmrCluster) GetDmrClusterName() string {
	if o == nil || o.DmrClusterName == nil {
		var ret string
		return ret
	}
	return *o.DmrClusterName
}

// GetDmrClusterNameOk returns a tuple with the DmrClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrCluster) GetDmrClusterNameOk() (*string, bool) {
	if o == nil || o.DmrClusterName == nil {
		return nil, false
	}
	return o.DmrClusterName, true
}

// HasDmrClusterName returns a boolean if a field has been set.
func (o *DmrCluster) HasDmrClusterName() bool {
	if o != nil && o.DmrClusterName != nil {
		return true
	}

	return false
}

// SetDmrClusterName gets a reference to the given string and assigns it to the DmrClusterName field.
func (o *DmrCluster) SetDmrClusterName(v string) {
	o.DmrClusterName = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *DmrCluster) GetEnabled() bool {
	if o == nil || o.Enabled == nil {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrCluster) GetEnabledOk() (*bool, bool) {
	if o == nil || o.Enabled == nil {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *DmrCluster) HasEnabled() bool {
	if o != nil && o.Enabled != nil {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *DmrCluster) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetNodeName returns the NodeName field value if set, zero value otherwise.
func (o *DmrCluster) GetNodeName() string {
	if o == nil || o.NodeName == nil {
		var ret string
		return ret
	}
	return *o.NodeName
}

// GetNodeNameOk returns a tuple with the NodeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrCluster) GetNodeNameOk() (*string, bool) {
	if o == nil || o.NodeName == nil {
		return nil, false
	}
	return o.NodeName, true
}

// HasNodeName returns a boolean if a field has been set.
func (o *DmrCluster) HasNodeName() bool {
	if o != nil && o.NodeName != nil {
		return true
	}

	return false
}

// SetNodeName gets a reference to the given string and assigns it to the NodeName field.
func (o *DmrCluster) SetNodeName(v string) {
	o.NodeName = &v
}

// GetTlsServerCertEnforceTrustedCommonNameEnabled returns the TlsServerCertEnforceTrustedCommonNameEnabled field value if set, zero value otherwise.
func (o *DmrCluster) GetTlsServerCertEnforceTrustedCommonNameEnabled() bool {
	if o == nil || o.TlsServerCertEnforceTrustedCommonNameEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TlsServerCertEnforceTrustedCommonNameEnabled
}

// GetTlsServerCertEnforceTrustedCommonNameEnabledOk returns a tuple with the TlsServerCertEnforceTrustedCommonNameEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrCluster) GetTlsServerCertEnforceTrustedCommonNameEnabledOk() (*bool, bool) {
	if o == nil || o.TlsServerCertEnforceTrustedCommonNameEnabled == nil {
		return nil, false
	}
	return o.TlsServerCertEnforceTrustedCommonNameEnabled, true
}

// HasTlsServerCertEnforceTrustedCommonNameEnabled returns a boolean if a field has been set.
func (o *DmrCluster) HasTlsServerCertEnforceTrustedCommonNameEnabled() bool {
	if o != nil && o.TlsServerCertEnforceTrustedCommonNameEnabled != nil {
		return true
	}

	return false
}

// SetTlsServerCertEnforceTrustedCommonNameEnabled gets a reference to the given bool and assigns it to the TlsServerCertEnforceTrustedCommonNameEnabled field.
func (o *DmrCluster) SetTlsServerCertEnforceTrustedCommonNameEnabled(v bool) {
	o.TlsServerCertEnforceTrustedCommonNameEnabled = &v
}

// GetTlsServerCertMaxChainDepth returns the TlsServerCertMaxChainDepth field value if set, zero value otherwise.
func (o *DmrCluster) GetTlsServerCertMaxChainDepth() int64 {
	if o == nil || o.TlsServerCertMaxChainDepth == nil {
		var ret int64
		return ret
	}
	return *o.TlsServerCertMaxChainDepth
}

// GetTlsServerCertMaxChainDepthOk returns a tuple with the TlsServerCertMaxChainDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrCluster) GetTlsServerCertMaxChainDepthOk() (*int64, bool) {
	if o == nil || o.TlsServerCertMaxChainDepth == nil {
		return nil, false
	}
	return o.TlsServerCertMaxChainDepth, true
}

// HasTlsServerCertMaxChainDepth returns a boolean if a field has been set.
func (o *DmrCluster) HasTlsServerCertMaxChainDepth() bool {
	if o != nil && o.TlsServerCertMaxChainDepth != nil {
		return true
	}

	return false
}

// SetTlsServerCertMaxChainDepth gets a reference to the given int64 and assigns it to the TlsServerCertMaxChainDepth field.
func (o *DmrCluster) SetTlsServerCertMaxChainDepth(v int64) {
	o.TlsServerCertMaxChainDepth = &v
}

// GetTlsServerCertValidateDateEnabled returns the TlsServerCertValidateDateEnabled field value if set, zero value otherwise.
func (o *DmrCluster) GetTlsServerCertValidateDateEnabled() bool {
	if o == nil || o.TlsServerCertValidateDateEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TlsServerCertValidateDateEnabled
}

// GetTlsServerCertValidateDateEnabledOk returns a tuple with the TlsServerCertValidateDateEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrCluster) GetTlsServerCertValidateDateEnabledOk() (*bool, bool) {
	if o == nil || o.TlsServerCertValidateDateEnabled == nil {
		return nil, false
	}
	return o.TlsServerCertValidateDateEnabled, true
}

// HasTlsServerCertValidateDateEnabled returns a boolean if a field has been set.
func (o *DmrCluster) HasTlsServerCertValidateDateEnabled() bool {
	if o != nil && o.TlsServerCertValidateDateEnabled != nil {
		return true
	}

	return false
}

// SetTlsServerCertValidateDateEnabled gets a reference to the given bool and assigns it to the TlsServerCertValidateDateEnabled field.
func (o *DmrCluster) SetTlsServerCertValidateDateEnabled(v bool) {
	o.TlsServerCertValidateDateEnabled = &v
}

// GetTlsServerCertValidateNameEnabled returns the TlsServerCertValidateNameEnabled field value if set, zero value otherwise.
func (o *DmrCluster) GetTlsServerCertValidateNameEnabled() bool {
	if o == nil || o.TlsServerCertValidateNameEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TlsServerCertValidateNameEnabled
}

// GetTlsServerCertValidateNameEnabledOk returns a tuple with the TlsServerCertValidateNameEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DmrCluster) GetTlsServerCertValidateNameEnabledOk() (*bool, bool) {
	if o == nil || o.TlsServerCertValidateNameEnabled == nil {
		return nil, false
	}
	return o.TlsServerCertValidateNameEnabled, true
}

// HasTlsServerCertValidateNameEnabled returns a boolean if a field has been set.
func (o *DmrCluster) HasTlsServerCertValidateNameEnabled() bool {
	if o != nil && o.TlsServerCertValidateNameEnabled != nil {
		return true
	}

	return false
}

// SetTlsServerCertValidateNameEnabled gets a reference to the given bool and assigns it to the TlsServerCertValidateNameEnabled field.
func (o *DmrCluster) SetTlsServerCertValidateNameEnabled(v bool) {
	o.TlsServerCertValidateNameEnabled = &v
}

func (o DmrCluster) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AuthenticationBasicEnabled != nil {
		toSerialize["authenticationBasicEnabled"] = o.AuthenticationBasicEnabled
	}
	if o.AuthenticationBasicPassword != nil {
		toSerialize["authenticationBasicPassword"] = o.AuthenticationBasicPassword
	}
	if o.AuthenticationBasicType != nil {
		toSerialize["authenticationBasicType"] = o.AuthenticationBasicType
	}
	if o.AuthenticationClientCertContent != nil {
		toSerialize["authenticationClientCertContent"] = o.AuthenticationClientCertContent
	}
	if o.AuthenticationClientCertEnabled != nil {
		toSerialize["authenticationClientCertEnabled"] = o.AuthenticationClientCertEnabled
	}
	if o.AuthenticationClientCertPassword != nil {
		toSerialize["authenticationClientCertPassword"] = o.AuthenticationClientCertPassword
	}
	if o.DirectOnlyEnabled != nil {
		toSerialize["directOnlyEnabled"] = o.DirectOnlyEnabled
	}
	if o.DmrClusterName != nil {
		toSerialize["dmrClusterName"] = o.DmrClusterName
	}
	if o.Enabled != nil {
		toSerialize["enabled"] = o.Enabled
	}
	if o.NodeName != nil {
		toSerialize["nodeName"] = o.NodeName
	}
	if o.TlsServerCertEnforceTrustedCommonNameEnabled != nil {
		toSerialize["tlsServerCertEnforceTrustedCommonNameEnabled"] = o.TlsServerCertEnforceTrustedCommonNameEnabled
	}
	if o.TlsServerCertMaxChainDepth != nil {
		toSerialize["tlsServerCertMaxChainDepth"] = o.TlsServerCertMaxChainDepth
	}
	if o.TlsServerCertValidateDateEnabled != nil {
		toSerialize["tlsServerCertValidateDateEnabled"] = o.TlsServerCertValidateDateEnabled
	}
	if o.TlsServerCertValidateNameEnabled != nil {
		toSerialize["tlsServerCertValidateNameEnabled"] = o.TlsServerCertValidateNameEnabled
	}
	return json.Marshal(toSerialize)
}

type NullableDmrCluster struct {
	value *DmrCluster
	isSet bool
}

func (v NullableDmrCluster) Get() *DmrCluster {
	return v.value
}

func (v *NullableDmrCluster) Set(val *DmrCluster) {
	v.value = val
	v.isSet = true
}

func (v NullableDmrCluster) IsSet() bool {
	return v.isSet
}

func (v *NullableDmrCluster) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDmrCluster(val *DmrCluster) *NullableDmrCluster {
	return &NullableDmrCluster{value: val, isSet: true}
}

func (v NullableDmrCluster) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDmrCluster) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


