package operations

import (
	"context"
	"terraform-provider-eas/internal/client"
	"terraform-provider-eas/internal/writeonly"

	"github.com/fintreal/eas-sdk-go/eas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Create(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*client.EASClient)

	name := d.Get("name").(string)
	value := writeonly.GetValue(d)
	visibility := d.Get("visibility").(string)

	set := d.Get("environments").(*schema.Set).List()

	var environments []string
	for _, v := range set {
		str := v.(string)
		environments = append(environments, str)
	}

	input := eas.CreateAccountVariableData{
		AccountId:    client.AccountId,
		Name:         name,
		Value:        value,
		Visibility:   visibility,
		Environments: environments,
	}

	data, err := client.AccountVariable.Create(input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(data.Id)

	var diags diag.Diagnostics

	if err := d.Set("created_at", data.CreatedAt); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("updated_at", data.UpdatedAt); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
