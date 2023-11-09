package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `pg:"type:uuid,default:gen_random_uuid()" json:"id"`
	Name      string     `json:"name"`
	Purchases []Purchase `pg:"-" json:"purchases,omitempty"`
}

type Event struct {
	ID        uuid.UUID  `pg:"type:uuid,default:gen_random_uuid()" json:"id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Status    string     `json:"status"`
	Purchases []Purchase `pg:"-" json:"purchases,omitempty"`
}

type Purchase struct {
	ID      uuid.UUID `pg:"type:uuid,default:gen_random_uuid()" json:"id"`
	UserID  uuid.UUID `pg:"type:uuid" json:"userId"`
	EventID uuid.UUID `pg:"type:uuid" json:"eventId"`
	Status  string    `json:"status"`
	User    *User     `pg:"rel:has-one" json:"user"`
	Event   *Event    `pg:"rel:has-one" json:"event"`
}

var db *pg.DB

func main() {

	var homeCertDir = os.Getenv("TICKETS_CERTS")

	caCert, err := ioutil.ReadFile(homeCertDir + "/jlevi-ca.crt")
	if err != nil {
		log.Fatalf("Unable to read CA cert: %v", err)
	}

	// Create a certificate pool and add the CA certificate
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("Failed to append CA certificate to pool")
	}

	// Set up TLS configuration
	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		InsecureSkipVerify: true, // Do not skip server certificate verification
	}

	db = pg.Connect(&pg.Options{
		Addr:            "host:26257",
		User:            "",
		Password:        "",
		Database:        "tickets",
		PoolSize:        20,
		ApplicationName: "gopg-app",

		// TLS settings
		TLSConfig: tlsConfig,
	})
	defer db.Close()

	// Set application name using Exec
	q, err := db.Exec("SET application_name = 'gopg-app'")
	if err != nil {
		panic(err)
		fmt.Println(q)
	}

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://192.168.86.202:3000"} // Replace with your React frontend's address
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}

	r.Use(cors.New(config))
	r.GET("/user/:userID/purchases", getUserPurchases)
	r.GET("/user/:userID/purchases/cancellations", getUserCancelledPurchases)
	r.GET("/search/user/:userID", searchUsers)
	r.Run(":3001") // Listen and serve on 0.0.0.0:3001
}

func getUserPurchases(c *gin.Context) {
	userID := c.Param("userID")
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var purchases []Purchase
	query := db.Model(&purchases).
		ColumnExpr("purchase.id, purchase.user_id, purchase.event_id, purchase.status, users.id AS user__id, users.name AS user__name, events.id AS event__id, events.name AS event__name, events.type AS event__type, events.status AS event__status").
		Join("LEFT JOIN users AS users ON users.id = purchase.user_id").
		Join("LEFT JOIN events AS events ON events.id = purchase.event_id").
		Where("purchase.user_id = ?", uuidUserID)
	err = query.Select(&purchases)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, purchases)
}

func getUserCancelledPurchases(c *gin.Context) {
	userID := c.Param("userID")
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var purchases []Purchase
	err = db.Model(&purchases).
		ColumnExpr("purchase.id, purchase.user_id, purchase.event_id, purchase.status, users.id AS user__id, users.name AS user__name, events.id AS event__id, events.name AS event__name, events.type AS event__type, events.status AS event__status").
		Join("LEFT JOIN users AS users ON users.id = purchase.user_id").
		Join("LEFT JOIN events AS events ON events.id = purchase.event_id").
		Where("purchase.user_id = ? AND purchase.status = 'cancelled'", uuidUserID).
		Select(&purchases)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, purchases)
}

func searchUsers(c *gin.Context) {
	// userid := c.Query("id")
	userID := c.Param("userID")
	uuidUserID, err := uuid.Parse(userID)
	var users []User
	err = db.Model(&users).
		Where("id = ?", uuidUserID).
		Select(&users)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// fmt.Println(u)
		return
	}

	if users == nil {
		users = []User{}
	}
	c.JSON(http.StatusOK, users)
}
