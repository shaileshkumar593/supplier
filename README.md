# Overview

## Requirements

1.  Go
2. Postgres
3.  Redis

---

## Setup

1. Clone/Fork this repo
2. Clean-up packages, run **go mod tidy**
3. Update the configuration values based on your local settings in **env.local**

---

## Migrations

1. Run `go run migrations/migration.go migrate`

*For new migrations that you want to add on your service*
*Below will automatically add to your migration directory/db driver specified on config* **(example: migrations/svc/mysql or migrations/svc/postgres)**

1. Run `go run migrations/migration.go create <filename>`

---

## Starting your application

1. Run `go run main.go`

---

## Test

Open [http://127.0.0.1:7001/heartbeat](http://127.0.0.1:7001/heartbeat) on your web browser
 
---

## Running Test Cases

Run `sh test.sh`

## Do changes to .env and helm prod, dev

## How to run Secret Manager on different Environment 
## Running Locally (Windows, Mac, Linux)
## Authenticate locally using gcloud  on git bash 
    
   1. gcloud auth application-default login
   2. Then run the code 

## for running on gcp simply deploy the code 

##
 use APPENV = LOCAL,"",DEVELOPMENT,PRODUCTION

 To run in local use LOCAL or ""
  