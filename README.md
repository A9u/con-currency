# XE_Currency
The XE Currency Data API is a REST-ful (or REST-like, depending how strictly you interpret REST) web-services API.


Please ensure you have config.json file to run the project.

1. # Add file config.json
2. {
    "xe_account": {
       "xe_url": "https://xecdapi.xe.com/v1/convert_from.json/",
       "xe_account_id": "xe_account_id",
       "xe_account_key": "xe_account_key"
    },
    "postgres": {
        "host":"localhost",
        "port":"5432",
        "user":"user",  // change username
        "password":"password", //change password
        "dbname":"dbname",
        "sslmode":"disable"
    },
    "currency":[
        "AED", "CUP", "AFN"
    ]
}
3. Run using command: go run main.go



# Benchmark_InitJob-4     2758207762 ns/op
PASS
ok      XE_Currency     2.766s

# For 10 Currency Total time taken:=2.31708774s
