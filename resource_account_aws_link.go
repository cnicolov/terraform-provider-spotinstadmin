package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kinvey/terraform-provider-spotinstadmin/services/accounts"
)

func resourceAccountAWSLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountAWSLinkCreate,
		Read:   resourceAccountAWSLinkRead,
		Update: resourceAccountAWSLinkUpdate,
		Delete: resourceAccountAWSLinkDelete,

		Schema: map[string]*schema.Schema{
			linkAccountResourceAccountIDAttrKey: {
				Type:        schema.TypeString,
				Description: "Account ID to link",
				Required:    true,
			},

			linkAccountResourceRoleArnAttrKey: {
				Type:        schema.TypeString,
				Description: "AWS Role arn to assume",
				Required:    true,
			},
			linkAccountResourceProviderExternalIdAttrKey: {
				Type:        schema.TypeString,
				Description: "AWS Provider External ID",
				Computed:    true,
			},
			linkAccountResourceOrganizationIdAttrKey: {
				Type:        schema.TypeString,
				Description: "Organization ID",
				Computed:    true,
			},
		},
	}
}

func resourceAccountAWSLinkCreate(d *schema.ResourceData, m interface{}) error {
	accountsService := m.(*Meta).accountsService
	accountID := d.Get(linkAccountResourceAccountIDAttrKey).(string)
	iamRole := d.Get(linkAccountResourceRoleArnAttrKey).(string)

	err := accountsService.Link(accountID, iamRole)

	if err != nil {
		return err
	}

	return resourceAccountAWSLinkRead(d, m)
}

func resourceAccountAWSLinkRead(d *schema.ResourceData, m interface{}) error {
	accountsService := m.(*Meta).accountsService
	obj, err := accountsService.Get(d.Id())
	if err != nil {
		if accounts.IsAccountNotFoundErr(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	
	d.SetId(obj.ProviderExternalID)
	d.Set(accountResourceNameAttrKey, obj.Name)
	d.Set("organization_id", obj.OrganizationID)
	d.Set(linkAccountResourceProviderExternalIdAttrKey, obj.ProviderExternalID)

	return nil
}

func resourceAccountAWSLinkUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceAccountAWSLinkRead(d, m)
}

func resourceAccountAWSLinkDelete(d *schema.ResourceData, m interface{}) error {
	// accountsService := m.(*Meta).accountsService
	// return accountsService.Delete(d.Id())
	return nil
}
