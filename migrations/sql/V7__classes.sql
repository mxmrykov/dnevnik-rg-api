drop function if exists classes.get_coach_schedule(coach_ integer, class_date_ text);
create or replace function classes.get_coach_schedule(coach_ integer, class_date_ text)
    returns table
            (
                key        bigint,
                pupil      integer[],
                coach      integer,
                class_date text,
                class_time varchar(5),
                class_dur  varchar(5)
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select c.key, c.pupil, c.coach, c.class_date, c.class_time, c.class_dur
        from classes.classes as c
        where c.coach = coach_
          and c.class_date = class_date_;
end;
$$;

drop function if exists classes.create_class_if_not_exists(
    key_ bigint,
    pupil_ integer[],
    coach_ integer,
    class_date_ text,
    class_time_ varchar(5),
    class_dur_ varchar(5),
    price_ smallint,
    class_type_ varchar(10),
    pupils_count_ integer,
    isopentosignup_ boolean
);
create or replace function classes.create_class_if_not_exists(
    key_ bigint,
    pupil_ integer[],
    coach_ integer,
    class_date_ text,
    class_time_ varchar(5),
    class_dur_ varchar(5),
    price_ smallint,
    class_type_ varchar(10),
    pupils_count_ integer,
    isopentosignup_ boolean
)
    returns table
            (
                key bigint
            )
    security definer
    language plpgsql
as
$$
begin
    insert into classes.classes as c (key, pupil, coach, class_date, class_time, class_dur, price, scheduled,
                                      classtype,
                                      pupilcount, isopentosignup)
    values (key_, pupil_, coach_, class_date_, class_time_, class_dur_, price_, true, class_type_,
            pupils_count_,
            isopentosignup_);
    return query
        select cl.key
        from classes.classes as cl
        where cl.class_date = class_date_
          and cl.class_time = class_time_
          and cl.coach = coach_;
end;
$$;

drop function if exists classes.if_class_available(coach_ integer, class_date_ text, class_time_ varchar(5));
create or replace function classes.if_class_available(coach_ integer, class_date_ text, class_time_ varchar(5))
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
        select count(*)
        from classes.classes
        where coach = coach_
          and class_date = class_date_
          and class_time = class_time_;
end;
$$;

drop function if exists classes.get_classes_for_today_admin(class_date_ text);
create or replace function classes.get_classes_for_today_admin(class_date_ text)
    returns table
            (
                key            bigint,
                pupil          integer[],
                coach          integer,
                class_time     varchar(5),
                class_dur      varchar(5),
                classtype      varchar(10),
                pupilcount     integer,
                scheduled      boolean,
                isopentosignup boolean
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select cl.key,
               cl.pupil,
               cl.coach,
               cl.class_time,
               cl.class_dur,
               cl.classtype,
               cl.pupilcount,
               cl.scheduled,
               cl.isopentosignup
        from classes.classes as cl
        where class_date = class_date_
        order by class_time;
end;
$$;

drop function if exists classes.cancel_class(class_id_ bigint);
create or replace function classes.cancel_class(class_id_ bigint)
    returns void
    security definer
    language plpgsql
as
$$
begin
    update classes.classes set scheduled = false where key = class_id_;
end;
$$;

drop function if exists classes.delete_class(class_id_ bigint);
create or replace function classes.delete_class(class_id_ bigint)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into classes.deleted_classes (key, pupil, coach, class_date, class_time, class_dur, presence, price, mark,
                                         review, scheduled, classtype, pupilcount, isopentosignup)
        (select c.key,
                c.pupil,
                c.coach,
                c.class_date,
                c.class_time,
                c.class_dur,
                c.presence,
                c.price,
                c.mark,
                c.review,
                c.scheduled,
                c.classtype,
                c.pupilcount,
                c.isopentosignup
         from classes.classes as c
         where c.key = class_id_);
    delete from classes.classes as cl where cl.key = class_id_;
end;
$$;

create or replace function classes.get_classes_for_month_admin(today text, lastDate text)
    returns table
            (
                key        bigint,
                class_date text,
                class_time varchar(5),
                class_dur  varchar(5)
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select cl.key, cl.class_date, cl.class_time, cl.class_dur
        from classes.classes cl
        where cl.class_date >= today
          and cl.class_date < lastDate
        group by cl.class_date, cl.class_time, cl.class_dur, cl.key;
end;
$$;

create or replace function classes.get_class_info(class_id bigint)
    returns table
            (
                key            bigint,
                pupil          integer[],
                coach          integer,
                class_date     text,
                class_time     varchar(5),
                class_dur      varchar(5),
                classtype      varchar(10),
                pupilcount     integer,
                scheduled      boolean,
                isopentosignup boolean
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select cl.key,
               cl.pupil,
               cl.coach,
               cl.class_date,
               cl.class_time,
               cl.class_dur,
               cl.classtype,
               cl.pupilcount,
               cl.scheduled,
               cl.isopentosignup
        from classes.classes as cl
        where cl.key = class_id;
end;
$$;

drop function if exists classes.get_classes_for_today_pupil(class_date_ text, user_id int);
create or replace function classes.get_classes_for_today_pupil(class_date_ text, user_id int)
    returns table
            (
                key            bigint,
                class_date     text,
                class_time     varchar(5),
                class_dur      varchar(5),
                coach          text,
                classtype      varchar(10),
                pupilcount     integer,
                scheduled      boolean,
                capacity       int,
                isopentosignup boolean
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select cl.key,
               cl.class_date,
               cl.class_time,
               cl.class_dur,
               (select c.fio from users.coaches as c where c.key = cl.coach),
               cl.classtype,
               cl.pupilcount,
               cl.scheduled,
               array_length(cl.pupil, 1) as capacity,
               cl.isopentosignup
        from classes.classes as cl
        where cl.class_date = class_date_
          and user_id = any (cl.pupil)
        order by class_time;
end;
$$;

drop function if exists classes.get_admin_classes_history(datefrom text);
create or replace function classes.get_admin_classes_history(datefrom text)
    returns table
            (
                key            bigint,
                class_date     text,
                class_time     varchar(5),
                class_dur      varchar(5),
                coach          text,
                coach_key      int,
                classtype      varchar(10),
                pupils         text[],
                pupils_keys    int[],
                scheduled      boolean,
                capacity       int,
                isopentosignup boolean
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select cl.key,
               cl.class_date,
               cl.class_time,
               cl.class_dur,
               (select c.fio from users.coaches as c where c.key = cl.coach),
               cl.coach                            as coach_key,
               cl.classtype,
               ARRAY(SELECT u.fio
                     FROM users.pupils AS u
                     WHERE u.key = ANY (cl.pupil)) AS pupils,
               cl.pupil                            as pupils_keys,
               cl.scheduled,
               cl.pupilcount                       as capacity,
               cl.isopentosignup
        from classes.classes as cl
        where cl.class_date < datefrom
        order by class_date, class_time desc;
end;
$$;