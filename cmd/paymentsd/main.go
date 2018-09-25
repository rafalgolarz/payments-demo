package main

import (
	"net/http"

	"github.com/rafalgolarz/payments-demo/pkg/storage"

	"github.com/gin-gonic/gin"
	log "github.com/rafalgolarz/payments-demo/pkg/log/logrus"
)

// getPayments returns collection of Payments.
func getPayments(c *gin.Context) {
	paymentsDB, err := setup(paymentsStorage)

	//connect to db
	if err != nil {
		logHandler.Error("problem connecting to database", log.Fields{"dbname": paymentsStorage.Cfg.Db, "func": "getPayments"})
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Problem connecting to db"})
		return
	}
	defer paymentsDB.Close()

	payments, err := paymentsDB.GetPayments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Problem retrieving payments"})
		return
	}
	c.JSON(http.StatusOK, payments)

}

// getPaymentByID returns single Payments resource.
func getPaymentByID(c *gin.Context) {
	paymentsDB, err := setup(paymentsStorage)

	//connect to db
	if err != nil {
		logHandler.Error("problem connecting to database", log.Fields{"dbname": paymentsStorage.Cfg.Db, "func": "getPaymentsByID"})
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Problem connecting to db"})
		return
	}
	defer paymentsDB.Close()

	payments, err := paymentsDB.GetPayment(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Could not find a payment"})
		return
	}
	c.JSON(http.StatusOK, payments)

}

// addPayment adds new Payment to collection.
func addPayment(c *gin.Context) {
	paymentsDB, err := setup(paymentsStorage)

	//connect to db
	if err != nil {
		logHandler.Error("problem connecting to database", log.Fields{"dbname": paymentsStorage.Cfg.Db, "func": "addPayment"})
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Problem connecting to db"})
		return
	}
	defer paymentsDB.Close()

	var p storage.Payments
	err = c.BindJSON(&p)

	err = paymentsDB.CreatePayment(&p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not add a payment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Payment created"})
}

// updatePaymentByID updates existing Payment by UUID
func updatePaymentByID(c *gin.Context) {

	paymentsDB, err := setup(paymentsStorage)

	//connect to db
	if err != nil {
		logHandler.Error("problem connecting to database", log.Fields{"dbname": paymentsStorage.Cfg.Db, "func": "updatePaymentByID"})
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Problem connecting to db"})
		return
	}
	defer paymentsDB.Close()

	var p storage.Payments
	err = c.BindJSON(&p)

	err = paymentsDB.UpdatePayment(c.Param("id"), &p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Could not update the payment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Payment updated"})

}

// deletePaymentByID removes Payment by its UUID from collection
func deletePaymentByID(c *gin.Context) {
	paymentsDB, err := setup(paymentsStorage)

	//connect to db
	if err != nil {
		logHandler.Error("problem connecting to database", log.Fields{"dbname": paymentsStorage.Cfg.Db, "func": "deletePaymentsByID"})
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Problem connecting to db"})
		return
	}
	defer paymentsDB.Close()

	err = paymentsDB.DeletePayment(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Could not find a payment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Payment deleted"})

}

func main() {

	router := gin.Default()
	router.Use(CORSMiddleware())
	v1 := router.Group("/v1")
	{
		v1.GET("/payments", getPayments)
		v1.GET("/payments/:id", getPaymentByID)
		v1.POST("/payments", addPayment)
		v1.PUT("/payments/:id", updatePaymentByID)
		v1.DELETE("/payments/:id", deletePaymentByID)

	}
	router.Run(apiPort)

}

// CORSMiddleware added to control access
// Allow all (*) should not be used on production
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
