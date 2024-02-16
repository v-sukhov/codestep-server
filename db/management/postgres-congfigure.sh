# Взято отсюда: https://ubuntu.com/server/docs/databases-postgresql

sudo -u postgres psql template1

ALTER USER postgres with encrypted password 'your_password';

# Добавить в /etc/postgresql/*/main/pg_hba.conf

hostssl all postgres localhost scram-sha-256
hostssl codestep_db codestep localhost scram-sha-256
hostssl codestep_db codestep_readonly 0.0.0.0/0 scram-sha-256

# Рестарт postgres

# если нет systemctl то 
# sudo service postgresql restart
sudo systemctl restart postgresql

(или sudo systemctl restart postgresql)

# Создать пользователя codestep и БД codestep_db

psql --host localhost --username postgres --password

create user codestep password 'codestep';
create database codestep_db owner = 'codestep' encoding = 'UTF8';

# Создать пользователя codestep_readonly с правами только на чтение

psql --host localhost --username postgres --password -d codestep_db

create user codestep_readonly password 'codestep_readonly';
grant select on all tables in schema "public" to codestep_readonly;

# Восстановить дамп БД (см. https://www.postgresql.org/docs/current/backup-dump.html):
sudo -u postgres pg_restore -d codestep_db dumpfile