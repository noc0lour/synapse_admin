package types

// https://github.com/matrix-org/synapse/blob/master/docs/admin_api/rooms.md
// Responses

type RoomList struct {
	Rooms  []Room `json:"rooms"`
	Offset int `json:"offset"`
	Total  int    `json:"total_rooms"`
}
type Room struct {
	Id                 string `json:"room_id"`
	Name               string `json:"name"`
	Alias              string `json:"canonical_alias"`
	JoinedMembers      int    `json:"joined_members"`
	JoinedLocalMembers int    `json:"joined_local_members"`
	Version            string `json:"version"`
	Creator            string `json:"creator"`
	Encryption         string `json:"encryption"`
	Federatable        bool    `json:"federatable"`
	Public             bool    `json:"public"`
	JoinRules          string `json:"join_rules"`
	GuestAccess        string `json:"GuestAccess"`
	HistoryVisibility  string `json:"HistoryVisibility"`
	StateEvents        int `json:"state_events"`
}
