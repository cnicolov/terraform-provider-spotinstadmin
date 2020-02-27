package main

import (
	"github.com/cnicolov/terraform-provider-spotinstadmin/services/accounts"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountCreate,
		Read:   resourceAccountRead,
		Update: resourceAccountUpdate,
		Delete: resourceAccountDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"aws_role_arn": {
				Type:        schema.TypeString,
				Description: "AWS Role arn to assume",
				Required:    true,
			},
			"aws_external_id": {
				Type:        schema.TypeString,
				Description: "ExternalID to use",
				Required:    true,
			},
		},
	}
}

func resourceAccountCreate(d *schema.ResourceData, m interface{}) error {
	accountsService := m.(*Meta).accountsService
	name := d.Get("name").(string)
	iamRole := d.Get("aws_role_arn").(string)
	externalID := d.Get("aws_external_id").(string)

	out, err := accountsService.Create(name, iamRole, externalID)

	if err != nil {
		return err
	}

	d.SetId(out.ID)

	d.Set("organization_id", out.OrganizationID)

	return resourceAccountRead(d, m)
}

func resourceAccountRead(d *schema.ResourceData, m interface{}) error {
	accountsService := m.(*Meta).accountsService
	obj, err := accountsService.Get(d.Id())
	if err != nil {
		if accounts.IsAccountNotFoundErr(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", obj.Name)
	d.Set("organization_id", obj.OrganizationID)
	d.Set("provider_external_id", obj.ProviderExternalID)

	return nil
}

func resourceAccountUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceAccountRead(d, m)
}

func resourceAccountDelete(d *schema.ResourceData, m interface{}) error {
	accountsService := m.(*Meta).accountsService
	return accountsService.Delete(d.Id())
}
