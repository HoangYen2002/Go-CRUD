package controllers

import (
	"go-crud/initialize"
	"go-crud/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Models []Models

type Link struct {
	Rel  string
	Path string
}

type Meta struct {
	Description string
	Keywords    []string
}

func BookGetting(c *gin.Context) {
	var books []models.Book
	initialize.DB.Find(&books)
	c.IndentedJSON(200, gin.H{
		"books": books,
	})
}

func BookCreate(c *gin.Context) {
	var body struct {
		Title      string `json:"title"`
		AuthorId   uint   `json:"author_id"`
		CategoryId uint   `json:"category_id"`
		PostId     uint   `json:"post_id"`
	}
	c.BindJSON(&body)
	books := models.Book{Title: body.Title, AuthorId: body.AuthorId, CategoryId: body.CategoryId, PostId: body.PostId}
	result := initialize.DB.Create(&books)
	if result.Error != nil {
		c.Status(400)
		return
	}
	c.IndentedJSON(200, gin.H{
		"book":    books,
		"message": "Book created successfully",
	})
}

func BookShowByID(c *gin.Context) {
	id := c.Param("id")
	//authorId := c.Param("authorId")
	//categoryId := c.Param("categoryId")
	var book models.Book
	var author models.Author
	var category models.Category
	initialize.DB.First(&book, id)
	initialize.DB.First(&author, book.AuthorId)

	author_ID := strconv.Itoa(int(book.AuthorId))
	initialize.DB.First(&category, book.CategoryId)
	category_ID := strconv.Itoa(int(book.CategoryId))
	c.IndentedJSON(200, gin.H{
		"message": "Book found successfully",
		"error":   false,
		"data":    book,
		"links_related": []gin.H{
			{"method": "GET"},
			{
				"self_URL": "http://localhost:3000/api/book/" + id,
			},
			{
				"authors_URL": "http://localhost:3000/api/author/" + author_ID,
			},
			{
				"categories_URL": "http://localhost:3000/api/category/" + category_ID,
			},
		},
	})
}

func BookBuilder(book *models.Book, rel string) *Link {
	return &Link{
		Rel:  rel,
		Path: "/api/book_detail/" + strconv.Itoa(int(book.ID)),
	}
}

func AuthorBuilder(author *models.Author, rel string) *Link {
	return &Link{
		Rel:  rel,
		Path: "/api/author_detail/" + strconv.Itoa(int(author.ID)),
	}
}

func CategoryBuilder(category *models.Category, rel string) *Link {
	return &Link{
		Rel:  rel,
		Path: "/api/category_detail/" + strconv.Itoa(int(category.ID)),
	}
}

func GetAuthorID(author *models.Author) string {
	author_model := strconv.Itoa(int(author.ID))
	return author_model
}

func GetCategoryID(category *models.Category) string {
	category_model := strconv.Itoa(int(category.ID))
	return category_model
}

func GetBookID(book *models.Book) string {
	book_model := strconv.Itoa(int(book.ID))
	return book_model
}

// var configMap_ = map[models]string {
// 	&models.Book:     "/api/book_detail/",
// 	&models.Author:   "/api/author_detail/",
// 	&models.Category: "/api/category_detail/",
// }

// var configMap[model]string =  {
// 	model.book: "/api/book_detail/",
// }
// func buildDetailLink(model any, routeMap) string {
// 	return {
// 		path: routeMap[model] + strconv.Itoa(int(model.ID)),
// 	}
// }

var configMap = map[interface{}]string{
	&models.Book{}:     "/api/book_detail/",
	&models.Author{}:   "/api/author_detail/",
	&models.Category{}: "/api/category_detail/",
}

var routeMap = map[string]string{
	"book":     "book/",
	"author":   "author/",
	"category": "category/",
}

// func buildDetailLink_( configMap map[interface{}]string, routeMap map[string]string) string {
// 	return {
// 		path := routeMap[model] + strconv.Itoa(int(model.ID)),

// 	}
// }

func buildDetailLink(model interface{}, configMap map[interface{}]string, routeMap map[string]string, id string) string {
	return configMap[0] + id
}

// func buildDetailLink(model interface{},
// 	configMap map[interface{}]string, routeMap map[interface{}]string, id string) string {
// 	return configMap[model] + routeMap[model] + id
// }

func BookDetail(routeMap map[string]string) func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		var books []models.Book
		var authors []models.Author
		var categories []models.Category

		initialize.DB.Find(&books)
		// pp.Fprint(&books)
		initialize.DB.Find(&authors)
		initialize.DB.Find(&categories)
		// initialize.DB.First(&books, id)

		c.IndentedJSON(200, gin.H{
			"data": books,
			"Links": gin.H{
				"_Self": gin.H{
					"method": "GET",
					"self":   buildDetailLink(&models.Book{}, configMap, routeMap, id),
				},
				// "author": gin.H{
				// 	"method": "GET",
				// 	"author": buildDetailLink_(authors, routeMap),
				// },
			},
		})
	}
}

func BookDelete(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	initialize.DB.First(&book, id)
	initialize.DB.Delete(&book)
	c.IndentedJSON(200, gin.H{
		"message": "Book deleted successfully",
	})
}

func BookUpdate(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	initialize.DB.First(&book, id)
	var body struct {
		Title      string          `json:"title"`
		AuthorId   uint            `json:"author_id"`
		CategoryId uint            `json:"category_id"`
		Category   models.Category `json:"category"`
	}
	c.BindJSON(&body)
	initialize.DB.Model(&book).Updates(models.Book{Title: body.Title, AuthorId: body.AuthorId})
	c.IndentedJSON(200, gin.H{
		"message":     "Book updated successfully",
		"bookUpdated": book,
	})
}
