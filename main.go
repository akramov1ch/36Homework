package main

import (
    "database/sql"
    "log"
    "net/http"
	"fmt"

    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
    var err error
    connStr := "user=postgres password=vakhaboff dbname=shaxboz sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Database connected!")
}



type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

func getAlbums(c *gin.Context) {
    var albums []album
    rows, err := db.Query("SELECT id, title, artist, price FROM albums")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var alb album
        if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
            log.Fatal(err)
        }
        albums = append(albums, alb)
    }
    c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
    var newAlbum album
    if err := c.BindJSON(&newAlbum); err != nil {
        return
    }

    _, err := db.Exec("INSERT INTO albums (title, artist, price) VALUES ($1, $2, $3)",
        newAlbum.Title, newAlbum.Artist, newAlbum.Price)
    if err != nil {
        log.Fatal(err)
    }

    c.IndentedJSON(http.StatusCreated, newAlbum)
}

func deleteAlbums(c *gin.Context) {
    id := c.Param("id")
    if _, err := db.Exec("DELETE FROM albums WHERE id = $1", id); err != nil {
        log.Fatal(err)
    }
    c.IndentedJSON(http.StatusOK, gin.H{"message": "album deleted successfully"})
}

func updateAlbums(c *gin.Context) {
    id := c.Param("id")
    var upAlbum album
    if err := c.BindJSON(&upAlbum); err != nil {
        return
    }

    if _, err := db.Exec("UPDATE albums SET title = $1, artist = $2, price = $3 WHERE id = $4",
        upAlbum.Title, upAlbum.Artist, upAlbum.Price, id); err != nil {
        log.Fatal(err)
    }
    c.IndentedJSON(http.StatusOK, gin.H{"message": "album updated successfully"})
}

func main() {
    router := gin.Default()
    initDB()

    router.GET("/albums", getAlbums)
    router.POST("/albums", postAlbums)
    router.DELETE("/albums/:id", deleteAlbums)
    router.PUT("/albums/:id", updateAlbums)

    if err := router.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
