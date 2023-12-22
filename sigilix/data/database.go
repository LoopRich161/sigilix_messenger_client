package data

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
)

//import _ "github.com/mattn/go-sqlite3"
import _ "github.com/mutecomm/go-sqlcipher"

var initSqls = []string{
	`CREATE TABLE IF NOT EXISTS chats (
    		chat_id INTEGER PRIMARY KEY,
    		other_user_id INTEGER NOT NULL,
    		last_message_id INTEGER DEFAULT 0,
    		am_i_initiator INTEGER NOT NULL,
    		accepted INTEGER DEFAULT 0 NOT NULL,
    		other_user_rsa_public BLOB,
    		other_user_ecdsa_public BLOB,
    		my_rsa_private BLOB NOT NULL,
    		title TEXT DEFAULT '' NOT NULL
	);`,
	`CREATE TABLE IF NOT EXISTS messages (
    		message_id INTEGER,
    		chat_id INTEGER NOT NULL,
    		sender_id INTEGER NOT NULL,
    		content TEXT NOT NULL,
    		FOREIGN KEY(chat_id) REFERENCES chats(chat_id) ON DELETE CASCADE
	);`,
	//`CREATE TABLE IF NOT EXISTS config (
	//		user_id INTEGER PRIMARY KEY,
	//		username TEXT NOT NULL,
	//		search_by_username INTEGER NOT NULL,
	//		initial_rsa_private_key BLOB NOT NULL,
	//		initial_ecdsa_private_key BLOB NOT NULL,
	//		password_hash BLOB NOT NULL
	//);`,
	// set messages primary key to (chat_id, message_id)
	`CREATE UNIQUE INDEX IF NOT EXISTS messages_chat_id_message_id ON messages (chat_id, message_id);`,
	// additional indexes
	`CREATE INDEX IF NOT EXISTS chats_other_user_id ON chats (other_user_id);`,
	`CREATE INDEX IF NOT EXISTS messages_sender_id ON messages (sender_id);`,
	`CREATE INDEX IF NOT EXISTS messages_chat_id ON messages (chat_id);`,
}

type SqliteDB struct {
	db *sql.DB
}

func NewSqliteDB(filename string) (*SqliteDB, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		return nil, err
	}
	for _, sqlq := range initSqls {
		_, err = db.Exec(sqlq)
		if err != nil {
			return nil, err
		}
	}
	return &SqliteDB{db: db}, nil
}
func NewSqliteDBWithPassword(filename string, password string) (*SqliteDB, error) {
	password = strings.ToUpper(password)
	key := url.QueryEscape(password)
	ur := fmt.Sprintf("%s?_pragma_key=%s&_pragma_cipher_page_size=4096", filename, key)
	db, err := sql.Open("sqlite3", ur)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		return nil, err
	}
	for _, sqlq := range initSqls {
		_, err = db.Exec(sqlq)
		if err != nil {
			return nil, err
		}
	}
	return &SqliteDB{db: db}, nil
}

func (s *SqliteDB) Close() error {
	return s.db.Close()
}

func (s *SqliteDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

func (s *SqliteDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

func (s *SqliteDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.db.QueryRow(query, args...)
}

func (s *SqliteDB) NewChat() *Chat {
	return &Chat{db: s}
}

func (s *SqliteDB) NewMessage() *Message {
	return &Message{db: s}
}

func (s *SqliteDB) GetAllChats() ([]*Chat, error) {
	rows, err := s.Query("SELECT chat_id, other_user_id, last_message_id, am_i_initiator, accepted, other_user_rsa_public, other_user_ecdsa_public, my_rsa_private, title FROM chats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	chats := make([]*Chat, 0)
	for rows.Next() {
		chat := &Chat{
			db: s,
		}
		err = rows.Scan(&chat.ChatId, &chat.OtherUserId, &chat.LastMessageId, &chat.AmIInitiator, &chat.Accepted, &chat.OtherUserRsaPublic, &chat.OtherUserEcdsaPublic, &chat.MyRsaPrivate, &chat.Title)
		if err != nil {
			return nil, err
		}
		chatMessages, err := s.GetMessages(chat.ChatId)
		if err != nil {
			return nil, err
		}
		chat.Messages = chatMessages
		chats = append(chats, chat)
	}
	return chats, nil
}

func (s *SqliteDB) GetChat(chatId uint64) (*Chat, error) {
	row := s.QueryRow("SELECT chat_id, other_user_id, last_message_id, am_i_initiator, accepted, other_user_rsa_public, other_user_ecdsa_public, my_rsa_private, title FROM chats WHERE chat_id = ?", chatId)
	chat := &Chat{
		db: s,
	}
	err := row.Scan(&chat.ChatId, &chat.OtherUserId, &chat.LastMessageId, &chat.AmIInitiator, &chat.Accepted, &chat.OtherUserRsaPublic, &chat.OtherUserEcdsaPublic, &chat.MyRsaPrivate, &chat.Title)
	if err != nil {
		return nil, err
	}
	chatMessages, err := s.GetMessages(chatId)
	if err != nil {
		return nil, err
	}
	chat.Messages = chatMessages
	return chat, nil
}

func (s *SqliteDB) GetMessages(chatId uint64) ([]*Message, error) {
	rows, err := s.Query("SELECT message_id, chat_id, sender_id, content FROM messages WHERE chat_id = ?", chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := make([]*Message, 0)
	for rows.Next() {
		message := &Message{
			db: s,
		}
		err = rows.Scan(&message.MessageId, &message.ChatId, &message.SenderId, &message.Content)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
