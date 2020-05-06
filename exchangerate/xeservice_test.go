package exchangerate

import (
	"con-currency/config"
	"con-currency/exchangerate"
	"con-currency/model"
	"testing"

	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

func TestGetExRateFromAPI(t *testing.T) {
	assert := assert.New(t)
	err := config.Init("../config")
	if err != nil {
		t.Errorf("InitJob = %d; want ", err)
		return
	}

	xeResponse := `{
    "terms": "http://www.xe.com/legal/dfs.php",
    "privacy": "http://www.xe.com/privacy.php",
    "from": "EUR",
    "amount": 1,
    "timestamp": "2020-04-27T00:00:00Z",
    "to": [
        {
            "quotecurrency": "INR",
            "mid": 76.2834881272
        }
    ]
}`

	defer gock.Off()
	gock.New("https://xecdapi.xe.com/v1/convert_from.json").
		MatchParams(map[string]string{
			"from": "EUR",
			"to":   "INR",
		}).
		MatchHeader("Authorization", "Basic c2JlbnRlcnMzNTY4NDk1MDE6cDE5bnUxN3BvcnRzbWxzdXBkZmUwbjRnY3Q=").
		Reply(200).
		JSON([]byte(xeResponse))

	xeservice := exchangerate.New()
	rates, e := xeservice.Get("EUR")

	assert.NotNil(rates)
	assert.Nil(e)
	assert.True(gock.IsDone())
	assert.IsType([]model.CurrencyRate{}, rates)
}

func TestGetExRateFromAPIFailure(t *testing.T) {
	assert := assert.New(t)

	err := config.Init("../config")
	if err != nil {
		t.Errorf("InitJob = %d; want ", err)
		return
	}
	xeResponse := `{
    "code": 7,
    "message": "No USDD found on 2020-04-27T00:00:00Z",
    "documentation_url": "https://xecdapi.xe.com/docs/v1/"
}`

	defer gock.Off()
	gock.New("https://xecdapi.xe.com/v1/convert_from.json").
		MatchParams(map[string]string{
			"from": "UST",
			"to":   "INR",
		}).
		Reply(500).
		JSON([]byte(xeResponse))

	xeservice := exchangerate.New()
	r, e := xeservice.Get("UST")
	assert.Nil(r)
	assert.NotNil(e)
}
