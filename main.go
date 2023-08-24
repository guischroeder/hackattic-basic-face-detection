package main

import  (
   "net/http"

   "github.com/gin-gonic/gin"
)

func main() {
   router := gin.Default()
   router.GET("/hello", func (c *gin.Context) {
      c.IndentedJSON(http.StatusOK, "Hello")
   })

   router.Run()
}

