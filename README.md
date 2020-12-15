# jsonIso
jsonIso is a golang app that handle transaction request and gonna processed to standard ISO:8583, it included database connection to save any request that it get.

## Installation
Get the program
```bash
git clone https://github.com/j03hanafi/jsonIso
cd jsonIso
```
Prepare package
```bash
go get github.com/go-sql-driver/mysql
go get github.com/gorilla/mux
go get github.com/mofax/iso8583
go get github.com/rivo/uniseg
go get github.com/stretchr/testify
```

Running program
```bash
go run *.go
```
Build and running program
```bash
go build -o {program name}
./{program name}
```

## Available Request (For Now)
#### GET
- `/payment` : Get all transaction data
- `/payment/{processingCode}` : Get transaction data based on processingCode
- `/payment/{processingCode}/iso` : Get ISO8385 format of transaction data based on processingCode
- `/payment/{processingCode}/iso8583/{element}` : Get each Data Element from ISO8583 format message

#### POST
- `/payment/iso` : Convert JSON data to ISO8583 format based on JSON request sent on body
- `/payment/json` : Convert ISO8583 message to JSON format based on ISO8583 message request sent on body
- `/payment/json/file` : Convert ISO8583 message to JSON format based on file's name sent on body and read the content of the file as ISO8583 message

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
