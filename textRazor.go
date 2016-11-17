package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-querystring/query"
	"go4.org/errorutil"
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

// Fetch method
func (t *TextRazorRequest) Fetch(client *http.Client) (io.Reader, error) {
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

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return nil, ErrHTTPBadRequest
	case http.StatusUnauthorized:
		return nil, ErrHTTPUnauthorized
	case http.StatusRequestEntityTooLarge:
		return nil, ErrHTTPRequestEntityTooLarge
	}

	return resp.Body, nil
}

// Analysis method
func (t *TextRazorRequest) Analysis(r io.Reader) (*TextRazorResult, error) {
	var tr TextRazorResult
	tr.URL = t.URL

	dj := json.NewDecoder(r)
	if err := dj.Decode(&tr); err != nil {
		extra := ""
		if serr, ok := err.(*json.SyntaxError); ok {

			if s, ok := r.(io.Seeker); ok {
				if _, serr := s.Seek(0, os.SEEK_SET); serr != nil {
					log.Fatalf("seek error: %v", serr)
				}
			}

			line, col, highlight := errorutil.HighlightBytePosition(r, serr.Offset)
			extra = fmt.Sprintf(":\nError at line %d, column %d (file offset %d):\n%s",
				line, col, serr.Offset, highlight)
		} else if serr, ok := err.(*json.UnmarshalTypeError); ok {

			if s, ok := r.(io.Seeker); ok {
				if _, serr := s.Seek(0, os.SEEK_SET); serr != nil {
					log.Fatalf("seek error: %v", serr)
				}
			}

			line, col, highlight := errorutil.HighlightBytePosition(r, serr.Offset)
			extra = fmt.Sprintf(":\nError at line %d, column %d (file offset %d):\n%s",
				line, col, serr.Offset, highlight)
		}
		return nil, fmt.Errorf("error parsing JSON object %s\n%v",
			extra, err)
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
	Type            []string
	FreebaseTypes   []string
	FreebaseID      string
	MatchingTokens  []int64
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
	ContextScore  float64
	EntailedTree  TextRazorEntailedTree
	WordPositions []int64
	PriorScore    float64
	Score         float64
}

// TextRazorEntailedTree struct
type TextRazorEntailedTree struct {
	Word string
}

// TextRazorRelationParam struct
type TextRazorRelationParam struct {
	WordPositions []int64
	Relation      string
}

// TextRazorNounPhrase struct
type TextRazorNounPhrase struct {
	WordPositions []int64
}

// TextRazorProperty struct
type TextRazorProperty struct {
	WordPositions     []int64
	PropertyPositions []int64
}

// TextRazorRelation struct
type TextRazorRelation struct {
	Params        []TextRazorRelationParam
	WordPositions []int64
}

// TextRazorWord struct
type TextRazorWord struct {
	StartingPos      int
	EndingPos        int
	Lemma            string
	ParentPosition   int64
	PartOfSpeech     string
	Senses           []TextRazorSense
	Position         int64
	RelationToParent string
	Stem             string
	Token            string
}

// TextRazorSentence struct
type TextRazorSentence struct {
	Words []TextRazorWord
}

// TextRazorSense struct
type TextRazorSense struct {
	Synset string
	Score  float64
}
