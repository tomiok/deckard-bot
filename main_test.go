package main

import (
	"fmt"
	"log"
	"testing"
)

const blade = "blade"

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
	res, err := getTitleInput("/movie blade")

	if err != nil {
		t.Error()
	}

	fmt.Println(res.cmd)
	fmt.Println(res.title)
}