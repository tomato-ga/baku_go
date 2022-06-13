package thread

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/PuerkitoBio/goquery"
// )

// func Thread_text() {

// 	res, err := http.Get("https://bakusai.com/thr_res/acode=3/ctgid=136/bid=2027/tid=10474189/tp=1/")
// 	// TODO ↑のURLにとってきたスレッドURLを入れる

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != 200 {
// 		log.Fatalf(res.Status)
// 	}
// 	response, err := goquery.NewDocumentFromReader(res.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// コメント取得

// 	comment := response.Find(".article")
// 	comment.Each(func(index int, item *goquery.Selection) {

// 		comment := item.Text()

// 		fmt.Println(comment)
// 		comments[index] = comment
// 	})
// }
