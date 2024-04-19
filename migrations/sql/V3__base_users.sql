insert into users.admins (udid, key, fio, date_reg, logo_uri, role)
values (1, 1713005383, 'Рыков Максим Алексеевич', '2024-01-05T16:31:15+03:00', 'https://dnevnik-rg.ru/admin-logo.png',
        'ADMIN');
insert into passwords.passwords (udid, key, checksum, token, last_update)
values (1, 1713005383, 'NjUzMjk5OThmNmIzZWVlNzZmMWJjMjJmYmNkOWViZWI=',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOjE3MTMwMDUzODMsImNoZWNrX3N1bSI6Ik5qVXpNams1T1RobU5tSXpaV1ZsTnpabU1XSmpNakptWW1Oa09XVmlaV0k9Iiwicm9sZSI6IkFETUlOIiwiZXhwIjoxNzE4MTg5MzgzLCJpc3MiOiJsb2NhbGhvc3Q6ODAwMCJ9.ggHKSCabhdNxo070hTJ7gA85JQ9TF95QqdxeoRxaxas',
        '2024-01-05T16:31:15+03:00');
INSERT INTO users.coaches (udid, key, fio, date_reg, home_city, training_city, birthday, about, logo_uri, role)
VALUES (2, 1713556382, 'Дубова Дарья Вадимовна', '2024-04-19T22:53:02+03:00', 'Воронеж', 'Москва',
        '1999-01-29T00:00:00Z',
        'Мастер спорта международного класса по художественной гимнастике, чемпионка юношеских Олимпийских игр в 2014 году, Чемпионка Европы 2013 года, победительница международных турниров. Чемпионка Мира, бронзовый призер Чемпионата Европы по эстетической гимнастике. Стаж работы тренером 7 лет. Преподает онлайн с 2020-го года.',
        'https://dnevnik-rg.ru/coach-logo.png', 'COACH');
INSERT INTO passwords.passwords (udid, key, checksum, token, last_update)
VALUES (2, 1713556382, 'NzBhYWNiMzA2MmUwNjAzOTlhMTZiZjljMTQzYzU4NDI=',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOjE3MTM1NTYzODIsImNoZWNrX3N1bSI6Ik56QmhZV05pTXpBMk1tVXdOakF6T1RsaE1UWmlaamxqTVRRell6VTROREk9Iiwicm9sZSI6IkNPQUNIIiwiZXhwIjoxNzE4NzQwMzgyLCJpc3MiOiJsb2NhbGhvc3Q6ODAwMCJ9.X8M4ulMZ1Dhsk7Fn5nPQumU3Zjz5jpS2HgSo-M-Sdyg',
        '2024-04-19T22:53:02+03:00');
