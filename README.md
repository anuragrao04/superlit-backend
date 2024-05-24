# superlit-backend

Backend For Superlit, Written In Go

## Setup

1. Install golang
2. Run the following to install all dependencies

```shell
go mod tidy
```

3. Populate the `.env` file. A template is provided in `.env.example`
4. You have to generate a private key for the JWT web tokens. You can do this by executing the following command:

```shell
openssl ecparam -genkey -name prime256v1 -noout -out private_key.pem
```

Note: Everything in the `.env` file and the private key are a secret and are not supposed to be shared.

5. You can run using:

```shell
go run main.go
```

## Commit Conventions

- Follow [conventionalcommits.org](https://www.conventionalcommits.org/en/v1.0.0/)
