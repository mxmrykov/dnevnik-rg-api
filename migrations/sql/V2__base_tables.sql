create table if not exists users.admins
(
    UDID     integer primary key default nextval('users.admins_id_seq'),
    key      integer,
    fio      text,
    date_reg text,
    logo_uri text,
    role     varchar(10)
);

create table if not exists users.coaches
(
    UDID          integer primary key default nextval('users.coaches_id_seq'),
    key           integer,
    fio           text,
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
    UDID          integer primary key default nextval('users.pupils_id_seq'),
    key           integer,
    fio           text,
    date_reg      text,
    coach         integer,
    home_city     varchar(30),
    training_city varchar(30),
    birthday      text,
    about         text,
    coach_review  text,
    logo_uri      text,
    role          varchar(10)
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
    UDID        integer primary key default nextval('passwords.passwords_id_seq'),
    key         integer,
    checksum    varchar(64),
    token       text,
    last_update text
);

create table if not exists passwords.deleted_passwords as
select *
from passwords.passwords
where false;

create table if not exists classes.classes
(
    UDID           integer primary key default nextval('classes.classes_id_seq'),
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

create table if not exists auth.auth_history
(
    id           integer primary key default nextval('auth.auth_id_seq'),
    user_        int,
    attempt_type text,
    ip           text,
    tm           timestamp
);

create table if not exists classes.deleted_classes as
select *
from classes.classes
where false;