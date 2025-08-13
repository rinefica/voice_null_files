CREATE TABLE info_data
(
    uuid    varchar not null unique,
    data varchar,
    additional varchar,
    type varchar,
    user_id integer,
    PRIMARY KEY (uuid),
    constraint fk_info_data_users
        foreign key (user_id)
            REFERENCES users (id)
);