package thread

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const domain string = "https://bakusai.com/"

var get_ichiran_count int = 0
var thread_urls = make(map[string]string)
var comments = make(map[int]string)

func Thread_ichiran() map[string]string {

	res, err := http.Get("https://bakusai.com/thr_tl/acode=3/ctgid=136/bid=2027/")

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

	nexts, ok := response.Find(".paging_nextlink_btn > a").Attr("href")
	fmt.Println(nexts)

	if ok {
		fmt.Println("[next!!!!!!!!!!!!!!]:\n", nexts, ok)
		//Nexts_ichiran(nexts)
		get_ichiran_count++
	}

	threads.Find("a").Each(func(index int, item *goquery.Selection) {

		href, _ := item.Attr("href")
		thread_title := item.Text()

		fmt.Println(thread_title)
		fmt.Println(href)
		thread_urls[thread_title] = href
	})

	return thread_urls
}

func Nexts_ichiran(nexts string) {

	next_url := domain + nexts
	fmt.Println(next_url)

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
