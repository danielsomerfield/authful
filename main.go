package main

import "github.com/danielsomerfield/authful/server"

func main() {
	server := server.NewAuthServer(8080)
	server.Start()
	server.Wait()
}
