package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	infoLog  = "INFO:"
	errorLog = "ERROR:"
)

type FullData struct {
	DomainName string     `json:"domain"`
	ScanDate   string     `json:"scan_date"`
	Registrar  Registrar  `json:"registrar"`
	Registrant Registrant `json:"registrant"`
}

type Data struct {
	Registrar  Registrar  `json:"registrar"`
	Registrant Registrant `json:"registrant"`
}

type Registrar struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Referral_url string `json:"referral_url"`
}

type Registrant struct {
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Street       string `json:"street"`
	City         string `json:"city"`
	Province     string `json:"province"`
	Postal_code  string `json:"postal_code"`
	Country      string `json:"country"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}

func lastIndex(jsonFileName string, jsonExtension string) (int, error) {
	lastIndex := 0

	workingDirectory, err := os.Getwd()
	if err != nil {
		return 0, err
	}

	f, err := os.Open(workingDirectory)
	if err != nil {
		return 0, err
	}

	files, err := f.Readdir(0)
	if err != nil {
		return 0, err
	}

	pattern := fmt.Sprintf("%s-*.%s", jsonFileName, jsonExtension)
	for _, file := range files {
		match, err := filepath.Match(pattern, file.Name())
		if err != nil {
			return 0, err
		}

		if match {
			tempIndex, err := strconv.Atoi(strings.Trim(file.Name(), pattern))
			if err != nil {
				return 0, err
			}

			if tempIndex > lastIndex {
				lastIndex = tempIndex
			}
		}
	}

	if lastIndex > 0 {
		return lastIndex, nil
	}

	return 1, nil
}

func domainFileSizeCheck(domainFileName string) error {
	fileInfo, err := os.Stat(domainFileName)
	if err != nil {
		return err
	}

	if fileInfo.Size() == 0 {
		return errors.New("domain file is empty")
	}

	return nil
}

func errorLogWithExit(err error) {
	if err != nil {
		log.Printf("%s %s", errorLog, err)
		os.Exit(1)
	}
}

func main() {
	data := Data{}
	client := &http.Client{}
	err := godotenv.Load(".dScan")
	errorLogWithExit(err)

	jsonFileName := os.Getenv("DSCAN_JSON_FILE_NAME")
	jsonExtension := os.Getenv("DSCAN_JSON_EXTENSION")
	jsonSizeLimit, err := strconv.ParseInt(os.Getenv("DSCAN_JSON_SIZE_LIMIT"), 10, 64)
	errorLogWithExit(err)

	jsonFileLastIndex, err := lastIndex(jsonFileName, jsonExtension)
	errorLogWithExit(err)

	logFile, err := os.OpenFile(os.Getenv("DSCAN_LOG_FILE_NAME"), os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	errorLogWithExit(err)

	log.SetOutput(logFile)
	defer logFile.Close()

	err = domainFileSizeCheck(os.Getenv("DSCAN_DOMAIN_FILE"))
	errorLogWithExit(err)

	domainFile, err := os.OpenFile(os.Getenv("DSCAN_DOMAIN_FILE"), os.O_RDONLY, os.ModePerm)
	errorLogWithExit(err)

	defer domainFile.Close()
	df := bufio.NewReader(domainFile)

	jsonFullFileName := fmt.Sprintf("%s-%d.%s", jsonFileName, jsonFileLastIndex, jsonExtension)
	jsonFile, err := os.OpenFile(jsonFullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	errorLogWithExit(err)

	defer jsonFile.Close()

	for {
		fileStats, _ := jsonFile.Stat()
		if fileStats.Size() > jsonSizeLimit {
			jsonFile.Close()
			jsonFileLastIndex++
			jsonFullFileName := fmt.Sprintf("%s-%d.%s", jsonFileName, jsonFileLastIndex, jsonExtension)

			jsonFile, err = os.OpenFile(jsonFullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			errorLogWithExit(err)
		}

		line, _, err := df.ReadLine()
		opLine := string(line)
		if err != nil {
			log.Printf("%s %s", infoLog, err)
			if err == io.EOF {
				break
			}
		}

		requestUrl := "https://whoisjsonapi.com/v1/www." + opLine
		req, err := http.NewRequest("GET", requestUrl, nil)
		if err != nil {
			log.Printf("%s domainName: %s", errorLog, opLine)
		}
		req.Header.Set("Authorization", os.Getenv("DSCAN_WHOIS_API_TOKEN"))

		response, err := client.Do(req)
		if err != nil {
			jsonFile.Close()
			rmError := os.Remove(jsonFullFileName)
			if rmError != nil {
				log.Printf("%s %s", errorLog, err)
			}
			errorLogWithExit(err)
		}

		defer response.Body.Close()

		responseBody, err := io.ReadAll(response.Body)
		errorLogWithExit(err)

		err = json.Unmarshal(responseBody, &data)
		if err != nil {
			errorLogWithExit(err)
		} else {
			currentTime := time.Now()

			var fullData = FullData{
				DomainName: opLine,
				ScanDate:   currentTime.Format("2006-01-02"),
				Registrar:  data.Registrar,
				Registrant: data.Registrant,
			}

			fullDataJSON, err := json.Marshal(fullData)
			errorLogWithExit(err)

			_, err = jsonFile.WriteString(string(fullDataJSON) + "\n")
			errorLogWithExit(err)
		}
	}
}
