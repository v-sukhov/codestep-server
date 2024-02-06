# Взято отсюда: https://ubuntu.com/server/docs/databases-postgresql

sudo -u postgres psql template1

ALTER USER postgres with encrypted password 'your_password';

# Добавить в /etc/postgresql/*/main/pg_hba.conf

hostssl all postgres localhost scram-sha-256
hostssl codestep_db codestep localhost scram-sha-256

# Рестарт postgres

sudo service postgresql restart

(или sudo systemctl restart postgresql)

# Создать пользователя codestep и БД codestep_db
psql --host localhost --username postgres --password

create user codestep superuser password 'codestep';
create database codestep_db owner = 'codestep' encoding = 'UTF8';

# Восстановить дамп БД (см. https://www.postgresql.org/docs/current/backup-dump.html):
sudo -u postgres pg_restore -d codestep_db dumpfile