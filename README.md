# Home-social-network

## Производительность индексов

### Предусловия

Используется один golang instance со следующими настройками пула соединений с БД.
```
DB_CONN_LIFE_TIME=3
DB_MAX_OPEN_CONNECTION=10
DB_MAX_IDLE_CONNECTION=10
```

Приложению требуется поиск анкет по имени и фамилии. Для выдачи анкет используется cursor-based пагинация по id анкеты

```mysql
select id, username, password, firstname, lastname, birthdate, gender, interests, city 
from profiles 
where (firstname like 'Ad%' and lastname like 'Kem%') and id > 401989
limit 10
```
Сортировку анкет в выдаче решено производить на стороне приложения. О причинах такого решения будет сказано ниже.
Будем тестировать получение первой страницы выдачи (`id > 0`).

### Поиск без использования индекса (Вариант №1)

План выполнения запроса
```
+----+-------------+----------+------------+-------+---------------+---------+---------+------+--------+----------+-------------+
| id | select_type | table    | partitions | type  | possible_keys | key     | key_len | ref  | rows   | filtered | Extra       |
+----+-------------+----------+------------+-------+---------------+---------+---------+------+--------+----------+-------------+
|  1 | SIMPLE      | profiles | NULL       | range | PRIMARY       | PRIMARY | 8       | NULL | 497205 |     1.23 | Using where |
+----+-------------+----------+------------+-------+---------------+---------+---------+------+--------+----------+-------------+
```

В данной ситуации используется первичный ключ (из-за присутствия условия `id > 0`).
* Спускаемся по BTree первичного ключа к первому листу, соответствующему условию `id > 0`.
* Двигаемся по листьям, попутно проверяя условие поиска по имени и фамилии.
* Как только нашли 100 записей, возвращаем ответ.

Как видно пришлось просканировать половину таблицы ~ 500k строк.

Проведем тестирование
```
wrk -t1 -c1 -d20s --timeout 90s -s ./indexes/generator.lua --latency http://localhost:8080
wrk -t10 -c10 -d20s --timeout 90s -s ./indexes/generator.lua --latency http://localhost:8080
wrk -t100 -c100 -d20s --timeout 90s -s ./indexes/generator.lua --latency http://localhost:8080
wrk -t100 -c1000 -d20s --timeout 90s -s ./indexes/generator.lua --latency http://localhost:8080
```

Результаты нагрузки в табличной форме

| Load  | RPS     | Latency (avg/stdev)     | Max latency
| :---  |    :----:   |   :----:         | :----: |
| 1     | 1.80       | 554.60ms (+/-) 5.75ms | 570.38ms |
| 10     | 8.66        | 1.12s (+/-) 346.47ms  | 2.26s |
| 100     | 8.51         | 7.69s (+/-) 4.32s | 19.73s | 
| 1000     | 8.71          | 10.70s (+/-) 5.34s | 20.15s |

### С композитным индексом (Вариант №2)

Протестируем следующий композитный индекс на полях `firstname` и `lastname`, так как эти поля используются при поиске и именно в этом порядке.
Помимо этого в моих данных селективность имени выше селективности фамилии.

```mysql
create index firstname_lastname_idx on profiles (firstname, lastname) using btree;
```

План выполнения запроса
```
+----+-------------+----------+------------+-------+--------------------------------+------------------------+---------+------+-------+----------+-----------------------+
| id | select_type | table    | partitions | type  | possible_keys                  | key                    | key_len | ref  | rows  | filtered | Extra                 |
+----+-------------+----------+------------+-------+--------------------------------+------------------------+---------+------+-------+----------+-----------------------+
|  1 | SIMPLE      | profiles | NULL       | range | PRIMARY,firstname_lastname_idx | firstname_lastname_idx | 332     | NULL | 17894 |     5.55 | Using index condition |
+----+-------------+----------+------------+-------+--------------------------------+------------------------+---------+------+-------+----------+-----------------------+
```

Видно, что в данной ситуации используется вновь созданный индекс и количество просмотренных строк упало до ~18k.

Результаты нагрузки в табличной форме

| Load  | RPS     | Latency (avg/stdev)     | Max latency
| :---  |    :----:   |   :----:         | :----: |
| 1     | 1535.00       | 703.46us (+/-) 558.20us | 13.74ms |
| 10     | 7294.13        | 1.41ms (+/-) 806.89us  | 29.70ms |
| 100     | 8346.82         | 13.50ms (+/-) 10.99ms | 133.66ms | 
| 1000     | 8530.81          | 221.60ms (+/-) 383.64ms | 4.78s |

Детальный план
```json
{
  "query_block": {
    "select_id": 1,
    "cost_info": {
      "query_cost": "10522.17"
    },
    "table": {
      "table_name": "profiles",
      "access_type": "range",
      "possible_keys": [
        "PRIMARY",
        "firstname_lastname_idx"
      ],
      "key": "firstname_lastname_idx",
      "used_key_parts": [
        "firstname"
      ],
      "key_length": "332",
      "rows_examined_per_scan": 17894,
      "rows_produced_per_join": 994,
      "filtered": "5.55",
      "index_condition": "((`social_db`.`profiles`.`firstname` like 'Ad%') and (`social_db`.`profiles`.`lastname` like 'Kem%') and (`social_db`.`profiles`.`id` > 0))",
      "cost_info": {
        "read_cost": "10422.77",
        "eval_cost": "99.40",
        "prefix_cost": "10522.17",
        "data_read_per_join": "660K"
      },
      "used_columns": [
        "id",
        "username",
        "password",
        "firstname",
        "lastname",
        "gender",
        "interests",
        "city",
        "birthdate"
      ]
    }
  }
}
```

Из детального отчета видно
* Используется только первая часть индекса.
* `read_cost` упала в 10 раз, так как сканируется меньше строк.

### Индекс на поле `firstname` (Вариант №3)

Попробуем использовать индекс с одним полем `firstname`. В данном случае выбрано поле `firstname` так как его селективность в моем наборе выше.

```mysql
create index firstname_idx on profiles (firstname) using btree;
```

План выполнения запроса
```
+----+-------------+----------+------------+-------+-----------------------+---------------+---------+------+-------+----------+------------------------------------+
| id | select_type | table    | partitions | type  | possible_keys         | key           | key_len | ref  | rows  | filtered | Extra                              |
+----+-------------+----------+------------+-------+-----------------------+---------------+---------+------+-------+----------+------------------------------------+
|  1 | SIMPLE      | profiles | NULL       | range | PRIMARY,firstname_idx | firstname_idx | 130     | NULL | 17692 |     5.55 | Using index condition; Using where |
+----+-------------+----------+------------+-------+-----------------------+---------------+---------+------+-------+----------+------------------------------------+
```

Детальный план выполнения запроса. Из плана видно, что стоимость выполнения операции немного выше варианта №2.
```json
{
  "query_block": {
    "select_id": 1,
    "cost_info": {
      "query_cost": "10630.48"
    },
    "table": {
      "table_name": "profiles",
      "access_type": "range",
      "possible_keys": [
        "PRIMARY",
        "firstname_idx"
      ],
      "key": "firstname_idx",
      "used_key_parts": [
        "firstname"
      ],
      "key_length": "130",
      "rows_examined_per_scan": 17692,
      "rows_produced_per_join": 982,
      "filtered": "5.55",
      "index_condition": "((`social_db`.`profiles`.`firstname` like 'Ad%') and (`social_db`.`profiles`.`id` > 0))",
      "cost_info": {
        "read_cost": "10532.20",
        "eval_cost": "98.28",
        "prefix_cost": "10630.48",
        "data_read_per_join": "652K"
      },
      "used_columns": [
        "id",
        "username",
        "password",
        "firstname",
        "lastname",
        "gender",
        "interests",
        "city",
        "birthdate"
      ],
      "attached_condition": "(`social_db`.`profiles`.`lastname` like 'Kem%')"
    }
  }
}
```

Результаты нагрузки в табличной форме

| Load  | RPS     | Latency (avg/stdev)     | Max latency
| :---  |    :----:   |   :----:         | :----: |
| 1     | 1056.36       | 3.29ms (+/-) 9.16ms | 139.62ms |
| 10     | 4242.39        | 8.33ms (+/-) 23.64ms  | 326.80ms |
| 100     | 6041.21         | 18.76ms (+/-) 15.63ms | 221.23ms | 
| 1000     | 5735.85          | 291.44ms (+/-) 466.19ms | 5.17s |

Видно, что по результатам тестирования вариант №3 оказался хуже варианта №2. Наиболее вероятным объяснением этого является использование ICP оптимизации в варианте №2.

* Мы спускаемся по Btree дереву вторичного индекса используя условие на поле `firstname`
* Далее идем по листьям индекса отфильтровывая записи по `lastname` (`lastname` содержится в композитном индексе).
* Далее читаем полученные строки из таблицы (фильтрации на этом этапе уже отсутствует).

В случае же с вариантом №3
* Мы спускаемся по Btree дереву вторичного индекса используя условие на поле `firstname`
* Далее читаем полученные строки из таблицы и отфильтровываем по `lastname` (в данном случае возможно больше random io)

### Графики

![png](./indexes/rps.png)

Для второго графика проведено логарифмическое масштабирование по оси Y

![png](./indexes/latency.png)

### Вывод

Лучший результат по метрикам `RPS` и `latency` показал вариант №2
```mysql
create index firstname_lastname_idx on profiles (firstname, lastname) using btree;
```