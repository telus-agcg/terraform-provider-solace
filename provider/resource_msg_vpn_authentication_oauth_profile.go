package provider

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func NewMsgVpnAuthenticationOauthProfileResource() resource.Resource {
	return &solaceResource[MsgVpnAuthenticationOauthProfile]{spr: &oauthProfileResource{}}
}

var _ solaceProviderResource[MsgVpnAuthenticationOauthProfile] = &oauthProfileResource{}

type oauthProfileResource struct {
	*solaceProvider
}

func (r oauthProfileResource) Name() string {
	return "oauthprofile"
}

func (r oauthProfileResource) Schema() schema.Schema {
	return MsgVpnAuthenticationOauthProfileResourceSchema("msg_vpn_name", "oauth_profile_name")
}

func (r *oauthProfileResource) SetProvider(provider *solaceProvider) {
	r.solaceProvider = provider
}

func (r oauthProfileResource) NewData() *MsgVpnAuthenticationOauthProfile {
	return &MsgVpnAuthenticationOauthProfile{}
}

func (r oauthProfileResource) Create(data *MsgVpnAuthenticationOauthProfile, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.MsgVpnApi.CreateMsgVpnAuthenticationOauthProfile(r.Context, *data.MsgVpnName).Body(*data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r oauthProfileResource) Read(data *MsgVpnAuthenticationOauthProfile, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.MsgVpnApi.GetMsgVpnAuthenticationOauthProfile(r.Context, *data.MsgVpnName, *data.OauthProfileName)
	apiResponse, httpResponse, err := apiReq.Execute()
	if err == nil && apiResponse != nil && apiResponse.Data != nil {
		data.ToTF(apiResponse.Data)
	}
	return httpResponse, err
}

func (r oauthProfileResource) Update(cur *MsgVpnAuthenticationOauthProfile, pln *MsgVpnAuthenticationOauthProfile, diag *diag.Diagnostics) (*http.Response, error) {
	_, httpResponse, err := r.Client.MsgVpnApi.
		UpdateMsgVpnAuthenticationOauthProfile(r.Context, *pln.MsgVpnName, *pln.OauthProfileName).
		Body(*pln.ToApi()).
		Execute()

	return httpResponse, err
}

func (r oauthProfileResource) Delete(data *MsgVpnAuthenticationOauthProfile, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.MsgVpnApi.DeleteMsgVpnAuthenticationOauthProfile(r.Context, *data.MsgVpnName, *data.OauthProfileName)
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

var oauthProfileImportRegexp *regexp.Regexp = regexp.MustCompile(fmt.Sprintf(
	"^([^\\s%s]+)\\/([^\\s%s]+)$", regexp.QuoteMeta("*?/"), regexp.QuoteMeta("'<>*?&;)")))

func (r oauthProfileResource) Import(id string, data *MsgVpnAuthenticationOauthProfile, diag *diag.Diagnostics) {
	match := oauthProfileImportRegexp.FindStringSubmatch(id)
	if match != nil {
		data.MsgVpnName = &match[1]
		data.OauthProfileName = &match[2]
	} else {
		diag.AddError("Expected <vpn-name>/<oauth-profile-name>", id+" does not match "+oauthProfileImportRegexp.String())
		return
	}
}
