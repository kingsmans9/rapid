package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spectrocloud/rapid-agent/pkg/store"
)

func Test_site_ListSites(t *testing.T) {
	req, err := http.NewRequest("GET", "/sites/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		store store.Store
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantStatus   int
		wantResponse string
	}{
		{
			name: "expect 200",
			fields: fields{
				store: &mockStore{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: req,
			},
			wantStatus:   http.StatusOK,
			wantResponse: `{"site_id":"site1","site_name":"site1","email_address":"site1@site.com","site_description":"site1","isv_site_id":"site1"}`,
		}, {
			name: "expect 500",
			fields: fields{
				store: &mockErrStore{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: req,
			},
			wantStatus:   http.StatusInternalServerError,
			wantResponse: `{"error":"error listing sites"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &site{
				store: tt.fields.store,
			}
			s.ListSites(tt.args.w, tt.args.r)

			if tt.wantStatus != tt.args.w.(*httptest.ResponseRecorder).Code {
				t.Errorf("want status %d, got %d", tt.wantStatus, tt.args.w.(*httptest.ResponseRecorder).Code)
			}
		})
	}
}
