insert into users.admins (udid, key, fio, date_reg, logo_uri, role)
values (1, 1713005383, 'Рыков Максим Алексеевич', '2024-01-05T16:31:15+03:00', 'https://dnevnik-rg.ru/admin-logo.png',
        'ADMIN');
insert into passwords.passwords (udid, key, checksum, token, last_update)
values (1, 1713005383, 'NjUzMjk5OThmNmIzZWVlNzZmMWJjMjJmYmNkOWViZWI=',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOjE3MTMwMDUzODMsImNoZWNrX3N1bSI6Ik5qVXpNams1T1RobU5tSXpaV1ZsTnpabU1XSmpNakptWW1Oa09XVmlaV0k9Iiwicm9sZSI6IkFETUlOIiwiZXhwIjoxNzE4MTg5MzgzLCJpc3MiOiJsb2NhbGhvc3Q6ODAwMCJ9.ggHKSCabhdNxo070hTJ7gA85JQ9TF95QqdxeoRxaxas',
        '2024-01-05T16:31:15+03:00');


INSERT INTO users.coaches (udid, key, fio, date_reg, home_city, training_city, birthday, about, logo_uri, role)
VALUES (1, 1713556382, 'Дубова Дарья Вадимовна', '2024-04-19T22:53:02+03:00', 'Воронеж', 'Москва',
        '1999-01-29T00:00:00Z',
        'Мастер спорта международного класса по художественной гимнастике, чемпионка юношеских Олимпийских игр в 2014 году, Чемпионка Европы 2013 года, победительница международных турниров. Чемпионка Мира, бронзовый призер Чемпионата Европы по эстетической гимнастике. Стаж работы тренером 7 лет. Преподает онлайн с 2020-го года.',
        'https://dnevnik-rg.ru/coach-logo.png', 'COACH');
INSERT INTO passwords.passwords (udid, key, checksum, token, last_update)
VALUES (2, 1713556382, 'NzBhYWNiMzA2MmUwNjAzOTlhMTZiZjljMTQzYzU4NDI=',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOjE3MTM1NTYzODIsImNoZWNrX3N1bSI6Ik56QmhZV05pTXpBMk1tVXdOakF6T1RsaE1UWmlaamxqTVRRell6VTROREk9Iiwicm9sZSI6IkNPQUNIIiwiZXhwIjoxNzE4NzQwMzgyLCJpc3MiOiJsb2NhbGhvc3Q6ODAwMCJ9.X8M4ulMZ1Dhsk7Fn5nPQumU3Zjz5jpS2HgSo-M-Sdyg',
        '2024-04-19T22:53:02+03:00');



INSERT INTO users.pupils (udid, key, fio, date_reg, coach, home_city, training_city, birthday, about, coach_review,
                          logo_uri, role)
VALUES (1, 1713604472, 'Гераськова Милена Евгеньевна', '2024-04-20T09:14:32Z', 1713556382, 'Тамбов', 'Тамбов',
        '2013-04-16T00:00:00Z',
        'Занимаемся с 2020-го года. За годы проведенных тренировок, проработали разные ошибки, включая работу с предметом, элементы и программы.',
        'Умная девочка! За 3 года работы онлайн прошли и выучили многое, начиная от элементов тела , заканчивая предметом. На тренировках с Миленой много разговариваем. Чтобы спортсменка понимала, что, для чего и почему.Сделала большой прорыв с работой предмета. Базовую предметную подготовку с булавами освоили с нуля.',
        'https://dnevnik-rg.ru/pupil-logo.png', 'PUPIL');
INSERT INTO passwords.passwords (udid, key, checksum, token, last_update)
VALUES (3, 1713604472, 'ZDQyMGJkMWI5YjI3MzM4NzIxNDVkNmM2N2EzNDcwMzc=',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOjE3MTM2MDQ0NzIsImNoZWNrX3N1bSI6IlpEUXlNR0prTVdJNVlqSTNNek00TnpJeE5EVmtObU0yTjJFek5EY3dNemM9Iiwicm9sZSI6IlBVUElMIiwiZXhwIjoxNzE4Nzg4NDcyLCJpc3MiOiJsb2NhbGhvc3Q6ODAwMCJ9.tmMyAr8Ail_vVmX14DL21t0xsUL-lNc3Cz6pTc2KfJg',
        '2024-04-20T09:14:32Z');
