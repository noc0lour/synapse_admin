package main

import (
	// "flag"
	"fmt"
	sa_client "github.com/noc0lour/synapse_admin/pkg/client"
	sa_room "github.com/noc0lour/synapse_admin/pkg/room"
	sa_user "github.com/noc0lour/synapse_admin/pkg/user"
	gomatrix "maunium.net/go/mautrix"
)

func main() {
	// Parse commandline arguments
	// Setup admin access
	cli, err := gomatrix.NewClient("https://matrix.example.com", "", "")
	user := gomatrix.UserIdentifier{Type: "m.id.user", User: "@bar"}
	resp, err := cli.Login(&gomatrix.ReqLogin{
		Type:       "m.login.password",
		Identifier: user, Password: "foo",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	cli.SetCredentials(resp.UserID, resp.AccessToken)
	go cli.Sync()
	real_cli := sa_client.Client{Client: cli}
	user_cli := sa_user.Client{Client: &real_cli}
	users, err := user_cli.ListUsers()
	if err != nil {
		fmt.Println("ListUsers() returned ", err)
		return
	}
	fmt.Printf("%+v", users)

	room_cli := sa_room.Client{Client: &real_cli}
	rooms, err := room_cli.ListRooms()
	if err != nil {
		fmt.Println("ListRooms() returned ", err)
		return
	}
	fmt.Printf("%+v", rooms)

}
