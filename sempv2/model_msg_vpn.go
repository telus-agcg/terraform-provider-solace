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

// MsgVpn struct for MsgVpn
type MsgVpn struct {
	// The name of another Message VPN which this Message VPN is an alias for. When this Message VPN is enabled, the alias has no effect. When this Message VPN is disabled, Clients (but not Bridges and routing Links) logging into this Message VPN are automatically logged in to the other Message VPN, and authentication and authorization take place in the context of the other Message VPN.  Aliases may form a non-circular chain, cascading one to the next. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.14.
	Alias *string `json:"alias,omitempty"`
	// Enable or disable basic authentication for clients connecting to the Message VPN. Basic authentication is authentication that involves the use of a username and password to prove identity. If a user provides credentials for a different authentication scheme, this setting is not applicable. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	AuthenticationBasicEnabled *bool `json:"authenticationBasicEnabled,omitempty"`
	// The name of the RADIUS or LDAP Profile to use for basic authentication. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"default\"`.
	AuthenticationBasicProfileName *string `json:"authenticationBasicProfileName,omitempty"`
	// The RADIUS domain to use for basic authentication. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.
	AuthenticationBasicRadiusDomain *string `json:"authenticationBasicRadiusDomain,omitempty"`
	// The type of basic authentication to use for clients connecting to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"radius\"`. The allowed values and their meaning are:  <pre> \"internal\" - Internal database. Authentication is against Client Usernames. \"ldap\" - LDAP authentication. An LDAP profile name must be provided. \"radius\" - RADIUS authentication. A RADIUS profile name must be provided. \"none\" - No authentication. Anonymous login allowed. </pre>
	AuthenticationBasicType *string `json:"authenticationBasicType,omitempty"`
	// Enable or disable allowing a client to specify a Client Username via the API connect method. When disabled, the certificate CN (Common Name) is always used. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	AuthenticationClientCertAllowApiProvidedUsernameEnabled *bool `json:"authenticationClientCertAllowApiProvidedUsernameEnabled,omitempty"`
	// Enable or disable client certificate authentication in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	AuthenticationClientCertEnabled *bool `json:"authenticationClientCertEnabled,omitempty"`
	// The maximum depth for a client certificate chain. The depth of a chain is defined as the number of signing CA certificates that are present in the chain back to a trusted self-signed root CA certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.
	AuthenticationClientCertMaxChainDepth *int64 `json:"authenticationClientCertMaxChainDepth,omitempty"`
	// The desired behavior for client certificate revocation checking. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"allow-valid\"`. The allowed values and their meaning are:  <pre> \"allow-all\" - Allow the client to authenticate, the result of client certificate revocation check is ignored. \"allow-unknown\" - Allow the client to authenticate even if the revocation status of his certificate cannot be determined. \"allow-valid\" - Allow the client to authenticate only when the revocation check returned an explicit positive response. </pre>  Available since 2.6.
	AuthenticationClientCertRevocationCheckMode *string `json:"authenticationClientCertRevocationCheckMode,omitempty"`
	// The field from the client certificate to use as the client username. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"common-name\"`. The allowed values and their meaning are:  <pre> \"certificate-thumbprint\" - The username is computed as the SHA-1 hash over the entire DER-encoded contents of the client certificate. \"common-name\" - The username is extracted from the certificate's first instance of the Common Name attribute in the Subject DN. \"common-name-last\" - The username is extracted from the certificate's last instance of the Common Name attribute in the Subject DN. \"subject-alternate-name-msupn\" - The username is extracted from the certificate's Other Name type of the Subject Alternative Name and must have the msUPN signature. \"uid\" - The username is extracted from the certificate's first instance of the User Identifier attribute in the Subject DN. \"uid-last\" - The username is extracted from the certificate's last instance of the User Identifier attribute in the Subject DN. </pre>  Available since 2.6.
	AuthenticationClientCertUsernameSource *string `json:"authenticationClientCertUsernameSource,omitempty"`
	// Enable or disable validation of the \"Not Before\" and \"Not After\" validity dates in the client certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	AuthenticationClientCertValidateDateEnabled *bool `json:"authenticationClientCertValidateDateEnabled,omitempty"`
	// Enable or disable allowing a client to specify a Client Username via the API connect method. When disabled, the Kerberos Principal name is always used. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	AuthenticationKerberosAllowApiProvidedUsernameEnabled *bool `json:"authenticationKerberosAllowApiProvidedUsernameEnabled,omitempty"`
	// Enable or disable Kerberos authentication in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	AuthenticationKerberosEnabled *bool `json:"authenticationKerberosEnabled,omitempty"`
	// The name of the profile to use when the client does not supply a profile name. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Available since 2.25.
	AuthenticationOauthDefaultProfileName *string `json:"authenticationOauthDefaultProfileName,omitempty"`
	// The name of the provider to use when the client does not supply a provider name. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`. Deprecated since 2.25. authenticationOauthDefaultProviderName and authenticationOauthProviders replaced by authenticationOauthDefaultProfileName and authenticationOauthProfiles.
	AuthenticationOauthDefaultProviderName *string `json:"authenticationOauthDefaultProviderName,omitempty"`
	// Enable or disable OAuth authentication. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.13.
	AuthenticationOauthEnabled *bool `json:"authenticationOauthEnabled,omitempty"`
	// The name of the attribute that is retrieved from the LDAP server as part of the LDAP search when authorizing a client connecting to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"memberOf\"`.
	AuthorizationLdapGroupMembershipAttributeName *string `json:"authorizationLdapGroupMembershipAttributeName,omitempty"`
	// Enable or disable client-username domain trimming for LDAP lookups of client connections. When enabled, the value of $CLIENT_USERNAME (when used for searching) will be truncated at the first occurance of the @ character. For example, if the client-username is in the form of an email address, then the domain portion will be removed. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.13.
	AuthorizationLdapTrimClientUsernameDomainEnabled *bool `json:"authorizationLdapTrimClientUsernameDomainEnabled,omitempty"`
	// The name of the LDAP Profile to use for client authorization. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.
	AuthorizationProfileName *string `json:"authorizationProfileName,omitempty"`
	// The type of authorization to use for clients connecting to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"internal\"`. The allowed values and their meaning are:  <pre> \"ldap\" - LDAP authorization. \"internal\" - Internal authorization. </pre>
	AuthorizationType *string `json:"authorizationType,omitempty"`
	// Enable or disable validation of the Common Name (CN) in the server certificate from the remote broker. If enabled, the Common Name is checked against the list of Trusted Common Names configured for the Bridge. Common Name validation is not performed if Server Certificate Name Validation is enabled, even if Common Name validation is enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Deprecated since 2.18. Common Name validation has been replaced by Server Certificate Name validation.
	BridgingTlsServerCertEnforceTrustedCommonNameEnabled *bool `json:"bridgingTlsServerCertEnforceTrustedCommonNameEnabled,omitempty"`
	// The maximum depth for a server certificate chain. The depth of a chain is defined as the number of signing CA certificates that are present in the chain back to a trusted self-signed root CA certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.
	BridgingTlsServerCertMaxChainDepth *int64 `json:"bridgingTlsServerCertMaxChainDepth,omitempty"`
	// Enable or disable validation of the \"Not Before\" and \"Not After\" validity dates in the server certificate. When disabled, a certificate will be accepted even if the certificate is not valid based on these dates. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	BridgingTlsServerCertValidateDateEnabled *bool `json:"bridgingTlsServerCertValidateDateEnabled,omitempty"`
	// Enable or disable the standard TLS authentication mechanism of verifying the name used to connect to the bridge. If enabled, the name used to connect to the bridge is checked against the names specified in the certificate returned by the remote router. Legacy Common Name validation is not performed if Server Certificate Name Validation is enabled, even if Common Name validation is also enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`. Available since 2.18.
	BridgingTlsServerCertValidateNameEnabled *bool `json:"bridgingTlsServerCertValidateNameEnabled,omitempty"`
	// Enable or disable managing of cache instances over the message bus. The default value is `true`.
	DistributedCacheManagementEnabled *bool `json:"distributedCacheManagementEnabled,omitempty"`
	// Enable or disable Dynamic Message Routing (DMR) for the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.11.
	DmrEnabled *bool `json:"dmrEnabled,omitempty"`
	// Enable or disable the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	Enabled                        *bool                  `json:"enabled,omitempty"`
	EventConnectionCountThreshold  *EventThreshold        `json:"eventConnectionCountThreshold,omitempty"`
	EventEgressFlowCountThreshold  *EventThreshold        `json:"eventEgressFlowCountThreshold,omitempty"`
	EventEgressMsgRateThreshold    *EventThresholdByValue `json:"eventEgressMsgRateThreshold,omitempty"`
	EventEndpointCountThreshold    *EventThreshold        `json:"eventEndpointCountThreshold,omitempty"`
	EventIngressFlowCountThreshold *EventThreshold        `json:"eventIngressFlowCountThreshold,omitempty"`
	EventIngressMsgRateThreshold   *EventThresholdByValue `json:"eventIngressMsgRateThreshold,omitempty"`
	// The threshold, in kilobytes, after which a message is considered to be large for the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1024`.
	EventLargeMsgThreshold *int64 `json:"eventLargeMsgThreshold,omitempty"`
	// A prefix applied to all published Events in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.
	EventLogTag                 *string         `json:"eventLogTag,omitempty"`
	EventMsgSpoolUsageThreshold *EventThreshold `json:"eventMsgSpoolUsageThreshold,omitempty"`
	// Enable or disable Client level Event message publishing. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	EventPublishClientEnabled *bool `json:"eventPublishClientEnabled,omitempty"`
	// Enable or disable Message VPN level Event message publishing. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	EventPublishMsgVpnEnabled *bool `json:"eventPublishMsgVpnEnabled,omitempty"`
	// Subscription level Event message publishing mode. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"off\"`. The allowed values and their meaning are:  <pre> \"off\" - Disable client level event message publishing. \"on-with-format-v1\" - Enable client level event message publishing with format v1. \"on-with-no-unsubscribe-events-on-disconnect-format-v1\" - As \"on-with-format-v1\", but unsubscribe events are not generated when a client disconnects. Unsubscribe events are still raised when a client explicitly unsubscribes from its subscriptions. \"on-with-format-v2\" - Enable client level event message publishing with format v2. \"on-with-no-unsubscribe-events-on-disconnect-format-v2\" - As \"on-with-format-v2\", but unsubscribe events are not generated when a client disconnects. Unsubscribe events are still raised when a client explicitly unsubscribes from its subscriptions. </pre>
	EventPublishSubscriptionMode *string `json:"eventPublishSubscriptionMode,omitempty"`
	// Enable or disable Event publish topics in MQTT format. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	EventPublishTopicFormatMqttEnabled *bool `json:"eventPublishTopicFormatMqttEnabled,omitempty"`
	// Enable or disable Event publish topics in SMF format. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	EventPublishTopicFormatSmfEnabled                *bool           `json:"eventPublishTopicFormatSmfEnabled,omitempty"`
	EventServiceAmqpConnectionCountThreshold         *EventThreshold `json:"eventServiceAmqpConnectionCountThreshold,omitempty"`
	EventServiceMqttConnectionCountThreshold         *EventThreshold `json:"eventServiceMqttConnectionCountThreshold,omitempty"`
	EventServiceRestIncomingConnectionCountThreshold *EventThreshold `json:"eventServiceRestIncomingConnectionCountThreshold,omitempty"`
	EventServiceSmfConnectionCountThreshold          *EventThreshold `json:"eventServiceSmfConnectionCountThreshold,omitempty"`
	EventServiceWebConnectionCountThreshold          *EventThreshold `json:"eventServiceWebConnectionCountThreshold,omitempty"`
	EventSubscriptionCountThreshold                  *EventThreshold `json:"eventSubscriptionCountThreshold,omitempty"`
	EventTransactedSessionCountThreshold             *EventThreshold `json:"eventTransactedSessionCountThreshold,omitempty"`
	EventTransactionCountThreshold                   *EventThreshold `json:"eventTransactionCountThreshold,omitempty"`
	// Enable or disable the export of subscriptions in the Message VPN to other routers in the network over Neighbor links. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	ExportSubscriptionsEnabled *bool `json:"exportSubscriptionsEnabled,omitempty"`
	// Enable or disable JNDI access for clients in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.2.
	JndiEnabled *bool `json:"jndiEnabled,omitempty"`
	// The maximum number of client connections to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform.
	MaxConnectionCount *int64 `json:"maxConnectionCount,omitempty"`
	// The maximum number of transmit flows that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.
	MaxEgressFlowCount *int64 `json:"maxEgressFlowCount,omitempty"`
	// The maximum number of Queues and Topic Endpoints that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.
	MaxEndpointCount *int64 `json:"maxEndpointCount,omitempty"`
	// The maximum number of receive flows that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.
	MaxIngressFlowCount *int64 `json:"maxIngressFlowCount,omitempty"`
	// The maximum message spool usage by the Message VPN, in megabytes. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.
	MaxMsgSpoolUsage *int64 `json:"maxMsgSpoolUsage,omitempty"`
	// The maximum number of local client subscriptions that can be added to the Message VPN. This limit is not enforced when a subscription is added using a management interface, such as CLI or SEMP. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.
	MaxSubscriptionCount *int64 `json:"maxSubscriptionCount,omitempty"`
	// The maximum number of transacted sessions that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.
	MaxTransactedSessionCount *int64 `json:"maxTransactedSessionCount,omitempty"`
	// The maximum number of transactions that can be created in the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.
	MaxTransactionCount *int64 `json:"maxTransactionCount,omitempty"`
	// The maximum total memory usage of the MQTT Retain feature for this Message VPN, in MB. If the maximum memory is reached, any arriving retain messages that require more memory are discarded. A value of -1 indicates that the memory is bounded only by the global max memory limit. A value of 0 prevents MQTT Retain from becoming operational. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `-1`. Available since 2.11.
	MqttRetainMaxMemory *int32 `json:"mqttRetainMaxMemory,omitempty"`
	// The name of the Message VPN.
	MsgVpnName *string `json:"msgVpnName,omitempty"`
	// The acknowledgement (ACK) propagation interval for the replication Bridge, in number of replicated messages. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `20`.
	ReplicationAckPropagationIntervalMsgCount *int64 `json:"replicationAckPropagationIntervalMsgCount,omitempty"`
	// The Client Username the replication Bridge uses to login to the remote Message VPN. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	ReplicationBridgeAuthenticationBasicClientUsername *string `json:"replicationBridgeAuthenticationBasicClientUsername,omitempty"`
	// The password for the Client Username. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"\"`.
	ReplicationBridgeAuthenticationBasicPassword *string `json:"replicationBridgeAuthenticationBasicPassword,omitempty"`
	// The PEM formatted content for the client certificate used by this bridge to login to the Remote Message VPN. It must consist of a private key and between one and three certificates comprising the certificate trust chain. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changing this attribute requires an HTTPS connection. The default value is `\"\"`. Available since 2.9.
	ReplicationBridgeAuthenticationClientCertContent *string `json:"replicationBridgeAuthenticationClientCertContent,omitempty"`
	// The password for the client certificate. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changing this attribute requires an HTTPS connection. The default value is `\"\"`. Available since 2.9.
	ReplicationBridgeAuthenticationClientCertPassword *string `json:"replicationBridgeAuthenticationClientCertPassword,omitempty"`
	// The authentication scheme for the replication Bridge in the Message VPN. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"basic\"`. The allowed values and their meaning are:  <pre> \"basic\" - Basic Authentication Scheme (via username and password). \"client-certificate\" - Client Certificate Authentication Scheme (via certificate file or content). </pre>
	ReplicationBridgeAuthenticationScheme *string `json:"replicationBridgeAuthenticationScheme,omitempty"`
	// Enable or disable use of compression for the replication Bridge. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.
	ReplicationBridgeCompressedDataEnabled *bool `json:"replicationBridgeCompressedDataEnabled,omitempty"`
	// The size of the window used for guaranteed messages published to the replication Bridge, in messages. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `255`.
	ReplicationBridgeEgressFlowWindowSize *int64 `json:"replicationBridgeEgressFlowWindowSize,omitempty"`
	// The number of seconds that must pass before retrying the replication Bridge connection. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `3`.
	ReplicationBridgeRetryDelay *int64 `json:"replicationBridgeRetryDelay,omitempty"`
	// Enable or disable use of encryption (TLS) for the replication Bridge connection. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.
	ReplicationBridgeTlsEnabled *bool `json:"replicationBridgeTlsEnabled,omitempty"`
	// The Client Profile for the unidirectional replication Bridge in the Message VPN. It is used only for the TCP parameters. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"#client-profile\"`.
	ReplicationBridgeUnidirectionalClientProfileName *string `json:"replicationBridgeUnidirectionalClientProfileName,omitempty"`
	// Enable or disable replication for the Message VPN. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `false`.
	ReplicationEnabled *bool `json:"replicationEnabled,omitempty"`
	// The behavior to take when enabling replication for the Message VPN, depending on the existence of the replication Queue. This attribute is absent from a GET and not updated when absent in a PUT, subject to the exceptions in note 4. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"fail-on-existing-queue\"`. The allowed values and their meaning are:  <pre> \"fail-on-existing-queue\" - The data replication queue must not already exist. \"force-use-existing-queue\" - The data replication queue must already exist. Any data messages on the Queue will be forwarded to interested applications. IMPORTANT: Before using this mode be certain that the messages are not stale or otherwise unsuitable to be forwarded. This mode can only be specified when the existing queue is configured the same as is currently specified under replication configuration otherwise the enabling of replication will fail. \"force-recreate-queue\" - The data replication queue must already exist. Any data messages on the Queue will be discarded. IMPORTANT: Before using this mode be certain that the messages on the existing data replication queue are not needed by interested applications. </pre>
	ReplicationEnabledQueueBehavior *string `json:"replicationEnabledQueueBehavior,omitempty"`
	// The maximum message spool usage by the replication Bridge local Queue (quota), in megabytes. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `60000`.
	ReplicationQueueMaxMsgSpoolUsage *int64 `json:"replicationQueueMaxMsgSpoolUsage,omitempty"`
	// Enable or disable whether messages discarded on the replication Bridge local Queue are rejected back to the sender. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	ReplicationQueueRejectMsgToSenderOnDiscardEnabled *bool `json:"replicationQueueRejectMsgToSenderOnDiscardEnabled,omitempty"`
	// Enable or disable whether guaranteed messages published to synchronously replicated Topics are rejected back to the sender when synchronous replication becomes ineligible. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	ReplicationRejectMsgWhenSyncIneligibleEnabled *bool `json:"replicationRejectMsgWhenSyncIneligibleEnabled,omitempty"`
	// The replication role for the Message VPN. Changes to this attribute are synchronized to HA mates via config-sync. The default value is `\"standby\"`. The allowed values and their meaning are:  <pre> \"active\" - Assume the Active role in replication for the Message VPN. \"standby\" - Assume the Standby role in replication for the Message VPN. </pre>
	ReplicationRole *string `json:"replicationRole,omitempty"`
	// The transaction replication mode for all transactions within the Message VPN. Changing this value during operation will not affect existing transactions; it is only used upon starting a transaction. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"async\"`. The allowed values and their meaning are:  <pre> \"sync\" - Messages are acknowledged when replicated (spooled remotely). \"async\" - Messages are acknowledged when pending replication (spooled locally). </pre>
	ReplicationTransactionMode *string `json:"replicationTransactionMode,omitempty"`
	// Enable or disable validation of the Common Name (CN) in the server certificate from the remote REST Consumer. If enabled, the Common Name is checked against the list of Trusted Common Names configured for the REST Consumer. Common Name validation is not performed if Server Certificate Name Validation is enabled, even if Common Name validation is enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Deprecated since 2.17. Common Name validation has been replaced by Server Certificate Name validation.
	RestTlsServerCertEnforceTrustedCommonNameEnabled *bool `json:"restTlsServerCertEnforceTrustedCommonNameEnabled,omitempty"`
	// The maximum depth for a REST Consumer server certificate chain. The depth of a chain is defined as the number of signing CA certificates that are present in the chain back to a trusted self-signed root CA certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `3`.
	RestTlsServerCertMaxChainDepth *int64 `json:"restTlsServerCertMaxChainDepth,omitempty"`
	// Enable or disable validation of the \"Not Before\" and \"Not After\" validity dates in the REST Consumer server certificate. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	RestTlsServerCertValidateDateEnabled *bool `json:"restTlsServerCertValidateDateEnabled,omitempty"`
	// Enable or disable the standard TLS authentication mechanism of verifying the name used to connect to the remote REST Consumer. If enabled, the name used to connect to the remote REST Consumer is checked against the names specified in the certificate returned by the remote router. Legacy Common Name validation is not performed if Server Certificate Name Validation is enabled, even if Common Name validation is also enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`. Available since 2.17.
	RestTlsServerCertValidateNameEnabled *bool `json:"restTlsServerCertValidateNameEnabled,omitempty"`
	// Enable or disable \"admin client\" SEMP over the message bus commands for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	SempOverMsgBusAdminClientEnabled *bool `json:"sempOverMsgBusAdminClientEnabled,omitempty"`
	// Enable or disable \"admin distributed-cache\" SEMP over the message bus commands for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	SempOverMsgBusAdminDistributedCacheEnabled *bool `json:"sempOverMsgBusAdminDistributedCacheEnabled,omitempty"`
	// Enable or disable \"admin\" SEMP over the message bus commands for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	SempOverMsgBusAdminEnabled *bool `json:"sempOverMsgBusAdminEnabled,omitempty"`
	// Enable or disable SEMP over the message bus for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	SempOverMsgBusEnabled *bool `json:"sempOverMsgBusEnabled,omitempty"`
	// Enable or disable \"show\" SEMP over the message bus commands for the current Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	SempOverMsgBusShowEnabled *bool `json:"sempOverMsgBusShowEnabled,omitempty"`
	// The maximum number of AMQP client connections that can be simultaneously connected to the Message VPN. This value may be higher than supported by the platform. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform. Available since 2.7.
	ServiceAmqpMaxConnectionCount *int64 `json:"serviceAmqpMaxConnectionCount,omitempty"`
	// Enable or disable the plain-text AMQP service in the Message VPN. Disabling causes clients connected to the corresponding listen-port to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.7.
	ServiceAmqpPlainTextEnabled *bool `json:"serviceAmqpPlainTextEnabled,omitempty"`
	// The port number for plain-text AMQP clients that connect to the Message VPN. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.7.
	ServiceAmqpPlainTextListenPort *int64 `json:"serviceAmqpPlainTextListenPort,omitempty"`
	// Enable or disable the use of encryption (TLS) for the AMQP service in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.7.
	ServiceAmqpTlsEnabled *bool `json:"serviceAmqpTlsEnabled,omitempty"`
	// The port number for AMQP clients that connect to the Message VPN over TLS. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.7.
	ServiceAmqpTlsListenPort *int64 `json:"serviceAmqpTlsListenPort,omitempty"`
	// Determines when to request a client certificate from an incoming MQTT client connecting via a TLS port. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"when-enabled-in-message-vpn\"`. The allowed values and their meaning are:  <pre> \"always\" - Always ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"never\" - Never ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"when-enabled-in-message-vpn\" - Only ask for a client-certificate if client certificate authentication is enabled under \"message-vpn >  authentication > client-certificate > shutdown\". </pre>  Available since 2.21.
	ServiceMqttAuthenticationClientCertRequest *string `json:"serviceMqttAuthenticationClientCertRequest,omitempty"`
	// The maximum number of MQTT client connections that can be simultaneously connected to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform. Available since 2.1.
	ServiceMqttMaxConnectionCount *int64 `json:"serviceMqttMaxConnectionCount,omitempty"`
	// Enable or disable the plain-text MQTT service in the Message VPN. Disabling causes clients currently connected to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.1.
	ServiceMqttPlainTextEnabled *bool `json:"serviceMqttPlainTextEnabled,omitempty"`
	// The port number for plain-text MQTT clients that connect to the Message VPN. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.1.
	ServiceMqttPlainTextListenPort *int64 `json:"serviceMqttPlainTextListenPort,omitempty"`
	// Enable or disable the use of encryption (TLS) for the MQTT service in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.1.
	ServiceMqttTlsEnabled *bool `json:"serviceMqttTlsEnabled,omitempty"`
	// The port number for MQTT clients that connect to the Message VPN over TLS. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.1.
	ServiceMqttTlsListenPort *int64 `json:"serviceMqttTlsListenPort,omitempty"`
	// Enable or disable the use of encrypted WebSocket (WebSocket over TLS) for the MQTT service in the Message VPN. Disabling causes clients currently connected by encrypted WebSocket to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.1.
	ServiceMqttTlsWebSocketEnabled *bool `json:"serviceMqttTlsWebSocketEnabled,omitempty"`
	// The port number for MQTT clients that connect to the Message VPN using WebSocket over TLS. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.1.
	ServiceMqttTlsWebSocketListenPort *int64 `json:"serviceMqttTlsWebSocketListenPort,omitempty"`
	// Enable or disable the use of WebSocket for the MQTT service in the Message VPN. Disabling causes clients currently connected by WebSocket to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.1.
	ServiceMqttWebSocketEnabled *bool `json:"serviceMqttWebSocketEnabled,omitempty"`
	// The port number for plain-text MQTT clients that connect to the Message VPN using WebSocket. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.1.
	ServiceMqttWebSocketListenPort *int64 `json:"serviceMqttWebSocketListenPort,omitempty"`
	// Determines when to request a client certificate from an incoming REST Producer connecting via a TLS port. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"when-enabled-in-message-vpn\"`. The allowed values and their meaning are:  <pre> \"always\" - Always ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"never\" - Never ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"when-enabled-in-message-vpn\" - Only ask for a client-certificate if client certificate authentication is enabled under \"message-vpn >  authentication > client-certificate > shutdown\". </pre>  Available since 2.21.
	ServiceRestIncomingAuthenticationClientCertRequest *string `json:"serviceRestIncomingAuthenticationClientCertRequest,omitempty"`
	// The handling of Authorization headers for incoming REST connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"drop\"`. The allowed values and their meaning are:  <pre> \"drop\" - Do not attach the Authorization header to the message as a user property. This configuration is most secure. \"forward\" - Forward the Authorization header, attaching it to the message as a user property in the same way as other headers. For best security, use the drop setting. \"legacy\" - If the Authorization header was used for authentication to the broker, do not attach it to the message. If the Authorization header was not used for authentication to the broker, attach it to the message as a user property in the same way as other headers. For best security, use the drop setting. </pre>  Available since 2.19.
	ServiceRestIncomingAuthorizationHeaderHandling *string `json:"serviceRestIncomingAuthorizationHeaderHandling,omitempty"`
	// The maximum number of REST incoming client connections that can be simultaneously connected to the Message VPN. This value may be higher than supported by the platform. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform.
	ServiceRestIncomingMaxConnectionCount *int64 `json:"serviceRestIncomingMaxConnectionCount,omitempty"`
	// Enable or disable the plain-text REST service for incoming clients in the Message VPN. Disabling causes clients currently connected to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	ServiceRestIncomingPlainTextEnabled *bool `json:"serviceRestIncomingPlainTextEnabled,omitempty"`
	// The port number for incoming plain-text REST clients that connect to the Message VPN. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.
	ServiceRestIncomingPlainTextListenPort *int64 `json:"serviceRestIncomingPlainTextListenPort,omitempty"`
	// Enable or disable the use of encryption (TLS) for the REST service for incoming clients in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	ServiceRestIncomingTlsEnabled *bool `json:"serviceRestIncomingTlsEnabled,omitempty"`
	// The port number for incoming REST clients that connect to the Message VPN over TLS. The port must be unique across the message backbone. A value of 0 means that the listen-port is unassigned and cannot be enabled. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.
	ServiceRestIncomingTlsListenPort *int64 `json:"serviceRestIncomingTlsListenPort,omitempty"`
	// The REST service mode for incoming REST clients that connect to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"messaging\"`. The allowed values and their meaning are:  <pre> \"gateway\" - Act as a message gateway through which REST messages are propagated. \"messaging\" - Act as a message broker on which REST messages are queued. </pre>  Available since 2.6.
	ServiceRestMode *string `json:"serviceRestMode,omitempty"`
	// The maximum number of REST Consumer (outgoing) client connections that can be simultaneously connected to the Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.
	ServiceRestOutgoingMaxConnectionCount *int64 `json:"serviceRestOutgoingMaxConnectionCount,omitempty"`
	// The maximum number of SMF client connections that can be simultaneously connected to the Message VPN. This value may be higher than supported by the platform. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default varies by platform.
	ServiceSmfMaxConnectionCount *int64 `json:"serviceSmfMaxConnectionCount,omitempty"`
	// Enable or disable the plain-text SMF service in the Message VPN. Disabling causes clients currently connected to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	ServiceSmfPlainTextEnabled *bool `json:"serviceSmfPlainTextEnabled,omitempty"`
	// Enable or disable the use of encryption (TLS) for the SMF service in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	ServiceSmfTlsEnabled *bool `json:"serviceSmfTlsEnabled,omitempty"`
	// Determines when to request a client certificate from a Web Transport client connecting via a TLS port. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"when-enabled-in-message-vpn\"`. The allowed values and their meaning are:  <pre> \"always\" - Always ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"never\" - Never ask for a client certificate regardless of the \"message-vpn > authentication > client-certificate > shutdown\" configuration. \"when-enabled-in-message-vpn\" - Only ask for a client-certificate if client certificate authentication is enabled under \"message-vpn >  authentication > client-certificate > shutdown\". </pre>  Available since 2.21.
	ServiceWebAuthenticationClientCertRequest *string `json:"serviceWebAuthenticationClientCertRequest,omitempty"`
	// The maximum number of Web Transport client connections that can be simultaneously connected to the Message VPN. This value may be higher than supported by the platform. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default is the maximum value supported by the platform.
	ServiceWebMaxConnectionCount *int64 `json:"serviceWebMaxConnectionCount,omitempty"`
	// Enable or disable the plain-text Web Transport service in the Message VPN. Disabling causes clients currently connected to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	ServiceWebPlainTextEnabled *bool `json:"serviceWebPlainTextEnabled,omitempty"`
	// Enable or disable the use of TLS for the Web Transport service in the Message VPN. Disabling causes clients currently connected over TLS to be disconnected. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.
	ServiceWebTlsEnabled *bool `json:"serviceWebTlsEnabled,omitempty"`
	// Enable or disable the allowing of TLS SMF clients to downgrade their connections to plain-text connections. Changing this will not affect existing connections. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.
	TlsAllowDowngradeToPlainTextEnabled *bool `json:"tlsAllowDowngradeToPlainTextEnabled,omitempty"`
}

// NewMsgVpn instantiates a new MsgVpn object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMsgVpn() *MsgVpn {
	this := MsgVpn{}
	return &this
}

// NewMsgVpnWithDefaults instantiates a new MsgVpn object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMsgVpnWithDefaults() *MsgVpn {
	this := MsgVpn{}
	return &this
}

// GetAlias returns the Alias field value if set, zero value otherwise.
func (o *MsgVpn) GetAlias() string {
	if o == nil || o.Alias == nil {
		var ret string
		return ret
	}
	return *o.Alias
}

// GetAliasOk returns a tuple with the Alias field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAliasOk() (*string, bool) {
	if o == nil || o.Alias == nil {
		return nil, false
	}
	return o.Alias, true
}

// HasAlias returns a boolean if a field has been set.
func (o *MsgVpn) HasAlias() bool {
	if o != nil && o.Alias != nil {
		return true
	}

	return false
}

// SetAlias gets a reference to the given string and assigns it to the Alias field.
func (o *MsgVpn) SetAlias(v string) {
	o.Alias = &v
}

// GetAuthenticationBasicEnabled returns the AuthenticationBasicEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationBasicEnabled() bool {
	if o == nil || o.AuthenticationBasicEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AuthenticationBasicEnabled
}

// GetAuthenticationBasicEnabledOk returns a tuple with the AuthenticationBasicEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationBasicEnabledOk() (*bool, bool) {
	if o == nil || o.AuthenticationBasicEnabled == nil {
		return nil, false
	}
	return o.AuthenticationBasicEnabled, true
}

// HasAuthenticationBasicEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationBasicEnabled() bool {
	if o != nil && o.AuthenticationBasicEnabled != nil {
		return true
	}

	return false
}

// SetAuthenticationBasicEnabled gets a reference to the given bool and assigns it to the AuthenticationBasicEnabled field.
func (o *MsgVpn) SetAuthenticationBasicEnabled(v bool) {
	o.AuthenticationBasicEnabled = &v
}

// GetAuthenticationBasicProfileName returns the AuthenticationBasicProfileName field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationBasicProfileName() string {
	if o == nil || o.AuthenticationBasicProfileName == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationBasicProfileName
}

// GetAuthenticationBasicProfileNameOk returns a tuple with the AuthenticationBasicProfileName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationBasicProfileNameOk() (*string, bool) {
	if o == nil || o.AuthenticationBasicProfileName == nil {
		return nil, false
	}
	return o.AuthenticationBasicProfileName, true
}

// HasAuthenticationBasicProfileName returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationBasicProfileName() bool {
	if o != nil && o.AuthenticationBasicProfileName != nil {
		return true
	}

	return false
}

// SetAuthenticationBasicProfileName gets a reference to the given string and assigns it to the AuthenticationBasicProfileName field.
func (o *MsgVpn) SetAuthenticationBasicProfileName(v string) {
	o.AuthenticationBasicProfileName = &v
}

// GetAuthenticationBasicRadiusDomain returns the AuthenticationBasicRadiusDomain field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationBasicRadiusDomain() string {
	if o == nil || o.AuthenticationBasicRadiusDomain == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationBasicRadiusDomain
}

// GetAuthenticationBasicRadiusDomainOk returns a tuple with the AuthenticationBasicRadiusDomain field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationBasicRadiusDomainOk() (*string, bool) {
	if o == nil || o.AuthenticationBasicRadiusDomain == nil {
		return nil, false
	}
	return o.AuthenticationBasicRadiusDomain, true
}

// HasAuthenticationBasicRadiusDomain returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationBasicRadiusDomain() bool {
	if o != nil && o.AuthenticationBasicRadiusDomain != nil {
		return true
	}

	return false
}

// SetAuthenticationBasicRadiusDomain gets a reference to the given string and assigns it to the AuthenticationBasicRadiusDomain field.
func (o *MsgVpn) SetAuthenticationBasicRadiusDomain(v string) {
	o.AuthenticationBasicRadiusDomain = &v
}

// GetAuthenticationBasicType returns the AuthenticationBasicType field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationBasicType() string {
	if o == nil || o.AuthenticationBasicType == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationBasicType
}

// GetAuthenticationBasicTypeOk returns a tuple with the AuthenticationBasicType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationBasicTypeOk() (*string, bool) {
	if o == nil || o.AuthenticationBasicType == nil {
		return nil, false
	}
	return o.AuthenticationBasicType, true
}

// HasAuthenticationBasicType returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationBasicType() bool {
	if o != nil && o.AuthenticationBasicType != nil {
		return true
	}

	return false
}

// SetAuthenticationBasicType gets a reference to the given string and assigns it to the AuthenticationBasicType field.
func (o *MsgVpn) SetAuthenticationBasicType(v string) {
	o.AuthenticationBasicType = &v
}

// GetAuthenticationClientCertAllowApiProvidedUsernameEnabled returns the AuthenticationClientCertAllowApiProvidedUsernameEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationClientCertAllowApiProvidedUsernameEnabled() bool {
	if o == nil || o.AuthenticationClientCertAllowApiProvidedUsernameEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AuthenticationClientCertAllowApiProvidedUsernameEnabled
}

// GetAuthenticationClientCertAllowApiProvidedUsernameEnabledOk returns a tuple with the AuthenticationClientCertAllowApiProvidedUsernameEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationClientCertAllowApiProvidedUsernameEnabledOk() (*bool, bool) {
	if o == nil || o.AuthenticationClientCertAllowApiProvidedUsernameEnabled == nil {
		return nil, false
	}
	return o.AuthenticationClientCertAllowApiProvidedUsernameEnabled, true
}

// HasAuthenticationClientCertAllowApiProvidedUsernameEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationClientCertAllowApiProvidedUsernameEnabled() bool {
	if o != nil && o.AuthenticationClientCertAllowApiProvidedUsernameEnabled != nil {
		return true
	}

	return false
}

// SetAuthenticationClientCertAllowApiProvidedUsernameEnabled gets a reference to the given bool and assigns it to the AuthenticationClientCertAllowApiProvidedUsernameEnabled field.
func (o *MsgVpn) SetAuthenticationClientCertAllowApiProvidedUsernameEnabled(v bool) {
	o.AuthenticationClientCertAllowApiProvidedUsernameEnabled = &v
}

// GetAuthenticationClientCertEnabled returns the AuthenticationClientCertEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationClientCertEnabled() bool {
	if o == nil || o.AuthenticationClientCertEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AuthenticationClientCertEnabled
}

// GetAuthenticationClientCertEnabledOk returns a tuple with the AuthenticationClientCertEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationClientCertEnabledOk() (*bool, bool) {
	if o == nil || o.AuthenticationClientCertEnabled == nil {
		return nil, false
	}
	return o.AuthenticationClientCertEnabled, true
}

// HasAuthenticationClientCertEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationClientCertEnabled() bool {
	if o != nil && o.AuthenticationClientCertEnabled != nil {
		return true
	}

	return false
}

// SetAuthenticationClientCertEnabled gets a reference to the given bool and assigns it to the AuthenticationClientCertEnabled field.
func (o *MsgVpn) SetAuthenticationClientCertEnabled(v bool) {
	o.AuthenticationClientCertEnabled = &v
}

// GetAuthenticationClientCertMaxChainDepth returns the AuthenticationClientCertMaxChainDepth field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationClientCertMaxChainDepth() int64 {
	if o == nil || o.AuthenticationClientCertMaxChainDepth == nil {
		var ret int64
		return ret
	}
	return *o.AuthenticationClientCertMaxChainDepth
}

// GetAuthenticationClientCertMaxChainDepthOk returns a tuple with the AuthenticationClientCertMaxChainDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationClientCertMaxChainDepthOk() (*int64, bool) {
	if o == nil || o.AuthenticationClientCertMaxChainDepth == nil {
		return nil, false
	}
	return o.AuthenticationClientCertMaxChainDepth, true
}

// HasAuthenticationClientCertMaxChainDepth returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationClientCertMaxChainDepth() bool {
	if o != nil && o.AuthenticationClientCertMaxChainDepth != nil {
		return true
	}

	return false
}

// SetAuthenticationClientCertMaxChainDepth gets a reference to the given int64 and assigns it to the AuthenticationClientCertMaxChainDepth field.
func (o *MsgVpn) SetAuthenticationClientCertMaxChainDepth(v int64) {
	o.AuthenticationClientCertMaxChainDepth = &v
}

// GetAuthenticationClientCertRevocationCheckMode returns the AuthenticationClientCertRevocationCheckMode field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationClientCertRevocationCheckMode() string {
	if o == nil || o.AuthenticationClientCertRevocationCheckMode == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationClientCertRevocationCheckMode
}

// GetAuthenticationClientCertRevocationCheckModeOk returns a tuple with the AuthenticationClientCertRevocationCheckMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationClientCertRevocationCheckModeOk() (*string, bool) {
	if o == nil || o.AuthenticationClientCertRevocationCheckMode == nil {
		return nil, false
	}
	return o.AuthenticationClientCertRevocationCheckMode, true
}

// HasAuthenticationClientCertRevocationCheckMode returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationClientCertRevocationCheckMode() bool {
	if o != nil && o.AuthenticationClientCertRevocationCheckMode != nil {
		return true
	}

	return false
}

// SetAuthenticationClientCertRevocationCheckMode gets a reference to the given string and assigns it to the AuthenticationClientCertRevocationCheckMode field.
func (o *MsgVpn) SetAuthenticationClientCertRevocationCheckMode(v string) {
	o.AuthenticationClientCertRevocationCheckMode = &v
}

// GetAuthenticationClientCertUsernameSource returns the AuthenticationClientCertUsernameSource field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationClientCertUsernameSource() string {
	if o == nil || o.AuthenticationClientCertUsernameSource == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationClientCertUsernameSource
}

// GetAuthenticationClientCertUsernameSourceOk returns a tuple with the AuthenticationClientCertUsernameSource field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationClientCertUsernameSourceOk() (*string, bool) {
	if o == nil || o.AuthenticationClientCertUsernameSource == nil {
		return nil, false
	}
	return o.AuthenticationClientCertUsernameSource, true
}

// HasAuthenticationClientCertUsernameSource returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationClientCertUsernameSource() bool {
	if o != nil && o.AuthenticationClientCertUsernameSource != nil {
		return true
	}

	return false
}

// SetAuthenticationClientCertUsernameSource gets a reference to the given string and assigns it to the AuthenticationClientCertUsernameSource field.
func (o *MsgVpn) SetAuthenticationClientCertUsernameSource(v string) {
	o.AuthenticationClientCertUsernameSource = &v
}

// GetAuthenticationClientCertValidateDateEnabled returns the AuthenticationClientCertValidateDateEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationClientCertValidateDateEnabled() bool {
	if o == nil || o.AuthenticationClientCertValidateDateEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AuthenticationClientCertValidateDateEnabled
}

// GetAuthenticationClientCertValidateDateEnabledOk returns a tuple with the AuthenticationClientCertValidateDateEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationClientCertValidateDateEnabledOk() (*bool, bool) {
	if o == nil || o.AuthenticationClientCertValidateDateEnabled == nil {
		return nil, false
	}
	return o.AuthenticationClientCertValidateDateEnabled, true
}

// HasAuthenticationClientCertValidateDateEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationClientCertValidateDateEnabled() bool {
	if o != nil && o.AuthenticationClientCertValidateDateEnabled != nil {
		return true
	}

	return false
}

// SetAuthenticationClientCertValidateDateEnabled gets a reference to the given bool and assigns it to the AuthenticationClientCertValidateDateEnabled field.
func (o *MsgVpn) SetAuthenticationClientCertValidateDateEnabled(v bool) {
	o.AuthenticationClientCertValidateDateEnabled = &v
}

// GetAuthenticationKerberosAllowApiProvidedUsernameEnabled returns the AuthenticationKerberosAllowApiProvidedUsernameEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationKerberosAllowApiProvidedUsernameEnabled() bool {
	if o == nil || o.AuthenticationKerberosAllowApiProvidedUsernameEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AuthenticationKerberosAllowApiProvidedUsernameEnabled
}

// GetAuthenticationKerberosAllowApiProvidedUsernameEnabledOk returns a tuple with the AuthenticationKerberosAllowApiProvidedUsernameEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationKerberosAllowApiProvidedUsernameEnabledOk() (*bool, bool) {
	if o == nil || o.AuthenticationKerberosAllowApiProvidedUsernameEnabled == nil {
		return nil, false
	}
	return o.AuthenticationKerberosAllowApiProvidedUsernameEnabled, true
}

// HasAuthenticationKerberosAllowApiProvidedUsernameEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationKerberosAllowApiProvidedUsernameEnabled() bool {
	if o != nil && o.AuthenticationKerberosAllowApiProvidedUsernameEnabled != nil {
		return true
	}

	return false
}

// SetAuthenticationKerberosAllowApiProvidedUsernameEnabled gets a reference to the given bool and assigns it to the AuthenticationKerberosAllowApiProvidedUsernameEnabled field.
func (o *MsgVpn) SetAuthenticationKerberosAllowApiProvidedUsernameEnabled(v bool) {
	o.AuthenticationKerberosAllowApiProvidedUsernameEnabled = &v
}

// GetAuthenticationKerberosEnabled returns the AuthenticationKerberosEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationKerberosEnabled() bool {
	if o == nil || o.AuthenticationKerberosEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AuthenticationKerberosEnabled
}

// GetAuthenticationKerberosEnabledOk returns a tuple with the AuthenticationKerberosEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationKerberosEnabledOk() (*bool, bool) {
	if o == nil || o.AuthenticationKerberosEnabled == nil {
		return nil, false
	}
	return o.AuthenticationKerberosEnabled, true
}

// HasAuthenticationKerberosEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationKerberosEnabled() bool {
	if o != nil && o.AuthenticationKerberosEnabled != nil {
		return true
	}

	return false
}

// SetAuthenticationKerberosEnabled gets a reference to the given bool and assigns it to the AuthenticationKerberosEnabled field.
func (o *MsgVpn) SetAuthenticationKerberosEnabled(v bool) {
	o.AuthenticationKerberosEnabled = &v
}

// GetAuthenticationOauthDefaultProfileName returns the AuthenticationOauthDefaultProfileName field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationOauthDefaultProfileName() string {
	if o == nil || o.AuthenticationOauthDefaultProfileName == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationOauthDefaultProfileName
}

// GetAuthenticationOauthDefaultProfileNameOk returns a tuple with the AuthenticationOauthDefaultProfileName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationOauthDefaultProfileNameOk() (*string, bool) {
	if o == nil || o.AuthenticationOauthDefaultProfileName == nil {
		return nil, false
	}
	return o.AuthenticationOauthDefaultProfileName, true
}

// HasAuthenticationOauthDefaultProfileName returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationOauthDefaultProfileName() bool {
	if o != nil && o.AuthenticationOauthDefaultProfileName != nil {
		return true
	}

	return false
}

// SetAuthenticationOauthDefaultProfileName gets a reference to the given string and assigns it to the AuthenticationOauthDefaultProfileName field.
func (o *MsgVpn) SetAuthenticationOauthDefaultProfileName(v string) {
	o.AuthenticationOauthDefaultProfileName = &v
}

// GetAuthenticationOauthDefaultProviderName returns the AuthenticationOauthDefaultProviderName field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationOauthDefaultProviderName() string {
	if o == nil || o.AuthenticationOauthDefaultProviderName == nil {
		var ret string
		return ret
	}
	return *o.AuthenticationOauthDefaultProviderName
}

// GetAuthenticationOauthDefaultProviderNameOk returns a tuple with the AuthenticationOauthDefaultProviderName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationOauthDefaultProviderNameOk() (*string, bool) {
	if o == nil || o.AuthenticationOauthDefaultProviderName == nil {
		return nil, false
	}
	return o.AuthenticationOauthDefaultProviderName, true
}

// HasAuthenticationOauthDefaultProviderName returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationOauthDefaultProviderName() bool {
	if o != nil && o.AuthenticationOauthDefaultProviderName != nil {
		return true
	}

	return false
}

// SetAuthenticationOauthDefaultProviderName gets a reference to the given string and assigns it to the AuthenticationOauthDefaultProviderName field.
func (o *MsgVpn) SetAuthenticationOauthDefaultProviderName(v string) {
	o.AuthenticationOauthDefaultProviderName = &v
}

// GetAuthenticationOauthEnabled returns the AuthenticationOauthEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthenticationOauthEnabled() bool {
	if o == nil || o.AuthenticationOauthEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AuthenticationOauthEnabled
}

// GetAuthenticationOauthEnabledOk returns a tuple with the AuthenticationOauthEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthenticationOauthEnabledOk() (*bool, bool) {
	if o == nil || o.AuthenticationOauthEnabled == nil {
		return nil, false
	}
	return o.AuthenticationOauthEnabled, true
}

// HasAuthenticationOauthEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthenticationOauthEnabled() bool {
	if o != nil && o.AuthenticationOauthEnabled != nil {
		return true
	}

	return false
}

// SetAuthenticationOauthEnabled gets a reference to the given bool and assigns it to the AuthenticationOauthEnabled field.
func (o *MsgVpn) SetAuthenticationOauthEnabled(v bool) {
	o.AuthenticationOauthEnabled = &v
}

// GetAuthorizationLdapGroupMembershipAttributeName returns the AuthorizationLdapGroupMembershipAttributeName field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthorizationLdapGroupMembershipAttributeName() string {
	if o == nil || o.AuthorizationLdapGroupMembershipAttributeName == nil {
		var ret string
		return ret
	}
	return *o.AuthorizationLdapGroupMembershipAttributeName
}

// GetAuthorizationLdapGroupMembershipAttributeNameOk returns a tuple with the AuthorizationLdapGroupMembershipAttributeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthorizationLdapGroupMembershipAttributeNameOk() (*string, bool) {
	if o == nil || o.AuthorizationLdapGroupMembershipAttributeName == nil {
		return nil, false
	}
	return o.AuthorizationLdapGroupMembershipAttributeName, true
}

// HasAuthorizationLdapGroupMembershipAttributeName returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthorizationLdapGroupMembershipAttributeName() bool {
	if o != nil && o.AuthorizationLdapGroupMembershipAttributeName != nil {
		return true
	}

	return false
}

// SetAuthorizationLdapGroupMembershipAttributeName gets a reference to the given string and assigns it to the AuthorizationLdapGroupMembershipAttributeName field.
func (o *MsgVpn) SetAuthorizationLdapGroupMembershipAttributeName(v string) {
	o.AuthorizationLdapGroupMembershipAttributeName = &v
}

// GetAuthorizationLdapTrimClientUsernameDomainEnabled returns the AuthorizationLdapTrimClientUsernameDomainEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthorizationLdapTrimClientUsernameDomainEnabled() bool {
	if o == nil || o.AuthorizationLdapTrimClientUsernameDomainEnabled == nil {
		var ret bool
		return ret
	}
	return *o.AuthorizationLdapTrimClientUsernameDomainEnabled
}

// GetAuthorizationLdapTrimClientUsernameDomainEnabledOk returns a tuple with the AuthorizationLdapTrimClientUsernameDomainEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthorizationLdapTrimClientUsernameDomainEnabledOk() (*bool, bool) {
	if o == nil || o.AuthorizationLdapTrimClientUsernameDomainEnabled == nil {
		return nil, false
	}
	return o.AuthorizationLdapTrimClientUsernameDomainEnabled, true
}

// HasAuthorizationLdapTrimClientUsernameDomainEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthorizationLdapTrimClientUsernameDomainEnabled() bool {
	if o != nil && o.AuthorizationLdapTrimClientUsernameDomainEnabled != nil {
		return true
	}

	return false
}

// SetAuthorizationLdapTrimClientUsernameDomainEnabled gets a reference to the given bool and assigns it to the AuthorizationLdapTrimClientUsernameDomainEnabled field.
func (o *MsgVpn) SetAuthorizationLdapTrimClientUsernameDomainEnabled(v bool) {
	o.AuthorizationLdapTrimClientUsernameDomainEnabled = &v
}

// GetAuthorizationProfileName returns the AuthorizationProfileName field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthorizationProfileName() string {
	if o == nil || o.AuthorizationProfileName == nil {
		var ret string
		return ret
	}
	return *o.AuthorizationProfileName
}

// GetAuthorizationProfileNameOk returns a tuple with the AuthorizationProfileName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthorizationProfileNameOk() (*string, bool) {
	if o == nil || o.AuthorizationProfileName == nil {
		return nil, false
	}
	return o.AuthorizationProfileName, true
}

// HasAuthorizationProfileName returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthorizationProfileName() bool {
	if o != nil && o.AuthorizationProfileName != nil {
		return true
	}

	return false
}

// SetAuthorizationProfileName gets a reference to the given string and assigns it to the AuthorizationProfileName field.
func (o *MsgVpn) SetAuthorizationProfileName(v string) {
	o.AuthorizationProfileName = &v
}

// GetAuthorizationType returns the AuthorizationType field value if set, zero value otherwise.
func (o *MsgVpn) GetAuthorizationType() string {
	if o == nil || o.AuthorizationType == nil {
		var ret string
		return ret
	}
	return *o.AuthorizationType
}

// GetAuthorizationTypeOk returns a tuple with the AuthorizationType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetAuthorizationTypeOk() (*string, bool) {
	if o == nil || o.AuthorizationType == nil {
		return nil, false
	}
	return o.AuthorizationType, true
}

// HasAuthorizationType returns a boolean if a field has been set.
func (o *MsgVpn) HasAuthorizationType() bool {
	if o != nil && o.AuthorizationType != nil {
		return true
	}

	return false
}

// SetAuthorizationType gets a reference to the given string and assigns it to the AuthorizationType field.
func (o *MsgVpn) SetAuthorizationType(v string) {
	o.AuthorizationType = &v
}

// GetBridgingTlsServerCertEnforceTrustedCommonNameEnabled returns the BridgingTlsServerCertEnforceTrustedCommonNameEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetBridgingTlsServerCertEnforceTrustedCommonNameEnabled() bool {
	if o == nil || o.BridgingTlsServerCertEnforceTrustedCommonNameEnabled == nil {
		var ret bool
		return ret
	}
	return *o.BridgingTlsServerCertEnforceTrustedCommonNameEnabled
}

// GetBridgingTlsServerCertEnforceTrustedCommonNameEnabledOk returns a tuple with the BridgingTlsServerCertEnforceTrustedCommonNameEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetBridgingTlsServerCertEnforceTrustedCommonNameEnabledOk() (*bool, bool) {
	if o == nil || o.BridgingTlsServerCertEnforceTrustedCommonNameEnabled == nil {
		return nil, false
	}
	return o.BridgingTlsServerCertEnforceTrustedCommonNameEnabled, true
}

// HasBridgingTlsServerCertEnforceTrustedCommonNameEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasBridgingTlsServerCertEnforceTrustedCommonNameEnabled() bool {
	if o != nil && o.BridgingTlsServerCertEnforceTrustedCommonNameEnabled != nil {
		return true
	}

	return false
}

// SetBridgingTlsServerCertEnforceTrustedCommonNameEnabled gets a reference to the given bool and assigns it to the BridgingTlsServerCertEnforceTrustedCommonNameEnabled field.
func (o *MsgVpn) SetBridgingTlsServerCertEnforceTrustedCommonNameEnabled(v bool) {
	o.BridgingTlsServerCertEnforceTrustedCommonNameEnabled = &v
}

// GetBridgingTlsServerCertMaxChainDepth returns the BridgingTlsServerCertMaxChainDepth field value if set, zero value otherwise.
func (o *MsgVpn) GetBridgingTlsServerCertMaxChainDepth() int64 {
	if o == nil || o.BridgingTlsServerCertMaxChainDepth == nil {
		var ret int64
		return ret
	}
	return *o.BridgingTlsServerCertMaxChainDepth
}

// GetBridgingTlsServerCertMaxChainDepthOk returns a tuple with the BridgingTlsServerCertMaxChainDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetBridgingTlsServerCertMaxChainDepthOk() (*int64, bool) {
	if o == nil || o.BridgingTlsServerCertMaxChainDepth == nil {
		return nil, false
	}
	return o.BridgingTlsServerCertMaxChainDepth, true
}

// HasBridgingTlsServerCertMaxChainDepth returns a boolean if a field has been set.
func (o *MsgVpn) HasBridgingTlsServerCertMaxChainDepth() bool {
	if o != nil && o.BridgingTlsServerCertMaxChainDepth != nil {
		return true
	}

	return false
}

// SetBridgingTlsServerCertMaxChainDepth gets a reference to the given int64 and assigns it to the BridgingTlsServerCertMaxChainDepth field.
func (o *MsgVpn) SetBridgingTlsServerCertMaxChainDepth(v int64) {
	o.BridgingTlsServerCertMaxChainDepth = &v
}

// GetBridgingTlsServerCertValidateDateEnabled returns the BridgingTlsServerCertValidateDateEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetBridgingTlsServerCertValidateDateEnabled() bool {
	if o == nil || o.BridgingTlsServerCertValidateDateEnabled == nil {
		var ret bool
		return ret
	}
	return *o.BridgingTlsServerCertValidateDateEnabled
}

// GetBridgingTlsServerCertValidateDateEnabledOk returns a tuple with the BridgingTlsServerCertValidateDateEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetBridgingTlsServerCertValidateDateEnabledOk() (*bool, bool) {
	if o == nil || o.BridgingTlsServerCertValidateDateEnabled == nil {
		return nil, false
	}
	return o.BridgingTlsServerCertValidateDateEnabled, true
}

// HasBridgingTlsServerCertValidateDateEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasBridgingTlsServerCertValidateDateEnabled() bool {
	if o != nil && o.BridgingTlsServerCertValidateDateEnabled != nil {
		return true
	}

	return false
}

// SetBridgingTlsServerCertValidateDateEnabled gets a reference to the given bool and assigns it to the BridgingTlsServerCertValidateDateEnabled field.
func (o *MsgVpn) SetBridgingTlsServerCertValidateDateEnabled(v bool) {
	o.BridgingTlsServerCertValidateDateEnabled = &v
}

// GetBridgingTlsServerCertValidateNameEnabled returns the BridgingTlsServerCertValidateNameEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetBridgingTlsServerCertValidateNameEnabled() bool {
	if o == nil || o.BridgingTlsServerCertValidateNameEnabled == nil {
		var ret bool
		return ret
	}
	return *o.BridgingTlsServerCertValidateNameEnabled
}

// GetBridgingTlsServerCertValidateNameEnabledOk returns a tuple with the BridgingTlsServerCertValidateNameEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetBridgingTlsServerCertValidateNameEnabledOk() (*bool, bool) {
	if o == nil || o.BridgingTlsServerCertValidateNameEnabled == nil {
		return nil, false
	}
	return o.BridgingTlsServerCertValidateNameEnabled, true
}

// HasBridgingTlsServerCertValidateNameEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasBridgingTlsServerCertValidateNameEnabled() bool {
	if o != nil && o.BridgingTlsServerCertValidateNameEnabled != nil {
		return true
	}

	return false
}

// SetBridgingTlsServerCertValidateNameEnabled gets a reference to the given bool and assigns it to the BridgingTlsServerCertValidateNameEnabled field.
func (o *MsgVpn) SetBridgingTlsServerCertValidateNameEnabled(v bool) {
	o.BridgingTlsServerCertValidateNameEnabled = &v
}

// GetDistributedCacheManagementEnabled returns the DistributedCacheManagementEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetDistributedCacheManagementEnabled() bool {
	if o == nil || o.DistributedCacheManagementEnabled == nil {
		var ret bool
		return ret
	}
	return *o.DistributedCacheManagementEnabled
}

// GetDistributedCacheManagementEnabledOk returns a tuple with the DistributedCacheManagementEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetDistributedCacheManagementEnabledOk() (*bool, bool) {
	if o == nil || o.DistributedCacheManagementEnabled == nil {
		return nil, false
	}
	return o.DistributedCacheManagementEnabled, true
}

// HasDistributedCacheManagementEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasDistributedCacheManagementEnabled() bool {
	if o != nil && o.DistributedCacheManagementEnabled != nil {
		return true
	}

	return false
}

// SetDistributedCacheManagementEnabled gets a reference to the given bool and assigns it to the DistributedCacheManagementEnabled field.
func (o *MsgVpn) SetDistributedCacheManagementEnabled(v bool) {
	o.DistributedCacheManagementEnabled = &v
}

// GetDmrEnabled returns the DmrEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetDmrEnabled() bool {
	if o == nil || o.DmrEnabled == nil {
		var ret bool
		return ret
	}
	return *o.DmrEnabled
}

// GetDmrEnabledOk returns a tuple with the DmrEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetDmrEnabledOk() (*bool, bool) {
	if o == nil || o.DmrEnabled == nil {
		return nil, false
	}
	return o.DmrEnabled, true
}

// HasDmrEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasDmrEnabled() bool {
	if o != nil && o.DmrEnabled != nil {
		return true
	}

	return false
}

// SetDmrEnabled gets a reference to the given bool and assigns it to the DmrEnabled field.
func (o *MsgVpn) SetDmrEnabled(v bool) {
	o.DmrEnabled = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *MsgVpn) GetEnabled() bool {
	if o == nil || o.Enabled == nil {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEnabledOk() (*bool, bool) {
	if o == nil || o.Enabled == nil {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasEnabled() bool {
	if o != nil && o.Enabled != nil {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *MsgVpn) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetEventConnectionCountThreshold returns the EventConnectionCountThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventConnectionCountThreshold() EventThreshold {
	if o == nil || o.EventConnectionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventConnectionCountThreshold
}

// GetEventConnectionCountThresholdOk returns a tuple with the EventConnectionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventConnectionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventConnectionCountThreshold == nil {
		return nil, false
	}
	return o.EventConnectionCountThreshold, true
}

// HasEventConnectionCountThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventConnectionCountThreshold() bool {
	if o != nil && o.EventConnectionCountThreshold != nil {
		return true
	}

	return false
}

// SetEventConnectionCountThreshold gets a reference to the given EventThreshold and assigns it to the EventConnectionCountThreshold field.
func (o *MsgVpn) SetEventConnectionCountThreshold(v EventThreshold) {
	o.EventConnectionCountThreshold = &v
}

// GetEventEgressFlowCountThreshold returns the EventEgressFlowCountThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventEgressFlowCountThreshold() EventThreshold {
	if o == nil || o.EventEgressFlowCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventEgressFlowCountThreshold
}

// GetEventEgressFlowCountThresholdOk returns a tuple with the EventEgressFlowCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventEgressFlowCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventEgressFlowCountThreshold == nil {
		return nil, false
	}
	return o.EventEgressFlowCountThreshold, true
}

// HasEventEgressFlowCountThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventEgressFlowCountThreshold() bool {
	if o != nil && o.EventEgressFlowCountThreshold != nil {
		return true
	}

	return false
}

// SetEventEgressFlowCountThreshold gets a reference to the given EventThreshold and assigns it to the EventEgressFlowCountThreshold field.
func (o *MsgVpn) SetEventEgressFlowCountThreshold(v EventThreshold) {
	o.EventEgressFlowCountThreshold = &v
}

// GetEventEgressMsgRateThreshold returns the EventEgressMsgRateThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventEgressMsgRateThreshold() EventThresholdByValue {
	if o == nil || o.EventEgressMsgRateThreshold == nil {
		var ret EventThresholdByValue
		return ret
	}
	return *o.EventEgressMsgRateThreshold
}

// GetEventEgressMsgRateThresholdOk returns a tuple with the EventEgressMsgRateThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventEgressMsgRateThresholdOk() (*EventThresholdByValue, bool) {
	if o == nil || o.EventEgressMsgRateThreshold == nil {
		return nil, false
	}
	return o.EventEgressMsgRateThreshold, true
}

// HasEventEgressMsgRateThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventEgressMsgRateThreshold() bool {
	if o != nil && o.EventEgressMsgRateThreshold != nil {
		return true
	}

	return false
}

// SetEventEgressMsgRateThreshold gets a reference to the given EventThresholdByValue and assigns it to the EventEgressMsgRateThreshold field.
func (o *MsgVpn) SetEventEgressMsgRateThreshold(v EventThresholdByValue) {
	o.EventEgressMsgRateThreshold = &v
}

// GetEventEndpointCountThreshold returns the EventEndpointCountThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventEndpointCountThreshold() EventThreshold {
	if o == nil || o.EventEndpointCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventEndpointCountThreshold
}

// GetEventEndpointCountThresholdOk returns a tuple with the EventEndpointCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventEndpointCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventEndpointCountThreshold == nil {
		return nil, false
	}
	return o.EventEndpointCountThreshold, true
}

// HasEventEndpointCountThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventEndpointCountThreshold() bool {
	if o != nil && o.EventEndpointCountThreshold != nil {
		return true
	}

	return false
}

// SetEventEndpointCountThreshold gets a reference to the given EventThreshold and assigns it to the EventEndpointCountThreshold field.
func (o *MsgVpn) SetEventEndpointCountThreshold(v EventThreshold) {
	o.EventEndpointCountThreshold = &v
}

// GetEventIngressFlowCountThreshold returns the EventIngressFlowCountThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventIngressFlowCountThreshold() EventThreshold {
	if o == nil || o.EventIngressFlowCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventIngressFlowCountThreshold
}

// GetEventIngressFlowCountThresholdOk returns a tuple with the EventIngressFlowCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventIngressFlowCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventIngressFlowCountThreshold == nil {
		return nil, false
	}
	return o.EventIngressFlowCountThreshold, true
}

// HasEventIngressFlowCountThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventIngressFlowCountThreshold() bool {
	if o != nil && o.EventIngressFlowCountThreshold != nil {
		return true
	}

	return false
}

// SetEventIngressFlowCountThreshold gets a reference to the given EventThreshold and assigns it to the EventIngressFlowCountThreshold field.
func (o *MsgVpn) SetEventIngressFlowCountThreshold(v EventThreshold) {
	o.EventIngressFlowCountThreshold = &v
}

// GetEventIngressMsgRateThreshold returns the EventIngressMsgRateThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventIngressMsgRateThreshold() EventThresholdByValue {
	if o == nil || o.EventIngressMsgRateThreshold == nil {
		var ret EventThresholdByValue
		return ret
	}
	return *o.EventIngressMsgRateThreshold
}

// GetEventIngressMsgRateThresholdOk returns a tuple with the EventIngressMsgRateThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventIngressMsgRateThresholdOk() (*EventThresholdByValue, bool) {
	if o == nil || o.EventIngressMsgRateThreshold == nil {
		return nil, false
	}
	return o.EventIngressMsgRateThreshold, true
}

// HasEventIngressMsgRateThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventIngressMsgRateThreshold() bool {
	if o != nil && o.EventIngressMsgRateThreshold != nil {
		return true
	}

	return false
}

// SetEventIngressMsgRateThreshold gets a reference to the given EventThresholdByValue and assigns it to the EventIngressMsgRateThreshold field.
func (o *MsgVpn) SetEventIngressMsgRateThreshold(v EventThresholdByValue) {
	o.EventIngressMsgRateThreshold = &v
}

// GetEventLargeMsgThreshold returns the EventLargeMsgThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventLargeMsgThreshold() int64 {
	if o == nil || o.EventLargeMsgThreshold == nil {
		var ret int64
		return ret
	}
	return *o.EventLargeMsgThreshold
}

// GetEventLargeMsgThresholdOk returns a tuple with the EventLargeMsgThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventLargeMsgThresholdOk() (*int64, bool) {
	if o == nil || o.EventLargeMsgThreshold == nil {
		return nil, false
	}
	return o.EventLargeMsgThreshold, true
}

// HasEventLargeMsgThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventLargeMsgThreshold() bool {
	if o != nil && o.EventLargeMsgThreshold != nil {
		return true
	}

	return false
}

// SetEventLargeMsgThreshold gets a reference to the given int64 and assigns it to the EventLargeMsgThreshold field.
func (o *MsgVpn) SetEventLargeMsgThreshold(v int64) {
	o.EventLargeMsgThreshold = &v
}

// GetEventLogTag returns the EventLogTag field value if set, zero value otherwise.
func (o *MsgVpn) GetEventLogTag() string {
	if o == nil || o.EventLogTag == nil {
		var ret string
		return ret
	}
	return *o.EventLogTag
}

// GetEventLogTagOk returns a tuple with the EventLogTag field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventLogTagOk() (*string, bool) {
	if o == nil || o.EventLogTag == nil {
		return nil, false
	}
	return o.EventLogTag, true
}

// HasEventLogTag returns a boolean if a field has been set.
func (o *MsgVpn) HasEventLogTag() bool {
	if o != nil && o.EventLogTag != nil {
		return true
	}

	return false
}

// SetEventLogTag gets a reference to the given string and assigns it to the EventLogTag field.
func (o *MsgVpn) SetEventLogTag(v string) {
	o.EventLogTag = &v
}

// GetEventMsgSpoolUsageThreshold returns the EventMsgSpoolUsageThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventMsgSpoolUsageThreshold() EventThreshold {
	if o == nil || o.EventMsgSpoolUsageThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventMsgSpoolUsageThreshold
}

// GetEventMsgSpoolUsageThresholdOk returns a tuple with the EventMsgSpoolUsageThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventMsgSpoolUsageThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventMsgSpoolUsageThreshold == nil {
		return nil, false
	}
	return o.EventMsgSpoolUsageThreshold, true
}

// HasEventMsgSpoolUsageThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventMsgSpoolUsageThreshold() bool {
	if o != nil && o.EventMsgSpoolUsageThreshold != nil {
		return true
	}

	return false
}

// SetEventMsgSpoolUsageThreshold gets a reference to the given EventThreshold and assigns it to the EventMsgSpoolUsageThreshold field.
func (o *MsgVpn) SetEventMsgSpoolUsageThreshold(v EventThreshold) {
	o.EventMsgSpoolUsageThreshold = &v
}

// GetEventPublishClientEnabled returns the EventPublishClientEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetEventPublishClientEnabled() bool {
	if o == nil || o.EventPublishClientEnabled == nil {
		var ret bool
		return ret
	}
	return *o.EventPublishClientEnabled
}

// GetEventPublishClientEnabledOk returns a tuple with the EventPublishClientEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventPublishClientEnabledOk() (*bool, bool) {
	if o == nil || o.EventPublishClientEnabled == nil {
		return nil, false
	}
	return o.EventPublishClientEnabled, true
}

// HasEventPublishClientEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasEventPublishClientEnabled() bool {
	if o != nil && o.EventPublishClientEnabled != nil {
		return true
	}

	return false
}

// SetEventPublishClientEnabled gets a reference to the given bool and assigns it to the EventPublishClientEnabled field.
func (o *MsgVpn) SetEventPublishClientEnabled(v bool) {
	o.EventPublishClientEnabled = &v
}

// GetEventPublishMsgVpnEnabled returns the EventPublishMsgVpnEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetEventPublishMsgVpnEnabled() bool {
	if o == nil || o.EventPublishMsgVpnEnabled == nil {
		var ret bool
		return ret
	}
	return *o.EventPublishMsgVpnEnabled
}

// GetEventPublishMsgVpnEnabledOk returns a tuple with the EventPublishMsgVpnEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventPublishMsgVpnEnabledOk() (*bool, bool) {
	if o == nil || o.EventPublishMsgVpnEnabled == nil {
		return nil, false
	}
	return o.EventPublishMsgVpnEnabled, true
}

// HasEventPublishMsgVpnEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasEventPublishMsgVpnEnabled() bool {
	if o != nil && o.EventPublishMsgVpnEnabled != nil {
		return true
	}

	return false
}

// SetEventPublishMsgVpnEnabled gets a reference to the given bool and assigns it to the EventPublishMsgVpnEnabled field.
func (o *MsgVpn) SetEventPublishMsgVpnEnabled(v bool) {
	o.EventPublishMsgVpnEnabled = &v
}

// GetEventPublishSubscriptionMode returns the EventPublishSubscriptionMode field value if set, zero value otherwise.
func (o *MsgVpn) GetEventPublishSubscriptionMode() string {
	if o == nil || o.EventPublishSubscriptionMode == nil {
		var ret string
		return ret
	}
	return *o.EventPublishSubscriptionMode
}

// GetEventPublishSubscriptionModeOk returns a tuple with the EventPublishSubscriptionMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventPublishSubscriptionModeOk() (*string, bool) {
	if o == nil || o.EventPublishSubscriptionMode == nil {
		return nil, false
	}
	return o.EventPublishSubscriptionMode, true
}

// HasEventPublishSubscriptionMode returns a boolean if a field has been set.
func (o *MsgVpn) HasEventPublishSubscriptionMode() bool {
	if o != nil && o.EventPublishSubscriptionMode != nil {
		return true
	}

	return false
}

// SetEventPublishSubscriptionMode gets a reference to the given string and assigns it to the EventPublishSubscriptionMode field.
func (o *MsgVpn) SetEventPublishSubscriptionMode(v string) {
	o.EventPublishSubscriptionMode = &v
}

// GetEventPublishTopicFormatMqttEnabled returns the EventPublishTopicFormatMqttEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetEventPublishTopicFormatMqttEnabled() bool {
	if o == nil || o.EventPublishTopicFormatMqttEnabled == nil {
		var ret bool
		return ret
	}
	return *o.EventPublishTopicFormatMqttEnabled
}

// GetEventPublishTopicFormatMqttEnabledOk returns a tuple with the EventPublishTopicFormatMqttEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventPublishTopicFormatMqttEnabledOk() (*bool, bool) {
	if o == nil || o.EventPublishTopicFormatMqttEnabled == nil {
		return nil, false
	}
	return o.EventPublishTopicFormatMqttEnabled, true
}

// HasEventPublishTopicFormatMqttEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasEventPublishTopicFormatMqttEnabled() bool {
	if o != nil && o.EventPublishTopicFormatMqttEnabled != nil {
		return true
	}

	return false
}

// SetEventPublishTopicFormatMqttEnabled gets a reference to the given bool and assigns it to the EventPublishTopicFormatMqttEnabled field.
func (o *MsgVpn) SetEventPublishTopicFormatMqttEnabled(v bool) {
	o.EventPublishTopicFormatMqttEnabled = &v
}

// GetEventPublishTopicFormatSmfEnabled returns the EventPublishTopicFormatSmfEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetEventPublishTopicFormatSmfEnabled() bool {
	if o == nil || o.EventPublishTopicFormatSmfEnabled == nil {
		var ret bool
		return ret
	}
	return *o.EventPublishTopicFormatSmfEnabled
}

// GetEventPublishTopicFormatSmfEnabledOk returns a tuple with the EventPublishTopicFormatSmfEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventPublishTopicFormatSmfEnabledOk() (*bool, bool) {
	if o == nil || o.EventPublishTopicFormatSmfEnabled == nil {
		return nil, false
	}
	return o.EventPublishTopicFormatSmfEnabled, true
}

// HasEventPublishTopicFormatSmfEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasEventPublishTopicFormatSmfEnabled() bool {
	if o != nil && o.EventPublishTopicFormatSmfEnabled != nil {
		return true
	}

	return false
}

// SetEventPublishTopicFormatSmfEnabled gets a reference to the given bool and assigns it to the EventPublishTopicFormatSmfEnabled field.
func (o *MsgVpn) SetEventPublishTopicFormatSmfEnabled(v bool) {
	o.EventPublishTopicFormatSmfEnabled = &v
}

// GetEventServiceAmqpConnectionCountThreshold returns the EventServiceAmqpConnectionCountThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventServiceAmqpConnectionCountThreshold() EventThreshold {
	if o == nil || o.EventServiceAmqpConnectionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventServiceAmqpConnectionCountThreshold
}

// GetEventServiceAmqpConnectionCountThresholdOk returns a tuple with the EventServiceAmqpConnectionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventServiceAmqpConnectionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventServiceAmqpConnectionCountThreshold == nil {
		return nil, false
	}
	return o.EventServiceAmqpConnectionCountThreshold, true
}

// HasEventServiceAmqpConnectionCountThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventServiceAmqpConnectionCountThreshold() bool {
	if o != nil && o.EventServiceAmqpConnectionCountThreshold != nil {
		return true
	}

	return false
}

// SetEventServiceAmqpConnectionCountThreshold gets a reference to the given EventThreshold and assigns it to the EventServiceAmqpConnectionCountThreshold field.
func (o *MsgVpn) SetEventServiceAmqpConnectionCountThreshold(v EventThreshold) {
	o.EventServiceAmqpConnectionCountThreshold = &v
}

// GetEventServiceMqttConnectionCountThreshold returns the EventServiceMqttConnectionCountThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventServiceMqttConnectionCountThreshold() EventThreshold {
	if o == nil || o.EventServiceMqttConnectionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventServiceMqttConnectionCountThreshold
}

// GetEventServiceMqttConnectionCountThresholdOk returns a tuple with the EventServiceMqttConnectionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventServiceMqttConnectionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventServiceMqttConnectionCountThreshold == nil {
		return nil, false
	}
	return o.EventServiceMqttConnectionCountThreshold, true
}

// HasEventServiceMqttConnectionCountThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventServiceMqttConnectionCountThreshold() bool {
	if o != nil && o.EventServiceMqttConnectionCountThreshold != nil {
		return true
	}

	return false
}

// SetEventServiceMqttConnectionCountThreshold gets a reference to the given EventThreshold and assigns it to the EventServiceMqttConnectionCountThreshold field.
func (o *MsgVpn) SetEventServiceMqttConnectionCountThreshold(v EventThreshold) {
	o.EventServiceMqttConnectionCountThreshold = &v
}

// GetEventServiceRestIncomingConnectionCountThreshold returns the EventServiceRestIncomingConnectionCountThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventServiceRestIncomingConnectionCountThreshold() EventThreshold {
	if o == nil || o.EventServiceRestIncomingConnectionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventServiceRestIncomingConnectionCountThreshold
}

// GetEventServiceRestIncomingConnectionCountThresholdOk returns a tuple with the EventServiceRestIncomingConnectionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventServiceRestIncomingConnectionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventServiceRestIncomingConnectionCountThreshold == nil {
		return nil, false
	}
	return o.EventServiceRestIncomingConnectionCountThreshold, true
}

// HasEventServiceRestIncomingConnectionCountThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventServiceRestIncomingConnectionCountThreshold() bool {
	if o != nil && o.EventServiceRestIncomingConnectionCountThreshold != nil {
		return true
	}

	return false
}

// SetEventServiceRestIncomingConnectionCountThreshold gets a reference to the given EventThreshold and assigns it to the EventServiceRestIncomingConnectionCountThreshold field.
func (o *MsgVpn) SetEventServiceRestIncomingConnectionCountThreshold(v EventThreshold) {
	o.EventServiceRestIncomingConnectionCountThreshold = &v
}

// GetEventServiceSmfConnectionCountThreshold returns the EventServiceSmfConnectionCountThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventServiceSmfConnectionCountThreshold() EventThreshold {
	if o == nil || o.EventServiceSmfConnectionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventServiceSmfConnectionCountThreshold
}

// GetEventServiceSmfConnectionCountThresholdOk returns a tuple with the EventServiceSmfConnectionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventServiceSmfConnectionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventServiceSmfConnectionCountThreshold == nil {
		return nil, false
	}
	return o.EventServiceSmfConnectionCountThreshold, true
}

// HasEventServiceSmfConnectionCountThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventServiceSmfConnectionCountThreshold() bool {
	if o != nil && o.EventServiceSmfConnectionCountThreshold != nil {
		return true
	}

	return false
}

// SetEventServiceSmfConnectionCountThreshold gets a reference to the given EventThreshold and assigns it to the EventServiceSmfConnectionCountThreshold field.
func (o *MsgVpn) SetEventServiceSmfConnectionCountThreshold(v EventThreshold) {
	o.EventServiceSmfConnectionCountThreshold = &v
}

// GetEventServiceWebConnectionCountThreshold returns the EventServiceWebConnectionCountThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventServiceWebConnectionCountThreshold() EventThreshold {
	if o == nil || o.EventServiceWebConnectionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventServiceWebConnectionCountThreshold
}

// GetEventServiceWebConnectionCountThresholdOk returns a tuple with the EventServiceWebConnectionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventServiceWebConnectionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventServiceWebConnectionCountThreshold == nil {
		return nil, false
	}
	return o.EventServiceWebConnectionCountThreshold, true
}

// HasEventServiceWebConnectionCountThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventServiceWebConnectionCountThreshold() bool {
	if o != nil && o.EventServiceWebConnectionCountThreshold != nil {
		return true
	}

	return false
}

// SetEventServiceWebConnectionCountThreshold gets a reference to the given EventThreshold and assigns it to the EventServiceWebConnectionCountThreshold field.
func (o *MsgVpn) SetEventServiceWebConnectionCountThreshold(v EventThreshold) {
	o.EventServiceWebConnectionCountThreshold = &v
}

// GetEventSubscriptionCountThreshold returns the EventSubscriptionCountThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventSubscriptionCountThreshold() EventThreshold {
	if o == nil || o.EventSubscriptionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventSubscriptionCountThreshold
}

// GetEventSubscriptionCountThresholdOk returns a tuple with the EventSubscriptionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventSubscriptionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventSubscriptionCountThreshold == nil {
		return nil, false
	}
	return o.EventSubscriptionCountThreshold, true
}

// HasEventSubscriptionCountThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventSubscriptionCountThreshold() bool {
	if o != nil && o.EventSubscriptionCountThreshold != nil {
		return true
	}

	return false
}

// SetEventSubscriptionCountThreshold gets a reference to the given EventThreshold and assigns it to the EventSubscriptionCountThreshold field.
func (o *MsgVpn) SetEventSubscriptionCountThreshold(v EventThreshold) {
	o.EventSubscriptionCountThreshold = &v
}

// GetEventTransactedSessionCountThreshold returns the EventTransactedSessionCountThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventTransactedSessionCountThreshold() EventThreshold {
	if o == nil || o.EventTransactedSessionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventTransactedSessionCountThreshold
}

// GetEventTransactedSessionCountThresholdOk returns a tuple with the EventTransactedSessionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventTransactedSessionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventTransactedSessionCountThreshold == nil {
		return nil, false
	}
	return o.EventTransactedSessionCountThreshold, true
}

// HasEventTransactedSessionCountThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventTransactedSessionCountThreshold() bool {
	if o != nil && o.EventTransactedSessionCountThreshold != nil {
		return true
	}

	return false
}

// SetEventTransactedSessionCountThreshold gets a reference to the given EventThreshold and assigns it to the EventTransactedSessionCountThreshold field.
func (o *MsgVpn) SetEventTransactedSessionCountThreshold(v EventThreshold) {
	o.EventTransactedSessionCountThreshold = &v
}

// GetEventTransactionCountThreshold returns the EventTransactionCountThreshold field value if set, zero value otherwise.
func (o *MsgVpn) GetEventTransactionCountThreshold() EventThreshold {
	if o == nil || o.EventTransactionCountThreshold == nil {
		var ret EventThreshold
		return ret
	}
	return *o.EventTransactionCountThreshold
}

// GetEventTransactionCountThresholdOk returns a tuple with the EventTransactionCountThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetEventTransactionCountThresholdOk() (*EventThreshold, bool) {
	if o == nil || o.EventTransactionCountThreshold == nil {
		return nil, false
	}
	return o.EventTransactionCountThreshold, true
}

// HasEventTransactionCountThreshold returns a boolean if a field has been set.
func (o *MsgVpn) HasEventTransactionCountThreshold() bool {
	if o != nil && o.EventTransactionCountThreshold != nil {
		return true
	}

	return false
}

// SetEventTransactionCountThreshold gets a reference to the given EventThreshold and assigns it to the EventTransactionCountThreshold field.
func (o *MsgVpn) SetEventTransactionCountThreshold(v EventThreshold) {
	o.EventTransactionCountThreshold = &v
}

// GetExportSubscriptionsEnabled returns the ExportSubscriptionsEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetExportSubscriptionsEnabled() bool {
	if o == nil || o.ExportSubscriptionsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ExportSubscriptionsEnabled
}

// GetExportSubscriptionsEnabledOk returns a tuple with the ExportSubscriptionsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetExportSubscriptionsEnabledOk() (*bool, bool) {
	if o == nil || o.ExportSubscriptionsEnabled == nil {
		return nil, false
	}
	return o.ExportSubscriptionsEnabled, true
}

// HasExportSubscriptionsEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasExportSubscriptionsEnabled() bool {
	if o != nil && o.ExportSubscriptionsEnabled != nil {
		return true
	}

	return false
}

// SetExportSubscriptionsEnabled gets a reference to the given bool and assigns it to the ExportSubscriptionsEnabled field.
func (o *MsgVpn) SetExportSubscriptionsEnabled(v bool) {
	o.ExportSubscriptionsEnabled = &v
}

// GetJndiEnabled returns the JndiEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetJndiEnabled() bool {
	if o == nil || o.JndiEnabled == nil {
		var ret bool
		return ret
	}
	return *o.JndiEnabled
}

// GetJndiEnabledOk returns a tuple with the JndiEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetJndiEnabledOk() (*bool, bool) {
	if o == nil || o.JndiEnabled == nil {
		return nil, false
	}
	return o.JndiEnabled, true
}

// HasJndiEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasJndiEnabled() bool {
	if o != nil && o.JndiEnabled != nil {
		return true
	}

	return false
}

// SetJndiEnabled gets a reference to the given bool and assigns it to the JndiEnabled field.
func (o *MsgVpn) SetJndiEnabled(v bool) {
	o.JndiEnabled = &v
}

// GetMaxConnectionCount returns the MaxConnectionCount field value if set, zero value otherwise.
func (o *MsgVpn) GetMaxConnectionCount() int64 {
	if o == nil || o.MaxConnectionCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxConnectionCount
}

// GetMaxConnectionCountOk returns a tuple with the MaxConnectionCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetMaxConnectionCountOk() (*int64, bool) {
	if o == nil || o.MaxConnectionCount == nil {
		return nil, false
	}
	return o.MaxConnectionCount, true
}

// HasMaxConnectionCount returns a boolean if a field has been set.
func (o *MsgVpn) HasMaxConnectionCount() bool {
	if o != nil && o.MaxConnectionCount != nil {
		return true
	}

	return false
}

// SetMaxConnectionCount gets a reference to the given int64 and assigns it to the MaxConnectionCount field.
func (o *MsgVpn) SetMaxConnectionCount(v int64) {
	o.MaxConnectionCount = &v
}

// GetMaxEgressFlowCount returns the MaxEgressFlowCount field value if set, zero value otherwise.
func (o *MsgVpn) GetMaxEgressFlowCount() int64 {
	if o == nil || o.MaxEgressFlowCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxEgressFlowCount
}

// GetMaxEgressFlowCountOk returns a tuple with the MaxEgressFlowCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetMaxEgressFlowCountOk() (*int64, bool) {
	if o == nil || o.MaxEgressFlowCount == nil {
		return nil, false
	}
	return o.MaxEgressFlowCount, true
}

// HasMaxEgressFlowCount returns a boolean if a field has been set.
func (o *MsgVpn) HasMaxEgressFlowCount() bool {
	if o != nil && o.MaxEgressFlowCount != nil {
		return true
	}

	return false
}

// SetMaxEgressFlowCount gets a reference to the given int64 and assigns it to the MaxEgressFlowCount field.
func (o *MsgVpn) SetMaxEgressFlowCount(v int64) {
	o.MaxEgressFlowCount = &v
}

// GetMaxEndpointCount returns the MaxEndpointCount field value if set, zero value otherwise.
func (o *MsgVpn) GetMaxEndpointCount() int64 {
	if o == nil || o.MaxEndpointCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxEndpointCount
}

// GetMaxEndpointCountOk returns a tuple with the MaxEndpointCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetMaxEndpointCountOk() (*int64, bool) {
	if o == nil || o.MaxEndpointCount == nil {
		return nil, false
	}
	return o.MaxEndpointCount, true
}

// HasMaxEndpointCount returns a boolean if a field has been set.
func (o *MsgVpn) HasMaxEndpointCount() bool {
	if o != nil && o.MaxEndpointCount != nil {
		return true
	}

	return false
}

// SetMaxEndpointCount gets a reference to the given int64 and assigns it to the MaxEndpointCount field.
func (o *MsgVpn) SetMaxEndpointCount(v int64) {
	o.MaxEndpointCount = &v
}

// GetMaxIngressFlowCount returns the MaxIngressFlowCount field value if set, zero value otherwise.
func (o *MsgVpn) GetMaxIngressFlowCount() int64 {
	if o == nil || o.MaxIngressFlowCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxIngressFlowCount
}

// GetMaxIngressFlowCountOk returns a tuple with the MaxIngressFlowCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetMaxIngressFlowCountOk() (*int64, bool) {
	if o == nil || o.MaxIngressFlowCount == nil {
		return nil, false
	}
	return o.MaxIngressFlowCount, true
}

// HasMaxIngressFlowCount returns a boolean if a field has been set.
func (o *MsgVpn) HasMaxIngressFlowCount() bool {
	if o != nil && o.MaxIngressFlowCount != nil {
		return true
	}

	return false
}

// SetMaxIngressFlowCount gets a reference to the given int64 and assigns it to the MaxIngressFlowCount field.
func (o *MsgVpn) SetMaxIngressFlowCount(v int64) {
	o.MaxIngressFlowCount = &v
}

// GetMaxMsgSpoolUsage returns the MaxMsgSpoolUsage field value if set, zero value otherwise.
func (o *MsgVpn) GetMaxMsgSpoolUsage() int64 {
	if o == nil || o.MaxMsgSpoolUsage == nil {
		var ret int64
		return ret
	}
	return *o.MaxMsgSpoolUsage
}

// GetMaxMsgSpoolUsageOk returns a tuple with the MaxMsgSpoolUsage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetMaxMsgSpoolUsageOk() (*int64, bool) {
	if o == nil || o.MaxMsgSpoolUsage == nil {
		return nil, false
	}
	return o.MaxMsgSpoolUsage, true
}

// HasMaxMsgSpoolUsage returns a boolean if a field has been set.
func (o *MsgVpn) HasMaxMsgSpoolUsage() bool {
	if o != nil && o.MaxMsgSpoolUsage != nil {
		return true
	}

	return false
}

// SetMaxMsgSpoolUsage gets a reference to the given int64 and assigns it to the MaxMsgSpoolUsage field.
func (o *MsgVpn) SetMaxMsgSpoolUsage(v int64) {
	o.MaxMsgSpoolUsage = &v
}

// GetMaxSubscriptionCount returns the MaxSubscriptionCount field value if set, zero value otherwise.
func (o *MsgVpn) GetMaxSubscriptionCount() int64 {
	if o == nil || o.MaxSubscriptionCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxSubscriptionCount
}

// GetMaxSubscriptionCountOk returns a tuple with the MaxSubscriptionCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetMaxSubscriptionCountOk() (*int64, bool) {
	if o == nil || o.MaxSubscriptionCount == nil {
		return nil, false
	}
	return o.MaxSubscriptionCount, true
}

// HasMaxSubscriptionCount returns a boolean if a field has been set.
func (o *MsgVpn) HasMaxSubscriptionCount() bool {
	if o != nil && o.MaxSubscriptionCount != nil {
		return true
	}

	return false
}

// SetMaxSubscriptionCount gets a reference to the given int64 and assigns it to the MaxSubscriptionCount field.
func (o *MsgVpn) SetMaxSubscriptionCount(v int64) {
	o.MaxSubscriptionCount = &v
}

// GetMaxTransactedSessionCount returns the MaxTransactedSessionCount field value if set, zero value otherwise.
func (o *MsgVpn) GetMaxTransactedSessionCount() int64 {
	if o == nil || o.MaxTransactedSessionCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxTransactedSessionCount
}

// GetMaxTransactedSessionCountOk returns a tuple with the MaxTransactedSessionCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetMaxTransactedSessionCountOk() (*int64, bool) {
	if o == nil || o.MaxTransactedSessionCount == nil {
		return nil, false
	}
	return o.MaxTransactedSessionCount, true
}

// HasMaxTransactedSessionCount returns a boolean if a field has been set.
func (o *MsgVpn) HasMaxTransactedSessionCount() bool {
	if o != nil && o.MaxTransactedSessionCount != nil {
		return true
	}

	return false
}

// SetMaxTransactedSessionCount gets a reference to the given int64 and assigns it to the MaxTransactedSessionCount field.
func (o *MsgVpn) SetMaxTransactedSessionCount(v int64) {
	o.MaxTransactedSessionCount = &v
}

// GetMaxTransactionCount returns the MaxTransactionCount field value if set, zero value otherwise.
func (o *MsgVpn) GetMaxTransactionCount() int64 {
	if o == nil || o.MaxTransactionCount == nil {
		var ret int64
		return ret
	}
	return *o.MaxTransactionCount
}

// GetMaxTransactionCountOk returns a tuple with the MaxTransactionCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetMaxTransactionCountOk() (*int64, bool) {
	if o == nil || o.MaxTransactionCount == nil {
		return nil, false
	}
	return o.MaxTransactionCount, true
}

// HasMaxTransactionCount returns a boolean if a field has been set.
func (o *MsgVpn) HasMaxTransactionCount() bool {
	if o != nil && o.MaxTransactionCount != nil {
		return true
	}

	return false
}

// SetMaxTransactionCount gets a reference to the given int64 and assigns it to the MaxTransactionCount field.
func (o *MsgVpn) SetMaxTransactionCount(v int64) {
	o.MaxTransactionCount = &v
}

// GetMqttRetainMaxMemory returns the MqttRetainMaxMemory field value if set, zero value otherwise.
func (o *MsgVpn) GetMqttRetainMaxMemory() int32 {
	if o == nil || o.MqttRetainMaxMemory == nil {
		var ret int32
		return ret
	}
	return *o.MqttRetainMaxMemory
}

// GetMqttRetainMaxMemoryOk returns a tuple with the MqttRetainMaxMemory field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetMqttRetainMaxMemoryOk() (*int32, bool) {
	if o == nil || o.MqttRetainMaxMemory == nil {
		return nil, false
	}
	return o.MqttRetainMaxMemory, true
}

// HasMqttRetainMaxMemory returns a boolean if a field has been set.
func (o *MsgVpn) HasMqttRetainMaxMemory() bool {
	if o != nil && o.MqttRetainMaxMemory != nil {
		return true
	}

	return false
}

// SetMqttRetainMaxMemory gets a reference to the given int32 and assigns it to the MqttRetainMaxMemory field.
func (o *MsgVpn) SetMqttRetainMaxMemory(v int32) {
	o.MqttRetainMaxMemory = &v
}

// GetMsgVpnName returns the MsgVpnName field value if set, zero value otherwise.
func (o *MsgVpn) GetMsgVpnName() string {
	if o == nil || o.MsgVpnName == nil {
		var ret string
		return ret
	}
	return *o.MsgVpnName
}

// GetMsgVpnNameOk returns a tuple with the MsgVpnName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetMsgVpnNameOk() (*string, bool) {
	if o == nil || o.MsgVpnName == nil {
		return nil, false
	}
	return o.MsgVpnName, true
}

// HasMsgVpnName returns a boolean if a field has been set.
func (o *MsgVpn) HasMsgVpnName() bool {
	if o != nil && o.MsgVpnName != nil {
		return true
	}

	return false
}

// SetMsgVpnName gets a reference to the given string and assigns it to the MsgVpnName field.
func (o *MsgVpn) SetMsgVpnName(v string) {
	o.MsgVpnName = &v
}

// GetReplicationAckPropagationIntervalMsgCount returns the ReplicationAckPropagationIntervalMsgCount field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationAckPropagationIntervalMsgCount() int64 {
	if o == nil || o.ReplicationAckPropagationIntervalMsgCount == nil {
		var ret int64
		return ret
	}
	return *o.ReplicationAckPropagationIntervalMsgCount
}

// GetReplicationAckPropagationIntervalMsgCountOk returns a tuple with the ReplicationAckPropagationIntervalMsgCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationAckPropagationIntervalMsgCountOk() (*int64, bool) {
	if o == nil || o.ReplicationAckPropagationIntervalMsgCount == nil {
		return nil, false
	}
	return o.ReplicationAckPropagationIntervalMsgCount, true
}

// HasReplicationAckPropagationIntervalMsgCount returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationAckPropagationIntervalMsgCount() bool {
	if o != nil && o.ReplicationAckPropagationIntervalMsgCount != nil {
		return true
	}

	return false
}

// SetReplicationAckPropagationIntervalMsgCount gets a reference to the given int64 and assigns it to the ReplicationAckPropagationIntervalMsgCount field.
func (o *MsgVpn) SetReplicationAckPropagationIntervalMsgCount(v int64) {
	o.ReplicationAckPropagationIntervalMsgCount = &v
}

// GetReplicationBridgeAuthenticationBasicClientUsername returns the ReplicationBridgeAuthenticationBasicClientUsername field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationBridgeAuthenticationBasicClientUsername() string {
	if o == nil || o.ReplicationBridgeAuthenticationBasicClientUsername == nil {
		var ret string
		return ret
	}
	return *o.ReplicationBridgeAuthenticationBasicClientUsername
}

// GetReplicationBridgeAuthenticationBasicClientUsernameOk returns a tuple with the ReplicationBridgeAuthenticationBasicClientUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationBridgeAuthenticationBasicClientUsernameOk() (*string, bool) {
	if o == nil || o.ReplicationBridgeAuthenticationBasicClientUsername == nil {
		return nil, false
	}
	return o.ReplicationBridgeAuthenticationBasicClientUsername, true
}

// HasReplicationBridgeAuthenticationBasicClientUsername returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationBridgeAuthenticationBasicClientUsername() bool {
	if o != nil && o.ReplicationBridgeAuthenticationBasicClientUsername != nil {
		return true
	}

	return false
}

// SetReplicationBridgeAuthenticationBasicClientUsername gets a reference to the given string and assigns it to the ReplicationBridgeAuthenticationBasicClientUsername field.
func (o *MsgVpn) SetReplicationBridgeAuthenticationBasicClientUsername(v string) {
	o.ReplicationBridgeAuthenticationBasicClientUsername = &v
}

// GetReplicationBridgeAuthenticationBasicPassword returns the ReplicationBridgeAuthenticationBasicPassword field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationBridgeAuthenticationBasicPassword() string {
	if o == nil || o.ReplicationBridgeAuthenticationBasicPassword == nil {
		var ret string
		return ret
	}
	return *o.ReplicationBridgeAuthenticationBasicPassword
}

// GetReplicationBridgeAuthenticationBasicPasswordOk returns a tuple with the ReplicationBridgeAuthenticationBasicPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationBridgeAuthenticationBasicPasswordOk() (*string, bool) {
	if o == nil || o.ReplicationBridgeAuthenticationBasicPassword == nil {
		return nil, false
	}
	return o.ReplicationBridgeAuthenticationBasicPassword, true
}

// HasReplicationBridgeAuthenticationBasicPassword returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationBridgeAuthenticationBasicPassword() bool {
	if o != nil && o.ReplicationBridgeAuthenticationBasicPassword != nil {
		return true
	}

	return false
}

// SetReplicationBridgeAuthenticationBasicPassword gets a reference to the given string and assigns it to the ReplicationBridgeAuthenticationBasicPassword field.
func (o *MsgVpn) SetReplicationBridgeAuthenticationBasicPassword(v string) {
	o.ReplicationBridgeAuthenticationBasicPassword = &v
}

// GetReplicationBridgeAuthenticationClientCertContent returns the ReplicationBridgeAuthenticationClientCertContent field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationBridgeAuthenticationClientCertContent() string {
	if o == nil || o.ReplicationBridgeAuthenticationClientCertContent == nil {
		var ret string
		return ret
	}
	return *o.ReplicationBridgeAuthenticationClientCertContent
}

// GetReplicationBridgeAuthenticationClientCertContentOk returns a tuple with the ReplicationBridgeAuthenticationClientCertContent field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationBridgeAuthenticationClientCertContentOk() (*string, bool) {
	if o == nil || o.ReplicationBridgeAuthenticationClientCertContent == nil {
		return nil, false
	}
	return o.ReplicationBridgeAuthenticationClientCertContent, true
}

// HasReplicationBridgeAuthenticationClientCertContent returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationBridgeAuthenticationClientCertContent() bool {
	if o != nil && o.ReplicationBridgeAuthenticationClientCertContent != nil {
		return true
	}

	return false
}

// SetReplicationBridgeAuthenticationClientCertContent gets a reference to the given string and assigns it to the ReplicationBridgeAuthenticationClientCertContent field.
func (o *MsgVpn) SetReplicationBridgeAuthenticationClientCertContent(v string) {
	o.ReplicationBridgeAuthenticationClientCertContent = &v
}

// GetReplicationBridgeAuthenticationClientCertPassword returns the ReplicationBridgeAuthenticationClientCertPassword field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationBridgeAuthenticationClientCertPassword() string {
	if o == nil || o.ReplicationBridgeAuthenticationClientCertPassword == nil {
		var ret string
		return ret
	}
	return *o.ReplicationBridgeAuthenticationClientCertPassword
}

// GetReplicationBridgeAuthenticationClientCertPasswordOk returns a tuple with the ReplicationBridgeAuthenticationClientCertPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationBridgeAuthenticationClientCertPasswordOk() (*string, bool) {
	if o == nil || o.ReplicationBridgeAuthenticationClientCertPassword == nil {
		return nil, false
	}
	return o.ReplicationBridgeAuthenticationClientCertPassword, true
}

// HasReplicationBridgeAuthenticationClientCertPassword returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationBridgeAuthenticationClientCertPassword() bool {
	if o != nil && o.ReplicationBridgeAuthenticationClientCertPassword != nil {
		return true
	}

	return false
}

// SetReplicationBridgeAuthenticationClientCertPassword gets a reference to the given string and assigns it to the ReplicationBridgeAuthenticationClientCertPassword field.
func (o *MsgVpn) SetReplicationBridgeAuthenticationClientCertPassword(v string) {
	o.ReplicationBridgeAuthenticationClientCertPassword = &v
}

// GetReplicationBridgeAuthenticationScheme returns the ReplicationBridgeAuthenticationScheme field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationBridgeAuthenticationScheme() string {
	if o == nil || o.ReplicationBridgeAuthenticationScheme == nil {
		var ret string
		return ret
	}
	return *o.ReplicationBridgeAuthenticationScheme
}

// GetReplicationBridgeAuthenticationSchemeOk returns a tuple with the ReplicationBridgeAuthenticationScheme field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationBridgeAuthenticationSchemeOk() (*string, bool) {
	if o == nil || o.ReplicationBridgeAuthenticationScheme == nil {
		return nil, false
	}
	return o.ReplicationBridgeAuthenticationScheme, true
}

// HasReplicationBridgeAuthenticationScheme returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationBridgeAuthenticationScheme() bool {
	if o != nil && o.ReplicationBridgeAuthenticationScheme != nil {
		return true
	}

	return false
}

// SetReplicationBridgeAuthenticationScheme gets a reference to the given string and assigns it to the ReplicationBridgeAuthenticationScheme field.
func (o *MsgVpn) SetReplicationBridgeAuthenticationScheme(v string) {
	o.ReplicationBridgeAuthenticationScheme = &v
}

// GetReplicationBridgeCompressedDataEnabled returns the ReplicationBridgeCompressedDataEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationBridgeCompressedDataEnabled() bool {
	if o == nil || o.ReplicationBridgeCompressedDataEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ReplicationBridgeCompressedDataEnabled
}

// GetReplicationBridgeCompressedDataEnabledOk returns a tuple with the ReplicationBridgeCompressedDataEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationBridgeCompressedDataEnabledOk() (*bool, bool) {
	if o == nil || o.ReplicationBridgeCompressedDataEnabled == nil {
		return nil, false
	}
	return o.ReplicationBridgeCompressedDataEnabled, true
}

// HasReplicationBridgeCompressedDataEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationBridgeCompressedDataEnabled() bool {
	if o != nil && o.ReplicationBridgeCompressedDataEnabled != nil {
		return true
	}

	return false
}

// SetReplicationBridgeCompressedDataEnabled gets a reference to the given bool and assigns it to the ReplicationBridgeCompressedDataEnabled field.
func (o *MsgVpn) SetReplicationBridgeCompressedDataEnabled(v bool) {
	o.ReplicationBridgeCompressedDataEnabled = &v
}

// GetReplicationBridgeEgressFlowWindowSize returns the ReplicationBridgeEgressFlowWindowSize field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationBridgeEgressFlowWindowSize() int64 {
	if o == nil || o.ReplicationBridgeEgressFlowWindowSize == nil {
		var ret int64
		return ret
	}
	return *o.ReplicationBridgeEgressFlowWindowSize
}

// GetReplicationBridgeEgressFlowWindowSizeOk returns a tuple with the ReplicationBridgeEgressFlowWindowSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationBridgeEgressFlowWindowSizeOk() (*int64, bool) {
	if o == nil || o.ReplicationBridgeEgressFlowWindowSize == nil {
		return nil, false
	}
	return o.ReplicationBridgeEgressFlowWindowSize, true
}

// HasReplicationBridgeEgressFlowWindowSize returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationBridgeEgressFlowWindowSize() bool {
	if o != nil && o.ReplicationBridgeEgressFlowWindowSize != nil {
		return true
	}

	return false
}

// SetReplicationBridgeEgressFlowWindowSize gets a reference to the given int64 and assigns it to the ReplicationBridgeEgressFlowWindowSize field.
func (o *MsgVpn) SetReplicationBridgeEgressFlowWindowSize(v int64) {
	o.ReplicationBridgeEgressFlowWindowSize = &v
}

// GetReplicationBridgeRetryDelay returns the ReplicationBridgeRetryDelay field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationBridgeRetryDelay() int64 {
	if o == nil || o.ReplicationBridgeRetryDelay == nil {
		var ret int64
		return ret
	}
	return *o.ReplicationBridgeRetryDelay
}

// GetReplicationBridgeRetryDelayOk returns a tuple with the ReplicationBridgeRetryDelay field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationBridgeRetryDelayOk() (*int64, bool) {
	if o == nil || o.ReplicationBridgeRetryDelay == nil {
		return nil, false
	}
	return o.ReplicationBridgeRetryDelay, true
}

// HasReplicationBridgeRetryDelay returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationBridgeRetryDelay() bool {
	if o != nil && o.ReplicationBridgeRetryDelay != nil {
		return true
	}

	return false
}

// SetReplicationBridgeRetryDelay gets a reference to the given int64 and assigns it to the ReplicationBridgeRetryDelay field.
func (o *MsgVpn) SetReplicationBridgeRetryDelay(v int64) {
	o.ReplicationBridgeRetryDelay = &v
}

// GetReplicationBridgeTlsEnabled returns the ReplicationBridgeTlsEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationBridgeTlsEnabled() bool {
	if o == nil || o.ReplicationBridgeTlsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ReplicationBridgeTlsEnabled
}

// GetReplicationBridgeTlsEnabledOk returns a tuple with the ReplicationBridgeTlsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationBridgeTlsEnabledOk() (*bool, bool) {
	if o == nil || o.ReplicationBridgeTlsEnabled == nil {
		return nil, false
	}
	return o.ReplicationBridgeTlsEnabled, true
}

// HasReplicationBridgeTlsEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationBridgeTlsEnabled() bool {
	if o != nil && o.ReplicationBridgeTlsEnabled != nil {
		return true
	}

	return false
}

// SetReplicationBridgeTlsEnabled gets a reference to the given bool and assigns it to the ReplicationBridgeTlsEnabled field.
func (o *MsgVpn) SetReplicationBridgeTlsEnabled(v bool) {
	o.ReplicationBridgeTlsEnabled = &v
}

// GetReplicationBridgeUnidirectionalClientProfileName returns the ReplicationBridgeUnidirectionalClientProfileName field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationBridgeUnidirectionalClientProfileName() string {
	if o == nil || o.ReplicationBridgeUnidirectionalClientProfileName == nil {
		var ret string
		return ret
	}
	return *o.ReplicationBridgeUnidirectionalClientProfileName
}

// GetReplicationBridgeUnidirectionalClientProfileNameOk returns a tuple with the ReplicationBridgeUnidirectionalClientProfileName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationBridgeUnidirectionalClientProfileNameOk() (*string, bool) {
	if o == nil || o.ReplicationBridgeUnidirectionalClientProfileName == nil {
		return nil, false
	}
	return o.ReplicationBridgeUnidirectionalClientProfileName, true
}

// HasReplicationBridgeUnidirectionalClientProfileName returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationBridgeUnidirectionalClientProfileName() bool {
	if o != nil && o.ReplicationBridgeUnidirectionalClientProfileName != nil {
		return true
	}

	return false
}

// SetReplicationBridgeUnidirectionalClientProfileName gets a reference to the given string and assigns it to the ReplicationBridgeUnidirectionalClientProfileName field.
func (o *MsgVpn) SetReplicationBridgeUnidirectionalClientProfileName(v string) {
	o.ReplicationBridgeUnidirectionalClientProfileName = &v
}

// GetReplicationEnabled returns the ReplicationEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationEnabled() bool {
	if o == nil || o.ReplicationEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ReplicationEnabled
}

// GetReplicationEnabledOk returns a tuple with the ReplicationEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationEnabledOk() (*bool, bool) {
	if o == nil || o.ReplicationEnabled == nil {
		return nil, false
	}
	return o.ReplicationEnabled, true
}

// HasReplicationEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationEnabled() bool {
	if o != nil && o.ReplicationEnabled != nil {
		return true
	}

	return false
}

// SetReplicationEnabled gets a reference to the given bool and assigns it to the ReplicationEnabled field.
func (o *MsgVpn) SetReplicationEnabled(v bool) {
	o.ReplicationEnabled = &v
}

// GetReplicationEnabledQueueBehavior returns the ReplicationEnabledQueueBehavior field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationEnabledQueueBehavior() string {
	if o == nil || o.ReplicationEnabledQueueBehavior == nil {
		var ret string
		return ret
	}
	return *o.ReplicationEnabledQueueBehavior
}

// GetReplicationEnabledQueueBehaviorOk returns a tuple with the ReplicationEnabledQueueBehavior field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationEnabledQueueBehaviorOk() (*string, bool) {
	if o == nil || o.ReplicationEnabledQueueBehavior == nil {
		return nil, false
	}
	return o.ReplicationEnabledQueueBehavior, true
}

// HasReplicationEnabledQueueBehavior returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationEnabledQueueBehavior() bool {
	if o != nil && o.ReplicationEnabledQueueBehavior != nil {
		return true
	}

	return false
}

// SetReplicationEnabledQueueBehavior gets a reference to the given string and assigns it to the ReplicationEnabledQueueBehavior field.
func (o *MsgVpn) SetReplicationEnabledQueueBehavior(v string) {
	o.ReplicationEnabledQueueBehavior = &v
}

// GetReplicationQueueMaxMsgSpoolUsage returns the ReplicationQueueMaxMsgSpoolUsage field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationQueueMaxMsgSpoolUsage() int64 {
	if o == nil || o.ReplicationQueueMaxMsgSpoolUsage == nil {
		var ret int64
		return ret
	}
	return *o.ReplicationQueueMaxMsgSpoolUsage
}

// GetReplicationQueueMaxMsgSpoolUsageOk returns a tuple with the ReplicationQueueMaxMsgSpoolUsage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationQueueMaxMsgSpoolUsageOk() (*int64, bool) {
	if o == nil || o.ReplicationQueueMaxMsgSpoolUsage == nil {
		return nil, false
	}
	return o.ReplicationQueueMaxMsgSpoolUsage, true
}

// HasReplicationQueueMaxMsgSpoolUsage returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationQueueMaxMsgSpoolUsage() bool {
	if o != nil && o.ReplicationQueueMaxMsgSpoolUsage != nil {
		return true
	}

	return false
}

// SetReplicationQueueMaxMsgSpoolUsage gets a reference to the given int64 and assigns it to the ReplicationQueueMaxMsgSpoolUsage field.
func (o *MsgVpn) SetReplicationQueueMaxMsgSpoolUsage(v int64) {
	o.ReplicationQueueMaxMsgSpoolUsage = &v
}

// GetReplicationQueueRejectMsgToSenderOnDiscardEnabled returns the ReplicationQueueRejectMsgToSenderOnDiscardEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationQueueRejectMsgToSenderOnDiscardEnabled() bool {
	if o == nil || o.ReplicationQueueRejectMsgToSenderOnDiscardEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ReplicationQueueRejectMsgToSenderOnDiscardEnabled
}

// GetReplicationQueueRejectMsgToSenderOnDiscardEnabledOk returns a tuple with the ReplicationQueueRejectMsgToSenderOnDiscardEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationQueueRejectMsgToSenderOnDiscardEnabledOk() (*bool, bool) {
	if o == nil || o.ReplicationQueueRejectMsgToSenderOnDiscardEnabled == nil {
		return nil, false
	}
	return o.ReplicationQueueRejectMsgToSenderOnDiscardEnabled, true
}

// HasReplicationQueueRejectMsgToSenderOnDiscardEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationQueueRejectMsgToSenderOnDiscardEnabled() bool {
	if o != nil && o.ReplicationQueueRejectMsgToSenderOnDiscardEnabled != nil {
		return true
	}

	return false
}

// SetReplicationQueueRejectMsgToSenderOnDiscardEnabled gets a reference to the given bool and assigns it to the ReplicationQueueRejectMsgToSenderOnDiscardEnabled field.
func (o *MsgVpn) SetReplicationQueueRejectMsgToSenderOnDiscardEnabled(v bool) {
	o.ReplicationQueueRejectMsgToSenderOnDiscardEnabled = &v
}

// GetReplicationRejectMsgWhenSyncIneligibleEnabled returns the ReplicationRejectMsgWhenSyncIneligibleEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationRejectMsgWhenSyncIneligibleEnabled() bool {
	if o == nil || o.ReplicationRejectMsgWhenSyncIneligibleEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ReplicationRejectMsgWhenSyncIneligibleEnabled
}

// GetReplicationRejectMsgWhenSyncIneligibleEnabledOk returns a tuple with the ReplicationRejectMsgWhenSyncIneligibleEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationRejectMsgWhenSyncIneligibleEnabledOk() (*bool, bool) {
	if o == nil || o.ReplicationRejectMsgWhenSyncIneligibleEnabled == nil {
		return nil, false
	}
	return o.ReplicationRejectMsgWhenSyncIneligibleEnabled, true
}

// HasReplicationRejectMsgWhenSyncIneligibleEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationRejectMsgWhenSyncIneligibleEnabled() bool {
	if o != nil && o.ReplicationRejectMsgWhenSyncIneligibleEnabled != nil {
		return true
	}

	return false
}

// SetReplicationRejectMsgWhenSyncIneligibleEnabled gets a reference to the given bool and assigns it to the ReplicationRejectMsgWhenSyncIneligibleEnabled field.
func (o *MsgVpn) SetReplicationRejectMsgWhenSyncIneligibleEnabled(v bool) {
	o.ReplicationRejectMsgWhenSyncIneligibleEnabled = &v
}

// GetReplicationRole returns the ReplicationRole field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationRole() string {
	if o == nil || o.ReplicationRole == nil {
		var ret string
		return ret
	}
	return *o.ReplicationRole
}

// GetReplicationRoleOk returns a tuple with the ReplicationRole field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationRoleOk() (*string, bool) {
	if o == nil || o.ReplicationRole == nil {
		return nil, false
	}
	return o.ReplicationRole, true
}

// HasReplicationRole returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationRole() bool {
	if o != nil && o.ReplicationRole != nil {
		return true
	}

	return false
}

// SetReplicationRole gets a reference to the given string and assigns it to the ReplicationRole field.
func (o *MsgVpn) SetReplicationRole(v string) {
	o.ReplicationRole = &v
}

// GetReplicationTransactionMode returns the ReplicationTransactionMode field value if set, zero value otherwise.
func (o *MsgVpn) GetReplicationTransactionMode() string {
	if o == nil || o.ReplicationTransactionMode == nil {
		var ret string
		return ret
	}
	return *o.ReplicationTransactionMode
}

// GetReplicationTransactionModeOk returns a tuple with the ReplicationTransactionMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetReplicationTransactionModeOk() (*string, bool) {
	if o == nil || o.ReplicationTransactionMode == nil {
		return nil, false
	}
	return o.ReplicationTransactionMode, true
}

// HasReplicationTransactionMode returns a boolean if a field has been set.
func (o *MsgVpn) HasReplicationTransactionMode() bool {
	if o != nil && o.ReplicationTransactionMode != nil {
		return true
	}

	return false
}

// SetReplicationTransactionMode gets a reference to the given string and assigns it to the ReplicationTransactionMode field.
func (o *MsgVpn) SetReplicationTransactionMode(v string) {
	o.ReplicationTransactionMode = &v
}

// GetRestTlsServerCertEnforceTrustedCommonNameEnabled returns the RestTlsServerCertEnforceTrustedCommonNameEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetRestTlsServerCertEnforceTrustedCommonNameEnabled() bool {
	if o == nil || o.RestTlsServerCertEnforceTrustedCommonNameEnabled == nil {
		var ret bool
		return ret
	}
	return *o.RestTlsServerCertEnforceTrustedCommonNameEnabled
}

// GetRestTlsServerCertEnforceTrustedCommonNameEnabledOk returns a tuple with the RestTlsServerCertEnforceTrustedCommonNameEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetRestTlsServerCertEnforceTrustedCommonNameEnabledOk() (*bool, bool) {
	if o == nil || o.RestTlsServerCertEnforceTrustedCommonNameEnabled == nil {
		return nil, false
	}
	return o.RestTlsServerCertEnforceTrustedCommonNameEnabled, true
}

// HasRestTlsServerCertEnforceTrustedCommonNameEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasRestTlsServerCertEnforceTrustedCommonNameEnabled() bool {
	if o != nil && o.RestTlsServerCertEnforceTrustedCommonNameEnabled != nil {
		return true
	}

	return false
}

// SetRestTlsServerCertEnforceTrustedCommonNameEnabled gets a reference to the given bool and assigns it to the RestTlsServerCertEnforceTrustedCommonNameEnabled field.
func (o *MsgVpn) SetRestTlsServerCertEnforceTrustedCommonNameEnabled(v bool) {
	o.RestTlsServerCertEnforceTrustedCommonNameEnabled = &v
}

// GetRestTlsServerCertMaxChainDepth returns the RestTlsServerCertMaxChainDepth field value if set, zero value otherwise.
func (o *MsgVpn) GetRestTlsServerCertMaxChainDepth() int64 {
	if o == nil || o.RestTlsServerCertMaxChainDepth == nil {
		var ret int64
		return ret
	}
	return *o.RestTlsServerCertMaxChainDepth
}

// GetRestTlsServerCertMaxChainDepthOk returns a tuple with the RestTlsServerCertMaxChainDepth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetRestTlsServerCertMaxChainDepthOk() (*int64, bool) {
	if o == nil || o.RestTlsServerCertMaxChainDepth == nil {
		return nil, false
	}
	return o.RestTlsServerCertMaxChainDepth, true
}

// HasRestTlsServerCertMaxChainDepth returns a boolean if a field has been set.
func (o *MsgVpn) HasRestTlsServerCertMaxChainDepth() bool {
	if o != nil && o.RestTlsServerCertMaxChainDepth != nil {
		return true
	}

	return false
}

// SetRestTlsServerCertMaxChainDepth gets a reference to the given int64 and assigns it to the RestTlsServerCertMaxChainDepth field.
func (o *MsgVpn) SetRestTlsServerCertMaxChainDepth(v int64) {
	o.RestTlsServerCertMaxChainDepth = &v
}

// GetRestTlsServerCertValidateDateEnabled returns the RestTlsServerCertValidateDateEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetRestTlsServerCertValidateDateEnabled() bool {
	if o == nil || o.RestTlsServerCertValidateDateEnabled == nil {
		var ret bool
		return ret
	}
	return *o.RestTlsServerCertValidateDateEnabled
}

// GetRestTlsServerCertValidateDateEnabledOk returns a tuple with the RestTlsServerCertValidateDateEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetRestTlsServerCertValidateDateEnabledOk() (*bool, bool) {
	if o == nil || o.RestTlsServerCertValidateDateEnabled == nil {
		return nil, false
	}
	return o.RestTlsServerCertValidateDateEnabled, true
}

// HasRestTlsServerCertValidateDateEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasRestTlsServerCertValidateDateEnabled() bool {
	if o != nil && o.RestTlsServerCertValidateDateEnabled != nil {
		return true
	}

	return false
}

// SetRestTlsServerCertValidateDateEnabled gets a reference to the given bool and assigns it to the RestTlsServerCertValidateDateEnabled field.
func (o *MsgVpn) SetRestTlsServerCertValidateDateEnabled(v bool) {
	o.RestTlsServerCertValidateDateEnabled = &v
}

// GetRestTlsServerCertValidateNameEnabled returns the RestTlsServerCertValidateNameEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetRestTlsServerCertValidateNameEnabled() bool {
	if o == nil || o.RestTlsServerCertValidateNameEnabled == nil {
		var ret bool
		return ret
	}
	return *o.RestTlsServerCertValidateNameEnabled
}

// GetRestTlsServerCertValidateNameEnabledOk returns a tuple with the RestTlsServerCertValidateNameEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetRestTlsServerCertValidateNameEnabledOk() (*bool, bool) {
	if o == nil || o.RestTlsServerCertValidateNameEnabled == nil {
		return nil, false
	}
	return o.RestTlsServerCertValidateNameEnabled, true
}

// HasRestTlsServerCertValidateNameEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasRestTlsServerCertValidateNameEnabled() bool {
	if o != nil && o.RestTlsServerCertValidateNameEnabled != nil {
		return true
	}

	return false
}

// SetRestTlsServerCertValidateNameEnabled gets a reference to the given bool and assigns it to the RestTlsServerCertValidateNameEnabled field.
func (o *MsgVpn) SetRestTlsServerCertValidateNameEnabled(v bool) {
	o.RestTlsServerCertValidateNameEnabled = &v
}

// GetSempOverMsgBusAdminClientEnabled returns the SempOverMsgBusAdminClientEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetSempOverMsgBusAdminClientEnabled() bool {
	if o == nil || o.SempOverMsgBusAdminClientEnabled == nil {
		var ret bool
		return ret
	}
	return *o.SempOverMsgBusAdminClientEnabled
}

// GetSempOverMsgBusAdminClientEnabledOk returns a tuple with the SempOverMsgBusAdminClientEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetSempOverMsgBusAdminClientEnabledOk() (*bool, bool) {
	if o == nil || o.SempOverMsgBusAdminClientEnabled == nil {
		return nil, false
	}
	return o.SempOverMsgBusAdminClientEnabled, true
}

// HasSempOverMsgBusAdminClientEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasSempOverMsgBusAdminClientEnabled() bool {
	if o != nil && o.SempOverMsgBusAdminClientEnabled != nil {
		return true
	}

	return false
}

// SetSempOverMsgBusAdminClientEnabled gets a reference to the given bool and assigns it to the SempOverMsgBusAdminClientEnabled field.
func (o *MsgVpn) SetSempOverMsgBusAdminClientEnabled(v bool) {
	o.SempOverMsgBusAdminClientEnabled = &v
}

// GetSempOverMsgBusAdminDistributedCacheEnabled returns the SempOverMsgBusAdminDistributedCacheEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetSempOverMsgBusAdminDistributedCacheEnabled() bool {
	if o == nil || o.SempOverMsgBusAdminDistributedCacheEnabled == nil {
		var ret bool
		return ret
	}
	return *o.SempOverMsgBusAdminDistributedCacheEnabled
}

// GetSempOverMsgBusAdminDistributedCacheEnabledOk returns a tuple with the SempOverMsgBusAdminDistributedCacheEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetSempOverMsgBusAdminDistributedCacheEnabledOk() (*bool, bool) {
	if o == nil || o.SempOverMsgBusAdminDistributedCacheEnabled == nil {
		return nil, false
	}
	return o.SempOverMsgBusAdminDistributedCacheEnabled, true
}

// HasSempOverMsgBusAdminDistributedCacheEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasSempOverMsgBusAdminDistributedCacheEnabled() bool {
	if o != nil && o.SempOverMsgBusAdminDistributedCacheEnabled != nil {
		return true
	}

	return false
}

// SetSempOverMsgBusAdminDistributedCacheEnabled gets a reference to the given bool and assigns it to the SempOverMsgBusAdminDistributedCacheEnabled field.
func (o *MsgVpn) SetSempOverMsgBusAdminDistributedCacheEnabled(v bool) {
	o.SempOverMsgBusAdminDistributedCacheEnabled = &v
}

// GetSempOverMsgBusAdminEnabled returns the SempOverMsgBusAdminEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetSempOverMsgBusAdminEnabled() bool {
	if o == nil || o.SempOverMsgBusAdminEnabled == nil {
		var ret bool
		return ret
	}
	return *o.SempOverMsgBusAdminEnabled
}

// GetSempOverMsgBusAdminEnabledOk returns a tuple with the SempOverMsgBusAdminEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetSempOverMsgBusAdminEnabledOk() (*bool, bool) {
	if o == nil || o.SempOverMsgBusAdminEnabled == nil {
		return nil, false
	}
	return o.SempOverMsgBusAdminEnabled, true
}

// HasSempOverMsgBusAdminEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasSempOverMsgBusAdminEnabled() bool {
	if o != nil && o.SempOverMsgBusAdminEnabled != nil {
		return true
	}

	return false
}

// SetSempOverMsgBusAdminEnabled gets a reference to the given bool and assigns it to the SempOverMsgBusAdminEnabled field.
func (o *MsgVpn) SetSempOverMsgBusAdminEnabled(v bool) {
	o.SempOverMsgBusAdminEnabled = &v
}

// GetSempOverMsgBusEnabled returns the SempOverMsgBusEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetSempOverMsgBusEnabled() bool {
	if o == nil || o.SempOverMsgBusEnabled == nil {
		var ret bool
		return ret
	}
	return *o.SempOverMsgBusEnabled
}

// GetSempOverMsgBusEnabledOk returns a tuple with the SempOverMsgBusEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetSempOverMsgBusEnabledOk() (*bool, bool) {
	if o == nil || o.SempOverMsgBusEnabled == nil {
		return nil, false
	}
	return o.SempOverMsgBusEnabled, true
}

// HasSempOverMsgBusEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasSempOverMsgBusEnabled() bool {
	if o != nil && o.SempOverMsgBusEnabled != nil {
		return true
	}

	return false
}

// SetSempOverMsgBusEnabled gets a reference to the given bool and assigns it to the SempOverMsgBusEnabled field.
func (o *MsgVpn) SetSempOverMsgBusEnabled(v bool) {
	o.SempOverMsgBusEnabled = &v
}

// GetSempOverMsgBusShowEnabled returns the SempOverMsgBusShowEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetSempOverMsgBusShowEnabled() bool {
	if o == nil || o.SempOverMsgBusShowEnabled == nil {
		var ret bool
		return ret
	}
	return *o.SempOverMsgBusShowEnabled
}

// GetSempOverMsgBusShowEnabledOk returns a tuple with the SempOverMsgBusShowEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetSempOverMsgBusShowEnabledOk() (*bool, bool) {
	if o == nil || o.SempOverMsgBusShowEnabled == nil {
		return nil, false
	}
	return o.SempOverMsgBusShowEnabled, true
}

// HasSempOverMsgBusShowEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasSempOverMsgBusShowEnabled() bool {
	if o != nil && o.SempOverMsgBusShowEnabled != nil {
		return true
	}

	return false
}

// SetSempOverMsgBusShowEnabled gets a reference to the given bool and assigns it to the SempOverMsgBusShowEnabled field.
func (o *MsgVpn) SetSempOverMsgBusShowEnabled(v bool) {
	o.SempOverMsgBusShowEnabled = &v
}

// GetServiceAmqpMaxConnectionCount returns the ServiceAmqpMaxConnectionCount field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceAmqpMaxConnectionCount() int64 {
	if o == nil || o.ServiceAmqpMaxConnectionCount == nil {
		var ret int64
		return ret
	}
	return *o.ServiceAmqpMaxConnectionCount
}

// GetServiceAmqpMaxConnectionCountOk returns a tuple with the ServiceAmqpMaxConnectionCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceAmqpMaxConnectionCountOk() (*int64, bool) {
	if o == nil || o.ServiceAmqpMaxConnectionCount == nil {
		return nil, false
	}
	return o.ServiceAmqpMaxConnectionCount, true
}

// HasServiceAmqpMaxConnectionCount returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceAmqpMaxConnectionCount() bool {
	if o != nil && o.ServiceAmqpMaxConnectionCount != nil {
		return true
	}

	return false
}

// SetServiceAmqpMaxConnectionCount gets a reference to the given int64 and assigns it to the ServiceAmqpMaxConnectionCount field.
func (o *MsgVpn) SetServiceAmqpMaxConnectionCount(v int64) {
	o.ServiceAmqpMaxConnectionCount = &v
}

// GetServiceAmqpPlainTextEnabled returns the ServiceAmqpPlainTextEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceAmqpPlainTextEnabled() bool {
	if o == nil || o.ServiceAmqpPlainTextEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceAmqpPlainTextEnabled
}

// GetServiceAmqpPlainTextEnabledOk returns a tuple with the ServiceAmqpPlainTextEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceAmqpPlainTextEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceAmqpPlainTextEnabled == nil {
		return nil, false
	}
	return o.ServiceAmqpPlainTextEnabled, true
}

// HasServiceAmqpPlainTextEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceAmqpPlainTextEnabled() bool {
	if o != nil && o.ServiceAmqpPlainTextEnabled != nil {
		return true
	}

	return false
}

// SetServiceAmqpPlainTextEnabled gets a reference to the given bool and assigns it to the ServiceAmqpPlainTextEnabled field.
func (o *MsgVpn) SetServiceAmqpPlainTextEnabled(v bool) {
	o.ServiceAmqpPlainTextEnabled = &v
}

// GetServiceAmqpPlainTextListenPort returns the ServiceAmqpPlainTextListenPort field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceAmqpPlainTextListenPort() int64 {
	if o == nil || o.ServiceAmqpPlainTextListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceAmqpPlainTextListenPort
}

// GetServiceAmqpPlainTextListenPortOk returns a tuple with the ServiceAmqpPlainTextListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceAmqpPlainTextListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceAmqpPlainTextListenPort == nil {
		return nil, false
	}
	return o.ServiceAmqpPlainTextListenPort, true
}

// HasServiceAmqpPlainTextListenPort returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceAmqpPlainTextListenPort() bool {
	if o != nil && o.ServiceAmqpPlainTextListenPort != nil {
		return true
	}

	return false
}

// SetServiceAmqpPlainTextListenPort gets a reference to the given int64 and assigns it to the ServiceAmqpPlainTextListenPort field.
func (o *MsgVpn) SetServiceAmqpPlainTextListenPort(v int64) {
	o.ServiceAmqpPlainTextListenPort = &v
}

// GetServiceAmqpTlsEnabled returns the ServiceAmqpTlsEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceAmqpTlsEnabled() bool {
	if o == nil || o.ServiceAmqpTlsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceAmqpTlsEnabled
}

// GetServiceAmqpTlsEnabledOk returns a tuple with the ServiceAmqpTlsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceAmqpTlsEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceAmqpTlsEnabled == nil {
		return nil, false
	}
	return o.ServiceAmqpTlsEnabled, true
}

// HasServiceAmqpTlsEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceAmqpTlsEnabled() bool {
	if o != nil && o.ServiceAmqpTlsEnabled != nil {
		return true
	}

	return false
}

// SetServiceAmqpTlsEnabled gets a reference to the given bool and assigns it to the ServiceAmqpTlsEnabled field.
func (o *MsgVpn) SetServiceAmqpTlsEnabled(v bool) {
	o.ServiceAmqpTlsEnabled = &v
}

// GetServiceAmqpTlsListenPort returns the ServiceAmqpTlsListenPort field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceAmqpTlsListenPort() int64 {
	if o == nil || o.ServiceAmqpTlsListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceAmqpTlsListenPort
}

// GetServiceAmqpTlsListenPortOk returns a tuple with the ServiceAmqpTlsListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceAmqpTlsListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceAmqpTlsListenPort == nil {
		return nil, false
	}
	return o.ServiceAmqpTlsListenPort, true
}

// HasServiceAmqpTlsListenPort returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceAmqpTlsListenPort() bool {
	if o != nil && o.ServiceAmqpTlsListenPort != nil {
		return true
	}

	return false
}

// SetServiceAmqpTlsListenPort gets a reference to the given int64 and assigns it to the ServiceAmqpTlsListenPort field.
func (o *MsgVpn) SetServiceAmqpTlsListenPort(v int64) {
	o.ServiceAmqpTlsListenPort = &v
}

// GetServiceMqttAuthenticationClientCertRequest returns the ServiceMqttAuthenticationClientCertRequest field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceMqttAuthenticationClientCertRequest() string {
	if o == nil || o.ServiceMqttAuthenticationClientCertRequest == nil {
		var ret string
		return ret
	}
	return *o.ServiceMqttAuthenticationClientCertRequest
}

// GetServiceMqttAuthenticationClientCertRequestOk returns a tuple with the ServiceMqttAuthenticationClientCertRequest field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceMqttAuthenticationClientCertRequestOk() (*string, bool) {
	if o == nil || o.ServiceMqttAuthenticationClientCertRequest == nil {
		return nil, false
	}
	return o.ServiceMqttAuthenticationClientCertRequest, true
}

// HasServiceMqttAuthenticationClientCertRequest returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceMqttAuthenticationClientCertRequest() bool {
	if o != nil && o.ServiceMqttAuthenticationClientCertRequest != nil {
		return true
	}

	return false
}

// SetServiceMqttAuthenticationClientCertRequest gets a reference to the given string and assigns it to the ServiceMqttAuthenticationClientCertRequest field.
func (o *MsgVpn) SetServiceMqttAuthenticationClientCertRequest(v string) {
	o.ServiceMqttAuthenticationClientCertRequest = &v
}

// GetServiceMqttMaxConnectionCount returns the ServiceMqttMaxConnectionCount field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceMqttMaxConnectionCount() int64 {
	if o == nil || o.ServiceMqttMaxConnectionCount == nil {
		var ret int64
		return ret
	}
	return *o.ServiceMqttMaxConnectionCount
}

// GetServiceMqttMaxConnectionCountOk returns a tuple with the ServiceMqttMaxConnectionCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceMqttMaxConnectionCountOk() (*int64, bool) {
	if o == nil || o.ServiceMqttMaxConnectionCount == nil {
		return nil, false
	}
	return o.ServiceMqttMaxConnectionCount, true
}

// HasServiceMqttMaxConnectionCount returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceMqttMaxConnectionCount() bool {
	if o != nil && o.ServiceMqttMaxConnectionCount != nil {
		return true
	}

	return false
}

// SetServiceMqttMaxConnectionCount gets a reference to the given int64 and assigns it to the ServiceMqttMaxConnectionCount field.
func (o *MsgVpn) SetServiceMqttMaxConnectionCount(v int64) {
	o.ServiceMqttMaxConnectionCount = &v
}

// GetServiceMqttPlainTextEnabled returns the ServiceMqttPlainTextEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceMqttPlainTextEnabled() bool {
	if o == nil || o.ServiceMqttPlainTextEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceMqttPlainTextEnabled
}

// GetServiceMqttPlainTextEnabledOk returns a tuple with the ServiceMqttPlainTextEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceMqttPlainTextEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceMqttPlainTextEnabled == nil {
		return nil, false
	}
	return o.ServiceMqttPlainTextEnabled, true
}

// HasServiceMqttPlainTextEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceMqttPlainTextEnabled() bool {
	if o != nil && o.ServiceMqttPlainTextEnabled != nil {
		return true
	}

	return false
}

// SetServiceMqttPlainTextEnabled gets a reference to the given bool and assigns it to the ServiceMqttPlainTextEnabled field.
func (o *MsgVpn) SetServiceMqttPlainTextEnabled(v bool) {
	o.ServiceMqttPlainTextEnabled = &v
}

// GetServiceMqttPlainTextListenPort returns the ServiceMqttPlainTextListenPort field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceMqttPlainTextListenPort() int64 {
	if o == nil || o.ServiceMqttPlainTextListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceMqttPlainTextListenPort
}

// GetServiceMqttPlainTextListenPortOk returns a tuple with the ServiceMqttPlainTextListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceMqttPlainTextListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceMqttPlainTextListenPort == nil {
		return nil, false
	}
	return o.ServiceMqttPlainTextListenPort, true
}

// HasServiceMqttPlainTextListenPort returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceMqttPlainTextListenPort() bool {
	if o != nil && o.ServiceMqttPlainTextListenPort != nil {
		return true
	}

	return false
}

// SetServiceMqttPlainTextListenPort gets a reference to the given int64 and assigns it to the ServiceMqttPlainTextListenPort field.
func (o *MsgVpn) SetServiceMqttPlainTextListenPort(v int64) {
	o.ServiceMqttPlainTextListenPort = &v
}

// GetServiceMqttTlsEnabled returns the ServiceMqttTlsEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceMqttTlsEnabled() bool {
	if o == nil || o.ServiceMqttTlsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceMqttTlsEnabled
}

// GetServiceMqttTlsEnabledOk returns a tuple with the ServiceMqttTlsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceMqttTlsEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceMqttTlsEnabled == nil {
		return nil, false
	}
	return o.ServiceMqttTlsEnabled, true
}

// HasServiceMqttTlsEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceMqttTlsEnabled() bool {
	if o != nil && o.ServiceMqttTlsEnabled != nil {
		return true
	}

	return false
}

// SetServiceMqttTlsEnabled gets a reference to the given bool and assigns it to the ServiceMqttTlsEnabled field.
func (o *MsgVpn) SetServiceMqttTlsEnabled(v bool) {
	o.ServiceMqttTlsEnabled = &v
}

// GetServiceMqttTlsListenPort returns the ServiceMqttTlsListenPort field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceMqttTlsListenPort() int64 {
	if o == nil || o.ServiceMqttTlsListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceMqttTlsListenPort
}

// GetServiceMqttTlsListenPortOk returns a tuple with the ServiceMqttTlsListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceMqttTlsListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceMqttTlsListenPort == nil {
		return nil, false
	}
	return o.ServiceMqttTlsListenPort, true
}

// HasServiceMqttTlsListenPort returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceMqttTlsListenPort() bool {
	if o != nil && o.ServiceMqttTlsListenPort != nil {
		return true
	}

	return false
}

// SetServiceMqttTlsListenPort gets a reference to the given int64 and assigns it to the ServiceMqttTlsListenPort field.
func (o *MsgVpn) SetServiceMqttTlsListenPort(v int64) {
	o.ServiceMqttTlsListenPort = &v
}

// GetServiceMqttTlsWebSocketEnabled returns the ServiceMqttTlsWebSocketEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceMqttTlsWebSocketEnabled() bool {
	if o == nil || o.ServiceMqttTlsWebSocketEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceMqttTlsWebSocketEnabled
}

// GetServiceMqttTlsWebSocketEnabledOk returns a tuple with the ServiceMqttTlsWebSocketEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceMqttTlsWebSocketEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceMqttTlsWebSocketEnabled == nil {
		return nil, false
	}
	return o.ServiceMqttTlsWebSocketEnabled, true
}

// HasServiceMqttTlsWebSocketEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceMqttTlsWebSocketEnabled() bool {
	if o != nil && o.ServiceMqttTlsWebSocketEnabled != nil {
		return true
	}

	return false
}

// SetServiceMqttTlsWebSocketEnabled gets a reference to the given bool and assigns it to the ServiceMqttTlsWebSocketEnabled field.
func (o *MsgVpn) SetServiceMqttTlsWebSocketEnabled(v bool) {
	o.ServiceMqttTlsWebSocketEnabled = &v
}

// GetServiceMqttTlsWebSocketListenPort returns the ServiceMqttTlsWebSocketListenPort field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceMqttTlsWebSocketListenPort() int64 {
	if o == nil || o.ServiceMqttTlsWebSocketListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceMqttTlsWebSocketListenPort
}

// GetServiceMqttTlsWebSocketListenPortOk returns a tuple with the ServiceMqttTlsWebSocketListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceMqttTlsWebSocketListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceMqttTlsWebSocketListenPort == nil {
		return nil, false
	}
	return o.ServiceMqttTlsWebSocketListenPort, true
}

// HasServiceMqttTlsWebSocketListenPort returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceMqttTlsWebSocketListenPort() bool {
	if o != nil && o.ServiceMqttTlsWebSocketListenPort != nil {
		return true
	}

	return false
}

// SetServiceMqttTlsWebSocketListenPort gets a reference to the given int64 and assigns it to the ServiceMqttTlsWebSocketListenPort field.
func (o *MsgVpn) SetServiceMqttTlsWebSocketListenPort(v int64) {
	o.ServiceMqttTlsWebSocketListenPort = &v
}

// GetServiceMqttWebSocketEnabled returns the ServiceMqttWebSocketEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceMqttWebSocketEnabled() bool {
	if o == nil || o.ServiceMqttWebSocketEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceMqttWebSocketEnabled
}

// GetServiceMqttWebSocketEnabledOk returns a tuple with the ServiceMqttWebSocketEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceMqttWebSocketEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceMqttWebSocketEnabled == nil {
		return nil, false
	}
	return o.ServiceMqttWebSocketEnabled, true
}

// HasServiceMqttWebSocketEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceMqttWebSocketEnabled() bool {
	if o != nil && o.ServiceMqttWebSocketEnabled != nil {
		return true
	}

	return false
}

// SetServiceMqttWebSocketEnabled gets a reference to the given bool and assigns it to the ServiceMqttWebSocketEnabled field.
func (o *MsgVpn) SetServiceMqttWebSocketEnabled(v bool) {
	o.ServiceMqttWebSocketEnabled = &v
}

// GetServiceMqttWebSocketListenPort returns the ServiceMqttWebSocketListenPort field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceMqttWebSocketListenPort() int64 {
	if o == nil || o.ServiceMqttWebSocketListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceMqttWebSocketListenPort
}

// GetServiceMqttWebSocketListenPortOk returns a tuple with the ServiceMqttWebSocketListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceMqttWebSocketListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceMqttWebSocketListenPort == nil {
		return nil, false
	}
	return o.ServiceMqttWebSocketListenPort, true
}

// HasServiceMqttWebSocketListenPort returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceMqttWebSocketListenPort() bool {
	if o != nil && o.ServiceMqttWebSocketListenPort != nil {
		return true
	}

	return false
}

// SetServiceMqttWebSocketListenPort gets a reference to the given int64 and assigns it to the ServiceMqttWebSocketListenPort field.
func (o *MsgVpn) SetServiceMqttWebSocketListenPort(v int64) {
	o.ServiceMqttWebSocketListenPort = &v
}

// GetServiceRestIncomingAuthenticationClientCertRequest returns the ServiceRestIncomingAuthenticationClientCertRequest field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceRestIncomingAuthenticationClientCertRequest() string {
	if o == nil || o.ServiceRestIncomingAuthenticationClientCertRequest == nil {
		var ret string
		return ret
	}
	return *o.ServiceRestIncomingAuthenticationClientCertRequest
}

// GetServiceRestIncomingAuthenticationClientCertRequestOk returns a tuple with the ServiceRestIncomingAuthenticationClientCertRequest field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceRestIncomingAuthenticationClientCertRequestOk() (*string, bool) {
	if o == nil || o.ServiceRestIncomingAuthenticationClientCertRequest == nil {
		return nil, false
	}
	return o.ServiceRestIncomingAuthenticationClientCertRequest, true
}

// HasServiceRestIncomingAuthenticationClientCertRequest returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceRestIncomingAuthenticationClientCertRequest() bool {
	if o != nil && o.ServiceRestIncomingAuthenticationClientCertRequest != nil {
		return true
	}

	return false
}

// SetServiceRestIncomingAuthenticationClientCertRequest gets a reference to the given string and assigns it to the ServiceRestIncomingAuthenticationClientCertRequest field.
func (o *MsgVpn) SetServiceRestIncomingAuthenticationClientCertRequest(v string) {
	o.ServiceRestIncomingAuthenticationClientCertRequest = &v
}

// GetServiceRestIncomingAuthorizationHeaderHandling returns the ServiceRestIncomingAuthorizationHeaderHandling field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceRestIncomingAuthorizationHeaderHandling() string {
	if o == nil || o.ServiceRestIncomingAuthorizationHeaderHandling == nil {
		var ret string
		return ret
	}
	return *o.ServiceRestIncomingAuthorizationHeaderHandling
}

// GetServiceRestIncomingAuthorizationHeaderHandlingOk returns a tuple with the ServiceRestIncomingAuthorizationHeaderHandling field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceRestIncomingAuthorizationHeaderHandlingOk() (*string, bool) {
	if o == nil || o.ServiceRestIncomingAuthorizationHeaderHandling == nil {
		return nil, false
	}
	return o.ServiceRestIncomingAuthorizationHeaderHandling, true
}

// HasServiceRestIncomingAuthorizationHeaderHandling returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceRestIncomingAuthorizationHeaderHandling() bool {
	if o != nil && o.ServiceRestIncomingAuthorizationHeaderHandling != nil {
		return true
	}

	return false
}

// SetServiceRestIncomingAuthorizationHeaderHandling gets a reference to the given string and assigns it to the ServiceRestIncomingAuthorizationHeaderHandling field.
func (o *MsgVpn) SetServiceRestIncomingAuthorizationHeaderHandling(v string) {
	o.ServiceRestIncomingAuthorizationHeaderHandling = &v
}

// GetServiceRestIncomingMaxConnectionCount returns the ServiceRestIncomingMaxConnectionCount field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceRestIncomingMaxConnectionCount() int64 {
	if o == nil || o.ServiceRestIncomingMaxConnectionCount == nil {
		var ret int64
		return ret
	}
	return *o.ServiceRestIncomingMaxConnectionCount
}

// GetServiceRestIncomingMaxConnectionCountOk returns a tuple with the ServiceRestIncomingMaxConnectionCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceRestIncomingMaxConnectionCountOk() (*int64, bool) {
	if o == nil || o.ServiceRestIncomingMaxConnectionCount == nil {
		return nil, false
	}
	return o.ServiceRestIncomingMaxConnectionCount, true
}

// HasServiceRestIncomingMaxConnectionCount returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceRestIncomingMaxConnectionCount() bool {
	if o != nil && o.ServiceRestIncomingMaxConnectionCount != nil {
		return true
	}

	return false
}

// SetServiceRestIncomingMaxConnectionCount gets a reference to the given int64 and assigns it to the ServiceRestIncomingMaxConnectionCount field.
func (o *MsgVpn) SetServiceRestIncomingMaxConnectionCount(v int64) {
	o.ServiceRestIncomingMaxConnectionCount = &v
}

// GetServiceRestIncomingPlainTextEnabled returns the ServiceRestIncomingPlainTextEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceRestIncomingPlainTextEnabled() bool {
	if o == nil || o.ServiceRestIncomingPlainTextEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceRestIncomingPlainTextEnabled
}

// GetServiceRestIncomingPlainTextEnabledOk returns a tuple with the ServiceRestIncomingPlainTextEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceRestIncomingPlainTextEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceRestIncomingPlainTextEnabled == nil {
		return nil, false
	}
	return o.ServiceRestIncomingPlainTextEnabled, true
}

// HasServiceRestIncomingPlainTextEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceRestIncomingPlainTextEnabled() bool {
	if o != nil && o.ServiceRestIncomingPlainTextEnabled != nil {
		return true
	}

	return false
}

// SetServiceRestIncomingPlainTextEnabled gets a reference to the given bool and assigns it to the ServiceRestIncomingPlainTextEnabled field.
func (o *MsgVpn) SetServiceRestIncomingPlainTextEnabled(v bool) {
	o.ServiceRestIncomingPlainTextEnabled = &v
}

// GetServiceRestIncomingPlainTextListenPort returns the ServiceRestIncomingPlainTextListenPort field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceRestIncomingPlainTextListenPort() int64 {
	if o == nil || o.ServiceRestIncomingPlainTextListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceRestIncomingPlainTextListenPort
}

// GetServiceRestIncomingPlainTextListenPortOk returns a tuple with the ServiceRestIncomingPlainTextListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceRestIncomingPlainTextListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceRestIncomingPlainTextListenPort == nil {
		return nil, false
	}
	return o.ServiceRestIncomingPlainTextListenPort, true
}

// HasServiceRestIncomingPlainTextListenPort returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceRestIncomingPlainTextListenPort() bool {
	if o != nil && o.ServiceRestIncomingPlainTextListenPort != nil {
		return true
	}

	return false
}

// SetServiceRestIncomingPlainTextListenPort gets a reference to the given int64 and assigns it to the ServiceRestIncomingPlainTextListenPort field.
func (o *MsgVpn) SetServiceRestIncomingPlainTextListenPort(v int64) {
	o.ServiceRestIncomingPlainTextListenPort = &v
}

// GetServiceRestIncomingTlsEnabled returns the ServiceRestIncomingTlsEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceRestIncomingTlsEnabled() bool {
	if o == nil || o.ServiceRestIncomingTlsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceRestIncomingTlsEnabled
}

// GetServiceRestIncomingTlsEnabledOk returns a tuple with the ServiceRestIncomingTlsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceRestIncomingTlsEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceRestIncomingTlsEnabled == nil {
		return nil, false
	}
	return o.ServiceRestIncomingTlsEnabled, true
}

// HasServiceRestIncomingTlsEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceRestIncomingTlsEnabled() bool {
	if o != nil && o.ServiceRestIncomingTlsEnabled != nil {
		return true
	}

	return false
}

// SetServiceRestIncomingTlsEnabled gets a reference to the given bool and assigns it to the ServiceRestIncomingTlsEnabled field.
func (o *MsgVpn) SetServiceRestIncomingTlsEnabled(v bool) {
	o.ServiceRestIncomingTlsEnabled = &v
}

// GetServiceRestIncomingTlsListenPort returns the ServiceRestIncomingTlsListenPort field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceRestIncomingTlsListenPort() int64 {
	if o == nil || o.ServiceRestIncomingTlsListenPort == nil {
		var ret int64
		return ret
	}
	return *o.ServiceRestIncomingTlsListenPort
}

// GetServiceRestIncomingTlsListenPortOk returns a tuple with the ServiceRestIncomingTlsListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceRestIncomingTlsListenPortOk() (*int64, bool) {
	if o == nil || o.ServiceRestIncomingTlsListenPort == nil {
		return nil, false
	}
	return o.ServiceRestIncomingTlsListenPort, true
}

// HasServiceRestIncomingTlsListenPort returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceRestIncomingTlsListenPort() bool {
	if o != nil && o.ServiceRestIncomingTlsListenPort != nil {
		return true
	}

	return false
}

// SetServiceRestIncomingTlsListenPort gets a reference to the given int64 and assigns it to the ServiceRestIncomingTlsListenPort field.
func (o *MsgVpn) SetServiceRestIncomingTlsListenPort(v int64) {
	o.ServiceRestIncomingTlsListenPort = &v
}

// GetServiceRestMode returns the ServiceRestMode field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceRestMode() string {
	if o == nil || o.ServiceRestMode == nil {
		var ret string
		return ret
	}
	return *o.ServiceRestMode
}

// GetServiceRestModeOk returns a tuple with the ServiceRestMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceRestModeOk() (*string, bool) {
	if o == nil || o.ServiceRestMode == nil {
		return nil, false
	}
	return o.ServiceRestMode, true
}

// HasServiceRestMode returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceRestMode() bool {
	if o != nil && o.ServiceRestMode != nil {
		return true
	}

	return false
}

// SetServiceRestMode gets a reference to the given string and assigns it to the ServiceRestMode field.
func (o *MsgVpn) SetServiceRestMode(v string) {
	o.ServiceRestMode = &v
}

// GetServiceRestOutgoingMaxConnectionCount returns the ServiceRestOutgoingMaxConnectionCount field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceRestOutgoingMaxConnectionCount() int64 {
	if o == nil || o.ServiceRestOutgoingMaxConnectionCount == nil {
		var ret int64
		return ret
	}
	return *o.ServiceRestOutgoingMaxConnectionCount
}

// GetServiceRestOutgoingMaxConnectionCountOk returns a tuple with the ServiceRestOutgoingMaxConnectionCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceRestOutgoingMaxConnectionCountOk() (*int64, bool) {
	if o == nil || o.ServiceRestOutgoingMaxConnectionCount == nil {
		return nil, false
	}
	return o.ServiceRestOutgoingMaxConnectionCount, true
}

// HasServiceRestOutgoingMaxConnectionCount returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceRestOutgoingMaxConnectionCount() bool {
	if o != nil && o.ServiceRestOutgoingMaxConnectionCount != nil {
		return true
	}

	return false
}

// SetServiceRestOutgoingMaxConnectionCount gets a reference to the given int64 and assigns it to the ServiceRestOutgoingMaxConnectionCount field.
func (o *MsgVpn) SetServiceRestOutgoingMaxConnectionCount(v int64) {
	o.ServiceRestOutgoingMaxConnectionCount = &v
}

// GetServiceSmfMaxConnectionCount returns the ServiceSmfMaxConnectionCount field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceSmfMaxConnectionCount() int64 {
	if o == nil || o.ServiceSmfMaxConnectionCount == nil {
		var ret int64
		return ret
	}
	return *o.ServiceSmfMaxConnectionCount
}

// GetServiceSmfMaxConnectionCountOk returns a tuple with the ServiceSmfMaxConnectionCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceSmfMaxConnectionCountOk() (*int64, bool) {
	if o == nil || o.ServiceSmfMaxConnectionCount == nil {
		return nil, false
	}
	return o.ServiceSmfMaxConnectionCount, true
}

// HasServiceSmfMaxConnectionCount returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceSmfMaxConnectionCount() bool {
	if o != nil && o.ServiceSmfMaxConnectionCount != nil {
		return true
	}

	return false
}

// SetServiceSmfMaxConnectionCount gets a reference to the given int64 and assigns it to the ServiceSmfMaxConnectionCount field.
func (o *MsgVpn) SetServiceSmfMaxConnectionCount(v int64) {
	o.ServiceSmfMaxConnectionCount = &v
}

// GetServiceSmfPlainTextEnabled returns the ServiceSmfPlainTextEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceSmfPlainTextEnabled() bool {
	if o == nil || o.ServiceSmfPlainTextEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceSmfPlainTextEnabled
}

// GetServiceSmfPlainTextEnabledOk returns a tuple with the ServiceSmfPlainTextEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceSmfPlainTextEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceSmfPlainTextEnabled == nil {
		return nil, false
	}
	return o.ServiceSmfPlainTextEnabled, true
}

// HasServiceSmfPlainTextEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceSmfPlainTextEnabled() bool {
	if o != nil && o.ServiceSmfPlainTextEnabled != nil {
		return true
	}

	return false
}

// SetServiceSmfPlainTextEnabled gets a reference to the given bool and assigns it to the ServiceSmfPlainTextEnabled field.
func (o *MsgVpn) SetServiceSmfPlainTextEnabled(v bool) {
	o.ServiceSmfPlainTextEnabled = &v
}

// GetServiceSmfTlsEnabled returns the ServiceSmfTlsEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceSmfTlsEnabled() bool {
	if o == nil || o.ServiceSmfTlsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceSmfTlsEnabled
}

// GetServiceSmfTlsEnabledOk returns a tuple with the ServiceSmfTlsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceSmfTlsEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceSmfTlsEnabled == nil {
		return nil, false
	}
	return o.ServiceSmfTlsEnabled, true
}

// HasServiceSmfTlsEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceSmfTlsEnabled() bool {
	if o != nil && o.ServiceSmfTlsEnabled != nil {
		return true
	}

	return false
}

// SetServiceSmfTlsEnabled gets a reference to the given bool and assigns it to the ServiceSmfTlsEnabled field.
func (o *MsgVpn) SetServiceSmfTlsEnabled(v bool) {
	o.ServiceSmfTlsEnabled = &v
}

// GetServiceWebAuthenticationClientCertRequest returns the ServiceWebAuthenticationClientCertRequest field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceWebAuthenticationClientCertRequest() string {
	if o == nil || o.ServiceWebAuthenticationClientCertRequest == nil {
		var ret string
		return ret
	}
	return *o.ServiceWebAuthenticationClientCertRequest
}

// GetServiceWebAuthenticationClientCertRequestOk returns a tuple with the ServiceWebAuthenticationClientCertRequest field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceWebAuthenticationClientCertRequestOk() (*string, bool) {
	if o == nil || o.ServiceWebAuthenticationClientCertRequest == nil {
		return nil, false
	}
	return o.ServiceWebAuthenticationClientCertRequest, true
}

// HasServiceWebAuthenticationClientCertRequest returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceWebAuthenticationClientCertRequest() bool {
	if o != nil && o.ServiceWebAuthenticationClientCertRequest != nil {
		return true
	}

	return false
}

// SetServiceWebAuthenticationClientCertRequest gets a reference to the given string and assigns it to the ServiceWebAuthenticationClientCertRequest field.
func (o *MsgVpn) SetServiceWebAuthenticationClientCertRequest(v string) {
	o.ServiceWebAuthenticationClientCertRequest = &v
}

// GetServiceWebMaxConnectionCount returns the ServiceWebMaxConnectionCount field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceWebMaxConnectionCount() int64 {
	if o == nil || o.ServiceWebMaxConnectionCount == nil {
		var ret int64
		return ret
	}
	return *o.ServiceWebMaxConnectionCount
}

// GetServiceWebMaxConnectionCountOk returns a tuple with the ServiceWebMaxConnectionCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceWebMaxConnectionCountOk() (*int64, bool) {
	if o == nil || o.ServiceWebMaxConnectionCount == nil {
		return nil, false
	}
	return o.ServiceWebMaxConnectionCount, true
}

// HasServiceWebMaxConnectionCount returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceWebMaxConnectionCount() bool {
	if o != nil && o.ServiceWebMaxConnectionCount != nil {
		return true
	}

	return false
}

// SetServiceWebMaxConnectionCount gets a reference to the given int64 and assigns it to the ServiceWebMaxConnectionCount field.
func (o *MsgVpn) SetServiceWebMaxConnectionCount(v int64) {
	o.ServiceWebMaxConnectionCount = &v
}

// GetServiceWebPlainTextEnabled returns the ServiceWebPlainTextEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceWebPlainTextEnabled() bool {
	if o == nil || o.ServiceWebPlainTextEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceWebPlainTextEnabled
}

// GetServiceWebPlainTextEnabledOk returns a tuple with the ServiceWebPlainTextEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceWebPlainTextEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceWebPlainTextEnabled == nil {
		return nil, false
	}
	return o.ServiceWebPlainTextEnabled, true
}

// HasServiceWebPlainTextEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceWebPlainTextEnabled() bool {
	if o != nil && o.ServiceWebPlainTextEnabled != nil {
		return true
	}

	return false
}

// SetServiceWebPlainTextEnabled gets a reference to the given bool and assigns it to the ServiceWebPlainTextEnabled field.
func (o *MsgVpn) SetServiceWebPlainTextEnabled(v bool) {
	o.ServiceWebPlainTextEnabled = &v
}

// GetServiceWebTlsEnabled returns the ServiceWebTlsEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetServiceWebTlsEnabled() bool {
	if o == nil || o.ServiceWebTlsEnabled == nil {
		var ret bool
		return ret
	}
	return *o.ServiceWebTlsEnabled
}

// GetServiceWebTlsEnabledOk returns a tuple with the ServiceWebTlsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetServiceWebTlsEnabledOk() (*bool, bool) {
	if o == nil || o.ServiceWebTlsEnabled == nil {
		return nil, false
	}
	return o.ServiceWebTlsEnabled, true
}

// HasServiceWebTlsEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasServiceWebTlsEnabled() bool {
	if o != nil && o.ServiceWebTlsEnabled != nil {
		return true
	}

	return false
}

// SetServiceWebTlsEnabled gets a reference to the given bool and assigns it to the ServiceWebTlsEnabled field.
func (o *MsgVpn) SetServiceWebTlsEnabled(v bool) {
	o.ServiceWebTlsEnabled = &v
}

// GetTlsAllowDowngradeToPlainTextEnabled returns the TlsAllowDowngradeToPlainTextEnabled field value if set, zero value otherwise.
func (o *MsgVpn) GetTlsAllowDowngradeToPlainTextEnabled() bool {
	if o == nil || o.TlsAllowDowngradeToPlainTextEnabled == nil {
		var ret bool
		return ret
	}
	return *o.TlsAllowDowngradeToPlainTextEnabled
}

// GetTlsAllowDowngradeToPlainTextEnabledOk returns a tuple with the TlsAllowDowngradeToPlainTextEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MsgVpn) GetTlsAllowDowngradeToPlainTextEnabledOk() (*bool, bool) {
	if o == nil || o.TlsAllowDowngradeToPlainTextEnabled == nil {
		return nil, false
	}
	return o.TlsAllowDowngradeToPlainTextEnabled, true
}

// HasTlsAllowDowngradeToPlainTextEnabled returns a boolean if a field has been set.
func (o *MsgVpn) HasTlsAllowDowngradeToPlainTextEnabled() bool {
	if o != nil && o.TlsAllowDowngradeToPlainTextEnabled != nil {
		return true
	}

	return false
}

// SetTlsAllowDowngradeToPlainTextEnabled gets a reference to the given bool and assigns it to the TlsAllowDowngradeToPlainTextEnabled field.
func (o *MsgVpn) SetTlsAllowDowngradeToPlainTextEnabled(v bool) {
	o.TlsAllowDowngradeToPlainTextEnabled = &v
}

func (o MsgVpn) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Alias != nil {
		toSerialize["alias"] = o.Alias
	}
	if o.AuthenticationBasicEnabled != nil {
		toSerialize["authenticationBasicEnabled"] = o.AuthenticationBasicEnabled
	}
	if o.AuthenticationBasicProfileName != nil {
		toSerialize["authenticationBasicProfileName"] = o.AuthenticationBasicProfileName
	}
	if o.AuthenticationBasicRadiusDomain != nil {
		toSerialize["authenticationBasicRadiusDomain"] = o.AuthenticationBasicRadiusDomain
	}
	if o.AuthenticationBasicType != nil {
		toSerialize["authenticationBasicType"] = o.AuthenticationBasicType
	}
	if o.AuthenticationClientCertAllowApiProvidedUsernameEnabled != nil {
		toSerialize["authenticationClientCertAllowApiProvidedUsernameEnabled"] = o.AuthenticationClientCertAllowApiProvidedUsernameEnabled
	}
	if o.AuthenticationClientCertEnabled != nil {
		toSerialize["authenticationClientCertEnabled"] = o.AuthenticationClientCertEnabled
	}
	if o.AuthenticationClientCertMaxChainDepth != nil {
		toSerialize["authenticationClientCertMaxChainDepth"] = o.AuthenticationClientCertMaxChainDepth
	}
	if o.AuthenticationClientCertRevocationCheckMode != nil {
		toSerialize["authenticationClientCertRevocationCheckMode"] = o.AuthenticationClientCertRevocationCheckMode
	}
	if o.AuthenticationClientCertUsernameSource != nil {
		toSerialize["authenticationClientCertUsernameSource"] = o.AuthenticationClientCertUsernameSource
	}
	if o.AuthenticationClientCertValidateDateEnabled != nil {
		toSerialize["authenticationClientCertValidateDateEnabled"] = o.AuthenticationClientCertValidateDateEnabled
	}
	if o.AuthenticationKerberosAllowApiProvidedUsernameEnabled != nil {
		toSerialize["authenticationKerberosAllowApiProvidedUsernameEnabled"] = o.AuthenticationKerberosAllowApiProvidedUsernameEnabled
	}
	if o.AuthenticationKerberosEnabled != nil {
		toSerialize["authenticationKerberosEnabled"] = o.AuthenticationKerberosEnabled
	}
	if o.AuthenticationOauthDefaultProfileName != nil {
		toSerialize["authenticationOauthDefaultProfileName"] = o.AuthenticationOauthDefaultProfileName
	}
	if o.AuthenticationOauthDefaultProviderName != nil {
		toSerialize["authenticationOauthDefaultProviderName"] = o.AuthenticationOauthDefaultProviderName
	}
	if o.AuthenticationOauthEnabled != nil {
		toSerialize["authenticationOauthEnabled"] = o.AuthenticationOauthEnabled
	}
	if o.AuthorizationLdapGroupMembershipAttributeName != nil {
		toSerialize["authorizationLdapGroupMembershipAttributeName"] = o.AuthorizationLdapGroupMembershipAttributeName
	}
	if o.AuthorizationLdapTrimClientUsernameDomainEnabled != nil {
		toSerialize["authorizationLdapTrimClientUsernameDomainEnabled"] = o.AuthorizationLdapTrimClientUsernameDomainEnabled
	}
	if o.AuthorizationProfileName != nil {
		toSerialize["authorizationProfileName"] = o.AuthorizationProfileName
	}
	if o.AuthorizationType != nil {
		toSerialize["authorizationType"] = o.AuthorizationType
	}
	if o.BridgingTlsServerCertEnforceTrustedCommonNameEnabled != nil {
		toSerialize["bridgingTlsServerCertEnforceTrustedCommonNameEnabled"] = o.BridgingTlsServerCertEnforceTrustedCommonNameEnabled
	}
	if o.BridgingTlsServerCertMaxChainDepth != nil {
		toSerialize["bridgingTlsServerCertMaxChainDepth"] = o.BridgingTlsServerCertMaxChainDepth
	}
	if o.BridgingTlsServerCertValidateDateEnabled != nil {
		toSerialize["bridgingTlsServerCertValidateDateEnabled"] = o.BridgingTlsServerCertValidateDateEnabled
	}
	if o.BridgingTlsServerCertValidateNameEnabled != nil {
		toSerialize["bridgingTlsServerCertValidateNameEnabled"] = o.BridgingTlsServerCertValidateNameEnabled
	}
	if o.DistributedCacheManagementEnabled != nil {
		toSerialize["distributedCacheManagementEnabled"] = o.DistributedCacheManagementEnabled
	}
	if o.DmrEnabled != nil {
		toSerialize["dmrEnabled"] = o.DmrEnabled
	}
	if o.Enabled != nil {
		toSerialize["enabled"] = o.Enabled
	}
	if o.EventConnectionCountThreshold != nil {
		toSerialize["eventConnectionCountThreshold"] = o.EventConnectionCountThreshold
	}
	if o.EventEgressFlowCountThreshold != nil {
		toSerialize["eventEgressFlowCountThreshold"] = o.EventEgressFlowCountThreshold
	}
	if o.EventEgressMsgRateThreshold != nil {
		toSerialize["eventEgressMsgRateThreshold"] = o.EventEgressMsgRateThreshold
	}
	if o.EventEndpointCountThreshold != nil {
		toSerialize["eventEndpointCountThreshold"] = o.EventEndpointCountThreshold
	}
	if o.EventIngressFlowCountThreshold != nil {
		toSerialize["eventIngressFlowCountThreshold"] = o.EventIngressFlowCountThreshold
	}
	if o.EventIngressMsgRateThreshold != nil {
		toSerialize["eventIngressMsgRateThreshold"] = o.EventIngressMsgRateThreshold
	}
	if o.EventLargeMsgThreshold != nil {
		toSerialize["eventLargeMsgThreshold"] = o.EventLargeMsgThreshold
	}
	if o.EventLogTag != nil {
		toSerialize["eventLogTag"] = o.EventLogTag
	}
	if o.EventMsgSpoolUsageThreshold != nil {
		toSerialize["eventMsgSpoolUsageThreshold"] = o.EventMsgSpoolUsageThreshold
	}
	if o.EventPublishClientEnabled != nil {
		toSerialize["eventPublishClientEnabled"] = o.EventPublishClientEnabled
	}
	if o.EventPublishMsgVpnEnabled != nil {
		toSerialize["eventPublishMsgVpnEnabled"] = o.EventPublishMsgVpnEnabled
	}
	if o.EventPublishSubscriptionMode != nil {
		toSerialize["eventPublishSubscriptionMode"] = o.EventPublishSubscriptionMode
	}
	if o.EventPublishTopicFormatMqttEnabled != nil {
		toSerialize["eventPublishTopicFormatMqttEnabled"] = o.EventPublishTopicFormatMqttEnabled
	}
	if o.EventPublishTopicFormatSmfEnabled != nil {
		toSerialize["eventPublishTopicFormatSmfEnabled"] = o.EventPublishTopicFormatSmfEnabled
	}
	if o.EventServiceAmqpConnectionCountThreshold != nil {
		toSerialize["eventServiceAmqpConnectionCountThreshold"] = o.EventServiceAmqpConnectionCountThreshold
	}
	if o.EventServiceMqttConnectionCountThreshold != nil {
		toSerialize["eventServiceMqttConnectionCountThreshold"] = o.EventServiceMqttConnectionCountThreshold
	}
	if o.EventServiceRestIncomingConnectionCountThreshold != nil {
		toSerialize["eventServiceRestIncomingConnectionCountThreshold"] = o.EventServiceRestIncomingConnectionCountThreshold
	}
	if o.EventServiceSmfConnectionCountThreshold != nil {
		toSerialize["eventServiceSmfConnectionCountThreshold"] = o.EventServiceSmfConnectionCountThreshold
	}
	if o.EventServiceWebConnectionCountThreshold != nil {
		toSerialize["eventServiceWebConnectionCountThreshold"] = o.EventServiceWebConnectionCountThreshold
	}
	if o.EventSubscriptionCountThreshold != nil {
		toSerialize["eventSubscriptionCountThreshold"] = o.EventSubscriptionCountThreshold
	}
	if o.EventTransactedSessionCountThreshold != nil {
		toSerialize["eventTransactedSessionCountThreshold"] = o.EventTransactedSessionCountThreshold
	}
	if o.EventTransactionCountThreshold != nil {
		toSerialize["eventTransactionCountThreshold"] = o.EventTransactionCountThreshold
	}
	if o.ExportSubscriptionsEnabled != nil {
		toSerialize["exportSubscriptionsEnabled"] = o.ExportSubscriptionsEnabled
	}
	if o.JndiEnabled != nil {
		toSerialize["jndiEnabled"] = o.JndiEnabled
	}
	if o.MaxConnectionCount != nil {
		toSerialize["maxConnectionCount"] = o.MaxConnectionCount
	}
	if o.MaxEgressFlowCount != nil {
		toSerialize["maxEgressFlowCount"] = o.MaxEgressFlowCount
	}
	if o.MaxEndpointCount != nil {
		toSerialize["maxEndpointCount"] = o.MaxEndpointCount
	}
	if o.MaxIngressFlowCount != nil {
		toSerialize["maxIngressFlowCount"] = o.MaxIngressFlowCount
	}
	if o.MaxMsgSpoolUsage != nil {
		toSerialize["maxMsgSpoolUsage"] = o.MaxMsgSpoolUsage
	}
	if o.MaxSubscriptionCount != nil {
		toSerialize["maxSubscriptionCount"] = o.MaxSubscriptionCount
	}
	if o.MaxTransactedSessionCount != nil {
		toSerialize["maxTransactedSessionCount"] = o.MaxTransactedSessionCount
	}
	if o.MaxTransactionCount != nil {
		toSerialize["maxTransactionCount"] = o.MaxTransactionCount
	}
	if o.MqttRetainMaxMemory != nil {
		toSerialize["mqttRetainMaxMemory"] = o.MqttRetainMaxMemory
	}
	if o.MsgVpnName != nil {
		toSerialize["msgVpnName"] = o.MsgVpnName
	}
	if o.ReplicationAckPropagationIntervalMsgCount != nil {
		toSerialize["replicationAckPropagationIntervalMsgCount"] = o.ReplicationAckPropagationIntervalMsgCount
	}
	if o.ReplicationBridgeAuthenticationBasicClientUsername != nil {
		toSerialize["replicationBridgeAuthenticationBasicClientUsername"] = o.ReplicationBridgeAuthenticationBasicClientUsername
	}
	if o.ReplicationBridgeAuthenticationBasicPassword != nil {
		toSerialize["replicationBridgeAuthenticationBasicPassword"] = o.ReplicationBridgeAuthenticationBasicPassword
	}
	if o.ReplicationBridgeAuthenticationClientCertContent != nil {
		toSerialize["replicationBridgeAuthenticationClientCertContent"] = o.ReplicationBridgeAuthenticationClientCertContent
	}
	if o.ReplicationBridgeAuthenticationClientCertPassword != nil {
		toSerialize["replicationBridgeAuthenticationClientCertPassword"] = o.ReplicationBridgeAuthenticationClientCertPassword
	}
	if o.ReplicationBridgeAuthenticationScheme != nil {
		toSerialize["replicationBridgeAuthenticationScheme"] = o.ReplicationBridgeAuthenticationScheme
	}
	if o.ReplicationBridgeCompressedDataEnabled != nil {
		toSerialize["replicationBridgeCompressedDataEnabled"] = o.ReplicationBridgeCompressedDataEnabled
	}
	if o.ReplicationBridgeEgressFlowWindowSize != nil {
		toSerialize["replicationBridgeEgressFlowWindowSize"] = o.ReplicationBridgeEgressFlowWindowSize
	}
	if o.ReplicationBridgeRetryDelay != nil {
		toSerialize["replicationBridgeRetryDelay"] = o.ReplicationBridgeRetryDelay
	}
	if o.ReplicationBridgeTlsEnabled != nil {
		toSerialize["replicationBridgeTlsEnabled"] = o.ReplicationBridgeTlsEnabled
	}
	if o.ReplicationBridgeUnidirectionalClientProfileName != nil {
		toSerialize["replicationBridgeUnidirectionalClientProfileName"] = o.ReplicationBridgeUnidirectionalClientProfileName
	}
	if o.ReplicationEnabled != nil {
		toSerialize["replicationEnabled"] = o.ReplicationEnabled
	}
	if o.ReplicationEnabledQueueBehavior != nil {
		toSerialize["replicationEnabledQueueBehavior"] = o.ReplicationEnabledQueueBehavior
	}
	if o.ReplicationQueueMaxMsgSpoolUsage != nil {
		toSerialize["replicationQueueMaxMsgSpoolUsage"] = o.ReplicationQueueMaxMsgSpoolUsage
	}
	if o.ReplicationQueueRejectMsgToSenderOnDiscardEnabled != nil {
		toSerialize["replicationQueueRejectMsgToSenderOnDiscardEnabled"] = o.ReplicationQueueRejectMsgToSenderOnDiscardEnabled
	}
	if o.ReplicationRejectMsgWhenSyncIneligibleEnabled != nil {
		toSerialize["replicationRejectMsgWhenSyncIneligibleEnabled"] = o.ReplicationRejectMsgWhenSyncIneligibleEnabled
	}
	if o.ReplicationRole != nil {
		toSerialize["replicationRole"] = o.ReplicationRole
	}
	if o.ReplicationTransactionMode != nil {
		toSerialize["replicationTransactionMode"] = o.ReplicationTransactionMode
	}
	if o.RestTlsServerCertEnforceTrustedCommonNameEnabled != nil {
		toSerialize["restTlsServerCertEnforceTrustedCommonNameEnabled"] = o.RestTlsServerCertEnforceTrustedCommonNameEnabled
	}
	if o.RestTlsServerCertMaxChainDepth != nil {
		toSerialize["restTlsServerCertMaxChainDepth"] = o.RestTlsServerCertMaxChainDepth
	}
	if o.RestTlsServerCertValidateDateEnabled != nil {
		toSerialize["restTlsServerCertValidateDateEnabled"] = o.RestTlsServerCertValidateDateEnabled
	}
	if o.RestTlsServerCertValidateNameEnabled != nil {
		toSerialize["restTlsServerCertValidateNameEnabled"] = o.RestTlsServerCertValidateNameEnabled
	}
	if o.SempOverMsgBusAdminClientEnabled != nil {
		toSerialize["sempOverMsgBusAdminClientEnabled"] = o.SempOverMsgBusAdminClientEnabled
	}
	if o.SempOverMsgBusAdminDistributedCacheEnabled != nil {
		toSerialize["sempOverMsgBusAdminDistributedCacheEnabled"] = o.SempOverMsgBusAdminDistributedCacheEnabled
	}
	if o.SempOverMsgBusAdminEnabled != nil {
		toSerialize["sempOverMsgBusAdminEnabled"] = o.SempOverMsgBusAdminEnabled
	}
	if o.SempOverMsgBusEnabled != nil {
		toSerialize["sempOverMsgBusEnabled"] = o.SempOverMsgBusEnabled
	}
	if o.SempOverMsgBusShowEnabled != nil {
		toSerialize["sempOverMsgBusShowEnabled"] = o.SempOverMsgBusShowEnabled
	}
	if o.ServiceAmqpMaxConnectionCount != nil {
		toSerialize["serviceAmqpMaxConnectionCount"] = o.ServiceAmqpMaxConnectionCount
	}
	if o.ServiceAmqpPlainTextEnabled != nil {
		toSerialize["serviceAmqpPlainTextEnabled"] = o.ServiceAmqpPlainTextEnabled
	}
	if o.ServiceAmqpPlainTextListenPort != nil {
		toSerialize["serviceAmqpPlainTextListenPort"] = o.ServiceAmqpPlainTextListenPort
	}
	if o.ServiceAmqpTlsEnabled != nil {
		toSerialize["serviceAmqpTlsEnabled"] = o.ServiceAmqpTlsEnabled
	}
	if o.ServiceAmqpTlsListenPort != nil {
		toSerialize["serviceAmqpTlsListenPort"] = o.ServiceAmqpTlsListenPort
	}
	if o.ServiceMqttAuthenticationClientCertRequest != nil {
		toSerialize["serviceMqttAuthenticationClientCertRequest"] = o.ServiceMqttAuthenticationClientCertRequest
	}
	if o.ServiceMqttMaxConnectionCount != nil {
		toSerialize["serviceMqttMaxConnectionCount"] = o.ServiceMqttMaxConnectionCount
	}
	if o.ServiceMqttPlainTextEnabled != nil {
		toSerialize["serviceMqttPlainTextEnabled"] = o.ServiceMqttPlainTextEnabled
	}
	if o.ServiceMqttPlainTextListenPort != nil {
		toSerialize["serviceMqttPlainTextListenPort"] = o.ServiceMqttPlainTextListenPort
	}
	if o.ServiceMqttTlsEnabled != nil {
		toSerialize["serviceMqttTlsEnabled"] = o.ServiceMqttTlsEnabled
	}
	if o.ServiceMqttTlsListenPort != nil {
		toSerialize["serviceMqttTlsListenPort"] = o.ServiceMqttTlsListenPort
	}
	if o.ServiceMqttTlsWebSocketEnabled != nil {
		toSerialize["serviceMqttTlsWebSocketEnabled"] = o.ServiceMqttTlsWebSocketEnabled
	}
	if o.ServiceMqttTlsWebSocketListenPort != nil {
		toSerialize["serviceMqttTlsWebSocketListenPort"] = o.ServiceMqttTlsWebSocketListenPort
	}
	if o.ServiceMqttWebSocketEnabled != nil {
		toSerialize["serviceMqttWebSocketEnabled"] = o.ServiceMqttWebSocketEnabled
	}
	if o.ServiceMqttWebSocketListenPort != nil {
		toSerialize["serviceMqttWebSocketListenPort"] = o.ServiceMqttWebSocketListenPort
	}
	if o.ServiceRestIncomingAuthenticationClientCertRequest != nil {
		toSerialize["serviceRestIncomingAuthenticationClientCertRequest"] = o.ServiceRestIncomingAuthenticationClientCertRequest
	}
	if o.ServiceRestIncomingAuthorizationHeaderHandling != nil {
		toSerialize["serviceRestIncomingAuthorizationHeaderHandling"] = o.ServiceRestIncomingAuthorizationHeaderHandling
	}
	if o.ServiceRestIncomingMaxConnectionCount != nil {
		toSerialize["serviceRestIncomingMaxConnectionCount"] = o.ServiceRestIncomingMaxConnectionCount
	}
	if o.ServiceRestIncomingPlainTextEnabled != nil {
		toSerialize["serviceRestIncomingPlainTextEnabled"] = o.ServiceRestIncomingPlainTextEnabled
	}
	if o.ServiceRestIncomingPlainTextListenPort != nil {
		toSerialize["serviceRestIncomingPlainTextListenPort"] = o.ServiceRestIncomingPlainTextListenPort
	}
	if o.ServiceRestIncomingTlsEnabled != nil {
		toSerialize["serviceRestIncomingTlsEnabled"] = o.ServiceRestIncomingTlsEnabled
	}
	if o.ServiceRestIncomingTlsListenPort != nil {
		toSerialize["serviceRestIncomingTlsListenPort"] = o.ServiceRestIncomingTlsListenPort
	}
	if o.ServiceRestMode != nil {
		toSerialize["serviceRestMode"] = o.ServiceRestMode
	}
	if o.ServiceRestOutgoingMaxConnectionCount != nil {
		toSerialize["serviceRestOutgoingMaxConnectionCount"] = o.ServiceRestOutgoingMaxConnectionCount
	}
	if o.ServiceSmfMaxConnectionCount != nil {
		toSerialize["serviceSmfMaxConnectionCount"] = o.ServiceSmfMaxConnectionCount
	}
	if o.ServiceSmfPlainTextEnabled != nil {
		toSerialize["serviceSmfPlainTextEnabled"] = o.ServiceSmfPlainTextEnabled
	}
	if o.ServiceSmfTlsEnabled != nil {
		toSerialize["serviceSmfTlsEnabled"] = o.ServiceSmfTlsEnabled
	}
	if o.ServiceWebAuthenticationClientCertRequest != nil {
		toSerialize["serviceWebAuthenticationClientCertRequest"] = o.ServiceWebAuthenticationClientCertRequest
	}
	if o.ServiceWebMaxConnectionCount != nil {
		toSerialize["serviceWebMaxConnectionCount"] = o.ServiceWebMaxConnectionCount
	}
	if o.ServiceWebPlainTextEnabled != nil {
		toSerialize["serviceWebPlainTextEnabled"] = o.ServiceWebPlainTextEnabled
	}
	if o.ServiceWebTlsEnabled != nil {
		toSerialize["serviceWebTlsEnabled"] = o.ServiceWebTlsEnabled
	}
	if o.TlsAllowDowngradeToPlainTextEnabled != nil {
		toSerialize["tlsAllowDowngradeToPlainTextEnabled"] = o.TlsAllowDowngradeToPlainTextEnabled
	}
	return json.Marshal(toSerialize)
}

type NullableMsgVpn struct {
	value *MsgVpn
	isSet bool
}

func (v NullableMsgVpn) Get() *MsgVpn {
	return v.value
}

func (v *NullableMsgVpn) Set(val *MsgVpn) {
	v.value = val
	v.isSet = true
}

func (v NullableMsgVpn) IsSet() bool {
	return v.isSet
}

func (v *NullableMsgVpn) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMsgVpn(val *MsgVpn) *NullableMsgVpn {
	return &NullableMsgVpn{value: val, isSet: true}
}

func (v NullableMsgVpn) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMsgVpn) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
