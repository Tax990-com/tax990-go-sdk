package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apierrors "github.com/tax990/sdk-go/tax990/errors"
	httpclient "github.com/tax990/sdk-go/tax990/http"
	"github.com/tax990/sdk-go/tax990/resources"
	"github.com/tax990/sdk-go/tests/fixtures"
)

func newTestForm990N(handler http.HandlerFunc) (resources.Form990NResource, *httptest.Server) {
	server := httptest.NewServer(handler)
	c := httpclient.NewClient(httpclient.Options{BaseURL: server.URL})
	return resources.NewForm990N(c), server
}

func TestForm990N_Create_Success(t *testing.T) {
	resource, server := newTestForm990N(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/form990n/create", r.URL.Path)
		assert.NotEmpty(t, r.Header.Get("idempotency-key"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Contains(t, body, "Form990NRecords")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fixtures.CreateSuccessResponseBody))
	})
	defer server.Close()

	resp, err := resource.Create(context.Background(), fixtures.SampleCreatePayload(), "")

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	require.NotNil(t, resp.Form990NRecords)
	require.Len(t, resp.Form990NRecords.SuccessRecords, 1)
	assert.Equal(t, "Created", resp.Form990NRecords.SuccessRecords[0].RecordStatus)
}

func TestForm990N_Create_UsesProvidedIdempotencyKey(t *testing.T) {
	const customKey = "my-custom-idempotency-key"
	resource, server := newTestForm990N(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, customKey, r.Header.Get("idempotency-key"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fixtures.CreateSuccessResponseBody))
	})
	defer server.Close()

	_, err := resource.Create(context.Background(), fixtures.SampleCreatePayload(), customKey)
	require.NoError(t, err)
}

func TestForm990N_Submit_IsAliasForCreate(t *testing.T) {
	called := false
	resource, server := newTestForm990N(func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "/v1/form990n/create", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fixtures.CreateSuccessResponseBody))
	})
	defer server.Close()

	_, err := resource.Submit(context.Background(), fixtures.SampleCreatePayload(), "")
	require.NoError(t, err)
	assert.True(t, called)
}

func TestForm990N_Get(t *testing.T) {
	resource, server := newTestForm990N(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/form990n/get", r.URL.Path)
		assert.Equal(t, fixtures.SubmissionID, r.URL.Query().Get("SubmissionId"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"StatusCode":200,"StatusNm":"OK","StatusMessage":"OK","CorrelationId":"c1","SubmissionId":"` + fixtures.SubmissionID + `","Form990NRecords":{"SuccessRecords":[],"ErrorRecords":null},"Errors":null}`))
	})
	defer server.Close()

	resp, err := resource.Get(context.Background(), fixtures.SubmissionID, "")
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestForm990N_Get_PassesRecordID(t *testing.T) {
	resource, server := newTestForm990N(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, fixtures.RecordID, r.URL.Query().Get("RecordId"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"StatusCode":200,"StatusNm":"OK","StatusMessage":"OK","CorrelationId":"c2","SubmissionId":null,"Form990NRecords":null,"Errors":null}`))
	})
	defer server.Close()

	_, err := resource.Get(context.Background(), fixtures.SubmissionID, fixtures.RecordID)
	require.NoError(t, err)
}

func TestForm990N_List(t *testing.T) {
	resource, server := newTestForm990N(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/form990n/list", r.URL.Path)
		assert.Equal(t, fixtures.BusinessID, r.URL.Query().Get("BusinessId"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"StatusCode":200,"StatusNm":"OK","StatusMessage":"OK","CorrelationId":"c3","SubmissionId":null,"Form990NRecords":{"SuccessRecords":[],"ErrorRecords":null},"Errors":null}`))
	})
	defer server.Close()

	resp, err := resource.List(context.Background(), "", fixtures.BusinessID)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestForm990N_Delete(t *testing.T) {
	resource, server := newTestForm990N(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v1/form990n/delete", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"StatusCode":200,"StatusNm":"OK","StatusMessage":"Deleted.","CorrelationId":"c4","SubmissionId":null,"Form990NRecords":null,"Errors":null}`))
	})
	defer server.Close()

	resp, err := resource.Delete(context.Background(), fixtures.SubmissionID, "")
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestForm990N_Validate(t *testing.T) {
	resource, server := newTestForm990N(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/form990n/validate", r.URL.Path)
		assert.Equal(t, fixtures.RecordID, r.URL.Query().Get("RecordIds"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"StatusCode":200,"StatusNm":"OK","StatusMessage":"Valid.","CorrelationId":"c5","SubmissionId":null,"Form990NRecords":{"SuccessRecords":[],"ErrorRecords":null},"Errors":null}`))
	})
	defer server.Close()

	resp, err := resource.Validate(context.Background(), fixtures.SubmissionID, []string{fixtures.RecordID})
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestForm990N_Transmit(t *testing.T) {
	resource, server := newTestForm990N(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/form990n/transmit", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"StatusCode":200,"StatusNm":"OK","StatusMessage":"Transmitted.","CorrelationId":"c6","SubmissionId":null,"Form990NRecords":{"SuccessRecords":[],"ErrorRecords":null},"Errors":null}`))
	})
	defer server.Close()

	resp, err := resource.Transmit(context.Background(), &fixtures.SampleTransmitPayload)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestForm990N_Status(t *testing.T) {
	resource, server := newTestForm990N(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/form990n/status", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"StatusCode":200,"StatusNm":"OK","StatusMessage":"OK.","CorrelationId":"c7","SubmissionId":null,"Form990NRecords":{"SuccessRecords":[],"ErrorRecords":null},"Errors":null}`))
	})
	defer server.Close()

	resp, err := resource.Status(context.Background(), fixtures.SubmissionID, nil)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestForm990N_Create_ValidationError(t *testing.T) {
	resource, server := newTestForm990N(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fixtures.ValidationErrorResponseBody))
	})
	defer server.Close()

	_, err := resource.Create(context.Background(), fixtures.SampleCreatePayload(), "")

	require.Error(t, err)
	var valErr *apierrors.ValidationError
	assert.ErrorAs(t, err, &valErr)
	assert.Equal(t, 400, valErr.StatusCode)
}

func TestForm990N_Get_NotFound(t *testing.T) {
	resource, server := newTestForm990N(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fixtures.NotFoundResponseBody))
	})
	defer server.Close()

	_, err := resource.Get(context.Background(), "nonexistent-id", "")

	require.Error(t, err)
	var notFound *apierrors.NotFoundError
	assert.ErrorAs(t, err, &notFound)
	assert.Equal(t, 404, notFound.StatusCode)
}
