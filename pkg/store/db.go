package store

type TinyQuizStore interface {
	//UserStore
	Migrate() error
	Transaction(func(TinyQuizStore) error) error
}

// UserStore is an interface that defines the methods that must be implemented by a user store.
//type UserStore interface {
//	// CreateUser creates a new user.
//	CreateUser(ctx context.Context, user *entity.User) error
//	// GetUserByEmail returns the user with the given email.
//	GetUserByEmail(email string) (*entity.User, error)
//	// GetUserByID returns the user with the given ID.
//	GetUserByID(id string) (*entity.User, error)
//	// UpdateUser updates the user with the given ID.
//	UpdateUser(id string, user *entity.User) error
//	// DeleteUser deletes the user with the given ID.
//	DeleteUser(id string) error
//}
