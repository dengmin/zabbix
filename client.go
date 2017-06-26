package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type API struct {
	url    string
	user   string
	passwd string
	id     int
	auth   string
	Client *http.Client
}

type ZabbixRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Auth    string      `json:"auth,omitempty"`
}

type ZabbixResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Result  interface{} `json:"result"`
	Error   ZabbixError `json:"error"`
}

type ZabbixError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (api *API) GetToken() string {
	return api.auth
}

func (ze *ZabbixError) Error() string {
	return ze.Data
}

func ZabbixApi(server, username, password string) (*API, error) {
	return &API{server, username, password, 0, "", &http.Client{}}, nil
}

func (api *API) Call(method string, params interface{}) (ZabbixResponse, error) {
	id := api.id
	api.id = api.id + 1
	reqparams := ZabbixRequest{"2.0", id, method, params, api.auth}
	zabbix_request, err := json.Marshal(reqparams)
	fmt.Printf("request Params %+v\n", reqparams)
	if err != nil {
		return ZabbixResponse{}, err
	}
	request, err := http.NewRequest("POST", api.url, bytes.NewBuffer(zabbix_request))
	if err != nil {
		fmt.Printf("Request Error: %s\n", err)
		return ZabbixResponse{}, err
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := api.Client.Do(request)
	if err != nil {
		fmt.Printf("Request Error: %s\n", err)
		return ZabbixResponse{}, err
	}
	defer response.Body.Close()

	var result ZabbixResponse
	var buf bytes.Buffer

	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(buf.Bytes(), &result)
	return result, nil
}

func (api *API) Login() (bool, error) {
	// 用户登录
	params := make(map[string]string)
	params["user"] = api.user
	params["password"] = api.passwd
	response, err := api.Call("user.login", params)
	if err != nil {
		fmt.Printf("Login Error: %s\n", err)
		return false, err
	}
	if response.Error.Code != 0 {
		return false, &response.Error
	}
	api.auth = response.Result.(string)
	return true, nil
}

func (api *API) Logout() (bool, error) {
	params := make(map[string]string)
	response, err := api.Call("user.logout", params)
	if err != nil {
		fmt.Printf("Logout Error: %s\n", err)
		return false, err
	}
	if response.Error.Code != 0 {
		return false, &response.Error
	}
	return true, nil
}
