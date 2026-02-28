package webserver

// HandleError is a unified error handling utility for Go.
func HandleError(err error) {
    if err != nil {
        // Handle the error (e.g., log it, return a response, etc.)
        log.Fatalf("Error: %v", err)
    }
}