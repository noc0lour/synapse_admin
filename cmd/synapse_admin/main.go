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
	Null          bool   `help:"Separate results with zero-byte" short:0`
	User          struct {
		List       struct{} `cmd help:"List local user accounts"`
		Deactivate struct {
			UserId string `help: "UserID"`
		} `cmd help:"User deactivation API"`
		Whois struct {
			UserId string `arg help: "UserID"`
		} `cmd help: "Query whois information of a single user"`
		ListLastSeen struct {
			Since int `help: Unix time since when to check for Users`
		} `cmd help:"List local user accounts"`
	} `cmd help:"User administration"`
	Room struct {
		List       struct{} `cmd help: "List rooms"`
		MemberList struct {
			RoomId string `help: "RoomID"`
		} `cmd help: List room members`
	} `cmd help:"Room administration"`
}

func main() {
	// Parse commandline arguments
	ctx := kong.Parse(&args)
	arg_split := "\n"
	if args.Null {
		arg_split = "\000"
	}
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
				w, _ := user_cli.WhoisUser(u.Name)
				LastSeen := sa_user.GetLastSeen(w)
				fmt.Printf("%v,%v,%v,%v,%v,%v,%v,%v,%s", i, u.Name, u.IsGuest, u.Admin, u.UserType, u.Deactivated, u.DisplayName, LastSeen, arg_split)
			}
		case "deactivate":
			fmt.Printf("Deactivate() called for user %s\n", args.User.Deactivate.UserId)
			err := user_cli.DeactivateUser(args.User.Deactivate.UserId)
			if err != nil {
				fmt.Printf("Deactivate() returned %v", err)
			}
		case "whois <user-id>":
			whois, err := user_cli.WhoisUser(args.User.Whois.UserId)
			if err != nil {
				fmt.Printf("Whois() returned %v", err)
			}
			fmt.Printf("%+v\n", whois)
		case "list-last-seen":
			users, err := user_cli.ListLastSeen(args.User.ListLastSeen.Since)
			if err != nil {
				fmt.Printf("ListLastSeen() returned %v", err)
			}
			fmt.Printf("%v\n", strings.Join(users, "\n"))
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
			fmt.Printf("Room number, roomID, room name, room alias, number members, number local members, creator\n")
			for i, r := range rooms {
				fmt.Printf("%v,%v,%v,%v,%v,%v,%v%s", i, r.Id, r.Name, r.Alias, r.JoinedMembers, r.JoinedLocalMembers, r.Creator, arg_split)
				// fmt.Printf("%+v\n", r)
			}
		case "member-list":
			// fmt.Printf("RoomID: %v\n", args.Room.MemberList.RoomId)
			members, err := room_cli.ListRoomMembers(args.Room.MemberList.RoomId)
			if err != nil {
				fmt.Println("ListRoomMembers() returned ", err)
				return
			}
			for _, r := range members {
				fmt.Printf("%v%s", r, arg_split)
			}
		}
	} else {
		panic("Panic!")
	}

}
