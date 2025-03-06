package consensus

import "errors"

var (
	errViewChangeMajority = errors.New("majority does not agree with viewchange")
	errNewViewTimeout     = errors.New("new leader missed new-view message")
)
