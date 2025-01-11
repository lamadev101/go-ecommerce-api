package router

import (
	"net/http"

	"github.com/lamadev101/ecommerce-api/constant"
	"github.com/lamadev101/ecommerce-api/controller"
)

// server status check
var statusCheckRoutes = Routes{
	Route{"Status Check", http.MethodGet, constant.STATUS_CHECK_ROUTE, controller.StatusCheck},
}

var userRoutes = Routes{
	Route{"Verify Email", http.MethodPost, constant.VERIFY_EMAIL_ROUTE, controller.VerifyEmail},
	Route{"Verify OTP", http.MethodPost, constant.VERIFY_OTP_ROUTE, controller.VerifyOtp},
	Route{"Register User", http.MethodPost, constant.USER_REGISTER_ROUTE, controller.RegisterUser},
	Route{"Login User", http.MethodPost, constant.USER_LOGIN_ROUTE, controller.LoginUser},
	Route{"Change Password", http.MethodPut, constant.CHANGE_PASSWORD_ROUTE, controller.ChangePassword},
	Route{"Reset Password", http.MethodPut, constant.RESET_PASSWORD_ROUTE, controller.ResetPassword},
}

var ecommerceRoutes = Routes{
	Route{"Create Product", http.MethodPost, constant.CREATE_PRODUCT, controller.CreateProduct},
	Route{"Create Banner", http.MethodPost, constant.CREATE_BANNER, controller.CreateBanner},
}

var ecommerceGlobalRoutes = Routes{
	Route{"List Product", http.MethodGet, constant.PRODUCTS_ROUTE, controller.ListProduct},
	Route{"Get Product By Slug", http.MethodGet, constant.PRODUCT_ROUTE, controller.GetProductBySlug},
}
