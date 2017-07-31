CREATE TABLE 'tbl_hardwares' (
    'guid'              VARCHAR(64) PRIMARY KEY,
    'platform'          VARCHAR(64) NOT NULL UNIQUE,
    'supplier'          VARCHAR(64) NOT NULL,
    'type'              VARCHAR(64) NOT NULL,
    'unit'              VARCHAR(64) NOT NULL,
    'price'             VARCHAR(64) NULL,
    'class'             VARCHAR(64) NOT NULL,
    'publish_date'      VARCHAR(64) NOT NULL,
    'description'       VARCHAR(64) NULL
);
