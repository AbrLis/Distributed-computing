# Distributed-computing
Учебный репозиторий для задания по созданию сервиса распредённых вычислений

## Оглавление
- [Описание задачи](#описание-задачи)
- [Запуск проекта](#запуск-проекта)
- [API](#API)
- [Описание структуры проекта в целом](#описание-структуры-проекта-в-целом)
- [Что реализовано, а что осталось за кадром](#что-реализовано-а-что-осталось-за-кадром)
- [Дополнение, если у вас Windows и нет команды curl](#дополнение-если-у-вас-windows-и-нет-команды-curl)


## Описание задачи

### Распределенный вычислитель арифметических выражений

#### Общее описание
Пользователь хочет считать арифметические выражения. Он вводит строку 2 + 2 * 2 и хочет получить в ответ 6. Но наши операции сложения и умножения (также деления и вычитания) выполняются "очень-очень" долго. Поэтому вариант, при котором пользователь делает http-запрос и получает в качетсве ответа результат, невозможна. Более того: вычисление каждой такой операции в нашей "альтернативной реальности" занимает "гигантские" вычислительные мощности. Соответственно, каждое действие мы должны уметь выполнять отдельно и масштабировать эту систему можем добавлением вычислительных мощностей в нашу систему в виде новых "машин". Поэтому пользователь, присылая выражение, получает в ответ идентификатор выражения и может с какой-то периодичностью уточнять у сервера "не посчиталость ли выражение"? Если выражение наконец будет вычислено - то он получит результат. Помните, что некоторые части арфиметического выражения можно вычислять параллельно.

#### Front-end часть

GUI, который можно представить как 4 страницы

Форма ввода арифметического выражения. Пользователь вводит арифметическое выражение и отправляет POST http-запрос с этим выражением на back-end. Примечание: Запросы должны быть идемпотентными. К запросам добавляется уникальный идентификатор. Если пользователь отправляет запрос с идентификатором, который уже отправлялся и был принят к обработке - ответ 200. Возможные варианты ответа:
 - `200.` Выражение успешно принято, распаршено и принято к обработке
 - `400.` Выражение невалидно
 - `500.` Что-то не так на back-end. В качестве ответа нужно возвращать id принятного к выполнению выражения.


 - Страница со списком выражений в виде списка с выражениями. Каждая запись на странице содержит статус, выражение, дату его создания и дату заверщения вычисления. Страница получает данные GET http-запрсом с back-end-а
 - Страница со списком операций в виде пар: имя операции + время его выполнения (доступное для редактирования поле). Как уже оговаривалось в условии задачи, наши операции выполняются "как будто бы очень долго". Страница получает данные GET http-запрсом с back-end-а. Пользователь может настроить время выполения операции и сохранить изменения.
 - Страница со списком вычислительных можностей. Страница получает данные GET http-запросом с сервера в виде пар: имя вычислительного ресурса + выполняемая на нём операция.

Требования:
- Оркестратор может перезапускаться без потери состояния. Все выражения храним в СУБД.
- Оркестратор должен отслеживать задачи, которые выполняются слишком долго (вычислитель тоже может уйти со связи) и делать их повторно доступными для вычислений.

#### Back-end часть
Состоит из 2 элементов:

- Сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения. Далее будем называть его оркестратором.
- Вычислитель, который может получить от оркестратора задачу, выполнить его и вернуть серверу результат. Далее будем называть его агентом.

**Оркестратор**

Сервер, который имеет следующие endpoint-ы:

- Добавление вычисления арифметического выражения.
- Получение списка выражений со статусами.
- Получение значения выражения по его идентификатору.
- Получение списка доступных операций со временем их выполения.
- Получение задачи для выполения.
- Приём результата обработки данных.

[Назад к оглавлению](#оглавление)
 
## Запуск проекта

- Склонировать репозиторий в папку проекта
- Запустить проект с помощью команды 
```shell
go run main.go
```
Проект запустится по адресу `http://localhost:3000`

[Назад к оглавлению](#оглавление)

## API

Реализованые эндпоинты и примеры их вызова:

----

- API для добавления выражений 
`POST запрос`
```shell
curl -X POST -H "Content-Type: text/plain" -d "2+2*2" http://localhost:3000/add-expression
```
Принимается текстовое выражение в формате математического выражения.

Поддерживаются операции сложения, вычитания, умножения, деления. Учитывается приоритет операций.

Скобочные выражения не поддерживаются.

Так же не поддерживаются выражения вида 2+-2 и отрицательные числа в том числе вида -2*-2. Такое выражение будет принято к обработке, но в результате обработке будет значиться ошибка парсинга.

Ответ будет представлен в виде текста вида:
```shell
Выражение добавлено в базу данных и принято к обработке. ID: 1707470348191170700
```
В дальнейшем будем использовать этот ID для получения результата вычислений.

----

- API для получения списка выражений
`GET запрос`
```shell
curl -X GET http://localhost:3000/get-expressions
```
Получает список всех выражений со статусами в формате JSON.

Результат будет представлен ответом вида:
```json
[
  {
    "id": "1707470161983588300",
    "expression": "2+2*4/2+64-33*4",
    "status": "Completed",
    "result": "-62"
  },
  {
    "id": "1707470348191170700",
    "expression": "2+-2",
    "status": "Error",
    "result": "error parse or calculate"
  }
]
```
Либо пустым списком, если на данный момент выражения отстутствуют в базе данных.

----

- API для получения значения выражения по ID

`GET запрос`
```shell
curl -X GET http://localhost:3000/get-value/1707470348191170700
```
где 1707470348191170700 - ID выражения.

Ответ будет представлен в виде JSON с данными выражения.

```json
{
  "id": "1707470348191170700",
  "expression": "2+2*4/2+64-33*4",
  "status": "Completed",
  "result": "-62"
}
```
либо текстом c кодом 404, если задача не найдена.
```shell
Задача не найдена
```
----

- API для получения полного списка операций +, -, *, / со времение их выполнения в секундах.

`GET запрос`
```shell
curl -X GET http://localhost:3000/get-operations
```
Ответ будет представлен в виде JSON.
```json
{
  "add": "5s",
  "sub": "3s",
  "mult": "4s",
  "div": "6s"
}
```

----

- API для получения таймаутов ответов вычислителей (мониторинг отклика)

`GET запрос`
```shell
curl -X GET http://localhost:3000/monitoring
```
Ответ будет представлен в виде JSON.
```json
{
  "1": "1.177 sec",
  "2": "1.209 sec",
  "3": "1.209 sec",
  "4": "1.225 sec",
  "5": "1.209 sec"
}
```

где ключом является номер вычислителя, а значением - время с последнего отклика в секундах.

[Назад к оглавлению](#оглавление)


----

## Описание структуры проекта в целом

- Проект реализован в виде нескольких сервисов взаимодействующих между собой посредством API.

- База данных в этом проекте представлена как map структура с методами взаимодействия и не сохраняется между перезапусками.

- API представляет собой интерфейс для взаимодействия между сервисами и клиентом предоставляя тому возможность получить результаты и задать выражение на вычисление.

- Демон вычислителя - это отдельный сервис, который служит прослойком между API и запускаемыми им же вычислителями. Он следит за работой очереди вычислителей, запрашивает задачи через API. Получив задачу, демон, выставляет её в очередь исполнения, и ожидает свободных вычилитетей. Если вычислитель свободен, то из очереди ожидания, задача переводится в очередь исполнения и передаётся на обработку.

- Вычислители это отдельные go рутины запускаемые демоном и ожидающие задачи в канале исполения. Каждая задача выполняется отдельным вычислителем. После завершения вычислений, этот процесс передаёт результаты своей работы через API в базу данных, сообщает демону что он свободен и переходит в режим ожидания.

[Назад к оглавлению](#оглавление)

----

## Что реализовано, а что осталось за кадром
### Реализовано:

Сервер (оркестратор)

Агенты

- API для добавления выражений
- API для получения списка выражений
- API для получения значения выражения
- API для получения списка операций

### Не реализовано:

Frontend:

- Форма ввода выражения
- Список выражений
- Список операций
- Список вычислительных ресурсов

Backend:

- Не реализован функционал отслеживания задач, которые выполняются слишком долго.
- Перезапуск оркестратора без потери состояния:
- Не реализована возможность перезапуска оркестра

[Назад к оглавлению](#оглавление)


----

## Дополнение, если у вас Windows и нет команды curl

Установить эту утилиту можно с официального [сайта](https://curl.se/)

И так же ссылка на сам файл для Windows [downloads](https://curl.se/download.html#Win64)

[Назад к оглавлению](#оглавление)