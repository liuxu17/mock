package helper

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/util/constants"
	"strconv"
)

// post json data use http client
func HttpClientPostJsonData(uri string, requestBody *bytes.Buffer) (int, []byte, error) {
	url := conf.BlockChainNodeServerUrl + uri
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
	res, err := http.Get(conf.BlockChainNodeServerUrl + uri)
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

func ConvertStrToInt64(s string) (int64, error)  {
	return strconv.ParseInt(s, 10, 64)
}
