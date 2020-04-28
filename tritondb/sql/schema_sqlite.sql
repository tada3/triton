
PRAGMA encoding = 'UTF-8';

CREATE TABLE city_list (
    id int,
    name varchar(255) COLLATE nocase,
    country varchar(10),
    lon double(9,6),
    lat double(9,6),    
    PRIMARY KEY (id)
);

CREATE INDEX idx_name ON city_list(name);

CREATE TABLE preferred_city (
    id INTEGER PRIMARY KEY,
    name varchar(255) COLLATE nocase,
    country varchar(10),
    priority int
);

CREATE TABLE translation (
    id INTEGER PRIMARY KEY,
    src varchar(255) COLLATE nocase,
    dst varchar(255) COLLATE nocase
);

CREATE INDEX idx_src ON translation(src);

CREATE TABLE country_city (
  countryName varchar(255) NOT NULL,
  officialName varchar(255) DEFAULT NULL,
  countryCode varchar(10) DEFAULT NULL,
  cityName varchar(255) DEFAULT NULL,
  isCountry int DEFAULT 0,
  PRIMARY KEY (countryName, cityName)
);

CREATE TABLE country_code (
  countryCode varchar(10) NOT NULL,
  officialName varchar(255) DEFAULT NULL,
  PRIMARY KEY (countryCode)
);

CREATE TABLE poi_city (
  id INTEGER PRIMARY KEY,
  name varchar(255) NOT NULL COLLATE nocase,
  countryCode varchar(10) DEFAULT NULL,
  name2 varchar(255) DEFAULT NULL COLLATE nocase,
  cityName varchar(255) NOT NULL COLLATE nocase,
  precedence int DEFAULT 100
);
