CREATE TABLE 'tbl_harddeliverys' (
        'guid'                 VARCHAR(64) PRIMARY KEY,
        'order_id'             VARCHAR(64) NOT NULL,
        'date'                 VARCHAR(64) NOT NULL,
        'id'                   VARCHAR(64) NOT NULL,
        'price'                VARCHAR(64) NULL,
        'project_id'           VARCHAR(64) NOT NULL,
        'soft_model'           VARCHAR(64) NOT NULL,
        'version'              VARCHAR(64) NOT NULL,
        'algorithm'            VARCHAR(64) NULL,
        'period'               VARCHAR(64) NULL,
        'user'                 VARCHAR(64) NULL,
        'delivery_type'        VARCHAR(64) NULL,
        'tracking_number'      VARCHAR(64) NULL,
        'description'          VARCHAR(64) NULL
);
