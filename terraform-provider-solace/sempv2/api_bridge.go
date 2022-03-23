/*
SEMP (Solace Element Management Protocol)

SEMP (starting in `v2`, see note 1) is a RESTful API for configuring, monitoring, and administering a Solace PubSub+ broker.  SEMP uses URIs to address manageable **resources** of the Solace PubSub+ broker. Resources are individual **objects**, **collections** of objects, or (exclusively in the action API) **actions**. This document applies to the following API:   API|Base Path|Purpose|Comments :---|:---|:---|:--- Configuration|/SEMP/v2/config|Reading and writing config state|See note 2    The following APIs are also available:   API|Base Path|Purpose|Comments :---|:---|:---|:--- Action|/SEMP/v2/action|Performing actions|See note 2 Monitoring|/SEMP/v2/monitor|Querying operational parameters|See note 2    Resources are always nouns, with individual objects being singular and collections being plural.  Objects within a collection are identified by an `obj-id`, which follows the collection name with the form `collection-name/obj-id`.  Actions within an object are identified by an `action-id`, which follows the object name with the form `obj-id/action-id`.  Some examples:  ``` /SEMP/v2/config/msgVpns                        ; MsgVpn collection /SEMP/v2/config/msgVpns/a                      ; MsgVpn object named \"a\" /SEMP/v2/config/msgVpns/a/queues               ; Queue collection in MsgVpn \"a\" /SEMP/v2/config/msgVpns/a/queues/b             ; Queue object named \"b\" in MsgVpn \"a\" /SEMP/v2/action/msgVpns/a/queues/b/startReplay ; Action that starts a replay on Queue \"b\" in MsgVpn \"a\" /SEMP/v2/monitor/msgVpns/a/clients             ; Client collection in MsgVpn \"a\" /SEMP/v2/monitor/msgVpns/a/clients/c           ; Client object named \"c\" in MsgVpn \"a\" ```  ## Collection Resources  Collections are unordered lists of objects (unless described as otherwise), and are described by JSON arrays. Each item in the array represents an object in the same manner as the individual object would normally be represented. In the configuration API, the creation of a new object is done through its collection resource.  ## Object and Action Resources  Objects are composed of attributes, actions, collections, and other objects. They are described by JSON objects as name/value pairs. The collections and actions of an object are not contained directly in the object's JSON content; rather the content includes an attribute containing a URI which points to the collections and actions. These contained resources must be managed through this URI. At a minimum, every object has one or more identifying attributes, and its own `uri` attribute which contains the URI pointing to itself.  Actions are also composed of attributes, and are described by JSON objects as name/value pairs. Unlike objects, however, they are not members of a collection and cannot be retrieved, only performed. Actions only exist in the action API.  Attributes in an object or action may have any combination of the following properties:   Property|Meaning|Comments :---|:---|:--- Identifying|Attribute is involved in unique identification of the object, and appears in its URI| Required|Attribute must be provided in the request| Read-Only|Attribute can only be read, not written.|See note 3 Write-Only|Attribute can only be written, not read, unless the attribute is also opaque|See the documentation for the opaque property Requires-Disable|Attribute can only be changed when object is disabled| Deprecated|Attribute is deprecated, and will disappear in the next SEMP version| Opaque|Attribute can be set or retrieved in opaque form when the `opaquePassword` query parameter is present|See the `opaquePassword` query parameter documentation    In some requests, certain attributes may only be provided in certain combinations with other attributes:   Relationship|Meaning :---|:--- Requires|Attribute may only be changed by a request if a particular attribute or combination of attributes is also provided in the request Conflicts|Attribute may only be provided in a request if a particular attribute or combination of attributes is not also provided in the request    In the monitoring API, any non-identifying attribute may not be returned in a GET.  ## HTTP Methods  The following HTTP methods manipulate resources in accordance with these general principles. Note that some methods are only used in certain APIs:   Method|Resource|Meaning|Request Body|Response Body|Missing Request Attributes :---|:---|:---|:---|:---|:--- POST|Collection|Create object|Initial attribute values|Object attributes and metadata|Set to default PUT|Object|Create or replace object (see note 5)|New attribute values|Object attributes and metadata|Set to default, with certain exceptions (see note 4) PUT|Action|Performs action|Action arguments|Action metadata|N/A PATCH|Object|Update object|New attribute values|Object attributes and metadata|unchanged DELETE|Object|Delete object|Empty|Object metadata|N/A GET|Object|Get object|Empty|Object attributes and metadata|N/A GET|Collection|Get collection|Empty|Object attributes and collection metadata|N/A    ## Common Query Parameters  The following are some common query parameters that are supported by many method/URI combinations. Individual URIs may document additional parameters. Note that multiple query parameters can be used together in a single URI, separated by the ampersand character. For example:  ``` ; Request for the MsgVpns collection using two hypothetical query parameters ; \"q1\" and \"q2\" with values \"val1\" and \"val2\" respectively /SEMP/v2/config/msgVpns?q1=val1&q2=val2 ```  ### select  Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. Use this query parameter to limit the size of the returned data for each returned object, return only those fields that are desired, or exclude fields that are not desired.  The value of `select` is a comma-separated list of attribute names. If the list contains attribute names that are not prefaced by `-`, only those attributes are included in the response. If the list contains attribute names that are prefaced by `-`, those attributes are excluded from the response. If the list contains both types, then the difference of the first set of attributes and the second set of attributes is returned. If the list is empty (i.e. `select=`), no attributes are returned.  All attributes that are prefaced by `-` must follow all attributes that are not prefaced by `-`. In addition, each attribute name in the list must match at least one attribute in the object.  Names may include the `*` wildcard (zero or more characters). Nested attribute names are supported using periods (e.g. `parentName.childName`).  Some examples:  ``` ; List of all MsgVpn names /SEMP/v2/config/msgVpns?select=msgVpnName ; List of all MsgVpn and their attributes except for their names /SEMP/v2/config/msgVpns?select=-msgVpnName ; Authentication attributes of MsgVpn \"finance\" /SEMP/v2/config/msgVpns/finance?select=authentication* ; All attributes of MsgVpn \"finance\" except for authentication attributes /SEMP/v2/config/msgVpns/finance?select=-authentication* ; Access related attributes of Queue \"orderQ\" of MsgVpn \"finance\" /SEMP/v2/config/msgVpns/finance/queues/orderQ?select=owner,permission ```  ### where  Include in the response only objects where certain conditions are true. Use this query parameter to limit which objects are returned to those whose attribute values meet the given conditions.  The value of `where` is a comma-separated list of expressions. All expressions must be true for the object to be included in the response. Each expression takes the form:  ``` expression  = attribute-name OP value OP          = '==' | '!=' | '&lt;' | '&gt;' | '&lt;=' | '&gt;=' ```  `value` may be a number, string, `true`, or `false`, as appropriate for the type of `attribute-name`. Greater-than and less-than comparisons only work for numbers. A `*` in a string `value` is interpreted as a wildcard (zero or more characters). Some examples:  ``` ; Only enabled MsgVpns /SEMP/v2/config/msgVpns?where=enabled==true ; Only MsgVpns using basic non-LDAP authentication /SEMP/v2/config/msgVpns?where=authenticationBasicEnabled==true,authenticationBasicType!=ldap ; Only MsgVpns that allow more than 100 client connections /SEMP/v2/config/msgVpns?where=maxConnectionCount>100 ; Only MsgVpns with msgVpnName starting with \"B\": /SEMP/v2/config/msgVpns?where=msgVpnName==B* ```  ### count  Limit the count of objects in the response. This can be useful to limit the size of the response for large collections. The minimum value for `count` is `1` and the default is `10`. There is also a per-collection maximum value to limit request handling time.  `count` does not guarantee that a minimum number of objects will be returned. A page may contain fewer than `count` objects or even be empty. Additional objects may nonetheless be available for retrieval on subsequent pages. See the `cursor` query parameter documentation for more information on paging.  For example: ``` ; Up to 25 MsgVpns /SEMP/v2/config/msgVpns?count=25 ```  ### cursor  The cursor, or position, for the next page of objects. Cursors are opaque data that should not be created or interpreted by SEMP clients, and should only be used as described below.  When a request is made for a collection and there may be additional objects available for retrieval that are not included in the initial response, the response will include a `cursorQuery` field containing a cursor. The value of this field can be specified in the `cursor` query parameter of a subsequent request to retrieve the next page of objects. For convenience, an appropriate URI is constructed automatically by the broker and included in the `nextPageUri` field of the response. This URI can be used directly to retrieve the next page of objects.  Applications must continue to follow the `nextPageUri` if one is provided in order to retrieve the full set of objects associated with the request, even if a page contains fewer than the requested number of objects (see the `count` query parameter documentation) or is empty.  ### opaquePassword  Attributes with the opaque property are also write-only and so cannot normally be retrieved in a GET. However, when a password is provided in the `opaquePassword` query parameter, attributes with the opaque property are retrieved in a GET in opaque form, encrypted with this password. The query parameter can also be used on a POST, PATCH, or PUT to set opaque attributes using opaque attribute values retrieved in a GET, so long as:  1. the same password that was used to retrieve the opaque attribute values is provided; and  2. the broker to which the request is being sent has the same major and minor SEMP version as the broker that produced the opaque attribute values.  The password provided in the query parameter must be a minimum of 8 characters and a maximum of 128 characters.  The query parameter can only be used in the configuration API, and only over HTTPS.  ## Authentication  When a client makes its first SEMPv2 request, it must supply a username and password using HTTP Basic authentication, or an OAuth token or tokens using HTTP Bearer authentication.  When HTTP Basic authentication is used, the broker returns a cookie containing a session key. The client can omit the username and password from subsequent requests, because the broker can use the session cookie for authentication instead. When the session expires or is deleted, the client must provide the username and password again, and the broker creates a new session.  There are a limited number of session slots available on the broker. The broker returns 529 No SEMP Session Available if it is not able to allocate a session.  If certain attributes—such as a user's password—are changed, the broker automatically deletes the affected sessions. These attributes are documented below. However, changes in external user configuration data stored on a RADIUS or LDAP server do not trigger the broker to delete the associated session(s), therefore you must do this manually, if required.  A client can retrieve its current session information using the /about/user endpoint and delete its own session using the /about/user/logout endpoint. A client with appropriate permissions can also manage all sessions using the /sessions endpoint.  Sessions are not created when authenticating with an OAuth token or tokens using HTTP Bearer authentication. If a session cookie is provided, it is ignored.  ## Help  Visit [our website](https://solace.com) to learn more about Solace.  You can also download the SEMP API specifications by clicking [here](https://solace.com/downloads/).  If you need additional support, please contact us at [support@solace.com](mailto:support@solace.com).  ## Notes  Note|Description :---:|:--- 1|This specification defines SEMP starting in \"v2\", and not the original SEMP \"v1\" interface. Request and response formats between \"v1\" and \"v2\" are entirely incompatible, although both protocols share a common port configuration on the Solace PubSub+ broker. They are differentiated by the initial portion of the URI path, one of either \"/SEMP/\" or \"/SEMP/v2/\" 2|This API is partially implemented. Only a subset of all objects are available. 3|Read-only attributes may appear in POST and PUT/PATCH requests. However, if a read-only attribute is not marked as identifying, it will be ignored during a PUT/PATCH. 4|On a PUT, if the SEMP user is not authorized to modify the attribute, its value is left unchanged rather than set to default. In addition, the values of write-only attributes are not set to their defaults on a PUT, except in the following two cases: there is a mutual requires relationship with another non-write-only attribute, both attributes are absent from the request, and the non-write-only attribute is not currently set to its default value; or the attribute is also opaque and the `opaquePassword` query parameter is provided in the request. 5|On a PUT, if the object does not exist, it is created first.

API version: 2.26
Contact: support@solace.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package sempv2

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Linger please
var (
	_ context.Context
)

// BridgeApiService BridgeApi service
type BridgeApiService service

type BridgeApiApiCreateMsgVpnBridgeRequest struct {
	ctx            context.Context
	ApiService     *BridgeApiService
	msgVpnName     string
	body           *MsgVpnBridge
	opaquePassword *string
	select_        *[]string
}

// The Bridge object&#39;s attributes.
func (r BridgeApiApiCreateMsgVpnBridgeRequest) Body(body MsgVpnBridge) BridgeApiApiCreateMsgVpnBridgeRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiCreateMsgVpnBridgeRequest) OpaquePassword(opaquePassword string) BridgeApiApiCreateMsgVpnBridgeRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiCreateMsgVpnBridgeRequest) Select_(select_ []string) BridgeApiApiCreateMsgVpnBridgeRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiCreateMsgVpnBridgeRequest) Execute() (*MsgVpnBridgeResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnBridgeExecute(r)
}

/*
CreateMsgVpnBridge Create a Bridge object.

Create a Bridge object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

Bridges can be used to link two Message VPNs so that messages published to one Message VPN that match the topic subscriptions set for the bridge are also delivered to the linked Message VPN.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
bridgeName|x|x||||
bridgeVirtualRouter|x|x||||
msgVpnName|x||x|||
remoteAuthenticationBasicPassword||||x||x
remoteAuthenticationClientCertContent||||x||x
remoteAuthenticationClientCertPassword||||x||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
MsgVpnBridge|remoteAuthenticationBasicClientUsername|remoteAuthenticationBasicPassword|
MsgVpnBridge|remoteAuthenticationBasicPassword|remoteAuthenticationBasicClientUsername|
MsgVpnBridge|remoteAuthenticationClientCertPassword|remoteAuthenticationClientCertContent|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @return BridgeApiApiCreateMsgVpnBridgeRequest
*/
func (a *BridgeApiService) CreateMsgVpnBridge(ctx context.Context, msgVpnName string) BridgeApiApiCreateMsgVpnBridgeRequest {
	return BridgeApiApiCreateMsgVpnBridgeRequest{
		ApiService: a,
		ctx:        ctx,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeResponse
func (a *BridgeApiService) CreateMsgVpnBridgeExecute(r BridgeApiApiCreateMsgVpnBridgeRequest) (*MsgVpnBridgeResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.CreateMsgVpnBridge")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.body == nil {
		return localVarReturnValue, nil, reportError("body is required and must be specified")
	}

	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.body
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiCreateMsgVpnBridgeRemoteMsgVpnRequest struct {
	ctx                 context.Context
	ApiService          *BridgeApiService
	msgVpnName          string
	bridgeName          string
	bridgeVirtualRouter string
	body                *MsgVpnBridgeRemoteMsgVpn
	opaquePassword      *string
	select_             *[]string
}

// The Remote Message VPN object&#39;s attributes.
func (r BridgeApiApiCreateMsgVpnBridgeRemoteMsgVpnRequest) Body(body MsgVpnBridgeRemoteMsgVpn) BridgeApiApiCreateMsgVpnBridgeRemoteMsgVpnRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiCreateMsgVpnBridgeRemoteMsgVpnRequest) OpaquePassword(opaquePassword string) BridgeApiApiCreateMsgVpnBridgeRemoteMsgVpnRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiCreateMsgVpnBridgeRemoteMsgVpnRequest) Select_(select_ []string) BridgeApiApiCreateMsgVpnBridgeRemoteMsgVpnRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiCreateMsgVpnBridgeRemoteMsgVpnRequest) Execute() (*MsgVpnBridgeRemoteMsgVpnResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnBridgeRemoteMsgVpnExecute(r)
}

/*
CreateMsgVpnBridgeRemoteMsgVpn Create a Remote Message VPN object.

Create a Remote Message VPN object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

The Remote Message VPN is the Message VPN that the Bridge connects to.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
bridgeName|x||x|||
bridgeVirtualRouter|x||x|||
msgVpnName|x||x|||
password||||x||x
remoteMsgVpnInterface|x|||||
remoteMsgVpnLocation|x|x||||
remoteMsgVpnName|x|x||||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
MsgVpnBridgeRemoteMsgVpn|clientUsername|password|
MsgVpnBridgeRemoteMsgVpn|password|clientUsername|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @return BridgeApiApiCreateMsgVpnBridgeRemoteMsgVpnRequest
*/
func (a *BridgeApiService) CreateMsgVpnBridgeRemoteMsgVpn(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string) BridgeApiApiCreateMsgVpnBridgeRemoteMsgVpnRequest {
	return BridgeApiApiCreateMsgVpnBridgeRemoteMsgVpnRequest{
		ApiService:          a,
		ctx:                 ctx,
		msgVpnName:          msgVpnName,
		bridgeName:          bridgeName,
		bridgeVirtualRouter: bridgeVirtualRouter,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeRemoteMsgVpnResponse
func (a *BridgeApiService) CreateMsgVpnBridgeRemoteMsgVpnExecute(r BridgeApiApiCreateMsgVpnBridgeRemoteMsgVpnRequest) (*MsgVpnBridgeRemoteMsgVpnResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeRemoteMsgVpnResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.CreateMsgVpnBridgeRemoteMsgVpn")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.body == nil {
		return localVarReturnValue, nil, reportError("body is required and must be specified")
	}

	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.body
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiCreateMsgVpnBridgeRemoteSubscriptionRequest struct {
	ctx                 context.Context
	ApiService          *BridgeApiService
	msgVpnName          string
	bridgeName          string
	bridgeVirtualRouter string
	body                *MsgVpnBridgeRemoteSubscription
	opaquePassword      *string
	select_             *[]string
}

// The Remote Subscription object&#39;s attributes.
func (r BridgeApiApiCreateMsgVpnBridgeRemoteSubscriptionRequest) Body(body MsgVpnBridgeRemoteSubscription) BridgeApiApiCreateMsgVpnBridgeRemoteSubscriptionRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiCreateMsgVpnBridgeRemoteSubscriptionRequest) OpaquePassword(opaquePassword string) BridgeApiApiCreateMsgVpnBridgeRemoteSubscriptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiCreateMsgVpnBridgeRemoteSubscriptionRequest) Select_(select_ []string) BridgeApiApiCreateMsgVpnBridgeRemoteSubscriptionRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiCreateMsgVpnBridgeRemoteSubscriptionRequest) Execute() (*MsgVpnBridgeRemoteSubscriptionResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnBridgeRemoteSubscriptionExecute(r)
}

/*
CreateMsgVpnBridgeRemoteSubscription Create a Remote Subscription object.

Create a Remote Subscription object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Remote Subscription is a topic subscription used by the Message VPN Bridge to attract messages from the remote message broker.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
bridgeName|x||x|||
bridgeVirtualRouter|x||x|||
deliverAlwaysEnabled||x||||
msgVpnName|x||x|||
remoteSubscriptionTopic|x|x||||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @return BridgeApiApiCreateMsgVpnBridgeRemoteSubscriptionRequest
*/
func (a *BridgeApiService) CreateMsgVpnBridgeRemoteSubscription(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string) BridgeApiApiCreateMsgVpnBridgeRemoteSubscriptionRequest {
	return BridgeApiApiCreateMsgVpnBridgeRemoteSubscriptionRequest{
		ApiService:          a,
		ctx:                 ctx,
		msgVpnName:          msgVpnName,
		bridgeName:          bridgeName,
		bridgeVirtualRouter: bridgeVirtualRouter,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeRemoteSubscriptionResponse
func (a *BridgeApiService) CreateMsgVpnBridgeRemoteSubscriptionExecute(r BridgeApiApiCreateMsgVpnBridgeRemoteSubscriptionRequest) (*MsgVpnBridgeRemoteSubscriptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeRemoteSubscriptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.CreateMsgVpnBridgeRemoteSubscription")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.body == nil {
		return localVarReturnValue, nil, reportError("body is required and must be specified")
	}

	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.body
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiCreateMsgVpnBridgeTlsTrustedCommonNameRequest struct {
	ctx                 context.Context
	ApiService          *BridgeApiService
	msgVpnName          string
	bridgeName          string
	bridgeVirtualRouter string
	body                *MsgVpnBridgeTlsTrustedCommonName
	opaquePassword      *string
	select_             *[]string
}

// The Trusted Common Name object&#39;s attributes.
func (r BridgeApiApiCreateMsgVpnBridgeTlsTrustedCommonNameRequest) Body(body MsgVpnBridgeTlsTrustedCommonName) BridgeApiApiCreateMsgVpnBridgeTlsTrustedCommonNameRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiCreateMsgVpnBridgeTlsTrustedCommonNameRequest) OpaquePassword(opaquePassword string) BridgeApiApiCreateMsgVpnBridgeTlsTrustedCommonNameRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiCreateMsgVpnBridgeTlsTrustedCommonNameRequest) Select_(select_ []string) BridgeApiApiCreateMsgVpnBridgeTlsTrustedCommonNameRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiCreateMsgVpnBridgeTlsTrustedCommonNameRequest) Execute() (*MsgVpnBridgeTlsTrustedCommonNameResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnBridgeTlsTrustedCommonNameExecute(r)
}

/*
CreateMsgVpnBridgeTlsTrustedCommonName Create a Trusted Common Name object.

Create a Trusted Common Name object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

The Trusted Common Names for the Bridge are used by encrypted transports to verify the name in the certificate presented by the remote node. They must include the common name of the remote node's server certificate or client certificate, depending upon the initiator of the connection.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
bridgeName|x||x||x|
bridgeVirtualRouter|x||x||x|
msgVpnName|x||x||x|
tlsTrustedCommonName|x|x|||x|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been deprecated since 2.18. Common Name validation has been replaced by Server Certificate Name validation.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @return BridgeApiApiCreateMsgVpnBridgeTlsTrustedCommonNameRequest

Deprecated
*/
func (a *BridgeApiService) CreateMsgVpnBridgeTlsTrustedCommonName(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string) BridgeApiApiCreateMsgVpnBridgeTlsTrustedCommonNameRequest {
	return BridgeApiApiCreateMsgVpnBridgeTlsTrustedCommonNameRequest{
		ApiService:          a,
		ctx:                 ctx,
		msgVpnName:          msgVpnName,
		bridgeName:          bridgeName,
		bridgeVirtualRouter: bridgeVirtualRouter,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeTlsTrustedCommonNameResponse
// Deprecated
func (a *BridgeApiService) CreateMsgVpnBridgeTlsTrustedCommonNameExecute(r BridgeApiApiCreateMsgVpnBridgeTlsTrustedCommonNameRequest) (*MsgVpnBridgeTlsTrustedCommonNameResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeTlsTrustedCommonNameResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.CreateMsgVpnBridgeTlsTrustedCommonName")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.body == nil {
		return localVarReturnValue, nil, reportError("body is required and must be specified")
	}

	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.body
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiDeleteMsgVpnBridgeRequest struct {
	ctx                 context.Context
	ApiService          *BridgeApiService
	msgVpnName          string
	bridgeName          string
	bridgeVirtualRouter string
}

func (r BridgeApiApiDeleteMsgVpnBridgeRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnBridgeExecute(r)
}

/*
DeleteMsgVpnBridge Delete a Bridge object.

Delete a Bridge object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

Bridges can be used to link two Message VPNs so that messages published to one Message VPN that match the topic subscriptions set for the bridge are also delivered to the linked Message VPN.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @return BridgeApiApiDeleteMsgVpnBridgeRequest
*/
func (a *BridgeApiService) DeleteMsgVpnBridge(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string) BridgeApiApiDeleteMsgVpnBridgeRequest {
	return BridgeApiApiDeleteMsgVpnBridgeRequest{
		ApiService:          a,
		ctx:                 ctx,
		msgVpnName:          msgVpnName,
		bridgeName:          bridgeName,
		bridgeVirtualRouter: bridgeVirtualRouter,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *BridgeApiService) DeleteMsgVpnBridgeExecute(r BridgeApiApiDeleteMsgVpnBridgeRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.DeleteMsgVpnBridge")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiDeleteMsgVpnBridgeRemoteMsgVpnRequest struct {
	ctx                   context.Context
	ApiService            *BridgeApiService
	msgVpnName            string
	bridgeName            string
	bridgeVirtualRouter   string
	remoteMsgVpnName      string
	remoteMsgVpnLocation  string
	remoteMsgVpnInterface string
}

func (r BridgeApiApiDeleteMsgVpnBridgeRemoteMsgVpnRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnBridgeRemoteMsgVpnExecute(r)
}

/*
DeleteMsgVpnBridgeRemoteMsgVpn Delete a Remote Message VPN object.

Delete a Remote Message VPN object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

The Remote Message VPN is the Message VPN that the Bridge connects to.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @param remoteMsgVpnName The name of the remote Message VPN.
 @param remoteMsgVpnLocation The location of the remote Message VPN as either an FQDN with port, IP address with port, or virtual router name (starting with \"v:\").
 @param remoteMsgVpnInterface The physical interface on the local Message VPN host for connecting to the remote Message VPN. By default, an interface is chosen automatically (recommended), but if specified, `remoteMsgVpnLocation` must not be a virtual router name.
 @return BridgeApiApiDeleteMsgVpnBridgeRemoteMsgVpnRequest
*/
func (a *BridgeApiService) DeleteMsgVpnBridgeRemoteMsgVpn(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string, remoteMsgVpnName string, remoteMsgVpnLocation string, remoteMsgVpnInterface string) BridgeApiApiDeleteMsgVpnBridgeRemoteMsgVpnRequest {
	return BridgeApiApiDeleteMsgVpnBridgeRemoteMsgVpnRequest{
		ApiService:            a,
		ctx:                   ctx,
		msgVpnName:            msgVpnName,
		bridgeName:            bridgeName,
		bridgeVirtualRouter:   bridgeVirtualRouter,
		remoteMsgVpnName:      remoteMsgVpnName,
		remoteMsgVpnLocation:  remoteMsgVpnLocation,
		remoteMsgVpnInterface: remoteMsgVpnInterface,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *BridgeApiService) DeleteMsgVpnBridgeRemoteMsgVpnExecute(r BridgeApiApiDeleteMsgVpnBridgeRemoteMsgVpnRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.DeleteMsgVpnBridgeRemoteMsgVpn")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteMsgVpnName"+"}", url.PathEscape(parameterToString(r.remoteMsgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteMsgVpnLocation"+"}", url.PathEscape(parameterToString(r.remoteMsgVpnLocation, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteMsgVpnInterface"+"}", url.PathEscape(parameterToString(r.remoteMsgVpnInterface, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiDeleteMsgVpnBridgeRemoteSubscriptionRequest struct {
	ctx                     context.Context
	ApiService              *BridgeApiService
	msgVpnName              string
	bridgeName              string
	bridgeVirtualRouter     string
	remoteSubscriptionTopic string
}

func (r BridgeApiApiDeleteMsgVpnBridgeRemoteSubscriptionRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnBridgeRemoteSubscriptionExecute(r)
}

/*
DeleteMsgVpnBridgeRemoteSubscription Delete a Remote Subscription object.

Delete a Remote Subscription object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Remote Subscription is a topic subscription used by the Message VPN Bridge to attract messages from the remote message broker.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @param remoteSubscriptionTopic The topic of the Bridge remote subscription.
 @return BridgeApiApiDeleteMsgVpnBridgeRemoteSubscriptionRequest
*/
func (a *BridgeApiService) DeleteMsgVpnBridgeRemoteSubscription(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string, remoteSubscriptionTopic string) BridgeApiApiDeleteMsgVpnBridgeRemoteSubscriptionRequest {
	return BridgeApiApiDeleteMsgVpnBridgeRemoteSubscriptionRequest{
		ApiService:              a,
		ctx:                     ctx,
		msgVpnName:              msgVpnName,
		bridgeName:              bridgeName,
		bridgeVirtualRouter:     bridgeVirtualRouter,
		remoteSubscriptionTopic: remoteSubscriptionTopic,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *BridgeApiService) DeleteMsgVpnBridgeRemoteSubscriptionExecute(r BridgeApiApiDeleteMsgVpnBridgeRemoteSubscriptionRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.DeleteMsgVpnBridgeRemoteSubscription")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions/{remoteSubscriptionTopic}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteSubscriptionTopic"+"}", url.PathEscape(parameterToString(r.remoteSubscriptionTopic, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiDeleteMsgVpnBridgeTlsTrustedCommonNameRequest struct {
	ctx                  context.Context
	ApiService           *BridgeApiService
	msgVpnName           string
	bridgeName           string
	bridgeVirtualRouter  string
	tlsTrustedCommonName string
}

func (r BridgeApiApiDeleteMsgVpnBridgeTlsTrustedCommonNameRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnBridgeTlsTrustedCommonNameExecute(r)
}

/*
DeleteMsgVpnBridgeTlsTrustedCommonName Delete a Trusted Common Name object.

Delete a Trusted Common Name object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

The Trusted Common Names for the Bridge are used by encrypted transports to verify the name in the certificate presented by the remote node. They must include the common name of the remote node's server certificate or client certificate, depending upon the initiator of the connection.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been deprecated since 2.18. Common Name validation has been replaced by Server Certificate Name validation.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @param tlsTrustedCommonName The expected trusted common name of the remote certificate.
 @return BridgeApiApiDeleteMsgVpnBridgeTlsTrustedCommonNameRequest

Deprecated
*/
func (a *BridgeApiService) DeleteMsgVpnBridgeTlsTrustedCommonName(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string, tlsTrustedCommonName string) BridgeApiApiDeleteMsgVpnBridgeTlsTrustedCommonNameRequest {
	return BridgeApiApiDeleteMsgVpnBridgeTlsTrustedCommonNameRequest{
		ApiService:           a,
		ctx:                  ctx,
		msgVpnName:           msgVpnName,
		bridgeName:           bridgeName,
		bridgeVirtualRouter:  bridgeVirtualRouter,
		tlsTrustedCommonName: tlsTrustedCommonName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
// Deprecated
func (a *BridgeApiService) DeleteMsgVpnBridgeTlsTrustedCommonNameExecute(r BridgeApiApiDeleteMsgVpnBridgeTlsTrustedCommonNameRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.DeleteMsgVpnBridgeTlsTrustedCommonName")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames/{tlsTrustedCommonName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"tlsTrustedCommonName"+"}", url.PathEscape(parameterToString(r.tlsTrustedCommonName, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiGetMsgVpnBridgeRequest struct {
	ctx                 context.Context
	ApiService          *BridgeApiService
	msgVpnName          string
	bridgeName          string
	bridgeVirtualRouter string
	opaquePassword      *string
	select_             *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeRequest) OpaquePassword(opaquePassword string) BridgeApiApiGetMsgVpnBridgeRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeRequest) Select_(select_ []string) BridgeApiApiGetMsgVpnBridgeRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiGetMsgVpnBridgeRequest) Execute() (*MsgVpnBridgeResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnBridgeExecute(r)
}

/*
GetMsgVpnBridge Get a Bridge object.

Get a Bridge object.

Bridges can be used to link two Message VPNs so that messages published to one Message VPN that match the topic subscriptions set for the bridge are also delivered to the linked Message VPN.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
bridgeName|x|||
bridgeVirtualRouter|x|||
msgVpnName|x|||
remoteAuthenticationBasicPassword||x||x
remoteAuthenticationClientCertContent||x||x
remoteAuthenticationClientCertPassword||x||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @return BridgeApiApiGetMsgVpnBridgeRequest
*/
func (a *BridgeApiService) GetMsgVpnBridge(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string) BridgeApiApiGetMsgVpnBridgeRequest {
	return BridgeApiApiGetMsgVpnBridgeRequest{
		ApiService:          a,
		ctx:                 ctx,
		msgVpnName:          msgVpnName,
		bridgeName:          bridgeName,
		bridgeVirtualRouter: bridgeVirtualRouter,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeResponse
func (a *BridgeApiService) GetMsgVpnBridgeExecute(r BridgeApiApiGetMsgVpnBridgeRequest) (*MsgVpnBridgeResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.GetMsgVpnBridge")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnRequest struct {
	ctx                   context.Context
	ApiService            *BridgeApiService
	msgVpnName            string
	bridgeName            string
	bridgeVirtualRouter   string
	remoteMsgVpnName      string
	remoteMsgVpnLocation  string
	remoteMsgVpnInterface string
	opaquePassword        *string
	select_               *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnRequest) OpaquePassword(opaquePassword string) BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnRequest) Select_(select_ []string) BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnRequest) Execute() (*MsgVpnBridgeRemoteMsgVpnResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnBridgeRemoteMsgVpnExecute(r)
}

/*
GetMsgVpnBridgeRemoteMsgVpn Get a Remote Message VPN object.

Get a Remote Message VPN object.

The Remote Message VPN is the Message VPN that the Bridge connects to.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
bridgeName|x|||
bridgeVirtualRouter|x|||
msgVpnName|x|||
password||x||x
remoteMsgVpnInterface|x|||
remoteMsgVpnLocation|x|||
remoteMsgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @param remoteMsgVpnName The name of the remote Message VPN.
 @param remoteMsgVpnLocation The location of the remote Message VPN as either an FQDN with port, IP address with port, or virtual router name (starting with \"v:\").
 @param remoteMsgVpnInterface The physical interface on the local Message VPN host for connecting to the remote Message VPN. By default, an interface is chosen automatically (recommended), but if specified, `remoteMsgVpnLocation` must not be a virtual router name.
 @return BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnRequest
*/
func (a *BridgeApiService) GetMsgVpnBridgeRemoteMsgVpn(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string, remoteMsgVpnName string, remoteMsgVpnLocation string, remoteMsgVpnInterface string) BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnRequest {
	return BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnRequest{
		ApiService:            a,
		ctx:                   ctx,
		msgVpnName:            msgVpnName,
		bridgeName:            bridgeName,
		bridgeVirtualRouter:   bridgeVirtualRouter,
		remoteMsgVpnName:      remoteMsgVpnName,
		remoteMsgVpnLocation:  remoteMsgVpnLocation,
		remoteMsgVpnInterface: remoteMsgVpnInterface,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeRemoteMsgVpnResponse
func (a *BridgeApiService) GetMsgVpnBridgeRemoteMsgVpnExecute(r BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnRequest) (*MsgVpnBridgeRemoteMsgVpnResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeRemoteMsgVpnResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.GetMsgVpnBridgeRemoteMsgVpn")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteMsgVpnName"+"}", url.PathEscape(parameterToString(r.remoteMsgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteMsgVpnLocation"+"}", url.PathEscape(parameterToString(r.remoteMsgVpnLocation, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteMsgVpnInterface"+"}", url.PathEscape(parameterToString(r.remoteMsgVpnInterface, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnsRequest struct {
	ctx                 context.Context
	ApiService          *BridgeApiService
	msgVpnName          string
	bridgeName          string
	bridgeVirtualRouter string
	opaquePassword      *string
	where               *[]string
	select_             *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnsRequest) OpaquePassword(opaquePassword string) BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnsRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnsRequest) Where(where []string) BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnsRequest {
	r.where = &where
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnsRequest) Select_(select_ []string) BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnsRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnsRequest) Execute() (*MsgVpnBridgeRemoteMsgVpnsResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnBridgeRemoteMsgVpnsExecute(r)
}

/*
GetMsgVpnBridgeRemoteMsgVpns Get a list of Remote Message VPN objects.

Get a list of Remote Message VPN objects.

The Remote Message VPN is the Message VPN that the Bridge connects to.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
bridgeName|x|||
bridgeVirtualRouter|x|||
msgVpnName|x|||
password||x||x
remoteMsgVpnInterface|x|||
remoteMsgVpnLocation|x|||
remoteMsgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @return BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnsRequest
*/
func (a *BridgeApiService) GetMsgVpnBridgeRemoteMsgVpns(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string) BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnsRequest {
	return BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnsRequest{
		ApiService:          a,
		ctx:                 ctx,
		msgVpnName:          msgVpnName,
		bridgeName:          bridgeName,
		bridgeVirtualRouter: bridgeVirtualRouter,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeRemoteMsgVpnsResponse
func (a *BridgeApiService) GetMsgVpnBridgeRemoteMsgVpnsExecute(r BridgeApiApiGetMsgVpnBridgeRemoteMsgVpnsRequest) (*MsgVpnBridgeRemoteMsgVpnsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeRemoteMsgVpnsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.GetMsgVpnBridgeRemoteMsgVpns")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.where != nil {
		localVarQueryParams.Add("where", parameterToString(*r.where, "csv"))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionRequest struct {
	ctx                     context.Context
	ApiService              *BridgeApiService
	msgVpnName              string
	bridgeName              string
	bridgeVirtualRouter     string
	remoteSubscriptionTopic string
	opaquePassword          *string
	select_                 *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionRequest) OpaquePassword(opaquePassword string) BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionRequest) Select_(select_ []string) BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionRequest) Execute() (*MsgVpnBridgeRemoteSubscriptionResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnBridgeRemoteSubscriptionExecute(r)
}

/*
GetMsgVpnBridgeRemoteSubscription Get a Remote Subscription object.

Get a Remote Subscription object.

A Remote Subscription is a topic subscription used by the Message VPN Bridge to attract messages from the remote message broker.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
bridgeName|x|||
bridgeVirtualRouter|x|||
msgVpnName|x|||
remoteSubscriptionTopic|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @param remoteSubscriptionTopic The topic of the Bridge remote subscription.
 @return BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionRequest
*/
func (a *BridgeApiService) GetMsgVpnBridgeRemoteSubscription(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string, remoteSubscriptionTopic string) BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionRequest {
	return BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionRequest{
		ApiService:              a,
		ctx:                     ctx,
		msgVpnName:              msgVpnName,
		bridgeName:              bridgeName,
		bridgeVirtualRouter:     bridgeVirtualRouter,
		remoteSubscriptionTopic: remoteSubscriptionTopic,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeRemoteSubscriptionResponse
func (a *BridgeApiService) GetMsgVpnBridgeRemoteSubscriptionExecute(r BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionRequest) (*MsgVpnBridgeRemoteSubscriptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeRemoteSubscriptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.GetMsgVpnBridgeRemoteSubscription")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions/{remoteSubscriptionTopic}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteSubscriptionTopic"+"}", url.PathEscape(parameterToString(r.remoteSubscriptionTopic, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest struct {
	ctx                 context.Context
	ApiService          *BridgeApiService
	msgVpnName          string
	bridgeName          string
	bridgeVirtualRouter string
	count               *int32
	cursor              *string
	opaquePassword      *string
	where               *[]string
	select_             *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest) Count(count int32) BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest {
	r.count = &count
	return r
}

// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest) Cursor(cursor string) BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest {
	r.cursor = &cursor
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest) OpaquePassword(opaquePassword string) BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest) Where(where []string) BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest {
	r.where = &where
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest) Select_(select_ []string) BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest) Execute() (*MsgVpnBridgeRemoteSubscriptionsResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnBridgeRemoteSubscriptionsExecute(r)
}

/*
GetMsgVpnBridgeRemoteSubscriptions Get a list of Remote Subscription objects.

Get a list of Remote Subscription objects.

A Remote Subscription is a topic subscription used by the Message VPN Bridge to attract messages from the remote message broker.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
bridgeName|x|||
bridgeVirtualRouter|x|||
msgVpnName|x|||
remoteSubscriptionTopic|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @return BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest
*/
func (a *BridgeApiService) GetMsgVpnBridgeRemoteSubscriptions(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string) BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest {
	return BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest{
		ApiService:          a,
		ctx:                 ctx,
		msgVpnName:          msgVpnName,
		bridgeName:          bridgeName,
		bridgeVirtualRouter: bridgeVirtualRouter,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeRemoteSubscriptionsResponse
func (a *BridgeApiService) GetMsgVpnBridgeRemoteSubscriptionsExecute(r BridgeApiApiGetMsgVpnBridgeRemoteSubscriptionsRequest) (*MsgVpnBridgeRemoteSubscriptionsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeRemoteSubscriptionsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.GetMsgVpnBridgeRemoteSubscriptions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteSubscriptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.count != nil {
		localVarQueryParams.Add("count", parameterToString(*r.count, ""))
	}
	if r.cursor != nil {
		localVarQueryParams.Add("cursor", parameterToString(*r.cursor, ""))
	}
	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.where != nil {
		localVarQueryParams.Add("where", parameterToString(*r.where, "csv"))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNameRequest struct {
	ctx                  context.Context
	ApiService           *BridgeApiService
	msgVpnName           string
	bridgeName           string
	bridgeVirtualRouter  string
	tlsTrustedCommonName string
	opaquePassword       *string
	select_              *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNameRequest) OpaquePassword(opaquePassword string) BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNameRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNameRequest) Select_(select_ []string) BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNameRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNameRequest) Execute() (*MsgVpnBridgeTlsTrustedCommonNameResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnBridgeTlsTrustedCommonNameExecute(r)
}

/*
GetMsgVpnBridgeTlsTrustedCommonName Get a Trusted Common Name object.

Get a Trusted Common Name object.

The Trusted Common Names for the Bridge are used by encrypted transports to verify the name in the certificate presented by the remote node. They must include the common name of the remote node's server certificate or client certificate, depending upon the initiator of the connection.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
bridgeName|x||x|
bridgeVirtualRouter|x||x|
msgVpnName|x||x|
tlsTrustedCommonName|x||x|



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been deprecated since 2.18. Common Name validation has been replaced by Server Certificate Name validation.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @param tlsTrustedCommonName The expected trusted common name of the remote certificate.
 @return BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNameRequest

Deprecated
*/
func (a *BridgeApiService) GetMsgVpnBridgeTlsTrustedCommonName(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string, tlsTrustedCommonName string) BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNameRequest {
	return BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNameRequest{
		ApiService:           a,
		ctx:                  ctx,
		msgVpnName:           msgVpnName,
		bridgeName:           bridgeName,
		bridgeVirtualRouter:  bridgeVirtualRouter,
		tlsTrustedCommonName: tlsTrustedCommonName,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeTlsTrustedCommonNameResponse
// Deprecated
func (a *BridgeApiService) GetMsgVpnBridgeTlsTrustedCommonNameExecute(r BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNameRequest) (*MsgVpnBridgeTlsTrustedCommonNameResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeTlsTrustedCommonNameResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.GetMsgVpnBridgeTlsTrustedCommonName")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames/{tlsTrustedCommonName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"tlsTrustedCommonName"+"}", url.PathEscape(parameterToString(r.tlsTrustedCommonName, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNamesRequest struct {
	ctx                 context.Context
	ApiService          *BridgeApiService
	msgVpnName          string
	bridgeName          string
	bridgeVirtualRouter string
	opaquePassword      *string
	where               *[]string
	select_             *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNamesRequest) OpaquePassword(opaquePassword string) BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNamesRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNamesRequest) Where(where []string) BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNamesRequest {
	r.where = &where
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNamesRequest) Select_(select_ []string) BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNamesRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNamesRequest) Execute() (*MsgVpnBridgeTlsTrustedCommonNamesResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnBridgeTlsTrustedCommonNamesExecute(r)
}

/*
GetMsgVpnBridgeTlsTrustedCommonNames Get a list of Trusted Common Name objects.

Get a list of Trusted Common Name objects.

The Trusted Common Names for the Bridge are used by encrypted transports to verify the name in the certificate presented by the remote node. They must include the common name of the remote node's server certificate or client certificate, depending upon the initiator of the connection.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
bridgeName|x||x|
bridgeVirtualRouter|x||x|
msgVpnName|x||x|
tlsTrustedCommonName|x||x|



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been deprecated since 2.18. Common Name validation has been replaced by Server Certificate Name validation.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @return BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNamesRequest

Deprecated
*/
func (a *BridgeApiService) GetMsgVpnBridgeTlsTrustedCommonNames(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string) BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNamesRequest {
	return BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNamesRequest{
		ApiService:          a,
		ctx:                 ctx,
		msgVpnName:          msgVpnName,
		bridgeName:          bridgeName,
		bridgeVirtualRouter: bridgeVirtualRouter,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeTlsTrustedCommonNamesResponse
// Deprecated
func (a *BridgeApiService) GetMsgVpnBridgeTlsTrustedCommonNamesExecute(r BridgeApiApiGetMsgVpnBridgeTlsTrustedCommonNamesRequest) (*MsgVpnBridgeTlsTrustedCommonNamesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeTlsTrustedCommonNamesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.GetMsgVpnBridgeTlsTrustedCommonNames")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/tlsTrustedCommonNames"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.where != nil {
		localVarQueryParams.Add("where", parameterToString(*r.where, "csv"))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiGetMsgVpnBridgesRequest struct {
	ctx            context.Context
	ApiService     *BridgeApiService
	msgVpnName     string
	count          *int32
	cursor         *string
	opaquePassword *string
	where          *[]string
	select_        *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgesRequest) Count(count int32) BridgeApiApiGetMsgVpnBridgesRequest {
	r.count = &count
	return r
}

// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgesRequest) Cursor(cursor string) BridgeApiApiGetMsgVpnBridgesRequest {
	r.cursor = &cursor
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgesRequest) OpaquePassword(opaquePassword string) BridgeApiApiGetMsgVpnBridgesRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgesRequest) Where(where []string) BridgeApiApiGetMsgVpnBridgesRequest {
	r.where = &where
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiGetMsgVpnBridgesRequest) Select_(select_ []string) BridgeApiApiGetMsgVpnBridgesRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiGetMsgVpnBridgesRequest) Execute() (*MsgVpnBridgesResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnBridgesExecute(r)
}

/*
GetMsgVpnBridges Get a list of Bridge objects.

Get a list of Bridge objects.

Bridges can be used to link two Message VPNs so that messages published to one Message VPN that match the topic subscriptions set for the bridge are also delivered to the linked Message VPN.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
bridgeName|x|||
bridgeVirtualRouter|x|||
msgVpnName|x|||
remoteAuthenticationBasicPassword||x||x
remoteAuthenticationClientCertContent||x||x
remoteAuthenticationClientCertPassword||x||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @return BridgeApiApiGetMsgVpnBridgesRequest
*/
func (a *BridgeApiService) GetMsgVpnBridges(ctx context.Context, msgVpnName string) BridgeApiApiGetMsgVpnBridgesRequest {
	return BridgeApiApiGetMsgVpnBridgesRequest{
		ApiService: a,
		ctx:        ctx,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return MsgVpnBridgesResponse
func (a *BridgeApiService) GetMsgVpnBridgesExecute(r BridgeApiApiGetMsgVpnBridgesRequest) (*MsgVpnBridgesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.GetMsgVpnBridges")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.count != nil {
		localVarQueryParams.Add("count", parameterToString(*r.count, ""))
	}
	if r.cursor != nil {
		localVarQueryParams.Add("cursor", parameterToString(*r.cursor, ""))
	}
	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.where != nil {
		localVarQueryParams.Add("where", parameterToString(*r.where, "csv"))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiReplaceMsgVpnBridgeRequest struct {
	ctx                 context.Context
	ApiService          *BridgeApiService
	msgVpnName          string
	bridgeName          string
	bridgeVirtualRouter string
	body                *MsgVpnBridge
	opaquePassword      *string
	select_             *[]string
}

// The Bridge object&#39;s attributes.
func (r BridgeApiApiReplaceMsgVpnBridgeRequest) Body(body MsgVpnBridge) BridgeApiApiReplaceMsgVpnBridgeRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiReplaceMsgVpnBridgeRequest) OpaquePassword(opaquePassword string) BridgeApiApiReplaceMsgVpnBridgeRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiReplaceMsgVpnBridgeRequest) Select_(select_ []string) BridgeApiApiReplaceMsgVpnBridgeRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiReplaceMsgVpnBridgeRequest) Execute() (*MsgVpnBridgeResponse, *http.Response, error) {
	return r.ApiService.ReplaceMsgVpnBridgeExecute(r)
}

/*
ReplaceMsgVpnBridge Replace a Bridge object.

Replace a Bridge object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

Bridges can be used to link two Message VPNs so that messages published to one Message VPN that match the topic subscriptions set for the bridge are also delivered to the linked Message VPN.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
bridgeName|x|x||||
bridgeVirtualRouter|x|x||||
maxTtl||||x||
msgVpnName|x|x||||
remoteAuthenticationBasicClientUsername||||x||
remoteAuthenticationBasicPassword|||x|x||x
remoteAuthenticationClientCertContent|||x|x||x
remoteAuthenticationClientCertPassword|||x|x||
remoteAuthenticationScheme||||x||
remoteDeliverToOnePriority||||x||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
MsgVpnBridge|remoteAuthenticationBasicClientUsername|remoteAuthenticationBasicPassword|
MsgVpnBridge|remoteAuthenticationBasicPassword|remoteAuthenticationBasicClientUsername|
MsgVpnBridge|remoteAuthenticationClientCertPassword|remoteAuthenticationClientCertContent|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @return BridgeApiApiReplaceMsgVpnBridgeRequest
*/
func (a *BridgeApiService) ReplaceMsgVpnBridge(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string) BridgeApiApiReplaceMsgVpnBridgeRequest {
	return BridgeApiApiReplaceMsgVpnBridgeRequest{
		ApiService:          a,
		ctx:                 ctx,
		msgVpnName:          msgVpnName,
		bridgeName:          bridgeName,
		bridgeVirtualRouter: bridgeVirtualRouter,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeResponse
func (a *BridgeApiService) ReplaceMsgVpnBridgeExecute(r BridgeApiApiReplaceMsgVpnBridgeRequest) (*MsgVpnBridgeResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.ReplaceMsgVpnBridge")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.body == nil {
		return localVarReturnValue, nil, reportError("body is required and must be specified")
	}

	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.body
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiReplaceMsgVpnBridgeRemoteMsgVpnRequest struct {
	ctx                   context.Context
	ApiService            *BridgeApiService
	msgVpnName            string
	bridgeName            string
	bridgeVirtualRouter   string
	remoteMsgVpnName      string
	remoteMsgVpnLocation  string
	remoteMsgVpnInterface string
	body                  *MsgVpnBridgeRemoteMsgVpn
	opaquePassword        *string
	select_               *[]string
}

// The Remote Message VPN object&#39;s attributes.
func (r BridgeApiApiReplaceMsgVpnBridgeRemoteMsgVpnRequest) Body(body MsgVpnBridgeRemoteMsgVpn) BridgeApiApiReplaceMsgVpnBridgeRemoteMsgVpnRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiReplaceMsgVpnBridgeRemoteMsgVpnRequest) OpaquePassword(opaquePassword string) BridgeApiApiReplaceMsgVpnBridgeRemoteMsgVpnRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiReplaceMsgVpnBridgeRemoteMsgVpnRequest) Select_(select_ []string) BridgeApiApiReplaceMsgVpnBridgeRemoteMsgVpnRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiReplaceMsgVpnBridgeRemoteMsgVpnRequest) Execute() (*MsgVpnBridgeRemoteMsgVpnResponse, *http.Response, error) {
	return r.ApiService.ReplaceMsgVpnBridgeRemoteMsgVpnExecute(r)
}

/*
ReplaceMsgVpnBridgeRemoteMsgVpn Replace a Remote Message VPN object.

Replace a Remote Message VPN object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

The Remote Message VPN is the Message VPN that the Bridge connects to.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
bridgeName|x|x||||
bridgeVirtualRouter|x|x||||
clientUsername||||x||
compressedDataEnabled||||x||
egressFlowWindowSize||||x||
msgVpnName|x|x||||
password|||x|x||x
remoteMsgVpnInterface|x|x||||
remoteMsgVpnLocation|x|x||||
remoteMsgVpnName|x|x||||
tlsEnabled||||x||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
MsgVpnBridgeRemoteMsgVpn|clientUsername|password|
MsgVpnBridgeRemoteMsgVpn|password|clientUsername|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @param remoteMsgVpnName The name of the remote Message VPN.
 @param remoteMsgVpnLocation The location of the remote Message VPN as either an FQDN with port, IP address with port, or virtual router name (starting with \"v:\").
 @param remoteMsgVpnInterface The physical interface on the local Message VPN host for connecting to the remote Message VPN. By default, an interface is chosen automatically (recommended), but if specified, `remoteMsgVpnLocation` must not be a virtual router name.
 @return BridgeApiApiReplaceMsgVpnBridgeRemoteMsgVpnRequest
*/
func (a *BridgeApiService) ReplaceMsgVpnBridgeRemoteMsgVpn(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string, remoteMsgVpnName string, remoteMsgVpnLocation string, remoteMsgVpnInterface string) BridgeApiApiReplaceMsgVpnBridgeRemoteMsgVpnRequest {
	return BridgeApiApiReplaceMsgVpnBridgeRemoteMsgVpnRequest{
		ApiService:            a,
		ctx:                   ctx,
		msgVpnName:            msgVpnName,
		bridgeName:            bridgeName,
		bridgeVirtualRouter:   bridgeVirtualRouter,
		remoteMsgVpnName:      remoteMsgVpnName,
		remoteMsgVpnLocation:  remoteMsgVpnLocation,
		remoteMsgVpnInterface: remoteMsgVpnInterface,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeRemoteMsgVpnResponse
func (a *BridgeApiService) ReplaceMsgVpnBridgeRemoteMsgVpnExecute(r BridgeApiApiReplaceMsgVpnBridgeRemoteMsgVpnRequest) (*MsgVpnBridgeRemoteMsgVpnResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeRemoteMsgVpnResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.ReplaceMsgVpnBridgeRemoteMsgVpn")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteMsgVpnName"+"}", url.PathEscape(parameterToString(r.remoteMsgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteMsgVpnLocation"+"}", url.PathEscape(parameterToString(r.remoteMsgVpnLocation, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteMsgVpnInterface"+"}", url.PathEscape(parameterToString(r.remoteMsgVpnInterface, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.body == nil {
		return localVarReturnValue, nil, reportError("body is required and must be specified")
	}

	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.body
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiUpdateMsgVpnBridgeRequest struct {
	ctx                 context.Context
	ApiService          *BridgeApiService
	msgVpnName          string
	bridgeName          string
	bridgeVirtualRouter string
	body                *MsgVpnBridge
	opaquePassword      *string
	select_             *[]string
}

// The Bridge object&#39;s attributes.
func (r BridgeApiApiUpdateMsgVpnBridgeRequest) Body(body MsgVpnBridge) BridgeApiApiUpdateMsgVpnBridgeRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiUpdateMsgVpnBridgeRequest) OpaquePassword(opaquePassword string) BridgeApiApiUpdateMsgVpnBridgeRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiUpdateMsgVpnBridgeRequest) Select_(select_ []string) BridgeApiApiUpdateMsgVpnBridgeRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiUpdateMsgVpnBridgeRequest) Execute() (*MsgVpnBridgeResponse, *http.Response, error) {
	return r.ApiService.UpdateMsgVpnBridgeExecute(r)
}

/*
UpdateMsgVpnBridge Update a Bridge object.

Update a Bridge object. Any attribute missing from the request will be left unchanged.

Bridges can be used to link two Message VPNs so that messages published to one Message VPN that match the topic subscriptions set for the bridge are also delivered to the linked Message VPN.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
bridgeName|x|x||||
bridgeVirtualRouter|x|x||||
maxTtl||||x||
msgVpnName|x|x||||
remoteAuthenticationBasicClientUsername||||x||
remoteAuthenticationBasicPassword|||x|x||x
remoteAuthenticationClientCertContent|||x|x||x
remoteAuthenticationClientCertPassword|||x|x||
remoteAuthenticationScheme||||x||
remoteDeliverToOnePriority||||x||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
MsgVpnBridge|remoteAuthenticationBasicClientUsername|remoteAuthenticationBasicPassword|
MsgVpnBridge|remoteAuthenticationBasicPassword|remoteAuthenticationBasicClientUsername|
MsgVpnBridge|remoteAuthenticationClientCertPassword|remoteAuthenticationClientCertContent|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @return BridgeApiApiUpdateMsgVpnBridgeRequest
*/
func (a *BridgeApiService) UpdateMsgVpnBridge(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string) BridgeApiApiUpdateMsgVpnBridgeRequest {
	return BridgeApiApiUpdateMsgVpnBridgeRequest{
		ApiService:          a,
		ctx:                 ctx,
		msgVpnName:          msgVpnName,
		bridgeName:          bridgeName,
		bridgeVirtualRouter: bridgeVirtualRouter,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeResponse
func (a *BridgeApiService) UpdateMsgVpnBridgeExecute(r BridgeApiApiUpdateMsgVpnBridgeRequest) (*MsgVpnBridgeResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.UpdateMsgVpnBridge")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.body == nil {
		return localVarReturnValue, nil, reportError("body is required and must be specified")
	}

	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.body
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type BridgeApiApiUpdateMsgVpnBridgeRemoteMsgVpnRequest struct {
	ctx                   context.Context
	ApiService            *BridgeApiService
	msgVpnName            string
	bridgeName            string
	bridgeVirtualRouter   string
	remoteMsgVpnName      string
	remoteMsgVpnLocation  string
	remoteMsgVpnInterface string
	body                  *MsgVpnBridgeRemoteMsgVpn
	opaquePassword        *string
	select_               *[]string
}

// The Remote Message VPN object&#39;s attributes.
func (r BridgeApiApiUpdateMsgVpnBridgeRemoteMsgVpnRequest) Body(body MsgVpnBridgeRemoteMsgVpn) BridgeApiApiUpdateMsgVpnBridgeRemoteMsgVpnRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r BridgeApiApiUpdateMsgVpnBridgeRemoteMsgVpnRequest) OpaquePassword(opaquePassword string) BridgeApiApiUpdateMsgVpnBridgeRemoteMsgVpnRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r BridgeApiApiUpdateMsgVpnBridgeRemoteMsgVpnRequest) Select_(select_ []string) BridgeApiApiUpdateMsgVpnBridgeRemoteMsgVpnRequest {
	r.select_ = &select_
	return r
}

func (r BridgeApiApiUpdateMsgVpnBridgeRemoteMsgVpnRequest) Execute() (*MsgVpnBridgeRemoteMsgVpnResponse, *http.Response, error) {
	return r.ApiService.UpdateMsgVpnBridgeRemoteMsgVpnExecute(r)
}

/*
UpdateMsgVpnBridgeRemoteMsgVpn Update a Remote Message VPN object.

Update a Remote Message VPN object. Any attribute missing from the request will be left unchanged.

The Remote Message VPN is the Message VPN that the Bridge connects to.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
bridgeName|x|x||||
bridgeVirtualRouter|x|x||||
clientUsername||||x||
compressedDataEnabled||||x||
egressFlowWindowSize||||x||
msgVpnName|x|x||||
password|||x|x||x
remoteMsgVpnInterface|x|x||||
remoteMsgVpnLocation|x|x||||
remoteMsgVpnName|x|x||||
tlsEnabled||||x||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
MsgVpnBridgeRemoteMsgVpn|clientUsername|password|
MsgVpnBridgeRemoteMsgVpn|password|clientUsername|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param bridgeName The name of the Bridge.
 @param bridgeVirtualRouter The virtual router of the Bridge.
 @param remoteMsgVpnName The name of the remote Message VPN.
 @param remoteMsgVpnLocation The location of the remote Message VPN as either an FQDN with port, IP address with port, or virtual router name (starting with \"v:\").
 @param remoteMsgVpnInterface The physical interface on the local Message VPN host for connecting to the remote Message VPN. By default, an interface is chosen automatically (recommended), but if specified, `remoteMsgVpnLocation` must not be a virtual router name.
 @return BridgeApiApiUpdateMsgVpnBridgeRemoteMsgVpnRequest
*/
func (a *BridgeApiService) UpdateMsgVpnBridgeRemoteMsgVpn(ctx context.Context, msgVpnName string, bridgeName string, bridgeVirtualRouter string, remoteMsgVpnName string, remoteMsgVpnLocation string, remoteMsgVpnInterface string) BridgeApiApiUpdateMsgVpnBridgeRemoteMsgVpnRequest {
	return BridgeApiApiUpdateMsgVpnBridgeRemoteMsgVpnRequest{
		ApiService:            a,
		ctx:                   ctx,
		msgVpnName:            msgVpnName,
		bridgeName:            bridgeName,
		bridgeVirtualRouter:   bridgeVirtualRouter,
		remoteMsgVpnName:      remoteMsgVpnName,
		remoteMsgVpnLocation:  remoteMsgVpnLocation,
		remoteMsgVpnInterface: remoteMsgVpnInterface,
	}
}

// Execute executes the request
//  @return MsgVpnBridgeRemoteMsgVpnResponse
func (a *BridgeApiService) UpdateMsgVpnBridgeRemoteMsgVpnExecute(r BridgeApiApiUpdateMsgVpnBridgeRemoteMsgVpnRequest) (*MsgVpnBridgeRemoteMsgVpnResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnBridgeRemoteMsgVpnResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BridgeApiService.UpdateMsgVpnBridgeRemoteMsgVpn")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/bridges/{bridgeName},{bridgeVirtualRouter}/remoteMsgVpns/{remoteMsgVpnName},{remoteMsgVpnLocation},{remoteMsgVpnInterface}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeName"+"}", url.PathEscape(parameterToString(r.bridgeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"bridgeVirtualRouter"+"}", url.PathEscape(parameterToString(r.bridgeVirtualRouter, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteMsgVpnName"+"}", url.PathEscape(parameterToString(r.remoteMsgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteMsgVpnLocation"+"}", url.PathEscape(parameterToString(r.remoteMsgVpnLocation, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteMsgVpnInterface"+"}", url.PathEscape(parameterToString(r.remoteMsgVpnInterface, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.body == nil {
		return localVarReturnValue, nil, reportError("body is required and must be specified")
	}

	if r.opaquePassword != nil {
		localVarQueryParams.Add("opaquePassword", parameterToString(*r.opaquePassword, ""))
	}
	if r.select_ != nil {
		localVarQueryParams.Add("select", parameterToString(*r.select_, "csv"))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.body
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v SempMetaOnlyResponse
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
