# Домашнее задание №4

1. Ускорить Checkout.listCart (т.е. уменьшить время ответа этой ручки)
    1. При использовании worker pool запрашивать не более 5 sku одновременно
    2. Worker pool нужно написать самостоятельно. Обязательное требование - читаемость и покрытие кода комментариями
2. Во всем сервисе при общении с Product Service необходимо использовать рейт лимит на клиентской стороне (10 RPS)
    1. Допускается использование библиотечных рейт лимитеров
3. Во всех слоях сервиса необходимо прокинуть контекст в интерфейсах, если этого не было сделано ранее
4. Аннулирование заказов старше 10 минут в фоне (рекомендуется применять воркер пул для общения с базой)

Задание со звездочкой - написать собственный рейт-лимитер (читаемый код + комментарии обязательны). За это задание дается алмаз.

# Домашнее задание №3

1. Для каждого сервиса(где необходимо что-то сохранять/брать) поднять отдельную БД в __docker-compose__.
2. Сделать миграции в каждом сервисе (достаточно папки миграций и скрипта).
3. Создать необходимые таблицы.
4. Реализовать логику репозитория в кажом сервисе.
5. В качестве query builder-а можно использовать любую либу(согласовать индивидуально с тьютором). Рекомендуеnся https://github.com/Masterminds/squirrel.
6. Драйвер для работы с postgresql: только __pgx__ (pool).
7. В одном из сервсиов сделать транзакционность запросов (как на воркшопе).

Задание с *:
1. Для каждой БД полнять свой балансировщик (pgbouncer или odyssey, можно и то и то). Сервисы ходят не на прямую в БД, а через балансировщик

# Домашнее задание №2

Во второй домашке ваша задача перевести всё взаимодействие между вашими сервисами на протокол gRPC. То есть взаимодействие по http мы полностью выпиливаем и оставляем только gRPC.  Для вашего удобства и удобства тьютора в каждом проекте заведите Makefile (если ещё нет) и там укажите полезные команды: генерация кода из proto файла и скачивание нужных зависимостей.

Теперь кратко:

1. Переводим всё взаимодействие на gRPC.
2. В Makefile реализуем команды generate (если есть, что. генерить), vendor-proto (если используете вендоринг)

P. S. Gateway и proto-валидацию прикручивать НЕ нужно.

Ссылка на код из workshop (ветка master):

[https://gitlab.ozon.dev/go/classroom-5/Week-2/workshop](https://gitlab.ozon.dev/go/classroom-5/Week-2/workshop)

# Домашнее задание №1

- Создать скелеты трёх сервисов по описанию АПИ из файла contracts.md
- Структуру проекта сделать с учетом разбиения на слои, бизнес-логику писать отвязанной от реализаций клиентов и хендлеров
- Все хендлеры отвечают просто заглушками
- Сделать удобный враппер для сервера по тому принципу, по которому делали на воркшопе
- Придумать самостоятельно удобный враппер для клиента
- Все межсервисные вызовы выполняются. Если хендлер по описанию из contracts.md должен ходить в другой сервис, он должен у вас это успешно делать в коде.
- Общение сервисов по http-json-rpc
- должны успешно проходить make precommit и make run-all в корневой папке
- Наладить общение с product-service (в хендлере Checkout.listCart). Токен для общения с product-service получить, написав в личку @pav5000
