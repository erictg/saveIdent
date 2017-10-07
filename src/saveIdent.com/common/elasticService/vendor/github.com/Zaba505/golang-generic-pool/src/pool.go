package genericPool

import (
	"time"
	"reflect"
	"io"
	"encoding/json"
	"sync"
)


type NoFishieError struct {
	Message string
}

func (noFishieErr NoFishieError) Error() string {
	return noFishieErr.Message
}

type FishiesPool interface {
	GetFishie() interface{}
	PutFishieBack(fishie interface{})
	RemoveFishie(fishieCheck interface{})
	Do(fn interface{})
	Len() int
}

type GenericPool struct {
	fishies chan interface{}
	mu sync.Mutex
}

func New(numOfFishies int, fishieType interface{}, fishieFieldGenerator func(fishieNum int) map[int]interface{}) *GenericPool {

	var wg sync.WaitGroup

	var mu sync.Mutex
	fishies := make(chan interface{}, numOfFishies)

	for i := 0; i < numOfFishies; i++ {

		fishieFields := fishieFieldGenerator(i)

		// This dynamically creates a user defined fishie type from the generic interface they passed in
		emptyFishie := reflect.New(reflect.TypeOf(fishieType)).Elem()

		if fishieFields != nil {
			for j := 0; j < emptyFishie.NumField(); j++ {
				if fieldVal, ok := fishieFields[j + 1]; ok {
					emptyFishie.Field(j).Set(reflect.ValueOf(fieldVal))
				}
			}
		}

		wg.Add(1)

		go func(fishieInterface interface{}, fishieChan chan interface{}) {
			defer wg.Done()

			fishies <- fishieInterface

		}(emptyFishie.Interface(), fishies)

	}

	wg.Wait()

	return &GenericPool{fishies, mu}

}


// Gets a fishie from the pool
// It has a 2 second timeout if there isn't any free clients
// Therefore you should always check if you caught a fish or not i.e. if fishie != nil
func (gP *GenericPool) GetFishie() interface{} {
	select {
		case fishie := <-gP.fishies:
			return fishie
		case <-time.After(time.Second * 2):
			return nil
	}
}


// Puts a fishie back into the pool
func (gP *GenericPool) PutFishieBack(fishie interface{}) {

	go func(fish interface{}) {
		gP.fishies <- fish
	}(fishie)

}


// Removes a fishie from the pool
// Must be a function which takes in a fishie and returns a bool
func (gP *GenericPool) RemoveFishie(fishieCheck interface{}) {
	go func() {

		fishieCheckFunc := reflect.ValueOf(fishieCheck)

		for fishie := range gP.fishies {

			result := fishieCheckFunc.Call([]reflect.Value{ reflect.ValueOf(fishie) })

			if res, ok := result[0].Interface().(bool); ok && !res {
				gP.PutFishieBack(fishie)
			} else {
				break
			}

		}
	}()
}


// Runs an action on all the fishies in the pool at the time it gets called
// Must be a void function which takes in a fishie
func (gP *GenericPool) Do(fn interface{}) {

	gP.mu.Lock()

	numOfFishies := gP.Len()

	tempFishiesChan := make(chan interface{}, numOfFishies)
	var tempFishiesSlice []interface{}

	doFunc := reflect.ValueOf(fn)

	for i := 0; i < numOfFishies; i++ {

		fishie := <-gP.fishies
		doFunc.Call([]reflect.Value{ reflect.ValueOf(fishie) })

		tempFishiesSlice = append(tempFishiesSlice, fishie)

	}

	go func(tempFishChan chan interface{}, tempFishSlice []interface{}) {
		for _, tempFishie := range tempFishiesSlice {
			tempFishiesChan <- tempFishie
		}

		close(tempFishiesChan)
	}(tempFishiesChan, tempFishiesSlice)

	for fishie := range tempFishiesChan {
		gP.PutFishieBack(fishie)
	}

	gP.mu.Unlock()

}


// Returns the number of fishies in the pool
func (gP *GenericPool) Len() int {
	return len(gP.fishies)
}


// Helper function for visualizing fishies if necssary
// Primary use case is for debugging your code
func PrettyFishie(fishie interface{}, w io.Writer) {

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	encoder.Encode(&fishie)

}