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

// MsgVpnJndiConnectionFactory struct for MsgVpnJndiConnectionFactory
type MsgVpnJndiConnectionFactory struct {
	// Enable or disable whether new JMS connections can use the same Client identifier (ID) as an existing connection. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.3.
	AllowDuplicateClientIdEnabled *bool `json:"allowDuplicateClientIdEnabled,omitempty"`
	// The description of the Client. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.
	ClientDescription *string `json:"clientDescription,omitempty"`
	// The Client identifier (ID). If not specified, a unique value for it will be generated. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.
	ClientId *string `json:"clientId,omitempty"`
	// The name of the JMS Connection Factory.
	ConnectionFactoryName *string `json:"connectionFactoryName,omitempty"`
	// Enable or disable overriding by the Subscriber (Consumer) of the deliver-to-one (DTO) property on messages. When enabled, the Subscriber can receive all DTO tagged messages. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	DtoReceiveOverrideEnabled *bool `json:"dtoReceiveOverrideEnabled,omitempty"`
	// The priority for receiving deliver-to-one (DTO) messages by the Subscriber (Consumer) if the messages are published on the local broker that the Subscriber is directly connected to. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1`.
	DtoReceiveSubscriberLocalPriority *int32 `json:"dtoReceiveSubscriberLocalPriority,omitempty"`
	// The priority for receiving deliver-to-one (DTO) messages by the Subscriber (Consumer) if the messages are published on a remote broker. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1`.
	DtoReceiveSubscriberNetworkPriority *int32 `json:"dtoReceiveSubscriberNetworkPriority,omitempty"`
	// Enable or disable the deliver-to-one (DTO) property on messages sent by the Publisher (Producer). Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	DtoSendEnabled *bool `json:"dtoSendEnabled,omitempty"`
	// Enable or disable whether a durable endpoint will be dynamically created on the broker when the client calls \"Session.createDurableSubscriber()\" or \"Session.createQueue()\". The created endpoint respects the message time-to-live (TTL) according to the \"dynamicEndpointRespectTtlEnabled\" property. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	DynamicEndpointCreateDurableEnabled *bool `json:"dynamicEndpointCreateDurableEnabled,omitempty"`
	// Enable or disable whether dynamically created durable and non-durable endpoints respect the message time-to-live (TTL) property. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	DynamicEndpointRespectTtlEnabled *bool `json:"dynamicEndpointRespectTtlEnabled,omitempty"`
	// The timeout for sending the acknowledgement (ACK) for guaranteed messages received by the Subscriber (Consumer), in milliseconds. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.
	GuaranteedReceiveAckTimeout *int32 `json:"guaranteedReceiveAckTimeout,omitempty"`
	// The maximum number of attempts to reconnect to the host or list of hosts after the guaranteed  messaging connection has been lost. The value \"-1\" means to retry forever. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `-1`. Available since 2.14.
	GuaranteedReceiveReconnectRetryCount *int32 `json:"guaranteedReceiveReconnectRetryCount,omitempty"`
	// The amount of time to wait before making another attempt to connect or reconnect to the host after the guaranteed messaging connection has been lost, in milliseconds. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3000`. Available since 2.14.
	GuaranteedReceiveReconnectRetryWait *int32 `json:"guaranteedReceiveReconnectRetryWait,omitempty"`
	// The size of the window for guaranteed messages received by the Subscriber (Consumer), in messages. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `18`.
	GuaranteedReceiveWindowSize *int32 `json:"guaranteedReceiveWindowSize,omitempty"`
	// The threshold for sending the acknowledgement (ACK) for guaranteed messages received by the Subscriber (Consumer) as a percentage of `guaranteedReceiveWindowSize`. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `60`.
	GuaranteedReceiveWindowSizeAckThreshold *int32 `json:"guaranteedReceiveWindowSizeAckThreshold,omitempty"`
	// The timeout for receiving the acknowledgement (ACK) for guaranteed messages sent by the Publisher (Producer), in milliseconds. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `2000`.
	GuaranteedSendAckTimeout *int32 `json:"guaranteedSendAckTimeout,omitempty"`
	// The size of the window for non-persistent guaranteed messages sent by the Publisher (Producer), in messages. For persistent messages the window size is fixed at 1. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `255`.
	GuaranteedSendWindowSize *int32 `json:"guaranteedSendWindowSize,omitempty"`
	// The default delivery mode for messages sent by the Publisher (Producer). Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"persistent\"`. The allowed values and their meaning are:  <pre> \"persistent\" - The broker spools messages (persists in the Message Spool) as part of the send operation. \"non-persistent\" - The broker does not spool messages (does not persist in the Message Spool) as part of the send operation. </pre> 
	MessagingDefaultDeliveryMode *string `json:"messagingDefaultDeliveryMode,omitempty"`
	// Enable or disable whether messages sent by the Publisher (Producer) are Dead Message Queue (DMQ) eligible by default. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	MessagingDefaultDmqEligibleEnabled *bool `json:"messagingDefaultDmqEligibleEnabled,omitempty"`
	// Enable or disable whether messages sent by the Publisher (Producer) are Eliding eligible by default. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	MessagingDefaultElidingEligibleEnabled *bool `json:"messagingDefaultElidingEligibleEnabled,omitempty"`
	// Enable or disable inclusion (adding or replacing) of the JMSXUserID property in messages sent by the Publisher (Producer). Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	MessagingJmsxUserIdEnabled *bool `json:"messagingJmsxUserIdEnabled,omitempty"`
	// Enable or disable encoding of JMS text messages in Publisher (Producer) messages as XML payload. When disabled, JMS text messages are encoded as a binary attachment. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	MessagingTextInXmlPayloadEnabled *bool `json:"messagingTextInXmlPayloadEnabled,omitempty"`
	// The name of the Message VPN.
	MsgVpnName *string `json:"msgVpnName,omitempty"`
	// The ZLIB compression level for the connection to the broker. The value \"0\" means no compression, and the value \"-1\" means the compression level is specified in the JNDI Properties file. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `-1`.
	TransportCompressionLevel *int32 `json:"transportCompressionLevel,omitempty"`
	// The maximum number of retry attempts to establish an initial connection to the host or list of hosts. The value \"0\" means a single attempt (no retries), and the value \"-1\" means to retry forever. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.
	TransportConnectRetryCount *int32 `json:"transportConnectRetryCount,omitempty"`
	// The maximum number of retry attempts to establish an initial connection to each host on the list of hosts. The value \"0\" means a single attempt (no retries), and the value \"-1\" means to retry forever. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.
	TransportConnectRetryPerHostCount *int32 `json:"transportConnectRetryPerHostCount,omitempty"`
	// The timeout for establishing an initial connection to the broker, in milliseconds. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `30000`.
	TransportConnectTimeout *int32 `json:"transportConnectTimeout,omitempty"`
	// Enable or disable usage of the Direct Transport mode for sending non-persistent messages. When disabled, the Guaranteed Transport mode is used. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	TransportDirectTransportEnabled *bool `json:"transportDirectTransportEnabled,omitempty"`
	// The maximum number of consecutive application-level keepalive messages sent without the broker response before the connection to the broker is closed. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.
	TransportKeepaliveCount *int32 `json:"transportKeepaliveCount,omitempty"`
	// Enable or disable usage of application-level keepalive messages to maintain a connection with the broker. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	TransportKeepaliveEnabled *bool `json:"transportKeepaliveEnabled,omitempty"`
	// The interval between application-level keepalive messages, in milliseconds. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3000`.
	TransportKeepaliveInterval *int32 `json:"transportKeepaliveInterval,omitempty"`
	// Enable or disable delivery of asynchronous messages directly from the I/O thread. Contact support before enabling this property. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	TransportMsgCallbackOnIoThreadEnabled *bool `json:"transportMsgCallbackOnIoThreadEnabled,omitempty"`
	// Enable or disable optimization for the Direct Transport delivery mode. If enabled, the client application is limited to one Publisher (Producer) and one non-durable Subscriber (Consumer). Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	TransportOptimizeDirectEnabled *bool `json:"transportOptimizeDirectEnabled,omitempty"`
	// The connection port number on the broker for SMF clients. The value \"-1\" means the port is specified in the JNDI Properties file. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `-1`.
	TransportPort *int32 `json:"transportPort,omitempty"`
	// The timeout for reading a reply from the broker, in milliseconds. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `10000`.
	TransportReadTimeout *int32 `json:"transportReadTimeout,omitempty"`
	// The size of the receive socket buffer, in bytes. It corresponds to the SO_RCVBUF socket option. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `65536`.
	TransportReceiveBufferSize *int32 `json:"transportReceiveBufferSize,omitempty"`
	// The maximum number of attempts to reconnect to the host or list of hosts after the connection has been lost. The value \"-1\" means to retry forever. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.
	TransportReconnectRetryCount *int32 `json:"transportReconnectRetryCount,omitempty"`
	// The amount of time before making another attempt to connect or reconnect to the host after the connection has been lost, in milliseconds. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3000`.
	TransportReconnectRetryWait *int32 `json:"transportReconnectRetryWait,omitempty"`
	// The size of the send socket buffer, in bytes. It corresponds to the SO_SNDBUF socket option. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `65536`.
	TransportSendBufferSize *int32 `json:"transportSendBufferSize,omitempty"`
	// Enable or disable the TCP_NODELAY option. When enabled, Nagle's algorithm for TCP/IP congestion control (RFC 896) is disabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	TransportTcpNoDelayEnabled *bool `json:"transportTcpNoDelayEnabled,omitempty"`
	// Enable or disable this as an XA Connection Factory. When enabled, the Connection Factory can be cast to \"XAConnectionFactory\", \"XAQueueConnectionFactory\" or \"XATopicConnectionFactory\". Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	XaEnabled *bool `json:"xaEnabled,omitempty"`
}

// NewMsgVpnJndiConnectionFactory instantiates a new MsgVpnJndiConnectionFactory object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMsgVpnJndiConnectionFactory() *MsgVpnJndiConnectionFactory {
	this := MsgVpnJndiConnectionFactory{}
	return &this
}

// NewMsgVpnJndiConnectionFactoryWithDefaults instantiates a new MsgVpnJndiConnectionFactory object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMsgVpnJndiConnectionFactoryWithDefaults() *MsgVpnJndiConnectionFactory {
	this := MsgVpnJndiConnectionFactory{}
	return &this
}

// GetAllowDuplicateClientIdEnabled returns the AllowDuplicateClientIdEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetAllowDuplicateClientIdEnabled() bool {
	if o == nil || o.AllowDuplicateClientIdEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AllowDuplicateClientIdEnabled
}

// GetAllowDuplicateClientIdEnabledOk returns a tuple with the AllowDuplicateClientIdEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetAllowDuplicateClientIdEnabledOk() (*bool, bool) {
	if o == nil || o.AllowDuplicateClientIdEnabled == nil {
		return nil, false
	}
	return o.AllowDuplicateClientIdEnabled, true
}

// HasAllowDuplicateClientIdEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasAllowDuplicateClientIdEnabled() bool {
	if o != nil && o.AllowDuplicateClientIdEnabled != nil {
		return true
	}

	return false
}

// SetAllowDuplicateClientIdEnabled gets a reference to the given bool and assigns it to the AllowDuplicateClientIdEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetAllowDuplicateClientIdEnabled(v bool) {
	o.AllowDuplicateClientIdEnabled = &v
}

// GetClientDescription returns the ClientDescription field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetClientDescription() string {
	if o == nil || o.ClientDescription == nil {
		var ret string
		return ret
	}
	return *o.ClientDescription
}

// GetClientDescriptionOk returns a tuple with the ClientDescription field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetClientDescriptionOk() (*string, bool) {
	if o == nil || o.ClientDescription == nil {
		return nil, false
	}
	return o.ClientDescription, true
}

// HasClientDescription returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasClientDescription() bool {
	if o != nil && o.ClientDescription != nil {
		return true
	}

	return false
}

// SetClientDescription gets a reference to the given string and assigns it to the ClientDescription field.
func (o *MsgVpnJndiConnectionFactory) SetClientDescription(v string) {
	o.ClientDescription = &v
}

// GetClientId returns the ClientId field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetClientId() string {
	if o == nil || o.ClientId == nil {
		var ret string
		return ret
	}
	return *o.ClientId
}

// GetClientIdOk returns a tuple with the ClientId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetClientIdOk() (*string, bool) {
	if o == nil || o.ClientId == nil {
		return nil, false
	}
	return o.ClientId, true
}

// HasClientId returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasClientId() bool {
	if o != nil && o.ClientId != nil {
		return true
	}

	return false
}

// SetClientId gets a reference to the given string and assigns it to the ClientId field.
func (o *MsgVpnJndiConnectionFactory) SetClientId(v string) {
	o.ClientId = &v
}

// GetConnectionFactoryName returns the ConnectionFactoryName field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetConnectionFactoryName() string {
	if o == nil || o.ConnectionFactoryName == nil {
		var ret string
		return ret
	}
	return *o.ConnectionFactoryName
}

// GetConnectionFactoryNameOk returns a tuple with the ConnectionFactoryName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetConnectionFactoryNameOk() (*string, bool) {
	if o == nil || o.ConnectionFactoryName == nil {
		return nil, false
	}
	return o.ConnectionFactoryName, true
}

// HasConnectionFactoryName returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasConnectionFactoryName() bool {
	if o != nil && o.ConnectionFactoryName != nil {
		return true
	}

	return false
}

// SetConnectionFactoryName gets a reference to the given string and assigns it to the ConnectionFactoryName field.
func (o *MsgVpnJndiConnectionFactory) SetConnectionFactoryName(v string) {
	o.ConnectionFactoryName = &v
}

// GetDtoReceiveOverrideEnabled returns the DtoReceiveOverrideEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetDtoReceiveOverrideEnabled() bool {
	if o == nil || o.DtoReceiveOverrideEnabled == nil {
		var ret bool
		return ret
	}
	return *o.DtoReceiveOverrideEnabled
}

// GetDtoReceiveOverrideEnabledOk returns a tuple with the DtoReceiveOverrideEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetDtoReceiveOverrideEnabledOk() (*bool, bool) {
	if o == nil || o.DtoReceiveOverrideEnabled == nil {
		return nil, false
	}
	return o.DtoReceiveOverrideEnabled, true
}

// HasDtoReceiveOverrideEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasDtoReceiveOverrideEnabled() bool {
	if o != nil && o.DtoReceiveOverrideEnabled != nil {
		return true
	}

	return false
}

// SetDtoReceiveOverrideEnabled gets a reference to the given bool and assigns it to the DtoReceiveOverrideEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetDtoReceiveOverrideEnabled(v bool) {
	o.DtoReceiveOverrideEnabled = &v
}

// GetDtoReceiveSubscriberLocalPriority returns the DtoReceiveSubscriberLocalPriority field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetDtoReceiveSubscriberLocalPriority() int32 {
	if o == nil || o.DtoReceiveSubscriberLocalPriority == nil {
		var ret int32
		return ret
	}
	return *o.DtoReceiveSubscriberLocalPriority
}

// GetDtoReceiveSubscriberLocalPriorityOk returns a tuple with the DtoReceiveSubscriberLocalPriority field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetDtoReceiveSubscriberLocalPriorityOk() (*int32, bool) {
	if o == nil || o.DtoReceiveSubscriberLocalPriority == nil {
		return nil, false
	}
	return o.DtoReceiveSubscriberLocalPriority, true
}

// HasDtoReceiveSubscriberLocalPriority returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasDtoReceiveSubscriberLocalPriority() bool {
	if o != nil && o.DtoReceiveSubscriberLocalPriority != nil {
		return true
	}

	return false
}

// SetDtoReceiveSubscriberLocalPriority gets a reference to the given int32 and assigns it to the DtoReceiveSubscriberLocalPriority field.
func (o *MsgVpnJndiConnectionFactory) SetDtoReceiveSubscriberLocalPriority(v int32) {
	o.DtoReceiveSubscriberLocalPriority = &v
}

// GetDtoReceiveSubscriberNetworkPriority returns the DtoReceiveSubscriberNetworkPriority field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetDtoReceiveSubscriberNetworkPriority() int32 {
	if o == nil || o.DtoReceiveSubscriberNetworkPriority == nil {
		var ret int32
		return ret
	}
	return *o.DtoReceiveSubscriberNetworkPriority
}

// GetDtoReceiveSubscriberNetworkPriorityOk returns a tuple with the DtoReceiveSubscriberNetworkPriority field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetDtoReceiveSubscriberNetworkPriorityOk() (*int32, bool) {
	if o == nil || o.DtoReceiveSubscriberNetworkPriority == nil {
		return nil, false
	}
	return o.DtoReceiveSubscriberNetworkPriority, true
}

// HasDtoReceiveSubscriberNetworkPriority returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasDtoReceiveSubscriberNetworkPriority() bool {
	if o != nil && o.DtoReceiveSubscriberNetworkPriority != nil {
		return true
	}

	return false
}

// SetDtoReceiveSubscriberNetworkPriority gets a reference to the given int32 and assigns it to the DtoReceiveSubscriberNetworkPriority field.
func (o *MsgVpnJndiConnectionFactory) SetDtoReceiveSubscriberNetworkPriority(v int32) {
	o.DtoReceiveSubscriberNetworkPriority = &v
}

// GetDtoSendEnabled returns the DtoSendEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetDtoSendEnabled() bool {
	if o == nil || o.DtoSendEnabled == nil {
		var ret bool
		return ret
	}
	return *o.DtoSendEnabled
}

// GetDtoSendEnabledOk returns a tuple with the DtoSendEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetDtoSendEnabledOk() (*bool, bool) {
	if o == nil || o.DtoSendEnabled == nil {
		return nil, false
	}
	return o.DtoSendEnabled, true
}

// HasDtoSendEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasDtoSendEnabled() bool {
	if o != nil && o.DtoSendEnabled != nil {
		return true
	}

	return false
}

// SetDtoSendEnabled gets a reference to the given bool and assigns it to the DtoSendEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetDtoSendEnabled(v bool) {
	o.DtoSendEnabled = &v
}

// GetDynamicEndpointCreateDurableEnabled returns the DynamicEndpointCreateDurableEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetDynamicEndpointCreateDurableEnabled() bool {
	if o == nil || o.DynamicEndpointCreateDurableEnabled == nil {
		var ret bool
		return ret
	}
	return *o.DynamicEndpointCreateDurableEnabled
}

// GetDynamicEndpointCreateDurableEnabledOk returns a tuple with the DynamicEndpointCreateDurableEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetDynamicEndpointCreateDurableEnabledOk() (*bool, bool) {
	if o == nil || o.DynamicEndpointCreateDurableEnabled == nil {
		return nil, false
	}
	return o.DynamicEndpointCreateDurableEnabled, true
}

// HasDynamicEndpointCreateDurableEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasDynamicEndpointCreateDurableEnabled() bool {
	if o != nil && o.DynamicEndpointCreateDurableEnabled != nil {
		return true
	}

	return false
}

// SetDynamicEndpointCreateDurableEnabled gets a reference to the given bool and assigns it to the DynamicEndpointCreateDurableEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetDynamicEndpointCreateDurableEnabled(v bool) {
	o.DynamicEndpointCreateDurableEnabled = &v
}

// GetDynamicEndpointRespectTtlEnabled returns the DynamicEndpointRespectTtlEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetDynamicEndpointRespectTtlEnabled() bool {
	if o == nil || o.DynamicEndpointRespectTtlEnabled == nil {
		var ret bool
		return ret
	}
	return *o.DynamicEndpointRespectTtlEnabled
}

// GetDynamicEndpointRespectTtlEnabledOk returns a tuple with the DynamicEndpointRespectTtlEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetDynamicEndpointRespectTtlEnabledOk() (*bool, bool) {
	if o == nil || o.DynamicEndpointRespectTtlEnabled == nil {
		return nil, false
	}
	return o.DynamicEndpointRespectTtlEnabled, true
}

// HasDynamicEndpointRespectTtlEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasDynamicEndpointRespectTtlEnabled() bool {
	if o != nil && o.DynamicEndpointRespectTtlEnabled != nil {
		return true
	}

	return false
}

// SetDynamicEndpointRespectTtlEnabled gets a reference to the given bool and assigns it to the DynamicEndpointRespectTtlEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetDynamicEndpointRespectTtlEnabled(v bool) {
	o.DynamicEndpointRespectTtlEnabled = &v
}

// GetGuaranteedReceiveAckTimeout returns the GuaranteedReceiveAckTimeout field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetGuaranteedReceiveAckTimeout() int32 {
	if o == nil || o.GuaranteedReceiveAckTimeout == nil {
		var ret int32
		return ret
	}
	return *o.GuaranteedReceiveAckTimeout
}

// GetGuaranteedReceiveAckTimeoutOk returns a tuple with the GuaranteedReceiveAckTimeout field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetGuaranteedReceiveAckTimeoutOk() (*int32, bool) {
	if o == nil || o.GuaranteedReceiveAckTimeout == nil {
		return nil, false
	}
	return o.GuaranteedReceiveAckTimeout, true
}

// HasGuaranteedReceiveAckTimeout returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasGuaranteedReceiveAckTimeout() bool {
	if o != nil && o.GuaranteedReceiveAckTimeout != nil {
		return true
	}

	return false
}

// SetGuaranteedReceiveAckTimeout gets a reference to the given int32 and assigns it to the GuaranteedReceiveAckTimeout field.
func (o *MsgVpnJndiConnectionFactory) SetGuaranteedReceiveAckTimeout(v int32) {
	o.GuaranteedReceiveAckTimeout = &v
}

// GetGuaranteedReceiveReconnectRetryCount returns the GuaranteedReceiveReconnectRetryCount field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetGuaranteedReceiveReconnectRetryCount() int32 {
	if o == nil || o.GuaranteedReceiveReconnectRetryCount == nil {
		var ret int32
		return ret
	}
	return *o.GuaranteedReceiveReconnectRetryCount
}

// GetGuaranteedReceiveReconnectRetryCountOk returns a tuple with the GuaranteedReceiveReconnectRetryCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetGuaranteedReceiveReconnectRetryCountOk() (*int32, bool) {
	if o == nil || o.GuaranteedReceiveReconnectRetryCount == nil {
		return nil, false
	}
	return o.GuaranteedReceiveReconnectRetryCount, true
}

// HasGuaranteedReceiveReconnectRetryCount returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasGuaranteedReceiveReconnectRetryCount() bool {
	if o != nil && o.GuaranteedReceiveReconnectRetryCount != nil {
		return true
	}

	return false
}

// SetGuaranteedReceiveReconnectRetryCount gets a reference to the given int32 and assigns it to the GuaranteedReceiveReconnectRetryCount field.
func (o *MsgVpnJndiConnectionFactory) SetGuaranteedReceiveReconnectRetryCount(v int32) {
	o.GuaranteedReceiveReconnectRetryCount = &v
}

// GetGuaranteedReceiveReconnectRetryWait returns the GuaranteedReceiveReconnectRetryWait field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetGuaranteedReceiveReconnectRetryWait() int32 {
	if o == nil || o.GuaranteedReceiveReconnectRetryWait == nil {
		var ret int32
		return ret
	}
	return *o.GuaranteedReceiveReconnectRetryWait
}

// GetGuaranteedReceiveReconnectRetryWaitOk returns a tuple with the GuaranteedReceiveReconnectRetryWait field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetGuaranteedReceiveReconnectRetryWaitOk() (*int32, bool) {
	if o == nil || o.GuaranteedReceiveReconnectRetryWait == nil {
		return nil, false
	}
	return o.GuaranteedReceiveReconnectRetryWait, true
}

// HasGuaranteedReceiveReconnectRetryWait returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasGuaranteedReceiveReconnectRetryWait() bool {
	if o != nil && o.GuaranteedReceiveReconnectRetryWait != nil {
		return true
	}

	return false
}

// SetGuaranteedReceiveReconnectRetryWait gets a reference to the given int32 and assigns it to the GuaranteedReceiveReconnectRetryWait field.
func (o *MsgVpnJndiConnectionFactory) SetGuaranteedReceiveReconnectRetryWait(v int32) {
	o.GuaranteedReceiveReconnectRetryWait = &v
}

// GetGuaranteedReceiveWindowSize returns the GuaranteedReceiveWindowSize field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetGuaranteedReceiveWindowSize() int32 {
	if o == nil || o.GuaranteedReceiveWindowSize == nil {
		var ret int32
		return ret
	}
	return *o.GuaranteedReceiveWindowSize
}

// GetGuaranteedReceiveWindowSizeOk returns a tuple with the GuaranteedReceiveWindowSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetGuaranteedReceiveWindowSizeOk() (*int32, bool) {
	if o == nil || o.GuaranteedReceiveWindowSize == nil {
		return nil, false
	}
	return o.GuaranteedReceiveWindowSize, true
}

// HasGuaranteedReceiveWindowSize returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasGuaranteedReceiveWindowSize() bool {
	if o != nil && o.GuaranteedReceiveWindowSize != nil {
		return true
	}

	return false
}

// SetGuaranteedReceiveWindowSize gets a reference to the given int32 and assigns it to the GuaranteedReceiveWindowSize field.
func (o *MsgVpnJndiConnectionFactory) SetGuaranteedReceiveWindowSize(v int32) {
	o.GuaranteedReceiveWindowSize = &v
}

// GetGuaranteedReceiveWindowSizeAckThreshold returns the GuaranteedReceiveWindowSizeAckThreshold field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetGuaranteedReceiveWindowSizeAckThreshold() int32 {
	if o == nil || o.GuaranteedReceiveWindowSizeAckThreshold == nil {
		var ret int32
		return ret
	}
	return *o.GuaranteedReceiveWindowSizeAckThreshold
}

// GetGuaranteedReceiveWindowSizeAckThresholdOk returns a tuple with the GuaranteedReceiveWindowSizeAckThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetGuaranteedReceiveWindowSizeAckThresholdOk() (*int32, bool) {
	if o == nil || o.GuaranteedReceiveWindowSizeAckThreshold == nil {
		return nil, false
	}
	return o.GuaranteedReceiveWindowSizeAckThreshold, true
}

// HasGuaranteedReceiveWindowSizeAckThreshold returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasGuaranteedReceiveWindowSizeAckThreshold() bool {
	if o != nil && o.GuaranteedReceiveWindowSizeAckThreshold != nil {
		return true
	}

	return false
}

// SetGuaranteedReceiveWindowSizeAckThreshold gets a reference to the given int32 and assigns it to the GuaranteedReceiveWindowSizeAckThreshold field.
func (o *MsgVpnJndiConnectionFactory) SetGuaranteedReceiveWindowSizeAckThreshold(v int32) {
	o.GuaranteedReceiveWindowSizeAckThreshold = &v
}

// GetGuaranteedSendAckTimeout returns the GuaranteedSendAckTimeout field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetGuaranteedSendAckTimeout() int32 {
	if o == nil || o.GuaranteedSendAckTimeout == nil {
		var ret int32
		return ret
	}
	return *o.GuaranteedSendAckTimeout
}

// GetGuaranteedSendAckTimeoutOk returns a tuple with the GuaranteedSendAckTimeout field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetGuaranteedSendAckTimeoutOk() (*int32, bool) {
	if o == nil || o.GuaranteedSendAckTimeout == nil {
		return nil, false
	}
	return o.GuaranteedSendAckTimeout, true
}

// HasGuaranteedSendAckTimeout returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasGuaranteedSendAckTimeout() bool {
	if o != nil && o.GuaranteedSendAckTimeout != nil {
		return true
	}

	return false
}

// SetGuaranteedSendAckTimeout gets a reference to the given int32 and assigns it to the GuaranteedSendAckTimeout field.
func (o *MsgVpnJndiConnectionFactory) SetGuaranteedSendAckTimeout(v int32) {
	o.GuaranteedSendAckTimeout = &v
}

// GetGuaranteedSendWindowSize returns the GuaranteedSendWindowSize field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetGuaranteedSendWindowSize() int32 {
	if o == nil || o.GuaranteedSendWindowSize == nil {
		var ret int32
		return ret
	}
	return *o.GuaranteedSendWindowSize
}

// GetGuaranteedSendWindowSizeOk returns a tuple with the GuaranteedSendWindowSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetGuaranteedSendWindowSizeOk() (*int32, bool) {
	if o == nil || o.GuaranteedSendWindowSize == nil {
		return nil, false
	}
	return o.GuaranteedSendWindowSize, true
}

// HasGuaranteedSendWindowSize returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasGuaranteedSendWindowSize() bool {
	if o != nil && o.GuaranteedSendWindowSize != nil {
		return true
	}

	return false
}

// SetGuaranteedSendWindowSize gets a reference to the given int32 and assigns it to the GuaranteedSendWindowSize field.
func (o *MsgVpnJndiConnectionFactory) SetGuaranteedSendWindowSize(v int32) {
	o.GuaranteedSendWindowSize = &v
}

// GetMessagingDefaultDeliveryMode returns the MessagingDefaultDeliveryMode field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetMessagingDefaultDeliveryMode() string {
	if o == nil || o.MessagingDefaultDeliveryMode == nil {
		var ret string
		return ret
	}
	return *o.MessagingDefaultDeliveryMode
}

// GetMessagingDefaultDeliveryModeOk returns a tuple with the MessagingDefaultDeliveryMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetMessagingDefaultDeliveryModeOk() (*string, bool) {
	if o == nil || o.MessagingDefaultDeliveryMode == nil {
		return nil, false
	}
	return o.MessagingDefaultDeliveryMode, true
}

// HasMessagingDefaultDeliveryMode returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasMessagingDefaultDeliveryMode() bool {
	if o != nil && o.MessagingDefaultDeliveryMode != nil {
		return true
	}

	return false
}

// SetMessagingDefaultDeliveryMode gets a reference to the given string and assigns it to the MessagingDefaultDeliveryMode field.
func (o *MsgVpnJndiConnectionFactory) SetMessagingDefaultDeliveryMode(v string) {
	o.MessagingDefaultDeliveryMode = &v
}

// GetMessagingDefaultDmqEligibleEnabled returns the MessagingDefaultDmqEligibleEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetMessagingDefaultDmqEligibleEnabled() bool {
	if o == nil || o.MessagingDefaultDmqEligibleEnabled == nil {
		var ret bool
		return ret
	}
	return *o.MessagingDefaultDmqEligibleEnabled
}

// GetMessagingDefaultDmqEligibleEnabledOk returns a tuple with the MessagingDefaultDmqEligibleEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetMessagingDefaultDmqEligibleEnabledOk() (*bool, bool) {
	if o == nil || o.MessagingDefaultDmqEligibleEnabled == nil {
		return nil, false
	}
	return o.MessagingDefaultDmqEligibleEnabled, true
}

// HasMessagingDefaultDmqEligibleEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasMessagingDefaultDmqEligibleEnabled() bool {
	if o != nil && o.MessagingDefaultDmqEligibleEnabled != nil {
		return true
	}

	return false
}

// SetMessagingDefaultDmqEligibleEnabled gets a reference to the given bool and assigns it to the MessagingDefaultDmqEligibleEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetMessagingDefaultDmqEligibleEnabled(v bool) {
	o.MessagingDefaultDmqEligibleEnabled = &v
}

// GetMessagingDefaultElidingEligibleEnabled returns the MessagingDefaultElidingEligibleEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetMessagingDefaultElidingEligibleEnabled() bool {
	if o == nil || o.MessagingDefaultElidingEligibleEnabled == nil {
		var ret bool
		return ret
	}
	return *o.MessagingDefaultElidingEligibleEnabled
}

// GetMessagingDefaultElidingEligibleEnabledOk returns a tuple with the MessagingDefaultElidingEligibleEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetMessagingDefaultElidingEligibleEnabledOk() (*bool, bool) {
	if o == nil || o.MessagingDefaultElidingEligibleEnabled == nil {
		return nil, false
	}
	return o.MessagingDefaultElidingEligibleEnabled, true
}

// HasMessagingDefaultElidingEligibleEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasMessagingDefaultElidingEligibleEnabled() bool {
	if o != nil && o.MessagingDefaultElidingEligibleEnabled != nil {
		return true
	}

	return false
}

// SetMessagingDefaultElidingEligibleEnabled gets a reference to the given bool and assigns it to the MessagingDefaultElidingEligibleEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetMessagingDefaultElidingEligibleEnabled(v bool) {
	o.MessagingDefaultElidingEligibleEnabled = &v
}

// GetMessagingJmsxUserIdEnabled returns the MessagingJmsxUserIdEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetMessagingJmsxUserIdEnabled() bool {
	if o == nil || o.MessagingJmsxUserIdEnabled == nil {
		var ret bool
		return ret
	}
	return *o.MessagingJmsxUserIdEnabled
}

// GetMessagingJmsxUserIdEnabledOk returns a tuple with the MessagingJmsxUserIdEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetMessagingJmsxUserIdEnabledOk() (*bool, bool) {
	if o == nil || o.MessagingJmsxUserIdEnabled == nil {
		return nil, false
	}
	return o.MessagingJmsxUserIdEnabled, true
}

// HasMessagingJmsxUserIdEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasMessagingJmsxUserIdEnabled() bool {
	if o != nil && o.MessagingJmsxUserIdEnabled != nil {
		return true
	}

	return false
}

// SetMessagingJmsxUserIdEnabled gets a reference to the given bool and assigns it to the MessagingJmsxUserIdEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetMessagingJmsxUserIdEnabled(v bool) {
	o.MessagingJmsxUserIdEnabled = &v
}

// GetMessagingTextInXmlPayloadEnabled returns the MessagingTextInXmlPayloadEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetMessagingTextInXmlPayloadEnabled() bool {
	if o == nil || o.MessagingTextInXmlPayloadEnabled == nil {
		var ret bool
		return ret
	}
	return *o.MessagingTextInXmlPayloadEnabled
}

// GetMessagingTextInXmlPayloadEnabledOk returns a tuple with the MessagingTextInXmlPayloadEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetMessagingTextInXmlPayloadEnabledOk() (*bool, bool) {
	if o == nil || o.MessagingTextInXmlPayloadEnabled == nil {
		return nil, false
	}
	return o.MessagingTextInXmlPayloadEnabled, true
}

// HasMessagingTextInXmlPayloadEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasMessagingTextInXmlPayloadEnabled() bool {
	if o != nil && o.MessagingTextInXmlPayloadEnabled != nil {
		return true
	}

	return false
}

// SetMessagingTextInXmlPayloadEnabled gets a reference to the given bool and assigns it to the MessagingTextInXmlPayloadEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetMessagingTextInXmlPayloadEnabled(v bool) {
	o.MessagingTextInXmlPayloadEnabled = &v
}

// GetMsgVpnName returns the MsgVpnName field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetMsgVpnName() string {
	if o == nil || o.MsgVpnName == nil {
		var ret string
		return ret
	}
	return *o.MsgVpnName
}

// GetMsgVpnNameOk returns a tuple with the MsgVpnName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetMsgVpnNameOk() (*string, bool) {
	if o == nil || o.MsgVpnName == nil {
		return nil, false
	}
	return o.MsgVpnName, true
}

// HasMsgVpnName returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasMsgVpnName() bool {
	if o != nil && o.MsgVpnName != nil {
		return true
	}

	return false
}

// SetMsgVpnName gets a reference to the given string and assigns it to the MsgVpnName field.
func (o *MsgVpnJndiConnectionFactory) SetMsgVpnName(v string) {
	o.MsgVpnName = &v
}

// GetTransportCompressionLevel returns the TransportCompressionLevel field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportCompressionLevel() int32 {
	if o == nil || o.TransportCompressionLevel == nil {
		var ret int32
		return ret
	}
	return *o.TransportCompressionLevel
}

// GetTransportCompressionLevelOk returns a tuple with the TransportCompressionLevel field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportCompressionLevelOk() (*int32, bool) {
	if o == nil || o.TransportCompressionLevel == nil {
		return nil, false
	}
	return o.TransportCompressionLevel, true
}

// HasTransportCompressionLevel returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportCompressionLevel() bool {
	if o != nil && o.TransportCompressionLevel != nil {
		return true
	}

	return false
}

// SetTransportCompressionLevel gets a reference to the given int32 and assigns it to the TransportCompressionLevel field.
func (o *MsgVpnJndiConnectionFactory) SetTransportCompressionLevel(v int32) {
	o.TransportCompressionLevel = &v
}

// GetTransportConnectRetryCount returns the TransportConnectRetryCount field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportConnectRetryCount() int32 {
	if o == nil || o.TransportConnectRetryCount == nil {
		var ret int32
		return ret
	}
	return *o.TransportConnectRetryCount
}

// GetTransportConnectRetryCountOk returns a tuple with the TransportConnectRetryCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportConnectRetryCountOk() (*int32, bool) {
	if o == nil || o.TransportConnectRetryCount == nil {
		return nil, false
	}
	return o.TransportConnectRetryCount, true
}

// HasTransportConnectRetryCount returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportConnectRetryCount() bool {
	if o != nil && o.TransportConnectRetryCount != nil {
		return true
	}

	return false
}

// SetTransportConnectRetryCount gets a reference to the given int32 and assigns it to the TransportConnectRetryCount field.
func (o *MsgVpnJndiConnectionFactory) SetTransportConnectRetryCount(v int32) {
	o.TransportConnectRetryCount = &v
}

// GetTransportConnectRetryPerHostCount returns the TransportConnectRetryPerHostCount field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportConnectRetryPerHostCount() int32 {
	if o == nil || o.TransportConnectRetryPerHostCount == nil {
		var ret int32
		return ret
	}
	return *o.TransportConnectRetryPerHostCount
}

// GetTransportConnectRetryPerHostCountOk returns a tuple with the TransportConnectRetryPerHostCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportConnectRetryPerHostCountOk() (*int32, bool) {
	if o == nil || o.TransportConnectRetryPerHostCount == nil {
		return nil, false
	}
	return o.TransportConnectRetryPerHostCount, true
}

// HasTransportConnectRetryPerHostCount returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportConnectRetryPerHostCount() bool {
	if o != nil && o.TransportConnectRetryPerHostCount != nil {
		return true
	}

	return false
}

// SetTransportConnectRetryPerHostCount gets a reference to the given int32 and assigns it to the TransportConnectRetryPerHostCount field.
func (o *MsgVpnJndiConnectionFactory) SetTransportConnectRetryPerHostCount(v int32) {
	o.TransportConnectRetryPerHostCount = &v
}

// GetTransportConnectTimeout returns the TransportConnectTimeout field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportConnectTimeout() int32 {
	if o == nil || o.TransportConnectTimeout == nil {
		var ret int32
		return ret
	}
	return *o.TransportConnectTimeout
}

// GetTransportConnectTimeoutOk returns a tuple with the TransportConnectTimeout field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportConnectTimeoutOk() (*int32, bool) {
	if o == nil || o.TransportConnectTimeout == nil {
		return nil, false
	}
	return o.TransportConnectTimeout, true
}

// HasTransportConnectTimeout returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportConnectTimeout() bool {
	if o != nil && o.TransportConnectTimeout != nil {
		return true
	}

	return false
}

// SetTransportConnectTimeout gets a reference to the given int32 and assigns it to the TransportConnectTimeout field.
func (o *MsgVpnJndiConnectionFactory) SetTransportConnectTimeout(v int32) {
	o.TransportConnectTimeout = &v
}

// GetTransportDirectTransportEnabled returns the TransportDirectTransportEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportDirectTransportEnabled() bool {
	if o == nil || o.TransportDirectTransportEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TransportDirectTransportEnabled
}

// GetTransportDirectTransportEnabledOk returns a tuple with the TransportDirectTransportEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportDirectTransportEnabledOk() (*bool, bool) {
	if o == nil || o.TransportDirectTransportEnabled == nil {
		return nil, false
	}
	return o.TransportDirectTransportEnabled, true
}

// HasTransportDirectTransportEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportDirectTransportEnabled() bool {
	if o != nil && o.TransportDirectTransportEnabled != nil {
		return true
	}

	return false
}

// SetTransportDirectTransportEnabled gets a reference to the given bool and assigns it to the TransportDirectTransportEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetTransportDirectTransportEnabled(v bool) {
	o.TransportDirectTransportEnabled = &v
}

// GetTransportKeepaliveCount returns the TransportKeepaliveCount field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportKeepaliveCount() int32 {
	if o == nil || o.TransportKeepaliveCount == nil {
		var ret int32
		return ret
	}
	return *o.TransportKeepaliveCount
}

// GetTransportKeepaliveCountOk returns a tuple with the TransportKeepaliveCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportKeepaliveCountOk() (*int32, bool) {
	if o == nil || o.TransportKeepaliveCount == nil {
		return nil, false
	}
	return o.TransportKeepaliveCount, true
}

// HasTransportKeepaliveCount returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportKeepaliveCount() bool {
	if o != nil && o.TransportKeepaliveCount != nil {
		return true
	}

	return false
}

// SetTransportKeepaliveCount gets a reference to the given int32 and assigns it to the TransportKeepaliveCount field.
func (o *MsgVpnJndiConnectionFactory) SetTransportKeepaliveCount(v int32) {
	o.TransportKeepaliveCount = &v
}

// GetTransportKeepaliveEnabled returns the TransportKeepaliveEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportKeepaliveEnabled() bool {
	if o == nil || o.TransportKeepaliveEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TransportKeepaliveEnabled
}

// GetTransportKeepaliveEnabledOk returns a tuple with the TransportKeepaliveEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportKeepaliveEnabledOk() (*bool, bool) {
	if o == nil || o.TransportKeepaliveEnabled == nil {
		return nil, false
	}
	return o.TransportKeepaliveEnabled, true
}

// HasTransportKeepaliveEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportKeepaliveEnabled() bool {
	if o != nil && o.TransportKeepaliveEnabled != nil {
		return true
	}

	return false
}

// SetTransportKeepaliveEnabled gets a reference to the given bool and assigns it to the TransportKeepaliveEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetTransportKeepaliveEnabled(v bool) {
	o.TransportKeepaliveEnabled = &v
}

// GetTransportKeepaliveInterval returns the TransportKeepaliveInterval field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportKeepaliveInterval() int32 {
	if o == nil || o.TransportKeepaliveInterval == nil {
		var ret int32
		return ret
	}
	return *o.TransportKeepaliveInterval
}

// GetTransportKeepaliveIntervalOk returns a tuple with the TransportKeepaliveInterval field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportKeepaliveIntervalOk() (*int32, bool) {
	if o == nil || o.TransportKeepaliveInterval == nil {
		return nil, false
	}
	return o.TransportKeepaliveInterval, true
}

// HasTransportKeepaliveInterval returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportKeepaliveInterval() bool {
	if o != nil && o.TransportKeepaliveInterval != nil {
		return true
	}

	return false
}

// SetTransportKeepaliveInterval gets a reference to the given int32 and assigns it to the TransportKeepaliveInterval field.
func (o *MsgVpnJndiConnectionFactory) SetTransportKeepaliveInterval(v int32) {
	o.TransportKeepaliveInterval = &v
}

// GetTransportMsgCallbackOnIoThreadEnabled returns the TransportMsgCallbackOnIoThreadEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportMsgCallbackOnIoThreadEnabled() bool {
	if o == nil || o.TransportMsgCallbackOnIoThreadEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TransportMsgCallbackOnIoThreadEnabled
}

// GetTransportMsgCallbackOnIoThreadEnabledOk returns a tuple with the TransportMsgCallbackOnIoThreadEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportMsgCallbackOnIoThreadEnabledOk() (*bool, bool) {
	if o == nil || o.TransportMsgCallbackOnIoThreadEnabled == nil {
		return nil, false
	}
	return o.TransportMsgCallbackOnIoThreadEnabled, true
}

// HasTransportMsgCallbackOnIoThreadEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportMsgCallbackOnIoThreadEnabled() bool {
	if o != nil && o.TransportMsgCallbackOnIoThreadEnabled != nil {
		return true
	}

	return false
}

// SetTransportMsgCallbackOnIoThreadEnabled gets a reference to the given bool and assigns it to the TransportMsgCallbackOnIoThreadEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetTransportMsgCallbackOnIoThreadEnabled(v bool) {
	o.TransportMsgCallbackOnIoThreadEnabled = &v
}

// GetTransportOptimizeDirectEnabled returns the TransportOptimizeDirectEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportOptimizeDirectEnabled() bool {
	if o == nil || o.TransportOptimizeDirectEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TransportOptimizeDirectEnabled
}

// GetTransportOptimizeDirectEnabledOk returns a tuple with the TransportOptimizeDirectEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportOptimizeDirectEnabledOk() (*bool, bool) {
	if o == nil || o.TransportOptimizeDirectEnabled == nil {
		return nil, false
	}
	return o.TransportOptimizeDirectEnabled, true
}

// HasTransportOptimizeDirectEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportOptimizeDirectEnabled() bool {
	if o != nil && o.TransportOptimizeDirectEnabled != nil {
		return true
	}

	return false
}

// SetTransportOptimizeDirectEnabled gets a reference to the given bool and assigns it to the TransportOptimizeDirectEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetTransportOptimizeDirectEnabled(v bool) {
	o.TransportOptimizeDirectEnabled = &v
}

// GetTransportPort returns the TransportPort field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportPort() int32 {
	if o == nil || o.TransportPort == nil {
		var ret int32
		return ret
	}
	return *o.TransportPort
}

// GetTransportPortOk returns a tuple with the TransportPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportPortOk() (*int32, bool) {
	if o == nil || o.TransportPort == nil {
		return nil, false
	}
	return o.TransportPort, true
}

// HasTransportPort returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportPort() bool {
	if o != nil && o.TransportPort != nil {
		return true
	}

	return false
}

// SetTransportPort gets a reference to the given int32 and assigns it to the TransportPort field.
func (o *MsgVpnJndiConnectionFactory) SetTransportPort(v int32) {
	o.TransportPort = &v
}

// GetTransportReadTimeout returns the TransportReadTimeout field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportReadTimeout() int32 {
	if o == nil || o.TransportReadTimeout == nil {
		var ret int32
		return ret
	}
	return *o.TransportReadTimeout
}

// GetTransportReadTimeoutOk returns a tuple with the TransportReadTimeout field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportReadTimeoutOk() (*int32, bool) {
	if o == nil || o.TransportReadTimeout == nil {
		return nil, false
	}
	return o.TransportReadTimeout, true
}

// HasTransportReadTimeout returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportReadTimeout() bool {
	if o != nil && o.TransportReadTimeout != nil {
		return true
	}

	return false
}

// SetTransportReadTimeout gets a reference to the given int32 and assigns it to the TransportReadTimeout field.
func (o *MsgVpnJndiConnectionFactory) SetTransportReadTimeout(v int32) {
	o.TransportReadTimeout = &v
}

// GetTransportReceiveBufferSize returns the TransportReceiveBufferSize field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportReceiveBufferSize() int32 {
	if o == nil || o.TransportReceiveBufferSize == nil {
		var ret int32
		return ret
	}
	return *o.TransportReceiveBufferSize
}

// GetTransportReceiveBufferSizeOk returns a tuple with the TransportReceiveBufferSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportReceiveBufferSizeOk() (*int32, bool) {
	if o == nil || o.TransportReceiveBufferSize == nil {
		return nil, false
	}
	return o.TransportReceiveBufferSize, true
}

// HasTransportReceiveBufferSize returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportReceiveBufferSize() bool {
	if o != nil && o.TransportReceiveBufferSize != nil {
		return true
	}

	return false
}

// SetTransportReceiveBufferSize gets a reference to the given int32 and assigns it to the TransportReceiveBufferSize field.
func (o *MsgVpnJndiConnectionFactory) SetTransportReceiveBufferSize(v int32) {
	o.TransportReceiveBufferSize = &v
}

// GetTransportReconnectRetryCount returns the TransportReconnectRetryCount field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportReconnectRetryCount() int32 {
	if o == nil || o.TransportReconnectRetryCount == nil {
		var ret int32
		return ret
	}
	return *o.TransportReconnectRetryCount
}

// GetTransportReconnectRetryCountOk returns a tuple with the TransportReconnectRetryCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportReconnectRetryCountOk() (*int32, bool) {
	if o == nil || o.TransportReconnectRetryCount == nil {
		return nil, false
	}
	return o.TransportReconnectRetryCount, true
}

// HasTransportReconnectRetryCount returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportReconnectRetryCount() bool {
	if o != nil && o.TransportReconnectRetryCount != nil {
		return true
	}

	return false
}

// SetTransportReconnectRetryCount gets a reference to the given int32 and assigns it to the TransportReconnectRetryCount field.
func (o *MsgVpnJndiConnectionFactory) SetTransportReconnectRetryCount(v int32) {
	o.TransportReconnectRetryCount = &v
}

// GetTransportReconnectRetryWait returns the TransportReconnectRetryWait field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportReconnectRetryWait() int32 {
	if o == nil || o.TransportReconnectRetryWait == nil {
		var ret int32
		return ret
	}
	return *o.TransportReconnectRetryWait
}

// GetTransportReconnectRetryWaitOk returns a tuple with the TransportReconnectRetryWait field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportReconnectRetryWaitOk() (*int32, bool) {
	if o == nil || o.TransportReconnectRetryWait == nil {
		return nil, false
	}
	return o.TransportReconnectRetryWait, true
}

// HasTransportReconnectRetryWait returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportReconnectRetryWait() bool {
	if o != nil && o.TransportReconnectRetryWait != nil {
		return true
	}

	return false
}

// SetTransportReconnectRetryWait gets a reference to the given int32 and assigns it to the TransportReconnectRetryWait field.
func (o *MsgVpnJndiConnectionFactory) SetTransportReconnectRetryWait(v int32) {
	o.TransportReconnectRetryWait = &v
}

// GetTransportSendBufferSize returns the TransportSendBufferSize field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportSendBufferSize() int32 {
	if o == nil || o.TransportSendBufferSize == nil {
		var ret int32
		return ret
	}
	return *o.TransportSendBufferSize
}

// GetTransportSendBufferSizeOk returns a tuple with the TransportSendBufferSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportSendBufferSizeOk() (*int32, bool) {
	if o == nil || o.TransportSendBufferSize == nil {
		return nil, false
	}
	return o.TransportSendBufferSize, true
}

// HasTransportSendBufferSize returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportSendBufferSize() bool {
	if o != nil && o.TransportSendBufferSize != nil {
		return true
	}

	return false
}

// SetTransportSendBufferSize gets a reference to the given int32 and assigns it to the TransportSendBufferSize field.
func (o *MsgVpnJndiConnectionFactory) SetTransportSendBufferSize(v int32) {
	o.TransportSendBufferSize = &v
}

// GetTransportTcpNoDelayEnabled returns the TransportTcpNoDelayEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetTransportTcpNoDelayEnabled() bool {
	if o == nil || o.TransportTcpNoDelayEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TransportTcpNoDelayEnabled
}

// GetTransportTcpNoDelayEnabledOk returns a tuple with the TransportTcpNoDelayEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetTransportTcpNoDelayEnabledOk() (*bool, bool) {
	if o == nil || o.TransportTcpNoDelayEnabled == nil {
		return nil, false
	}
	return o.TransportTcpNoDelayEnabled, true
}

// HasTransportTcpNoDelayEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasTransportTcpNoDelayEnabled() bool {
	if o != nil && o.TransportTcpNoDelayEnabled != nil {
		return true
	}

	return false
}

// SetTransportTcpNoDelayEnabled gets a reference to the given bool and assigns it to the TransportTcpNoDelayEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetTransportTcpNoDelayEnabled(v bool) {
	o.TransportTcpNoDelayEnabled = &v
}

// GetXaEnabled returns the XaEnabled field value if set, zero value otherwise.
func (o *MsgVpnJndiConnectionFactory) GetXaEnabled() bool {
	if o == nil || o.XaEnabled == nil {
		var ret bool
		return ret
	}
	return *o.XaEnabled
}

// GetXaEnabledOk returns a tuple with the XaEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpnJndiConnectionFactory) GetXaEnabledOk() (*bool, bool) {
	if o == nil || o.XaEnabled == nil {
		return nil, false
	}
	return o.XaEnabled, true
}

// HasXaEnabled returns a boolean if a field has been set.
func (o *MsgVpnJndiConnectionFactory) HasXaEnabled() bool {
	if o != nil && o.XaEnabled != nil {
		return true
	}

	return false
}

// SetXaEnabled gets a reference to the given bool and assigns it to the XaEnabled field.
func (o *MsgVpnJndiConnectionFactory) SetXaEnabled(v bool) {
	o.XaEnabled = &v
}

func (o MsgVpnJndiConnectionFactory) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AllowDuplicateClientIdEnabled != nil {
		toSerialize["allowDuplicateClientIdEnabled"] = o.AllowDuplicateClientIdEnabled
	}
	if o.ClientDescription != nil {
		toSerialize["clientDescription"] = o.ClientDescription
	}
	if o.ClientId != nil {
		toSerialize["clientId"] = o.ClientId
	}
	if o.ConnectionFactoryName != nil {
		toSerialize["connectionFactoryName"] = o.ConnectionFactoryName
	}
	if o.DtoReceiveOverrideEnabled != nil {
		toSerialize["dtoReceiveOverrideEnabled"] = o.DtoReceiveOverrideEnabled
	}
	if o.DtoReceiveSubscriberLocalPriority != nil {
		toSerialize["dtoReceiveSubscriberLocalPriority"] = o.DtoReceiveSubscriberLocalPriority
	}
	if o.DtoReceiveSubscriberNetworkPriority != nil {
		toSerialize["dtoReceiveSubscriberNetworkPriority"] = o.DtoReceiveSubscriberNetworkPriority
	}
	if o.DtoSendEnabled != nil {
		toSerialize["dtoSendEnabled"] = o.DtoSendEnabled
	}
	if o.DynamicEndpointCreateDurableEnabled != nil {
		toSerialize["dynamicEndpointCreateDurableEnabled"] = o.DynamicEndpointCreateDurableEnabled
	}
	if o.DynamicEndpointRespectTtlEnabled != nil {
		toSerialize["dynamicEndpointRespectTtlEnabled"] = o.DynamicEndpointRespectTtlEnabled
	}
	if o.GuaranteedReceiveAckTimeout != nil {
		toSerialize["guaranteedReceiveAckTimeout"] = o.GuaranteedReceiveAckTimeout
	}
	if o.GuaranteedReceiveReconnectRetryCount != nil {
		toSerialize["guaranteedReceiveReconnectRetryCount"] = o.GuaranteedReceiveReconnectRetryCount
	}
	if o.GuaranteedReceiveReconnectRetryWait != nil {
		toSerialize["guaranteedReceiveReconnectRetryWait"] = o.GuaranteedReceiveReconnectRetryWait
	}
	if o.GuaranteedReceiveWindowSize != nil {
		toSerialize["guaranteedReceiveWindowSize"] = o.GuaranteedReceiveWindowSize
	}
	if o.GuaranteedReceiveWindowSizeAckThreshold != nil {
		toSerialize["guaranteedReceiveWindowSizeAckThreshold"] = o.GuaranteedReceiveWindowSizeAckThreshold
	}
	if o.GuaranteedSendAckTimeout != nil {
		toSerialize["guaranteedSendAckTimeout"] = o.GuaranteedSendAckTimeout
	}
	if o.GuaranteedSendWindowSize != nil {
		toSerialize["guaranteedSendWindowSize"] = o.GuaranteedSendWindowSize
	}
	if o.MessagingDefaultDeliveryMode != nil {
		toSerialize["messagingDefaultDeliveryMode"] = o.MessagingDefaultDeliveryMode
	}
	if o.MessagingDefaultDmqEligibleEnabled != nil {
		toSerialize["messagingDefaultDmqEligibleEnabled"] = o.MessagingDefaultDmqEligibleEnabled
	}
	if o.MessagingDefaultElidingEligibleEnabled != nil {
		toSerialize["messagingDefaultElidingEligibleEnabled"] = o.MessagingDefaultElidingEligibleEnabled
	}
	if o.MessagingJmsxUserIdEnabled != nil {
		toSerialize["messagingJmsxUserIdEnabled"] = o.MessagingJmsxUserIdEnabled
	}
	if o.MessagingTextInXmlPayloadEnabled != nil {
		toSerialize["messagingTextInXmlPayloadEnabled"] = o.MessagingTextInXmlPayloadEnabled
	}
	if o.MsgVpnName != nil {
		toSerialize["msgVpnName"] = o.MsgVpnName
	}
	if o.TransportCompressionLevel != nil {
		toSerialize["transportCompressionLevel"] = o.TransportCompressionLevel
	}
	if o.TransportConnectRetryCount != nil {
		toSerialize["transportConnectRetryCount"] = o.TransportConnectRetryCount
	}
	if o.TransportConnectRetryPerHostCount != nil {
		toSerialize["transportConnectRetryPerHostCount"] = o.TransportConnectRetryPerHostCount
	}
	if o.TransportConnectTimeout != nil {
		toSerialize["transportConnectTimeout"] = o.TransportConnectTimeout
	}
	if o.TransportDirectTransportEnabled != nil {
		toSerialize["transportDirectTransportEnabled"] = o.TransportDirectTransportEnabled
	}
	if o.TransportKeepaliveCount != nil {
		toSerialize["transportKeepaliveCount"] = o.TransportKeepaliveCount
	}
	if o.TransportKeepaliveEnabled != nil {
		toSerialize["transportKeepaliveEnabled"] = o.TransportKeepaliveEnabled
	}
	if o.TransportKeepaliveInterval != nil {
		toSerialize["transportKeepaliveInterval"] = o.TransportKeepaliveInterval
	}
	if o.TransportMsgCallbackOnIoThreadEnabled != nil {
		toSerialize["transportMsgCallbackOnIoThreadEnabled"] = o.TransportMsgCallbackOnIoThreadEnabled
	}
	if o.TransportOptimizeDirectEnabled != nil {
		toSerialize["transportOptimizeDirectEnabled"] = o.TransportOptimizeDirectEnabled
	}
	if o.TransportPort != nil {
		toSerialize["transportPort"] = o.TransportPort
	}
	if o.TransportReadTimeout != nil {
		toSerialize["transportReadTimeout"] = o.TransportReadTimeout
	}
	if o.TransportReceiveBufferSize != nil {
		toSerialize["transportReceiveBufferSize"] = o.TransportReceiveBufferSize
	}
	if o.TransportReconnectRetryCount != nil {
		toSerialize["transportReconnectRetryCount"] = o.TransportReconnectRetryCount
	}
	if o.TransportReconnectRetryWait != nil {
		toSerialize["transportReconnectRetryWait"] = o.TransportReconnectRetryWait
	}
	if o.TransportSendBufferSize != nil {
		toSerialize["transportSendBufferSize"] = o.TransportSendBufferSize
	}
	if o.TransportTcpNoDelayEnabled != nil {
		toSerialize["transportTcpNoDelayEnabled"] = o.TransportTcpNoDelayEnabled
	}
	if o.XaEnabled != nil {
		toSerialize["xaEnabled"] = o.XaEnabled
	}
	return json.Marshal(toSerialize)
}

type NullableMsgVpnJndiConnectionFactory struct {
	value *MsgVpnJndiConnectionFactory
	isSet bool
}

func (v NullableMsgVpnJndiConnectionFactory) Get() *MsgVpnJndiConnectionFactory {
	return v.value
}

func (v *NullableMsgVpnJndiConnectionFactory) Set(val *MsgVpnJndiConnectionFactory) {
	v.value = val
	v.isSet = true
}

func (v NullableMsgVpnJndiConnectionFactory) IsSet() bool {
	return v.isSet
}

func (v *NullableMsgVpnJndiConnectionFactory) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMsgVpnJndiConnectionFactory(val *MsgVpnJndiConnectionFactory) *NullableMsgVpnJndiConnectionFactory {
	return &NullableMsgVpnJndiConnectionFactory{value: val, isSet: true}
}

func (v NullableMsgVpnJndiConnectionFactory) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMsgVpnJndiConnectionFactory) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


