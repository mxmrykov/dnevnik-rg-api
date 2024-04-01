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
	    key bigint,
	    pupil integer[],
	    coach integer,
	    class_date text,
	    class_time varchar(5),
	    class_dur varchar(5),
	    presence boolean,
	    price smallint,
	    mark smallint,
	    review text,
	    scheduled boolean,
		classType varchar(10),
	    pupilCount int,
		isOpenToSignup boolean
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
    a.key, a.fio, a.date_reg, a.logo_uri, a.role, p.checksum, p.last_update, p.token
	FROM admins a LEFT JOIN passwords p on a.key = p.key WHERE a.key = $1;`
	GetAdminFull = `
	SELECT
    a.key, a.fio, a.date_reg, a.logo_uri, a.role, p.checksum, p.last_update, p.token
	FROM admins a LEFT JOIN passwords p on a.key = p.key WHERE a.key = $1;`
	GetAllAdminsExcept = `
	SELECT key, fio, logo_uri FROM admins WHERE key != $1; 
	`
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
	GetAllCoachesExcept = `
	SELECT key, fio, logo_uri FROM coach WHERE key != $1;
	`
	GetCoachPupils = `
	SELECT key, fio, logo_uri FROM pupil WHERE coach = $1;
	`
	IsCoachExists = `
	SELECT COUNT(*) FROM coach WHERE key = $1;`
	IsAdminExistsByName = `
	SELECT COUNT(*) FROM admins WHERE fio = $1;`
	DeleteCoach = `
	DELETE FROM coach WHERE key = $1;
	`
	CreatePupil = `
	INSERT INTO pupil (key, fio, date_reg, coach, home_city, training_city, birthday, about, coach_review, logo_uri, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) on conflict do nothing
	`
	GetPupil = `
	SELECT * FROM pupil WHERE key = $1;
	`
	GetPupilFull = `
	SELECT pupil.key, fio, date_reg, coach, home_city, training_city, birthday, about, coach_review, logo_uri, role
    FROM pupil WHERE key = $1;
	`
	IsPupilExists = `
	SELECT COUNT(*) FROM pupil WHERE key = $1;`
	IsPupilExistsByName = `
	SELECT COUNT(*) FROM pupil WHERE fio = $1;`
	DeletePupil = `
	DELETE FROM pupil WHERE key = $1;
	`
	GetAllPupils = `
	SELECT * FROM pupil
	`
	GetAllCoaches = `
	SELECT * FROM coach
	`
	GetAllAdmins = `
	SELECT * FROM admins
	`
	GetPasswordCheck = `
	SELECT checksum, token FROM passwords WHERE key = $1;
	`
	Auth = `
	SELECT 'ADMIN' AS role
	WHERE EXISTS (SELECT * FROM admins WHERE admins.key = $1)
	UNION
	SELECT 'COACH' AS role
	WHERE EXISTS (SELECT * FROM coach WHERE coach.key = $1)
	UNION
	SELECT 'PUPIL' AS role
	WHERE EXISTS (SELECT * FROM pupil WHERE pupil.key = $1);
	`
	GetCoachNearestBirthdays = `
	SELECT key, fio, birthday FROM pupil WHERE coach = $1;
	`
	GetCoachSchedule = `
	SELECT key, pupil, coach, class_date, class_time, class_dur FROM classes WHERE coach = $1 AND class_date = $2;
	`
	UpdateOldToken = `
	UPDATE passwords SET token = $2, last_update = $3 WHERE key = $1;
	`
	CreateClassIfNotExists = `
	INSERT Into classes (key, pupil, coach, class_date, class_time, class_dur, price, scheduled, classtype, pupilcount, isopentosignup)
	VALUES ($1, $2, $3, $4, $5, $6, $7, true, $8, $9, $10) returning key;
	`
	IfClassAvail = `
	SELECT COUNT(*) FROM classes WHERE coach = $1 AND class_date = $2 AND class_time = $3;
	`
	GetClassesForTodayAdmin = `
	SELECT key, pupil, coach, class_time, class_dur, classtype, pupilcount, scheduled, isopentosignup FROM classes WHERE class_date = $1 ORDER by class_time;
	`
	GetPupilsName = `
	SELECT fio FROM pupil WHERE key = any($1)
	`
	CancelClass = `
	UPDATE classes SET scheduled = false WHERE key = $1;
	`
)
