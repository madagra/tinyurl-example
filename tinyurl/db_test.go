package main

import "testing"

func initTestDb(t *testing.T) LocalDbClient {
	var dbClient = GetLocalDbClient()

	// this function is registered to be called at the
	// of every test which uses this function via the
	// Cleanup() interface
	t.Cleanup(func() {
		for key := range dbClient.urlKeysDB {
			delete(dbClient.urlKeysDB, key)
		}
	})

	return *dbClient
}
