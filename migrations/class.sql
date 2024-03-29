CREATE TABLE IF NOT EXISTS classes
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
)