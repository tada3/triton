create table city_list (
    id int,
    name varchar(255) COLLATE utf8_general_ci,
    country varchar(10),
    lon double(9,6),
    lat double(9,6),    
    PRIMARY KEY (id)
) CHARACTER SET utf8 COLLATE utf8_bin;


create table city_list2 (
    id int,
    name varchar(255) COLLATE utf8_general_ci,
    country varchar(10),
    lon double,
    lat double,
    PRIMARY KEY (id)
) CHARACTER SET utf8 COLLATE utf8_bin;
