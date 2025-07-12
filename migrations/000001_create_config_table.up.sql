CREATE TABLE config(
    run bool DEFAULT FALSE,
    worker_count INT DEFAULT 3,
    timer_interval INTERVAL DEFAULT '3 minutes'
);

INSERT INTO config DEFAULT VALUES;
