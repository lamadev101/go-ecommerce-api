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
		utils.RespondWithBadRequestError(c, err.Error())
		return
	}

	if req.Email == "" {
		utils.RespondWithBadRequestError(c, constant.EMAIL_VALIDATION_FAILED)
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
				utils.RespondWithBadRequestError(c, checkEmail.Error())
				return
			}
			req.CreatedAt = time.Now().Unix()
			// update the OTP in the database
			database.Mgr.UpateUserOTP(req, constant.USER_VERIFICATION_COLLECTION)
			utils.RespondWithSuccessMsg(c, "OTP sent successfully")
			return
		} else {
			utils.RespondWithBadRequestError(c, constant.OTP_ALREADY_SENT)
			return
		}
	}

	req, checkEmail := utils.SendEmail(req)

	if checkEmail != nil {
		utils.RespondWithBadRequestError(c, checkEmail.Error())
		return
	}

	req.CreatedAt = time.Now().Unix()
	// insert the OTP in the database
	_, err := database.Mgr.Insert(req, constant.USER_VERIFICATION_COLLECTION)
	if err != nil {
		utils.RespondWithBadRequestError(c, err.Error())
		return
	}
	utils.RespondWithSuccessMsg(c, "OTP sent successfully")
}

func VerifyOtp(c *gin.Context) {
	var req types.UserVerification

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithBadRequestError(c, err.Error())
		return
	}

	if req.Email == "" {
		utils.RespondWithBadRequestError(c, constant.EMAIL_VALIDATION_FAILED)
		return
	}
	if req.Otp <= 0 {
		utils.RespondWithBadRequestError(c, constant.OTP_VALIDATION_FAILED)
		return
	}

	res := database.Mgr.GetUserByEmail(req.Email, constant.USER_VERIFICATION_COLLECTION)
	// if status or email is already verified
	if res.Status {
		utils.RespondWithBadRequestError(c, "This user is already verified")
		return
	}

	sec := res.CreatedAt + constant.OTP_VALIDATION_TIME

	if res.Otp != req.Otp {
		utils.RespondWithBadRequestError(c, constant.OTP_VALIDATION_FAILED)
		return
	}
	// otp expired
	if sec < time.Now().Unix() {
		utils.RespondWithBadRequestError(c, constant.OTP_EXPIRED)
		return
	}

	// if all good then we will verified the email
	req.Status = true
	req.CreatedAt = time.Now().Unix()
	err := database.Mgr.UpdateEmailVerifiedStatus(req, constant.USER_VERIFICATION_COLLECTION)

	if err != nil {
		utils.RespondWithBadRequestError(c, constant.OTP_VALIDATION_FAILED)
		return
	}
	utils.RespondWithSuccessMsg(c, "OTP verified successfully")
}

func RegisterUser(c *gin.Context) {
	var userReq types.UserClient
	var dbUser types.User

	if err := c.ShouldBindJSON(&userReq); err != nil {
		utils.RespondWithBadRequestError(c, err.Error())
		return
	}

	// payload error handle
	if err := utils.RegisterUserValidation(userReq); err != nil {
		utils.RespondWithBadRequestError(c, err.Error())
		return
	}

	// check if email is verified or not
	if res := database.Mgr.GetUserByEmail(userReq.Email, constant.USER_VERIFICATION_COLLECTION); !res.Status {
		utils.RespondWithBadRequestError(c, constant.EMAIL_NOT_VERIFIED)
		return
	}

	// Check duplicate user
	if res := database.Mgr.GetUserByEmail(userReq.Email, constant.USERS_COLLECTION); res.Email != "" {
		utils.RespondWithBadRequestError(c, constant.ALREADY_REGISTERED_WITH_EMAIL)
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
		utils.RespondWithInternalServerError(c, err.Error())
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
		utils.RespondWithInternalServerError(c, err.Error())
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
		utils.RespondWithBadRequestError(c, err.Error())
		return
	}

	user := database.Mgr.GetSingleRecordByEmailForUser(loginReq.Email, constant.USERS_COLLECTION)
	if user.Email == "" {
		utils.RespondWithBadRequestError(c, constant.NOT_REGISTERED_USER)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		utils.RespondWithBadRequestError(c, constant.CREADENTIAL_DOES_NOT_MATCH)
		return
	}

	jwtWrapper := auth.JwtWrapper{
		SecretKey:      os.Getenv("JWT_SECRET_KEY"),
		Issuer:         os.Getenv("JWT_ISSUER"),
		ExpirationTime: 48,
	}

	token, err := jwtWrapper.GenerateToken(user.Id, user.Email, user.UserType)
	if err != nil {
		utils.RespondWithInternalServerError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func ChangePassword(c *gin.Context) {
	var changePasswordReq types.ChangePassword

	if err := c.ShouldBindJSON(&changePasswordReq); err != nil {
		utils.RespondWithBadRequestError(c, err.Error())
		return
	}

	if changePasswordReq.OldPassword == changePasswordReq.NewPassword {
		utils.RespondWithBadRequestError(c, constant.ALREADY_USED_PASSWORD)
		return
	}

	if changePasswordReq.ConfirmPassword != changePasswordReq.NewPassword {
		utils.RespondWithBadRequestError(c, constant.CONFIRM_PASSOWRD_MATCH)
		return
	}

	user := database.Mgr.GetSingleRecordByEmailForUser(changePasswordReq.Email, constant.USERS_COLLECTION)

	if user.Email == "" {
		utils.RespondWithBadRequestError(c, constant.EMAIL_NOT_VERIFIED)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePasswordReq.OldPassword)); err != nil {
		utils.RespondWithBadRequestError(c, constant.PREV_PASSWORD_DOES_NOT_MATCH)
		return
	}

	if err := utils.ChangePasswordValidation(changePasswordReq); err != nil {
		utils.RespondWithBadRequestError(c, err.Error())
		return
	}

	// Assign new password
	hashPassword := utils.GenerateHashPassword(changePasswordReq.NewPassword)

	if err := database.Mgr.UpdateByEmail(changePasswordReq.Email, hashPassword, constant.USERS_COLLECTION); err != nil {
		utils.RespondWithInternalServerError(c, err.Error())
		return
	}
	utils.RespondWithSuccessMsg(c, constant.PASSWORD_CHANGE_SUCCESSFULLY)
}

func ResetPassword(c *gin.Context) {
	var resetPasswordReq types.Login

	if err := c.ShouldBindJSON(&resetPasswordReq); err != nil {
		utils.RespondWithBadRequestError(c, err.Error())
		return
	}

	user := database.Mgr.GetSingleRecordByEmailForUser(resetPasswordReq.Email, constant.USERS_COLLECTION)

	if user.Email == "" {
		utils.RespondWithBadRequestError(c, constant.EMAIL_NOT_VERIFIED)
		return
	}

	if resetPasswordReq.Password == "" {
		utils.RespondWithBadRequestError(c, "Password is required")
		return
	}

	hashPassword := utils.GenerateHashPassword(resetPasswordReq.Password)

	if err := database.Mgr.UpdateByEmail(resetPasswordReq.Email, hashPassword, constant.USERS_COLLECTION); err != nil {
		utils.RespondWithInternalServerError(c, err.Error())
		return
	}
	utils.RespondWithSuccessMsg(c, constant.PASSWORD_RESET_SUCCESSFULLY)
}

func TwoFactorVerification(c *gin.Context) {

}
