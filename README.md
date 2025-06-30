# BACKEND ENVIRONMENT EWALLET

This project implement basic setup for backend of ewallet application using Go. This project utilizes Gin-Gonic as web framework, Pgx for database integration with PostgreSQL. This project mainly use MVC pattern for better separation of concerns, which helps in building more organized and mantainable application.

## API Endpoints Overview

```go
//authentication path
//first checkpoint to get into this app, make sure to register before do anyhting else
/auth/register	//used for create new user
/auth/login		//used for login and get access to transaction

//transaction path
//protected by token authorization, obtained from auth/login
/transactions/topup		//first transaction to do, since default balance user is set to 0 in the first time
/transactions/transfer	//transfer balance to other users, success if balance is enough, make sure to topup in advance
/transactions/history	//retrieve all history transaction of transfers and topups
```

## How to run this project
1. Clone this project
```sh
git clone https://github.com/yasirmaxstyle/fgo24-be-weeklytask.git .
```
2. Install `gow`for hot reload running
```go
go install github.com/mitranim/gow@latest

//move path for enabling command gow
export GOPATH="$HOME/go"
export PATH="$GOPATH/bin:$PATH"
```
3. Run postgres on Docker
```sh
docker pull postgres
docker run -e PASSWORD_POSTGRES=1 -p 5432:5432 -d postgres
```

## Technologies and Dependencies
1. Go
2. PostgreSQL
3. Docker
4. [Gin](https://github.com/gin-gonic/gin) (Web Framework)
5. [PGX](https://github.com/jackc/pgx) (database integration)
6. [Gow](https://github.com/mitranim/gow) (hot reload)

## How to take part in this project
You are free to fork this project, make improvement and submit a pull request to improve this project. If you find this useful or if you have suggestion, you can start discussing through my social media below.
- [Instagram](https://www.instagram.com/yasirmaxstyle/)
- [LinkedIn](https://www.linkedin.com/in/muhamad-yasir-806230117/)

## License
This project is under MIT License
