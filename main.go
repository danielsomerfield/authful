package main // import "github.com/danielsomerfield/authful"
import "github.com/danielsomerfield/authful/server"

func main() {
	server.NewAuthServer().Start()
}