package main

// ReportRecorder stores an TextRazorResult in a MySQL table
type ReportRecorder interface {
	// StoreTopics If there is an error executing any of the inserts, all pervious inserts
	// for this TextRazorResult is reolledback, ensuring we dont have a partial
	// TextRazorResult written to the database.
	StoreTopics(r *TextRazorResult) error
}
