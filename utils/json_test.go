package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAndValidateCIDGeneration(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		expected string
	}{
		{
			"empty string", "", "",
		},
		{
			"empty json", "{}", "bafyreigbtj4x7ip5legnfznufuopl4sg4knzc2cof6duas4b3q2fy6swua",
		},

		{
			"test record", "{\"build_artifact_cid\":\"QmP8jTG1m9GSDJLCbeWhVSVgEzCPPwXRdCRuJtQ5Tz9Kc9\",\"repo_registration_record_cid\":\"QmSnuWmxptJZdLJpKRarxBMS2Ju2oANVrgbr2xWbie9b2D\",\"tls_cert_cid\":\"QmbWqxBEKC3P8tqsKc98xmWNzrzDtRLMiMPL8wBuTGsMnR\",\"type\":\"WebsiteRegistrationRecord\",\"url\":\"https://cerc.io\",\"version\":\"0.0.1\"}",
			"bafyreiek4hnoqmits66bjyxswapplweuoqe4en2ux6u772o4y3askpd3ny",
		},
	}

	for _, tc := range testCases {
		newImpl, err := CIDFromJSONBytes([]byte(tc.content))
		require.NoError(t, err)
		require.Equal(t, tc.expected, newImpl)
	}
}
