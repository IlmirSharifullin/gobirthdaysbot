Телеграм бот - оповещения о днях рождений друзей и знакомых пользователя

Стек: go-telegram-bot-api, sqlite/postgreSQL(в процессе имплементации), redis(в процессе имплементации, на данный момент kv-хранилище - map),

При разработке придерживаюсь правильной архитектуры: структура проекта, использование контекста, DAL, логирования.

Запуск через docker compose:

```docker compose up --build -d```

**TODO:**

1. implement redis, postgresql
2. CRUD birthdays
3. add groups of people (friends, family, etc.)
4. CRUD groups
5. ...