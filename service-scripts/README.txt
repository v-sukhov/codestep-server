1. Создать пользователя codestep (см. create-user-codestep.sh)

2. Скопирвать codestep.service в /etc/systemd/system/

3. Выполнить:

	sudo systemctl daemon-reload
	sudo systemctl enable codestep