package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/kaifei-bianjie/mock/util/constants"
	"github.com/kaifei-bianjie/mock/util/helper"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

func GenKeyName(namePrefix string, num int) string {
	uid := uuid.NewV4().String()
	return fmt.Sprintf("%s_%v_%v", namePrefix, uid, num)
}

// create key
func CreateAccount(name, password, seed string) (string, error) {
	req := types.KeyCreateReq{
		Name:     name,
		Password: password,
		Seed:     seed,
	}

	uri := constants.UriKeyCreate

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	reqBody := bytes.NewBuffer(reqBytes)

	statusCode, resBytes, err := helper.HttpClientPostJsonData(uri, reqBody)

	if err != nil {
		return "", err
	}

	if statusCode == constants.StatusCodeOk {
		res := types.KeyCreateRes{}
		if err := json.Unmarshal(resBytes, &res); err != nil {
			return "", nil
		}
		return res.Address, nil
	} else if statusCode == constants.StatusCodeConflict {
		return "", fmt.Errorf("%v", string(resBytes))
	} else {
		errRes := types.ErrorRes{}
		if err := json.Unmarshal(resBytes, &errRes); err != nil {
			return "", err
		}
		return "", fmt.Errorf("err code: %v, err msg: %v", errRes.Code, errRes.ErrorMessage)
	}
}

// create key
func CreateAccountByCmd(name, password string, home string) (string, error) {
	cmdStr := constants.KeysAddCmd + name + " --home=" + home
	cmd := getCmd(cmdStr, nil)
	//cmd = exec.Command(constants.KeysAddCmd + name + " --home=" + home)
	stdin, _ := cmd.StdinPipe()
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
	}
	log.Printf("Executing command `%v` with arguments: `%v`", cmdStr, constants.KeyPassword)
	//log.Printf("Waiting for command to finish...")

	stdin.Write([]byte(constants.KeyPassword + "\n"))
	stdin.Write([]byte(constants.KeyPassword + "\n"))

	errmsg, _ := ioutil.ReadAll(stderr)
	output, _ := ioutil.ReadAll(stdout)

	if err := cmd.Wait(); err != nil {
		log.Printf("Command finished with error: %v, %v", err.Error(), string(errmsg))
		return "", err
	}
	msg := string(output)
	//log.Printf("Command finished with response: %v", msg)

	if strings.Contains(msg, "It is the only way to recover your account") {
		index := strings.Index(msg, "local") + 6
		address := string(msg[index : index+42])
		log.Printf("Successfully create account %v", name)
		return address, nil
	}

	return "", fmt.Errorf("the responseBody is wrong during the check process")
}

// get account info
func GetAccountInfo(address string) (types.AccountInfoRes, error) {
	var (
		accountInfo types.AccountInfoRes
	)
	uri := fmt.Sprintf(constants.UriAccountInfo, address)
	statusCode, resByte, err := helper.HttpClientGetData(uri)

	if err != nil {
		return accountInfo, err
	}

	if statusCode == constants.StatusCodeOk {
		if err := json.Unmarshal(resByte, &accountInfo); err != nil {
			return accountInfo, err
		}
		return accountInfo, nil
	} else {
		return accountInfo, fmt.Errorf("status code is not ok, code: %v", statusCode)
	}
}

// get account address by name
func GetAccAddr(name string, home string) (string, error) {
	cmdStr := constants.KeysShowCmd + name + " --home=" + home
	cmd := getCmd(cmdStr, nil)
	//cmd = exec.Command(constants.KeysAddCmd + name + " --home=" + home)
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		log.Println(err)
	}

	errmsg, _ := ioutil.ReadAll(stderr)
	output, _ := ioutil.ReadAll(stdout)

	if err := cmd.Wait(); err != nil {
		log.Printf("Command finished with error: %v, %v", err.Error(), string(errmsg))
		return "", err
	}
	msg := string(output)
	//log.Printf("Command finished with response: %v", msg)

	if strings.Contains(msg, "ADDRESS:") {
		index := strings.Index(msg, "local") + 6
		address := string(msg[index : index+42])
		log.Printf("Found the account %v with address %v", name, address)
		return address, nil
	}

	return "", fmt.Errorf("the responseBody is wrong during the check process")
	/*	var (
			accountInfo types.KeyInfo
		)
		uri := fmt.Sprintf(constants.UriKeyInfo, name)
		statusCode, resByte, err := helper.HttpClientGetData(uri)

		if err != nil {
			return accountInfo, err
		}

		if statusCode == constants.StatusCodeOk {
			if err := json.Unmarshal(resByte, &accountInfo); err != nil {
				return accountInfo, err
			}
			return accountInfo, nil
		} else {
			return accountInfo, fmt.Errorf("status code is not ok, code: %v", statusCode)
		}*/
}

func getCmd(command string, escapeParams []string) *exec.Cmd {
	split := strings.Split(command, " ")

	escape := false
	tempStr := ""
	var cmdArray []string
	for i := 0; i < len(split); i++ {
		s := split[i]

		if !escape {
			for _, e := range escapeParams {
				escape = strings.Index(s, e) == 0
				if escape {
					break
				}
			}
		}

		if escape {
			if tempStr == "" {
				println(tempStr)

				tempStr = s
				continue
			}

			// TODO 只根据"--"开头来判断是否已结束当前escape的参数，可能会导致bug
			// TODO Repair of bug, Conside Restructure
			if strings.Index(s, "--") == 0 {
				cmdArray = append(cmdArray, escapeQuotes(tempStr))
				escape = false
				tempStr = ""
				i--
				continue
			}
			tempStr = tempStr + " " + s
		} else {
			cmdArray = append(cmdArray, s)
		}
	}

	if escape {
		cmdArray = append(cmdArray, escapeQuotes(tempStr))
	}

	var cmd *exec.Cmd
	if len(cmdArray) == 1 {
		cmd = exec.Command(cmdArray[0])
	} else {
		cmd = exec.Command(cmdArray[0], cmdArray[1:]...)
	}

	// fmt.Println(cmd.Args)

	return cmd
}

func escapeQuotes(tempStr string) string {
	tempStr = strings.Replace(tempStr, "\\\"", "$escape$", -1)
	tempStr = strings.Replace(tempStr, "\"", "", -1)
	tempStr = strings.Replace(tempStr, "$escape$", "\"", -1)
	return tempStr
}
