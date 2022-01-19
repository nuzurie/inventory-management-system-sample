CREATE TABLE IF NOT EXISTS item (
    id text PRIMARY KEY,
    name text NOT NULL,
    description text,
    price float,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

CREATE TABLE IF NOT EXISTS inventory (
     id text PRIMARY KEY,
     quantity int,
     updated_at timestamp without time zone,
     item_id text,
     FOREIGN KEY (item_id)
         REFERENCES item(id)
);

INSERT INTO item (id, name, description, price, created_at, updated_at)
VALUES ('abcdef', 'creative name 1', 'some keywords to search for', 1.99, NOW(), now());

INSERT INTO item (id, name, description, price, created_at, updated_at)
VALUES ('jdasoidjp', 'additional item', 'these items were added with sql script from docker file', 12.99, NOW(), NOW());

INSERT INTO item (id, name, description, price, created_at, updated_at)
VALUES ('apodad', 'creme brulee', 'dockerfile in backend repo db folder', 99.99, NOW(), NOW());

INSERT INTO item (id, name, description, price, created_at, updated_at)
VALUES ('oaifoad', 'expensive', 'some keywords to search for', 0.99, NOW(), NOW());

INSERT INTO inventory (id, quantity, updated_at, item_id)
VALUES ('asdad', 30, now(), 'apodad');

INSERT INTO inventory (id, quantity, updated_at, item_id)
VALUES ('asawqe', 302, now(), 'jdasoidjp');


