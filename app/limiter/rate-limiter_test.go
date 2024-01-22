package limiter

import (
	"net/http"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestShouldBeExtractClientInfoFromRequestWithToken(t *testing.T) {
	// Arrange
	rateLimiter := RateLimiter{}
	mockedRequest := &http.Request{
		Header: map[string][]string{
			"Api-Key": {"token-abc"},
		},
	}
	// Act
	clientInfo := rateLimiter.ExtractClientInfoFromRequest(mockedRequest)
	// Assert
	assert.Equal(t, clientInfo.Key, "token-abc")
	assert.Equal(t, clientInfo.RequestLimit, int64(10))
	assert.Equal(t, clientInfo.RequestInterval, int64(10))
}

func TestShouldBeExtractClientInfoFromRequestWithIP(t *testing.T) {
	// Arrange
	rateLimiter := RateLimiter{}
	mockedRequest := &http.Request{
		Header: map[string][]string{
			"X-Real-Ip": {"123.123.123"},
		},
	}
	// Act
	clientInfo := rateLimiter.ExtractClientInfoFromRequest(mockedRequest)
	// Assert
	assert.Equal(t, clientInfo.Key, "123.123.123")
	assert.Equal(t, clientInfo.RequestLimit, viper.GetInt64("REQUEST_LIMIT"))
	assert.Equal(t, clientInfo.RequestInterval, viper.GetInt64("REQUEST_SECONDS_INTERVAL"))
}
