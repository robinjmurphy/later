package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"

	"github.com/garyburd/go-oauth/oauth"
)

const accessTokenUrl = "https://www.readability.com/api/rest/v1/oauth/access_token/"

type Credentials struct {
	AccessToken       string
	AccessTokenSecret string
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

func login(key string, secret string) (Credentials, error) {
	credentials := Credentials{}
	err := credentials.load()
	if credentials.AccessToken != "" {
		return credentials, nil
	}
	username := read("Username: ")
	password := read("password: ")
	credentials, err = requestAccessToken(key, secret, username, password)
	if err != nil {
		log.Fatal(err)
	}
	err = credentials.save()
	if err != nil {
		log.Fatal(err)
	}
	return credentials, nil
}

func requestAccessToken(key string, secret string, username string, password string) (Credentials, error) {
	client := oauth.Client{Credentials: oauth.Credentials{Token: key, Secret: secret}}
	data := url.Values{
		"x_auth_username": {username},
		"x_auth_password": {password},
		"x_auth_mode":     {"client_auth"},
	}
	client.SignForm(nil, "POST", accessTokenUrl, data)
	resp, err := http.PostForm(accessTokenUrl, data)
	if err != nil {
		return Credentials{}, err
	}
	if resp.StatusCode != 200 {
		return Credentials{}, errors.New("Login failed. Please check your username and password.")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if (err) != nil {
		return Credentials{}, err
	}
	formData, err := url.ParseQuery(string(body))
	if (err) != nil {
		return Credentials{}, err
	}
	return Credentials{
		AccessToken:       formData.Get("oauth_token"),
		AccessTokenSecret: formData.Get("oauth_token_secret"),
	}, nil
}

func main() {
	key := os.Getenv("READABILITY_API_KEY")
	secret := os.Getenv("READABILITY_API_SECRET")
	if key == "" || secret == "" {
		printMissingCredentialsMessage()
	}
	_, err := login(key, secret)
	if err != nil {
		log.Fatal(err)
	}
}
