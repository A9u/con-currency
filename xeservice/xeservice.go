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

//XEService is model of the XE credentials
type XEService struct {
	URL      string
	Username string
	Password string
}

//GetExchangeRater is an interface
// type GetExchangeRater interface {
// 	GetExchangeRate(string) (model.XEcurrency, error)
// }

// New creates a instance of XEService
func New() (xeService XEService) {
	xeService.Username = config.GetConfigString("api_config.xe_account_id")
	xeService.Password = config.GetConfigString("api_config.xe_account_key")
	xeService.URL = config.GetConfigString("api_config.xe_url")
	return
}

//GetExchangeRate fetches the currency rate with respect to other currencies
func (xeService XEService) GetExchangeRate(currency string) (xeResp model.XEcurrency, err error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", xeService.URL, nil)
	if err != nil {
		logger.WithField("err", err.Error()).Error("http New Request Failed")
		return
	}

	query := req.URL.Query()

	query.Add("to", "INR") // '*' for all
	query.Add("from", currency)

	req.URL.RawQuery = query.Encode()

	req.SetBasicAuth(xeService.Username, xeService.Password)

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
		// warning
		logger.WithField("msg", "XE return error").Info("XE error")

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
