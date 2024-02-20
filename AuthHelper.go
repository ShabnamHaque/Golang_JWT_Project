package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("Unauthorised to access this resource!")
	}
	return err
}
func MatchUserTypetoUid(c *gin.Context, userId string) (err error) {

	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("Unauthorised to access this resource.")
		return err
		// any other user other than him/hersefl not allowed
	}
	//admin and user itself allowed.
	err = CheckUserType(c, userType)
	return err
}
