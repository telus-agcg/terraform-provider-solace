package provider

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func NewMsgVpnQueueSubscriptionResource() resource.Resource {
	return &solaceResource[MsgVpnQueueSubscription]{spr: &queueSubscriptionResource{}}
}

var _ solaceProviderResource[MsgVpnQueueSubscription] = &queueSubscriptionResource{}

type queueSubscriptionResource struct {
	*solaceProvider
}

func (r queueSubscriptionResource) Name() string {
	return "queue_subscription"
}

func (r queueSubscriptionResource) Schema() schema.Schema {
	return MsgVpnQueueSubscriptionResourceSchema("msg_vpn_name", "queue_name", "subscription_topic")
}

func (r *queueSubscriptionResource) SetProvider(provider *solaceProvider) {
	r.solaceProvider = provider
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
