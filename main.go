package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	BASE_URL           = "https://bakusai.com/"
	BASE_THREADTOP_URL = "https://bakusai.com/thr_tl/acode=3/ctgid=136/bid=2027/"
)

var (
	get_ichiran_count int = 0
	thread_urls       []string
	comments          = make(map[int]string)
)

func main() {

	// var threadUrlpool []string
	next_url := ThreadnextURL(BASE_THREADTOP_URL)    //・スレッドの2ページ目のURLを取得
	thread_urls := Threadichiran(BASE_THREADTOP_URL) //スレッド一覧の1ページ目のURLを全部取得
	fmt.Println("[main]:", thread_urls, "スレッドは", len(thread_urls), "件です")

	// スレッドのURLを一覧ページから取得する
	for {
		url := ThreadnextURL(BASE_URL + next_url)
		time.Sleep(1)
		thread_urls := Threadichiran(BASE_URL + url)
		time.Sleep(1)
		fmt.Println("[main] forループ中", thread_urls, "スレッドは", len(thread_urls), "件です")

		if url == "" {
			break
		} else if get_ichiran_count >= 4 {
			break
		}
	}

	// スレッドURLへアクセスして中身を取得する

}

//スレッド一覧からURLを取得する
func Threadichiran(BASE_THREADTOP_URL string) []string {

	res, err := http.Get(BASE_THREADTOP_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf(res.Status)
	}

	response, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	threads := response.Find("div.lSideColumn")

	// next_url, ok := response.Find(".paging_nextlink_btn > a").Attr("href")
	// fmt.Println(next_url)

	// if ok {
	// 	fmt.Println("[next!!!!!!!!!!!!!!]:\n", next_url, ok)
	// 	get_ichiran_count++
	// }

	threads.Find("a").Each(func(index int, item *goquery.Selection) {

		href, _ := item.Attr("href")

		//fmt.Println("[Threadichiran]:", href)
		thread_urls = append(thread_urls, href)
	})

	return thread_urls
}

func Thread_text() {

	res, err := http.Get("https://bakusai.com/thr_res/acode=3/ctgid=136/bid=2027/tid=10474189/tp=1/")
	// TODO ↑のURLにとってきたスレッドURLを入れる

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf(res.Status)
	}
	response, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// コメント取得

	comment := response.Find(".article")
	comment.Each(func(index int, item *goquery.Selection) {

		comment := item.Text()

		fmt.Println(comment)
		comments[index] = comment
	})
}

// //次へのリンクを取得する
// func ThreadnextUrl(current string) (url string) {
// 	fmt.Printf("############ url: %s\n", current)
// 	doc, err := goquery.NewDocument(current)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var exist bool
// 	doc.Find("######").Each(func(_ int, s *goquery.Selection) {

// 		url, exist = s.Attr("href")
// 	})
// 	if !exist {
// 		return ""
// 	}
// 	return
// }

func ThreadnextURL(BASE_THREADTOP_URL string) string {

	res, err := http.Get(BASE_THREADTOP_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf(res.Status)
	}

	response, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	next_url, exist := response.Find(".paging_nextlink_btn > a").Attr("href")

	if exist {
		fmt.Println("[ThradnextURL]", next_url)
		get_ichiran_count++
		fmt.Println("カウントは", get_ichiran_count, "です")
	}

	return next_url
}
