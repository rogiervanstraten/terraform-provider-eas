package writeonly

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestValidateValueOrValueWo_BothEmpty(t *testing.T) {
	req := schema.ValidateResourceConfigFuncRequest{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"value":    cty.NullVal(cty.String),
			"value_wo": cty.NullVal(cty.String),
		}),
	}
	resp := &schema.ValidateResourceConfigFuncResponse{}

	ValidateValueOrValueWo(context.Background(), req, resp)

	if len(resp.Diagnostics) != 1 {
		t.Errorf("expected 1 diagnostic, got %d", len(resp.Diagnostics))
	}
}

func TestValidateValueOrValueWo_ValueSet(t *testing.T) {
	req := schema.ValidateResourceConfigFuncRequest{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"value":    cty.StringVal("test"),
			"value_wo": cty.NullVal(cty.String),
		}),
	}
	resp := &schema.ValidateResourceConfigFuncResponse{}

	ValidateValueOrValueWo(context.Background(), req, resp)

	if len(resp.Diagnostics) != 0 {
		t.Errorf("expected 0 diagnostics, got %d", len(resp.Diagnostics))
	}
}

func TestValidateValueOrValueWo_ValueWoSet(t *testing.T) {
	req := schema.ValidateResourceConfigFuncRequest{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"value":    cty.NullVal(cty.String),
			"value_wo": cty.StringVal("secret"),
		}),
	}
	resp := &schema.ValidateResourceConfigFuncResponse{}

	ValidateValueOrValueWo(context.Background(), req, resp)

	if len(resp.Diagnostics) != 0 {
		t.Errorf("expected 0 diagnostics, got %d", len(resp.Diagnostics))
	}
}

func TestValidateValueWoVersion_MissingVersion(t *testing.T) {
	req := schema.ValidateResourceConfigFuncRequest{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"value_wo":         cty.StringVal("secret"),
			"value_wo_version": cty.NullVal(cty.Number),
		}),
	}
	resp := &schema.ValidateResourceConfigFuncResponse{}

	ValidateValueWoVersion(context.Background(), req, resp)

	if len(resp.Diagnostics) != 1 {
		t.Errorf("expected 1 warning diagnostic, got %d", len(resp.Diagnostics))
	}
}

func TestValidateValueWoVersion_VersionSet(t *testing.T) {
	req := schema.ValidateResourceConfigFuncRequest{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"value_wo":         cty.StringVal("secret"),
			"value_wo_version": cty.NumberIntVal(1),
		}),
	}
	resp := &schema.ValidateResourceConfigFuncResponse{}

	ValidateValueWoVersion(context.Background(), req, resp)

	if len(resp.Diagnostics) != 0 {
		t.Errorf("expected 0 diagnostics, got %d", len(resp.Diagnostics))
	}
}

func TestPreferWriteOnlyForSensitive_SecretVisibility(t *testing.T) {
	req := schema.ValidateResourceConfigFuncRequest{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"value":      cty.StringVal("secret-value"),
			"value_wo":   cty.NullVal(cty.String),
			"visibility": cty.StringVal("SECRET"),
		}),
	}
	resp := &schema.ValidateResourceConfigFuncResponse{}

	PreferWriteOnlyForSensitive(context.Background(), req, resp)

	if len(resp.Diagnostics) != 1 {
		t.Errorf("expected 1 warning for SECRET visibility, got %d", len(resp.Diagnostics))
	}
}

func TestPreferWriteOnlyForSensitive_SensitiveVisibility(t *testing.T) {
	req := schema.ValidateResourceConfigFuncRequest{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"value":      cty.StringVal("sensitive-value"),
			"value_wo":   cty.NullVal(cty.String),
			"visibility": cty.StringVal("SENSITIVE"),
		}),
	}
	resp := &schema.ValidateResourceConfigFuncResponse{}

	PreferWriteOnlyForSensitive(context.Background(), req, resp)

	if len(resp.Diagnostics) != 1 {
		t.Errorf("expected 1 warning for SENSITIVE visibility, got %d", len(resp.Diagnostics))
	}
}

func TestPreferWriteOnlyForSensitive_PublicVisibility(t *testing.T) {
	req := schema.ValidateResourceConfigFuncRequest{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"value":      cty.StringVal("public-value"),
			"value_wo":   cty.NullVal(cty.String),
			"visibility": cty.StringVal("PUBLIC"),
		}),
	}
	resp := &schema.ValidateResourceConfigFuncResponse{}

	PreferWriteOnlyForSensitive(context.Background(), req, resp)

	if len(resp.Diagnostics) != 0 {
		t.Errorf("expected 0 warnings for PUBLIC visibility, got %d", len(resp.Diagnostics))
	}
}

func TestPreferWriteOnlyForSensitive_AlreadyUsingWriteOnly(t *testing.T) {
	req := schema.ValidateResourceConfigFuncRequest{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"value":      cty.NullVal(cty.String),
			"value_wo":   cty.StringVal("secret"),
			"visibility": cty.StringVal("SECRET"),
		}),
	}
	resp := &schema.ValidateResourceConfigFuncResponse{}

	PreferWriteOnlyForSensitive(context.Background(), req, resp)

	if len(resp.Diagnostics) != 0 {
		t.Errorf("expected 0 warnings when already using value_wo, got %d", len(resp.Diagnostics))
	}
}
