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
	ExtractorEntities        string = "entities"
	ExtractorTopics          string = "topics"
	ExtractorWords           string = "words"
	ExtractorPhrases         string = "phrases"
	ExtractorDependancyTrees string = "dependency-trees"
	ExtractorRelations       string = "relations"
	ExtractorEntailments     string = "entailments"
	ExtractorSenses          string = "senses"
)

// cleanup mode constants
const (
	ModeRaw       string = "raw"
	ModeStripTags string = "stripTags"
	ModeCleanHTML string = "cleanHTML"
)

var (
	// ErrHTTPBadRequest error
	ErrHTTPBadRequest = errors.New("Bad Request")
	// ErrHTTPUnauthorized error
	ErrHTTPUnauthorized = errors.New("Unauthorized")
	// ErrHTTPRequestEntityTooLarge error
	ErrHTTPRequestEntityTooLarge = errors.New("Request Entity Too Large")
)

// TextRazorRequest struct
type TextRazorRequest struct {
	Text                 string `form:"text,omitempty"                url:"text,omitempty"                yaml:"text,omitempty"`
	URL                  string `form:"url,omitempty"                 url:"url,omitempty"                 yaml:"url,omitempty"`
	APIKey               string `form:"apiKey"                        url:"apiKey"                        yaml:"apiKey,omitempty"`     // required field
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

// NewTextRazorRequest is a TextRazorRequest constructor
func NewTextRazorRequest(key string) *TextRazorRequest {
	return &TextRazorRequest{
		APIKey: key,
	}
}

// Analysis method
func (t *TextRazorRequest) Analysis(client *http.Client) (*TextRazorResult, error) {
	v, err := query.Values(t)
	if err != nil {
		return nil, err
	}
	s := v.Encode()
	logInfo.Println(s)

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
	tr.URL = t.URL
	err = json.Unmarshal(data, &tr)
	if err != nil {
		logInfo.Printf("%s\n", data)
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return nil, ErrHTTPBadRequest
	case http.StatusUnauthorized:
		return nil, ErrHTTPUnauthorized
	case http.StatusRequestEntityTooLarge:
		return nil, ErrHTTPRequestEntityTooLarge
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

// SetExtractors method
func (t *TextRazorRequest) SetExtractors(e ...string) {
	t.Extractors = strings.Join(e, ",")
}

// TextRazorResult struct
type TextRazorResult struct {
	URL              string
	Time             float64
	Response         TextRazorResponse
	Ok               bool
	Error            string
	Message          string
	CustomAnnotation string
	CleanedText      string
	RawText          string
}

// TextRazorResponse struct
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

// TypeRazorEntity struct
type TypeRazorEntity struct {
	EntityID        string
	EntityEnglishID string
	ConfidenceScore float64
	Type            string
	FreebaseTypes   string
	FreebaseID      string
	MatchingTokens  string
	MatchedText     string
	Data            string
	RelevanceScore  float64
	WikiLink        string
}

// TextRazorTopic struct
type TextRazorTopic struct {
	ID       int
	Label    string
	Score    float64
	WikiLink string
}

// TextRazorEntailment struct
type TextRazorEntailment struct {
	ContextScore  int
	EntailedTree  string
	WordPositions string
	PriorScore    string
	Score         float64
}

// TextRazorRelationParam struct
type TextRazorRelationParam struct {
	WordPositions string
	Relation      string
}

// TextRazorNounPhrase struct
type TextRazorNounPhrase struct {
	WordPositions string
}

// TextRazorProperty struct
type TextRazorProperty struct {
	WordPositions     string
	PropertyPositions string
}

// TextRazorRelation struct
type TextRazorRelation struct {
	Params        string
	WordPositions string
}

// TextRazorWord struct
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

// TextRazorSentence struct
type TextRazorSentence struct {
	Words string
}
