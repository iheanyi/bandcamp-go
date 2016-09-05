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
  fullCodeBlock := "albumData = " + jsCode
  fmt.Println(jsCode)

  vm := otto.New()
  vm.Run(fullCodeBlock);
  vm.Run(`
  albumDataStr = JSON.stringify(albumData);
  `)

  /* TO-DO: Fix Decoding of JSON from Otto VM into an actual Go structure. 
  Mad close to getting in working though. */
  if value, err := vm.Get("albumDataStr"); err == nil {
    fmt.Println("This will be decoding.");
    //fmt.Println(value.ToString())
    if valueStr, err := value.ToString(); err == nil {
      jsonByteArray := []byte(valueStr)
      var albumMap interface{}
      jsonErr := json.Unmarshal(jsonByteArray, &albumMap)
      if (jsonErr == nil) {
        fmt.Println(albumMap)
      }
    }
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
