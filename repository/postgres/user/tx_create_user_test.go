package user

// import (
// 	"context"
// 	"testing"

// 	"github.com/aradwann/eenergy/util"
// 	"github.com/stretchr/testify/require"
// )

// func TestCreateUserTx(t *testing.T) {

// 	n := 5

// 	errs := make(chan error)
// 	results := make(chan CreateUserTxResult)

// 	// run n concurrent transfer transaction
// 	for i := 0; i < n; i++ {
// 		go func() {
// 			result, err := testStore.CreateUserTx(context.Background(), CreateUserTxParams{
// 				CreateUserParams: CreateUserParams{
// 					Username:       util.RandomOwner(),
// 					HashedPassword: util.RandomString(12),
// 					FullName:       util.RandomString(8),
// 					Email:          util.RandomEmail(),
// 				},

// 				AfterCreate: func(user User) error {

// 					return nil
// 				},
// 			})

// 			errs <- err
// 			results <- result
// 		}()
// 	}

// 	for i := 0; i < n; i++ {
// 		err := <-errs
// 		require.NoError(t, err)

// 		result := <-results
// 		require.NotEmpty(t, result)

// 	}

// }
