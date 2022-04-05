-- /home/krylon/go/src/github.com/blicero/recovery/database/avg.sql
-- created on 05. 04. 2022 by Benjamin Walkenhorst
-- (c) 2022 Benjamin Walkenhorst
-- Use at your own risk!

SELECT
        m.id,
        m.timestamp,
        m.score,
        (SELECT avg(m1.score) AS ravg
         FROM mood m1
         WHERE m1.timestamp BETWEEN m.timestamp - 259200 AND m.timestamp)
               AS ravg,
        m.note
FROM mood m
ORDER BY timestamp;


