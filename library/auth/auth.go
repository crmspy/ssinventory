package auth

/*
Code By Nurul Hidayat
crmspy@gmail.com

authentication using jwt with database support
*/

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
    "github.com/dgrijalva/jwt-go"
)
const (
    mySigningKey = "WOW,MuchShibe,ToDogge"
)

func AuthError(code int, message string,c *gin.Context ) {
	resp := map[string]string{"error": message}
	c.JSON(code, resp)
	c.Abort()
  }
func Auth(c *gin.Context) {
	var myToken string = c.GetHeader("Authorization");
	if myToken != "" {
		token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		})
	
		if err == nil && token.Valid {
			c.Next()
		}else{
		AuthError(401,"wrong token",c)
		}
	}else{
		AuthError(401,"wrong token",c)
	}

}

func GetProfile(c *gin.Context) {
	//var myToken string = c.GetHeader("Authorization");
	//createdToken, err := ExampleNew([]byte(myToken))
	//ExampleParse(createdToken, mySigningKey)
	//if err != nil {
		c.JSON(200, gin.H{
			"message": "Hello see my profile",
		})
    //}
}
func GetKey(c *gin.Context) {
	createdToken, err := ExampleNew([]byte(mySigningKey))
	//ExampleParse(createdToken, mySigningKey)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "Creating TOken failed",
		})
    }else{
		c.JSON(200, gin.H{
			"message": createdToken,
		})
	}
 
}
func ExampleNew(mySigningKey []byte) (string, error) {
    // Create the token
    token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	Claims := make(jwt.MapClaims)
	Claims["foo"] = "barrac is oqs"
	Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	token.Claims = Claims
    // Sign and get the complete encoded token as a string
    tokenString, err := token.SignedString(mySigningKey)
    return tokenString, err
}

func ExampleParse(myToken string, myKey string){
    token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
        return []byte(myKey), nil
    })

    if err == nil && token.Valid {
		user := token.Claims.(jwt.MapClaims)
		fmt.Println("Your token is valid.  I like your style.",user["foo"])
    } else {
		fmt.Println("This token is terrible!  I cannot accept this.")
    }
}
