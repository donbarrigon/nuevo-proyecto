package migration

func ActivitiesUp() {
	CreateCollection("activities", func(collection string) {
		CreateIndex(collection, 1, "user_id")
		CreateIndex(collection, 1, "document_id")
		CreateIndex(collection, 1, "collection")
		CreateIndex(collection, -1, "created_at")
	})
}

func ActivitiesDown() {
	DropCollection("activities")
}
