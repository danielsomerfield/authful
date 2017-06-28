package handlers

import (
	"fmt"
	"time"
	"net/http"
	"log"
)

type TestServer struct {
	httpServer http.Server
}

func init() {
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"status\":\"ready\"}"))
	})
}

func RunTestServer(pattern string, handler func(w http.ResponseWriter, r *http.Request)) TestServer {
	port := 8080
	httpServer := http.Server{Addr: fmt.Sprintf(":%v", port)}
	http.HandleFunc(pattern, handler)
	go startServer(httpServer)
	waitForServer()
	return TestServer{
		httpServer: httpServer,
	}
}

func (server TestServer) Shutdown() {
	server.httpServer.Shutdown(nil)
}

func startServer(httpServer http.Server) {
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %+v", err)
	}
}

func waitForServer() error {
	var err error = nil
	var resp *http.Response = nil

	for i := 0; i < 50; i++ {
		resp, err = http.Get("http://localhost:8080/ready")

		if err == nil {
			if resp.StatusCode == 200 {
				return nil
			} else {
				err = fmt.Errorf("Expected status code 200 but was %s", resp.StatusCode)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	return err
}