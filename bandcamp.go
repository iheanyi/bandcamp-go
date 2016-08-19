package main

import (
  "fmt"
  "net/http"
  "log"
  "strings"
  "os"
  "encoding/json"
  "github.com/robertkrimen/otto"
  "github.com/PuerkitoBio/goquery"
)

// Struct for trAlbumData variable
type AlbumData struct {
  artFullsizeUrl string
  trackinfo []string
}

// Struct for BandData variable
type BandData struct {
  name string
}

// Struct for EmbedData variable
type EmbedData struct {
  albumTitle string
  albumEmbedData string
}

func ExecuteCode(jsCode string) {
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

  /* TO-DO: Fix Decoding of JSON from Otto VM into an actual Go structure. 
  Mad close to getting in working though. */
  if value, err := vm.Get("albumData"); err == nil {
    if valueStr, err := value.ToString(); err == nil {
      dec := json.NewDecoder(strings.NewReader(valueStr))
      fmt.Println(dec)
    }
  } else {
    log.Fatal("Error occurred!")
    log.Fatal(err)
  }
}

func FetchPage(url string) {
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

      ExecuteCode(albumData)
    }
  })

  defer resp.Body.Close()
}

func main() {
  albumUrl := os.Args[1]
  FetchPage(albumUrl);
}
