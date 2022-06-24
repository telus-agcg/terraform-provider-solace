# Go API client for sempv2

SEMP (starting in `v2`, see note 1) is a RESTful API for configuring, monitoring, and administering a Solace PubSub+ broker.

SEMP uses URIs to address manageable **resources** of the Solace PubSub+ broker. Resources are individual **objects**, **collections** of objects, or (exclusively in the action API) **actions**. This document applies to the following API:


API|Base Path|Purpose|Comments
:---|:---|:---|:---
Configuration|/SEMP/v2/config|Reading and writing config state|See note 2



The following APIs are also available:


API|Base Path|Purpose|Comments
:---|:---|:---|:---
Action|/SEMP/v2/action|Performing actions|See note 2
Monitoring|/SEMP/v2/monitor|Querying operational parameters|See note 2



Resources are always nouns, with individual objects being singular and collections being plural.

Objects within a collection are identified by an `obj-id`, which follows the collection name with the form `collection-name/obj-id`.

Actions within an object are identified by an `action-id`, which follows the object name with the form `obj-id/action-id`.

Some examples:

```
/SEMP/v2/config/msgVpns                        ; MsgVpn collection
/SEMP/v2/config/msgVpns/a                      ; MsgVpn object named \"a\"
/SEMP/v2/config/msgVpns/a/queues               ; Queue collection in MsgVpn \"a\"
/SEMP/v2/config/msgVpns/a/queues/b             ; Queue object named \"b\" in MsgVpn \"a\"
/SEMP/v2/action/msgVpns/a/queues/b/startReplay ; Action that starts a replay on Queue \"b\" in MsgVpn \"a\"
/SEMP/v2/monitor/msgVpns/a/clients             ; Client collection in MsgVpn \"a\"
/SEMP/v2/monitor/msgVpns/a/clients/c           ; Client object named \"c\" in MsgVpn \"a\"
```

## Collection Resources

Collections are unordered lists of objects (unless described as otherwise), and are described by JSON arrays. Each item in the array represents an object in the same manner as the individual object would normally be represented. In the configuration API, the creation of a new object is done through its collection resource.

## Object and Action Resources

Objects are composed of attributes, actions, collections, and other objects. They are described by JSON objects as name/value pairs. The collections and actions of an object are not contained directly in the object's JSON content; rather the content includes an attribute containing a URI which points to the collections and actions. These contained resources must be managed through this URI. At a minimum, every object has one or more identifying attributes, and its own `uri` attribute which contains the URI pointing to itself.

Actions are also composed of attributes, and are described by JSON objects as name/value pairs. Unlike objects, however, they are not members of a collection and cannot be retrieved, only performed. Actions only exist in the action API.

Attributes in an object or action may have any combination of the following properties:


Property|Meaning|Comments
:---|:---|:---
Identifying|Attribute is involved in unique identification of the object, and appears in its URI|
Required|Attribute must be provided in the request|
Read-Only|Attribute can only be read, not written.|See note 3
Write-Only|Attribute can only be written, not read, unless the attribute is also opaque|See the documentation for the opaque property
Requires-Disable|Attribute can only be changed when object is disabled|
Deprecated|Attribute is deprecated, and will disappear in the next SEMP version|
Opaque|Attribute can be set or retrieved in opaque form when the `opaquePassword` query parameter is present|See the `opaquePassword` query parameter documentation



In some requests, certain attributes may only be provided in certain combinations with other attributes:


Relationship|Meaning
:---|:---
Requires|Attribute may only be changed by a request if a particular attribute or combination of attributes is also provided in the request
Conflicts|Attribute may only be provided in a request if a particular attribute or combination of attributes is not also provided in the request



In the monitoring API, any non-identifying attribute may not be returned in a GET.

## HTTP Methods

The following HTTP methods manipulate resources in accordance with these general principles. Note that some methods are only used in certain APIs:


Method|Resource|Meaning|Request Body|Response Body|Missing Request Attributes
:---|:---|:---|:---|:---|:---
POST|Collection|Create object|Initial attribute values|Object attributes and metadata|Set to default
PUT|Object|Create or replace object (see note 5)|New attribute values|Object attributes and metadata|Set to default, with certain exceptions (see note 4)
PUT|Action|Performs action|Action arguments|Action metadata|N/A
PATCH|Object|Update object|New attribute values|Object attributes and metadata|unchanged
DELETE|Object|Delete object|Empty|Object metadata|N/A
GET|Object|Get object|Empty|Object attributes and metadata|N/A
GET|Collection|Get collection|Empty|Object attributes and collection metadata|N/A



## Common Query Parameters

The following are some common query parameters that are supported by many method/URI combinations. Individual URIs may document additional parameters. Note that multiple query parameters can be used together in a single URI, separated by the ampersand character. For example:

```
; Request for the MsgVpns collection using two hypothetical query parameters
; \"q1\" and \"q2\" with values \"val1\" and \"val2\" respectively
/SEMP/v2/config/msgVpns?q1=val1&q2=val2
```

### select

Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. Use this query parameter to limit the size of the returned data for each returned object, return only those fields that are desired, or exclude fields that are not desired.

The value of `select` is a comma-separated list of attribute names. If the list contains attribute names that are not prefaced by `-`, only those attributes are included in the response. If the list contains attribute names that are prefaced by `-`, those attributes are excluded from the response. If the list contains both types, then the difference of the first set of attributes and the second set of attributes is returned. If the list is empty (i.e. `select=`), no attributes are returned.

All attributes that are prefaced by `-` must follow all attributes that are not prefaced by `-`. In addition, each attribute name in the list must match at least one attribute in the object.

Names may include the `*` wildcard (zero or more characters). Nested attribute names are supported using periods (e.g. `parentName.childName`).

Some examples:

```
; List of all MsgVpn names
/SEMP/v2/config/msgVpns?select=msgVpnName
; List of all MsgVpn and their attributes except for their names
/SEMP/v2/config/msgVpns?select=-msgVpnName
; Authentication attributes of MsgVpn \"finance\"
/SEMP/v2/config/msgVpns/finance?select=authentication*
; All attributes of MsgVpn \"finance\" except for authentication attributes
/SEMP/v2/config/msgVpns/finance?select=-authentication*
; Access related attributes of Queue \"orderQ\" of MsgVpn \"finance\"
/SEMP/v2/config/msgVpns/finance/queues/orderQ?select=owner,permission
```

### where

Include in the response only objects where certain conditions are true. Use this query parameter to limit which objects are returned to those whose attribute values meet the given conditions.

The value of `where` is a comma-separated list of expressions. All expressions must be true for the object to be included in the response. Each expression takes the form:

```
expression  = attribute-name OP value
OP          = '==' | '!=' | '&lt;' | '&gt;' | '&lt;=' | '&gt;='
```

`value` may be a number, string, `true`, or `false`, as appropriate for the type of `attribute-name`. Greater-than and less-than comparisons only work for numbers. A `*` in a string `value` is interpreted as a wildcard (zero or more characters). Some examples:

```
; Only enabled MsgVpns
/SEMP/v2/config/msgVpns?where=enabled==true
; Only MsgVpns using basic non-LDAP authentication
/SEMP/v2/config/msgVpns?where=authenticationBasicEnabled==true,authenticationBasicType!=ldap
; Only MsgVpns that allow more than 100 client connections
/SEMP/v2/config/msgVpns?where=maxConnectionCount>100
; Only MsgVpns with msgVpnName starting with \"B\":
/SEMP/v2/config/msgVpns?where=msgVpnName==B*
```

### count

Limit the count of objects in the response. This can be useful to limit the size of the response for large collections. The minimum value for `count` is `1` and the default is `10`. There is also a per-collection maximum value to limit request handling time.

`count` does not guarantee that a minimum number of objects will be returned. A page may contain fewer than `count` objects or even be empty. Additional objects may nonetheless be available for retrieval on subsequent pages. See the `cursor` query parameter documentation for more information on paging.

For example:
```
; Up to 25 MsgVpns
/SEMP/v2/config/msgVpns?count=25
```

### cursor

The cursor, or position, for the next page of objects. Cursors are opaque data that should not be created or interpreted by SEMP clients, and should only be used as described below.

When a request is made for a collection and there may be additional objects available for retrieval that are not included in the initial response, the response will include a `cursorQuery` field containing a cursor. The value of this field can be specified in the `cursor` query parameter of a subsequent request to retrieve the next page of objects. For convenience, an appropriate URI is constructed automatically by the broker and included in the `nextPageUri` field of the response. This URI can be used directly to retrieve the next page of objects.

Applications must continue to follow the `nextPageUri` if one is provided in order to retrieve the full set of objects associated with the request, even if a page contains fewer than the requested number of objects (see the `count` query parameter documentation) or is empty.

### opaquePassword

Attributes with the opaque property are also write-only and so cannot normally be retrieved in a GET. However, when a password is provided in the `opaquePassword` query parameter, attributes with the opaque property are retrieved in a GET in opaque form, encrypted with this password. The query parameter can also be used on a POST, PATCH, or PUT to set opaque attributes using opaque attribute values retrieved in a GET, so long as:

1. the same password that was used to retrieve the opaque attribute values is provided; and

2. the broker to which the request is being sent has the same major and minor SEMP version as the broker that produced the opaque attribute values.

The password provided in the query parameter must be a minimum of 8 characters and a maximum of 128 characters.

The query parameter can only be used in the configuration API, and only over HTTPS.

## Authentication

When a client makes its first SEMPv2 request, it must supply a username and password using HTTP Basic authentication, or an OAuth token or tokens using HTTP Bearer authentication.

When HTTP Basic authentication is used, the broker returns a cookie containing a session key. The client can omit the username and password from subsequent requests, because the broker can use the session cookie for authentication instead. When the session expires or is deleted, the client must provide the username and password again, and the broker creates a new session.

There are a limited number of session slots available on the broker. The broker returns 529 No SEMP Session Available if it is not able to allocate a session.

If certain attributes—such as a user's password—are changed, the broker automatically deletes the affected sessions. These attributes are documented below. However, changes in external user configuration data stored on a RADIUS or LDAP server do not trigger the broker to delete the associated session(s), therefore you must do this manually, if required.

A client can retrieve its current session information using the /about/user endpoint and delete its own session using the /about/user/logout endpoint. A client with appropriate permissions can also manage all sessions using the /sessions endpoint.

Sessions are not created when authenticating with an OAuth token or tokens using HTTP Bearer authentication. If a session cookie is provided, it is ignored.

## Help

Visit [our website](https://solace.com) to learn more about Solace.

You can also download the SEMP API specifications by clicking [here](https://solace.com/downloads/).

If you need additional support, please contact us at [support@solace.com](mailto:support@solace.com).

## Notes

Note|Description
:---:|:---
1|This specification defines SEMP starting in \"v2\", and not the original SEMP \"v1\" interface. Request and response formats between \"v1\" and \"v2\" are entirely incompatible, although both protocols share a common port configuration on the Solace PubSub+ broker. They are differentiated by the initial portion of the URI path, one of either \"/SEMP/\" or \"/SEMP/v2/\"
2|This API is partially implemented. Only a subset of all objects are available.
3|Read-only attributes may appear in POST and PUT/PATCH requests. However, if a read-only attribute is not marked as identifying, it will be ignored during a PUT/PATCH.
4|On a PUT, if the SEMP user is not authorized to modify the attribute, its value is left unchanged rather than set to default. In addition, the values of write-only attributes are not set to their defaults on a PUT, except in the following two cases: there is a mutual requires relationship with another non-write-only attribute, both attributes are absent from the request, and the non-write-only attribute is not currently set to its default value; or the attribute is also opaque and the `opaquePassword` query parameter is provided in the request.
5|On a PUT, if the object does not exist, it is created first.



## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 2.26
- Package version: 1.0.0
- Build package: org.openapitools.codegen.languages.GoClientCodegen
For more information, please visit [http://www.solace.com](http://www.solace.com)

## Installation

Install the following dependencies:

```shell
go get github.com/stretchr/testify/assert
go get golang.org/x/oauth2
go get golang.org/x/net/context
```

Put the package under your project folder and add the following in import:

```golang
import sempv2 "github.com/GIT_USER_ID/GIT_REPO_ID/sempv2"
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```golang
os.Setenv("HTTP_PROXY", "http://proxy_name:proxy_port")
```

## Configuration of Server URL

Default configuration comes with `Servers` field that contains server objects as defined in the OpenAPI specification.

### Select Server Configuration

For using other server than the one defined on index 0 set context value `sw.ContextServerIndex` of type `int`.

```golang
ctx := context.WithValue(context.Background(), sempv2.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `sw.ContextServerVariables` of type `map[string]string`.

```golang
ctx := context.WithValue(context.Background(), sempv2.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

### URLs Configuration per Operation

Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"{classname}Service.{nickname}"` string.
Similar rules for overriding default operation server index and variables applies by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```
ctx := context.WithValue(context.Background(), sempv2.ContextOperationServerIndices, map[string]int{
	"{classname}Service.{nickname}": 2,
})
ctx = context.WithValue(context.Background(), sempv2.ContextOperationServerVariables, map[string]map[string]string{
	"{classname}Service.{nickname}": {
		"port": "8443",
	},
})
```

## Documentation for API Endpoints

All URIs are relative to *http://www.solace.com/SEMP/v2/config*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AboutApi* | [**GetAbout**](docs/AboutApi.md#getabout) | **Get** /about | Get an About object.
*AboutApi* | [**GetAboutApi**](docs/AboutApi.md#getaboutapi) | **Get** /about/api | Get an API Description object.
*AboutApi* | [**GetAboutUser**](docs/AboutApi.md#getaboutuser) | **Get** /about/user | Get a User object.
*AboutApi* | [**GetAboutUserMsgVpn**](docs/AboutApi.md#getaboutusermsgvpn) | **Get** /about/user/msgVpns/{msgVpnName} | Get a User Message VPN object.
*AboutApi* | [**GetAboutUserMsgVpns**](docs/AboutApi.md#getaboutusermsgvpns) | **Get** /about/user/msgVpns | Get a list of User Message VPN objects.
*AclProfileApi* | [**CreateMsgVpnAclProfile**](docs/AclProfileApi.md#createmsgvpnaclprofile) | **Post** /msgVpns/{msgVpnName}/aclProfiles | Create an ACL Profile object.
*AclProfileApi* | [**CreateMsgVpnAclProfileClientConnectException**](docs/AclProfileApi.md#createmsgvpnaclprofileclientconnectexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions | Create a Client Connect Exception object.
*AclProfileApi* | [**CreateMsgVpnAclProfilePublishException**](docs/AclProfileApi.md#createmsgvpnaclprofilepublishexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions | Create a Publish Topic Exception object.
*AclProfileApi* | [**CreateMsgVpnAclProfilePublishTopicException**](docs/AclProfileApi.md#createmsgvpnaclprofilepublishtopicexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions | Create a Publish Topic Exception object.
*AclProfileApi* | [**CreateMsgVpnAclProfileSubscribeException**](docs/AclProfileApi.md#createmsgvpnaclprofilesubscribeexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions | Create a Subscribe Topic Exception object.
*AclProfileApi* | [**CreateMsgVpnAclProfileSubscribeShareNameException**](docs/AclProfileApi.md#createmsgvpnaclprofilesubscribesharenameexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions | Create a Subscribe Share Name Exception object.
*AclProfileApi* | [**CreateMsgVpnAclProfileSubscribeTopicException**](docs/AclProfileApi.md#createmsgvpnaclprofilesubscribetopicexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions | Create a Subscribe Topic Exception object.
*AclProfileApi* | [**DeleteMsgVpnAclProfile**](docs/AclProfileApi.md#deletemsgvpnaclprofile) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName} | Delete an ACL Profile object.
*AclProfileApi* | [**DeleteMsgVpnAclProfileClientConnectException**](docs/AclProfileApi.md#deletemsgvpnaclprofileclientconnectexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions/{clientConnectExceptionAddress} | Delete a Client Connect Exception object.
*AclProfileApi* | [**DeleteMsgVpnAclProfilePublishException**](docs/AclProfileApi.md#deletemsgvpnaclprofilepublishexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions/{topicSyntax},{publishExceptionTopic} | Delete a Publish Topic Exception object.
*AclProfileApi* | [**DeleteMsgVpnAclProfilePublishTopicException**](docs/AclProfileApi.md#deletemsgvpnaclprofilepublishtopicexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions/{publishTopicExceptionSyntax},{publishTopicException} | Delete a Publish Topic Exception object.
*AclProfileApi* | [**DeleteMsgVpnAclProfileSubscribeException**](docs/AclProfileApi.md#deletemsgvpnaclprofilesubscribeexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions/{topicSyntax},{subscribeExceptionTopic} | Delete a Subscribe Topic Exception object.
*AclProfileApi* | [**DeleteMsgVpnAclProfileSubscribeShareNameException**](docs/AclProfileApi.md#deletemsgvpnaclprofilesubscribesharenameexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions/{subscribeShareNameExceptionSyntax},{subscribeShareNameException} | Delete a Subscribe Share Name Exception object.
*AclProfileApi* | [**DeleteMsgVpnAclProfileSubscribeTopicException**](docs/AclProfileApi.md#deletemsgvpnaclprofilesubscribetopicexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions/{subscribeTopicExceptionSyntax},{subscribeTopicException} | Delete a Subscribe Topic Exception object.
*AclProfileApi* | [**GetMsgVpnAclProfile**](docs/AclProfileApi.md#getmsgvpnaclprofile) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName} | Get an ACL Profile object.
*AclProfileApi* | [**GetMsgVpnAclProfileClientConnectException**](docs/AclProfileApi.md#getmsgvpnaclprofileclientconnectexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions/{clientConnectExceptionAddress} | Get a Client Connect Exception object.
*AclProfileApi* | [**GetMsgVpnAclProfileClientConnectExceptions**](docs/AclProfileApi.md#getmsgvpnaclprofileclientconnectexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions | Get a list of Client Connect Exception objects.
*AclProfileApi* | [**GetMsgVpnAclProfilePublishException**](docs/AclProfileApi.md#getmsgvpnaclprofilepublishexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions/{topicSyntax},{publishExceptionTopic} | Get a Publish Topic Exception object.
*AclProfileApi* | [**GetMsgVpnAclProfilePublishExceptions**](docs/AclProfileApi.md#getmsgvpnaclprofilepublishexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions | Get a list of Publish Topic Exception objects.
*AclProfileApi* | [**GetMsgVpnAclProfilePublishTopicException**](docs/AclProfileApi.md#getmsgvpnaclprofilepublishtopicexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions/{publishTopicExceptionSyntax},{publishTopicException} | Get a Publish Topic Exception object.
*AclProfileApi* | [**GetMsgVpnAclProfilePublishTopicExceptions**](docs/AclProfileApi.md#getmsgvpnaclprofilepublishtopicexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions | Get a list of Publish Topic Exception objects.
*AclProfileApi* | [**GetMsgVpnAclProfileSubscribeException**](docs/AclProfileApi.md#getmsgvpnaclprofilesubscribeexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions/{topicSyntax},{subscribeExceptionTopic} | Get a Subscribe Topic Exception object.
*AclProfileApi* | [**GetMsgVpnAclProfileSubscribeExceptions**](docs/AclProfileApi.md#getmsgvpnaclprofilesubscribeexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions | Get a list of Subscribe Topic Exception objects.
*AclProfileApi* | [**GetMsgVpnAclProfileSubscribeShareNameException**](docs/AclProfileApi.md#getmsgvpnaclprofilesubscribesharenameexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions/{subscribeShareNameExceptionSyntax},{subscribeShareNameException} | Get a Subscribe Share Name Exception object.
*AclProfileApi* | [**GetMsgVpnAclProfileSubscribeShareNameExceptions**](docs/AclProfileApi.md#getmsgvpnaclprofilesubscribesharenameexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions | Get a list of Subscribe Share Name Exception objects.
*AclProfileApi* | [**GetMsgVpnAclProfileSubscribeTopicException**](docs/AclProfileApi.md#getmsgvpnaclprofilesubscribetopicexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions/{subscribeTopicExceptionSyntax},{subscribeTopicException} | Get a Subscribe Topic Exception object.
*AclProfileApi* | [**GetMsgVpnAclProfileSubscribeTopicExceptions**](docs/AclProfileApi.md#getmsgvpnaclprofilesubscribetopicexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions | Get a list of Subscribe Topic Exception objects.
*AclProfileApi* | [**GetMsgVpnAclProfiles**](docs/AclProfileApi.md#getmsgvpnaclprofiles) | **Get** /msgVpns/{msgVpnName}/aclProfiles | Get a list of ACL Profile objects.
*AclProfileApi* | [**ReplaceMsgVpnAclProfile**](docs/AclProfileApi.md#replacemsgvpnaclprofile) | **Put** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName} | Replace an ACL Profile object.
*AclProfileApi* | [**UpdateMsgVpnAclProfile**](docs/AclProfileApi.md#updatemsgvpnaclprofile) | **Patch** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName} | Update an ACL Profile object.
*AllApi* | [**CreateCertAuthority**](docs/AllApi.md#createcertauthority) | **Post** /certAuthorities | Create a Certificate Authority object.
*AllApi* | [**CreateCertAuthorityOcspTlsTrustedCommonName**](docs/AllApi.md#createcertauthorityocsptlstrustedcommonname) | **Post** /certAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames | Create an OCSP Responder Trusted Common Name object.
*AllApi* | [**CreateClientCertAuthority**](docs/AllApi.md#createclientcertauthority) | **Post** /clientCertAuthorities | Create a Client Certificate Authority object.
*AllApi* | [**CreateClientCertAuthorityOcspTlsTrustedCommonName**](docs/AllApi.md#createclientcertauthorityocsptlstrustedcommonname) | **Post** /clientCertAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames | Create an OCSP Responder Trusted Common Name object.
*AllApi* | [**CreateDmrCluster**](docs/AllApi.md#createdmrcluster) | **Post** /dmrClusters | Create a Cluster object.
*AllApi* | [**CreateDmrClusterLink**](docs/AllApi.md#createdmrclusterlink) | **Post** /dmrClusters/{dmrClusterName}/links | Create a Link object.
*AllApi* | [**CreateDmrClusterLinkRemoteAddress**](docs/AllApi.md#createdmrclusterlinkremoteaddress) | **Post** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/remoteAddresses | Create a Remote Address object.
*AllApi* | [**CreateDmrClusterLinkTlsTrustedCommonName**](docs/AllApi.md#createdmrclusterlinktlstrustedcommonname) | **Post** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/tlsTrustedCommonNames | Create a Trusted Common Name object.
*AllApi* | [**CreateDomainCertAuthority**](docs/AllApi.md#createdomaincertauthority) | **Post** /domainCertAuthorities | Create a Domain Certificate Authority object.
*AllApi* | [**CreateMsgVpn**](docs/AllApi.md#createmsgvpn) | **Post** /msgVpns | Create a Message VPN object.
*AllApi* | [**CreateMsgVpnAclProfile**](docs/AllApi.md#createmsgvpnaclprofile) | **Post** /msgVpns/{msgVpnName}/aclProfiles | Create an ACL Profile object.
*AllApi* | [**CreateMsgVpnAclProfileClientConnectException**](docs/AllApi.md#createmsgvpnaclprofileclientconnectexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions | Create a Client Connect Exception object.
*AllApi* | [**CreateMsgVpnAclProfilePublishException**](docs/AllApi.md#createmsgvpnaclprofilepublishexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions | Create a Publish Topic Exception object.
*AllApi* | [**CreateMsgVpnAclProfilePublishTopicException**](docs/AllApi.md#createmsgvpnaclprofilepublishtopicexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions | Create a Publish Topic Exception object.
*AllApi* | [**CreateMsgVpnAclProfileSubscribeException**](docs/AllApi.md#createmsgvpnaclprofilesubscribeexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions | Create a Subscribe Topic Exception object.
*AllApi* | [**CreateMsgVpnAclProfileSubscribeShareNameException**](docs/AllApi.md#createmsgvpnaclprofilesubscribesharenameexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions | Create a Subscribe Share Name Exception object.
*AllApi* | [**CreateMsgVpnAclProfileSubscribeTopicException**](docs/AllApi.md#createmsgvpnaclprofilesubscribetopicexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions | Create a Subscribe Topic Exception object.
*AllApi* | [**CreateMsgVpnAuthenticationOauthProfile**](docs/AllApi.md#createmsgvpnauthenticationoauthprofile) | **Post** /msgVpns/{msgVpnName}/authenticationOauthProfiles | Create an OAuth Profile object.
*AllApi* | [**CreateMsgVpnAuthenticationOauthProfileClientRequiredClaim**](docs/AllApi.md#createmsgvpnauthenticationoauthprofileclientrequiredclaim) | **Post** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/clientRequiredClaims | Create a Required Claim object.
*AllApi* | [**CreateMsgVpnAuthenticationOauthProfileResourceServerRequiredClaim**](docs/AllApi.md#createmsgvpnauthenticationoauthprofileresourceserverrequiredclaim) | **Post** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/resourceServerRequiredClaims | Create a Required Claim object.
*AllApi* | [**CreateMsgVpnAuthenticationOauthProvider**](docs/AllApi.md#createmsgvpnauthenticationoauthprovider) | **Post** /msgVpns/{msgVpnName}/authenticationOauthProviders | Create an OAuth Provider object.
*AllApi* | [**CreateMsgVpnAuthorizationGroup**](docs/AllApi.md#createmsgvpnauthorizationgroup) | **Post** /msgVpns/{msgVpnName}/authorizationGroups | Create an Authorization Group object.
*AllApi* | [**CreateMsgVpnBridge**](docs/AllApi.md#createmsgvpnbridge) | **Post** /msgVpns/{msgVpnName}/bridges | Create a Bridge object.
*AllApi* | [**CreateMsgVpnBridgeRemoteMsgVpn**](docs/AllApi.md#createmsgvpnbridgeremotemsgvpn) | **Post** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns | Create a Remote Message VPN object.
*AllApi* | [**CreateMsgVpnBridgeRemoteSubscription**](docs/AllApi.md#createmsgvpnbridgeremotesubscription) | **Post** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions | Create a Remote Subscription object.
*AllApi* | [**CreateMsgVpnBridgeTlsTrustedCommonName**](docs/AllApi.md#createmsgvpnbridgetlstrustedcommonname) | **Post** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames | Create a Trusted Common Name object.
*AllApi* | [**CreateMsgVpnClientProfile**](docs/AllApi.md#createmsgvpnclientprofile) | **Post** /msgVpns/{msgVpnName}/clientProfiles | Create a Client Profile object.
*AllApi* | [**CreateMsgVpnClientUsername**](docs/AllApi.md#createmsgvpnclientusername) | **Post** /msgVpns/{msgVpnName}/clientUsernames | Create a Client Username object.
*AllApi* | [**CreateMsgVpnDistributedCache**](docs/AllApi.md#createmsgvpndistributedcache) | **Post** /msgVpns/{msgVpnName}/distributedCaches | Create a Distributed Cache object.
*AllApi* | [**CreateMsgVpnDistributedCacheCluster**](docs/AllApi.md#createmsgvpndistributedcachecluster) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters | Create a Cache Cluster object.
*AllApi* | [**CreateMsgVpnDistributedCacheClusterGlobalCachingHomeCluster**](docs/AllApi.md#createmsgvpndistributedcacheclusterglobalcachinghomecluster) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters | Create a Home Cache Cluster object.
*AllApi* | [**CreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix**](docs/AllApi.md#createmsgvpndistributedcacheclusterglobalcachinghomeclustertopicprefix) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes | Create a Topic Prefix object.
*AllApi* | [**CreateMsgVpnDistributedCacheClusterInstance**](docs/AllApi.md#createmsgvpndistributedcacheclusterinstance) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances | Create a Cache Instance object.
*AllApi* | [**CreateMsgVpnDistributedCacheClusterTopic**](docs/AllApi.md#createmsgvpndistributedcacheclustertopic) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics | Create a Topic object.
*AllApi* | [**CreateMsgVpnDmrBridge**](docs/AllApi.md#createmsgvpndmrbridge) | **Post** /msgVpns/{msgVpnName}/dmrBridges | Create a DMR Bridge object.
*AllApi* | [**CreateMsgVpnJndiConnectionFactory**](docs/AllApi.md#createmsgvpnjndiconnectionfactory) | **Post** /msgVpns/{msgVpnName}/jndiConnectionFactories | Create a JNDI Connection Factory object.
*AllApi* | [**CreateMsgVpnJndiQueue**](docs/AllApi.md#createmsgvpnjndiqueue) | **Post** /msgVpns/{msgVpnName}/jndiQueues | Create a JNDI Queue object.
*AllApi* | [**CreateMsgVpnJndiTopic**](docs/AllApi.md#createmsgvpnjnditopic) | **Post** /msgVpns/{msgVpnName}/jndiTopics | Create a JNDI Topic object.
*AllApi* | [**CreateMsgVpnMqttRetainCache**](docs/AllApi.md#createmsgvpnmqttretaincache) | **Post** /msgVpns/{msgVpnName}/mqttRetainCaches | Create an MQTT Retain Cache object.
*AllApi* | [**CreateMsgVpnMqttSession**](docs/AllApi.md#createmsgvpnmqttsession) | **Post** /msgVpns/{msgVpnName}/mqttSessions | Create an MQTT Session object.
*AllApi* | [**CreateMsgVpnMqttSessionSubscription**](docs/AllApi.md#createmsgvpnmqttsessionsubscription) | **Post** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions | Create a Subscription object.
*AllApi* | [**CreateMsgVpnQueue**](docs/AllApi.md#createmsgvpnqueue) | **Post** /msgVpns/{msgVpnName}/queues | Create a Queue object.
*AllApi* | [**CreateMsgVpnQueueSubscription**](docs/AllApi.md#createmsgvpnqueuesubscription) | **Post** /msgVpns/{msgVpnName}/queues/{queueName}/subscriptions | Create a Queue Subscription object.
*AllApi* | [**CreateMsgVpnQueueTemplate**](docs/AllApi.md#createmsgvpnqueuetemplate) | **Post** /msgVpns/{msgVpnName}/queueTemplates | Create a Queue Template object.
*AllApi* | [**CreateMsgVpnReplayLog**](docs/AllApi.md#createmsgvpnreplaylog) | **Post** /msgVpns/{msgVpnName}/replayLogs | Create a Replay Log object.
*AllApi* | [**CreateMsgVpnReplicatedTopic**](docs/AllApi.md#createmsgvpnreplicatedtopic) | **Post** /msgVpns/{msgVpnName}/replicatedTopics | Create a Replicated Topic object.
*AllApi* | [**CreateMsgVpnRestDeliveryPoint**](docs/AllApi.md#createmsgvpnrestdeliverypoint) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints | Create a REST Delivery Point object.
*AllApi* | [**CreateMsgVpnRestDeliveryPointQueueBinding**](docs/AllApi.md#createmsgvpnrestdeliverypointqueuebinding) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings | Create a Queue Binding object.
*AllApi* | [**CreateMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/AllApi.md#createmsgvpnrestdeliverypointqueuebindingrequestheader) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders | Create a Request Header object.
*AllApi* | [**CreateMsgVpnRestDeliveryPointRestConsumer**](docs/AllApi.md#createmsgvpnrestdeliverypointrestconsumer) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers | Create a REST Consumer object.
*AllApi* | [**CreateMsgVpnRestDeliveryPointRestConsumerOauthJwtClaim**](docs/AllApi.md#createmsgvpnrestdeliverypointrestconsumeroauthjwtclaim) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/oauthJwtClaims | Create a Claim object.
*AllApi* | [**CreateMsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonName**](docs/AllApi.md#createmsgvpnrestdeliverypointrestconsumertlstrustedcommonname) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/tlsTrustedCommonNames | Create a Trusted Common Name object.
*AllApi* | [**CreateMsgVpnSequencedTopic**](docs/AllApi.md#createmsgvpnsequencedtopic) | **Post** /msgVpns/{msgVpnName}/sequencedTopics | Create a Sequenced Topic object.
*AllApi* | [**CreateMsgVpnTopicEndpoint**](docs/AllApi.md#createmsgvpntopicendpoint) | **Post** /msgVpns/{msgVpnName}/topicEndpoints | Create a Topic Endpoint object.
*AllApi* | [**CreateMsgVpnTopicEndpointTemplate**](docs/AllApi.md#createmsgvpntopicendpointtemplate) | **Post** /msgVpns/{msgVpnName}/topicEndpointTemplates | Create a Topic Endpoint Template object.
*AllApi* | [**CreateOauthProfile**](docs/AllApi.md#createoauthprofile) | **Post** /oauthProfiles | Create an OAuth Profile object.
*AllApi* | [**CreateOauthProfileAccessLevelGroup**](docs/AllApi.md#createoauthprofileaccesslevelgroup) | **Post** /oauthProfiles/{oauthProfileName}/accessLevelGroups | Create a Group Access Level object.
*AllApi* | [**CreateOauthProfileAccessLevelGroupMsgVpnAccessLevelException**](docs/AllApi.md#createoauthprofileaccesslevelgroupmsgvpnaccesslevelexception) | **Post** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions | Create a Message VPN Access-Level Exception object.
*AllApi* | [**CreateOauthProfileClientAllowedHost**](docs/AllApi.md#createoauthprofileclientallowedhost) | **Post** /oauthProfiles/{oauthProfileName}/clientAllowedHosts | Create an Allowed Host Value object.
*AllApi* | [**CreateOauthProfileClientAuthorizationParameter**](docs/AllApi.md#createoauthprofileclientauthorizationparameter) | **Post** /oauthProfiles/{oauthProfileName}/clientAuthorizationParameters | Create an Authorization Parameter object.
*AllApi* | [**CreateOauthProfileClientRequiredClaim**](docs/AllApi.md#createoauthprofileclientrequiredclaim) | **Post** /oauthProfiles/{oauthProfileName}/clientRequiredClaims | Create a Required Claim object.
*AllApi* | [**CreateOauthProfileDefaultMsgVpnAccessLevelException**](docs/AllApi.md#createoauthprofiledefaultmsgvpnaccesslevelexception) | **Post** /oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions | Create a Message VPN Access-Level Exception object.
*AllApi* | [**CreateOauthProfileResourceServerRequiredClaim**](docs/AllApi.md#createoauthprofileresourceserverrequiredclaim) | **Post** /oauthProfiles/{oauthProfileName}/resourceServerRequiredClaims | Create a Required Claim object.
*AllApi* | [**CreateVirtualHostname**](docs/AllApi.md#createvirtualhostname) | **Post** /virtualHostnames | Create a Virtual Hostname object.
*AllApi* | [**DeleteCertAuthority**](docs/AllApi.md#deletecertauthority) | **Delete** /certAuthorities/{certAuthorityName} | Delete a Certificate Authority object.
*AllApi* | [**DeleteCertAuthorityOcspTlsTrustedCommonName**](docs/AllApi.md#deletecertauthorityocsptlstrustedcommonname) | **Delete** /certAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames/{ocspTlsTrustedCommonName} | Delete an OCSP Responder Trusted Common Name object.
*AllApi* | [**DeleteClientCertAuthority**](docs/AllApi.md#deleteclientcertauthority) | **Delete** /clientCertAuthorities/{certAuthorityName} | Delete a Client Certificate Authority object.
*AllApi* | [**DeleteClientCertAuthorityOcspTlsTrustedCommonName**](docs/AllApi.md#deleteclientcertauthorityocsptlstrustedcommonname) | **Delete** /clientCertAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames/{ocspTlsTrustedCommonName} | Delete an OCSP Responder Trusted Common Name object.
*AllApi* | [**DeleteDmrCluster**](docs/AllApi.md#deletedmrcluster) | **Delete** /dmrClusters/{dmrClusterName} | Delete a Cluster object.
*AllApi* | [**DeleteDmrClusterLink**](docs/AllApi.md#deletedmrclusterlink) | **Delete** /dmrClusters/{dmrClusterName}/links/{remoteNodeName} | Delete a Link object.
*AllApi* | [**DeleteDmrClusterLinkRemoteAddress**](docs/AllApi.md#deletedmrclusterlinkremoteaddress) | **Delete** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/remoteAddresses/{remoteAddress} | Delete a Remote Address object.
*AllApi* | [**DeleteDmrClusterLinkTlsTrustedCommonName**](docs/AllApi.md#deletedmrclusterlinktlstrustedcommonname) | **Delete** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Delete a Trusted Common Name object.
*AllApi* | [**DeleteDomainCertAuthority**](docs/AllApi.md#deletedomaincertauthority) | **Delete** /domainCertAuthorities/{certAuthorityName} | Delete a Domain Certificate Authority object.
*AllApi* | [**DeleteMsgVpn**](docs/AllApi.md#deletemsgvpn) | **Delete** /msgVpns/{msgVpnName} | Delete a Message VPN object.
*AllApi* | [**DeleteMsgVpnAclProfile**](docs/AllApi.md#deletemsgvpnaclprofile) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName} | Delete an ACL Profile object.
*AllApi* | [**DeleteMsgVpnAclProfileClientConnectException**](docs/AllApi.md#deletemsgvpnaclprofileclientconnectexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions/{clientConnectExceptionAddress} | Delete a Client Connect Exception object.
*AllApi* | [**DeleteMsgVpnAclProfilePublishException**](docs/AllApi.md#deletemsgvpnaclprofilepublishexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions/{topicSyntax},{publishExceptionTopic} | Delete a Publish Topic Exception object.
*AllApi* | [**DeleteMsgVpnAclProfilePublishTopicException**](docs/AllApi.md#deletemsgvpnaclprofilepublishtopicexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions/{publishTopicExceptionSyntax},{publishTopicException} | Delete a Publish Topic Exception object.
*AllApi* | [**DeleteMsgVpnAclProfileSubscribeException**](docs/AllApi.md#deletemsgvpnaclprofilesubscribeexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions/{topicSyntax},{subscribeExceptionTopic} | Delete a Subscribe Topic Exception object.
*AllApi* | [**DeleteMsgVpnAclProfileSubscribeShareNameException**](docs/AllApi.md#deletemsgvpnaclprofilesubscribesharenameexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions/{subscribeShareNameExceptionSyntax},{subscribeShareNameException} | Delete a Subscribe Share Name Exception object.
*AllApi* | [**DeleteMsgVpnAclProfileSubscribeTopicException**](docs/AllApi.md#deletemsgvpnaclprofilesubscribetopicexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions/{subscribeTopicExceptionSyntax},{subscribeTopicException} | Delete a Subscribe Topic Exception object.
*AllApi* | [**DeleteMsgVpnAuthenticationOauthProfile**](docs/AllApi.md#deletemsgvpnauthenticationoauthprofile) | **Delete** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName} | Delete an OAuth Profile object.
*AllApi* | [**DeleteMsgVpnAuthenticationOauthProfileClientRequiredClaim**](docs/AllApi.md#deletemsgvpnauthenticationoauthprofileclientrequiredclaim) | **Delete** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/clientRequiredClaims/{clientRequiredClaimName} | Delete a Required Claim object.
*AllApi* | [**DeleteMsgVpnAuthenticationOauthProfileResourceServerRequiredClaim**](docs/AllApi.md#deletemsgvpnauthenticationoauthprofileresourceserverrequiredclaim) | **Delete** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/resourceServerRequiredClaims/{resourceServerRequiredClaimName} | Delete a Required Claim object.
*AllApi* | [**DeleteMsgVpnAuthenticationOauthProvider**](docs/AllApi.md#deletemsgvpnauthenticationoauthprovider) | **Delete** /msgVpns/{msgVpnName}/authenticationOauthProviders/{oauthProviderName} | Delete an OAuth Provider object.
*AllApi* | [**DeleteMsgVpnAuthorizationGroup**](docs/AllApi.md#deletemsgvpnauthorizationgroup) | **Delete** /msgVpns/{msgVpnName}/authorizationGroups/{authorizationGroupName} | Delete an Authorization Group object.
*AllApi* | [**DeleteMsgVpnBridge**](docs/AllApi.md#deletemsgvpnbridge) | **Delete** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter} | Delete a Bridge object.
*AllApi* | [**DeleteMsgVpnBridgeRemoteMsgVpn**](docs/AllApi.md#deletemsgvpnbridgeremotemsgvpn) | **Delete** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface} | Delete a Remote Message VPN object.
*AllApi* | [**DeleteMsgVpnBridgeRemoteSubscription**](docs/AllApi.md#deletemsgvpnbridgeremotesubscription) | **Delete** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions/{remoteSubscriptionTopic} | Delete a Remote Subscription object.
*AllApi* | [**DeleteMsgVpnBridgeTlsTrustedCommonName**](docs/AllApi.md#deletemsgvpnbridgetlstrustedcommonname) | **Delete** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Delete a Trusted Common Name object.
*AllApi* | [**DeleteMsgVpnClientProfile**](docs/AllApi.md#deletemsgvpnclientprofile) | **Delete** /msgVpns/{msgVpnName}/clientProfiles/{clientProfileName} | Delete a Client Profile object.
*AllApi* | [**DeleteMsgVpnClientUsername**](docs/AllApi.md#deletemsgvpnclientusername) | **Delete** /msgVpns/{msgVpnName}/clientUsernames/{clientUsername} | Delete a Client Username object.
*AllApi* | [**DeleteMsgVpnDistributedCache**](docs/AllApi.md#deletemsgvpndistributedcache) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName} | Delete a Distributed Cache object.
*AllApi* | [**DeleteMsgVpnDistributedCacheCluster**](docs/AllApi.md#deletemsgvpndistributedcachecluster) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName} | Delete a Cache Cluster object.
*AllApi* | [**DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeCluster**](docs/AllApi.md#deletemsgvpndistributedcacheclusterglobalcachinghomecluster) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName} | Delete a Home Cache Cluster object.
*AllApi* | [**DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix**](docs/AllApi.md#deletemsgvpndistributedcacheclusterglobalcachinghomeclustertopicprefix) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes/{topicPrefix} | Delete a Topic Prefix object.
*AllApi* | [**DeleteMsgVpnDistributedCacheClusterInstance**](docs/AllApi.md#deletemsgvpndistributedcacheclusterinstance) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName} | Delete a Cache Instance object.
*AllApi* | [**DeleteMsgVpnDistributedCacheClusterTopic**](docs/AllApi.md#deletemsgvpndistributedcacheclustertopic) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics/{topic} | Delete a Topic object.
*AllApi* | [**DeleteMsgVpnDmrBridge**](docs/AllApi.md#deletemsgvpndmrbridge) | **Delete** /msgVpns/{msgVpnName}/dmrBridges/{remoteNodeName} | Delete a DMR Bridge object.
*AllApi* | [**DeleteMsgVpnJndiConnectionFactory**](docs/AllApi.md#deletemsgvpnjndiconnectionfactory) | **Delete** /msgVpns/{msgVpnName}/jndiConnectionFactories/{connectionFactoryName} | Delete a JNDI Connection Factory object.
*AllApi* | [**DeleteMsgVpnJndiQueue**](docs/AllApi.md#deletemsgvpnjndiqueue) | **Delete** /msgVpns/{msgVpnName}/jndiQueues/{queueName} | Delete a JNDI Queue object.
*AllApi* | [**DeleteMsgVpnJndiTopic**](docs/AllApi.md#deletemsgvpnjnditopic) | **Delete** /msgVpns/{msgVpnName}/jndiTopics/{topicName} | Delete a JNDI Topic object.
*AllApi* | [**DeleteMsgVpnMqttRetainCache**](docs/AllApi.md#deletemsgvpnmqttretaincache) | **Delete** /msgVpns/{msgVpnName}/mqttRetainCaches/{cacheName} | Delete an MQTT Retain Cache object.
*AllApi* | [**DeleteMsgVpnMqttSession**](docs/AllApi.md#deletemsgvpnmqttsession) | **Delete** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter} | Delete an MQTT Session object.
*AllApi* | [**DeleteMsgVpnMqttSessionSubscription**](docs/AllApi.md#deletemsgvpnmqttsessionsubscription) | **Delete** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic} | Delete a Subscription object.
*AllApi* | [**DeleteMsgVpnQueue**](docs/AllApi.md#deletemsgvpnqueue) | **Delete** /msgVpns/{msgVpnName}/queues/{queueName} | Delete a Queue object.
*AllApi* | [**DeleteMsgVpnQueueSubscription**](docs/AllApi.md#deletemsgvpnqueuesubscription) | **Delete** /msgVpns/{msgVpnName}/queues/{queueName}/subscriptions/{subscriptionTopic} | Delete a Queue Subscription object.
*AllApi* | [**DeleteMsgVpnQueueTemplate**](docs/AllApi.md#deletemsgvpnqueuetemplate) | **Delete** /msgVpns/{msgVpnName}/queueTemplates/{queueTemplateName} | Delete a Queue Template object.
*AllApi* | [**DeleteMsgVpnReplayLog**](docs/AllApi.md#deletemsgvpnreplaylog) | **Delete** /msgVpns/{msgVpnName}/replayLogs/{replayLogName} | Delete a Replay Log object.
*AllApi* | [**DeleteMsgVpnReplicatedTopic**](docs/AllApi.md#deletemsgvpnreplicatedtopic) | **Delete** /msgVpns/{msgVpnName}/replicatedTopics/{replicatedTopic} | Delete a Replicated Topic object.
*AllApi* | [**DeleteMsgVpnRestDeliveryPoint**](docs/AllApi.md#deletemsgvpnrestdeliverypoint) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName} | Delete a REST Delivery Point object.
*AllApi* | [**DeleteMsgVpnRestDeliveryPointQueueBinding**](docs/AllApi.md#deletemsgvpnrestdeliverypointqueuebinding) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName} | Delete a Queue Binding object.
*AllApi* | [**DeleteMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/AllApi.md#deletemsgvpnrestdeliverypointqueuebindingrequestheader) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders/{headerName} | Delete a Request Header object.
*AllApi* | [**DeleteMsgVpnRestDeliveryPointRestConsumer**](docs/AllApi.md#deletemsgvpnrestdeliverypointrestconsumer) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName} | Delete a REST Consumer object.
*AllApi* | [**DeleteMsgVpnRestDeliveryPointRestConsumerOauthJwtClaim**](docs/AllApi.md#deletemsgvpnrestdeliverypointrestconsumeroauthjwtclaim) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/oauthJwtClaims/{oauthJwtClaimName} | Delete a Claim object.
*AllApi* | [**DeleteMsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonName**](docs/AllApi.md#deletemsgvpnrestdeliverypointrestconsumertlstrustedcommonname) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Delete a Trusted Common Name object.
*AllApi* | [**DeleteMsgVpnSequencedTopic**](docs/AllApi.md#deletemsgvpnsequencedtopic) | **Delete** /msgVpns/{msgVpnName}/sequencedTopics/{sequencedTopic} | Delete a Sequenced Topic object.
*AllApi* | [**DeleteMsgVpnTopicEndpoint**](docs/AllApi.md#deletemsgvpntopicendpoint) | **Delete** /msgVpns/{msgVpnName}/topicEndpoints/{topicEndpointName} | Delete a Topic Endpoint object.
*AllApi* | [**DeleteMsgVpnTopicEndpointTemplate**](docs/AllApi.md#deletemsgvpntopicendpointtemplate) | **Delete** /msgVpns/{msgVpnName}/topicEndpointTemplates/{topicEndpointTemplateName} | Delete a Topic Endpoint Template object.
*AllApi* | [**DeleteOauthProfile**](docs/AllApi.md#deleteoauthprofile) | **Delete** /oauthProfiles/{oauthProfileName} | Delete an OAuth Profile object.
*AllApi* | [**DeleteOauthProfileAccessLevelGroup**](docs/AllApi.md#deleteoauthprofileaccesslevelgroup) | **Delete** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName} | Delete a Group Access Level object.
*AllApi* | [**DeleteOauthProfileAccessLevelGroupMsgVpnAccessLevelException**](docs/AllApi.md#deleteoauthprofileaccesslevelgroupmsgvpnaccesslevelexception) | **Delete** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions/{msgVpnName} | Delete a Message VPN Access-Level Exception object.
*AllApi* | [**DeleteOauthProfileClientAllowedHost**](docs/AllApi.md#deleteoauthprofileclientallowedhost) | **Delete** /oauthProfiles/{oauthProfileName}/clientAllowedHosts/{allowedHost} | Delete an Allowed Host Value object.
*AllApi* | [**DeleteOauthProfileClientAuthorizationParameter**](docs/AllApi.md#deleteoauthprofileclientauthorizationparameter) | **Delete** /oauthProfiles/{oauthProfileName}/clientAuthorizationParameters/{authorizationParameterName} | Delete an Authorization Parameter object.
*AllApi* | [**DeleteOauthProfileClientRequiredClaim**](docs/AllApi.md#deleteoauthprofileclientrequiredclaim) | **Delete** /oauthProfiles/{oauthProfileName}/clientRequiredClaims/{clientRequiredClaimName} | Delete a Required Claim object.
*AllApi* | [**DeleteOauthProfileDefaultMsgVpnAccessLevelException**](docs/AllApi.md#deleteoauthprofiledefaultmsgvpnaccesslevelexception) | **Delete** /oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions/{msgVpnName} | Delete a Message VPN Access-Level Exception object.
*AllApi* | [**DeleteOauthProfileResourceServerRequiredClaim**](docs/AllApi.md#deleteoauthprofileresourceserverrequiredclaim) | **Delete** /oauthProfiles/{oauthProfileName}/resourceServerRequiredClaims/{resourceServerRequiredClaimName} | Delete a Required Claim object.
*AllApi* | [**DeleteVirtualHostname**](docs/AllApi.md#deletevirtualhostname) | **Delete** /virtualHostnames/{virtualHostname} | Delete a Virtual Hostname object.
*AllApi* | [**GetAbout**](docs/AllApi.md#getabout) | **Get** /about | Get an About object.
*AllApi* | [**GetAboutApi**](docs/AllApi.md#getaboutapi) | **Get** /about/api | Get an API Description object.
*AllApi* | [**GetAboutUser**](docs/AllApi.md#getaboutuser) | **Get** /about/user | Get a User object.
*AllApi* | [**GetAboutUserMsgVpn**](docs/AllApi.md#getaboutusermsgvpn) | **Get** /about/user/msgVpns/{msgVpnName} | Get a User Message VPN object.
*AllApi* | [**GetAboutUserMsgVpns**](docs/AllApi.md#getaboutusermsgvpns) | **Get** /about/user/msgVpns | Get a list of User Message VPN objects.
*AllApi* | [**GetBroker**](docs/AllApi.md#getbroker) | **Get** / | Get a Broker object.
*AllApi* | [**GetCertAuthorities**](docs/AllApi.md#getcertauthorities) | **Get** /certAuthorities | Get a list of Certificate Authority objects.
*AllApi* | [**GetCertAuthority**](docs/AllApi.md#getcertauthority) | **Get** /certAuthorities/{certAuthorityName} | Get a Certificate Authority object.
*AllApi* | [**GetCertAuthorityOcspTlsTrustedCommonName**](docs/AllApi.md#getcertauthorityocsptlstrustedcommonname) | **Get** /certAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames/{ocspTlsTrustedCommonName} | Get an OCSP Responder Trusted Common Name object.
*AllApi* | [**GetCertAuthorityOcspTlsTrustedCommonNames**](docs/AllApi.md#getcertauthorityocsptlstrustedcommonnames) | **Get** /certAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames | Get a list of OCSP Responder Trusted Common Name objects.
*AllApi* | [**GetClientCertAuthorities**](docs/AllApi.md#getclientcertauthorities) | **Get** /clientCertAuthorities | Get a list of Client Certificate Authority objects.
*AllApi* | [**GetClientCertAuthority**](docs/AllApi.md#getclientcertauthority) | **Get** /clientCertAuthorities/{certAuthorityName} | Get a Client Certificate Authority object.
*AllApi* | [**GetClientCertAuthorityOcspTlsTrustedCommonName**](docs/AllApi.md#getclientcertauthorityocsptlstrustedcommonname) | **Get** /clientCertAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames/{ocspTlsTrustedCommonName} | Get an OCSP Responder Trusted Common Name object.
*AllApi* | [**GetClientCertAuthorityOcspTlsTrustedCommonNames**](docs/AllApi.md#getclientcertauthorityocsptlstrustedcommonnames) | **Get** /clientCertAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames | Get a list of OCSP Responder Trusted Common Name objects.
*AllApi* | [**GetDmrCluster**](docs/AllApi.md#getdmrcluster) | **Get** /dmrClusters/{dmrClusterName} | Get a Cluster object.
*AllApi* | [**GetDmrClusterLink**](docs/AllApi.md#getdmrclusterlink) | **Get** /dmrClusters/{dmrClusterName}/links/{remoteNodeName} | Get a Link object.
*AllApi* | [**GetDmrClusterLinkRemoteAddress**](docs/AllApi.md#getdmrclusterlinkremoteaddress) | **Get** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/remoteAddresses/{remoteAddress} | Get a Remote Address object.
*AllApi* | [**GetDmrClusterLinkRemoteAddresses**](docs/AllApi.md#getdmrclusterlinkremoteaddresses) | **Get** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/remoteAddresses | Get a list of Remote Address objects.
*AllApi* | [**GetDmrClusterLinkTlsTrustedCommonName**](docs/AllApi.md#getdmrclusterlinktlstrustedcommonname) | **Get** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Get a Trusted Common Name object.
*AllApi* | [**GetDmrClusterLinkTlsTrustedCommonNames**](docs/AllApi.md#getdmrclusterlinktlstrustedcommonnames) | **Get** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/tlsTrustedCommonNames | Get a list of Trusted Common Name objects.
*AllApi* | [**GetDmrClusterLinks**](docs/AllApi.md#getdmrclusterlinks) | **Get** /dmrClusters/{dmrClusterName}/links | Get a list of Link objects.
*AllApi* | [**GetDmrClusters**](docs/AllApi.md#getdmrclusters) | **Get** /dmrClusters | Get a list of Cluster objects.
*AllApi* | [**GetDomainCertAuthorities**](docs/AllApi.md#getdomaincertauthorities) | **Get** /domainCertAuthorities | Get a list of Domain Certificate Authority objects.
*AllApi* | [**GetDomainCertAuthority**](docs/AllApi.md#getdomaincertauthority) | **Get** /domainCertAuthorities/{certAuthorityName} | Get a Domain Certificate Authority object.
*AllApi* | [**GetMsgVpn**](docs/AllApi.md#getmsgvpn) | **Get** /msgVpns/{msgVpnName} | Get a Message VPN object.
*AllApi* | [**GetMsgVpnAclProfile**](docs/AllApi.md#getmsgvpnaclprofile) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName} | Get an ACL Profile object.
*AllApi* | [**GetMsgVpnAclProfileClientConnectException**](docs/AllApi.md#getmsgvpnaclprofileclientconnectexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions/{clientConnectExceptionAddress} | Get a Client Connect Exception object.
*AllApi* | [**GetMsgVpnAclProfileClientConnectExceptions**](docs/AllApi.md#getmsgvpnaclprofileclientconnectexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions | Get a list of Client Connect Exception objects.
*AllApi* | [**GetMsgVpnAclProfilePublishException**](docs/AllApi.md#getmsgvpnaclprofilepublishexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions/{topicSyntax},{publishExceptionTopic} | Get a Publish Topic Exception object.
*AllApi* | [**GetMsgVpnAclProfilePublishExceptions**](docs/AllApi.md#getmsgvpnaclprofilepublishexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions | Get a list of Publish Topic Exception objects.
*AllApi* | [**GetMsgVpnAclProfilePublishTopicException**](docs/AllApi.md#getmsgvpnaclprofilepublishtopicexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions/{publishTopicExceptionSyntax},{publishTopicException} | Get a Publish Topic Exception object.
*AllApi* | [**GetMsgVpnAclProfilePublishTopicExceptions**](docs/AllApi.md#getmsgvpnaclprofilepublishtopicexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions | Get a list of Publish Topic Exception objects.
*AllApi* | [**GetMsgVpnAclProfileSubscribeException**](docs/AllApi.md#getmsgvpnaclprofilesubscribeexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions/{topicSyntax},{subscribeExceptionTopic} | Get a Subscribe Topic Exception object.
*AllApi* | [**GetMsgVpnAclProfileSubscribeExceptions**](docs/AllApi.md#getmsgvpnaclprofilesubscribeexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions | Get a list of Subscribe Topic Exception objects.
*AllApi* | [**GetMsgVpnAclProfileSubscribeShareNameException**](docs/AllApi.md#getmsgvpnaclprofilesubscribesharenameexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions/{subscribeShareNameExceptionSyntax},{subscribeShareNameException} | Get a Subscribe Share Name Exception object.
*AllApi* | [**GetMsgVpnAclProfileSubscribeShareNameExceptions**](docs/AllApi.md#getmsgvpnaclprofilesubscribesharenameexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions | Get a list of Subscribe Share Name Exception objects.
*AllApi* | [**GetMsgVpnAclProfileSubscribeTopicException**](docs/AllApi.md#getmsgvpnaclprofilesubscribetopicexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions/{subscribeTopicExceptionSyntax},{subscribeTopicException} | Get a Subscribe Topic Exception object.
*AllApi* | [**GetMsgVpnAclProfileSubscribeTopicExceptions**](docs/AllApi.md#getmsgvpnaclprofilesubscribetopicexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions | Get a list of Subscribe Topic Exception objects.
*AllApi* | [**GetMsgVpnAclProfiles**](docs/AllApi.md#getmsgvpnaclprofiles) | **Get** /msgVpns/{msgVpnName}/aclProfiles | Get a list of ACL Profile objects.
*AllApi* | [**GetMsgVpnAuthenticationOauthProfile**](docs/AllApi.md#getmsgvpnauthenticationoauthprofile) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName} | Get an OAuth Profile object.
*AllApi* | [**GetMsgVpnAuthenticationOauthProfileClientRequiredClaim**](docs/AllApi.md#getmsgvpnauthenticationoauthprofileclientrequiredclaim) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/clientRequiredClaims/{clientRequiredClaimName} | Get a Required Claim object.
*AllApi* | [**GetMsgVpnAuthenticationOauthProfileClientRequiredClaims**](docs/AllApi.md#getmsgvpnauthenticationoauthprofileclientrequiredclaims) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/clientRequiredClaims | Get a list of Required Claim objects.
*AllApi* | [**GetMsgVpnAuthenticationOauthProfileResourceServerRequiredClaim**](docs/AllApi.md#getmsgvpnauthenticationoauthprofileresourceserverrequiredclaim) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/resourceServerRequiredClaims/{resourceServerRequiredClaimName} | Get a Required Claim object.
*AllApi* | [**GetMsgVpnAuthenticationOauthProfileResourceServerRequiredClaims**](docs/AllApi.md#getmsgvpnauthenticationoauthprofileresourceserverrequiredclaims) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/resourceServerRequiredClaims | Get a list of Required Claim objects.
*AllApi* | [**GetMsgVpnAuthenticationOauthProfiles**](docs/AllApi.md#getmsgvpnauthenticationoauthprofiles) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles | Get a list of OAuth Profile objects.
*AllApi* | [**GetMsgVpnAuthenticationOauthProvider**](docs/AllApi.md#getmsgvpnauthenticationoauthprovider) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProviders/{oauthProviderName} | Get an OAuth Provider object.
*AllApi* | [**GetMsgVpnAuthenticationOauthProviders**](docs/AllApi.md#getmsgvpnauthenticationoauthproviders) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProviders | Get a list of OAuth Provider objects.
*AllApi* | [**GetMsgVpnAuthorizationGroup**](docs/AllApi.md#getmsgvpnauthorizationgroup) | **Get** /msgVpns/{msgVpnName}/authorizationGroups/{authorizationGroupName} | Get an Authorization Group object.
*AllApi* | [**GetMsgVpnAuthorizationGroups**](docs/AllApi.md#getmsgvpnauthorizationgroups) | **Get** /msgVpns/{msgVpnName}/authorizationGroups | Get a list of Authorization Group objects.
*AllApi* | [**GetMsgVpnBridge**](docs/AllApi.md#getmsgvpnbridge) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter} | Get a Bridge object.
*AllApi* | [**GetMsgVpnBridgeRemoteMsgVpn**](docs/AllApi.md#getmsgvpnbridgeremotemsgvpn) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface} | Get a Remote Message VPN object.
*AllApi* | [**GetMsgVpnBridgeRemoteMsgVpns**](docs/AllApi.md#getmsgvpnbridgeremotemsgvpns) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns | Get a list of Remote Message VPN objects.
*AllApi* | [**GetMsgVpnBridgeRemoteSubscription**](docs/AllApi.md#getmsgvpnbridgeremotesubscription) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions/{remoteSubscriptionTopic} | Get a Remote Subscription object.
*AllApi* | [**GetMsgVpnBridgeRemoteSubscriptions**](docs/AllApi.md#getmsgvpnbridgeremotesubscriptions) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions | Get a list of Remote Subscription objects.
*AllApi* | [**GetMsgVpnBridgeTlsTrustedCommonName**](docs/AllApi.md#getmsgvpnbridgetlstrustedcommonname) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Get a Trusted Common Name object.
*AllApi* | [**GetMsgVpnBridgeTlsTrustedCommonNames**](docs/AllApi.md#getmsgvpnbridgetlstrustedcommonnames) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames | Get a list of Trusted Common Name objects.
*AllApi* | [**GetMsgVpnBridges**](docs/AllApi.md#getmsgvpnbridges) | **Get** /msgVpns/{msgVpnName}/bridges | Get a list of Bridge objects.
*AllApi* | [**GetMsgVpnClientProfile**](docs/AllApi.md#getmsgvpnclientprofile) | **Get** /msgVpns/{msgVpnName}/clientProfiles/{clientProfileName} | Get a Client Profile object.
*AllApi* | [**GetMsgVpnClientProfiles**](docs/AllApi.md#getmsgvpnclientprofiles) | **Get** /msgVpns/{msgVpnName}/clientProfiles | Get a list of Client Profile objects.
*AllApi* | [**GetMsgVpnClientUsername**](docs/AllApi.md#getmsgvpnclientusername) | **Get** /msgVpns/{msgVpnName}/clientUsernames/{clientUsername} | Get a Client Username object.
*AllApi* | [**GetMsgVpnClientUsernames**](docs/AllApi.md#getmsgvpnclientusernames) | **Get** /msgVpns/{msgVpnName}/clientUsernames | Get a list of Client Username objects.
*AllApi* | [**GetMsgVpnDistributedCache**](docs/AllApi.md#getmsgvpndistributedcache) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName} | Get a Distributed Cache object.
*AllApi* | [**GetMsgVpnDistributedCacheCluster**](docs/AllApi.md#getmsgvpndistributedcachecluster) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName} | Get a Cache Cluster object.
*AllApi* | [**GetMsgVpnDistributedCacheClusterGlobalCachingHomeCluster**](docs/AllApi.md#getmsgvpndistributedcacheclusterglobalcachinghomecluster) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName} | Get a Home Cache Cluster object.
*AllApi* | [**GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix**](docs/AllApi.md#getmsgvpndistributedcacheclusterglobalcachinghomeclustertopicprefix) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes/{topicPrefix} | Get a Topic Prefix object.
*AllApi* | [**GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixes**](docs/AllApi.md#getmsgvpndistributedcacheclusterglobalcachinghomeclustertopicprefixes) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes | Get a list of Topic Prefix objects.
*AllApi* | [**GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusters**](docs/AllApi.md#getmsgvpndistributedcacheclusterglobalcachinghomeclusters) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters | Get a list of Home Cache Cluster objects.
*AllApi* | [**GetMsgVpnDistributedCacheClusterInstance**](docs/AllApi.md#getmsgvpndistributedcacheclusterinstance) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName} | Get a Cache Instance object.
*AllApi* | [**GetMsgVpnDistributedCacheClusterInstances**](docs/AllApi.md#getmsgvpndistributedcacheclusterinstances) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances | Get a list of Cache Instance objects.
*AllApi* | [**GetMsgVpnDistributedCacheClusterTopic**](docs/AllApi.md#getmsgvpndistributedcacheclustertopic) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics/{topic} | Get a Topic object.
*AllApi* | [**GetMsgVpnDistributedCacheClusterTopics**](docs/AllApi.md#getmsgvpndistributedcacheclustertopics) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics | Get a list of Topic objects.
*AllApi* | [**GetMsgVpnDistributedCacheClusters**](docs/AllApi.md#getmsgvpndistributedcacheclusters) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters | Get a list of Cache Cluster objects.
*AllApi* | [**GetMsgVpnDistributedCaches**](docs/AllApi.md#getmsgvpndistributedcaches) | **Get** /msgVpns/{msgVpnName}/distributedCaches | Get a list of Distributed Cache objects.
*AllApi* | [**GetMsgVpnDmrBridge**](docs/AllApi.md#getmsgvpndmrbridge) | **Get** /msgVpns/{msgVpnName}/dmrBridges/{remoteNodeName} | Get a DMR Bridge object.
*AllApi* | [**GetMsgVpnDmrBridges**](docs/AllApi.md#getmsgvpndmrbridges) | **Get** /msgVpns/{msgVpnName}/dmrBridges | Get a list of DMR Bridge objects.
*AllApi* | [**GetMsgVpnJndiConnectionFactories**](docs/AllApi.md#getmsgvpnjndiconnectionfactories) | **Get** /msgVpns/{msgVpnName}/jndiConnectionFactories | Get a list of JNDI Connection Factory objects.
*AllApi* | [**GetMsgVpnJndiConnectionFactory**](docs/AllApi.md#getmsgvpnjndiconnectionfactory) | **Get** /msgVpns/{msgVpnName}/jndiConnectionFactories/{connectionFactoryName} | Get a JNDI Connection Factory object.
*AllApi* | [**GetMsgVpnJndiQueue**](docs/AllApi.md#getmsgvpnjndiqueue) | **Get** /msgVpns/{msgVpnName}/jndiQueues/{queueName} | Get a JNDI Queue object.
*AllApi* | [**GetMsgVpnJndiQueues**](docs/AllApi.md#getmsgvpnjndiqueues) | **Get** /msgVpns/{msgVpnName}/jndiQueues | Get a list of JNDI Queue objects.
*AllApi* | [**GetMsgVpnJndiTopic**](docs/AllApi.md#getmsgvpnjnditopic) | **Get** /msgVpns/{msgVpnName}/jndiTopics/{topicName} | Get a JNDI Topic object.
*AllApi* | [**GetMsgVpnJndiTopics**](docs/AllApi.md#getmsgvpnjnditopics) | **Get** /msgVpns/{msgVpnName}/jndiTopics | Get a list of JNDI Topic objects.
*AllApi* | [**GetMsgVpnMqttRetainCache**](docs/AllApi.md#getmsgvpnmqttretaincache) | **Get** /msgVpns/{msgVpnName}/mqttRetainCaches/{cacheName} | Get an MQTT Retain Cache object.
*AllApi* | [**GetMsgVpnMqttRetainCaches**](docs/AllApi.md#getmsgvpnmqttretaincaches) | **Get** /msgVpns/{msgVpnName}/mqttRetainCaches | Get a list of MQTT Retain Cache objects.
*AllApi* | [**GetMsgVpnMqttSession**](docs/AllApi.md#getmsgvpnmqttsession) | **Get** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter} | Get an MQTT Session object.
*AllApi* | [**GetMsgVpnMqttSessionSubscription**](docs/AllApi.md#getmsgvpnmqttsessionsubscription) | **Get** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic} | Get a Subscription object.
*AllApi* | [**GetMsgVpnMqttSessionSubscriptions**](docs/AllApi.md#getmsgvpnmqttsessionsubscriptions) | **Get** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions | Get a list of Subscription objects.
*AllApi* | [**GetMsgVpnMqttSessions**](docs/AllApi.md#getmsgvpnmqttsessions) | **Get** /msgVpns/{msgVpnName}/mqttSessions | Get a list of MQTT Session objects.
*AllApi* | [**GetMsgVpnQueue**](docs/AllApi.md#getmsgvpnqueue) | **Get** /msgVpns/{msgVpnName}/queues/{queueName} | Get a Queue object.
*AllApi* | [**GetMsgVpnQueueSubscription**](docs/AllApi.md#getmsgvpnqueuesubscription) | **Get** /msgVpns/{msgVpnName}/queues/{queueName}/subscriptions/{subscriptionTopic} | Get a Queue Subscription object.
*AllApi* | [**GetMsgVpnQueueSubscriptions**](docs/AllApi.md#getmsgvpnqueuesubscriptions) | **Get** /msgVpns/{msgVpnName}/queues/{queueName}/subscriptions | Get a list of Queue Subscription objects.
*AllApi* | [**GetMsgVpnQueueTemplate**](docs/AllApi.md#getmsgvpnqueuetemplate) | **Get** /msgVpns/{msgVpnName}/queueTemplates/{queueTemplateName} | Get a Queue Template object.
*AllApi* | [**GetMsgVpnQueueTemplates**](docs/AllApi.md#getmsgvpnqueuetemplates) | **Get** /msgVpns/{msgVpnName}/queueTemplates | Get a list of Queue Template objects.
*AllApi* | [**GetMsgVpnQueues**](docs/AllApi.md#getmsgvpnqueues) | **Get** /msgVpns/{msgVpnName}/queues | Get a list of Queue objects.
*AllApi* | [**GetMsgVpnReplayLog**](docs/AllApi.md#getmsgvpnreplaylog) | **Get** /msgVpns/{msgVpnName}/replayLogs/{replayLogName} | Get a Replay Log object.
*AllApi* | [**GetMsgVpnReplayLogs**](docs/AllApi.md#getmsgvpnreplaylogs) | **Get** /msgVpns/{msgVpnName}/replayLogs | Get a list of Replay Log objects.
*AllApi* | [**GetMsgVpnReplicatedTopic**](docs/AllApi.md#getmsgvpnreplicatedtopic) | **Get** /msgVpns/{msgVpnName}/replicatedTopics/{replicatedTopic} | Get a Replicated Topic object.
*AllApi* | [**GetMsgVpnReplicatedTopics**](docs/AllApi.md#getmsgvpnreplicatedtopics) | **Get** /msgVpns/{msgVpnName}/replicatedTopics | Get a list of Replicated Topic objects.
*AllApi* | [**GetMsgVpnRestDeliveryPoint**](docs/AllApi.md#getmsgvpnrestdeliverypoint) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName} | Get a REST Delivery Point object.
*AllApi* | [**GetMsgVpnRestDeliveryPointQueueBinding**](docs/AllApi.md#getmsgvpnrestdeliverypointqueuebinding) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName} | Get a Queue Binding object.
*AllApi* | [**GetMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/AllApi.md#getmsgvpnrestdeliverypointqueuebindingrequestheader) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders/{headerName} | Get a Request Header object.
*AllApi* | [**GetMsgVpnRestDeliveryPointQueueBindingRequestHeaders**](docs/AllApi.md#getmsgvpnrestdeliverypointqueuebindingrequestheaders) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders | Get a list of Request Header objects.
*AllApi* | [**GetMsgVpnRestDeliveryPointQueueBindings**](docs/AllApi.md#getmsgvpnrestdeliverypointqueuebindings) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings | Get a list of Queue Binding objects.
*AllApi* | [**GetMsgVpnRestDeliveryPointRestConsumer**](docs/AllApi.md#getmsgvpnrestdeliverypointrestconsumer) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName} | Get a REST Consumer object.
*AllApi* | [**GetMsgVpnRestDeliveryPointRestConsumerOauthJwtClaim**](docs/AllApi.md#getmsgvpnrestdeliverypointrestconsumeroauthjwtclaim) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/oauthJwtClaims/{oauthJwtClaimName} | Get a Claim object.
*AllApi* | [**GetMsgVpnRestDeliveryPointRestConsumerOauthJwtClaims**](docs/AllApi.md#getmsgvpnrestdeliverypointrestconsumeroauthjwtclaims) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/oauthJwtClaims | Get a list of Claim objects.
*AllApi* | [**GetMsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonName**](docs/AllApi.md#getmsgvpnrestdeliverypointrestconsumertlstrustedcommonname) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Get a Trusted Common Name object.
*AllApi* | [**GetMsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonNames**](docs/AllApi.md#getmsgvpnrestdeliverypointrestconsumertlstrustedcommonnames) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/tlsTrustedCommonNames | Get a list of Trusted Common Name objects.
*AllApi* | [**GetMsgVpnRestDeliveryPointRestConsumers**](docs/AllApi.md#getmsgvpnrestdeliverypointrestconsumers) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers | Get a list of REST Consumer objects.
*AllApi* | [**GetMsgVpnRestDeliveryPoints**](docs/AllApi.md#getmsgvpnrestdeliverypoints) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints | Get a list of REST Delivery Point objects.
*AllApi* | [**GetMsgVpnSequencedTopic**](docs/AllApi.md#getmsgvpnsequencedtopic) | **Get** /msgVpns/{msgVpnName}/sequencedTopics/{sequencedTopic} | Get a Sequenced Topic object.
*AllApi* | [**GetMsgVpnSequencedTopics**](docs/AllApi.md#getmsgvpnsequencedtopics) | **Get** /msgVpns/{msgVpnName}/sequencedTopics | Get a list of Sequenced Topic objects.
*AllApi* | [**GetMsgVpnTopicEndpoint**](docs/AllApi.md#getmsgvpntopicendpoint) | **Get** /msgVpns/{msgVpnName}/topicEndpoints/{topicEndpointName} | Get a Topic Endpoint object.
*AllApi* | [**GetMsgVpnTopicEndpointTemplate**](docs/AllApi.md#getmsgvpntopicendpointtemplate) | **Get** /msgVpns/{msgVpnName}/topicEndpointTemplates/{topicEndpointTemplateName} | Get a Topic Endpoint Template object.
*AllApi* | [**GetMsgVpnTopicEndpointTemplates**](docs/AllApi.md#getmsgvpntopicendpointtemplates) | **Get** /msgVpns/{msgVpnName}/topicEndpointTemplates | Get a list of Topic Endpoint Template objects.
*AllApi* | [**GetMsgVpnTopicEndpoints**](docs/AllApi.md#getmsgvpntopicendpoints) | **Get** /msgVpns/{msgVpnName}/topicEndpoints | Get a list of Topic Endpoint objects.
*AllApi* | [**GetMsgVpns**](docs/AllApi.md#getmsgvpns) | **Get** /msgVpns | Get a list of Message VPN objects.
*AllApi* | [**GetOauthProfile**](docs/AllApi.md#getoauthprofile) | **Get** /oauthProfiles/{oauthProfileName} | Get an OAuth Profile object.
*AllApi* | [**GetOauthProfileAccessLevelGroup**](docs/AllApi.md#getoauthprofileaccesslevelgroup) | **Get** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName} | Get a Group Access Level object.
*AllApi* | [**GetOauthProfileAccessLevelGroupMsgVpnAccessLevelException**](docs/AllApi.md#getoauthprofileaccesslevelgroupmsgvpnaccesslevelexception) | **Get** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions/{msgVpnName} | Get a Message VPN Access-Level Exception object.
*AllApi* | [**GetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptions**](docs/AllApi.md#getoauthprofileaccesslevelgroupmsgvpnaccesslevelexceptions) | **Get** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions | Get a list of Message VPN Access-Level Exception objects.
*AllApi* | [**GetOauthProfileAccessLevelGroups**](docs/AllApi.md#getoauthprofileaccesslevelgroups) | **Get** /oauthProfiles/{oauthProfileName}/accessLevelGroups | Get a list of Group Access Level objects.
*AllApi* | [**GetOauthProfileClientAllowedHost**](docs/AllApi.md#getoauthprofileclientallowedhost) | **Get** /oauthProfiles/{oauthProfileName}/clientAllowedHosts/{allowedHost} | Get an Allowed Host Value object.
*AllApi* | [**GetOauthProfileClientAllowedHosts**](docs/AllApi.md#getoauthprofileclientallowedhosts) | **Get** /oauthProfiles/{oauthProfileName}/clientAllowedHosts | Get a list of Allowed Host Value objects.
*AllApi* | [**GetOauthProfileClientAuthorizationParameter**](docs/AllApi.md#getoauthprofileclientauthorizationparameter) | **Get** /oauthProfiles/{oauthProfileName}/clientAuthorizationParameters/{authorizationParameterName} | Get an Authorization Parameter object.
*AllApi* | [**GetOauthProfileClientAuthorizationParameters**](docs/AllApi.md#getoauthprofileclientauthorizationparameters) | **Get** /oauthProfiles/{oauthProfileName}/clientAuthorizationParameters | Get a list of Authorization Parameter objects.
*AllApi* | [**GetOauthProfileClientRequiredClaim**](docs/AllApi.md#getoauthprofileclientrequiredclaim) | **Get** /oauthProfiles/{oauthProfileName}/clientRequiredClaims/{clientRequiredClaimName} | Get a Required Claim object.
*AllApi* | [**GetOauthProfileClientRequiredClaims**](docs/AllApi.md#getoauthprofileclientrequiredclaims) | **Get** /oauthProfiles/{oauthProfileName}/clientRequiredClaims | Get a list of Required Claim objects.
*AllApi* | [**GetOauthProfileDefaultMsgVpnAccessLevelException**](docs/AllApi.md#getoauthprofiledefaultmsgvpnaccesslevelexception) | **Get** /oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions/{msgVpnName} | Get a Message VPN Access-Level Exception object.
*AllApi* | [**GetOauthProfileDefaultMsgVpnAccessLevelExceptions**](docs/AllApi.md#getoauthprofiledefaultmsgvpnaccesslevelexceptions) | **Get** /oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions | Get a list of Message VPN Access-Level Exception objects.
*AllApi* | [**GetOauthProfileResourceServerRequiredClaim**](docs/AllApi.md#getoauthprofileresourceserverrequiredclaim) | **Get** /oauthProfiles/{oauthProfileName}/resourceServerRequiredClaims/{resourceServerRequiredClaimName} | Get a Required Claim object.
*AllApi* | [**GetOauthProfileResourceServerRequiredClaims**](docs/AllApi.md#getoauthprofileresourceserverrequiredclaims) | **Get** /oauthProfiles/{oauthProfileName}/resourceServerRequiredClaims | Get a list of Required Claim objects.
*AllApi* | [**GetOauthProfiles**](docs/AllApi.md#getoauthprofiles) | **Get** /oauthProfiles | Get a list of OAuth Profile objects.
*AllApi* | [**GetSystemInformation**](docs/AllApi.md#getsysteminformation) | **Get** /systemInformation | Get a System Information object.
*AllApi* | [**GetVirtualHostname**](docs/AllApi.md#getvirtualhostname) | **Get** /virtualHostnames/{virtualHostname} | Get a Virtual Hostname object.
*AllApi* | [**GetVirtualHostnames**](docs/AllApi.md#getvirtualhostnames) | **Get** /virtualHostnames | Get a list of Virtual Hostname objects.
*AllApi* | [**ReplaceCertAuthority**](docs/AllApi.md#replacecertauthority) | **Put** /certAuthorities/{certAuthorityName} | Replace a Certificate Authority object.
*AllApi* | [**ReplaceClientCertAuthority**](docs/AllApi.md#replaceclientcertauthority) | **Put** /clientCertAuthorities/{certAuthorityName} | Replace a Client Certificate Authority object.
*AllApi* | [**ReplaceDmrCluster**](docs/AllApi.md#replacedmrcluster) | **Put** /dmrClusters/{dmrClusterName} | Replace a Cluster object.
*AllApi* | [**ReplaceDmrClusterLink**](docs/AllApi.md#replacedmrclusterlink) | **Put** /dmrClusters/{dmrClusterName}/links/{remoteNodeName} | Replace a Link object.
*AllApi* | [**ReplaceDomainCertAuthority**](docs/AllApi.md#replacedomaincertauthority) | **Put** /domainCertAuthorities/{certAuthorityName} | Replace a Domain Certificate Authority object.
*AllApi* | [**ReplaceMsgVpn**](docs/AllApi.md#replacemsgvpn) | **Put** /msgVpns/{msgVpnName} | Replace a Message VPN object.
*AllApi* | [**ReplaceMsgVpnAclProfile**](docs/AllApi.md#replacemsgvpnaclprofile) | **Put** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName} | Replace an ACL Profile object.
*AllApi* | [**ReplaceMsgVpnAuthenticationOauthProfile**](docs/AllApi.md#replacemsgvpnauthenticationoauthprofile) | **Put** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName} | Replace an OAuth Profile object.
*AllApi* | [**ReplaceMsgVpnAuthenticationOauthProvider**](docs/AllApi.md#replacemsgvpnauthenticationoauthprovider) | **Put** /msgVpns/{msgVpnName}/authenticationOauthProviders/{oauthProviderName} | Replace an OAuth Provider object.
*AllApi* | [**ReplaceMsgVpnAuthorizationGroup**](docs/AllApi.md#replacemsgvpnauthorizationgroup) | **Put** /msgVpns/{msgVpnName}/authorizationGroups/{authorizationGroupName} | Replace an Authorization Group object.
*AllApi* | [**ReplaceMsgVpnBridge**](docs/AllApi.md#replacemsgvpnbridge) | **Put** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter} | Replace a Bridge object.
*AllApi* | [**ReplaceMsgVpnBridgeRemoteMsgVpn**](docs/AllApi.md#replacemsgvpnbridgeremotemsgvpn) | **Put** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface} | Replace a Remote Message VPN object.
*AllApi* | [**ReplaceMsgVpnClientProfile**](docs/AllApi.md#replacemsgvpnclientprofile) | **Put** /msgVpns/{msgVpnName}/clientProfiles/{clientProfileName} | Replace a Client Profile object.
*AllApi* | [**ReplaceMsgVpnClientUsername**](docs/AllApi.md#replacemsgvpnclientusername) | **Put** /msgVpns/{msgVpnName}/clientUsernames/{clientUsername} | Replace a Client Username object.
*AllApi* | [**ReplaceMsgVpnDistributedCache**](docs/AllApi.md#replacemsgvpndistributedcache) | **Put** /msgVpns/{msgVpnName}/distributedCaches/{cacheName} | Replace a Distributed Cache object.
*AllApi* | [**ReplaceMsgVpnDistributedCacheCluster**](docs/AllApi.md#replacemsgvpndistributedcachecluster) | **Put** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName} | Replace a Cache Cluster object.
*AllApi* | [**ReplaceMsgVpnDistributedCacheClusterInstance**](docs/AllApi.md#replacemsgvpndistributedcacheclusterinstance) | **Put** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName} | Replace a Cache Instance object.
*AllApi* | [**ReplaceMsgVpnDmrBridge**](docs/AllApi.md#replacemsgvpndmrbridge) | **Put** /msgVpns/{msgVpnName}/dmrBridges/{remoteNodeName} | Replace a DMR Bridge object.
*AllApi* | [**ReplaceMsgVpnJndiConnectionFactory**](docs/AllApi.md#replacemsgvpnjndiconnectionfactory) | **Put** /msgVpns/{msgVpnName}/jndiConnectionFactories/{connectionFactoryName} | Replace a JNDI Connection Factory object.
*AllApi* | [**ReplaceMsgVpnJndiQueue**](docs/AllApi.md#replacemsgvpnjndiqueue) | **Put** /msgVpns/{msgVpnName}/jndiQueues/{queueName} | Replace a JNDI Queue object.
*AllApi* | [**ReplaceMsgVpnJndiTopic**](docs/AllApi.md#replacemsgvpnjnditopic) | **Put** /msgVpns/{msgVpnName}/jndiTopics/{topicName} | Replace a JNDI Topic object.
*AllApi* | [**ReplaceMsgVpnMqttRetainCache**](docs/AllApi.md#replacemsgvpnmqttretaincache) | **Put** /msgVpns/{msgVpnName}/mqttRetainCaches/{cacheName} | Replace an MQTT Retain Cache object.
*AllApi* | [**ReplaceMsgVpnMqttSession**](docs/AllApi.md#replacemsgvpnmqttsession) | **Put** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter} | Replace an MQTT Session object.
*AllApi* | [**ReplaceMsgVpnMqttSessionSubscription**](docs/AllApi.md#replacemsgvpnmqttsessionsubscription) | **Put** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic} | Replace a Subscription object.
*AllApi* | [**ReplaceMsgVpnQueue**](docs/AllApi.md#replacemsgvpnqueue) | **Put** /msgVpns/{msgVpnName}/queues/{queueName} | Replace a Queue object.
*AllApi* | [**ReplaceMsgVpnQueueTemplate**](docs/AllApi.md#replacemsgvpnqueuetemplate) | **Put** /msgVpns/{msgVpnName}/queueTemplates/{queueTemplateName} | Replace a Queue Template object.
*AllApi* | [**ReplaceMsgVpnReplayLog**](docs/AllApi.md#replacemsgvpnreplaylog) | **Put** /msgVpns/{msgVpnName}/replayLogs/{replayLogName} | Replace a Replay Log object.
*AllApi* | [**ReplaceMsgVpnReplicatedTopic**](docs/AllApi.md#replacemsgvpnreplicatedtopic) | **Put** /msgVpns/{msgVpnName}/replicatedTopics/{replicatedTopic} | Replace a Replicated Topic object.
*AllApi* | [**ReplaceMsgVpnRestDeliveryPoint**](docs/AllApi.md#replacemsgvpnrestdeliverypoint) | **Put** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName} | Replace a REST Delivery Point object.
*AllApi* | [**ReplaceMsgVpnRestDeliveryPointQueueBinding**](docs/AllApi.md#replacemsgvpnrestdeliverypointqueuebinding) | **Put** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName} | Replace a Queue Binding object.
*AllApi* | [**ReplaceMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/AllApi.md#replacemsgvpnrestdeliverypointqueuebindingrequestheader) | **Put** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders/{headerName} | Replace a Request Header object.
*AllApi* | [**ReplaceMsgVpnRestDeliveryPointRestConsumer**](docs/AllApi.md#replacemsgvpnrestdeliverypointrestconsumer) | **Put** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName} | Replace a REST Consumer object.
*AllApi* | [**ReplaceMsgVpnTopicEndpoint**](docs/AllApi.md#replacemsgvpntopicendpoint) | **Put** /msgVpns/{msgVpnName}/topicEndpoints/{topicEndpointName} | Replace a Topic Endpoint object.
*AllApi* | [**ReplaceMsgVpnTopicEndpointTemplate**](docs/AllApi.md#replacemsgvpntopicendpointtemplate) | **Put** /msgVpns/{msgVpnName}/topicEndpointTemplates/{topicEndpointTemplateName} | Replace a Topic Endpoint Template object.
*AllApi* | [**ReplaceOauthProfile**](docs/AllApi.md#replaceoauthprofile) | **Put** /oauthProfiles/{oauthProfileName} | Replace an OAuth Profile object.
*AllApi* | [**ReplaceOauthProfileAccessLevelGroup**](docs/AllApi.md#replaceoauthprofileaccesslevelgroup) | **Put** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName} | Replace a Group Access Level object.
*AllApi* | [**ReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelException**](docs/AllApi.md#replaceoauthprofileaccesslevelgroupmsgvpnaccesslevelexception) | **Put** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions/{msgVpnName} | Replace a Message VPN Access-Level Exception object.
*AllApi* | [**ReplaceOauthProfileClientAuthorizationParameter**](docs/AllApi.md#replaceoauthprofileclientauthorizationparameter) | **Put** /oauthProfiles/{oauthProfileName}/clientAuthorizationParameters/{authorizationParameterName} | Replace an Authorization Parameter object.
*AllApi* | [**ReplaceOauthProfileDefaultMsgVpnAccessLevelException**](docs/AllApi.md#replaceoauthprofiledefaultmsgvpnaccesslevelexception) | **Put** /oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions/{msgVpnName} | Replace a Message VPN Access-Level Exception object.
*AllApi* | [**ReplaceVirtualHostname**](docs/AllApi.md#replacevirtualhostname) | **Put** /virtualHostnames/{virtualHostname} | Replace a Virtual Hostname object.
*AllApi* | [**UpdateBroker**](docs/AllApi.md#updatebroker) | **Patch** / | Update a Broker object.
*AllApi* | [**UpdateCertAuthority**](docs/AllApi.md#updatecertauthority) | **Patch** /certAuthorities/{certAuthorityName} | Update a Certificate Authority object.
*AllApi* | [**UpdateClientCertAuthority**](docs/AllApi.md#updateclientcertauthority) | **Patch** /clientCertAuthorities/{certAuthorityName} | Update a Client Certificate Authority object.
*AllApi* | [**UpdateDmrCluster**](docs/AllApi.md#updatedmrcluster) | **Patch** /dmrClusters/{dmrClusterName} | Update a Cluster object.
*AllApi* | [**UpdateDmrClusterLink**](docs/AllApi.md#updatedmrclusterlink) | **Patch** /dmrClusters/{dmrClusterName}/links/{remoteNodeName} | Update a Link object.
*AllApi* | [**UpdateDomainCertAuthority**](docs/AllApi.md#updatedomaincertauthority) | **Patch** /domainCertAuthorities/{certAuthorityName} | Update a Domain Certificate Authority object.
*AllApi* | [**UpdateMsgVpn**](docs/AllApi.md#updatemsgvpn) | **Patch** /msgVpns/{msgVpnName} | Update a Message VPN object.
*AllApi* | [**UpdateMsgVpnAclProfile**](docs/AllApi.md#updatemsgvpnaclprofile) | **Patch** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName} | Update an ACL Profile object.
*AllApi* | [**UpdateMsgVpnAuthenticationOauthProfile**](docs/AllApi.md#updatemsgvpnauthenticationoauthprofile) | **Patch** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName} | Update an OAuth Profile object.
*AllApi* | [**UpdateMsgVpnAuthenticationOauthProvider**](docs/AllApi.md#updatemsgvpnauthenticationoauthprovider) | **Patch** /msgVpns/{msgVpnName}/authenticationOauthProviders/{oauthProviderName} | Update an OAuth Provider object.
*AllApi* | [**UpdateMsgVpnAuthorizationGroup**](docs/AllApi.md#updatemsgvpnauthorizationgroup) | **Patch** /msgVpns/{msgVpnName}/authorizationGroups/{authorizationGroupName} | Update an Authorization Group object.
*AllApi* | [**UpdateMsgVpnBridge**](docs/AllApi.md#updatemsgvpnbridge) | **Patch** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter} | Update a Bridge object.
*AllApi* | [**UpdateMsgVpnBridgeRemoteMsgVpn**](docs/AllApi.md#updatemsgvpnbridgeremotemsgvpn) | **Patch** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface} | Update a Remote Message VPN object.
*AllApi* | [**UpdateMsgVpnClientProfile**](docs/AllApi.md#updatemsgvpnclientprofile) | **Patch** /msgVpns/{msgVpnName}/clientProfiles/{clientProfileName} | Update a Client Profile object.
*AllApi* | [**UpdateMsgVpnClientUsername**](docs/AllApi.md#updatemsgvpnclientusername) | **Patch** /msgVpns/{msgVpnName}/clientUsernames/{clientUsername} | Update a Client Username object.
*AllApi* | [**UpdateMsgVpnDistributedCache**](docs/AllApi.md#updatemsgvpndistributedcache) | **Patch** /msgVpns/{msgVpnName}/distributedCaches/{cacheName} | Update a Distributed Cache object.
*AllApi* | [**UpdateMsgVpnDistributedCacheCluster**](docs/AllApi.md#updatemsgvpndistributedcachecluster) | **Patch** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName} | Update a Cache Cluster object.
*AllApi* | [**UpdateMsgVpnDistributedCacheClusterInstance**](docs/AllApi.md#updatemsgvpndistributedcacheclusterinstance) | **Patch** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName} | Update a Cache Instance object.
*AllApi* | [**UpdateMsgVpnDmrBridge**](docs/AllApi.md#updatemsgvpndmrbridge) | **Patch** /msgVpns/{msgVpnName}/dmrBridges/{remoteNodeName} | Update a DMR Bridge object.
*AllApi* | [**UpdateMsgVpnJndiConnectionFactory**](docs/AllApi.md#updatemsgvpnjndiconnectionfactory) | **Patch** /msgVpns/{msgVpnName}/jndiConnectionFactories/{connectionFactoryName} | Update a JNDI Connection Factory object.
*AllApi* | [**UpdateMsgVpnJndiQueue**](docs/AllApi.md#updatemsgvpnjndiqueue) | **Patch** /msgVpns/{msgVpnName}/jndiQueues/{queueName} | Update a JNDI Queue object.
*AllApi* | [**UpdateMsgVpnJndiTopic**](docs/AllApi.md#updatemsgvpnjnditopic) | **Patch** /msgVpns/{msgVpnName}/jndiTopics/{topicName} | Update a JNDI Topic object.
*AllApi* | [**UpdateMsgVpnMqttRetainCache**](docs/AllApi.md#updatemsgvpnmqttretaincache) | **Patch** /msgVpns/{msgVpnName}/mqttRetainCaches/{cacheName} | Update an MQTT Retain Cache object.
*AllApi* | [**UpdateMsgVpnMqttSession**](docs/AllApi.md#updatemsgvpnmqttsession) | **Patch** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter} | Update an MQTT Session object.
*AllApi* | [**UpdateMsgVpnMqttSessionSubscription**](docs/AllApi.md#updatemsgvpnmqttsessionsubscription) | **Patch** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic} | Update a Subscription object.
*AllApi* | [**UpdateMsgVpnQueue**](docs/AllApi.md#updatemsgvpnqueue) | **Patch** /msgVpns/{msgVpnName}/queues/{queueName} | Update a Queue object.
*AllApi* | [**UpdateMsgVpnQueueTemplate**](docs/AllApi.md#updatemsgvpnqueuetemplate) | **Patch** /msgVpns/{msgVpnName}/queueTemplates/{queueTemplateName} | Update a Queue Template object.
*AllApi* | [**UpdateMsgVpnReplayLog**](docs/AllApi.md#updatemsgvpnreplaylog) | **Patch** /msgVpns/{msgVpnName}/replayLogs/{replayLogName} | Update a Replay Log object.
*AllApi* | [**UpdateMsgVpnReplicatedTopic**](docs/AllApi.md#updatemsgvpnreplicatedtopic) | **Patch** /msgVpns/{msgVpnName}/replicatedTopics/{replicatedTopic} | Update a Replicated Topic object.
*AllApi* | [**UpdateMsgVpnRestDeliveryPoint**](docs/AllApi.md#updatemsgvpnrestdeliverypoint) | **Patch** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName} | Update a REST Delivery Point object.
*AllApi* | [**UpdateMsgVpnRestDeliveryPointQueueBinding**](docs/AllApi.md#updatemsgvpnrestdeliverypointqueuebinding) | **Patch** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName} | Update a Queue Binding object.
*AllApi* | [**UpdateMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/AllApi.md#updatemsgvpnrestdeliverypointqueuebindingrequestheader) | **Patch** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders/{headerName} | Update a Request Header object.
*AllApi* | [**UpdateMsgVpnRestDeliveryPointRestConsumer**](docs/AllApi.md#updatemsgvpnrestdeliverypointrestconsumer) | **Patch** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName} | Update a REST Consumer object.
*AllApi* | [**UpdateMsgVpnTopicEndpoint**](docs/AllApi.md#updatemsgvpntopicendpoint) | **Patch** /msgVpns/{msgVpnName}/topicEndpoints/{topicEndpointName} | Update a Topic Endpoint object.
*AllApi* | [**UpdateMsgVpnTopicEndpointTemplate**](docs/AllApi.md#updatemsgvpntopicendpointtemplate) | **Patch** /msgVpns/{msgVpnName}/topicEndpointTemplates/{topicEndpointTemplateName} | Update a Topic Endpoint Template object.
*AllApi* | [**UpdateOauthProfile**](docs/AllApi.md#updateoauthprofile) | **Patch** /oauthProfiles/{oauthProfileName} | Update an OAuth Profile object.
*AllApi* | [**UpdateOauthProfileAccessLevelGroup**](docs/AllApi.md#updateoauthprofileaccesslevelgroup) | **Patch** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName} | Update a Group Access Level object.
*AllApi* | [**UpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelException**](docs/AllApi.md#updateoauthprofileaccesslevelgroupmsgvpnaccesslevelexception) | **Patch** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions/{msgVpnName} | Update a Message VPN Access-Level Exception object.
*AllApi* | [**UpdateOauthProfileClientAuthorizationParameter**](docs/AllApi.md#updateoauthprofileclientauthorizationparameter) | **Patch** /oauthProfiles/{oauthProfileName}/clientAuthorizationParameters/{authorizationParameterName} | Update an Authorization Parameter object.
*AllApi* | [**UpdateOauthProfileDefaultMsgVpnAccessLevelException**](docs/AllApi.md#updateoauthprofiledefaultmsgvpnaccesslevelexception) | **Patch** /oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions/{msgVpnName} | Update a Message VPN Access-Level Exception object.
*AllApi* | [**UpdateVirtualHostname**](docs/AllApi.md#updatevirtualhostname) | **Patch** /virtualHostnames/{virtualHostname} | Update a Virtual Hostname object.
*AuthenticationOauthProfileApi* | [**CreateMsgVpnAuthenticationOauthProfile**](docs/AuthenticationOauthProfileApi.md#createmsgvpnauthenticationoauthprofile) | **Post** /msgVpns/{msgVpnName}/authenticationOauthProfiles | Create an OAuth Profile object.
*AuthenticationOauthProfileApi* | [**CreateMsgVpnAuthenticationOauthProfileClientRequiredClaim**](docs/AuthenticationOauthProfileApi.md#createmsgvpnauthenticationoauthprofileclientrequiredclaim) | **Post** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/clientRequiredClaims | Create a Required Claim object.
*AuthenticationOauthProfileApi* | [**CreateMsgVpnAuthenticationOauthProfileResourceServerRequiredClaim**](docs/AuthenticationOauthProfileApi.md#createmsgvpnauthenticationoauthprofileresourceserverrequiredclaim) | **Post** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/resourceServerRequiredClaims | Create a Required Claim object.
*AuthenticationOauthProfileApi* | [**DeleteMsgVpnAuthenticationOauthProfile**](docs/AuthenticationOauthProfileApi.md#deletemsgvpnauthenticationoauthprofile) | **Delete** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName} | Delete an OAuth Profile object.
*AuthenticationOauthProfileApi* | [**DeleteMsgVpnAuthenticationOauthProfileClientRequiredClaim**](docs/AuthenticationOauthProfileApi.md#deletemsgvpnauthenticationoauthprofileclientrequiredclaim) | **Delete** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/clientRequiredClaims/{clientRequiredClaimName} | Delete a Required Claim object.
*AuthenticationOauthProfileApi* | [**DeleteMsgVpnAuthenticationOauthProfileResourceServerRequiredClaim**](docs/AuthenticationOauthProfileApi.md#deletemsgvpnauthenticationoauthprofileresourceserverrequiredclaim) | **Delete** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/resourceServerRequiredClaims/{resourceServerRequiredClaimName} | Delete a Required Claim object.
*AuthenticationOauthProfileApi* | [**GetMsgVpnAuthenticationOauthProfile**](docs/AuthenticationOauthProfileApi.md#getmsgvpnauthenticationoauthprofile) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName} | Get an OAuth Profile object.
*AuthenticationOauthProfileApi* | [**GetMsgVpnAuthenticationOauthProfileClientRequiredClaim**](docs/AuthenticationOauthProfileApi.md#getmsgvpnauthenticationoauthprofileclientrequiredclaim) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/clientRequiredClaims/{clientRequiredClaimName} | Get a Required Claim object.
*AuthenticationOauthProfileApi* | [**GetMsgVpnAuthenticationOauthProfileClientRequiredClaims**](docs/AuthenticationOauthProfileApi.md#getmsgvpnauthenticationoauthprofileclientrequiredclaims) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/clientRequiredClaims | Get a list of Required Claim objects.
*AuthenticationOauthProfileApi* | [**GetMsgVpnAuthenticationOauthProfileResourceServerRequiredClaim**](docs/AuthenticationOauthProfileApi.md#getmsgvpnauthenticationoauthprofileresourceserverrequiredclaim) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/resourceServerRequiredClaims/{resourceServerRequiredClaimName} | Get a Required Claim object.
*AuthenticationOauthProfileApi* | [**GetMsgVpnAuthenticationOauthProfileResourceServerRequiredClaims**](docs/AuthenticationOauthProfileApi.md#getmsgvpnauthenticationoauthprofileresourceserverrequiredclaims) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/resourceServerRequiredClaims | Get a list of Required Claim objects.
*AuthenticationOauthProfileApi* | [**GetMsgVpnAuthenticationOauthProfiles**](docs/AuthenticationOauthProfileApi.md#getmsgvpnauthenticationoauthprofiles) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles | Get a list of OAuth Profile objects.
*AuthenticationOauthProfileApi* | [**ReplaceMsgVpnAuthenticationOauthProfile**](docs/AuthenticationOauthProfileApi.md#replacemsgvpnauthenticationoauthprofile) | **Put** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName} | Replace an OAuth Profile object.
*AuthenticationOauthProfileApi* | [**UpdateMsgVpnAuthenticationOauthProfile**](docs/AuthenticationOauthProfileApi.md#updatemsgvpnauthenticationoauthprofile) | **Patch** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName} | Update an OAuth Profile object.
*AuthenticationOauthProviderApi* | [**CreateMsgVpnAuthenticationOauthProvider**](docs/AuthenticationOauthProviderApi.md#createmsgvpnauthenticationoauthprovider) | **Post** /msgVpns/{msgVpnName}/authenticationOauthProviders | Create an OAuth Provider object.
*AuthenticationOauthProviderApi* | [**DeleteMsgVpnAuthenticationOauthProvider**](docs/AuthenticationOauthProviderApi.md#deletemsgvpnauthenticationoauthprovider) | **Delete** /msgVpns/{msgVpnName}/authenticationOauthProviders/{oauthProviderName} | Delete an OAuth Provider object.
*AuthenticationOauthProviderApi* | [**GetMsgVpnAuthenticationOauthProvider**](docs/AuthenticationOauthProviderApi.md#getmsgvpnauthenticationoauthprovider) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProviders/{oauthProviderName} | Get an OAuth Provider object.
*AuthenticationOauthProviderApi* | [**GetMsgVpnAuthenticationOauthProviders**](docs/AuthenticationOauthProviderApi.md#getmsgvpnauthenticationoauthproviders) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProviders | Get a list of OAuth Provider objects.
*AuthenticationOauthProviderApi* | [**ReplaceMsgVpnAuthenticationOauthProvider**](docs/AuthenticationOauthProviderApi.md#replacemsgvpnauthenticationoauthprovider) | **Put** /msgVpns/{msgVpnName}/authenticationOauthProviders/{oauthProviderName} | Replace an OAuth Provider object.
*AuthenticationOauthProviderApi* | [**UpdateMsgVpnAuthenticationOauthProvider**](docs/AuthenticationOauthProviderApi.md#updatemsgvpnauthenticationoauthprovider) | **Patch** /msgVpns/{msgVpnName}/authenticationOauthProviders/{oauthProviderName} | Update an OAuth Provider object.
*AuthorizationGroupApi* | [**CreateMsgVpnAuthorizationGroup**](docs/AuthorizationGroupApi.md#createmsgvpnauthorizationgroup) | **Post** /msgVpns/{msgVpnName}/authorizationGroups | Create an Authorization Group object.
*AuthorizationGroupApi* | [**DeleteMsgVpnAuthorizationGroup**](docs/AuthorizationGroupApi.md#deletemsgvpnauthorizationgroup) | **Delete** /msgVpns/{msgVpnName}/authorizationGroups/{authorizationGroupName} | Delete an Authorization Group object.
*AuthorizationGroupApi* | [**GetMsgVpnAuthorizationGroup**](docs/AuthorizationGroupApi.md#getmsgvpnauthorizationgroup) | **Get** /msgVpns/{msgVpnName}/authorizationGroups/{authorizationGroupName} | Get an Authorization Group object.
*AuthorizationGroupApi* | [**GetMsgVpnAuthorizationGroups**](docs/AuthorizationGroupApi.md#getmsgvpnauthorizationgroups) | **Get** /msgVpns/{msgVpnName}/authorizationGroups | Get a list of Authorization Group objects.
*AuthorizationGroupApi* | [**ReplaceMsgVpnAuthorizationGroup**](docs/AuthorizationGroupApi.md#replacemsgvpnauthorizationgroup) | **Put** /msgVpns/{msgVpnName}/authorizationGroups/{authorizationGroupName} | Replace an Authorization Group object.
*AuthorizationGroupApi* | [**UpdateMsgVpnAuthorizationGroup**](docs/AuthorizationGroupApi.md#updatemsgvpnauthorizationgroup) | **Patch** /msgVpns/{msgVpnName}/authorizationGroups/{authorizationGroupName} | Update an Authorization Group object.
*BridgeApi* | [**CreateMsgVpnBridge**](docs/BridgeApi.md#createmsgvpnbridge) | **Post** /msgVpns/{msgVpnName}/bridges | Create a Bridge object.
*BridgeApi* | [**CreateMsgVpnBridgeRemoteMsgVpn**](docs/BridgeApi.md#createmsgvpnbridgeremotemsgvpn) | **Post** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns | Create a Remote Message VPN object.
*BridgeApi* | [**CreateMsgVpnBridgeRemoteSubscription**](docs/BridgeApi.md#createmsgvpnbridgeremotesubscription) | **Post** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions | Create a Remote Subscription object.
*BridgeApi* | [**CreateMsgVpnBridgeTlsTrustedCommonName**](docs/BridgeApi.md#createmsgvpnbridgetlstrustedcommonname) | **Post** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames | Create a Trusted Common Name object.
*BridgeApi* | [**DeleteMsgVpnBridge**](docs/BridgeApi.md#deletemsgvpnbridge) | **Delete** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter} | Delete a Bridge object.
*BridgeApi* | [**DeleteMsgVpnBridgeRemoteMsgVpn**](docs/BridgeApi.md#deletemsgvpnbridgeremotemsgvpn) | **Delete** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface} | Delete a Remote Message VPN object.
*BridgeApi* | [**DeleteMsgVpnBridgeRemoteSubscription**](docs/BridgeApi.md#deletemsgvpnbridgeremotesubscription) | **Delete** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions/{remoteSubscriptionTopic} | Delete a Remote Subscription object.
*BridgeApi* | [**DeleteMsgVpnBridgeTlsTrustedCommonName**](docs/BridgeApi.md#deletemsgvpnbridgetlstrustedcommonname) | **Delete** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Delete a Trusted Common Name object.
*BridgeApi* | [**GetMsgVpnBridge**](docs/BridgeApi.md#getmsgvpnbridge) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter} | Get a Bridge object.
*BridgeApi* | [**GetMsgVpnBridgeRemoteMsgVpn**](docs/BridgeApi.md#getmsgvpnbridgeremotemsgvpn) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface} | Get a Remote Message VPN object.
*BridgeApi* | [**GetMsgVpnBridgeRemoteMsgVpns**](docs/BridgeApi.md#getmsgvpnbridgeremotemsgvpns) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns | Get a list of Remote Message VPN objects.
*BridgeApi* | [**GetMsgVpnBridgeRemoteSubscription**](docs/BridgeApi.md#getmsgvpnbridgeremotesubscription) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions/{remoteSubscriptionTopic} | Get a Remote Subscription object.
*BridgeApi* | [**GetMsgVpnBridgeRemoteSubscriptions**](docs/BridgeApi.md#getmsgvpnbridgeremotesubscriptions) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions | Get a list of Remote Subscription objects.
*BridgeApi* | [**GetMsgVpnBridgeTlsTrustedCommonName**](docs/BridgeApi.md#getmsgvpnbridgetlstrustedcommonname) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Get a Trusted Common Name object.
*BridgeApi* | [**GetMsgVpnBridgeTlsTrustedCommonNames**](docs/BridgeApi.md#getmsgvpnbridgetlstrustedcommonnames) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames | Get a list of Trusted Common Name objects.
*BridgeApi* | [**GetMsgVpnBridges**](docs/BridgeApi.md#getmsgvpnbridges) | **Get** /msgVpns/{msgVpnName}/bridges | Get a list of Bridge objects.
*BridgeApi* | [**ReplaceMsgVpnBridge**](docs/BridgeApi.md#replacemsgvpnbridge) | **Put** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter} | Replace a Bridge object.
*BridgeApi* | [**ReplaceMsgVpnBridgeRemoteMsgVpn**](docs/BridgeApi.md#replacemsgvpnbridgeremotemsgvpn) | **Put** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface} | Replace a Remote Message VPN object.
*BridgeApi* | [**UpdateMsgVpnBridge**](docs/BridgeApi.md#updatemsgvpnbridge) | **Patch** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter} | Update a Bridge object.
*BridgeApi* | [**UpdateMsgVpnBridgeRemoteMsgVpn**](docs/BridgeApi.md#updatemsgvpnbridgeremotemsgvpn) | **Patch** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface} | Update a Remote Message VPN object.
*CertAuthorityApi* | [**CreateCertAuthority**](docs/CertAuthorityApi.md#createcertauthority) | **Post** /certAuthorities | Create a Certificate Authority object.
*CertAuthorityApi* | [**CreateCertAuthorityOcspTlsTrustedCommonName**](docs/CertAuthorityApi.md#createcertauthorityocsptlstrustedcommonname) | **Post** /certAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames | Create an OCSP Responder Trusted Common Name object.
*CertAuthorityApi* | [**DeleteCertAuthority**](docs/CertAuthorityApi.md#deletecertauthority) | **Delete** /certAuthorities/{certAuthorityName} | Delete a Certificate Authority object.
*CertAuthorityApi* | [**DeleteCertAuthorityOcspTlsTrustedCommonName**](docs/CertAuthorityApi.md#deletecertauthorityocsptlstrustedcommonname) | **Delete** /certAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames/{ocspTlsTrustedCommonName} | Delete an OCSP Responder Trusted Common Name object.
*CertAuthorityApi* | [**GetCertAuthorities**](docs/CertAuthorityApi.md#getcertauthorities) | **Get** /certAuthorities | Get a list of Certificate Authority objects.
*CertAuthorityApi* | [**GetCertAuthority**](docs/CertAuthorityApi.md#getcertauthority) | **Get** /certAuthorities/{certAuthorityName} | Get a Certificate Authority object.
*CertAuthorityApi* | [**GetCertAuthorityOcspTlsTrustedCommonName**](docs/CertAuthorityApi.md#getcertauthorityocsptlstrustedcommonname) | **Get** /certAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames/{ocspTlsTrustedCommonName} | Get an OCSP Responder Trusted Common Name object.
*CertAuthorityApi* | [**GetCertAuthorityOcspTlsTrustedCommonNames**](docs/CertAuthorityApi.md#getcertauthorityocsptlstrustedcommonnames) | **Get** /certAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames | Get a list of OCSP Responder Trusted Common Name objects.
*CertAuthorityApi* | [**ReplaceCertAuthority**](docs/CertAuthorityApi.md#replacecertauthority) | **Put** /certAuthorities/{certAuthorityName} | Replace a Certificate Authority object.
*CertAuthorityApi* | [**UpdateCertAuthority**](docs/CertAuthorityApi.md#updatecertauthority) | **Patch** /certAuthorities/{certAuthorityName} | Update a Certificate Authority object.
*ClientCertAuthorityApi* | [**CreateClientCertAuthority**](docs/ClientCertAuthorityApi.md#createclientcertauthority) | **Post** /clientCertAuthorities | Create a Client Certificate Authority object.
*ClientCertAuthorityApi* | [**CreateClientCertAuthorityOcspTlsTrustedCommonName**](docs/ClientCertAuthorityApi.md#createclientcertauthorityocsptlstrustedcommonname) | **Post** /clientCertAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames | Create an OCSP Responder Trusted Common Name object.
*ClientCertAuthorityApi* | [**DeleteClientCertAuthority**](docs/ClientCertAuthorityApi.md#deleteclientcertauthority) | **Delete** /clientCertAuthorities/{certAuthorityName} | Delete a Client Certificate Authority object.
*ClientCertAuthorityApi* | [**DeleteClientCertAuthorityOcspTlsTrustedCommonName**](docs/ClientCertAuthorityApi.md#deleteclientcertauthorityocsptlstrustedcommonname) | **Delete** /clientCertAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames/{ocspTlsTrustedCommonName} | Delete an OCSP Responder Trusted Common Name object.
*ClientCertAuthorityApi* | [**GetClientCertAuthorities**](docs/ClientCertAuthorityApi.md#getclientcertauthorities) | **Get** /clientCertAuthorities | Get a list of Client Certificate Authority objects.
*ClientCertAuthorityApi* | [**GetClientCertAuthority**](docs/ClientCertAuthorityApi.md#getclientcertauthority) | **Get** /clientCertAuthorities/{certAuthorityName} | Get a Client Certificate Authority object.
*ClientCertAuthorityApi* | [**GetClientCertAuthorityOcspTlsTrustedCommonName**](docs/ClientCertAuthorityApi.md#getclientcertauthorityocsptlstrustedcommonname) | **Get** /clientCertAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames/{ocspTlsTrustedCommonName} | Get an OCSP Responder Trusted Common Name object.
*ClientCertAuthorityApi* | [**GetClientCertAuthorityOcspTlsTrustedCommonNames**](docs/ClientCertAuthorityApi.md#getclientcertauthorityocsptlstrustedcommonnames) | **Get** /clientCertAuthorities/{certAuthorityName}/ocspTlsTrustedCommonNames | Get a list of OCSP Responder Trusted Common Name objects.
*ClientCertAuthorityApi* | [**ReplaceClientCertAuthority**](docs/ClientCertAuthorityApi.md#replaceclientcertauthority) | **Put** /clientCertAuthorities/{certAuthorityName} | Replace a Client Certificate Authority object.
*ClientCertAuthorityApi* | [**UpdateClientCertAuthority**](docs/ClientCertAuthorityApi.md#updateclientcertauthority) | **Patch** /clientCertAuthorities/{certAuthorityName} | Update a Client Certificate Authority object.
*ClientProfileApi* | [**CreateMsgVpnClientProfile**](docs/ClientProfileApi.md#createmsgvpnclientprofile) | **Post** /msgVpns/{msgVpnName}/clientProfiles | Create a Client Profile object.
*ClientProfileApi* | [**DeleteMsgVpnClientProfile**](docs/ClientProfileApi.md#deletemsgvpnclientprofile) | **Delete** /msgVpns/{msgVpnName}/clientProfiles/{clientProfileName} | Delete a Client Profile object.
*ClientProfileApi* | [**GetMsgVpnClientProfile**](docs/ClientProfileApi.md#getmsgvpnclientprofile) | **Get** /msgVpns/{msgVpnName}/clientProfiles/{clientProfileName} | Get a Client Profile object.
*ClientProfileApi* | [**GetMsgVpnClientProfiles**](docs/ClientProfileApi.md#getmsgvpnclientprofiles) | **Get** /msgVpns/{msgVpnName}/clientProfiles | Get a list of Client Profile objects.
*ClientProfileApi* | [**ReplaceMsgVpnClientProfile**](docs/ClientProfileApi.md#replacemsgvpnclientprofile) | **Put** /msgVpns/{msgVpnName}/clientProfiles/{clientProfileName} | Replace a Client Profile object.
*ClientProfileApi* | [**UpdateMsgVpnClientProfile**](docs/ClientProfileApi.md#updatemsgvpnclientprofile) | **Patch** /msgVpns/{msgVpnName}/clientProfiles/{clientProfileName} | Update a Client Profile object.
*ClientUsernameApi* | [**CreateMsgVpnClientUsername**](docs/ClientUsernameApi.md#createmsgvpnclientusername) | **Post** /msgVpns/{msgVpnName}/clientUsernames | Create a Client Username object.
*ClientUsernameApi* | [**DeleteMsgVpnClientUsername**](docs/ClientUsernameApi.md#deletemsgvpnclientusername) | **Delete** /msgVpns/{msgVpnName}/clientUsernames/{clientUsername} | Delete a Client Username object.
*ClientUsernameApi* | [**GetMsgVpnClientUsername**](docs/ClientUsernameApi.md#getmsgvpnclientusername) | **Get** /msgVpns/{msgVpnName}/clientUsernames/{clientUsername} | Get a Client Username object.
*ClientUsernameApi* | [**GetMsgVpnClientUsernames**](docs/ClientUsernameApi.md#getmsgvpnclientusernames) | **Get** /msgVpns/{msgVpnName}/clientUsernames | Get a list of Client Username objects.
*ClientUsernameApi* | [**ReplaceMsgVpnClientUsername**](docs/ClientUsernameApi.md#replacemsgvpnclientusername) | **Put** /msgVpns/{msgVpnName}/clientUsernames/{clientUsername} | Replace a Client Username object.
*ClientUsernameApi* | [**UpdateMsgVpnClientUsername**](docs/ClientUsernameApi.md#updatemsgvpnclientusername) | **Patch** /msgVpns/{msgVpnName}/clientUsernames/{clientUsername} | Update a Client Username object.
*DistributedCacheApi* | [**CreateMsgVpnDistributedCache**](docs/DistributedCacheApi.md#createmsgvpndistributedcache) | **Post** /msgVpns/{msgVpnName}/distributedCaches | Create a Distributed Cache object.
*DistributedCacheApi* | [**CreateMsgVpnDistributedCacheCluster**](docs/DistributedCacheApi.md#createmsgvpndistributedcachecluster) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters | Create a Cache Cluster object.
*DistributedCacheApi* | [**CreateMsgVpnDistributedCacheClusterGlobalCachingHomeCluster**](docs/DistributedCacheApi.md#createmsgvpndistributedcacheclusterglobalcachinghomecluster) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters | Create a Home Cache Cluster object.
*DistributedCacheApi* | [**CreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix**](docs/DistributedCacheApi.md#createmsgvpndistributedcacheclusterglobalcachinghomeclustertopicprefix) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes | Create a Topic Prefix object.
*DistributedCacheApi* | [**CreateMsgVpnDistributedCacheClusterInstance**](docs/DistributedCacheApi.md#createmsgvpndistributedcacheclusterinstance) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances | Create a Cache Instance object.
*DistributedCacheApi* | [**CreateMsgVpnDistributedCacheClusterTopic**](docs/DistributedCacheApi.md#createmsgvpndistributedcacheclustertopic) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics | Create a Topic object.
*DistributedCacheApi* | [**DeleteMsgVpnDistributedCache**](docs/DistributedCacheApi.md#deletemsgvpndistributedcache) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName} | Delete a Distributed Cache object.
*DistributedCacheApi* | [**DeleteMsgVpnDistributedCacheCluster**](docs/DistributedCacheApi.md#deletemsgvpndistributedcachecluster) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName} | Delete a Cache Cluster object.
*DistributedCacheApi* | [**DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeCluster**](docs/DistributedCacheApi.md#deletemsgvpndistributedcacheclusterglobalcachinghomecluster) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName} | Delete a Home Cache Cluster object.
*DistributedCacheApi* | [**DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix**](docs/DistributedCacheApi.md#deletemsgvpndistributedcacheclusterglobalcachinghomeclustertopicprefix) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes/{topicPrefix} | Delete a Topic Prefix object.
*DistributedCacheApi* | [**DeleteMsgVpnDistributedCacheClusterInstance**](docs/DistributedCacheApi.md#deletemsgvpndistributedcacheclusterinstance) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName} | Delete a Cache Instance object.
*DistributedCacheApi* | [**DeleteMsgVpnDistributedCacheClusterTopic**](docs/DistributedCacheApi.md#deletemsgvpndistributedcacheclustertopic) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics/{topic} | Delete a Topic object.
*DistributedCacheApi* | [**GetMsgVpnDistributedCache**](docs/DistributedCacheApi.md#getmsgvpndistributedcache) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName} | Get a Distributed Cache object.
*DistributedCacheApi* | [**GetMsgVpnDistributedCacheCluster**](docs/DistributedCacheApi.md#getmsgvpndistributedcachecluster) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName} | Get a Cache Cluster object.
*DistributedCacheApi* | [**GetMsgVpnDistributedCacheClusterGlobalCachingHomeCluster**](docs/DistributedCacheApi.md#getmsgvpndistributedcacheclusterglobalcachinghomecluster) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName} | Get a Home Cache Cluster object.
*DistributedCacheApi* | [**GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix**](docs/DistributedCacheApi.md#getmsgvpndistributedcacheclusterglobalcachinghomeclustertopicprefix) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes/{topicPrefix} | Get a Topic Prefix object.
*DistributedCacheApi* | [**GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixes**](docs/DistributedCacheApi.md#getmsgvpndistributedcacheclusterglobalcachinghomeclustertopicprefixes) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes | Get a list of Topic Prefix objects.
*DistributedCacheApi* | [**GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusters**](docs/DistributedCacheApi.md#getmsgvpndistributedcacheclusterglobalcachinghomeclusters) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters | Get a list of Home Cache Cluster objects.
*DistributedCacheApi* | [**GetMsgVpnDistributedCacheClusterInstance**](docs/DistributedCacheApi.md#getmsgvpndistributedcacheclusterinstance) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName} | Get a Cache Instance object.
*DistributedCacheApi* | [**GetMsgVpnDistributedCacheClusterInstances**](docs/DistributedCacheApi.md#getmsgvpndistributedcacheclusterinstances) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances | Get a list of Cache Instance objects.
*DistributedCacheApi* | [**GetMsgVpnDistributedCacheClusterTopic**](docs/DistributedCacheApi.md#getmsgvpndistributedcacheclustertopic) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics/{topic} | Get a Topic object.
*DistributedCacheApi* | [**GetMsgVpnDistributedCacheClusterTopics**](docs/DistributedCacheApi.md#getmsgvpndistributedcacheclustertopics) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics | Get a list of Topic objects.
*DistributedCacheApi* | [**GetMsgVpnDistributedCacheClusters**](docs/DistributedCacheApi.md#getmsgvpndistributedcacheclusters) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters | Get a list of Cache Cluster objects.
*DistributedCacheApi* | [**GetMsgVpnDistributedCaches**](docs/DistributedCacheApi.md#getmsgvpndistributedcaches) | **Get** /msgVpns/{msgVpnName}/distributedCaches | Get a list of Distributed Cache objects.
*DistributedCacheApi* | [**ReplaceMsgVpnDistributedCache**](docs/DistributedCacheApi.md#replacemsgvpndistributedcache) | **Put** /msgVpns/{msgVpnName}/distributedCaches/{cacheName} | Replace a Distributed Cache object.
*DistributedCacheApi* | [**ReplaceMsgVpnDistributedCacheCluster**](docs/DistributedCacheApi.md#replacemsgvpndistributedcachecluster) | **Put** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName} | Replace a Cache Cluster object.
*DistributedCacheApi* | [**ReplaceMsgVpnDistributedCacheClusterInstance**](docs/DistributedCacheApi.md#replacemsgvpndistributedcacheclusterinstance) | **Put** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName} | Replace a Cache Instance object.
*DistributedCacheApi* | [**UpdateMsgVpnDistributedCache**](docs/DistributedCacheApi.md#updatemsgvpndistributedcache) | **Patch** /msgVpns/{msgVpnName}/distributedCaches/{cacheName} | Update a Distributed Cache object.
*DistributedCacheApi* | [**UpdateMsgVpnDistributedCacheCluster**](docs/DistributedCacheApi.md#updatemsgvpndistributedcachecluster) | **Patch** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName} | Update a Cache Cluster object.
*DistributedCacheApi* | [**UpdateMsgVpnDistributedCacheClusterInstance**](docs/DistributedCacheApi.md#updatemsgvpndistributedcacheclusterinstance) | **Patch** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName} | Update a Cache Instance object.
*DmrBridgeApi* | [**CreateMsgVpnDmrBridge**](docs/DmrBridgeApi.md#createmsgvpndmrbridge) | **Post** /msgVpns/{msgVpnName}/dmrBridges | Create a DMR Bridge object.
*DmrBridgeApi* | [**DeleteMsgVpnDmrBridge**](docs/DmrBridgeApi.md#deletemsgvpndmrbridge) | **Delete** /msgVpns/{msgVpnName}/dmrBridges/{remoteNodeName} | Delete a DMR Bridge object.
*DmrBridgeApi* | [**GetMsgVpnDmrBridge**](docs/DmrBridgeApi.md#getmsgvpndmrbridge) | **Get** /msgVpns/{msgVpnName}/dmrBridges/{remoteNodeName} | Get a DMR Bridge object.
*DmrBridgeApi* | [**GetMsgVpnDmrBridges**](docs/DmrBridgeApi.md#getmsgvpndmrbridges) | **Get** /msgVpns/{msgVpnName}/dmrBridges | Get a list of DMR Bridge objects.
*DmrBridgeApi* | [**ReplaceMsgVpnDmrBridge**](docs/DmrBridgeApi.md#replacemsgvpndmrbridge) | **Put** /msgVpns/{msgVpnName}/dmrBridges/{remoteNodeName} | Replace a DMR Bridge object.
*DmrBridgeApi* | [**UpdateMsgVpnDmrBridge**](docs/DmrBridgeApi.md#updatemsgvpndmrbridge) | **Patch** /msgVpns/{msgVpnName}/dmrBridges/{remoteNodeName} | Update a DMR Bridge object.
*DmrClusterApi* | [**CreateDmrCluster**](docs/DmrClusterApi.md#createdmrcluster) | **Post** /dmrClusters | Create a Cluster object.
*DmrClusterApi* | [**CreateDmrClusterLink**](docs/DmrClusterApi.md#createdmrclusterlink) | **Post** /dmrClusters/{dmrClusterName}/links | Create a Link object.
*DmrClusterApi* | [**CreateDmrClusterLinkRemoteAddress**](docs/DmrClusterApi.md#createdmrclusterlinkremoteaddress) | **Post** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/remoteAddresses | Create a Remote Address object.
*DmrClusterApi* | [**CreateDmrClusterLinkTlsTrustedCommonName**](docs/DmrClusterApi.md#createdmrclusterlinktlstrustedcommonname) | **Post** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/tlsTrustedCommonNames | Create a Trusted Common Name object.
*DmrClusterApi* | [**DeleteDmrCluster**](docs/DmrClusterApi.md#deletedmrcluster) | **Delete** /dmrClusters/{dmrClusterName} | Delete a Cluster object.
*DmrClusterApi* | [**DeleteDmrClusterLink**](docs/DmrClusterApi.md#deletedmrclusterlink) | **Delete** /dmrClusters/{dmrClusterName}/links/{remoteNodeName} | Delete a Link object.
*DmrClusterApi* | [**DeleteDmrClusterLinkRemoteAddress**](docs/DmrClusterApi.md#deletedmrclusterlinkremoteaddress) | **Delete** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/remoteAddresses/{remoteAddress} | Delete a Remote Address object.
*DmrClusterApi* | [**DeleteDmrClusterLinkTlsTrustedCommonName**](docs/DmrClusterApi.md#deletedmrclusterlinktlstrustedcommonname) | **Delete** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Delete a Trusted Common Name object.
*DmrClusterApi* | [**GetDmrCluster**](docs/DmrClusterApi.md#getdmrcluster) | **Get** /dmrClusters/{dmrClusterName} | Get a Cluster object.
*DmrClusterApi* | [**GetDmrClusterLink**](docs/DmrClusterApi.md#getdmrclusterlink) | **Get** /dmrClusters/{dmrClusterName}/links/{remoteNodeName} | Get a Link object.
*DmrClusterApi* | [**GetDmrClusterLinkRemoteAddress**](docs/DmrClusterApi.md#getdmrclusterlinkremoteaddress) | **Get** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/remoteAddresses/{remoteAddress} | Get a Remote Address object.
*DmrClusterApi* | [**GetDmrClusterLinkRemoteAddresses**](docs/DmrClusterApi.md#getdmrclusterlinkremoteaddresses) | **Get** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/remoteAddresses | Get a list of Remote Address objects.
*DmrClusterApi* | [**GetDmrClusterLinkTlsTrustedCommonName**](docs/DmrClusterApi.md#getdmrclusterlinktlstrustedcommonname) | **Get** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Get a Trusted Common Name object.
*DmrClusterApi* | [**GetDmrClusterLinkTlsTrustedCommonNames**](docs/DmrClusterApi.md#getdmrclusterlinktlstrustedcommonnames) | **Get** /dmrClusters/{dmrClusterName}/links/{remoteNodeName}/tlsTrustedCommonNames | Get a list of Trusted Common Name objects.
*DmrClusterApi* | [**GetDmrClusterLinks**](docs/DmrClusterApi.md#getdmrclusterlinks) | **Get** /dmrClusters/{dmrClusterName}/links | Get a list of Link objects.
*DmrClusterApi* | [**GetDmrClusters**](docs/DmrClusterApi.md#getdmrclusters) | **Get** /dmrClusters | Get a list of Cluster objects.
*DmrClusterApi* | [**ReplaceDmrCluster**](docs/DmrClusterApi.md#replacedmrcluster) | **Put** /dmrClusters/{dmrClusterName} | Replace a Cluster object.
*DmrClusterApi* | [**ReplaceDmrClusterLink**](docs/DmrClusterApi.md#replacedmrclusterlink) | **Put** /dmrClusters/{dmrClusterName}/links/{remoteNodeName} | Replace a Link object.
*DmrClusterApi* | [**UpdateDmrCluster**](docs/DmrClusterApi.md#updatedmrcluster) | **Patch** /dmrClusters/{dmrClusterName} | Update a Cluster object.
*DmrClusterApi* | [**UpdateDmrClusterLink**](docs/DmrClusterApi.md#updatedmrclusterlink) | **Patch** /dmrClusters/{dmrClusterName}/links/{remoteNodeName} | Update a Link object.
*DomainCertAuthorityApi* | [**CreateDomainCertAuthority**](docs/DomainCertAuthorityApi.md#createdomaincertauthority) | **Post** /domainCertAuthorities | Create a Domain Certificate Authority object.
*DomainCertAuthorityApi* | [**DeleteDomainCertAuthority**](docs/DomainCertAuthorityApi.md#deletedomaincertauthority) | **Delete** /domainCertAuthorities/{certAuthorityName} | Delete a Domain Certificate Authority object.
*DomainCertAuthorityApi* | [**GetDomainCertAuthorities**](docs/DomainCertAuthorityApi.md#getdomaincertauthorities) | **Get** /domainCertAuthorities | Get a list of Domain Certificate Authority objects.
*DomainCertAuthorityApi* | [**GetDomainCertAuthority**](docs/DomainCertAuthorityApi.md#getdomaincertauthority) | **Get** /domainCertAuthorities/{certAuthorityName} | Get a Domain Certificate Authority object.
*DomainCertAuthorityApi* | [**ReplaceDomainCertAuthority**](docs/DomainCertAuthorityApi.md#replacedomaincertauthority) | **Put** /domainCertAuthorities/{certAuthorityName} | Replace a Domain Certificate Authority object.
*DomainCertAuthorityApi* | [**UpdateDomainCertAuthority**](docs/DomainCertAuthorityApi.md#updatedomaincertauthority) | **Patch** /domainCertAuthorities/{certAuthorityName} | Update a Domain Certificate Authority object.
*JndiApi* | [**CreateMsgVpnJndiConnectionFactory**](docs/JndiApi.md#createmsgvpnjndiconnectionfactory) | **Post** /msgVpns/{msgVpnName}/jndiConnectionFactories | Create a JNDI Connection Factory object.
*JndiApi* | [**CreateMsgVpnJndiQueue**](docs/JndiApi.md#createmsgvpnjndiqueue) | **Post** /msgVpns/{msgVpnName}/jndiQueues | Create a JNDI Queue object.
*JndiApi* | [**CreateMsgVpnJndiTopic**](docs/JndiApi.md#createmsgvpnjnditopic) | **Post** /msgVpns/{msgVpnName}/jndiTopics | Create a JNDI Topic object.
*JndiApi* | [**DeleteMsgVpnJndiConnectionFactory**](docs/JndiApi.md#deletemsgvpnjndiconnectionfactory) | **Delete** /msgVpns/{msgVpnName}/jndiConnectionFactories/{connectionFactoryName} | Delete a JNDI Connection Factory object.
*JndiApi* | [**DeleteMsgVpnJndiQueue**](docs/JndiApi.md#deletemsgvpnjndiqueue) | **Delete** /msgVpns/{msgVpnName}/jndiQueues/{queueName} | Delete a JNDI Queue object.
*JndiApi* | [**DeleteMsgVpnJndiTopic**](docs/JndiApi.md#deletemsgvpnjnditopic) | **Delete** /msgVpns/{msgVpnName}/jndiTopics/{topicName} | Delete a JNDI Topic object.
*JndiApi* | [**GetMsgVpnJndiConnectionFactories**](docs/JndiApi.md#getmsgvpnjndiconnectionfactories) | **Get** /msgVpns/{msgVpnName}/jndiConnectionFactories | Get a list of JNDI Connection Factory objects.
*JndiApi* | [**GetMsgVpnJndiConnectionFactory**](docs/JndiApi.md#getmsgvpnjndiconnectionfactory) | **Get** /msgVpns/{msgVpnName}/jndiConnectionFactories/{connectionFactoryName} | Get a JNDI Connection Factory object.
*JndiApi* | [**GetMsgVpnJndiQueue**](docs/JndiApi.md#getmsgvpnjndiqueue) | **Get** /msgVpns/{msgVpnName}/jndiQueues/{queueName} | Get a JNDI Queue object.
*JndiApi* | [**GetMsgVpnJndiQueues**](docs/JndiApi.md#getmsgvpnjndiqueues) | **Get** /msgVpns/{msgVpnName}/jndiQueues | Get a list of JNDI Queue objects.
*JndiApi* | [**GetMsgVpnJndiTopic**](docs/JndiApi.md#getmsgvpnjnditopic) | **Get** /msgVpns/{msgVpnName}/jndiTopics/{topicName} | Get a JNDI Topic object.
*JndiApi* | [**GetMsgVpnJndiTopics**](docs/JndiApi.md#getmsgvpnjnditopics) | **Get** /msgVpns/{msgVpnName}/jndiTopics | Get a list of JNDI Topic objects.
*JndiApi* | [**ReplaceMsgVpnJndiConnectionFactory**](docs/JndiApi.md#replacemsgvpnjndiconnectionfactory) | **Put** /msgVpns/{msgVpnName}/jndiConnectionFactories/{connectionFactoryName} | Replace a JNDI Connection Factory object.
*JndiApi* | [**ReplaceMsgVpnJndiQueue**](docs/JndiApi.md#replacemsgvpnjndiqueue) | **Put** /msgVpns/{msgVpnName}/jndiQueues/{queueName} | Replace a JNDI Queue object.
*JndiApi* | [**ReplaceMsgVpnJndiTopic**](docs/JndiApi.md#replacemsgvpnjnditopic) | **Put** /msgVpns/{msgVpnName}/jndiTopics/{topicName} | Replace a JNDI Topic object.
*JndiApi* | [**UpdateMsgVpnJndiConnectionFactory**](docs/JndiApi.md#updatemsgvpnjndiconnectionfactory) | **Patch** /msgVpns/{msgVpnName}/jndiConnectionFactories/{connectionFactoryName} | Update a JNDI Connection Factory object.
*JndiApi* | [**UpdateMsgVpnJndiQueue**](docs/JndiApi.md#updatemsgvpnjndiqueue) | **Patch** /msgVpns/{msgVpnName}/jndiQueues/{queueName} | Update a JNDI Queue object.
*JndiApi* | [**UpdateMsgVpnJndiTopic**](docs/JndiApi.md#updatemsgvpnjnditopic) | **Patch** /msgVpns/{msgVpnName}/jndiTopics/{topicName} | Update a JNDI Topic object.
*MqttRetainCacheApi* | [**CreateMsgVpnMqttRetainCache**](docs/MqttRetainCacheApi.md#createmsgvpnmqttretaincache) | **Post** /msgVpns/{msgVpnName}/mqttRetainCaches | Create an MQTT Retain Cache object.
*MqttRetainCacheApi* | [**DeleteMsgVpnMqttRetainCache**](docs/MqttRetainCacheApi.md#deletemsgvpnmqttretaincache) | **Delete** /msgVpns/{msgVpnName}/mqttRetainCaches/{cacheName} | Delete an MQTT Retain Cache object.
*MqttRetainCacheApi* | [**GetMsgVpnMqttRetainCache**](docs/MqttRetainCacheApi.md#getmsgvpnmqttretaincache) | **Get** /msgVpns/{msgVpnName}/mqttRetainCaches/{cacheName} | Get an MQTT Retain Cache object.
*MqttRetainCacheApi* | [**GetMsgVpnMqttRetainCaches**](docs/MqttRetainCacheApi.md#getmsgvpnmqttretaincaches) | **Get** /msgVpns/{msgVpnName}/mqttRetainCaches | Get a list of MQTT Retain Cache objects.
*MqttRetainCacheApi* | [**ReplaceMsgVpnMqttRetainCache**](docs/MqttRetainCacheApi.md#replacemsgvpnmqttretaincache) | **Put** /msgVpns/{msgVpnName}/mqttRetainCaches/{cacheName} | Replace an MQTT Retain Cache object.
*MqttRetainCacheApi* | [**UpdateMsgVpnMqttRetainCache**](docs/MqttRetainCacheApi.md#updatemsgvpnmqttretaincache) | **Patch** /msgVpns/{msgVpnName}/mqttRetainCaches/{cacheName} | Update an MQTT Retain Cache object.
*MqttSessionApi* | [**CreateMsgVpnMqttSession**](docs/MqttSessionApi.md#createmsgvpnmqttsession) | **Post** /msgVpns/{msgVpnName}/mqttSessions | Create an MQTT Session object.
*MqttSessionApi* | [**CreateMsgVpnMqttSessionSubscription**](docs/MqttSessionApi.md#createmsgvpnmqttsessionsubscription) | **Post** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions | Create a Subscription object.
*MqttSessionApi* | [**DeleteMsgVpnMqttSession**](docs/MqttSessionApi.md#deletemsgvpnmqttsession) | **Delete** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter} | Delete an MQTT Session object.
*MqttSessionApi* | [**DeleteMsgVpnMqttSessionSubscription**](docs/MqttSessionApi.md#deletemsgvpnmqttsessionsubscription) | **Delete** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic} | Delete a Subscription object.
*MqttSessionApi* | [**GetMsgVpnMqttSession**](docs/MqttSessionApi.md#getmsgvpnmqttsession) | **Get** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter} | Get an MQTT Session object.
*MqttSessionApi* | [**GetMsgVpnMqttSessionSubscription**](docs/MqttSessionApi.md#getmsgvpnmqttsessionsubscription) | **Get** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic} | Get a Subscription object.
*MqttSessionApi* | [**GetMsgVpnMqttSessionSubscriptions**](docs/MqttSessionApi.md#getmsgvpnmqttsessionsubscriptions) | **Get** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions | Get a list of Subscription objects.
*MqttSessionApi* | [**GetMsgVpnMqttSessions**](docs/MqttSessionApi.md#getmsgvpnmqttsessions) | **Get** /msgVpns/{msgVpnName}/mqttSessions | Get a list of MQTT Session objects.
*MqttSessionApi* | [**ReplaceMsgVpnMqttSession**](docs/MqttSessionApi.md#replacemsgvpnmqttsession) | **Put** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter} | Replace an MQTT Session object.
*MqttSessionApi* | [**ReplaceMsgVpnMqttSessionSubscription**](docs/MqttSessionApi.md#replacemsgvpnmqttsessionsubscription) | **Put** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic} | Replace a Subscription object.
*MqttSessionApi* | [**UpdateMsgVpnMqttSession**](docs/MqttSessionApi.md#updatemsgvpnmqttsession) | **Patch** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter} | Update an MQTT Session object.
*MqttSessionApi* | [**UpdateMsgVpnMqttSessionSubscription**](docs/MqttSessionApi.md#updatemsgvpnmqttsessionsubscription) | **Patch** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic} | Update a Subscription object.
*MsgVpnApi* | [**CreateMsgVpn**](docs/MsgVpnApi.md#createmsgvpn) | **Post** /msgVpns | Create a Message VPN object.
*MsgVpnApi* | [**CreateMsgVpnAclProfile**](docs/MsgVpnApi.md#createmsgvpnaclprofile) | **Post** /msgVpns/{msgVpnName}/aclProfiles | Create an ACL Profile object.
*MsgVpnApi* | [**CreateMsgVpnAclProfileClientConnectException**](docs/MsgVpnApi.md#createmsgvpnaclprofileclientconnectexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions | Create a Client Connect Exception object.
*MsgVpnApi* | [**CreateMsgVpnAclProfilePublishException**](docs/MsgVpnApi.md#createmsgvpnaclprofilepublishexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions | Create a Publish Topic Exception object.
*MsgVpnApi* | [**CreateMsgVpnAclProfilePublishTopicException**](docs/MsgVpnApi.md#createmsgvpnaclprofilepublishtopicexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions | Create a Publish Topic Exception object.
*MsgVpnApi* | [**CreateMsgVpnAclProfileSubscribeException**](docs/MsgVpnApi.md#createmsgvpnaclprofilesubscribeexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions | Create a Subscribe Topic Exception object.
*MsgVpnApi* | [**CreateMsgVpnAclProfileSubscribeShareNameException**](docs/MsgVpnApi.md#createmsgvpnaclprofilesubscribesharenameexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions | Create a Subscribe Share Name Exception object.
*MsgVpnApi* | [**CreateMsgVpnAclProfileSubscribeTopicException**](docs/MsgVpnApi.md#createmsgvpnaclprofilesubscribetopicexception) | **Post** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions | Create a Subscribe Topic Exception object.
*MsgVpnApi* | [**CreateMsgVpnAuthenticationOauthProfile**](docs/MsgVpnApi.md#createmsgvpnauthenticationoauthprofile) | **Post** /msgVpns/{msgVpnName}/authenticationOauthProfiles | Create an OAuth Profile object.
*MsgVpnApi* | [**CreateMsgVpnAuthenticationOauthProfileClientRequiredClaim**](docs/MsgVpnApi.md#createmsgvpnauthenticationoauthprofileclientrequiredclaim) | **Post** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/clientRequiredClaims | Create a Required Claim object.
*MsgVpnApi* | [**CreateMsgVpnAuthenticationOauthProfileResourceServerRequiredClaim**](docs/MsgVpnApi.md#createmsgvpnauthenticationoauthprofileresourceserverrequiredclaim) | **Post** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/resourceServerRequiredClaims | Create a Required Claim object.
*MsgVpnApi* | [**CreateMsgVpnAuthenticationOauthProvider**](docs/MsgVpnApi.md#createmsgvpnauthenticationoauthprovider) | **Post** /msgVpns/{msgVpnName}/authenticationOauthProviders | Create an OAuth Provider object.
*MsgVpnApi* | [**CreateMsgVpnAuthorizationGroup**](docs/MsgVpnApi.md#createmsgvpnauthorizationgroup) | **Post** /msgVpns/{msgVpnName}/authorizationGroups | Create an Authorization Group object.
*MsgVpnApi* | [**CreateMsgVpnBridge**](docs/MsgVpnApi.md#createmsgvpnbridge) | **Post** /msgVpns/{msgVpnName}/bridges | Create a Bridge object.
*MsgVpnApi* | [**CreateMsgVpnBridgeRemoteMsgVpn**](docs/MsgVpnApi.md#createmsgvpnbridgeremotemsgvpn) | **Post** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns | Create a Remote Message VPN object.
*MsgVpnApi* | [**CreateMsgVpnBridgeRemoteSubscription**](docs/MsgVpnApi.md#createmsgvpnbridgeremotesubscription) | **Post** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions | Create a Remote Subscription object.
*MsgVpnApi* | [**CreateMsgVpnBridgeTlsTrustedCommonName**](docs/MsgVpnApi.md#createmsgvpnbridgetlstrustedcommonname) | **Post** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames | Create a Trusted Common Name object.
*MsgVpnApi* | [**CreateMsgVpnClientProfile**](docs/MsgVpnApi.md#createmsgvpnclientprofile) | **Post** /msgVpns/{msgVpnName}/clientProfiles | Create a Client Profile object.
*MsgVpnApi* | [**CreateMsgVpnClientUsername**](docs/MsgVpnApi.md#createmsgvpnclientusername) | **Post** /msgVpns/{msgVpnName}/clientUsernames | Create a Client Username object.
*MsgVpnApi* | [**CreateMsgVpnDistributedCache**](docs/MsgVpnApi.md#createmsgvpndistributedcache) | **Post** /msgVpns/{msgVpnName}/distributedCaches | Create a Distributed Cache object.
*MsgVpnApi* | [**CreateMsgVpnDistributedCacheCluster**](docs/MsgVpnApi.md#createmsgvpndistributedcachecluster) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters | Create a Cache Cluster object.
*MsgVpnApi* | [**CreateMsgVpnDistributedCacheClusterGlobalCachingHomeCluster**](docs/MsgVpnApi.md#createmsgvpndistributedcacheclusterglobalcachinghomecluster) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters | Create a Home Cache Cluster object.
*MsgVpnApi* | [**CreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix**](docs/MsgVpnApi.md#createmsgvpndistributedcacheclusterglobalcachinghomeclustertopicprefix) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes | Create a Topic Prefix object.
*MsgVpnApi* | [**CreateMsgVpnDistributedCacheClusterInstance**](docs/MsgVpnApi.md#createmsgvpndistributedcacheclusterinstance) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances | Create a Cache Instance object.
*MsgVpnApi* | [**CreateMsgVpnDistributedCacheClusterTopic**](docs/MsgVpnApi.md#createmsgvpndistributedcacheclustertopic) | **Post** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics | Create a Topic object.
*MsgVpnApi* | [**CreateMsgVpnDmrBridge**](docs/MsgVpnApi.md#createmsgvpndmrbridge) | **Post** /msgVpns/{msgVpnName}/dmrBridges | Create a DMR Bridge object.
*MsgVpnApi* | [**CreateMsgVpnJndiConnectionFactory**](docs/MsgVpnApi.md#createmsgvpnjndiconnectionfactory) | **Post** /msgVpns/{msgVpnName}/jndiConnectionFactories | Create a JNDI Connection Factory object.
*MsgVpnApi* | [**CreateMsgVpnJndiQueue**](docs/MsgVpnApi.md#createmsgvpnjndiqueue) | **Post** /msgVpns/{msgVpnName}/jndiQueues | Create a JNDI Queue object.
*MsgVpnApi* | [**CreateMsgVpnJndiTopic**](docs/MsgVpnApi.md#createmsgvpnjnditopic) | **Post** /msgVpns/{msgVpnName}/jndiTopics | Create a JNDI Topic object.
*MsgVpnApi* | [**CreateMsgVpnMqttRetainCache**](docs/MsgVpnApi.md#createmsgvpnmqttretaincache) | **Post** /msgVpns/{msgVpnName}/mqttRetainCaches | Create an MQTT Retain Cache object.
*MsgVpnApi* | [**CreateMsgVpnMqttSession**](docs/MsgVpnApi.md#createmsgvpnmqttsession) | **Post** /msgVpns/{msgVpnName}/mqttSessions | Create an MQTT Session object.
*MsgVpnApi* | [**CreateMsgVpnMqttSessionSubscription**](docs/MsgVpnApi.md#createmsgvpnmqttsessionsubscription) | **Post** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions | Create a Subscription object.
*MsgVpnApi* | [**CreateMsgVpnQueue**](docs/MsgVpnApi.md#createmsgvpnqueue) | **Post** /msgVpns/{msgVpnName}/queues | Create a Queue object.
*MsgVpnApi* | [**CreateMsgVpnQueueSubscription**](docs/MsgVpnApi.md#createmsgvpnqueuesubscription) | **Post** /msgVpns/{msgVpnName}/queues/{queueName}/subscriptions | Create a Queue Subscription object.
*MsgVpnApi* | [**CreateMsgVpnQueueTemplate**](docs/MsgVpnApi.md#createmsgvpnqueuetemplate) | **Post** /msgVpns/{msgVpnName}/queueTemplates | Create a Queue Template object.
*MsgVpnApi* | [**CreateMsgVpnReplayLog**](docs/MsgVpnApi.md#createmsgvpnreplaylog) | **Post** /msgVpns/{msgVpnName}/replayLogs | Create a Replay Log object.
*MsgVpnApi* | [**CreateMsgVpnReplicatedTopic**](docs/MsgVpnApi.md#createmsgvpnreplicatedtopic) | **Post** /msgVpns/{msgVpnName}/replicatedTopics | Create a Replicated Topic object.
*MsgVpnApi* | [**CreateMsgVpnRestDeliveryPoint**](docs/MsgVpnApi.md#createmsgvpnrestdeliverypoint) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints | Create a REST Delivery Point object.
*MsgVpnApi* | [**CreateMsgVpnRestDeliveryPointQueueBinding**](docs/MsgVpnApi.md#createmsgvpnrestdeliverypointqueuebinding) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings | Create a Queue Binding object.
*MsgVpnApi* | [**CreateMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/MsgVpnApi.md#createmsgvpnrestdeliverypointqueuebindingrequestheader) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders | Create a Request Header object.
*MsgVpnApi* | [**CreateMsgVpnRestDeliveryPointRestConsumer**](docs/MsgVpnApi.md#createmsgvpnrestdeliverypointrestconsumer) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers | Create a REST Consumer object.
*MsgVpnApi* | [**CreateMsgVpnRestDeliveryPointRestConsumerOauthJwtClaim**](docs/MsgVpnApi.md#createmsgvpnrestdeliverypointrestconsumeroauthjwtclaim) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/oauthJwtClaims | Create a Claim object.
*MsgVpnApi* | [**CreateMsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonName**](docs/MsgVpnApi.md#createmsgvpnrestdeliverypointrestconsumertlstrustedcommonname) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/tlsTrustedCommonNames | Create a Trusted Common Name object.
*MsgVpnApi* | [**CreateMsgVpnSequencedTopic**](docs/MsgVpnApi.md#createmsgvpnsequencedtopic) | **Post** /msgVpns/{msgVpnName}/sequencedTopics | Create a Sequenced Topic object.
*MsgVpnApi* | [**CreateMsgVpnTopicEndpoint**](docs/MsgVpnApi.md#createmsgvpntopicendpoint) | **Post** /msgVpns/{msgVpnName}/topicEndpoints | Create a Topic Endpoint object.
*MsgVpnApi* | [**CreateMsgVpnTopicEndpointTemplate**](docs/MsgVpnApi.md#createmsgvpntopicendpointtemplate) | **Post** /msgVpns/{msgVpnName}/topicEndpointTemplates | Create a Topic Endpoint Template object.
*MsgVpnApi* | [**DeleteMsgVpn**](docs/MsgVpnApi.md#deletemsgvpn) | **Delete** /msgVpns/{msgVpnName} | Delete a Message VPN object.
*MsgVpnApi* | [**DeleteMsgVpnAclProfile**](docs/MsgVpnApi.md#deletemsgvpnaclprofile) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName} | Delete an ACL Profile object.
*MsgVpnApi* | [**DeleteMsgVpnAclProfileClientConnectException**](docs/MsgVpnApi.md#deletemsgvpnaclprofileclientconnectexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions/{clientConnectExceptionAddress} | Delete a Client Connect Exception object.
*MsgVpnApi* | [**DeleteMsgVpnAclProfilePublishException**](docs/MsgVpnApi.md#deletemsgvpnaclprofilepublishexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions/{topicSyntax},{publishExceptionTopic} | Delete a Publish Topic Exception object.
*MsgVpnApi* | [**DeleteMsgVpnAclProfilePublishTopicException**](docs/MsgVpnApi.md#deletemsgvpnaclprofilepublishtopicexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions/{publishTopicExceptionSyntax},{publishTopicException} | Delete a Publish Topic Exception object.
*MsgVpnApi* | [**DeleteMsgVpnAclProfileSubscribeException**](docs/MsgVpnApi.md#deletemsgvpnaclprofilesubscribeexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions/{topicSyntax},{subscribeExceptionTopic} | Delete a Subscribe Topic Exception object.
*MsgVpnApi* | [**DeleteMsgVpnAclProfileSubscribeShareNameException**](docs/MsgVpnApi.md#deletemsgvpnaclprofilesubscribesharenameexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions/{subscribeShareNameExceptionSyntax},{subscribeShareNameException} | Delete a Subscribe Share Name Exception object.
*MsgVpnApi* | [**DeleteMsgVpnAclProfileSubscribeTopicException**](docs/MsgVpnApi.md#deletemsgvpnaclprofilesubscribetopicexception) | **Delete** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions/{subscribeTopicExceptionSyntax},{subscribeTopicException} | Delete a Subscribe Topic Exception object.
*MsgVpnApi* | [**DeleteMsgVpnAuthenticationOauthProfile**](docs/MsgVpnApi.md#deletemsgvpnauthenticationoauthprofile) | **Delete** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName} | Delete an OAuth Profile object.
*MsgVpnApi* | [**DeleteMsgVpnAuthenticationOauthProfileClientRequiredClaim**](docs/MsgVpnApi.md#deletemsgvpnauthenticationoauthprofileclientrequiredclaim) | **Delete** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/clientRequiredClaims/{clientRequiredClaimName} | Delete a Required Claim object.
*MsgVpnApi* | [**DeleteMsgVpnAuthenticationOauthProfileResourceServerRequiredClaim**](docs/MsgVpnApi.md#deletemsgvpnauthenticationoauthprofileresourceserverrequiredclaim) | **Delete** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/resourceServerRequiredClaims/{resourceServerRequiredClaimName} | Delete a Required Claim object.
*MsgVpnApi* | [**DeleteMsgVpnAuthenticationOauthProvider**](docs/MsgVpnApi.md#deletemsgvpnauthenticationoauthprovider) | **Delete** /msgVpns/{msgVpnName}/authenticationOauthProviders/{oauthProviderName} | Delete an OAuth Provider object.
*MsgVpnApi* | [**DeleteMsgVpnAuthorizationGroup**](docs/MsgVpnApi.md#deletemsgvpnauthorizationgroup) | **Delete** /msgVpns/{msgVpnName}/authorizationGroups/{authorizationGroupName} | Delete an Authorization Group object.
*MsgVpnApi* | [**DeleteMsgVpnBridge**](docs/MsgVpnApi.md#deletemsgvpnbridge) | **Delete** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter} | Delete a Bridge object.
*MsgVpnApi* | [**DeleteMsgVpnBridgeRemoteMsgVpn**](docs/MsgVpnApi.md#deletemsgvpnbridgeremotemsgvpn) | **Delete** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface} | Delete a Remote Message VPN object.
*MsgVpnApi* | [**DeleteMsgVpnBridgeRemoteSubscription**](docs/MsgVpnApi.md#deletemsgvpnbridgeremotesubscription) | **Delete** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions/{remoteSubscriptionTopic} | Delete a Remote Subscription object.
*MsgVpnApi* | [**DeleteMsgVpnBridgeTlsTrustedCommonName**](docs/MsgVpnApi.md#deletemsgvpnbridgetlstrustedcommonname) | **Delete** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Delete a Trusted Common Name object.
*MsgVpnApi* | [**DeleteMsgVpnClientProfile**](docs/MsgVpnApi.md#deletemsgvpnclientprofile) | **Delete** /msgVpns/{msgVpnName}/clientProfiles/{clientProfileName} | Delete a Client Profile object.
*MsgVpnApi* | [**DeleteMsgVpnClientUsername**](docs/MsgVpnApi.md#deletemsgvpnclientusername) | **Delete** /msgVpns/{msgVpnName}/clientUsernames/{clientUsername} | Delete a Client Username object.
*MsgVpnApi* | [**DeleteMsgVpnDistributedCache**](docs/MsgVpnApi.md#deletemsgvpndistributedcache) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName} | Delete a Distributed Cache object.
*MsgVpnApi* | [**DeleteMsgVpnDistributedCacheCluster**](docs/MsgVpnApi.md#deletemsgvpndistributedcachecluster) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName} | Delete a Cache Cluster object.
*MsgVpnApi* | [**DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeCluster**](docs/MsgVpnApi.md#deletemsgvpndistributedcacheclusterglobalcachinghomecluster) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName} | Delete a Home Cache Cluster object.
*MsgVpnApi* | [**DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix**](docs/MsgVpnApi.md#deletemsgvpndistributedcacheclusterglobalcachinghomeclustertopicprefix) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes/{topicPrefix} | Delete a Topic Prefix object.
*MsgVpnApi* | [**DeleteMsgVpnDistributedCacheClusterInstance**](docs/MsgVpnApi.md#deletemsgvpndistributedcacheclusterinstance) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName} | Delete a Cache Instance object.
*MsgVpnApi* | [**DeleteMsgVpnDistributedCacheClusterTopic**](docs/MsgVpnApi.md#deletemsgvpndistributedcacheclustertopic) | **Delete** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics/{topic} | Delete a Topic object.
*MsgVpnApi* | [**DeleteMsgVpnDmrBridge**](docs/MsgVpnApi.md#deletemsgvpndmrbridge) | **Delete** /msgVpns/{msgVpnName}/dmrBridges/{remoteNodeName} | Delete a DMR Bridge object.
*MsgVpnApi* | [**DeleteMsgVpnJndiConnectionFactory**](docs/MsgVpnApi.md#deletemsgvpnjndiconnectionfactory) | **Delete** /msgVpns/{msgVpnName}/jndiConnectionFactories/{connectionFactoryName} | Delete a JNDI Connection Factory object.
*MsgVpnApi* | [**DeleteMsgVpnJndiQueue**](docs/MsgVpnApi.md#deletemsgvpnjndiqueue) | **Delete** /msgVpns/{msgVpnName}/jndiQueues/{queueName} | Delete a JNDI Queue object.
*MsgVpnApi* | [**DeleteMsgVpnJndiTopic**](docs/MsgVpnApi.md#deletemsgvpnjnditopic) | **Delete** /msgVpns/{msgVpnName}/jndiTopics/{topicName} | Delete a JNDI Topic object.
*MsgVpnApi* | [**DeleteMsgVpnMqttRetainCache**](docs/MsgVpnApi.md#deletemsgvpnmqttretaincache) | **Delete** /msgVpns/{msgVpnName}/mqttRetainCaches/{cacheName} | Delete an MQTT Retain Cache object.
*MsgVpnApi* | [**DeleteMsgVpnMqttSession**](docs/MsgVpnApi.md#deletemsgvpnmqttsession) | **Delete** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter} | Delete an MQTT Session object.
*MsgVpnApi* | [**DeleteMsgVpnMqttSessionSubscription**](docs/MsgVpnApi.md#deletemsgvpnmqttsessionsubscription) | **Delete** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic} | Delete a Subscription object.
*MsgVpnApi* | [**DeleteMsgVpnQueue**](docs/MsgVpnApi.md#deletemsgvpnqueue) | **Delete** /msgVpns/{msgVpnName}/queues/{queueName} | Delete a Queue object.
*MsgVpnApi* | [**DeleteMsgVpnQueueSubscription**](docs/MsgVpnApi.md#deletemsgvpnqueuesubscription) | **Delete** /msgVpns/{msgVpnName}/queues/{queueName}/subscriptions/{subscriptionTopic} | Delete a Queue Subscription object.
*MsgVpnApi* | [**DeleteMsgVpnQueueTemplate**](docs/MsgVpnApi.md#deletemsgvpnqueuetemplate) | **Delete** /msgVpns/{msgVpnName}/queueTemplates/{queueTemplateName} | Delete a Queue Template object.
*MsgVpnApi* | [**DeleteMsgVpnReplayLog**](docs/MsgVpnApi.md#deletemsgvpnreplaylog) | **Delete** /msgVpns/{msgVpnName}/replayLogs/{replayLogName} | Delete a Replay Log object.
*MsgVpnApi* | [**DeleteMsgVpnReplicatedTopic**](docs/MsgVpnApi.md#deletemsgvpnreplicatedtopic) | **Delete** /msgVpns/{msgVpnName}/replicatedTopics/{replicatedTopic} | Delete a Replicated Topic object.
*MsgVpnApi* | [**DeleteMsgVpnRestDeliveryPoint**](docs/MsgVpnApi.md#deletemsgvpnrestdeliverypoint) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName} | Delete a REST Delivery Point object.
*MsgVpnApi* | [**DeleteMsgVpnRestDeliveryPointQueueBinding**](docs/MsgVpnApi.md#deletemsgvpnrestdeliverypointqueuebinding) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName} | Delete a Queue Binding object.
*MsgVpnApi* | [**DeleteMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/MsgVpnApi.md#deletemsgvpnrestdeliverypointqueuebindingrequestheader) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders/{headerName} | Delete a Request Header object.
*MsgVpnApi* | [**DeleteMsgVpnRestDeliveryPointRestConsumer**](docs/MsgVpnApi.md#deletemsgvpnrestdeliverypointrestconsumer) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName} | Delete a REST Consumer object.
*MsgVpnApi* | [**DeleteMsgVpnRestDeliveryPointRestConsumerOauthJwtClaim**](docs/MsgVpnApi.md#deletemsgvpnrestdeliverypointrestconsumeroauthjwtclaim) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/oauthJwtClaims/{oauthJwtClaimName} | Delete a Claim object.
*MsgVpnApi* | [**DeleteMsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonName**](docs/MsgVpnApi.md#deletemsgvpnrestdeliverypointrestconsumertlstrustedcommonname) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Delete a Trusted Common Name object.
*MsgVpnApi* | [**DeleteMsgVpnSequencedTopic**](docs/MsgVpnApi.md#deletemsgvpnsequencedtopic) | **Delete** /msgVpns/{msgVpnName}/sequencedTopics/{sequencedTopic} | Delete a Sequenced Topic object.
*MsgVpnApi* | [**DeleteMsgVpnTopicEndpoint**](docs/MsgVpnApi.md#deletemsgvpntopicendpoint) | **Delete** /msgVpns/{msgVpnName}/topicEndpoints/{topicEndpointName} | Delete a Topic Endpoint object.
*MsgVpnApi* | [**DeleteMsgVpnTopicEndpointTemplate**](docs/MsgVpnApi.md#deletemsgvpntopicendpointtemplate) | **Delete** /msgVpns/{msgVpnName}/topicEndpointTemplates/{topicEndpointTemplateName} | Delete a Topic Endpoint Template object.
*MsgVpnApi* | [**GetMsgVpn**](docs/MsgVpnApi.md#getmsgvpn) | **Get** /msgVpns/{msgVpnName} | Get a Message VPN object.
*MsgVpnApi* | [**GetMsgVpnAclProfile**](docs/MsgVpnApi.md#getmsgvpnaclprofile) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName} | Get an ACL Profile object.
*MsgVpnApi* | [**GetMsgVpnAclProfileClientConnectException**](docs/MsgVpnApi.md#getmsgvpnaclprofileclientconnectexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions/{clientConnectExceptionAddress} | Get a Client Connect Exception object.
*MsgVpnApi* | [**GetMsgVpnAclProfileClientConnectExceptions**](docs/MsgVpnApi.md#getmsgvpnaclprofileclientconnectexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions | Get a list of Client Connect Exception objects.
*MsgVpnApi* | [**GetMsgVpnAclProfilePublishException**](docs/MsgVpnApi.md#getmsgvpnaclprofilepublishexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions/{topicSyntax},{publishExceptionTopic} | Get a Publish Topic Exception object.
*MsgVpnApi* | [**GetMsgVpnAclProfilePublishExceptions**](docs/MsgVpnApi.md#getmsgvpnaclprofilepublishexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions | Get a list of Publish Topic Exception objects.
*MsgVpnApi* | [**GetMsgVpnAclProfilePublishTopicException**](docs/MsgVpnApi.md#getmsgvpnaclprofilepublishtopicexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions/{publishTopicExceptionSyntax},{publishTopicException} | Get a Publish Topic Exception object.
*MsgVpnApi* | [**GetMsgVpnAclProfilePublishTopicExceptions**](docs/MsgVpnApi.md#getmsgvpnaclprofilepublishtopicexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions | Get a list of Publish Topic Exception objects.
*MsgVpnApi* | [**GetMsgVpnAclProfileSubscribeException**](docs/MsgVpnApi.md#getmsgvpnaclprofilesubscribeexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions/{topicSyntax},{subscribeExceptionTopic} | Get a Subscribe Topic Exception object.
*MsgVpnApi* | [**GetMsgVpnAclProfileSubscribeExceptions**](docs/MsgVpnApi.md#getmsgvpnaclprofilesubscribeexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions | Get a list of Subscribe Topic Exception objects.
*MsgVpnApi* | [**GetMsgVpnAclProfileSubscribeShareNameException**](docs/MsgVpnApi.md#getmsgvpnaclprofilesubscribesharenameexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions/{subscribeShareNameExceptionSyntax},{subscribeShareNameException} | Get a Subscribe Share Name Exception object.
*MsgVpnApi* | [**GetMsgVpnAclProfileSubscribeShareNameExceptions**](docs/MsgVpnApi.md#getmsgvpnaclprofilesubscribesharenameexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions | Get a list of Subscribe Share Name Exception objects.
*MsgVpnApi* | [**GetMsgVpnAclProfileSubscribeTopicException**](docs/MsgVpnApi.md#getmsgvpnaclprofilesubscribetopicexception) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions/{subscribeTopicExceptionSyntax},{subscribeTopicException} | Get a Subscribe Topic Exception object.
*MsgVpnApi* | [**GetMsgVpnAclProfileSubscribeTopicExceptions**](docs/MsgVpnApi.md#getmsgvpnaclprofilesubscribetopicexceptions) | **Get** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions | Get a list of Subscribe Topic Exception objects.
*MsgVpnApi* | [**GetMsgVpnAclProfiles**](docs/MsgVpnApi.md#getmsgvpnaclprofiles) | **Get** /msgVpns/{msgVpnName}/aclProfiles | Get a list of ACL Profile objects.
*MsgVpnApi* | [**GetMsgVpnAuthenticationOauthProfile**](docs/MsgVpnApi.md#getmsgvpnauthenticationoauthprofile) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName} | Get an OAuth Profile object.
*MsgVpnApi* | [**GetMsgVpnAuthenticationOauthProfileClientRequiredClaim**](docs/MsgVpnApi.md#getmsgvpnauthenticationoauthprofileclientrequiredclaim) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/clientRequiredClaims/{clientRequiredClaimName} | Get a Required Claim object.
*MsgVpnApi* | [**GetMsgVpnAuthenticationOauthProfileClientRequiredClaims**](docs/MsgVpnApi.md#getmsgvpnauthenticationoauthprofileclientrequiredclaims) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/clientRequiredClaims | Get a list of Required Claim objects.
*MsgVpnApi* | [**GetMsgVpnAuthenticationOauthProfileResourceServerRequiredClaim**](docs/MsgVpnApi.md#getmsgvpnauthenticationoauthprofileresourceserverrequiredclaim) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/resourceServerRequiredClaims/{resourceServerRequiredClaimName} | Get a Required Claim object.
*MsgVpnApi* | [**GetMsgVpnAuthenticationOauthProfileResourceServerRequiredClaims**](docs/MsgVpnApi.md#getmsgvpnauthenticationoauthprofileresourceserverrequiredclaims) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName}/resourceServerRequiredClaims | Get a list of Required Claim objects.
*MsgVpnApi* | [**GetMsgVpnAuthenticationOauthProfiles**](docs/MsgVpnApi.md#getmsgvpnauthenticationoauthprofiles) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProfiles | Get a list of OAuth Profile objects.
*MsgVpnApi* | [**GetMsgVpnAuthenticationOauthProvider**](docs/MsgVpnApi.md#getmsgvpnauthenticationoauthprovider) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProviders/{oauthProviderName} | Get an OAuth Provider object.
*MsgVpnApi* | [**GetMsgVpnAuthenticationOauthProviders**](docs/MsgVpnApi.md#getmsgvpnauthenticationoauthproviders) | **Get** /msgVpns/{msgVpnName}/authenticationOauthProviders | Get a list of OAuth Provider objects.
*MsgVpnApi* | [**GetMsgVpnAuthorizationGroup**](docs/MsgVpnApi.md#getmsgvpnauthorizationgroup) | **Get** /msgVpns/{msgVpnName}/authorizationGroups/{authorizationGroupName} | Get an Authorization Group object.
*MsgVpnApi* | [**GetMsgVpnAuthorizationGroups**](docs/MsgVpnApi.md#getmsgvpnauthorizationgroups) | **Get** /msgVpns/{msgVpnName}/authorizationGroups | Get a list of Authorization Group objects.
*MsgVpnApi* | [**GetMsgVpnBridge**](docs/MsgVpnApi.md#getmsgvpnbridge) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter} | Get a Bridge object.
*MsgVpnApi* | [**GetMsgVpnBridgeRemoteMsgVpn**](docs/MsgVpnApi.md#getmsgvpnbridgeremotemsgvpn) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface} | Get a Remote Message VPN object.
*MsgVpnApi* | [**GetMsgVpnBridgeRemoteMsgVpns**](docs/MsgVpnApi.md#getmsgvpnbridgeremotemsgvpns) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns | Get a list of Remote Message VPN objects.
*MsgVpnApi* | [**GetMsgVpnBridgeRemoteSubscription**](docs/MsgVpnApi.md#getmsgvpnbridgeremotesubscription) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions/{remoteSubscriptionTopic} | Get a Remote Subscription object.
*MsgVpnApi* | [**GetMsgVpnBridgeRemoteSubscriptions**](docs/MsgVpnApi.md#getmsgvpnbridgeremotesubscriptions) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions | Get a list of Remote Subscription objects.
*MsgVpnApi* | [**GetMsgVpnBridgeTlsTrustedCommonName**](docs/MsgVpnApi.md#getmsgvpnbridgetlstrustedcommonname) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Get a Trusted Common Name object.
*MsgVpnApi* | [**GetMsgVpnBridgeTlsTrustedCommonNames**](docs/MsgVpnApi.md#getmsgvpnbridgetlstrustedcommonnames) | **Get** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames | Get a list of Trusted Common Name objects.
*MsgVpnApi* | [**GetMsgVpnBridges**](docs/MsgVpnApi.md#getmsgvpnbridges) | **Get** /msgVpns/{msgVpnName}/bridges | Get a list of Bridge objects.
*MsgVpnApi* | [**GetMsgVpnClientProfile**](docs/MsgVpnApi.md#getmsgvpnclientprofile) | **Get** /msgVpns/{msgVpnName}/clientProfiles/{clientProfileName} | Get a Client Profile object.
*MsgVpnApi* | [**GetMsgVpnClientProfiles**](docs/MsgVpnApi.md#getmsgvpnclientprofiles) | **Get** /msgVpns/{msgVpnName}/clientProfiles | Get a list of Client Profile objects.
*MsgVpnApi* | [**GetMsgVpnClientUsername**](docs/MsgVpnApi.md#getmsgvpnclientusername) | **Get** /msgVpns/{msgVpnName}/clientUsernames/{clientUsername} | Get a Client Username object.
*MsgVpnApi* | [**GetMsgVpnClientUsernames**](docs/MsgVpnApi.md#getmsgvpnclientusernames) | **Get** /msgVpns/{msgVpnName}/clientUsernames | Get a list of Client Username objects.
*MsgVpnApi* | [**GetMsgVpnDistributedCache**](docs/MsgVpnApi.md#getmsgvpndistributedcache) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName} | Get a Distributed Cache object.
*MsgVpnApi* | [**GetMsgVpnDistributedCacheCluster**](docs/MsgVpnApi.md#getmsgvpndistributedcachecluster) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName} | Get a Cache Cluster object.
*MsgVpnApi* | [**GetMsgVpnDistributedCacheClusterGlobalCachingHomeCluster**](docs/MsgVpnApi.md#getmsgvpndistributedcacheclusterglobalcachinghomecluster) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName} | Get a Home Cache Cluster object.
*MsgVpnApi* | [**GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix**](docs/MsgVpnApi.md#getmsgvpndistributedcacheclusterglobalcachinghomeclustertopicprefix) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes/{topicPrefix} | Get a Topic Prefix object.
*MsgVpnApi* | [**GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixes**](docs/MsgVpnApi.md#getmsgvpndistributedcacheclusterglobalcachinghomeclustertopicprefixes) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes | Get a list of Topic Prefix objects.
*MsgVpnApi* | [**GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusters**](docs/MsgVpnApi.md#getmsgvpndistributedcacheclusterglobalcachinghomeclusters) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters | Get a list of Home Cache Cluster objects.
*MsgVpnApi* | [**GetMsgVpnDistributedCacheClusterInstance**](docs/MsgVpnApi.md#getmsgvpndistributedcacheclusterinstance) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName} | Get a Cache Instance object.
*MsgVpnApi* | [**GetMsgVpnDistributedCacheClusterInstances**](docs/MsgVpnApi.md#getmsgvpndistributedcacheclusterinstances) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances | Get a list of Cache Instance objects.
*MsgVpnApi* | [**GetMsgVpnDistributedCacheClusterTopic**](docs/MsgVpnApi.md#getmsgvpndistributedcacheclustertopic) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics/{topic} | Get a Topic object.
*MsgVpnApi* | [**GetMsgVpnDistributedCacheClusterTopics**](docs/MsgVpnApi.md#getmsgvpndistributedcacheclustertopics) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics | Get a list of Topic objects.
*MsgVpnApi* | [**GetMsgVpnDistributedCacheClusters**](docs/MsgVpnApi.md#getmsgvpndistributedcacheclusters) | **Get** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters | Get a list of Cache Cluster objects.
*MsgVpnApi* | [**GetMsgVpnDistributedCaches**](docs/MsgVpnApi.md#getmsgvpndistributedcaches) | **Get** /msgVpns/{msgVpnName}/distributedCaches | Get a list of Distributed Cache objects.
*MsgVpnApi* | [**GetMsgVpnDmrBridge**](docs/MsgVpnApi.md#getmsgvpndmrbridge) | **Get** /msgVpns/{msgVpnName}/dmrBridges/{remoteNodeName} | Get a DMR Bridge object.
*MsgVpnApi* | [**GetMsgVpnDmrBridges**](docs/MsgVpnApi.md#getmsgvpndmrbridges) | **Get** /msgVpns/{msgVpnName}/dmrBridges | Get a list of DMR Bridge objects.
*MsgVpnApi* | [**GetMsgVpnJndiConnectionFactories**](docs/MsgVpnApi.md#getmsgvpnjndiconnectionfactories) | **Get** /msgVpns/{msgVpnName}/jndiConnectionFactories | Get a list of JNDI Connection Factory objects.
*MsgVpnApi* | [**GetMsgVpnJndiConnectionFactory**](docs/MsgVpnApi.md#getmsgvpnjndiconnectionfactory) | **Get** /msgVpns/{msgVpnName}/jndiConnectionFactories/{connectionFactoryName} | Get a JNDI Connection Factory object.
*MsgVpnApi* | [**GetMsgVpnJndiQueue**](docs/MsgVpnApi.md#getmsgvpnjndiqueue) | **Get** /msgVpns/{msgVpnName}/jndiQueues/{queueName} | Get a JNDI Queue object.
*MsgVpnApi* | [**GetMsgVpnJndiQueues**](docs/MsgVpnApi.md#getmsgvpnjndiqueues) | **Get** /msgVpns/{msgVpnName}/jndiQueues | Get a list of JNDI Queue objects.
*MsgVpnApi* | [**GetMsgVpnJndiTopic**](docs/MsgVpnApi.md#getmsgvpnjnditopic) | **Get** /msgVpns/{msgVpnName}/jndiTopics/{topicName} | Get a JNDI Topic object.
*MsgVpnApi* | [**GetMsgVpnJndiTopics**](docs/MsgVpnApi.md#getmsgvpnjnditopics) | **Get** /msgVpns/{msgVpnName}/jndiTopics | Get a list of JNDI Topic objects.
*MsgVpnApi* | [**GetMsgVpnMqttRetainCache**](docs/MsgVpnApi.md#getmsgvpnmqttretaincache) | **Get** /msgVpns/{msgVpnName}/mqttRetainCaches/{cacheName} | Get an MQTT Retain Cache object.
*MsgVpnApi* | [**GetMsgVpnMqttRetainCaches**](docs/MsgVpnApi.md#getmsgvpnmqttretaincaches) | **Get** /msgVpns/{msgVpnName}/mqttRetainCaches | Get a list of MQTT Retain Cache objects.
*MsgVpnApi* | [**GetMsgVpnMqttSession**](docs/MsgVpnApi.md#getmsgvpnmqttsession) | **Get** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter} | Get an MQTT Session object.
*MsgVpnApi* | [**GetMsgVpnMqttSessionSubscription**](docs/MsgVpnApi.md#getmsgvpnmqttsessionsubscription) | **Get** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic} | Get a Subscription object.
*MsgVpnApi* | [**GetMsgVpnMqttSessionSubscriptions**](docs/MsgVpnApi.md#getmsgvpnmqttsessionsubscriptions) | **Get** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions | Get a list of Subscription objects.
*MsgVpnApi* | [**GetMsgVpnMqttSessions**](docs/MsgVpnApi.md#getmsgvpnmqttsessions) | **Get** /msgVpns/{msgVpnName}/mqttSessions | Get a list of MQTT Session objects.
*MsgVpnApi* | [**GetMsgVpnQueue**](docs/MsgVpnApi.md#getmsgvpnqueue) | **Get** /msgVpns/{msgVpnName}/queues/{queueName} | Get a Queue object.
*MsgVpnApi* | [**GetMsgVpnQueueSubscription**](docs/MsgVpnApi.md#getmsgvpnqueuesubscription) | **Get** /msgVpns/{msgVpnName}/queues/{queueName}/subscriptions/{subscriptionTopic} | Get a Queue Subscription object.
*MsgVpnApi* | [**GetMsgVpnQueueSubscriptions**](docs/MsgVpnApi.md#getmsgvpnqueuesubscriptions) | **Get** /msgVpns/{msgVpnName}/queues/{queueName}/subscriptions | Get a list of Queue Subscription objects.
*MsgVpnApi* | [**GetMsgVpnQueueTemplate**](docs/MsgVpnApi.md#getmsgvpnqueuetemplate) | **Get** /msgVpns/{msgVpnName}/queueTemplates/{queueTemplateName} | Get a Queue Template object.
*MsgVpnApi* | [**GetMsgVpnQueueTemplates**](docs/MsgVpnApi.md#getmsgvpnqueuetemplates) | **Get** /msgVpns/{msgVpnName}/queueTemplates | Get a list of Queue Template objects.
*MsgVpnApi* | [**GetMsgVpnQueues**](docs/MsgVpnApi.md#getmsgvpnqueues) | **Get** /msgVpns/{msgVpnName}/queues | Get a list of Queue objects.
*MsgVpnApi* | [**GetMsgVpnReplayLog**](docs/MsgVpnApi.md#getmsgvpnreplaylog) | **Get** /msgVpns/{msgVpnName}/replayLogs/{replayLogName} | Get a Replay Log object.
*MsgVpnApi* | [**GetMsgVpnReplayLogs**](docs/MsgVpnApi.md#getmsgvpnreplaylogs) | **Get** /msgVpns/{msgVpnName}/replayLogs | Get a list of Replay Log objects.
*MsgVpnApi* | [**GetMsgVpnReplicatedTopic**](docs/MsgVpnApi.md#getmsgvpnreplicatedtopic) | **Get** /msgVpns/{msgVpnName}/replicatedTopics/{replicatedTopic} | Get a Replicated Topic object.
*MsgVpnApi* | [**GetMsgVpnReplicatedTopics**](docs/MsgVpnApi.md#getmsgvpnreplicatedtopics) | **Get** /msgVpns/{msgVpnName}/replicatedTopics | Get a list of Replicated Topic objects.
*MsgVpnApi* | [**GetMsgVpnRestDeliveryPoint**](docs/MsgVpnApi.md#getmsgvpnrestdeliverypoint) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName} | Get a REST Delivery Point object.
*MsgVpnApi* | [**GetMsgVpnRestDeliveryPointQueueBinding**](docs/MsgVpnApi.md#getmsgvpnrestdeliverypointqueuebinding) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName} | Get a Queue Binding object.
*MsgVpnApi* | [**GetMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/MsgVpnApi.md#getmsgvpnrestdeliverypointqueuebindingrequestheader) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders/{headerName} | Get a Request Header object.
*MsgVpnApi* | [**GetMsgVpnRestDeliveryPointQueueBindingRequestHeaders**](docs/MsgVpnApi.md#getmsgvpnrestdeliverypointqueuebindingrequestheaders) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders | Get a list of Request Header objects.
*MsgVpnApi* | [**GetMsgVpnRestDeliveryPointQueueBindings**](docs/MsgVpnApi.md#getmsgvpnrestdeliverypointqueuebindings) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings | Get a list of Queue Binding objects.
*MsgVpnApi* | [**GetMsgVpnRestDeliveryPointRestConsumer**](docs/MsgVpnApi.md#getmsgvpnrestdeliverypointrestconsumer) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName} | Get a REST Consumer object.
*MsgVpnApi* | [**GetMsgVpnRestDeliveryPointRestConsumerOauthJwtClaim**](docs/MsgVpnApi.md#getmsgvpnrestdeliverypointrestconsumeroauthjwtclaim) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/oauthJwtClaims/{oauthJwtClaimName} | Get a Claim object.
*MsgVpnApi* | [**GetMsgVpnRestDeliveryPointRestConsumerOauthJwtClaims**](docs/MsgVpnApi.md#getmsgvpnrestdeliverypointrestconsumeroauthjwtclaims) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/oauthJwtClaims | Get a list of Claim objects.
*MsgVpnApi* | [**GetMsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonName**](docs/MsgVpnApi.md#getmsgvpnrestdeliverypointrestconsumertlstrustedcommonname) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Get a Trusted Common Name object.
*MsgVpnApi* | [**GetMsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonNames**](docs/MsgVpnApi.md#getmsgvpnrestdeliverypointrestconsumertlstrustedcommonnames) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/tlsTrustedCommonNames | Get a list of Trusted Common Name objects.
*MsgVpnApi* | [**GetMsgVpnRestDeliveryPointRestConsumers**](docs/MsgVpnApi.md#getmsgvpnrestdeliverypointrestconsumers) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers | Get a list of REST Consumer objects.
*MsgVpnApi* | [**GetMsgVpnRestDeliveryPoints**](docs/MsgVpnApi.md#getmsgvpnrestdeliverypoints) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints | Get a list of REST Delivery Point objects.
*MsgVpnApi* | [**GetMsgVpnSequencedTopic**](docs/MsgVpnApi.md#getmsgvpnsequencedtopic) | **Get** /msgVpns/{msgVpnName}/sequencedTopics/{sequencedTopic} | Get a Sequenced Topic object.
*MsgVpnApi* | [**GetMsgVpnSequencedTopics**](docs/MsgVpnApi.md#getmsgvpnsequencedtopics) | **Get** /msgVpns/{msgVpnName}/sequencedTopics | Get a list of Sequenced Topic objects.
*MsgVpnApi* | [**GetMsgVpnTopicEndpoint**](docs/MsgVpnApi.md#getmsgvpntopicendpoint) | **Get** /msgVpns/{msgVpnName}/topicEndpoints/{topicEndpointName} | Get a Topic Endpoint object.
*MsgVpnApi* | [**GetMsgVpnTopicEndpointTemplate**](docs/MsgVpnApi.md#getmsgvpntopicendpointtemplate) | **Get** /msgVpns/{msgVpnName}/topicEndpointTemplates/{topicEndpointTemplateName} | Get a Topic Endpoint Template object.
*MsgVpnApi* | [**GetMsgVpnTopicEndpointTemplates**](docs/MsgVpnApi.md#getmsgvpntopicendpointtemplates) | **Get** /msgVpns/{msgVpnName}/topicEndpointTemplates | Get a list of Topic Endpoint Template objects.
*MsgVpnApi* | [**GetMsgVpnTopicEndpoints**](docs/MsgVpnApi.md#getmsgvpntopicendpoints) | **Get** /msgVpns/{msgVpnName}/topicEndpoints | Get a list of Topic Endpoint objects.
*MsgVpnApi* | [**GetMsgVpns**](docs/MsgVpnApi.md#getmsgvpns) | **Get** /msgVpns | Get a list of Message VPN objects.
*MsgVpnApi* | [**ReplaceMsgVpn**](docs/MsgVpnApi.md#replacemsgvpn) | **Put** /msgVpns/{msgVpnName} | Replace a Message VPN object.
*MsgVpnApi* | [**ReplaceMsgVpnAclProfile**](docs/MsgVpnApi.md#replacemsgvpnaclprofile) | **Put** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName} | Replace an ACL Profile object.
*MsgVpnApi* | [**ReplaceMsgVpnAuthenticationOauthProfile**](docs/MsgVpnApi.md#replacemsgvpnauthenticationoauthprofile) | **Put** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName} | Replace an OAuth Profile object.
*MsgVpnApi* | [**ReplaceMsgVpnAuthenticationOauthProvider**](docs/MsgVpnApi.md#replacemsgvpnauthenticationoauthprovider) | **Put** /msgVpns/{msgVpnName}/authenticationOauthProviders/{oauthProviderName} | Replace an OAuth Provider object.
*MsgVpnApi* | [**ReplaceMsgVpnAuthorizationGroup**](docs/MsgVpnApi.md#replacemsgvpnauthorizationgroup) | **Put** /msgVpns/{msgVpnName}/authorizationGroups/{authorizationGroupName} | Replace an Authorization Group object.
*MsgVpnApi* | [**ReplaceMsgVpnBridge**](docs/MsgVpnApi.md#replacemsgvpnbridge) | **Put** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter} | Replace a Bridge object.
*MsgVpnApi* | [**ReplaceMsgVpnBridgeRemoteMsgVpn**](docs/MsgVpnApi.md#replacemsgvpnbridgeremotemsgvpn) | **Put** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface} | Replace a Remote Message VPN object.
*MsgVpnApi* | [**ReplaceMsgVpnClientProfile**](docs/MsgVpnApi.md#replacemsgvpnclientprofile) | **Put** /msgVpns/{msgVpnName}/clientProfiles/{clientProfileName} | Replace a Client Profile object.
*MsgVpnApi* | [**ReplaceMsgVpnClientUsername**](docs/MsgVpnApi.md#replacemsgvpnclientusername) | **Put** /msgVpns/{msgVpnName}/clientUsernames/{clientUsername} | Replace a Client Username object.
*MsgVpnApi* | [**ReplaceMsgVpnDistributedCache**](docs/MsgVpnApi.md#replacemsgvpndistributedcache) | **Put** /msgVpns/{msgVpnName}/distributedCaches/{cacheName} | Replace a Distributed Cache object.
*MsgVpnApi* | [**ReplaceMsgVpnDistributedCacheCluster**](docs/MsgVpnApi.md#replacemsgvpndistributedcachecluster) | **Put** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName} | Replace a Cache Cluster object.
*MsgVpnApi* | [**ReplaceMsgVpnDistributedCacheClusterInstance**](docs/MsgVpnApi.md#replacemsgvpndistributedcacheclusterinstance) | **Put** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName} | Replace a Cache Instance object.
*MsgVpnApi* | [**ReplaceMsgVpnDmrBridge**](docs/MsgVpnApi.md#replacemsgvpndmrbridge) | **Put** /msgVpns/{msgVpnName}/dmrBridges/{remoteNodeName} | Replace a DMR Bridge object.
*MsgVpnApi* | [**ReplaceMsgVpnJndiConnectionFactory**](docs/MsgVpnApi.md#replacemsgvpnjndiconnectionfactory) | **Put** /msgVpns/{msgVpnName}/jndiConnectionFactories/{connectionFactoryName} | Replace a JNDI Connection Factory object.
*MsgVpnApi* | [**ReplaceMsgVpnJndiQueue**](docs/MsgVpnApi.md#replacemsgvpnjndiqueue) | **Put** /msgVpns/{msgVpnName}/jndiQueues/{queueName} | Replace a JNDI Queue object.
*MsgVpnApi* | [**ReplaceMsgVpnJndiTopic**](docs/MsgVpnApi.md#replacemsgvpnjnditopic) | **Put** /msgVpns/{msgVpnName}/jndiTopics/{topicName} | Replace a JNDI Topic object.
*MsgVpnApi* | [**ReplaceMsgVpnMqttRetainCache**](docs/MsgVpnApi.md#replacemsgvpnmqttretaincache) | **Put** /msgVpns/{msgVpnName}/mqttRetainCaches/{cacheName} | Replace an MQTT Retain Cache object.
*MsgVpnApi* | [**ReplaceMsgVpnMqttSession**](docs/MsgVpnApi.md#replacemsgvpnmqttsession) | **Put** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter} | Replace an MQTT Session object.
*MsgVpnApi* | [**ReplaceMsgVpnMqttSessionSubscription**](docs/MsgVpnApi.md#replacemsgvpnmqttsessionsubscription) | **Put** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic} | Replace a Subscription object.
*MsgVpnApi* | [**ReplaceMsgVpnQueue**](docs/MsgVpnApi.md#replacemsgvpnqueue) | **Put** /msgVpns/{msgVpnName}/queues/{queueName} | Replace a Queue object.
*MsgVpnApi* | [**ReplaceMsgVpnQueueTemplate**](docs/MsgVpnApi.md#replacemsgvpnqueuetemplate) | **Put** /msgVpns/{msgVpnName}/queueTemplates/{queueTemplateName} | Replace a Queue Template object.
*MsgVpnApi* | [**ReplaceMsgVpnReplayLog**](docs/MsgVpnApi.md#replacemsgvpnreplaylog) | **Put** /msgVpns/{msgVpnName}/replayLogs/{replayLogName} | Replace a Replay Log object.
*MsgVpnApi* | [**ReplaceMsgVpnReplicatedTopic**](docs/MsgVpnApi.md#replacemsgvpnreplicatedtopic) | **Put** /msgVpns/{msgVpnName}/replicatedTopics/{replicatedTopic} | Replace a Replicated Topic object.
*MsgVpnApi* | [**ReplaceMsgVpnRestDeliveryPoint**](docs/MsgVpnApi.md#replacemsgvpnrestdeliverypoint) | **Put** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName} | Replace a REST Delivery Point object.
*MsgVpnApi* | [**ReplaceMsgVpnRestDeliveryPointQueueBinding**](docs/MsgVpnApi.md#replacemsgvpnrestdeliverypointqueuebinding) | **Put** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName} | Replace a Queue Binding object.
*MsgVpnApi* | [**ReplaceMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/MsgVpnApi.md#replacemsgvpnrestdeliverypointqueuebindingrequestheader) | **Put** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders/{headerName} | Replace a Request Header object.
*MsgVpnApi* | [**ReplaceMsgVpnRestDeliveryPointRestConsumer**](docs/MsgVpnApi.md#replacemsgvpnrestdeliverypointrestconsumer) | **Put** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName} | Replace a REST Consumer object.
*MsgVpnApi* | [**ReplaceMsgVpnTopicEndpoint**](docs/MsgVpnApi.md#replacemsgvpntopicendpoint) | **Put** /msgVpns/{msgVpnName}/topicEndpoints/{topicEndpointName} | Replace a Topic Endpoint object.
*MsgVpnApi* | [**ReplaceMsgVpnTopicEndpointTemplate**](docs/MsgVpnApi.md#replacemsgvpntopicendpointtemplate) | **Put** /msgVpns/{msgVpnName}/topicEndpointTemplates/{topicEndpointTemplateName} | Replace a Topic Endpoint Template object.
*MsgVpnApi* | [**UpdateMsgVpn**](docs/MsgVpnApi.md#updatemsgvpn) | **Patch** /msgVpns/{msgVpnName} | Update a Message VPN object.
*MsgVpnApi* | [**UpdateMsgVpnAclProfile**](docs/MsgVpnApi.md#updatemsgvpnaclprofile) | **Patch** /msgVpns/{msgVpnName}/aclProfiles/{aclProfileName} | Update an ACL Profile object.
*MsgVpnApi* | [**UpdateMsgVpnAuthenticationOauthProfile**](docs/MsgVpnApi.md#updatemsgvpnauthenticationoauthprofile) | **Patch** /msgVpns/{msgVpnName}/authenticationOauthProfiles/{oauthProfileName} | Update an OAuth Profile object.
*MsgVpnApi* | [**UpdateMsgVpnAuthenticationOauthProvider**](docs/MsgVpnApi.md#updatemsgvpnauthenticationoauthprovider) | **Patch** /msgVpns/{msgVpnName}/authenticationOauthProviders/{oauthProviderName} | Update an OAuth Provider object.
*MsgVpnApi* | [**UpdateMsgVpnAuthorizationGroup**](docs/MsgVpnApi.md#updatemsgvpnauthorizationgroup) | **Patch** /msgVpns/{msgVpnName}/authorizationGroups/{authorizationGroupName} | Update an Authorization Group object.
*MsgVpnApi* | [**UpdateMsgVpnBridge**](docs/MsgVpnApi.md#updatemsgvpnbridge) | **Patch** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter} | Update a Bridge object.
*MsgVpnApi* | [**UpdateMsgVpnBridgeRemoteMsgVpn**](docs/MsgVpnApi.md#updatemsgvpnbridgeremotemsgvpn) | **Patch** /msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface} | Update a Remote Message VPN object.
*MsgVpnApi* | [**UpdateMsgVpnClientProfile**](docs/MsgVpnApi.md#updatemsgvpnclientprofile) | **Patch** /msgVpns/{msgVpnName}/clientProfiles/{clientProfileName} | Update a Client Profile object.
*MsgVpnApi* | [**UpdateMsgVpnClientUsername**](docs/MsgVpnApi.md#updatemsgvpnclientusername) | **Patch** /msgVpns/{msgVpnName}/clientUsernames/{clientUsername} | Update a Client Username object.
*MsgVpnApi* | [**UpdateMsgVpnDistributedCache**](docs/MsgVpnApi.md#updatemsgvpndistributedcache) | **Patch** /msgVpns/{msgVpnName}/distributedCaches/{cacheName} | Update a Distributed Cache object.
*MsgVpnApi* | [**UpdateMsgVpnDistributedCacheCluster**](docs/MsgVpnApi.md#updatemsgvpndistributedcachecluster) | **Patch** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName} | Update a Cache Cluster object.
*MsgVpnApi* | [**UpdateMsgVpnDistributedCacheClusterInstance**](docs/MsgVpnApi.md#updatemsgvpndistributedcacheclusterinstance) | **Patch** /msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName} | Update a Cache Instance object.
*MsgVpnApi* | [**UpdateMsgVpnDmrBridge**](docs/MsgVpnApi.md#updatemsgvpndmrbridge) | **Patch** /msgVpns/{msgVpnName}/dmrBridges/{remoteNodeName} | Update a DMR Bridge object.
*MsgVpnApi* | [**UpdateMsgVpnJndiConnectionFactory**](docs/MsgVpnApi.md#updatemsgvpnjndiconnectionfactory) | **Patch** /msgVpns/{msgVpnName}/jndiConnectionFactories/{connectionFactoryName} | Update a JNDI Connection Factory object.
*MsgVpnApi* | [**UpdateMsgVpnJndiQueue**](docs/MsgVpnApi.md#updatemsgvpnjndiqueue) | **Patch** /msgVpns/{msgVpnName}/jndiQueues/{queueName} | Update a JNDI Queue object.
*MsgVpnApi* | [**UpdateMsgVpnJndiTopic**](docs/MsgVpnApi.md#updatemsgvpnjnditopic) | **Patch** /msgVpns/{msgVpnName}/jndiTopics/{topicName} | Update a JNDI Topic object.
*MsgVpnApi* | [**UpdateMsgVpnMqttRetainCache**](docs/MsgVpnApi.md#updatemsgvpnmqttretaincache) | **Patch** /msgVpns/{msgVpnName}/mqttRetainCaches/{cacheName} | Update an MQTT Retain Cache object.
*MsgVpnApi* | [**UpdateMsgVpnMqttSession**](docs/MsgVpnApi.md#updatemsgvpnmqttsession) | **Patch** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter} | Update an MQTT Session object.
*MsgVpnApi* | [**UpdateMsgVpnMqttSessionSubscription**](docs/MsgVpnApi.md#updatemsgvpnmqttsessionsubscription) | **Patch** /msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic} | Update a Subscription object.
*MsgVpnApi* | [**UpdateMsgVpnQueue**](docs/MsgVpnApi.md#updatemsgvpnqueue) | **Patch** /msgVpns/{msgVpnName}/queues/{queueName} | Update a Queue object.
*MsgVpnApi* | [**UpdateMsgVpnQueueTemplate**](docs/MsgVpnApi.md#updatemsgvpnqueuetemplate) | **Patch** /msgVpns/{msgVpnName}/queueTemplates/{queueTemplateName} | Update a Queue Template object.
*MsgVpnApi* | [**UpdateMsgVpnReplayLog**](docs/MsgVpnApi.md#updatemsgvpnreplaylog) | **Patch** /msgVpns/{msgVpnName}/replayLogs/{replayLogName} | Update a Replay Log object.
*MsgVpnApi* | [**UpdateMsgVpnReplicatedTopic**](docs/MsgVpnApi.md#updatemsgvpnreplicatedtopic) | **Patch** /msgVpns/{msgVpnName}/replicatedTopics/{replicatedTopic} | Update a Replicated Topic object.
*MsgVpnApi* | [**UpdateMsgVpnRestDeliveryPoint**](docs/MsgVpnApi.md#updatemsgvpnrestdeliverypoint) | **Patch** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName} | Update a REST Delivery Point object.
*MsgVpnApi* | [**UpdateMsgVpnRestDeliveryPointQueueBinding**](docs/MsgVpnApi.md#updatemsgvpnrestdeliverypointqueuebinding) | **Patch** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName} | Update a Queue Binding object.
*MsgVpnApi* | [**UpdateMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/MsgVpnApi.md#updatemsgvpnrestdeliverypointqueuebindingrequestheader) | **Patch** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders/{headerName} | Update a Request Header object.
*MsgVpnApi* | [**UpdateMsgVpnRestDeliveryPointRestConsumer**](docs/MsgVpnApi.md#updatemsgvpnrestdeliverypointrestconsumer) | **Patch** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName} | Update a REST Consumer object.
*MsgVpnApi* | [**UpdateMsgVpnTopicEndpoint**](docs/MsgVpnApi.md#updatemsgvpntopicendpoint) | **Patch** /msgVpns/{msgVpnName}/topicEndpoints/{topicEndpointName} | Update a Topic Endpoint object.
*MsgVpnApi* | [**UpdateMsgVpnTopicEndpointTemplate**](docs/MsgVpnApi.md#updatemsgvpntopicendpointtemplate) | **Patch** /msgVpns/{msgVpnName}/topicEndpointTemplates/{topicEndpointTemplateName} | Update a Topic Endpoint Template object.
*OauthProfileApi* | [**CreateOauthProfile**](docs/OauthProfileApi.md#createoauthprofile) | **Post** /oauthProfiles | Create an OAuth Profile object.
*OauthProfileApi* | [**CreateOauthProfileAccessLevelGroup**](docs/OauthProfileApi.md#createoauthprofileaccesslevelgroup) | **Post** /oauthProfiles/{oauthProfileName}/accessLevelGroups | Create a Group Access Level object.
*OauthProfileApi* | [**CreateOauthProfileAccessLevelGroupMsgVpnAccessLevelException**](docs/OauthProfileApi.md#createoauthprofileaccesslevelgroupmsgvpnaccesslevelexception) | **Post** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions | Create a Message VPN Access-Level Exception object.
*OauthProfileApi* | [**CreateOauthProfileClientAllowedHost**](docs/OauthProfileApi.md#createoauthprofileclientallowedhost) | **Post** /oauthProfiles/{oauthProfileName}/clientAllowedHosts | Create an Allowed Host Value object.
*OauthProfileApi* | [**CreateOauthProfileClientAuthorizationParameter**](docs/OauthProfileApi.md#createoauthprofileclientauthorizationparameter) | **Post** /oauthProfiles/{oauthProfileName}/clientAuthorizationParameters | Create an Authorization Parameter object.
*OauthProfileApi* | [**CreateOauthProfileClientRequiredClaim**](docs/OauthProfileApi.md#createoauthprofileclientrequiredclaim) | **Post** /oauthProfiles/{oauthProfileName}/clientRequiredClaims | Create a Required Claim object.
*OauthProfileApi* | [**CreateOauthProfileDefaultMsgVpnAccessLevelException**](docs/OauthProfileApi.md#createoauthprofiledefaultmsgvpnaccesslevelexception) | **Post** /oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions | Create a Message VPN Access-Level Exception object.
*OauthProfileApi* | [**CreateOauthProfileResourceServerRequiredClaim**](docs/OauthProfileApi.md#createoauthprofileresourceserverrequiredclaim) | **Post** /oauthProfiles/{oauthProfileName}/resourceServerRequiredClaims | Create a Required Claim object.
*OauthProfileApi* | [**DeleteOauthProfile**](docs/OauthProfileApi.md#deleteoauthprofile) | **Delete** /oauthProfiles/{oauthProfileName} | Delete an OAuth Profile object.
*OauthProfileApi* | [**DeleteOauthProfileAccessLevelGroup**](docs/OauthProfileApi.md#deleteoauthprofileaccesslevelgroup) | **Delete** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName} | Delete a Group Access Level object.
*OauthProfileApi* | [**DeleteOauthProfileAccessLevelGroupMsgVpnAccessLevelException**](docs/OauthProfileApi.md#deleteoauthprofileaccesslevelgroupmsgvpnaccesslevelexception) | **Delete** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions/{msgVpnName} | Delete a Message VPN Access-Level Exception object.
*OauthProfileApi* | [**DeleteOauthProfileClientAllowedHost**](docs/OauthProfileApi.md#deleteoauthprofileclientallowedhost) | **Delete** /oauthProfiles/{oauthProfileName}/clientAllowedHosts/{allowedHost} | Delete an Allowed Host Value object.
*OauthProfileApi* | [**DeleteOauthProfileClientAuthorizationParameter**](docs/OauthProfileApi.md#deleteoauthprofileclientauthorizationparameter) | **Delete** /oauthProfiles/{oauthProfileName}/clientAuthorizationParameters/{authorizationParameterName} | Delete an Authorization Parameter object.
*OauthProfileApi* | [**DeleteOauthProfileClientRequiredClaim**](docs/OauthProfileApi.md#deleteoauthprofileclientrequiredclaim) | **Delete** /oauthProfiles/{oauthProfileName}/clientRequiredClaims/{clientRequiredClaimName} | Delete a Required Claim object.
*OauthProfileApi* | [**DeleteOauthProfileDefaultMsgVpnAccessLevelException**](docs/OauthProfileApi.md#deleteoauthprofiledefaultmsgvpnaccesslevelexception) | **Delete** /oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions/{msgVpnName} | Delete a Message VPN Access-Level Exception object.
*OauthProfileApi* | [**DeleteOauthProfileResourceServerRequiredClaim**](docs/OauthProfileApi.md#deleteoauthprofileresourceserverrequiredclaim) | **Delete** /oauthProfiles/{oauthProfileName}/resourceServerRequiredClaims/{resourceServerRequiredClaimName} | Delete a Required Claim object.
*OauthProfileApi* | [**GetOauthProfile**](docs/OauthProfileApi.md#getoauthprofile) | **Get** /oauthProfiles/{oauthProfileName} | Get an OAuth Profile object.
*OauthProfileApi* | [**GetOauthProfileAccessLevelGroup**](docs/OauthProfileApi.md#getoauthprofileaccesslevelgroup) | **Get** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName} | Get a Group Access Level object.
*OauthProfileApi* | [**GetOauthProfileAccessLevelGroupMsgVpnAccessLevelException**](docs/OauthProfileApi.md#getoauthprofileaccesslevelgroupmsgvpnaccesslevelexception) | **Get** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions/{msgVpnName} | Get a Message VPN Access-Level Exception object.
*OauthProfileApi* | [**GetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptions**](docs/OauthProfileApi.md#getoauthprofileaccesslevelgroupmsgvpnaccesslevelexceptions) | **Get** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions | Get a list of Message VPN Access-Level Exception objects.
*OauthProfileApi* | [**GetOauthProfileAccessLevelGroups**](docs/OauthProfileApi.md#getoauthprofileaccesslevelgroups) | **Get** /oauthProfiles/{oauthProfileName}/accessLevelGroups | Get a list of Group Access Level objects.
*OauthProfileApi* | [**GetOauthProfileClientAllowedHost**](docs/OauthProfileApi.md#getoauthprofileclientallowedhost) | **Get** /oauthProfiles/{oauthProfileName}/clientAllowedHosts/{allowedHost} | Get an Allowed Host Value object.
*OauthProfileApi* | [**GetOauthProfileClientAllowedHosts**](docs/OauthProfileApi.md#getoauthprofileclientallowedhosts) | **Get** /oauthProfiles/{oauthProfileName}/clientAllowedHosts | Get a list of Allowed Host Value objects.
*OauthProfileApi* | [**GetOauthProfileClientAuthorizationParameter**](docs/OauthProfileApi.md#getoauthprofileclientauthorizationparameter) | **Get** /oauthProfiles/{oauthProfileName}/clientAuthorizationParameters/{authorizationParameterName} | Get an Authorization Parameter object.
*OauthProfileApi* | [**GetOauthProfileClientAuthorizationParameters**](docs/OauthProfileApi.md#getoauthprofileclientauthorizationparameters) | **Get** /oauthProfiles/{oauthProfileName}/clientAuthorizationParameters | Get a list of Authorization Parameter objects.
*OauthProfileApi* | [**GetOauthProfileClientRequiredClaim**](docs/OauthProfileApi.md#getoauthprofileclientrequiredclaim) | **Get** /oauthProfiles/{oauthProfileName}/clientRequiredClaims/{clientRequiredClaimName} | Get a Required Claim object.
*OauthProfileApi* | [**GetOauthProfileClientRequiredClaims**](docs/OauthProfileApi.md#getoauthprofileclientrequiredclaims) | **Get** /oauthProfiles/{oauthProfileName}/clientRequiredClaims | Get a list of Required Claim objects.
*OauthProfileApi* | [**GetOauthProfileDefaultMsgVpnAccessLevelException**](docs/OauthProfileApi.md#getoauthprofiledefaultmsgvpnaccesslevelexception) | **Get** /oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions/{msgVpnName} | Get a Message VPN Access-Level Exception object.
*OauthProfileApi* | [**GetOauthProfileDefaultMsgVpnAccessLevelExceptions**](docs/OauthProfileApi.md#getoauthprofiledefaultmsgvpnaccesslevelexceptions) | **Get** /oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions | Get a list of Message VPN Access-Level Exception objects.
*OauthProfileApi* | [**GetOauthProfileResourceServerRequiredClaim**](docs/OauthProfileApi.md#getoauthprofileresourceserverrequiredclaim) | **Get** /oauthProfiles/{oauthProfileName}/resourceServerRequiredClaims/{resourceServerRequiredClaimName} | Get a Required Claim object.
*OauthProfileApi* | [**GetOauthProfileResourceServerRequiredClaims**](docs/OauthProfileApi.md#getoauthprofileresourceserverrequiredclaims) | **Get** /oauthProfiles/{oauthProfileName}/resourceServerRequiredClaims | Get a list of Required Claim objects.
*OauthProfileApi* | [**GetOauthProfiles**](docs/OauthProfileApi.md#getoauthprofiles) | **Get** /oauthProfiles | Get a list of OAuth Profile objects.
*OauthProfileApi* | [**ReplaceOauthProfile**](docs/OauthProfileApi.md#replaceoauthprofile) | **Put** /oauthProfiles/{oauthProfileName} | Replace an OAuth Profile object.
*OauthProfileApi* | [**ReplaceOauthProfileAccessLevelGroup**](docs/OauthProfileApi.md#replaceoauthprofileaccesslevelgroup) | **Put** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName} | Replace a Group Access Level object.
*OauthProfileApi* | [**ReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelException**](docs/OauthProfileApi.md#replaceoauthprofileaccesslevelgroupmsgvpnaccesslevelexception) | **Put** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions/{msgVpnName} | Replace a Message VPN Access-Level Exception object.
*OauthProfileApi* | [**ReplaceOauthProfileClientAuthorizationParameter**](docs/OauthProfileApi.md#replaceoauthprofileclientauthorizationparameter) | **Put** /oauthProfiles/{oauthProfileName}/clientAuthorizationParameters/{authorizationParameterName} | Replace an Authorization Parameter object.
*OauthProfileApi* | [**ReplaceOauthProfileDefaultMsgVpnAccessLevelException**](docs/OauthProfileApi.md#replaceoauthprofiledefaultmsgvpnaccesslevelexception) | **Put** /oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions/{msgVpnName} | Replace a Message VPN Access-Level Exception object.
*OauthProfileApi* | [**UpdateOauthProfile**](docs/OauthProfileApi.md#updateoauthprofile) | **Patch** /oauthProfiles/{oauthProfileName} | Update an OAuth Profile object.
*OauthProfileApi* | [**UpdateOauthProfileAccessLevelGroup**](docs/OauthProfileApi.md#updateoauthprofileaccesslevelgroup) | **Patch** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName} | Update a Group Access Level object.
*OauthProfileApi* | [**UpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelException**](docs/OauthProfileApi.md#updateoauthprofileaccesslevelgroupmsgvpnaccesslevelexception) | **Patch** /oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions/{msgVpnName} | Update a Message VPN Access-Level Exception object.
*OauthProfileApi* | [**UpdateOauthProfileClientAuthorizationParameter**](docs/OauthProfileApi.md#updateoauthprofileclientauthorizationparameter) | **Patch** /oauthProfiles/{oauthProfileName}/clientAuthorizationParameters/{authorizationParameterName} | Update an Authorization Parameter object.
*OauthProfileApi* | [**UpdateOauthProfileDefaultMsgVpnAccessLevelException**](docs/OauthProfileApi.md#updateoauthprofiledefaultmsgvpnaccesslevelexception) | **Patch** /oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions/{msgVpnName} | Update a Message VPN Access-Level Exception object.
*QueueApi* | [**CreateMsgVpnQueue**](docs/QueueApi.md#createmsgvpnqueue) | **Post** /msgVpns/{msgVpnName}/queues | Create a Queue object.
*QueueApi* | [**CreateMsgVpnQueueSubscription**](docs/QueueApi.md#createmsgvpnqueuesubscription) | **Post** /msgVpns/{msgVpnName}/queues/{queueName}/subscriptions | Create a Queue Subscription object.
*QueueApi* | [**DeleteMsgVpnQueue**](docs/QueueApi.md#deletemsgvpnqueue) | **Delete** /msgVpns/{msgVpnName}/queues/{queueName} | Delete a Queue object.
*QueueApi* | [**DeleteMsgVpnQueueSubscription**](docs/QueueApi.md#deletemsgvpnqueuesubscription) | **Delete** /msgVpns/{msgVpnName}/queues/{queueName}/subscriptions/{subscriptionTopic} | Delete a Queue Subscription object.
*QueueApi* | [**GetMsgVpnQueue**](docs/QueueApi.md#getmsgvpnqueue) | **Get** /msgVpns/{msgVpnName}/queues/{queueName} | Get a Queue object.
*QueueApi* | [**GetMsgVpnQueueSubscription**](docs/QueueApi.md#getmsgvpnqueuesubscription) | **Get** /msgVpns/{msgVpnName}/queues/{queueName}/subscriptions/{subscriptionTopic} | Get a Queue Subscription object.
*QueueApi* | [**GetMsgVpnQueueSubscriptions**](docs/QueueApi.md#getmsgvpnqueuesubscriptions) | **Get** /msgVpns/{msgVpnName}/queues/{queueName}/subscriptions | Get a list of Queue Subscription objects.
*QueueApi* | [**GetMsgVpnQueues**](docs/QueueApi.md#getmsgvpnqueues) | **Get** /msgVpns/{msgVpnName}/queues | Get a list of Queue objects.
*QueueApi* | [**ReplaceMsgVpnQueue**](docs/QueueApi.md#replacemsgvpnqueue) | **Put** /msgVpns/{msgVpnName}/queues/{queueName} | Replace a Queue object.
*QueueApi* | [**UpdateMsgVpnQueue**](docs/QueueApi.md#updatemsgvpnqueue) | **Patch** /msgVpns/{msgVpnName}/queues/{queueName} | Update a Queue object.
*QueueTemplateApi* | [**CreateMsgVpnQueueTemplate**](docs/QueueTemplateApi.md#createmsgvpnqueuetemplate) | **Post** /msgVpns/{msgVpnName}/queueTemplates | Create a Queue Template object.
*QueueTemplateApi* | [**DeleteMsgVpnQueueTemplate**](docs/QueueTemplateApi.md#deletemsgvpnqueuetemplate) | **Delete** /msgVpns/{msgVpnName}/queueTemplates/{queueTemplateName} | Delete a Queue Template object.
*QueueTemplateApi* | [**GetMsgVpnQueueTemplate**](docs/QueueTemplateApi.md#getmsgvpnqueuetemplate) | **Get** /msgVpns/{msgVpnName}/queueTemplates/{queueTemplateName} | Get a Queue Template object.
*QueueTemplateApi* | [**GetMsgVpnQueueTemplates**](docs/QueueTemplateApi.md#getmsgvpnqueuetemplates) | **Get** /msgVpns/{msgVpnName}/queueTemplates | Get a list of Queue Template objects.
*QueueTemplateApi* | [**ReplaceMsgVpnQueueTemplate**](docs/QueueTemplateApi.md#replacemsgvpnqueuetemplate) | **Put** /msgVpns/{msgVpnName}/queueTemplates/{queueTemplateName} | Replace a Queue Template object.
*QueueTemplateApi* | [**UpdateMsgVpnQueueTemplate**](docs/QueueTemplateApi.md#updatemsgvpnqueuetemplate) | **Patch** /msgVpns/{msgVpnName}/queueTemplates/{queueTemplateName} | Update a Queue Template object.
*ReplayLogApi* | [**CreateMsgVpnReplayLog**](docs/ReplayLogApi.md#createmsgvpnreplaylog) | **Post** /msgVpns/{msgVpnName}/replayLogs | Create a Replay Log object.
*ReplayLogApi* | [**DeleteMsgVpnReplayLog**](docs/ReplayLogApi.md#deletemsgvpnreplaylog) | **Delete** /msgVpns/{msgVpnName}/replayLogs/{replayLogName} | Delete a Replay Log object.
*ReplayLogApi* | [**GetMsgVpnReplayLog**](docs/ReplayLogApi.md#getmsgvpnreplaylog) | **Get** /msgVpns/{msgVpnName}/replayLogs/{replayLogName} | Get a Replay Log object.
*ReplayLogApi* | [**GetMsgVpnReplayLogs**](docs/ReplayLogApi.md#getmsgvpnreplaylogs) | **Get** /msgVpns/{msgVpnName}/replayLogs | Get a list of Replay Log objects.
*ReplayLogApi* | [**ReplaceMsgVpnReplayLog**](docs/ReplayLogApi.md#replacemsgvpnreplaylog) | **Put** /msgVpns/{msgVpnName}/replayLogs/{replayLogName} | Replace a Replay Log object.
*ReplayLogApi* | [**UpdateMsgVpnReplayLog**](docs/ReplayLogApi.md#updatemsgvpnreplaylog) | **Patch** /msgVpns/{msgVpnName}/replayLogs/{replayLogName} | Update a Replay Log object.
*ReplicatedTopicApi* | [**CreateMsgVpnReplicatedTopic**](docs/ReplicatedTopicApi.md#createmsgvpnreplicatedtopic) | **Post** /msgVpns/{msgVpnName}/replicatedTopics | Create a Replicated Topic object.
*ReplicatedTopicApi* | [**DeleteMsgVpnReplicatedTopic**](docs/ReplicatedTopicApi.md#deletemsgvpnreplicatedtopic) | **Delete** /msgVpns/{msgVpnName}/replicatedTopics/{replicatedTopic} | Delete a Replicated Topic object.
*ReplicatedTopicApi* | [**GetMsgVpnReplicatedTopic**](docs/ReplicatedTopicApi.md#getmsgvpnreplicatedtopic) | **Get** /msgVpns/{msgVpnName}/replicatedTopics/{replicatedTopic} | Get a Replicated Topic object.
*ReplicatedTopicApi* | [**GetMsgVpnReplicatedTopics**](docs/ReplicatedTopicApi.md#getmsgvpnreplicatedtopics) | **Get** /msgVpns/{msgVpnName}/replicatedTopics | Get a list of Replicated Topic objects.
*ReplicatedTopicApi* | [**ReplaceMsgVpnReplicatedTopic**](docs/ReplicatedTopicApi.md#replacemsgvpnreplicatedtopic) | **Put** /msgVpns/{msgVpnName}/replicatedTopics/{replicatedTopic} | Replace a Replicated Topic object.
*ReplicatedTopicApi* | [**UpdateMsgVpnReplicatedTopic**](docs/ReplicatedTopicApi.md#updatemsgvpnreplicatedtopic) | **Patch** /msgVpns/{msgVpnName}/replicatedTopics/{replicatedTopic} | Update a Replicated Topic object.
*RestDeliveryPointApi* | [**CreateMsgVpnRestDeliveryPoint**](docs/RestDeliveryPointApi.md#createmsgvpnrestdeliverypoint) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints | Create a REST Delivery Point object.
*RestDeliveryPointApi* | [**CreateMsgVpnRestDeliveryPointQueueBinding**](docs/RestDeliveryPointApi.md#createmsgvpnrestdeliverypointqueuebinding) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings | Create a Queue Binding object.
*RestDeliveryPointApi* | [**CreateMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/RestDeliveryPointApi.md#createmsgvpnrestdeliverypointqueuebindingrequestheader) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders | Create a Request Header object.
*RestDeliveryPointApi* | [**CreateMsgVpnRestDeliveryPointRestConsumer**](docs/RestDeliveryPointApi.md#createmsgvpnrestdeliverypointrestconsumer) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers | Create a REST Consumer object.
*RestDeliveryPointApi* | [**CreateMsgVpnRestDeliveryPointRestConsumerOauthJwtClaim**](docs/RestDeliveryPointApi.md#createmsgvpnrestdeliverypointrestconsumeroauthjwtclaim) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/oauthJwtClaims | Create a Claim object.
*RestDeliveryPointApi* | [**CreateMsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonName**](docs/RestDeliveryPointApi.md#createmsgvpnrestdeliverypointrestconsumertlstrustedcommonname) | **Post** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/tlsTrustedCommonNames | Create a Trusted Common Name object.
*RestDeliveryPointApi* | [**DeleteMsgVpnRestDeliveryPoint**](docs/RestDeliveryPointApi.md#deletemsgvpnrestdeliverypoint) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName} | Delete a REST Delivery Point object.
*RestDeliveryPointApi* | [**DeleteMsgVpnRestDeliveryPointQueueBinding**](docs/RestDeliveryPointApi.md#deletemsgvpnrestdeliverypointqueuebinding) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName} | Delete a Queue Binding object.
*RestDeliveryPointApi* | [**DeleteMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/RestDeliveryPointApi.md#deletemsgvpnrestdeliverypointqueuebindingrequestheader) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders/{headerName} | Delete a Request Header object.
*RestDeliveryPointApi* | [**DeleteMsgVpnRestDeliveryPointRestConsumer**](docs/RestDeliveryPointApi.md#deletemsgvpnrestdeliverypointrestconsumer) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName} | Delete a REST Consumer object.
*RestDeliveryPointApi* | [**DeleteMsgVpnRestDeliveryPointRestConsumerOauthJwtClaim**](docs/RestDeliveryPointApi.md#deletemsgvpnrestdeliverypointrestconsumeroauthjwtclaim) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/oauthJwtClaims/{oauthJwtClaimName} | Delete a Claim object.
*RestDeliveryPointApi* | [**DeleteMsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonName**](docs/RestDeliveryPointApi.md#deletemsgvpnrestdeliverypointrestconsumertlstrustedcommonname) | **Delete** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Delete a Trusted Common Name object.
*RestDeliveryPointApi* | [**GetMsgVpnRestDeliveryPoint**](docs/RestDeliveryPointApi.md#getmsgvpnrestdeliverypoint) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName} | Get a REST Delivery Point object.
*RestDeliveryPointApi* | [**GetMsgVpnRestDeliveryPointQueueBinding**](docs/RestDeliveryPointApi.md#getmsgvpnrestdeliverypointqueuebinding) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName} | Get a Queue Binding object.
*RestDeliveryPointApi* | [**GetMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/RestDeliveryPointApi.md#getmsgvpnrestdeliverypointqueuebindingrequestheader) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders/{headerName} | Get a Request Header object.
*RestDeliveryPointApi* | [**GetMsgVpnRestDeliveryPointQueueBindingRequestHeaders**](docs/RestDeliveryPointApi.md#getmsgvpnrestdeliverypointqueuebindingrequestheaders) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders | Get a list of Request Header objects.
*RestDeliveryPointApi* | [**GetMsgVpnRestDeliveryPointQueueBindings**](docs/RestDeliveryPointApi.md#getmsgvpnrestdeliverypointqueuebindings) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings | Get a list of Queue Binding objects.
*RestDeliveryPointApi* | [**GetMsgVpnRestDeliveryPointRestConsumer**](docs/RestDeliveryPointApi.md#getmsgvpnrestdeliverypointrestconsumer) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName} | Get a REST Consumer object.
*RestDeliveryPointApi* | [**GetMsgVpnRestDeliveryPointRestConsumerOauthJwtClaim**](docs/RestDeliveryPointApi.md#getmsgvpnrestdeliverypointrestconsumeroauthjwtclaim) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/oauthJwtClaims/{oauthJwtClaimName} | Get a Claim object.
*RestDeliveryPointApi* | [**GetMsgVpnRestDeliveryPointRestConsumerOauthJwtClaims**](docs/RestDeliveryPointApi.md#getmsgvpnrestdeliverypointrestconsumeroauthjwtclaims) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/oauthJwtClaims | Get a list of Claim objects.
*RestDeliveryPointApi* | [**GetMsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonName**](docs/RestDeliveryPointApi.md#getmsgvpnrestdeliverypointrestconsumertlstrustedcommonname) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/tlsTrustedCommonNames/{tlsTrustedCommonName} | Get a Trusted Common Name object.
*RestDeliveryPointApi* | [**GetMsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonNames**](docs/RestDeliveryPointApi.md#getmsgvpnrestdeliverypointrestconsumertlstrustedcommonnames) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName}/tlsTrustedCommonNames | Get a list of Trusted Common Name objects.
*RestDeliveryPointApi* | [**GetMsgVpnRestDeliveryPointRestConsumers**](docs/RestDeliveryPointApi.md#getmsgvpnrestdeliverypointrestconsumers) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers | Get a list of REST Consumer objects.
*RestDeliveryPointApi* | [**GetMsgVpnRestDeliveryPoints**](docs/RestDeliveryPointApi.md#getmsgvpnrestdeliverypoints) | **Get** /msgVpns/{msgVpnName}/restDeliveryPoints | Get a list of REST Delivery Point objects.
*RestDeliveryPointApi* | [**ReplaceMsgVpnRestDeliveryPoint**](docs/RestDeliveryPointApi.md#replacemsgvpnrestdeliverypoint) | **Put** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName} | Replace a REST Delivery Point object.
*RestDeliveryPointApi* | [**ReplaceMsgVpnRestDeliveryPointQueueBinding**](docs/RestDeliveryPointApi.md#replacemsgvpnrestdeliverypointqueuebinding) | **Put** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName} | Replace a Queue Binding object.
*RestDeliveryPointApi* | [**ReplaceMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/RestDeliveryPointApi.md#replacemsgvpnrestdeliverypointqueuebindingrequestheader) | **Put** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders/{headerName} | Replace a Request Header object.
*RestDeliveryPointApi* | [**ReplaceMsgVpnRestDeliveryPointRestConsumer**](docs/RestDeliveryPointApi.md#replacemsgvpnrestdeliverypointrestconsumer) | **Put** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName} | Replace a REST Consumer object.
*RestDeliveryPointApi* | [**UpdateMsgVpnRestDeliveryPoint**](docs/RestDeliveryPointApi.md#updatemsgvpnrestdeliverypoint) | **Patch** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName} | Update a REST Delivery Point object.
*RestDeliveryPointApi* | [**UpdateMsgVpnRestDeliveryPointQueueBinding**](docs/RestDeliveryPointApi.md#updatemsgvpnrestdeliverypointqueuebinding) | **Patch** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName} | Update a Queue Binding object.
*RestDeliveryPointApi* | [**UpdateMsgVpnRestDeliveryPointQueueBindingRequestHeader**](docs/RestDeliveryPointApi.md#updatemsgvpnrestdeliverypointqueuebindingrequestheader) | **Patch** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/queueBindings/{queueBindingName}/requestHeaders/{headerName} | Update a Request Header object.
*RestDeliveryPointApi* | [**UpdateMsgVpnRestDeliveryPointRestConsumer**](docs/RestDeliveryPointApi.md#updatemsgvpnrestdeliverypointrestconsumer) | **Patch** /msgVpns/{msgVpnName}/restDeliveryPoints/{restDeliveryPointName}/restConsumers/{restConsumerName} | Update a REST Consumer object.
*SystemInformationApi* | [**GetSystemInformation**](docs/SystemInformationApi.md#getsysteminformation) | **Get** /systemInformation | Get a System Information object.
*TopicEndpointApi* | [**CreateMsgVpnTopicEndpoint**](docs/TopicEndpointApi.md#createmsgvpntopicendpoint) | **Post** /msgVpns/{msgVpnName}/topicEndpoints | Create a Topic Endpoint object.
*TopicEndpointApi* | [**DeleteMsgVpnTopicEndpoint**](docs/TopicEndpointApi.md#deletemsgvpntopicendpoint) | **Delete** /msgVpns/{msgVpnName}/topicEndpoints/{topicEndpointName} | Delete a Topic Endpoint object.
*TopicEndpointApi* | [**GetMsgVpnTopicEndpoint**](docs/TopicEndpointApi.md#getmsgvpntopicendpoint) | **Get** /msgVpns/{msgVpnName}/topicEndpoints/{topicEndpointName} | Get a Topic Endpoint object.
*TopicEndpointApi* | [**GetMsgVpnTopicEndpoints**](docs/TopicEndpointApi.md#getmsgvpntopicendpoints) | **Get** /msgVpns/{msgVpnName}/topicEndpoints | Get a list of Topic Endpoint objects.
*TopicEndpointApi* | [**ReplaceMsgVpnTopicEndpoint**](docs/TopicEndpointApi.md#replacemsgvpntopicendpoint) | **Put** /msgVpns/{msgVpnName}/topicEndpoints/{topicEndpointName} | Replace a Topic Endpoint object.
*TopicEndpointApi* | [**UpdateMsgVpnTopicEndpoint**](docs/TopicEndpointApi.md#updatemsgvpntopicendpoint) | **Patch** /msgVpns/{msgVpnName}/topicEndpoints/{topicEndpointName} | Update a Topic Endpoint object.
*TopicEndpointTemplateApi* | [**CreateMsgVpnTopicEndpointTemplate**](docs/TopicEndpointTemplateApi.md#createmsgvpntopicendpointtemplate) | **Post** /msgVpns/{msgVpnName}/topicEndpointTemplates | Create a Topic Endpoint Template object.
*TopicEndpointTemplateApi* | [**DeleteMsgVpnTopicEndpointTemplate**](docs/TopicEndpointTemplateApi.md#deletemsgvpntopicendpointtemplate) | **Delete** /msgVpns/{msgVpnName}/topicEndpointTemplates/{topicEndpointTemplateName} | Delete a Topic Endpoint Template object.
*TopicEndpointTemplateApi* | [**GetMsgVpnTopicEndpointTemplate**](docs/TopicEndpointTemplateApi.md#getmsgvpntopicendpointtemplate) | **Get** /msgVpns/{msgVpnName}/topicEndpointTemplates/{topicEndpointTemplateName} | Get a Topic Endpoint Template object.
*TopicEndpointTemplateApi* | [**GetMsgVpnTopicEndpointTemplates**](docs/TopicEndpointTemplateApi.md#getmsgvpntopicendpointtemplates) | **Get** /msgVpns/{msgVpnName}/topicEndpointTemplates | Get a list of Topic Endpoint Template objects.
*TopicEndpointTemplateApi* | [**ReplaceMsgVpnTopicEndpointTemplate**](docs/TopicEndpointTemplateApi.md#replacemsgvpntopicendpointtemplate) | **Put** /msgVpns/{msgVpnName}/topicEndpointTemplates/{topicEndpointTemplateName} | Replace a Topic Endpoint Template object.
*TopicEndpointTemplateApi* | [**UpdateMsgVpnTopicEndpointTemplate**](docs/TopicEndpointTemplateApi.md#updatemsgvpntopicendpointtemplate) | **Patch** /msgVpns/{msgVpnName}/topicEndpointTemplates/{topicEndpointTemplateName} | Update a Topic Endpoint Template object.
*VirtualHostnameApi* | [**CreateVirtualHostname**](docs/VirtualHostnameApi.md#createvirtualhostname) | **Post** /virtualHostnames | Create a Virtual Hostname object.
*VirtualHostnameApi* | [**DeleteVirtualHostname**](docs/VirtualHostnameApi.md#deletevirtualhostname) | **Delete** /virtualHostnames/{virtualHostname} | Delete a Virtual Hostname object.
*VirtualHostnameApi* | [**GetVirtualHostname**](docs/VirtualHostnameApi.md#getvirtualhostname) | **Get** /virtualHostnames/{virtualHostname} | Get a Virtual Hostname object.
*VirtualHostnameApi* | [**GetVirtualHostnames**](docs/VirtualHostnameApi.md#getvirtualhostnames) | **Get** /virtualHostnames | Get a list of Virtual Hostname objects.
*VirtualHostnameApi* | [**ReplaceVirtualHostname**](docs/VirtualHostnameApi.md#replacevirtualhostname) | **Put** /virtualHostnames/{virtualHostname} | Replace a Virtual Hostname object.
*VirtualHostnameApi* | [**UpdateVirtualHostname**](docs/VirtualHostnameApi.md#updatevirtualhostname) | **Patch** /virtualHostnames/{virtualHostname} | Update a Virtual Hostname object.


## Documentation For Models

 - [AboutApi](docs/AboutApi.md)
 - [AboutApiLinks](docs/AboutApiLinks.md)
 - [AboutApiResponse](docs/AboutApiResponse.md)
 - [AboutLinks](docs/AboutLinks.md)
 - [AboutResponse](docs/AboutResponse.md)
 - [AboutUser](docs/AboutUser.md)
 - [AboutUserLinks](docs/AboutUserLinks.md)
 - [AboutUserMsgVpn](docs/AboutUserMsgVpn.md)
 - [AboutUserMsgVpnLinks](docs/AboutUserMsgVpnLinks.md)
 - [AboutUserMsgVpnResponse](docs/AboutUserMsgVpnResponse.md)
 - [AboutUserMsgVpnsResponse](docs/AboutUserMsgVpnsResponse.md)
 - [AboutUserResponse](docs/AboutUserResponse.md)
 - [Broker](docs/Broker.md)
 - [BrokerLinks](docs/BrokerLinks.md)
 - [BrokerResponse](docs/BrokerResponse.md)
 - [CertAuthoritiesResponse](docs/CertAuthoritiesResponse.md)
 - [CertAuthority](docs/CertAuthority.md)
 - [CertAuthorityLinks](docs/CertAuthorityLinks.md)
 - [CertAuthorityOcspTlsTrustedCommonName](docs/CertAuthorityOcspTlsTrustedCommonName.md)
 - [CertAuthorityOcspTlsTrustedCommonNameLinks](docs/CertAuthorityOcspTlsTrustedCommonNameLinks.md)
 - [CertAuthorityOcspTlsTrustedCommonNameResponse](docs/CertAuthorityOcspTlsTrustedCommonNameResponse.md)
 - [CertAuthorityOcspTlsTrustedCommonNamesResponse](docs/CertAuthorityOcspTlsTrustedCommonNamesResponse.md)
 - [CertAuthorityResponse](docs/CertAuthorityResponse.md)
 - [ClientCertAuthoritiesResponse](docs/ClientCertAuthoritiesResponse.md)
 - [ClientCertAuthority](docs/ClientCertAuthority.md)
 - [ClientCertAuthorityLinks](docs/ClientCertAuthorityLinks.md)
 - [ClientCertAuthorityOcspTlsTrustedCommonName](docs/ClientCertAuthorityOcspTlsTrustedCommonName.md)
 - [ClientCertAuthorityOcspTlsTrustedCommonNameLinks](docs/ClientCertAuthorityOcspTlsTrustedCommonNameLinks.md)
 - [ClientCertAuthorityOcspTlsTrustedCommonNameResponse](docs/ClientCertAuthorityOcspTlsTrustedCommonNameResponse.md)
 - [ClientCertAuthorityOcspTlsTrustedCommonNamesResponse](docs/ClientCertAuthorityOcspTlsTrustedCommonNamesResponse.md)
 - [ClientCertAuthorityResponse](docs/ClientCertAuthorityResponse.md)
 - [DmrCluster](docs/DmrCluster.md)
 - [DmrClusterLink](docs/DmrClusterLink.md)
 - [DmrClusterLinkLinks](docs/DmrClusterLinkLinks.md)
 - [DmrClusterLinkRemoteAddress](docs/DmrClusterLinkRemoteAddress.md)
 - [DmrClusterLinkRemoteAddressLinks](docs/DmrClusterLinkRemoteAddressLinks.md)
 - [DmrClusterLinkRemoteAddressResponse](docs/DmrClusterLinkRemoteAddressResponse.md)
 - [DmrClusterLinkRemoteAddressesResponse](docs/DmrClusterLinkRemoteAddressesResponse.md)
 - [DmrClusterLinkResponse](docs/DmrClusterLinkResponse.md)
 - [DmrClusterLinkTlsTrustedCommonName](docs/DmrClusterLinkTlsTrustedCommonName.md)
 - [DmrClusterLinkTlsTrustedCommonNameLinks](docs/DmrClusterLinkTlsTrustedCommonNameLinks.md)
 - [DmrClusterLinkTlsTrustedCommonNameResponse](docs/DmrClusterLinkTlsTrustedCommonNameResponse.md)
 - [DmrClusterLinkTlsTrustedCommonNamesResponse](docs/DmrClusterLinkTlsTrustedCommonNamesResponse.md)
 - [DmrClusterLinks](docs/DmrClusterLinks.md)
 - [DmrClusterLinksResponse](docs/DmrClusterLinksResponse.md)
 - [DmrClusterResponse](docs/DmrClusterResponse.md)
 - [DmrClustersResponse](docs/DmrClustersResponse.md)
 - [DomainCertAuthoritiesResponse](docs/DomainCertAuthoritiesResponse.md)
 - [DomainCertAuthority](docs/DomainCertAuthority.md)
 - [DomainCertAuthorityLinks](docs/DomainCertAuthorityLinks.md)
 - [DomainCertAuthorityResponse](docs/DomainCertAuthorityResponse.md)
 - [EventThreshold](docs/EventThreshold.md)
 - [EventThresholdByPercent](docs/EventThresholdByPercent.md)
 - [EventThresholdByValue](docs/EventThresholdByValue.md)
 - [MsgVpn](docs/MsgVpn.md)
 - [MsgVpnAclProfile](docs/MsgVpnAclProfile.md)
 - [MsgVpnAclProfileClientConnectException](docs/MsgVpnAclProfileClientConnectException.md)
 - [MsgVpnAclProfileClientConnectExceptionLinks](docs/MsgVpnAclProfileClientConnectExceptionLinks.md)
 - [MsgVpnAclProfileClientConnectExceptionResponse](docs/MsgVpnAclProfileClientConnectExceptionResponse.md)
 - [MsgVpnAclProfileClientConnectExceptionsResponse](docs/MsgVpnAclProfileClientConnectExceptionsResponse.md)
 - [MsgVpnAclProfileLinks](docs/MsgVpnAclProfileLinks.md)
 - [MsgVpnAclProfilePublishException](docs/MsgVpnAclProfilePublishException.md)
 - [MsgVpnAclProfilePublishExceptionLinks](docs/MsgVpnAclProfilePublishExceptionLinks.md)
 - [MsgVpnAclProfilePublishExceptionResponse](docs/MsgVpnAclProfilePublishExceptionResponse.md)
 - [MsgVpnAclProfilePublishExceptionsResponse](docs/MsgVpnAclProfilePublishExceptionsResponse.md)
 - [MsgVpnAclProfilePublishTopicException](docs/MsgVpnAclProfilePublishTopicException.md)
 - [MsgVpnAclProfilePublishTopicExceptionLinks](docs/MsgVpnAclProfilePublishTopicExceptionLinks.md)
 - [MsgVpnAclProfilePublishTopicExceptionResponse](docs/MsgVpnAclProfilePublishTopicExceptionResponse.md)
 - [MsgVpnAclProfilePublishTopicExceptionsResponse](docs/MsgVpnAclProfilePublishTopicExceptionsResponse.md)
 - [MsgVpnAclProfileResponse](docs/MsgVpnAclProfileResponse.md)
 - [MsgVpnAclProfileSubscribeException](docs/MsgVpnAclProfileSubscribeException.md)
 - [MsgVpnAclProfileSubscribeExceptionLinks](docs/MsgVpnAclProfileSubscribeExceptionLinks.md)
 - [MsgVpnAclProfileSubscribeExceptionResponse](docs/MsgVpnAclProfileSubscribeExceptionResponse.md)
 - [MsgVpnAclProfileSubscribeExceptionsResponse](docs/MsgVpnAclProfileSubscribeExceptionsResponse.md)
 - [MsgVpnAclProfileSubscribeShareNameException](docs/MsgVpnAclProfileSubscribeShareNameException.md)
 - [MsgVpnAclProfileSubscribeShareNameExceptionLinks](docs/MsgVpnAclProfileSubscribeShareNameExceptionLinks.md)
 - [MsgVpnAclProfileSubscribeShareNameExceptionResponse](docs/MsgVpnAclProfileSubscribeShareNameExceptionResponse.md)
 - [MsgVpnAclProfileSubscribeShareNameExceptionsResponse](docs/MsgVpnAclProfileSubscribeShareNameExceptionsResponse.md)
 - [MsgVpnAclProfileSubscribeTopicException](docs/MsgVpnAclProfileSubscribeTopicException.md)
 - [MsgVpnAclProfileSubscribeTopicExceptionLinks](docs/MsgVpnAclProfileSubscribeTopicExceptionLinks.md)
 - [MsgVpnAclProfileSubscribeTopicExceptionResponse](docs/MsgVpnAclProfileSubscribeTopicExceptionResponse.md)
 - [MsgVpnAclProfileSubscribeTopicExceptionsResponse](docs/MsgVpnAclProfileSubscribeTopicExceptionsResponse.md)
 - [MsgVpnAclProfilesResponse](docs/MsgVpnAclProfilesResponse.md)
 - [MsgVpnAuthenticationOauthProfile](docs/MsgVpnAuthenticationOauthProfile.md)
 - [MsgVpnAuthenticationOauthProfileClientRequiredClaim](docs/MsgVpnAuthenticationOauthProfileClientRequiredClaim.md)
 - [MsgVpnAuthenticationOauthProfileClientRequiredClaimLinks](docs/MsgVpnAuthenticationOauthProfileClientRequiredClaimLinks.md)
 - [MsgVpnAuthenticationOauthProfileClientRequiredClaimResponse](docs/MsgVpnAuthenticationOauthProfileClientRequiredClaimResponse.md)
 - [MsgVpnAuthenticationOauthProfileClientRequiredClaimsResponse](docs/MsgVpnAuthenticationOauthProfileClientRequiredClaimsResponse.md)
 - [MsgVpnAuthenticationOauthProfileLinks](docs/MsgVpnAuthenticationOauthProfileLinks.md)
 - [MsgVpnAuthenticationOauthProfileResourceServerRequiredClaim](docs/MsgVpnAuthenticationOauthProfileResourceServerRequiredClaim.md)
 - [MsgVpnAuthenticationOauthProfileResourceServerRequiredClaimLinks](docs/MsgVpnAuthenticationOauthProfileResourceServerRequiredClaimLinks.md)
 - [MsgVpnAuthenticationOauthProfileResourceServerRequiredClaimResponse](docs/MsgVpnAuthenticationOauthProfileResourceServerRequiredClaimResponse.md)
 - [MsgVpnAuthenticationOauthProfileResourceServerRequiredClaimsResponse](docs/MsgVpnAuthenticationOauthProfileResourceServerRequiredClaimsResponse.md)
 - [MsgVpnAuthenticationOauthProfileResponse](docs/MsgVpnAuthenticationOauthProfileResponse.md)
 - [MsgVpnAuthenticationOauthProfilesResponse](docs/MsgVpnAuthenticationOauthProfilesResponse.md)
 - [MsgVpnAuthenticationOauthProvider](docs/MsgVpnAuthenticationOauthProvider.md)
 - [MsgVpnAuthenticationOauthProviderLinks](docs/MsgVpnAuthenticationOauthProviderLinks.md)
 - [MsgVpnAuthenticationOauthProviderResponse](docs/MsgVpnAuthenticationOauthProviderResponse.md)
 - [MsgVpnAuthenticationOauthProvidersResponse](docs/MsgVpnAuthenticationOauthProvidersResponse.md)
 - [MsgVpnAuthorizationGroup](docs/MsgVpnAuthorizationGroup.md)
 - [MsgVpnAuthorizationGroupLinks](docs/MsgVpnAuthorizationGroupLinks.md)
 - [MsgVpnAuthorizationGroupResponse](docs/MsgVpnAuthorizationGroupResponse.md)
 - [MsgVpnAuthorizationGroupsResponse](docs/MsgVpnAuthorizationGroupsResponse.md)
 - [MsgVpnBridge](docs/MsgVpnBridge.md)
 - [MsgVpnBridgeLinks](docs/MsgVpnBridgeLinks.md)
 - [MsgVpnBridgeRemoteMsgVpn](docs/MsgVpnBridgeRemoteMsgVpn.md)
 - [MsgVpnBridgeRemoteMsgVpnLinks](docs/MsgVpnBridgeRemoteMsgVpnLinks.md)
 - [MsgVpnBridgeRemoteMsgVpnResponse](docs/MsgVpnBridgeRemoteMsgVpnResponse.md)
 - [MsgVpnBridgeRemoteMsgVpnsResponse](docs/MsgVpnBridgeRemoteMsgVpnsResponse.md)
 - [MsgVpnBridgeRemoteSubscription](docs/MsgVpnBridgeRemoteSubscription.md)
 - [MsgVpnBridgeRemoteSubscriptionLinks](docs/MsgVpnBridgeRemoteSubscriptionLinks.md)
 - [MsgVpnBridgeRemoteSubscriptionResponse](docs/MsgVpnBridgeRemoteSubscriptionResponse.md)
 - [MsgVpnBridgeRemoteSubscriptionsResponse](docs/MsgVpnBridgeRemoteSubscriptionsResponse.md)
 - [MsgVpnBridgeResponse](docs/MsgVpnBridgeResponse.md)
 - [MsgVpnBridgeTlsTrustedCommonName](docs/MsgVpnBridgeTlsTrustedCommonName.md)
 - [MsgVpnBridgeTlsTrustedCommonNameLinks](docs/MsgVpnBridgeTlsTrustedCommonNameLinks.md)
 - [MsgVpnBridgeTlsTrustedCommonNameResponse](docs/MsgVpnBridgeTlsTrustedCommonNameResponse.md)
 - [MsgVpnBridgeTlsTrustedCommonNamesResponse](docs/MsgVpnBridgeTlsTrustedCommonNamesResponse.md)
 - [MsgVpnBridgesResponse](docs/MsgVpnBridgesResponse.md)
 - [MsgVpnClientProfile](docs/MsgVpnClientProfile.md)
 - [MsgVpnClientProfileLinks](docs/MsgVpnClientProfileLinks.md)
 - [MsgVpnClientProfileResponse](docs/MsgVpnClientProfileResponse.md)
 - [MsgVpnClientProfilesResponse](docs/MsgVpnClientProfilesResponse.md)
 - [MsgVpnClientUsername](docs/MsgVpnClientUsername.md)
 - [MsgVpnClientUsernameLinks](docs/MsgVpnClientUsernameLinks.md)
 - [MsgVpnClientUsernameResponse](docs/MsgVpnClientUsernameResponse.md)
 - [MsgVpnClientUsernamesResponse](docs/MsgVpnClientUsernamesResponse.md)
 - [MsgVpnDistributedCache](docs/MsgVpnDistributedCache.md)
 - [MsgVpnDistributedCacheCluster](docs/MsgVpnDistributedCacheCluster.md)
 - [MsgVpnDistributedCacheClusterGlobalCachingHomeCluster](docs/MsgVpnDistributedCacheClusterGlobalCachingHomeCluster.md)
 - [MsgVpnDistributedCacheClusterGlobalCachingHomeClusterLinks](docs/MsgVpnDistributedCacheClusterGlobalCachingHomeClusterLinks.md)
 - [MsgVpnDistributedCacheClusterGlobalCachingHomeClusterResponse](docs/MsgVpnDistributedCacheClusterGlobalCachingHomeClusterResponse.md)
 - [MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix](docs/MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix.md)
 - [MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixLinks](docs/MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixLinks.md)
 - [MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixResponse](docs/MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixResponse.md)
 - [MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesResponse](docs/MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesResponse.md)
 - [MsgVpnDistributedCacheClusterGlobalCachingHomeClustersResponse](docs/MsgVpnDistributedCacheClusterGlobalCachingHomeClustersResponse.md)
 - [MsgVpnDistributedCacheClusterInstance](docs/MsgVpnDistributedCacheClusterInstance.md)
 - [MsgVpnDistributedCacheClusterInstanceLinks](docs/MsgVpnDistributedCacheClusterInstanceLinks.md)
 - [MsgVpnDistributedCacheClusterInstanceResponse](docs/MsgVpnDistributedCacheClusterInstanceResponse.md)
 - [MsgVpnDistributedCacheClusterInstancesResponse](docs/MsgVpnDistributedCacheClusterInstancesResponse.md)
 - [MsgVpnDistributedCacheClusterLinks](docs/MsgVpnDistributedCacheClusterLinks.md)
 - [MsgVpnDistributedCacheClusterResponse](docs/MsgVpnDistributedCacheClusterResponse.md)
 - [MsgVpnDistributedCacheClusterTopic](docs/MsgVpnDistributedCacheClusterTopic.md)
 - [MsgVpnDistributedCacheClusterTopicLinks](docs/MsgVpnDistributedCacheClusterTopicLinks.md)
 - [MsgVpnDistributedCacheClusterTopicResponse](docs/MsgVpnDistributedCacheClusterTopicResponse.md)
 - [MsgVpnDistributedCacheClusterTopicsResponse](docs/MsgVpnDistributedCacheClusterTopicsResponse.md)
 - [MsgVpnDistributedCacheClustersResponse](docs/MsgVpnDistributedCacheClustersResponse.md)
 - [MsgVpnDistributedCacheLinks](docs/MsgVpnDistributedCacheLinks.md)
 - [MsgVpnDistributedCacheResponse](docs/MsgVpnDistributedCacheResponse.md)
 - [MsgVpnDistributedCachesResponse](docs/MsgVpnDistributedCachesResponse.md)
 - [MsgVpnDmrBridge](docs/MsgVpnDmrBridge.md)
 - [MsgVpnDmrBridgeLinks](docs/MsgVpnDmrBridgeLinks.md)
 - [MsgVpnDmrBridgeResponse](docs/MsgVpnDmrBridgeResponse.md)
 - [MsgVpnDmrBridgesResponse](docs/MsgVpnDmrBridgesResponse.md)
 - [MsgVpnJndiConnectionFactoriesResponse](docs/MsgVpnJndiConnectionFactoriesResponse.md)
 - [MsgVpnJndiConnectionFactory](docs/MsgVpnJndiConnectionFactory.md)
 - [MsgVpnJndiConnectionFactoryLinks](docs/MsgVpnJndiConnectionFactoryLinks.md)
 - [MsgVpnJndiConnectionFactoryResponse](docs/MsgVpnJndiConnectionFactoryResponse.md)
 - [MsgVpnJndiQueue](docs/MsgVpnJndiQueue.md)
 - [MsgVpnJndiQueueLinks](docs/MsgVpnJndiQueueLinks.md)
 - [MsgVpnJndiQueueResponse](docs/MsgVpnJndiQueueResponse.md)
 - [MsgVpnJndiQueuesResponse](docs/MsgVpnJndiQueuesResponse.md)
 - [MsgVpnJndiTopic](docs/MsgVpnJndiTopic.md)
 - [MsgVpnJndiTopicLinks](docs/MsgVpnJndiTopicLinks.md)
 - [MsgVpnJndiTopicResponse](docs/MsgVpnJndiTopicResponse.md)
 - [MsgVpnJndiTopicsResponse](docs/MsgVpnJndiTopicsResponse.md)
 - [MsgVpnLinks](docs/MsgVpnLinks.md)
 - [MsgVpnMqttRetainCache](docs/MsgVpnMqttRetainCache.md)
 - [MsgVpnMqttRetainCacheLinks](docs/MsgVpnMqttRetainCacheLinks.md)
 - [MsgVpnMqttRetainCacheResponse](docs/MsgVpnMqttRetainCacheResponse.md)
 - [MsgVpnMqttRetainCachesResponse](docs/MsgVpnMqttRetainCachesResponse.md)
 - [MsgVpnMqttSession](docs/MsgVpnMqttSession.md)
 - [MsgVpnMqttSessionLinks](docs/MsgVpnMqttSessionLinks.md)
 - [MsgVpnMqttSessionResponse](docs/MsgVpnMqttSessionResponse.md)
 - [MsgVpnMqttSessionSubscription](docs/MsgVpnMqttSessionSubscription.md)
 - [MsgVpnMqttSessionSubscriptionLinks](docs/MsgVpnMqttSessionSubscriptionLinks.md)
 - [MsgVpnMqttSessionSubscriptionResponse](docs/MsgVpnMqttSessionSubscriptionResponse.md)
 - [MsgVpnMqttSessionSubscriptionsResponse](docs/MsgVpnMqttSessionSubscriptionsResponse.md)
 - [MsgVpnMqttSessionsResponse](docs/MsgVpnMqttSessionsResponse.md)
 - [MsgVpnQueue](docs/MsgVpnQueue.md)
 - [MsgVpnQueueLinks](docs/MsgVpnQueueLinks.md)
 - [MsgVpnQueueResponse](docs/MsgVpnQueueResponse.md)
 - [MsgVpnQueueSubscription](docs/MsgVpnQueueSubscription.md)
 - [MsgVpnQueueSubscriptionLinks](docs/MsgVpnQueueSubscriptionLinks.md)
 - [MsgVpnQueueSubscriptionResponse](docs/MsgVpnQueueSubscriptionResponse.md)
 - [MsgVpnQueueSubscriptionsResponse](docs/MsgVpnQueueSubscriptionsResponse.md)
 - [MsgVpnQueueTemplate](docs/MsgVpnQueueTemplate.md)
 - [MsgVpnQueueTemplateLinks](docs/MsgVpnQueueTemplateLinks.md)
 - [MsgVpnQueueTemplateResponse](docs/MsgVpnQueueTemplateResponse.md)
 - [MsgVpnQueueTemplatesResponse](docs/MsgVpnQueueTemplatesResponse.md)
 - [MsgVpnQueuesResponse](docs/MsgVpnQueuesResponse.md)
 - [MsgVpnReplayLog](docs/MsgVpnReplayLog.md)
 - [MsgVpnReplayLogLinks](docs/MsgVpnReplayLogLinks.md)
 - [MsgVpnReplayLogResponse](docs/MsgVpnReplayLogResponse.md)
 - [MsgVpnReplayLogsResponse](docs/MsgVpnReplayLogsResponse.md)
 - [MsgVpnReplicatedTopic](docs/MsgVpnReplicatedTopic.md)
 - [MsgVpnReplicatedTopicLinks](docs/MsgVpnReplicatedTopicLinks.md)
 - [MsgVpnReplicatedTopicResponse](docs/MsgVpnReplicatedTopicResponse.md)
 - [MsgVpnReplicatedTopicsResponse](docs/MsgVpnReplicatedTopicsResponse.md)
 - [MsgVpnResponse](docs/MsgVpnResponse.md)
 - [MsgVpnRestDeliveryPoint](docs/MsgVpnRestDeliveryPoint.md)
 - [MsgVpnRestDeliveryPointLinks](docs/MsgVpnRestDeliveryPointLinks.md)
 - [MsgVpnRestDeliveryPointQueueBinding](docs/MsgVpnRestDeliveryPointQueueBinding.md)
 - [MsgVpnRestDeliveryPointQueueBindingLinks](docs/MsgVpnRestDeliveryPointQueueBindingLinks.md)
 - [MsgVpnRestDeliveryPointQueueBindingRequestHeader](docs/MsgVpnRestDeliveryPointQueueBindingRequestHeader.md)
 - [MsgVpnRestDeliveryPointQueueBindingRequestHeaderLinks](docs/MsgVpnRestDeliveryPointQueueBindingRequestHeaderLinks.md)
 - [MsgVpnRestDeliveryPointQueueBindingRequestHeaderResponse](docs/MsgVpnRestDeliveryPointQueueBindingRequestHeaderResponse.md)
 - [MsgVpnRestDeliveryPointQueueBindingRequestHeadersResponse](docs/MsgVpnRestDeliveryPointQueueBindingRequestHeadersResponse.md)
 - [MsgVpnRestDeliveryPointQueueBindingResponse](docs/MsgVpnRestDeliveryPointQueueBindingResponse.md)
 - [MsgVpnRestDeliveryPointQueueBindingsResponse](docs/MsgVpnRestDeliveryPointQueueBindingsResponse.md)
 - [MsgVpnRestDeliveryPointResponse](docs/MsgVpnRestDeliveryPointResponse.md)
 - [MsgVpnRestDeliveryPointRestConsumer](docs/MsgVpnRestDeliveryPointRestConsumer.md)
 - [MsgVpnRestDeliveryPointRestConsumerLinks](docs/MsgVpnRestDeliveryPointRestConsumerLinks.md)
 - [MsgVpnRestDeliveryPointRestConsumerOauthJwtClaim](docs/MsgVpnRestDeliveryPointRestConsumerOauthJwtClaim.md)
 - [MsgVpnRestDeliveryPointRestConsumerOauthJwtClaimLinks](docs/MsgVpnRestDeliveryPointRestConsumerOauthJwtClaimLinks.md)
 - [MsgVpnRestDeliveryPointRestConsumerOauthJwtClaimResponse](docs/MsgVpnRestDeliveryPointRestConsumerOauthJwtClaimResponse.md)
 - [MsgVpnRestDeliveryPointRestConsumerOauthJwtClaimsResponse](docs/MsgVpnRestDeliveryPointRestConsumerOauthJwtClaimsResponse.md)
 - [MsgVpnRestDeliveryPointRestConsumerResponse](docs/MsgVpnRestDeliveryPointRestConsumerResponse.md)
 - [MsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonName](docs/MsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonName.md)
 - [MsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonNameLinks](docs/MsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonNameLinks.md)
 - [MsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonNameResponse](docs/MsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonNameResponse.md)
 - [MsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonNamesResponse](docs/MsgVpnRestDeliveryPointRestConsumerTlsTrustedCommonNamesResponse.md)
 - [MsgVpnRestDeliveryPointRestConsumersResponse](docs/MsgVpnRestDeliveryPointRestConsumersResponse.md)
 - [MsgVpnRestDeliveryPointsResponse](docs/MsgVpnRestDeliveryPointsResponse.md)
 - [MsgVpnSequencedTopic](docs/MsgVpnSequencedTopic.md)
 - [MsgVpnSequencedTopicLinks](docs/MsgVpnSequencedTopicLinks.md)
 - [MsgVpnSequencedTopicResponse](docs/MsgVpnSequencedTopicResponse.md)
 - [MsgVpnSequencedTopicsResponse](docs/MsgVpnSequencedTopicsResponse.md)
 - [MsgVpnTopicEndpoint](docs/MsgVpnTopicEndpoint.md)
 - [MsgVpnTopicEndpointLinks](docs/MsgVpnTopicEndpointLinks.md)
 - [MsgVpnTopicEndpointResponse](docs/MsgVpnTopicEndpointResponse.md)
 - [MsgVpnTopicEndpointTemplate](docs/MsgVpnTopicEndpointTemplate.md)
 - [MsgVpnTopicEndpointTemplateLinks](docs/MsgVpnTopicEndpointTemplateLinks.md)
 - [MsgVpnTopicEndpointTemplateResponse](docs/MsgVpnTopicEndpointTemplateResponse.md)
 - [MsgVpnTopicEndpointTemplatesResponse](docs/MsgVpnTopicEndpointTemplatesResponse.md)
 - [MsgVpnTopicEndpointsResponse](docs/MsgVpnTopicEndpointsResponse.md)
 - [MsgVpnsResponse](docs/MsgVpnsResponse.md)
 - [OauthProfile](docs/OauthProfile.md)
 - [OauthProfileAccessLevelGroup](docs/OauthProfileAccessLevelGroup.md)
 - [OauthProfileAccessLevelGroupLinks](docs/OauthProfileAccessLevelGroupLinks.md)
 - [OauthProfileAccessLevelGroupMsgVpnAccessLevelException](docs/OauthProfileAccessLevelGroupMsgVpnAccessLevelException.md)
 - [OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionLinks](docs/OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionLinks.md)
 - [OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse](docs/OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse.md)
 - [OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsResponse](docs/OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsResponse.md)
 - [OauthProfileAccessLevelGroupResponse](docs/OauthProfileAccessLevelGroupResponse.md)
 - [OauthProfileAccessLevelGroupsResponse](docs/OauthProfileAccessLevelGroupsResponse.md)
 - [OauthProfileClientAllowedHost](docs/OauthProfileClientAllowedHost.md)
 - [OauthProfileClientAllowedHostLinks](docs/OauthProfileClientAllowedHostLinks.md)
 - [OauthProfileClientAllowedHostResponse](docs/OauthProfileClientAllowedHostResponse.md)
 - [OauthProfileClientAllowedHostsResponse](docs/OauthProfileClientAllowedHostsResponse.md)
 - [OauthProfileClientAuthorizationParameter](docs/OauthProfileClientAuthorizationParameter.md)
 - [OauthProfileClientAuthorizationParameterLinks](docs/OauthProfileClientAuthorizationParameterLinks.md)
 - [OauthProfileClientAuthorizationParameterResponse](docs/OauthProfileClientAuthorizationParameterResponse.md)
 - [OauthProfileClientAuthorizationParametersResponse](docs/OauthProfileClientAuthorizationParametersResponse.md)
 - [OauthProfileClientRequiredClaim](docs/OauthProfileClientRequiredClaim.md)
 - [OauthProfileClientRequiredClaimLinks](docs/OauthProfileClientRequiredClaimLinks.md)
 - [OauthProfileClientRequiredClaimResponse](docs/OauthProfileClientRequiredClaimResponse.md)
 - [OauthProfileClientRequiredClaimsResponse](docs/OauthProfileClientRequiredClaimsResponse.md)
 - [OauthProfileDefaultMsgVpnAccessLevelException](docs/OauthProfileDefaultMsgVpnAccessLevelException.md)
 - [OauthProfileDefaultMsgVpnAccessLevelExceptionLinks](docs/OauthProfileDefaultMsgVpnAccessLevelExceptionLinks.md)
 - [OauthProfileDefaultMsgVpnAccessLevelExceptionResponse](docs/OauthProfileDefaultMsgVpnAccessLevelExceptionResponse.md)
 - [OauthProfileDefaultMsgVpnAccessLevelExceptionsResponse](docs/OauthProfileDefaultMsgVpnAccessLevelExceptionsResponse.md)
 - [OauthProfileLinks](docs/OauthProfileLinks.md)
 - [OauthProfileResourceServerRequiredClaim](docs/OauthProfileResourceServerRequiredClaim.md)
 - [OauthProfileResourceServerRequiredClaimLinks](docs/OauthProfileResourceServerRequiredClaimLinks.md)
 - [OauthProfileResourceServerRequiredClaimResponse](docs/OauthProfileResourceServerRequiredClaimResponse.md)
 - [OauthProfileResourceServerRequiredClaimsResponse](docs/OauthProfileResourceServerRequiredClaimsResponse.md)
 - [OauthProfileResponse](docs/OauthProfileResponse.md)
 - [OauthProfilesResponse](docs/OauthProfilesResponse.md)
 - [SempError](docs/SempError.md)
 - [SempMeta](docs/SempMeta.md)
 - [SempMetaOnlyResponse](docs/SempMetaOnlyResponse.md)
 - [SempPaging](docs/SempPaging.md)
 - [SempRequest](docs/SempRequest.md)
 - [SystemInformation](docs/SystemInformation.md)
 - [SystemInformationLinks](docs/SystemInformationLinks.md)
 - [SystemInformationResponse](docs/SystemInformationResponse.md)
 - [VirtualHostname](docs/VirtualHostname.md)
 - [VirtualHostnameLinks](docs/VirtualHostnameLinks.md)
 - [VirtualHostnameResponse](docs/VirtualHostnameResponse.md)
 - [VirtualHostnamesResponse](docs/VirtualHostnamesResponse.md)


## Documentation For Authorization



### basicAuth

- **Type**: HTTP basic authentication

Example

```golang
auth := context.WithValue(context.Background(), sw.ContextBasicAuth, sw.BasicAuth{
    UserName: "username",
    Password: "password",
})
r, err := client.Service.Operation(auth, args)
```


## Documentation for Utility Methods

Due to the fact that model structure members are all pointers, this package contains
a number of utility functions to easily obtain pointers to values of basic types.
Each of these functions takes a value of the given basic type and returns a pointer to it:

* `PtrBool`
* `PtrInt`
* `PtrInt32`
* `PtrInt64`
* `PtrFloat`
* `PtrFloat32`
* `PtrFloat64`
* `PtrString`
* `PtrTime`

## Author

support@solace.com

