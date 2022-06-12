package main

import (
	thread "baku/thread_text"
	"fmt"
)

func main() {

	// TODO スレッド一覧取得と次のページのURL取得とスレッド中身取得をそれぞれ関数化して、mainで実行する
	// https://l-chika.hatenablog.com/entry/2017/10/18/192930

	m := thread.Thread_ichiran()
	fmt.Printf("%T %v\n", m, len(m))

	thread.Thread_text()

}
