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

// ClientProfileApiService ClientProfileApi service
type ClientProfileApiService service

type ClientProfileApiApiCreateMsgVpnClientProfileRequest struct {
	ctx context.Context
	ApiService *ClientProfileApiService
	msgVpnName string
	body *MsgVpnClientProfile
	opaquePassword *string
	select_ *[]string
}

// The Client Profile object&#39;s attributes.
func (r ClientProfileApiApiCreateMsgVpnClientProfileRequest) Body(body MsgVpnClientProfile) ClientProfileApiApiCreateMsgVpnClientProfileRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r ClientProfileApiApiCreateMsgVpnClientProfileRequest) OpaquePassword(opaquePassword string) ClientProfileApiApiCreateMsgVpnClientProfileRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r ClientProfileApiApiCreateMsgVpnClientProfileRequest) Select_(select_ []string) ClientProfileApiApiCreateMsgVpnClientProfileRequest {
	r.select_ = &select_
	return r
}

func (r ClientProfileApiApiCreateMsgVpnClientProfileRequest) Execute() (*MsgVpnClientProfileResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnClientProfileExecute(r)
}

/*
CreateMsgVpnClientProfile Create a Client Profile object.

Create a Client Profile object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

Client Profiles are used to assign common configuration properties to clients that have been successfully authorized.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
allowCutThroughForwardingEnabled|||||x|
apiQueueManagementCopyFromOnCreateName|||||x|
apiTopicEndpointManagementCopyFromOnCreateName|||||x|
clientProfileName|x|x||||
msgVpnName|x||x|||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
EventThreshold|clearPercent|setPercent|clearValue, setValue
EventThreshold|clearValue|setValue|clearPercent, setPercent
EventThreshold|setPercent|clearPercent|clearValue, setValue
EventThreshold|setValue|clearValue|clearPercent, setPercent
EventThresholdByPercent|clearPercent|setPercent|
EventThresholdByPercent|setPercent|clearPercent|



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @return ClientProfileApiApiCreateMsgVpnClientProfileRequest
*/
func (a *ClientProfileApiService) CreateMsgVpnClientProfile(ctx context.Context, msgVpnName string) ClientProfileApiApiCreateMsgVpnClientProfileRequest {
	return ClientProfileApiApiCreateMsgVpnClientProfileRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return MsgVpnClientProfileResponse
func (a *ClientProfileApiService) CreateMsgVpnClientProfileExecute(r ClientProfileApiApiCreateMsgVpnClientProfileRequest) (*MsgVpnClientProfileResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnClientProfileResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientProfileApiService.CreateMsgVpnClientProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/clientProfiles"
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

type ClientProfileApiApiDeleteMsgVpnClientProfileRequest struct {
	ctx context.Context
	ApiService *ClientProfileApiService
	msgVpnName string
	clientProfileName string
}


func (r ClientProfileApiApiDeleteMsgVpnClientProfileRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnClientProfileExecute(r)
}

/*
DeleteMsgVpnClientProfile Delete a Client Profile object.

Delete a Client Profile object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

Client Profiles are used to assign common configuration properties to clients that have been successfully authorized.

A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param clientProfileName The name of the Client Profile.
 @return ClientProfileApiApiDeleteMsgVpnClientProfileRequest
*/
func (a *ClientProfileApiService) DeleteMsgVpnClientProfile(ctx context.Context, msgVpnName string, clientProfileName string) ClientProfileApiApiDeleteMsgVpnClientProfileRequest {
	return ClientProfileApiApiDeleteMsgVpnClientProfileRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		clientProfileName: clientProfileName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *ClientProfileApiService) DeleteMsgVpnClientProfileExecute(r ClientProfileApiApiDeleteMsgVpnClientProfileRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientProfileApiService.DeleteMsgVpnClientProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/clientProfiles/{clientProfileName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clientProfileName"+"}", url.PathEscape(parameterToString(r.clientProfileName, "")), -1)

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

type ClientProfileApiApiGetMsgVpnClientProfileRequest struct {
	ctx context.Context
	ApiService *ClientProfileApiService
	msgVpnName string
	clientProfileName string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r ClientProfileApiApiGetMsgVpnClientProfileRequest) OpaquePassword(opaquePassword string) ClientProfileApiApiGetMsgVpnClientProfileRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r ClientProfileApiApiGetMsgVpnClientProfileRequest) Select_(select_ []string) ClientProfileApiApiGetMsgVpnClientProfileRequest {
	r.select_ = &select_
	return r
}

func (r ClientProfileApiApiGetMsgVpnClientProfileRequest) Execute() (*MsgVpnClientProfileResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnClientProfileExecute(r)
}

/*
GetMsgVpnClientProfile Get a Client Profile object.

Get a Client Profile object.

Client Profiles are used to assign common configuration properties to clients that have been successfully authorized.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
allowCutThroughForwardingEnabled|||x|
apiQueueManagementCopyFromOnCreateName|||x|
apiTopicEndpointManagementCopyFromOnCreateName|||x|
clientProfileName|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param clientProfileName The name of the Client Profile.
 @return ClientProfileApiApiGetMsgVpnClientProfileRequest
*/
func (a *ClientProfileApiService) GetMsgVpnClientProfile(ctx context.Context, msgVpnName string, clientProfileName string) ClientProfileApiApiGetMsgVpnClientProfileRequest {
	return ClientProfileApiApiGetMsgVpnClientProfileRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		clientProfileName: clientProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnClientProfileResponse
func (a *ClientProfileApiService) GetMsgVpnClientProfileExecute(r ClientProfileApiApiGetMsgVpnClientProfileRequest) (*MsgVpnClientProfileResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnClientProfileResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientProfileApiService.GetMsgVpnClientProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/clientProfiles/{clientProfileName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clientProfileName"+"}", url.PathEscape(parameterToString(r.clientProfileName, "")), -1)

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

type ClientProfileApiApiGetMsgVpnClientProfilesRequest struct {
	ctx context.Context
	ApiService *ClientProfileApiService
	msgVpnName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r ClientProfileApiApiGetMsgVpnClientProfilesRequest) Count(count int32) ClientProfileApiApiGetMsgVpnClientProfilesRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r ClientProfileApiApiGetMsgVpnClientProfilesRequest) Cursor(cursor string) ClientProfileApiApiGetMsgVpnClientProfilesRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r ClientProfileApiApiGetMsgVpnClientProfilesRequest) OpaquePassword(opaquePassword string) ClientProfileApiApiGetMsgVpnClientProfilesRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r ClientProfileApiApiGetMsgVpnClientProfilesRequest) Where(where []string) ClientProfileApiApiGetMsgVpnClientProfilesRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r ClientProfileApiApiGetMsgVpnClientProfilesRequest) Select_(select_ []string) ClientProfileApiApiGetMsgVpnClientProfilesRequest {
	r.select_ = &select_
	return r
}

func (r ClientProfileApiApiGetMsgVpnClientProfilesRequest) Execute() (*MsgVpnClientProfilesResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnClientProfilesExecute(r)
}

/*
GetMsgVpnClientProfiles Get a list of Client Profile objects.

Get a list of Client Profile objects.

Client Profiles are used to assign common configuration properties to clients that have been successfully authorized.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
allowCutThroughForwardingEnabled|||x|
apiQueueManagementCopyFromOnCreateName|||x|
apiTopicEndpointManagementCopyFromOnCreateName|||x|
clientProfileName|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @return ClientProfileApiApiGetMsgVpnClientProfilesRequest
*/
func (a *ClientProfileApiService) GetMsgVpnClientProfiles(ctx context.Context, msgVpnName string) ClientProfileApiApiGetMsgVpnClientProfilesRequest {
	return ClientProfileApiApiGetMsgVpnClientProfilesRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return MsgVpnClientProfilesResponse
func (a *ClientProfileApiService) GetMsgVpnClientProfilesExecute(r ClientProfileApiApiGetMsgVpnClientProfilesRequest) (*MsgVpnClientProfilesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnClientProfilesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientProfileApiService.GetMsgVpnClientProfiles")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/clientProfiles"
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

type ClientProfileApiApiReplaceMsgVpnClientProfileRequest struct {
	ctx context.Context
	ApiService *ClientProfileApiService
	msgVpnName string
	clientProfileName string
	body *MsgVpnClientProfile
	opaquePassword *string
	select_ *[]string
}

// The Client Profile object&#39;s attributes.
func (r ClientProfileApiApiReplaceMsgVpnClientProfileRequest) Body(body MsgVpnClientProfile) ClientProfileApiApiReplaceMsgVpnClientProfileRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r ClientProfileApiApiReplaceMsgVpnClientProfileRequest) OpaquePassword(opaquePassword string) ClientProfileApiApiReplaceMsgVpnClientProfileRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r ClientProfileApiApiReplaceMsgVpnClientProfileRequest) Select_(select_ []string) ClientProfileApiApiReplaceMsgVpnClientProfileRequest {
	r.select_ = &select_
	return r
}

func (r ClientProfileApiApiReplaceMsgVpnClientProfileRequest) Execute() (*MsgVpnClientProfileResponse, *http.Response, error) {
	return r.ApiService.ReplaceMsgVpnClientProfileExecute(r)
}

/*
ReplaceMsgVpnClientProfile Replace a Client Profile object.

Replace a Client Profile object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

Client Profiles are used to assign common configuration properties to clients that have been successfully authorized.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
allowCutThroughForwardingEnabled|||||x|
apiQueueManagementCopyFromOnCreateName|||||x|
apiTopicEndpointManagementCopyFromOnCreateName|||||x|
clientProfileName|x|x||||
msgVpnName|x|x||||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
EventThreshold|clearPercent|setPercent|clearValue, setValue
EventThreshold|clearValue|setValue|clearPercent, setPercent
EventThreshold|setPercent|clearPercent|clearValue, setValue
EventThreshold|setValue|clearValue|clearPercent, setPercent
EventThresholdByPercent|clearPercent|setPercent|
EventThresholdByPercent|setPercent|clearPercent|



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param clientProfileName The name of the Client Profile.
 @return ClientProfileApiApiReplaceMsgVpnClientProfileRequest
*/
func (a *ClientProfileApiService) ReplaceMsgVpnClientProfile(ctx context.Context, msgVpnName string, clientProfileName string) ClientProfileApiApiReplaceMsgVpnClientProfileRequest {
	return ClientProfileApiApiReplaceMsgVpnClientProfileRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		clientProfileName: clientProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnClientProfileResponse
func (a *ClientProfileApiService) ReplaceMsgVpnClientProfileExecute(r ClientProfileApiApiReplaceMsgVpnClientProfileRequest) (*MsgVpnClientProfileResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnClientProfileResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientProfileApiService.ReplaceMsgVpnClientProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/clientProfiles/{clientProfileName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clientProfileName"+"}", url.PathEscape(parameterToString(r.clientProfileName, "")), -1)

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

type ClientProfileApiApiUpdateMsgVpnClientProfileRequest struct {
	ctx context.Context
	ApiService *ClientProfileApiService
	msgVpnName string
	clientProfileName string
	body *MsgVpnClientProfile
	opaquePassword *string
	select_ *[]string
}

// The Client Profile object&#39;s attributes.
func (r ClientProfileApiApiUpdateMsgVpnClientProfileRequest) Body(body MsgVpnClientProfile) ClientProfileApiApiUpdateMsgVpnClientProfileRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r ClientProfileApiApiUpdateMsgVpnClientProfileRequest) OpaquePassword(opaquePassword string) ClientProfileApiApiUpdateMsgVpnClientProfileRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r ClientProfileApiApiUpdateMsgVpnClientProfileRequest) Select_(select_ []string) ClientProfileApiApiUpdateMsgVpnClientProfileRequest {
	r.select_ = &select_
	return r
}

func (r ClientProfileApiApiUpdateMsgVpnClientProfileRequest) Execute() (*MsgVpnClientProfileResponse, *http.Response, error) {
	return r.ApiService.UpdateMsgVpnClientProfileExecute(r)
}

/*
UpdateMsgVpnClientProfile Update a Client Profile object.

Update a Client Profile object. Any attribute missing from the request will be left unchanged.

Client Profiles are used to assign common configuration properties to clients that have been successfully authorized.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
allowCutThroughForwardingEnabled|||||x|
apiQueueManagementCopyFromOnCreateName|||||x|
apiTopicEndpointManagementCopyFromOnCreateName|||||x|
clientProfileName|x|x||||
msgVpnName|x|x||||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
EventThreshold|clearPercent|setPercent|clearValue, setValue
EventThreshold|clearValue|setValue|clearPercent, setPercent
EventThreshold|setPercent|clearPercent|clearValue, setValue
EventThreshold|setValue|clearValue|clearPercent, setPercent
EventThresholdByPercent|clearPercent|setPercent|
EventThresholdByPercent|setPercent|clearPercent|



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param clientProfileName The name of the Client Profile.
 @return ClientProfileApiApiUpdateMsgVpnClientProfileRequest
*/
func (a *ClientProfileApiService) UpdateMsgVpnClientProfile(ctx context.Context, msgVpnName string, clientProfileName string) ClientProfileApiApiUpdateMsgVpnClientProfileRequest {
	return ClientProfileApiApiUpdateMsgVpnClientProfileRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		clientProfileName: clientProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnClientProfileResponse
func (a *ClientProfileApiService) UpdateMsgVpnClientProfileExecute(r ClientProfileApiApiUpdateMsgVpnClientProfileRequest) (*MsgVpnClientProfileResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPatch
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnClientProfileResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientProfileApiService.UpdateMsgVpnClientProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/clientProfiles/{clientProfileName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clientProfileName"+"}", url.PathEscape(parameterToString(r.clientProfileName, "")), -1)

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
