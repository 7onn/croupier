package user

import (
	u "croupier/pkg/utils"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//Token !
type Token struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

//Account !
type Account struct {
	gorm.Model
	Email     string `json:"email"`
	Password  string `json:"password"`
	Login     string `json:"login"`
	BirthDate string `json:"birth_date"`
	Active    bool   `json:"active"`

	Phone            string `json:"phone"`
	ConfirmationCode string `json:"-"`
	ConfirmationDate int    `json:"-"`

	Token string `json:"token";sql:"-"`
}

//Validate !
func (account *Account) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "dude, come on, it's your own money at stake; Password must be at least 6 digits length"), false
	}

	temp := &Account{}

	err := DB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use."), false
	}

	err = DB().Table("accounts").Where("login = ?", account.Login).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Login already in use."), false
	}

	err = DB().Table("accounts").Where("phone = ?", account.Phone).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Phone already in use."), false
	}

	return u.Message(false, "Requirement passed"), true
}

//Create !
func (account *Account) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	DB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_TOKEN")))
	account.Token = tokenString

	account.Password = ""

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

//Login !
func Login(email, login, password string) map[string]interface{} {

	account := &Account{}
	err := DB().Table("accounts").Where("email = ? OR login = ?", email, login).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Neither login nor email address were found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials. Please try again")
	}

	account.Password = ""

	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_TOKEN")))
	account.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

//GetUser !
func GetUser(u uint) *Account {

	acc := &Account{}
	DB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" {
		return nil
	}

	acc.Password = ""
	return acc
}
