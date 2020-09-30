package main

import (
	// "flag"
	"fmt"
	"github.com/alecthomas/kong"
	sa_client "github.com/noc0lour/synapse_admin/pkg/client"
	sa_room "github.com/noc0lour/synapse_admin/pkg/room"
	sa_user "github.com/noc0lour/synapse_admin/pkg/user"
	gomatrix "maunium.net/go/mautrix"
	"strings"
)

var args struct {
	Loginuser     string `help:"User name for authentication" short:u`
	Loginpassword string `help:"User password for authentication" short:w`
	Server        string `help:"Matrix home server url" short:s`
	User          struct {
		List struct{} `cmd`
	} `cmd help:"User administration"`
	Room struct {
		List struct{} `cmd`
	} `cmd help:"Room administration"`
}

func main() {
	// Parse commandline arguments
	ctx := kong.Parse(&args)
	// Setup admin access
	cli, err := gomatrix.NewClient(args.Server, "", "")
	user := gomatrix.UserIdentifier{Type: "m.id.user", User: args.Loginuser}
	resp, err := cli.Login(&gomatrix.ReqLogin{
		Type:       "m.login.password",
		Identifier: user, Password: args.Loginpassword,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	// Refetch token if needed
	cli.SetCredentials(resp.UserID, resp.AccessToken)
	go cli.Sync()

	// Initialize base client
	real_cli := sa_client.Client{Client: cli}

	cmd := ctx.Command()

	// Branch for user or room subcommand
	if strings.HasPrefix(ctx.Command(), "user") {
		user_cmd := cmd[strings.IndexByte(cmd, ' ')+1:]
		user_cli := sa_user.Client{Client: &real_cli}
		switch user_cmd {
		case "list":
			users, err := user_cli.ListUsers()
			if err != nil {
				fmt.Println("ListUsers() returned ", err)
				return
			}
			for i, u := range users {
				fmt.Printf("%v,%v,%v,%v,%v,%v,%v\n", i, u.Name, u.IsGuest, u.Admin, u.UserType, u.Deactivated, u.DisplayName)
			}

		default:
			panic(user_cmd)
		}

	} else if strings.HasPrefix(ctx.Command(), "room") {
		room_cmd := cmd[strings.IndexByte(cmd, ' ')+1:]
		room_cli := sa_room.Client{Client: &real_cli}
		switch room_cmd {
		case "list":
			rooms, err := room_cli.ListRooms()
			if err != nil {
				fmt.Println("ListRooms() returned ", err)
				return
			}
			for i, r := range rooms {
					fmt.Printf("%v,%v,%v,%v,%v,%v,%v\n", i, r.Id, r.Name, r.Alias, r.JoinedMembers, r.JoinedLocalMembers, r.Creator)
				// fmt.Printf("%+v\n", r)
			}

		}
	} else {
		panic("Panic!")
	}

}
