## Репликация

### Настройка асинхронной репликации.

Подготовим конфигурационные файлы для [мастера](async/master/my.conf), [первого](async/slave1.conf) и [второго](async/slave2.conf) слейвов.

Далее запустим экземпляры СУБД

```
docker-compose -f replication/async/docker-compose.yaml up -d
```

Зайдем в мастер контейнер

```
docker exec -it social-db bash
mysql -uroot -ppassword -hmaster
```

Создаем пользователя для репликации и выдаем ему права

```
create user 'replica'@'%' IDENTIFIED BY 'secret';
grant replication slave on *.* to 'replica'@'%';
```

Далее определим *MASTER_LOG_POS*

```
show master status;
```

Получаем результат

```
+---------------+----------+--------------+------------------+-------------------+
| File          | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+---------------+----------+--------------+------------------+-------------------+
| binlog.000003 |      657 |              |                  |                   |
+---------------+----------+--------------+------------------+-------------------+
```

Далее идем на первый слейв

```
docker exec -it social-db-slave1 bash
mysql -uroot -ppassword -h127.0.0.1
```

Добавляем информацию о мастере

```
CHANGE MASTER TO
    MASTER_HOST='master',
    MASTER_USER='replica',
    MASTER_PASSWORD='secret',
    MASTER_LOG_FILE='mysql-bin.000003',
    MASTER_LOG_POS=657;
```

Запускаем слейв

```
start slave;
```

Проверим статус слейва

```
show slave status\G;
```

Повторяем шаги для второго слейва.

Теперь создадим таблицы и индексы в мастере
```
create table if not exists profiles (
    id                bigint not null auto_increment primary key,
    username          varchar(30) not null constraint username_empty_field check (username != '') unique,
    password          text not null constraint password_empty_field check (password != ''),
    firstname         varchar(30) not null,
    lastname          varchar(50) not null,
    gender            integer not null,
    interests         text not null,
    city              varchar(50) not null,
    birthdate         date not null
) engine = InnoDB;
create index firstname_lastname_idx on profiles (firstname, lastname) using btree;
create table if not exists followers (
    follower_id     bigint not null,
    user_id         bigint not null,
    is_active       boolean not null,
    created_at      timestamp default current_timestamp(),
    primary key (user_id, follower_id)
) engine = InnoDB;
```

Проверим, что на слейвах появились таблицы

```
show tables;
+---------------------+
| Tables_in_social_db |
+---------------------+
| followers           |
| profiles            |
+---------------------+
```

Теперь загрузим данные профайлов. После выгрузки вспоминаем, что лучше бы сначала импортировать данные, а уже потом создавать индексы.

```
load data infile '/out.csv' into profiles fields terminated by ',' ignore 1 rows;
```

Проверяем, что данные доехали на слейвы

```
select count(*) from profiles;
+----------+
| count(*) |
+----------+
|  1000000 |
+----------+
```

### Добавим ProxySQL

Проверим, что ProxySQL прочитал [конфиг](./proxysql.cnf) и добавил сервера и пользователей.
```
select * from mysql_servers;
+--------------+----------+------+-----------+--------+--------+-------------+-----------------+---------------------+---------+----------------+---------+
| hostgroup_id | hostname | port | gtid_port | status | weight | compression | max_connections | max_replication_lag | use_ssl | max_latency_ms | comment |
+--------------+----------+------+-----------+--------+--------+-------------+-----------------+---------------------+---------+----------------+---------+
| 0            | master   | 3306 | 0         | ONLINE | 1      | 0           | 200             | 0                   | 0       | 0              |         |
| 2            | slave1   | 3306 | 0         | ONLINE | 1      | 0           | 200             | 0                   | 0       | 0              |         |
| 2            | slave2   | 3306 | 0         | ONLINE | 1      | 0           | 200             | 0                   | 0       | 0              |         |
+--------------+----------+------+-----------+--------+--------+-------------+-----------------+---------------------+---------+----------------+---------+

select * from mysql_users;
+----------+----------+--------+---------+-------------------+----------------+---------------+------------------------+--------------+---------+----------+-----------------+---------+
| username | password | active | use_ssl | default_hostgroup | default_schema | schema_locked | transaction_persistent | fast_forward | backend | frontend | max_connections | comment |
+----------+----------+--------+---------+-------------------+----------------+---------------+------------------------+--------------+---------+----------+-----------------+---------+
| root     | password | 1      | 0       | 0                 | social_db      | 0             | 1                      | 0            | 1       | 1        | 1000            |         |
| user     | password | 1      | 0       | 0                 | social_db      | 0             | 0                      | 0            | 1       | 1        | 1000            |         |
+----------+----------+--------+---------+-------------------+----------------+---------------+------------------------+--------------+---------+----------+-----------------+---------+
```

Далее добавим правило направляющее весь трафик на мастер

```
insert into mysql_query_rules (active, match_pattern, destination_hostgroup, cache_ttl, username) values (1, '*', 0, NULL, 'user');
LOAD MYSQL QUERY RULES TO RUNTIME;
SAVE MYSQL QUERY RULES TO DISK;
```

Проверим работоспособность конфигурации.

```
SELECT hostgroup_id,hostname,port,status FROM runtime_mysql_servers;
+--------------+----------+------+--------+
| hostgroup_id | hostname | port | status |
+--------------+----------+------+--------+
| 0            | master   | 3306 | ONLINE |
| 2            | slave2   | 3306 | ONLINE |
| 2            | slave1   | 3306 | ONLINE |
+--------------+----------+------+--------+
```

Далее выполним запрос к БД через frontend ProxySQL

```
mysql -uuser -ppassword -h0.0.0.0 -P6033 -D social_db

select * from profiles limit 5;

+----+-----------------+-----------------+-----------+-----------+--------+-------------+----------------+------------+
| id | username        | password        | firstname | lastname  | gender | interests   | city           | birthdate  |
+----+-----------------+-----------------+-----------+-----------+--------+-------------+----------------+------------+
|  1 | Mac.Turner      | BNC2JXCAKVlsc1P | Lawson    | Wilkinson |      1 | programming | Lehi           | 1990-09-09 |
|  2 | Lolita.Johnston | 62OfLiDAXWWV2Oo | Amanda    | Robel     |      1 | programming | Highland       | 1990-09-09 |
|  3 | Lacy25          | CIsBlFAjh4UZmjr | Priscilla | Kerluke   |      1 | programming | Bethesda       | 1990-09-09 |
|  4 | Demarcus.Yost   | 1ocY7TTyHIXmbF_ | Bertha    | Emmerich  |      1 | programming | Danville       | 1990-09-09 |
|  5 | Pearlie46       | imRod7BYE7JvKCO | Aida      | Cremin    |      1 | programming | Citrus Heights | 1990-09-09 |
+----+-----------------+-----------------+-----------+-----------+--------+-------------+----------------+------------+
```

Схема стенда перед первым прогоном нагрузочных тестов.

![png](./assets/schema.png)