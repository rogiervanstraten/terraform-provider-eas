package appvariable

import (
	"strings"
	"terraform-provider-eas/internal/writeonly"
	"terraform-provider-eas/provider/appvariable/operations"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages an app-level environment variable in EAS.",
		ReadContext:   operations.Read,
		CreateContext: operations.Create,
		UpdateContext: operations.Update,
		DeleteContext: operations.Delete,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Description: "The id of the app for the environment variable",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of the app for the environment variable",
				Type:        schema.TypeString,
				Required:    true,
			},
			"id": {
				Description: "The id of the environment variable",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"value": {
				Description: "The value of the environment variable. Prefer 'value_wo' for ephemeral secrets on Terraform 1.11+.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"value_wo": {
				Description: "Write-only value for the environment variable. Not stored in state. Requires Terraform 1.11+. Use this with ephemeral resources for secrets.",
				Type:        schema.TypeString,
				Optional:    true,
				WriteOnly:   true,
			},
			"value_wo_version": {
				Description: "Version identifier for the write-only value. Changing this value triggers an update to write the new value_wo. Required when using value_wo (e.g., set to aws_secretsmanager_secret_version.example.version_id).",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"visibility": {
				Description:  "The visibility of the app for the environment variable",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLIC", "SENSITIVE", "SECRET"}, false),
			},
			"environments": {
				Description: "The environments of the app for the environment variable",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"development", "preview", "production"}, true),
					StateFunc: func(val any) string {
						return strings.ToLower(val.(string))
					},
				},
				Required: true,
			},
		},
		ValidateRawResourceConfigFuncs: []schema.ValidateRawResourceConfigFunc{
			writeonly.PreferWriteOnlyForSensitive,
			writeonly.ValidateValueOrValueWo,
			writeonly.ValidateValueWoVersion,
		},
	}
}
