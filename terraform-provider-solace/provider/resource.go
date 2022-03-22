package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type solaceProviderResource[Tdat any] interface {
	// NewData returns a new data struct for the resource. This struct
	// needs to have `tfsdk:` tags for Terraform to reflect the config
	// into it.
	NewData() *Tdat

	// Create a new resource
	Create(*Tdat, *diag.Diagnostics) (*http.Response, error)

	// Read an existing resource
	Read(*Tdat, *diag.Diagnostics) (*http.Response, error)

	// Update an existing resource
	Update(curState *Tdat, plnState *Tdat, diag *diag.Diagnostics) (*http.Response, error)

	// Delete a resource
	Delete(*Tdat, *diag.Diagnostics) (*http.Response, error)

	// Import a resource
	Import(*Tdat, *diag.Diagnostics)
}

var _ tfsdk.Resource = resource[struct{}]{}

type resource[Tdat any] struct {
	spr solaceProviderResource[Tdat]
}

func NewResource[Tdat any](spr solaceProviderResource[Tdat]) *resource[Tdat] {
	return &resource[Tdat]{spr: spr}
}

func (r resource[Tdat]) DataFromCtx(ctx context.Context, config TFConfigGetter, diag *diag.Diagnostics) *Tdat {
	data := r.spr.NewData()
	diags := config.Get(ctx, data)
	diag.Append(diags...)

	return data
}

func (r resource[Tdat]) DataToCtx(ctx context.Context, data *Tdat, config TFConfigSetter, diag *diag.Diagnostics) {
	diags := config.Set(ctx, data)
	diag.Append(diags...)
}

func (r resource[Tdat]) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	data := r.DataFromCtx(ctx, &req.Plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.spr.Create(data, &resp.Diagnostics)
	HandleSolaceApiError(err, &resp.Diagnostics)
	if !resp.Diagnostics.HasError() {
		r.DataToCtx(ctx, data, &resp.State, &resp.Diagnostics)
	}
}

func (r resource[Tdat]) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	data := r.DataFromCtx(ctx, &req.State, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.spr.Read(data, &resp.Diagnostics)
	HandleSolaceApiError(err, &resp.Diagnostics, SempNotFound)
	if IsSempNotFound(err) {
		resp.State.RemoveResource(ctx)
	} else if !resp.Diagnostics.HasError() {
		r.DataToCtx(ctx, data, &resp.State, &resp.Diagnostics)
	}

}

func (r resource[Tdat]) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	curState := r.DataFromCtx(ctx, &req.State, &resp.Diagnostics)
	plnState := r.DataFromCtx(ctx, &req.Plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.spr.Update(curState, plnState, &resp.Diagnostics)
	HandleSolaceApiError(err, &resp.Diagnostics)
	if !resp.Diagnostics.HasError() {
		// Update the data with the response from the API
		r.DataToCtx(ctx, plnState, &resp.State, &resp.Diagnostics)
	}
}

func (r resource[Tdat]) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	data := r.DataFromCtx(ctx, &req.State, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.spr.Delete(data, &resp.Diagnostics)
	HandleSolaceApiError(err, &resp.Diagnostics)
	if !resp.Diagnostics.HasError() {
		resp.State.RemoveResource(ctx)
	}
}

func (r resource[Tdat]) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: msg-vpn/name. Got: %q", req.ID),
		)
		return
	}

	data := r.spr.NewData()
	r.spr.Import(data, &resp.Diagnostics)

	r.DataToCtx(ctx, data, &resp.State, &resp.Diagnostics)
}
