package main

import (
	"bufio"
	"bytes"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/bmaupin/go-epub"
	"github.com/go-shiori/go-readability"
	"golang.org/x/net/html"
)

func forEachImgSrc(node *html.Node, f func(html.Attribute) *html.Attribute) {
	if node.Type == html.ElementNode && strings.EqualFold(node.Data, "img") {
		for i, attr := range node.Attr {
			if strings.EqualFold(attr.Key, "src") {
				if attr := f(attr); attr != nil {
					node.Attr[i] = *attr
				}
			}
		}
	}
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		forEachImgSrc(n, f)
	}
}

func addArticle(e *epub.Epub, url string) error {
	print("↓↓↓ ", url, "\n")

	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return err
	}

	article.Content = "<p><h2>" + article.Title + "</h2></p>" + article.Content

	doc, err := html.Parse(bytes.NewBufferString(article.Content))
	if err != nil {
		return err
	}

	forEachImgSrc(doc.FirstChild, func(attr html.Attribute) *html.Attribute {
		print("↓↓↓ ", attr.Val, "\n")

		path, err2 := e.AddImage(attr.Val, "")
		if err2 != nil {
			err = err2
			return &attr
		}

		attr.Val = path
		return &attr
	})

	var b bytes.Buffer
	if err := html.Render(&b, doc.FirstChild); err != nil {
		return err
	}

	_, err = e.AddSection(b.String(), article.Title, "", "")
	return err
}

func run() int {
	ignore := flag.Bool("ignore-errors", false, "ignore errors during article downloading")
	flag.Parse()

	if len(os.Args) < 2 {
		print("Usage: ", os.Args[0], " --ignore-errors <filename>\n")
		return 0
	}

	e := epub.NewEpub(time.Now().Format("January 2 15:04:05"))

	s := bufio.NewScanner(os.Stdin)
	var i int
	for s.Scan() {
		if len(s.Text()) > 0 {
			err := addArticle(e, s.Text())
			if err != nil {
				print("Error: ", err.Error(), "\n")
				if *ignore {
					continue
				}
				return 1
			}
			i++
		}
	}

	if i > 0 {
		err := e.Write(os.Args[1])
		if err != nil {
			print("Error: ", err.Error(), "\n")
			return 1
		}
	}
	return 0
}

func main() {
	os.Exit(run())
}
