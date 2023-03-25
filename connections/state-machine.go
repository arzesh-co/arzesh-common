package connections

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type MachineStates struct {
	Id                string            `json:"id"`
	Account           string            `json:"acnt_uuid"`
	MachineKey        string            `json:"machine_key"`
	RefID             string            `json:"ref_uuid"`
	CurrentStateID    string            `json:"current_state_uuid"`
	CurrentStateKey   string            `json:"current_state_key"`
	CurrentStateTitle map[string]string `json:"current_state_title"`
	VersionNo         int16             `json:"version_no"`
	CreatedAt         int64             `json:"created_at"`
	CreatedBy         string            `json:"created_by"`
	Status            string            `json:"status"`
}

type StartMachineResponse struct {
	Data  MachineStates  `json:"data"`
	Error map[string]any `json:"error"`
}
type MachineBody struct {
	RefId string `json:"ref_uuid"`
}

func startMachineApi(account string, key string, ref string) (*MachineStates, string) {
	QBody := &MachineBody{}
	QBody.RefId = ref
	body, _ := json.Marshal(QBody)
	req, err := http.NewRequest(
		"POST", os.Getenv("stateMachineApi")+"/api/state-machine/machines/key/"+key+"/start", bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, "REF.CANNOT_CONNECT"
	}
	req.Header.Set("Accept", "/")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("account_uuid", account)
	req.Close = true
	client := &http.Client{Transport: &http.Transport{}}
	res, err := client.Do(req)
	if err != nil {
		return nil, "REF.CANNOT_CONNECT"
	}
	defer res.Body.Close()
	Machine := &StartMachineResponse{}
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "REF.INVALIDATION_ERROR"
	}
	err = json.Unmarshal(responseBody, Machine)
	if err != nil {
		return nil, "REF.INVALIDATION_ERROR"
	}

	errorKey, ok := Machine.Error["status_key"]
	if ok {
		return nil, errorKey.(string)
	}

	return &Machine.Data, ""
}

func StartMachine(account string, key string, ref string) (*MachineStates, string) {
	if account == "" {
		return nil, "REF.INVALIDATION_ERROR"
	}
	if key == "" {
		return nil, "REF.INVALIDATION_ERROR"
	}
	if ref == "" {
		return nil, "REF.INVALIDATION_ERROR"
	}
	return startMachineApi(account, key, ref)
}
