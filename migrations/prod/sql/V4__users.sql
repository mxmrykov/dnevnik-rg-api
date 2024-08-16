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
    select uc.key,
           uc.fio,
           uc.date_reg,
           uc.home_city,
           uc.training_city,
           uc.birthday,
           uc.about,
           uc.logo_uri,
           uc.role
    from users.coaches as uc
    where uc.key = key_;
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
    return query
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
                date_reg text,
                logo_uri text,
                role     varchar(10)
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select t.key, t.fio, t.date_reg, t.logo_uri, t.role from users.admins as t where t.key != key_;
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
    return query
        select * from users.admins;
end;
$$;

drop function if exists users.if_admin_exists(key_ integer);
create or replace function users.if_admin_exists(key_ integer)
    returns table
            (
                ex bool
            )
    security definer
    language plpgsql
as
$$
begin
    return query
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
            (select * from users.pupils up where up.key = key_);
    delete from users.pupils dp where dp.key = key_;
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
    return query
        select c.key,
               c.fio,
               c.date_reg,
               c.home_city,
               c.training_city,
               c.birthday,
               c.about,
               c.logo_uri,
               c.role,
               p.checksum,
               p.token,
               p.last_update
        from users.coaches as c
                 left join passwords.passwords p on c.key = p.key
        where c.key = key_;
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
    return query
        select up.key, up.fio, up.logo_uri from users.pupils as up where coach = key_;
end;
$$;

drop function if exists users.if_coach_exists(key_ integer);
create or replace function users.if_coach_exists(key_ integer)
    returns table
            (
                count bigint
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select count(*) from users.coaches where key = key_;
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
        select * from users.pupils up where up.key = key_;
end;
$$;

drop function if exists users.get_pupil_full(key_ integer);
create or replace function users.get_pupil_full(key_ integer)
    returns table
            (
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
               p.fio,
               p.date_reg,
               p.coach,
               p.home_city,
               p.training_city,
               p.birthday,
               p.about,
               p.coach_review,
               p.logo_uri,
               p.role,
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
    return query select u.key, u.fio, u.birthday from users.pupils as u where u.coach = coach_;
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

drop function if exists archive.archive_pupil(pupil_id int);
create or replace function archive.archive_pupil(pupil_id int)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into archive.archived_pupils
    (UDID, key, fio, date_reg, coach, home_city, training_city, birthday, about, coach_review, logo_uri, role)
            (select * from users.pupils up where up.key = pupil_id);
    delete from users.pupils dp where dp.key = pupil_id;
end;
$$;

drop function if exists archive.archive_coach(int);
create or replace function archive.archive_coach(coach_id int)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into archive.archived_coaches
    (UDID, key, fio, date_reg, home_city, training_city, birthday, about, logo_uri, role)
        (select * from users.coaches up where up.key = coach_id);
    delete from users.coaches dp where dp.key = coach_id;
end;
$$;

drop function if exists archive.get_pupils();
create or replace function archive.get_pupils()
    returns table
            (
                key int,
                fio text
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select ap.key, ap.fio from archive.archived_pupils ap;
end;
$$;

drop function if exists archive.get_coaches();
create or replace function archive.get_coaches()
    returns table
            (
                key int,
                fio text
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select ac.key, ac.fio from archive.archived_coaches ac;
end;
$$;

drop function if exists archive.dearchivate_pupil(int);
create or replace function archive.dearchivate_pupil(pupil_id int)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into users.pupils
    (UDID, key, fio, date_reg, coach, home_city, training_city, birthday, about, coach_review, logo_uri, role)
            (select * from archive.archived_pupils up where up.key = pupil_id);
    delete from archive.archived_pupils dp where dp.key = pupil_id;
end;
$$;

drop function if exists archive.dearchivate_coach(int);
create or replace function archive.dearchivate_coach(coach_id int)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into users.coaches
    (UDID, key, fio, date_reg, home_city, training_city, birthday, about, logo_uri, role)
        (select * from archive.archived_coaches up where up.key = coach_id);
    delete from archive.archived_coaches dp where dp.key = coach_id;
end;
$$;
