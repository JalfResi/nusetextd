package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
	"gopkg.in/yaml.v2"
)

// extractor constants
const (
	EXTRACTOR_ENTITIES        string = "entities"
	EXTRACTOR_TOPICS          string = "topics"
	EXTRACTOR_WORDS           string = "words"
	EXTRACTOR_PHRASES         string = "phrases"
	EXTRACTOR_DEPENDENCYTREES string = "dependency-trees"
	EXTRACTOR_RELATIONS       string = "relations"
	EXTRACTOR_ENTAILMENTS     string = "entailments"
	EXTRACTOR_SENSES          string = "senses"
)

// cleanup mode constants
const (
	MODE_RAW       string = "raw"
	MODE_STRIPTAGS string = "stripTags"
	MODE_CLEANHTML string = "cleanHTML"
)

var (
	httpBadRequest      = errors.New("Bad Request")
	httpUnauthorized    = errors.New("Unauthorized")
	httpRequestTooLarge = errors.New("Request Too Large")
)

type TextRazorRequest struct {
	Text                 string `form:"text,omitempty"                url:"text,omitempty"                yaml:"text,omitempty"`
	Url                  string `form:"url,omitempty"                 url:"url,omitempty"                 yaml:"url,omitempty"`
	ApiKey               string `form:"apiKey"                        url:"apiKey"                        yaml:"apiKey,omitempty"`     // required field
	Extractors           string `form:"extractors,omitempty"          url:"extractors,omitempty"          yaml:"extractors,omitempty"` // use extractor constants
	Rules                string `form:"rules,omitempty"               url:"rules,omitempty"               yaml:"rules,omitempty"`
	CleanupMode          string `form:"cleanup.mode,omitempty"        url:"cleanup.mode,omitempty"        yaml:"cleanup.mode,omitempty"` // cleanup mode constants
	CleanupReturnCleaned bool   `form:"cleanup.returnCleaned"         url:"cleanup.returnCleaned"         yaml:"cleanup.returnCleaned,omitempty"`
	CleanupReturnRaw     bool   `form:"cleanup.returnRaw"             url:"cleanup.returnRaw"             yaml:"cleanup.returnRaw,omitempty"`
	CleanupUseMetadata   bool   `form:"cleanup.useMetadata"           url:"cleanup.useMetadata"           yaml:"cleanup.useMetadata,omitempty"`
	DownloadUserAgent    string `form:"download.userAgent,omitempty"  url:"download.userAgent,omitempty"  yaml:"download.userAgent,omitempty"`
	LanguageOverride     string `form:"languageOverride,omitempty"    url:"languageOverride,omitempty"    yaml:"languageOverride,omitempty"`
	EntitiesFilter       string `form:"entities.filter,omitempty"     url:"entities.filter,omitempty"     yaml:"entities.filter,omitempty"`
	EntitiesAllowOverlap bool   `form:"entities.allowOverlap"         url:"entities.allowOverlap"         yaml:"entities,omitempty"`
	EntitiesEnrichment   string `form:"entities.enrichment,omitempty" url:"entities.enrichment,omitempty" yaml:"entities.enrichment,omitempty"`
}

func NewTextRazorRequest(key string) *TextRazorRequest {
	return &TextRazorRequest{
		ApiKey: key,
	}
}

func (t *TextRazorRequest) Analysis(client *http.Client) (*TextRazorResult, error) {
	v, err := query.Values(t)
	if err != nil {
		return nil, err
	}
	s := v.Encode()

	req, err := http.NewRequest("POST", "https://api.textrazor.com/", bytes.NewBufferString(s))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept-encoding: ", "gzip")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tr TextRazorResult
	err = json.Unmarshal(data, &tr)
	if err != nil {
		logInfo.Printf("%s\n", data)
		return nil, err
	}

	switch resp.StatusCode {
	case 400:
		return nil, httpBadRequest
	case 401:
		return nil, httpUnauthorized
	case 413:
		return nil, httpRequestTooLarge
	}

	if !tr.Ok {
		return nil, errors.New(tr.Error)
	}

	return &tr, nil
}

func (t *TextRazorRequest) String() string {
	b, err := yaml.Marshal(t)
	if err != nil {
		logError.Fatal(err)
	}
	return string(b)
}

func (t *TextRazorRequest) SetExtractors(e ...string) {
	t.Extractors = strings.Join(e, ",")
}

type TextRazorResult struct {
	Time             float64
	Response         TextRazorResponse
	Ok               bool
	Error            string
	Message          string
	CustomAnnotation string
	CleanedText      string
	RawText          string
}

type TextRazorResponse struct {
	Entailments        []TextRazorEntailment
	Entities           []TypeRazorEntity
	CoarseTopics       []TextRazorTopic
	Topics             []TextRazorTopic
	NounPhrases        []TextRazorNounPhrase
	Properties         []TextRazorProperty
	Relations          []TextRazorRelation
	Sentences          []TextRazorSentence
	MatchingRules      string
	Language           string
	LanguageIsReliable bool
}

type TypeRazorEntity struct {
	EntityId        string
	EntityEnglishId string
	ConfidenceScore float64
	Type            string
	FreebaseTypes   string
	FreebaseId      string
	MatchingTokens  string
	MatchedText     string
	Data            string
	RelevanceScore  int
	WikiLink        string
}

type TextRazorTopic struct {
	Id       int
	Label    string
	Score    float64
	WikiLink string
}

type TextRazorEntailment struct {
	ContextScore  int
	EntailedTree  string
	WordPositions string
	PriorScore    string
	Score         float64
}

type TextRazorRelationParam struct {
	WordPositions string
	Relation      string
}

type TextRazorNounPhrase struct {
	WordPositions string
}

type TextRazorProperty struct {
	WordPositions     string
	PropertyPositions string
}

type TextRazorRelation struct {
	Params        string
	WordPositions string
}

type TextRazorWord struct {
	StartingPos      int
	EndingPos        int
	Lemma            string
	ParentPosition   string
	PartOfSpeech     string
	Senses           string
	Position         string
	RelationToParent string
	Stem             string
	Token            string
}

type TextRazorSentence struct {
	Words string
}
