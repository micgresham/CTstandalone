package main


import (
    "fmt"
    "os"
//    "io/ioutil"
//    "io"
//    "bytes"
//    "bufio"
//    "mime/multipart"
//    "net/http"
//    "time"
//    "log"
//    "strings"
//    "github.com/buger/jsonparser"
    "github.com/akamensky/argparse"
//    "sigs.k8s.io/yaml"
//    "gopkg.in/yaml.v2"
//    "crypto/aes"
//    "crypto/cipher"
//    "crypto/md5"
//    "crypto/rand"
//    "encoding/hex"
//    "encoding/base64"
    "github.com/micgresham/goCentral"
)

var appName = "CTcentral_check"
var appVer = "1.0"
var appAuthor = "Michael Gresham"
var appAuthorEmail = "michael.gresham@ghpe.com"
var pgmDescription = fmt.Sprintf("%s: Example program to access Central using the API.",appName)
var central_info goCentral.Central_struct



func main() {

  filename:= "CTconfig.yml"

  goCentral.Passphrase = "“You can use logic to justify almost anything. That’s its power. And its flaw. –Captain Cathryn Janeway"


  parser := argparse.NewParser(appName,pgmDescription)
  //token := parser.String("","token", &argparse.Options{Help: "Central API token",Required: true})
  //url := parser.String("","url", &argparse.Options{Help: "Central API URL",Required: true})
  initDB := parser.Flag("","initDB", &argparse.Options{Help: "Initialize secure storage"})


  fmt.Println("-------------------------------------")
  fmt.Printf("%s Version: %s\r\n",appName, appVer)
  fmt.Printf("Author: %s (%s)\r\n",appAuthor, appAuthorEmail)
  fmt.Println("-------------------------------------")

  err := parser.Parse(os.Args)
  if err != nil {
	fmt.Println(parser.Usage(err))
	return
  }
  
  if *initDB {
    goCentral.Init_DB(filename)
  } 

  fmt.Println("Running as normal.")
  central_info = goCentral.Read_DB(filename)


//======================================================
// test if valid token, refresh the token if needed
//======================================================
  fmt.Printf("\n-------------------------------------------------\n")
  fmt.Printf("Test Central access and renew token if needed\n")
  fmt.Printf("-------------------------------------------------\n")
  respCode, new_token, new_refresh_token := goCentral.Test_central(central_info)
  if (respCode != 200) { 
    fmt.Printf("Central access failed with response code: %d\n",respCode)
//    os.Exit(3)
  } else {
    fmt.Print("Central access OK.  Token verified.")
    fmt.Printf("Response code: %d\n",respCode)
    central_info.Token = new_token
    central_info.Refresh_token = new_refresh_token
    goCentral.Write_DB(filename,central_info)
  }
  fmt.Printf("\n----------------------------------------------------------------------\n")
  fmt.Printf("Refresh Central access token.  Also used to get remain time on token\n")
  fmt.Printf("------------------------------------------------------------------------\n")
  respCode, new_token, new_refresh_token,expires_in := goCentral.RefreshApiToken(central_info)
  if (respCode != 200) { 
    fmt.Printf("Central access failed with response code: %d\n",respCode)
//    os.Exit(3)
  } else {
    fmt.Print("Token good/refreshed.")
    fmt.Printf("Response code: %d\n",respCode)
    fmt.Printf("Token : %s\n",new_token)
    fmt.Printf("Refresh Token : %s\n",new_refresh_token)
    fmt.Printf("Token expires in : %d\n",expires_in)
    central_info.Token = new_token
    central_info.Refresh_token = new_refresh_token
    goCentral.Write_DB(filename,central_info)
  }
   
  central_info = goCentral.Read_DB(filename)
  fmt.Printf("---------------------------\n")
  fmt.Printf("Central Info Decrypted\n")
  fmt.Printf("---------------------------\n")
  fmt.Printf("Central URL: %s\n",central_info.Base_url)
  fmt.Printf("Central Customer ID: %s\n",central_info.Customer_id)
  fmt.Printf("Central Client ID: %s\n",central_info.Client_id)
  fmt.Printf("Central Client Secret: %s\n",central_info.Client_secret)
  fmt.Printf("Central Token: %s\n",central_info.Token)
  fmt.Printf("Central Refresh Token: %s\n",central_info.Refresh_token)
}


