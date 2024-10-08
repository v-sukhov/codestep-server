# codestep-server
Server side for codestep project

## Разработка

1. Установить go
2. Установить postgres
3. Установить liquibase.
4. Установить go-swagger:
https://github.com/go-swagger/go-swagger

## Установка без дистрибутива

1. Установить и настроить postgres - см. db/management
2. Установить liquibase:
	https://docs.liquibase.com/start/install/liquibase-linux-debian-ubuntu.html
3. Установить пакет сервера в /var/codestep/
4. Задать параметры в server.conf
5. Настроить сервис - см. service-script
6. Чтобы можно было использовать 443 порт нужно выполнить
	sudo sysctl net.ipv4.ip_unprivileged_port_start=443
	(см. https://stackoverflow.com/questions/413807/is-there-a-way-for-non-root-processes-to-bind-to-privileged-ports-on-linux)

## Установка из deb-дистрибутива
1. Установить и настроить postgres - см. db/management
2. Установить liquibase:
	https://docs.liquibase.com/start/install/liquibase-linux-debian-ubuntu.html
3. Чтобы можно было использовать 443 порт нужно выполнить
	sudo sysctl net.ipv4.ip_unprivileged_port_start=443
	(см. https://stackoverflow.com/questions/413807/is-there-a-way-for-non-root-processes-to-bind-to-privileged-ports-on-linux)	
	
	Ещё можно выполнить 
		
	sudo ufw allow 22
	sudo ufw allow 443/tcp
	sudo ufw enable
	
	Но возможно ufw не требуется.
	
4. Установить codestep...amd64.deb
5. Задать параметры в server.conf
6. В /var/codestep/liquibase выполнить:

	liquibase update
	
7. Рестартовать сервис:

	sudo systemctl restart codestep

## Выпуск и обновление сертификата Let's encrypt от Яндекс-Облака

1. Получить OAuth-токен, как написано здесь (он не должен меняться вообще-то):

https://yandex.cloud/ru/docs/iam/operations/iam-token/create#api_1

2. Как написано там же получить IAM-токен с помощью полученного на шаге 1 токена. 
Для этого использовать сохранённый POST-запрос в Postman. Идентификатор сертификата не должен меняться.
Его можно скопировать, зайдя в сертификат в Яндекс-Облаке (в Certificate Manager).

3. Как написано здесь:

https://yandex.cloud/ru/docs/certificate-manager/api-ref/CertificateContent/get

сделать запрос (через Postman) и получить сертификат и секретный ключ. В качестве авторизационного
Bearer-токена использовать токен, полученный на шаге 2.

4. Скопировать из ответа цепочку сертификатов и вставить в файл code-step.pem (удалить кавычки,
запятые, и перевести \n в переводы строк). Скопировать секретный ключ code-step-private-key.pem
(выполнить аналогичные преобразования).


## Материалы

### Материалы по go-swagger:
https://goswagger.io/use/spec.html
https://medium.com/@pedram.esmaeeli/generate-swagger-specification-from-go-source-code-648615f7b9d9

### Материалы по работе с sql.db и Postgres:
http://go-database-sql.org/retrieving.html
https://habr.com/ru/company/oleg-bunin/blog/461935/

### Интересный ролик по истории Go:
https://youtu.be/ql-uncsqoAU

