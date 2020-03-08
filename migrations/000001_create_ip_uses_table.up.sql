CREATE TABLE ip_uses
(
    ip       VARCHAR(39) NOT NULL,
    last_use timestamp NOT NULL,
    PRIMARY KEY (ip)
);