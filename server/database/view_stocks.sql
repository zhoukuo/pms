CREATE VIEW view_stocks AS
SELECT DISTINCT tbl_stocks.*, tbl_hardwares.class
FROM tbl_stocks
LEFT JOIN tbl_hardwares ON tbl_stocks.platform=tbl_hardwares.platform
;