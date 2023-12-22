package messenger_client

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/apepenkov/wails_sigilix_interface/sigilix/crypto_utils"
	"github.com/apepenkov/wails_sigilix_interface/sigilix/custom_types"
	"github.com/apepenkov/wails_sigilix_interface/sigilix/data"
	"github.com/apepenkov/wails_sigilix_interface/sigilix/http_client"
	"log"
	"os"
	"strconv"
	"time"
)

const configFilename = "config.json"

type MessengerClient struct {
	config   *data.Config
	database *data.SqliteDB
	http     *http_client.SigilixHttpClient
	apiUrl   string
	unlocked bool
}

func NewClient(apiUrl string) *MessengerClient {
	return &MessengerClient{
		apiUrl: apiUrl,
	}
}

func (c *MessengerClient) connectSqlite(filename string) error {
	db, err := data.NewSqliteDB(filename)
	if err != nil {
		return err
	}
	c.database = db
	return nil
}

func (c *MessengerClient) connectSqlitePassword(filename string, passwordHash []byte) error {
	db, err := data.NewSqliteDBWithPassword(filename, hex.EncodeToString(passwordHash))
	if err != nil {
		return err
	}
	c.database = db
	return nil
}

func (c *MessengerClient) IsSignedUp() bool {
	_, err := os.Stat(configFilename)
	if err != nil {
		return false
	}
	return true
}

func (c *MessengerClient) IsUnlocked() bool {
	return c.unlocked
}

func Sha256x(x int, data []byte) []byte {
	for i := 0; i < x; i++ {
		hashed := sha256.Sum256(data)
		data = hashed[:]
	}
	return data
}

func (c *MessengerClient) SignUp(password string) error {
	passHash := Sha256x(100, []byte(password))

	conf := &data.Config{}
	ecdsaPrivate, err := crypto_utils.GenerateKey()
	if err != nil {
		return err
	}
	rsa, err := crypto_utils.NewRSAKeyPair()
	if err != nil {
		return err
	}

	conf.Username = ""
	conf.SearchByUsername = false
	conf.UserId = crypto_utils.GenerateUserIdByPublicKey(ecdsaPrivate.Public().(*ecdsa.PublicKey))
	conf.InitialRsaRivateKey = crypto_utils.RsaPrivateToBytes(rsa)
	conf.InitialECDSAPrivateKey = crypto_utils.PrivateKeyToBytes(ecdsaPrivate)
	conf.PaswordHash = passHash

	err = conf.SaveToFile(configFilename)
	if err != nil {
		return err
	}
	return nil
}

func (c *MessengerClient) Unlock(password string) error {
	if !c.IsSignedUp() {
		return errors.New("not signed up")
	}
	passHash := Sha256x(100, []byte(password))
	conf, err := data.LoadConfigFromFiles(configFilename, passHash)
	if err != nil {
		return err
	}
	c.config = conf
	c.http = http_client.NewSigilixHttpClient(c.apiUrl, conf.MustEcdsaPrivateKey(), conf.UserId)

	login, err := c.http.Login(conf.MustEcdsaPublicKey(), conf.MustRsaPublicKey())
	if err != nil {
		return err
	}
	if login.UserId != conf.UserId {
		return errors.New("wrong user id")
	}
	//err = c.connectSqlite(fmt.Sprintf("sigilix_%d.db", conf.UserId))
	err = c.connectSqlitePassword(fmt.Sprintf("sigilix_%d.db", conf.UserId), passHash)
	if err != nil {
		return err
	}
	c.unlocked = true
	return nil
}

func (c *MessengerClient) GetChats() ([]*data.Chat, error) {
	if !c.unlocked {
		return nil, errors.New("not unlocked")
	}
	return c.database.GetAllChats()
}

func (c *MessengerClient) GetChatMessages(chatId uint64) ([]*data.Message, error) {
	if !c.unlocked {
		return nil, errors.New("not unlocked")
	}
	return c.database.GetMessages(chatId)
}

func (c *MessengerClient) SendMessage(chatId uint64, text string) (*data.Message, error) {
	if !c.unlocked {
		return nil, errors.New("not unlocked")
	}
	chat, err := c.database.GetChat(chatId)

	if chat == nil {
		return nil, errors.New("chat not found")
	}

	if !chat.Accepted {
		return nil, errors.New("chat not accepted")
	}

	if err != nil {
		return nil, err
	}

	rsaPub, err := chat.OtherUserRsaPublicKey()
	if err != nil {
		return nil, err
	}

	msg, err := c.http.SendMessage(chatId, text, rsaPub)
	if err != nil {
		return nil, err
	}

	message := c.database.NewMessage()
	message.ChatId = chatId
	message.Content = text
	message.MessageId = msg.MessageId
	message.SenderId = c.config.UserId

	err = message.Save()
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (c *MessengerClient) InitChatFromInitializer(userId uint64) (*data.Chat, error) {
	if !c.unlocked {
		return nil, errors.New("not unlocked")
	}
	chat, err := c.http.InitChatFromInitializer(userId)
	if err != nil {
		return nil, err
	}

	chatBucket := c.database.NewChat()
	chatBucket.ChatId = chat.ChatId
	chatBucket.OtherUserId = userId
	chatBucket.LastMessageId = 0
	chatBucket.AmIInitiator = true
	chatBucket.Accepted = false
	chatBucket.OtherUserRsaPublic = nil
	chatBucket.OtherUserEcdsaPublic = nil
	chatBucket.MyRsaPrivate = c.config.InitialRsaRivateKey

	chatBucket.Title = fmt.Sprintf("Chat with %d, %s", userId, time.Now().Format("2006-01-02 15:04:05"))
	err = chatBucket.Save()
	if err != nil {
		return nil, err
	}

	return chatBucket, nil
}

func (c *MessengerClient) InitChatFromReceiver(chatId uint64) (*data.Chat, error) {
	if !c.unlocked {
		return nil, errors.New("not unlocked")
	}
	chat, err := c.http.InitChatFromReceiver(chatId)
	if err != nil {
		return nil, err
	}
	existingChat, err := c.database.GetChat(chat.ChatId)
	if err != nil {
		return nil, err
	}
	if existingChat == nil {
		return nil, errors.New("chat not found")
	}
	existingChat.Accepted = true
	err = existingChat.Update()
	if err != nil {
		return nil, err
	}
	return existingChat, nil
}

func (c *MessengerClient) SearchByUsername(username string) (uint64, error) {
	if !c.unlocked {
		return 0, errors.New("not unlocked")
	}
	search, err := c.http.SearchByUsername(username)
	if err != nil {
		return 0, err
	}
	if search.PublicInfo == nil {
		return 0, nil
	}
	return search.PublicInfo.UserId, nil
}

func (c *MessengerClient) GetChat(chatId uint64) (*data.Chat, error) {
	if !c.unlocked {
		return nil, errors.New("not unlocked")
	}
	return c.database.GetChat(chatId)
}

type WebNotificationType string

const (
	NewIncomingChat WebNotificationType = "new_incoming_chat"
	NewMessage      WebNotificationType = "new_message"
	ChatAccepted    WebNotificationType = "chat_accepted"
)

type WebNotification interface {
	NotificationType() WebNotificationType
}

type IncomingChatNotification struct {
	Chat *data.Chat `json:"chat"`
}

func (i *IncomingChatNotification) NotificationType() WebNotificationType { return NewIncomingChat }

type NewMessageNotification struct {
	ChatId  uint64        `json:"chat_id,omitempty"`
	Message *data.Message `json:"message,omitempty"`
}

func (i *NewMessageNotification) NotificationType() WebNotificationType { return NewMessage }

type ChatAcceptedNotification struct {
	Chat *data.Chat `json:"chat"`
}

func (i *ChatAcceptedNotification) NotificationType() WebNotificationType { return ChatAccepted }

type WebNotificationWithTypeInfo struct {
	Notification WebNotification     `json:"notification"`
	Type         WebNotificationType `json:"type"`
}

func (c *MessengerClient) PullNotificationsAndUpdateData() ([]*WebNotificationWithTypeInfo, error) {
	if !c.unlocked {
		return nil, errors.New("not unlocked")
	}
	notifications, err := c.http.FetchNotifications(100)
	if err != nil {
		return nil, err
	}
	toReturn := make([]WebNotification, 0, len(notifications))
	for _, notification := range notifications {
		inner := notification.Notification
		switch inner.(type) {
		case *custom_types.InitChatFromInitializerNotification:
			notif := inner.(*custom_types.InitChatFromInitializerNotification)
			chat := c.database.NewChat()
			chat.ChatId = notif.ChatId
			chat.OtherUserId = notif.InitializerUserInfo.UserId
			chat.LastMessageId = 0
			chat.AmIInitiator = false
			chat.Accepted = false
			chat.OtherUserRsaPublic = notif.InitializerUserInfo.InitialRsaPublicKey
			chat.OtherUserEcdsaPublic = notif.InitializerUserInfo.EcdsaPublicKey
			chat.MyRsaPrivate = c.config.InitialRsaRivateKey
			chat.Title = fmt.Sprintf("Chat with %d, %s", notif.InitializerUserInfo.UserId, time.Now().Format("2006-01-02 15:04:05"))
			err = chat.Save()
			toReturn = append(toReturn, &IncomingChatNotification{
				Chat: chat,
			})
			break
		case *custom_types.InitChatFromReceiverNotification:
			notif := inner.(*custom_types.InitChatFromReceiverNotification)
			chat, err := c.database.GetChat(notif.ChatId)
			if err != nil {
				log.Printf("error getting chat: %s", err.Error())
				continue
			}
			if chat == nil {
				log.Printf("chat not found")
				continue
			}
			chat.Accepted = true
			chat.OtherUserEcdsaPublic = notif.ReceiverUserInfo.EcdsaPublicKey
			chat.OtherUserRsaPublic = notif.ReceiverUserInfo.InitialRsaPublicKey
			err = chat.Update()
			if err != nil {
				log.Printf("error updating chat: %s", err.Error())
				continue
			}
			toReturn = append(toReturn, &ChatAcceptedNotification{
				Chat: chat,
			})

			break
		case *custom_types.UpdateChatRsaKeyNotification:
			notif := inner.(*custom_types.UpdateChatRsaKeyNotification)
			chat, err := c.database.GetChat(notif.ChatId)
			if err != nil {
				log.Printf("error getting chat: %s", err.Error())
				continue
			}
			if chat == nil {
				log.Printf("chat not found")
				continue
			}
			chat.OtherUserRsaPublic = notif.RsaPublicKey
			err = chat.Update()
			if err != nil {
				log.Printf("error updating chat: %s", err.Error())
				continue
			}
			break
		case *custom_types.SendMessageNotification:
			notif := inner.(*custom_types.SendMessageNotification)
			chat, err := c.database.GetChat(notif.ChatId)
			if err != nil {
				log.Printf("error getting chat: %s", err.Error())
				continue
			}
			if chat == nil {
				log.Printf("chat not found")
				continue
			}

			otherEcPub, err := chat.OtherUserEcdsaPublicKey()
			if err != nil {
				log.Printf("error getting other user ecdsa public key: %s", err.Error())
				continue
			}
			myRsaPriv, err := chat.MyRsaPrivateKey()
			if err != nil {
				log.Printf("error getting my rsa private key: %s", err.Error())
				continue
			}

			messageContent, err := notif.ValidateAndDecrypt(otherEcPub, myRsaPriv)
			if err != nil {
				log.Printf("error decrypting message: %s", err.Error())
				continue
			}
			message := c.database.NewMessage()
			message.ChatId = notif.ChatId
			message.Content = string(messageContent)
			message.MessageId = notif.MessageId
			message.SenderId = notif.SenderUserId
			err = message.Save()
			if err != nil {
				log.Printf("error saving message: %s", err.Error())
				continue
			}

			toReturn = append(toReturn, &NewMessageNotification{
				ChatId:  notif.ChatId,
				Message: message,
			})
			break
		case *custom_types.SendFileNotification:
			break
		default:
			return nil, errors.New("unknown notification type")

		}
	}
	newToreturn := make([]*WebNotificationWithTypeInfo, 0, len(toReturn))
	for _, notif := range toReturn {
		newToreturn = append(newToreturn, &WebNotificationWithTypeInfo{
			Notification: notif,
			Type:         notif.NotificationType(),
		})
	}
	return newToreturn, nil
}

func (c *MessengerClient) GetUsername() string {
	return c.config.Username
}

func (c *MessengerClient) GetUserId() uint64 {
	return c.config.UserId
}

func (c *MessengerClient) SetUsernameConfig(username string, searchable bool) error {
	if !c.unlocked {
		return errors.New("not unlocked")
	}
	_, err := c.http.SetUsernameConfig(username, searchable)
	if err != nil {
		return err
	}
	c.config.Username = username
	c.config.SearchByUsername = searchable
	err = c.config.SaveToFile(configFilename)
	if err != nil {
		return err
	}
	return nil
}

func (c *MessengerClient) TryRequestChat(userIdOrUsername string) (*data.Chat, error) {
	asInt, err := strconv.ParseUint(userIdOrUsername, 10, 64)
	if err == nil {
		return c.InitChatFromInitializer(asInt)
	}
	asInt, err = c.SearchByUsername(userIdOrUsername)
	if err != nil {
		return nil, err
	}
	return c.InitChatFromInitializer(asInt)
}

func (c *MessengerClient) DeleteChat(chatId uint64) error {
	if !c.unlocked {
		return errors.New("not unlocked")
	}
	chat, err := c.database.GetChat(chatId)
	if err != nil {
		return err
	}
	err = chat.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (c *MessengerClient) RenameChat(chatId uint64, newName string) error {
	if !c.unlocked {
		return errors.New("not unlocked")
	}
	chat, err := c.database.GetChat(chatId)
	if err != nil {
		return err
	}
	chat.Title = newName
	err = chat.Update()
	if err != nil {
		return err
	}
	return nil
}
