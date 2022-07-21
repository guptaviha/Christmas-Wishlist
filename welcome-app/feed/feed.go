package feed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

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

type Post struct {
	PostID      int
	Title       string
	Author      string
	Description string
	Timestamp   string
}

func (s *Server) GetFeed(ctx context.Context, in *FeedRequest) (*FeedResponse, error) {
	log.Printf("Receieved following details from Client: \nusername: %s", in.Username)
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

	followers := map[string]bool{}

	for _, b := range users {
		if b.Username == in.Username {
			for _, follower := range b.Follows {
				followers[follower] = true
			}
		}
	}

	resp, err = http.Get("http://127.0.0.1:12380/posts")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	// we initialize our Users array
	var posts []Post
	filterPostIDs := []int32{}
	filterPostTitles := []string{}
	filterPostAuthor := []string{}
	filterPostDescription := []string{}
	filterPostTimestamps := []string{}
	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(body, &posts)

	for _, p := range posts {
		if followers[p.Author] {
			filterPostIDs = append(filterPostIDs, int32(p.PostID))
			filterPostTitles = append(filterPostTitles, p.Title)
			filterPostAuthor = append(filterPostAuthor, p.Author)
			filterPostDescription = append(filterPostDescription, p.Description)
			filterPostTimestamps = append(filterPostTimestamps, p.Timestamp)
		}
	}

	return &FeedResponse{Postid: filterPostIDs,
		Title:       filterPostTitles,
		Author:      filterPostAuthor,
		Description: filterPostDescription,
		Timestamp:   filterPostTimestamps}, nil
}
func (s *Server) PostToServer(ctx context.Context, in *PostData) (*PostDataResponse, error) {

	resp, err := http.Get("http://127.0.0.1:12380/posts")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	// we initialize our Users array
	var posts []Post

	json.Unmarshal(body, &posts)

	posts = append(posts, Post{
		PostID:      int(in.Postid),
		Title:       in.Title,
		Author:      in.Author,
		Description: in.Description,
		Timestamp:   in.Timestamp,
	})
	byteValue, err := json.Marshal(posts)
	if err != nil {
		fmt.Println(err)
	}
	cmd := exec.Command("curl", "-L", "http://127.0.0.1:12380/posts", "-XPUT", "-d "+string(byteValue))
	cmd.Run()

	return &PostDataResponse{Success: true}, nil
}
