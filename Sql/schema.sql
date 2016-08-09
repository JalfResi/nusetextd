CREATE TABLE IF NOT EXISTS articles (
    hash BINARY(16) NOT NULL,
    url TEXT,
    PRIMARY KEY (hash)
);

CREATE TABLE IF NOT EXISTS topics (
    hash BINARY(16) NOT NULL,
    label TINYTEXT,
    score DOUBLE,
    wikiLink TEXT,
    wikidataId INT
    PRIMARY KEY (hash)    
);

CREATE TABLE IF NOT EXISTS article_has_topics (
    articleHash BINARY(16) NOT NULL,
    topicHash BINARY(16) NOT NULL,
    PRIMARY KEY (articleHash, topicHash),
    FOREIGN KEY (articleHash) REFERENCES articles(hash),
    FOREIGN KEY (topicHash) REFERENCES topics(hash)    
);

