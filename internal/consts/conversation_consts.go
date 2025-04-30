package consts

type ConversationRole int16

const (
	CONVERSATION_OWNER   ConversationRole = 0
	CONVERSATION_ADMIN   ConversationRole = 1
	CONVERSATION_COADMIN ConversationRole = 2
	CONVERSATION_MEMBER  ConversationRole = 3
)

var ConversationRoles = []interface{}{
	CONVERSATION_OWNER,
	CONVERSATION_ADMIN,
	CONVERSATION_COADMIN,
	CONVERSATION_MEMBER,
}
