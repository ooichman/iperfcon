#!/bin/bash

# checking the variable and replacing them

cat > /home/apache/oidc-main.conf << EOF
<VirtualHost *:*>
	ServerAdmin webmaster@infra.local
	DocumentRoot /var/www/html
	ErrorLog /dev/stderr
        TransferLog /dev/stdout
ScriptAlias "/iperf" "/var/www/cgi-bin/"
<Directory "/var/www/cgi-bin/">
	require all granted
	Options +ExecCGI
        AddHandler cgi-script .cgi
</Directory>
</VirtualHost>
EOF

# starting the HTTPD
/usr/sbin/httpd -DFOREGROUND
