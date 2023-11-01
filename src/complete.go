package main

import (
	"embed"
	"fmt"
	"log"
)

//go:embed autocomplete/*
var static embed.FS

func load(path string) (string, error) {
	file, err := static.ReadFile(path)
	if err != nil {
		log.Fatal("Cant load static file.\n", err)
		return "", err
	}

	return string(file), nil
}

func bashComplete() {
	script, err := load("autocomplete/bash_autocomplete")
	if err != nil {
		log.Fatal("Error: loading bash script is failed.\n", err)
	}
	fmt.Println(script)
}

func zshComplete() {
	script, err := load("autocomplete/zsh_autocomplete")
	if err != nil {
		log.Fatal("Error: loading zsh script is failed.\n", err)
	}
	fmt.Println(script)
}

func fishComplete() {
	script, err := load("autocomplete/fish_autocomplete.fish")
	if err != nil {
		log.Fatal("Error: loading fish script is failed.\n", err)
	}
	fmt.Println(script)
}
