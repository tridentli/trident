# Include this in your server {} block
#
# Example configuration (drop in /etc/conf/
# 8<------------------
# 
# server {
#	listen 443;
#	listen [::]:443;
#	server_name trident.example.net;
#	access_log /var/log/nginx/trident-access.log
#	include
# ------------------>8
#
#


	# Our webroot (contains static, non-sensitive files, source if public ;)
	root /usr/share/trident/webroot/;

	######################################################
	# Static files
	######################################################
	location /css/epiceditor/themes/ {
		alias /usr/share/epiceditor/themes/;
	}

	location /css/ {
		try_files /usr/share/trident/webroot$uri /usr/share/pitchfork/webroot$uri =404;
	}

	location /favicon.ico {
		try_files /usr/share/trident/webroot$uri /usr/share/pitchfork/webroot$uri =404;
	}

	location /gfx/ {
		try_files /usr/share/trident/webroot$uri /usr/share/pitchfork/webroot$uri =404;
	}

	location /js/epiceditor.min.js {
		alias /usr/share/epiceditor/js/epiceditor.min.js;
	}

	location /js/ {
		try_files /usr/share/trident/webroot$uri /usr/share/pitchfork/webroot$uri =404;
	}

	######################################################
	# Forward all requests to the Trident Daemon
	######################################################
	location / {
		client_max_body_size    0;
		proxy_set_header	Host $host;
		proxy_http_version	1.1;
		proxy_pass		http://trident-daemon;
	}

