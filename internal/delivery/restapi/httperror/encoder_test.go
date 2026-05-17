package httperror

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainerrors "go-clean-api/internal/domain/errors"
	"go-clean-api/pkg/validate"
)

func TestHandleErrorResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   errorResponse
	}{
		{
			name:       "ドメインエラー_TypeNotFound",
			err:        domainerrors.TypeNotFound.New("user not found"),
			wantStatus: http.StatusNotFound,
			wantBody: errorResponse{
				Type:  domainerrors.TypeNotFound.Name(),
				Error: "user not found",
			},
		},
		{
			name:       "ドメインエラー_TypeAlreadyExists",
			err:        domainerrors.TypeAlreadyExists.New("email already exists"),
			wantStatus: http.StatusConflict,
			wantBody: errorResponse{
				Type:  domainerrors.TypeAlreadyExists.Name(),
				Error: "email already exists",
			},
		},
		{
			name:       "ドメインエラー_TypeInvalidCredentials",
			err:        domainerrors.TypeInvalidCredentials.New("invalid email or password"),
			wantStatus: http.StatusUnauthorized,
			wantBody: errorResponse{
				Type:  domainerrors.TypeInvalidCredentials.Name(),
				Error: "invalid email or password",
			},
		},
		{
			name: "ドメインエラー_TypeValidation_詳細付き",
			err: domainerrors.TypeValidation.New(
				"validation failed",
				domainerrors.NewFieldError("email", "required"),
			),
			wantStatus: http.StatusBadRequest,
			wantBody: errorResponse{
				Type:  domainerrors.TypeValidation.Name(),
				Error: "validation failed",
				Details: []domainerrors.FieldError{
					domainerrors.NewFieldError("email", "required"),
				},
			},
		},
		{
			name: "ドメインエラー_Wrap",
			err: domainerrors.TypeUnexpected.Wrap(
				errors.New("underlying"),
				"failed to save user",
			),
			wantStatus: http.StatusInternalServerError,
			wantBody: errorResponse{
				Type:  domainerrors.TypeUnexpected.Name(),
				Error: "failed to save user",
			},
		},
		{
			name: "validate_Errors",
			err: &validate.Errors{
				errors.New("email is required"),
			},
			wantStatus: http.StatusBadRequest,
			wantBody: errorResponse{
				Type:  domainerrors.TypeValidation.Name(),
				Error: msgErrorInvalidParameter,
				Details: []error{
					errors.New("email is required"),
				},
			},
		},
		{
			name:       "echo_HTTPError_原因なし",
			err:        echo.NewHTTPError(http.StatusUnprocessableEntity, "bind failed"),
			wantStatus: http.StatusUnprocessableEntity,
			wantBody: errorResponse{
				Type:  domainerrors.TypeValidation.Name(),
				Error: msgErrorInvalidParameter,
				Details: []domainerrors.FieldError{
					domainerrors.NewFieldError("", "code=422, message=bind failed"),
				},
			},
		},
		{
			name: "echo_HTTPError_原因あり",
			err: echo.NewHTTPError(http.StatusBadRequest, "bad request").
				Wrap(errors.New("invalid json")),
			wantStatus: http.StatusBadRequest,
			wantBody: errorResponse{
				Type:  domainerrors.TypeValidation.Name(),
				Error: msgErrorInvalidParameter,
				Details: []domainerrors.FieldError{
					domainerrors.NewFieldError("", "invalid json"),
				},
			},
		},
		{
			name:       "想定外エラー",
			err:        errors.New("something broke"),
			wantStatus: http.StatusInternalServerError,
			wantBody: errorResponse{
				Type:  domainerrors.TypeUnexpected.Name(),
				Error: msgErrorUnexpected,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			status, resp := handleErrorResponse(tt.err)

			assert.Equal(t, tt.wantStatus, status)
			require.NotNil(t, resp)
			assert.Equal(t, tt.wantBody.Type, resp.Type)
			assert.Equal(t, tt.wantBody.Error, resp.Error)
			assertDetailsEqual(t, tt.wantBody.Details, resp.Details)
		})
	}
}

func TestEncode_nilはnilを返す(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	c := echo.New().NewContext(req, rec)

	err := Encode(c, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Body.Bytes())
}

func TestEncode_ドメインエラーはJSONで返す(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", http.NoBody)
	c := echo.New().NewContext(req, rec)

	derr := domainerrors.TypeAlreadyExists.New("email already exists")
	err := Encode(c, derr)
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)

	var body errorResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.Equal(t, domainerrors.TypeAlreadyExists.Name(), body.Type)
	assert.Equal(t, "email already exists", body.Error)
}

func assertDetailsEqual(t *testing.T, want, got any) {
	t.Helper()

	if want == nil {
		assert.True(t, isEmptyDetails(got), "details should be nil or empty, got %#v", got)
		return
	}
	require.NotNil(t, got)

	wantJSON, err := json.Marshal(want)
	require.NoError(t, err)
	gotJSON, err := json.Marshal(got)
	require.NoError(t, err)
	assert.JSONEq(t, string(wantJSON), string(gotJSON))
}

func isEmptyDetails(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		return rv.Len() == 0
	default:
		return false
	}
}
