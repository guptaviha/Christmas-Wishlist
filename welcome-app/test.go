package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"welcome-app/login"
)

func main() {

	resp, err := http.Get("http://127.0.0.1:12380/users")
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var users []login.Users
	json.Unmarshal(body, &users)

}
