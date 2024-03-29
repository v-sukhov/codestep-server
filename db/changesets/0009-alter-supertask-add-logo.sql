--liquibase formatted sql

--changeset sukhov:15

ALTER TABLE T_SUPERTASK_VERSION ADD COLUMN SUPERTASK_LOGO_HREF VARCHAR(512);


--changeset sukhov:16

ALTER TABLE T_SUPERTASK_VERSION ALTER COLUMN SUPERTASK_LOGO_HREF SET DEFAULT '/img/default/default_logo.png';
ALTER TABLE T_SUPERTASK_VERSION ALTER COLUMN SUPERTASK_LOGO_HREF SET NOT NULL;
ALTER TABLE T_CONTEST ALTER COLUMN CONTEST_LOGO_HREF SET DEFAULT '/img/default/default_logo.png';
ALTER TABLE T_CONTEST ALTER COLUMN CONTEST_LOGO_HREF SET NOT NULL;

--changeset sukhov:17

COMMENT ON COLUMN T_SUPERTASK_VERSION.SUPERTASK_LOGO_HREF IS 'Логотип суперзадачи';
COMMENT ON COLUMN T_CONTEST.CONTEST_LOGO_HREF IS 'Логотип контеста';