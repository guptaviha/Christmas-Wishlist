package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	"golang.org/x/net/context"
)

type User struct {
	Username string
	Name     string
	Password string
	Follows  []string
}

type Server struct {
}

func (s *Server) SignUpUser(ctx context.Context, in *SignUpRequest) (*SignUpResponse, error) {
	if in.Password != in.ConfirmPassword {
		return &SignUpResponse{Success: false}, nil
	}

	resp, err := http.Get("http://127.0.0.1:12380/users")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	users := []User{}

	json.Unmarshal(body, &users)

	for _, u := range users {
		if u.Username == in.Username {
			return &SignUpResponse{Success: false}, nil
		}
	}

	users = append(users, User{
		Username: in.Username,
		Name:     in.Name,
		Password: in.Password,
		Follows:  []string{},
	})

	dataBytes, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
	}

	cmd := exec.Command("curl", "-L", "http://127.0.0.1:12380/users", "-XPUT", "-d "+string(dataBytes))
	cmd.Run()

	return &SignUpResponse{Success: true}, nil
}

func (s *Server) GetFollowing(ctx context.Context, in *FollowerRequest) (*FollowerResponse, error) {

	resp, err := http.Get("http://127.0.0.1:12380/users")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	users := []User{}

	json.Unmarshal(body, &users)

	notFollowed := []string{}
	Followed := []string{}

	for _, u := range users {
		if in.Username == u.Username {
			Followed = u.Follows
		}
	}

	for _, u := range users {
		if !CheckIfFollowed(u.Username, Followed) {
			notFollowed = append(notFollowed, u.Username)
		}
	}

	return &FollowerResponse{
		Followers:    Followed,
		NotFollowers: notFollowed,
	}, nil
}

func CheckIfFollowed(username string, follows []string) bool {
	for _, s := range follows {
		if s == username {
			return true
		}
	}
	return false
}

func (s *Server) UpdateFollower(ctx context.Context, in *UpdateFollowersRequest) (*UpdateFollowersResponse, error) {

	resp, err := http.Get("http://127.0.0.1:12380/users")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	users := []User{}

	json.Unmarshal(body, &users)

	for i, u := range users {
		if u.Username == in.Username {
			if in.IsFollow {
				users[i].Follows = append(users[i].Follows, in.Newuser)
			} else {
				NewFollows := []string{}
				for _, f := range u.Follows {
					if in.Newuser != f {
						NewFollows = append(NewFollows, f)
					}
				}
				users[i].Follows = NewFollows
			}
		}
	}

	dataBytes, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(dataBytes))

	cmd := exec.Command("curl", "-L", "http://127.0.0.1:12380/users", "-XPUT", "-d "+string(dataBytes))
	cmd.Run()

	return &UpdateFollowersResponse{Success: true}, nil
}
