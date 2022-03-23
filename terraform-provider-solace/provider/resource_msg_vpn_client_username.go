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

var _ tfsdk.ResourceType = clientUsernameResourceType{}

type clientUsernameResourceType struct {
}

func (t clientUsernameResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return NewResource[MsgVpnClientUsername](
		clientUsernameResource{provider: provider}), diags
}

func (t clientUsernameResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return MsgVpnClientUsernameSchema("msg_vpn_name", "client_username"), nil
}

var _ solaceProviderResource[MsgVpnClientUsername] = clientUsernameResource{}

type clientUsernameResource struct {
	provider
}

func (r clientUsernameResource) NewData() *MsgVpnClientUsername {
	return &MsgVpnClientUsername{}
}

func (r clientUsernameResource) Create(data *MsgVpnClientUsername, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientUsernameApi.CreateMsgVpnClientUsername(r.Context, *data.MsgVpnName).Body(*data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r clientUsernameResource) Read(data *MsgVpnClientUsername, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientUsernameApi.GetMsgVpnClientUsername(r.Context, *data.MsgVpnName, *data.ClientUsername)
	apiResponse, httpResponse, err := apiReq.Execute()
	if err == nil && apiResponse != nil && apiResponse.Data != nil {
		data.ToTF(apiResponse.Data)
	}
	return httpResponse, err
}

func (r clientUsernameResource) Update(cur *MsgVpnClientUsername, pln *MsgVpnClientUsername, diag *diag.Diagnostics) (*http.Response, error) {
	requiresShutdown := false
	HasChanged(cur.ClientProfileName, pln.ClientProfileName, &requiresShutdown)
	HasChanged(cur.AclProfileName, pln.AclProfileName, &requiresShutdown)

	apiPlan := pln.ToApi()
	if requiresShutdown {
		apiPlan.Enabled = sempv2.PtrBool(false)
	}

	_, httpResponse, err := r.Client.ClientUsernameApi.
		UpdateMsgVpnClientUsername(r.Context, *pln.MsgVpnName, *pln.ClientUsername).
		Body(*apiPlan).
		Execute()

	// If the client-username needed shut down before *and* the desired state
	// is enabled, turn it back on
	if err == nil && requiresShutdown && pln.Enabled != nil && *pln.Enabled {
		body := sempv2.MsgVpnClientUsername{Enabled: sempv2.PtrBool(true)}
		_, httpResponse, err = r.Client.ClientUsernameApi.UpdateMsgVpnClientUsername(
			r.Context, *pln.MsgVpnName, *pln.ClientUsername).Body(body).Execute()
	}

	return httpResponse, err
}

func (r clientUsernameResource) Delete(data *MsgVpnClientUsername, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientUsernameApi.DeleteMsgVpnClientUsername(r.Context, *data.MsgVpnName, *data.ClientUsername)
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

var msgVpnClientUsernameImportRegexp *regexp.Regexp = regexp.MustCompile(fmt.Sprintf(
	"^([^\\s%s]+)\\/([^\\s%s]+)$", regexp.QuoteMeta("*?/"), regexp.QuoteMeta("*?")))

func (r clientUsernameResource) Import(id string, data *MsgVpnClientUsername, diag *diag.Diagnostics) {
	match := msgVpnClientUsernameImportRegexp.FindStringSubmatch(id)
	if match != nil {
		data.MsgVpnName = &match[1]
		data.ClientUsername = &match[2]
	} else {
		diag.AddError("Expected <vpn-name>/<client-username>", id+" does not match "+msgVpnClientUsernameImportRegexp.String())
		return
	}
}
