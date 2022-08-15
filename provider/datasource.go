package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type solaceProviderDataSource[Tdat any] interface {
	// NewData returns a new data struct for the datasource. This struct
	// needs to have `tfsdk:` tags for Terraform to reflect the config
	// into it.
	NewData() *Tdat

	// Read an existing resource
	Read(*Tdat, *diag.Diagnostics) (*http.Response, error)
}

var _ datasource.DataSource = dataSource[struct{}]{}

type dataSource[Tdat any] struct {
	spds solaceProviderDataSource[Tdat]
}

func NewDataSource[Tdat any](spds solaceProviderDataSource[Tdat]) *dataSource[Tdat] {
	return &dataSource[Tdat]{spds: spds}
}

func (r dataSource[Tdat]) DataFromCtx(ctx context.Context, config TFConfigGetter, diag *diag.Diagnostics) *Tdat {
	data := r.spds.NewData()
	diags := config.Get(ctx, data)
	diag.Append(diags...)

	return data
}

func (r dataSource[Tdat]) DataToCtx(ctx context.Context, data *Tdat, config TFConfigSetter, diag *diag.Diagnostics) {
	diags := config.Set(ctx, data)
	diag.Append(diags...)
}

func (ds dataSource[Tdat]) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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
