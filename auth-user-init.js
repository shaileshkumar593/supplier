// mongo-init.js
db.createUser({
  user: "mongosupplier",
  pwd: "mongo",
  roles: [
    { role: "readWrite", db: "mongo_swallow_supplier" }
  ]
});