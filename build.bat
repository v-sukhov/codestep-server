cd ./src
go build -o ..\build\codestep-server.exe
cd ..
copy /Y server.conf.template .\build
