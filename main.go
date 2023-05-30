package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type UserProfile struct {
	ID       int
	Comments []string
	Likes    int
	Friends  []int
}

func main() {
	start := time.Now()
	userProfile, err := handleGetUserProfile(11)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("userProfile: ", userProfile)
	fmt.Println("fetching the user profile took: ", time.Since(start))
}

type Response struct {
	data any
	err  error
}

func handleGetUserProfile(id int) (*UserProfile, error) {
	var (
		resCh = make(chan Response, 3)
		wg    = sync.WaitGroup{}
	)
	// we are making 3 requests inside their own goroutine
	go getComments(id, resCh, &wg)
	go getLikes(id, resCh, &wg)
	go getFriends(id, resCh, &wg)

	//add the number of requests which is 3
	wg.Add(3)

	//block until our wg 'counter' is 0 and then unblock
	wg.Wait()
	close(resCh)

	userProfile := &UserProfile{}
	for response := range resCh {

		if response.err != nil {
			return nil, response.err
		}

		switch msg := response.data.(type) {
		case int:
			userProfile.Likes = msg
		case []int:
			userProfile.Friends = msg
		case []string:
			userProfile.Comments = msg
		}
	}
	return userProfile, nil
}

func getComments(id int, resCh chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 200)
	comments := []string{
		"This looks good",
		"LGTM!",
		"Insightful piece!",
	}
	resCh <- Response{
		data: comments,
		err:  nil,
	}
	wg.Done()
}

func getLikes(id int, resCh chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 200)
	resCh <- Response{
		data: 33,
		err:  nil,
	}
	wg.Done()
}

func getFriends(id int, resCh chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 100)
	resCh <- Response{
		data: []int{11, 34, 78, 455},
		err:  nil,
	}
	wg.Done()
}
