package migration

func AccessTokensUp() {
	CreateCollection("access_tokens", func(collection string) {
		CreateIndex(collection, 1, "user_id")
		CreateUniqueIndex(collection, 1, "token")
		CreateIndex(collection, 1, "permissions")
	})
}

func AccessTokensDown() {
	DropCollection("access_tokens")
}
