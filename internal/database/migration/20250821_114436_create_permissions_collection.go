package migration

func PermissionsUp() {
	CreateCollection("permissions", func(collection string) {
		CreateUniqueIndex(collection, 1, "name")
	})

}

func PermissionsDown() {
	DropCollection("permissions")
}
