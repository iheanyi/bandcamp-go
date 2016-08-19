package main

import (
  "fmt"
  "net/http"
  "log"
  "strings"
  "os"
  "github.com/robertkrimen/otto"
  "github.com/PuerkitoBio/goquery"
)

func executeCode(jsCode string) {
  /*
   Executes arbitrary JavaScript within the Bandcamp page.
  */
  fmt.Println(jsCode)
  fullCodeBlock := "albumData = " + jsCode
  fmt.Println(fullCodeBlock)

  vm := otto.New()
  vm.Run(fullCodeBlock);
  vm.Run(`
    console.log(albumData.stringify());
  `)
}

func fetchPage() {
  url := os.Args[1]

  fmt.Println("Here's the URL that we'll be parsing:", url)
  resp, err := http.Get(os.Args[1])

  doc, err := goquery.NewDocument(url);

  if err != nil {
    log.Fatal(err)
  }

  doc.Find(".yui-skin-sam script").Each(func(i int, s *goquery.Selection) {
    if (i == 1) {
      nodeText := s.Text()
      albumDataDef  := strings.Split(nodeText, "var TralbumData = ")[1]
      albumData := strings.Split(albumDataDef, ";")[0]

      executeCode(albumData)
    }
  })

  defer resp.Body.Close()
}

func main() {
  fetchPage();
}
