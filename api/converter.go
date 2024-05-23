package api

// import (
// 	"github.com/aradwann/eenergy/pb"
// 	db "github.com/aradwann/eenergy/repository/store"
// 	"google.golang.org/protobuf/types/known/timestamppb"
// )

// func convertUser(user db.User) *pb.User {
// 	return &pb.User{
// 		Username:          user.Username,
// 		FullName:          user.FullName,
// 		Email:             user.Email,
// 		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
// 		CreatedAt:         timestamppb.New(user.CreatedAt),
// 	}
// }
// func convertTransfer(transfer db.Transfer) *pb.Transfer {
// 	return &pb.Transfer{
// 		Id:            transfer.ID,
// 		FromAccountId: transfer.FromAccountID,
// 		ToAccountId:   transfer.ToAccountID,
// 		Amount:        transfer.Amount,
// 		CreatedAt:     timestamppb.New(transfer.CreatedAt),
// 	}
// }

// func convertAccount(account db.Account) *pb.Account {
// 	return &pb.Account{
// 		Id:        account.ID,
// 		Owner:     account.Owner,
// 		Balance:   account.Balance,
// 		Unit:      account.Unit,
// 		CreatedAt: timestamppb.New(account.CreatedAt),
// 	}
// }

// func convertAccounts(accounts []db.Account) []*pb.Account {
// 	var accs []*pb.Account

// 	for _, acc := range accounts {
// 		accs = append(accs, convertAccount(acc))
// 	}
// 	return accs
// }

// func convertEntries(entries []db.Entry) []*pb.Entry {
// 	var ents []*pb.Entry

// 	for _, entry := range entries {
// 		ents = append(ents, convertEntry(entry))
// 	}
// 	return ents
// }

// func convertEntry(entry db.Entry) *pb.Entry {
// 	return &pb.Entry{
// 		Id:        entry.ID,
// 		AccountId: entry.AccountID,
// 		Amount:    entry.Amount,
// 		CreatedAt: timestamppb.New(entry.CreatedAt),
// 	}
// }
