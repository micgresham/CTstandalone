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
    "github.com/micgresham/goCentral"
)


var appName = "CTauto_commit"
var appVer = "1.0"
var appAuthor = "Michael Gresham"
var appAuthorEmail = "michael.gresham@hpe.com"
var pgmDescription = fmt.Sprintf("%s: Enable/Disable autocommit for a list of device serial numbers.",appName)
var central_info goCentral.Central_struct
var useSecureStorage = true

var p_check_dict = map[string]interface{}{}

func set_autocommit(central_info goCentral.Central_struct, serial string, state string) int {

  access_token := central_info.Token
  base_url := central_info.Base_url
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


func do_commit(central_info goCentral.Central_struct, fname string, state string) int {

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

func doesFileExist(fileName string) bool {
   _ , error := os.Stat(fileName)

// check if error is "file not exists"
   if os.IsNotExist(error) {
//     fmt.Printf("%v file does not exist\n", fileName)
     return true
   } else {
//     fmt.Printf("%v file exist\n", fileName)
     return false
   }
}

func main() {

  parser := argparse.NewParser(appName,pgmDescription)
  token := parser.String("","token", &argparse.Options{Help: "Central API token if not using encrypted storage."})
  base_url := parser.String("","url", &argparse.Options{Help: "Central API URL if not using encrypted storage."})
  initDB := parser.Flag("","initDB", &argparse.Options{Help: "Initialize secure storage"})

  infile := parser.String("","infile", &argparse.Options{Help: "Input file consisting of a single device serial on each line"})
  state := parser.String("","state", &argparse.Options{Help: "Autocommit state: enable or disable"})
  test := parser.Flag("t", "test", &argparse.Options{Help: "Enable test mode. No variables will be changed"})
  config := parser.String("c", "config", &argparse.Options{Help: "Config file location"})

  //encrypted storage setup
  SSfilename:= "CTconfig.yml"
  goCentral.Passphrase = "“You can use logic to justify almost anything. That’s its power. And its flaw. –Captain Cathryn Janeway"

  fmt.Println("-------------------------------------")
  fmt.Printf("%s Version: %s\r\n",appName, appVer)
  fmt.Printf("Author: %s (%s)\r\n",appAuthor, appAuthorEmail)
  fmt.Println("-------------------------------------")


  err := parser.Parse(os.Args)
  if err != nil {
        fmt.Println(parser.Usage(err))
        return
  }

  if *test {
    fmt.Println("--------------------------------------------------")
    fmt.Println("TEST MODE - NO VARIABLE CHANGE WILL BE PERFORMED")
    fmt.Println("--------------------------------------------------")
  }

  if *config != "" {
    SSfilename = *config
    fmt.Println("-------------------------------------")
    fmt.Println("Loading config file from ",SSfilename)
    fmt.Println("-------------------------------------")
  }

  //initialize the secure storage if requested
  if *initDB {
     goCentral.Init_DB(SSfilename)
     os.Exit(0) //we do not do anything after the secure storage initialization
  }

  if doesFileExist(SSfilename) {

     //we are not using secure storage
     useSecureStorage = false

     //if the user provided a token AND a URL we will use it
     if (*token != "") {
        if (*base_url == "") {
           fmt.Println("Token supplied, but Central URL is missing.  Both are required if using the command line options.")
           os.Exit(1)
        }
     } else { //ask for the token

     fmt.Print("\nProvide the Central API URL: ")
     fmt.Scanln(base_url)
     fmt.Print("Provide the Central token: ")
     fmt.Scanln(token)

     central_info.Base_url = *base_url
     central_info.Customer_id = ""
     central_info.Client_id = ""
     central_info.Client_secret = ""
     central_info.Token = *token
     central_info.Refresh_token = ""

     }
  } else {
    fmt.Println("Reading secure storage")

    central_info = goCentral.Read_DB(SSfilename)
    fmt.Printf("---------------------------\n")
    fmt.Printf("Central Info Decrypted\n")
    fmt.Printf("---------------------------\n")
    fmt.Printf("Central URL: %s\n",central_info.Base_url)
//    fmt.Printf("Central Customer ID: %s\n",central_info.Customer_id)
//    fmt.Printf("Central Client ID: %s\n",central_info.Client_id)
//    fmt.Printf("Central Client Secret: %s\n",central_info.Client_secret)
    fmt.Printf("Central Token: %s\n",central_info.Token)
//    fmt.Printf("Central Refresh Token: %s\n",central_info.Refresh_token)
  }


if doesFileExist(*infile) {
    fmt.Print("\nProvide the input file name: ")
    fmt.Scanln(infile)
}

if (*state == "") {
   fmt.Print("\nProvide autocommit state to be applied (ON|OFF): ")
   fmt.Scanln(state)
}


  fmt.Println("Input file name:",*infile)
  fmt.Println("Autocommit state to be applied:",*state)

//======================================================
// test if valid token
//======================================================
  respCode, new_token, new_refresh_token := goCentral.Test_central(central_info)
  if (respCode != 200) {
     fmt.Printf("\nCentral access failed with response code: %d\n",respCode)
     os.Exit(3)
  } else {
     fmt.Print("Central access OK.  Token verified.")
     fmt.Printf("Response code: %d\n",respCode)
     central_info.Token = new_token
     if useSecureStorage {
       central_info.Refresh_token = new_refresh_token
       goCentral.Write_DB(SSfilename,central_info)
     }
  }
  if (respCode != 200) {
   fmt.Printf("\nCentral access failed with response code: %d\n",respCode)
   os.Exit(3)
  }


//---------------------------------------
//  os.Exit(0)
//---------------------------------------


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
