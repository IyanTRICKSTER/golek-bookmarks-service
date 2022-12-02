db.createUser({
  user: "ayocodedb",
  pwd: "secret",
  roles: [{ role: "readWrite", db: "golek_tagging" }],
  mechanisms: ["SCRAM-SHA-1"],
});
