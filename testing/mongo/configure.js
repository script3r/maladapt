

conn = new Mongo('mongodb://127.0.0.1:27017');
username = 'mongo';
password = 'mongo';
database = 'maladapt';


// create the ufa db and ufa user
db = conn.getDB(database);
db.dropUser(username);
db.createUser({
  'user': username,
  'pwd': password,
  'roles':[
    {
      role: 'readWrite', db: 'maladapt'
    }
  ],
  'passwordDigestor': 'server'
});


