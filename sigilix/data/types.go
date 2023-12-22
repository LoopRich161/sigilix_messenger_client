package data

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"github.com/apepenkov/wails_sigilix_interface/sigilix/crypto_utils"
	"github.com/apepenkov/wails_sigilix_interface/sigilix/custom_types"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"os"
)

type Chat struct {
	ChatId               uint64 `json:"chat_id"`
	OtherUserId          uint64 `json:"other_user_id"`
	LastMessageId        uint64 `json:"last_message_id"`
	AmIInitiator         bool   `json:"am_i_initiator"`
	Accepted             bool   `json:"accepted"`
	OtherUserRsaPublic   []byte
	OtherUserEcdsaPublic []byte
	MyRsaPrivate         []byte
	Title                string `json:"title"`

	Messages []*Message `json:"messages"`

	db *SqliteDB
}

func (c *Chat) OtherUserRsaPublicKey() (*rsa.PublicKey, error) {
	return crypto_utils.PublicRSAKeyFromBytes(c.OtherUserRsaPublic)
}

func (c *Chat) OtherUserEcdsaPublicKey() (*ecdsa.PublicKey, error) {
	return crypto_utils.PublicECDSAKeyFromBytes(c.OtherUserEcdsaPublic)
}

func (c *Chat) MyRsaPrivateKey() (*rsa.PrivateKey, error) {
	return crypto_utils.RsaPrivateFromBytes(c.MyRsaPrivate)
}

func (c *Chat) Save() error {
	_, err := c.db.Exec(
		"INSERT INTO chats (chat_id, other_user_id, last_message_id, am_i_initiator, accepted, other_user_rsa_public, other_user_ecdsa_public, my_rsa_private, title) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		c.ChatId, c.OtherUserId, c.LastMessageId, c.AmIInitiator, c.Accepted, c.OtherUserRsaPublic, c.OtherUserEcdsaPublic, c.MyRsaPrivate, c.Title,
	)
	return err
}

func (c *Chat) Update() error {
	_, err := c.db.Exec(
		"UPDATE chats SET other_user_id = ?, last_message_id = ?, am_i_initiator = ?, accepted = ?, other_user_rsa_public = ?, other_user_ecdsa_public = ?, my_rsa_private = ?, title = ? WHERE chat_id = ?",
		c.OtherUserId, c.LastMessageId, c.AmIInitiator, c.Accepted, c.OtherUserRsaPublic, c.OtherUserEcdsaPublic, c.MyRsaPrivate, c.Title, c.ChatId,
	)
	return err
}

func (c *Chat) Delete() error {
	_, err := c.db.Exec("DELETE FROM chats WHERE chat_id = ?", c.ChatId)
	return err
}

type Message struct {
	MessageId uint64 `json:"message_id"`
	ChatId    uint64 `json:"chat_id"`
	SenderId  uint64 `json:"sender_id"`
	Content   string `json:"content"`

	db *SqliteDB
}

func (m *Message) Save() error {
	_, err := m.db.Exec(
		"INSERT INTO messages (message_id, chat_id, sender_id, content) VALUES (?, ?, ?, ?)",
		m.MessageId, m.ChatId, m.SenderId, m.Content,
	)

	return err
}

func (m *Message) Update() error {
	_, err := m.db.Exec(
		"UPDATE messages SET chat_id = ?, sender_id = ?, content = ? WHERE message_id = ?",
		m.ChatId, m.SenderId, m.Content, m.MessageId,
	)
	return err
}

func (m *Message) Delete() error {
	_, err := m.db.Exec("DELETE FROM messages WHERE message_id = ?", m.MessageId)
	return err
}

type Config struct {
	UserId                 uint64                   `json:"user_id"`
	Username               string                   `json:"username"`
	SearchByUsername       bool                     `json:"search_by_username"`
	InitialRsaRivateKey    custom_types.Base64Bytes `json:"initial_rsa_rivate_key"`
	InitialECDSAPrivateKey custom_types.Base64Bytes `json:"initial_ecdsa_private_key"`
	PaswordHash            custom_types.Base64Bytes `json:"pasword_hash"`
}

func EncryptDataWithBytes(data []byte, password []byte) ([]byte, error) {
	salt := []byte("salt")

	key := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func DecryptDataWithBytes(data []byte, password []byte) ([]byte, error) {
	salt := []byte("salt")

	key := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func LoadConfigFromFiles(filename string, password []byte) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	read, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	decryptedBytes, err := DecryptDataWithBytes(read, password)

	if err != nil {
		return nil, err
	}

	var config *Config

	err = json.Unmarshal(decryptedBytes, &config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) SaveToFile(filename string) error {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}
	encryptedBytes, err := EncryptDataWithBytes(jsonBytes, c.PaswordHash)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, encryptedBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) RsaPrivateKey() (*rsa.PrivateKey, error) {
	return crypto_utils.RsaPrivateFromBytes(c.InitialRsaRivateKey)
}

func (c *Config) EcdsaPrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto_utils.PrivateKeyFromBytes(c.InitialECDSAPrivateKey)
}

func (c *Config) MustRsaPrivateKey() *rsa.PrivateKey {
	key, err := c.RsaPrivateKey()
	if err != nil {
		panic(err)
	}
	return key
}

func (c *Config) MustEcdsaPrivateKey() *ecdsa.PrivateKey {
	key, err := c.EcdsaPrivateKey()
	if err != nil {
		panic(err)
	}
	return key
}

func (c *Config) RsaPublicKey() (*rsa.PublicKey, error) {
	k, err := c.RsaPrivateKey()
	if err != nil {
		return nil, err
	}
	return k.Public().(*rsa.PublicKey), nil
}

func (c *Config) EcdsaPublicKey() (*ecdsa.PublicKey, error) {
	k, err := c.EcdsaPrivateKey()
	if err != nil {
		return nil, err
	}
	return k.Public().(*ecdsa.PublicKey), nil
}

func (c *Config) MustRsaPublicKey() *rsa.PublicKey {
	return c.MustRsaPrivateKey().Public().(*rsa.PublicKey)
}

func (c *Config) MustEcdsaPublicKey() *ecdsa.PublicKey {
	return c.MustEcdsaPrivateKey().Public().(*ecdsa.PublicKey)
}
