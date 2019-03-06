package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func main() {
	resp, err := http.Get("http://myanimelist.net/users.php")
	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()
	parser := html.NewTokenizer(resp.Body)
	for {
		tt := parser.Next()
		if tt == html.ErrorToken {
			return
		}
		token := parser.Token()
		if len(token.Attr) > 0 {
			// [{ href /profile/ReapStick}]
			fmt.Println(token.Attr)
		}
	}
}
