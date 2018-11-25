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
  countryCode varchar(10) DEFAULT NULL,
  cityName varchar(255) DEFAULT NULL,
  PRIMARY KEY (countryName)
) CHARACTER SET utf8 COLLATE utf8_bin;

CREATE TABLE country_code (
  countryCode varchar(10) NOT NULL,
  officialName varchar(255) DEFAULT NULL,
  PRIMARY KEY (countryCode)
) CHARACTER SET utf8 COLLATE utf8_bin;

ALTER TABLE country_city ADD COLUMN countryCode varchar(10) DEFAULT NULL AFTER officialName;

UPDATE country_city AS city, country_code AS code SET city.countryCode = code.countryCode WHERE city.officialName = code.officialName;
UPDATE country_city AS city, country_code AS code SET city.countryCode = code.countryCode WHERE city.countryName = code.officialName;


UPDATE country_city AS city, country_code AS code SET city.countryCode = code.countryCode WHERE code.officialName like concat('%', city.countryName, '%') AND city.countryName <> code.officialName AND city.officialName <> code.officialName AND city.countryName NOT IN ('アイルランド', 'インド', 'タイ', 'ギニア', 'サモア', 'マリ', 'リビア', 'フランス');

UPDATE country_city SET countryCode = 'TH' WHERE countryName = 'タイ';
UPDATE country_city SET countryCode = 'GN' WHERE countryName = 'ギニア';
UPDATE country_city SET countryCode = 'GW' WHERE countryName = 'ギニアビサウ';
UPDATE country_city SET countryCode = 'WS' WHERE countryName = 'サモア';
UPDATE country_city SET countryCode = 'ML' WHERE countryName = 'マリ';
UPDATE country_city SET countryCode = 'LY' WHERE countryName = 'リビア';
UPDATE country_city SET countryCode = 'FR' WHERE countryName = 'フランス';

UPDATE country_city SET countryCode = 'GB' WHERE countryName = 'イギリス';
UPDATE country_city SET countryCode = 'KH' WHERE countryName = 'カンボジア';
UPDATE country_city SET countryCode = 'GN' WHERE countryName = 'ギニア';
UPDATE country_city SET countryCode = 'GR' WHERE countryName = 'ギリシャ';
UPDATE country_city SET countryCode = 'MV' WHERE countryName = 'モルディブ';
UPDATE country_city SET countryCode = 'RS' WHERE countryName = 'セルビア';
UPDATE country_city SET countryCode = 'SC' WHERE countryName = 'セーシェル';

INSERT INTO country_city VALUES ('グアム','グアム','GU','ハガニア')
INSERT INTO country_city VALUES ('ボルネオ島','ボルネオ島',NULL,'クチン')
INSERT INTO country_city VALUES ('セントーサ島','セントーサ島',NULL,'シンガポール');
INSERT INTO country_city VALUES ('北海道','北海道',NULL,'札幌');
INSERT INTO country_city VALUES ('アフリカ','アフリカ',NULL,'ナイロビ');
INSERT INTO country_city VALUES ('ワイキキ','ワイキキ',NULL,'ホノルル');





 select city.countryName, city.officialName, code.officialName from country_city AS city, country_code AS code WHERE city.countryName = code.officialName OR city.officialName = code.officialName;


select city.countryName, city.officialName, code.officialName from country_city AS city, country_code AS code WHERE code.officialName like concat('%', city.countryName, '%') AND city.countryName <> code.officialName AND city.officialName <> code.officialName AND city.countryName NOT IN ('アイルランド', 'インド', 'タイ', 'ギニア', 'サモア', 'マリ', 'リビア', 'フランス');


select distinct city.countryName from country_city AS city, country_code AS code WHERE code.officialName like concat('%', city.countryName, '%') AND city.countryName <> code.officialName AND city.officialName <> code.officialName AND city.countryName NOT IN ('アイルランド', 'インド', 'タイ', 'ギニア', 'サモア', 'マリ', 'リビア');


select * from country_city WHERE countryCode is null AND officialName not like '%州' AND officialName not like '%省' AND officialName not like '%市';