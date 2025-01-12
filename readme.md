# ILCS

## Installation

To run this project need:

1. [Taskfile](https://taskfile.dev/installation/)
2. [Go](https://go.dev/doc/install) 1.23
3. [Postman](https://www.postman.com/downloads/)
4. [Docker](https://docs.docker.com/engine/install/)

run this command on terminal it will download dependencies

```bash
task setup
```

```bash
go mod tidy
```

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`GOOSE_DRIVER`

`GOOSE_DBSTRING`

`GOOSE_MIGRATION_DIR`

`PORT`

`JWT_SECRET`

`REDIS_ADDR`

## Run Locally

Run with docker

```bash
docker-compose up -d
```

Run without docker

```bash
task run
```

## Documentation

import using postman this json

```json
{
  "info": {
    "_postman_id": "4f29f18e-fb52-4fb8-8a83-73b16acdddda",
    "name": "ILCS",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
    "_exporter_id": "14623263"
  },
  "item": [
    {
      "name": "Create Todo",
      "request": {
        "auth": {
          "type": "bearer",
          "bearer": [
            {
              "key": "token",
              "value": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzY3Nzg3NzR9.3UF9wtg--CjqFDr1awHM-qJpZRnH-l4JEzM5NSBCl-w",
              "type": "string"
            }
          ]
        },
        "method": "POST",
        "header": [],
        "body": {
          "mode": "raw",
          "raw": "{\r\n    \"title\": \"stream\",\r\n    \"description\": \"melakukan streaming\",\r\n    \"due_date\": \"2025-01-12\"\r\n}",
          "options": {
            "raw": {
              "language": "json"
            }
          }
        },
        "url": {
          "raw": "http://localhost:7575/api/v1/tasks",
          "protocol": "http",
          "host": ["localhost"],
          "port": "7575",
          "path": ["api", "v1", "tasks"]
        }
      },
      "response": []
    },
    {
      "name": "Get List Todo",
      "request": {
        "auth": {
          "type": "bearer",
          "bearer": [
            {
              "key": "token",
              "value": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzY3Nzg3NzR9.3UF9wtg--CjqFDr1awHM-qJpZRnH-l4JEzM5NSBCl-w",
              "type": "string"
            }
          ]
        },
        "method": "GET",
        "header": [],
        "url": {
          "raw": "http://localhost:7575/api/v1/tasks?page=1&limit=10",
          "protocol": "http",
          "host": ["localhost"],
          "port": "7575",
          "path": ["api", "v1", "tasks"],
          "query": [
            {
              "key": "page",
              "value": "1"
            },
            {
              "key": "limit",
              "value": "10"
            },
            {
              "key": "search",
              "value": "makan",
              "disabled": true
            },
            {
              "key": "status",
              "value": "completed",
              "disabled": true
            }
          ]
        }
      },
      "response": []
    },
    {
      "name": "Get Todo By Id",
      "request": {
        "auth": {
          "type": "bearer",
          "bearer": [
            {
              "key": "token",
              "value": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzY3Nzg3NzR9.3UF9wtg--CjqFDr1awHM-qJpZRnH-l4JEzM5NSBCl-w",
              "type": "string"
            }
          ]
        },
        "method": "GET",
        "header": [],
        "url": {
          "raw": "http://localhost:7575/api/v1/tasks/:id",
          "protocol": "http",
          "host": ["localhost"],
          "port": "7575",
          "path": ["api", "v1", "tasks", ":id"],
          "variable": [
            {
              "key": "id",
              "value": "01945895-1d9d-7de4-8623-80f32435fe73"
            }
          ]
        }
      },
      "response": []
    },
    {
      "name": "Update Todo",
      "request": {
        "auth": {
          "type": "bearer",
          "bearer": [
            {
              "key": "token",
              "value": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzY3Nzg3NzR9.3UF9wtg--CjqFDr1awHM-qJpZRnH-l4JEzM5NSBCl-w",
              "type": "string"
            }
          ]
        },
        "method": "PUT",
        "header": [],
        "body": {
          "mode": "raw",
          "raw": "{\r\n    \"title\": \"game\",\r\n    \"description\": \"memainkan game\",\r\n    \"status\": \"completed\",\r\n    \"due_date\": \"2025-02-12\"\r\n}",
          "options": {
            "raw": {
              "language": "json"
            }
          }
        },
        "url": {
          "raw": "http://localhost:7575/api/v1/tasks/:id",
          "protocol": "http",
          "host": ["localhost"],
          "port": "7575",
          "path": ["api", "v1", "tasks", ":id"],
          "variable": [
            {
              "key": "id",
              "value": "01945895-1d9d-7de4-8623-80f32435fe73"
            }
          ]
        }
      },
      "response": []
    },
    {
      "name": "Delete Todo",
      "request": {
        "auth": {
          "type": "bearer",
          "bearer": [
            {
              "key": "token",
              "value": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzY3Nzg3NzR9.3UF9wtg--CjqFDr1awHM-qJpZRnH-l4JEzM5NSBCl-w",
              "type": "string"
            }
          ]
        },
        "method": "DELETE",
        "header": [],
        "url": {
          "raw": "http://localhost:7575/api/v1/tasks/:id",
          "protocol": "http",
          "host": ["localhost"],
          "port": "7575",
          "path": ["api", "v1", "tasks", ":id"],
          "variable": [
            {
              "key": "id",
              "value": "019455ed-27b2-770b-b74c-67d1ff797cb0"
            }
          ]
        }
      },
      "response": []
    },
    {
      "name": "Get Token",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "http://localhost:7575/api/v1/token",
          "protocol": "http",
          "host": ["localhost"],
          "port": "7575",
          "path": ["api", "v1", "token"]
        }
      },
      "response": []
    }
  ]
}
```

## License

[MIT](https://choosealicense.com/licenses/mit/)

## Authors

- [@sutantodadang](https://www.github.com/sutantodadang)
