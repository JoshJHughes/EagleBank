package users

type UserStore interface {
	Get(id UserID) (User, error)
	Put(user User) error
	Delete(id UserID) error
}

type UserService struct {
	userStore UserStore
}

func NewUserService(userStore UserStore) UserService {
	return UserService{
		userStore: userStore,
	}
}

func (svc UserService) CreateUser(req CreateUserRequest) (User, error) {
	id, err := NewRandUserID()
	if err != nil {
		return User{}, err
	}
	usr, err := NewUser(id, req.Name, req.Address, req.PhoneNumber, req.Email)
	if err != nil {
		return User{}, err
	}
	err = svc.userStore.Put(usr)
	if err != nil {
		return User{}, err
	}
	return usr, nil
}
