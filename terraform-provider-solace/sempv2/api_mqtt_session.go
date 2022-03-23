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

// MqttSessionApiService MqttSessionApi service
type MqttSessionApiService service

type MqttSessionApiApiCreateMsgVpnMqttSessionRequest struct {
	ctx            context.Context
	ApiService     *MqttSessionApiService
	msgVpnName     string
	body           *MsgVpnMqttSession
	opaquePassword *string
	select_        *[]string
}

// The MQTT Session object&#39;s attributes.
func (r MqttSessionApiApiCreateMsgVpnMqttSessionRequest) Body(body MsgVpnMqttSession) MqttSessionApiApiCreateMsgVpnMqttSessionRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r MqttSessionApiApiCreateMsgVpnMqttSessionRequest) OpaquePassword(opaquePassword string) MqttSessionApiApiCreateMsgVpnMqttSessionRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r MqttSessionApiApiCreateMsgVpnMqttSessionRequest) Select_(select_ []string) MqttSessionApiApiCreateMsgVpnMqttSessionRequest {
	r.select_ = &select_
	return r
}

func (r MqttSessionApiApiCreateMsgVpnMqttSessionRequest) Execute() (*MsgVpnMqttSessionResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnMqttSessionExecute(r)
}

/*
CreateMsgVpnMqttSession Create an MQTT Session object.

Create an MQTT Session object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

An MQTT Session object is a virtual representation of an MQTT client connection. An MQTT session holds the state of an MQTT client (that is, it is used to contain a client's QoS 0 and QoS 1 subscription sets and any undelivered QoS 1 messages).


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
mqttSessionClientId|x|x||||
mqttSessionVirtualRouter|x|x||||
msgVpnName|x||x|||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
EventThreshold|clearPercent|setPercent|clearValue, setValue
EventThreshold|clearValue|setValue|clearPercent, setPercent
EventThreshold|setPercent|clearPercent|clearValue, setValue
EventThreshold|setValue|clearValue|clearPercent, setPercent



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.1.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @return MqttSessionApiApiCreateMsgVpnMqttSessionRequest
*/
func (a *MqttSessionApiService) CreateMsgVpnMqttSession(ctx context.Context, msgVpnName string) MqttSessionApiApiCreateMsgVpnMqttSessionRequest {
	return MqttSessionApiApiCreateMsgVpnMqttSessionRequest{
		ApiService: a,
		ctx:        ctx,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return MsgVpnMqttSessionResponse
func (a *MqttSessionApiService) CreateMsgVpnMqttSessionExecute(r MqttSessionApiApiCreateMsgVpnMqttSessionRequest) (*MsgVpnMqttSessionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnMqttSessionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MqttSessionApiService.CreateMsgVpnMqttSession")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/mqttSessions"
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

type MqttSessionApiApiCreateMsgVpnMqttSessionSubscriptionRequest struct {
	ctx                      context.Context
	ApiService               *MqttSessionApiService
	msgVpnName               string
	mqttSessionClientId      string
	mqttSessionVirtualRouter string
	body                     *MsgVpnMqttSessionSubscription
	opaquePassword           *string
	select_                  *[]string
}

// The Subscription object&#39;s attributes.
func (r MqttSessionApiApiCreateMsgVpnMqttSessionSubscriptionRequest) Body(body MsgVpnMqttSessionSubscription) MqttSessionApiApiCreateMsgVpnMqttSessionSubscriptionRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r MqttSessionApiApiCreateMsgVpnMqttSessionSubscriptionRequest) OpaquePassword(opaquePassword string) MqttSessionApiApiCreateMsgVpnMqttSessionSubscriptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r MqttSessionApiApiCreateMsgVpnMqttSessionSubscriptionRequest) Select_(select_ []string) MqttSessionApiApiCreateMsgVpnMqttSessionSubscriptionRequest {
	r.select_ = &select_
	return r
}

func (r MqttSessionApiApiCreateMsgVpnMqttSessionSubscriptionRequest) Execute() (*MsgVpnMqttSessionSubscriptionResponse, *http.Response, error) {
	return r.ApiService.CreateMsgVpnMqttSessionSubscriptionExecute(r)
}

/*
CreateMsgVpnMqttSessionSubscription Create a Subscription object.

Create a Subscription object. Any attribute missing from the request will be set to its default value. The creation of instances of this object are synchronized to HA mates and replication sites via config-sync.

An MQTT session contains a client's QoS 0 and QoS 1 subscription sets. On creation, a subscription defaults to QoS 0.


Attribute|Identifying|Required|Read-Only|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
mqttSessionClientId|x||x|||
mqttSessionVirtualRouter|x||x|||
msgVpnName|x||x|||
subscriptionTopic|x|x||||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.1.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param mqttSessionClientId The Client ID of the MQTT Session, which corresponds to the ClientId provided in the MQTT CONNECT packet.
 @param mqttSessionVirtualRouter The virtual router of the MQTT Session.
 @return MqttSessionApiApiCreateMsgVpnMqttSessionSubscriptionRequest
*/
func (a *MqttSessionApiService) CreateMsgVpnMqttSessionSubscription(ctx context.Context, msgVpnName string, mqttSessionClientId string, mqttSessionVirtualRouter string) MqttSessionApiApiCreateMsgVpnMqttSessionSubscriptionRequest {
	return MqttSessionApiApiCreateMsgVpnMqttSessionSubscriptionRequest{
		ApiService:               a,
		ctx:                      ctx,
		msgVpnName:               msgVpnName,
		mqttSessionClientId:      mqttSessionClientId,
		mqttSessionVirtualRouter: mqttSessionVirtualRouter,
	}
}

// Execute executes the request
//  @return MsgVpnMqttSessionSubscriptionResponse
func (a *MqttSessionApiService) CreateMsgVpnMqttSessionSubscriptionExecute(r MqttSessionApiApiCreateMsgVpnMqttSessionSubscriptionRequest) (*MsgVpnMqttSessionSubscriptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnMqttSessionSubscriptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MqttSessionApiService.CreateMsgVpnMqttSessionSubscription")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionClientId"+"}", url.PathEscape(parameterToString(r.mqttSessionClientId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionVirtualRouter"+"}", url.PathEscape(parameterToString(r.mqttSessionVirtualRouter, "")), -1)

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

type MqttSessionApiApiDeleteMsgVpnMqttSessionRequest struct {
	ctx                      context.Context
	ApiService               *MqttSessionApiService
	msgVpnName               string
	mqttSessionClientId      string
	mqttSessionVirtualRouter string
}

func (r MqttSessionApiApiDeleteMsgVpnMqttSessionRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnMqttSessionExecute(r)
}

/*
DeleteMsgVpnMqttSession Delete an MQTT Session object.

Delete an MQTT Session object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

An MQTT Session object is a virtual representation of an MQTT client connection. An MQTT session holds the state of an MQTT client (that is, it is used to contain a client's QoS 0 and QoS 1 subscription sets and any undelivered QoS 1 messages).

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.1.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param mqttSessionClientId The Client ID of the MQTT Session, which corresponds to the ClientId provided in the MQTT CONNECT packet.
 @param mqttSessionVirtualRouter The virtual router of the MQTT Session.
 @return MqttSessionApiApiDeleteMsgVpnMqttSessionRequest
*/
func (a *MqttSessionApiService) DeleteMsgVpnMqttSession(ctx context.Context, msgVpnName string, mqttSessionClientId string, mqttSessionVirtualRouter string) MqttSessionApiApiDeleteMsgVpnMqttSessionRequest {
	return MqttSessionApiApiDeleteMsgVpnMqttSessionRequest{
		ApiService:               a,
		ctx:                      ctx,
		msgVpnName:               msgVpnName,
		mqttSessionClientId:      mqttSessionClientId,
		mqttSessionVirtualRouter: mqttSessionVirtualRouter,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *MqttSessionApiService) DeleteMsgVpnMqttSessionExecute(r MqttSessionApiApiDeleteMsgVpnMqttSessionRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MqttSessionApiService.DeleteMsgVpnMqttSession")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionClientId"+"}", url.PathEscape(parameterToString(r.mqttSessionClientId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionVirtualRouter"+"}", url.PathEscape(parameterToString(r.mqttSessionVirtualRouter, "")), -1)

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

type MqttSessionApiApiDeleteMsgVpnMqttSessionSubscriptionRequest struct {
	ctx                      context.Context
	ApiService               *MqttSessionApiService
	msgVpnName               string
	mqttSessionClientId      string
	mqttSessionVirtualRouter string
	subscriptionTopic        string
}

func (r MqttSessionApiApiDeleteMsgVpnMqttSessionSubscriptionRequest) Execute() (*SempMetaOnlyResponse, *http.Response, error) {
	return r.ApiService.DeleteMsgVpnMqttSessionSubscriptionExecute(r)
}

/*
DeleteMsgVpnMqttSessionSubscription Delete a Subscription object.

Delete a Subscription object. The deletion of instances of this object are synchronized to HA mates and replication sites via config-sync.

An MQTT session contains a client's QoS 0 and QoS 1 subscription sets. On creation, a subscription defaults to QoS 0.

A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.1.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param mqttSessionClientId The Client ID of the MQTT Session, which corresponds to the ClientId provided in the MQTT CONNECT packet.
 @param mqttSessionVirtualRouter The virtual router of the MQTT Session.
 @param subscriptionTopic The MQTT subscription topic.
 @return MqttSessionApiApiDeleteMsgVpnMqttSessionSubscriptionRequest
*/
func (a *MqttSessionApiService) DeleteMsgVpnMqttSessionSubscription(ctx context.Context, msgVpnName string, mqttSessionClientId string, mqttSessionVirtualRouter string, subscriptionTopic string) MqttSessionApiApiDeleteMsgVpnMqttSessionSubscriptionRequest {
	return MqttSessionApiApiDeleteMsgVpnMqttSessionSubscriptionRequest{
		ApiService:               a,
		ctx:                      ctx,
		msgVpnName:               msgVpnName,
		mqttSessionClientId:      mqttSessionClientId,
		mqttSessionVirtualRouter: mqttSessionVirtualRouter,
		subscriptionTopic:        subscriptionTopic,
	}
}

// Execute executes the request
//  @return SempMetaOnlyResponse
func (a *MqttSessionApiService) DeleteMsgVpnMqttSessionSubscriptionExecute(r MqttSessionApiApiDeleteMsgVpnMqttSessionSubscriptionRequest) (*SempMetaOnlyResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *SempMetaOnlyResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MqttSessionApiService.DeleteMsgVpnMqttSessionSubscription")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionClientId"+"}", url.PathEscape(parameterToString(r.mqttSessionClientId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionVirtualRouter"+"}", url.PathEscape(parameterToString(r.mqttSessionVirtualRouter, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"subscriptionTopic"+"}", url.PathEscape(parameterToString(r.subscriptionTopic, "")), -1)

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

type MqttSessionApiApiGetMsgVpnMqttSessionRequest struct {
	ctx                      context.Context
	ApiService               *MqttSessionApiService
	msgVpnName               string
	mqttSessionClientId      string
	mqttSessionVirtualRouter string
	opaquePassword           *string
	select_                  *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r MqttSessionApiApiGetMsgVpnMqttSessionRequest) OpaquePassword(opaquePassword string) MqttSessionApiApiGetMsgVpnMqttSessionRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r MqttSessionApiApiGetMsgVpnMqttSessionRequest) Select_(select_ []string) MqttSessionApiApiGetMsgVpnMqttSessionRequest {
	r.select_ = &select_
	return r
}

func (r MqttSessionApiApiGetMsgVpnMqttSessionRequest) Execute() (*MsgVpnMqttSessionResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnMqttSessionExecute(r)
}

/*
GetMsgVpnMqttSession Get an MQTT Session object.

Get an MQTT Session object.

An MQTT Session object is a virtual representation of an MQTT client connection. An MQTT session holds the state of an MQTT client (that is, it is used to contain a client's QoS 0 and QoS 1 subscription sets and any undelivered QoS 1 messages).


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
mqttSessionClientId|x|||
mqttSessionVirtualRouter|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.1.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param mqttSessionClientId The Client ID of the MQTT Session, which corresponds to the ClientId provided in the MQTT CONNECT packet.
 @param mqttSessionVirtualRouter The virtual router of the MQTT Session.
 @return MqttSessionApiApiGetMsgVpnMqttSessionRequest
*/
func (a *MqttSessionApiService) GetMsgVpnMqttSession(ctx context.Context, msgVpnName string, mqttSessionClientId string, mqttSessionVirtualRouter string) MqttSessionApiApiGetMsgVpnMqttSessionRequest {
	return MqttSessionApiApiGetMsgVpnMqttSessionRequest{
		ApiService:               a,
		ctx:                      ctx,
		msgVpnName:               msgVpnName,
		mqttSessionClientId:      mqttSessionClientId,
		mqttSessionVirtualRouter: mqttSessionVirtualRouter,
	}
}

// Execute executes the request
//  @return MsgVpnMqttSessionResponse
func (a *MqttSessionApiService) GetMsgVpnMqttSessionExecute(r MqttSessionApiApiGetMsgVpnMqttSessionRequest) (*MsgVpnMqttSessionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnMqttSessionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MqttSessionApiService.GetMsgVpnMqttSession")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionClientId"+"}", url.PathEscape(parameterToString(r.mqttSessionClientId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionVirtualRouter"+"}", url.PathEscape(parameterToString(r.mqttSessionVirtualRouter, "")), -1)

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

type MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionRequest struct {
	ctx                      context.Context
	ApiService               *MqttSessionApiService
	msgVpnName               string
	mqttSessionClientId      string
	mqttSessionVirtualRouter string
	subscriptionTopic        string
	opaquePassword           *string
	select_                  *[]string
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionRequest) OpaquePassword(opaquePassword string) MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionRequest) Select_(select_ []string) MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionRequest {
	r.select_ = &select_
	return r
}

func (r MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionRequest) Execute() (*MsgVpnMqttSessionSubscriptionResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnMqttSessionSubscriptionExecute(r)
}

/*
GetMsgVpnMqttSessionSubscription Get a Subscription object.

Get a Subscription object.

An MQTT session contains a client's QoS 0 and QoS 1 subscription sets. On creation, a subscription defaults to QoS 0.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
mqttSessionClientId|x|||
mqttSessionVirtualRouter|x|||
msgVpnName|x|||
subscriptionTopic|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.1.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param mqttSessionClientId The Client ID of the MQTT Session, which corresponds to the ClientId provided in the MQTT CONNECT packet.
 @param mqttSessionVirtualRouter The virtual router of the MQTT Session.
 @param subscriptionTopic The MQTT subscription topic.
 @return MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionRequest
*/
func (a *MqttSessionApiService) GetMsgVpnMqttSessionSubscription(ctx context.Context, msgVpnName string, mqttSessionClientId string, mqttSessionVirtualRouter string, subscriptionTopic string) MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionRequest {
	return MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionRequest{
		ApiService:               a,
		ctx:                      ctx,
		msgVpnName:               msgVpnName,
		mqttSessionClientId:      mqttSessionClientId,
		mqttSessionVirtualRouter: mqttSessionVirtualRouter,
		subscriptionTopic:        subscriptionTopic,
	}
}

// Execute executes the request
//  @return MsgVpnMqttSessionSubscriptionResponse
func (a *MqttSessionApiService) GetMsgVpnMqttSessionSubscriptionExecute(r MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionRequest) (*MsgVpnMqttSessionSubscriptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnMqttSessionSubscriptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MqttSessionApiService.GetMsgVpnMqttSessionSubscription")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionClientId"+"}", url.PathEscape(parameterToString(r.mqttSessionClientId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionVirtualRouter"+"}", url.PathEscape(parameterToString(r.mqttSessionVirtualRouter, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"subscriptionTopic"+"}", url.PathEscape(parameterToString(r.subscriptionTopic, "")), -1)

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

type MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest struct {
	ctx                      context.Context
	ApiService               *MqttSessionApiService
	msgVpnName               string
	mqttSessionClientId      string
	mqttSessionVirtualRouter string
	count                    *int32
	cursor                   *string
	opaquePassword           *string
	where                    *[]string
	select_                  *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest) Count(count int32) MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest {
	r.count = &count
	return r
}

// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest) Cursor(cursor string) MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest {
	r.cursor = &cursor
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest) OpaquePassword(opaquePassword string) MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest) Where(where []string) MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest {
	r.where = &where
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest) Select_(select_ []string) MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest {
	r.select_ = &select_
	return r
}

func (r MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest) Execute() (*MsgVpnMqttSessionSubscriptionsResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnMqttSessionSubscriptionsExecute(r)
}

/*
GetMsgVpnMqttSessionSubscriptions Get a list of Subscription objects.

Get a list of Subscription objects.

An MQTT session contains a client's QoS 0 and QoS 1 subscription sets. On creation, a subscription defaults to QoS 0.


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
mqttSessionClientId|x|||
mqttSessionVirtualRouter|x|||
msgVpnName|x|||
subscriptionTopic|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.1.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param mqttSessionClientId The Client ID of the MQTT Session, which corresponds to the ClientId provided in the MQTT CONNECT packet.
 @param mqttSessionVirtualRouter The virtual router of the MQTT Session.
 @return MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest
*/
func (a *MqttSessionApiService) GetMsgVpnMqttSessionSubscriptions(ctx context.Context, msgVpnName string, mqttSessionClientId string, mqttSessionVirtualRouter string) MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest {
	return MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest{
		ApiService:               a,
		ctx:                      ctx,
		msgVpnName:               msgVpnName,
		mqttSessionClientId:      mqttSessionClientId,
		mqttSessionVirtualRouter: mqttSessionVirtualRouter,
	}
}

// Execute executes the request
//  @return MsgVpnMqttSessionSubscriptionsResponse
func (a *MqttSessionApiService) GetMsgVpnMqttSessionSubscriptionsExecute(r MqttSessionApiApiGetMsgVpnMqttSessionSubscriptionsRequest) (*MsgVpnMqttSessionSubscriptionsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnMqttSessionSubscriptionsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MqttSessionApiService.GetMsgVpnMqttSessionSubscriptions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionClientId"+"}", url.PathEscape(parameterToString(r.mqttSessionClientId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionVirtualRouter"+"}", url.PathEscape(parameterToString(r.mqttSessionVirtualRouter, "")), -1)

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

type MqttSessionApiApiGetMsgVpnMqttSessionsRequest struct {
	ctx            context.Context
	ApiService     *MqttSessionApiService
	msgVpnName     string
	count          *int32
	cursor         *string
	opaquePassword *string
	where          *[]string
	select_        *[]string
}

// Limit the count of objects in the response. See the documentation for the &#x60;count&#x60; parameter.
func (r MqttSessionApiApiGetMsgVpnMqttSessionsRequest) Count(count int32) MqttSessionApiApiGetMsgVpnMqttSessionsRequest {
	r.count = &count
	return r
}

// The cursor, or position, for the next page of objects. See the documentation for the &#x60;cursor&#x60; parameter.
func (r MqttSessionApiApiGetMsgVpnMqttSessionsRequest) Cursor(cursor string) MqttSessionApiApiGetMsgVpnMqttSessionsRequest {
	r.cursor = &cursor
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r MqttSessionApiApiGetMsgVpnMqttSessionsRequest) OpaquePassword(opaquePassword string) MqttSessionApiApiGetMsgVpnMqttSessionsRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only objects where certain conditions are true. See the the documentation for the &#x60;where&#x60; parameter.
func (r MqttSessionApiApiGetMsgVpnMqttSessionsRequest) Where(where []string) MqttSessionApiApiGetMsgVpnMqttSessionsRequest {
	r.where = &where
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r MqttSessionApiApiGetMsgVpnMqttSessionsRequest) Select_(select_ []string) MqttSessionApiApiGetMsgVpnMqttSessionsRequest {
	r.select_ = &select_
	return r
}

func (r MqttSessionApiApiGetMsgVpnMqttSessionsRequest) Execute() (*MsgVpnMqttSessionsResponse, *http.Response, error) {
	return r.ApiService.GetMsgVpnMqttSessionsExecute(r)
}

/*
GetMsgVpnMqttSessions Get a list of MQTT Session objects.

Get a list of MQTT Session objects.

An MQTT Session object is a virtual representation of an MQTT client connection. An MQTT session holds the state of an MQTT client (that is, it is used to contain a client's QoS 0 and QoS 1 subscription sets and any undelivered QoS 1 messages).


Attribute|Identifying|Write-Only|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:
mqttSessionClientId|x|||
mqttSessionVirtualRouter|x|||
msgVpnName|x|||



A SEMP client authorized with a minimum access scope/level of "vpn/read-only" is required to perform this operation.

This has been available since 2.1.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @return MqttSessionApiApiGetMsgVpnMqttSessionsRequest
*/
func (a *MqttSessionApiService) GetMsgVpnMqttSessions(ctx context.Context, msgVpnName string) MqttSessionApiApiGetMsgVpnMqttSessionsRequest {
	return MqttSessionApiApiGetMsgVpnMqttSessionsRequest{
		ApiService: a,
		ctx:        ctx,
		msgVpnName: msgVpnName,
	}
}

// Execute executes the request
//  @return MsgVpnMqttSessionsResponse
func (a *MqttSessionApiService) GetMsgVpnMqttSessionsExecute(r MqttSessionApiApiGetMsgVpnMqttSessionsRequest) (*MsgVpnMqttSessionsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnMqttSessionsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MqttSessionApiService.GetMsgVpnMqttSessions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/mqttSessions"
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

type MqttSessionApiApiReplaceMsgVpnMqttSessionRequest struct {
	ctx                      context.Context
	ApiService               *MqttSessionApiService
	msgVpnName               string
	mqttSessionClientId      string
	mqttSessionVirtualRouter string
	body                     *MsgVpnMqttSession
	opaquePassword           *string
	select_                  *[]string
}

// The MQTT Session object&#39;s attributes.
func (r MqttSessionApiApiReplaceMsgVpnMqttSessionRequest) Body(body MsgVpnMqttSession) MqttSessionApiApiReplaceMsgVpnMqttSessionRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r MqttSessionApiApiReplaceMsgVpnMqttSessionRequest) OpaquePassword(opaquePassword string) MqttSessionApiApiReplaceMsgVpnMqttSessionRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r MqttSessionApiApiReplaceMsgVpnMqttSessionRequest) Select_(select_ []string) MqttSessionApiApiReplaceMsgVpnMqttSessionRequest {
	r.select_ = &select_
	return r
}

func (r MqttSessionApiApiReplaceMsgVpnMqttSessionRequest) Execute() (*MsgVpnMqttSessionResponse, *http.Response, error) {
	return r.ApiService.ReplaceMsgVpnMqttSessionExecute(r)
}

/*
ReplaceMsgVpnMqttSession Replace an MQTT Session object.

Replace an MQTT Session object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

An MQTT Session object is a virtual representation of an MQTT client connection. An MQTT session holds the state of an MQTT client (that is, it is used to contain a client's QoS 0 and QoS 1 subscription sets and any undelivered QoS 1 messages).


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
mqttSessionClientId|x|x||||
mqttSessionVirtualRouter|x|x||||
msgVpnName|x|x||||
owner||||x||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
EventThreshold|clearPercent|setPercent|clearValue, setValue
EventThreshold|clearValue|setValue|clearPercent, setPercent
EventThreshold|setPercent|clearPercent|clearValue, setValue
EventThreshold|setValue|clearValue|clearPercent, setPercent



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.1.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param mqttSessionClientId The Client ID of the MQTT Session, which corresponds to the ClientId provided in the MQTT CONNECT packet.
 @param mqttSessionVirtualRouter The virtual router of the MQTT Session.
 @return MqttSessionApiApiReplaceMsgVpnMqttSessionRequest
*/
func (a *MqttSessionApiService) ReplaceMsgVpnMqttSession(ctx context.Context, msgVpnName string, mqttSessionClientId string, mqttSessionVirtualRouter string) MqttSessionApiApiReplaceMsgVpnMqttSessionRequest {
	return MqttSessionApiApiReplaceMsgVpnMqttSessionRequest{
		ApiService:               a,
		ctx:                      ctx,
		msgVpnName:               msgVpnName,
		mqttSessionClientId:      mqttSessionClientId,
		mqttSessionVirtualRouter: mqttSessionVirtualRouter,
	}
}

// Execute executes the request
//  @return MsgVpnMqttSessionResponse
func (a *MqttSessionApiService) ReplaceMsgVpnMqttSessionExecute(r MqttSessionApiApiReplaceMsgVpnMqttSessionRequest) (*MsgVpnMqttSessionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnMqttSessionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MqttSessionApiService.ReplaceMsgVpnMqttSession")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionClientId"+"}", url.PathEscape(parameterToString(r.mqttSessionClientId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionVirtualRouter"+"}", url.PathEscape(parameterToString(r.mqttSessionVirtualRouter, "")), -1)

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

type MqttSessionApiApiReplaceMsgVpnMqttSessionSubscriptionRequest struct {
	ctx                      context.Context
	ApiService               *MqttSessionApiService
	msgVpnName               string
	mqttSessionClientId      string
	mqttSessionVirtualRouter string
	subscriptionTopic        string
	body                     *MsgVpnMqttSessionSubscription
	opaquePassword           *string
	select_                  *[]string
}

// The Subscription object&#39;s attributes.
func (r MqttSessionApiApiReplaceMsgVpnMqttSessionSubscriptionRequest) Body(body MsgVpnMqttSessionSubscription) MqttSessionApiApiReplaceMsgVpnMqttSessionSubscriptionRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r MqttSessionApiApiReplaceMsgVpnMqttSessionSubscriptionRequest) OpaquePassword(opaquePassword string) MqttSessionApiApiReplaceMsgVpnMqttSessionSubscriptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r MqttSessionApiApiReplaceMsgVpnMqttSessionSubscriptionRequest) Select_(select_ []string) MqttSessionApiApiReplaceMsgVpnMqttSessionSubscriptionRequest {
	r.select_ = &select_
	return r
}

func (r MqttSessionApiApiReplaceMsgVpnMqttSessionSubscriptionRequest) Execute() (*MsgVpnMqttSessionSubscriptionResponse, *http.Response, error) {
	return r.ApiService.ReplaceMsgVpnMqttSessionSubscriptionExecute(r)
}

/*
ReplaceMsgVpnMqttSessionSubscription Replace a Subscription object.

Replace a Subscription object. Any attribute missing from the request will be set to its default value, subject to the exceptions in note 4.

An MQTT session contains a client's QoS 0 and QoS 1 subscription sets. On creation, a subscription defaults to QoS 0.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
mqttSessionClientId|x|x||||
mqttSessionVirtualRouter|x|x||||
msgVpnName|x|x||||
subscriptionTopic|x|x||||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.1.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param mqttSessionClientId The Client ID of the MQTT Session, which corresponds to the ClientId provided in the MQTT CONNECT packet.
 @param mqttSessionVirtualRouter The virtual router of the MQTT Session.
 @param subscriptionTopic The MQTT subscription topic.
 @return MqttSessionApiApiReplaceMsgVpnMqttSessionSubscriptionRequest
*/
func (a *MqttSessionApiService) ReplaceMsgVpnMqttSessionSubscription(ctx context.Context, msgVpnName string, mqttSessionClientId string, mqttSessionVirtualRouter string, subscriptionTopic string) MqttSessionApiApiReplaceMsgVpnMqttSessionSubscriptionRequest {
	return MqttSessionApiApiReplaceMsgVpnMqttSessionSubscriptionRequest{
		ApiService:               a,
		ctx:                      ctx,
		msgVpnName:               msgVpnName,
		mqttSessionClientId:      mqttSessionClientId,
		mqttSessionVirtualRouter: mqttSessionVirtualRouter,
		subscriptionTopic:        subscriptionTopic,
	}
}

// Execute executes the request
//  @return MsgVpnMqttSessionSubscriptionResponse
func (a *MqttSessionApiService) ReplaceMsgVpnMqttSessionSubscriptionExecute(r MqttSessionApiApiReplaceMsgVpnMqttSessionSubscriptionRequest) (*MsgVpnMqttSessionSubscriptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnMqttSessionSubscriptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MqttSessionApiService.ReplaceMsgVpnMqttSessionSubscription")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionClientId"+"}", url.PathEscape(parameterToString(r.mqttSessionClientId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionVirtualRouter"+"}", url.PathEscape(parameterToString(r.mqttSessionVirtualRouter, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"subscriptionTopic"+"}", url.PathEscape(parameterToString(r.subscriptionTopic, "")), -1)

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

type MqttSessionApiApiUpdateMsgVpnMqttSessionRequest struct {
	ctx                      context.Context
	ApiService               *MqttSessionApiService
	msgVpnName               string
	mqttSessionClientId      string
	mqttSessionVirtualRouter string
	body                     *MsgVpnMqttSession
	opaquePassword           *string
	select_                  *[]string
}

// The MQTT Session object&#39;s attributes.
func (r MqttSessionApiApiUpdateMsgVpnMqttSessionRequest) Body(body MsgVpnMqttSession) MqttSessionApiApiUpdateMsgVpnMqttSessionRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r MqttSessionApiApiUpdateMsgVpnMqttSessionRequest) OpaquePassword(opaquePassword string) MqttSessionApiApiUpdateMsgVpnMqttSessionRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r MqttSessionApiApiUpdateMsgVpnMqttSessionRequest) Select_(select_ []string) MqttSessionApiApiUpdateMsgVpnMqttSessionRequest {
	r.select_ = &select_
	return r
}

func (r MqttSessionApiApiUpdateMsgVpnMqttSessionRequest) Execute() (*MsgVpnMqttSessionResponse, *http.Response, error) {
	return r.ApiService.UpdateMsgVpnMqttSessionExecute(r)
}

/*
UpdateMsgVpnMqttSession Update an MQTT Session object.

Update an MQTT Session object. Any attribute missing from the request will be left unchanged.

An MQTT Session object is a virtual representation of an MQTT client connection. An MQTT session holds the state of an MQTT client (that is, it is used to contain a client's QoS 0 and QoS 1 subscription sets and any undelivered QoS 1 messages).


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
mqttSessionClientId|x|x||||
mqttSessionVirtualRouter|x|x||||
msgVpnName|x|x||||
owner||||x||



The following attributes in the request may only be provided in certain combinations with other attributes:


Class|Attribute|Requires|Conflicts
:---|:---|:---|:---
EventThreshold|clearPercent|setPercent|clearValue, setValue
EventThreshold|clearValue|setValue|clearPercent, setPercent
EventThreshold|setPercent|clearPercent|clearValue, setValue
EventThreshold|setValue|clearValue|clearPercent, setPercent



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.1.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param mqttSessionClientId The Client ID of the MQTT Session, which corresponds to the ClientId provided in the MQTT CONNECT packet.
 @param mqttSessionVirtualRouter The virtual router of the MQTT Session.
 @return MqttSessionApiApiUpdateMsgVpnMqttSessionRequest
*/
func (a *MqttSessionApiService) UpdateMsgVpnMqttSession(ctx context.Context, msgVpnName string, mqttSessionClientId string, mqttSessionVirtualRouter string) MqttSessionApiApiUpdateMsgVpnMqttSessionRequest {
	return MqttSessionApiApiUpdateMsgVpnMqttSessionRequest{
		ApiService:               a,
		ctx:                      ctx,
		msgVpnName:               msgVpnName,
		mqttSessionClientId:      mqttSessionClientId,
		mqttSessionVirtualRouter: mqttSessionVirtualRouter,
	}
}

// Execute executes the request
//  @return MsgVpnMqttSessionResponse
func (a *MqttSessionApiService) UpdateMsgVpnMqttSessionExecute(r MqttSessionApiApiUpdateMsgVpnMqttSessionRequest) (*MsgVpnMqttSessionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnMqttSessionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MqttSessionApiService.UpdateMsgVpnMqttSession")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionClientId"+"}", url.PathEscape(parameterToString(r.mqttSessionClientId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionVirtualRouter"+"}", url.PathEscape(parameterToString(r.mqttSessionVirtualRouter, "")), -1)

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

type MqttSessionApiApiUpdateMsgVpnMqttSessionSubscriptionRequest struct {
	ctx                      context.Context
	ApiService               *MqttSessionApiService
	msgVpnName               string
	mqttSessionClientId      string
	mqttSessionVirtualRouter string
	subscriptionTopic        string
	body                     *MsgVpnMqttSessionSubscription
	opaquePassword           *string
	select_                  *[]string
}

// The Subscription object&#39;s attributes.
func (r MqttSessionApiApiUpdateMsgVpnMqttSessionSubscriptionRequest) Body(body MsgVpnMqttSessionSubscription) MqttSessionApiApiUpdateMsgVpnMqttSessionSubscriptionRequest {
	r.body = &body
	return r
}

// Accept opaque attributes in the request or return opaque attributes in the response, encrypted with the specified password. See the documentation for the &#x60;opaquePassword&#x60; parameter.
func (r MqttSessionApiApiUpdateMsgVpnMqttSessionSubscriptionRequest) OpaquePassword(opaquePassword string) MqttSessionApiApiUpdateMsgVpnMqttSessionSubscriptionRequest {
	r.opaquePassword = &opaquePassword
	return r
}

// Include in the response only selected attributes of the object, or exclude from the response selected attributes of the object. See the documentation for the &#x60;select&#x60; parameter.
func (r MqttSessionApiApiUpdateMsgVpnMqttSessionSubscriptionRequest) Select_(select_ []string) MqttSessionApiApiUpdateMsgVpnMqttSessionSubscriptionRequest {
	r.select_ = &select_
	return r
}

func (r MqttSessionApiApiUpdateMsgVpnMqttSessionSubscriptionRequest) Execute() (*MsgVpnMqttSessionSubscriptionResponse, *http.Response, error) {
	return r.ApiService.UpdateMsgVpnMqttSessionSubscriptionExecute(r)
}

/*
UpdateMsgVpnMqttSessionSubscription Update a Subscription object.

Update a Subscription object. Any attribute missing from the request will be left unchanged.

An MQTT session contains a client's QoS 0 and QoS 1 subscription sets. On creation, a subscription defaults to QoS 0.


Attribute|Identifying|Read-Only|Write-Only|Requires-Disable|Deprecated|Opaque
:---|:---:|:---:|:---:|:---:|:---:|:---:
mqttSessionClientId|x|x||||
mqttSessionVirtualRouter|x|x||||
msgVpnName|x|x||||
subscriptionTopic|x|x||||



A SEMP client authorized with a minimum access scope/level of "vpn/read-write" is required to perform this operation.

This has been available since 2.1.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param msgVpnName The name of the Message VPN.
 @param mqttSessionClientId The Client ID of the MQTT Session, which corresponds to the ClientId provided in the MQTT CONNECT packet.
 @param mqttSessionVirtualRouter The virtual router of the MQTT Session.
 @param subscriptionTopic The MQTT subscription topic.
 @return MqttSessionApiApiUpdateMsgVpnMqttSessionSubscriptionRequest
*/
func (a *MqttSessionApiService) UpdateMsgVpnMqttSessionSubscription(ctx context.Context, msgVpnName string, mqttSessionClientId string, mqttSessionVirtualRouter string, subscriptionTopic string) MqttSessionApiApiUpdateMsgVpnMqttSessionSubscriptionRequest {
	return MqttSessionApiApiUpdateMsgVpnMqttSessionSubscriptionRequest{
		ApiService:               a,
		ctx:                      ctx,
		msgVpnName:               msgVpnName,
		mqttSessionClientId:      mqttSessionClientId,
		mqttSessionVirtualRouter: mqttSessionVirtualRouter,
		subscriptionTopic:        subscriptionTopic,
	}
}

// Execute executes the request
//  @return MsgVpnMqttSessionSubscriptionResponse
func (a *MqttSessionApiService) UpdateMsgVpnMqttSessionSubscriptionExecute(r MqttSessionApiApiUpdateMsgVpnMqttSessionSubscriptionRequest) (*MsgVpnMqttSessionSubscriptionResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *MsgVpnMqttSessionSubscriptionResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MqttSessionApiService.UpdateMsgVpnMqttSessionSubscription")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/msgVpns/{msgVpnName}/mqttSessions/{mqttSessionClientId},{mqttSessionVirtualRouter}/subscriptions/{subscriptionTopic}"
	localVarPath = strings.Replace(localVarPath, "{"+"msgVpnName"+"}", url.PathEscape(parameterToString(r.msgVpnName, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionClientId"+"}", url.PathEscape(parameterToString(r.mqttSessionClientId, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"mqttSessionVirtualRouter"+"}", url.PathEscape(parameterToString(r.mqttSessionVirtualRouter, "")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"subscriptionTopic"+"}", url.PathEscape(parameterToString(r.subscriptionTopic, "")), -1)

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
