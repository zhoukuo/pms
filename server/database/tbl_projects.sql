CREATE TABLE 'tbl_projects' (
    'guid'              VARCHAR(64) PRIMARY KEY,
    'id'                VARCHAR(64) NOT NULL UNIQUE,
    'name'              VARCHAR(64) NOT NULL,
    'manager'           VARCHAR(64) NULL,
    'dept'              VARCHAR(64) NULL,
    'customer'          VARCHAR(64) NULL,
    'description'       VARCHAR(64) NULL
);
