package main 

import (
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
)

func NewServer(cache *Cache, port string) *Server {
	server := Server{cache: cache}

	server.listen(port)

	return &server
}

type Server struct {
	cache *Cache
	server *http.Server
}

func (s *Server) Close() {
	s.server.Shutdown(nil)
}

func (s *Server) listen(port string) {
	log.Println("Starting server")

	http.HandleFunc("/remove", s.apiRouter)
	
	s.server = &http.Server{Addr: ":"+port}
	
	go func() {
		err := s.server.ListenAndServe()

		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	log.Println("Server listening on port " + port)
}

type deleteRequest struct { 
	Key string `json:"key"`
}

func (s *Server) apiRouter(res http.ResponseWriter, req *http.Request) {
	log.Println("Received request")

	defer req.Body.Close()

	str, _ := ioutil.ReadAll(req.Body)

	var data deleteRequest
	
	json.Unmarshal(str, &data)

	if !keyValid(data.Key) {
		log.Println("Invalid key received: " + data.Key)
		return
	}

	log.Println("Sending delete: " + data.Key)
	s.cache.Events <- CacheEvent{ Key: data.Key, Op: DELETE }

}

func keyValid(key string) bool {
	return len(key) > 6
}