package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"reflect"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"
	"os"
	"errors"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func TestParseAddItemRequest(t *testing.T) {
	t.Parallel()

	type wants struct {
		req *AddItemRequest
		err bool
	}

	// STEP 6-1: define test cases
	cases := map[string]struct {
		args map[string]string
		wants
	}{
		"ok: valid request": {
			args: map[string]string{
				"name":     "name", // fill here
				"category": "category", // fill here
			},
			wants: wants{
				req: &AddItemRequest{
					Name: "name", // fill here
					Category: "category", // fill here
				},
				err: false,
			},
		},
		"ng: empty request": {
			args: map[string]string{},
			wants: wants{
				req: nil,
				err: true,
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// prepare request body
			values := url.Values{}
			for k, v := range tt.args {
				values.Set(k, v)
			}

			// prepare HTTP request
			req, err := http.NewRequest("POST", "http://localhost:9000/items", strings.NewReader(values.Encode()))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// execute test target
			got, err := parseAddItemRequest(req)

			// confirm the result
			if err != nil {
				if !tt.err {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}
			if diff := cmp.Diff(tt.wants.req, got); diff != "" {
				t.Errorf("unexpected request (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHelloHandler(t *testing.T) {
	t.Parallel()

	// Please comment out for STEP 6-2
	// predefine what we want
	type wants struct {
	 	code int               // desired HTTP status code
	 	body map[string]string // desired body
	}
	want := wants{
		code: http.StatusOK,
		body: map[string]string{"message": "Hello, world!"},
	}

	// set up test
	req := httptest.NewRequest("GET", "/hello", nil)
	res := httptest.NewRecorder()

	h := &Handlers{}
	h.Hello(res, req)

	// STEP 6-2: confirm the status code
	if res.Code != want.code {
		t.Errorf("unexpected status code: want %d, got %d", want.code, res.Code)
	}

	// STEP 6-2: confirm response body
	var resBody map[string]string
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	if err != nil {
		t.Errorf("Error decoding JSON: %v", err)
	}

	if reflect.DeepEqual(resBody, want.body) {
		t.Errorf("unexpected body: diff=%s", cmp.Diff(resBody, want.body))
	}
}

func TestAddItem(t *testing.T) {
	t.Parallel()

	type wants struct {
		code int
	}
	cases := map[string]struct {
		args     map[string]string
		injector func(m *MockItemRepository)
		wants
	}{
		"ok: correctly inserted": {
			args: map[string]string{
				"name":     "used iPhone 16e",
				"category": "phone",
			},
			injector: func(m *MockItemRepository) {
				// STEP 6-3: define mock expectation
				m.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().GetItem(gomock.Any(), gomock.Any()).Return(&Item{}, nil).AnyTimes()
				// succeeded to insert
			},
			wants: wants{
				code: http.StatusOK,
			},
		},
		"ng: failed to insert": {
			args: map[string]string{
				"name":     "used iPhone 16e",
				"category": "phone",
			},
			injector: func(m *MockItemRepository) {
				// STEP 6-3: define mock expectation
				m.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(errors.New("insert failed"))
				m.EXPECT().GetItem(gomock.Any(), gomock.Any()).Return(&Item{}, nil).AnyTimes()
				// failed to insert
			},
			wants: wants{
				code: http.StatusInternalServerError,
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			mockIR := NewMockItemRepository(ctrl)
			tt.injector(mockIR)
			h := &Handlers{itemRepo: mockIR}

			values := url.Values{}
			for k, v := range tt.args {
				values.Set(k, v)
			}
			req := httptest.NewRequest("POST", "/items", strings.NewReader(values.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()
			h.AddItem(rr, req)

			if tt.wants.code != rr.Code {
				t.Errorf("expected status code %d, got %d", tt.wants.code, rr.Code)
			}
			if tt.wants.code >= 400 {
				return
			}

			for _, v := range tt.args {
				if !strings.Contains(rr.Body.String(), v) {
					t.Errorf("response body does not contain %s, got: %s", v, rr.Body.String())
				}
			}
		})
	}
}

// STEP 6-4: uncomment this test
func TestAddItemE2e(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	db, dbPath, closers, err := setupDB(t)
 	if err != nil {
 		t.Fatalf("failed to set up database: %v", err)
 	}
 	t.Cleanup(func() {
 		for _, c := range closers {
 			c()
 		}
 	})

 	type wants struct {
 		code int
 	}
 	cases := map[string]struct {
 		args map[string]string
 		wants
 	}{
 		"ok: correctly inserted": {
 			args: map[string]string{
 				"name":     "used iPhone 16e",
 				"category": "phone",
 			},
 			wants: wants{
 				code: http.StatusOK,
 			},
 		},
 		"ng: failed to insert": {
 			args: map[string]string{
 				"name":     "",
 				"category": "phone",
 			},
 			wants: wants{
 				code: http.StatusBadRequest,
 			},
 		},
 	}

 	for name, tt := range cases {
 		t.Run(name, func(t *testing.T) {
 			h := &Handlers{itemRepo: &itemRepository{dbPath: dbPath}}

 			values := url.Values{}
 			for k, v := range tt.args {
 				values.Set(k, v)
 			}
 			req := httptest.NewRequest("POST", "/items", strings.NewReader(values.Encode()))
 			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

 			rr := httptest.NewRecorder()
 			h.AddItem(rr, req)

 			// check response
 			if tt.wants.code != rr.Code {
 				t.Errorf("expected status code %d, got %d", tt.wants.code, rr.Code)
 			}
 			if tt.wants.code >= 400 {
 				return
 			}
			/*
 			for _, v := range tt.args {
 				if !strings.Contains(rr.Body.String(), v) {
 					t.Errorf("response body does not contain %s, got: %s", v, rr.Body.String())
 				}
 			}
			*/
 			// STEP 6-4: check inserted data
			var count int
			err := db.QueryRow("SELECT COUNT(*) FROM items WHERE name = ? AND category = ?",tt.args["name"], tt.args["category"]).Scan(&count)
			if err != nil {
				t.Fatalf("failed to query database: %v", err)
			}
			if count != 1 {
				t.Errorf("expected 1 item in database, got %d", count)
			}
 		})
 	}
 }

 func setupDB(t *testing.T) (db *sql.DB, dbPath string, closers []func(), e error) {
 	t.Helper()

 	defer func() {
 		if e != nil {
 			for _, c := range closers {
 				c()
 			}
 		}
 	}()

 	// create a temporary file for e2e testing
 	f, err := os.CreateTemp(".", "*.sqlite3")
 	if err != nil {
 		return nil, "", nil, err
 	}
 	closers = append(closers, func() {
 		f.Close()
 		os.Remove(f.Name())
 	})

 	// set up tables
 	db, err = sql.Open("sqlite3", f.Name())
 	if err != nil {
 		return nil, "", nil, err
 	}
 	closers = append(closers, func() {
 		db.Close()
 	})

 	// TODO: replace it with real SQL statements.
 	cmd := `CREATE TABLE IF NOT EXISTS items (
 		id INTEGER PRIMARY KEY AUTOINCREMENT,
 		name TEXT NOT NULL,
 		category TEXT NOT NULL
 	)`
 	_, err = db.Exec(cmd)
 	if err != nil {
 		return nil, "", nil, err
 	}

 	return db, f.Name(), closers, nil
 }
