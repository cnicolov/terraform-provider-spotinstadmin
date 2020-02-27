package main

import (
	"github.com/cnicolov/terraform-provider-spotinstadmin/services/accounts"
	"github.com/cnicolov/terraform-provider-spotinstadmin/services/users"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider ...
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			providerTokenAttrKey: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(envSpotinstTokenKey, nil),
			},
			providerEmailAttrKey: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(envSpotinstEmailKey, nil),
			},
			providerPasswordAttrKey: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(envSpotinstPasswordKey, nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			accountResourceName:          resourceAccount(),
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
	username := d.Get(providerEmailAttrKey).(string)
	password := d.Get(providerPasswordAttrKey).(string)
	consoleToken, err := users.GetConsoleToken(username, password)

	if err != nil {
		return nil, err
	}

	return &Meta{
		accountsService: accounts.New(apiToken),
		usersService:    users.New(consoleToken),
	}, nil
}
