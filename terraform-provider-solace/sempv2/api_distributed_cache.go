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

// DistributedCacheApiService DistributedCacheApi service
type DistributedCacheApiService service

type DistributedCacheApiApiCreateMsgVpnDistributedCacheRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	body *MsgVpnDistributedCache
	opaquePassword *string
	select_ *[]string
}

// The Distributed Cache object&#39;s attributes.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheRequest) Body(body MsgVpnDistributedCache) DistributedCacheApiApiCreateMsgVpnDistributedCacheRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiCreateMsgVpnDistributedCacheRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheRequest) Select_(select_ []string) DistributedCacheApiApiCreateMsgVpnDistributedCacheRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheRequest) Execute() (*MsgVpnDistributedCacheResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnDistributedCacheExecute(r)
}

/*
CreateMsgVpnDistributedCache Create a Distributed Cache object.

Create a Distributed Cache object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Distributed Cache is a collection of one or more Cache Clusters that belong to the same Message VPN. Each Cache Cluster in a Distributed Cache is configured to subscribe to a different set of topics. This effectively divides up the configured topic space, to provide scaling to very large topic spaces or very high cached message throughput.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
cacheName|x|x||||
msgVpnName|x||x|||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
MsgVpnDistributedCache|scheduledDeleteMsgDayList|scheduledDeleteMsgTimeList|
MsgVpnDistributedCache|scheduledDeleteMsgTimeList|scheduledDeleteMsgDayList|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @return DistributedCacheApiApiCreateMsgVpnDistributedCacheRequest
*/
func (a *DistributedCacheApiService) CreateMsgVpnDistributedCache(ctx context.Context, msgVpnName string) DistributedCacheApiApiCreateMsgVpnDistributedCacheRequest {
	return DistributedCacheApiApiCreateMsgVpnDistributedCacheRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheResponse
func (a *DistributedCacheApiService) CreateMsgVpnDistributedCacheExecute(r DistributedCacheApiApiCreateMsgVpnDistributedCacheRequest) (*MsgVpnDistributedCacheResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.CreateMsgVpnDistributedCache")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches"
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

type DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	body *MsgVpnDistributedCacheCluster
	opaquePassword *string
	select_ *[]string
}

// The Cache Cluster object&#39;s attributes.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterRequest) Body(body MsgVpnDistributedCacheCluster) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterRequest) Select_(select_ []string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterRequest) Execute() (*MsgVpnDistributedCacheClusterResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnDistributedCacheClusterExecute(r)
}

/*
CreateMsgVpnDistributedCacheCluster Create a Cache Cluster object.

Create a Cache Cluster object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Cache Cluster is a collection of one or more Cache Instances that subscribe to exactly the same topics. Cache Instances are grouped together in a Cache Cluster for the purpose of fault tolerance and load balancing. As published messages are received, the message broker message bus sends these live data messages to the Cache Instances in the Cache Cluster. This enables client cache requests to be served by any of Cache Instances in the Cache Cluster.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
cacheName|x||x|||
clusterName|x|x||||
msgVpnName|x||x|||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
EventThresholdByPercent|clearPercent|setPercent|
EventThresholdByPercent|setPercent|clearPercent|
EventThresholdByValue|clearValue|setValue|
EventThresholdByValue|setValue|clearValue|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @return DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterRequest
*/
func (a *DistributedCacheApiService) CreateMsgVpnDistributedCacheCluster(ctx context.Context, msgVpnName string, cacheName string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterRequest {
	return DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterResponse
func (a *DistributedCacheApiService) CreateMsgVpnDistributedCacheClusterExecute(r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterRequest) (*MsgVpnDistributedCacheClusterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.CreateMsgVpnDistributedCacheCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)

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

type DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	body *MsgVpnDistributedCacheClusterGlobalCachingHomeCluster
	opaquePassword *string
	select_ *[]string
}

// The Home Cache Cluster object&#39;s attributes.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest) Body(body MsgVpnDistributedCacheClusterGlobalCachingHomeCluster) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest) Select_(select_ []string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest) Execute() (*MsgVpnDistributedCacheClusterGlobalCachingHomeClusterResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterExecute(r)
}

/*
CreateMsgVpnDistributedCacheClusterGlobalCachingHomeCluster Create a Home Cache Cluster object.

Create a Home Cache Cluster object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Home Cache Cluster is a Cache Cluster that is the "definitive" Cache Cluster for a given topic in the context of the Global Caching feature.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
cacheName|x||x|||
clusterName|x||x|||
homeClusterName|x|x||||
msgVpnName|x||x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @return DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest
*/
func (a *DistributedCacheApiService) CreateMsgVpnDistributedCacheClusterGlobalCachingHomeCluster(ctx context.Context, msgVpnName string, cacheName string, clusterName string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest {
	return DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterGlobalCachingHomeClusterResponse
func (a *DistributedCacheApiService) CreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterExecute(r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest) (*MsgVpnDistributedCacheClusterGlobalCachingHomeClusterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterGlobalCachingHomeClusterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.CreateMsgVpnDistributedCacheClusterGlobalCachingHomeCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)

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

type DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	homeClusterName string
	body *MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix
	opaquePassword *string
	select_ *[]string
}

// The Topic Prefix object&#39;s attributes.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest) Body(body MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest) Select_(select_ []string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest) Execute() (*MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixExecute(r)
}

/*
CreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix Create a Topic Prefix object.

Create a Topic Prefix object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Topic Prefix is a prefix for a global topic that is available from the containing Home Cache Cluster.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
cacheName|x||x|||
clusterName|x||x|||
homeClusterName|x||x|||
msgVpnName|x||x|||
topicPrefix|x|x||||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @param homeClusterName The name of the remote Home Cache Cluster.
 @return DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest
*/
func (a *DistributedCacheApiService) CreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix(ctx context.Context, msgVpnName string, cacheName string, clusterName string, homeClusterName string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest {
	return DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
		homeClusterName: homeClusterName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixResponse
func (a *DistributedCacheApiService) CreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixExecute(r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest) (*MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.CreateMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"homeClusterName"+"}", url.PathEscape(parameterToString(r.homeClusterName, "")), -1)

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

type DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterInstanceRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	body *MsgVpnDistributedCacheClusterInstance
	opaquePassword *string
	select_ *[]string
}

// The Cache Instance object&#39;s attributes.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterInstanceRequest) Body(body MsgVpnDistributedCacheClusterInstance) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterInstanceRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterInstanceRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterInstanceRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterInstanceRequest) Select_(select_ []string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterInstanceRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterInstanceRequest) Execute() (*MsgVpnDistributedCacheClusterInstanceResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnDistributedCacheClusterInstanceExecute(r)
}

/*
CreateMsgVpnDistributedCacheClusterInstance Create a Cache Instance object.

Create a Cache Instance object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Cache Instance is a single Cache process that belongs to a single Cache Cluster. A Cache Instance object provisioned on the broker is used to disseminate configuration information to the Cache process. Cache Instances listen for and cache live data messages that match the topic subscriptions configured for their parent Cache Cluster.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
cacheName|x||x|||
clusterName|x||x|||
instanceName|x|x||||
msgVpnName|x||x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @return DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterInstanceRequest
*/
func (a *DistributedCacheApiService) CreateMsgVpnDistributedCacheClusterInstance(ctx context.Context, msgVpnName string, cacheName string, clusterName string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterInstanceRequest {
	return DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterInstanceRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterInstanceResponse
func (a *DistributedCacheApiService) CreateMsgVpnDistributedCacheClusterInstanceExecute(r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterInstanceRequest) (*MsgVpnDistributedCacheClusterInstanceResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterInstanceResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.CreateMsgVpnDistributedCacheClusterInstance")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)

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

type DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterTopicRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	body *MsgVpnDistributedCacheClusterTopic
	opaquePassword *string
	select_ *[]string
}

// The Topic object&#39;s attributes.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterTopicRequest) Body(body MsgVpnDistributedCacheClusterTopic) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterTopicRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterTopicRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterTopicRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterTopicRequest) Select_(select_ []string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterTopicRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterTopicRequest) Execute() (*MsgVpnDistributedCacheClusterTopicResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnDistributedCacheClusterTopicExecute(r)
}

/*
CreateMsgVpnDistributedCacheClusterTopic Create a Topic object.

Create a Topic object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

The Cache Instances that belong to the containing Cache Cluster will cache any messages published to topics that match a Topic Subscription.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
cacheName|x||x|||
clusterName|x||x|||
msgVpnName|x||x|||
topic|x|x||||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @return DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterTopicRequest
*/
func (a *DistributedCacheApiService) CreateMsgVpnDistributedCacheClusterTopic(ctx context.Context, msgVpnName string, cacheName string, clusterName string) DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterTopicRequest {
	return DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterTopicRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterTopicResponse
func (a *DistributedCacheApiService) CreateMsgVpnDistributedCacheClusterTopicExecute(r DistributedCacheApiApiCreateMsgVpnDistributedCacheClusterTopicRequest) (*MsgVpnDistributedCacheClusterTopicResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterTopicResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.CreateMsgVpnDistributedCacheClusterTopic")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)

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

type DistributedCacheApiApiDeleteMsgVpnDistributedCacheRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
}


func (r DistributedCacheApiApiDeleteMsgVpnDistributedCacheRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnDistributedCacheExecute(r)
}

/*
DeleteMsgVpnDistributedCache Delete a Distributed Cache object.

Delete a Distributed Cache object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Distributed Cache is a collection of one or more Cache Clusters that belong to the same Message VPN. Each Cache Cluster in a Distributed Cache is configured to subscribe to a different set of topics. This effectively divides up the configured topic space, to provide scaling to very large topic spaces or very high cached message throughput.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @return DistributedCacheApiApiDeleteMsgVpnDistributedCacheRequest
*/
func (a *DistributedCacheApiService) DeleteMsgVpnDistributedCache(ctx context.Context, msgVpnName string, cacheName string) DistributedCacheApiApiDeleteMsgVpnDistributedCacheRequest {
	return DistributedCacheApiApiDeleteMsgVpnDistributedCacheRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *DistributedCacheApiService) DeleteMsgVpnDistributedCacheExecute(r DistributedCacheApiApiDeleteMsgVpnDistributedCacheRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.DeleteMsgVpnDistributedCache")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)

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

type DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
}


func (r DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnDistributedCacheClusterExecute(r)
}

/*
DeleteMsgVpnDistributedCacheCluster Delete a Cache Cluster object.

Delete a Cache Cluster object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Cache Cluster is a collection of one or more Cache Instances that subscribe to exactly the same topics. Cache Instances are grouped together in a Cache Cluster for the purpose of fault tolerance and load balancing. As published messages are received, the message broker message bus sends these live data messages to the Cache Instances in the Cache Cluster. This enables client cache requests to be served by any of Cache Instances in the Cache Cluster.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @return DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterRequest
*/
func (a *DistributedCacheApiService) DeleteMsgVpnDistributedCacheCluster(ctx context.Context, msgVpnName string, cacheName string, clusterName string) DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterRequest {
	return DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *DistributedCacheApiService) DeleteMsgVpnDistributedCacheClusterExecute(r DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.DeleteMsgVpnDistributedCacheCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)

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

type DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	homeClusterName string
}


func (r DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterExecute(r)
}

/*
DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeCluster Delete a Home Cache Cluster object.

Delete a Home Cache Cluster object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Home Cache Cluster is a Cache Cluster that is the "definitive" Cache Cluster for a given topic in the context of the Global Caching feature.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @param homeClusterName The name of the remote Home Cache Cluster.
 @return DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest
*/
func (a *DistributedCacheApiService) DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeCluster(ctx context.Context, msgVpnName string, cacheName string, clusterName string, homeClusterName string) DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest {
	return DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
		homeClusterName: homeClusterName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *DistributedCacheApiService) DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterExecute(r DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"homeClusterName"+"}", url.PathEscape(parameterToString(r.homeClusterName, "")), -1)

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

type DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	homeClusterName string
	topicPrefix string
}


func (r DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixExecute(r)
}

/*
DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix Delete a Topic Prefix object.

Delete a Topic Prefix object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Topic Prefix is a prefix for a global topic that is available from the containing Home Cache Cluster.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @param homeClusterName The name of the remote Home Cache Cluster.
 @param topicPrefix A topic prefix for global topics available from the remote Home Cache Cluster. A wildcard (/>) is implied at the end of the prefix.
 @return DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest
*/
func (a *DistributedCacheApiService) DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix(ctx context.Context, msgVpnName string, cacheName string, clusterName string, homeClusterName string, topicPrefix string) DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest {
	return DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
		homeClusterName: homeClusterName,
		topicPrefix: topicPrefix,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *DistributedCacheApiService) DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixExecute(r DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.DeleteMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes/{topicPrefix}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"homeClusterName"+"}", url.PathEscape(parameterToString(r.homeClusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"topicPrefix"+"}", url.PathEscape(parameterToString(r.topicPrefix, "")), -1)

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

type DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterInstanceRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	instanceName string
}


func (r DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterInstanceRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnDistributedCacheClusterInstanceExecute(r)
}

/*
DeleteMsgVpnDistributedCacheClusterInstance Delete a Cache Instance object.

Delete a Cache Instance object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

A Cache Instance is a single Cache process that belongs to a single Cache Cluster. A Cache Instance object provisioned on the broker is used to disseminate configuration information to the Cache process. Cache Instances listen for and cache live data messages that match the topic subscriptions configured for their parent Cache Cluster.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @param instanceName The name of the Cache Instance.
 @return DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterInstanceRequest
*/
func (a *DistributedCacheApiService) DeleteMsgVpnDistributedCacheClusterInstance(ctx context.Context, msgVpnName string, cacheName string, clusterName string, instanceName string) DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterInstanceRequest {
	return DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterInstanceRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
		instanceName: instanceName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *DistributedCacheApiService) DeleteMsgVpnDistributedCacheClusterInstanceExecute(r DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterInstanceRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.DeleteMsgVpnDistributedCacheClusterInstance")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"instanceName"+"}", url.PathEscape(parameterToString(r.instanceName, "")), -1)

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

type DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterTopicRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	topic string
}


func (r DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterTopicRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnDistributedCacheClusterTopicExecute(r)
}

/*
DeleteMsgVpnDistributedCacheClusterTopic Delete a Topic object.

Delete a Topic object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

The Cache Instances that belong to the containing Cache Cluster will cache any messages published to topics that match a Topic Subscription.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @param topic The value of the Topic in the form a/b/c.
 @return DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterTopicRequest
*/
func (a *DistributedCacheApiService) DeleteMsgVpnDistributedCacheClusterTopic(ctx context.Context, msgVpnName string, cacheName string, clusterName string, topic string) DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterTopicRequest {
	return DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterTopicRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
		topic: topic,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *DistributedCacheApiService) DeleteMsgVpnDistributedCacheClusterTopicExecute(r DistributedCacheApiApiDeleteMsgVpnDistributedCacheClusterTopicRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.DeleteMsgVpnDistributedCacheClusterTopic")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics/{topic}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"topic"+"}", url.PathEscape(parameterToString(r.topic, "")), -1)

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

type DistributedCacheApiApiGetMsgVpnDistributedCacheRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiGetMsgVpnDistributedCacheRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheRequest) Select_(select_ []string) DistributedCacheApiApiGetMsgVpnDistributedCacheRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiGetMsgVpnDistributedCacheRequest) Execute() (*MsgVpnDistributedCacheResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnDistributedCacheExecute(r)
}

/*
GetMsgVpnDistributedCache Get a Distributed Cache object.

Get a Distributed Cache object.

A Distributed Cache is a collection of one or more Cache Clusters that belong to the same Message VPN. Each Cache Cluster in a Distributed Cache is configured to subscribe to a different set of topics. This effectively divides up the configured topic space, to provide scaling to very large topic spaces or very high cached message throughput.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
cacheName|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @return DistributedCacheApiApiGetMsgVpnDistributedCacheRequest
*/
func (a *DistributedCacheApiService) GetMsgVpnDistributedCache(ctx context.Context, msgVpnName string, cacheName string) DistributedCacheApiApiGetMsgVpnDistributedCacheRequest {
	return DistributedCacheApiApiGetMsgVpnDistributedCacheRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheResponse
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheExecute(r DistributedCacheApiApiGetMsgVpnDistributedCacheRequest) (*MsgVpnDistributedCacheResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.GetMsgVpnDistributedCache")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)

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

type DistributedCacheApiApiGetMsgVpnDistributedCacheClusterRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterRequest) Select_(select_ []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterRequest) Execute() (*MsgVpnDistributedCacheClusterResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnDistributedCacheClusterExecute(r)
}

/*
GetMsgVpnDistributedCacheCluster Get a Cache Cluster object.

Get a Cache Cluster object.

A Cache Cluster is a collection of one or more Cache Instances that subscribe to exactly the same topics. Cache Instances are grouped together in a Cache Cluster for the purpose of fault tolerance and load balancing. As published messages are received, the message broker message bus sends these live data messages to the Cache Instances in the Cache Cluster. This enables client cache requests to be served by any of Cache Instances in the Cache Cluster.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
cacheName|x|||
clusterName|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterRequest
*/
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheCluster(ctx context.Context, msgVpnName string, cacheName string, clusterName string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterRequest {
	return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterResponse
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterExecute(r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterRequest) (*MsgVpnDistributedCacheClusterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.GetMsgVpnDistributedCacheCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)

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

type DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	homeClusterName string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest) Select_(select_ []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest) Execute() (*MsgVpnDistributedCacheClusterGlobalCachingHomeClusterResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterExecute(r)
}

/*
GetMsgVpnDistributedCacheClusterGlobalCachingHomeCluster Get a Home Cache Cluster object.

Get a Home Cache Cluster object.

A Home Cache Cluster is a Cache Cluster that is the "definitive" Cache Cluster for a given topic in the context of the Global Caching feature.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
cacheName|x|||
clusterName|x|||
homeClusterName|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @param homeClusterName The name of the remote Home Cache Cluster.
 @return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest
*/
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterGlobalCachingHomeCluster(ctx context.Context, msgVpnName string, cacheName string, clusterName string, homeClusterName string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest {
	return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
		homeClusterName: homeClusterName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterGlobalCachingHomeClusterResponse
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterExecute(r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterRequest) (*MsgVpnDistributedCacheClusterGlobalCachingHomeClusterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterGlobalCachingHomeClusterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.GetMsgVpnDistributedCacheClusterGlobalCachingHomeCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"homeClusterName"+"}", url.PathEscape(parameterToString(r.homeClusterName, "")), -1)

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

type DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	homeClusterName string
	topicPrefix string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest) Select_(select_ []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest) Execute() (*MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixExecute(r)
}

/*
GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix Get a Topic Prefix object.

Get a Topic Prefix object.

A Topic Prefix is a prefix for a global topic that is available from the containing Home Cache Cluster.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
cacheName|x|||
clusterName|x|||
homeClusterName|x|||
msgVpnName|x|||
topicPrefix|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @param homeClusterName The name of the remote Home Cache Cluster.
 @param topicPrefix A topic prefix for global topics available from the remote Home Cache Cluster. A wildcard (/>) is implied at the end of the prefix.
 @return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest
*/
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix(ctx context.Context, msgVpnName string, cacheName string, clusterName string, homeClusterName string, topicPrefix string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest {
	return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
		homeClusterName: homeClusterName,
		topicPrefix: topicPrefix,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixResponse
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixExecute(r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixRequest) (*MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefix")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes/{topicPrefix}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"homeClusterName"+"}", url.PathEscape(parameterToString(r.homeClusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"topicPrefix"+"}", url.PathEscape(parameterToString(r.topicPrefix, "")), -1)

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

type DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	homeClusterName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest) Count(count int32) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest) Cursor(cursor string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest) Where(where []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest) Select_(select_ []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest) Execute() (*MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesExecute(r)
}

/*
GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixes Get a list of Topic Prefix objects.

Get a list of Topic Prefix objects.

A Topic Prefix is a prefix for a global topic that is available from the containing Home Cache Cluster.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
cacheName|x|||
clusterName|x|||
homeClusterName|x|||
msgVpnName|x|||
topicPrefix|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @param homeClusterName The name of the remote Home Cache Cluster.
 @return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest
*/
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixes(ctx context.Context, msgVpnName string, cacheName string, clusterName string, homeClusterName string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest {
	return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
		homeClusterName: homeClusterName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesResponse
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesExecute(r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesRequest) (*MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusterTopicPrefixes")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters/{homeClusterName}/topicPrefixes"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"homeClusterName"+"}", url.PathEscape(parameterToString(r.homeClusterName, "")), -1)

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

type DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest) Count(count int32) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest) Cursor(cursor string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest) Where(where []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest) Select_(select_ []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest) Execute() (*MsgVpnDistributedCacheClusterGlobalCachingHomeClustersResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersExecute(r)
}

/*
GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusters Get a list of Home Cache Cluster objects.

Get a list of Home Cache Cluster objects.

A Home Cache Cluster is a Cache Cluster that is the "definitive" Cache Cluster for a given topic in the context of the Global Caching feature.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
cacheName|x|||
clusterName|x|||
homeClusterName|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest
*/
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusters(ctx context.Context, msgVpnName string, cacheName string, clusterName string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest {
	return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterGlobalCachingHomeClustersResponse
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersExecute(r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterGlobalCachingHomeClustersRequest) (*MsgVpnDistributedCacheClusterGlobalCachingHomeClustersResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterGlobalCachingHomeClustersResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.GetMsgVpnDistributedCacheClusterGlobalCachingHomeClusters")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/globalCachingHomeClusters"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)

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

type DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstanceRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	instanceName string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstanceRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstanceRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstanceRequest) Select_(select_ []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstanceRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstanceRequest) Execute() (*MsgVpnDistributedCacheClusterInstanceResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnDistributedCacheClusterInstanceExecute(r)
}

/*
GetMsgVpnDistributedCacheClusterInstance Get a Cache Instance object.

Get a Cache Instance object.

A Cache Instance is a single Cache process that belongs to a single Cache Cluster. A Cache Instance object provisioned on the broker is used to disseminate configuration information to the Cache process. Cache Instances listen for and cache live data messages that match the topic subscriptions configured for their parent Cache Cluster.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
cacheName|x|||
clusterName|x|||
instanceName|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @param instanceName The name of the Cache Instance.
 @return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstanceRequest
*/
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterInstance(ctx context.Context, msgVpnName string, cacheName string, clusterName string, instanceName string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstanceRequest {
	return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstanceRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
		instanceName: instanceName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterInstanceResponse
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterInstanceExecute(r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstanceRequest) (*MsgVpnDistributedCacheClusterInstanceResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterInstanceResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.GetMsgVpnDistributedCacheClusterInstance")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"instanceName"+"}", url.PathEscape(parameterToString(r.instanceName, "")), -1)

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

type DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest) Count(count int32) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest) Cursor(cursor string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest) Where(where []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest) Select_(select_ []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest) Execute() (*MsgVpnDistributedCacheClusterInstancesResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnDistributedCacheClusterInstancesExecute(r)
}

/*
GetMsgVpnDistributedCacheClusterInstances Get a list of Cache Instance objects.

Get a list of Cache Instance objects.

A Cache Instance is a single Cache process that belongs to a single Cache Cluster. A Cache Instance object provisioned on the broker is used to disseminate configuration information to the Cache process. Cache Instances listen for and cache live data messages that match the topic subscriptions configured for their parent Cache Cluster.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
cacheName|x|||
clusterName|x|||
instanceName|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest
*/
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterInstances(ctx context.Context, msgVpnName string, cacheName string, clusterName string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest {
	return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterInstancesResponse
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterInstancesExecute(r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterInstancesRequest) (*MsgVpnDistributedCacheClusterInstancesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterInstancesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.GetMsgVpnDistributedCacheClusterInstances")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)

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

type DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	topic string
	opaquePassword *string
	select_ *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicRequest) Select_(select_ []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicRequest) Execute() (*MsgVpnDistributedCacheClusterTopicResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnDistributedCacheClusterTopicExecute(r)
}

/*
GetMsgVpnDistributedCacheClusterTopic Get a Topic object.

Get a Topic object.

The Cache Instances that belong to the containing Cache Cluster will cache any messages published to topics that match a Topic Subscription.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
cacheName|x|||
clusterName|x|||
msgVpnName|x|||
topic|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @param topic The value of the Topic in the form a/b/c.
 @return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicRequest
*/
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterTopic(ctx context.Context, msgVpnName string, cacheName string, clusterName string, topic string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicRequest {
	return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
		topic: topic,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterTopicResponse
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterTopicExecute(r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicRequest) (*MsgVpnDistributedCacheClusterTopicResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterTopicResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.GetMsgVpnDistributedCacheClusterTopic")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics/{topic}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"topic"+"}", url.PathEscape(parameterToString(r.topic, "")), -1)

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

type DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest) Count(count int32) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest) Cursor(cursor string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest) Where(where []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest) Select_(select_ []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest) Execute() (*MsgVpnDistributedCacheClusterTopicsResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnDistributedCacheClusterTopicsExecute(r)
}

/*
GetMsgVpnDistributedCacheClusterTopics Get a list of Topic objects.

Get a list of Topic objects.

The Cache Instances that belong to the containing Cache Cluster will cache any messages published to topics that match a Topic Subscription.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
cacheName|x|||
clusterName|x|||
msgVpnName|x|||
topic|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest
*/
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterTopics(ctx context.Context, msgVpnName string, cacheName string, clusterName string) DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest {
	return DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterTopicsResponse
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusterTopicsExecute(r DistributedCacheApiApiGetMsgVpnDistributedCacheClusterTopicsRequest) (*MsgVpnDistributedCacheClusterTopicsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterTopicsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.GetMsgVpnDistributedCacheClusterTopics")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/topics"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)

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

type DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest) Count(count int32) DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest) Cursor(cursor string) DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest) Where(where []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest) Select_(select_ []string) DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest) Execute() (*MsgVpnDistributedCacheClustersResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnDistributedCacheClustersExecute(r)
}

/*
GetMsgVpnDistributedCacheClusters Get a list of Cache Cluster objects.

Get a list of Cache Cluster objects.

A Cache Cluster is a collection of one or more Cache Instances that subscribe to exactly the same topics. Cache Instances are grouped together in a Cache Cluster for the purpose of fault tolerance and load balancing. As published messages are received, the message broker message bus sends these live data messages to the Cache Instances in the Cache Cluster. This enables client cache requests to be served by any of Cache Instances in the Cache Cluster.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
cacheName|x|||
clusterName|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @return DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest
*/
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClusters(ctx context.Context, msgVpnName string, cacheName string) DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest {
	return DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClustersResponse
func (a *DistributedCacheApiService) GetMsgVpnDistributedCacheClustersExecute(r DistributedCacheApiApiGetMsgVpnDistributedCacheClustersRequest) (*MsgVpnDistributedCacheClustersResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClustersResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.GetMsgVpnDistributedCacheClusters")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)

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

type DistributedCacheApiApiGetMsgVpnDistributedCachesRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	count *int32
	cursor *string
	opaquePassword *string
	where *[]string
	select_ *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCachesRequest) Count(count int32) DistributedCacheApiApiGetMsgVpnDistributedCachesRequest {
	r.count = &count
	return r
}
// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCachesRequest) Cursor(cursor string) DistributedCacheApiApiGetMsgVpnDistributedCachesRequest {
	r.cursor = &cursor
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCachesRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiGetMsgVpnDistributedCachesRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCachesRequest) Where(where []string) DistributedCacheApiApiGetMsgVpnDistributedCachesRequest {
	r.where = &where
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiGetMsgVpnDistributedCachesRequest) Select_(select_ []string) DistributedCacheApiApiGetMsgVpnDistributedCachesRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiGetMsgVpnDistributedCachesRequest) Execute() (*MsgVpnDistributedCachesResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnDistributedCachesExecute(r)
}

/*
GetMsgVpnDistributedCaches Get a list of Distributed Cache objects.

Get a list of Distributed Cache objects.

A Distributed Cache is a collection of one or more Cache Clusters that belong to the same Message VPN. Each Cache Cluster in a Distributed Cache is configured to subscribe to a different set of topics. This effectively divides up the configured topic space, to provide scaling to very large topic spaces or very high cached message throughput.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
cacheName|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @return DistributedCacheApiApiGetMsgVpnDistributedCachesRequest
*/
func (a *DistributedCacheApiService) GetMsgVpnDistributedCaches(ctx context.Context, msgVpnName string) DistributedCacheApiApiGetMsgVpnDistributedCachesRequest {
	return DistributedCacheApiApiGetMsgVpnDistributedCachesRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCachesResponse
func (a *DistributedCacheApiService) GetMsgVpnDistributedCachesExecute(r DistributedCacheApiApiGetMsgVpnDistributedCachesRequest) (*MsgVpnDistributedCachesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCachesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.GetMsgVpnDistributedCaches")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches"
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

type DistributedCacheApiApiReplaceMsgVpnDistributedCacheRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	body *MsgVpnDistributedCache
	opaquePassword *string
	select_ *[]string
}

// The Distributed Cache object&#39;s attributes.
func (r DistributedCacheApiApiReplaceMsgVpnDistributedCacheRequest) Body(body MsgVpnDistributedCache) DistributedCacheApiApiReplaceMsgVpnDistributedCacheRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiReplaceMsgVpnDistributedCacheRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiReplaceMsgVpnDistributedCacheRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiReplaceMsgVpnDistributedCacheRequest) Select_(select_ []string) DistributedCacheApiApiReplaceMsgVpnDistributedCacheRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiReplaceMsgVpnDistributedCacheRequest) Execute() (*MsgVpnDistributedCacheResponse, *http.Response, error) {
	return r.ApiService.ReplaceMsgVpnDistributedCacheExecute(r)
}

/*
ReplaceMsgVpnDistributedCache Replace a Distributed Cache object.

Replace a Distributed Cache object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

A Distributed Cache is a collection of one or more Cache Clusters that belong to the same Message VPN. Each Cache Cluster in a Distributed Cache is configured to subscribe to a different set of topics. This effectively divides up the configured topic space, to provide scaling to very large topic spaces or very high cached message throughput.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
cacheName|x|x||||
msgVpnName|x|x||||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
MsgVpnDistributedCache|scheduledDeleteMsgDayList|scheduledDeleteMsgTimeList|
MsgVpnDistributedCache|scheduledDeleteMsgTimeList|scheduledDeleteMsgDayList|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @return DistributedCacheApiApiReplaceMsgVpnDistributedCacheRequest
*/
func (a *DistributedCacheApiService) ReplaceMsgVpnDistributedCache(ctx context.Context, msgVpnName string, cacheName string) DistributedCacheApiApiReplaceMsgVpnDistributedCacheRequest {
	return DistributedCacheApiApiReplaceMsgVpnDistributedCacheRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheResponse
func (a *DistributedCacheApiService) ReplaceMsgVpnDistributedCacheExecute(r DistributedCacheApiApiReplaceMsgVpnDistributedCacheRequest) (*MsgVpnDistributedCacheResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.ReplaceMsgVpnDistributedCache")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)

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

type DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	body *MsgVpnDistributedCacheCluster
	opaquePassword *string
	select_ *[]string
}

// The Cache Cluster object&#39;s attributes.
func (r DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterRequest) Body(body MsgVpnDistributedCacheCluster) DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterRequest) Select_(select_ []string) DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterRequest) Execute() (*MsgVpnDistributedCacheClusterResponse, *http.Response, error) {
	return r.ApiService.ReplaceMsgVpnDistributedCacheClusterExecute(r)
}

/*
ReplaceMsgVpnDistributedCacheCluster Replace a Cache Cluster object.

Replace a Cache Cluster object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

A Cache Cluster is a collection of one or more Cache Instances that subscribe to exactly the same topics. Cache Instances are grouped together in a Cache Cluster for the purpose of fault tolerance and load balancing. As published messages are received, the message broker message bus sends these live data messages to the Cache Instances in the Cache Cluster. This enables client cache requests to be served by any of Cache Instances in the Cache Cluster.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
cacheName|x|x||||
clusterName|x|x||||
msgVpnName|x|x||||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
EventThresholdByPercent|clearPercent|setPercent|
EventThresholdByPercent|setPercent|clearPercent|
EventThresholdByValue|clearValue|setValue|
EventThresholdByValue|setValue|clearValue|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @return DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterRequest
*/
func (a *DistributedCacheApiService) ReplaceMsgVpnDistributedCacheCluster(ctx context.Context, msgVpnName string, cacheName string, clusterName string) DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterRequest {
	return DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterResponse
func (a *DistributedCacheApiService) ReplaceMsgVpnDistributedCacheClusterExecute(r DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterRequest) (*MsgVpnDistributedCacheClusterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.ReplaceMsgVpnDistributedCacheCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)

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

type DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterInstanceRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	instanceName string
	body *MsgVpnDistributedCacheClusterInstance
	opaquePassword *string
	select_ *[]string
}

// The Cache Instance object&#39;s attributes.
func (r DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterInstanceRequest) Body(body MsgVpnDistributedCacheClusterInstance) DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterInstanceRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterInstanceRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterInstanceRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterInstanceRequest) Select_(select_ []string) DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterInstanceRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterInstanceRequest) Execute() (*MsgVpnDistributedCacheClusterInstanceResponse, *http.Response, error) {
	return r.ApiService.ReplaceMsgVpnDistributedCacheClusterInstanceExecute(r)
}

/*
ReplaceMsgVpnDistributedCacheClusterInstance Replace a Cache Instance object.

Replace a Cache Instance object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

A Cache Instance is a single Cache process that belongs to a single Cache Cluster. A Cache Instance object provisioned on the broker is used to disseminate configuration information to the Cache process. Cache Instances listen for and cache live data messages that match the topic subscriptions configured for their parent Cache Cluster.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
cacheName|x|x||||
clusterName|x|x||||
instanceName|x|x||||
msgVpnName|x|x||||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @param instanceName The name of the Cache Instance.
 @return DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterInstanceRequest
*/
func (a *DistributedCacheApiService) ReplaceMsgVpnDistributedCacheClusterInstance(ctx context.Context, msgVpnName string, cacheName string, clusterName string, instanceName string) DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterInstanceRequest {
	return DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterInstanceRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
		instanceName: instanceName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterInstanceResponse
func (a *DistributedCacheApiService) ReplaceMsgVpnDistributedCacheClusterInstanceExecute(r DistributedCacheApiApiReplaceMsgVpnDistributedCacheClusterInstanceRequest) (*MsgVpnDistributedCacheClusterInstanceResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterInstanceResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.ReplaceMsgVpnDistributedCacheClusterInstance")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"instanceName"+"}", url.PathEscape(parameterToString(r.instanceName, "")), -1)

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

type DistributedCacheApiApiUpdateMsgVpnDistributedCacheRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	body *MsgVpnDistributedCache
	opaquePassword *string
	select_ *[]string
}

// The Distributed Cache object&#39;s attributes.
func (r DistributedCacheApiApiUpdateMsgVpnDistributedCacheRequest) Body(body MsgVpnDistributedCache) DistributedCacheApiApiUpdateMsgVpnDistributedCacheRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiUpdateMsgVpnDistributedCacheRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiUpdateMsgVpnDistributedCacheRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiUpdateMsgVpnDistributedCacheRequest) Select_(select_ []string) DistributedCacheApiApiUpdateMsgVpnDistributedCacheRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiUpdateMsgVpnDistributedCacheRequest) Execute() (*MsgVpnDistributedCacheResponse, *http.Response, error) {
	return r.ApiService.UpdateMsgVpnDistributedCacheExecute(r)
}

/*
UpdateMsgVpnDistributedCache Update a Distributed Cache object.

Update a Distributed Cache object. Any attribute missing from the request will be left unchanged.

A Distributed Cache is a collection of one or more Cache Clusters that belong to the same Message VPN. Each Cache Cluster in a Distributed Cache is configured to subscribe to a different set of topics. This effectively divides up the configured topic space, to provide scaling to very large topic spaces or very high cached message throughput.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
cacheName|x|x||||
msgVpnName|x|x||||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
MsgVpnDistributedCache|scheduledDeleteMsgDayList|scheduledDeleteMsgTimeList|
MsgVpnDistributedCache|scheduledDeleteMsgTimeList|scheduledDeleteMsgDayList|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @return DistributedCacheApiApiUpdateMsgVpnDistributedCacheRequest
*/
func (a *DistributedCacheApiService) UpdateMsgVpnDistributedCache(ctx context.Context, msgVpnName string, cacheName string) DistributedCacheApiApiUpdateMsgVpnDistributedCacheRequest {
	return DistributedCacheApiApiUpdateMsgVpnDistributedCacheRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheResponse
func (a *DistributedCacheApiService) UpdateMsgVpnDistributedCacheExecute(r DistributedCacheApiApiUpdateMsgVpnDistributedCacheRequest) (*MsgVpnDistributedCacheResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPatch
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.UpdateMsgVpnDistributedCache")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)

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

type DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	body *MsgVpnDistributedCacheCluster
	opaquePassword *string
	select_ *[]string
}

// The Cache Cluster object&#39;s attributes.
func (r DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterRequest) Body(body MsgVpnDistributedCacheCluster) DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterRequest) Select_(select_ []string) DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterRequest) Execute() (*MsgVpnDistributedCacheClusterResponse, *http.Response, error) {
	return r.ApiService.UpdateMsgVpnDistributedCacheClusterExecute(r)
}

/*
UpdateMsgVpnDistributedCacheCluster Update a Cache Cluster object.

Update a Cache Cluster object. Any attribute missing from the request will be left unchanged.

A Cache Cluster is a collection of one or more Cache Instances that subscribe to exactly the same topics. Cache Instances are grouped together in a Cache Cluster for the purpose of fault tolerance and load balancing. As published messages are received, the message broker message bus sends these live data messages to the Cache Instances in the Cache Cluster. This enables client cache requests to be served by any of Cache Instances in the Cache Cluster.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
cacheName|x|x||||
clusterName|x|x||||
msgVpnName|x|x||||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
EventThresholdByPercent|clearPercent|setPercent|
EventThresholdByPercent|setPercent|clearPercent|
EventThresholdByValue|clearValue|setValue|
EventThresholdByValue|setValue|clearValue|



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @return DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterRequest
*/
func (a *DistributedCacheApiService) UpdateMsgVpnDistributedCacheCluster(ctx context.Context, msgVpnName string, cacheName string, clusterName string) DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterRequest {
	return DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterResponse
func (a *DistributedCacheApiService) UpdateMsgVpnDistributedCacheClusterExecute(r DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterRequest) (*MsgVpnDistributedCacheClusterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPatch
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.UpdateMsgVpnDistributedCacheCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)

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

type DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterInstanceRequest struct {
	ctx context.Context
	ApiService *DistributedCacheApiService
	msgVpnName string
	cacheName string
	clusterName string
	instanceName string
	body *MsgVpnDistributedCacheClusterInstance
	opaquePassword *string
	select_ *[]string
}

// The Cache Instance object&#39;s attributes.
func (r DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterInstanceRequest) Body(body MsgVpnDistributedCacheClusterInstance) DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterInstanceRequest {
	r.body = &body
	return r
}
// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterInstanceRequest) OpaquePassword(opaquePassword string) DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterInstanceRequest {
	r.opaquePassword = &opaquePassword
	return r
}
// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterInstanceRequest) Select_(select_ []string) DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterInstanceRequest {
	r.select_ = &select_
	return r
}

func (r DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterInstanceRequest) Execute() (*MsgVpnDistributedCacheClusterInstanceResponse, *http.Response, error) {
	return r.ApiService.UpdateMsgVpnDistributedCacheClusterInstanceExecute(r)
}

/*
UpdateMsgVpnDistributedCacheClusterInstance Update a Cache Instance object.

Update a Cache Instance object. Any attribute missing from the request will be left unchanged.

A Cache Instance is a single Cache process that belongs to a single Cache Cluster. A Cache Instance object provisioned on the broker is used to disseminate configuration information to the Cache process. Cache Instances listen for and cache live data messages that match the topic subscriptions configured for their parent Cache Cluster.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
cacheName|x|x||||
clusterName|x|x||||
instanceName|x|x||||
msgVpnName|x|x||||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param cacheName The name of the Distributed Cache.
 @param clusterName The name of the Cache Cluster.
 @param instanceName The name of the Cache Instance.
 @return DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterInstanceRequest
*/
func (a *DistributedCacheApiService) UpdateMsgVpnDistributedCacheClusterInstance(ctx context.Context, msgVpnName string, cacheName string, clusterName string, instanceName string) DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterInstanceRequest {
	return DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterInstanceRequest{
		ApiService: a,
		ctx: ctx,
		msgVpnName: msgVpnName,
		cacheName: cacheName,
		clusterName: clusterName,
		instanceName: instanceName,
	}
}

// Execute executes the request
//  @return MsgVpnDistributedCacheClusterInstanceResponse
func (a *DistributedCacheApiService) UpdateMsgVpnDistributedCacheClusterInstanceExecute(r DistributedCacheApiApiUpdateMsgVpnDistributedCacheClusterInstanceRequest) (*MsgVpnDistributedCacheClusterInstanceResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPatch
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *MsgVpnDistributedCacheClusterInstanceResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DistributedCacheApiService.UpdateMsgVpnDistributedCacheClusterInstance")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/distributedCaches/{cacheName}/clusters/{clusterName}/instances/{instanceName}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"cacheName"+"}", url.PathEscape(parameterToString(r.cacheName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"clusterName"+"}", url.PathEscape(parameterToString(r.clusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"instanceName"+"}", url.PathEscape(parameterToString(r.instanceName, "")), -1)

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
