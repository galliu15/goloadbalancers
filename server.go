package main

import (
	"net/http"

	"github.com/jangie/bestofnlb/bestof"
	"github.com/vulcand/oxy/forward"
)

//Test harness
func main() {
	var fwd, _ = forward.New()
	var bal = bestof.NewBalancer(
		[]string{"http://testa:8080", "http://testb:8080", "http://testc:8080"},
		bestof.GoRandom{},
		2,
		fwd,
	)
	http.Handle("/", bal)
	http.ListenAndServe(":8090", nil)
}
