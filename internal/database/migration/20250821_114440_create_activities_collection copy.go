package migration

func ActivitiesUp() {
	CreateCollection("activities", func(collection string) {
		CreateIndex(collection, 1, "user_id")
		CreateIndex(collection, 1, "document_id")
		CreateIndex(collection, 1, "collection")
	})
}

func ActivitiesDown() {
	DropCollection("activities")
}
