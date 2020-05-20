package exchangerate

import (
	"con-currency/config"
	"con-currency/model"
	"encoding/json"
	"errors"

	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

//XEService is model of the XE credentials
type xeServiceConfig struct {
	url      string
	username string
	password string
}

type xeResponse struct {
	Terms     string    `json:"terms"`
	Privacy   string    `json:"privacy"`
	From      string    `json:"from"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	To        []struct {
		Quotecurrency string  `json:"quotecurrency"`
		Mid           float64 `json:"mid"`
	} `json:"to"`
}

type XeService struct {
	config xeServiceConfig
}

// New creates a instance of XeService
func New() Converter {
	xeConfig := xeServiceConfig{}
	xeConfig.username = config.GetString("api_config.xe_account_id")
	xeConfig.password = config.GetString("api_config.xe_account_key")
	xeConfig.url = config.GetString("api_config.xe_url")

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

	query.Add("to", "*") // '*' for all
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
		logger.WithField("msg", errResp.Message).Info("XE error")

		return nil, errors.New(errResp.Message)
	}

	// for successful response

	err = json.Unmarshal(respBody, &xeResp)
	if err != nil {

		logger.WithField("err", err.Error()).Error("Unmarshal Failed")
		return
	}

	for _, value := range xeResp.To {
		rates = append(rates, model.CurrencyRate{xeResp.From, value.Quotecurrency, value.Mid, xeResp.Timestamp})
	}
	return rates, nil
}
