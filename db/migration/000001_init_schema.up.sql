CREATE TABLE Profile (
    id SERIAL PRIMARY KEY,
    photo VARCHAR(512) NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_deleted BOOLEAN NOT NULL,
    token VARCHAR(512) NOT NULL DEFAULT '',
    created_at INT NOT NULL,
    updated_at INT NOT NULL
);

CREATE TABLE Storage(
    id SERIAL PRIMARY KEY,
    supervisor_id INT NOT NULL REFERENCES Profile(id),
    name VARCHAR(255) NOT NULL,
    description VARCHAR(512) NULL,
    is_deleted BOOLEAN NOT NULL,
    created_at INT NOT NULL,
    updated_at INT NOT NULL,
    UNIQUE(supervisor_id, name)
);

CREATE TABLE Base_Paper(
    id SERIAL PRIMARY KEY,
    storage_id INT NOT NULL REFERENCES Storage(id),
    gsm INT NOT NULL,
    width INT NOT NULL,
    io INT NOT NULL,
    material_number INT NOT NULL,
    quantity INT NOT NULL,
    location VARCHAR(10) NOT NULL DEFAULT '',
    is_deleted BOOLEAN NOT NULL,
    created_at INT NOT NULL,
    updated_at INT NOT NULL,
    UNIQUE(storage_id, gsm, width, io, material_number, location)
);

CREATE TABLE Storage_Member(
    id SERIAL PRIMARY KEY,
    storage_id INT NOT NULL REFERENCES Storage(id),
    member_id INT NOT NULL REFERENCES Profile(id),
    is_admin BOOLEAN NOT NULL,
    is_active BOOLEAN NOT NULL,
    is_deleted BOOLEAN NOT NULL,
    created_at INT NOT NULL,
    updated_at INT NOT NULL,
    UNIQUE (storage_id, member_id)
);

CREATE TYPE History_Status AS ENUM ('stored', 'moved','deleted', 'delivered');

CREATE TABLE History (
    id SERIAL PRIMARY KEY,
    base_paper_id INT NOT NULL REFERENCES Base_Paper(id),
    storage_id INT NOT NULL,
    member_id INT NOT NULL,
    affected INT NOT NULL,
    status History_Status NOT NULL,
    created_at INT NOT NULL,
    CONSTRAINT fk_history FOREIGN KEY (storage_id, member_id) REFERENCES Storage_Member (storage_id, member_id)
)
