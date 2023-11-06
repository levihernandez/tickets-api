package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Purchases []Purchase
}

type Event struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Purchases []Purchase
}

type Purchase struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"userId"`
	EventID uuid.UUID `json:"eventId"`
	Status  string    `json:"status"`
	User    *User     `json:"user"`
	Event   *Event    `json:"event"`
}

var db *pgxpool.Pool

func main() {
	
    var homeCertDir = os.Getenv("TICKETS_CERTS")
    var err error
	
    // Insecure	
    // db, err = pgxpool.Connect(context.Background(), "postgres://root@<host>:26257/tickets?application_name=pgx-crdb-app")
    // Secure
    db, err = pgxpool.Connect(context.Background(), "postgres://julian:<password>@<host>:26257/tickets?sslmode=require&sslcert="+homeCertDir+"/certs/client.julian.crt&sslkey="+homeCertDir+"/certs/client.julian.key&sslrootcert="+homeCertDir+"/certs/ca.crt&application_name=pgx-crdb-app")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	// Connection Pool manage in Go
	/*
	- if random number generator is 5%
 	- then drop existing connection manager
  	- and create a new connection
   	This is performed by K6 Scripts
 	*/
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://192.168.86.202:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}

	r.Use(cors.New(config))

	r.GET("/user/:userID/purchases", getUserPurchases)
	r.GET("/user/:userID/purchases/cancellations", getUserCancelledPurchases)
	r.GET("/search/users", searchUsers)
	r.Run(":3002")
}

func getUserPurchases(c *gin.Context) {
	userID := c.Param("userID")
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

    rows, err := db.Query(context.Background(),
    `SELECT purchase.id, purchase.user_id, purchase.event_id, purchase.status, users.id AS user__id, users.name AS user__name,
events.id AS event__id, events.name AS event__name, events.type AS event__type, events.status AS event__status
 FROM purchases AS purchase LEFT
   JOIN users AS users
    ON users.id = purchase.user_id LEFT
   JOIN events AS events
    ON events.id = purchase.event_id
   WHERE (purchase.user_id = $1)`, uuidUserID)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer rows.Close()

	var purchases []Purchase
	for rows.Next() {
		var purchase Purchase
		var user User
		var event Event
		err := rows.Scan(&purchase.ID, &purchase.UserID, &purchase.EventID, &purchase.Status,
			&user.ID, &user.Name,
			&event.ID, &event.Name, &event.Type, &event.Status)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		purchase.User = &user
		purchase.Event = &event
		purchases = append(purchases, purchase)
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

	rows, err := db.Query(context.Background(),
    `SELECT purchase.id, purchase.user_id, purchase.event_id, purchase.status, users.id AS user__id, users.name AS user__name,
events.id AS event__id, events.name AS event__name, events.type AS event__type, events.status AS event__status
 FROM purchases AS purchase LEFT
   JOIN users AS users
    ON users.id = purchase.user_id LEFT
   JOIN events AS events
    ON events.id = purchase.event_id
        WHERE (purchase.user_id = $1) AND purchase.status = 'cancelled'`, uuidUserID)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer rows.Close()

	var purchases []Purchase
	for rows.Next() {
		var purchase Purchase
		var user User
		var event Event
		err := rows.Scan(&purchase.ID, &purchase.UserID, &purchase.EventID, &purchase.Status,
			&user.ID, &user.Name,
			&event.ID, &event.Name, &event.Type, &event.Status)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		purchase.User = &user
		purchase.Event = &event
		purchases = append(purchases, purchase)
	}
	c.JSON(http.StatusOK, purchases)
}

func searchUsers(c *gin.Context) {
	name := c.Query("name")
	rows, err := db.Query(context.Background(), `SELECT * FROM users WHERE name ILIKE $1`, "%"+name+"%")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}
	if users == nil {
    users = []User{}
  }
	c.JSON(http.StatusOK, users)
}
