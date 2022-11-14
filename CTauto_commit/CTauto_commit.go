package main


import (
    "fmt"
    "os"
//    "io/ioutil"
//    "io"
    "bytes"
    "bufio"
//    "mime/multipart"
    "net/http"
    "time"
//    "log"
//    "strings"
//    "github.com/buger/jsonparser"
    "github.com/akamensky/argparse"
//    "sigs.k8s.io/yaml"
)

type central struct {
    base_url string
    customer_id string
    token string
}

var appName = "CTauto_commit"
var appVer = "1.0"
var appAuthor = "Michael Gresham"
var appAuthorEmail = "michael.gresham@hpe.com"

var p_check_dict = map[string]interface{}{}

func test_central(central_info central) int {

  access_token := central_info.token
  base_url := central_info.base_url
  api_function_url := fmt.Sprintf("%s/configuration/v2/groups",base_url)

  c := http.Client{Timeout: time.Duration(10) * time.Second}
  req, err := http.NewRequest("GET", api_function_url, nil)
  if err != nil {
      fmt.Printf("error %s", err)
      return(0)
  }
  q := req.URL.Query()
  q.Add("limit","1")
  q.Add("offset","0")
  req.URL.RawQuery = q.Encode()

  req.Header.Add("Content-Type", `application/json`)
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",fmt.Sprintf(access_token)))
  req.Header.Add("limit","1")
  resp, err := c.Do(req)
  if err != nil {
      fmt.Printf("error %s", err)
      return(0)
  }

  defer resp.Body.Close()
  return(resp.StatusCode)
}

func set_autocommit(central_info central, serial string, state string) int {

  access_token := central_info.token
  base_url := central_info.base_url
  api_function_url := fmt.Sprintf("%sconfiguration/v1/auto_commit_state/devices",base_url)
  jsonPrep := fmt.Sprintf("{\"serials\": [ \"%s\" ], \"auto_commit_state\": \"%s\"}",serial,state)
  jsonStr := []byte(jsonPrep)

//  fmt.Printf("Json string : %s",jsonPrep)

  c := http.Client{Timeout: time.Duration(10) * time.Second}
  req, err := http.NewRequest("POST", api_function_url, bytes.NewBuffer(jsonStr))
  if err != nil {
      fmt.Printf("error %s", err)
      return(0)
  }
//  q := req.URL.Query()
//  req.URL.RawQuery = q.Encode()

  req.Header.Add("Content-Type", `application/json`)
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",fmt.Sprintf(access_token)))
//  fmt.Println("\n\n",req)
  resp, err := c.Do(req)
  if err != nil {
      fmt.Printf("error %s", err)
      return(0)
  }

  defer resp.Body.Close()
//  fmt.Println("\n\n",resp)
  return(resp.StatusCode)
}


func do_commit(central_info central, fname string, state string) int {

  f := func() *os.File {
  f, err := os.Open(fname)
  if err != nil {
     panic(err)
   }
   return f
  }()
  fileScanner := bufio.NewScanner(f)
  fileScanner.Split(bufio.ScanLines)
  
  var fileLines []string

  for fileScanner.Scan() {
    fileLines = append(fileLines, fileScanner.Text())
  }

  count := 0

  f.Close()

  for _, line := range fileLines {
        fmt.Printf("Setting commit for %s to %s - ",line,state)
        fmt.Println(set_autocommit(central_info,line, state))
	count = count + 1 
  }
  return(count)
}

func main() {

  pgmDescription:= fmt.Sprintf("%s: Enable/Disable autocommit for a list of device serial numbers.",appName)
  parser := argparse.NewParser("test_api",pgmDescription)
  token := parser.String("","token", &argparse.Options{Help: "Central API token",Required: true})
  url := parser.String("","url", &argparse.Options{Help: "Central API URL",Required: true})
  infile := parser.String("","infile", &argparse.Options{Help: "Input file consisting of a single device serial on each line",Required: true})
  state := parser.String("","state", &argparse.Options{Help: "Autocommit state: enable or disable",Required: true})
  test := parser.Flag("t", "test", &argparse.Options{Help: "Enable test mode. No variables will be changed"})


  fmt.Println("-------------------------------------")
  fmt.Printf("%s Version: %s\r\n",appName, appVer)
  fmt.Printf("Author: %s (%s)\r\n",appAuthor, appAuthorEmail)
  fmt.Println("-------------------------------------")

  err := parser.Parse(os.Args)
  if err != nil {
	fmt.Println(parser.Usage(err))
	return
  }

  central_info := central {
    base_url: fmt.Sprintf(*url),
    customer_id: "",
    token: fmt.Sprintf(*token),
  }

  fmt.Println("Input file name:",*infile)
  fmt.Println("Autocommit state to be applied:",*state)

//======================================================
// test if valid token
//======================================================
  respCode := test_central(central_info)
  fmt.Printf("Central Status: %s(%d)\r\n",http.StatusText(respCode),respCode)
  if (respCode != 200) { os.Exit(3)}

//======================================================
// perform autocommit state change 
//======================================================
 
  count := 0

  fmt.Println("--------------------------------")
  if *test {
    fmt.Println("TEST MODE - NO AUTOCOMMIT STATE CHANGE")
  } else {
    fmt.Println("Performing autocommit state change")
    count = do_commit(central_info, *infile, *state)
  }
  fmt.Println("--------------------------------")

//  temp_string := ""
//  temp_dict := map[interface{}]interface{}{}
  

  fmt.Println("Total devices: ", count)

}
