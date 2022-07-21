package login

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
)

type Server struct {
}

type Users struct {
	Username string
	Name     string
	Password string
	Follows  []string
}

func (s *Server) Authenticate(ctx context.Context, in *LoginDetails) (*LoginResponse, error) {

	resp, err := http.Get("http://127.0.0.1:12380/users")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	// we initialize our Users array
	var users []Users

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(body, &users)

	done := false

	for _, b := range users {
		if b.Username == in.Username && b.Password == in.Password {
			// validate = true
			fmt.Println("Validated successfully!")
			done = true
		}
	}

	return &LoginResponse{Name: in.Username, Done: done}, nil
}
