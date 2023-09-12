package tweets

import (
	"testing"
)

func Test_createTweetRequest_validate(t *testing.T) {
	t.Run("text_ok", func(t *testing.T) {
		r := createTweetRequest{
			Text: "text ok",
		}
		if err := r.validate(); err != nil {
			t.Errorf("excpt: err==nil, got:%s", err)
		}

	})
	t.Run("text_fail", func(t *testing.T) {
		expectedErrorSubstring := "Key: 'createTweetRequest.Text' Error:Field validation for 'Text' failed on the 'text' tag"
		r := createTweetRequest{
			Text: "text vmropbmrpobmrtopbmtpomboptmbpotmboptmpobmtpobmptombpotmbpomtpobmtpombpotmbpomtpobmspotmboptmbpot,rgpa,mbpomtbapomtbpmtbopsmtbpomtpobmtpombpotmbopaalrgmaoprbmoptmbnportsnmbpoisnrmt[pbonsmtpo[bnmspotinbmsiotnbmsintbopisnmbiookdfsefrlkngklrglrmepaomtrpobmpondgsgastkgjnstrlkhjoirwhjoistjrhgoi'sjtoi'nbsoitnboisntbionatoigbnstoinboistnboisntboinstoibnsoitnboistnboisntobinsoitbnoistnboinboistnboisntoibnsoirtbnosirntb[santdfgaerngaioenroinaeroibnorinboaienboisnrboianrobinreoibnoitnb[po",
		}
		err := r.validate()
		if err.Error() != expectedErrorSubstring {
			t.Errorf("Expected error message: %s, Actual error message: %s", expectedErrorSubstring, err.Error())
		}

	})
}
