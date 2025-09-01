package migration

func ActivitiesUp() {
	CreateCollection("histories", func(collection string) {
		CreateIndex(collection, 1, "user_id")
		CreateIndex(collection, 1, "document_id")
		CreateIndex(collection, 1, "collection")
		CreateIndex(collection, -1, "occurred_at")
	})
}

func ActivitiesDown() {
	DropCollection("activities")
}
