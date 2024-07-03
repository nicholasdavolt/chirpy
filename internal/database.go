package database

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp        `json:"chirps"`
	Users         map[int]User         `json:"users"`
	RefreshTokens map[int]RefreshToken `json:"refreshTokens"`
}

type Chirp struct {
	Id        int    `json:"id"`
	Body      string `json:"body"`
	Author_Id int    `json:"author_id"`
}

type User struct {
	Id            int    `json:"id"`
	Email         string `json:"email"`
	Password      []byte `json:"password"`
	Is_Chirpy_Red bool   `json:"is_chirpy_red"`
}

type RefreshToken struct {
	UserId      int    `json:"userId"`
	TokenString string `json:"tokenString"`
	Expiration  string `json:"expiration"`
}

func NewDB(path string) (*DB, error) {

	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()

	return db, err

}

func (db *DB) CreateUser(email string, password []byte) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1

	user := User{
		Id:            id,
		Email:         email,
		Password:      password,
		Is_Chirpy_Red: false,
	}

	for _, dbUser := range dbStructure.Users {
		if string(dbUser.Email) == string(user.Email) {
			return User{}, errors.New("User Already Exists")
		}
	}

	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) WriteRefreshToken(refreshTokenString, expiration string, id int) error {
	dbStructure, err := db.loadDB()

	if err != nil {
		return err
	}

	dbId := len(dbStructure.RefreshTokens) + 1

	refreshToken := RefreshToken{
		UserId:      id,
		TokenString: refreshTokenString,
		Expiration:  expiration,
	}

	dbStructure.RefreshTokens[dbId] = refreshToken

	err = db.writeDB(dbStructure)

	if err != nil {
		return err
	}

	return nil

}

func (db *DB) UpdateUser(idString, email string, password []byte) (User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return User{}, err
	}

	id, err := strconv.ParseInt(idString, 10, 0)

	if err != nil {
		return User{}, err
	}

	user := User{
		Id:            int(id),
		Email:         email,
		Password:      password,
		Is_Chirpy_Red: dbStructure.Users[int(id)].Is_Chirpy_Red,
	}

	dbStructure.Users[int(id)] = user

	err = db.writeDB(dbStructure)

	if err != nil {
		return User{}, err
	}

	return user, nil

}

func (db *DB) UpdateChirpyRed(id int) error {
	dbStructure, err := db.loadDB()

	if err != nil {
		return err
	}

	_, ok := dbStructure.Users[id]

	if ok {
		for _, user := range dbStructure.Users {
			if user.Id == id {
				newUser := User{
					Id:            user.Id,
					Email:         user.Email,
					Password:      user.Password,
					Is_Chirpy_Red: true,
				}

				dbStructure.Users[user.Id] = newUser
			}

		}

	} else {
		return errors.New("user not found")
	}

	err = db.writeDB(dbStructure)

	if err != nil {
		return err
	}

	return nil

}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStructure.Users))

	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) GetRefreshTokens() ([]RefreshToken, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return nil, err
	}

	tokens := make([]RefreshToken, 0, len(dbStructure.RefreshTokens))

	for _, token := range dbStructure.RefreshTokens {
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (db *DB) RevokeRefreshToken(tokenString string) error {
	dbStructure, err := db.loadDB()

	if err != nil {
		return err
	}

	id := 0

	for i, dbToken := range dbStructure.RefreshTokens {

		if dbToken.TokenString == tokenString {
			id = i
		}
	}

	token := RefreshToken{
		UserId:      0,
		Expiration:  "",
		TokenString: "",
	}

	dbStructure.RefreshTokens[id] = token

	err = db.writeDB(dbStructure)

	if err != nil {
		return err
	}

	return nil

}

func (db *DB) DeleteChirp(chirpId int) error {
	dbStructure, err := db.loadDB()

	if err != nil {
		return err
	}

	id := 0

	for i, chirp := range dbStructure.Chirps {

		if chirp.Id == chirpId {
			id = i
		}
	}

	chirp := Chirp{
		Id:        0,
		Body:      "",
		Author_Id: 0,
	}

	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)

	if err != nil {
		return errors.New("could not update db")
	}

	return nil

}

func (db *DB) CreateChirp(body string, userID int) (Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1

	chirp := Chirp{
		Id:        id,
		Body:      body,
		Author_Id: userID,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)

	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil

}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))

	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil

}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{}

	data, err := os.ReadFile(db.path)

	if err != nil {
		return dbStructure, err
	}

	err = json.Unmarshal(data, &dbStructure)

	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil

}

func (db *DB) ensureDB() error {

	_, err := os.ReadFile(db.path)

	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}

	return err
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps:        map[int]Chirp{},
		Users:         map[int]User{},
		RefreshTokens: map[int]RefreshToken{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dat, err := json.Marshal(dbStructure)

	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)

	if err != nil {
		return err
	}

	return nil

}
