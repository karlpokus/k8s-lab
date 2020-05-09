conn = new Mongo();
db = conn.getDB("blog");

// add user
if (db.users.count() == 0) {
	var user = {name:"bob", pwd:"dylan"};
	insert(user, "users")
}

// add some pets
if (db.posts.count() == 0) {
	var posts = [
		{title:"title1", body:"lorem ipsum 1"},
		{title:"title2", body:"lorem ipsum 2"}
	];
	insert(posts, "posts")
}

function insert(payload, collection) {
	db[collection].insert(payload, {}, function(){
		print("mongo-seed: " + collection + " updated")
	})
}
