package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type solaceProviderDataSource[Tdat any] interface {
	// The name of this datasource, without the provider_ prefix
	Name() string

	// The schema for this datasource
	Schema() schema.Schema

	// To inject the provider, which holds the SEMPv2 client
	SetProvider(*solaceProvider)

	// NewData returns a new data struct for the datasource. This struct
	// needs to have `tfsdk:` tags for Terraform to reflect the config
	// into it.
	NewData() *Tdat

	// Read an existing resource
	Read(*Tdat, *diag.Diagnostics) (*http.Response, error)
}

var _ datasource.DataSourceWithConfigure = &dataSource[struct{}]{}

type dataSource[Tdat any] struct {
	spds solaceProviderDataSource[Tdat]
}

func (r *dataSource[Tdat]) DataFromCtx(ctx context.Context, config TFConfigGetter, diag *diag.Diagnostics) *Tdat {
	data := r.spds.NewData()
	diags := config.Get(ctx, data)
	diag.Append(diags...)

	return data
}

func (r *dataSource[Tdat]) DataToCtx(ctx context.Context, data *Tdat, config TFConfigSetter, diag *diag.Diagnostics) {
	diags := config.Set(ctx, data)
	diag.Append(diags...)
}

func (ds *dataSource[Tdat]) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "solace_" + ds.spds.Name()
}

func (ds *dataSource[Tdat]) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ds.spds.Schema()
}

func (r *dataSource[Tdat]) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if provider, ok := req.ProviderData.(*solaceProvider); ok {
		r.spds.SetProvider(provider)
	}
}

func (ds *dataSource[Tdat]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := ds.DataFromCtx(ctx, &req.Config, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := ds.spds.Read(data, &resp.Diagnostics)
	HandleSolaceApiError(err, &resp.Diagnostics, SempNotFound)
	if IsSempNotFound(err) {
		resp.State.RemoveResource(ctx)
	} else if !resp.Diagnostics.HasError() {
		ds.DataToCtx(ctx, data, &resp.State, &resp.Diagnostics)
	}
}
