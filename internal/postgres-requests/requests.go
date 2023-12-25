package postgres_requests

const (
	InitTablePupils = `
	CREATE TABLE IF NOT EXISTS pupil (
	    UDID bigserial PRIMARY KEY,
	    key integer,
	    FIO text,
	    date_reg text,
	    coach integer,
	    home_city varchar(30),
	    training_city varchar(30),
	    birthday text,
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
	    date_reg text,
	    home_city varchar(30),
	    training_city varchar(30),
	    birthday text,
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
	    class_date text,
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
	CREATE TABLE IF NOT EXISTS passwords (
	    UDID bigserial PRIMARY KEY,
	    key integer,
	    checksum varchar(8),
	    token text,
	    last_update text
	)
	`
	InitTableAdmins = `
	CREATE TABLE IF NOT EXISTS admins (
	    UDID bigserial PRIMARY KEY,
	    key integer,
	    FIO text,
	    date_reg text,
	    logo_uri text,
	    role varchar(10)
	)
	`
	NewAdmin = `
	INSERT INTO admins (key, fio, date_reg, logo_uri, role) VALUES ($1, $2, $3, $4, $5)
	`
	DeleteAdmin = `
	DELETE FROM admins WHERE key = $1
	`
	NewPassword = `
	INSERT INTO passwords (key, checksum, token, last_update) VALUES ($1, $2, $3, $4)
	`
	DeletePassword = `
	DELETE FROM passwords WHERE key = $1
	`
	GetAdmin = `
	SELECT
    a.key, a.fio, a.date_reg, a.logo_uri, p.checksum, p.last_update, p.token
	FROM admins a LEFT JOIN passwords p on a.key = p.key WHERE a.key = $1;`
	IsAdminExists = `
	SELECT COUNT(*) FROM admins WHERE key = $1;`
	CreateCoach = `
	INSERT INTO coach (key, fio, date_reg, home_city, training_city, birthday, about, logo_uri, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	GetCoach = `
	SELECT * FROM coach WHERE key = $1;
	`
	GetCoachFull = `
	SELECT coach.key, fio, date_reg, home_city, training_city, birthday, about, logo_uri, role, checksum, token, last_update
    FROM coach
	LEFT JOIN public.passwords p on coach.key = p.key WHERE coach.key = $1;
	`
	IsCoachExists = `
	SELECT COUNT(*) FROM coach WHERE key = $1;`
	IsAdminExistsByName = `
	SELECT COUNT(*) FROM admins WHERE fio = $1;`
	UpdateCoach = ``
	DeleteCoach = `
	DELETE FROM coach WHERE key = $1;
	`
)
