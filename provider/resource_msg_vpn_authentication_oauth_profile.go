package provider

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ provider.ResourceType = queueResourceType{}

type oauthProfileResourceType struct {
}

func (t oauthProfileResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	solaceProvider, diags := convertProviderType(in)

	return NewResource[MsgVpnAuthenticationOauthProfile](
		oauthProfileResource{solaceProvider: solaceProvider}), diags
}

func (t oauthProfileResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return MsgVpnAuthenticationOauthProfileSchema("msg_vpn_name", "oauth_profile_name"), nil
}

var _ solaceProviderResource[MsgVpnAuthenticationOauthProfile] = oauthProfileResource{}

type oauthProfileResource struct {
	solaceProvider
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
