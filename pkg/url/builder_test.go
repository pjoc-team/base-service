package url

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompactUrl(t *testing.T) {
	assert.Equal(t, "http://google.com/hello?a=b", CompactUrl("http://google.com/", "hello", "?a=b"))
	assert.Equal(t, "http://google.com/hello?a=b", CompactUrl("http://google.com/", "/hello", "?a=b"))
	assert.Equal(t, "http://google.com/hello/?a=b", CompactUrl("http://google.com/", "hello/", "a=b"))
	assert.Equal(t, "http://google.com/hello/?a=b", CompactUrl("http://google.com", "hello/", "a=b"))
	assert.Equal(t, "http://google.com/hello/?a=b", CompactUrl("http://google.com", "hello/", "a=b"))
	assert.Equal(t, "http://google.com/hello/?a=b", CompactUrl("http://google.com", "hello/", "a=b"))
	assert.Equal(t, "http://google.com/hello/?a=b", CompactUrl("http://google.com", "hello/", "a=b"))

	assert.Equal(t, "hello/?a=b", CompactUrl("", "hello/", "a=b"))
	assert.Equal(t, "http://google.com/?a=b", CompactUrl("http://google.com/", "", "a=b"))
	assert.Equal(t, "hello/?a=b", CompactUrl("", "hello/", "a=b"))
	assert.Equal(t, "http://google.com/hello", CompactUrl("http://google.com", "hello", ""))

}
