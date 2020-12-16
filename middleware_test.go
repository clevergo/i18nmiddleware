package i18nmiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"clevergo.tech/clevergo"
	"github.com/alecthomas/assert"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func TestWithFormField(t *testing.T) {
	m := &middleware{}
	field := "_lang"
	WithFormField(field)(m)
	assert.Equal(t, field, m.formField)
}

func newTestBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.ParseMessageFileBytes([]byte(`{"home": "Home"}`), "en.json")
	bundle.ParseMessageFileBytes([]byte(`{"home": "主页"}`), "zh-CN.json")
	return bundle
}

func testHandle(c *clevergo.Context) error {
	localizer := Localizer(c)
	s, _, _ := localizer.LocalizeWithTag(&i18n.LocalizeConfig{
		MessageID: "home",
	})
	return c.String(http.StatusOK, s)
}

func TestNew(t *testing.T) {
	fn := New(newTestBundle())
	handle := fn(testHandle)
	runTests(t, handle, []testCase{
		{httptest.NewRequest(http.MethodGet, "/", nil), "Home"},
		{httptest.NewRequest(http.MethodGet, "/?lang=en", nil), "Home"},
		{httptest.NewRequest(http.MethodGet, "/?lang=zh", nil), "主页"},
		{httptest.NewRequest(http.MethodGet, "/?_lang=zh", nil), "Home"}, // incorrect form field
	})
}

func TestNewWithOption(t *testing.T) {
	fn := New(newTestBundle(), WithFormField("_lang"))
	handle := fn(testHandle)
	runTests(t, handle, []testCase{
		{httptest.NewRequest(http.MethodGet, "/", nil), "Home"},
		{httptest.NewRequest(http.MethodGet, "/?_lang=en", nil), "Home"},
		{httptest.NewRequest(http.MethodGet, "/?_lang=zh", nil), "主页"},
		{httptest.NewRequest(http.MethodGet, "/?lang=zh", nil), "Home"}, // incorrect form field
	})
}

type testCase struct {
	req  *http.Request
	body string
}

func runTests(t *testing.T, handle clevergo.Handle, testCases []testCase) {
	for _, testCase := range testCases {
		resp := httptest.NewRecorder()
		c := &clevergo.Context{
			Request:  testCase.req,
			Response: resp,
		}
		handle(c)
		assert.Equal(t, testCase.body, resp.Body.String())
	}
}

func TestLocalizer(t *testing.T) {
	var localizer *i18n.Localizer
	c := &clevergo.Context{
		Request: httptest.NewRequest(http.MethodGet, "/", nil),
	}
	c.WithValue(LocalizerKey, localizer)
	assert.Equal(t, localizer, Localizer(c))
}
