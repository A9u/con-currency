package test

import (
	"con-currency/config"
	"con-currency/model"
	"con-currency/xeservice"
	"reflect"
	"testing"

	"github.com/nbio/st"
	gock "gopkg.in/h2non/gock.v1"
)

func TestGetExRateFromAPI(t *testing.T) {
	err := config.InitConfig("../config")
	if err != nil {
		t.Errorf("InitJob = %d; want ", err)
		return
	}

	xeResponse := `{
    "terms": "http://www.xe.com/legal/dfs.php",
    "privacy": "http://www.xe.com/privacy.php",
    "from": "USD",
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
		MatchHeader("Authorization", "Basic c3BiZW50MzQzNDM2MDUwOmVma3YyczEyM2RsOTNzc2U3MXRtajhyM243").
		MatchParams(map[string]string{
			"from": "USD",
			"to":   "INR",
		}).
		Reply(200).
		JSON([]byte(xeResponse))

	r, e := xeservice.GetExRateFromAPI("USD")
	st.Expect(t, reflect.TypeOf(r), reflect.TypeOf(model.XEcurrency{}))
	st.Expect(t, e, nil)

	st.Expect(t, gock.IsDone(), true)
}

func TestGetExRateFromAPIFailure(t *testing.T) {
	err := config.InitConfig("../config")
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
		Reply(301).
		JSON([]byte(xeResponse))

	r, e := xeservice.GetExRateFromAPI("UST")
	st.Expect(t, r, model.XEcurrency{})
	st.Expect(t, (e != nil), true)

	st.Expect(t, gock.IsDone(), true)
}
