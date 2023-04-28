package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type solaceProviderResource[Tdat any] interface {
	// The name of this resource, without the provider_ prefix
	Name() string

	// The schema for this resource
	Schema() schema.Schema

	// To inject the provider, which holds the SEMPv2 client
	SetProvider(*solaceProvider)

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
	Import(string, *Tdat, *diag.Diagnostics)
}

var _ resource.ResourceWithConfigure = &solaceResource[struct{}]{}

type solaceResource[Tdat any] struct {
	spr solaceProviderResource[Tdat]
}

func (r *solaceResource[Tdat]) DataFromCtx(ctx context.Context, config TFConfigGetter, diag *diag.Diagnostics) *Tdat {
	data := r.spr.NewData()
	diags := config.Get(ctx, data)
	diag.Append(diags...)

	return data
}

func (r *solaceResource[Tdat]) DataToCtx(ctx context.Context, data *Tdat, config TFConfigSetter, diag *diag.Diagnostics) {
	diags := config.Set(ctx, data)
	diag.Append(diags...)
}

func (r *solaceResource[Tdat]) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "solace_" + r.spr.Name()
}

func (r *solaceResource[Tdat]) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.spr.Schema()
}

func (r *solaceResource[Tdat]) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if provider, ok := req.ProviderData.(*solaceProvider); ok {
		r.spr.SetProvider(provider)
	}
}

func (r *solaceResource[Tdat]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

func (r *solaceResource[Tdat]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

func (r *solaceResource[Tdat]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

func (r *solaceResource[Tdat]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

func (r *solaceResource[Tdat]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	data := r.spr.NewData()
	r.spr.Import(req.ID, data, &resp.Diagnostics)

	r.DataToCtx(ctx, data, &resp.State, &resp.Diagnostics)
}
