--liquibase formatted sql

--changeset sukhov:4

CREATE TABLE T_SUPERTASK
(
	SUPERTASK_ID		SERIAL	PRIMARY KEY,
	SUPERTASK_STATUS_ID	SMALLINT
);

COMMENT ON TABLE T_SUPERTASK IS 'Суперзадачи';
COMMENT ON COLUMN T_SUPERTASK.SUPERTASK_ID IS 'ID суперзадачи';

--

CREATE TABLE T_SUPERTASK_STATUS
(
	SUPERTASK_STATUS_ID		SMALLINT,
	SUPERTASK_STATUS_NAME	VARCHAR(256),
	SUPERTASK_STATUS_DESC	VARCHAR(512)
);

COMMENT ON TABLE T_SUPERTASK_STATUS IS 'Статусы суперзадачи';
COMMENT ON COLUMN T_SUPERTASK_STATUS.SUPERTASK_STATUS_ID IS 'ID статуса суперзадачи';
COMMENT ON COLUMN T_SUPERTASK_STATUS.SUPERTASK_STATUS_NAME IS 'Название статуса суперзадачи';
COMMENT ON COLUMN T_SUPERTASK_STATUS.SUPERTASK_STATUS_DESC IS 'Описание статуса суперзадачи';

INSERT INTO T_SUPERTASK_STATUS(SUPERTASK_STATUS_ID, SUPERTASK_STATUS_NAME, SUPERTASK_STATUS_DESC)
VALUES (0, 'Удалённая', 'Задача удалена');

INSERT INTO T_SUPERTASK_STATUS(SUPERTASK_STATUS_ID, SUPERTASK_STATUS_NAME, SUPERTASK_STATUS_DESC)
VALUES (1, 'Активная', 'Задача доступна для разработки и использования');


--

CREATE TABLE T_SUPERTASK_RIGHT_TYPE
(
	SUPERTASK_RIGHT_TYPE_ID		SMALLINT PRIMARY KEY,
	SUPERTASK_RIGHT_TYPE_NAME	VARCHAR(256),
	SUPERTASK_RIGHT_TYPE_DESC	VARCHAR(512)
);

COMMENT ON TABLE T_SUPERTASK_RIGHT_TYPE IS 'Типы прав пользователелей на суперзадачи';
COMMENT ON COLUMN T_SUPERTASK_RIGHT_TYPE.SUPERTASK_RIGHT_TYPE_ID IS 'ID типа права';
COMMENT ON COLUMN T_SUPERTASK_RIGHT_TYPE.SUPERTASK_RIGHT_TYPE_NAME IS 'Название типа права пользователелей на суперзадачи';
COMMENT ON COLUMN T_SUPERTASK_RIGHT_TYPE.SUPERTASK_RIGHT_TYPE_DESC IS 'Описание типа права пользователелей на суперзадачи';

INSERT INTO T_SUPERTASK_RIGHT_TYPE(SUPERTASK_RIGHT_TYPE_ID, SUPERTASK_RIGHT_TYPE_NAME, SUPERTASK_RIGHT_TYPE_DESC)
VALUES (1, 'Владелец', 'Полные права на разработку задачи и права на предоставление прав другим пользователям');

INSERT INTO T_SUPERTASK_RIGHT_TYPE(SUPERTASK_RIGHT_TYPE_ID, SUPERTASK_RIGHT_TYPE_NAME, SUPERTASK_RIGHT_TYPE_DESC)
VALUES (2, 'Разработчик', 'Полные права на разработку задачи');

INSERT INTO T_SUPERTASK_RIGHT_TYPE(SUPERTASK_RIGHT_TYPE_ID, SUPERTASK_RIGHT_TYPE_NAME, SUPERTASK_RIGHT_TYPE_DESC)
VALUES (3, 'Дизайнер', 'Права на изменение визуальных компонент задачи');

INSERT INTO T_SUPERTASK_RIGHT_TYPE(SUPERTASK_RIGHT_TYPE_ID, SUPERTASK_RIGHT_TYPE_NAME, SUPERTASK_RIGHT_TYPE_DESC)
VALUES (4, 'Наблюдатель', 'Права на просмотр кода задачи и клонирование');

--

CREATE TABLE T_SUPERTASK_USER_RIGHT
(
	SUPERTASK_ID	INTEGER,
	USER_ID			INTEGER,
	SUPERTASK_RIGHT_TYPE_ID	INTEGER
);

COMMENT ON TABLE T_SUPERTASK_USER_RIGHT IS 'Права пользователей на суперзадачи';
COMMENT ON COLUMN T_SUPERTASK_USER_RIGHT.SUPERTASK_ID IS 'ID суперзадачи';
COMMENT ON COLUMN T_SUPERTASK_USER_RIGHT.USER_ID IS 'ID пользователя';
COMMENT ON COLUMN T_SUPERTASK_USER_RIGHT.SUPERTASK_RIGHT_TYPE_ID IS 'ID типа прав';

CREATE INDEX T_SUPERTASK_USER_RIGHT_USER_INDEX ON T_SUPERTASK_USER_RIGHT(USER_ID);

--

CREATE TABLE T_SUPERTASK_VERSION
(
	SUPERTASK_ID			INTEGER,
	VERSION_NUMBER			INTEGER,
	PARENT_VERSION_NUMBER	INTEGER,
	COMMITED				BOOLEAN,
	AUTHOR_USER_ID			INTEGER,
	COMMIT_MESSAGE			VARCHAR(512),
	SAVE_DTM				TIMESTAMP,
	SUPERTASK_NAME			VARCHAR(256),
	SUPERTASK_DESC			VARCHAR(512),
	SUPERTASK_OBJECT_JSON	TEXT
);

COMMENT ON TABLE T_SUPERTASK_VERSION IS 'Сохраненныё версии суперзадачи. Каждый пользователь может просто сохранить суперзадачу - такая версия называется рабочей. Также пользователь может закоммитить версию - в этом случае его рабочая версия превращается в закоммиченную и в таком виде остаётся неизменной с новым номером версии.';
COMMENT ON COLUMN T_SUPERTASK_VERSION.SUPERTASK_ID IS 'ID суперзадачи';
COMMENT ON COLUMN T_SUPERTASK_VERSION.VERSION_NUMBER IS 'Номер версии (по порядку). У незакоммиченной версии - номер версии, от которой произошло копирование';
COMMENT ON COLUMN T_SUPERTASK_VERSION.PARENT_VERSION_NUMBER IS 'Номер родительской версии (от которой ответвилась данная версия).';
COMMENT ON COLUMN T_SUPERTASK_VERSION.COMMITED IS 'Признак закоммиченной версии';
COMMENT ON COLUMN T_SUPERTASK_VERSION.AUTHOR_USER_ID IS 'ID пользователя автора версии';
COMMENT ON COLUMN T_SUPERTASK_VERSION.COMMIT_MESSAGE IS 'Комментарий к коммиту';
COMMENT ON COLUMN T_SUPERTASK_VERSION.SAVE_DTM IS 'Время сохранения версии';
COMMENT ON COLUMN T_SUPERTASK_VERSION.SUPERTASK_NAME IS 'Название суперзадачи';
COMMENT ON COLUMN T_SUPERTASK_VERSION.SUPERTASK_DESC IS 'Описание суперзадачи';
COMMENT ON COLUMN T_SUPERTASK_VERSION.SUPERTASK_OBJECT_JSON IS 'Объект суперзадачи в формате JSON';

CREATE INDEX T_SUPERTASK_VERSION_SUPERTASK_ID ON T_SUPERTASK_VERSION(SUPERTASK_ID);

--rollback drop table T_SUPERTASK; drop table T_SUPERTASK_STATUS; drop table T_SUPERTASK_RIGHT_TYPE; drop table T_SUPERTASK_USER_RIGHT; drop table T_SUPERTASK_VERSION;


--changeset sukhov:5

CREATE TABLE T_SUPERTASK_LAST_VERSION
(
	SUPERTASK_ID		INTEGER	PRIMARY KEY,
	LAST_VERSION_NUMBER INTEGER
);

COMMENT ON TABLE T_SUPERTASK_LAST_VERSION IS 'Номер последней версии суперзадачи. Таблица нужна для двух целей: 1) лок на чтение при коммите без лока T_SUPERTASK_VERSION чтобы не влиять на производительность запросов от участников соревнований, 2) быстрый запрос последней закоммиченной версии суперзадачи. Таблица не должна участвовать в запросах от участников соревнований.';
COMMENT ON COLUMN T_SUPERTASK_LAST_VERSION.SUPERTASK_ID IS 'ID суперзадачи';
COMMENT ON COLUMN T_SUPERTASK_LAST_VERSION.LAST_VERSION_NUMBER IS 'Номер последней версии суперзадачи';

--

CREATE INDEX T_SUPERTASK_USER_RIGHT_SUPERTASK_INDEX ON T_SUPERTASK_USER_RIGHT(SUPERTASK_ID);
CREATE INDEX T_SUPERTASK_VERSION_SUPERTASK_INDEX ON T_SUPERTASK_VERSION(SUPERTASK_ID);

--changeset sukhov:6

ALTER TABLE T_SUPERTASK ADD COLUMN CREATE_DTM TIMESTAMP DEFAULT NOW();
ALTER TABLE T_SUPERTASK ADD COLUMN STATUS_CHANGE_DTM TIMESTAMP DEFAULT NOW();
COMMENT ON COLUMN T_SUPERTASK.CREATE_DTM IS 'Время создания суперзадачи';
COMMENT ON COLUMN T_SUPERTASK.STATUS_CHANGE_DTM IS 'Время перевода суперзадачи в текущий статус';