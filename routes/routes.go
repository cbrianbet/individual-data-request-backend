package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/palladiumkenya/individual-data-request-backend/controllers"
)

func Handlers(router *gin.Engine) {
	// Define routes here
	router.GET("/api_health", controllers.GetApiHealth) // Test endpoint
	router.POST("/send_mail", controllers.SendMail)     // Send test email

	router.GET("/requests", controllers.GetRequests)                      // get requests
	router.GET("/internal_approval/:id", controllers.GetInternalApproval) // get approval page data
	router.POST("/internal_approval/action", controllers.ApproverAction)  // approve or reject requests
	router.GET("/approval/:type/:id", controllers.GetApproval)            // get approval page data

	router.POST("/new_review_thread", controllers.CreateReviewThread)      // create review thread
	router.POST("/add_review", controllers.AddReview)                      // add review
	router.GET("/get_reviews/:thread_id", controllers.GetReviewsForThread) // get reviews

}
