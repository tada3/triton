# Triton
A Clova skill that tells you weather information around the world.


SELECT IFNULL(countryCode, ''),cityName from country_city WHERE countryName = ? OR officialName = ? ORDER BY RAND() LIMIT 1