# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

## Environment variable

By default, the server listens on port **7540**.

To use a different port, set the `TODO_PORT` environment variable before starting the server.

```bash
# Export the variable
export TODO_PORT=7540
go run .

# Or set it only for the current command
TODO_PORT=7540 go run .
```
