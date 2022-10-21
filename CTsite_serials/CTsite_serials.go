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
)

type central struct {
    base_url string
    customer_id string
    token string
}

var appName = "CTsite_serial"
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

func get_apSerials(central_info central, siteName string)[]string {

  my_serials := []string{}

  access_token := central_info.token
  base_url := central_info.base_url
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

func get_switchSerials(central_info central, siteName string)[]string {

  my_serials := []string{}

  access_token := central_info.token
  base_url := central_info.base_url
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

func get_gatewaySerials(central_info central, siteName string)[]string {

  my_serials := []string{}

  access_token := central_info.token
  base_url := central_info.base_url
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

func get_mcSerials(central_info central, siteName string)[]string {

  my_serials := []string{}

  access_token := central_info.token
  base_url := central_info.base_url
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



func main() {

  pgmDescription:= fmt.Sprintf("%s: Retrieve a list of device serial numbers for a given site.",appName)
  parser := argparse.NewParser("CTsite_serials",pgmDescription)
  token := parser.String("","token", &argparse.Options{Help: "Central API token",Required: true})
  url := parser.String("","url", &argparse.Options{Help: "Central API URL",Required: true})
  siteName := parser.String("","site", &argparse.Options{Help: "Target site name",Required: true})
  devType := parser.String("","dev_type", &argparse.Options{Help: "Device type requested: ALL|AP|SW|GW|MC",Required: false})

  ap_serials := []string{}
  sw_serials := []string{}
  gw_serials := []string{}
  mc_serials := []string{}

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

  fmt.Println("Target site name:",*siteName)
  fmt.Println("Requested device type:",*devType)

//======================================================
// test if valid token
//======================================================
  respCode := test_central(central_info)
  fmt.Printf("Central Status: %s(%d)\r\n",http.StatusText(respCode),respCode)
  if (respCode != 200) { os.Exit(3)}

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
