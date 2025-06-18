# install golang migrate
[HOW TO INSTALL](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

# SQLC 
[how to install](https://docs.sqlc.dev/en/stable/overview/install.html)

or
```sh
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

# generate pem
openssl genrsa -out key/private.key 2048

openssl rsa -in key/private.key -pubout -out key/public.pem
