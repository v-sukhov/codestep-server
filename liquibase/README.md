Liquibase migrations for codestep database.

Connection settings are stored in `liquibase.properties` (independent from `server.conf`).

Run manually from this directory:

	liquibase update

Or enable automatic update on server start in server.conf:

	liquibase_update = yes

Before first run, copy `liquibase.properties.template` to `liquibase.properties` and set url, username and password.
