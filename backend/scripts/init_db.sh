#!/bin/bash
DB_NAME=${MYSQL_DATABASE}_test
DB_USER=${MYSQL_USER}
DB_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}

echo "###";
echo "# CREATE TEST DB";
echo "###";
mysql -u root -p"${DB_ROOT_PASSWORD}" -e "CREATE DATABASE IF NOT EXISTS ${DB_NAME};"
mysql -u root -p"${DB_ROOT_PASSWORD}" -e "GRANT ALL PRIVILEGES ON ${DB_NAME}.* TO '${DB_USER}'@'%';"
mysql -u root -p"${DB_ROOT_PASSWORD}" -e "FLUSH PRIVILEGES;"