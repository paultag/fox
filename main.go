package main

import (
	"log"
	"os"

	"github.com/armon/go-socks5"
	"gopkg.in/yaml.v3"
)

type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type AuthBackend struct {
	Users []User `yaml:"users"`
}

func (abe AuthBackend) Valid(user, password string) bool {
	for _, u := range abe.Users {
		if u.Username == user {
			return u.Password == password
		}
	}
	return false
}

type Config struct {
	Listen string      `yaml:"listen"`
	Auth   AuthBackend `yaml:"auth"`
}

func main() {
	config := Config{}

	fd, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	if err := yaml.NewDecoder(fd).Decode(&config); err != nil {
		panic(err)
	}
	fd.Close()

	conf := &socks5.Config{
		Credentials: config.Auth,
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	addr := config.Listen
	log.Printf("Listening on %s", addr)
	if err := server.ListenAndServe("tcp", addr); err != nil {
		panic(err)
	}
}
