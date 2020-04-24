# con-currency
    con-currency fetches currency data and save it to database.
    

    Please ensure that you create a 'config.json' file to run the project in its root location.
    
    Config file example →

{

    "api_config": {
    
        "xe_url": "https://xecdapi.xe.com/v1/convert_from.json/",
        
        "xe_account_id": "xe_account_id",
        
        "xe_account_key": "xe_account_key"
        
    },
    
    "db_config": {
    
        "host": "host",
        
        "port": "port",
        
        "user": "user",    
        
        "password":"password", 
        
        "dbname":"dbname", 
        
        "sslmode":"disable"
    }, 
    
    "currency_list":[ "USD" ] 
}

# Run using command: go run main.go



# Benchmark: 
took 15 seconds for 162 currencies.
