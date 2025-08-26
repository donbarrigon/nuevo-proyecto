package migration

func RolesUp() {
	CreateCollection("roles", func(collection string) {
		CreateUniqueIndex(collection, 1, "name")
	})

}

func RolesDown() {
	DropCollection("roles")
}
