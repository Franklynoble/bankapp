#!/bin/sh
#this is to enable golang migrate. since we use  the alpine verion, the bash shell is not availabe
#set -e  the script would return immediatle if the command  returns a non zero status
set -e

#firstly, run db migration; pass in the path to all the db sequel files, the database url using  the emv variable
echo "run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up 

echo "start the app"
#takes all paremter pass to the script and  run it
exec "$@"
