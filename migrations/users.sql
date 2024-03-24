CREATE TABLE IF NOT EXISTS pupil
(
    UDID          bigserial PRIMARY KEY,
    key           integer,
    FIO           text,
    date_reg      text,
    coach         integer,
    home_city     varchar(30),
    training_city varchar(30),
    birthday      text,
    about         text,
    coach_review  text,
    logo_uri      text,
    role          varchar(10)
);

CREATE TABLE IF NOT EXISTS coach
(
    UDID          bigserial PRIMARY KEY,
    key           integer,
    FIO           text,
    date_reg      text,
    home_city     varchar(30),
    training_city varchar(30),
    birthday      text,
    about         text,
    logo_uri      text,
    role          varchar(10)
);

CREATE TABLE IF NOT EXISTS admins
(
    UDID     bigserial PRIMARY KEY,
    key      integer,
    FIO      text,
    date_reg text,
    logo_uri text,
    role     varchar(10)
);

INSERT INTO pupil (udid, key, fio, date_reg, coach, home_city, training_city, birthday, about, coach_review,
                   logo_uri, role)
VALUES (4, 1704480003, 'Газарян Эвелина Владимировна', '2024-01-05T21:40:03+03:00', 1704463559, 'Краснодар',
        'Краснодар', '2013-10-24T00:00:00Z',
        'Тренируемся в спортивной школе 6 раз в неделю по 4 часа. Познакомились на сборах по художественной гимнастике в центре Винер в Новогорске , Эвелине тогда было 6  лет , Вы вели предметную подготовку групповых упражнений . Эвелине очень понравились ваши занятия и мы взяли дополнительные уроки с вами. После чего не смогли расстаться и стали заниматься онлайн . Сейчас уже 9 лет, мы освоили предметы , выучили много элементов , занимаемся над растяжкой , координацией , устойчивостью . Дарья Вадимовна очень добрый , чуткий , требовательный тренер . Знает подход к ребенку . Тренировки с Дарьей Вадимовной всегда в удовольствие и главное с результатом ❤️Вместе с 2020-го года.',
        'Очень талантливая девочка! За 3 года совместных тренировок, освоили с нуля предметную подготовку с обручем, мячом, булавами и работаем над лентой. Отработали сложные, соединенные элементы. Продолжаем стремиться к лучшему. Отрабатываем программы.',
        'https://dnevnik-rg.ru/pupil-logo.png', 'PUPIL');
INSERT INTO pupil (udid, key, fio, date_reg, coach, home_city, training_city, birthday, about, coach_review,
                   logo_uri, role)
VALUES (5, 1704568249, 'Зуева Полина Александровна', '2024-01-06T22:10:49+03:00', 1704463559, 'Находка', 'Находка',
        '2011-04-07T00:00:00Z',
        'Узнали про вас из инстаграм. С вами начали заниматься с 20.05.2022. Прогресс в работе с предметом, стопы и подъем стали более красивыми. И конечно же моральная поддержка от вас много значит❤️',
        'Проработали чистоту исполнения предметом и тела в упражнениях. Растягиваем стопы и тянем шпагаты.',
        'https://dnevnik-rg.ru/pupil-logo.png', 'PUPIL');
INSERT INTO coach (udid, key, fio, date_reg, home_city, training_city, birthday, about, logo_uri, role)
VALUES (4, 1704463559, 'Дубова Дарья Вадимовна', '2024-01-05T17:05:59+03:00', 'Воронеж', 'Москва',
        '1999-01-29T00:00:00Z',
        'Мастер спорта международного класса по художественной гимнастике, чемпионка юношеских Олимпийских игр в 2014 году, Чемпионка Европы 2013 года, победительница международных турниров. Чемпионка Мира, бронзовый призер Чемпионата Европы по эстетической гимнастике. Стаж работы тренером 7 лет. Преподает онлайн с 2020-го года.',
        'https://dnevnik-rg.ru/coach-logo.png', 'COACH');
INSERT INTO admins (udid, key, fio, date_reg, logo_uri, role)
VALUES (1, 1704461475, 'Рыков Максим Алексеевич', '2024-01-05T16:31:15+03:00', 'https://dnevnik-rg.ru/admin-logo.png',
        'ADMIN');
INSERT INTO admins (udid, key, fio, date_reg, logo_uri, role)
VALUES (4, 1704461712, 'Дубова Дарья Вадимовна', '2024-01-05T16:35:12+03:00', 'https://dnevnik-rg.ru/admin-logo.png',
        'ADMIN');