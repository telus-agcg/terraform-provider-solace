package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.ResourceType = aclProfileResourceType{}

type aclProfilePublishExceptionResourceType struct {
}

func (t aclProfilePublishExceptionResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return NewResource[MsgVpnAclProfilePublishException](
		aclProfilePublishExceptionResource{provider: provider}), diags
}

func (t aclProfilePublishExceptionResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return MsgVpnAclProfilePublishExceptionSchema("msg_vpn_name", "acl_profile_name", "topic_syntax", "publish_exception_topic"), nil
}

var _ solaceProviderResource[MsgVpnAclProfile] = aclProfileResource{}

type aclProfilePublishExceptionResource struct {
	provider
}

func (r aclProfilePublishExceptionResource) NewData() *MsgVpnAclProfilePublishException {
	return &MsgVpnAclProfilePublishException{}
}

func (r aclProfilePublishExceptionResource) Create(data *MsgVpnAclProfilePublishException, diag *diag.Diagnostics) (*http.Response, error) {
	_, httpResponse, err := r.Client.AclProfileApi.
		CreateMsgVpnAclProfilePublishException(r.Context, *data.MsgVpnName, *data.AclProfileName).
		Body(*data.ToApi()).
		Execute()
	return httpResponse, err
}

func (r aclProfilePublishExceptionResource) Read(data *MsgVpnAclProfilePublishException, diag *diag.Diagnostics) (*http.Response, error) {
	apiResponse, httpResponse, err := r.Client.AclProfileApi.
		GetMsgVpnAclProfilePublishException(
			r.Context, *data.MsgVpnName, *data.AclProfileName, *data.TopicSyntax, *data.PublishExceptionTopic).
		Execute()
	if err == nil && apiResponse != nil && apiResponse.Data != nil {
		data.ToTF(apiResponse.Data)
	}
	return httpResponse, err
}

func (r aclProfilePublishExceptionResource) Update(_ *MsgVpnAclProfilePublishException, data *MsgVpnAclProfilePublishException, diag *diag.Diagnostics) (*http.Response, error) {
	diag.AddError("Update not supported", "Acl Profile Client Connect Exceptions cannot be updated")
	return nil, nil
}

func (r aclProfilePublishExceptionResource) Delete(data *MsgVpnAclProfilePublishException, diag *diag.Diagnostics) (*http.Response, error) {
	_, httpResponse, err := r.Client.AclProfileApi.
		DeleteMsgVpnAclProfilePublishException(
			r.Context, *data.MsgVpnName, *data.AclProfileName, *data.TopicSyntax, *data.PublishExceptionTopic).
		Execute()
	return httpResponse, err
}

func (r aclProfilePublishExceptionResource) Import(*MsgVpnAclProfilePublishException, *diag.Diagnostics) {
}
