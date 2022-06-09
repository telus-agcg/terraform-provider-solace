package provider

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.ResourceType = queueSubscriptionResourceType{}

type queueSubscriptionResourceType struct {
}

func (t queueSubscriptionResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return NewResource[MsgVpnQueueSubscription](
		queueSubscriptionResource{provider: provider}), diags
}

func (t queueSubscriptionResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return MsgVpnQueueSubscriptionSchema("msg_vpn_name", "queue_name", "subscription_topic"), nil
}

var _ solaceProviderResource[MsgVpnQueueSubscription] = queueSubscriptionResource{}

type queueSubscriptionResource struct {
	provider
}

func (r queueSubscriptionResource) NewData() *MsgVpnQueueSubscription {
	return &MsgVpnQueueSubscription{}
}

func (r queueSubscriptionResource) Create(data *MsgVpnQueueSubscription, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.QueueApi.CreateMsgVpnQueueSubscription(r.Context, *data.MsgVpnName, *data.QueueName).Body(*data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r queueSubscriptionResource) Read(data *MsgVpnQueueSubscription, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.QueueApi.GetMsgVpnQueueSubscription(r.Context, *data.MsgVpnName, *data.QueueName, *data.SubscriptionTopic)
	apiResponse, httpResponse, err := apiReq.Execute()
	if err == nil && apiResponse != nil && apiResponse.Data != nil {
		data.ToTF(apiResponse.Data)
	}
	return httpResponse, err
}

func (r queueSubscriptionResource) Update(cur *MsgVpnQueueSubscription, pln *MsgVpnQueueSubscription, diag *diag.Diagnostics) (*http.Response, error) {
	diag.AddError("Update not supported", "Queue subscriptions cannot be updated")
	return nil, nil
}

func (r queueSubscriptionResource) Delete(data *MsgVpnQueueSubscription, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.QueueApi.DeleteMsgVpnQueueSubscription(r.Context, *data.MsgVpnName, *data.QueueName, *data.SubscriptionTopic)
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

var msgVpnQueueSubscriptionImportRegexp *regexp.Regexp = regexp.MustCompile(fmt.Sprintf(
	"^([^\\s%s]+)\\/([^\\s%s]+)\\/([^\\s]+)$", regexp.QuoteMeta("*?/"), regexp.QuoteMeta("'<>*?&;)")))

func (r queueSubscriptionResource) Import(id string, data *MsgVpnQueueSubscription, diag *diag.Diagnostics) {
	match := msgVpnQueueSubscriptionImportRegexp.FindStringSubmatch(id)
	if match != nil {
		data.MsgVpnName = &match[1]
		data.QueueName = &match[2]
		data.SubscriptionTopic = &match[3]
	} else {
		diag.AddError("Expected <vpn-name>/<queue-name>/<subscription>", id+" does not match "+msgVpnQueueImportRegexp.String())
		return
	}
}
