package main

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceProgrammaticUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			"access_token": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},

		Create: resourceProgrammaticUserCreate,
		Read:   resourceProgrammaticUserRead,
		//		Update: resourceProgrammaticUserUpdate,
		Delete: resourceProgrammaticUserDelete,
	}
}

func resourceProgrammaticUserRead(d *schema.ResourceData, m interface{}) error {
	usersService := m.(*Meta).usersService

	accountID := d.Get("account_id").(string)
	log.Println(accountID)
	userName := d.Id()

	log.Printf("IN_RESOURCE_READ: %v-%v\n", userName, accountID)
	obj, err := usersService.Get(userName, accountID)

	if err != nil {
		d.SetId("")
		return err
	}

	log.Println(obj)

	actualName := strings.ToLower(obj.CoreUser.FirstName)

	if actualName == d.Id() {
		d.SetId(actualName)
		return d.Set("account_id", obj.AccountID)
	}

	d.SetId("")
	return nil
}

// func resourceProgrammaticUserUpdate(d *schema.ResourceData, m interface{}) error {
// }
func resourceProgrammaticUserCreate(d *schema.ResourceData, m interface{}) error {
	usersService := m.(*Meta).usersService

	username := d.Get("name").(string)
	description := d.Get("description").(string)
	accountID := d.Get("account_id").(string)

	log.Printf("IN_RESOURCE_CREATE: %v\n", accountID)
	user, err := usersService.Create(username, description, accountID)

	if err != nil {
		return err
	}

	d.SetId(strings.ToLower(user.CoreUser.FirstName))
	if err := d.Set("access_token", user.AccessToken); err != nil {
		return err
	}

	if err := d.Set("account_id", accountID); err != nil {
		return err
	}

	return resourceProgrammaticUserRead(d, m)
}

func resourceProgrammaticUserDelete(d *schema.ResourceData, m interface{}) error {
	usersService := m.(*Meta).usersService
	username := d.Get("name").(string)
	accountID := d.Get("account_id").(string)
	return usersService.Delete(username, accountID)
}
