cd ./src
go build -o ..\build\codestep-server.exe
cd ..
if not exist build\db mkdir build\db
copy /Y server.conf.template .\build
xcopy /E /I /Y db\changesets build\db\changesets\
copy /Y db\postgresql-42.2.14.jar build\db\
copy /Y db\dbchangelog.xml build\db\
copy /Y db\liquibase.properties.template build\db\
