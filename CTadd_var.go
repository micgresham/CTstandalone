package main


import (
    "fmt"
    "os"
    "io/ioutil"
    "io"
    "bytes"
    "bufio"
    "mime/multipart"
    "net/http"
    "time"
    "log"
    "strings"
    "github.com/buger/jsonparser"
    "github.com/akamensky/argparse"
    "sigs.k8s.io/yaml"
    "github.com/micgresham/goCentral"
)
/*
type central struct {
    base_url string
    customer_id string
    token string
}
*/

var appName = "CTadd_var"
var appVer = "1.5"
var appAuthor = "Michael Gresham"
var appAuthorEmail = "michael.gresham@hpe.com"
var pgmDescription = fmt.Sprintf("%s: Add/Update a variable in Central for a list of device serial numbers.",appName)
var central_info goCentral.Central_struct
var useSecureStorage = true

var p_check_dict = map[string]interface{}{}

func get_variables(central_info goCentral.Central_struct, serial string) string {

  access_token := central_info.Token
  base_url := central_info.Base_url
  api_function_url := fmt.Sprintf("%sconfiguration/v1/devices/%s/template_variables",base_url,serial)

  c := http.Client{Timeout: time.Duration(10) * time.Second}
  req, err := http.NewRequest("GET", api_function_url, nil)
  if err != nil {
      fmt.Printf("error %s", err)
      return("")
  }
  q := req.URL.Query()
  q.Add("device_serial",serial)
  req.URL.RawQuery = q.Encode()

  req.Header.Add("Content-Type", `application/json`)
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",fmt.Sprintf(access_token)))
  resp, err := c.Do(req)
  if err != nil {
      fmt.Printf("error %s", err)
      return("")
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
//  fmt.Printf("%s",body)
//  fmt.Printf("**************\n")
  _sys_serial, err := jsonparser.GetString(body, "data", "variables", "_sys_serial")
  if err != nil {
      fmt.Printf("error %s", err)
      return("")
  }

  data, _, _, err := jsonparser.Get(body, "data")
  if err != nil {
      fmt.Printf("error %s", err)
      return("")
  }
  data_dict := fmt.Sprintf("{ \"%s\" : %s }",_sys_serial,data)

  return(string(data_dict))
}

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

func set_variables(central_info goCentral.Central_struct, fname string) string {

  access_token := central_info.Token
  base_url := central_info.Base_url
  api_function_url := fmt.Sprintf("%s/configuration/v1/devices/template_variables",base_url)

  extraFields := map[string]string{}
  data, w, err := createMultipartFormData("variables",fname,fname, extraFields )
    if err != nil {
       return("ERROR Create Multipart Form Data")
    }

  c := http.Client{Timeout: time.Duration(10) * time.Second}
  req, err := http.NewRequest("PATCH", api_function_url, &data)
  if err != nil {
      fmt.Printf("error %s", err)
      return("ERROR")
  }
  q := req.URL.Query()
  req.URL.RawQuery = q.Encode()

  req.Header.Set("Content-Type", w.FormDataContentType())
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",fmt.Sprintf(access_token)))
  resp, err := c.Do(req)
  if err != nil {
      fmt.Printf("error %s", err)
      return("")
  }

  defer resp.Body.Close()
//  body, err := ioutil.ReadAll(resp.Body)
//  fmt.Println(string(body))
  
  return(resp.Status)
}



func p_check(central_info goCentral.Central_struct, fname string, oprefix string, variable string, value string) {

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

  f.Close()

  ofilename := fmt.Sprintf("%s-%s",oprefix,fname)
  of, err := os.Create(ofilename)
    if err != nil {
        fmt.Println(err)
        return
    }

  for _, line := range fileLines {
//        fmt.Println(line)
        return_dict := []byte(get_variables(central_info, strings.TrimSpace(line)))
//        fmt.Println(string(return_dict))
        _sys_lan_mac, err := jsonparser.GetString(return_dict, line,"variables","_sys_lan_mac")
//         fmt.Println("*******************")
//         fmt.Println(line,_sys_lan_mac)
//         fmt.Println("*******************")
         
         tmp_json := fmt.Sprintf("{\"_sys_serial\":\"%s\", \"_sys_lan_mac\":\"%s\", \"%s\":\"%s\"}",line, _sys_lan_mac, variable, value)
         p_check_dict[line] = tmp_json

//        jsonparser.ObjectEach(return_dict, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
//              fmt.Printf("'%s'\n Value: '%s'\n Type: %s\n", string(key), string(value), dataType)
//              return nil
//            }, line,"[0]", "variables")
//	if err != nil {
//		fmt.Printf("err: %v\n", err)
//		return
//	}
        yaml_vars, err := yaml.JSONToYAML(return_dict)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
//        fmt.Println(string(yaml_vars))
        n, err := of.WriteString(string(yaml_vars))
        if err != nil {
          fmt.Println(err)
          return
        }
        if n == 0 {
          fmt.Println("No bytes written")
          return
        }
  } 
  err = of.Close()
   if err != nil {
       fmt.Println(err)
       return
  }
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
  variable := parser.String("","variable", &argparse.Options{Help: "Variable to create/update"})
  value := parser.String("","value", &argparse.Options{Help: "Value to assign to the variable"})
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

if (*variable == "") {
   fmt.Print("\nProvide the Central variable to be changed: ")
   fmt.Scanln(variable)
}

if (*value == "") {
   fmt.Print("\nProvide the Central variable's new value: ")
   fmt.Scanln(value)
}
   
fmt.Println("Input file name:",*infile)
fmt.Println("Variable to be changed:",*variable)
fmt.Println("Value for variable:",*value)

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
// perform pre check
//======================================================
  fmt.Println("--------------------------------")
  fmt.Println("Performing Pre-Check")
  fmt.Println("--------------------------------")
  p_check(central_info,*infile,"pre-check",*variable,*value)


//======================================================
// perform variable add/update
//======================================================

  fmt.Println("--------------------------------")
  if *test {
    fmt.Println("TEST MODE - NO VARIABLE CHANGE")
  } else {
    fmt.Println("Performing variable add/update")
  }
  fmt.Println("--------------------------------")

  count := 0
  count_batches := 0  
  temp_string := ""
  temp_dict := map[interface{}]interface{}{}
  
  for key, element := range p_check_dict {

    fmt.Println("Setting device",key, "=>", "Values:", element)
    temp_dict[key] = element  

    if count == 1000 {
      fmt.Println("STOP, Hammer time")

      //create temp file
      tfile, err := ioutil.TempFile(os.TempDir(), "CTadd_var")
      if err != nil {
        log.Fatal(err)
      }
      count_batches = count_batches + 1
      if _, err := tfile.Write([]byte("{")); err != nil {
        log.Fatal(err)
      }
      
      for key2, element2 := range temp_dict {
        count = count + 1
        if (count == 1) {
          temp_string = fmt.Sprintf("\"%s\": %v",key2,element2)
        } else {
          temp_string = fmt.Sprintf(",\"%s\": %v",key2,element2)
        }
        if _, err := tfile.Write([]byte(temp_string)); err != nil {
          log.Fatal(err)
        }
 
      } //end for
      if _, err := tfile.Write([]byte("}")); err != nil {
        log.Fatal(err)
      }
      tfile.Close()
      response := ""
      if *test {
        fmt.Println("SIMULATION set_variables")
        response = "NO ACTION"
      } else {
        response = set_variables(central_info, tfile.Name())
      }
      fmt.Println(response)
      temp_dict = map[interface{}]interface{}{}
      fmt.Println(tfile.Name())

     //this will delete the temp file on program exit
//      defer os.Remove(tfile.Name())
      fmt.Println("Total devices in file: ", count)
      count = 1
    } //endif

  } // end for

  //create temp file
  tfile, err := ioutil.TempFile(os.TempDir(), "CTadd_var")
  if err != nil {
    log.Fatal(err)
  }
  count_batches = count_batches + 1

  if _, err := tfile.Write([]byte("{")); err != nil {
    log.Fatal(err)
  }
  
  for key2, element2 := range temp_dict {
    count = count + 1
    if (count == 1) {
      temp_string = fmt.Sprintf("\"%s\": %v",key2,element2)
    } else {
      temp_string = fmt.Sprintf(",\"%s\": %v",key2,element2)
    }
    if _, err := tfile.Write([]byte(temp_string)); err != nil {
      log.Fatal(err)
    }
  } 
  if _, err := tfile.Write([]byte("}")); err != nil {
    log.Fatal(err)
  }
  tfile.Close()
  fmt.Println("Tmp file name:",tfile.Name())
  response := ""
  if *test {
    fmt.Println("SIMULATION set_variables")
    response = "NO ACTION"
  } else {
    response = set_variables(central_info, tfile.Name())
  }
  fmt.Println("Status:",response)

//this will delete the temp file on program exit
//      defer os.Remove(tfile.Name())
  fmt.Println("Total devices in file: ", count)

  fmt.Println("Total devices: ", count*count_batches)

//======================================================
// perform post check
//======================================================
  fmt.Println("--------------------------------")
  fmt.Println("Performing Post-Check")
  fmt.Println("--------------------------------")
  p_check(central_info,*infile,"post-check",*variable,*value)
   
}
