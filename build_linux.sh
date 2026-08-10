mkdir -p build
cd ./src
go build -o ../build/codestep-server
cd ..
mkdir -p build/liquibase
cp -f server.conf.template ./build
cp -r liquibase/changesets build/liquibase/
cp -f liquibase/postgresql-42.2.14.jar liquibase/dbchangelog.xml liquibase/liquibase.properties.template build/liquibase/
