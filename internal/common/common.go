package common

const (
	HeaderContentType = "Content-Type"
)

const (
	MIMEApplicationJSON            = "application/json"
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationForm            = "application/x-www-form-urlencoded"
	MIMEApplicationFormCharsetUTF8 = MIMEApplicationForm + "; " + charsetUTF8
	charsetUTF8                    = "charset=UTF-8"
)

const (
	ChannelGlobal        = "global"
	ChannelNotifications = "notifications"
)
