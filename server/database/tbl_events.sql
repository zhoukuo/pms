CREATE TABLE 'tbl_events' (
    'guid'              VARCHAR(64) PRIMARY KEY,
    'id'                VARCHAR(64) NOT NULL,
    'time'              VARCHAR(64) NOT NULL,
    'business'          VARCHAR(64) NOT NULL,
    'action'            VARCHAR(64) NOT NULL,
    'operator'          VARCHAR(64) NOT NULL,
    'description'       VARCHAR(64) NULL
);
