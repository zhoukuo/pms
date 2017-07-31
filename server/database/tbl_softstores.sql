CREATE TABLE 'tbl_softstores' (
    'guid'                  VARCHAR(64) PRIMARY KEY,
    'model'                 VARCHAR(64) NOT NULL,
    'version'               VARCHAR(64) NOT NULL,
    'count'                 VARCHAR(64) NULL,
    'publish_date'          VARCHAR(64) NOT NULL,
    'store_date'            VARCHAR(64) NOT NULL,
    'class'                 VARCHAR(64) NOT NULL,
    'description'           VARCHAR(64) NULL
);
