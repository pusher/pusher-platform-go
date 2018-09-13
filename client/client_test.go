package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestClientRequestSuccess(t *testing.T) {
	jwt := "some.jwt.string"

	var (
		recievedPostBody    string
		recievedQueryParams url.Values
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/qp", func(w http.ResponseWriter, r *http.Request) {
		recievedQueryParams = r.URL.Query()
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		recievedPostBody = string(bodyBytes)
		w.WriteHeader(http.StatusCreated)
	})

	server := httptest.NewTLSServer(mux)
	defer server.Close()

	uri, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Failed to parse server url with error: %+v", err)
	}

	client := New(Options{
		Host: uri.Host,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	t.Run("GET request", func(t *testing.T) {
		resp, err := client.Request(context.Background(), RequestOptions{
			Method: http.MethodGet,
			Jwt:    &jwt,
		})
		if err != nil {
			t.Fatalf("Failed to request backend resource with error: %+v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected a 200 OK, but got %v", resp.StatusCode)
		}
	})

	t.Run("GET with query params", func(t *testing.T) {
		resp, err := client.Request(context.Background(), RequestOptions{
			Method: http.MethodGet,
			Path:   "/qp",
			QueryParams: &url.Values{
				"foo": []string{"bar"},
			},
			Jwt: &jwt,
		})
		if err != nil {
			t.Fatalf("Failed to request backend resource with error: %+v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected a 200 OK, but got %v", resp.StatusCode)
		}

		if recievedQueryParams.Get("foo") != "bar" {
			t.Fatalf("Expected key `foo` to have value `bar` in query params")
		}
	})

	t.Run("POST request", func(t *testing.T) {
		resp, err := client.Request(context.Background(), RequestOptions{
			Method: http.MethodPost,
			Path:   "/post",
			Body:   bytes.NewReader([]byte("hello")),
			Jwt:    &jwt,
		})
		if err != nil {
			t.Fatalf("Failed to request backend resource with error: %+v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected a 201, but got %v", resp.StatusCode)
		}

		if recievedPostBody != "hello" {
			t.Fatalf("Expected post body to be `%s`, but got `%s`", "hello", recievedPostBody)
		}
	})
}

func TestClientRequestErrors(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyIn0.iYtIg5fTkceCsK2TE2cELE-U10_joFdQ3S5tswA6jT0"

	mux := http.NewServeMux()
	mux.HandleFunc("/not_json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("invalid json"))
	})
	mux.HandleFunc("/bad_request", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid_name"}`))
	})

	server := httptest.NewTLSServer(mux)
	defer server.CloseClientConnections()

	uri, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Failed to parse server url with error: %+v", err)
	}

	client := New(Options{
		Host:               uri.Host,
		DontFollowRedirect: true,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	t.Run("Invalid JSON bodies in errors are detected", func(t *testing.T) {
		_, err := client.Request(context.Background(), RequestOptions{
			Method: http.MethodGet,
			Path:   "/not_json",
			Jwt:    &jwt,
		})

		switch err := err.(type) {
		case BodyNotJSONError:
			return
		default:
			t.Fatalf("Expected BodyNotJsonError but got %s", err)
		}
	})

	t.Run("Responses with status 4xx should have an error", func(t *testing.T) {
		_, err := client.Request(context.Background(), RequestOptions{
			Method: http.MethodGet,
			Path:   "/bad_request",
			Jwt:    &jwt,
		})

		switch err := (err).(type) {
		case *ErrorResponse:
			if err.Status != http.StatusBadRequest {
				t.Fatalf("Expected a 400 status, but got %v", err.Status)
			}

			errorInfo := err.Info.(map[string]interface{})
			if errorType := errorInfo["error"].(string); errorType != "invalid_name" {
				t.Fatalf("Expected error to be invalid_name, but got %s", errorType)
			}
		default:
			t.Fatalf("Expected ErrorResponse but got %s", err)
		}
	})
}

type testClient struct {
	name                   string
	dontFollowRedirect     bool
	path                   string
	expectedStatus         int
	expectedLocationHeader string
	expectedBody           string
}

func TestClientDoesNotFollowRedirect(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyIn0.iYtIg5fTkceCsK2TE2cELE-U10_joFdQ3S5tswA6jT0"

	mux := http.NewServeMux()
	mux.HandleFunc("/original", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/redirect")
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte("original"))
	})
	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("redirect"))
	})

	server := httptest.NewTLSServer(mux)
	defer server.CloseClientConnections()

	uri, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Failed to parse server url with error: %+v", err)
	}

	testCases := []testClient{
		testClient{
			name:                   "by default follow redirect",
			dontFollowRedirect:     false, // default value
			path:                   "/original",
			expectedStatus:         http.StatusOK,
			expectedLocationHeader: "",
			expectedBody:           "redirect",
		},
		testClient{
			name:                   "disable follow redirect",
			dontFollowRedirect:     true,
			path:                   "/original",
			expectedStatus:         http.StatusTemporaryRedirect,
			expectedLocationHeader: "/redirect",
			expectedBody:           "original",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			client := New(Options{
				Host:               uri.Host,
				DontFollowRedirect: testCase.dontFollowRedirect,
				TLSConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			})

			res, err := client.Request(context.Background(), RequestOptions{
				Method: "GET",
				Path:   testCase.path,
				Jwt:    &jwt,
			})
			if err != nil {
				t.Fatalf("Failed to GET resource from %s: %+v", err, testCase.path)
			}

			if res.StatusCode != testCase.expectedStatus {
				t.Fatalf(
					"Expected status code to be %d, but got %d",
					testCase.expectedStatus,
					res.StatusCode,
				)
			}

			if res.Header.Get("Location") != testCase.expectedLocationHeader {
				t.Fatalf(
					"Expected to receive Location header with value `%s` but got `%s`",
					testCase.expectedLocationHeader,
					res.Header.Get("Location"),
				)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Unable to read response body, error: %v\n", err)
			}

			actual := strings.TrimSpace(string(body))
			if actual != testCase.expectedBody {
				t.Fatalf("Expected response body to be `%s` but got `%s`", testCase.expectedBody, actual)
			}
		})
	}
}

func TestClientFailsWhenCipherSuitesAreDeclaredExplicitly(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyIn0.iYtIg5fTkceCsK2TE2cELE-U10_joFdQ3S5tswA6jT0"

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.CloseClientConnections()

	uri, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Failed to parse server url with error: %+v", err)
	}

	client := New(Options{
		Host: uri.Host,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
			CipherSuites: []uint16{
				tls.TLS_RSA_WITH_RC4_128_SHA,
				tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			},
		},
	})

	res, err := client.Request(context.Background(), RequestOptions{
		Method: http.MethodGet,
		Path:   "/",
		Jwt:    &jwt,
	})
	if err != nil {
		t.Fatalf("Failed to GET the resource from the server: %+v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code to be 200, but got %d", res.StatusCode)
	}
}
