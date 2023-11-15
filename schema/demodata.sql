LOCK TABLES users WRITE;

INSERT INTO users (uuid, name, email, email_verified_at, created_at, last_login_at, deactivated_at, updated_at)
VALUES ('1-13-37', 'Demo User', 'demo.user@example.com', '2023-11-15 00:00:00', '2023-11-15 00:12:00', NULL, NULL, NULL);

UNLOCK TABLES;



LOCK TABLES password_credentials WRITE;

# password == "test": $2a$12$CiSlW.twXmObKXTX9jHrvOeElzaRsYQM202Zldoqszgp8QZwek6Ee
INSERT INTO password_credentials (user_uuid, password_hash, created_at, updated_at, last_asserted_at)
VALUES ('1-13-37', '$2a$12$CiSlW.twXmObKXTX9jHrvOeElzaRsYQM202Zldoqszgp8QZwek6Ee', '2023-11-15 00:00:20', NULL, NULL);

UNLOCK TABLES;