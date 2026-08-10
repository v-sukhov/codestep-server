cd ./src
go build -o ../build/codestep-server
cd ..
mkdir -p build/liquibase
cp -f server.conf.template ./build
cp -r db/changesets build/liquibase/
cp -f db/postgresql-42.2.14.jar db/dbchangelog.xml db/liquibase.properties.template build/liquibase/
