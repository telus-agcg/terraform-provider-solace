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

// DmrClusterApiService DmrClusterApi service
type DmrClusterApiService service

type DmrClusterApiApiCreateDmrClusterRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	body           *DmrCluster
	opaquePassword *string
	select_        *[]string
}

// The Cluster object&#39;s attributes.
func (r DmrClusterApiApiCreateDmrClusterRequest) Body(body DmrCluster) DmrClusterApiApiCreateDmrClusterRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiCreateDmrClusterRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiCreateDmrClusterRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiCreateDmrClusterRequest) Select_(select_ []string) DmrClusterApiApiCreateDmrClusterRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiCreateDmrClusterRequest) Execute() (*DmrClusterResponse, *http.Response, error) {
	return r.ApiService.CreateDmrClusterExecute(r)
}

/*
CreateDmrCluster Create a Cluster object.

Create a Cluster object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates via config-sync.

A Cluster is a provisioned object on a message broker that contains global DMR configuration parameters.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
authenticationBasicPassword||||x||x
authenticationClientCertContent||||x||x
authenticationClientCertPassword||||x||
dmrClusterName|x|x||||
nodeName|||x|||
tlsServerCertEnforceTrustedCommonNameEnabled|||||x|



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
DmrCluster|authenticationClientCertPassword|authenticationClientCertContent|



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return DmrClusterApiApiCreateDmrClusterRequest
*/
func (a *DmrClusterApiService) CreateDmrCluster(ctx context.Context) DmrClusterApiApiCreateDmrClusterRequest {
	return DmrClusterApiApiCreateDmrClusterRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//  @return DmrClusterResponse
func (a *DmrClusterApiService) CreateDmrClusterExecute(r DmrClusterApiApiCreateDmrClusterRequest) (*DmrClusterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.CreateDmrCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters"

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

type DmrClusterApiApiCreateDmrClusterLinkRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	body           *DmrClusterLink
	opaquePassword *string
	select_        *[]string
}

// The Link object&#39;s attributes.
func (r DmrClusterApiApiCreateDmrClusterLinkRequest) Body(body DmrClusterLink) DmrClusterApiApiCreateDmrClusterLinkRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiCreateDmrClusterLinkRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiCreateDmrClusterLinkRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiCreateDmrClusterLinkRequest) Select_(select_ []string) DmrClusterApiApiCreateDmrClusterLinkRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiCreateDmrClusterLinkRequest) Execute() (*DmrClusterLinkResponse, *http.Response, error) {
	return r.ApiService.CreateDmrClusterLinkExecute(r)
}

/*
CreateDmrClusterLink Create a Link object.

Create a Link object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates via config-sync.

A Link connects nodes (either within a Cluster or between two different Clusters) and allows them to exchange topology information, subscriptions and data.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
authenticationBasicPassword||||x||x
dmrClusterName|x||x|||
remoteNodeName|x|x||||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
EventThreshold|clearPercent|setPercent|clearValue, setValue
EventThreshold|clearValue|setValue|clearPercent, setPercent
EventThreshold|setPercent|clearPercent|clearValue, setValue
EventThreshold|setValue|clearValue|clearPercent, setPercent



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @return DmrClusterApiApiCreateDmrClusterLinkRequest
*/
func (a *DmrClusterApiService) CreateDmrClusterLink(ctx context.Context, dmrClusterName string) DmrClusterApiApiCreateDmrClusterLinkRequest {
	return DmrClusterApiApiCreateDmrClusterLinkRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
	}
}

// Execute executes the request
//  @return DmrClusterLinkResponse
func (a *DmrClusterApiService) CreateDmrClusterLinkExecute(r DmrClusterApiApiCreateDmrClusterLinkRequest) (*DmrClusterLinkResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterLinkResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.CreateDmrClusterLink")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}/links"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)

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

type DmrClusterApiApiCreateDmrClusterLinkRemoteAddressRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	remoteNodeName string
	body           *DmrClusterLinkRemoteAddress
	opaquePassword *string
	select_        *[]string
}

// The Remote Address object&#39;s attributes.
func (r DmrClusterApiApiCreateDmrClusterLinkRemoteAddressRequest) Body(body DmrClusterLinkRemoteAddress) DmrClusterApiApiCreateDmrClusterLinkRemoteAddressRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiCreateDmrClusterLinkRemoteAddressRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiCreateDmrClusterLinkRemoteAddressRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiCreateDmrClusterLinkRemoteAddressRequest) Select_(select_ []string) DmrClusterApiApiCreateDmrClusterLinkRemoteAddressRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiCreateDmrClusterLinkRemoteAddressRequest) Execute() (*DmrClusterLinkRemoteAddressResponse, *http.Response, error) {
	return r.ApiService.CreateDmrClusterLinkRemoteAddressExecute(r)
}

/*
CreateDmrClusterLinkRemoteAddress Create a Remote Address object.

Create a Remote Address object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates via config-sync.

Each Remote Address, consisting of a FQDN or IP address and optional port, is used to connect to the remote node for this Link. Up to 4 addresses may be provided for each Link, and will be tried on a round-robin basis.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
dmrClusterName|x||x|||
remoteAddress|x|x||||
remoteNodeName|x||x|||



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @param remoteNodeName The name of the node at the remote end of the Link.
 @return DmrClusterApiApiCreateDmrClusterLinkRemoteAddressRequest
*/
func (a *DmrClusterApiService) CreateDmrClusterLinkRemoteAddress(ctx context.Context, dmrClusterName string, remoteNodeName string) DmrClusterApiApiCreateDmrClusterLinkRemoteAddressRequest {
	return DmrClusterApiApiCreateDmrClusterLinkRemoteAddressRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
		remoteNodeName: remoteNodeName,
	}
}

// Execute executes the request
//  @return DmrClusterLinkRemoteAddressResponse
func (a *DmrClusterApiService) CreateDmrClusterLinkRemoteAddressExecute(r DmrClusterApiApiCreateDmrClusterLinkRemoteAddressRequest) (*DmrClusterLinkRemoteAddressResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterLinkRemoteAddressResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.CreateDmrClusterLinkRemoteAddress")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}/links/{remoteNodeName}/remoteAddresses"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteNodeName"+"}", url.PathEscape(parameterToString(r.remoteNodeName, "")), -1)

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

type DmrClusterApiApiCreateDmrClusterLinkTlsTrustedCommonNameRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	remoteNodeName string
	body           *DmrClusterLinkTlsTrustedCommonName
	opaquePassword *string
	select_        *[]string
}

// The Trusted Common Name object&#39;s attributes.
func (r DmrClusterApiApiCreateDmrClusterLinkTlsTrustedCommonNameRequest) Body(body DmrClusterLinkTlsTrustedCommonName) DmrClusterApiApiCreateDmrClusterLinkTlsTrustedCommonNameRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiCreateDmrClusterLinkTlsTrustedCommonNameRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiCreateDmrClusterLinkTlsTrustedCommonNameRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiCreateDmrClusterLinkTlsTrustedCommonNameRequest) Select_(select_ []string) DmrClusterApiApiCreateDmrClusterLinkTlsTrustedCommonNameRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiCreateDmrClusterLinkTlsTrustedCommonNameRequest) Execute() (*DmrClusterLinkTlsTrustedCommonNameResponse, *http.Response, error) {
	return r.ApiService.CreateDmrClusterLinkTlsTrustedCommonNameExecute(r)
}

/*
CreateDmrClusterLinkTlsTrustedCommonName Create a Trusted Common Name object.

Create a Trusted Common Name object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates via config-sync.

The Trusted Common Names for the Link are used by encrypted transports to verify the name in the certificate presented by the remote node. They must include the common name of the remote node's server certificate or client certificate, depending upon the initiator of the connection.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
dmrClusterName|x||x||x|
remoteNodeName|x||x||x|
tlsTrustedCommonName|x|x|||x|



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been deprecated since 2.18. Common Name validation has been replaced by Server Certificate Name validation.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @param remoteNodeName The name of the node at the remote end of the Link.
 @return DmrClusterApiApiCreateDmrClusterLinkTlsTrustedCommonNameRequest

Deprecated
*/
func (a *DmrClusterApiService) CreateDmrClusterLinkTlsTrustedCommonName(ctx context.Context, dmrClusterName string, remoteNodeName string) DmrClusterApiApiCreateDmrClusterLinkTlsTrustedCommonNameRequest {
	return DmrClusterApiApiCreateDmrClusterLinkTlsTrustedCommonNameRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
		remoteNodeName: remoteNodeName,
	}
}

// Execute executes the request
//  @return DmrClusterLinkTlsTrustedCommonNameResponse
// Deprecated
func (a *DmrClusterApiService) CreateDmrClusterLinkTlsTrustedCommonNameExecute(r DmrClusterApiApiCreateDmrClusterLinkTlsTrustedCommonNameRequest) (*DmrClusterLinkTlsTrustedCommonNameResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterLinkTlsTrustedCommonNameResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.CreateDmrClusterLinkTlsTrustedCommonName")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}/links/{remoteNodeName}/tlsTrustedCommonNames"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteNodeName"+"}", url.PathEscape(parameterToString(r.remoteNodeName, "")), -1)

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

type DmrClusterApiApiDeleteDmrClusterRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
}

func (r DmrClusterApiApiDeleteDmrClusterRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteDmrClusterExecute(r)
}

/*
DeleteDmrCluster Delete a Cluster object.

Delete a Cluster object. The deletion of instances of this object are synchronized to HA mates via config-sync.

A Cluster is a provisioned object on a message broker that contains global DMR configuration parameters.

A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @return DmrClusterApiApiDeleteDmrClusterRequest
*/
func (a *DmrClusterApiService) DeleteDmrCluster(ctx context.Context, dmrClusterName string) DmrClusterApiApiDeleteDmrClusterRequest {
	return DmrClusterApiApiDeleteDmrClusterRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *DmrClusterApiService) DeleteDmrClusterExecute(r DmrClusterApiApiDeleteDmrClusterRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.DeleteDmrCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)

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

type DmrClusterApiApiDeleteDmrClusterLinkRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	remoteNodeName string
}

func (r DmrClusterApiApiDeleteDmrClusterLinkRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteDmrClusterLinkExecute(r)
}

/*
DeleteDmrClusterLink Delete a Link object.

Delete a Link object. The deletion of instances of this object are synchronized to HA mates via config-sync.

A Link connects nodes (either within a Cluster or between two different Clusters) and allows them to exchange topology information, subscriptions and data.

A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @param remoteNodeName The name of the node at the remote end of the Link.
 @return DmrClusterApiApiDeleteDmrClusterLinkRequest
*/
func (a *DmrClusterApiService) DeleteDmrClusterLink(ctx context.Context, dmrClusterName string, remoteNodeName string) DmrClusterApiApiDeleteDmrClusterLinkRequest {
	return DmrClusterApiApiDeleteDmrClusterLinkRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
		remoteNodeName: remoteNodeName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *DmrClusterApiService) DeleteDmrClusterLinkExecute(r DmrClusterApiApiDeleteDmrClusterLinkRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.DeleteDmrClusterLink")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}/links/{remoteNodeName}"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteNodeName"+"}", url.PathEscape(parameterToString(r.remoteNodeName, "")), -1)

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

type DmrClusterApiApiDeleteDmrClusterLinkRemoteAddressRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	remoteNodeName string
	remoteAddress  string
}

func (r DmrClusterApiApiDeleteDmrClusterLinkRemoteAddressRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteDmrClusterLinkRemoteAddressExecute(r)
}

/*
DeleteDmrClusterLinkRemoteAddress Delete a Remote Address object.

Delete a Remote Address object. The deletion of instances of this object are synchronized to HA mates via config-sync.

Each Remote Address, consisting of a FQDN or IP address and optional port, is used to connect to the remote node for this Link. Up to 4 addresses may be provided for each Link, and will be tried on a round-robin basis.

A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @param remoteNodeName The name of the node at the remote end of the Link.
 @param remoteAddress The FQDN or IP address (and optional port) of the remote node. If a port is not provided, it will vary based on the transport encoding: 55555 (plain-text), 55443 (encrypted), or 55003 (compressed).
 @return DmrClusterApiApiDeleteDmrClusterLinkRemoteAddressRequest
*/
func (a *DmrClusterApiService) DeleteDmrClusterLinkRemoteAddress(ctx context.Context, dmrClusterName string, remoteNodeName string, remoteAddress string) DmrClusterApiApiDeleteDmrClusterLinkRemoteAddressRequest {
	return DmrClusterApiApiDeleteDmrClusterLinkRemoteAddressRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
		remoteNodeName: remoteNodeName,
		remoteAddress:  remoteAddress,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *DmrClusterApiService) DeleteDmrClusterLinkRemoteAddressExecute(r DmrClusterApiApiDeleteDmrClusterLinkRemoteAddressRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.DeleteDmrClusterLinkRemoteAddress")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}/links/{remoteNodeName}/remoteAddresses/{remoteAddress}"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteNodeName"+"}", url.PathEscape(parameterToString(r.remoteNodeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteAddress"+"}", url.PathEscape(parameterToString(r.remoteAddress, "")), -1)

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

type DmrClusterApiApiDeleteDmrClusterLinkTlsTrustedCommonNameRequest struct {
	ctx                  context.Context
	ApiService           *DmrClusterApiService
	dmrClusterName       string
	remoteNodeName       string
	tlsTrustedCommonName string
}

func (r DmrClusterApiApiDeleteDmrClusterLinkTlsTrustedCommonNameRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteDmrClusterLinkTlsTrustedCommonNameExecute(r)
}

/*
DeleteDmrClusterLinkTlsTrustedCommonName Delete a Trusted Common Name object.

Delete a Trusted Common Name object. The deletion of instances of this object are synchronized to HA mates via config-sync.

The Trusted Common Names for the Link are used by encrypted transports to verify the name in the certificate presented by the remote node. They must include the common name of the remote node's server certificate or client certificate, depending upon the initiator of the connection.

A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been deprecated since 2.18. Common Name validation has been replaced by Server Certificate Name validation.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @param remoteNodeName The name of the node at the remote end of the Link.
 @param tlsTrustedCommonName The expected trusted common name of the remote certificate.
 @return DmrClusterApiApiDeleteDmrClusterLinkTlsTrustedCommonNameRequest

Deprecated
*/
func (a *DmrClusterApiService) DeleteDmrClusterLinkTlsTrustedCommonName(ctx context.Context, dmrClusterName string, remoteNodeName string, tlsTrustedCommonName string) DmrClusterApiApiDeleteDmrClusterLinkTlsTrustedCommonNameRequest {
	return DmrClusterApiApiDeleteDmrClusterLinkTlsTrustedCommonNameRequest{
		ApiService:           a,
		ctx:                  ctx,
		dmrClusterName:       dmrClusterName,
		remoteNodeName:       remoteNodeName,
		tlsTrustedCommonName: tlsTrustedCommonName,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
// Deprecated
func (a *DmrClusterApiService) DeleteDmrClusterLinkTlsTrustedCommonNameExecute(r DmrClusterApiApiDeleteDmrClusterLinkTlsTrustedCommonNameRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.DeleteDmrClusterLinkTlsTrustedCommonName")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}/links/{remoteNodeName}/tlsTrustedCommonNames/{tlsTrustedCommonName}"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteNodeName"+"}", url.PathEscape(parameterToString(r.remoteNodeName, "")), -1)
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

type DmrClusterApiApiGetDmrClusterRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	opaquePassword *string
	select_        *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiGetDmrClusterRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterRequest) Select_(select_ []string) DmrClusterApiApiGetDmrClusterRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiGetDmrClusterRequest) Execute() (*DmrClusterResponse, *http.Response, error) {
	return r.ApiService.GetDmrClusterExecute(r)
}

/*
GetDmrCluster Get a Cluster object.

Get a Cluster object.

A Cluster is a provisioned object on a message broker that contains global DMR configuration parameters.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
authenticationBasicPassword||x||x
authenticationClientCertContent||x||x
authenticationClientCertPassword||x||
dmrClusterName|x|||
tlsServerCertEnforceTrustedCommonNameEnabled|||x|



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @return DmrClusterApiApiGetDmrClusterRequest
*/
func (a *DmrClusterApiService) GetDmrCluster(ctx context.Context, dmrClusterName string) DmrClusterApiApiGetDmrClusterRequest {
	return DmrClusterApiApiGetDmrClusterRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
	}
}

// Execute executes the request
//  @return DmrClusterResponse
func (a *DmrClusterApiService) GetDmrClusterExecute(r DmrClusterApiApiGetDmrClusterRequest) (*DmrClusterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.GetDmrCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)

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

type DmrClusterApiApiGetDmrClusterLinkRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	remoteNodeName string
	opaquePassword *string
	select_        *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinkRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiGetDmrClusterLinkRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinkRequest) Select_(select_ []string) DmrClusterApiApiGetDmrClusterLinkRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiGetDmrClusterLinkRequest) Execute() (*DmrClusterLinkResponse, *http.Response, error) {
	return r.ApiService.GetDmrClusterLinkExecute(r)
}

/*
GetDmrClusterLink Get a Link object.

Get a Link object.

A Link connects nodes (either within a Cluster or between two different Clusters) and allows them to exchange topology information, subscriptions and data.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
authenticationBasicPassword||x||x
dmrClusterName|x|||
remoteNodeName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @param remoteNodeName The name of the node at the remote end of the Link.
 @return DmrClusterApiApiGetDmrClusterLinkRequest
*/
func (a *DmrClusterApiService) GetDmrClusterLink(ctx context.Context, dmrClusterName string, remoteNodeName string) DmrClusterApiApiGetDmrClusterLinkRequest {
	return DmrClusterApiApiGetDmrClusterLinkRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
		remoteNodeName: remoteNodeName,
	}
}

// Execute executes the request
//  @return DmrClusterLinkResponse
func (a *DmrClusterApiService) GetDmrClusterLinkExecute(r DmrClusterApiApiGetDmrClusterLinkRequest) (*DmrClusterLinkResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterLinkResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.GetDmrClusterLink")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}/links/{remoteNodeName}"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteNodeName"+"}", url.PathEscape(parameterToString(r.remoteNodeName, "")), -1)

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

type DmrClusterApiApiGetDmrClusterLinkRemoteAddressRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	remoteNodeName string
	remoteAddress  string
	opaquePassword *string
	select_        *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinkRemoteAddressRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiGetDmrClusterLinkRemoteAddressRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinkRemoteAddressRequest) Select_(select_ []string) DmrClusterApiApiGetDmrClusterLinkRemoteAddressRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiGetDmrClusterLinkRemoteAddressRequest) Execute() (*DmrClusterLinkRemoteAddressResponse, *http.Response, error) {
	return r.ApiService.GetDmrClusterLinkRemoteAddressExecute(r)
}

/*
GetDmrClusterLinkRemoteAddress Get a Remote Address object.

Get a Remote Address object.

Each Remote Address, consisting of a FQDN or IP address and optional port, is used to connect to the remote node for this Link. Up to 4 addresses may be provided for each Link, and will be tried on a round-robin basis.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
dmrClusterName|x|||
remoteAddress|x|||
remoteNodeName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @param remoteNodeName The name of the node at the remote end of the Link.
 @param remoteAddress The FQDN or IP address (and optional port) of the remote node. If a port is not provided, it will vary based on the transport encoding: 55555 (plain-text), 55443 (encrypted), or 55003 (compressed).
 @return DmrClusterApiApiGetDmrClusterLinkRemoteAddressRequest
*/
func (a *DmrClusterApiService) GetDmrClusterLinkRemoteAddress(ctx context.Context, dmrClusterName string, remoteNodeName string, remoteAddress string) DmrClusterApiApiGetDmrClusterLinkRemoteAddressRequest {
	return DmrClusterApiApiGetDmrClusterLinkRemoteAddressRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
		remoteNodeName: remoteNodeName,
		remoteAddress:  remoteAddress,
	}
}

// Execute executes the request
//  @return DmrClusterLinkRemoteAddressResponse
func (a *DmrClusterApiService) GetDmrClusterLinkRemoteAddressExecute(r DmrClusterApiApiGetDmrClusterLinkRemoteAddressRequest) (*DmrClusterLinkRemoteAddressResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterLinkRemoteAddressResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.GetDmrClusterLinkRemoteAddress")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}/links/{remoteNodeName}/remoteAddresses/{remoteAddress}"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteNodeName"+"}", url.PathEscape(parameterToString(r.remoteNodeName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteAddress"+"}", url.PathEscape(parameterToString(r.remoteAddress, "")), -1)

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

type DmrClusterApiApiGetDmrClusterLinkRemoteAddressesRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	remoteNodeName string
	opaquePassword *string
	where          *[]string
	select_        *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinkRemoteAddressesRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiGetDmrClusterLinkRemoteAddressesRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinkRemoteAddressesRequest) Where(where []string) DmrClusterApiApiGetDmrClusterLinkRemoteAddressesRequest {
	r.where = &where
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinkRemoteAddressesRequest) Select_(select_ []string) DmrClusterApiApiGetDmrClusterLinkRemoteAddressesRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiGetDmrClusterLinkRemoteAddressesRequest) Execute() (*DmrClusterLinkRemoteAddressesResponse, *http.Response, error) {
	return r.ApiService.GetDmrClusterLinkRemoteAddressesExecute(r)
}

/*
GetDmrClusterLinkRemoteAddresses Get a list of Remote Address objects.

Get a list of Remote Address objects.

Each Remote Address, consisting of a FQDN or IP address and optional port, is used to connect to the remote node for this Link. Up to 4 addresses may be provided for each Link, and will be tried on a round-robin basis.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
dmrClusterName|x|||
remoteAddress|x|||
remoteNodeName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @param remoteNodeName The name of the node at the remote end of the Link.
 @return DmrClusterApiApiGetDmrClusterLinkRemoteAddressesRequest
*/
func (a *DmrClusterApiService) GetDmrClusterLinkRemoteAddresses(ctx context.Context, dmrClusterName string, remoteNodeName string) DmrClusterApiApiGetDmrClusterLinkRemoteAddressesRequest {
	return DmrClusterApiApiGetDmrClusterLinkRemoteAddressesRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
		remoteNodeName: remoteNodeName,
	}
}

// Execute executes the request
//  @return DmrClusterLinkRemoteAddressesResponse
func (a *DmrClusterApiService) GetDmrClusterLinkRemoteAddressesExecute(r DmrClusterApiApiGetDmrClusterLinkRemoteAddressesRequest) (*DmrClusterLinkRemoteAddressesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterLinkRemoteAddressesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.GetDmrClusterLinkRemoteAddresses")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}/links/{remoteNodeName}/remoteAddresses"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteNodeName"+"}", url.PathEscape(parameterToString(r.remoteNodeName, "")), -1)

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

type DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNameRequest struct {
	ctx                  context.Context
	ApiService           *DmrClusterApiService
	dmrClusterName       string
	remoteNodeName       string
	tlsTrustedCommonName string
	opaquePassword       *string
	select_              *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNameRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNameRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNameRequest) Select_(select_ []string) DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNameRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNameRequest) Execute() (*DmrClusterLinkTlsTrustedCommonNameResponse, *http.Response, error) {
	return r.ApiService.GetDmrClusterLinkTlsTrustedCommonNameExecute(r)
}

/*
GetDmrClusterLinkTlsTrustedCommonName Get a Trusted Common Name object.

Get a Trusted Common Name object.

The Trusted Common Names for the Link are used by encrypted transports to verify the name in the certificate presented by the remote node. They must include the common name of the remote node's server certificate or client certificate, depending upon the initiator of the connection.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
dmrClusterName|x||x|
remoteNodeName|x||x|
tlsTrustedCommonName|x||x|



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been deprecated since 2.18. Common Name validation has been replaced by Server Certificate Name validation.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @param remoteNodeName The name of the node at the remote end of the Link.
 @param tlsTrustedCommonName The expected trusted common name of the remote certificate.
 @return DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNameRequest

Deprecated
*/
func (a *DmrClusterApiService) GetDmrClusterLinkTlsTrustedCommonName(ctx context.Context, dmrClusterName string, remoteNodeName string, tlsTrustedCommonName string) DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNameRequest {
	return DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNameRequest{
		ApiService:           a,
		ctx:                  ctx,
		dmrClusterName:       dmrClusterName,
		remoteNodeName:       remoteNodeName,
		tlsTrustedCommonName: tlsTrustedCommonName,
	}
}

// Execute executes the request
//  @return DmrClusterLinkTlsTrustedCommonNameResponse
// Deprecated
func (a *DmrClusterApiService) GetDmrClusterLinkTlsTrustedCommonNameExecute(r DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNameRequest) (*DmrClusterLinkTlsTrustedCommonNameResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterLinkTlsTrustedCommonNameResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.GetDmrClusterLinkTlsTrustedCommonName")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}/links/{remoteNodeName}/tlsTrustedCommonNames/{tlsTrustedCommonName}"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteNodeName"+"}", url.PathEscape(parameterToString(r.remoteNodeName, "")), -1)
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

type DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNamesRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	remoteNodeName string
	opaquePassword *string
	where          *[]string
	select_        *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNamesRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNamesRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNamesRequest) Where(where []string) DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNamesRequest {
	r.where = &where
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNamesRequest) Select_(select_ []string) DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNamesRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNamesRequest) Execute() (*DmrClusterLinkTlsTrustedCommonNamesResponse, *http.Response, error) {
	return r.ApiService.GetDmrClusterLinkTlsTrustedCommonNamesExecute(r)
}

/*
GetDmrClusterLinkTlsTrustedCommonNames Get a list of Trusted Common Name objects.

Get a list of Trusted Common Name objects.

The Trusted Common Names for the Link are used by encrypted transports to verify the name in the certificate presented by the remote node. They must include the common name of the remote node's server certificate or client certificate, depending upon the initiator of the connection.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
dmrClusterName|x||x|
remoteNodeName|x||x|
tlsTrustedCommonName|x||x|



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been deprecated since 2.18. Common Name validation has been replaced by Server Certificate Name validation.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @param remoteNodeName The name of the node at the remote end of the Link.
 @return DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNamesRequest

Deprecated
*/
func (a *DmrClusterApiService) GetDmrClusterLinkTlsTrustedCommonNames(ctx context.Context, dmrClusterName string, remoteNodeName string) DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNamesRequest {
	return DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNamesRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
		remoteNodeName: remoteNodeName,
	}
}

// Execute executes the request
//  @return DmrClusterLinkTlsTrustedCommonNamesResponse
// Deprecated
func (a *DmrClusterApiService) GetDmrClusterLinkTlsTrustedCommonNamesExecute(r DmrClusterApiApiGetDmrClusterLinkTlsTrustedCommonNamesRequest) (*DmrClusterLinkTlsTrustedCommonNamesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterLinkTlsTrustedCommonNamesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.GetDmrClusterLinkTlsTrustedCommonNames")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}/links/{remoteNodeName}/tlsTrustedCommonNames"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteNodeName"+"}", url.PathEscape(parameterToString(r.remoteNodeName, "")), -1)

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

type DmrClusterApiApiGetDmrClusterLinksRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	count          *int32
	cursor         *string
	opaquePassword *string
	where          *[]string
	select_        *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinksRequest) Count(count int32) DmrClusterApiApiGetDmrClusterLinksRequest {
	r.count = &count
	return r
}

// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinksRequest) Cursor(cursor string) DmrClusterApiApiGetDmrClusterLinksRequest {
	r.cursor = &cursor
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinksRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiGetDmrClusterLinksRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinksRequest) Where(where []string) DmrClusterApiApiGetDmrClusterLinksRequest {
	r.where = &where
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiGetDmrClusterLinksRequest) Select_(select_ []string) DmrClusterApiApiGetDmrClusterLinksRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiGetDmrClusterLinksRequest) Execute() (*DmrClusterLinksResponse, *http.Response, error) {
	return r.ApiService.GetDmrClusterLinksExecute(r)
}

/*
GetDmrClusterLinks Get a list of Link objects.

Get a list of Link objects.

A Link connects nodes (either within a Cluster or between two different Clusters) and allows them to exchange topology information, subscriptions and data.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
authenticationBasicPassword||x||x
dmrClusterName|x|||
remoteNodeName|x|||



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @return DmrClusterApiApiGetDmrClusterLinksRequest
*/
func (a *DmrClusterApiService) GetDmrClusterLinks(ctx context.Context, dmrClusterName string) DmrClusterApiApiGetDmrClusterLinksRequest {
	return DmrClusterApiApiGetDmrClusterLinksRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
	}
}

// Execute executes the request
//  @return DmrClusterLinksResponse
func (a *DmrClusterApiService) GetDmrClusterLinksExecute(r DmrClusterApiApiGetDmrClusterLinksRequest) (*DmrClusterLinksResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterLinksResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.GetDmrClusterLinks")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}/links"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)

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

type DmrClusterApiApiGetDmrClustersRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	count          *int32
	cursor         *string
	opaquePassword *string
	where          *[]string
	select_        *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r DmrClusterApiApiGetDmrClustersRequest) Count(count int32) DmrClusterApiApiGetDmrClustersRequest {
	r.count = &count
	return r
}

// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r DmrClusterApiApiGetDmrClustersRequest) Cursor(cursor string) DmrClusterApiApiGetDmrClustersRequest {
	r.cursor = &cursor
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiGetDmrClustersRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiGetDmrClustersRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r DmrClusterApiApiGetDmrClustersRequest) Where(where []string) DmrClusterApiApiGetDmrClustersRequest {
	r.where = &where
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiGetDmrClustersRequest) Select_(select_ []string) DmrClusterApiApiGetDmrClustersRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiGetDmrClustersRequest) Execute() (*DmrClustersResponse, *http.Response, error) {
	return r.ApiService.GetDmrClustersExecute(r)
}

/*
GetDmrClusters Get a list of Cluster objects.

Get a list of Cluster objects.

A Cluster is a provisioned object on a message broker that contains global DMR configuration parameters.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
authenticationBasicPassword||x||x
authenticationClientCertContent||x||x
authenticationClientCertPassword||x||
dmrClusterName|x|||
tlsServerCertEnforceTrustedCommonNameEnabled|||x|



A SEMP client authorized with a minimum access scope/level of "global/read-only" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return DmrClusterApiApiGetDmrClustersRequest
*/
func (a *DmrClusterApiService) GetDmrClusters(ctx context.Context) DmrClusterApiApiGetDmrClustersRequest {
	return DmrClusterApiApiGetDmrClustersRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//  @return DmrClustersResponse
func (a *DmrClusterApiService) GetDmrClustersExecute(r DmrClusterApiApiGetDmrClustersRequest) (*DmrClustersResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClustersResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.GetDmrClusters")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters"

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

type DmrClusterApiApiReplaceDmrClusterRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	body           *DmrCluster
	opaquePassword *string
	select_        *[]string
}

// The Cluster object&#39;s attributes.
func (r DmrClusterApiApiReplaceDmrClusterRequest) Body(body DmrCluster) DmrClusterApiApiReplaceDmrClusterRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiReplaceDmrClusterRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiReplaceDmrClusterRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiReplaceDmrClusterRequest) Select_(select_ []string) DmrClusterApiApiReplaceDmrClusterRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiReplaceDmrClusterRequest) Execute() (*DmrClusterResponse, *http.Response, error) {
	return r.ApiService.ReplaceDmrClusterExecute(r)
}

/*
ReplaceDmrCluster Replace a Cluster object.

Replace a Cluster object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

A Cluster is a provisioned object on a message broker that contains global DMR configuration parameters.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
authenticationBasicPassword|||x|x||x
authenticationClientCertContent|||x|x||x
authenticationClientCertPassword|||x|x||
directOnlyEnabled||x||||
dmrClusterName|x|x||||
nodeName||x||||
tlsServerCertEnforceTrustedCommonNameEnabled|||||x|



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
DmrCluster|authenticationClientCertPassword|authenticationClientCertContent|



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @return DmrClusterApiApiReplaceDmrClusterRequest
*/
func (a *DmrClusterApiService) ReplaceDmrCluster(ctx context.Context, dmrClusterName string) DmrClusterApiApiReplaceDmrClusterRequest {
	return DmrClusterApiApiReplaceDmrClusterRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
	}
}

// Execute executes the request
//  @return DmrClusterResponse
func (a *DmrClusterApiService) ReplaceDmrClusterExecute(r DmrClusterApiApiReplaceDmrClusterRequest) (*DmrClusterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.ReplaceDmrCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)

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

type DmrClusterApiApiReplaceDmrClusterLinkRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	remoteNodeName string
	body           *DmrClusterLink
	opaquePassword *string
	select_        *[]string
}

// The Link object&#39;s attributes.
func (r DmrClusterApiApiReplaceDmrClusterLinkRequest) Body(body DmrClusterLink) DmrClusterApiApiReplaceDmrClusterLinkRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiReplaceDmrClusterLinkRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiReplaceDmrClusterLinkRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiReplaceDmrClusterLinkRequest) Select_(select_ []string) DmrClusterApiApiReplaceDmrClusterLinkRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiReplaceDmrClusterLinkRequest) Execute() (*DmrClusterLinkResponse, *http.Response, error) {
	return r.ApiService.ReplaceDmrClusterLinkExecute(r)
}

/*
ReplaceDmrClusterLink Replace a Link object.

Replace a Link object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

A Link connects nodes (either within a Cluster or between two different Clusters) and allows them to exchange topology information, subscriptions and data.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
authenticationBasicPassword|||x|x||x
authenticationScheme||||x||
dmrClusterName|x|x||||
egressFlowWindowSize||||x||
initiator||||x||
remoteNodeName|x|x||||
span||||x||
transportCompressedEnabled||||x||
transportTlsEnabled||||x||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
EventThreshold|clearPercent|setPercent|clearValue, setValue
EventThreshold|clearValue|setValue|clearPercent, setPercent
EventThreshold|setPercent|clearPercent|clearValue, setValue
EventThreshold|setValue|clearValue|clearPercent, setPercent



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @param remoteNodeName The name of the node at the remote end of the Link.
 @return DmrClusterApiApiReplaceDmrClusterLinkRequest
*/
func (a *DmrClusterApiService) ReplaceDmrClusterLink(ctx context.Context, dmrClusterName string, remoteNodeName string) DmrClusterApiApiReplaceDmrClusterLinkRequest {
	return DmrClusterApiApiReplaceDmrClusterLinkRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
		remoteNodeName: remoteNodeName,
	}
}

// Execute executes the request
//  @return DmrClusterLinkResponse
func (a *DmrClusterApiService) ReplaceDmrClusterLinkExecute(r DmrClusterApiApiReplaceDmrClusterLinkRequest) (*DmrClusterLinkResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterLinkResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.ReplaceDmrClusterLink")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}/links/{remoteNodeName}"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteNodeName"+"}", url.PathEscape(parameterToString(r.remoteNodeName, "")), -1)

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

type DmrClusterApiApiUpdateDmrClusterRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	body           *DmrCluster
	opaquePassword *string
	select_        *[]string
}

// The Cluster object&#39;s attributes.
func (r DmrClusterApiApiUpdateDmrClusterRequest) Body(body DmrCluster) DmrClusterApiApiUpdateDmrClusterRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiUpdateDmrClusterRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiUpdateDmrClusterRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiUpdateDmrClusterRequest) Select_(select_ []string) DmrClusterApiApiUpdateDmrClusterRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiUpdateDmrClusterRequest) Execute() (*DmrClusterResponse, *http.Response, error) {
	return r.ApiService.UpdateDmrClusterExecute(r)
}

/*
UpdateDmrCluster Update a Cluster object.

Update a Cluster object. Any attribute missing from the request will be left unchanged.

A Cluster is a provisioned object on a message broker that contains global DMR configuration parameters.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
authenticationBasicPassword|||x|x||x
authenticationClientCertContent|||x|x||x
authenticationClientCertPassword|||x|x||
directOnlyEnabled||x||||
dmrClusterName|x|x||||
nodeName||x||||
tlsServerCertEnforceTrustedCommonNameEnabled|||||x|



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
DmrCluster|authenticationClientCertPassword|authenticationClientCertContent|



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @return DmrClusterApiApiUpdateDmrClusterRequest
*/
func (a *DmrClusterApiService) UpdateDmrCluster(ctx context.Context, dmrClusterName string) DmrClusterApiApiUpdateDmrClusterRequest {
	return DmrClusterApiApiUpdateDmrClusterRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
	}
}

// Execute executes the request
//  @return DmrClusterResponse
func (a *DmrClusterApiService) UpdateDmrClusterExecute(r DmrClusterApiApiUpdateDmrClusterRequest) (*DmrClusterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.UpdateDmrCluster")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)

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

type DmrClusterApiApiUpdateDmrClusterLinkRequest struct {
	ctx            context.Context
	ApiService     *DmrClusterApiService
	dmrClusterName string
	remoteNodeName string
	body           *DmrClusterLink
	opaquePassword *string
	select_        *[]string
}

// The Link object&#39;s attributes.
func (r DmrClusterApiApiUpdateDmrClusterLinkRequest) Body(body DmrClusterLink) DmrClusterApiApiUpdateDmrClusterLinkRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r DmrClusterApiApiUpdateDmrClusterLinkRequest) OpaquePassword(opaquePassword string) DmrClusterApiApiUpdateDmrClusterLinkRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r DmrClusterApiApiUpdateDmrClusterLinkRequest) Select_(select_ []string) DmrClusterApiApiUpdateDmrClusterLinkRequest {
	r.select_ = &select_
	return r
}

func (r DmrClusterApiApiUpdateDmrClusterLinkRequest) Execute() (*DmrClusterLinkResponse, *http.Response, error) {
	return r.ApiService.UpdateDmrClusterLinkExecute(r)
}

/*
UpdateDmrClusterLink Update a Link object.

Update a Link object. Any attribute missing from the request will be left unchanged.

A Link connects nodes (either within a Cluster or between two different Clusters) and allows them to exchange topology information, subscriptions and data.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
authenticationBasicPassword|||x|x||x
authenticationScheme||||x||
dmrClusterName|x|x||||
egressFlowWindowSize||||x||
initiator||||x||
remoteNodeName|x|x||||
span||||x||
transportCompressedEnabled||||x||
transportTlsEnabled||||x||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
EventThreshold|clearPercent|setPercent|clearValue, setValue
EventThreshold|clearValue|setValue|clearPercent, setPercent
EventThreshold|setPercent|clearPercent|clearValue, setValue
EventThreshold|setValue|clearValue|clearPercent, setPercent



A SEMP client authorized with a minimum access scope/level of "global/read-write" is required to perform this operation.

This has been available since 2.11.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param dmrClusterName The name of the Cluster.
 @param remoteNodeName The name of the node at the remote end of the Link.
 @return DmrClusterApiApiUpdateDmrClusterLinkRequest
*/
func (a *DmrClusterApiService) UpdateDmrClusterLink(ctx context.Context, dmrClusterName string, remoteNodeName string) DmrClusterApiApiUpdateDmrClusterLinkRequest {
	return DmrClusterApiApiUpdateDmrClusterLinkRequest{
		ApiService:     a,
		ctx:            ctx,
		dmrClusterName: dmrClusterName,
		remoteNodeName: remoteNodeName,
	}
}

// Execute executes the request
//  @return DmrClusterLinkResponse
func (a *DmrClusterApiService) UpdateDmrClusterLinkExecute(r DmrClusterApiApiUpdateDmrClusterLinkRequest) (*DmrClusterLinkResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DmrClusterLinkResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DmrClusterApiService.UpdateDmrClusterLink")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/dmrClusters/{dmrClusterName}/links/{remoteNodeName}"
	localVarPath = strings.Replace(localVarPath, "{"+"dmrClusterName"+"}", url.PathEscape(parameterToString(r.dmrClusterName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"remoteNodeName"+"}", url.PathEscape(parameterToString(r.remoteNodeName, "")), -1)

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
