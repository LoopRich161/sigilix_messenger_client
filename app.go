package main

import (
	"context"
	"fmt"
	"github.com/apepenkov/wails_sigilix_interface/sigilix/data"
	"github.com/apepenkov/wails_sigilix_interface/sigilix/messenger_client"
)

// App struct
type App struct {
	Client *messenger_client.MessengerClient
	ctx    context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	c := messenger_client.NewClient("https://sigilix.aperlaqf.work/api/")
	//err := c.ConnectSqlite("sigilix.db")
	//if err != nil {
	//	panic(err)
	//}
	ap := &App{
		Client: c,
	}
	//c.Unlock("123")
	//c.TryRequestChat("apepenkov2")
	//c.GetChats()
	//a, _ := c.PullNotificationsAndUpdateData()
	//_ = a

	return ap
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after the front-end dom has been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}

func (a *App) IsSignedUp() bool {
	return a.Client.IsSignedUp()
}

func (a *App) IsUnlocked() bool {
	return a.Client.IsUnlocked()
}

func (a *App) Unlock(password string) error {
	return a.Client.Unlock(password)
}

func (a *App) SignUp(password string) error {
	return a.Client.SignUp(password)
}

func (a *App) GetState() string {
	if !a.IsSignedUp() {
		return "signup"
	}
	if !a.IsUnlocked() {
		return "login"
	}

	return "messenger"
}

func (a *App) GetChats() ([]*data.Chat, error) {
	return a.Client.GetChats()
}

func (a *App) GetChatMessages(chatId uint64) ([]*data.Message, error) {
	return a.Client.GetChatMessages(chatId)
}

func (a *App) SendMessage(chatId uint64, text string) (*data.Message, error) {
	return a.Client.SendMessage(chatId, text)
}

func (a *App) GetUsername() string {
	return a.Client.GetUsername()
}

func (a *App) GetUserId() uint64 {
	return a.Client.GetUserId()
}

func (a *App) InitChatFromInitializer(userId uint64) (*data.Chat, error) {
	return a.Client.InitChatFromInitializer(userId)
}

func (a *App) InitChatFromReceiver(chatId uint64) (*data.Chat, error) {
	return a.Client.InitChatFromReceiver(chatId)
}

func (a *App) SearchByUsername(searchUsername string) (uint64, error) {
	return a.Client.SearchByUsername(searchUsername)
}

func (a *App) SetUsernameConfig(setUsername string, searchable bool) error {
	return a.Client.SetUsernameConfig(setUsername, searchable)
}

func (a *App) GetChat(chatId uint64) (*data.Chat, error) {
	return a.Client.GetChat(chatId)
}

func (a *App) PullNotificationsAndUpdateData() ([]*messenger_client.WebNotificationWithTypeInfo, error) {
	return a.Client.PullNotificationsAndUpdateData()
}
func (a *App) TryRequestChat(userIdOrUsername string) (*data.Chat, error) {
	return a.Client.TryRequestChat(userIdOrUsername)
}

func (a *App) RenameChat(chatId uint64, newName string) error {
	return a.Client.RenameChat(chatId, newName)
}

func (a *App) DeleteChat(chatId uint64) error {
	return a.Client.DeleteChat(chatId)
}
