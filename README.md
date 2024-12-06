# Medods task API

<hr>

## Требования для запуска

- Docker compose 

### Для начала
```
cp .env.example .env
```

заполнить `.env` значениями, которые нужны.

В `.env.example` заполнены уже рабочие значение.


### Запуск в dev режиме:

```
docker compose -f compose.dev.yml up
```

Апи запустится по адресу `http://localhost`, port указывается в `.env` файле

### Запуск в prod режиме:

```
docker compose -f compose.prod.yml up
```

Апи запустится по адресу `http://localhost`, port указывается в `.env` файле

### API Routes

```
PATH: /auth/login
Query:
- userId: uuid

Response body:
{
    "refresh": "jwt_token",
    "access": "jwt_token"
}

```

```
PATH: /auth/refresh
Query:
- refreshToken: jwtToken
- accessToken: jwtToken

Response body:
{
    "refresh": "jwt_token",
    "access": "jwt_token"
}
```


## Комментарии по заданию
- Структура проекта следует Hexagonal Architecture
- Суффиксы модулей:
  - `srv` - `service`
  - `repo` - `repository`
  - `hdl` - `handler`
- Структура бд не сложная, одна таблица с полями `hashedRefreshTokenId` (unique, not null) и `token_family` (primary key)
  - В базе данных хранится не сам рефреш токен, а его уникальный айди (захешированный) вместе с его familyToken 
  - По сути нет смысла хешировать айди токена, но решил сделать т.к. написано в тз
- Данные в `accessToken`
  - Алгоритм SHA-512  
  - userId - Айди юзера, не хранить его трудно
  - linkerId - Мое решение для связи refresh и access токена. У пары одинаковый айди тип uuid
  - userIp - Айпи, откуда пришел запрос
  - exp - Стандартное поле для истечения токена
  - iat - Стандартное поле по времени выдача токена
- Данные в `refreshToken`
  - Алгоритм SHA-512
  - jti - UUID, используется для хранение в бд
  - tokenFamily - Мое решение для защиты от повторного использования. При рефреше копируется tokenFamily в новый рефреш токен и в бд проверяется, соответвует ли айди токена с tokenFamily
  - linkerId - Описано в accessToken
  - exp - Стандартное поле для истечения токена
  - iat - Стандартное поле по времени выдача токена