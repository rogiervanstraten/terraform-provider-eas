package operations

import (
	"context"
	"terraform-provider-eas/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Delete removes an account-level environment variable.
//
// Required SDK signature:
//
//	func (c *AccountVariableClient) Delete(id string) (*AccountVariable, error)
func Delete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*client.EASClient)
	id := d.Id()

	_, err := client.AccountVariable.Delete(id)

	var diags diag.Diagnostics
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
