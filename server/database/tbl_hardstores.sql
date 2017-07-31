CREATE TABLE 'tbl_hardstores' (
    'guid'             VARCHAR(64) PRIMARY KEY,
    'id'               VARCHAR(64) NOT NULL UNIQUE,
    'order_id'         VARCHAR(64) NOT NULL,
    'date'             VARCHAR(64) NOT NULL,
    'price'            VARCHAR(64) NOT NULL,
    'sap'              VARCHAR(64) NOT NULL,
    'receiving'        VARCHAR(64) NOT NULL,
    'description'      VARCHAR(64) NULL
);
