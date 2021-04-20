package db

// Initialise initialises all database related operations
func Initialise() {
	// Connect to the database
	connect()

	// Run automatic migrations for all schemas
	migrate()
}
