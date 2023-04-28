package provider

import (
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func NewMsgVpnAclProfilePublishExceptionResource() resource.Resource {
	return &solaceResource[MsgVpnAclProfilePublishException]{spr: &aclProfilePublishExceptionResource{}}
}

var _ solaceProviderResource[MsgVpnAclProfilePublishException] = &aclProfilePublishExceptionResource{}

type aclProfilePublishExceptionResource struct {
	*solaceProvider
}

func (r aclProfilePublishExceptionResource) Name() string {
	return "aclprofile_publish_exception"
}

func (r aclProfilePublishExceptionResource) Schema() schema.Schema {
	return MsgVpnAclProfilePublishExceptionResourceSchema("msg_vpn_name", "acl_profile_name", "topic_syntax", "publish_exception_topic")
}

func (r *aclProfilePublishExceptionResource) SetProvider(provider *solaceProvider) {
	r.solaceProvider = provider
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

var msgVpnAclProfilePublishExceptionImportRegexp *regexp.Regexp = regexp.MustCompile(
	`^([^\s\*\?\/]+)\/([0-9a-zA-Z_\-]+)\/(smf|mqtt)\/([^\s]+)$`)

func (r aclProfilePublishExceptionResource) Import(id string, data *MsgVpnAclProfilePublishException, diag *diag.Diagnostics) {
	match := msgVpnAclProfilePublishExceptionImportRegexp.FindStringSubmatch(id)
	if match != nil {
		data.MsgVpnName = &match[1]
		data.AclProfileName = &match[2]
		data.TopicSyntax = &match[3]
		data.PublishExceptionTopic = &match[4]
	} else {
		diag.AddError("Expected <vpn-name>/<acl-profile>/<topic syntax>/<topic name>",
			id+" does not match "+msgVpnAclProfilePublishExceptionImportRegexp.String())
		return
	}
}
