CREATE VIEW view_harddeliverys AS 
SELECT
tbl_harddeliverys.*,
tbl_stocks.platform, tbl_stocks.model, tbl_stocks.status, tbl_projects.name, tbl_projects.manager, tbl_projects.dept, tbl_projects.customer
FROM tbl_harddeliverys 
LEFT JOIN tbl_stocks ON 
tbl_harddeliverys.id = tbl_stocks.id
LEFT JOIN tbl_projects ON 
tbl_harddeliverys.project_id = tbl_projects.id
;
