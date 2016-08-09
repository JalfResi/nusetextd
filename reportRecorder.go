package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// ReportRecorder stores an TextRazorResult in a MySQL table
type ReportRecorder struct {
	db *sql.DB
}

// NewReportRecorder is a ReportRecorder constructor
func NewReportRecorder(mysqlHost, mysqlUsername, mysqlPassword string) *ReportRecorder {
	return &ReportRecorder{}
}

// StoreTopics If there is an error executing any of the inserts, all pervious inserts
// for this TextRazorResult is reolledback, ensuring we dont have a partial
// TextRazorResult written to the database.
func (rr *ReportRecorder) StoreTopics(r *TextRazorResult) error {
	// Should be able to get the article url from the TextRazorResult
	// write this into MySQL linking table:
	//
	// | articleURLHash | articleUrl |
	//
	// | articleURLHash | TopicHash |
	//
	// | TopicHash | Label |
	//

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", "nusetext", "10rapid", "163.172.149.161", "nuseagent")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmtArticles, err := tx.Prepare("INSERT IGNORE INTO articles (hash, url) VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		return err
	}
	defer stmtArticles.Close()

	stmtTopics, err := tx.Prepare("INSERT IGNORE INTO topics (hash, label, score, wikiLink, wikidataId) VALUES( ?, ?, ?, ?, ? )") // ? = placeholder
	if err != nil {
		return err
	}
	defer stmtTopics.Close()

	stmtArticlesHasTopics, err := tx.Prepare("INSERT IGNORE INTO article_has_topics (articleHash, topicHash) VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		return err
	}
	defer stmtArticlesHasTopics.Close()

	for _, topic := range r.Response.Topics {

		articleURLHash := generateHash(r.URL)
		topicHash := generateHash(topic.Label)

		_, err = stmtArticles.Exec(articleURLHash, r.URL)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = stmtTopics.Exec(topicHash, topic.Label, topic.Score, topic.WikiLink, topic.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = stmtArticlesHasTopics.Exec(articleURLHash, topicHash)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}
