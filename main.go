package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ResponseData struct {
	Timestamp time.Time `json:"timestamp"`
	To        []struct {
		Quotecurrency string  `json:"quotecurrency"`
		Mid           float64 `json:"mid"`
	} `json:"to"`
}

var (
	xe_account_id  string
	xe_account_key string
	xe_url         string
	con_String     string
)

const (
	conString         = "host=localhost port=5432 user=postgres password= 12345 dbname=postgres sslmode=disable"
	queryString       = `INSERT INTO exchange_rates (from_currency,to_currency,rate,created_at,updated_at) values`
	updateQueryString = `ON CONFLICT ON CONSTRAINT unq
		DO UPDATE SET rate =excluded.rate,updated_at = excluded.updated_at where exchange_rates.rate is distinct from excluded.rate`
	createTableString = `CREATE TABLE IF NOT EXISTS public.exchange_rates(
	       from_currency character varying(3) COLLATE pg_catalog."default",
	       to_currency character varying(3) COLLATE pg_catalog."default",
	       rate double precision,
	       created_at timestamp with time zone,
	   	updated_at timestamp with time zone,
	   	CONSTRAINT unq UNIQUE (from_currency, to_currency)
	   )`
)

func main() {
	start := time.Now()
	var wg sync.WaitGroup

	// list of currency, TODO → to be read from external file
	// currencyArray := []string{"AED", "CUP", "AFN", "ETB", "ALL", "AMD", "AOA", "ARS", "AZN", "BAM", "BBD", "BDT", "BGN", "IQD", "BMD", "IRR", "BIF", "BRL", "BSD", "BTN", "BYN", "CAD", "BZD", "KPW", "JOD", "COP", "CRC", "CVE", "CZK", "DOP", "DZD", "EGP", "GBP", "GEL", "AWG", "GHS", "GIP", "GTQ", "GYD", "HKD", "HNL", "HRK", "HUF", "CUC", "ILS", "IMP", "INR", "BOB", "JEP", "JMD", "KES", "KGS", "FKP", "CHF", "ERN", "GGP", "BND", "CDF", "IDR", "CLP", "GNF", "JPY", "KMF", "SPL", "PYG", "TZS", "MRU", "KYD", "KZT", "MDL", "LKR", "LRD", "LSL", "RUB", "MGA", "SHP", "MMK", "MNT", "MOP", "MUR", "MVR", "MWK", "MXN", "MYR", "NAD", "NGN", "NIO", "NPR", "NZD", "PAB", "PEN", "PGK", "PHP", "PKR", "PLN", "KWD", "RON", "RSD", "SYP", "LYD", "SAR", "SBD", "SDG", "SEK", "SGD", "TWD", "SOS", "SRD", "SZL", "TJS", "TMT", "STN", "TOP", "TRY", "TVD", "TND", "MKD", "UAH", "UGX", "UYU", "UZS", "USD", "LAK", "RWF", "KRW", "BHD", "OMR", "BWP", "XCD", "CNY", "YER", "ZAR", "ZMW", "ANG", "FJD", "GMD", "HTG", "KHR", "LBP", "MAD", "MZN", "QAR", "SCR", "SLL", "THB", "TTD", "AUD", "DKK", "NOK", "SVC", "VEF", "WST", "ZWD", "EUR", "VES", "XOF", "XPF", "DJF", "ISK", "VUV", "XAF", "VND"}
	currencyArray := []string{"INR"}

	// loads api credentials from api.env file
	err := godotenv.Load("api.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	xe_account_id = os.Getenv("xe_account_id")
	xe_account_key = os.Getenv("xe_account_key")
	xe_url = os.Getenv("xe_url")
	con_String = os.Getenv("con_String")

	// open lazy connection to database
	db, err := sql.Open("postgres", con_String)
	if err != nil {
		log.Fatal("Conn Open → ", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(0)

	// create table if not exist
	_, err = db.Exec(createTableString)
	if err != nil {
		log.Fatal("Create table → ", err)
	}

	//spawn goroutines to consume api concurrently
	for _, symbol := range currencyArray {
		go fetchData(symbol, &wg, db)
		wg.Add(1)
	}

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Operation took %s", elapsed)
}

func sum() {

}

// fetchData: Consumes api and pass the result to 'updateDb' function
func fetchData(symbol string, wg *sync.WaitGroup, db *sql.DB) {

	defer wg.Done()

	result := ResponseData{}
	client := &http.Client{}

	// creating request
	// req, err := http.NewRequest("GET", xe_url+"?to=*&from="+symbol, nil)
	req, err := http.NewRequest("GET", xe_url+"?to=INR&from="+symbol, nil)
	req.SetBasicAuth(xe_account_id, xe_account_key)
	if err != nil {
		fmt.Println("NewRequest → ", err)
	}

	// send requestk
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Send request → ", err)
	}

	if resp.StatusCode == http.StatusOK {
		// reading response
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Reading resp → ", err)
		}

		// closing response body
		if resp.Body != nil {
			resp.Body.Close()
		}

		// mapping byte array to struct
		err = json.Unmarshal(bodyText, &result)
		if err != nil {
			log.Print("Unmarshal → ", err)
		}

		// calling updateDb to fill database
		updateDB(symbol, result, db)
	} else {
		fmt.Println("API Error, Please check!")
		return
	}

}

// updateDB: Insert or update database
func updateDB(symbol string, result ResponseData, db *sql.DB) {
	values := []interface{}{}
	query := queryString

	// building query string
	for i, res := range result.To {

		values = append(values, symbol, res.Quotecurrency, res.Mid, result.Timestamp, result.Timestamp)

		numFields := 5
		n := i * numFields

		query += `(`
		for j := 0; j < numFields; j++ {
			query += `$` + strconv.Itoa(n+j+1) + `,`
		}
		query = query[:len(query)-1] + `),`

	}

	query = query[:len(query)-1]
	query += updateQueryString

	// firing query
	res, err := db.Exec(query, values...)
	if err != nil || res == nil {
		log.Fatal("Firing query → ", err)
	}

}
