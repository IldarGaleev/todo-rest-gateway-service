# ToDo list rest API gateway

## Environment variables

|       Key       | Values               |   Default   | Description                   |
|:---------------:|----------------------|:-----------:|-------------------------------|
|   `ENV_MODE`    | `local`,`dev`,`prod` |   `prod`    | Production mode               |
| `GRPC_HOSTNAME` | `int`                | `localhost` | gRPC server hostname          |
|   `GRPC_PORT`   | `str`                |   `9090`    | gRPC server tcp port          |
| `API_HOSTNAME`  | `str`                | `localhost` | API server listening hostname |
|   `API_PORT`    | `int`                |   `8080`    | API server listening port     |

## YAML config file 

[config file template](template.config.yml)

```yaml
env-mode: 'local' # 'dev','prod'

grpc-client:
  port: 9090
  hostname: "localhost"

api:
  port: 8080
  hostname: "localhost"
```





