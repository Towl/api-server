# API Server

## Dependencies

- [logger](https://github.com/towl/logger)

## Configuration

Expected environment variables:

| Name | Default | Exemples | Description |
| ---- | ------- | -------- | ----------- |
| `API_SERVER_PORT` | | 8083 | Port on which the server will listen |
| `API_SERVER_HOST` | | 0.0.0.0 | Address on which the server will listen |
| `API_SERVER_LOG_LEVEL` | `info` | `trace`, `debug`, `info`, `warning`, `error` | Level of logs |
| `API_SERVER_LOG_FILE` | `/dev/stdout` | `/tmp/server.log` | File in which logs will be written |
| `API_SERVER_LOG_SYSLOG` | `false` | `true`, `false` | If true, most important (some info and `warning`or above) logs will also be written in stdout |
| `API_SERVER_WORKING_DIR` | `.` | `/var/www` | Root folder in which the server will find resources |
| `API_SERVER_PATTERN` | | <code>^/api/(object1&#124;object2)$</code> | Pattern of urls leading to API backend, more info bellow. |

## Backend
