package bestof

import (
	"fmt"
	"net/http"
	"net/url"
)

//TODO:
// - #semaphore concept: See http://www.golangpatterns.info/concurrency/semaphores
// - #Map of 'backends' to semaphore<empty struct>
// - #Random number provider
// - #Ability to insert random number provider for testing purposes
// - Code for request proxying
// - Increment and decrement semaphore appropriately during proxy
// - If number of backends is zero, return 502
// - #If number of backends is one, follow next immediately, do not follow algorithm
// - #N should never be zero or greater than number of backends, default to two or number of backends
// - Configuration should allow for selection of N, should default to two if not present
// - #If N is equal to number of backends, fall back to JSQ, do not go through selection
// - #If N is less than number of backends, randomly select N backends and choose the least loaded

//Balancer is a bookkeeping struct
type Balancer struct {
	balancees       map[url.URL]Semaphore
	randomGenerator RandomInt
	next            http.Handler
	choices         int
	keys            []url.URL
}

//constructor must:
//- set up keys to be keys of balancees Map
//- set up randomGenerator to choose between 0 and

func (b *Balancer) nextServer() (*url.URL, error) {
	//Special case: If balancees are nil or empty, return an error.
	if b.balancees == nil || len(b.balancees) == 0 {
		return nil, fmt.Errorf("Number of balancees is zero, cannot handle")
	}
	//Special case: If balancees is 1, there is no need to balance
	if len(b.balancees) == 1 {
		for key := range b.balancees {
			return &key, nil
		}
	}
	var normalizedChoices = b.choices
	//Special case: If choices is <= 0, default to 2
	if normalizedChoices <= 0 {
		normalizedChoices = 2
	}
	//Special case: If choices > number of balancees, default to number of backends
	if normalizedChoices > len(b.balancees) {
		normalizedChoices = len(b.balancees)
	}
	var potentialChoices = []url.URL{}

	//shuffle keys, we'll choose the first N from the shuffled result
	for i := range b.keys {
		j, _ := b.randomGenerator.nextInt(0, i+1)
		b.keys[i], b.keys[j] = b.keys[j], b.keys[i]
	}

	if normalizedChoices == len(b.balancees) {
		potentialChoices = b.keys
	} else {
		for i := 0; i < normalizedChoices; i++ {
			potentialChoices = append(potentialChoices, b.keys[i])
		}
	}

	var bestChoice *url.URL
	var leastConns = -1
	for _, key := range potentialChoices {
		if leastConns == -1 {
			leastConns = b.balancees[key].length()
			bestChoice = &key
			continue
		}
		if leastConns > b.balancees[key].length() {
			leastConns = b.balancees[key].length()
			bestChoice = &key
		}
	}
	return bestChoice, nil
}
