package main

import (
  "fmt"
  "net/http"
  "log"
  "strings"
  //"io"
  "os"
  "github.com/robertkrimen/otto"
  "github.com/PuerkitoBio/goquery"
)

func executeCode(jsCode string) {
  /*
  
  */
  fmt.Println(jsCode)
  vm := otto.New()
  vm.Run(jsCode);
  vm.Run(`
    console.log(TralbumData);
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
