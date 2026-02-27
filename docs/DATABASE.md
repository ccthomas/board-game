# Database

Board uses a [PostgreSQL](https://www.postgresql.org) database, with migrations managed using [golang-migrate](https://github.com/golang-migrate/migrate/tree/master)

## Create Migration

Run following command from base/root of the project
```cmd
docker run --rm -v $(pwd)/backend/db/migrations:/migrations migrate/migrate \
  create -ext sql -dir /migrations -seq <name-of-schema>
```

## Run Migrations

You can currently run migrations by running the backend go application. There are other ways using go/migrate locally or docker commands to run migrations as well. Not documented here.