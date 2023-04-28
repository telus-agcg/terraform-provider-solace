package provider

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
