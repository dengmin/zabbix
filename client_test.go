package zabbix

import "testing"
import "fmt"

func TestLogin(t *testing.T) {
	api, _ := ZabbixApi("http://10.57.17.31/api_jsonrpc.php", "min.deng", "themis@160318")
	ok, err := api.Login()
	if !ok {
		serr := err.(*ZabbixError)
		fmt.Printf("Login error: %s, %s, %s", serr.Code, serr.Message, serr.Data)
	} else {
		fmt.Println(api.GetToken())
	}
}
