package main


import (
    "fmt"
    "os"
//    "io/ioutil"
    "io"
    "bytes"
//    "bufio"
    "mime/multipart"
//    "net/http"
//    "time"
    "log"
//    "strings"
//    "strconv"
//    "github.com/buger/jsonparser"
    "github.com/akamensky/argparse"
//    "sigs.k8s.io/yaml"
    "github.com/micgresham/goCentral"
    "github.com/manifoldco/promptui"
)

var appName = "example2"
var appVer = "1.0"
var appAuthor = "George P. Burdell"
var appAuthorEmail = "gpburdell@gatech.edu"
var pgmDescription = fmt.Sprintf("%s: Description of what the program does.",appName)
var central_info goCentral.Central_struct
var useSecureStorage = true

func createMultipartFormData(fileFieldName, filePath string, fileName string, extraFormFields map[string]string) (b bytes.Buffer, w *multipart.Writer, err error) {
  w = multipart.NewWriter(&b)
  var fw io.Writer
  file, err := os.Open(filePath)

  if fw, err = w.CreateFormFile(fileFieldName, fileName); err != nil {
      return
  }
  if _, err = io.Copy(fw, file); err != nil {
      return
  }

  for k, v := range extraFormFields {
      w.WriteField(k, v)
  }

  w.Close()

  return
}

//--------------------------------------

//--------------------------------------
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

func yesNo(promptText string) bool {
    prompt := promptui.Select{
        Label: promptText,
        Items: []string{"Yes", "No"},
    }
    _, result, err := prompt.Run()
    if err != nil {
        log.Fatalf("Prompt failed %v\n", err)
    }
    return result == "Yes"
}

func promptString(promptText string, defaultText string) string {
    prompt := promptui.Prompt{
	Label:   promptText,
	Default: defaultText,
    }

    result, err := prompt.Run()

    if err != nil {
	fmt.Printf("Prompt failed %v\n", err)
        return ""
    }
    return result
}

func promptSelect(promptText string,promptItems []string) string {
    prompt := promptui.Select{
        Label: promptText,
        Items: promptItems,
    }

    _, result, err := prompt.Run()

    if err != nil {
	fmt.Printf("Prompt failed %v\n", err)
        return ""
    }
   return result
}

func main() {

  parser := argparse.NewParser(appName,pgmDescription)
  token := parser.String("","token", &argparse.Options{Help: "Central API token if not using encrypted storage."})
  base_url := parser.String("","url", &argparse.Options{Help: "Central API URL if not using encrypted storage."})
  initDB := parser.Flag("","initDB", &argparse.Options{Help: "Initialize secure storage"})
  test := parser.Flag("t", "test", &argparse.Options{Help: "Enable test mode. No variables will be changed"})
  infile := parser.String("","infile", &argparse.Options{Help: "Input file consisting of a single device serial on each line"})

  //encrypted storage setup
//  SSfilename:= "../CTcentral_check/CTconfig.yml"
  SSfilename:= "CTconfig.yml"

  //  SSfilename:= "CTconfig.yml"
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
    fmt.Println("TEST MODE - NO CHANGES WILL BE PERFORMED")
    fmt.Println("--------------------------------------------------")
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

     *base_url = promptString("Provide the Central API URL","")
     *token = promptString("Provide the Central token","")

     central_info.Base_url = *base_url
     central_info.Customer_id = ""
     central_info.Client_id = ""
     central_info.Client_secret = ""
     central_info.Token = *token
     central_info.Refresh_token = ""
    fmt.Printf("\n")

     }
  } else {
    fmt.Println("Reading secure storage")

    central_info = goCentral.Read_DB(SSfilename)
    fmt.Printf("---------------------------\n")
    fmt.Printf("Central Info Decrypted\n")
    fmt.Printf("---------------------------\n")
    fmt.Printf("Central URL: %s\n",central_info.Base_url)
    fmt.Printf("Central Token: %s\n",central_info.Token)
  }

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

//======================================================
// perform pre check
//======================================================
//  fmt.Println("--------------------------------")
//  fmt.Println("Performing Pre-Check")
//  fmt.Println("--------------------------------")


//======================================================
// perform action
//======================================================


  if doesFileExist(*infile) {
    fmt.Print("\nProvide the input file name: ")
    fmt.Scanln(infile)
  }

  //if the file does not exist, notify the user and exit
  if doesFileExist(*infile) {
	  fmt.Println("Filename does not exist. Exiting....")
	  os.Exit(3)
  }

  fmt.Println("\n-----------------------------")
  fmt.Println("Input file name:",*infile)
  fmt.Println("-----------------------------")



//======================================================
// perform post check
//======================================================
//  fmt.Println("--------------------------------")
//  fmt.Println("Performing Post-Check")
//  fmt.Println("--------------------------------")
   
}



