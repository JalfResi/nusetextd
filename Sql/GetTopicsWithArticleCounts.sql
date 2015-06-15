-- List Topics with Article counts
SELECT 
	t.label AS Topic, 
	COUNT(*) AS Articles
FROM topics t
LEFT JOIN article_has_topics aht ON aht.topicHash = t.hash
GROUP BY t.hash
ORDER BY t.label;