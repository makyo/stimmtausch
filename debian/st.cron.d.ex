#
# Regular cron jobs for the st package
#
0 4	* * *	root	[ -x /usr/bin/st_maintenance ] && /usr/bin/st_maintenance
