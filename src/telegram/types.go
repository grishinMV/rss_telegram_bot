package telegram

type UpdateResponse struct {
	Success bool     `json:"ok"`
	Result  []Update `json:"result"`
}

type SendMessageResponse struct {
	Success bool    `json:"ok"`
	Result  Message `json:"result"`
}

type User struct {
	Id                      int    `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name,omitempty"`
	Username                string `json:"username,omitempty"`
	IsPremium               bool   `json:"is_premium,omitempty"`
	AddedToAttachmentMenu   int    `json:"added_to_attachment_menu,omitempty"`
	CanJoinGroups           int    `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages int    `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries   int    `json:"supports_inline_queries,omitempty"`
}

type ChatPhoto struct {
	SmallFileId       string `json:"small_file_id"`
	SmallFileUniqueId string `json:"small_file_unique_id"`
	BigFileId         string `json:"big_file_id"`
	BigFileUniqueId   string `json:"big_file_unique_id"`
}

type Chat struct {
	Id                                 int       `json:"id"`
	Type                               string    `json:"type"`
	Title                              string    `json:"title,omitempty"`
	Username                           string    `json:"username,omitempty"`
	FirstName                          string    `json:"first_name,omitempty"`
	LastName                           string    `json:"last_name,omitempty"`
	IsForum                            bool      `json:"is_forum,omitempty"`
	Photo                              ChatPhoto `json:"photo,omitempty"`
	ActiveUsernames                    []string  `json:"active_usernames,omitempty"`
	EmojiStatusCustomEmojiId           string    `json:"emoji_status_custom_emoji_id,omitempty"`
	Bio                                string    `json:"bio,omitempty"`
	HasPrivateForwards                 bool      `json:"has_private_forwards,omitempty"`
	HasRestrictedVoiceAndVideoMessages bool      `json:"has_restricted_voice_and_video_messages,omitempty"`
	JoinToSendMessages                 bool      `json:"join_to_send_messages,omitempty"`
	JoinByRequest                      bool      `json:"join_by_request,omitempty"`
	Description                        string    `json:"description,omitempty"`
	InviteLink                         string    `json:"invite_link,omitempty"`
}

type Message struct {
	Id              int    `json:"message_id"`
	MessageThreadId int    `json:"message_thread_id,omitempty"`
	From            User   `json:"from,omitempty"`
	SenderChat      User   `json:"sender_chat,omitempty"`
	Date            int    `json:"date,omitempty"`
	Chat            Chat   `json:"chat,omitempty"`
	Text            string `json:"text,omitempty"`
}

type Update struct {
	UpdateID          int      `json:"update_id"`
	Message           *Message `json:"message,omitempty"`
	EditedMessage     *Message `json:"edited_message,omitempty"`
	ChannelPost       *Message `json:"channel_post,omitempty"`
	EditedChannelPost *Message `json:"edited_channel_post,omitempty"`
}
