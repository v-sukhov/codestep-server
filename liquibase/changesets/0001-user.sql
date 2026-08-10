--liquibase formatted sql

--changeset sukhov:1

create table T_USER
(
	user_id			SERIAL			PRIMARY KEY,
	login			VARCHAR(256)	NOT NULL,
	email			VARCHAR(256),
	password_sha256	BYTEA,
	user_type		SMALLINT,
	surname			VARCHAR(256),
	first_name		VARCHAR(256),
	second_name		VARCHAR(256),
	display_name	VARCHAR(256)
);

COMMENT ON TABLE T_USER IS 'Все пользователи системы (включая регулярных пользователей и участников)';
COMMENT ON COLUMN T_USER.user_id IS 'Уникальный идентификатор';
COMMENT ON COLUMN T_USER.login IS 'Логин пользователя. Для регулярных пользователей должен совпадать с почтовым адресом (кроме специальных случаев)';
COMMENT ON COLUMN T_USER.email IS 'EMail пользователя. Для регулярных пользователей должен совпадать с логином (кроме специальных случаев)';
COMMENT ON COLUMN T_USER.password_sha256 IS 'Пароль, зашифрованный при помощи секретного salt-слова и хэш-алгоритма';
COMMENT ON COLUMN T_USER.user_type IS 'Тип пользователя: 1 - регулярный пользователь, 2 - участник';
COMMENT ON COLUMN T_USER.surname IS '';
COMMENT ON COLUMN T_USER.first_name IS '';
COMMENT ON COLUMN T_USER.second_name IS '';
COMMENT ON COLUMN T_USER.display_name IS '';

--changeset sukhov:2

CREATE UNIQUE INDEX T_USER_LOGIN_INDEX ON T_USER(LOGIN);

--rollback drop index T_USER_LOGIN_INDEX;

--changeset sukhov:3

CREATE TABLE T_USER_RIGHTS
(
	user_id 		INTEGER		PRIMARY KEY,
	is_admin		BOOLEAN	NOT NULL DEFAULT FALSE,
	is_developer	BOOLEAN	NOT NULL DEFAULT FALSE,
	is_jury			BOOLEAN	NOT NULL DEFAULT FALSE
);

COMMENT ON TABLE T_USER_RIGHTS IS 'Права регулярного пользователя';
COMMENT ON COLUMN T_USER_RIGHTS.user_id IS 'ID';
COMMENT ON COLUMN T_USER_RIGHTS.is_admin IS 'Права администратора';
COMMENT ON COLUMN T_USER_RIGHTS.is_developer IS 'Права разработчика задач';
COMMENT ON COLUMN T_USER_RIGHTS.is_jury IS 'Права жюри';

--rollback drop table T_USER_RIGHTS;


