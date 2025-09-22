package listing_test

import (
	"errors"
	"scyllaDbAssignment/internal/cloud"
	"scyllaDbAssignment/internal/listing"
	"scyllaDbAssignment/pkg"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRun_FiltersCorrectly(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := cloud.NewMockClientInterface(ctrl)

	mockClient.EXPECT().
		ListVersions().
		Return([]pkg.Version{
			{Name: "1.0.0", CloudState: "ENABLED"},
			{Name: "2.0.0", CloudState: "DISABLED"},
			{Name: "3.0.0", CloudState: "ANY"},
		}, nil)

	params := listing.ListParams{
		GT: "1.0.0",
		LT: "3.0.0",
	}

	got, err := listing.Run(mockClient, params)

	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "2.0.0", got[0].Name)
	assert.Equal(t, "DISABLED", got[0].CloudState)
}

func TestRun_ClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	returnedError := errors.New("boom")

	mockClient := cloud.NewMockClientInterface(ctrl)

	mockClient.EXPECT().
		ListVersions().
		Return(nil, returnedError)

	params := listing.ListParams{}

	got, err := listing.Run(mockClient, params)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assert.Nil(t, got)
	assert.Equal(t, returnedError, err)
}
