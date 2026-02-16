package writeonly

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetValue(d *schema.ResourceData) string {
	valueWo, diags := d.GetRawConfigAt(cty.GetAttrPath("value_wo"))
	if !diags.HasError() && !valueWo.IsNull() && valueWo.IsKnown() {
		return valueWo.AsString()
	}

	if v, ok := d.GetOk("value"); ok {
		return v.(string)
	}

	return ""
}

func IsUsingWriteOnly(d *schema.ResourceData) bool {
	valueWo, _ := d.GetRawConfigAt(cty.GetAttrPath("value_wo"))
	return !valueWo.IsNull()
}

func ValidateValueOrValueWo(_ context.Context, req schema.ValidateResourceConfigFuncRequest, resp *schema.ValidateResourceConfigFuncResponse) {
	cfg := req.RawConfig

	value := cfg.GetAttr("value")
	valueWo := cfg.GetAttr("value_wo")

	valueIsNull := value.IsNull() || (value.IsKnown() && value.AsString() == "")
	valueWoIsNull := valueWo.IsNull()

	if valueIsNull && valueWoIsNull {
		resp.Diagnostics = append(resp.Diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing required value",
			Detail:   "Either 'value' or 'value_wo' must be provided.",
		})
	}
}

func ValidateValueWoVersion(_ context.Context, req schema.ValidateResourceConfigFuncRequest, resp *schema.ValidateResourceConfigFuncResponse) {
	cfg := req.RawConfig

	valueWo := cfg.GetAttr("value_wo")
	valueWoVersion := cfg.GetAttr("value_wo_version")

	if !valueWo.IsNull() && valueWoVersion.IsNull() {
		resp.Diagnostics = append(resp.Diagnostics, diag.Diagnostic{
			Severity:      diag.Warning,
			Summary:       "Missing value_wo_version",
			Detail:        "When using 'value_wo', it is recommended to also set 'value_wo_version' to trigger updates when the secret changes. Without it, Terraform cannot detect when the underlying secret value has changed.",
			AttributePath: cty.GetAttrPath("value_wo_version"),
		})
	}
}

func PreferWriteOnlyForSensitive(_ context.Context, req schema.ValidateResourceConfigFuncRequest, resp *schema.ValidateResourceConfigFuncResponse) {
	cfg := req.RawConfig

	value := cfg.GetAttr("value")
	valueWo := cfg.GetAttr("value_wo")
	visibility := cfg.GetAttr("visibility")

	if value.IsNull() || !valueWo.IsNull() {
		return
	}

	if !visibility.IsKnown() || visibility.IsNull() {
		return
	}

	visibilityStr := visibility.AsString()
	if visibilityStr != "SECRET" && visibilityStr != "SENSITIVE" {
		return
	}

	resp.Diagnostics = append(resp.Diagnostics, diag.Diagnostic{
		Severity:      diag.Warning,
		Summary:       "Use write-only attribute for sensitive values",
		Detail:        "For SECRET or SENSITIVE variables, consider using 'value_wo' instead of 'value'. Write-only attributes are not stored in state, providing better security for secrets. Requires Terraform 1.11+.",
		AttributePath: cty.GetAttrPath("value"),
	})
}
