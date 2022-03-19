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

// OauthProfileApiService OauthProfileApi service
type OauthProfileApiService service

type OauthProfileApiApiCreateOauthProfileRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	body *OauthProfile
	opaquePassword *string
	select_ *[]string
}

// The OAuth Profile object&#39;s attributes.
func (r OauthProfileApiApiCreateOauthProfileRequest) Body(body OauthProfile) OauthProfileApiApiCreateOauthProfileRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiCreateOauthProfileRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileRequest) Select_(select_ []string) OauthProfileApiApiCreateOauthProfileRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiCreateOauthProfileRequest) Execute() (*OauthProfileResponse, *http.Response, error) {
	return r.ApiService.CreateOauthProfileExecute(r)
}

/*
CreateOauthProfile Create an OAuth Profile object.

Create an OAuth Profile object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates via config-sync.

OAuth profiles specify how to securely authenticate to an OAuth provider.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
clientSecret||||x||x
oauthProfileName|x|x||||



A SEMP client authorized with a minimum access scope/level of "global/admin" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return OauthProfileApiApiCreateOauthProfileRequest
*/
func (a *OauthProfileApiService) CreateOauthProfile(ctx context.Context) OauthProfileApiApiCreateOauthProfileRequest {
	return OauthProfileApiApiCreateOauthProfileRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return OauthProfileResponse
func (a *OauthProfileApiService) CreateOauthProfileExecute(r OauthProfileApiApiCreateOauthProfileRequest) (*OauthProfileResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.CreateOauthProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles"

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

type OauthProfileApiApiCreateOauthProfileAccessLevelGroupRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	body *OauthProfileAccessLevelGroup
	opaquePassword *string
	select_ *[]string
}

// The Group Access Level object&#39;s attributes.
func (r OauthProfileApiApiCreateOauthProfileAccessLevelGroupRequest) Body(body OauthProfileAccessLevelGroup) OauthProfileApiApiCreateOauthProfileAccessLevelGroupRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileAccessLevelGroupRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiCreateOauthProfileAccessLevelGroupRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileAccessLevelGroupRequest) Select_(select_ []string) OauthProfileApiApiCreateOauthProfileAccessLevelGroupRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiCreateOauthProfileAccessLevelGroupRequest) Execute() (*OauthProfileAccessLevelGroupResponse, *http.Response, error) {
	return r.ApiService.CreateOauthProfileAccessLevelGroupExecute(r)
}

/*
CreateOauthProfileAccessLevelGroup Create a Group Access Level object.

Create a Group Access Level object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates via config-sync.

The name of a group as it exists on the OAuth server being used to authenticate SEMP users.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
groupName|x|x||||
oauthProfileName|x||x|||



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation. Requests which include the following attributes require greater access scope/level:


Attribute|Access Scope/Level
:---|:---:
globalAccessLevel|global/admin



This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiCreateOauthProfileAccessLevelGroupRequest
*/
func (a *OauthProfileApiService) CreateOauthProfileAccessLevelGroup(ctx context.Context, oauthProfileName string) OauthProfileApiApiCreateOauthProfileAccessLevelGroupRequest {
	return OauthProfileApiApiCreateOauthProfileAccessLevelGroupRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileAccessLevelGroupResponse
func (a *OauthProfileApiService) CreateOauthProfileAccessLevelGroupExecute(r OauthProfileApiApiCreateOauthProfileAccessLevelGroupRequest) (*OauthProfileAccessLevelGroupResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileAccessLevelGroupResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.CreateOauthProfileAccessLevelGroup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/accessLevelGroups"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiCreateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	groupName string
	body *OauthProfileAccessLevelGroupMsgVpnAccessLevelException
	opaquePassword *string
	select_ *[]string
}

// The Message VPN Access-Level Exception object&#39;s attributes.
func (r OauthProfileApiApiCreateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) Body(body OauthProfileAccessLevelGroupMsgVpnAccessLevelException) OauthProfileApiApiCreateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiCreateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) Select_(select_ []string) OauthProfileApiApiCreateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiCreateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) Execute() (*OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	return r.ApiService.CreateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionExecute(r)
}

/*
CreateOauthProfileAccessLevelGroupMsgVpnAccessLevelException Create a Message VPN Access-Level Exception object.

Create a Message VPN Access-Level Exception object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates via config-sync.

Message VPN access-level exceptions for members of this group.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
groupName|x||x|||
msgVpnName|x|x||||
oauthProfileName|x||x|||



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param groupName The name of the group.
 @return OauthProfileApiApiCreateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest
*/
func (a *OauthProfileApiService) CreateOauthProfileAccessLevelGroupMsgVpnAccessLevelException(ctx context.Context, oauthProfileName string, groupName string) OauthProfileApiApiCreateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	return OauthProfileApiApiCreateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		groupName: groupName,
	}
}

// Execute executes the request
//  @return OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse
func (a *OauthProfileApiService) CreateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionExecute(r OauthProfileApiApiCreateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) (*OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.CreateOauthProfileAccessLevelGroupMsgVpnAccessLevelException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"groupName"+"}", url.PathEscape(parameterToString(r.groupName, "")), -1)

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

type OauthProfileApiApiCreateOauthProfileClientAllowedHostRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	body *OauthProfileClientAllowedHost
	opaquePassword *string
	select_ *[]string
}

// The Allowed Host Value object&#39;s attributes.
func (r OauthProfileApiApiCreateOauthProfileClientAllowedHostRequest) Body(body OauthProfileClientAllowedHost) OauthProfileApiApiCreateOauthProfileClientAllowedHostRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileClientAllowedHostRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiCreateOauthProfileClientAllowedHostRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileClientAllowedHostRequest) Select_(select_ []string) OauthProfileApiApiCreateOauthProfileClientAllowedHostRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiCreateOauthProfileClientAllowedHostRequest) Execute() (*OauthProfileClientAllowedHostResponse, *http.Response, error) {
	return r.ApiService.CreateOauthProfileClientAllowedHostExecute(r)
}

/*
CreateOauthProfileClientAllowedHost Create an Allowed Host Value object.

Create an Allowed Host Value object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates via config-sync.

A valid hostname for this broker in OAuth redirects.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
allowedHost|x|x||||
oauthProfileName|x||x|||



A SEMP client authorized with a minimum access scope/level of "global/admin" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiCreateOauthProfileClientAllowedHostRequest
*/
func (a *OauthProfileApiService) CreateOauthProfileClientAllowedHost(ctx context.Context, oauthProfileName string) OauthProfileApiApiCreateOauthProfileClientAllowedHostRequest {
	return OauthProfileApiApiCreateOauthProfileClientAllowedHostRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileClientAllowedHostResponse
func (a *OauthProfileApiService) CreateOauthProfileClientAllowedHostExecute(r OauthProfileApiApiCreateOauthProfileClientAllowedHostRequest) (*OauthProfileClientAllowedHostResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileClientAllowedHostResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.CreateOauthProfileClientAllowedHost")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/clientAllowedHosts"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiCreateOauthProfileClientAuthorizationParameterRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	body *OauthProfileClientAuthorizationParameter
	opaquePassword *string
	select_ *[]string
}

// The Authorization Parameter object&#39;s attributes.
func (r OauthProfileApiApiCreateOauthProfileClientAuthorizationParameterRequest) Body(body OauthProfileClientAuthorizationParameter) OauthProfileApiApiCreateOauthProfileClientAuthorizationParameterRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileClientAuthorizationParameterRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiCreateOauthProfileClientAuthorizationParameterRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileClientAuthorizationParameterRequest) Select_(select_ []string) OauthProfileApiApiCreateOauthProfileClientAuthorizationParameterRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiCreateOauthProfileClientAuthorizationParameterRequest) Execute() (*OauthProfileClientAuthorizationParameterResponse, *http.Response, error) {
	return r.ApiService.CreateOauthProfileClientAuthorizationParameterExecute(r)
}

/*
CreateOauthProfileClientAuthorizationParameter Create an Authorization Parameter object.

Create an Authorization Parameter object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates via config-sync.

Additional parameters to be passed to the OAuth authorization endpoint.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
authorizationParameterName|x|x||||
oauthProfileName|x||x|||



A SEMP client authorized with a minimum access scope/level of "global/admin" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiCreateOauthProfileClientAuthorizationParameterRequest
*/
func (a *OauthProfileApiService) CreateOauthProfileClientAuthorizationParameter(ctx context.Context, oauthProfileName string) OauthProfileApiApiCreateOauthProfileClientAuthorizationParameterRequest {
	return OauthProfileApiApiCreateOauthProfileClientAuthorizationParameterRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileClientAuthorizationParameterResponse
func (a *OauthProfileApiService) CreateOauthProfileClientAuthorizationParameterExecute(r OauthProfileApiApiCreateOauthProfileClientAuthorizationParameterRequest) (*OauthProfileClientAuthorizationParameterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileClientAuthorizationParameterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.CreateOauthProfileClientAuthorizationParameter")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/clientAuthorizationParameters"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiCreateOauthProfileClientRequiredClaimRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	body *OauthProfileClientRequiredClaim
	opaquePassword *string
	select_ *[]string
}

// The Required Claim object&#39;s attributes.
func (r OauthProfileApiApiCreateOauthProfileClientRequiredClaimRequest) Body(body OauthProfileClientRequiredClaim) OauthProfileApiApiCreateOauthProfileClientRequiredClaimRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileClientRequiredClaimRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiCreateOauthProfileClientRequiredClaimRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileClientRequiredClaimRequest) Select_(select_ []string) OauthProfileApiApiCreateOauthProfileClientRequiredClaimRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiCreateOauthProfileClientRequiredClaimRequest) Execute() (*OauthProfileClientRequiredClaimResponse, *http.Response, error) {
	return r.ApiService.CreateOauthProfileClientRequiredClaimExecute(r)
}

/*
CreateOauthProfileClientRequiredClaim Create a Required Claim object.

Create a Required Claim object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates via config-sync.

Additional claims to be verified in the ID token.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
clientRequiredClaimName|x|x||||
clientRequiredClaimValue||x||||
oauthProfileName|x||x|||



A SEMP client authorized with a minimum access scope/level of "global/admin" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiCreateOauthProfileClientRequiredClaimRequest
*/
func (a *OauthProfileApiService) CreateOauthProfileClientRequiredClaim(ctx context.Context, oauthProfileName string) OauthProfileApiApiCreateOauthProfileClientRequiredClaimRequest {
	return OauthProfileApiApiCreateOauthProfileClientRequiredClaimRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileClientRequiredClaimResponse
func (a *OauthProfileApiService) CreateOauthProfileClientRequiredClaimExecute(r OauthProfileApiApiCreateOauthProfileClientRequiredClaimRequest) (*OauthProfileClientRequiredClaimResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileClientRequiredClaimResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.CreateOauthProfileClientRequiredClaim")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/clientRequiredClaims"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiCreateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	body *OauthProfileDefaultMsgVpnAccessLevelException
	opaquePassword *string
	select_ *[]string
}

// The Message VPN Access-Level Exception object&#39;s attributes.
func (r OauthProfileApiApiCreateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) Body(body OauthProfileDefaultMsgVpnAccessLevelException) OauthProfileApiApiCreateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiCreateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) Select_(select_ []string) OauthProfileApiApiCreateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiCreateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) Execute() (*OauthProfileDefaultMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	return r.ApiService.CreateOauthProfileDefaultMsgVpnAccessLevelExceptionExecute(r)
}

/*
CreateOauthProfileDefaultMsgVpnAccessLevelException Create a Message VPN Access-Level Exception object.

Create a Message VPN Access-Level Exception object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates via config-sync.

Default message VPN access-level exceptions.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
msgVpnName|x|x||||
oauthProfileName|x||x|||



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiCreateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest
*/
func (a *OauthProfileApiService) CreateOauthProfileDefaultMsgVpnAccessLevelException(ctx context.Context, oauthProfileName string) OauthProfileApiApiCreateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	return OauthProfileApiApiCreateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileDefaultMsgVpnAccessLevelExceptionResponse
func (a *OauthProfileApiService) CreateOauthProfileDefaultMsgVpnAccessLevelExceptionExecute(r OauthProfileApiApiCreateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) (*OauthProfileDefaultMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileDefaultMsgVpnAccessLevelExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.CreateOauthProfileDefaultMsgVpnAccessLevelException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiCreateOauthProfileResourceServerRequiredClaimRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	body *OauthProfileResourceServerRequiredClaim
	opaquePassword *string
	select_ *[]string
}

// The Required Claim object&#39;s attributes.
func (r OauthProfileApiApiCreateOauthProfileResourceServerRequiredClaimRequest) Body(body OauthProfileResourceServerRequiredClaim) OauthProfileApiApiCreateOauthProfileResourceServerRequiredClaimRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileResourceServerRequiredClaimRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiCreateOauthProfileResourceServerRequiredClaimRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiCreateOauthProfileResourceServerRequiredClaimRequest) Select_(select_ []string) OauthProfileApiApiCreateOauthProfileResourceServerRequiredClaimRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiCreateOauthProfileResourceServerRequiredClaimRequest) Execute() (*OauthProfileResourceServerRequiredClaimResponse, *http.Response, error) {
	return r.ApiService.CreateOauthProfileResourceServerRequiredClaimExecute(r)
}

/*
CreateOauthProfileResourceServerRequiredClaim Create a Required Claim object.

Create a Required Claim object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates via config-sync.

Additional claims to be verified in the access token.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
oauthProfileName|x||x|||
resourceServerRequiredClaimName|x|x||||
resourceServerRequiredClaimValue||x||||



A SEMP client authorized with a minimum access scope/level of "global/admin" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiCreateOauthProfileResourceServerRequiredClaimRequest
*/
func (a *OauthProfileApiService) CreateOauthProfileResourceServerRequiredClaim(ctx context.Context, oauthProfileName string) OauthProfileApiApiCreateOauthProfileResourceServerRequiredClaimRequest {
	return OauthProfileApiApiCreateOauthProfileResourceServerRequiredClaimRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileResourceServerRequiredClaimResponse
func (a *OauthProfileApiService) CreateOauthProfileResourceServerRequiredClaimExecute(r OauthProfileApiApiCreateOauthProfileResourceServerRequiredClaimRequest) (*OauthProfileResourceServerRequiredClaimResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileResourceServerRequiredClaimResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.CreateOauthProfileResourceServerRequiredClaim")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/resourceServerRequiredClaims"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiDeleteOauthProfileRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
}


func (r OauthProfileApiApiDeleteOauthProfileRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteOauthProfileExecute(r)
}

/*
DeleteOauthProfile Delete an OAuth Profile object.

Delete an OAuth Profile object. The deletion of instances of this object are synchronized to HA mates via config-sync.

OAuth profiles specify how to securely authenticate to an OAuth provider.

A SEMP client authorized with a minimum access scope/level of "global/admin" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiDeleteOauthProfileRequest
*/
func (a *OauthProfileApiService) DeleteOauthProfile(ctx context.Context, oauthProfileName string) OauthProfileApiApiDeleteOauthProfileRequest {
	return OauthProfileApiApiDeleteOauthProfileRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *OauthProfileApiService) DeleteOauthProfileExecute(r OauthProfileApiApiDeleteOauthProfileRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.DeleteOauthProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiDeleteOauthProfileAccessLevelGroupRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	groupName string
}


func (r OauthProfileApiApiDeleteOauthProfileAccessLevelGroupRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteOauthProfileAccessLevelGroupExecute(r)
}

/*
DeleteOauthProfileAccessLevelGroup Delete a Group Access Level object.

Delete a Group Access Level object. The deletion of instances of this object are synchronized to HA mates via config-sync.

The name of a group as it exists on the OAuth server being used to authenticate SEMP users.

A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param groupName The name of the group.
 @return OauthProfileApiApiDeleteOauthProfileAccessLevelGroupRequest
*/
func (a *OauthProfileApiService) DeleteOauthProfileAccessLevelGroup(ctx context.Context, oauthProfileName string, groupName string) OauthProfileApiApiDeleteOauthProfileAccessLevelGroupRequest {
	return OauthProfileApiApiDeleteOauthProfileAccessLevelGroupRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		groupName: groupName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *OauthProfileApiService) DeleteOauthProfileAccessLevelGroupExecute(r OauthProfileApiApiDeleteOauthProfileAccessLevelGroupRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.DeleteOauthProfileAccessLevelGroup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"groupName"+"}", url.PathEscape(parameterToString(r.groupName, "")), -1)

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

type OauthProfileApiApiDeleteOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	groupName string
	msgVpnName string
}


func (r OauthProfileApiApiDeleteOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionExecute(r)
}

/*
DeleteOauthProfileAccessLevelGroupMsgVpnAccessLevelException Delete a Message VPN Access-Level Exception object.

Delete a Message VPN Access-Level Exception object. The deletion of instances of this object are synchronized to HA mates via config-sync.

Message VPN access-level exceptions for members of this group.

A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param groupName The name of the group.
 @param msgVpnName The name of the message VPN.
 @return OauthProfileApiApiDeleteOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest
*/
func (a *OauthProfileApiService) DeleteOauthProfileAccessLevelGroupMsgVpnAccessLevelException(ctx context.Context, oauthProfileName string, groupName string, msgVpnName string) OauthProfileApiApiDeleteOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	return OauthProfileApiApiDeleteOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		groupName: groupName,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *OauthProfileApiService) DeleteOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionExecute(r OauthProfileApiApiDeleteOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.DeleteOauthProfileAccessLevelGroupMsgVpnAccessLevelException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions/{msgVpnName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"groupName"+"}", url.PathEscape(parameterToString(r.groupName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)

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

type OauthProfileApiApiDeleteOauthProfileClientAllowedHostRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	allowedHost string
}


func (r OauthProfileApiApiDeleteOauthProfileClientAllowedHostRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteOauthProfileClientAllowedHostExecute(r)
}

/*
DeleteOauthProfileClientAllowedHost Delete an Allowed Host Value object.

Delete an Allowed Host Value object. The deletion of instances of this object are synchronized to HA mates via config-sync.

A valid hostname for this broker in OAuth redirects.

A SEMP client authorized with a minimum access scope/level of "global/admin" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param allowedHost An allowed value for the Host header.
 @return OauthProfileApiApiDeleteOauthProfileClientAllowedHostRequest
*/
func (a *OauthProfileApiService) DeleteOauthProfileClientAllowedHost(ctx context.Context, oauthProfileName string, allowedHost string) OauthProfileApiApiDeleteOauthProfileClientAllowedHostRequest {
	return OauthProfileApiApiDeleteOauthProfileClientAllowedHostRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		allowedHost: allowedHost,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *OauthProfileApiService) DeleteOauthProfileClientAllowedHostExecute(r OauthProfileApiApiDeleteOauthProfileClientAllowedHostRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.DeleteOauthProfileClientAllowedHost")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/clientAllowedHosts/{allowedHost}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"allowedHost"+"}", url.PathEscape(parameterToString(r.allowedHost, "")), -1)

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

type OauthProfileApiApiDeleteOauthProfileClientAuthorizationParameterRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	authorizationParameterName string
}


func (r OauthProfileApiApiDeleteOauthProfileClientAuthorizationParameterRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteOauthProfileClientAuthorizationParameterExecute(r)
}

/*
DeleteOauthProfileClientAuthorizationParameter Delete an Authorization Parameter object.

Delete an Authorization Parameter object. The deletion of instances of this object are synchronized to HA mates via config-sync.

Additional parameters to be passed to the OAuth authorization endpoint.

A SEMP client authorized with a minimum access scope/level of "global/admin" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param authorizationParameterName The name of the authorization parameter.
 @return OauthProfileApiApiDeleteOauthProfileClientAuthorizationParameterRequest
*/
func (a *OauthProfileApiService) DeleteOauthProfileClientAuthorizationParameter(ctx context.Context, oauthProfileName string, authorizationParameterName string) OauthProfileApiApiDeleteOauthProfileClientAuthorizationParameterRequest {
	return OauthProfileApiApiDeleteOauthProfileClientAuthorizationParameterRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		authorizationParameterName: authorizationParameterName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *OauthProfileApiService) DeleteOauthProfileClientAuthorizationParameterExecute(r OauthProfileApiApiDeleteOauthProfileClientAuthorizationParameterRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.DeleteOauthProfileClientAuthorizationParameter")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/clientAuthorizationParameters/{authorizationParameterName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"authorizationParameterName"+"}", url.PathEscape(parameterToString(r.authorizationParameterName, "")), -1)

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

type OauthProfileApiApiDeleteOauthProfileClientRequiredClaimRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	clientRequiredClaimName string
}


func (r OauthProfileApiApiDeleteOauthProfileClientRequiredClaimRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteOauthProfileClientRequiredClaimExecute(r)
}

/*
DeleteOauthProfileClientRequiredClaim Delete a Required Claim object.

Delete a Required Claim object. The deletion of instances of this object are synchronized to HA mates via config-sync.

Additional claims to be verified in the ID token.

A SEMP client authorized with a minimum access scope/level of "global/admin" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param clientRequiredClaimName The name of the ID token claim to verify.
 @return OauthProfileApiApiDeleteOauthProfileClientRequiredClaimRequest
*/
func (a *OauthProfileApiService) DeleteOauthProfileClientRequiredClaim(ctx context.Context, oauthProfileName string, clientRequiredClaimName string) OauthProfileApiApiDeleteOauthProfileClientRequiredClaimRequest {
	return OauthProfileApiApiDeleteOauthProfileClientRequiredClaimRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		clientRequiredClaimName: clientRequiredClaimName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *OauthProfileApiService) DeleteOauthProfileClientRequiredClaimExecute(r OauthProfileApiApiDeleteOauthProfileClientRequiredClaimRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.DeleteOauthProfileClientRequiredClaim")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/clientRequiredClaims/{clientRequiredClaimName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clientRequiredClaimName"+"}", url.PathEscape(parameterToString(r.clientRequiredClaimName, "")), -1)

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

type OauthProfileApiApiDeleteOauthProfileDefaultMsgVpnAccessLevelExceptionRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	msgVpnName string
}


func (r OauthProfileApiApiDeleteOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteOauthProfileDefaultMsgVpnAccessLevelExceptionExecute(r)
}

/*
DeleteOauthProfileDefaultMsgVpnAccessLevelException Delete a Message VPN Access-Level Exception object.

Delete a Message VPN Access-Level Exception object. The deletion of instances of this object are synchronized to HA mates via config-sync.

Default message VPN access-level exceptions.

A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param msgVpnName The name of the message VPN.
 @return OauthProfileApiApiDeleteOauthProfileDefaultMsgVpnAccessLevelExceptionRequest
*/
func (a *OauthProfileApiService) DeleteOauthProfileDefaultMsgVpnAccessLevelException(ctx context.Context, oauthProfileName string, msgVpnName string) OauthProfileApiApiDeleteOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	return OauthProfileApiApiDeleteOauthProfileDefaultMsgVpnAccessLevelExceptionRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *OauthProfileApiService) DeleteOauthProfileDefaultMsgVpnAccessLevelExceptionExecute(r OauthProfileApiApiDeleteOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.DeleteOauthProfileDefaultMsgVpnAccessLevelException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions/{msgVpnName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)

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

type OauthProfileApiApiDeleteOauthProfileResourceServerRequiredClaimRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	resourceServerRequiredClaimName string
}


func (r OauthProfileApiApiDeleteOauthProfileResourceServerRequiredClaimRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteOauthProfileResourceServerRequiredClaimExecute(r)
}

/*
DeleteOauthProfileResourceServerRequiredClaim Delete a Required Claim object.

Delete a Required Claim object. The deletion of instances of this object are synchronized to HA mates via config-sync.

Additional claims to be verified in the access token.

A SEMP client authorized with a minimum access scope/level of "global/admin" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param resourceServerRequiredClaimName The name of the access token claim to verify.
 @return OauthProfileApiApiDeleteOauthProfileResourceServerRequiredClaimRequest
*/
func (a *OauthProfileApiService) DeleteOauthProfileResourceServerRequiredClaim(ctx context.Context, oauthProfileName string, resourceServerRequiredClaimName string) OauthProfileApiApiDeleteOauthProfileResourceServerRequiredClaimRequest {
	return OauthProfileApiApiDeleteOauthProfileResourceServerRequiredClaimRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		resourceServerRequiredClaimName: resourceServerRequiredClaimName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *OauthProfileApiService) DeleteOauthProfileResourceServerRequiredClaimExecute(r OauthProfileApiApiDeleteOauthProfileResourceServerRequiredClaimRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.DeleteOauthProfileResourceServerRequiredClaim")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/resourceServerRequiredClaims/{resourceServerRequiredClaimName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"resourceServerRequiredClaimName"+"}", url.PathEscape(parameterToString(r.resourceServerRequiredClaimName, "")), -1)

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

type OauthProfileApiApiGetOauthProfileRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileRequest) Execute() (*OauthProfileResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileExecute(r)
}

/*
GetOauthProfile Get an OAuth Profile object.

Get an OAuth Profile object.

OAuth profiles specify how to securely authenticate to an OAuth provider.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
clientSecret||x||x
oauthProfileName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiGetOauthProfileRequest
*/
func (a *OauthProfileApiService) GetOauthProfile(ctx context.Context, oauthProfileName string) OauthProfileApiApiGetOauthProfileRequest {
	return OauthProfileApiApiGetOauthProfileRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileResponse
func (a *OauthProfileApiService) GetOauthProfileExecute(r OauthProfileApiApiGetOauthProfileRequest) (*OauthProfileResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiGetOauthProfileAccessLevelGroupRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	groupName string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileAccessLevelGroupRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileAccessLevelGroupRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupRequest) Execute() (*OauthProfileAccessLevelGroupResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileAccessLevelGroupExecute(r)
}

/*
GetOauthProfileAccessLevelGroup Get a Group Access Level object.

Get a Group Access Level object.

The name of a group as it exists on the OAuth server being used to authenticate SEMP users.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
groupName|x|||
oauthProfileName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param groupName The name of the group.
 @return OauthProfileApiApiGetOauthProfileAccessLevelGroupRequest
*/
func (a *OauthProfileApiService) GetOauthProfileAccessLevelGroup(ctx context.Context, oauthProfileName string, groupName string) OauthProfileApiApiGetOauthProfileAccessLevelGroupRequest {
	return OauthProfileApiApiGetOauthProfileAccessLevelGroupRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		groupName: groupName,
	}
}

// Execute executes the request
//  @return OauthProfileAccessLevelGroupResponse
func (a *OauthProfileApiService) GetOauthProfileAccessLevelGroupExecute(r OauthProfileApiApiGetOauthProfileAccessLevelGroupRequest) (*OauthProfileAccessLevelGroupResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileAccessLevelGroupResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfileAccessLevelGroup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"groupName"+"}", url.PathEscape(parameterToString(r.groupName, "")), -1)

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

type OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	groupName string
	msgVpnName string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) Execute() (*OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionExecute(r)
}

/*
GetOauthProfileAccessLevelGroupMsgVpnAccessLevelException Get a Message VPN Access-Level Exception object.

Get a Message VPN Access-Level Exception object.

Message VPN access-level exceptions for members of this group.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
groupName|x|||
msgVpnName|x|||
oauthProfileName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param groupName The name of the group.
 @param msgVpnName The name of the message VPN.
 @return OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest
*/
func (a *OauthProfileApiService) GetOauthProfileAccessLevelGroupMsgVpnAccessLevelException(ctx context.Context, oauthProfileName string, groupName string, msgVpnName string) OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	return OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		groupName: groupName,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse
func (a *OauthProfileApiService) GetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionExecute(r OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) (*OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfileAccessLevelGroupMsgVpnAccessLevelException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions/{msgVpnName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"groupName"+"}", url.PathEscape(parameterToString(r.groupName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)

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

type OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	groupName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest) Count(count int32) OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest) Cursor(cursor string) OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest) Where(where []string) OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest) Execute() (*OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsExecute(r)
}

/*
GetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptions Get a list of Message VPN Access-Level Exception objects.

Get a list of Message VPN Access-Level Exception objects.

Message VPN access-level exceptions for members of this group.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
groupName|x|||
msgVpnName|x|||
oauthProfileName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param groupName The name of the group.
 @return OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest
*/
func (a *OauthProfileApiService) GetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptions(ctx context.Context, oauthProfileName string, groupName string) OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest {
	return OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		groupName: groupName,
	}
}

// Execute executes the request
//  @return OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsResponse
func (a *OauthProfileApiService) GetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsExecute(r OauthProfileApiApiGetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsRequest) (*OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"groupName"+"}", url.PathEscape(parameterToString(r.groupName, "")), -1)

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

type OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest) Count(count int32) OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest) Cursor(cursor string) OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest) Where(where []string) OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest) Execute() (*OauthProfileAccessLevelGroupsResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileAccessLevelGroupsExecute(r)
}

/*
GetOauthProfileAccessLevelGroups Get a list of Group Access Level objects.

Get a list of Group Access Level objects.

The name of a group as it exists on the OAuth server being used to authenticate SEMP users.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
groupName|x|||
oauthProfileName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest
*/
func (a *OauthProfileApiService) GetOauthProfileAccessLevelGroups(ctx context.Context, oauthProfileName string) OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest {
	return OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileAccessLevelGroupsResponse
func (a *OauthProfileApiService) GetOauthProfileAccessLevelGroupsExecute(r OauthProfileApiApiGetOauthProfileAccessLevelGroupsRequest) (*OauthProfileAccessLevelGroupsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileAccessLevelGroupsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfileAccessLevelGroups")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/accessLevelGroups"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiGetOauthProfileClientAllowedHostRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	allowedHost string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientAllowedHostRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileClientAllowedHostRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientAllowedHostRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileClientAllowedHostRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileClientAllowedHostRequest) Execute() (*OauthProfileClientAllowedHostResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileClientAllowedHostExecute(r)
}

/*
GetOauthProfileClientAllowedHost Get an Allowed Host Value object.

Get an Allowed Host Value object.

A valid hostname for this broker in OAuth redirects.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
allowedHost|x|||
oauthProfileName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param allowedHost An allowed value for the Host header.
 @return OauthProfileApiApiGetOauthProfileClientAllowedHostRequest
*/
func (a *OauthProfileApiService) GetOauthProfileClientAllowedHost(ctx context.Context, oauthProfileName string, allowedHost string) OauthProfileApiApiGetOauthProfileClientAllowedHostRequest {
	return OauthProfileApiApiGetOauthProfileClientAllowedHostRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		allowedHost: allowedHost,
	}
}

// Execute executes the request
//  @return OauthProfileClientAllowedHostResponse
func (a *OauthProfileApiService) GetOauthProfileClientAllowedHostExecute(r OauthProfileApiApiGetOauthProfileClientAllowedHostRequest) (*OauthProfileClientAllowedHostResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileClientAllowedHostResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfileClientAllowedHost")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/clientAllowedHosts/{allowedHost}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"allowedHost"+"}", url.PathEscape(parameterToString(r.allowedHost, "")), -1)

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

type OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest) Count(count int32) OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest) Cursor(cursor string) OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest) Where(where []string) OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest) Execute() (*OauthProfileClientAllowedHostsResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileClientAllowedHostsExecute(r)
}

/*
GetOauthProfileClientAllowedHosts Get a list of Allowed Host Value objects.

Get a list of Allowed Host Value objects.

A valid hostname for this broker in OAuth redirects.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
allowedHost|x|||
oauthProfileName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest
*/
func (a *OauthProfileApiService) GetOauthProfileClientAllowedHosts(ctx context.Context, oauthProfileName string) OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest {
	return OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileClientAllowedHostsResponse
func (a *OauthProfileApiService) GetOauthProfileClientAllowedHostsExecute(r OauthProfileApiApiGetOauthProfileClientAllowedHostsRequest) (*OauthProfileClientAllowedHostsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileClientAllowedHostsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfileClientAllowedHosts")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/clientAllowedHosts"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiGetOauthProfileClientAuthorizationParameterRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	authorizationParameterName string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientAuthorizationParameterRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileClientAuthorizationParameterRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientAuthorizationParameterRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileClientAuthorizationParameterRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileClientAuthorizationParameterRequest) Execute() (*OauthProfileClientAuthorizationParameterResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileClientAuthorizationParameterExecute(r)
}

/*
GetOauthProfileClientAuthorizationParameter Get an Authorization Parameter object.

Get an Authorization Parameter object.

Additional parameters to be passed to the OAuth authorization endpoint.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
authorizationParameterName|x|||
oauthProfileName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param authorizationParameterName The name of the authorization parameter.
 @return OauthProfileApiApiGetOauthProfileClientAuthorizationParameterRequest
*/
func (a *OauthProfileApiService) GetOauthProfileClientAuthorizationParameter(ctx context.Context, oauthProfileName string, authorizationParameterName string) OauthProfileApiApiGetOauthProfileClientAuthorizationParameterRequest {
	return OauthProfileApiApiGetOauthProfileClientAuthorizationParameterRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		authorizationParameterName: authorizationParameterName,
	}
}

// Execute executes the request
//  @return OauthProfileClientAuthorizationParameterResponse
func (a *OauthProfileApiService) GetOauthProfileClientAuthorizationParameterExecute(r OauthProfileApiApiGetOauthProfileClientAuthorizationParameterRequest) (*OauthProfileClientAuthorizationParameterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileClientAuthorizationParameterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfileClientAuthorizationParameter")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/clientAuthorizationParameters/{authorizationParameterName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"authorizationParameterName"+"}", url.PathEscape(parameterToString(r.authorizationParameterName, "")), -1)

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

type OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest) Count(count int32) OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest) Cursor(cursor string) OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest) Where(where []string) OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest) Execute() (*OauthProfileClientAuthorizationParametersResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileClientAuthorizationParametersExecute(r)
}

/*
GetOauthProfileClientAuthorizationParameters Get a list of Authorization Parameter objects.

Get a list of Authorization Parameter objects.

Additional parameters to be passed to the OAuth authorization endpoint.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
authorizationParameterName|x|||
oauthProfileName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest
*/
func (a *OauthProfileApiService) GetOauthProfileClientAuthorizationParameters(ctx context.Context, oauthProfileName string) OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest {
	return OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileClientAuthorizationParametersResponse
func (a *OauthProfileApiService) GetOauthProfileClientAuthorizationParametersExecute(r OauthProfileApiApiGetOauthProfileClientAuthorizationParametersRequest) (*OauthProfileClientAuthorizationParametersResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileClientAuthorizationParametersResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfileClientAuthorizationParameters")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/clientAuthorizationParameters"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiGetOauthProfileClientRequiredClaimRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	clientRequiredClaimName string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientRequiredClaimRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileClientRequiredClaimRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientRequiredClaimRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileClientRequiredClaimRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileClientRequiredClaimRequest) Execute() (*OauthProfileClientRequiredClaimResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileClientRequiredClaimExecute(r)
}

/*
GetOauthProfileClientRequiredClaim Get a Required Claim object.

Get a Required Claim object.

Additional claims to be verified in the ID token.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
clientRequiredClaimName|x|||
oauthProfileName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param clientRequiredClaimName The name of the ID token claim to verify.
 @return OauthProfileApiApiGetOauthProfileClientRequiredClaimRequest
*/
func (a *OauthProfileApiService) GetOauthProfileClientRequiredClaim(ctx context.Context, oauthProfileName string, clientRequiredClaimName string) OauthProfileApiApiGetOauthProfileClientRequiredClaimRequest {
	return OauthProfileApiApiGetOauthProfileClientRequiredClaimRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		clientRequiredClaimName: clientRequiredClaimName,
	}
}

// Execute executes the request
//  @return OauthProfileClientRequiredClaimResponse
func (a *OauthProfileApiService) GetOauthProfileClientRequiredClaimExecute(r OauthProfileApiApiGetOauthProfileClientRequiredClaimRequest) (*OauthProfileClientRequiredClaimResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileClientRequiredClaimResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfileClientRequiredClaim")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/clientRequiredClaims/{clientRequiredClaimName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clientRequiredClaimName"+"}", url.PathEscape(parameterToString(r.clientRequiredClaimName, "")), -1)

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

type OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest) Count(count int32) OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest) Cursor(cursor string) OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest) Where(where []string) OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest) Execute() (*OauthProfileClientRequiredClaimsResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileClientRequiredClaimsExecute(r)
}

/*
GetOauthProfileClientRequiredClaims Get a list of Required Claim objects.

Get a list of Required Claim objects.

Additional claims to be verified in the ID token.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
clientRequiredClaimName|x|||
oauthProfileName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest
*/
func (a *OauthProfileApiService) GetOauthProfileClientRequiredClaims(ctx context.Context, oauthProfileName string) OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest {
	return OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileClientRequiredClaimsResponse
func (a *OauthProfileApiService) GetOauthProfileClientRequiredClaimsExecute(r OauthProfileApiApiGetOauthProfileClientRequiredClaimsRequest) (*OauthProfileClientRequiredClaimsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileClientRequiredClaimsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfileClientRequiredClaims")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/clientRequiredClaims"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	msgVpnName string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) Execute() (*OauthProfileDefaultMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileDefaultMsgVpnAccessLevelExceptionExecute(r)
}

/*
GetOauthProfileDefaultMsgVpnAccessLevelException Get a Message VPN Access-Level Exception object.

Get a Message VPN Access-Level Exception object.

Default message VPN access-level exceptions.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
msgVpnName|x|||
oauthProfileName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param msgVpnName The name of the message VPN.
 @return OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionRequest
*/
func (a *OauthProfileApiService) GetOauthProfileDefaultMsgVpnAccessLevelException(ctx context.Context, oauthProfileName string, msgVpnName string) OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	return OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return OauthProfileDefaultMsgVpnAccessLevelExceptionResponse
func (a *OauthProfileApiService) GetOauthProfileDefaultMsgVpnAccessLevelExceptionExecute(r OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) (*OauthProfileDefaultMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileDefaultMsgVpnAccessLevelExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfileDefaultMsgVpnAccessLevelException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions/{msgVpnName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)

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

type OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest) Count(count int32) OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest) Cursor(cursor string) OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest) Where(where []string) OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest) Execute() (*OauthProfileDefaultMsgVpnAccessLevelExceptionsResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileDefaultMsgVpnAccessLevelExceptionsExecute(r)
}

/*
GetOauthProfileDefaultMsgVpnAccessLevelExceptions Get a list of Message VPN Access-Level Exception objects.

Get a list of Message VPN Access-Level Exception objects.

Default message VPN access-level exceptions.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
msgVpnName|x|||
oauthProfileName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest
*/
func (a *OauthProfileApiService) GetOauthProfileDefaultMsgVpnAccessLevelExceptions(ctx context.Context, oauthProfileName string) OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest {
	return OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileDefaultMsgVpnAccessLevelExceptionsResponse
func (a *OauthProfileApiService) GetOauthProfileDefaultMsgVpnAccessLevelExceptionsExecute(r OauthProfileApiApiGetOauthProfileDefaultMsgVpnAccessLevelExceptionsRequest) (*OauthProfileDefaultMsgVpnAccessLevelExceptionsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileDefaultMsgVpnAccessLevelExceptionsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfileDefaultMsgVpnAccessLevelExceptions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	resourceServerRequiredClaimName string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimRequest) Execute() (*OauthProfileResourceServerRequiredClaimResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileResourceServerRequiredClaimExecute(r)
}

/*
GetOauthProfileResourceServerRequiredClaim Get a Required Claim object.

Get a Required Claim object.

Additional claims to be verified in the access token.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
oauthProfileName|x|||
resourceServerRequiredClaimName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param resourceServerRequiredClaimName The name of the access token claim to verify.
 @return OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimRequest
*/
func (a *OauthProfileApiService) GetOauthProfileResourceServerRequiredClaim(ctx context.Context, oauthProfileName string, resourceServerRequiredClaimName string) OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimRequest {
	return OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		resourceServerRequiredClaimName: resourceServerRequiredClaimName,
	}
}

// Execute executes the request
//  @return OauthProfileResourceServerRequiredClaimResponse
func (a *OauthProfileApiService) GetOauthProfileResourceServerRequiredClaimExecute(r OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimRequest) (*OauthProfileResourceServerRequiredClaimResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileResourceServerRequiredClaimResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfileResourceServerRequiredClaim")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/resourceServerRequiredClaims/{resourceServerRequiredClaimName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"resourceServerRequiredClaimName"+"}", url.PathEscape(parameterToString(r.resourceServerRequiredClaimName, "")), -1)

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

type OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest) Count(count int32) OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest) Cursor(cursor string) OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest) Where(where []string) OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest) Execute() (*OauthProfileResourceServerRequiredClaimsResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfileResourceServerRequiredClaimsExecute(r)
}

/*
GetOauthProfileResourceServerRequiredClaims Get a list of Required Claim objects.

Get a list of Required Claim objects.

Additional claims to be verified in the access token.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
oauthProfileName|x|||
resourceServerRequiredClaimName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest
*/
func (a *OauthProfileApiService) GetOauthProfileResourceServerRequiredClaims(ctx context.Context, oauthProfileName string) OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest {
	return OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileResourceServerRequiredClaimsResponse
func (a *OauthProfileApiService) GetOauthProfileResourceServerRequiredClaimsExecute(r OauthProfileApiApiGetOauthProfileResourceServerRequiredClaimsRequest) (*OauthProfileResourceServerRequiredClaimsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileResourceServerRequiredClaimsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfileResourceServerRequiredClaims")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/resourceServerRequiredClaims"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiGetOauthProfilesRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfilesRequest) Count(count int32) OauthProfileApiApiGetOauthProfilesRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfilesRequest) Cursor(cursor string) OauthProfileApiApiGetOauthProfilesRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfilesRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiGetOauthProfilesRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfilesRequest) Where(where []string) OauthProfileApiApiGetOauthProfilesRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiGetOauthProfilesRequest) Select_(select_ []string) OauthProfileApiApiGetOauthProfilesRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiGetOauthProfilesRequest) Execute() (*OauthProfilesResponse, *http.Response, error) {
	return r.ApiService.GetOauthProfilesExecute(r)
}

/*
GetOauthProfiles Get a list of OAuth Profile objects.

Get a list of OAuth Profile objects.

OAuth profiles specify how to securely authenticate to an OAuth provider.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
clientSecret||x||x
oauthProfileName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return OauthProfileApiApiGetOauthProfilesRequest
*/
func (a *OauthProfileApiService) GetOauthProfiles(ctx context.Context) OauthProfileApiApiGetOauthProfilesRequest {
	return OauthProfileApiApiGetOauthProfilesRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return OauthProfilesResponse
func (a *OauthProfileApiService) GetOauthProfilesExecute(r OauthProfileApiApiGetOauthProfilesRequest) (*OauthProfilesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfilesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.GetOauthProfiles")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles"

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

type OauthProfileApiApiReplaceOauthProfileRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	body *OauthProfile
	opaquePassword *string
	select_ *[]string
}

// The OAuth Profile object&#39;s attributes.
func (r OauthProfileApiApiReplaceOauthProfileRequest) Body(body OauthProfile) OauthProfileApiApiReplaceOauthProfileRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiReplaceOauthProfileRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiReplaceOauthProfileRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiReplaceOauthProfileRequest) Select_(select_ []string) OauthProfileApiApiReplaceOauthProfileRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiReplaceOauthProfileRequest) Execute() (*OauthProfileResponse, *http.Response, error) {
	return r.ApiService.ReplaceOauthProfileExecute(r)
}

/*
ReplaceOauthProfile Replace an OAuth Profile object.

Replace an OAuth Profile object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

OAuth profiles specify how to securely authenticate to an OAuth provider.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
clientSecret|||x|||x
oauthProfileName|x|x||||



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation. Requests which include the following attributes require greater access scope/level:


Attribute|Access Scope/Level
:---|:---:
accessLevelGroupsClaimName|global/admin
clientId|global/admin
clientRedirectUri|global/admin
clientRequiredType|global/admin
clientScope|global/admin
clientSecret|global/admin
clientValidateTypeEnabled|global/admin
defaultGlobalAccessLevel|global/admin
displayName|global/admin
enabled|global/admin
endpointAuthorization|global/admin
endpointDiscovery|global/admin
endpointDiscoveryRefreshInterval|global/admin
endpointIntrospection|global/admin
endpointIntrospectionTimeout|global/admin
endpointJwks|global/admin
endpointJwksRefreshInterval|global/admin
endpointToken|global/admin
endpointTokenTimeout|global/admin
endpointUserinfo|global/admin
endpointUserinfoTimeout|global/admin
interactiveEnabled|global/admin
interactivePromptForExpiredSession|global/admin
interactivePromptForNewSession|global/admin
issuer|global/admin
oauthRole|global/admin
resourceServerParseAccessTokenEnabled|global/admin
resourceServerRequiredAudience|global/admin
resourceServerRequiredIssuer|global/admin
resourceServerRequiredScope|global/admin
resourceServerRequiredType|global/admin
resourceServerValidateAudienceEnabled|global/admin
resourceServerValidateIssuerEnabled|global/admin
resourceServerValidateScopeEnabled|global/admin
resourceServerValidateTypeEnabled|global/admin
sempEnabled|global/admin
usernameClaimName|global/admin



This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiReplaceOauthProfileRequest
*/
func (a *OauthProfileApiService) ReplaceOauthProfile(ctx context.Context, oauthProfileName string) OauthProfileApiApiReplaceOauthProfileRequest {
	return OauthProfileApiApiReplaceOauthProfileRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileResponse
func (a *OauthProfileApiService) ReplaceOauthProfileExecute(r OauthProfileApiApiReplaceOauthProfileRequest) (*OauthProfileResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.ReplaceOauthProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiReplaceOauthProfileAccessLevelGroupRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	groupName string
	body *OauthProfileAccessLevelGroup
	opaquePassword *string
	select_ *[]string
}

// The Group Access Level object&#39;s attributes.
func (r OauthProfileApiApiReplaceOauthProfileAccessLevelGroupRequest) Body(body OauthProfileAccessLevelGroup) OauthProfileApiApiReplaceOauthProfileAccessLevelGroupRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiReplaceOauthProfileAccessLevelGroupRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiReplaceOauthProfileAccessLevelGroupRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiReplaceOauthProfileAccessLevelGroupRequest) Select_(select_ []string) OauthProfileApiApiReplaceOauthProfileAccessLevelGroupRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiReplaceOauthProfileAccessLevelGroupRequest) Execute() (*OauthProfileAccessLevelGroupResponse, *http.Response, error) {
	return r.ApiService.ReplaceOauthProfileAccessLevelGroupExecute(r)
}

/*
ReplaceOauthProfileAccessLevelGroup Replace a Group Access Level object.

Replace a Group Access Level object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

The name of a group as it exists on the OAuth server being used to authenticate SEMP users.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
groupName|x|x||||
oauthProfileName|x|x||||



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation. Requests which include the following attributes require greater access scope/level:


Attribute|Access Scope/Level
:---|:---:
globalAccessLevel|global/admin



This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param groupName The name of the group.
 @return OauthProfileApiApiReplaceOauthProfileAccessLevelGroupRequest
*/
func (a *OauthProfileApiService) ReplaceOauthProfileAccessLevelGroup(ctx context.Context, oauthProfileName string, groupName string) OauthProfileApiApiReplaceOauthProfileAccessLevelGroupRequest {
	return OauthProfileApiApiReplaceOauthProfileAccessLevelGroupRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		groupName: groupName,
	}
}

// Execute executes the request
//  @return OauthProfileAccessLevelGroupResponse
func (a *OauthProfileApiService) ReplaceOauthProfileAccessLevelGroupExecute(r OauthProfileApiApiReplaceOauthProfileAccessLevelGroupRequest) (*OauthProfileAccessLevelGroupResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileAccessLevelGroupResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.ReplaceOauthProfileAccessLevelGroup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"groupName"+"}", url.PathEscape(parameterToString(r.groupName, "")), -1)

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

type OauthProfileApiApiReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	groupName string
	msgVpnName string
	body *OauthProfileAccessLevelGroupMsgVpnAccessLevelException
	opaquePassword *string
	select_ *[]string
}

// The Message VPN Access-Level Exception object&#39;s attributes.
func (r OauthProfileApiApiReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) Body(body OauthProfileAccessLevelGroupMsgVpnAccessLevelException) OauthProfileApiApiReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) Select_(select_ []string) OauthProfileApiApiReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) Execute() (*OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	return r.ApiService.ReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionExecute(r)
}

/*
ReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelException Replace a Message VPN Access-Level Exception object.

Replace a Message VPN Access-Level Exception object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

Message VPN access-level exceptions for members of this group.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
groupName|x|x||||
msgVpnName|x|x||||
oauthProfileName|x|x||||



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param groupName The name of the group.
 @param msgVpnName The name of the message VPN.
 @return OauthProfileApiApiReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest
*/
func (a *OauthProfileApiService) ReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelException(ctx context.Context, oauthProfileName string, groupName string, msgVpnName string) OauthProfileApiApiReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	return OauthProfileApiApiReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		groupName: groupName,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse
func (a *OauthProfileApiService) ReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionExecute(r OauthProfileApiApiReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) (*OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.ReplaceOauthProfileAccessLevelGroupMsgVpnAccessLevelException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions/{msgVpnName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"groupName"+"}", url.PathEscape(parameterToString(r.groupName, "")), -1)
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

type OauthProfileApiApiReplaceOauthProfileClientAuthorizationParameterRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	authorizationParameterName string
	body *OauthProfileClientAuthorizationParameter
	opaquePassword *string
	select_ *[]string
}

// The Authorization Parameter object&#39;s attributes.
func (r OauthProfileApiApiReplaceOauthProfileClientAuthorizationParameterRequest) Body(body OauthProfileClientAuthorizationParameter) OauthProfileApiApiReplaceOauthProfileClientAuthorizationParameterRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiReplaceOauthProfileClientAuthorizationParameterRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiReplaceOauthProfileClientAuthorizationParameterRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiReplaceOauthProfileClientAuthorizationParameterRequest) Select_(select_ []string) OauthProfileApiApiReplaceOauthProfileClientAuthorizationParameterRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiReplaceOauthProfileClientAuthorizationParameterRequest) Execute() (*OauthProfileClientAuthorizationParameterResponse, *http.Response, error) {
	return r.ApiService.ReplaceOauthProfileClientAuthorizationParameterExecute(r)
}

/*
ReplaceOauthProfileClientAuthorizationParameter Replace an Authorization Parameter object.

Replace an Authorization Parameter object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

Additional parameters to be passed to the OAuth authorization endpoint.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
authorizationParameterName|x|x||||
oauthProfileName|x|x||||



A SEMP client authorized with a minimum access scope/level of "global/admin" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param authorizationParameterName The name of the authorization parameter.
 @return OauthProfileApiApiReplaceOauthProfileClientAuthorizationParameterRequest
*/
func (a *OauthProfileApiService) ReplaceOauthProfileClientAuthorizationParameter(ctx context.Context, oauthProfileName string, authorizationParameterName string) OauthProfileApiApiReplaceOauthProfileClientAuthorizationParameterRequest {
	return OauthProfileApiApiReplaceOauthProfileClientAuthorizationParameterRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		authorizationParameterName: authorizationParameterName,
	}
}

// Execute executes the request
//  @return OauthProfileClientAuthorizationParameterResponse
func (a *OauthProfileApiService) ReplaceOauthProfileClientAuthorizationParameterExecute(r OauthProfileApiApiReplaceOauthProfileClientAuthorizationParameterRequest) (*OauthProfileClientAuthorizationParameterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileClientAuthorizationParameterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.ReplaceOauthProfileClientAuthorizationParameter")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/clientAuthorizationParameters/{authorizationParameterName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"authorizationParameterName"+"}", url.PathEscape(parameterToString(r.authorizationParameterName, "")), -1)

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

type OauthProfileApiApiReplaceOauthProfileDefaultMsgVpnAccessLevelExceptionRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	msgVpnName string
	body *OauthProfileDefaultMsgVpnAccessLevelException
	opaquePassword *string
	select_ *[]string
}

// The Message VPN Access-Level Exception object&#39;s attributes.
func (r OauthProfileApiApiReplaceOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) Body(body OauthProfileDefaultMsgVpnAccessLevelException) OauthProfileApiApiReplaceOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiReplaceOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiReplaceOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiReplaceOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) Select_(select_ []string) OauthProfileApiApiReplaceOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiReplaceOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) Execute() (*OauthProfileDefaultMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	return r.ApiService.ReplaceOauthProfileDefaultMsgVpnAccessLevelExceptionExecute(r)
}

/*
ReplaceOauthProfileDefaultMsgVpnAccessLevelException Replace a Message VPN Access-Level Exception object.

Replace a Message VPN Access-Level Exception object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

Default message VPN access-level exceptions.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
msgVpnName|x|x||||
oauthProfileName|x|x||||



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param msgVpnName The name of the message VPN.
 @return OauthProfileApiApiReplaceOauthProfileDefaultMsgVpnAccessLevelExceptionRequest
*/
func (a *OauthProfileApiService) ReplaceOauthProfileDefaultMsgVpnAccessLevelException(ctx context.Context, oauthProfileName string, msgVpnName string) OauthProfileApiApiReplaceOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	return OauthProfileApiApiReplaceOauthProfileDefaultMsgVpnAccessLevelExceptionRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return OauthProfileDefaultMsgVpnAccessLevelExceptionResponse
func (a *OauthProfileApiService) ReplaceOauthProfileDefaultMsgVpnAccessLevelExceptionExecute(r OauthProfileApiApiReplaceOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) (*OauthProfileDefaultMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileDefaultMsgVpnAccessLevelExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.ReplaceOauthProfileDefaultMsgVpnAccessLevelException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions/{msgVpnName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
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

type OauthProfileApiApiUpdateOauthProfileRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	body *OauthProfile
	opaquePassword *string
	select_ *[]string
}

// The OAuth Profile object&#39;s attributes.
func (r OauthProfileApiApiUpdateOauthProfileRequest) Body(body OauthProfile) OauthProfileApiApiUpdateOauthProfileRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiUpdateOauthProfileRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiUpdateOauthProfileRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiUpdateOauthProfileRequest) Select_(select_ []string) OauthProfileApiApiUpdateOauthProfileRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiUpdateOauthProfileRequest) Execute() (*OauthProfileResponse, *http.Response, error) {
	return r.ApiService.UpdateOauthProfileExecute(r)
}

/*
UpdateOauthProfile Update an OAuth Profile object.

Update an OAuth Profile object. Any attribute missing from the request will be left unchanged.

OAuth profiles specify how to securely authenticate to an OAuth provider.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
clientSecret|||x|||x
oauthProfileName|x|x||||



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation. Requests which include the following attributes require greater access scope/level:


Attribute|Access Scope/Level
:---|:---:
accessLevelGroupsClaimName|global/admin
clientId|global/admin
clientRedirectUri|global/admin
clientRequiredType|global/admin
clientScope|global/admin
clientSecret|global/admin
clientValidateTypeEnabled|global/admin
defaultGlobalAccessLevel|global/admin
displayName|global/admin
enabled|global/admin
endpointAuthorization|global/admin
endpointDiscovery|global/admin
endpointDiscoveryRefreshInterval|global/admin
endpointIntrospection|global/admin
endpointIntrospectionTimeout|global/admin
endpointJwks|global/admin
endpointJwksRefreshInterval|global/admin
endpointToken|global/admin
endpointTokenTimeout|global/admin
endpointUserinfo|global/admin
endpointUserinfoTimeout|global/admin
interactiveEnabled|global/admin
interactivePromptForExpiredSession|global/admin
interactivePromptForNewSession|global/admin
issuer|global/admin
oauthRole|global/admin
resourceServerParseAccessTokenEnabled|global/admin
resourceServerRequiredAudience|global/admin
resourceServerRequiredIssuer|global/admin
resourceServerRequiredScope|global/admin
resourceServerRequiredType|global/admin
resourceServerValidateAudienceEnabled|global/admin
resourceServerValidateIssuerEnabled|global/admin
resourceServerValidateScopeEnabled|global/admin
resourceServerValidateTypeEnabled|global/admin
sempEnabled|global/admin
usernameClaimName|global/admin



This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @return OauthProfileApiApiUpdateOauthProfileRequest
*/
func (a *OauthProfileApiService) UpdateOauthProfile(ctx context.Context, oauthProfileName string) OauthProfileApiApiUpdateOauthProfileRequest {
	return OauthProfileApiApiUpdateOauthProfileRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
	}
}

// Execute executes the request
//  @return OauthProfileResponse
func (a *OauthProfileApiService) UpdateOauthProfileExecute(r OauthProfileApiApiUpdateOauthProfileRequest) (*OauthProfileResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPatch
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.UpdateOauthProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)

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

type OauthProfileApiApiUpdateOauthProfileAccessLevelGroupRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	groupName string
	body *OauthProfileAccessLevelGroup
	opaquePassword *string
	select_ *[]string
}

// The Group Access Level object&#39;s attributes.
func (r OauthProfileApiApiUpdateOauthProfileAccessLevelGroupRequest) Body(body OauthProfileAccessLevelGroup) OauthProfileApiApiUpdateOauthProfileAccessLevelGroupRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiUpdateOauthProfileAccessLevelGroupRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiUpdateOauthProfileAccessLevelGroupRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiUpdateOauthProfileAccessLevelGroupRequest) Select_(select_ []string) OauthProfileApiApiUpdateOauthProfileAccessLevelGroupRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiUpdateOauthProfileAccessLevelGroupRequest) Execute() (*OauthProfileAccessLevelGroupResponse, *http.Response, error) {
	return r.ApiService.UpdateOauthProfileAccessLevelGroupExecute(r)
}

/*
UpdateOauthProfileAccessLevelGroup Update a Group Access Level object.

Update a Group Access Level object. Any attribute missing from the request will be left unchanged.

The name of a group as it exists on the OAuth server being used to authenticate SEMP users.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
groupName|x|x||||
oauthProfileName|x|x||||



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation. Requests which include the following attributes require greater access scope/level:


Attribute|Access Scope/Level
:---|:---:
globalAccessLevel|global/admin



This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param groupName The name of the group.
 @return OauthProfileApiApiUpdateOauthProfileAccessLevelGroupRequest
*/
func (a *OauthProfileApiService) UpdateOauthProfileAccessLevelGroup(ctx context.Context, oauthProfileName string, groupName string) OauthProfileApiApiUpdateOauthProfileAccessLevelGroupRequest {
	return OauthProfileApiApiUpdateOauthProfileAccessLevelGroupRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		groupName: groupName,
	}
}

// Execute executes the request
//  @return OauthProfileAccessLevelGroupResponse
func (a *OauthProfileApiService) UpdateOauthProfileAccessLevelGroupExecute(r OauthProfileApiApiUpdateOauthProfileAccessLevelGroupRequest) (*OauthProfileAccessLevelGroupResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPatch
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileAccessLevelGroupResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.UpdateOauthProfileAccessLevelGroup")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"groupName"+"}", url.PathEscape(parameterToString(r.groupName, "")), -1)

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

type OauthProfileApiApiUpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	groupName string
	msgVpnName string
	body *OauthProfileAccessLevelGroupMsgVpnAccessLevelException
	opaquePassword *string
	select_ *[]string
}

// The Message VPN Access-Level Exception object&#39;s attributes.
func (r OauthProfileApiApiUpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) Body(body OauthProfileAccessLevelGroupMsgVpnAccessLevelException) OauthProfileApiApiUpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiUpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiUpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiUpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) Select_(select_ []string) OauthProfileApiApiUpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiUpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) Execute() (*OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	return r.ApiService.UpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionExecute(r)
}

/*
UpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelException Update a Message VPN Access-Level Exception object.

Update a Message VPN Access-Level Exception object. Any attribute missing from the request will be left unchanged.

Message VPN access-level exceptions for members of this group.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
groupName|x|x||||
msgVpnName|x|x||||
oauthProfileName|x|x||||



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param groupName The name of the group.
 @param msgVpnName The name of the message VPN.
 @return OauthProfileApiApiUpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest
*/
func (a *OauthProfileApiService) UpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelException(ctx context.Context, oauthProfileName string, groupName string, msgVpnName string) OauthProfileApiApiUpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest {
	return OauthProfileApiApiUpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		groupName: groupName,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse
func (a *OauthProfileApiService) UpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionExecute(r OauthProfileApiApiUpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionRequest) (*OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPatch
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileAccessLevelGroupMsgVpnAccessLevelExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.UpdateOauthProfileAccessLevelGroupMsgVpnAccessLevelException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/accessLevelGroups/{groupName}/msgVpnAccessLevelExceptions/{msgVpnName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"groupName"+"}", url.PathEscape(parameterToString(r.groupName, "")), -1)
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

type OauthProfileApiApiUpdateOauthProfileClientAuthorizationParameterRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	authorizationParameterName string
	body *OauthProfileClientAuthorizationParameter
	opaquePassword *string
	select_ *[]string
}

// The Authorization Parameter object&#39;s attributes.
func (r OauthProfileApiApiUpdateOauthProfileClientAuthorizationParameterRequest) Body(body OauthProfileClientAuthorizationParameter) OauthProfileApiApiUpdateOauthProfileClientAuthorizationParameterRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiUpdateOauthProfileClientAuthorizationParameterRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiUpdateOauthProfileClientAuthorizationParameterRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiUpdateOauthProfileClientAuthorizationParameterRequest) Select_(select_ []string) OauthProfileApiApiUpdateOauthProfileClientAuthorizationParameterRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiUpdateOauthProfileClientAuthorizationParameterRequest) Execute() (*OauthProfileClientAuthorizationParameterResponse, *http.Response, error) {
	return r.ApiService.UpdateOauthProfileClientAuthorizationParameterExecute(r)
}

/*
UpdateOauthProfileClientAuthorizationParameter Update an Authorization Parameter object.

Update an Authorization Parameter object. Any attribute missing from the request will be left unchanged.

Additional parameters to be passed to the OAuth authorization endpoint.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
authorizationParameterName|x|x||||
oauthProfileName|x|x||||



A SEMP client authorized with a minimum access scope/level of "global/admin" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param authorizationParameterName The name of the authorization parameter.
 @return OauthProfileApiApiUpdateOauthProfileClientAuthorizationParameterRequest
*/
func (a *OauthProfileApiService) UpdateOauthProfileClientAuthorizationParameter(ctx context.Context, oauthProfileName string, authorizationParameterName string) OauthProfileApiApiUpdateOauthProfileClientAuthorizationParameterRequest {
	return OauthProfileApiApiUpdateOauthProfileClientAuthorizationParameterRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		authorizationParameterName: authorizationParameterName,
	}
}

// Execute executes the request
//  @return OauthProfileClientAuthorizationParameterResponse
func (a *OauthProfileApiService) UpdateOauthProfileClientAuthorizationParameterExecute(r OauthProfileApiApiUpdateOauthProfileClientAuthorizationParameterRequest) (*OauthProfileClientAuthorizationParameterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPatch
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileClientAuthorizationParameterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.UpdateOauthProfileClientAuthorizationParameter")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/clientAuthorizationParameters/{authorizationParameterName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"authorizationParameterName"+"}", url.PathEscape(parameterToString(r.authorizationParameterName, "")), -1)

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

type OauthProfileApiApiUpdateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest struct {
	ctx context.Context
	ApiService *OauthProfileApiService
	oauthProfileName string
	msgVpnName string
	body *OauthProfileDefaultMsgVpnAccessLevelException
	opaquePassword *string
	select_ *[]string
}

// The Message VPN Access-Level Exception object&#39;s attributes.
func (r OauthProfileApiApiUpdateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) Body(body OauthProfileDefaultMsgVpnAccessLevelException) OauthProfileApiApiUpdateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r OauthProfileApiApiUpdateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) OpaquePassword(opaquePassword string) OauthProfileApiApiUpdateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r OauthProfileApiApiUpdateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) Select_(select_ []string) OauthProfileApiApiUpdateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	r.select_ = &select_
	return r
}

func (r OauthProfileApiApiUpdateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) Execute() (*OauthProfileDefaultMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	return r.ApiService.UpdateOauthProfileDefaultMsgVpnAccessLevelExceptionExecute(r)
}

/*
UpdateOauthProfileDefaultMsgVpnAccessLevelException Update a Message VPN Access-Level Exception object.

Update a Message VPN Access-Level Exception object. Any attribute missing from the request will be left unchanged.

Default message VPN access-level exceptions.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
msgVpnName|x|x||||
oauthProfileName|x|x||||



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.24.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param oauthProfileName The name of the OAuth profile.
 @param msgVpnName The name of the message VPN.
 @return OauthProfileApiApiUpdateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest
*/
func (a *OauthProfileApiService) UpdateOauthProfileDefaultMsgVpnAccessLevelException(ctx context.Context, oauthProfileName string, msgVpnName string) OauthProfileApiApiUpdateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest {
	return OauthProfileApiApiUpdateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest{
		ApiService: a,
		ctx: ctx,
		oauthProfileName: oauthProfileName,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return OauthProfileDefaultMsgVpnAccessLevelExceptionResponse
func (a *OauthProfileApiService) UpdateOauthProfileDefaultMsgVpnAccessLevelExceptionExecute(r OauthProfileApiApiUpdateOauthProfileDefaultMsgVpnAccessLevelExceptionRequest) (*OauthProfileDefaultMsgVpnAccessLevelExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPatch
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *OauthProfileDefaultMsgVpnAccessLevelExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OauthProfileApiService.UpdateOauthProfileDefaultMsgVpnAccessLevelException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauthProfiles/{oauthProfileName}/defaultMsgVpnAccessLevelExceptions/{msgVpnName}"
	localVarPath = strings.Replace(localVarPath, "{"+"oauthProfileName"+"}", url.PathEscape(parameterToString(r.oauthProfileName, "")), -1)
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
