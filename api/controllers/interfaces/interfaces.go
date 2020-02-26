package interfaces

type Controllers struct {
	Get(*gin.Context)
	Post(*gin.Context)
}
