# Board Game Backend

## CLI Usage

If running from CLI, we need environment configs fron `.env`. The file is however configured for docker containers, and not cli. The `.env.backend_cli` contains overrides for connecting to docker container.

```
go run . "../local/.env" "../local/.env.cli
```