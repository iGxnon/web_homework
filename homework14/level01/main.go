package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	MainUrl = "http://xiaodiaodaya.cn/"

	// ContainerUrls 使用 map 保证不重复爬取
	ContainerUrls = make(map[string]struct{})
	InfoUrls      = make(map[string]struct{})
	N             = struct{}{}
)

func main() {
	fmt.Println("Start spider...")
	CrawContainer(MainUrl)
	CrawInfo()
}

func CrawContainer(containerUrl string) {
	if _, ok := ContainerUrls[containerUrl]; ok {
		return
	}
	// 装有笑话的url基本上满足以下公式
	if regexp.MustCompile(`html$`).MatchString(containerUrl) ||
		strings.Contains(containerUrl, "?id=") {
		InfoUrls[containerUrl] = N
		return
	}
	ContainerUrls[containerUrl] = N
	resp, err := http.Get(containerUrl)
	if err != nil {
		panic(err)
	}
	reg := regexp.MustCompile(`href="/[a-z0-9=./?_-]+"`)
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	findAll := reg.FindAll(all, -1)
	for _, url := range findAll {
		CrawContainer(MainUrl + string(url)[6:len(url)-1])
	}
}

func CrawInfo() {
	fmt.Println("Finish CrawContainer")
	fmt.Println("Start CrawInfo")
	jokes, err := os.OpenFile("./joke.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer jokes.Close()
	if err != nil {
		panic(err)
	}
	// 开很多 goroutine 爬
	waitGroup := sync.WaitGroup{}
	for url := range InfoUrls {
		waitGroup.Add(1)
		go func(u string) {
			resp, err := http.Get(u)
			all, err := ioutil.ReadAll(resp.Body)
			// 还贴心的准备了开始和结束的注释，他好温柔，我爆哭
			startSyb := []byte("<!--listS-->")
			endSyb := []byte("<!--listE-->")

			all = all[bytes.Index(all, startSyb)+len(startSyb) : bytes.Index(all, endSyb)]
			if err != nil {
				panic(err)
			}
			reg := regexp.MustCompile(`\d、`)
			for _, joke := range reg.Split(string(all), -1) {
				_, err = jokes.WriteString(strings.ReplaceAll(joke, "<br/>", "\n"))
				if err != nil {
					panic(err)
				}
			}
			waitGroup.Done()
		}(url)
	}
	waitGroup.Wait()
}
