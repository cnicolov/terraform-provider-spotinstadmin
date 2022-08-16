package main

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceProgrammaticUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			userResourceAccountIDAttrKey: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			userResourceNameAttrKey: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			userResourceDescriptionAttrKey: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			userResourceAccessTokenAttrKey: {
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

	accountID := d.Get(userResourceAccountIDAttrKey).(string)
	log.Println(accountID)
	userName := d.Id()

	log.Printf("[TRACE] IN_RESOURCE_READ: %v-%v\n", userName, accountID)
	obj, err := usersService.Get(userName)

	if err != nil {
		d.SetId("")
		return err
	}

	if obj == nil {
		d.SetId("")
		return nil
	}
	log.Println(obj)

	actualName := strings.ToLower(obj.UserName)

	if actualName == d.Id() {
		d.SetId(actualName)
		return d.Set(userResourceAccountIDAttrKey, accountID)
	}

	d.SetId("")
	return nil
}

// func resourceProgrammaticUserUpdate(d *schema.ResourceData, m interface{}) error {
// }
func resourceProgrammaticUserCreate(d *schema.ResourceData, m interface{}) error {
	usersService := m.(*Meta).usersService

	username := d.Get(userResourceNameAttrKey).(string)
	description := d.Get(userResourceDescriptionAttrKey).(string)
	accountID := d.Get(userResourceAccountIDAttrKey).(string)

	log.Printf("[TRACE] IN_RESOURCE_CREATE: %v\n", accountID)
	user, err := usersService.Create(username, description, accountID)

	if err != nil {
		return err
	}

	d.SetId(strings.ToLower(user.UserName))
	if err := d.Set(userResourceAccessTokenAttrKey, user.AccessToken); err != nil {
		return err
	}

	if err := d.Set(userResourceAccountIDAttrKey, accountID); err != nil {
		return err
	}

	return resourceProgrammaticUserRead(d, m)
}

func resourceProgrammaticUserDelete(d *schema.ResourceData, m interface{}) error {
	usersService := m.(*Meta).usersService
	username := d.Get(userResourceNameAttrKey).(string)
	return usersService.Delete(username)
}
