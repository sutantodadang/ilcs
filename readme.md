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
        "method": "POST",
        "header": [],
        "body": {
          "mode": "raw",
          "raw": "{\r\n    \"title\": \"game\",\r\n    \"description\": \"memainkan game\",\r\n    \"due_date\": \"2025-01-12\"\r\n}",
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
        "method": "GET",
        "header": [],
        "url": {
          "raw": "http://localhost:7575/api/v1/tasks?page=1&limit=10&search=makan&status=completed",
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
              "value": "makan"
            },
            {
              "key": "status",
              "value": "completed"
            }
          ]
        }
      },
      "response": []
    },
    {
      "name": "Get Todo By Id",
      "request": {
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
              "value": ""
            }
          ]
        }
      },
      "response": []
    },
    {
      "name": "Update Todo",
      "request": {
        "method": "PUT",
        "header": [],
        "body": {
          "mode": "raw",
          "raw": "{\r\n    \"title\": \"todo3\",\r\n    \"description\": \"desc\",\r\n    \"status\": \"completed\",\r\n    \"due_date\": \"2025-02-12\"\r\n}",
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
              "value": ""
            }
          ]
        }
      },
      "response": []
    },
    {
      "name": "Delete Todo",
      "request": {
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
              "value": ""
            }
          ]
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
