package rest

import (
	"math/rand"
	"net/http"
)

func (s *Server) handleOrderDelivered(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	w.Header().Add("Content-Type", "text/plain")

	// todo unmock

	if rand.Int31n(2) > 0 {
		return []byte("true"), nil
	}

	return []byte("false"), nil
}
