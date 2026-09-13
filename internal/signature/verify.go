package signature

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNoSignature      = errors.New("no signature  found")
	ErrInvalidTimestamp = errors.New("invalid timestamp")
	ErrInvalidSignature = errors.New("invalid signature")
)

func VerifySignature(r *http.Request, sigHead string, tolWin int64, body []byte, signsec string) error {
	sigh := r.Header.Get(sigHead)
	if len(sigh) == 0 {
		return ErrNoSignature
	}
	parts := strings.Split(sigh, ",")
	ts, tsv, tsf := strings.Cut(parts[0], "=")
	if !tsf {
		return ErrNoSignature
	}
	sig, sigv, sigf := strings.Cut(parts[1], "=")
	if !sigf {
		return ErrNoSignature
	}
	var timestampstr string
	var signaturestr string
	if ts == "t" {
		timestampstr = tsv
		signaturestr = sigv
	}
	if sig == "t" {
		timestampstr = sigv
		signaturestr = tsv
	}
	timestamp, err := strconv.ParseInt(timestampstr, 10, 64)
	if err != nil {
		return ErrInvalidTimestamp
	}

	diff := time.Now().Unix() - int64(timestamp)
	if diff > tolWin {
		return ErrInvalidTimestamp
	}

	matched := verifySignature(timestampstr, body, signsec, signaturestr)
	if !matched {
		return ErrInvalidSignature
	}
	return nil

}
