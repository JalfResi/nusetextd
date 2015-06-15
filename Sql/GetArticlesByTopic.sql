-- Get Articles by Topic
SELECT 
	a.url AS Article, 
	GROUP_CONCAT(t.label) AS Topics,
	a.createdDate
FROM articles a
left join article_has_topics aht on a.hash = aht.articleHash
left join topics t on aht.topicHash = t.hash
WHERE t.label LIKE '%Television%'
GROUP BY a.hash
;
