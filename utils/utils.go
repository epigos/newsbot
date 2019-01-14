package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Map a map interface
type Map map[string]interface{}

func (m *Map) String() string {
	bs, _ := json.Marshal(m)
	return string(bs)
}

// Get get item or return default
func (m Map) Get(key string, d interface{}) interface{} {
	if val, ok := m[key]; ok {
		return val
	}
	return d
}

// Set sets item
func (m Map) Set(key string, value interface{}) {
	m[key] = value
}

// Remove items from Map
func (m Map) Remove(keys ...string) Map {
	for _, key := range keys {
		delete(m, key)
	}
	return m
}

// FilterKeys returns new map specified keys from Map
func (m Map) FilterKeys(keys ...string) Map {
	n := Map{}
	for _, key := range keys {
		n[key] = m.Get(key, nil)
	}
	return n
}

// AppendIfMissing append to slice if i is missing
func AppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

// SliceUniqMap make slice contains unique elements
func SliceUniqMap(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// LinkToDoc retrieves a link and parse it to goquery document
func LinkToDoc(url string) (*goquery.Document, error) {
	// Load the URL
	doc := &goquery.Document{}
	res, err := http.Get(url)
	if err != nil {
		return doc, err
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err = goquery.NewDocumentFromReader(res.Body)
	return doc, err
}

// TrimSpacesList trim white spaces from items in list
func TrimSpacesList(ls []string) []string {
	out := []string{}
	for _, str := range ls {
		v := strings.TrimSpace(str)
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	return out
}

// ProcessTags cleanup news tags
func ProcessTags(t, sep string, repl ...string) []string {
	r := regexp.MustCompile("(?i)" + strings.Join(repl, "|"))
	out := r.ReplaceAllString(strings.ToLower(t), "")
	out = strings.Replace(out, "and", sep, -1)
	return TrimSpacesList(strings.Split(out, sep))
}

// ExtractMetaTags extract open graph tags from page
func ExtractMetaTags(doc *goquery.Document, prefix string) (Map, error) {
	m := Map{}
	// now get open graph tags
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("property"); strings.HasPrefix(name, prefix) {
			value, _ := s.Attr("content")
			name = strings.TrimSpace(strings.TrimPrefix(name, prefix))
			value = strings.TrimSpace(value)
			if value != "" && name != "" {
				m[name] = value
			}
		}
	})
	return m, nil
}

// WebTime internal utility methods
func WebTime(t time.Time) string {
	ftime := t.Format(time.RFC1123)
	if strings.HasSuffix(ftime, "UTC") {
		ftime = ftime[0:len(ftime)-3] + "GMT"
	}
	return ftime
}

// DirExists check if dir exists
func DirExists(dir string) bool {
	d, e := os.Stat(dir)
	switch {
	case e != nil:
		return false
	case !d.IsDir():
		return false
	}

	return true
}

// FileExists checks if file exists
func FileExists(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

// Urlencode is a helper method that converts a map into URL-encoded form data.
// It is a useful when constructing HTTP POST requests.
func Urlencode(data map[string]string) string {
	var buf bytes.Buffer
	for k, v := range data {
		buf.WriteString(url.QueryEscape(k))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(v))
		buf.WriteByte('&')
	}
	s := buf.String()
	return s[0 : len(s)-1]
}

var slugRegex = regexp.MustCompile(`(?i:[^a-z0-9\-_])`)

// Slug is a helper function that returns the URL slug for string s.
// It's used to return clean, URL-friendly strings that can be
// used in routing.
func Slug(s string) string {
	sep := "-"
	if s == "" {
		return ""
	}
	slug := slugRegex.ReplaceAllString(s, sep)
	if slug == "" {
		return ""
	}
	quoted := regexp.QuoteMeta(sep)
	sepRegex := regexp.MustCompile("(" + quoted + "){2,}")
	slug = sepRegex.ReplaceAllString(slug, sep)
	sepEnds := regexp.MustCompile("^" + quoted + "|" + quoted + "$")
	slug = sepEnds.ReplaceAllString(slug, "")
	return strings.ToLower(slug)
}

// NewCookie is a helper method that returns a new http.Cookie object.
// Duration is specified in seconds. If the duration is zero, the cookie is permanent.
// This can be used in conjunction with ctx.SetCookie.
func NewCookie(name string, value string, age int64) *http.Cookie {
	var utctime time.Time
	if age == 0 {
		// 2^31 - 1 seconds (roughly 2038)
		utctime = time.Unix(2147483647, 0)
	} else {
		utctime = time.Unix(time.Now().Unix()+age, 0)
	}
	return &http.Cookie{Name: name, Value: value, Expires: utctime}
}

// Base64EncodeMap encode map to base64
func Base64EncodeMap(m Map) string {
	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	encoded := base64.URLEncoding.EncodeToString(data)
	return encoded
}

// Base64DecodeMap decode encoded base64 string to map
func Base64DecodeMap(enc string) (Map, error) {
	var m Map
	decoded, err := base64.URLEncoding.DecodeString(enc)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(decoded, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

// HumanizeDuration humanizes time.Duration output to a meaningful value,
// golang's default ``time.Duration`` output is badly formatted and unreadable.
func HumanizeDuration(duration *time.Duration) string {
	if duration.Seconds() < 60.0 {
		return fmt.Sprintf("%d seconds", int64(duration.Seconds()))
	}
	if duration.Minutes() < 60.0 {
		//remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d min", int64(duration.Minutes()))
	}
	if duration.Hours() < 24.0 {
		remainingMinutes := math.Mod(duration.Minutes(), 60)
		// remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d hours %d min",
			int64(duration.Hours()), int64(remainingMinutes))
	}
	remainingHours := math.Mod(duration.Hours(), 24)
	remainingMinutes := math.Mod(duration.Minutes(), 60)
	// remainingSeconds := math.Mod(duration.Seconds(), 60)
	return fmt.Sprintf("%d days %d hours %d minutes",
		int64(duration.Hours()/24), int64(remainingHours),
		int64(remainingMinutes))
}

// IsDeployment indicate if app is in deployment mode
func IsDeployment() bool {
	env := os.Getenv("ENV")
	return env == dev || env == prod
}

// GetEnvironment returns env name
func GetEnvironment() string {
	switch os.Getenv("ENV") {
	case dev:
		return "development"
	case prod:
		return "production"
	default:
		return "local"
	}
}

// GetVersion get app version
func GetVersion() string {
	version, ok := os.LookupEnv("VERSION")
	if !ok {
		version = "default"
	}
	return version
}
