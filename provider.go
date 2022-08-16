package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kinvey/terraform-provider-spotinstadmin/services/accounts"
	"github.com/kinvey/terraform-provider-spotinstadmin/services/users"
)

// Provider ...
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			providerTokenAttrKey: {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(envSpotinstTokenKey, nil),
			},
			providerEmailAttrKey: {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(envSpotinstEmailKey, nil),
			},
			providerPasswordAttrKey: {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(envSpotinstPasswordKey, nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			accountResourceName:          resourceAccount(),
			linkAccountResourceName:      resourceAccountAWSLink(),
			programmaticUserResourceName: resourceProgrammaticUser(),
		},
		ConfigureFunc: providerConfigureFunc,
	}
}

// Meta ...
type Meta struct {
	accountsService *accounts.Service 
	usersService    *users.Service
}

func providerConfigureFunc(d *schema.ResourceData) (interface{}, error) {
	apiToken := d.Get(providerTokenAttrKey).(string)

	return &Meta{
		accountsService: accounts.New(apiToken),
		usersService:    users.New(apiToken),
	}, nil
}
