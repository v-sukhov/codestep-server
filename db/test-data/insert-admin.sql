insert into T_USER(
	login,
	password_sha256,
	user_type,
	surname,
	first_name,
	second_name,
	display_name
)
VALUES
(
	'admin',
	sha256('salt-1'),
	1,
	'admin_surname',
	'admin_firstname',
	'admin_secondname',
	'codestep admin'
)
