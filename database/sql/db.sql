drop DATABASE if EXISTS user_device_db;
create DATABASE  user_device_db;
use user_device_db;

create TABLE users(
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  name VARCHAR(255) not NULL,
  password VARCHAR(255) NOT NULL,
  address VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  created DATETIME DEFAULT CURRENT_TIMESTAMP

);

create index users_name_index on users(name);

create table roles(
  id int UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  name VARCHAR(255) UNIQUE
);

create table permissions(
  id int UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  name VARCHAR(255) UNIQUE
);

create table user_permission_relation_table(
  user_id int UNSIGNED NOT NULL,
  permission_id int UNSIGNED NOT NULL,
  CONSTRAINT FOREIGN KEY (user_id) REFERENCES users(id) on DELETE CASCADE,
  CONSTRAINT FOREIGN KEY (permission_id) REFERENCES permissions(id) on DELETE CASCADE,
  CONSTRAINT permission_unique UNIQUE (user_id, permission_id)
);

create INDEX user_permission_index on user_permission_relation_table(user_id, permission_id);

create TABLE user_roles_relation_table(
  user_id int UNSIGNED NOT NULL,
  role_id int UNSIGNED NOT NULL,
  CONSTRAINT FOREIGN KEY (user_id) REFERENCES users(id) on DELETE CASCADE,
  CONSTRAINT FOREIGN KEY (role_id) REFERENCES roles(id) on DELETE CASCADE,
  CONSTRAINT roles_unique UNIQUE (user_id, role_id)
);

CREATE INDEX user_roles_index on user_roles_relation_table(user_id, role_id);

create TABLE device(
  id int UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  user_id int UNSIGNED,
  FOREIGN KEY (user_id) REFERENCES users(id) on DELETE SET NULL,
  isReciever bool NOT NULL DEFAULT FALSE
);

create TABLE emergency_service_group(
  id int UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  name VARCHAR(255) NOT NULL,
  center_lat FLOAT,
  center_lon FLOAT
);

CREATE TABLE emergency_service_group_users_relation_table(
  group_id int UNSIGNED,
  users_id int UNSIGNED,
  FOREIGN KEY (group_id) REFERENCES emergency_service_group(id) on DELETE CASCADE,
  FOREIGN KEY (users_id) REFERENCES users(id) on DELETE CASCADE
);

create table DetectedEvent(
  id int UNSIGNED AUTO_INCREMENT PRIMARY KEY not NULL ,
  centerLat float not null,
  centerLon float not NULL
);