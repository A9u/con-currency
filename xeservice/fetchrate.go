package xeservice

import (
	"con-currency/config"
	"con-currency/model"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	logger "github.com/sirupsen/logrus"
)

//GetExRateFromAPI fetches the currency rate with respect to other currencies
func GetExRateFromAPI(currency string) (xeResp model.XEcurrency, err error) {

	url := config.GetConfig("api_config.xe_url")
	username := config.GetConfig("api_config.xe_account_id")
	password := config.GetConfig("api_config.xe_account_key")

	client := &http.Client{}
	//+"?to=INR&from="+currency
	req, err := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	q.Add("to", "INR") // '*' for all
	q.Add("from", currency)
	req.URL.RawQuery = q.Encode()
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		logger.WithField("err", err.Error()).Error("API Call Failed")
		return
	}
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Reading Response Failed")
		return
	}

	// for unsuccessful response
	if resp.StatusCode != http.StatusOK {
		errResp := model.ErrorResponse{}
		err = json.Unmarshal(r, &errResp)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Reading Response Failed")
			return
		}

		return model.XEcurrency{}, errors.New(errResp.Message)
	}

	// for successful response
	err = json.Unmarshal(r, &xeResp)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Unmarshal Failed")
		return
	}

	return
}
