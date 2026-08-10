if not exist builde mkdir build
cd ./src
go build -o ..\build\codestep-server.exe
cd ..
if not exist build\liquibase mkdir build\liquibase
copy /Y server.conf.template .\build
xcopy /E /I /Y db\changesets build\liquibase\changesets\
copy /Y db\postgresql-42.2.14.jar build\liquibase\
copy /Y db\dbchangelog.xml build\liquibase\
copy /Y db\liquibase.properties.template build\liquibase\
