CREATE TABLE 'tbl_stocks' (
    'guid'              VARCHAR(64) PRIMARY KEY,
    'id'                VARCHAR(64) NOT NULL UNIQUE,
    'platform'          VARCHAR(64) NOT NULL,
    'model'             VARCHAR(64) NOT NULL,
    'status'            VARCHAR(64) NOT NULL,
    'date'				VARCHAR(64) NOT NULL,
    'description'       VARCHAR(64) NULL
);
