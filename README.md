# Christmas Wishlist WebApp

An academic project for Distributed Systems (CS-GY-9223) at NYU Tandon with Prof Gustavo Sandoval

Objective: To develop a distributed and reliable backend in support of a simple social media application. 

Our application is a wishlist tool where users can login, create posts, follow other users, and view posts of the people they follow. The webserver, written in Go, interacts with the client using gRPCs. [CoreOS](https://github.com/etcd-io/etcd/tree/main/raft), an open source Raft implementation, is used in the backend to provide consensus for the raft nodes that are spun up locally. 

Please see the [Project Prompt](https://github.com/guptaviha/Christmas-Wishlist/blob/main/Project_Prompt.pdf) for more details.

## Instructions To Run App

Install go [Ref](https://go.dev/doc/install)
```brew install go```

Install goreman [Ref](https://github.com/mattn/goreman)
```go install github.com/mattn/goreman@latest```

Follow raftexample [Ref](https://github.com/etcd-io/etcd/tree/main/contrib/raftexample)

1. Setup PATH variables
```bash
    export GOPATH=<dir>
``` 
For example: 
```bash
    export GOPATH=/Users/vihagupta/Documents/code/Christmas-Wishlist/
```

```bash
    export PATH=$GOPATH/welcome-app/bin:$PATH
```

2. Check that goreman is setup correctly
```bash
    which goreman
```

3. Start goreman
```bash
    cd <dir>/welcome-app/src/go.etcd.io/etcd/contrib/raftexample

    go build -o raftexample

    goreman start
```

4. Add setup data - optional
```bash
    cd <dir>/welcome-app

    bash setup.sh
```

5. Start the app
```bash
    cd <dir>/welcome-app

    go run main.go
    go run server.go
```

## Instructions To Run Tests

Must run setup.sh first before running the following tests.

```bash
    cd welcome-app/login
    go test -v

    cd welcome-app/user
    go test -v

    cd welcome-app/feed
    go test -v    
```

## To Access the App

Open the browser and go to 

```
    http://localhost:8080/welcome
```

## Contributors:

<p float="left">

<a href="https://github.com/sidistic">
    <img src="https://github.com/sidistic.png?size=50" alt="Sidharth Sagar" width="50">
</a>

<a href="https://github.com/guptaviha">
    <img src="https://github.com/guptaviha.png?size=50" alt="Viha Gupta" width="50">
</a>

</p>
