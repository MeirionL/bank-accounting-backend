# boing-block

boing-block is a REST API that allows users to store safe bank account and transaction information to help view it easily in VSCode, without the ability for anyone to interact with the finances themselves. [Here](https://www.youtube.com/watch?v=1T8ZMjy2xSk) is my video explaining how the project works in under 3 minutes for my CS50 final project submission.

## Motivation

This started out from an idea to be able to help my family, through eventually designing a way to have quick access to only the information I need to help them manage their finances into a budget. boing-block allows for the essentials of bank account transfers information with proper authentication, without any other unnecessary add-ons to make it more complex to use. 

## Installation

Inside a Go module:

```bash
go get github.com/MeirionL/boing-block
```

I would also reccomend downloading a REST API client extension like the [Thunder Client](https://www.thunderclient.com/) for this. The below doccumentation will reference curl commands which you can translate to your client extention

## How to start

Inside the project directory create a `.env` file with the following 3 values:

* A `PORT` address number
* A `JWT_SECRET` that can be generated at [Online JWT Generator](https://www.javainuse.com/jwtgenerator)
* A `DB_URL` for a PostgreSQL database in this format: "postgres://<postgres_username>:<postgres_username_password>@localhost:<server_port>/<database_name>?sslmode_disable"

Start the server:

```bash
go build && ./boing-block
```

## __Requests__

| Request | Description |
| --- | --- |
| GET health | Show status of server |
| GET error | Show error handler is working |
| POST user | Create user with name, password and ID |
| GET users | Show info of all users |
| GET user by ID | Show specific user based on ID parameter |
| POST login | Take in user info to return valid refresh token and access token. Refresh token is valid for 60 days and access token is valid for an hour. **One of these tokens is required for all of the below commands** |
| PUT user | Update user info with new name and/or password |
| DELETE user | Delete user and all related accounts and transactions |
| POST refresh | Create new valid access token if valid refresh token provided |
| POST revoke token | Revoke currently used refresh token |
| GET revoked tokens | Show all revoked tokens for current user |
| POST account | Create account with name and current account details |
| GET accounts | Show current accounts for user. Can filter to show an account based on bank details or account ID |
| PUT account | Update account name and/or balance |
| DELETE account | DELETE account with entered account_ID parameter |
| GET balances | Show balances of all accounts for user |
| POST transaction | Create transaction with time, type, amount, pre and post balance and info about "other account" involved |
GET transactions | Show all transactions for desired account in order of most recent. Can filter to specific transactions by providing bank details, transaction ID, "others account" ID, transactions type, or a limit of how many you wish to show |
DELETE transaction | Delete specific transaction from account |
| GET "others account" | Show all accounts that desired user account has transactions with. Can filter to find individual "others accounts" by providing bank details or "other account" ID |

### GET healthz

    curl -i -X GET "http://localhost:<PORT>/healthz"

### GET error

    curl -i -X GET "http://localhost:<PORT>/err"

### POST user

    curl -i -H 'Accept: application/json' -H 'Content-Type: application/json' -X POST -d '{"name":"Foo","password":"Bar"}' "http://localhost:<PORT>/users"

### GET users

    curl -i -X GET "http://localhost:<PORT>/users"

### GET user by ID

    curl -i -X GET "http://localhost:<PORT>/users/<userID>"

### POST login

    curl -H 'Accept: application/json' -H 'Content-Type: application/json' -X POST -d '{"id": 123,"password":"Foo"}' "http://localhost:<PORT>/login"

### PUT user

    curl -i -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'Authorization: Bearer <ACCESS_TOKEN>' -X PUT -d '{"name":"Foo","password":"Bar"}' "http://localhost:<PORT>/auth/users"

### DELETE user

    curl -i -H 'Authorization: Bearer <ACCESS_TOKEN>' -X DELETE "http://localhost:<PORT>/auth/users"

### Refresh token

    curl -i -H 'Authorization: Bearer <REFRESH_TOKEN>' -X POST "http://localhost:<PORT>/auth/refresh"

### Revoke refresh token

    curl -i -H 'Authorization: Bearer <REFRESH_TOKEN>' -X POST "http://localhost:<PORT>/auth/revoke"

### GET revoked tokens

    curl -i -H 'Authorization: Bearer <ACCESS_TOKEN>' -X GET "http://localhost:<PORT>/auth/revoke"

### POST account

    curl -i -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'Authorization: Bearer <ACCESS_TOKEN>' -X POST -d '{"account_name":"Foo","balance":123.456, "account_number":"12345678", "sort_code":"12-34-56"}' "http://localhost:<PORT>/auth/accounts"

### GET accounts

    curl -i -H 'Authorization: Bearer <ACCESS_TOKEN>' -X GET "http://localhost:<PORT>/auth/accounts"

Optionally you can filter the returned accounts by adding these optional parameters on the end of the route:

    ?account_number=12345678&sort_code=12-34-56

    ?id=<ACCOUNT_ID>

### PUT account

    curl -i -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'Authorization: Bearer <ACCESS_TOKEN>' -X PUT -d '{"id":"<ACCOUNT_ID>", "account_name":"new_account_name","balance":123.456}' "http://localhost:<PORT>/auth/accounts"

### DELETE account

    curl -i -H 'Authorization: Bearer <ACCESS_TOKEN>' -X DELETE "http://localhost:<PORT>/auth/accounts/<ACCOUNT_ID>"

### GET accounts balances

    curl -i -H 'Authorization: Bearer <ACCESS_TOKEN>' -X GET "http://localhost:<PORT>/auth/balances"

### POST transaction

    curl -i -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'Authorization: Bearer <ACCESS_TOKEN>' -X POST -d '{"transaction_time":"YYYY-MM-DDTHH:MM:SSZ","transaction_type":"outgoing","amount":12.30,"pre_balance":45.60,"post_balance":78.90,"new_account":true,"account_name":"Foo","account_number":"12345678","sort_code":"12-34-56","account_id":"<ACCOUNT_ID>"}' "http://"localhost:<PORT>/auth/transactions

__Time format placeholders breakdown:__

* YYYY is the year
* MM is the month
* DD is the day
* T splits data and time and should be left as is
* HH is the hour in 24-hour format
* Second MM is the minutes
* SS is the seconds
* Z indicates UTC time and should be left as is

### GET transactions

    curl -i -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'Authorization: Bearer <ACCESS_TOKEN>' -X GET -d '{"account_id":"<ACCOUNT_ID>"}' "http://localhost:<PORT>/auth/transactions"

Optionally you can filter the returned transactions by adding these optional parameters on the end of the route:

    ?transaction_id=<TRANSACTION_ID>

    ?transaction_type=incoming

    ?account_number=12345678&sort_code=12-34-56

    ?limit=123

    ?others_account_id=<OTHER_ACCOUNT_ID>

### DELETE transaction

    curl -i -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'Authorization: Bearer <ACCESS_TOKEN>' -X DELETE -d '{"account_id":"<ACCOUNT_ID>"}' "http://localhost:<PORT>/auth/transactions/<TRANSACTION_ID>"

### GET others account

    curl -i -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'Authorization: Bearer <ACCESS_TOKEN>' -X GET -d '{"account_id":"<ACCOUNT_ID>"}'  "http://localhost:<PORT>/auth/others"

Optionally you can filter the returned transactions by adding these optional parameters on the end of the route:

    ?account_number=12345678&sort_code=12-34-56

    ?id=<OTHERS_ACCOUNT_ID>

## Contributing

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.