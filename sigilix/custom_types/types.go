package custom_types

type NotificationType string

const (
	notificationTypeInitChatFromInitializer NotificationType = "InitChatFromInitializer"
	notificationTypeInitChatFromReceiver    NotificationType = "InitChatFromReceiver"
	notificationTypeUpdateChatRsaKey        NotificationType = "UpdateChatRsaKey"
	notificationTypeSendMessage             NotificationType = "SendMessage"
	notificationTypeSendFile                NotificationType = "SendFile"
)

type Base64Bytes []byte

type SomeNotification interface {
	NotificationType() NotificationType
}

type SigilixStruct interface {
	ImplementSigilixStruct()
}

func (i *InitChatFromInitializerNotification) NotificationType() NotificationType {
	return notificationTypeInitChatFromInitializer
}
func (i *InitChatFromReceiverNotification) NotificationType() NotificationType {
	return notificationTypeInitChatFromReceiver
}

func (u *UpdateChatRsaKeyNotification) NotificationType() NotificationType {
	return notificationTypeUpdateChatRsaKey
}

func (s *SendMessageNotification) NotificationType() NotificationType {
	return notificationTypeSendMessage
}

func (s *SendFileNotification) NotificationType() NotificationType { return notificationTypeSendFile }
