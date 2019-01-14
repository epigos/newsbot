package utils

import (
	"time"

	textrank "github.com/DavidBelicza/TextRank"
	"github.com/jdkato/prose/summarize"
)

const (
	wordsPerMinute = 200
)

// TextAnalysis text analysis struct
type TextAnalysis struct {
	Text        *textrank.TextRank
	Description *textrank.TextRank
	Doc         *summarize.Document
}

func rankText(t string) *textrank.TextRank {
	// TextRank object
	tr := textrank.NewTextRank()
	// Default Rule for parsing.
	rule := textrank.NewDefaultRule()
	// Default Language for filtering stop words.
	language := textrank.NewDefaultLanguage()
	// Default algorithm for ranking text.
	algorithmDef := textrank.NewDefaultAlgorithm()

	// Add text.
	tr.Populate(t, language, rule)
	// Run the ranking.
	tr.Ranking(algorithmDef)
	return tr
}

// NewTextAnalysis returns new text analysis
func NewTextAnalysis(text, desc string) *TextAnalysis {

	ta := &TextAnalysis{
		Text:        rankText(text),
		Description: rankText(desc),
		Doc:         summarize.NewDocument(text),
	}

	return ta
}

// Tags return top words
func (t *TextAnalysis) Tags() []string {
	// Get all words order by weight.
	words := textrank.FindSingleWords(t.Description)
	out := make([]string, len(words))

	for i, word := range words {
		out[i] = word.Word
	}
	return SliceUniqMap(out)
}

// Sentences return top sentences
func (t *TextAnalysis) Sentences(length int) []string {
	// Get the most important 10 sentences. Importance by word occurrence.
	sentences := textrank.FindSentencesByWordQtyWeight(t.Text, length)

	out := make([]string, length)
	for i, sen := range sentences {
		out[i] = sen.Value
	}
	return TrimSpacesList(out)
}

// ReadingTime estimates how long an article will take to read
// based on 200 words per minutes
func (t *TextAnalysis) ReadingTime() *time.Duration {
	minutes := 60 * (t.Doc.NumWords / wordsPerMinute)
	rtime := time.Second * time.Duration(minutes)
	return &rtime
}
