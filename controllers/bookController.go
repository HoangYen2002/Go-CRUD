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

		// "data": []gin.H{
		// 	{
		// 		"book":          book,
		// 		"author_info":   author,
		// 		"category_info": category,
		// 	},
		// },

	})
}

// func LinkBuilder(c *gin.Context) {
// 	var book []models.Book
// 	initialize.DB.Find(&book)
// 	c.IndentedJSON(200, gin.H{
// 		"message": "Book found successfully",
// 		"data":    book,
// 	})
// }

func BookBuilder(book *models.Book, rel string) *Link {
	return &Link{

		Rel:  rel,
		Path: "/api/book/" + strconv.Itoa(int(book.ID)),
	}
}

func AuthorBuilder(author *models.Author, rel string) *Link {
	return &Link{
		Rel:  rel,
		Path: "/api/author/" + strconv.Itoa(int(author.ID)),
	}
}

func CategoryBuilder(category *models.Category, rel string) *Link {
	return &Link{
		Rel:  rel,
		Path: "/api/category/" + strconv.Itoa(int(category.ID)),
	}
}

func SpecifyModelInfo(path_name string, author *models.Author) string {
	return path_name + strconv.Itoa(int(author.ID))
}

func BuildMeta() *Meta {
	return &Meta{
		//Description: "This is a description",
		Keywords: []string{
			"book",
			"author",
			"category",
		},
	}
}

func BookDetail(c *gin.Context) {
	id := c.Param("id")
	var book []models.Book
	var author []models.Author
	var category []models.Category

	initialize.DB.Find(&book)
	initialize.DB.Find(&author)
	initialize.DB.Find(&category)

	initialize.DB.First(&book, id)

	meta := BuildMeta()
	path_builder := pathBuilder(meta.Keywords[0], "Self", c)
	links := []*Link{
		BookBuilder(&book[0], "self"),
		AuthorBuilder(&author[0], "author"),
		CategoryBuilder(&category[0], "category"),
	}
	c.IndentedJSON(200, gin.H{
		"Links": gin.H{
			"_Self": gin.H{
				"method": "GET",
				"self":   links[0].Path,
			},
			"Author": gin.H{
				"method": "GET",
				"author": links[1].Path,
			},
			"Category": gin.H{
				"method":   "GET",
				"category": links[2].Path,
			},
		},
		"Meta":        meta,
		"Test":        meta.Keywords[1],
		"PathBuilder": path_builder,
	})

}

func pathBuilder(entity string, path_name string, c *gin.Context) string {
	return "/api/" + entity + "/" + path_name
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

// func post_as_dict(posts, c *gin.Context){
// 	c.IndentedJSON(200, gin.H{
// 		"post_id": posts.ID,
// 		"post_title": posts.Title,
// 		"post_body": posts.Body,
// 	})
// }
