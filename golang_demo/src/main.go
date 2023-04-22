// 第一种解法
// 第一种解法是将URL放入通道中，先用10个协程池中的协程进行并行地获取用户的评论数据，然后并将提取到的电子邮件地址发送到另一个通道中；这种解法是使用了WaitGroup和通道关闭来确保所有协程都已完成并停止主线程的工作的
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Comment struct {
	PostId int    `json:"postId"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func main() {
	urls := make(chan string, 100)
	comments := make(chan []Comment, 100)
	emails := make(chan string, 500)

	for i := 1; i <= 100; i++ {
		urls <- fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d/comments", i)
	}
	close(urls)

	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}
				defer resp.Body.Close()

				var commentsData []Comment
				err = json.NewDecoder(resp.Body).Decode(&commentsData)
				if err != nil {
					fmt.Println(err)
					continue
				}

				comments <- commentsData
			}
		}()
	}

	go func() {
		wg.Wait()
		close(comments)
	}()

	for c := range comments {
		for _, comment := range c {
			emails <- comment.Email
		}
	}

	close(emails)

	for email := range emails {
		fmt.Println(email)
	}
}

// 第二种解法
// 这种方法使用的是Go语言的并发机制，每个URL都是单独的在协程中获取。然后使用WaitGroup来确保所有协程都已完成，最后在完成后关闭了channel通道以停止主线程
// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"sync"
// )

// type Comment struct {
// 	PostId int    `json:"postId"`
// 	Id     int    `json:"id"`
// 	Name   string `json:"name"`
// 	Email  string `json:"email"`
// 	Body   string `json:"body"`
// }

// func main() {
// 	var wg sync.WaitGroup
// 	emails := make(chan string, 500)

// 	for i := 1; i <= 100; i++ {
// 		url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d/comments", i)
// 		wg.Add(1)
// 		go func(url string) {
// 			defer wg.Done()
// 			resp, err := http.Get(url)
// 			if err != nil {
// 				fmt.Println(err)
// 				return
// 			}
// 			defer resp.Body.Close()

// 			var comments []Comment
// 			err = json.NewDecoder(resp.Body).Decode(&comments)
// 			if err != nil {
// 				fmt.Println(err)
// 				return
// 			}

// 			for _, c := range comments {
// 				emails <- c.Email
// 			}
// 		}(url)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(emails)
// 	}()

// 	for e := range emails {
// 		fmt.Println(e)
// 	}
// }
