package helper

import (
	"bufio"
	"bytes"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/util/constants"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// post json data use http client
func HttpClientPostJsonData(uri string, requestBody *bytes.Buffer) (int, []byte, error) {
	url := conf.NodeUrl + uri
	res, err := http.Post(url, constants.HeaderContentTypeJson, requestBody)
	defer res.Body.Close()

	if err != nil {
		return 0, nil, err
	}

	resByte, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return 0, nil, err
	}

	return res.StatusCode, resByte, nil
}

// get data use http client
func HttpClientGetData(uri string) (int, []byte, error) {
	res, err := http.Get(conf.NodeUrl + uri)
	defer res.Body.Close()

	if err != nil {
		return 0, nil, err
	}

	resByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, nil, err
	}

	return res.StatusCode, resByte, nil
}

func ConvertStrToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// check file whether exist
// return true if exist, otherwise return false
func CheckFileExist(filePath string) (bool, error) {
	exists := true
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			exists = false
		} else {
			// unknown err
			return false, err
		}
	}
	return exists, nil
}

// create folder if not exist
// return err if not successful create
func CreateFolder(folderPath string) error {
	folderExist, err := CheckFileExist(folderPath)
	if err != nil {
		return err
	}

	if !folderExist {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteFile(filePath string, content []byte) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	fileWrite := bufio.NewWriter(file)
	_, err = fileWrite.Write(content)
	if err != nil {
		return err
	}
	fileWrite.Flush()
	return nil
}
