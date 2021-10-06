package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"
)
var once sync.Once
var wg sync.WaitGroup
var lock sync.Mutex
var filePath = "dns_github.txt"
var file, _ = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)

type Loader struct {
	preUrl string
	headers map[string]string

}

var lloader *Loader

func (l *Loader) getHtml(suffix string) string {
	client := &http.Client{}
	request, err := http.NewRequest("GET", l.preUrl+suffix, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Add("user-agent", l.headers["user-agent"])
	response, err := client.Do(request)
	defer response.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	html := buf.String()
	return html
}

func (l *Loader) Analysis(suffix string) {
	reg, _ := regexp.Compile(`(\d+\.\d+\.\d+\.\d+)`)
	result := reg.FindAllString(l.getHtml(suffix), -1)
	fmt.Fprintf(file,getMaxSameElement(result) + " " + suffix + "\n")
	wg.Done()
}

func getMaxSameElement(arr []string) string {
	dict := make(map[string]int)
	maxNum := -1
	res := ""
	for _, value := range arr {
		if _, ok := dict[value]; ok {
			dict[value]++
		}else {
			dict[value] = 1
		}
		if dict[value] > maxNum {
			res = value
		}
	}
	return res
}

func New() *Loader {
	once.Do(func() {
		lloader = &Loader{
			preUrl:  "https://websites.ipaddress.com/",
			headers: map[string]string{
				"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36",
			},
		}
	})
	return lloader
}

func main() {
	loader := New()
	wg.Add(3)
	go loader.Analysis("github.com")
	go loader.Analysis("assets-cdn.github.com")
	go loader.Analysis("github.global.ssl.fastly.net")
	wg.Wait()
	file.Close()
}
