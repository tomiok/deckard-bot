package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	slash    = "/"
	movies   = "movies"
	start    = "start"
	movieCMD = slash + movies
	startCMD = slash + start

	telegramApiBaseUrl     string = "https://api.telegram.org/bot"
	telegramApiSendMessage string = "/sendMessage"
	telegramTokenEnv       string = "TELEGRAM_BOT_TOKEN"
)

var (
	errWrongCMD = errors.New("type /start to get help")
	errTypoCMD  = errors.New("please type the command correctly or use /start")
	errUnknown  = errors.New("unknown error")
)

type seed struct {
	cmd   string
	title string
	fn    func() string
}

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

type MovieResponse struct {
	Title  string `json:"title"`
	Year   string `json:"year"`
	ImdbID string `json:"imdbID"`
	Poster string `json:"poster"`
}

// HandleTelegramWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
func HandleTelegramWebHook(_ http.ResponseWriter, r *http.Request) {
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	log.Print(handleRequest(update))
}

func handleRequest(update *Update) string {
	seed, err := sanitizeInput(update.Message.Text)

	if err != nil {
		res, _ := sendTextToTelegramChat(update.Message.Chat.Id, "command error. Please use /start to start chatting")
		return res
	}

	response := prepareCommand(seed).fn()

	// Send the response back to Telegram
	_, err = sendTextToTelegramChat(update.Message.Chat.Id, response)
	if err != nil {
		return fmt.Sprintf("got error %s from telegram", err.Error())
	}

	return fmt.Sprintf("response successfully distributed to chat id %d", update.Message.Chat.Id)
}

func prepareCommand(seed *seed) *seed {
	switch seed.cmd {
	case start:
		seed.fn = func() string {
			return "Hi, please search a movie with /movies command. Use like /movies blade runner and see what happen"
		}
	case movies:
		seed.fn = func() string {
			req, _ := http.NewRequest("GET", "https://movies-lib-stg.herokuapp.com/query", nil)

			q := req.URL.Query()
			q.Add("s", seed.title)
			req.URL.RawQuery = q.Encode()

			res, err := http.Get(req.URL.String())

			if err != nil {
				log.Printf("err %v", err.Error())
				return "cannot bring the movie"
			}

			body := res.Body

			defer func() {
				_ = body.Close()
			}()

			var movies []MovieResponse
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(body)

			err = json.Unmarshal(buf.Bytes(), &movies)

			if err != nil {
				log.Printf("err %v", err.Error())
				return "cannot parse the API response"
			}
			return displayMoviesRes(movies)
		}
	}

	return seed
}

func sanitizeInput(input string) (*seed, error) {
	if input == "" {
		return nil, errWrongCMD
	}

	if strings.Index(input, slash) != 0 {
		return nil, errWrongCMD
	}

	if strings.Trim(input, " ") == startCMD {
		return &seed{cmd: start}, nil
	}

	words := strings.Fields(input)

	if len(words) <= 1 {
		return nil, errTypoCMD
	}
	cmd := words[0]
	movieTitle := strings.Join(words[1:], " ")

	switch cmd {
	case movieCMD:
		return &seed{cmd: movies, title: movieTitle}, nil
	default:
		return nil, errUnknown
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
func sendTextToTelegramChat(chatId int, res string) error {
	log.Printf("Sending message to chat_id: %d", chatId)

	var telegramApi = telegramApiBaseUrl + os.Getenv(telegramTokenEnv) + telegramApiSendMessage
	response, err := sendMessage(chatId, telegramApi, res)

	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	return nil
}

func sendMessage(chatID int, telegramApi, s string) (*http.Response, error) {
	return http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatID)},
			"text":    {s},
		})
}

func displayMoviesRes(movies []MovieResponse) string {
	var res = make([]string, len(movies))

	for i, movie := range movies {
		res[i] = formatMovieText(movie)
	}

	return strings.Join(res, "")
}

func formatMovieText(m MovieResponse) string {
	format := "Title: %s, Year: %s, imdbID: %s, Poster: %s\n"
	return fmt.Sprintf(format, m.Title, m.Year, m.ImdbID, m.Poster)
}
