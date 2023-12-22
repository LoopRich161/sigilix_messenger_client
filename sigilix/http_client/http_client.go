package http_client

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apepenkov/wails_sigilix_interface/sigilix/crypto_utils"
	"github.com/apepenkov/wails_sigilix_interface/sigilix/custom_types"
	"io"
	"net/http"
)

type SigilixHttpClient struct {
	httpClient   *http.Client
	baseUrl      string
	ecdsaPrivate *ecdsa.PrivateKey
	userId       uint64
}

func NewSigilixHttpClient(baseUrl string, ecdsaPrivate *ecdsa.PrivateKey, userId uint64) *SigilixHttpClient {
	return &SigilixHttpClient{
		httpClient:   &http.Client{},
		baseUrl:      baseUrl,
		ecdsaPrivate: ecdsaPrivate,
		userId:       userId,
	}
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("Error: %s, code: %v", e.Message, e.Code)
}

func (c *SigilixHttpClient) makeRequest(url string, body custom_types.SigilixStruct, writeTo custom_types.SigilixStruct) error {
	path := c.baseUrl + url

	encoded, err := json.Marshal(body)

	if err != nil {
		return err
	}

	signature, err := crypto_utils.SignMessageBase64(c.ecdsaPrivate, encoded)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", path, bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Sigilix-Signature", signature)
	req.Header.Set("X-Sigilix-User-Id", fmt.Sprintf("%d", c.userId))

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		var errorResponse ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to decode error response: %v. Server responded with status code %d", err, resp.StatusCode))
		}
		return &errorResponse
	}

	isDebug := true

	if !isDebug {
		if err := json.NewDecoder(resp.Body).Decode(writeTo); err != nil {
			return errors.New(fmt.Sprintf("failed to decode response: %v", err))
		}
	} else {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to read response body: %v", err))
		}
		err = json.Unmarshal(fullBody, writeTo)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to decode response: %v", err))
		}
	}

	return nil
}

func (c *SigilixHttpClient) Login(ecdsaPublicKey *ecdsa.PublicKey, rsaPublicKey *rsa.PublicKey) (*custom_types.LoginResponse, error) {
	rsaBytes, err := crypto_utils.PublicRSAKeyToBytes(rsaPublicKey)
	if err != nil {
		return nil, err
	}
	req := &custom_types.LoginRequest{
		ClientEcdaPublicKey: crypto_utils.PublicECDSAKeyToBytes(ecdsaPublicKey),
		ClientRsaPublicKey:  rsaBytes,
	}

	resp := &custom_types.LoginResponse{}

	err = c.makeRequest("users/login", req, resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *SigilixHttpClient) SetUsernameConfig(setUsername string, searchable bool) (*custom_types.SetUsernameConfigResponse, error) {
	req := &custom_types.SetUsernameConfigRequest{
		Username:                setUsername,
		SearchByUsernameAllowed: searchable,
	}

	resp := &custom_types.SetUsernameConfigResponse{}

	err := c.makeRequest("users/set_username_config", req, resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *SigilixHttpClient) SearchByUsername(searchUsername string) (*custom_types.SearchByUsernameResponse, error) {
	req := &custom_types.SearchByUsernameRequest{
		Username: searchUsername,
	}

	resp := &custom_types.SearchByUsernameResponse{}

	err := c.makeRequest("users/search_by_username", req, resp)

	if err != nil {
		return nil, err
	}

	if resp.PublicInfo == nil || resp.PublicInfo.UserId == 0 {
		return nil, nil
	}

	return resp, nil
}

func (c *SigilixHttpClient) InitChatFromInitializer(userId uint64) (*custom_types.InitChatFromInitializerResponse, error) {
	req := &custom_types.InitChatFromInitializerRequest{
		TargetUserId: userId,
	}

	resp := &custom_types.InitChatFromInitializerResponse{}

	err := c.makeRequest("messages/init_chat_from_initializer", req, resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *SigilixHttpClient) InitChatFromReceiver(chatId uint64) (*custom_types.InitChatFromReceiverResponse, error) {
	req := &custom_types.InitChatFromReceiverRequest{
		ChatId: chatId,
	}

	resp := &custom_types.InitChatFromReceiverResponse{}

	err := c.makeRequest("messages/init_chat_from_receiver", req, resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *SigilixHttpClient) UpdateChatRsaKey(chatId uint64, rsaPublicKey *rsa.PublicKey) (*custom_types.UpdateChatRsaKeyResponse, error) {
	rsaBytes, err := crypto_utils.PublicRSAKeyToBytes(rsaPublicKey)
	if err != nil {
		return nil, err
	}
	req := &custom_types.UpdateChatRsaKeyRequest{
		ChatId:       chatId,
		RsaPublicKey: rsaBytes,
	}

	resp := &custom_types.UpdateChatRsaKeyResponse{}

	err = c.makeRequest("messages/update_chat_rsa_key", req, resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *SigilixHttpClient) SendMessage(chatId uint64, message string, rsaPublicKey *rsa.PublicKey) (*custom_types.SendMessageResponse, error) {

	messageBytes := []byte(message)
	ecdsaSignature, err := crypto_utils.SignMessage(c.ecdsaPrivate, messageBytes)
	if err != nil {
		return nil, err
	}
	rsaEncryptedMessage, err := crypto_utils.EncryptMessage(rsaPublicKey, messageBytes)
	if err != nil {
		return nil, err
	}

	req := &custom_types.SendMessageRequest{
		ChatId:                chatId,
		EncryptedMessage:      rsaEncryptedMessage,
		MessageEcdsaSignature: ecdsaSignature,
	}

	resp := &custom_types.SendMessageResponse{}

	err = c.makeRequest("messages/send_message", req, resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *SigilixHttpClient) FetchNotifications(limit uint32) ([]*custom_types.IncomingNotification, error) {
	req := &custom_types.GetNotificationsRequest{
		Limit: limit,
	}

	resp := &custom_types.GetNotificationsResponse{}

	err := c.makeRequest("messages/get_notifications", req, resp)

	if err != nil {
		return nil, err
	}

	return resp.Notifications, nil
}
