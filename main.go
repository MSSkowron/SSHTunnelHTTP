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
		mu:      sync.RWMutex{},
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

	value, ok := tm.tunnels[id]
	if !ok {
		return value, false
	}

	return value, true
}

var tunnels = NewTunnelsMap()

func main() {
	go func() {
		http.HandleFunc("/", handleRequest)
		log.Fatalln(http.ListenAndServe(":3000", nil))
	}()

	ssh.Handle(func(s ssh.Session) {
		id := rand.Intn(math.MaxInt)
		tunnels.Put(id, make(chan Tunnel))

		log.Println("tunnel ID ->", id)

		t, _ := tunnels.Get(id)
		tunnel := <-t

		log.Println("tunnel is ready")

		if _, err := io.Copy(tunnel.w, s); err != nil {
			log.Fatalln(err)
		}

		close(tunnel.donech)
		s.Write([]byte("we are done!"))
	})

	log.Fatal(ssh.ListenAndServe(":2222", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idstr)

	tunnel, ok := tunnels.Get(id)
	if !ok {
		if _, err := w.Write([]byte("tunnel not found")); err != nil {
			log.Fatalln(err)
		}
	}

	donech := make(chan struct{})
	tunnel <- Tunnel{
		w:      w,
		donech: donech,
	}

	<-donech
}
