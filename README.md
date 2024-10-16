Телеграм бот - оповещения о днях рождений друзей и знакомых пользователя

Стек: go-telegram-bot-api, sqlite/postgreSQL(в процессе имплементации), redis(в процессе имплементации, на данный момент kv-хранилище - map),

При разработке придерживаюсь правильной архитектуры: структура проекта, использование контекста, DAL, логирования.

Запуск доступен только через билд cmd/mybot/main.go файла:

```go build cmd/mybot/main.go```

В будущем добавлю Dockerfile и docker-compose.yml для запуска в изолированном контейнере

*TODO:* 
1. Notifications (schedule).
2. implement redis, postgresql
3. docker
4. CRUD birthdays
5. add groups of people (friends, family, etc)
6. CRUD groups
7. ...