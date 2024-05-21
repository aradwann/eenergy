package db

// import (
// 	"context"
// 	"testing"

// 	"github.com/aradwann/eenergy/util"
// )

// func TestVerifyEmailTx(t *testing.T) {

// 	n := 5

// 	errs := make(chan error)
// 	results := make(chan VerifyEmailTxResult)

// 	// run n concurrent transfer transaction
// 	for i := 0; i < n; i++ {
// 		go func() {

// 			result, err := testStore.VerifyEmailTx(context.Background(), VerifyEmailTxParams{
// 				EmailId:    util.RandomInt(1, 29),
// 				SecretCode: util.RandomString(16),
// 			})

// 			errs <- err
// 			results <- result
// 		}()
// 	}

// 	for i := 0; i < n; i++ {
// 		// TODO: mock  and verify

// 		// err := <-errs
// 		// require.NoError(t, err)

// 		// result := <-results
// 		// require.NotEmpty(t, result)

// 	}

// }
