drop function if exists classes.get_coach_schedule(coach_ integer, class_date_ text);
create or replace function classes.get_coach_schedule(coach_ integer, class_date_ text)
    returns table
            (
                key        integer,
                pupil      integer,
                coach      integer,
                class_date text,
                class_time text,
                class_dur  text
            )
    security definer
    language plpgsql
as
$$
begin
    select key, pupil, coach, class_date, class_time, class_dur
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
                key integer
            )
    security definer
    language plpgsql
as
$$
begin
    insert into classes.classes (key, pupil, coach, class_date, class_time, class_dur, price, scheduled, classtype,
                                 pupilcount, isopentosignup)
    values (key_, pupil_, coach_, class_date_, class_time_, class_dur_, price_, true, class_type_, pupils_count_,
            isopentosignup_)
    returning key;
end;
$$;

drop function if exists classes.if_class_available(coach_ integer, class_date_ text, class_time_ varchar(5));
create or replace function classes.if_class_available(coach_ integer, class_date_ text, class_time_ varchar(5))
    returns table
            (
                count integer
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select count(*) from classes where coach = coach_ and class_date = class_date_ and class_time = class_time_;
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

drop function if exists classes.cancel_class(class_id_ integer);
create or replace function classes.cancel_class(class_id_ integer)
    returns void
    security definer
    language plpgsql
as
$$
begin
    update classes set scheduled = false where key = class_id_;
end;
$$;