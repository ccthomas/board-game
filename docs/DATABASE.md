# Database

Board uses a [PostgreSQL](https://www.postgresql.org) database, with migrations managed using [golang-migrate](https://github.com/golang-migrate/migrate/tree/master)

## Create Migration

Run following command from base/root of the project
```cmd
docker run --rm -v $(pwd)/backend/db/migrations:/migrations migrate/migrate \
  create -ext sql -dir /migrations -seq <name-of-schema>
```

## Run Migrations

Call the configuration database migration apis to run migrations. You can find documentation on the APIs under the TECHINICAL-API-DESIGN.md document. Look for the Configuration APIs for details on migrations

## TODO

Covnert ID feilds from 36bit VARCHAR to 16bit UUID.