package migration

func VerificationCodesUp() {
	CreateCollection("verification_codes", func(collection string) {
		CreateUniqueIndex(collection, 1, "user_id", "type", "code")
	})
}

func VerificationCodesDown() {
	DropCollection("verification_codes")
}
