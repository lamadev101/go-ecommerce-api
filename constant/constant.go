package constant

const (
	MDBURI   = "localhost:27017"
	DATABASE = "ecommerce"

	EMAIL_SENDER        = "ghisingkarma740@gmail.com"
	API_VERSION         = "v1"
	PORT                = "8080"
	OTP_VALIDATION_TIME = 60

	STATUS_CHECK_ROUTE = "/server-status"
)

// User roles
const (
	ADMIN_ROLE = "admin"
	USER_ROLE  = "user"
)

// MongoDB Collections
const (
	USERS_COLLECTION             = "users"
	ORDER_COLLECTION             = "order"
	ADDRESS_COLLECTION           = "address"
	PRODUCTS_COLLECTION          = "products"
	USER_VERIFICATION_COLLECTION = "user_verification"
)

// Routes for the API
const (
	// email verification routes
	VERIFY_EMAIL_ROUTE = "/verify-email"
	VERIFY_OTP_ROUTE   = "/verify-otp"
	SEND_EMAIL_ROUTE   = "/send-email"

	// user routes
	USER_REGISTER_ROUTE   = "/register"
	USER_LOGIN_ROUTE      = "/login"
	CHANGE_PASSWORD_ROUTE = "/change-password"
	RESET_PASSWORD_ROUTE  = "/reset-password"

	// Product routes
	PRODUCTS_ROUTE = "/products"
	CREATE_PRODUCT = "/create-product"
	UPDATE_PRODUCT = "/update-product"
	DELETE_PRODUCT = "/delete-product"

	// Address routes
	ADDRESS_ROUTE  = "/address"
	CREATE_ADDRESS = "/create-address"
	UPDATE_ADDRESS = "/update-address"
	DELETE_ADDRESS = "/delete-address"
)

// Http status messages
const (
	ALREADY_REGISTERED_WITH_EMAIL = "User already registered with this email"
	EMAIL_NOT_VERIFIED            = "Email not verified please verify your email first"
	EMAIL_VALIDATION_FAILED       = "Email validation failed"
	OTP_VALIDATION_FAILED         = "OTP wrong passed"
	OTP_EXPIRED                   = "OTP expired"
	OTP_ALREADY_SENT              = "OTP already sent to this email"
	NOT_REGISTERED_USER           = "Your are not registered user"
	CREADENTIAL_DOES_NOT_MATCH    = "Email or Password is incorrect"
	PREV_PASSWORD_DOES_NOT_MATCH  = "Old password is does not match with previous one"
	ALREADY_USED_PASSWORD         = "This password is already used. Please try new one"
	PASSWORD_CHANGE_SUCCESSFULLY  = "Password changed successfully"
	PASSWORD_RESET_SUCCESSFULLY   = "Password reset successfully"
	CONFIRM_PASSOWRD_MATCH        = "New password and confirm password should be same"
)
