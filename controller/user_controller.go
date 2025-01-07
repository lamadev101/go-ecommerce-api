package controller

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lamadev101/ecommerce-api/auth"
	"github.com/lamadev101/ecommerce-api/constant"
	"github.com/lamadev101/ecommerce-api/database"
	"github.com/lamadev101/ecommerce-api/types"
	"github.com/lamadev101/ecommerce-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func VerifyEmail(c *gin.Context) {
	// Verify email
	var req types.UserVerification
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EMAIL_VALIDATION_FAILED})
		return
	}
	res := database.Mgr.GetUserByEmail(req.Email, constant.USER_VERIFICATION_COLLECTION)
	fmt.Println("res", res)

	if res.Otp != 0 {
		sec := res.CreatedAt + constant.OTP_VALIDATION_TIME

		// Check if OTP is expired
		if sec < time.Now().Unix() {
			req, checkEmail := utils.SendEmail(req)
			if checkEmail != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": checkEmail.Error()})
				return
			}
			req.CreatedAt = time.Now().Unix()
			// update the OTP in the database
			database.Mgr.UpateUserOTP(req, constant.USER_VERIFICATION_COLLECTION)
			c.JSON(http.StatusOK, gin.H{"error": false, "message": "OTP sent successfully !!"})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OTP_ALREADY_SENT})
			return
		}
	}

	req, checkEmail := utils.SendEmail(req)

	if checkEmail != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": checkEmail.Error()})
		return
	}

	req.CreatedAt = time.Now().Unix()
	// insert the OTP in the database
	_, err := database.Mgr.Insert(req, constant.USER_VERIFICATION_COLLECTION)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "OTP sent successfully -:)"})
}

func VerifyOtp(c *gin.Context) {
	var req types.UserVerification

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EMAIL_VALIDATION_FAILED})
		return
	}
	if req.Otp <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OTP_VALIDATION_FAILED})
		return
	}

	res := database.Mgr.GetUserByEmail(req.Email, constant.USER_VERIFICATION_COLLECTION)
	// if status or email is already verified
	if res.Status {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "This user is already verified!"})
		return
	}

	sec := res.CreatedAt + constant.OTP_VALIDATION_TIME

	if res.Otp != req.Otp {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OTP_VALIDATION_FAILED})
		return
	}
	// otp expired
	if sec < time.Now().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OTP_EXPIRED})
		return
	}

	// if all good then we will verified the email
	req.Status = true
	req.CreatedAt = time.Now().Unix()
	err := database.Mgr.UpdateEmailVerifiedStatus(req, constant.USER_VERIFICATION_COLLECTION)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OTP_VALIDATION_FAILED})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success"})
}

func RegisterUser(c *gin.Context) {
	var userReq types.UserClient
	var dbUser types.User

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	// payload error handle
	if err := utils.RegisterUserValidation(userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	// check if email is verified or not
	if res := database.Mgr.GetUserByEmail(userReq.Email, constant.USER_VERIFICATION_COLLECTION); !res.Status {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EMAIL_NOT_VERIFIED})
		return
	}

	// Check duplicate user
	if res := database.Mgr.GetUserByEmail(userReq.Email, constant.USERS_COLLECTION); res.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.ALREADY_REGISTERED_WITH_EMAIL})
		return
	}

	// value assign
	dbUser.Name = userReq.Name
	dbUser.Email = userReq.Email
	dbUser.UserType = constant.USER_ROLE
	dbUser.Password = utils.GenerateHashPassword(userReq.Password)
	dbUser.CreatedAt = time.Now().Unix()
	dbUser.UpdatedAt = time.Now().Unix()

	id, err := database.Mgr.Insert(dbUser, constant.USERS_COLLECTION)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}

	// Jwt struct prepare
	jwtWrapper := auth.JwtWrapper{
		SecretKey:      os.Getenv("JWT_SECRET_KEY"),
		Issuer:         os.Getenv("JWT_ISSUER"),
		ExpirationTime: 48,
	}

	userId := id.(primitive.ObjectID)

	// token generation
	token, err := jwtWrapper.GenerateToken(userId, userReq.Email, dbUser.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}
	dbUser.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "success",
		"data":    dbUser,
		"token":   token,
	})

}

func LoginUser(c *gin.Context) {
	var loginReq types.Login

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	user := database.Mgr.GetSingleRecordByEmailForUser(loginReq.Email, constant.USERS_COLLECTION)
	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.NOT_REGISTERED_USER})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.CREADENTIAL_DOES_NOT_MATCH})
		return
	}

	jwtWrapper := auth.JwtWrapper{
		SecretKey:      os.Getenv("JWT_SECRET_KEY"),
		Issuer:         os.Getenv("JWT_ISSUER"),
		ExpirationTime: 48,
	}

	token, err := jwtWrapper.GenerateToken(user.Id, user.Email, user.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func ChangePassword(c *gin.Context) {
	var changePasswordReq types.ChangePassword

	if err := c.ShouldBindJSON(&changePasswordReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	if changePasswordReq.OldPassword == changePasswordReq.NewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.ALREADY_USED_PASSWORD})
		return
	}

	if changePasswordReq.ConfirmPassword != changePasswordReq.NewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.CONFIRM_PASSOWRD_MATCH})
		return
	}

	user := database.Mgr.GetSingleRecordByEmailForUser(changePasswordReq.Email, constant.USERS_COLLECTION)

	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EMAIL_NOT_VERIFIED})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePasswordReq.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.PREV_PASSWORD_DOES_NOT_MATCH})
		return
	}

	if err := utils.ChangePasswordValidation(changePasswordReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	// Assign new password
	hashPassword := utils.GenerateHashPassword(changePasswordReq.NewPassword)

	if err := database.Mgr.UpdateByEmail(changePasswordReq.Email, hashPassword, constant.USERS_COLLECTION); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": constant.PASSWORD_CHANGE_SUCCESSFULLY})
}

func ResetPassword(c *gin.Context) {
	var resetPasswordReq types.Login

	if err := c.ShouldBindJSON(&resetPasswordReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	user := database.Mgr.GetSingleRecordByEmailForUser(resetPasswordReq.Email, constant.USERS_COLLECTION)

	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EMAIL_NOT_VERIFIED})
		return
	}

	if resetPasswordReq.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Password is required!"})
		return
	}

	hashPassword := utils.GenerateHashPassword(resetPasswordReq.Password)

	if err := database.Mgr.UpdateByEmail(resetPasswordReq.Email, hashPassword, constant.USERS_COLLECTION); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": constant.PASSWORD_RESET_SUCCESSFULLY})
}

func TwoFactorVerification(c *gin.Context) {

}
