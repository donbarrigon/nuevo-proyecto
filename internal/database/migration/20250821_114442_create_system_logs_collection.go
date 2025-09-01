package migration

func SystemLogsUp() {
	CreateCollection("system_logs", func(collection string) {
		CreateIndex(collection, 1, "time")
	})
}

func SystemLogsDown() {
	DropCollection("system_logs")
}
