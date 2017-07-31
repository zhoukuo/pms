
//为导出功能设置数据源
function setExportDataAttr(fileName) {
    $(".export-csv").attr("data-table","#table");
    $(".export-tsv").attr("data-table","#table");
    $(".export-pdf").attr("data-table","#table");
    $(".export-png").attr("data-table","#table");
    $(".export-excel").attr("data-table","#table");
    $(".export-xlsx").attr("data-table","#table");
    $(".export-doc").attr("data-table","#table");
    $(".export-powerpoint").attr("data-table","#table");
    $(".export-txt").attr("data-table","#table");
    $(".export-xml").attr("data-table","#table");
    $(".export-sql").attr("data-table","#table");
    $(".export-json").attr("data-table","#table");

    $(".export-csv").attr("data-filename",fileName);
    $(".export-tsv").attr("data-filename",fileName);
    $(".export-pdf").attr("data-filename",fileName);
    $(".export-png").attr("data-filename",fileName);
    $(".export-excel").attr("data-filename",fileName);
    $(".export-xlsx").attr("data-filename",fileName);
    $(".export-doc").attr("data-filename",fileName);
    $(".export-powerpoint").attr("data-filename",fileName);
    $(".export-txt").attr("data-filename",fileName);
    $(".export-xml").attr("data-filename",fileName);
    $(".export-sql").attr("data-filename",fileName);
    $(".export-json").attr("data-filename",fileName);
};


