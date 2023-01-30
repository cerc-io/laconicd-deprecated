package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAndValidateCIDGeneration(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		expected string
	}{
		// empty string and empty json blows up
		// {
		//  	"empty string", "", "bafyreiengp2sbi6ez34a2jctv34bwyjl7yoliteleaswgcwtqzrhmpyt2m",
		// },
		// {
		//		"empty json", "{}", "bafyreihpfkdvib5muloxlj5b3tgdwibjdcu3zdsuhyft33z7gtgnlzlkpm",
		// },

		{
			"test record", "{\"build_artifact_cid\":\"QmP8jTG1m9GSDJLCbeWhVSVgEzCPPwXRdCRuJtQ5Tz9Kc9\",\"repo_registration_record_cid\":\"QmSnuWmxptJZdLJpKRarxBMS2Ju2oANVrgbr2xWbie9b2D\",\"tls_cert_cid\":\"QmbWqxBEKC3P8tqsKc98xmWNzrzDtRLMiMPL8wBuTGsMnR\",\"type\":\"WebsiteRegistrationRecord\",\"url\":\"https://cerc.io\",\"version\":\"0.0.1\"}",
			"bafyreiek4hnoqmits66bjyxswapplweuoqe4en2ux6u772o4y3askpd3ny",
		},
	}

	for _, tc := range testCases {
		deprecatedAndCorrect, _ := CIDFromJSONBytes([]byte(tc.content))
		newImpl, err := CIDFromJSONBytesUsingIpldPrime([]byte(tc.content))
		require.NoError(t, err)
		require.Equal(t, deprecatedAndCorrect, newImpl, tc.name)
		require.Equal(t, tc.expected, newImpl)
	}
}
