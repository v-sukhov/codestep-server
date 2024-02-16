--liquibase formatted sql

--changeset sukhov:18

ALTER TABLE T_SUPERTASK_VERSION ADD COLUMN TASKS_NUM INTEGER DEFAULT 1;
ALTER TABLE T_SUPERTASK_VERSION ADD COLUMN MAX_TOTAL_SCORE INTEGER DEFAULT 100;
ALTER TABLE T_SUPERTASK_VERSION ADD COLUMN MAX_TASK_SCORE VARCHAR(256) DEFAULT '100';

COMMENT ON COLUMN T_SUPERTASK_VERSION.TASKS_NUM IS 'Производное поле: кол-во задач в супер задаче';
COMMENT ON COLUMN T_SUPERTASK_VERSION.MAX_TOTAL_SCORE IS 'Производное поле: максимальный возможный балл за суперзадачу';
COMMENT ON COLUMN T_SUPERTASK_VERSION.MAX_TASK_SCORE IS 'Производное поле: максимальный возможный балл за каждую задачу - хранится в виде строки чисел, разделённых пробелом';
