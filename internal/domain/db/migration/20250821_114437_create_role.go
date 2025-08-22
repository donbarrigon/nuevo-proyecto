package migration

func RolesUp() {
	CreateCollection("roles")
	CreateUniqueIndex("roles", 1, "name")
}

func RolesDown() {
	DropCollection("roles")
}
