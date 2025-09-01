package migration

func StatesUp() {
	CreateCollection("states", func(collection string) {
		CreateIndex(collection, 1, "country_id")
		CreateTextIndex(collection, "name", "country_name")
	})

}

func StatesDown() {
	DropCollection("states")
}
