--liquibase formatted sql

--changeset sukhov:12

ALTER TABLE T_CONTEST_SUPERTASK ADD COLUMN SUPERTASK_VERSION_NUMBER INTEGER NOT NULL;
COMMENT ON COLUMN T_CONTEST_SUPERTASK.SUPERTASK_VERSION_NUMBER IS 'Версия суперзадачи';

CREATE UNIQUE INDEX ON T_CONTEST_SUPERTASK(CONTEST_ID, SUPERTASK_ID);