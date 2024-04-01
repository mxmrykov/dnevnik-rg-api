grant usage on schema users, passwords, classes, notifications, archive to c128f7;
grant execute on all functions in schema users, passwords, classes, notifications, archive to c128f7;
grant insert, select, update, delete on all tables in schema users, passwords, classes, notifications, archive to c128f7;


create schema if not exists users;
create schema if not exists passwords;
create schema if not exists classes;
create schema if not exists notifications;
create schema if not exists archive;

create table if not exists users.admins
(
    UDID     bigserial PRIMARY KEY,
    key      integer,
    FIO      text,
    date_reg text,
    logo_uri text,
    role     varchar(10)
);

create table if not exists users.coaches
(
    UDID          bigserial PRIMARY KEY,
    key           integer,
    FIO           text,
    date_reg      text,
    home_city     varchar(30),
    training_city varchar(30),
    birthday      text,
    about         text,
    logo_uri      text,
    role          varchar(10)
);

create table if not exists users.pupils
(
    UDID           bigserial PRIMARY KEY,
    key            integer,
    FIO            text,
    date_reg       text,
    coach          integer,
    home_city      varchar(30),
    traqining_city varchar(30),
    birthday       text,
    about          text,
    coach_review   text,
    logo_uri       text,
    role           varchar(10)
);

create table if not exists users.deleted_admins as
select *
from users.admins
where false;

create table if not exists users.deleted_coaches as
select *
from users.coaches
where false;

create table if not exists users.deleted_pupils as
select *
from users.pupils
where false;

create table if not exists passwords.passwords
(
    UDID        bigserial PRIMARY KEY,
    key         integer,
    checksum    varchar(64),
    token       text,
    last_update text
);

create table if not exists classes.classes
(
    UDID           bigserial PRIMARY KEY,
    key            bigint,
    pupil          integer[],
    coach          integer,
    class_date     text,
    class_time     varchar(5),
    class_dur      varchar(5),
    presence       boolean,
    price          smallint,
    mark           smallint,
    review         text,
    scheduled      boolean,
    classType      varchar(10),
    pupilCount     int,
    isOpenToSignup boolean
);

create table if not exists archive.archived_admins as
select *
from users.admins
where false;

create table if not exists archive.archived_coaches as
select *
from users.coaches
where false;

create table if not exists archive.archived_pupils as
select *
from users.pupils
where false;


-- init two main admins
INSERT INTO users.admins (udid, key, fio, date_reg, logo_uri, role)
VALUES (1, 1704461475, 'Рыков Максим Алексеевич', '2024-01-05T16:31:15+03:00', 'https://dnevnik-rg.ru/admin-logo.png',
        'ADMIN');
INSERT INTO users.admins (udid, key, fio, date_reg, logo_uri, role)
VALUES (4, 1704461712, 'Дубова Дарья Вадимовна', '2024-01-05T16:35:12+03:00', 'https://dnevnik-rg.ru/admin-logo.png',
        'ADMIN');
INSERT INTO passwords.passwords (udid, key, checksum, token, last_update)
VALUES (1, 1704461475, 'c4ef9c3',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOjE3MDQ0NjE0NzUsImNoZWNrX3N1bSI6ImM0ZWY5YzMiLCJyb2xlIjoiMjAyNC0wMS0wNVQxNjozMToxNSswMzowMCIsImV4cCI6MTcwOTY0NTQ3NSwiaXNzIjoibG9jYWxob3N0OjgwMDAifQ.8nXxKUXXkptMuo5x5Cvjxwq2lEF8EC2nCvlcuqepJLY',
        '2024-01-05T16:31:15+03:00');
INSERT INTO passwords.passwords (udid, key, checksum, token, last_update)
VALUES (4, 1704461712, 'b9f94c7',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOjE3MDQ0NjE3MTIsImNoZWNrX3N1bSI6ImI5Zjk0YzciLCJyb2xlIjoiMjAyNC0wMS0wNVQxNjozNToxMiswMzowMCIsImV4cCI6MTcwOTY0NTcxMiwiaXNzIjoibG9jYWxob3N0OjgwMDAifQ.kY4yLF9sODk68y3u_xIncqC5OZbZQcZ6x8WzvHF62bg',
        '2024-01-05T16:35:12+03:00');



create or replace function users.create_admin(
    key_ integer,
    fio_ text,
    date_reg_ text,
    logo_uri_ text,
    role_ text
)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into users.admins
        (key, FIO, date_reg, logo_uri, role)
    VALUES (key_, fio_, date_reg_, logo_uri_, role_)
    on conflict do nothing;
end;
$$;

create or replace function users.create_coach(
    key_ integer,
    fio_ text,
    date_reg_ text,
    home_city_ varchar(30),
    training_city_ varchar(30),
    birthday_ text,
    about_ text,
    logo_uri_ text,
    role_ varchar(10)
)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into users.coaches
    (key, FIO, date_reg, home_city, training_city, birthday, about, logo_uri, role)
    VALUES (key_, fio_, date_reg_, home_city_, training_city_, birthday_, about_, logo_uri_, role_);
end;
$$;

create or replace function get_coach(key_ integer)
    returns table
            (
                key_           integer,
                fio_           text,
                date_reg_      text,
                home_city_     varchar(30),
                training_city_ varchar(30),
                birthday_      text,
                about_         text,
                logo_uri_      text,
                role_          varchar(10)
            )
    security definer
    language plpgsql
as
$$
begin
    select * from users.coaches where key = key_ and
end;
$$;

create or replace function users.delete_admin(
    key_ integer
)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into users.deleted_admins (key, FIO, date_reg, logo_uri, role) (select * from users.admins where key = key_);
    delete from users.admins where key = key_;
end;
$$;

create or replace function users.get_admin(
    key_ integer
)
    returns table
            (
                key         integer,
                fio         text,
                date_reg    text,
                logo_uri    text,
                role        varchar(10),
                checksum    varchar(8),
                last_update text,
                token       text
            )
    security definer
    language plpgsql
as
$$
begin
    select a.key,
           a.fio,
           a.date_reg,
           a.logo_uri,
           a.role,
           p.checksum,
           p.last_update,
           p.token
    from users.admins a
             left join passwords.passwords p on a.key = p.key
    where a.key = key_;
end;
$$;

create or replace function users.list_admins_except(key_ integer)
    returns table
            (
                key      integer,
                fio      text,
                logo_uri text
            )
    security definer
    language plpgsql
as
$$
begin
    select key, fio, logo_uri from users.admins where key != key_;
end;
$$;

create or replace function users.if_admin_exists(key_ integer) returns boolean
    security definer
    language plpgsql
as
$$
begin
    select count(*) > 0 from users.admins where key = key_;
end;
$$;

create or replace function users.delete_coach(
    key_ integer
)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into users.deleted_coaches
    (key, FIO, date_reg, home_city, training_city, birthday, about, logo_uri, role)
            (select * from users.coaches where key = key_);
    delete from users.coaches where key = key_;
end;
$$;

create or replace function users.delete_pupil(
    key_ integer
)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into users.deleted_pupils
    (key, FIO, date_reg, coach, home_city, traqining_city, birthday, about, coach_review, logo_uri, role)
            (select * from users.pupils where key = key_);
    delete from users.pupils where key = key_;
end;
$$;