package matrix

type RoomMemberEventContent struct {
	Membership  string `json:"membership"`
	DisplayName string `json:"displayname"`
	AvatarUrl   string `json:"avatar_url,omitempty"`
}
