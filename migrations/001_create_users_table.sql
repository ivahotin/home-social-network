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
