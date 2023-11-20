CREATE TABLE users (
    uuid VARCHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NULL UNIQUE,
    email_verified_at DATETIME NULL,
    created_at DATETIME NOT NULL,
    last_login_at DATETIME NULL,
    deactivated_at DATETIME NULL,
    updated_at DATETIME NULL,

    INDEX (deactivated_at)
) CHARSET = utf8, ENGINE = InnoDB;

CREATE TABLE field_verifications (
    uuid VARCHAR(36) NOT NULL PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    user_uuid VARCHAR(36) DEFAULT NULL,
    field_name VARCHAR(50) NOT NULL,
    field_value VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,

    INDEX (created_at)
) CHARSET = utf8, ENGINE = InnoDB;

CREATE TABLE password_credentials (
    user_uuid VARCHAR(36) NOT NULL PRIMARY KEY,
    password_hash VARCHAR(100) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME DEFAULT NULL,
    last_asserted_at DATETIME DEFAULT NULL,

    INDEX (last_asserted_at),

    CONSTRAINT FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
) CHARSET = utf8, ENGINE = InnoDB;

CREATE TABLE sessions (
    id VARCHAR(255) NOT NULL PRIMARY KEY,
    user_uuid VARCHAR(36) NOT NULL,
    created_at DATETIME NOT NULL,
    last_seen_at DATETIME NULL,

    INDEX (last_seen_at),

    CONSTRAINT FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
) CHARSET = utf8, ENGINE = InnoDB;

CREATE TABLE channels (
    uuid VARCHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(50) NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME DEFAULT NULL
) CHARSET utf8, ENGINE InnoDB;

CREATE TABLE channel_members (
    channel_uuid VARCHAR(36) NOT NULL,
    user_uuid VARCHAR(36) NOT NULL,
    role VARCHAR(10) NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME DEFAULT NULL,

    PRIMARY KEY (channel_uuid, user_uuid),

    INDEX (channel_uuid),
    INDEX (user_uuid),

    CONSTRAINT FOREIGN KEY (channel_uuid) REFERENCES channels(uuid) ON DELETE CASCADE,
    CONSTRAINT FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
) CHARSET utf8, ENGINE InnoDB;

CREATE TABLE messages (
    uuid VARCHAR(36) NOT NULL PRIMARY KEY,
    channel_uuid VARCHAR(36) NOT NULL,
    user_uuid VARCHAR(36) NOT NULL,
    content TEXT NOT NULL,
    version SMALLINT NULL,
    sent_at DATETIME NOT NULL,
    deleted_at DATETIME NULL,
    updated_at DATETIME DEFAULT NULL,

    INDEX channel_idx (channel_uuid),
    INDEX created_idx (sent_at),

    CONSTRAINT FOREIGN KEY channel_fk (channel_uuid) REFERENCES channels(uuid) ON DELETE CASCADE,
    CONSTRAINT FOREIGN KEY user_fk (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
) CHARSET utf8, ENGINE InnoDB;

CREATE TABLE message_versions (
    message_uuid VARCHAR(36) NOT NULL,
    version SMALLINT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,

    INDEX message_idx (message_uuid),

    CONSTRAINT FOREIGN KEY message_fk (message_uuid) REFERENCES messages(uuid) ON DELETE CASCADE
) CHARSET utf8, ENGINE InnoDB;