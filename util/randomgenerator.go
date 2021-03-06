package util

import (
	"errors"
	"math/rand"
)

//RandomInt gives a nextInt from its range which may or may not actually be random
type RandomInt interface {
	NextInt(minimum int, maximum int) (int, error)
}

//GoRandom uses Go's internal random generator to return a random number between minimum and maximum
type GoRandom struct {
}

//TestingRandom gives the next value in its array as it's 'nextInt', looping back to the first on reaching the end
type TestingRandom struct {
	Values    []int
	index     int
	CallCount int
}

func (g GoRandom) NextInt(minimum int, maximum int) (int, error) {
	if maximum < minimum {
		return 0, errors.New("Illegal state: Minimum is greater than maximum")
	}
	if maximum < 0 || minimum < 0 {
		return 0, errors.New("Illegal state: Minimum and maximum must both be above zero")
	}
	var n = maximum - minimum
	return rand.Intn(n) + minimum, nil
}

func (t *TestingRandom) NextInt(minimum int, maximum int) (int, error) {
	//return 4, nil
	t.CallCount++

	if len(t.Values) == 0 {
		return 0, nil
	}

	if t.index > len(t.Values)-1 {
		t.index = 0
	}

	var value = t.Values[t.index]
	t.index++

	return value, nil
}
