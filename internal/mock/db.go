package mock

type Tx struct {
	//ListGreetingsFunc func(cloud.Context) ([]cloud.Greeting, error)
}

// func (tx Tx) SaveGreeting(ctx cloud.Context, greeting string) error {
// 	return nil
// }

// func (tx Tx) ListGreetings(ctx cloud.Context) ([]cloud.Greeting, error) {
// 	return tx.ListGreetingsFunc(ctx)
// }

// type Database struct{
// 	f(ctx context.Context,tx pg.Tx)
// }

// func (db Database) RunInTransaction(ctx cloud.Context, f func(cloud.Context, pg.Tx) error) error {
// 	return f(ctx, pg.Tx)
// }
