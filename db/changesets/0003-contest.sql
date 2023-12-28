--liquibase formatted sql

--changeset sukhov:7

CREATE TABLE T_CONTEST
(
	CONTEST_ID			SERIAL	PRIMARY KEY,
	CONTEST_NAME		VARCHAR(256),
	CONTEST_DESC		VARCHAR(1024),
	CONTEST_LOGO_HREF	VARCHAR(512)
);

COMMENT ON TABLE T_CONTEST IS 'Соревнования';
COMMENT ON COLUMN T_CONTEST.CONTEST_ID IS 'ID соревнования';
COMMENT ON COLUMN T_CONTEST.CONTEST_NAME IS 'Название соревнования';
COMMENT ON COLUMN T_CONTEST.CONTEST_DESC IS 'Описание соревнования';
COMMENT ON COLUMN T_CONTEST.CONTEST_LOGO_HREF IS 'Гиперссылка на логотип соревнования';

--rollback drop table T_CONTEST;

--

CREATE TABLE T_CONTEST_RIGHT_TYPE
(
	CONTEST_RIGHT_TYPE_ID	SMALLINT PRIMARY KEY,
	CONTEST_RIGHT_TYPE_NAME	VARCHAR(256),
	CONTEST_RIGHT_TYPE_DESC	VARCHAR(512)
);

COMMENT ON TABLE T_CONTEST_RIGHT_TYPE IS 'Типы прав пользователей на соревнования';
COMMENT ON COLUMN T_CONTEST_RIGHT_TYPE.CONTEST_RIGHT_TYPE_ID IS 'ID типа права';
COMMENT ON COLUMN T_CONTEST_RIGHT_TYPE.CONTEST_RIGHT_TYPE_NAME IS 'Название типа права';
COMMENT ON COLUMN T_CONTEST_RIGHT_TYPE.CONTEST_RIGHT_TYPE_DESC IS 'Описание типа права';

INSERT INTO T_CONTEST_RIGHT_TYPE(CONTEST_RIGHT_TYPE_ID, CONTEST_RIGHT_TYPE_NAME, CONTEST_RIGHT_TYPE_DESC)
VALUES (1, 'Владелец', 'Полные права на управление соревнованием и на предоставление прав другим пользователям');

INSERT INTO T_CONTEST_RIGHT_TYPE(CONTEST_RIGHT_TYPE_ID, CONTEST_RIGHT_TYPE_NAME, CONTEST_RIGHT_TYPE_DESC)
VALUES (2, 'Администратор', 'Полные права на управление соревнованием и на предоставление прав другим пользователям');

INSERT INTO T_CONTEST_RIGHT_TYPE(CONTEST_RIGHT_TYPE_ID, CONTEST_RIGHT_TYPE_NAME, CONTEST_RIGHT_TYPE_DESC)
VALUES (3, 'Жюри', 'Просмотр результатов соревнования и работ участников. Может быть ограничен определёнными точками проведения');

INSERT INTO T_CONTEST_RIGHT_TYPE(CONTEST_RIGHT_TYPE_ID, CONTEST_RIGHT_TYPE_NAME, CONTEST_RIGHT_TYPE_DESC)
VALUES (4, 'Участник', 'Участие в соревновании');

--rollback drop table T_CONTEST_RIGHT_TYPE;

--

CREATE TABLE T_CONTEST_USER_RIGHT
(
	USER_ID					INTEGER,
	CONTEST_ID				INTEGER,
	CONTEST_RIGHT_TYPE_ID	INTEGER
);

COMMENT ON TABLE T_CONTEST_USER_RIGHT IS '';
COMMENT ON COLUMN T_CONTEST_USER_RIGHT.USER_ID IS 'ID пользователя';
COMMENT ON COLUMN T_CONTEST_USER_RIGHT.CONTEST_ID IS 'ID соревнования';
COMMENT ON COLUMN T_CONTEST_USER_RIGHT.CONTEST_RIGHT_TYPE_ID IS 'ID типа прав';

--rollback drop table T_CONTEST_USER_RIGHT;

--

CREATE TABLE T_CONTEST_SUPERTASK
(
	CONTEST_ID		INTEGER,
	SUPERTASK_ID	INTEGER,
	ORDER_NUMBER	SMALLINT
);

COMMENT ON TABLE T_CONTEST_SUPERTASK IS 'Суперзадачи в соревнованиях';
COMMENT ON COLUMN T_CONTEST_SUPERTASK.CONTEST_ID IS 'ID соревнования';
COMMENT ON COLUMN T_CONTEST_SUPERTASK.SUPERTASK_ID IS 'ID суперзадачи';
COMMENT ON COLUMN T_CONTEST_SUPERTASK.ORDER_NUMBER IS 'Порядковый номер суперзадачи в соревновании (нумерация с 1)';

CREATE UNIQUE INDEX T_CONTEST_SUPERTASK_INDEX_CONTEST_SUPERTASK ON T_CONTEST_SUPERTASK(CONTEST_ID, SUPERTASK_ID);
CREATE INDEX T_CONTEST_SUPERTASK_INDEX_CONTEST ON T_CONTEST_SUPERTASK(CONTEST_ID);

--rollback drop table T_CONTEST_SUPERTASK;
