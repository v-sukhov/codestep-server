Liquibase migrations for codestep database.

Run manually from this directory:

	liquibase update

Or enable automatic update on server start in server.conf:

	liquibase_update = yes
	jdbc_driver_file = ./liquibase/postgresql-42.2.14.jar
