mongod --config /etc/mongo/mongo.yaml
mongo --ssl --sslCAFile /etc/ssl/mongodb.pem < /etc/mongo/configure.js 
