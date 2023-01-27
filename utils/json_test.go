package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAndValidateCIDGeneration(t *testing.T) {
	testCases := []struct {
		name    string
		content string
	}{
		{
			"empty string", "",
		},
		{
			"empty json", "{}",
		},

		{
			"test record", "\\xa6curlohttps://cerc.iodtypex\\x19WebsiteRegistrationRecordgversione0.0.1ltls_cert_cidx.QmbWqxBEKC3P8tqsKc98xmWNzrzDtRLMiMPL8wBuTGsMnRrbuild_artifact_cidx.QmP8jTG1m9GSDJLCbeWhVSVgEzCPPwXRdCRuJtQ5Tz9Kc9x\\x1crepo_registration_record_cidx.QmSnuWmxptJZdLJpKRarxBMS2Ju2oANVrgbr2xWbie9b2D",
		},
	}

	for _, tc := range testCases {
		deprecatedAndCorrect, _ := CIDFromJSONBytes([]byte(tc.content))
		newImpl, err := CIDFromJSONBytesUsingIpldPrime([]byte(tc.content))
		require.NoError(t, err)
		require.Equal(t, deprecatedAndCorrect, newImpl, tc.name)
	}
}
