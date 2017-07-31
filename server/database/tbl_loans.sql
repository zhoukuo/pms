CREATE TABLE 'tbl_loans' (
    'guid'                  VARCHAR(64) PRIMARY KEY,
    'date'                  VARCHAR(64) NOT NULL,
    'id'                    VARCHAR(64) NOT NULL,
    'project_id'            VARCHAR(64) NOT NULL,
    'user'                  VARCHAR(64) NOT NULL,
    'delivery_type'         VARCHAR(64) NOT NULL,
    'tracking_number'       VARCHAR(64) NOT NULL,
    'description'           VARCHAR(64) NULL
);
