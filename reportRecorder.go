package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Stores an AnalysisReport in a MySQL table

type ReportRecorderer interface {
	StoreReport(r *TopicsReport)
}

type ReportRecorder struct {
	db *sql.DB
}

func NewReportRecorder(mysqlHost, mysqlUsername, mysqlPassword string) *ReportRecorder {
	return &ReportRecorder{}
}

// If there is an error executing any of the inserts, all pervious inserts
// for this TopicsReport is reolledback, ensuring we dont have a partial
// topicReport written to the database.
func (rr *ReportRecorder) StoreTopicsReport(r *TopicsReport) error {
	// Should be able to get the article url from the TopicsReport
	// write this into MySQL linking table:
	//
	// | articleUrlHash | articleUrl |
	//
	// | articleUrlHash | TopicHash |
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

	for _, topic := range r.topics {

		articleUrlHash := generateHash(r.url)
		topicHash := generateHash(topic)

		_, err = stmtArticles.Exec(articleUrlHash, r.url)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = stmtTopics.Exec(topicHash, topic)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = stmtArticlesHasTopics.Exec(articleUrlHash, topicHash)
		if err != nil {
			tx.Rollback()
			return err
		}

		fmt.Println(topic)
	}

	tx.Commit()
	return nil
}
