
-- ///////////////////////////////////////////////////////

CREATE TABLE IF NOT EXISTS articles (
    hash BINARY(16) NOT NULL,
    url TEXT,
    PRIMARY KEY (hash)
);

DELIMITER $$

CREATE TRIGGER article_generate_hash  
BEFORE INSERT ON articles
FOR EACH ROW  
BEGIN  
  NEW.hash = MD5(NEW.url)
END $$

DELIMITER ;

-- ///////////////////////////////////////////////////////

CREATE TABLE IF NOT EXISTS topics (
    hash BINARY(16) NOT NULL,
    label TINYTEXT,
    score DOUBLE,
    wikiLink TEXT,
    wikidataId INT
    PRIMARY KEY (hash)    
);

DELIMITER $$

CREATE TRIGGER topic_generate_hash  
BEFORE INSERT ON topics
FOR EACH ROW  
BEGIN  
  NEW.hash = MD5(NEW.label)
END $$

DELIMITER ;

-- ///////////////////////////////////////////////////////

CREATE TABLE IF NOT EXISTS article_has_topics (
    articleHash BINARY(16) NOT NULL,
    topicHash BINARY(16) NOT NULL,
    PRIMARY KEY (articleHash, topicHash),
    FOREIGN KEY (articleHash) REFERENCES articles(hash),
    FOREIGN KEY (topicHash) REFERENCES topics(hash)    
);

-- ///////////////////////////////////////////////////////

CREATE TABLE IF NOT EXISTS entities (
    hash            BINARY(16) NOT NULL,
	entityId        TEXT,
	entityEnglishId TEXT,
	confidenceScore DOUBLE,
	`type`          TEXT,
	freebaseTypes   TEXT,
	freebaseId      TEXT,
	matchingTokens  TEXT,
	matchedText     TEXT,
	`data`          TEXT,
	relevanceScore  DOUBLE,
	wikiLink        TEXT
    PRIMARY KEY (entityHash)
);

DELIMITER $$

CREATE TRIGGER entity_generate_hash  
BEFORE INSERT ON entities
FOR EACH ROW  
BEGIN  
  NEW.hash = MD5(NEW.matchedText)
END $$

DELIMITER ;

-- ///////////////////////////////////////////////////////

CREATE TABLE IF NOT EXISTS article_has_entities (
    articleHash BINARY(16) NOT NULL,
    entityHash BINARY(16) NOT NULL,
    PRIMARY KEY (articleHash, entityHash),
    FOREIGN KEY (articleHash) REFERENCES articles(hash),
    FOREIGN KEY (entityHash) REFERENCES entities(hash)    
);
