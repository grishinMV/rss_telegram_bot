package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	apiUrl string
	token  string
	client http.Client
}

func NewTelegramClient(apiUrl string, token string, client http.Client) Client {
	return Client{
		apiUrl: apiUrl,
		token:  token,
		client: client,
	}
}

func (bot Client) GetMe() (*User, error) {
	response, err := bot.client.Get(bot.apiUrl + "/bot" + bot.token + "/getMe")
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	result := User{}
	err = json.Unmarshal(body, &result)

	return &result, err
}

func (bot Client) GetUpdates(offset int, limit int) (*UpdateResponse, error) {
	params := "offset=" + strconv.Itoa(offset) + "&limit=" + strconv.Itoa(limit)
	query := bot.apiUrl + "/bot" + bot.token + "/getUpdates?" + params
	response, err := bot.client.Get(query)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	result := UpdateResponse{}
	err = json.Unmarshal(body, &result)

	return &result, err
}

func (bot Client) SendTextMessage(chatId int, text string) (*SendMessageResponse, error) {
	params := "chat_id=" + strconv.Itoa(chatId) + "&text=" + text

	query := bot.apiUrl + "/bot" + bot.token + "/sendMessage?" + params
	response, err := bot.client.Get(query)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	result := SendMessageResponse{}
	err = json.Unmarshal(body, &result)

	return &result, err
}

func (bot Client) SendReplyMarkup(chatId int, text string, markup ReplyKeyboardMarkup) (*SendMessageResponse, error) {
	markupJson, err := json.Marshal(markup)
	if err != nil {
		return nil, err
	}

	text = url.QueryEscape(text)

	params := "chat_id=" + strconv.Itoa(chatId) + "&reply_markup=" + string(markupJson) + "&text=" + text

	query := bot.apiUrl + "/bot" + bot.token + "/sendMessage?" + params
	response, err := bot.client.Get(query)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	result := SendMessageResponse{}
	err = json.Unmarshal(body, &result)

	return &result, err
}
