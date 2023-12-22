package custom_types

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/apepenkov/wails_sigilix_interface/sigilix/crypto_utils"
)

type PublicUserInfo struct {
	UserId              uint64      `json:"user_id"`
	EcdsaPublicKey      Base64Bytes `json:"ecdsa_public_key"`
	Username            string      `json:"username"`
	InitialRsaPublicKey Base64Bytes `json:"initial_rsa_public_key"`
}

func (u *PublicUserInfo) ImplementSigilixStruct() {}

type PrivateUserInfo struct {
	PublicInfo              *PublicUserInfo `json:"public_info"`
	SearchByUsernameAllowed bool            `json:"search_by_username_allowed"`
}

func (u *PrivateUserInfo) ImplementSigilixStruct() {}

type LoginRequest struct {
	ClientEcdaPublicKey Base64Bytes `json:"client_ecdsa_public_key"`
	ClientRsaPublicKey  Base64Bytes `json:"client_rsa_public_key"`
}

func (u *LoginRequest) ImplementSigilixStruct() {}

type LoginResponse struct {
	PrivateInfo          *PrivateUserInfo `json:"private_info"`
	UserId               uint64           `json:"user_id"`
	ServerEcdsaPublicKey Base64Bytes      `json:"server_ecdsa_public_key"`
}

func (u *LoginResponse) ImplementSigilixStruct() {}

type SetUsernameConfigRequest struct {
	Username                string `json:"username"`
	SearchByUsernameAllowed bool   `json:"search_by_username_allowed"`
}

func (u *SetUsernameConfigRequest) ImplementSigilixStruct() {}

type SetUsernameConfigResponse struct {
	Success bool `json:"success"`
}

func (u *SetUsernameConfigResponse) ImplementSigilixStruct() {}

type SearchByUsernameRequest struct {
	Username string `json:"username"`
}

func (u *SearchByUsernameRequest) ImplementSigilixStruct() {}

type SearchByUsernameResponse struct {
	PublicInfo *PublicUserInfo `json:"public_info"`
}

func (u *SearchByUsernameResponse) ImplementSigilixStruct() {}

type InitChatFromInitializerRequest struct {
	TargetUserId uint64 `json:"target_user_id"`
}

func (u *InitChatFromInitializerRequest) ImplementSigilixStruct() {}

type InitChatFromInitializerResponse struct {
	ChatId uint64 `json:"chat_id"`
}

func (u *InitChatFromInitializerResponse) ImplementSigilixStruct() {}

type InitChatFromInitializerNotification struct {
	ChatId              uint64          `json:"chat_id"`
	InitializerUserInfo *PublicUserInfo `json:"initializer_user_info"`
}

func (i *InitChatFromInitializerNotification) ImplementSigilixStruct() {}

type InitChatFromReceiverRequest struct {
	ChatId uint64 `json:"chat_id"`
}

func (u *InitChatFromReceiverRequest) ImplementSigilixStruct() {}

type InitChatFromReceiverResponse struct {
	ChatId uint64 `json:"chat_id"`
}

func (u *InitChatFromReceiverResponse) ImplementSigilixStruct() {}

type InitChatFromReceiverNotification struct {
	ChatId           uint64          `json:"chat_id"`
	ReceiverUserInfo *PublicUserInfo `json:"receiver_user_info"`
}

func (i *InitChatFromReceiverNotification) ImplementSigilixStruct() {}

type UpdateChatRsaKeyRequest struct {
	ChatId       uint64      `json:"chat_id"`
	RsaPublicKey Base64Bytes `json:"rsa_public_key"`
}

func (u *UpdateChatRsaKeyRequest) ImplementSigilixStruct() {}

type UpdateChatRsaKeyResponse struct {
	ChatId uint64 `json:"chat_id"`
}

func (u *UpdateChatRsaKeyResponse) ImplementSigilixStruct() {}

type UpdateChatRsaKeyNotification struct {
	ChatId       uint64      `json:"chat_id"`
	UserId       uint64      `json:"user_id"`
	RsaPublicKey Base64Bytes `json:"rsa_public_key"`
}

func (u *UpdateChatRsaKeyNotification) ImplementSigilixStruct() {}

type SendMessageRequest struct {
	ChatId                uint64      `json:"chat_id"`
	EncryptedMessage      Base64Bytes `json:"encrypted_message"`
	MessageEcdsaSignature Base64Bytes `json:"message_ecdsa_signature"`
}

func (u *SendMessageRequest) ImplementSigilixStruct() {}

type SendMessageResponse struct {
	ChatId    uint64 `json:"chat_id"`
	MessageId uint64 `json:"message_id"`
}

func (u *SendMessageResponse) ImplementSigilixStruct() {}

type SendMessageNotification struct {
	ChatId                uint64      `json:"chat_id"`
	MessageId             uint64      `json:"message_id"`
	SenderUserId          uint64      `json:"sender_user_id"`
	EncryptedMessage      Base64Bytes `json:"encrypted_message"`
	MessageEcdsaSignature Base64Bytes `json:"message_ecdsa_signature"`
}

func (s *SendMessageNotification) ValidateAndDecrypt(ecdsaPublicKey *ecdsa.PublicKey, rsaPrivateKey *rsa.PrivateKey) ([]byte, error) {
	// Verify the signature
	// Decrypt the message

	decryptedMessage, err := crypto_utils.DecryptMessage(rsaPrivateKey, s.EncryptedMessage)
	if err != nil {
		return nil, err
	}

	// Verify the signature
	ok, err := crypto_utils.ValidateECDSASignature(ecdsaPublicKey, decryptedMessage, s.MessageEcdsaSignature)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("invalid signature")
	}

	return decryptedMessage, nil
}

func (s *SendMessageNotification) ImplementSigilixStruct() {}

type SendFileRequest struct {
	ChatId             uint64      `json:"chat_id"`
	EncryptedFile      Base64Bytes `json:"encrypted_file"`
	EncryptedMimeType  Base64Bytes `json:"encrypted_mime_type"`
	FileEcdsaSignature Base64Bytes `json:"file_ecdsa_signature"`
}

func (u *SendFileRequest) ImplementSigilixStruct() {}

type SendFileResponse struct {
	ChatId    uint64 `json:"chat_id"`
	MessageId uint64 `json:"message_id"`
}

func (u *SendFileResponse) ImplementSigilixStruct() {}

type SendFileNotification struct {
	ChatId             uint64      `json:"chat_id"`
	MessageId          uint64      `json:"message_id"`
	SenderUserId       uint64      `json:"sender_user_id"`
	EncryptedFile      Base64Bytes `json:"encrypted_file"`
	EncryptedMimeType  Base64Bytes `json:"encrypted_mime_type"`
	FileEcdsaSignature Base64Bytes `json:"file_ecdsa_signature"`
}

func (s *SendFileNotification) ImplementSigilixStruct() {}

type IncomingNotification struct {
	Notification   SomeNotification `json:"notification"`
	EcdsaSignature Base64Bytes      `json:"ecdsa_signature"`
}

func (i *IncomingNotification) ImplementSigilixStruct() {}

type NotificationWithTypeInfo struct {
	Notification SomeNotification `json:"notification"`
	Type         NotificationType `json:"type"`
}

func (u *NotificationWithTypeInfo) ImplementSigilixStruct() {}

func (i *IncomingNotification) MarshalJSON() ([]byte, error) {
	r := &NotificationWithTypeInfo{
		Notification: i.Notification,
		Type:         i.Notification.NotificationType(),
	}
	return json.Marshal(r)
}

func (i *IncomingNotification) UnmarshalJSON(data []byte) error {
	//r := &NotificationWithTypeInfo{}
	//err := json.Unmarshal(data, r)
	//if err != nil {
	//	return err
	//}
	//i.Notification = r.Notification
	//return nil
	var r map[string]interface{}
	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}
	// check the type
	notificationType, ok := r["type"].(string)
	if !ok {
		return fmt.Errorf("invalid type")
	}
	nType := NotificationType(notificationType)

	switch nType {
	case notificationTypeInitChatFromInitializer:
		i.Notification = &InitChatFromInitializerNotification{}
	case notificationTypeInitChatFromReceiver:
		i.Notification = &InitChatFromReceiverNotification{}
	case notificationTypeUpdateChatRsaKey:
		i.Notification = &UpdateChatRsaKeyNotification{}
	case notificationTypeSendMessage:
		i.Notification = &SendMessageNotification{}
	case notificationTypeSendFile:
		i.Notification = &SendFileNotification{}
	default:
		return fmt.Errorf("invalid type")
	}

	toReMarshal, ok := r["notification"].(map[string]interface{})

	if !ok {
		return fmt.Errorf("invalid notification")
	}

	toReMarshalBytes, err := json.Marshal(toReMarshal)
	if err != nil {
		return err
	}

	err = json.Unmarshal(toReMarshalBytes, i.Notification)
	if err != nil {
		return err
	}
	return nil
}

type GetNotificationsRequest struct {
	Limit uint32 `json:"limit"`
}

func (u *GetNotificationsRequest) ImplementSigilixStruct() {}

type GetNotificationsResponse struct {
	Notifications []*IncomingNotification `json:"notifications"`
	//Notifications []*NotificationWithTypeInfo `json:"notifications"`
}

func (u *GetNotificationsResponse) ImplementSigilixStruct() {}

func (b *Base64Bytes) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", base64.StdEncoding.EncodeToString(*b))), nil
}

func (b *Base64Bytes) UnmarshalJSON(data []byte) error {
	// Remove the quotes from the JSON string
	str := string(data[1 : len(data)-1])

	// Decode from base64
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}

	*b = decoded
	return nil
}
