package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.ResourceType = aclProfileResourceType{}

type aclProfileSubscribeExceptionResourceType struct {
}

func (t aclProfileSubscribeExceptionResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return NewResource[MsgVpnAclProfileSubscribeException](
		aclProfileSubscribeExceptionResource{provider: provider}), diags
}

func (t aclProfileSubscribeExceptionResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return MsgVpnAclProfileSubscribeExceptionSchema("msg_vpn_name", "acl_profile_name", "topic_syntax", "subscribe_exception_topic"), nil
}

var _ solaceProviderResource[MsgVpnAclProfile] = aclProfileResource{}

type aclProfileSubscribeExceptionResource struct {
	provider
}

func (r aclProfileSubscribeExceptionResource) NewData() *MsgVpnAclProfileSubscribeException {
	return &MsgVpnAclProfileSubscribeException{}
}

func (r aclProfileSubscribeExceptionResource) Create(data *MsgVpnAclProfileSubscribeException, diag *diag.Diagnostics) (*http.Response, error) {
	_, httpResponse, err := r.Client.AclProfileApi.
		CreateMsgVpnAclProfileSubscribeException(r.Context, *data.MsgVpnName, *data.AclProfileName).
		Body(*data.ToApi()).
		Execute()
	return httpResponse, err
}

func (r aclProfileSubscribeExceptionResource) Read(data *MsgVpnAclProfileSubscribeException, diag *diag.Diagnostics) (*http.Response, error) {
	apiResponse, httpResponse, err := r.Client.AclProfileApi.
		GetMsgVpnAclProfileSubscribeException(
			r.Context, *data.MsgVpnName, *data.AclProfileName, *data.TopicSyntax, *data.SubscribeExceptionTopic).
		Execute()
	if err == nil && apiResponse != nil && apiResponse.Data != nil {
		data.ToTF(apiResponse.Data)
	}
	return httpResponse, err
}

func (r aclProfileSubscribeExceptionResource) Update(_ *MsgVpnAclProfileSubscribeException, data *MsgVpnAclProfileSubscribeException, diag *diag.Diagnostics) (*http.Response, error) {
	diag.AddError("Update not supported", "Acl Profile Client Connect Exceptions cannot be updated")
	return nil, nil
}

func (r aclProfileSubscribeExceptionResource) Delete(data *MsgVpnAclProfileSubscribeException, diag *diag.Diagnostics) (*http.Response, error) {
	_, httpResponse, err := r.Client.AclProfileApi.
		DeleteMsgVpnAclProfileSubscribeException(
			r.Context, *data.MsgVpnName, *data.AclProfileName, *data.TopicSyntax, *data.SubscribeExceptionTopic).
		Execute()
	return httpResponse, err
}

func (r aclProfileSubscribeExceptionResource) Import(*MsgVpnAclProfileSubscribeException, *diag.Diagnostics) {
}
