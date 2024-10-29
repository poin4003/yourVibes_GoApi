package consts

type FriendStatus string

const (
	NOT_FRIEND             FriendStatus = "not_friend"
	IS_FRIEND              FriendStatus = "is_friend"
	SEND_FRIEND_REQUEST    FriendStatus = "send_friend_request"
	RECEIVE_FRIEND_REQUEST FriendStatus = "receive_friend_request"
)
