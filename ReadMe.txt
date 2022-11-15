X sessions

X srv/ws

X stat

X socket

X tokenBlock

mapping

app

events

****************************************
Проблема невозможности чтения с реплики после события insert/update
master
select write_lsn from pg_stat_replication;
select pg_current_wal_lsn();

replica
select received_lsn from pg_stat_wal_receiver;

select pg_wal_lsn_diff((select received_lsn from pg_stat_wal_receiver), '53E/1FF74908')
must be >0 on slave!
1) On Event add to params 'lsn': (SELECT pg_current_wal_lsn())
2) On slave:
javascript must add lsn parameter to select query (get_list)
	- if it is read query, from slave
	- got param lsn 
	- SELECT pg_wal_lsn_diff( (SELECT received_lsn FROM pg_stat_wal_receiver), LSN_PARAM) >= 0 - slave fits,
			otherwise it is behind, does not yet have server updates
