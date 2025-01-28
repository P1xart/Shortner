# Shortner
Проект, сокращающий длинную ссылку в 5-символьную.\
Имеет простейшую метрику переходов (Получение ссылки методом GET = 1 переход)\
Сервис не примет не валидную ссылку.

## Документация
### Создание сокращенной ссылки:
```
POST:/
{
  "link":"yoursourcelink.ru"
}
```
Пример запроса:
```
МЕТОД: POST
domain.com/
ТЕЛО:
{
  "link":"https://google.com/"
}
```
Пример ответа:
```
{
  "src_link":"https://google.com/"
  "short_link":"domain.com/5KdE4"
}
```
### Получение ссылки и метрики:
```
GET:/short_link
```
Пример запроса:
```
МЕТОД: GET
domain.com/5KdE4
```
Пример ответа:
```
{
  "src_link":"https://google.com/"
  "short_link":"domain.com/5KdE4"
  "count_visits":42
}
```


## Запуск:
1) В файле `shortner-service/internal/config/config.go` изменить константу *DOMEN* на своё значение.
2) `$ sudo docker compose up --build`


