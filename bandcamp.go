package main

import (
  "fmt"
  "net/http"
  "log"
  "strings"
  "os"
  "io"
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

var albumMap map[string]interface{}

func GenerateAlbumMap(jsCode string) map[string]interface{} {
  /*
  Executes arbitrary JavaScript within the Bandcamp page.
  */
  fullCodeBlock := "albumData = " + jsCode
  fmt.Println(jsCode)

  vm := otto.New()
  vm.Run(fullCodeBlock)
  vm.Run(`
  albumDataStr = JSON.stringify(albumData);
  `)

  /* TO-DO: Fix Decoding of JSON from Otto VM into an actual Go structure. 
  Mad close to getting in working though. */
  if value, err := vm.Get("albumDataStr"); err == nil {
    if valueStr, err := value.ToString(); err == nil {
      jsonByteArray := []byte(valueStr)
      jsonErr := json.Unmarshal(jsonByteArray, &albumMap)

      if jsonErr != nil {
        fmt.Println("Error encoding JSON from the JS.")
        panic(jsonErr)
      }
    }
  }

  return albumMap
}

func FetchPage(url string) {
  fmt.Println("Here's the URL that we'll be parsing: ", url)
  resp, err := http.Get(os.Args[1])

  doc, err := goquery.NewDocument(url)

  if err != nil {
    log.Fatal(err)
  }

  doc.Find(".yui-skin-sam script").Each(func(i int, s *goquery.Selection) {
    if (i == 1) {
      nodeText := s.Text()
      albumDataDef  := strings.Split(nodeText, "var TralbumData = ")[1]
      albumData := strings.Split(albumDataDef, ";")[0]

      albumInfo := GenerateAlbumMap(albumData)
      DownloadAlbumTracks(albumInfo)
    }
  })

  defer resp.Body.Close()
}

func DownloadAlbumTracks(albumInfo map[string]interface{}) {
  albumTracks := albumInfo["trackinfo"].([]interface{})
  currentAlbum := albumInfo["current"].(map[string]interface{})
  fmt.Println("Album Artist: ", albumInfo["artist"])
  fmt.Println("Album Title: ", currentAlbum["title"])
  fmt.Println("Album Release Date: ", albumInfo["album_release_date"])
  fmt.Println("------------------------------------")

  // Create Directory
  directoryName := fmt.Sprintf("albums/%s", currentAlbum["title"])
 
  // Making Directory
  os.MkdirAll(directoryName, 0700)

  for _, trackInstance := range albumTracks {
    track := trackInstance.(map[string]interface{})
    trackFile := track["file"].(map[string]interface{})
    trackUrl := fmt.Sprintf("https:%s", trackFile["mp3-128"])
    trackFileName := fmt.Sprintf("%s/%s.mp3", directoryName, track["title"])

    fmt.Println("Track Number: ", track["track_num"])
    fmt.Println("Track File: ", trackUrl)
    fmt.Println("Track Title: ", track["title"])
    fmt.Println("------------------------------------")

    trackOutFile, err := os.Create(trackFileName)

    if err != nil {
      fmt.Println("Error creating the filename.")
      panic(err)
    }

    defer trackOutFile.Close()

    resp, _ := http.Get(trackUrl)
    defer resp.Body.Close()

    _, downloadErr := io.Copy(trackOutFile, resp.Body)

    if downloadErr != nil {
      fmt.Println("Error downloading the file.")
      panic(downloadErr)
    }
  }
}

func main() {
  albumUrl := os.Args[1]
  FetchPage(albumUrl)
}
