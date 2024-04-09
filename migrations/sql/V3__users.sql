drop function if exists users.create_admin(
    key_ integer,
    fio_ text,
    date_reg_ text,
    logo_uri_ text,
    role_ text
);

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
        (key, fio, date_reg, logo_uri, role)
    VALUES (key_, fio_, date_reg_, logo_uri_, role_)
    on conflict do nothing;
end;
$$;

drop function if exists users.create_coach(
    key_ integer,
    fio_ text,
    date_reg_ text,
    home_city_ varchar(30),
    training_city_ varchar(30),
    birthday_ text,
    about_ text,
    logo_uri_ text,
    role_ varchar(10)
);

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
    (key, fio, date_reg, home_city, training_city, birthday, about, logo_uri, role)
    VALUES (key_, fio_, date_reg_, home_city_, training_city_, birthday_, about_, logo_uri_, role_);
end;
$$;

drop function if exists users.get_coach(key_ integer);
create or replace function users.get_coach(key_ integer)
    returns table
            (
                key           integer,
                fio           text,
                date_reg_     text,
                home_city     varchar(30),
                training_city varchar(30),
                birthday      text,
                about         text,
                logo_uri      text,
                role          varchar(10)
            )
    security definer
    language plpgsql
as
$$
begin
    select key,
           fio,
           date_reg_,
           home_city,
           training_city,
           birthday,
           about,
           logo_uri,
           role
    from users.coaches
    where key = key_;
end;
$$;

drop function if exists users.delete_admin(key_ integer);
create or replace function users.delete_admin(
    key_ integer
)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into users.deleted_admins (UDID, key, fio, date_reg, logo_uri, role) (select * from users.admins where key = key_);
    delete from users.admins where key = key_;
end;
$$;

drop function if exists users.get_admin(key_ integer);
create or replace function users.get_admin(key_ integer)
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

drop function if exists users.list_admins_except(key_ integer);
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

drop function if exists users.list_admins();
create or replace function users.list_admins()
    returns table
            (
                UDID     integer,
                key      integer,
                fio      text,
                date_reg text,
                logo_uri text,
                role     varchar(10)
            )
    security definer
    language plpgsql
as
$$
begin
    select * from users.admins;
end;
$$;

drop function if exists users.if_admin_exists(key_ integer);
create or replace function users.if_admin_exists(key_ integer) returns boolean
    security definer
    language plpgsql
as
$$
begin
    select count(*) > 0 from users.admins where key = key_;
end;
$$;

drop function if exists users.delete_coach(key_ integer);
create or replace function users.delete_coach(key_ integer)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into users.deleted_coaches
    (UDID, key, fio, date_reg, home_city, training_city, birthday, about, logo_uri, role)
            (select * from users.coaches where key = key_);
    delete from users.coaches where key = key_;
end;
$$;

drop function if exists users.get_all_pupils();
create or replace function users.get_all_pupils()
    returns table
            (
                UDID          integer,
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
            )
    security definer
    language plpgsql
as
$$
begin
    return query select * from users.pupils;
end;
$$;

drop function if exists users.delete_pupil(key_ integer);
create or replace function users.delete_pupil(key_ integer)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into users.deleted_pupils
    (UDID, key, fio, date_reg, coach, home_city, training_city, birthday, about, coach_review, logo_uri, role)
            (select * from users.pupils where key = key_);
    delete from users.pupils where key = key_;
end;
$$;

drop function if exists users.get_coach(key_ integer);
create or replace function users.get_coach(key_ integer)
    returns table
            (
                UDID          integer,
                key           integer,
                fio           text,
                date_reg      text,
                home_city     varchar(30),
                training_city varchar(30),
                birthday      text,
                about         text,
                logo_uri      text,
                role          varchar(10)
            )
    security definer
    language plpgsql
as
$$
begin
    select * from users.coaches WHERE key = key_;
end;
$$;

drop function if exists users.get_coach_full(key_ integer);
create or replace function users.get_coach_full(key_ integer)
    returns table
            (
                key           integer,
                fio           text,
                date_reg      text,
                home_city     varchar(30),
                training_city varchar(30),
                birthday      text,
                about         text,
                logo_uri      text,
                role          varchar(10),
                checksum      varchar(64),
                token         text,
                last_update   text
            )
    security definer
    language plpgsql
as
$$
begin
    SELECT c.key,
           fio,
           date_reg,
           home_city,
           training_city,
           birthday,
           about,
           logo_uri,
           role,
           p.checksum,
           p.token,
           p.last_update
    FROM users.coaches as c
             LEFT JOIN passwords.passwords p on c.key = p.key
    WHERE c.key = key_;
end;
$$;

drop function if exists users.get_coach_pupils(key_ integer);
create or replace function users.get_coach_pupils(key_ integer)
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
    select key, fio, logo_uri from users.pupils where coach = key_;
end;
$$;

drop function if exists users.if_coach_exists(key_ integer);
create or replace function users.if_coach_exists(key_ integer)
    returns boolean
    security definer
    language plpgsql
as
$$
begin
    select count(*) > 0 from users.coaches where key = key_;
end;
$$;

drop function if exists users.create_pupil(
    key_ integer,
    fio_ text,
    date_reg_ text,
    coach_ integer,
    home_city_ varchar(30),
    training_city_ varchar(30),
    birthday_ text,
    about_ text,
    coach_review_ text,
    logo_uri_ text,
    role_ varchar(10)
);
create or replace function users.create_pupil(
    key_ integer,
    fio_ text,
    date_reg_ text,
    coach_ integer,
    home_city_ varchar(30),
    training_city_ varchar(30),
    birthday_ text,
    about_ text,
    coach_review_ text,
    logo_uri_ text,
    role_ varchar(10)
)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into users.pupils
    (key, fio, date_reg, coach, home_city, training_city, birthday, about, coach_review, logo_uri, role)
    values (key_,
            fio_,
            date_reg_,
            coach_,
            home_city_,
            training_city_,
            birthday_,
            about_,
            coach_review_,
            logo_uri_,
            role_)
    on conflict do nothing;
end;
$$;

drop function if exists users.get_pupil(key_ integer);
create or replace function users.get_pupil(key_ integer)
    returns table
            (
                UDID          integer,
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
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select * from users.pupils where key = key_;
end;
$$;

drop function if exists users.get_pupil_full(key_ integer);
create or replace function users.get_pupil_full(key_ integer)
    returns table
            (
                UDID          integer,
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
                role          varchar(10),
                checksum      varchar(64),
                token         text,
                last_update   text
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select p.key,
               fio,
               date_reg,
               coach,
               home_city,
               training_city,
               birthday,
               about,
               coach_review,
               logo_uri,
               role,
               ps.checksum,
               ps.token,
               ps.last_update
        from users.pupils as p
                 left join passwords.passwords ps on p.key = ps.key
        where p.key = key_;
end;
$$;

drop function if exists users.get_all_coaches();
create or replace function users.get_all_coaches()
    returns table
            (
                UDID          integer,
                key           integer,
                fio           text,
                date_reg      text,
                home_city     varchar(30),
                training_city varchar(30),
                birthday      text,
                about         text,
                logo_uri      text,
                role          varchar(10)
            )
    security definer
    language plpgsql
as
$$
begin
    return query select * from users.coaches;
end;
$$;

drop function if exists users.get_nearest_pupils_bd(coach_ integer);
create or replace function users.get_nearest_pupils_bd(coach_ integer)
    returns table
            (
                key      integer,
                fio      text,
                birthday text
            )
    security definer
    language plpgsql
as
$$
begin
    return query select key, fio, birthday from users.pupils as u where u.coach = coach_;
end;
$$;

drop function if exists users.get_pupils_names(pupils_ integer[]);
create or replace function users.get_pupils_names(pupils_ integer[])
    returns table
            (
                pupil text[]
            )
    security definer
    language plpgsql
as
$$
declare
    pupil_name text;
    pupil_id   integer;
begin
    foreach pupil_id in array pupils_
        loop
            select u.fio
            into pupil_name
            from users.pupils as u
            where key = pupil_id;

            if found then
                pupil := array_append(pupil, pupil_name);
            end if;
        end loop;
    return next;
end;
$$;