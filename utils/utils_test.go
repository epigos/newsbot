package utils

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileFunctions(t *testing.T) {
	assert := assert.New(t)
	d := DirExists("../web")
	assert.True(d)
	d = DirExists("../webs")
	assert.False(d)
	d = DirExists("utils.go")
	assert.False(d)

	f := FileExists("utils.go")
	assert.True(f)
	f = FileExists("none.go")
	assert.False(f)
}

func TestURLFunctions(t *testing.T) {
	assert := assert.New(t)

	d := map[string]string{
		"foo": "bar",
		"x":   "y",
	}
	e := Urlencode(d)
	assert.NotNil(e)
	assert.Contains(e, "foo")
	assert.Contains(e, "y")

	s := Slug("this is a slug")
	assert.NotNil(s)
	assert.Equal(s, "this-is-a-slug")

	s = Slug("")
	assert.NotNil(s)
	assert.Equal(s, "")

	doc, err := LinkToDoc("http://example.com")
	assert.NoError(err)
	meta, err := ExtractMetaTags(doc, "")
	assert.NoError(err)
	assert.IsType(Map{}, meta)

	tm := time.Now()
	wt := WebTime(tm)
	assert.Contains(wt, strconv.Itoa(tm.Day()))
}

func TestTrimSpaces(t *testing.T) {
	assert := assert.New(t)
	ts := TrimSpacesList([]string{" text", "another "})
	assert.Equal(ts[0], "text")
	assert.Equal(ts[1], "another")

	ts = ProcessTags("Ghana News | Ghana Politics | Ghana Soccer | Ghana Showbiz", "|", "ghana")
	assert.Equal(ts[0], "news")
	assert.Equal(ts[1], "politics")

	ts = ProcessTags("Ghana News, Crime, Education, Events, Local News, Odd News, Travel and News Archive homepage", ",", "ghana")
	assert.Equal(ts[0], "news")
	assert.Equal(ts[1], "crime")
}

func TestCookie(t *testing.T) {
	assert := assert.New(t)

	c := NewCookie("foo", "bar", 3600)
	assert.IsType(&http.Cookie{}, c)
	assert.Equal(c.Value, "bar")
	assert.Equal(c.Expires, time.Unix(time.Now().Unix()+3600, 0))

	c = NewCookie("foo", "bar", 0)
	assert.IsType(&http.Cookie{}, c)
	assert.Equal(c.Value, "bar")
	assert.Equal(c.Expires, time.Unix(2147483647, 0))
}

func TestMap(t *testing.T) {
	assert := assert.New(t)
	m := Map{"key": "value"}

	assert.Equal(m.Get("key", nil), "value")
	assert.Equal(m.Get("one", 1), 1)
	m.Remove("key")
	assert.Equal(m, Map{})
}

func TestSliceFunctions(t *testing.T) {
	assert := assert.New(t)
	sl := []string{"1", "2", "3", "3"}
	unique := SliceUniqMap(sl)
	assert.Len(unique, 3)
	assert.NotEqual(sl, unique)

	ap := AppendIfMissing(sl, "2")
	assert.Equal(sl, ap)
}

func TestBase64Encoding(t *testing.T) {
	assert := assert.New(t)

	m := Map{"foo": "bar"}
	enc := Base64EncodeMap(m)

	dec, err := Base64DecodeMap(enc)
	assert.NoError(err)
	assert.Equal(dec, m)
}

func TestTextAnalysis(t *testing.T) {
	assert := assert.New(t)
	text := "Your long raw text, it could be a book. Lorem ipsum."

	ta := NewTextAnalysis(text, text)
	rtime := ta.ReadingTime()
	assert.Equal(rtime.String(), "3s")

	sumr := ta.Sentences(1)
	assert.Equal(len(sumr), 1)

	tags := ta.Tags()
	assert.Contains(tags, "lorem")
}
