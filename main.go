package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const telegramApiBaseUrl string = "https://api.telegram.org/bot"
const telegramApiSendMessage string = "/sendMessage"
const telegramTokenEnv string = "TELEGRAM_BOT_TOKEN"

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	Id int `json:"id"`
}

// HandleTelegramWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
func HandleTelegramWebHook(_ http.ResponseWriter, r *http.Request) {

	// Parse incoming request
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	seed, err := getTitleInput(update.Message.Text)

	if err != nil {
		log.Printf("errors getting command, %s", err.Error())
		return
	}

	res, err := fetchMovieInfo(seed.title)

	if err != nil {
		log.Printf("errors getting movies, %s", err.Error())
		return
	}

	// Send the punchline back to Telegram
	var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.Id, res)
	if errTelegram != nil {
		log.Printf("got error %s from telegram, response body is %s", errTelegram.Error(), telegramResponseBody)
		return
	}

	log.Printf("response %s successfully distributed to chat id %d", res, update.Message.Chat.Id)
}

func fetchMovieInfo(seed string) ([]MoviesResponse, error) {
	moviesUrl := fmt.Sprintf("https://movies-lib-stg.herokuapp.com/query?s=%s", seed)
	res, err := http.Get(moviesUrl)

	if err != nil {
		return nil, err
	}

	body := res.Body
	defer body.Close()

	var movies []MoviesResponse
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(body)

	err = json.Unmarshal(buf.Bytes(), &movies)

	if err != nil {
		return nil, err
	}

	return movies, nil
}

const (
	slash    = "/"
	movie    = "movie"
	movieCMD = slash + movie
)

var (
	movieLen = len(movieCMD)
)

type seed struct {
	cmd   string
	title string
}

func getTitleInput(input string) (*seed, error) {
	if input == "" || len(input) <= movieLen {
		return nil, errors.New("input is empty")
	}

	cmd := strings.Split(input, " ")

	if len(cmd) <= 1 {
		return nil, errors.New("please type the command correctly")
	}

	switch cmd[0] {
	case movieCMD:
		return &seed{cmd: movie, title: cmd[1]}, nil
	default:
		return nil, errors.New("unknown error")
	}
}

func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}
	return &update, nil
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendTextToTelegramChat(chatId int, movies []MoviesResponse) (string, error) {
	log.Printf("Sending message to chat_id: %d", chatId)

	if len(movies) == 0 {
		log.Print("no movies founded")
		return "", nil
	}

	text := movies[0].Title

	var telegramApi = telegramApiBaseUrl + os.Getenv(telegramTokenEnv) + telegramApiSendMessage
	log.Printf("api: %s", telegramApi)
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	log.Printf("chat id %d", chatId)

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}

type MoviesResponse struct {
	Title  string `json:"title"`
	Year   string `json:"year"`
	ImdbID string `json:"imdbID"`
	Poster string `json:"poster"`
}