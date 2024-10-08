package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newJSONRequest(t *testing.T, method, target string, body any) *http.Request {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		t.Fatalf("error marshaling %T as JSON: %s", body, err)
	}
	return httptest.NewRequest(method, target, &buf)
}

func anyFutureDate() string {
	t := time.Now().UTC().AddDate(rand.Intn(5), 1+rand.Intn(12), 0)
	return t.Format("01/2006")
}

func anyPastDate() string {
	t := time.Now().UTC().AddDate(-rand.Intn(5), -1-rand.Intn(12), -1)
	return t.Format("01/2006")
}

func TestValidationHandler_InvalidJSON(t *testing.T) {
	handler := ValidationHandler()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/validate", strings.NewReader("foobarbaz"))
	req.Header.Set("Content-Type", "text/plain")

	handler.ServeHTTP(rec, req)
	res := rec.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status code: want %d have %s",
			http.StatusBadRequest, res.Status)
	}

	var body validationResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("error parsing JSON response: %s", err)
	}

	if body.Valid {
		t.Fatalf("expected 'valid' to be false, got true")
	}

	if body.Error.Code != errGeneralError {
		t.Fatalf("API error code mismatch: want %d have %d", errGeneralError, body.Error.Code)
	}
}

func TestValidationHandler_CardValidation(t *testing.T) {
	tests := []struct {
		number  string
		expDate string
		code    apiErrorCode
	}{
		{"", anyFutureDate(), errMalformedNumber},
		{"abcdabcdabcdabcd", anyFutureDate(), errMalformedNumber},
		{"123", anyFutureDate(), errMalformedNumber},
		{"123123123123123123123123123123", anyFutureDate(), errMalformedNumber},

		// Malformed expiration dates.
		{"4539441071007551", "", errMalformedDate},
		{"4539984459069503", "-99/-100", errMalformedDate},

		// Expired cards.
		{"5439223588334647", anyPastDate(), errCardExpired},
		{"5246132897434423", anyPastDate(), errCardExpired},
		{"4111111111111111", anyPastDate(), errCardExpired},
		{"4012888888881881", anyPastDate(), errCardExpired},
		{"4539723775949752", anyPastDate(), errCardExpired},

		// Unknown issuers.
		{"9550998650131033", anyFutureDate(), errUnknownIssuer},
		{"9566111111111113", anyFutureDate(), errUnknownIssuer},
		{"7555555555554444", anyFutureDate(), errUnknownIssuer},
		{"1105105105105100", anyFutureDate(), errUnknownIssuer},
		{"8123856130355023", anyFutureDate(), errUnknownIssuer},

		// Invalid account numbers (based on Luhn's check).
		{"4111111111111121", anyFutureDate(), errInvalidAccountNumber},
		{"4111111111111611", anyFutureDate(), errInvalidAccountNumber},
		{"6212345678903265", anyFutureDate(), errInvalidAccountNumber},
		{"6212345678902232", anyFutureDate(), errInvalidAccountNumber},
		{"6212345678910128", anyFutureDate(), errInvalidAccountNumber},
		{"6212345678911036", anyFutureDate(), errInvalidAccountNumber},
	}

	handler := ValidationHandler()

	for _, tc := range tests {
		rec := httptest.NewRecorder()
		req := newJSONRequest(t, "POST", "/validate", creditCardInfo{
			CardNumber:     tc.number,
			ExpirationDate: tc.expDate,
		})

		handler.ServeHTTP(rec, req)
		res := rec.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("unexpected status code: want %d have %s",
				http.StatusUnprocessableEntity, res.Status)
		}

		var body validationResponse
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatalf("error parsing JSON response: %s", err)
		}

		if body.Valid {
			t.Fatalf("expected 'valid' to be false, got true")
		}

		if body.Error.Code != tc.code {
			t.Fatalf("API error code mismatch: want %d have %d", tc.code, body.Error.Code)
		}
	}
}

func TestValidationHandler(t *testing.T) {
	tests := []struct {
		number  string
		expDate string
	}{
		{"4111111111111111", anyFutureDate()},
		{"4111111111111111", anyFutureDate()},
		{"6212345678901265", anyFutureDate()},
		{"6212345678901232", anyFutureDate()},
		{"6212345678900028", anyFutureDate()},
		{"6212345678900036", anyFutureDate()},
		{"6212345678900085", anyFutureDate()},
		{"6212345678900093", anyFutureDate()},
		{"62123456789002", anyFutureDate()},
		{"621234567890003", anyFutureDate()},
		{"62123456789000003", anyFutureDate()},
		{"621234567890000002", anyFutureDate()},
		{"6212345678900000003", anyFutureDate()},
		{"371255422728692", anyFutureDate()},
		{"343030658955854", anyFutureDate()},
		{"343809910826775", anyFutureDate()},
		{"375114534473083", anyFutureDate()},
		{"346592318638144", anyFutureDate()},
		{"370558655421039", anyFutureDate()},
		{"375658534489747", anyFutureDate()},
		{"345923078383525", anyFutureDate()},
		{"346812451160932", anyFutureDate()},
		{"348236720246983", anyFutureDate()},
		{"378282246310005", anyFutureDate()},
		{"371449635398431", anyFutureDate()},
		{"378734493671000", anyFutureDate()},
		{"30569309025904", anyFutureDate()},
		{"38520000023237", anyFutureDate()},
		{"36683902117168", anyFutureDate()},
		{"36127951797629", anyFutureDate()},
		{"30351447989048", anyFutureDate()},
		{"36429462154559", anyFutureDate()},
		{"30006032187206", anyFutureDate()},
		{"30284611236484", anyFutureDate()},
		{"30195830207181", anyFutureDate()},
		{"38754423755923", anyFutureDate()},
		{"30027676109538", anyFutureDate()},
		{"38126374068061", anyFutureDate()},
		{"6011111111111117", anyFutureDate()},
		{"6011000990139424", anyFutureDate()},
		{"6011747305012486", anyFutureDate()},
		{"6011870537589871", anyFutureDate()},
		{"6011615793627497", anyFutureDate()},
		{"6011883712576324", anyFutureDate()},
		{"6011204723759652", anyFutureDate()},
		{"6011413989493944", anyFutureDate()},
		{"6011635930056772", anyFutureDate()},
		{"6011646259058190", anyFutureDate()},
		{"6011299144770809", anyFutureDate()},
		{"6011499150862066", anyFutureDate()},
		{"3530111333300000", anyFutureDate()},
		{"3566002020360505", anyFutureDate()},
		{"3550998650131033", anyFutureDate()},
		{"3566111111111113", anyFutureDate()},
		{"5555555555554444", anyFutureDate()},
		{"5105105105105100", anyFutureDate()},
		{"5123856130355023", anyFutureDate()},
		{"5429351731128749", anyFutureDate()},
		{"5408200584401528", anyFutureDate()},
		{"5369356574320040", anyFutureDate()},
		{"5483000279065275", anyFutureDate()},
		{"5447025105449174", anyFutureDate()},
		{"5532053336988980", anyFutureDate()},
		{"5594881983691894", anyFutureDate()},
		{"5439223588334647", anyFutureDate()},
		{"5246132897434423", anyFutureDate()},
		{"4111111111111111", anyFutureDate()},
		{"4012888888881881", anyFutureDate()},
		{"4539723775949752", anyFutureDate()},
		{"4539983514929271", anyFutureDate()},
		{"4556299996115273", anyFutureDate()},
		{"4556648760621678", anyFutureDate()},
		{"4556134715955412", anyFutureDate()},
		{"4539419864694884", anyFutureDate()},
		{"4024007173079426", anyFutureDate()},
		{"4539441071007551", anyFutureDate()},
		{"4539984459069503", anyFutureDate()},
	}

	handler := ValidationHandler()

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s %s", tc.number, tc.expDate), func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := newJSONRequest(t, "POST", "/validate", creditCardInfo{
				CardNumber:     tc.number,
				ExpirationDate: tc.expDate,
			})

			handler.ServeHTTP(rec, req)
			res := rec.Result()

			if res.StatusCode != http.StatusOK {
				t.Fatalf("unexpected status code: want %d have %s",
					http.StatusOK, res.Status)
			}

			var body validationResponse
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Fatalf("error parsing JSON response: %s", err)
			}

			if !body.Valid {
				t.Fatalf("expected 'valid' to be true, got false (%+v)", body.Error)
			}
		})
	}
}
