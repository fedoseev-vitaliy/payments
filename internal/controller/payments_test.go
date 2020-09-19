package controller

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/fedoseev-vitaliy/payments/internal/mocks"
)

//go:generate mockery -case=snake -dir=./../providers -outpkg=mocks -output=../mocks -name=.*Provider -recursive

func TestController_GetPaymentsURL(t *testing.T) {
	t.Parallel()

	aMock := &mocks.Provider{}
	gMock := &mocks.Provider{}

	c := New(aMock, gMock)
	productID := "testProduct"
	gURL := fmt.Sprintf("http://google.pay.com/payfor?product=%s", productID)
	aURL := fmt.Sprintf("http://apple.pay.com/payfor?product=%s", productID)

	aMock.On("GetPayURL", mock.Anything, productID).Return(aURL, nil).Once()
	gMock.On("GetPayURL", mock.Anything, productID).Return(gURL, nil).Once()

	urls, err := c.GetPaymentsURL(context.Background(), productID)
	require.NoError(t, err)
	require.Equal(t, urls.GPayURL, gURL)
	require.Equal(t, urls.APayURL, aURL)

	mock.AssertExpectationsForObjects(t, aMock, gMock)
}

func TestController_GetPaymentsURLNegative(t *testing.T) {
	t.Parallel()

	aMock := &mocks.Provider{}
	gMock := &mocks.Provider{}

	c := New(aMock, gMock)
	productID := "testProduct"
	gURL := fmt.Sprintf("http://google.pay.com/payfor?product=%s", productID)
	aURL := fmt.Sprintf("http://apple.pay.com/payfor?product=%s", productID)

	t.Run("aPay failed", func(t *testing.T) {
		aMock.On("GetPayURL", mock.Anything, productID).Return("", errors.New("opps apple failed")).Once()
		gMock.On("GetPayURL", mock.Anything, productID).Return(gURL, nil).Once()

		urls, err := c.GetPaymentsURL(context.Background(), productID)
		require.Error(t, err)
		require.Nil(t, urls)
	})

	t.Run("gPay failed", func(t *testing.T) {
		gMock.On("GetPayURL", mock.Anything, productID).Return("", errors.New("opps google failed")).Once()
		aMock.On("GetPayURL", mock.Anything, productID).Return(aURL, nil).Once()

		urls, err := c.GetPaymentsURL(context.Background(), productID)
		require.Error(t, err)
		require.Nil(t, urls)
	})

	mock.AssertExpectationsForObjects(t, aMock, gMock)
}
