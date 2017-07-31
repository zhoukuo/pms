CREATE VIEW view_stockscheck AS
SELECT DISTINCT tbl_hardstores.id, tbl_hardstores.order_id as store_order_id, tbl_hardstores.date as store_date, 
tbl_stocks.platform, tbl_stocks.status, tbl_stocks.date as lastchanged_date,
tbl_harddeliverys.order_id as delivery_order_id, tbl_harddeliverys.date as delivery_date
FROM tbl_hardstores
INNER JOIN tbl_stocks ON tbl_hardstores.id=tbl_stocks.id
LEFT JOIN tbl_harddeliverys ON tbl_hardstores.id=tbl_harddeliverys.id
;
