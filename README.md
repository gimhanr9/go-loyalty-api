
# Square Loyalty Challenge API
This is the REST API for the Square Loyalty Challenge built using Golang.




## Requirements
Two environment files

 1 .env.development

 2 .env.production (optional)

Variables needed for environment files:

`DB_PATH=./data/loyalty.db`

`SQUARE_ACCESS_TOKEN (from Square)`

`JWT_SECRET`

`PORT=8080`

`LOCATION_ID (from Square)`


## Setup & Run

1 Clone the project.

2 Add the environment files.

3 Go to the project directory and in the terminal run command `go mod download` (This will automatically add the sqlite database file)

4 Finally run the command `go run main.go`
