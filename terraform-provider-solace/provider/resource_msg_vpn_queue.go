package provider

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.ResourceType = queueResourceType{}

type queueResourceType struct {
}

func (t queueResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return NewResource[MsgVpnQueue](
		queueResource{provider: provider}), diags
}

func (t queueResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return MsgVpnQueueSchema("msg_vpn_name", "queue_name"), nil
}

var _ solaceProviderResource[MsgVpnQueue] = queueResource{}

type queueResource struct {
	provider
}

func (r queueResource) NewData() *MsgVpnQueue {
	return &MsgVpnQueue{}
}

func (r queueResource) Create(data *MsgVpnQueue, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.QueueApi.CreateMsgVpnQueue(r.Context, *data.MsgVpnName).Body(*data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r queueResource) Read(data *MsgVpnQueue, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.QueueApi.GetMsgVpnQueue(r.Context, *data.MsgVpnName, *data.QueueName)
	apiResponse, httpResponse, err := apiReq.Execute()
	if err == nil && apiResponse != nil && apiResponse.Data != nil {
		data.ToTF(apiResponse.Data)
	}
	return httpResponse, err
}

func (r queueResource) Update(cur *MsgVpnQueue, pln *MsgVpnQueue, diag *diag.Diagnostics) (*http.Response, error) {
	requiresShutdown := false
	HasChanged(cur.AccessType, pln.AccessType, &requiresShutdown)
	HasChanged(cur.Owner, pln.Owner, &requiresShutdown)
	HasChanged(cur.Permission, pln.Permission, &requiresShutdown)
	HasChanged(cur.RespectMsgPriorityEnabled, pln.RespectMsgPriorityEnabled, &requiresShutdown)

	apiPlan := pln.ToApi()
	if requiresShutdown {
		apiPlan.IngressEnabled = sempv2.PtrBool(false)
		apiPlan.EgressEnabled = sempv2.PtrBool(false)
	}

	_, httpResponse, err := r.Client.QueueApi.
		UpdateMsgVpnQueue(r.Context, *pln.MsgVpnName, *pln.QueueName).
		Body(*apiPlan).
		Execute()

	// If the queue needed shut down before, set the enabled flags
	// to the state required in the plan
	if err == nil && requiresShutdown {
		body := sempv2.MsgVpnQueue{
			IngressEnabled: pln.IngressEnabled,
			EgressEnabled:  pln.EgressEnabled,
		}
		_, httpResponse, err = r.Client.QueueApi.UpdateMsgVpnQueue(
			r.Context, *pln.MsgVpnName, *pln.QueueName).Body(body).Execute()
	}

	return httpResponse, err
}

func (r queueResource) Delete(data *MsgVpnQueue, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.QueueApi.DeleteMsgVpnQueue(r.Context, *data.MsgVpnName, *data.QueueName)
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

var msgVpnQueueImportRegexp *regexp.Regexp = regexp.MustCompile(fmt.Sprintf(
	"^([^\\s%s]+)\\/([^\\s%s]+)$", regexp.QuoteMeta("*?/"), regexp.QuoteMeta("'<>*?&;)")))

func (r queueResource) Import(id string, data *MsgVpnQueue, diag *diag.Diagnostics) {
	match := msgVpnQueueImportRegexp.FindStringSubmatch(id)
	if match != nil {
		data.MsgVpnName = &match[1]
		data.QueueName = &match[2]
	} else {
		diag.AddError("Expected <vpn-name>/<queue-name>", id+" does not match "+msgVpnQueueImportRegexp.String())
		return
	}
}
