package main

import (
	"github.com/cnicolov/terraform-provider-spotinstadmin/services/accounts"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider ...
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SPOTINST_TOKEN", nil),
			},
			"email": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SPOTINST_EMAIL", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SPOTINST_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"spotinstadmin_account": resourceAccount(),
			//			"spotinstadmin_programmatic_user": resourceProgrammaticUser(),
		},
		ConfigureFunc: providerConfigureFunc,
	}
}

// Meta ...
type Meta struct {
	accountsService *accounts.Service
	// UserService *UserService
}

func providerConfigureFunc(d *schema.ResourceData) (interface{}, error) {
	apiToken := d.Get("token").(string)
	// username := d.Get("email").(string)
	// 	password := d.Get("password").(string)
	// consoleToken, err := users.GetConsoleToken(username, password)

	// log.Println(consoleToken)

	// if err != nil {
	// 	return nil, err
	// }

	return &Meta{
		accountsService: accounts.New(apiToken),
		// UserService:    client.NewUserService(consoleToken),
	}, nil
}
