CREATE TABLE city_list (
    id int,
    name varchar(255) COLLATE utf8_general_ci,
    country varchar(10),
    lon double(9,6),
    lat double(9,6),    
    PRIMARY KEY (id)
) CHARACTER SET utf8 COLLATE utf8_bin;


CREATE TABLE country_city (
  countryName varchar(255) NOT NULL,
  officialName varchar(255) DEFAULT NULL,
  cityName varchar(255) DEFAULT NULL,
  PRIMARY KEY (countryName)
) CHARACTER SET utf8 COLLATE utf8_bin;
