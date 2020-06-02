package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/jlaffaye/ftp"
	"github.com/joho/godotenv"
)

// https://mholt.github.io/json-to-go/

type autogenerated struct {
	Beobachtungen []struct {
		Lng        float64 `json:"lng"`
		Lat        float64 `json:"lat"`
		Kopfid     int     `json:"kopfid"`
		Artname    string  `json:"artname"`
		Familie    string  `json:"familie"`
		Anzahl     int     `json:"anzahl"`
		Lebensraum string  `json:"lebensraum,omitempty"`
		Bundesland string  `json:"bundesland"`
		Taxon      string  `json:"taxon"`
		Datum      string  `json:"datum"`
		GUID       string  `json:"guid"`
		Quelle     int     `json:"quelle"`
		Fundort    string  `json:"fundort"`
		Gattung    string  `json:"gattung"`
		Methode    int     `json:"methode"`
		Ordnung    string  `json:"ordnung"`
		Geschaetzt string  `json:"geschaetzt,omitempty"`
	} `json:"beobachtungen"`
	Status  int    `json:"status"`
	Offset  int    `json:"offset"`
	Meldung string `json:"meldung"`
	More    bool   `json:"more"`
}

func readCSVFromURL(url string) []map[string]interface{} {
	log.Println("Start reading data from url.")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Response status:", resp.Status)

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var obj map[string]interface{}
	jsonErr := json.Unmarshal(body, &obj)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	var target []map[string]interface{}
	for _, v := range obj["beobachtungen"].([]interface{}) {
		target = append(target, v.(map[string]interface{}))

	}

	return target
}

func main() {
	var config string

	flag.StringVar(&config, "config", "", "Path to config.")
	flag.Parse()

	err := godotenv.Load(config)

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	url := os.Getenv("URL")
	action := os.Getenv("AKTION")
	ftpHost := os.Getenv("FTP_HOST")
	ftpUser := os.Getenv("FTP_USER")
	ftpPassword := os.Getenv("FTP_PASSWORD")
	ftpPath := os.Getenv("FTP_PATH")

	data := readCSVFromURL(url)
	buf := exportXLSX(data, action, nil)
	storeToFTP(ftpHost, ftpUser, ftpPassword, ftpPath, action+".xlsx", *buf)
	//_, errWrite := buf.WriteTo(os.Stdout)
	//errCheck(errWrite)

}

func storeToFTP(ftpHost string, ftpUser string, ftpPassword string, ftpPath string, action string, buf bytes.Buffer) {
	log.Println("Store to ftp.")
	c, err := ftp.Dial(ftpHost, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Login(ftpUser, ftpPassword)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Stor(ftpPath+action, &buf)
	if err != nil {
		panic(err)
	}

	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}

}

func exportXLSX(records []map[string]interface{}, sheetName string, header []string) *bytes.Buffer {

	log.Println("Start exporting to xlsx.")
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet(sheetName)

	if sheetName != "Sheet1" {
		xlsx.DeleteSheet("Sheet1")
	}

	if header == nil {
		header = createHeader(records)
	}

	colNames := make([]string, len(header))

	for i := range header {
		colNames[i], _ = excelize.ColumnNumberToName(i + 1)
	}
	// Set first row as header
	for i, entry := range header {
		xlsx.SetCellStr(sheetName, colNames[i]+"1", entry)
	}

	for rowIndex, entry := range records {
		rowS := strconv.Itoa(rowIndex + 2)
		for cellIndex, name := range header {
			if _, found := entry[name]; found {
				switch t := entry[name].(type) {
				case string:
					xlsx.SetCellStr(sheetName, colNames[cellIndex]+rowS, t)
				case float64:
					xlsx.SetCellStr(sheetName, colNames[cellIndex]+rowS, fmt.Sprint(t))
				case int:
					xlsx.SetCellStr(sheetName, colNames[cellIndex]+rowS, fmt.Sprint(t))
				default:
					log.Println(t)
				}
			} else {
				xlsx.SetCellStr(sheetName, colNames[cellIndex]+rowS, "")
			}
			cellIndex++
		}
	}
	xlsx.SetActiveSheet(index)

	buf, _ := xlsx.WriteToBuffer()
	log.Println("Finished exporting to xlsx.")
	return buf

}

func createHeader(entries []map[string]interface{}) []string {
	uniqueMap := make(map[string]string)
	for _, entry := range entries {
		for name := range entry {
			uniqueMap[name] = ""

		}
	}
	var names []string
	for k := range uniqueMap {
		names = append(names, k)

	}
	sort.Strings(names)
	return names
}

func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}