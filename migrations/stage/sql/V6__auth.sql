drop function if exists auth.authorize(key_ integer, sum varchar(64));
create function auth.authorize(key_ integer, sum varchar(64))
    returns table
            (
                token text,
                role  varchar(10)
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select p.token, u.role
        from passwords.passwords as p
                 left join auth.get_user_role(key_) as u on u.key = key_
        where p.key = key_
          and p.checksum = sum;
end;
$$;

drop function if exists auth.get_user_role(key_ integer);
create or replace function auth.get_user_role(key_ integer)
    returns table
            (
                role varchar(10),
                key  integer
            )
    security definer
    language plpgsql
as
$$
begin
    return query
        select a.role, a.key
        from users.admins as a
        where a.key = key_
        union
        select c.role, c.key
        from users.coaches as c
        where c.key = key_
        union
        select p.role, p.key
        from users.pupils as p
        where p.key = key_;
end;
$$;

drop function if exists auth.auto_upd_token(key_ integer, token_ text, last_update_ text);
create or replace function auth.auto_upd_token(key_ integer, token_ text, last_update_ text)
    returns void
    security definer
    language plpgsql
as
$$
begin
    update passwords.passwords set token = token_, last_update = last_update_ WHERE key = key_;
end;
$$;

drop function if exists auth.select_user_private(key_ int, ip_ text, auth_type_ text);
create function auth.select_user_private(key_ integer, ip_ text, auth_type_ text)
    returns TABLE
            (
                checksum character varying,
                token    text
            )
    security definer
    language plpgsql
as
$$
begin
    insert into auth.auth_history(user_, attempt_type, ip, tm) values (key_, auth_type_, ip_, now()::timestamp);
    return query
        select p.checksum, p.token from passwords.passwords as p where p.key = key_;
end;
$$;