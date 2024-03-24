CREATE TABLE IF NOT EXISTS passwords
(
    UDID        bigserial PRIMARY KEY,
    key         integer,
    checksum    varchar(8),
    token       text,
    last_update text
);

INSERT INTO passwords (udid, key, checksum, token, last_update)
VALUES (1, 1704461475, 'c4ef9c3',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOjE3MDQ0NjE0NzUsImNoZWNrX3N1bSI6ImM0ZWY5YzMiLCJyb2xlIjoiMjAyNC0wMS0wNVQxNjozMToxNSswMzowMCIsImV4cCI6MTcwOTY0NTQ3NSwiaXNzIjoibG9jYWxob3N0OjgwMDAifQ.8nXxKUXXkptMuo5x5Cvjxwq2lEF8EC2nCvlcuqepJLY',
        '2024-01-05T16:31:15+03:00');
INSERT INTO passwords (udid, key, checksum, token, last_update)
VALUES (4, 1704461712, 'b9f94c7',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOjE3MDQ0NjE3MTIsImNoZWNrX3N1bSI6ImI5Zjk0YzciLCJyb2xlIjoiMjAyNC0wMS0wNVQxNjozNToxMiswMzowMCIsImV4cCI6MTcwOTY0NTcxMiwiaXNzIjoibG9jYWxob3N0OjgwMDAifQ.kY4yLF9sODk68y3u_xIncqC5OZbZQcZ6x8WzvHF62bg',
        '2024-01-05T16:35:12+03:00');
INSERT INTO passwords (udid, key, checksum, token, last_update)
VALUES (8, 1704463559, '81b073d',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOjE3MDQ0NjM1NTksImNoZWNrX3N1bSI6IjgxYjA3M2QiLCJyb2xlIjoiMjAyNC0wMS0wNVQxNzowNTo1OSswMzowMCIsImV4cCI6MTcwOTY0NzU1OSwiaXNzIjoibG9jYWxob3N0OjgwMDAifQ.JZRBZxIFAkl66B6pbRBP--lX9GwzgjLKkU9Dd3edIfs',
        '2024-01-05T17:05:59+03:00');
INSERT INTO passwords (udid, key, checksum, token, last_update)
VALUES (15, 1704480003, 'f1981e4',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOjE3MDQ0ODAwMDMsImNoZWNrX3N1bSI6ImYxOTgxZTQiLCJyb2xlIjoiMjAyNC0wMS0wNVQyMTo0MDowMyswMzowMCIsImV4cCI6MTcwOTY2NDAwMywiaXNzIjoibG9jYWxob3N0OjgwMDAifQ.bg3h9bvZ1sJK2SGkbP440Fo9T_dsTDlSOLthC_wg1eU',
        '2024-01-05T21:40:03+03:00');
INSERT INTO passwords (udid, key, checksum, token, last_update)
VALUES (16, 1704568249, '7298332',
        'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOjE3MDQ1NjgyNDksImNoZWNrX3N1bSI6IjcyOTgzMzIiLCJyb2xlIjoiMjAyNC0wMS0wNlQyMjoxMDo0OSswMzowMCIsImV4cCI6MTcwOTc1MjI0OSwiaXNzIjoibG9jYWxob3N0OjgwMDAifQ.Ngg1IHPiboaO8QSX--i_JmuKisb-bvWOxOsaZmh0X1k',
        '2024-01-06T22:10:49+03:00');