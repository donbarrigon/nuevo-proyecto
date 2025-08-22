package migration

func PermissionsUp() {
	CreateCollection("permissions")
	CreateUniqueIndex("permissions", 1, "name")
}

func PermissionsDown() {
	DropCollection("permissions")
}
