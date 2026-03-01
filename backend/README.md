# Board Game Backend

## CLI Usage

If running from CLI, we need environment configs fron `.env`. The file is however configured for docker containers, and not cli. The `.env.backend_cli` contains overrides for connecting to docker container.

```
go run . "../local/.env" "../local/.env.cli"
```

## Testing

Run the following command to install Uber's mocking tool (Mac)
```cmd
go install go.uber.org/mock/mockgen@latest 

export PATH=$PATH:$(go env GOPATH)/bin 

mockgen --source=internal/X/Y.go --destination=internal/X/mock/mock_Y.go
```             
