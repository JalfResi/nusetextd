package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Stores an TextRazorResult in a MySQL table

type ReportRecorder struct {
	db *sql.DB
}

func NewReportRecorder(mysqlHost, mysqlUsername, mysqlPassword string) *ReportRecorder {
	return &ReportRecorder{}
}

// If there is an error executing any of the inserts, all pervious inserts
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

	stmtTopics, err := tx.Prepare("INSERT IGNORE INTO topics (hash, label) VALUES( ?, ? )") // ? = placeholder
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

		if topic.Score == 1.0 {

			articleURLHash := generateHash(r.URL)
			topicHash := generateHash(topic.Label)

			_, err = stmtArticles.Exec(articleURLHash, r.URL)
			if err != nil {
				tx.Rollback()
				return err
			}

			_, err = stmtTopics.Exec(topicHash, topic)
			if err != nil {
				tx.Rollback()
				return err
			}

			_, err = stmtArticlesHasTopics.Exec(articleURLHash, topicHash)
			if err != nil {
				tx.Rollback()
				return err
			}

			fmt.Println(topic)
		}
	}

	tx.Commit()
	return nil
}
