package main

import (
	"fmt"
	"sort"
	"sync"
)

func RunPipeline(cmds ...cmd) {
	if len(cmds) == 0 {
		return
	}

	channels := make([]chan interface{}, len(cmds)+1)
	for i := range channels {
		channels[i] = make(chan interface{})
	}

	var wg sync.WaitGroup

	for i, command := range cmds {
		wg.Add(1)
		go func(c cmd, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)
			c(in, out)
		}(command, channels[i], channels[i+1])
	}

	close(channels[0])
	wg.Wait()
}

func SelectUsers(in, out chan interface{}) {
	seen := make(map[uint64]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for email := range in {
		wg.Add(1)
		go func(e interface{}) {
			defer wg.Done()
			emailStr, ok := e.(string)
			if !ok {
				return
			}

			user := GetUser(emailStr)

			mu.Lock()
			if !seen[user.ID] {
				seen[user.ID] = true
				mu.Unlock()
				out <- user
			} else {
				mu.Unlock()
			}
		}(email)
	}

	wg.Wait()
}

func SelectMessages(in, out chan interface{}) {
	batchSize := 2
	batch := make([]User, 0, batchSize)
	var mu sync.Mutex
	var wg sync.WaitGroup

	processBatch := func(userBatch []User) {
		defer wg.Done()
		if len(userBatch) == 0 {
			return
		}

		msgIds, err := GetMessages(userBatch...)
		if err != nil {
			return
		}

		for _, msgID := range msgIds {
			out <- msgID
		}
	}

	for u := range in {
		user, ok := u.(User)
		if !ok {
			continue
		}

		mu.Lock()
		batch = append(batch, user)

		if len(batch) == batchSize {
			wg.Add(1)
			go processBatch(append([]User(nil), batch...))
			batch = batch[:0]
		}
		mu.Unlock()
	}

	if len(batch) > 0 {
		wg.Add(1)
		go processBatch(batch)
	}

	wg.Wait()
}

func CheckSpam(in, out chan interface{}) {
	semaphore := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for msgIDInterface := range in {
		msgID, ok := msgIDInterface.(MsgID)
		if !ok {
			continue
		}

		wg.Add(1)
		go func(id MsgID) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			hasSpam, err := HasSpam(id)
			if err != nil {
				hasSpam = true
			}

			out <- MsgData{
				ID:      id,
				HasSpam: hasSpam,
			}
		}(msgID)
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var results []MsgData

	for dataInterface := range in {
		data, ok := dataInterface.(MsgData)
		if !ok {
			continue
		}
		results = append(results, data)
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].HasSpam != results[j].HasSpam {
			return results[i].HasSpam && !results[j].HasSpam
		}
		return results[i].ID < results[j].ID
	})

	for _, data := range results {
		spamStr := "false"
		if data.HasSpam {
			spamStr = "true"
		}
		out <- fmt.Sprintf("%s %d", spamStr, data.ID)
	}
}
