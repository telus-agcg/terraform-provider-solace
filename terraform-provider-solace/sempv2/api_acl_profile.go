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

// AclProfileApiService AclProfileApi service
type AclProfileApiService service

type AclProfileApiApiCreateMsgVpnAclProfileRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	body *MsgVpnAclProfile
	opaquePassword *string
	select_ *[]string
}

// The ACL Profile object&#39;s attributes.
func (r AclProfileApiApiCreateMsgVpnAclProfileRequest) Body(body MsgVpnAclProfile) AclProfileApiApiCreateMsgVpnAclProfileRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiCreateMsgVpnAclProfileRequest) OpaquePassword(opaquePassword string) AclProfileApiApiCreateMsgVpnAclProfileRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiCreateMsgVpnAclProfileRequest) Select_(select_ []string) AclProfileApiApiCreateMsgVpnAclProfileRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiCreateMsgVpnAclProfileRequest) Execute() (*MsgVpnAclProfileResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnAclProfileExecute(r)
}

/*
CreateMsgVpnAclProfile Create an ACL Profile object.

Create an ACL Profile object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

An ACL Profile controls whether an authenticated client is permitted to establish a connection with the message broker or permitted to publish and subscribe to specific topics.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
aclProfileName|x|x||||
msgVpnName|x||x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @return AclProfileApiApiCreateMsgVpnAclProfileRequest
*/
func (a *AclProfileApiService) CreateMsgVpnAclProfile(ctx context.Context, msgVpnName string) AclProfileApiApiCreateMsgVpnAclProfileRequest {
	return AclProfileApiApiCreateMsgVpnAclProfileRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileResponse
func (a *AclProfileApiService) CreateMsgVpnAclProfileExecute(r AclProfileApiApiCreateMsgVpnAclProfileRequest) (*MsgVpnAclProfileResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.CreateMsgVpnAclProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles"
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

type AclProfileApiApiCreateMsgVpnAclProfileClientConnectExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	body *MsgVpnAclProfileClientConnectException
	opaquePassword *string
	select_ *[]string
}

// The Client Connect Exception object&#39;s attributes.
func (r AclProfileApiApiCreateMsgVpnAclProfileClientConnectExceptionRequest) Body(body MsgVpnAclProfileClientConnectException) AclProfileApiApiCreateMsgVpnAclProfileClientConnectExceptionRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiCreateMsgVpnAclProfileClientConnectExceptionRequest) OpaquePassword(opaquePassword string) AclProfileApiApiCreateMsgVpnAclProfileClientConnectExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiCreateMsgVpnAclProfileClientConnectExceptionRequest) Select_(select_ []string) AclProfileApiApiCreateMsgVpnAclProfileClientConnectExceptionRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiCreateMsgVpnAclProfileClientConnectExceptionRequest) Execute() (*MsgVpnAclProfileClientConnectExceptionResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnAclProfileClientConnectExceptionExecute(r)
}

/*
CreateMsgVpnAclProfileClientConnectException Create a Client Connect Exception object.

Create a Client Connect Exception object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Client Connect Exception is an exception to the default action to take when a client using the ACL Profile connects to the Message VPN. Exceptions must be expressed as an IP address/netmask in CIDR form.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
aclProfileName|x||x|||
clientConnectExceptionAddress|x|x||||
msgVpnName|x||x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiCreateMsgVpnAclProfileClientConnectExceptionRequest
*/
func (a *AclProfileApiService) CreateMsgVpnAclProfileClientConnectException(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiCreateMsgVpnAclProfileClientConnectExceptionRequest {
	return AclProfileApiApiCreateMsgVpnAclProfileClientConnectExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileClientConnectExceptionResponse
func (a *AclProfileApiService) CreateMsgVpnAclProfileClientConnectExceptionExecute(r AclProfileApiApiCreateMsgVpnAclProfileClientConnectExceptionRequest) (*MsgVpnAclProfileClientConnectExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileClientConnectExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.CreateMsgVpnAclProfileClientConnectException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiCreateMsgVpnAclProfilePublishExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	body *MsgVpnAclProfilePublishException
	opaquePassword *string
	select_ *[]string
}

// The Publish Topic Exception object&#39;s attributes.
func (r AclProfileApiApiCreateMsgVpnAclProfilePublishExceptionRequest) Body(body MsgVpnAclProfilePublishException) AclProfileApiApiCreateMsgVpnAclProfilePublishExceptionRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiCreateMsgVpnAclProfilePublishExceptionRequest) OpaquePassword(opaquePassword string) AclProfileApiApiCreateMsgVpnAclProfilePublishExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiCreateMsgVpnAclProfilePublishExceptionRequest) Select_(select_ []string) AclProfileApiApiCreateMsgVpnAclProfilePublishExceptionRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiCreateMsgVpnAclProfilePublishExceptionRequest) Execute() (*MsgVpnAclProfilePublishExceptionResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnAclProfilePublishExceptionExecute(r)
}

/*
CreateMsgVpnAclProfilePublishException Create a Publish Topic Exception object.

Create a Publish Topic Exception object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Publish Topic Exception is an exception to the default action to take when a client using the ACL Profile publishes to a topic in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
aclProfileName|x||x||x|
msgVpnName|x||x||x|
publishExceptionTopic|x|x|||x|
topicSyntax|x|x|||x|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been deprecated since 2.14. Replaced by publishTopicExceptions.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiCreateMsgVpnAclProfilePublishExceptionRequest

Deprecated
*/
func (a *AclProfileApiService) CreateMsgVpnAclProfilePublishException(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiCreateMsgVpnAclProfilePublishExceptionRequest {
	return AclProfileApiApiCreateMsgVpnAclProfilePublishExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfilePublishExceptionResponse
// Deprecated
func (a *AclProfileApiService) CreateMsgVpnAclProfilePublishExceptionExecute(r AclProfileApiApiCreateMsgVpnAclProfilePublishExceptionRequest) (*MsgVpnAclProfilePublishExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfilePublishExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.CreateMsgVpnAclProfilePublishException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiCreateMsgVpnAclProfilePublishTopicExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	body *MsgVpnAclProfilePublishTopicException
	opaquePassword *string
	select_ *[]string
}

// The Publish Topic Exception object&#39;s attributes.
func (r AclProfileApiApiCreateMsgVpnAclProfilePublishTopicExceptionRequest) Body(body MsgVpnAclProfilePublishTopicException) AclProfileApiApiCreateMsgVpnAclProfilePublishTopicExceptionRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiCreateMsgVpnAclProfilePublishTopicExceptionRequest) OpaquePassword(opaquePassword string) AclProfileApiApiCreateMsgVpnAclProfilePublishTopicExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiCreateMsgVpnAclProfilePublishTopicExceptionRequest) Select_(select_ []string) AclProfileApiApiCreateMsgVpnAclProfilePublishTopicExceptionRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiCreateMsgVpnAclProfilePublishTopicExceptionRequest) Execute() (*MsgVpnAclProfilePublishTopicExceptionResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnAclProfilePublishTopicExceptionExecute(r)
}

/*
CreateMsgVpnAclProfilePublishTopicException Create a Publish Topic Exception object.

Create a Publish Topic Exception object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Publish Topic Exception is an exception to the default action to take when a client using the ACL Profile publishes to a topic in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
aclProfileName|x||x|||
msgVpnName|x||x|||
publishTopicException|x|x||||
publishTopicExceptionSyntax|x|x||||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.14.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiCreateMsgVpnAclProfilePublishTopicExceptionRequest
*/
func (a *AclProfileApiService) CreateMsgVpnAclProfilePublishTopicException(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiCreateMsgVpnAclProfilePublishTopicExceptionRequest {
	return AclProfileApiApiCreateMsgVpnAclProfilePublishTopicExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfilePublishTopicExceptionResponse
func (a *AclProfileApiService) CreateMsgVpnAclProfilePublishTopicExceptionExecute(r AclProfileApiApiCreateMsgVpnAclProfilePublishTopicExceptionRequest) (*MsgVpnAclProfilePublishTopicExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfilePublishTopicExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.CreateMsgVpnAclProfilePublishTopicException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiCreateMsgVpnAclProfileSubscribeExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	body *MsgVpnAclProfileSubscribeException
	opaquePassword *string
	select_ *[]string
}

// The Subscribe Topic Exception object&#39;s attributes.
func (r AclProfileApiApiCreateMsgVpnAclProfileSubscribeExceptionRequest) Body(body MsgVpnAclProfileSubscribeException) AclProfileApiApiCreateMsgVpnAclProfileSubscribeExceptionRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiCreateMsgVpnAclProfileSubscribeExceptionRequest) OpaquePassword(opaquePassword string) AclProfileApiApiCreateMsgVpnAclProfileSubscribeExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiCreateMsgVpnAclProfileSubscribeExceptionRequest) Select_(select_ []string) AclProfileApiApiCreateMsgVpnAclProfileSubscribeExceptionRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiCreateMsgVpnAclProfileSubscribeExceptionRequest) Execute() (*MsgVpnAclProfileSubscribeExceptionResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnAclProfileSubscribeExceptionExecute(r)
}

/*
CreateMsgVpnAclProfileSubscribeException Create a Subscribe Topic Exception object.

Create a Subscribe Topic Exception object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Subscribe Topic Exception is an exception to the default action to take when a client using the ACL Profile subscribes to a topic in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
aclProfileName|x||x||x|
msgVpnName|x||x||x|
subscribeExceptionTopic|x|x|||x|
topicSyntax|x|x|||x|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been deprecated since 2.14. Replaced by subscribeTopicExceptions.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiCreateMsgVpnAclProfileSubscribeExceptionRequest

Deprecated
*/
func (a *AclProfileApiService) CreateMsgVpnAclProfileSubscribeException(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiCreateMsgVpnAclProfileSubscribeExceptionRequest {
	return AclProfileApiApiCreateMsgVpnAclProfileSubscribeExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileSubscribeExceptionResponse
// Deprecated
func (a *AclProfileApiService) CreateMsgVpnAclProfileSubscribeExceptionExecute(r AclProfileApiApiCreateMsgVpnAclProfileSubscribeExceptionRequest) (*MsgVpnAclProfileSubscribeExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileSubscribeExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.CreateMsgVpnAclProfileSubscribeException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiCreateMsgVpnAclProfileSubscribeShareNameExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	body *MsgVpnAclProfileSubscribeShareNameException
	opaquePassword *string
	select_ *[]string
}

// The Subscribe Share Name Exception object&#39;s attributes.
func (r AclProfileApiApiCreateMsgVpnAclProfileSubscribeShareNameExceptionRequest) Body(body MsgVpnAclProfileSubscribeShareNameException) AclProfileApiApiCreateMsgVpnAclProfileSubscribeShareNameExceptionRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiCreateMsgVpnAclProfileSubscribeShareNameExceptionRequest) OpaquePassword(opaquePassword string) AclProfileApiApiCreateMsgVpnAclProfileSubscribeShareNameExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiCreateMsgVpnAclProfileSubscribeShareNameExceptionRequest) Select_(select_ []string) AclProfileApiApiCreateMsgVpnAclProfileSubscribeShareNameExceptionRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiCreateMsgVpnAclProfileSubscribeShareNameExceptionRequest) Execute() (*MsgVpnAclProfileSubscribeShareNameExceptionResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnAclProfileSubscribeShareNameExceptionExecute(r)
}

/*
CreateMsgVpnAclProfileSubscribeShareNameException Create a Subscribe Share Name Exception object.

Create a Subscribe Share Name Exception object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Subscribe Share Name Exception is an exception to the default action to take when a client using the ACL Profile subscribes to a share-name subscription in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
aclProfileName|x||x|||
msgVpnName|x||x|||
subscribeShareNameException|x|x||||
subscribeShareNameExceptionSyntax|x|x||||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.14.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiCreateMsgVpnAclProfileSubscribeShareNameExceptionRequest
*/
func (a *AclProfileApiService) CreateMsgVpnAclProfileSubscribeShareNameException(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiCreateMsgVpnAclProfileSubscribeShareNameExceptionRequest {
	return AclProfileApiApiCreateMsgVpnAclProfileSubscribeShareNameExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileSubscribeShareNameExceptionResponse
func (a *AclProfileApiService) CreateMsgVpnAclProfileSubscribeShareNameExceptionExecute(r AclProfileApiApiCreateMsgVpnAclProfileSubscribeShareNameExceptionRequest) (*MsgVpnAclProfileSubscribeShareNameExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileSubscribeShareNameExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.CreateMsgVpnAclProfileSubscribeShareNameException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiCreateMsgVpnAclProfileSubscribeTopicExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	body *MsgVpnAclProfileSubscribeTopicException
	opaquePassword *string
	select_ *[]string
}

// The Subscribe Topic Exception object&#39;s attributes.
func (r AclProfileApiApiCreateMsgVpnAclProfileSubscribeTopicExceptionRequest) Body(body MsgVpnAclProfileSubscribeTopicException) AclProfileApiApiCreateMsgVpnAclProfileSubscribeTopicExceptionRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiCreateMsgVpnAclProfileSubscribeTopicExceptionRequest) OpaquePassword(opaquePassword string) AclProfileApiApiCreateMsgVpnAclProfileSubscribeTopicExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiCreateMsgVpnAclProfileSubscribeTopicExceptionRequest) Select_(select_ []string) AclProfileApiApiCreateMsgVpnAclProfileSubscribeTopicExceptionRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiCreateMsgVpnAclProfileSubscribeTopicExceptionRequest) Execute() (*MsgVpnAclProfileSubscribeTopicExceptionResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnAclProfileSubscribeTopicExceptionExecute(r)
}

/*
CreateMsgVpnAclProfileSubscribeTopicException Create a Subscribe Topic Exception object.

Create a Subscribe Topic Exception object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Subscribe Topic Exception is an exception to the default action to take when a client using the ACL Profile subscribes to a topic in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
aclProfileName|x||x|||
msgVpnName|x||x|||
subscribeTopicException|x|x||||
subscribeTopicExceptionSyntax|x|x||||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.14.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiCreateMsgVpnAclProfileSubscribeTopicExceptionRequest
*/
func (a *AclProfileApiService) CreateMsgVpnAclProfileSubscribeTopicException(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiCreateMsgVpnAclProfileSubscribeTopicExceptionRequest {
	return AclProfileApiApiCreateMsgVpnAclProfileSubscribeTopicExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileSubscribeTopicExceptionResponse
func (a *AclProfileApiService) CreateMsgVpnAclProfileSubscribeTopicExceptionExecute(r AclProfileApiApiCreateMsgVpnAclProfileSubscribeTopicExceptionRequest) (*MsgVpnAclProfileSubscribeTopicExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileSubscribeTopicExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.CreateMsgVpnAclProfileSubscribeTopicException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiDeleteMsgVpnAclProfileRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
}


func (r AclProfileApiApiDeleteMsgVpnAclProfileRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnAclProfileExecute(r)
}

/*
DeleteMsgVpnAclProfile Delete an ACL Profile object.

Delete an ACL Profile object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

An ACL Profile controls whether an authenticated client is permitted to establish a connection with the message broker or permitted to publish and subscribe to specific topics.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiDeleteMsgVpnAclProfileRequest
*/
func (a *AclProfileApiService) DeleteMsgVpnAclProfile(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiDeleteMsgVpnAclProfileRequest {
	return AclProfileApiApiDeleteMsgVpnAclProfileRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *AclProfileApiService) DeleteMsgVpnAclProfileExecute(r AclProfileApiApiDeleteMsgVpnAclProfileRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.DeleteMsgVpnAclProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiDeleteMsgVpnAclProfileClientConnectExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	clientConnectExceptionAddress string
}


func (r AclProfileApiApiDeleteMsgVpnAclProfileClientConnectExceptionRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnAclProfileClientConnectExceptionExecute(r)
}

/*
DeleteMsgVpnAclProfileClientConnectException Delete a Client Connect Exception object.

Delete a Client Connect Exception object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Client Connect Exception is an exception to the default action to take when a client using the ACL Profile connects to the Message VPN. Exceptions must be expressed as an IP address/netmask in CIDR form.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @param clientConnectExceptionAddress The IP address/netmask of the client connect exception in CIDR form.
 @return AclProfileApiApiDeleteMsgVpnAclProfileClientConnectExceptionRequest
*/
func (a *AclProfileApiService) DeleteMsgVpnAclProfileClientConnectException(ctx context.Context, msgVpnName string, aclProfileName string, clientConnectExceptionAddress string) AclProfileApiApiDeleteMsgVpnAclProfileClientConnectExceptionRequest {
	return AclProfileApiApiDeleteMsgVpnAclProfileClientConnectExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
		clientConnectExceptionAddress: clientConnectExceptionAddress,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *AclProfileApiService) DeleteMsgVpnAclProfileClientConnectExceptionExecute(r AclProfileApiApiDeleteMsgVpnAclProfileClientConnectExceptionRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.DeleteMsgVpnAclProfileClientConnectException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions/{clientConnectExceptionAddress}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clientConnectExceptionAddress"+"}", url.PathEscape(parameterToString(r.clientConnectExceptionAddress, "")), -1)

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

type AclProfileApiApiDeleteMsgVpnAclProfilePublishExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	topicSyntax string
	publishExceptionTopic string
}


func (r AclProfileApiApiDeleteMsgVpnAclProfilePublishExceptionRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnAclProfilePublishExceptionExecute(r)
}

/*
DeleteMsgVpnAclProfilePublishException Delete a Publish Topic Exception object.

Delete a Publish Topic Exception object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Publish Topic Exception is an exception to the default action to take when a client using the ACL Profile publishes to a topic in the Message VPN. Exceptions must be expressed as a topic.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been deprecated since 2.14. Replaced by publishTopicExceptions.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @param topicSyntax The syntax of the topic for the exception to the default action taken.
 @param publishExceptionTopic The topic for the exception to the default action taken. May include wildcard characters.
 @return AclProfileApiApiDeleteMsgVpnAclProfilePublishExceptionRequest

Deprecated
*/
func (a *AclProfileApiService) DeleteMsgVpnAclProfilePublishException(ctx context.Context, msgVpnName string, aclProfileName string, topicSyntax string, publishExceptionTopic string) AclProfileApiApiDeleteMsgVpnAclProfilePublishExceptionRequest {
	return AclProfileApiApiDeleteMsgVpnAclProfilePublishExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
		topicSyntax: topicSyntax,
		publishExceptionTopic: publishExceptionTopic,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
// Deprecated
func (a *AclProfileApiService) DeleteMsgVpnAclProfilePublishExceptionExecute(r AclProfileApiApiDeleteMsgVpnAclProfilePublishExceptionRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.DeleteMsgVpnAclProfilePublishException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions/{topicSyntax},{publishExceptionTopic}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"topicSyntax"+"}", url.PathEscape(parameterToString(r.topicSyntax, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"publishExceptionTopic"+"}", url.PathEscape(parameterToString(r.publishExceptionTopic, "")), -1)

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

type AclProfileApiApiDeleteMsgVpnAclProfilePublishTopicExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	publishTopicExceptionSyntax string
	publishTopicException string
}


func (r AclProfileApiApiDeleteMsgVpnAclProfilePublishTopicExceptionRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnAclProfilePublishTopicExceptionExecute(r)
}

/*
DeleteMsgVpnAclProfilePublishTopicException Delete a Publish Topic Exception object.

Delete a Publish Topic Exception object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Publish Topic Exception is an exception to the default action to take when a client using the ACL Profile publishes to a topic in the Message VPN. Exceptions must be expressed as a topic.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.14.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @param publishTopicExceptionSyntax The syntax of the topic for the exception to the default action taken.
 @param publishTopicException The topic for the exception to the default action taken. May include wildcard characters.
 @return AclProfileApiApiDeleteMsgVpnAclProfilePublishTopicExceptionRequest
*/
func (a *AclProfileApiService) DeleteMsgVpnAclProfilePublishTopicException(ctx context.Context, msgVpnName string, aclProfileName string, publishTopicExceptionSyntax string, publishTopicException string) AclProfileApiApiDeleteMsgVpnAclProfilePublishTopicExceptionRequest {
	return AclProfileApiApiDeleteMsgVpnAclProfilePublishTopicExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
		publishTopicExceptionSyntax: publishTopicExceptionSyntax,
		publishTopicException: publishTopicException,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *AclProfileApiService) DeleteMsgVpnAclProfilePublishTopicExceptionExecute(r AclProfileApiApiDeleteMsgVpnAclProfilePublishTopicExceptionRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.DeleteMsgVpnAclProfilePublishTopicException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions/{publishTopicExceptionSyntax},{publishTopicException}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"publishTopicExceptionSyntax"+"}", url.PathEscape(parameterToString(r.publishTopicExceptionSyntax, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"publishTopicException"+"}", url.PathEscape(parameterToString(r.publishTopicException, "")), -1)

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

type AclProfileApiApiDeleteMsgVpnAclProfileSubscribeExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	topicSyntax string
	subscribeExceptionTopic string
}


func (r AclProfileApiApiDeleteMsgVpnAclProfileSubscribeExceptionRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnAclProfileSubscribeExceptionExecute(r)
}

/*
DeleteMsgVpnAclProfileSubscribeException Delete a Subscribe Topic Exception object.

Delete a Subscribe Topic Exception object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Subscribe Topic Exception is an exception to the default action to take when a client using the ACL Profile subscribes to a topic in the Message VPN. Exceptions must be expressed as a topic.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been deprecated since 2.14. Replaced by subscribeTopicExceptions.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @param topicSyntax The syntax of the topic for the exception to the default action taken.
 @param subscribeExceptionTopic The topic for the exception to the default action taken. May include wildcard characters.
 @return AclProfileApiApiDeleteMsgVpnAclProfileSubscribeExceptionRequest

Deprecated
*/
func (a *AclProfileApiService) DeleteMsgVpnAclProfileSubscribeException(ctx context.Context, msgVpnName string, aclProfileName string, topicSyntax string, subscribeExceptionTopic string) AclProfileApiApiDeleteMsgVpnAclProfileSubscribeExceptionRequest {
	return AclProfileApiApiDeleteMsgVpnAclProfileSubscribeExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
		topicSyntax: topicSyntax,
		subscribeExceptionTopic: subscribeExceptionTopic,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
// Deprecated
func (a *AclProfileApiService) DeleteMsgVpnAclProfileSubscribeExceptionExecute(r AclProfileApiApiDeleteMsgVpnAclProfileSubscribeExceptionRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.DeleteMsgVpnAclProfileSubscribeException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions/{topicSyntax},{subscribeExceptionTopic}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"topicSyntax"+"}", url.PathEscape(parameterToString(r.topicSyntax, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"subscribeExceptionTopic"+"}", url.PathEscape(parameterToString(r.subscribeExceptionTopic, "")), -1)

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

type AclProfileApiApiDeleteMsgVpnAclProfileSubscribeShareNameExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	subscribeShareNameExceptionSyntax string
	subscribeShareNameException string
}


func (r AclProfileApiApiDeleteMsgVpnAclProfileSubscribeShareNameExceptionRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnAclProfileSubscribeShareNameExceptionExecute(r)
}

/*
DeleteMsgVpnAclProfileSubscribeShareNameException Delete a Subscribe Share Name Exception object.

Delete a Subscribe Share Name Exception object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Subscribe Share Name Exception is an exception to the default action to take when a client using the ACL Profile subscribes to a share-name subscription in the Message VPN. Exceptions must be expressed as a topic.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.14.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @param subscribeShareNameExceptionSyntax The syntax of the subscribe share name for the exception to the default action taken.
 @param subscribeShareNameException The subscribe share name exception to the default action taken. May include wildcard characters.
 @return AclProfileApiApiDeleteMsgVpnAclProfileSubscribeShareNameExceptionRequest
*/
func (a *AclProfileApiService) DeleteMsgVpnAclProfileSubscribeShareNameException(ctx context.Context, msgVpnName string, aclProfileName string, subscribeShareNameExceptionSyntax string, subscribeShareNameException string) AclProfileApiApiDeleteMsgVpnAclProfileSubscribeShareNameExceptionRequest {
	return AclProfileApiApiDeleteMsgVpnAclProfileSubscribeShareNameExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
		subscribeShareNameExceptionSyntax: subscribeShareNameExceptionSyntax,
		subscribeShareNameException: subscribeShareNameException,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *AclProfileApiService) DeleteMsgVpnAclProfileSubscribeShareNameExceptionExecute(r AclProfileApiApiDeleteMsgVpnAclProfileSubscribeShareNameExceptionRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.DeleteMsgVpnAclProfileSubscribeShareNameException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions/{subscribeShareNameExceptionSyntax},{subscribeShareNameException}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"subscribeShareNameExceptionSyntax"+"}", url.PathEscape(parameterToString(r.subscribeShareNameExceptionSyntax, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"subscribeShareNameException"+"}", url.PathEscape(parameterToString(r.subscribeShareNameException, "")), -1)

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

type AclProfileApiApiDeleteMsgVpnAclProfileSubscribeTopicExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	subscribeTopicExceptionSyntax string
	subscribeTopicException string
}


func (r AclProfileApiApiDeleteMsgVpnAclProfileSubscribeTopicExceptionRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnAclProfileSubscribeTopicExceptionExecute(r)
}

/*
DeleteMsgVpnAclProfileSubscribeTopicException Delete a Subscribe Topic Exception object.

Delete a Subscribe Topic Exception object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Subscribe Topic Exception is an exception to the default action to take when a client using the ACL Profile subscribes to a topic in the Message VPN. Exceptions must be expressed as a topic.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.14.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @param subscribeTopicExceptionSyntax The syntax of the topic for the exception to the default action taken.
 @param subscribeTopicException The topic for the exception to the default action taken. May include wildcard characters.
 @return AclProfileApiApiDeleteMsgVpnAclProfileSubscribeTopicExceptionRequest
*/
func (a *AclProfileApiService) DeleteMsgVpnAclProfileSubscribeTopicException(ctx context.Context, msgVpnName string, aclProfileName string, subscribeTopicExceptionSyntax string, subscribeTopicException string) AclProfileApiApiDeleteMsgVpnAclProfileSubscribeTopicExceptionRequest {
	return AclProfileApiApiDeleteMsgVpnAclProfileSubscribeTopicExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
		subscribeTopicExceptionSyntax: subscribeTopicExceptionSyntax,
		subscribeTopicException: subscribeTopicException,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *AclProfileApiService) DeleteMsgVpnAclProfileSubscribeTopicExceptionExecute(r AclProfileApiApiDeleteMsgVpnAclProfileSubscribeTopicExceptionRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.DeleteMsgVpnAclProfileSubscribeTopicException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions/{subscribeTopicExceptionSyntax},{subscribeTopicException}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"subscribeTopicExceptionSyntax"+"}", url.PathEscape(parameterToString(r.subscribeTopicExceptionSyntax, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"subscribeTopicException"+"}", url.PathEscape(parameterToString(r.subscribeTopicException, "")), -1)

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

type AclProfileApiApiGetMsgVpnAclProfileRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileRequest) OpaquePassword(opaquePassword string) AclProfileApiApiGetMsgVpnAclProfileRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileRequest) Select_(select_ []string) AclProfileApiApiGetMsgVpnAclProfileRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiGetMsgVpnAclProfileRequest) Execute() (*MsgVpnAclProfileResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnAclProfileExecute(r)
}

/*
GetMsgVpnAclProfile Get an ACL Profile object.

Get an ACL Profile object.

An ACL Profile controls whether an authenticated client is permitted to establish a connection with the message broker or permitted to publish and subscribe to specific topics.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
aclProfileName|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiGetMsgVpnAclProfileRequest
*/
func (a *AclProfileApiService) GetMsgVpnAclProfile(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiGetMsgVpnAclProfileRequest {
	return AclProfileApiApiGetMsgVpnAclProfileRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileResponse
func (a *AclProfileApiService) GetMsgVpnAclProfileExecute(r AclProfileApiApiGetMsgVpnAclProfileRequest) (*MsgVpnAclProfileResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.GetMsgVpnAclProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	clientConnectExceptionAddress string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionRequest) OpaquePassword(opaquePassword string) AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionRequest) Select_(select_ []string) AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionRequest) Execute() (*MsgVpnAclProfileClientConnectExceptionResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnAclProfileClientConnectExceptionExecute(r)
}

/*
GetMsgVpnAclProfileClientConnectException Get a Client Connect Exception object.

Get a Client Connect Exception object.

A Client Connect Exception is an exception to the default action to take when a client using the ACL Profile connects to the Message VPN. Exceptions must be expressed as an IP address/netmask in CIDR form.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
aclProfileName|x|||
clientConnectExceptionAddress|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @param clientConnectExceptionAddress The IP address/netmask of the client connect exception in CIDR form.
 @return AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionRequest
*/
func (a *AclProfileApiService) GetMsgVpnAclProfileClientConnectException(ctx context.Context, msgVpnName string, aclProfileName string, clientConnectExceptionAddress string) AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionRequest {
	return AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
		clientConnectExceptionAddress: clientConnectExceptionAddress,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileClientConnectExceptionResponse
func (a *AclProfileApiService) GetMsgVpnAclProfileClientConnectExceptionExecute(r AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionRequest) (*MsgVpnAclProfileClientConnectExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileClientConnectExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.GetMsgVpnAclProfileClientConnectException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions/{clientConnectExceptionAddress}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clientConnectExceptionAddress"+"}", url.PathEscape(parameterToString(r.clientConnectExceptionAddress, "")), -1)

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

type AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest) Count(count int32) AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest) Cursor(cursor string) AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest) OpaquePassword(opaquePassword string) AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest) Where(where []string) AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest) Select_(select_ []string) AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest) Execute() (*MsgVpnAclProfileClientConnectExceptionsResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnAclProfileClientConnectExceptionsExecute(r)
}

/*
GetMsgVpnAclProfileClientConnectExceptions Get a list of Client Connect Exception objects.

Get a list of Client Connect Exception objects.

A Client Connect Exception is an exception to the default action to take when a client using the ACL Profile connects to the Message VPN. Exceptions must be expressed as an IP address/netmask in CIDR form.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
aclProfileName|x|||
clientConnectExceptionAddress|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest
*/
func (a *AclProfileApiService) GetMsgVpnAclProfileClientConnectExceptions(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest {
	return AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileClientConnectExceptionsResponse
func (a *AclProfileApiService) GetMsgVpnAclProfileClientConnectExceptionsExecute(r AclProfileApiApiGetMsgVpnAclProfileClientConnectExceptionsRequest) (*MsgVpnAclProfileClientConnectExceptionsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileClientConnectExceptionsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.GetMsgVpnAclProfileClientConnectExceptions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/clientConnectExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiGetMsgVpnAclProfilePublishExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	topicSyntax string
	publishExceptionTopic string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilePublishExceptionRequest) OpaquePassword(opaquePassword string) AclProfileApiApiGetMsgVpnAclProfilePublishExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilePublishExceptionRequest) Select_(select_ []string) AclProfileApiApiGetMsgVpnAclProfilePublishExceptionRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiGetMsgVpnAclProfilePublishExceptionRequest) Execute() (*MsgVpnAclProfilePublishExceptionResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnAclProfilePublishExceptionExecute(r)
}

/*
GetMsgVpnAclProfilePublishException Get a Publish Topic Exception object.

Get a Publish Topic Exception object.

A Publish Topic Exception is an exception to the default action to take when a client using the ACL Profile publishes to a topic in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
aclProfileName|x||x|
msgVpnName|x||x|
publishExceptionTopic|x||x|
topicSyntax|x||x|



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been deprecated since 2.14. Replaced by publishTopicExceptions.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @param topicSyntax The syntax of the topic for the exception to the default action taken.
 @param publishExceptionTopic The topic for the exception to the default action taken. May include wildcard characters.
 @return AclProfileApiApiGetMsgVpnAclProfilePublishExceptionRequest

Deprecated
*/
func (a *AclProfileApiService) GetMsgVpnAclProfilePublishException(ctx context.Context, msgVpnName string, aclProfileName string, topicSyntax string, publishExceptionTopic string) AclProfileApiApiGetMsgVpnAclProfilePublishExceptionRequest {
	return AclProfileApiApiGetMsgVpnAclProfilePublishExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
		topicSyntax: topicSyntax,
		publishExceptionTopic: publishExceptionTopic,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfilePublishExceptionResponse
// Deprecated
func (a *AclProfileApiService) GetMsgVpnAclProfilePublishExceptionExecute(r AclProfileApiApiGetMsgVpnAclProfilePublishExceptionRequest) (*MsgVpnAclProfilePublishExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfilePublishExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.GetMsgVpnAclProfilePublishException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions/{topicSyntax},{publishExceptionTopic}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"topicSyntax"+"}", url.PathEscape(parameterToString(r.topicSyntax, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"publishExceptionTopic"+"}", url.PathEscape(parameterToString(r.publishExceptionTopic, "")), -1)

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

type AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest) Count(count int32) AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest) Cursor(cursor string) AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest) OpaquePassword(opaquePassword string) AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest) Where(where []string) AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest) Select_(select_ []string) AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest) Execute() (*MsgVpnAclProfilePublishExceptionsResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnAclProfilePublishExceptionsExecute(r)
}

/*
GetMsgVpnAclProfilePublishExceptions Get a list of Publish Topic Exception objects.

Get a list of Publish Topic Exception objects.

A Publish Topic Exception is an exception to the default action to take when a client using the ACL Profile publishes to a topic in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
aclProfileName|x||x|
msgVpnName|x||x|
publishExceptionTopic|x||x|
topicSyntax|x||x|



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been deprecated since 2.14. Replaced by publishTopicExceptions.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest

Deprecated
*/
func (a *AclProfileApiService) GetMsgVpnAclProfilePublishExceptions(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest {
	return AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfilePublishExceptionsResponse
// Deprecated
func (a *AclProfileApiService) GetMsgVpnAclProfilePublishExceptionsExecute(r AclProfileApiApiGetMsgVpnAclProfilePublishExceptionsRequest) (*MsgVpnAclProfilePublishExceptionsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfilePublishExceptionsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.GetMsgVpnAclProfilePublishExceptions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	publishTopicExceptionSyntax string
	publishTopicException string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionRequest) OpaquePassword(opaquePassword string) AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionRequest) Select_(select_ []string) AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionRequest) Execute() (*MsgVpnAclProfilePublishTopicExceptionResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnAclProfilePublishTopicExceptionExecute(r)
}

/*
GetMsgVpnAclProfilePublishTopicException Get a Publish Topic Exception object.

Get a Publish Topic Exception object.

A Publish Topic Exception is an exception to the default action to take when a client using the ACL Profile publishes to a topic in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
aclProfileName|x|||
msgVpnName|x|||
publishTopicException|x|||
publishTopicExceptionSyntax|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.14.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @param publishTopicExceptionSyntax The syntax of the topic for the exception to the default action taken.
 @param publishTopicException The topic for the exception to the default action taken. May include wildcard characters.
 @return AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionRequest
*/
func (a *AclProfileApiService) GetMsgVpnAclProfilePublishTopicException(ctx context.Context, msgVpnName string, aclProfileName string, publishTopicExceptionSyntax string, publishTopicException string) AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionRequest {
	return AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
		publishTopicExceptionSyntax: publishTopicExceptionSyntax,
		publishTopicException: publishTopicException,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfilePublishTopicExceptionResponse
func (a *AclProfileApiService) GetMsgVpnAclProfilePublishTopicExceptionExecute(r AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionRequest) (*MsgVpnAclProfilePublishTopicExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfilePublishTopicExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.GetMsgVpnAclProfilePublishTopicException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions/{publishTopicExceptionSyntax},{publishTopicException}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"publishTopicExceptionSyntax"+"}", url.PathEscape(parameterToString(r.publishTopicExceptionSyntax, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"publishTopicException"+"}", url.PathEscape(parameterToString(r.publishTopicException, "")), -1)

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

type AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest) Count(count int32) AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest) Cursor(cursor string) AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest) OpaquePassword(opaquePassword string) AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest) Where(where []string) AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest) Select_(select_ []string) AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest) Execute() (*MsgVpnAclProfilePublishTopicExceptionsResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnAclProfilePublishTopicExceptionsExecute(r)
}

/*
GetMsgVpnAclProfilePublishTopicExceptions Get a list of Publish Topic Exception objects.

Get a list of Publish Topic Exception objects.

A Publish Topic Exception is an exception to the default action to take when a client using the ACL Profile publishes to a topic in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
aclProfileName|x|||
msgVpnName|x|||
publishTopicException|x|||
publishTopicExceptionSyntax|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.14.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest
*/
func (a *AclProfileApiService) GetMsgVpnAclProfilePublishTopicExceptions(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest {
	return AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfilePublishTopicExceptionsResponse
func (a *AclProfileApiService) GetMsgVpnAclProfilePublishTopicExceptionsExecute(r AclProfileApiApiGetMsgVpnAclProfilePublishTopicExceptionsRequest) (*MsgVpnAclProfilePublishTopicExceptionsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfilePublishTopicExceptionsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.GetMsgVpnAclProfilePublishTopicExceptions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/publishTopicExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	topicSyntax string
	subscribeExceptionTopic string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionRequest) OpaquePassword(opaquePassword string) AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionRequest) Select_(select_ []string) AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionRequest) Execute() (*MsgVpnAclProfileSubscribeExceptionResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnAclProfileSubscribeExceptionExecute(r)
}

/*
GetMsgVpnAclProfileSubscribeException Get a Subscribe Topic Exception object.

Get a Subscribe Topic Exception object.

A Subscribe Topic Exception is an exception to the default action to take when a client using the ACL Profile subscribes to a topic in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
aclProfileName|x||x|
msgVpnName|x||x|
subscribeExceptionTopic|x||x|
topicSyntax|x||x|



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been deprecated since 2.14. Replaced by subscribeTopicExceptions.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @param topicSyntax The syntax of the topic for the exception to the default action taken.
 @param subscribeExceptionTopic The topic for the exception to the default action taken. May include wildcard characters.
 @return AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionRequest

Deprecated
*/
func (a *AclProfileApiService) GetMsgVpnAclProfileSubscribeException(ctx context.Context, msgVpnName string, aclProfileName string, topicSyntax string, subscribeExceptionTopic string) AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionRequest {
	return AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
		topicSyntax: topicSyntax,
		subscribeExceptionTopic: subscribeExceptionTopic,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileSubscribeExceptionResponse
// Deprecated
func (a *AclProfileApiService) GetMsgVpnAclProfileSubscribeExceptionExecute(r AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionRequest) (*MsgVpnAclProfileSubscribeExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileSubscribeExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.GetMsgVpnAclProfileSubscribeException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions/{topicSyntax},{subscribeExceptionTopic}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"topicSyntax"+"}", url.PathEscape(parameterToString(r.topicSyntax, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"subscribeExceptionTopic"+"}", url.PathEscape(parameterToString(r.subscribeExceptionTopic, "")), -1)

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

type AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest) Count(count int32) AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest) Cursor(cursor string) AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest) OpaquePassword(opaquePassword string) AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest) Where(where []string) AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest) Select_(select_ []string) AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest) Execute() (*MsgVpnAclProfileSubscribeExceptionsResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnAclProfileSubscribeExceptionsExecute(r)
}

/*
GetMsgVpnAclProfileSubscribeExceptions Get a list of Subscribe Topic Exception objects.

Get a list of Subscribe Topic Exception objects.

A Subscribe Topic Exception is an exception to the default action to take when a client using the ACL Profile subscribes to a topic in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
aclProfileName|x||x|
msgVpnName|x||x|
subscribeExceptionTopic|x||x|
topicSyntax|x||x|



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been deprecated since 2.14. Replaced by subscribeTopicExceptions.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest

Deprecated
*/
func (a *AclProfileApiService) GetMsgVpnAclProfileSubscribeExceptions(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest {
	return AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileSubscribeExceptionsResponse
// Deprecated
func (a *AclProfileApiService) GetMsgVpnAclProfileSubscribeExceptionsExecute(r AclProfileApiApiGetMsgVpnAclProfileSubscribeExceptionsRequest) (*MsgVpnAclProfileSubscribeExceptionsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileSubscribeExceptionsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.GetMsgVpnAclProfileSubscribeExceptions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	subscribeShareNameExceptionSyntax string
	subscribeShareNameException string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionRequest) OpaquePassword(opaquePassword string) AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionRequest) Select_(select_ []string) AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionRequest) Execute() (*MsgVpnAclProfileSubscribeShareNameExceptionResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnAclProfileSubscribeShareNameExceptionExecute(r)
}

/*
GetMsgVpnAclProfileSubscribeShareNameException Get a Subscribe Share Name Exception object.

Get a Subscribe Share Name Exception object.

A Subscribe Share Name Exception is an exception to the default action to take when a client using the ACL Profile subscribes to a share-name subscription in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
aclProfileName|x|||
msgVpnName|x|||
subscribeShareNameException|x|||
subscribeShareNameExceptionSyntax|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.14.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @param subscribeShareNameExceptionSyntax The syntax of the subscribe share name for the exception to the default action taken.
 @param subscribeShareNameException The subscribe share name exception to the default action taken. May include wildcard characters.
 @return AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionRequest
*/
func (a *AclProfileApiService) GetMsgVpnAclProfileSubscribeShareNameException(ctx context.Context, msgVpnName string, aclProfileName string, subscribeShareNameExceptionSyntax string, subscribeShareNameException string) AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionRequest {
	return AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
		subscribeShareNameExceptionSyntax: subscribeShareNameExceptionSyntax,
		subscribeShareNameException: subscribeShareNameException,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileSubscribeShareNameExceptionResponse
func (a *AclProfileApiService) GetMsgVpnAclProfileSubscribeShareNameExceptionExecute(r AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionRequest) (*MsgVpnAclProfileSubscribeShareNameExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileSubscribeShareNameExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.GetMsgVpnAclProfileSubscribeShareNameException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions/{subscribeShareNameExceptionSyntax},{subscribeShareNameException}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"subscribeShareNameExceptionSyntax"+"}", url.PathEscape(parameterToString(r.subscribeShareNameExceptionSyntax, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"subscribeShareNameException"+"}", url.PathEscape(parameterToString(r.subscribeShareNameException, "")), -1)

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

type AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest) Count(count int32) AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest) Cursor(cursor string) AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest) OpaquePassword(opaquePassword string) AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest) Where(where []string) AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest) Select_(select_ []string) AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest) Execute() (*MsgVpnAclProfileSubscribeShareNameExceptionsResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnAclProfileSubscribeShareNameExceptionsExecute(r)
}

/*
GetMsgVpnAclProfileSubscribeShareNameExceptions Get a list of Subscribe Share Name Exception objects.

Get a list of Subscribe Share Name Exception objects.

A Subscribe Share Name Exception is an exception to the default action to take when a client using the ACL Profile subscribes to a share-name subscription in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
aclProfileName|x|||
msgVpnName|x|||
subscribeShareNameException|x|||
subscribeShareNameExceptionSyntax|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.14.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest
*/
func (a *AclProfileApiService) GetMsgVpnAclProfileSubscribeShareNameExceptions(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest {
	return AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileSubscribeShareNameExceptionsResponse
func (a *AclProfileApiService) GetMsgVpnAclProfileSubscribeShareNameExceptionsExecute(r AclProfileApiApiGetMsgVpnAclProfileSubscribeShareNameExceptionsRequest) (*MsgVpnAclProfileSubscribeShareNameExceptionsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileSubscribeShareNameExceptionsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.GetMsgVpnAclProfileSubscribeShareNameExceptions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeShareNameExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	subscribeTopicExceptionSyntax string
	subscribeTopicException string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionRequest) OpaquePassword(opaquePassword string) AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionRequest) Select_(select_ []string) AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionRequest) Execute() (*MsgVpnAclProfileSubscribeTopicExceptionResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnAclProfileSubscribeTopicExceptionExecute(r)
}

/*
GetMsgVpnAclProfileSubscribeTopicException Get a Subscribe Topic Exception object.

Get a Subscribe Topic Exception object.

A Subscribe Topic Exception is an exception to the default action to take when a client using the ACL Profile subscribes to a topic in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
aclProfileName|x|||
msgVpnName|x|||
subscribeTopicException|x|||
subscribeTopicExceptionSyntax|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.14.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @param subscribeTopicExceptionSyntax The syntax of the topic for the exception to the default action taken.
 @param subscribeTopicException The topic for the exception to the default action taken. May include wildcard characters.
 @return AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionRequest
*/
func (a *AclProfileApiService) GetMsgVpnAclProfileSubscribeTopicException(ctx context.Context, msgVpnName string, aclProfileName string, subscribeTopicExceptionSyntax string, subscribeTopicException string) AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionRequest {
	return AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
		subscribeTopicExceptionSyntax: subscribeTopicExceptionSyntax,
		subscribeTopicException: subscribeTopicException,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileSubscribeTopicExceptionResponse
func (a *AclProfileApiService) GetMsgVpnAclProfileSubscribeTopicExceptionExecute(r AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionRequest) (*MsgVpnAclProfileSubscribeTopicExceptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileSubscribeTopicExceptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.GetMsgVpnAclProfileSubscribeTopicException")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions/{subscribeTopicExceptionSyntax},{subscribeTopicException}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"subscribeTopicExceptionSyntax"+"}", url.PathEscape(parameterToString(r.subscribeTopicExceptionSyntax, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"subscribeTopicException"+"}", url.PathEscape(parameterToString(r.subscribeTopicException, "")), -1)

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

type AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest) Count(count int32) AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest) Cursor(cursor string) AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest) OpaquePassword(opaquePassword string) AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest) Where(where []string) AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest) Select_(select_ []string) AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest) Execute() (*MsgVpnAclProfileSubscribeTopicExceptionsResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnAclProfileSubscribeTopicExceptionsExecute(r)
}

/*
GetMsgVpnAclProfileSubscribeTopicExceptions Get a list of Subscribe Topic Exception objects.

Get a list of Subscribe Topic Exception objects.

A Subscribe Topic Exception is an exception to the default action to take when a client using the ACL Profile subscribes to a topic in the Message VPN. Exceptions must be expressed as a topic.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
aclProfileName|x|||
msgVpnName|x|||
subscribeTopicException|x|||
subscribeTopicExceptionSyntax|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.14.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest
*/
func (a *AclProfileApiService) GetMsgVpnAclProfileSubscribeTopicExceptions(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest {
	return AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileSubscribeTopicExceptionsResponse
func (a *AclProfileApiService) GetMsgVpnAclProfileSubscribeTopicExceptionsExecute(r AclProfileApiApiGetMsgVpnAclProfileSubscribeTopicExceptionsRequest) (*MsgVpnAclProfileSubscribeTopicExceptionsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileSubscribeTopicExceptionsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.GetMsgVpnAclProfileSubscribeTopicExceptions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}/subscribeTopicExceptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiGetMsgVpnAclProfilesRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilesRequest) Count(count int32) AclProfileApiApiGetMsgVpnAclProfilesRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilesRequest) Cursor(cursor string) AclProfileApiApiGetMsgVpnAclProfilesRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilesRequest) OpaquePassword(opaquePassword string) AclProfileApiApiGetMsgVpnAclProfilesRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilesRequest) Where(where []string) AclProfileApiApiGetMsgVpnAclProfilesRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiGetMsgVpnAclProfilesRequest) Select_(select_ []string) AclProfileApiApiGetMsgVpnAclProfilesRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiGetMsgVpnAclProfilesRequest) Execute() (*MsgVpnAclProfilesResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnAclProfilesExecute(r)
}

/*
GetMsgVpnAclProfiles Get a list of ACL Profile objects.

Get a list of ACL Profile objects.

An ACL Profile controls whether an authenticated client is permitted to establish a connection with the message broker or permitted to publish and subscribe to specific topics.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
aclProfileName|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @return AclProfileApiApiGetMsgVpnAclProfilesRequest
*/
func (a *AclProfileApiService) GetMsgVpnAclProfiles(ctx context.Context, msgVpnName string) AclProfileApiApiGetMsgVpnAclProfilesRequest {
	return AclProfileApiApiGetMsgVpnAclProfilesRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfilesResponse
func (a *AclProfileApiService) GetMsgVpnAclProfilesExecute(r AclProfileApiApiGetMsgVpnAclProfilesRequest) (*MsgVpnAclProfilesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfilesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.GetMsgVpnAclProfiles")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles"
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

type AclProfileApiApiReplaceMsgVpnAclProfileRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	body *MsgVpnAclProfile
	opaquePassword *string
	select_ *[]string
}

// The ACL Profile object&#39;s attributes.
func (r AclProfileApiApiReplaceMsgVpnAclProfileRequest) Body(body MsgVpnAclProfile) AclProfileApiApiReplaceMsgVpnAclProfileRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiReplaceMsgVpnAclProfileRequest) OpaquePassword(opaquePassword string) AclProfileApiApiReplaceMsgVpnAclProfileRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiReplaceMsgVpnAclProfileRequest) Select_(select_ []string) AclProfileApiApiReplaceMsgVpnAclProfileRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiReplaceMsgVpnAclProfileRequest) Execute() (*MsgVpnAclProfileResponse, *http.Response, error) {
	return r.ApiService.ReplaceMsgVpnAclProfileExecute(r)
}

/*
ReplaceMsgVpnAclProfile Replace an ACL Profile object.

Replace an ACL Profile object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

An ACL Profile controls whether an authenticated client is permitted to establish a connection with the message broker or permitted to publish and subscribe to specific topics.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
aclProfileName|x|x||||
msgVpnName|x|x||||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiReplaceMsgVpnAclProfileRequest
*/
func (a *AclProfileApiService) ReplaceMsgVpnAclProfile(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiReplaceMsgVpnAclProfileRequest {
	return AclProfileApiApiReplaceMsgVpnAclProfileRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileResponse
func (a *AclProfileApiService) ReplaceMsgVpnAclProfileExecute(r AclProfileApiApiReplaceMsgVpnAclProfileRequest) (*MsgVpnAclProfileResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.ReplaceMsgVpnAclProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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

type AclProfileApiApiUpdateMsgVpnAclProfileRequest struct {
	ctx context.Context
	ApiService *AclProfileApiService
	msgVpnName string
	aclProfileName string
	body *MsgVpnAclProfile
	opaquePassword *string
	select_ *[]string
}

// The ACL Profile object&#39;s attributes.
func (r AclProfileApiApiUpdateMsgVpnAclProfileRequest) Body(body MsgVpnAclProfile) AclProfileApiApiUpdateMsgVpnAclProfileRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r AclProfileApiApiUpdateMsgVpnAclProfileRequest) OpaquePassword(opaquePassword string) AclProfileApiApiUpdateMsgVpnAclProfileRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r AclProfileApiApiUpdateMsgVpnAclProfileRequest) Select_(select_ []string) AclProfileApiApiUpdateMsgVpnAclProfileRequest {
	r.select_ = &select_
	return r
}

func (r AclProfileApiApiUpdateMsgVpnAclProfileRequest) Execute() (*MsgVpnAclProfileResponse, *http.Response, error) {
	return r.ApiService.UpdateMsgVpnAclProfileExecute(r)
}

/*
UpdateMsgVpnAclProfile Update an ACL Profile object.

Update an ACL Profile object. Any attribute missing from the request will be left unchanged.

An ACL Profile controls whether an authenticated client is permitted to establish a connection with the message broker or permitted to publish and subscribe to specific topics.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
aclProfileName|x|x||||
msgVpnName|x|x||||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.0.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param aclProfileName The name of the ACL Profile.
 @return AclProfileApiApiUpdateMsgVpnAclProfileRequest
*/
func (a *AclProfileApiService) UpdateMsgVpnAclProfile(ctx context.Context, msgVpnName string, aclProfileName string) AclProfileApiApiUpdateMsgVpnAclProfileRequest {
	return AclProfileApiApiUpdateMsgVpnAclProfileRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		aclProfileName: aclProfileName,
	}
}

// Execute executes the request
//  @return MsgVpnAclProfileResponse
func (a *AclProfileApiService) UpdateMsgVpnAclProfileExecute(r AclProfileApiApiUpdateMsgVpnAclProfileRequest) (*MsgVpnAclProfileResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPatch
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnAclProfileResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AclProfileApiService.UpdateMsgVpnAclProfile")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/aclProfiles/{aclProfileName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"aclProfileName"+"}", url.PathEscape(parameterToString(r.aclProfileName, "")), -1)

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
