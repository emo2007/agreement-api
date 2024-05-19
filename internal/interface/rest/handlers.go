package rest

import (
	"fmt"
	"math/rand"
	"net/http"
)

func (s *Server) handleOrderDelivered(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	w.Header().Add("Content-Type", "text/plain")

	// todo unmock

	return []byte(fmt.Sprintf("%d", rand.Int31n(2))), nil
}
