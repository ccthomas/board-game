# Development

## Getting Started

To run the Board Game Application, you only need [Docker](https://www.docker.com). Our local development docker compose has PGAdmin. The only development tool that must be installed is Hoppscotch. Although any API calling tool, including cli curl commands, are acceptable

If you'd like to run things locally, you will need the following
* [Go](https://go.dev) - Backend Service(s)
* [NodeJs](https://nodejs.org/en) - Front End Application
    * Recommended to use [NVM](https://github.com/nvm-sh/nvm)
* [PGAdmin](https://www.pgadmin.org) - administration and development platform for PostgreSQL
* [Postgres](https://www.postgresql.org) - SQL Database
* [Hoppscotch](https://hoppscotch.com/download) - makes it easy to create and test your APIs, helping you to ship products faster.
    * NOTE: Postman logs all API transactions in their cloud. For this reason, we have gone with Hoppscotch as an alternitive.

## Usage

### Database

To run the database, execute the following docker commands

```cmd
docker-compose -f ./local/docker-compose.yml  up -d postgres pgadmin
```

In the browser, go to `http://localhost:5050` to use PG Admin. The Email & Password can be found in the `./local/.env` file.

Register a new server and fill in the following information
* General
    * Name: Whatever you'd like, recommended "Board Game"
* Connection (use varuables in `./local/.env`)
    * Host name / address: DB_CONTAINER_NAME
    * Port: DB_PORT
    * Maintance Database: DB_NAME
    * Username: DB_USER
    * Password: DB_PASSWORD

Look to `docs/DATABASE.md` for usage around creating & running migrations.

### API

Download [Hoppscotch] and open the `Board Game Hoppscotch Collection.json` file. Here you can call all the exposed apis.

### Helpful Commands

* Nuke Docker containers, images, & volumes
  
  ```cmd
  docker system prune -a --volumes -f
  ```