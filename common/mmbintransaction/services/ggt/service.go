package customer_orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
)

var (
	// UsersEndpoint customer orchestrator get users endpoint
	UsersEndpoint = "/%s/v2/users/%s"
)

// GetUserDetails call customer orchestrator to retrieve user details
func GetUserDetails(log log.Logger, customerServiceHost, programCode, userHashID, opKey, opSecret string) (response map[string]interface{}, err error) {
	level.Info(log).Log("Method: getUserDetails")

	url := customerServiceHost + fmt.Sprintf(UsersEndpoint, programCode, userHashID)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(opKey, opSecret)

	if err != nil {
		level.Error(log).Log("Error in fetching Customer svc status of the user: ", err)
		err = errors.New("internal server error")
		return response, err
	}

	level.Info(log).Log("Sending request to customer orchestrator url ", url)

	var res *http.Response
	res, err = client.Do(req)

	if err != nil {
		level.Error(log).Log("Error in fetching Customer svc status of the user:", err)
		err = errors.New("external server error")
		return response, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		level.Error(log).Log("EStatus Response from Customer svc:", strconv.Itoa(res.StatusCode))
		return response, err
	}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		level.Error(log).Log("Error in decoding response body :", err)
		return response, err
	}

	level.Info(log).Log("Response from Customer svc: ", response)

	return response, err
}
