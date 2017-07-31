CREATE TABLE 'tbl_softwares' (
    'guid'              VARCHAR(64) PRIMARY KEY,
    'model'             VARCHAR(64) NOT NULL UNIQUE,
    'name'              VARCHAR(64) NOT NULL,
    'dept'              VARCHAR(64) NULL,
    'description'       VARCHAR(64) NULL
);
