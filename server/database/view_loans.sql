CREATE VIEW view_loans AS
SELECT tbl_loans.*,tbl_stocks.platform, tbl_stocks.model
FROM tbl_loans
LEFT JOIN tbl_stocks ON tbl_loans.id=tbl_stocks.id
;