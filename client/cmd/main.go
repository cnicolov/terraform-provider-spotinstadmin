package main

import (
	"fmt"

	"github.com/cnicolov/terraform-provider-spotinstadmin/services/users"
)

func main() {

	t, err := users.GetConsoleToken("kristiyan.nikolov@progress.com", "#Plankton1")

	fmt.Printf("qq: %v", t)

	svc := users.New(t)

	user, err := svc.Create("blaqqqqqqq", "somedesc", "act-df0b0053")

	if err != nil {
		panic(err)
	}

	fmt.Println(user)
}
