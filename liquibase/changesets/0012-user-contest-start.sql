--liquibase formatted sql

--changeset sukhov:20

CREATE TABLE T_USER_CONTEST_START (
    USER_ID INT NOT NULL,
    CONTEST_ID INT NOT NULL,
    USER_START_TIME TIMESTAMP WITH TIME ZONE,
    USER_FACT_START_TIME TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (user_id, contest_id)
);

COMMENT ON TABLE T_USER_CONTEST_START IS 'Время начала участия пользователя в соревновании. Хранит отсчётное и фактическое время начала соревнования пользователем';
COMMENT ON COLUMN T_USER_CONTEST_START.USER_ID IS 'ID пользователя';
COMMENT ON COLUMN T_USER_CONTEST_START.CONTEST_ID IS 'ID соревнования';
COMMENT ON COLUMN T_USER_CONTEST_START.USER_START_TIME IS 'Отсчётное время начала участия пользователя в соревновании';
COMMENT ON COLUMN T_USER_CONTEST_START.USER_FACT_START_TIME IS 'Время фактического начала участия пользователя в соревновании';

ALTER TABLE T_USER_CONTEST_START
    ADD CONSTRAINT t_user_contest_start_user_id_contest_id_unique UNIQUE (USER_ID, CONTEST_ID);

--rollback DROP TABLE T_USER_CONTEST_START;
