package main


import (
    "fmt"
    "os"
    "io/ioutil"
    "net/http"
    "time"
    "strings"
    "github.com/buger/jsonparser"
    "github.com/akamensky/argparse"
    "github.com/micgresham/goCentral"
)

var appName = "CTsite_serial"
var appVer = "1.5"
var appAuthor = "Michael Gresham"
var appAuthorEmail = "michael.gresham@hpe.com"
var pgmDescription = fmt.Sprintf("%s: Retrieve a list of device serial numbers for a given site.",appName)
var central_info goCentral.Central_struct
var useSecureStorage = true


var p_check_dict = map[string]interface{}{}

func get_apSerials(central_info goCentral.Central_struct, siteName string)[]string {

  my_serials := []string{}

  access_token := central_info.Token
  base_url := central_info.Base_url
  api_function_url := fmt.Sprintf("%s/monitoring/v2/aps",base_url)

  c := http.Client{Timeout: time.Duration(10) * time.Second}
  req, err := http.NewRequest("GET", api_function_url, nil)
  if err != nil {
      fmt.Printf("error %s", err)
      return(nil)
  }
  q := req.URL.Query()
  q.Add("site",strings.TrimSpace(siteName))
  q.Add("offset","0")
  req.URL.RawQuery = q.Encode()
 
//  fmt.Println(req)

  req.Header.Add("Content-Type", `application/json`)
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",fmt.Sprintf(access_token)))
  resp, err := c.Do(req)
  if err != nil {
      fmt.Printf("error %s", err)
      return(nil)
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
//  fmt.Println(resp)

  if err != nil {
      fmt.Printf("error %s", err)
      return(nil)
  }
  jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
    serial, _ := jsonparser.GetString(value, "serial")
    my_serials = append(my_serials, serial)
  }, "aps")

    return(my_serials)
}

func get_switchSerials(central_info goCentral.Central_struct, siteName string)[]string {

  my_serials := []string{}

  access_token := central_info.Token
  base_url := central_info.Base_url
  api_function_url := fmt.Sprintf("%s/monitoring/v1/switches",base_url)

  c := http.Client{Timeout: time.Duration(10) * time.Second}
  req, err := http.NewRequest("GET", api_function_url, nil)
  if err != nil {
      fmt.Printf("error %s", err)
      return(nil)
  }
  q := req.URL.Query()
  q.Add("site",strings.TrimSpace(siteName))
  q.Add("offset","0")
  req.URL.RawQuery = q.Encode()
 
//  fmt.Println(req)

  req.Header.Add("Content-Type", `application/json`)
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",fmt.Sprintf(access_token)))
  resp, err := c.Do(req)
  if err != nil {
      fmt.Printf("error %s", err)
      return(nil)
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
//  fmt.Println(resp)

  if err != nil {
      fmt.Printf("error %s", err)
      return(nil)
  }
  jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
    serial, _ := jsonparser.GetString(value, "serial")
    my_serials = append(my_serials, serial)
  }, "switches")

    return(my_serials)
}

func get_gatewaySerials(central_info goCentral.Central_struct, siteName string)[]string {

  my_serials := []string{}

  access_token := central_info.Token
  base_url := central_info.Base_url
  api_function_url := fmt.Sprintf("%s/monitoring/v1/gateways",base_url)

  c := http.Client{Timeout: time.Duration(10) * time.Second}
  req, err := http.NewRequest("GET", api_function_url, nil)
  if err != nil {
      fmt.Printf("error %s", err)
      return(nil)
  }
  q := req.URL.Query()
  q.Add("site",strings.TrimSpace(siteName))
  q.Add("offset","0")
  req.URL.RawQuery = q.Encode()
 
//  fmt.Println(req)

  req.Header.Add("Content-Type", `application/json`)
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",fmt.Sprintf(access_token)))
  resp, err := c.Do(req)
  if err != nil {
      fmt.Printf("error %s", err)
      return(nil)
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
//  fmt.Println(resp)

  if err != nil {
      fmt.Printf("error %s", err)
      return(nil)
  }
  jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
    serial, _ := jsonparser.GetString(value, "serial")
    my_serials = append(my_serials, serial)
  }, "gateways")

    return(my_serials)
}

func get_mcSerials(central_info goCentral.Central_struct, siteName string)[]string {

  my_serials := []string{}

  access_token := central_info.Token
  base_url := central_info.Base_url
  api_function_url := fmt.Sprintf("%s/monitoring/v1/mobility_controllers",base_url)

  c := http.Client{Timeout: time.Duration(10) * time.Second}
  req, err := http.NewRequest("GET", api_function_url, nil)
  if err != nil {
      fmt.Printf("error %s", err)
      return(nil)
  }
  q := req.URL.Query()
  q.Add("site",strings.TrimSpace(siteName))
  q.Add("offset","0")
  req.URL.RawQuery = q.Encode()
 
//  fmt.Println(req)

  req.Header.Add("Content-Type", `application/json`)
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",fmt.Sprintf(access_token)))
  resp, err := c.Do(req)
  if err != nil {
      fmt.Printf("error %s", err)
      return(nil)
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
//  fmt.Println(resp)

  if err != nil {
      fmt.Printf("error %s", err)
      return(nil)
  }
  jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
    serial, _ := jsonparser.GetString(value, "serial")
    my_serials = append(my_serials, serial)
  }, "mcs")

  return(my_serials)
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

  parser := argparse.NewParser("CTsite_serials",pgmDescription)
  token := parser.String("","token", &argparse.Options{Help: "Central API token if not using encrypted storage."})
  base_url := parser.String("","url", &argparse.Options{Help: "Central API URL if not using encrypted storage."})
  initDB := parser.Flag("","initDB", &argparse.Options{Help: "Initialize secure storage"})

  siteName := parser.String("","site", &argparse.Options{Help: "Target site name"})
  devType := parser.String("","dev_type", &argparse.Options{Help: "Device type requested: ALL|AP|SW|GW|MC",Required: false})

  ap_serials := []string{}
  sw_serials := []string{}
  gw_serials := []string{}
  mc_serials := []string{}

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



 if (*siteName == "") {
    fmt.Print("\nProvide the site name: ")
    fmt.Scanln(siteName)
 }

 if (*devType == "") {
   fmt.Print("\nProvide the device type requested (ALL|AP|SW|GW|MC): ")
   fmt.Scanln(devType)
 }
  fmt.Println("Target site name:",*siteName)
  fmt.Println("Requested device type:",*devType)

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


//prepare the output file
  ofilename := fmt.Sprintf("%s-serials.txt",*siteName)
  fmt.Printf("\nOutput file written to %s\n\n",ofilename)
  of, err := os.Create(ofilename)
    if err != nil {
        fmt.Println(err)
        return
    }

  if (*devType == "ALL") {
    fmt.Print("APs: ")
    ap_serials = get_apSerials(central_info, *siteName)
    fmt.Println(ap_serials)

    fmt.Print("Switches: ")
    sw_serials = get_switchSerials(central_info, *siteName)
    fmt.Println(sw_serials)

    fmt.Print("Gateways: ")
    gw_serials = get_gatewaySerials(central_info, *siteName)
    fmt.Println(gw_serials)

    fmt.Print("Mobility Controllers: ")
    mc_serials = get_mcSerials(central_info, *siteName)
    fmt.Println(mc_serials)
   } else if (*devType == "AP") {
    fmt.Print("APs: ")
    ap_serials = get_apSerials(central_info, *siteName)
    fmt.Println(ap_serials)
   }  else if (*devType == "SW") {
    fmt.Print("Switches: ")
    sw_serials = get_switchSerials(central_info, *siteName)
    fmt.Println(sw_serials)
   }  else if (*devType == "GW") {
    fmt.Print("Gateways: ")
    gw_serials = get_gatewaySerials(central_info, *siteName)
    fmt.Println(gw_serials)
   }  else if (*devType == "MC") {
    fmt.Print("Mobility Controllers: ")
    mc_serials = get_mcSerials(central_info, *siteName)
   }
   
  for _, serial := range ap_serials {
        _, err := of.WriteString(serial + "\n")
        if err != nil {
          fmt.Println(err)
          return
        }
    }
  for _, serial := range sw_serials {
        _, err := of.WriteString(serial + "\n")
        if err != nil {
          fmt.Println(err)
          return
        }
    }
  for _, serial := range gw_serials {
        _, err := of.WriteString(serial + "\n")
        if err != nil {
          fmt.Println(err)
          return
        }
    }
  for _, serial := range mc_serials {
        _, err := of.WriteString(serial + "\n")
        if err != nil {
          fmt.Println(err)
          return
        }
  }
  err = of.Close()
  if err != nil {
     fmt.Println(err)
     return
  }
}
