CREATE VIEW view_hardstores AS
SELECT tbl_hardstores.*, tbl_stocks.platform, tbl_stocks.model, tbl_stocks.status, tbl_hardwares.unit, tbl_hardwares.supplier, tbl_hardwares.type
FROM tbl_hardstores
LEFT JOIN tbl_stocks ON tbl_hardstores.id=tbl_stocks.id
LEFt JOIN tbl_hardwares ON tbl_stocks.platform=tbl_hardwares.platform
;