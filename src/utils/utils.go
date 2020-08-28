package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func GetEnvVar(key string) string {
	err := godotenv.Load("../.env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	return os.Getenv(key)
}

func FailOnError(err error) {
	if err != nil {
		fmt.Println("An error ocorrured!")
		fmt.Println(err)
	}
}

func Hash(key string) string {
	hash := sha1.New()
	fmt.Println(strings.ToLower(key))
	hash.Write([]byte(strings.ToLower(key)))
	bs := hash.Sum(nil)
	return hex.EncodeToString(bs)
}

func HashAlbum(name string, artist string) string {
	return Hash(name + "\u001F" + artist)
}

func HashTrack(name string, artist string, album string) string {
	return Hash(name + "\u001F" + artist + "\u0010" + album)
}
