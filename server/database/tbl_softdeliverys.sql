CREATE TABLE 'tbl_softdeliverys' (
	'guid'					VARCHAR(64) PRIMARY KEY,
    'date'                  VARCHAR(64) NOT NULL,
    'model'                 VARCHAR(64) NOT NULL,
    'count'                 VARCHAR(64) NOT NULL,
    'version'               VARCHAR(64) NOT NULL,
    'project_id'            VARCHAR(64) NOT NULL,
    'delivery_type'         VARCHAR(64) NOT NULL,
    'tracking_number'       VARCHAR(64) NULL,
    'description'           VARCHAR(64) NULL
);
