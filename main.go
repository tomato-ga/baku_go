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
	BASE_THREADTOP_URL = "thr_tl/acode=3/ctgid=136/bid=2027/"
)

var (
	get_ichiran_count  int = 0
	maxThreadPageCount int = 2
	thread_urls        []string
	comments           map[int]string
	next_url           string
)

func main() {

	// var threadUrlpool []string
	thread_urls := Threadichiran(BASE_URL + BASE_THREADTOP_URL)     //スレッド一覧の1ページ目のURLを全部取得
	next_url := ThreadichiranNextURL(BASE_URL + BASE_THREADTOP_URL) //・スレッドの2ページ目のURLを取得

	fmt.Println("[main]:", thread_urls, "スレッドは", len(thread_urls), "件です")

	// スレッドのURLを一覧ページから取得する
	for {
		thread_urls = Threadichiran(BASE_URL + next_url)
		time.Sleep(1)
		next_url = ThreadichiranNextURL(BASE_URL + next_url)
		time.Sleep(1)
		fmt.Println("[main] forループ中", thread_urls, "スレッドは", len(thread_urls), "件です")

		if next_url == "" {
			break
		} else if get_ichiran_count >= maxThreadPageCount {
			break
		}
	}

	// スレッドURLへアクセスして中身を取得する
	for _, u := range thread_urls {
		comments := ThreadGetText(BASE_URL + u)
		time.Sleep(1)
		np := ThreadGetNext(BASE_URL + u)
		time.Sleep(1)
		fmt.Println(np)
		fmt.Println(comments)

		for {
			comments = ThreadGetText(BASE_URL + np)
			np = ThreadGetNext(BASE_URL + np)
			if np == "" {
				break
			}
			fmt.Println(comments)
		}
	}

}

//スレッド一覧からURLを取得する
func Threadichiran(turl string) []string {

	res, err := http.Get(turl)
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

	threads.Find("a").Each(func(index int, item *goquery.Selection) {
		thread_title := item.Text()
		href, _ := item.Attr("href")
		fmt.Println(thread_title)

		//fmt.Println("[Threadichiran]:", href)
		thread_urls = append(thread_urls, href)
	})

	return thread_urls
}

func ThreadGetText(thread_parse_url string) (comments map[int]string) {

	res, err := http.Get(thread_parse_url)

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

	comments = make(map[int]string) // 初期化
	comment := response.Find(".article")
	comment.Each(func(index int, item *goquery.Selection) {
		comment := item.Text()
		comments[index] = comment
	})
	return comments
}

func ThreadGetNext(thread_parse_url string) (thread_next_page string) {
	res, err := http.Get(thread_parse_url)
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
	thread_next_page, exist := response.Find(".paging_nextlink_btn > a").Attr("href")

	if exist {
		fmt.Println("[ThradnextURL]", thread_next_page)
	}

	return thread_next_page
}

func ThreadichiranNextURL(nexts string) string {

	res, err := http.Get(nexts)
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
