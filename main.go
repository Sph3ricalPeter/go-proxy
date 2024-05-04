package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Name    string
	Host    string
	Port    uint16
	Targets map[string]string
}

func LoadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	dec := json.NewDecoder(file)

	var config Config
	err = dec.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

type Proxy struct {
	Name    string
	Targets map[string]string
}

func NewProxy(name string) *Proxy {
	return &Proxy{Name: name, Targets: make(map[string]string)}
}

func (p *Proxy) AddTarget(key, target string) {
	p.Targets[key] = target
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check if proxy target exists for given request host
	target, ok := p.Targets[r.URL.Path]
	if !ok {
		http.Error(w, "No proxy target for path", http.StatusNotFound)
		return
	}

	req, err := http.NewRequest("GET", target, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//  copy headers over to the new request
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	req.Header.Add("X-Proxied-By", p.Name)

	// send a GET request to the target
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// copy the proxied GET response body to this response writer
	io.Copy(w, resp.Body)
}

func main() {
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	p := NewProxy(config.Name)
	for host, target := range config.Targets {
		p.AddTarget(host, target)
	}

	http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), p)
}
