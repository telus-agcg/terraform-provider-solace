package provider

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TFConfigGetter interface {
	Get(ctx context.Context, target interface{}) diag.Diagnostics
}

type TFConfigSetter interface {
	Set(ctx context.Context, target interface{}) diag.Diagnostics
}

// HasChanged is intended to use in a series of comparisons
// to see whether one or more of a set of attributes have changed.
// v1 is compared to v2 (a change from nil to non-nil or vice versa
// is considered a change) if and only if 'changed' is not yet set
// to true. If the input parameter 'changed' is already true, the
// comparison between v1 and v2 will not be made
func HasChanged[T comparable](v1 *T, v2 *T, changed *bool) {
	if *changed {
		return
	}

	if v1 == nil && v2 == nil {
		return
	} else if v1 == nil && v2 != nil {
		*changed = true
	} else if v1 != nil && v2 == nil {
		*changed = true
	} else if *v1 != *v2 {
		*changed = true
	}
}

const SempNotFound = "NOT_FOUND"

func HandleSolaceApiWithIgnore(err error, diag *diag.Diagnostics, ignoreStatus ...string) {
	HandleSolaceApiErrorF(err, diag.AddWarning, ignoreStatus)
}

func HandleSolaceApiError(err error, diag *diag.Diagnostics, ignoreStatus ...string) {
	HandleSolaceApiErrorF(err, diag.AddError, ignoreStatus)
}

func HandleSolaceApiErrorF(err error, f func(string, string), ignoreStatus []string) {
	if err == nil {
		return
	}

	if sempErr, ok := err.(*sempv2.GenericOpenAPIError); ok {
		if sempMeta, ok := sempErr.Model().(sempv2.SempMetaOnlyResponse); ok {
			if sempMeta.Meta.Error != nil {
				if !contains(ignoreStatus, sempMeta.Meta.Error.Status) {
					sempMetaErr := sempMeta.Meta.Error
					errStr := fmt.Sprintf("[%v] %v (%v)",
						sempMetaErr.Status, sempMetaErr.Description, sempMetaErr.Code)
					f("Solace API error", errStr)
				}
			} else {
				f("Solace API error", "SEMPv2 meta only response without error")
			}
		} else {
			f("Solace API error", fmt.Sprintf("SEMPv2 response without meta: %s", string(sempErr.Body())))
		}
	} else {
		f("Solace API error", fmt.Sprintf("Unexpected error type %q: %s", reflect.TypeOf(err), err.Error()))
	}
}

func IsSempNotFound(err error) bool {
	if sempErr, ok := err.(*sempv2.GenericOpenAPIError); ok {
		if sempMeta, ok := sempErr.Model().(sempv2.SempMetaOnlyResponse); ok {
			if sempMeta.Meta.Error != nil {
				if sempMeta.Meta.Error.Status == SempNotFound {
					return true
				}
			}
		}
	}
	return false
}

func toTFInt64(val *int64) types.Int64 {
	if val == nil {
		return types.Int64{Null: true}
	} else {
		return types.Int64{Value: *val}
	}
}

// WithRequiredAttributes marks the specified attributes as required and
// changing them will require replacement of the resource
func WithRequiredAttributes(schema tfsdk.Schema, names []string) tfsdk.Schema {
	for _, name := range names {
		attr, ok := schema.Attributes[name]
		if !ok {
			log.Panicf("WithRequiredAttributes: Attribute %q not found in schema %q", name, schema.Description)
		}
		attr.Required = true
		attr.Optional = false
		attr.PlanModifiers = append(attr.PlanModifiers, resource.RequiresReplace())
		schema.Attributes[name] = attr
	}
	return schema
}

// int64ToTF sets 'dst.Value' to 'val' if 'dst.Null' is false
// This will only update values in the Terraform state if those
// values have been specified by the user in the .tf file
func int64ToTF(dst *types.Int64, val *int64) {
	if !dst.Null {
		dst.Value = *val
	}
}

// int32ToTF sets 'dst.Value' to 'val' if 'dst.Null' is false
// This will only update values in the Terraform state if those
// values have been specified by the user in the .tf file
func int32ToTF(dst *types.Int64, val *int32) {
	if !dst.Null {
		dst.Value = int64(*val)
	}
}

// tfInt64ToPtr returns a pointer to the value if the value is not null
// or unknown. If the value is null or unknown, this returns nil
func tfInt64ToPtr(val *types.Int64) (res *int64) {
	if !val.Null && !val.Unknown {
		res = &val.Value
	}
	return
}

// tfInt32ToPtr returns a pointer to the value if the value is not null
// or unknown. If the value is null or unknown, this returns nil
func tfInt32ToPtr(val *types.Int64) (res *int32) {
	if !val.Null && !val.Unknown {
		i32 := int32(val.Value)
		res = &i32
	}
	return
}

func toTFBool(val *bool) types.Bool {
	if val == nil {
		return types.Bool{Null: true}
	} else {
		return types.Bool{Value: *val}
	}
}

// BoolToTF sets 'dst.Value' to 'val' if 'dst.Null' is false
// This will only update values in the Terraform state if those
// values have been specified by the user in the .tf file
func boolToTF(dst *types.Bool, val *bool) {
	if !dst.Null {
		dst.Value = *val
	}
}

// tfBoolToPtr returns a pointer to the value if the value is not null
// or unknown. If the value is null or unknown, this returns nil
func tfBoolToPtr(val *types.Bool) (res *bool) {
	if !val.Null && !val.Unknown {
		res = &val.Value
	}
	return
}

func toTFString(val *string) types.String {
	if val == nil {
		return types.String{Null: true}
	} else {
		return types.String{Value: *val}
	}
}

// stringToTF sets 'dst.Value' to 'val' if 'dst.Null' is false
// This will only update values in the Terraform state if those
// values have been specified by the user in the .tf file
func stringToTF(dst *types.String, val *string) {
	if val == nil {
		*dst = types.String{Null: true}
	} else if !dst.Null {
		*dst = types.String{Value: *val}
	}
}

// tfStringToPtr returns a pointer to the value if the value is not null
// or unknown. If the value is null or unknown, this returns nil
func tfStringToPtr(val *types.String) (res *string) {
	if !val.Null && !val.Unknown {
		res = &val.Value
	}
	return
}

func AssignIfDstNotNil[T any](dst **T, src *T) {
	if *dst != nil {
		*dst = src
	}
}

func contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func IsNilOrEmpty(str *string) bool {
	if str == nil {
		return true
	}
	if *str == "" {
		return true
	}
	if strings.Trim(*str, " ") == "" {
		return true
	}
	return false
}

func ToPtr[T any](val T) *T {
	return &val
}
