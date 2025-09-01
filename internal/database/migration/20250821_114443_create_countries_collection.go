package migration

func CountriesUp() {
	CreateCollection("countries", func(collection string) {
		CreateUniqueIndex(collection, 1, "name")
		CreateTextIndex(collection, "name")
	})

}

func CountriesDown() {
	DropCollection("countries")
}
