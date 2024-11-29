# superlit-backend

Backend For Superlit, Written In Go

## Setup using Docker

1. populate the `.env` file. a template is provided in `.env.example`

note: everything in the `.env` file and the private key are a secret and are not supposed to be shared.

2. Clone the [frontend](https://github.com/anuragrao04/superlit-frontend) as well and place it adjacent to this directory:

```bash
project-root/
├── superlit-backend/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── ... (other backend files)
├── superlit-frontend/
│   ├── Dockerfile
│   ├── package.json
│   ├── package-lock.json
│   ├── src/
│   │   └── ... (frontend source files)
│   ├── nginx/
│   │   └── nginx.conf
│   └── ... (other frontend files)
```

Note: make sure to populate `.env` in the frontend repository as well

3. Run `docker-compose up`

4. You'll find the website running on `http://localhost`

## Setup Without Docker

1. Install golang
2. Run the following to install all dependencies

```shell
go mod tidy
```

3. populate the `.env` file. a template is provided in `.env.example`

note: everything in the `.env` file and the private key are a secret and are not supposed to be shared.

5. You can run using:

```shell
go run main.go
```

## Commit Conventions

- Follow [conventionalcommits.org](https://www.conventionalcommits.org/en/v1.0.0/)
