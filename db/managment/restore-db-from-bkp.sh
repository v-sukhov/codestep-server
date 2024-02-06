#!/bin/bash

if [$1 -eq 0]; then
	echo Restoring dump not found
	exit -1
fi

echo SAVING BACKUP
sudo -u postgres pg_dump -Fc codestep_db > dump-codestep_db-$(date +'%Y-%m-%d-%H-%M-%S')

if [$? -neq 0]; then
	echo Making backup failed. Exiting
	exit -1
fi

echo DROPING DATABASE
psql --host localhost --username postgres --password -c "DROP DATABASE codestep_db;"

if [$? -neq 0]; then
	echo Droping codestep_db database failed. Exiting
	exit -1
fi

echo CREATING DATABASE

if [$? -neq 0]; then
	echo Creating codestep_db database failed. Exiting
	exit -1
fi

psql --host localhost --username postgres --password -c "create database codestep_db owner = 'codestep' encoding = 'UTF8';"
sudo -u postgres pg_restore -d codestep_db $1