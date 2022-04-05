-- /home/krylon/go/src/github.com/blicero/recovery/database/avg.sql
-- created on 05. 04. 2022 by Benjamin Walkenhorst
-- (c) 2022 Benjamin Walkenhorst
-- Use at your own risk!

SELECT
        m.id,
        m.timestamp,
        CAST(ROUND(AVG(score) FILTER (
                                     WHERE timestamp BETWEEN m.timestamp - 259200
                                                         AND m.timestamp)
                   OVER (ORDER BY timestamp)) AS INTEGER)
           AS mavg,
        m.note
FROM mood m
ORDER BY timestamp;


