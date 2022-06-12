package thread

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const domain string = "https://bakusai.com/"

var get_ichiran_count int = 0
var m = make(map[string]string)

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
		m[thread_title] = href
	})

	return m
}

func Nexts_ichiran(nexts string) {

	next_url := domain + nexts
	fmt.Println(next_url)

}
