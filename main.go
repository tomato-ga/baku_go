package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	// TODO mongoDBと接続→structをJSONにして保存できるか確認する。ソートさせる
)

const (
	BASE_URL           = "https://bakusai.com/"
	BASE_THREADTOP_URL = "thr_tl/acode=3/ctgid=136/bid=2027/"
)

var (
	get_ichiran_count  int = 0
	maxThreadPageCount int = 1
	thread_urls        []string
	next_url           string
)

type thread_info struct {
	Shopname string
	Url      string
	Text     []string
}

func main() {

	now := time.Now()
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
		comm, _, _ := ThreadGetText(BASE_URL + u)
		time.Sleep(1)
		np := ThreadGetNext(BASE_URL + u)
		time.Sleep(1)
		fmt.Println("[main]", comm)
		if comm == nil {
			break
		}

		for {
			comm, shop_title, thread_parse_url := ThreadGetText(BASE_URL + np)
			shop_info := thread_info{Shopname: shop_title, Url: thread_parse_url, Text: comm}
			fmt.Println("[shop_info]:", shop_info)
			np = ThreadGetNext(BASE_URL + np) //np ネクストページの略
			if np == "" {
				break
			}
		}
	}
	fmt.Printf("経過: %vms\n", time.Since(now).Seconds())
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
		thread_urls = append(thread_urls, href)
	})

	return thread_urls
}

func ThreadGetText(thread_parse_url string) ([]string, string, string) {
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
	var comm []string // TODO #1からソートさせる #をとって数字をintにしてソート？ ブログ投稿するときには#入れればいい

	shop_title := response.Find(".title_thr_wrap ").Text()
	comment := response.Find(".article")

	comment.Each(func(index int, item *goquery.Selection) {
		comment := item.Text()
		comm = append(comm, comment)
		fmt.Println(comm)
	})

	// 正規表現でソートさせる
	for _, cc := range comm {
		var m_comm []map[int]string
		r_number := regexp.MustCompile(`(\d{1,4})`)                                     // #を含む数字だけ取得する
		r_time := regexp.MustCompile(`([0-9]{4}/[0-9]{2}/[0-9]{2}\ [0-9]{2}:[0-9]{2})`) // 2022/06/06の形式にマッチする日付だけ取得する
		res_number := r_number.FindString(cc)
		res_time := r_time.FindString(cc)

		sub := regexp.MustCompile(` `)
		split := sub.Split(cc, -1)
		r_time_delete := regexp.MustCompile(`([0-9]{2}:[0-9]{2})`)
		r_tokumei_delete := regexp.MustCompile(`(\[匿名さん\])`)
		res_timedelete_text := r_time_delete.ReplaceAllString(split[1], "")
		res_time_tokumeidelete_text := r_tokumei_delete.ReplaceAllString(res_timedelete_text, "")
		fmt.Println(split)
		fmt.Println(res_time_tokumeidelete_text)
		fmt.Println(res_number, res_time)
		res_mix := res_time + "," + res_time_tokumeidelete_text

		res_number_convert, _ := strconv.Atoi(res_number) // str→int変換
		fmt.Println(res_mix)
		fmt.Println(res_number_convert)

		//TODO m_comm = append(m_comm, map[res_number_convert]res_mix)
		fmt.Println(m_comm)

	}

	return comm, shop_title, thread_parse_url
}

func ThreadGetNext(thread_parse_url string) string {
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
		fmt.Println("[ThreadGetNext]", thread_next_page)
	} else if exist == false {
		return ""
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
