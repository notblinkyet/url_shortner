# URL Shortener

## Описание

URL Shortener — это сервис для генерации сокращенных ссылок, работающий через gRPC. На данный момент приложение создает случайные строки длиной 10 символов для сокращения URL.

proto файл и сгенерированный код находятся [здесь](https://github.com/notblinkyet/proto_url_shortner).

- Основная база данных: PostgreSQL
- Кэширование: Redis

### Примеры запросов

Создание короткой ссылки:
```
{
    "url": "vk.com"
}
```
Ответ:
```
{
    "short_url": "aq4cxYKy9R"
}
```
Получение оригинального URL:
```
{
    "short_url": "aq4cxYKy9R"
}
```
Ответ:
```
{
    "url": "vk.com"
}
```

## Установка и запуск

### Локальный запуск
```
git clone https://github.com/notblinkyet/url_shortner
export CONFIG_PATH=./configs/local.yaml
export POSTGRES_PASS=<your_postgres_password>
export REDIS_PASS=<your_redis_password>
```
Настройте конфигурацию в ./configs/local.yaml, затем запустите:

С помощью task:
```
task migrate_up
task run_app
```
Без task:
```
go run cmd/migrator/main.go
go run cmd/app/main.go
```
Приложение будет доступно по указанному в конфигурации хосту и порту.

### Запуск в Docker
```
git clone https://github.com/notblinkyet/url_shortner
```
Настройте переменные в Dockerfile.migrator, Dockerfile.app, docker-compose.app.yaml и ./configs/remote.yaml, затем выполните:
```
docker compose -f docker-compose.app.yaml up --build
```
Приложение будет доступно по хосту контейнера и порту, указанному в конфигурации.

## Структура проекта
```
url_shortner/
├── internal
│   ├── repository        # Работа с базой данных и кэшем
│   ├── services          # Бизнес-логика приложения
│   ├── transport         # gRPC сервер
│   ├── config            # Конфигурационные файлы
│   ├── errors            # Определения ошибок
│   └── app               # Инициализация приложения
├── cmd                   # Исполняемые файлы
├── migrations            # SQL-миграции для базы данных
├── configs               # Конфигурационные файлы
├── docker-compose.app.yaml
├── Dockerfile.migrator
├── Dockerfile.app
└── Taskfile.yaml
```

## Улучшения и планы

- Добавить авторизацию
- Улучшить логирование и мониторинг
- Добавить балансировщик

## Автор
[notblinkyet](https://github.com/notblinkyet)