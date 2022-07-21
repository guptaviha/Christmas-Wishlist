package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"welcome-app/feed"
	"welcome-app/login"
	"welcome-app/user"

	"google.golang.org/grpc"
)

//Create a struct that holds information to be displayed in our HTML file
type CurrUser struct {
	Username    string
	Followed    []string
	notFollowed []string
}

//Go application entrypoint
func main() {

	curruser := CurrUser{Username: "Guest"}
	currposts := []feed.Post{}

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	// landing page endpoint
	http.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {

		templates := template.Must(template.ParseFiles("templates/welcome.html"))

		if err := templates.ExecuteTemplate(w, "welcome.html", curruser); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// login page endpoint
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		templates := template.Must(template.ParseFiles("templates/login.html"))

		if err := templates.ExecuteTemplate(w, "login.html", curruser); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// verify-login endpoint
	http.HandleFunc("/verify-login", func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		var ip_username = r.Form["username"][0]
		var ip_password = r.Form["password"][0]

		var conn *grpc.ClientConn
		conn, err := grpc.Dial(":9000", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %s", err)
		}
		defer conn.Close()

		l := login.NewAuthServiceClient(conn)

		response, err := l.Authenticate(context.Background(), &login.LoginDetails{Username: ip_username, Password: ip_password})
		if err != nil {
			log.Fatalf("Error when calling Authenticate: %s", err)
		}
		log.Printf("Response from server: Name: %s", response.Name)

		if response.Done {
			curruser.Username = response.Name
			templates := template.Must(template.ParseFiles("templates/welcome.html"))

			if err := templates.ExecuteTemplate(w, "welcome.html", curruser); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			templates := template.Must(template.ParseFiles("templates/login.html"))

			if err := templates.ExecuteTemplate(w, "login.html", curruser); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})

	// signup page endpoint
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {

		templates := template.Must(template.ParseFiles("templates/signup.html"))

		if err := templates.ExecuteTemplate(w, "signup.html", curruser); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// add-user endpoint
	http.HandleFunc("/add-user", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		var ip_username = r.Form["username"][0]
		var ip_name = r.Form["name"][0]
		var ip_password = r.Form["password"][0]
		var ip_confirmPassword = r.Form["confirm_password"][0]

		var conn *grpc.ClientConn
		conn, err := grpc.Dial(":9000", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %s", err)
		}
		defer conn.Close()

		u := user.NewUserServiceClient(conn)

		response, err := u.SignUpUser(context.Background(), &user.SignUpRequest{
			Username:        ip_username,
			Password:        ip_password,
			Name:            ip_name,
			ConfirmPassword: ip_confirmPassword,
		})
		if err != nil {
			log.Fatalf("Error when calling SignUpUser: %s", err)
		}
		if !response.Success {
			templates := template.Must(template.ParseFiles("templates/signup.html"))

			if err := templates.ExecuteTemplate(w, "signup.html", curruser); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			templates := template.Must(template.ParseFiles("templates/login.html"))

			if err := templates.ExecuteTemplate(w, "login.html", curruser); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

	})

	// feed page endpoint
	http.HandleFunc("/feed", func(w http.ResponseWriter, r *http.Request) {
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(":9000", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %s", err)
		}
		defer conn.Close()

		f := feed.NewFeedServiceClient(conn)

		response, err := f.GetFeed(context.Background(), &feed.FeedRequest{Username: curruser.Username})
		if err != nil {
			log.Fatalf("Error when calling GetFeed: %s", err)
		}
		currposts = []feed.Post{}
		log.Println(curruser)
		log.Println(response)
		for i, _ := range response.Postid {
			currposts = append(currposts, feed.Post{
				PostID:      int(response.Postid[i]),
				Title:       response.Title[i],
				Author:      response.Author[i],
				Description: response.Description[i],
				Timestamp:   response.Timestamp[i],
			})
		}

		templates := template.Must(template.ParseFiles("templates/feed.html"))
		if err := templates.ExecuteTemplate(w, "feed.html", currposts); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	// add-post endpoint
	http.HandleFunc("/add-post", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		var ip_title = r.Form["title"][0]
		var ip_description = r.Form["description"][0]

		var conn *grpc.ClientConn
		conn, err := grpc.Dial(":9000", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %s", err)
		}
		defer conn.Close()

		f := feed.NewFeedServiceClient(conn)
		response, err := f.PostToServer(context.Background(), &feed.PostData{
			Postid:      0,
			Title:       ip_title,
			Description: ip_description,
			Author:      curruser.Username,
			Timestamp:   time.Now().Format("01-02-2006 15:04:05"),
		})
		if response.Success {
			templates := template.Must(template.ParseFiles("templates/feed.html"))

			if err := templates.ExecuteTemplate(w, "feed.html", currposts); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})

	// users page endpoint
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {

		var conn *grpc.ClientConn
		conn, err := grpc.Dial(":9000", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %s", err)
		}
		defer conn.Close()

		u := user.NewUserServiceClient(conn)

		response, err := u.GetFollowing(context.Background(), &user.FollowerRequest{
			Username: curruser.Username,
		})

		templates := template.Must(template.ParseFiles("templates/users.html"))

		if err := templates.ExecuteTemplate(w, "users.html", response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// unfollow endpoint
	http.HandleFunc("/unfollow", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		var ip_Newuser = r.Form["username"][0]

		var conn *grpc.ClientConn
		conn, err := grpc.Dial(":9000", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %s", err)
		}
		defer conn.Close()

		u := user.NewUserServiceClient(conn)

		res, err := u.UpdateFollower(context.Background(), &user.UpdateFollowersRequest{
			Username: curruser.Username,
			Newuser:  ip_Newuser,
			IsFollow: false,
		})
		if err != nil {
			log.Fatalf("Could not UpdateFollower: %s", err)
		}

		if res.Success {
			templates := template.Must(template.ParseFiles("templates/welcome.html"))

			if err := templates.ExecuteTemplate(w, "welcome.html", curruser); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			response, err := u.GetFollowing(context.Background(), &user.FollowerRequest{
				Username: curruser.Username,
			})
			if err != nil {
				log.Fatalf("Could not get following: %s", err)
			}
			templates := template.Must(template.ParseFiles("templates/users.html"))

			if err := templates.ExecuteTemplate(w, "users.html", response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

	})

	// follow endpoint
	http.HandleFunc("/follow", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		var ip_Newuser = r.Form["username"][0]

		var conn *grpc.ClientConn
		conn, err := grpc.Dial(":9000", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %s", err)
		}
		defer conn.Close()

		u := user.NewUserServiceClient(conn)

		res, err := u.UpdateFollower(context.Background(), &user.UpdateFollowersRequest{
			Username: curruser.Username,
			Newuser:  ip_Newuser,
			IsFollow: true,
		})
		if err != nil {
			log.Fatalf("Could not UpdateFollower: %s", err)
		}

		if res.Success {
			templates := template.Must(template.ParseFiles("templates/welcome.html"))

			if err := templates.ExecuteTemplate(w, "welcome.html", curruser); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			response, err := u.GetFollowing(context.Background(), &user.FollowerRequest{
				Username: curruser.Username,
			})
			if err != nil {
				log.Fatalf("Could not get following: %s", err)
			}
			templates := template.Must(template.ParseFiles("templates/users.html"))

			if err := templates.ExecuteTemplate(w, "users.html", response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

	})

	// signout page endpoint
	http.HandleFunc("/signout", func(w http.ResponseWriter, r *http.Request) {

		// Reset user to guest
		curruser.Username = "Guest"
		currposts = []feed.Post{}

		// Redirect to homepage
		templates := template.Must(template.ParseFiles("templates/welcome.html"))

		if err := templates.ExecuteTemplate(w, "welcome.html", curruser); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	fmt.Println(http.ListenAndServe(":8080", nil))
}
