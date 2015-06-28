package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/robinjmurphy/go-readability-api/readability"
)

type Credentials struct {
	AccessToken       string
	AccessTokenSecret string
}

func usage() {
	fmt.Println("usage: later <url>")
	flag.PrintDefaults()
	os.Exit(1)
}

func getConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir + "/.later", nil
}

func (credentials *Credentials) save() error {
	serialized, err := json.Marshal(credentials)
	if err != nil {
		return err
	}
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, serialized, 0644)
}

func (credentials *Credentials) load() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &credentials)
	if err != nil {
		return err
	}
	return nil
}

func read(message string) string {
	fmt.Print(message)
	var value string
	_, err := fmt.Scanln(&value)
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func printMissingCredentialsMessage() {
	fmt.Println("Ensure that READABILITY_API_KEY and READABILITY_API_SECRET are set. ")
	fmt.Println("See https://github.com/robinjmurphy/later#installation.")
	os.Exit(1)
}

func login(client *readability.Client) (readerClient *readability.ReaderClient, err error) {
	credentials := Credentials{}
	credentials.load()
	if credentials.AccessToken != "" {
		return client.NewReaderClient(credentials.AccessToken, credentials.AccessTokenSecret), nil
	}
	username := read("Username: ")
	password := read("Password: ")
	token, secret, err := client.Login(username, password)
	if err != nil {
		return readerClient, err
	}
	credentials.AccessToken = token
	credentials.AccessTokenSecret = secret
	err = credentials.save()
	if err != nil {
		return readerClient, err
	}
	return client.NewReaderClient(token, secret), nil
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		usage()
	}
	url := args[0]
	key := os.Getenv("READABILITY_API_KEY")
	secret := os.Getenv("READABILITY_API_SECRET")
	if key == "" || secret == "" {
		printMissingCredentialsMessage()
	}
	client := readability.NewClient(key, secret)
	reader, err := login(client)
	if err != nil {
		log.Fatal(err)
	}
	_, err = reader.AddBookmark(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully bookmarked %s", url)
}
