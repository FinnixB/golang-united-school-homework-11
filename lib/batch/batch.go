package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

var wg sync.WaitGroup

func worker(input chan int64, output chan user) {
	defer wg.Done()

	for id := range input {
		output <- getOne(id)
	}
}

// TODO: need to add code for errors processing
func getBatch(n int64, pool int64) (res []user) {
	input := make(chan int64, n)
	output := make(chan user, n)

	var i int64

	// run workers pool
	for i = 0; i < pool; i++ {
		wg.Add(1)
		go worker(input, output)
	}

	// add tasks for workers
	for i = 0; i < n; i++ {
		input <- i
	}
	close(input)

	wg.Wait()

	// here all workers done
	close(output)

	// store results
	for r := range output {
		res = append(res, r)
	}

	return res
}
