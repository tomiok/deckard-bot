package main

import (
	"errors"
	"strings"
	"testing"
)

const (
	blade = "blade"
)

func Test_prepareCommandMovies(t *testing.T) {
	res := prepareCommand(&seed{
		cmd:   movies,
		title: "blade",
	})

	movies := res.fn()

	if len(movies) < 100 {
		t.Error()
	}
}

func Test_prepareCommandMoviesLongTitle(t *testing.T) {
	res := prepareCommand(&seed{
		cmd:   movies,
		title: "tron legacy",
	})

	movies := res.fn()

	if len(movies) < 100 {
		t.Error()
	}
}

func Test_prepareCommandStart(t *testing.T) {
	res := prepareCommand(&seed{
		cmd: start,
	})

	message := res.fn()

	if message != "Hi, please search a movie with /movies command. Use like /movies blade runner and see what happen" {
		t.Error()
	}
}

func Test_getInput(t *testing.T) {
	res, err := sanitizeInput(movieCMD + " " + blade)

	if err != nil {
		t.Fatal(err.Error())
	}

	if movies != res.cmd {
		t.Fatal("command should be movie")
	}

	if res.title != blade {
		t.Fatal("title should be " + blade)
	}
}

func Test_getInputWithSpaces(t *testing.T) {
	movieTitle := "tron legacy"
	res, err := sanitizeInput(movieCMD + " " + movieTitle)
	if err != nil {
		t.Fatal(err.Error())
	}

	if movies != res.cmd {
		t.Fatal("command should be movie")
	}

	if res.title != movieTitle {
		t.Fatal("title should be " + movieTitle)
	}
}

func Test_inputStartCMD(t *testing.T) {
	res, err := sanitizeInput(startCMD)

	if err != nil {
		t.Fatal("err should be nil")
	}

	if res.cmd != start {
		t.Fatal("cmd should be start")
	}

	if res.title != "" {
		t.Fatal("title should be empty")
	}
}

func Test_errInput(t *testing.T) {
	_, err := sanitizeInput("start")

	if err == nil {
		t.Fatal("input is wrong, should return an error")
	}

	if !errors.Is(err, errWrongCMD) {
		t.Fatal("wrong error type")
	}
}

func Test_emptyInput(t *testing.T) {
	_, err := sanitizeInput("")

	if !errors.Is(err, errWrongCMD) {
		t.Error("wrong error type")
	}
}

func Test_nameDisplay(t *testing.T) {
	movie := MovieResponse{
		Title:  "Alien",
		Year:   "1968",
		ImdbID: "12345",
		Poster: "https://alien.jpg",
	}

	s := formatMovieText(movie)
	expected := "Title: Alien, Year: 1968, imdbID: 12345, Poster: https://alien.jpg\n"

	if expected != s {
		t.Error("wrong text")
	}
}

func Test_displayMovies(t *testing.T) {
	var movies = []MovieResponse{{
		Title:  "Blade Runner",
		Year:   "1984",
		ImdbID: "12345",
		Poster: "blade_runner.jpg",
	},
		{
			Title:  "Blade",
			Year:   "2011",
			ImdbID: "456789",
			Poster: "blade.jpg",
		},
	}

	res := displayMoviesRes(movies)

	if !strings.Contains(res, "12345") {
		t.Error("imdbID is not present")
	}
}

func Test_handleRequest_start(t *testing.T) {
	u := &Update{
		UpdateId: 1,
		Message: Message{
			Text: "/start",
			Chat: Chat{Id: 1},
		},
	}

	res := handleRequest(u)
	if res != "response successfully distributed to chat id 1" {
		t.Error()
	}
}
