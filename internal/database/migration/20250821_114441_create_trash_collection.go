package migration

func TrashUp() {
	CreateCollection("trash", func(collection string) {
		CreateIndex(collection, 1, "document._id")
		CreateIndex(collection, 1, "collection")
		CreateIndex(collection, -1, "deleted_at")
	})
}

func TrashDown() {
	DropCollection("trash")
}
