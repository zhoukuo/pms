#!/bin/bash

# Program to generate database file

timestamp=`date +"%Y/%m/%d %H:%M:%S"`
file="../pms.db"

if [ -f "$file" ]; then
	echo $timestamp 'remove old db file ...'
	rm $file
fi

echo $timestamp 'create new db file ...'
sqlite3 $file ".read tbl_user.sql" ".read tbl_hardstores.sql" ".read tbl_softstores.sql" ".read tbl_loans.sql" ".read tbl_returns.sql" ".read tbl_harddeliverys.sql" ".read tbl_softdeliverys.sql" ".read tbl_stocks.sql" ".read tbl_softwares.sql" ".read tbl_hardwares.sql" ".read tbl_projects.sql" ".read tbl_events.sql" ".read view_stocks.sql" ".read view_loans.sql" ".read view_returns.sql" ".read view_hardstores.sql" ".read view_harddeliverys.sql" ".read view_stockscheck.sql"

echo $timestamp 'pms.db created.'

