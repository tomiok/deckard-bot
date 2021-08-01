package main

import (
	"errors"
	"log"
	"strings"
	"testing"
)

const (
	blade = "blade"
)

func Test_fetchMovies(t *testing.T) {
	res, err := fetchMovieInfo(blade)

	if err != nil {
		log.Printf("%v", err.Error())
		t.Error()
	}

	if len(res) == 0 {
		t.Error("cannot be empty for title " + blade)
	}
}

func Test_getInput(t *testing.T) {
	res, err := sanitizeInput(movieCMD + " " + blade)

	if err != nil {
		t.Fatal(err.Error())
	}

	if movie != res.cmd {
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

	if movie != res.cmd {
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
