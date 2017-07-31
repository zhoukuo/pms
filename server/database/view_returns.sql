CREATE VIEW view_returns AS
SELECT tbl_returns.*,tbl_stocks.platform, tbl_stocks.model
FROM tbl_returns
LEFT JOIN tbl_stocks ON tbl_returns.id=tbl_stocks.id
;