package main

import (
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/gliderlabs/ssh"
)

type Tunnel struct {
	w      io.Writer
	donech chan struct{}
}

type TunnelsMap struct {
	mu      sync.RWMutex
	tunnels map[int]chan Tunnel
}

func NewTunnelsMap() *TunnelsMap {
	return &TunnelsMap{
		tunnels: make(map[int]chan Tunnel),
	}
}

func (tm *TunnelsMap) Put(id int, tunnel chan Tunnel) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.tunnels[id] = tunnel
}

func (tm *TunnelsMap) Get(id int) (chan Tunnel, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tunnel, ok := tm.tunnels[id]
	return tunnel, ok
}

var tunnels = NewTunnelsMap()

func main() {
	go startHTTPServer()

	ssh.Handle(handleSSHSession)

	log.Fatal(ssh.ListenAndServe(":2222", nil))
}

func startHTTPServer() {
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idstr)

	tunnel, ok := tunnels.Get(id)
	if !ok {
		http.Error(w, "tunnel not found", http.StatusNotFound)
		return
	}

	donech := make(chan struct{})
	tunnel <- Tunnel{
		w:      w,
		donech: donech,
	}

	<-donech
}

func handleSSHSession(s ssh.Session) {
	id := rand.Intn(math.MaxInt)
	tunnelCh := make(chan Tunnel)
	tunnels.Put(id, tunnelCh)

	log.Println("tunnel ID ->", id)

	tunnel, _ := tunnels.Get(id)
	t := <-tunnel

	log.Println("tunnel is ready")

	if _, err := io.Copy(t.w, s); err != nil {
		log.Fatalln(err)
	}

	close(t.donech)
	s.Write([]byte("we are done!"))
}
