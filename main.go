package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	fileTopVal := fmt.Sprintf("## %s Trending \n", time.Now().Format("2006-01-02"))
	fileDesc := "See what the GitHub community is most excited about today. \n"

	resp, err := http.Get("https://github.com/trending/go?since=daily")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var fileValue string
	doc.Find(".Box .Box-row").Each(func(i int, s *goquery.Selection) {
		boxInfo := s.Find("h1 a")
		author := strings.TrimSpace(boxInfo.Find("span").Text())
		title := strings.TrimSpace(boxInfo.Contents().Last().Text())
		link, ok := boxInfo.Attr("href")
		if !ok {
			log.Fatal("获取a链接错误！")
		}
		linkStr := fmt.Sprintf("https://github.com%s", strings.TrimSpace(link))
		desc := s.Find("p").Text()

		fileValue += fmt.Sprintf("[%s %s](%s) \n", author, title, linkStr)
		fileValue += fmt.Sprintf("%s \n", desc)
	})
	createFile(fileTopVal + fileDesc + "\n" + fileValue)
}

func createFile(str string) {
	curTime := time.Now().Format("2006-01-02")
	filePath := fmt.Sprintf("go/%s.md", curTime)
	_, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE, 0666) //打开文件
	if err != nil {
		fmt.Println("file open fail", err)
		return
	}
	// 写文件
	err1 := ioutil.WriteFile(filePath, []byte(str), 0666)
	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println("文件写入成功！")
}
