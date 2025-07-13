package main

import "strings"

func wordCleanUp(chirp string) string {
	words := strings.Fields(chirp)
	for i, word := range words {
		word := strings.ToLower(word)
		switch word {
		case "kerfuffle", "sharbert", "fornax":
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
