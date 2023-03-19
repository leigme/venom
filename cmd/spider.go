package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/gocolly/colly"
	"github.com/leigme/loki/app"
	loki "github.com/leigme/loki/cobra"
	"github.com/leigme/loki/file"
	"github.com/spf13/cobra"
	"log"
	"net/url"
	"path/filepath"
	"strings"
)

const (
	UserAgent   = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"
	Www6vhaoNet = "www.6vhao.net"
)

var encode = "utf8"

func init() {
	c := colly.NewCollector(
		colly.UserAgent(UserAgent),
		colly.MaxDepth(2),
	)
	sc := &SpiderCommand{
		SpiderParser: make(map[string]SpiderParse, 0),
		C:            c,
	}
	sc.SpiderParser[Www6vhaoNet] = Parse6vhao
	loki.Add(rootCmd, sc)
}

type SpiderInterface interface {
	Parse(string) SpiderParse
}

type SpiderCommand struct {
	SpiderParser map[string]SpiderParse
	C            *colly.Collector
}

type SpiderParse func(url string, c *colly.Collector) map[string]string

// http://www.6vhao.net/dlz/2023-01-15/44682.html

func (sc *SpiderCommand) Execute() loki.Exec {
	return func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			log.Fatal("inputUrl is nil")
		}
		inputUrl := args[0]
		result := sc.Parse(inputUrl)(inputUrl, sc.C)
		data, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		for k, v := range result {
			if strings.Contains(k, ".2160p.") {
				fmt.Println(v)
			}
		}
		err = file.Create(filepath.Join(app.WorkDir(), "spider", "data.json"), data)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (sc *SpiderCommand) Parse(inputUrl string) SpiderParse {
	return sc.SpiderParser[urlHost(inputUrl)]
}

func urlHost(inputUrl string) string {
	u, err := url.Parse(inputUrl)
	if err != nil {
		log.Fatal(err)
	}
	return u.Host
}

func Parse6vhao(inputUrl string, c *colly.Collector) map[string]string {
	c.OnHTML("meta[http-equiv]", func(element *colly.HTMLElement) {
		if strings.Contains(element.Attr("content"), "gb") {
			encode = "gbk"
		}
	})
	result := make(map[string]string, 0)
	c.OnHTML(
		"tbody", func(e *colly.HTMLElement) {
			e.ForEach("tr", func(i int, item *colly.HTMLElement) {
				href := item.ChildAttr("a[href]", "href")
				title := item.ChildText("td")
				title = mahonia.NewDecoder(encode).ConvertString(title)
				if strings.Contains(href, "magnet") {
					result[title] = href
				}
			})
		},
	)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	err := c.Visit(inputUrl)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
