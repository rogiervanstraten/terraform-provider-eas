package accountvariable

import (
	"strings"
	"terraform-provider-eas/internal/writeonly"
	"terraform-provider-eas/provider/accountvariable/operations"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages an account-level environment variable in EAS. Account variables are available across all apps within the account.",
		ReadContext:   operations.Read,
		CreateContext: operations.Create,
		UpdateContext: operations.Update,
		DeleteContext: operations.Delete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The id of the account environment variable",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the account environment variable",
				Type:        schema.TypeString,
				Required:    true,
			},
			"value": {
				Description: "The value of the account environment variable. Prefer 'value_wo' for ephemeral secrets on Terraform 1.11+.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"value_wo": {
				Description: "Write-only value for the account environment variable. Not stored in state. Requires Terraform 1.11+. Use this with ephemeral resources for secrets.",
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
				Description:  "The visibility of the account environment variable (PUBLIC, SENSITIVE, or SECRET)",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"PUBLIC", "SENSITIVE", "SECRET"}, false),
			},
			"environments": {
				Description: "The environments where the variable is available (development, preview, production)",
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
			"created_at": {
				Description: "The timestamp when the variable was created",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "The timestamp when the variable was last updated",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		ValidateRawResourceConfigFuncs: []schema.ValidateRawResourceConfigFunc{
			writeonly.PreferWriteOnlyForSensitive,
			writeonly.ValidateValueOrValueWo,
			writeonly.ValidateValueWoVersion,
		},
	}
}
