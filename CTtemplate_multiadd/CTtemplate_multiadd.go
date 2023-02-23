package main


import (
    "fmt"
    "os"
//    "io/ioutil"
    "io"
    "bytes"
//    "bufio"
    "mime/multipart"
    "net/http"
    "time"
    "log"
    "strings"
    "strconv"
//    "github.com/buger/jsonparser"
    "github.com/akamensky/argparse"
//    "sigs.k8s.io/yaml"
    "github.com/micgresham/goCentral"
    "github.com/manifoldco/promptui"
)

type templateGroup struct {
    gname string
    TGrange []int 
    devtype string
    version string
    model string
    TFname string
    wired bool
    wireless bool
    }


var appName = "CTtemplate_multiadd"
var appVer = "1.0"
var appAuthor = "Michael Gresham"
var appAuthorEmail = "michael.gresham@hpe.com"
var pgmDescription = fmt.Sprintf("%s: Add one or more template groups in Central.",appName)
var central_info goCentral.Central_struct
var TG templateGroup
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

func uploadTemplate(central_info goCentral.Central_struct, gname string, Tgroup templateGroup) string {

  access_token := central_info.Token
  base_url := central_info.Base_url
  api_function_url := fmt.Sprintf("%s/configuration/v1/groups/%s/templates",base_url,gname)

  extraFields := map[string]string{}
  data, w, err := createMultipartFormData("template",Tgroup.TFname,Tgroup.TFname, extraFields )
    if err != nil {
       return("ERROR Create Multipart Form Data")
    }

  c := http.Client{Timeout: time.Duration(10) * time.Second}
  req, err := http.NewRequest("POST", api_function_url, &data)
  if err != nil {
      fmt.Printf("error %s", err)
      return("ERROR")
  }
  q := req.URL.Query()
  q.Add("name",Tgroup.TFname)
  q.Add("device_type",Tgroup.devtype)
  q.Add("version",Tgroup.version)
  q.Add("model",Tgroup.model)
  req.URL.RawQuery = q.Encode()

//  fmt.Println(req)

  req.Header.Set("Content-Type", w.FormDataContentType())
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",fmt.Sprintf(access_token)))
  resp, err := c.Do(req)
  if err != nil {
      fmt.Printf("error %s", err)
      return("ERROR")
  }

  defer resp.Body.Close()
//  body, err := ioutil.ReadAll(resp.Body)
//  fmt.Println(resp)

  return(resp.Status)
}

func createGroup(central_info goCentral.Central_struct, gname string, Tgroup templateGroup) string {

  access_token := central_info.Token
  base_url := central_info.Base_url
  api_function_url := fmt.Sprintf("%s/configuration/v2/groups",base_url)

  jsonPrep := fmt.Sprintf("{\"group\":\"%s\", \"group_attributes\":{\"template_info\":{\"Wired\": %v, \"Wireless\": %v }}}",gname,Tgroup.wired,Tgroup.wireless)
  jsonStr := []byte(jsonPrep)

  c := http.Client{Timeout: time.Duration(10) * time.Second}
  req, err := http.NewRequest("POST", api_function_url, bytes.NewBuffer(jsonStr))
  if err != nil {
      fmt.Printf("error %s", err)
      return("ERROR")
  }
  q := req.URL.Query()
  q.Add("offset","0")
  req.URL.RawQuery = q.Encode()

//  fmt.Println(req)

  req.Header.Add("Content-Type", `application/json`)
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",fmt.Sprintf(access_token)))
  resp, err := c.Do(req)
  if err != nil {
      fmt.Printf("error %s", err)
      return("ERROR")
  }

  defer resp.Body.Close()
//  body, err := ioutil.ReadAll(resp.Body)
//  fmt.Println(resp)

  return(resp.Status)
}


func substituteNumbers(str string, numbers []int) (string, error) {
	var sb strings.Builder
	var placeholderCount int

	for _, c := range str {
		if c == '@' {
			placeholderCount++
			if placeholderCount%2 == 0 {
				if len(numbers) == 0 {
					return "", fmt.Errorf("not enough numbers to substitute")
				}
				sb.WriteString(string(fmt.Sprintf("%02d",numbers[0])))
				numbers = numbers[1:]
			} 
		//	else {
		//		sb.WriteRune(c)
		//	}
		} else {
			sb.WriteRune(c)
		}
	}

	if len(numbers) > 0 {
		return "", fmt.Errorf("too many numbers to substitute")
	}

	return sb.String(), nil
}


func expandNumbers(str string) ([]int, error) {
	var numbers []int

	ranges := strings.Split(str, ",")
	for _, r := range ranges {
		bounds := strings.Split(r, "-")
		if len(bounds) == 1 {
			num, err := strconv.Atoi(bounds[0])
			if err != nil {
				return nil, fmt.Errorf("invalid number format: %v", bounds[0])
			}
			numbers = append(numbers, num)
		} else if len(bounds) == 2 {
			start, err := strconv.Atoi(bounds[0])
			if err != nil {
				return nil, fmt.Errorf("invalid number format: %v", bounds[0])
			}
			end, err := strconv.Atoi(bounds[1])
			if err != nil {
				return nil, fmt.Errorf("invalid number format: %v", bounds[1])
			}
			for i := start; i <= end; i++ {
				numbers = append(numbers, i)
			}
		} else {
			return nil, fmt.Errorf("invalid range format: %v", r)
		}
	}

	return numbers, nil
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
//  test := parser.Flag("t", "test", &argparse.Options{Help: "Enable test mode. No variables will be changed"})

  //encrypted storage setup
  SSfilename:= "../CTcentral_check/CTconfig.yml"
//  SSfilename:= "CTconfig.yml"

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


  TG.TFname = promptString("Input template file name","")

  //if the file does not exist, notify the user and exit
  if doesFileExist(TG.TFname) {
	  fmt.Println("Filename does not exist. Exiting....")
	  os.Exit(3)
  }

  TG.gname = promptString("Template group name (use @@ for range substitution placeholder)","")
  TG.TGrange, _ = expandNumbers(promptString("Template group name range",""))

  TG.devtype = promptSelect("This group will contain devices of type:",[]string{"CX", "IAP", "MobilityController", "ArubaSwitch"})
 
  TG.version = promptString("This group firmware version limit","ALL")

  TG.model = promptString("This group limited to model","ALL")

  TG.wired = yesNo("Will this group contain wired devices? : ")
  TG.wireless = yesNo("Will this group contain wireless devices? : ")

  fmt.Println("Input file name:",TG.TFname)
  fmt.Println("TG Group Name:",TG.gname)
  fmt.Println("TG Range:",TG.TGrange)
  fmt.Println("TG Device Type:",TG.devtype)
  fmt.Println("TG version:",TG.version)
  fmt.Println("TG model:",TG.model)
  fmt.Println("TG wired:",TG.wired)
  fmt.Println("TG wireless:",TG.wireless)

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

  fmt.Println("\nThe following template groups will be created:")
  for _,rangeID := range TG.TGrange {
	  result, err := substituteNumbers(TG.gname, []int{rangeID})
	  if err != nil {
		  fmt.Println(err)
	  }
	  fmt.Print(result)
	  fmt.Println(" - ",createGroup(central_info, result, TG))
  }	  

  fmt.Println("\nUploading template to each group:")
  for _,rangeID := range TG.TGrange {
	  result, err := substituteNumbers(TG.gname, []int{rangeID})
	  if err != nil {
		  fmt.Println(err)
	  }
	  fmt.Print(result)
	  fmt.Println(" - ",uploadTemplate(central_info, result, TG))
  }	  


//======================================================
// perform post check
//======================================================
//  fmt.Println("--------------------------------")
//  fmt.Println("Performing Post-Check")
//  fmt.Println("--------------------------------")
   
}
