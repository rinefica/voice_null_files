CREATE TABLE files
(
    uuid    varchar not null unique,
    filename varchar,
    user_id integer,
    PRIMARY KEY (uuid),
    constraint fk_files_users
        foreign key (user_id)
            REFERENCES users (id)
);
