package wallet

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"unsafe"

	"github.com/qiniupd/qiniu-go-sdk/x/log.v7"
)

type stateInfo struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Worker string `json:"Worker"`
		Owner  string `json:"Owner"`
	}
}

// return owner worker addr ，相同返回true不同false
func ReqStateInfo(miner string, daemon string) (workerID string, ownerID string, ok bool, err error) {
	jsonInfo := `{ "jsonrpc": "2.0", "method": "Filecoin.StateMinerInfo", "params": ["` + miner + `",null], "id": 1 }`
	reader := bytes.NewReader([]byte(jsonInfo))
	request, err := http.NewRequest("POST", "http://"+daemon+":1234/rpc/v0", reader)
	if err != nil {
		//fmt.Println(err.Error())
		log.Error(err)
		return "", "", false, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		//fmt.Println(err.Error())
		log.Error(err)
		return "", "", false, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println(err.Error())
		log.Error(err)
		return "", "", false, err
	}
	str := (*string)(unsafe.Pointer(&respBytes))
	//log.Infof("StateMinerInfo: %v", *str)
	var state stateInfo
	if err := json.Unmarshal([]byte(*str), &state); err != nil {
		log.Error(err)
	}
	if state.Result.Worker == state.Result.Owner {
		return state.Result.Worker, state.Result.Owner, true, nil
	}
	return state.Result.Worker, state.Result.Owner, false, nil
}

type WalletBalanceStruct struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	Id      int    `json:"id"`
}

func ReqBalance(WalletID string, daemonIP string) (jsonStr string, err error, balance float64) {
	jsonstr := `{ "jsonrpc": "2.0", "method": "Filecoin.WalletBalance", "params": ["` + WalletID + `"], "id": 1 }`
	reader := bytes.NewReader([]byte(jsonstr))
	request, err := http.NewRequest("POST", "http://"+daemonIP+":1234/rpc/v0", reader)
	if err != nil {
		//fmt.Println(err.Error())
		log.Error(err)
		return "", err, 0
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		//fmt.Println(err.Error())
		log.Error(err)
		return "", err, 0
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println(err.Error())
		log.Error(err)
		return "", err, 0
	}
	str := (*string)(unsafe.Pointer(&respBytes))
	//log.Infof("WalletBalance: %v", *str)
	var walletString WalletBalanceStruct
	if err := json.Unmarshal([]byte(*str), &walletString); err != nil {
		log.Error(err)
	}
	resultInt, _ := strconv.ParseFloat(walletString.Result, 64)
	b := resultInt / 1000000000000000000
	return *str, nil, b
}
