create or replace function passwords.new_password_row(
    key_ integer,
    checksum_ varchar(64),
    token_ text,
    last_update_ text
)
    returns void
    security definer
    language plpgsql
as
$$
begin
    INSERT INTO passwords.passwords
        (key, checksum, token, last_update)
    VALUES (key_, checksum_, token_, last_update_);
end;
$$;

create or replace function passwords.delete_password_row(key_ integer)
    returns void
    security definer
    language plpgsql
as
$$
begin
    insert into passwords.deleted_passwords (key, checksum, token, last_update) (select * from passwords.passwords where key = key_);
    delete from passwords.passwords where key = key_;
end;
$$;