package exchangerate

import (
	"con-currency/config"
	"con-currency/model"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	logger "github.com/sirupsen/logrus"
)

//XEService is model of the XE credentials
type xeServiceConfig struct {
	url      string
	username string
	password string
}

type xeResponse struct {
	terms     string    `json:"terms"`
	privacy   string    `json:"privacy"`
	from      string    `json:"from"`
	amount    float64   `json:"amount"`
	timestamp time.Time `json:"timestamp"`
	to        []struct {
		quotecurrency string  `json:"quotecurrency"`
		mid           float64 `json:"mid"`
	} `json:"to"`
}

type XeService struct {
	config xeServiceConfig
}

// New creates a instance of XeService
func New() Converter {
	xeConfig := xeServiceConfig{}
	xeConfig.username = config.GetConfigString("api_config.xe_account_id")
	xeConfig.password = config.GetConfigString("api_config.xe_account_key")
	xeConfig.url = config.GetConfigString("api_config.xe_url")
	return &XeService{xeConfig}
}

//GetConverter fetches the currency rate with respect to other currencies
func (converter *XeService) Get(currency string) (rates []model.CurrencyRate, err error) {

	var xeResp xeResponse
	client := &http.Client{}
	req, err := http.NewRequest("GET", converter.config.url, nil)
	if err != nil {
		logger.WithField("err", err.Error()).Error("http New Request Failed")
		return
	}

	query := req.URL.Query()

	query.Add("to", "INR") // '*' for all
	query.Add("from", currency)

	req.URL.RawQuery = query.Encode()

	req.SetBasicAuth(converter.config.username, converter.config.password)

	resp, err := client.Do(req)
	if err != nil {
		logger.WithField("err", err.Error()).Error("API Call Failed")
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Reading Response Failed")
		return
	}

	// for unsuccessful response
	if resp.StatusCode != http.StatusOK {
		errResp := model.ErrorResponse{}
		err = json.Unmarshal(respBody, &errResp)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Reading Response Failed")
			return
		}
		// warning
		logger.WithField("msg", "XE return error").Info("XE error")

		return nil, errors.New(errResp.Message)
	}

	// for successful response
	err = json.Unmarshal(respBody, &xeResp)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Unmarshal Failed")
		return
	}

	for _, value := range xeResp.to {
		rates = append(rates, model.CurrencyRate{xeResp.from, value.quotecurrency, value.mid, xeResp.timestamp})
	}
	return rates, nil
}
