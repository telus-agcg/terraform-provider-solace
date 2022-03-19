package util

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type stringOneOfValidatorData struct {
	values []string
}

func StringOneOfValidator(values ...string) stringOneOfValidatorData {
	return stringOneOfValidatorData{values: values}
}

func (v stringOneOfValidatorData) Description(ctx context.Context) string {
	return fmt.Sprintf("string must of one of %v", v.values)
}

func (v stringOneOfValidatorData) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v stringOneOfValidatorData) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	// types.String must be the attr.Value produced by the attr.Type in the schema for this attribute
	// for generic validators, use
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
	// to convert into a known type.
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if str.Unknown || str.Null {
		return
	}

	contains := func(s []string, e string) bool {
		for _, a := range s {
			if a == e {
				return true
			}
		}
		return false
	}

	if !contains(v.values, str.Value) {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Value",
			fmt.Sprintf("%s is not one of %v.", str.Value, v.values),
		)

		return
	}
}
