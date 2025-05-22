package main

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(bytes), err
}

func checkPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func generateToken(length int) (string, error) {
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func formatLinks(unformattedLinks string) []Link {
	var links []Link

	parts := strings.SplitSeq(unformattedLinks, "#")
	for part := range parts {
		if part == "" {
			continue
		}

		sepIndex := strings.Index(part, ";")
		if sepIndex != -1 {
			title := part[:sepIndex]
			link := part[sepIndex+1:]

			links = append(links, Link{
				Title: title,
				Link:  link,
			})
		}
	}

	return links
}
