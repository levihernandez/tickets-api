package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"net/http"
	"github.com/gin-contrib/cors"
	"os"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

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

/*
    // Load CA certificate
    caCert, err := ioutil.ReadFile("path/to/ca.crt")
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
    	InsecureSkipVerify: false, // Do not skip server certificate verification
    }
*/
	
    // Create a new TLS configuration.
    tlsConfig := &tls.Config{
        // Set InsecureSkipVerify to true if you want to skip certificate validation.
        InsecureSkipVerify: true,
    }

    // Load the client certificate and key.
    cert, err := tls.LoadX509KeyPair(homeCertDir+"/certs/client.julian.crt", homeCertDir+"/certs/client.julian.key")
    if err != nil {
        panic(err)
    }

    // Load the CA certificate to verify the server's certificate.
    caCert, err := ioutil.ReadFile(homeCertDir+"/certs/ca.crt")
    if err != nil {
        panic(err)
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    // Add the client certificate and key to the TLS configuration.
    tlsConfig.Certificates = []tls.Certificate{cert}
    tlsConfig.RootCAs = caCertPool

	opts := &pg.Options{
		Addr:     "<host>:26257",
		User:     "",
		Password: "",
		Database: "tickets",
		PoolSize: 20,
		ApplicationName:  "gopg-crdb-app",

        // TLS settings
        TLSConfig: tlsConfig,
	}

	// db.SetMaxIdleConns(10)
	// db.SetMaxOpenConns(50)
	
	db = pg.Connect(opts)
	defer db.Close()

	// Set application name using Exec
	q , err := db.Exec("SET application_name = 'gopg-crdb-app'")
	if err != nil {
		panic(err)
		fmt.Println(q)
	}


	// Simulate pool disruption
	fmt.Println("Simulating connection pool disruption...")
	db.Close() // This will close all connections in the pool

	// Simulate some downtime
	time.Sleep(5 * time.Second)

	// "Restart" the connection pool by creating a new pool
	fmt.Println("Restarting the connection pool...")
	db = pg.Connect(opts)
	defer db.Close() // ensure closure of the new connection pool

	
	
	r := gin.Default()
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://192.168.86.202:3000"}  // Replace with your React frontend's address
    config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
    config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}

    r.Use(cors.New(config))
    r.GET("/user/:userID/purchases", getUserPurchases)
    r.GET("/user/:userID/purchases/cancellations", getUserCancelledPurchases)
    r.GET("/search/users", searchUsers)
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
        Where("purchase.user_id = ? AND purchase.status = 'cancelled'",uuidUserID).
        Select(&purchases)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, purchases)
}


func searchUsers(c *gin.Context) {
	name := c.Query("name")
	var users []User
  fmt.Println(name)
	err := db.Model(&users).
		Where("name ILIKE ?", "%"+name+"%").
		Select()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

  if users == nil {
    users = []User{}
  }
	c.JSON(http.StatusOK, users)
}
