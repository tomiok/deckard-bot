package main

import (
	"errors"
	"log"
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
