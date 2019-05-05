CREATE TABLE city_list (
    id int,
    name varchar(255) COLLATE utf8_general_ci,
    country varchar(10),
    lon double(9,6),
    lat double(9,6),    
    PRIMARY KEY (id),
    KEY idx_name (name)
) CHARACTER SET utf8 COLLATE utf8_bin;

ALTER TABLE city_list ADD INDEX idx_name(name);

CREATE TABLE preferred_city (
    id int AUTO_INCREMENT,
    name varchar(255) COLLATE utf8_general_ci,
    country varchar(10),
    priority int,
    PRIMARY KEY (id)
) CHARACTER SET utf8 COLLATE utf8_bin;



CREATE TABLE translation (
    id int AUTO_INCREMENT,
    src varchar(255) COLLATE utf8_general_ci,
    dst varchar(255) COLLATE utf8_general_ci,
    PRIMARY KEY (id),
    KEY idx_src (src)
) CHARACTER SET utf8 COLLATE utf8_bin;




CREATE TABLE country_city (
  countryName varchar(255) NOT NULL,
  officialName varchar(255) DEFAULT NULL,
  countryCode varchar(10) DEFAULT NULL,
  cityName varchar(255) DEFAULT NULL,
  isCountry int DEFAULT 0,
  PRIMARY KEY (countryName)
) CHARACTER SET utf8 COLLATE utf8_bin;

CREATE TABLE country_code (
  countryCode varchar(10) NOT NULL,
  officialName varchar(255) DEFAULT NULL,
  PRIMARY KEY (countryCode)
) CHARACTER SET utf8 COLLATE utf8_bin;

CREATE TABLE poi_city (
  id int AUTO_INCREMENT,
  name varchar(255) NOT NULL COLLATE utf8_general_ci,
  countryCode varchar(10) DEFAULT NULL,
  name2 varchar(255) DEFAULT NULL COLLATE utf8_general_ci,
  cityName varchar(255) NOT NULL COLLATE utf8_general_ci,
  precedence int DEFAULT 100,
  PRIMARY KEY (id)
) CHARACTER SET utf8 COLLATE utf8_bin;

ALTER TABLE country_city ADD COLUMN countryCode varchar(10) DEFAULT NULL AFTER officialName;
ALTER TABLE country_city ADD COLUMN isCountry int DEFAULT 0 AFTER cityName;


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


INSERT INTO country_city VALUES ('岩手県','岩手',NULL,'盛岡');
INSERT INTO country_city VALUES ('福島県','福島',NULL,'福島');
INSERT INTO country_city VALUES ('埼玉県','埼玉',NULL,'さいたま');
INSERT INTO country_city VALUES ('神奈川県','神奈川',NULL,'横浜');
INSERT INTO country_city VALUES ('沖縄県','沖縄',NULL,'那覇');

INSERT INTO country_city VALUES ('サハリン','サハリン州',NULL,'ユジノ-サハリンスク');
INSERT INTO country_city VALUES ('バリ','バリ島','ID','デンパサール', 0);
INSERT INTO country_city VALUES ('サイパン','サイパン島','MP','サイパン', 0);
INSERT INTO country_city VALUES ('北マリアナ諸島','北マリアナ諸島自治連邦区','MP','サイパン', 1);
INSERT INTO country_city VALUES ('北極','北極','XX','北極', 0);
INSERT INTO country_city VALUES ('南極','南極','AQ','北極', 1);
INSERT INTO country_city VALUES ('セブ','セブ島','PH','セブ', 0);
INSERT INTO country_city VALUES ('ヨーロッパ','ヨーロッパ','XX','ヨーロッパ', 0);
INSERT INTO country_city VALUES ('タスマニア','タスマニア島','AU','ホバート', 0);
INSERT INTO country_city VALUES ('カムチャッカ','カムチャッカ半島','RU','ペトロパブロフスクカムチャツキー', 0);
INSERT INTO country_city VALUES ('色丹','色丹島','JP','色丹村', 0);
INSERT INTO country_city VALUES ('歯舞','歯舞群島','JP','歯舞', 0);






INSERT INTO poi_city (name, countryCode, cityName) VALUES ('エッフェル塔','FR', 'パリ');
INSERT INTO poi_city (name, countryCode, cityName, precedence) VALUES ('ナイアガラの滝','CA', 'ナイアガラフォールズ', 200);
INSERT INTO poi_city (name, countryCode, cityName, precedence) VALUES ('ナイアガラの滝','US', 'ナイアガラフォールズ', 100);


INSERT INTO translation (src, dst) VALUES ('ナイアガラフォールズ', 'Niagara Falls');
INSERT INTO translation (src, dst) VALUES ('セブ', 'Cebu City');
INSERT INTO translation (src, dst) VALUES ('セブ', 'Cebu City');
 select city.countryName, city.officialName, code.officialName from country_city AS city, country_code AS code WHERE city.countryName = code.officialName OR city.officialName = code.officialName;


select city.countryName, city.officialName, code.officialName from country_city AS city, country_code AS code WHERE code.officialName like concat('%', city.countryName, '%') AND city.countryName <> code.officialName AND city.officialName <> code.officialName AND city.countryName NOT IN ('アイルランド', 'インド', 'タイ', 'ギニア', 'サモア', 'マリ', 'リビア', 'フランス');


select distinct city.countryName from country_city AS city, country_code AS code WHERE code.officialName like concat('%', city.countryName, '%') AND city.countryName <> code.officialName AND city.officialName <> code.officialName AND city.countryName NOT IN ('アイルランド', 'インド', 'タイ', 'ギニア', 'サモア', 'マリ', 'リビア');


select * from country_city WHERE countryCode is null AND officialName not like '%州' AND officialName not like '%省' AND officialName not like '%市';


SELECT name, TRIM(TRAILING '-shi' FROM name), country FROM city_list WHERE country = 'JP' AND name LIKE '%-shi';

12/23
UPDATE city_list SET name = TRIM(TRAILING '-shi' FROM name) WHERE country = 'JP' AND name LIKE '%-shi';

SELECT name, COUNT(name) FROM city_list GROUP BY name HAVING COUNT(name) > 30;


INSERT INTO preferred_city VALUES (2113015, 'Chiba', 'JP', 100);
INSERT INTO preferred_city VALUES (5392171, 'San Jose', 'US', 100);
INSERT INTO preferred_city VALUES (5393052, 'Santa Cruz', 'US', 100);
INSERT INTO preferred_city VALUES (4781708, 'Richmond', 'US', 100);
INSERT INTO preferred_city VALUES (2643743, 'London', 'GB', 100);
INSERT INTO preferred_city VALUES (1701668, 'Manila', 'PH', 100);
INSERT INTO preferred_city VALUES (6174041, 'Victoria', 'CA', 100);
INSERT INTO preferred_city VALUES (241131, 'Victoria', 'SC', 90);
INSERT INTO preferred_city VALUES (6173331, 'Vancouver', 'CA', 100);
INSERT INTO preferred_city VALUES (1819729, 'Hong Kong S.A.R', 'HK', 100);
INSERT INTO preferred_city VALUES (1583992, 'Da nang', 'VN', 100);
INSERT INTO preferred_city VALUES (3872395, 'San Antonio', 'CL', 100);
INSERT INTO preferred_city VALUES (1832909, 'Young', 'KR', 100);
INSERT INTO preferred_city VALUES (4726206, 'San antonio', 'US', 110);
INSERT INTO preferred_city VALUES (3169070, 'Rome', 'IT', 100);
INSERT INTO preferred_city VALUES (5368361, 'ロサンゼルス', 'US', 100);
INSERT INTO preferred_city VALUES (5391811, 'サンディエゴ', 'US', 100);
INSERT INTO preferred_city VALUES (5391959, 'サンフランシスコ', 'US', 100);
INSERT INTO preferred_city VALUES (4140963, 'ワシントン', 'US', 100);
INSERT INTO preferred_city VALUES (3128760, 'バルセロナ', 'ES', 100);

INSERT INTO preferred_city VALUES (4887398, 'シカゴ', 'US', 100);



6174041 | Victoria | CA
