--liquibase formatted sql

--changeset sukhov:13

create table T_USERGROUP
(
	usergroup_id	SERIAL				PRIMARY KEY,
	usergroup_name	VARCHAR(256)		NOT NULL,
	usergroup_desc	VARCHAR(512)
);

COMMENT ON TABLE T_USERGROUP IS 'Группы пользователей';
COMMENT ON COLUMN T_USERGROUP.usergroup_id IS 'ID группы пользователей';
COMMENT ON COLUMN T_USERGROUP.usergroup_name IS 'Название группы пользователей. Разрешены только маленькие латинские буквы, цифры и дефисы';
COMMENT ON COLUMN T_USERGROUP.usergroup_desc IS 'Описание группы пользователей';

CREATE UNIQUE INDEX T_USERGROUP_NAME_INDEX ON T_USERGROUP(usergroup_name);

--rollback drop T_USERGROUP;

--changeset sukhov:14

create table T_USER_USERGROUP
(
	usergroup_id	INTEGER REFERENCES T_USERGROUP(usergroup_id) ON DELETE CASCADE,
	user_id			INTEGER REFERENCES T_USER(user_id) ON DELETE CASCADE
);

COMMENT ON TABLE T_USER_USERGROUP IS 'Включение пользователей в группы пользователей';
COMMENT ON COLUMN T_USER_USERGROUP.usergroup_id IS 'ID группы пользователей';
COMMENT ON COLUMN T_USER_USERGROUP.user_id IS 'ID пользователя';

CREATE UNIQUE INDEX T_USER_USERGROUP_INDEX ON T_USER_USERGROUP(usergroup_id, user_id);
CREATE INDEX T_USER_USERGROUP_UG_INDEX ON T_USER_USERGROUP(usergroup_id);

--rollback drop T_USER_USERGROUP;