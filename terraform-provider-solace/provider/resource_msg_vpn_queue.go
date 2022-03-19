package provider

import (
	"context"
	"net/http"

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
	apiReq := r.Client.QueueApi.CreateMsgVpnQueue(r.Context, *data.MsgVpnName).Body(data.ToApi())
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

func (r queueResource) Update(data *MsgVpnQueue, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.QueueApi.UpdateMsgVpnQueue(r.Context, *data.MsgVpnName, *data.QueueName).Body(data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r queueResource) Delete(data *MsgVpnQueue, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.QueueApi.DeleteMsgVpnQueue(r.Context, *data.MsgVpnName, *data.QueueName)
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r queueResource) Import(*MsgVpnQueue, *diag.Diagnostics) {}
