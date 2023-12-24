package postgres_requests

const (
	InitTablePupils = `
	CREATE TABLE IF NOT EXISTS pupil (
	    UDID bigserial PRIMARY KEY,
	    key integer,
	    FIO text,
	    date_reg date,
	    coach integer,
	    home_city varchar(30),
	    training_city varchar(30),
	    birthday date,
	    about text,
	    coach_review text,
	    logo_uri text,
	    role varchar(10)
	)
	`
	InitTableCoaches = `
	CREATE TABLE IF NOT EXISTS coach (
	    UDID bigserial PRIMARY KEY,
	    key integer,
	    FIO text,
	    date_reg date,
	    home_city varchar(30),
	    training_city varchar(30),
	    birthday date,
	    about text,
	    logo_uri text,
	    role varchar(10)
	)
	`
	InitTableClasses = `
	CREATE TABLE IF NOT EXISTS classes (
	    UDID bigserial PRIMARY KEY,
	    key integer,
	    pupil integer,
	    coach integer,
	    class_date date,
	    class_time varchar(5),
	    class_dur varchar(5),
	    presence boolean,
	    price smallint,
	    mark smallint,
	    review text,
	    scheduled boolean
	)
	`
	InitTablePasswords = `
	CREATE TABLE IF NOT EXISTS classes (
	    UDID bigserial PRIMARY KEY,
	    key integer,
	    checksum varchar(8),
	    last_update date
	)
	`
	CreateAdmin = `
	CREATE TABLE IF NOT EXISTS admins (
	    UDID bigserial PRIMARY KEY,
	    key integer,
	    FIO text,
	    date_reg date,
	    logo_uri text,
	    role varchar(10)
	)
	`
)
