CREATE TABLE relationship
 (
	id SERIAL not null ,
	uid int not null,
	oid int not null,
	status smallint not null,
	type varchar(20),
	constraint pk_tb_relation_id primary key(id)
);

CREATE TABLE users (
  id SERIAL not null ,
  name VARCHAR(20),
  type varchar(20),
  constraint pk_tb_user_primary_id primary key(id)
);