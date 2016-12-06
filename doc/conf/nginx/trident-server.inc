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
#
#	ssl_certificate ...
#	...
#
#	include /etc/trident/nginx/trident-server.inc
# }
# ------------------>8

	# Our webroot (contains static, non-sensitive files, source if public ;)
	root /usr/share/trident/webroot/;

	######################################################
	# Static files
	######################################################
	location ~ ^/(css|gfx|js)/ {
		expires 7d;
		root /usr/share/;
		try_files /trident/webroot$uri /pitchfork/webroot$uri =404;
	}

	location /favicon.ico {
		expires 31d;
		root /usr/share/;
		try_files /trident/webroot$uri /pitchfork/webroot$uri =404;
	}

	######################################################
	# Forward all requests to the Trident Daemon
	######################################################
	location / {
		client_max_body_size    0;
		proxy_set_header	Host $host;
		proxy_http_version	1.1;
		proxy_set_header	X-Forwarded-For $proxy_add_x_forwarded_for;

		# No Connection header so that Keep-alive is performed where possible
		proxy_set_header	Connection "";

		# Pass it to Trident
		proxy_pass		http://trident-daemon;
	}

