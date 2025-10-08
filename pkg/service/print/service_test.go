package print

func ExamplePrintVersion() {
	PrintVersion("11", "2025-10-10", "#234")
	PrintVersion("", "", "")

	// Output:
	// Build version: "11"
	// Build date: "2025-10-10"
	// Build commit: "#234"
	// Build version: "N/A"
	// Build date: "N/A"
	// Build commit: "N/A"
}
