cd ./src
go build -o ../build/codestep-server
cd ..
mkdir -p build/db
cp -f server.conf.template ./build
cp -r db/changesets build/db/
cp -f db/postgresql-42.2.14.jar db/dbchangelog.xml db/liquibase.properties.template build/db/
