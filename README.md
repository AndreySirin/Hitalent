## Тестовое задание от Hitalent

## Установка
```shell

```
## Запуск
```shell
go mod tidy
docker compose up
docker compose down
```
## Методы API
### Создание вопроса
```shell
POST
http://localhost:8080/questions/

тело запроса
{
  "text": "как настроение?"
}
тело ответа
{
    "id": 1,
    "message": "Вопрос успешно создан"
}
```
### Удаление вопроса
```shell
DELETE
http://localhost:8080/questions/1

тело запроса
{}
тело ответа
{
    "id": 1,
    "message": "успешное удаление"
}
```
### Получение всех вопросов
```shell
GET
http://localhost:8080/questions

тело запроса
{}
[
    {
        "id": 1,
        "text": "как настроение?",
        "created_at": "2025-11-20T15:53:34.215681Z"
    }
]
```
### Создание ответа к вопросу
```shell
POST
http://localhost:8080/questions/1/answers/

тело запроса
{
  "text": "что ты сказал?!"
}
тело ответа
{
    "id": 1,
    "message": "ответ успешно добавлен"
}
```
### Получение вопроса и всех ответов на него
```shell
GET
http://localhost:8080/questions/1

тело запроса
{}
тело ответа
{
    "answers": [
        {
            "id": 1,
            "question_id": 1,
            "user_id": "ec3edc32-e34f-4671-aa3c-a6c9077b718e",
            "text": "что ты сказал?!",
            "created_at": "2025-11-20T15:53:45.143328Z"
        }
    ],
    "question": {
        "id": 1,
        "text": "как настроение?",
        "created_at": "2025-11-20T15:53:34.215681Z"
    }
}
```
### Получение конкретного ответа
```shell
GET
http://localhost:8080/answers/1

тело запроса
{}
тело ответа
{
    "id": 1,
    "question_id": 1,
    "user_id": "ec3edc32-e34f-4671-aa3c-a6c9077b718e",
    "text": "что ты сказал?!",
    "created_at": "2025-11-20T15:53:45.143328Z"
}
```
### Удаление конкретного ответа
```shell
DELETE
http://localhost:8080/answers/1

тело запроса
{}
тело ответа
{
    "id": 1,
    "message": "ответ успешно удален"
}
```
