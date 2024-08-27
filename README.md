# go-http-server
## How to run migrations
1. Download golang migrate:
`curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz`
`sudo mv migrate /usr/local/bin/migrate`
`sudo chmod +x /usr/local/bin/migrate`

2. Run Migration:
`migrate -database postgres://app_go:app12345@localhost:5432/production?sslmode=disable -path ./migrations up`