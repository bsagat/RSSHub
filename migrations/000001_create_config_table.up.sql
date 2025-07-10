CREATE TABLE config(
    Run bool DEFAULT FALSE,
    Worker_count INT DEFAULT 3,
    Timer_interval INTERVAL DEFAULT '3 minutes'
);

INSERT INTO config DEFAULT VALUES;
