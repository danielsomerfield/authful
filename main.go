package main

import "github.com/danielsomerfield/authful/server"
import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func main() {
	file := flag.String("config", "./authful.yml",
		"Path to the authful config file. Defaults to ./authful.yml")
	flag.Parse()
	config, err := fromFile(file)

	if err == nil {
		server := server.NewAuthServer(config)
		server.Start()
		server.Wait()
	} else {
		log.Printf("Failed to parse config file: %+v\n", err)
	}
}

func fromFile(path *string) (*server.Config, error) {
	bytes, err := ioutil.ReadFile(*path)
	if err != nil {
		return nil, err
	}
	return fromBytes(bytes)
}

func fromBytes(bytes []byte) (*server.Config, error) {
	config := &server.Config{
	}
	yaml.Unmarshal(bytes, config)
	return config, nil
}
