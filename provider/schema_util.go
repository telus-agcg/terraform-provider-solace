package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func StringPlanModifiersFor(name string, requiredAttributes []string) (modifiers []planmodifier.String) {
	if contains(requiredAttributes, name) {
		modifiers = append(modifiers, stringplanmodifier.RequiresReplace())
	}
	return
}
