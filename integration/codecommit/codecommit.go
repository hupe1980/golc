package codecommit

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

var (
	urlRgx = regexp.MustCompile(`^https://git-codecommit\.(.*)\.(amazonaws\.com|amazonaws\.com\.cn)/v1/repos/.*$`)
)

// Signer implements the AWS authenticated V4 Signature Specification for
// generating authenticated requests to CodeCommit repositories. The signature is
// generated using the V4 Signature Specification
// see https://docs.aws.amazon.com/general/latest/gr/sigv4_signing.html
type Signer struct {
	service     string
	region      string
	credentials aws.Credentials
	requestTime time.Time
}

// NewSigner creates a new V4 signer for signing CodeCommit URLs
func NewSigner(creds aws.Credentials) *Signer {
	return &Signer{
		service:     "codecommit",
		credentials: creds,
	}
}

// Sign will sign a CodeCommit clone URL using the AWS authenticated V4 Signature
// Specification. As CodeCommit is accessed directly through a git client over HTTPS,
// authentication details must be supplied to CodeCommit using Basic User Autentication.
func (s *Signer) Sign(cloneURL string) (string, error) {
	var err error
	if s.region, err = identifyRegion(cloneURL); err != nil {
		return "", err
	}

	// see https://github.com/aws/git-remote-codecommit/blob/master/git_remote_codecommit/__init__.py
	s.requestTime = time.Now().UTC()
	req, _ := http.NewRequest("GIT", cloneURL, http.NoBody)
	cr := s.canonicalRequest(req)
	sts := s.stringToSign(cr)
	sig := s.signature(sts)

	passw := fmt.Sprintf("%sZ%s", s.requestTime.Format("20060102T150405"), fmt.Sprintf("%x", sig))
	uname := url.QueryEscape(s.credentials.AccessKeyID + "%" + s.credentials.SessionToken)

	return strings.Replace(cloneURL, "https://", fmt.Sprintf("https://%s:%s@", uname, passw), 1), nil
}

func identifyRegion(url string) (string, error) {
	if m := urlRgx.FindStringSubmatch(url); len(m) > 1 {
		return m[1], nil
	}

	return "", errors.New("no region found in malformed codecommit URL")
}

// canonicalRequest generates a canonical request
// see https://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
func (s *Signer) canonicalRequest(req *http.Request) []byte {
	// CodeCommit doesn't support query parameters or a payload, so omit both from the request
	cr := new(bytes.Buffer)
	fmt.Fprintf(cr, "%s\n", req.Method)
	fmt.Fprintf(cr, "%s\n", req.URL.Path)
	fmt.Fprintf(cr, "%s\n", "")
	fmt.Fprintf(cr, "host:%s\n\n", req.URL.Host)
	fmt.Fprintf(cr, "%s\n", "host")
	fmt.Fprintf(cr, "%s", "")

	return cr.Bytes()
}

// stringToSign creates a string to sign
// see https://docs.aws.amazon.com/general/latest/gr/sigv4-create-string-to-sign.html
func (s *Signer) stringToSign(cr []byte) []byte {
	sts := new(bytes.Buffer)
	fmt.Fprint(sts, "AWS4-HMAC-SHA256\n")
	fmt.Fprintf(sts, "%s\n", s.requestTime.Format("20060102T150405"))
	fmt.Fprintf(sts, "%s/%s/%s/aws4_request\n", s.requestTime.Format("20060102"), s.region, s.service)
	fmt.Fprintf(sts, "%s", fmt.Sprintf("%x", v4Hash(cr)))

	return sts.Bytes()
}

// signature creates the V4 signature
// see https://docs.aws.amazon.com/general/latest/gr/sigv4-calculate-signature.html
func (s *Signer) signature(sts []byte) []byte {
	dsk := v4HMAC([]byte("AWS4"+s.credentials.SecretAccessKey), []byte(s.requestTime.Format("20060102")))
	dsk = v4HMAC(dsk, []byte(s.region))
	dsk = v4HMAC(dsk, []byte(s.service))
	dsk = v4HMAC(dsk, []byte("aws4_request"))

	return v4HMAC(dsk, sts)
}

func v4Hash(in []byte) []byte {
	h := sha256.New()
	h.Write(in)

	return h.Sum(nil)
}

func v4HMAC(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)

	return h.Sum(nil)
}
