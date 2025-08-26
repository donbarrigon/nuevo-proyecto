package migration

func CitiesUp() {
	CreateCollection("cities", func(collection string) {
		CreateIndex(collection, 1, "state_id")
		CreateIndex(collection, 1, "country_id")
		CreateTextIndex(collection, "name", "state_name", "country_name")
	})
}

func CitiesDown() {
	DropCollection("cities")
}
