package migration

func TrashUp() {
	CreateCollection("trash", func(collection string) {
		CreateIndex(collection, 1, "user_id")
		CreateIndex(collection, 1, "collection")
		CreateIndex(collection, -1, "deleted_at")
	})
}

func TrashDown() {
	DropCollection("trash")
}
