package testData

var UserDefinition = `
"User": {
  "type": "object",
  "properties": {
    "email": {
        "type": "string",
  "format": "email"
    },
    "token": {
      "type": "string"
    },
    "username": {
      "type": "string"
    },
    "bio": {
      "type": "string"
    },
    "image": {
      "type": "string"
    }
  },
  "required": [
    "email",
    "token",
    "username",
    "bio",
    "image"
  ]
}`

var UserRespDefinition = `{
	` + UserDefinition + `,
      "type": "object",
      "properties": {
        "user": {
          "$ref": "#/User"
        }
      },
      "required": [
        "user"
      ]
}`

var ProfileDefinition = `
"Profile": {
  "type": "object",
  "properties": {
    "username": {
      "type": "string"
    },
    "bio": {
      "type": "string"
    },
    "image": {
      "type": "string"
    },
    "following": {
      "type": "boolean"
    },
	"createdAt": {
      "type": "string",
      "format": "date-time"
    },
    "updatedAt": {
      "type": "string",
      "format": "date-time"
    }
  },
  "required": [
    "username",
    "bio",
    "image",
    "following",
	"createdAt",
	"updatedAt"
  ]
}`

var ProfileRespDefinition = `{
` + ProfileDefinition + `,
  "type": "object",
  "properties": {
    "profile": {
     "$ref": "#/Profile"
    }
  },
  "required": [
    "profile"
  ]
}`

var ArticleDefinition = `
` + ProfileDefinition + `,
"Article": {
  "type": "object",
  "properties": {
    "slug": {
      "type": "string"
    },
    "title": {
      "type": "string"
    },
    "description": {
      "type": "string"
    },
    "body": {
      "type": "string"
    },
    "tagList": {
      "type": "array",
      "items": {
        "type": "string"
      }
    },
    "createdAt": {
      "type": "string",
      "format": "date-time"
    },
    "updatedAt": {
      "type": "string",
      "format": "date-time"
    },
    "favorited": {
      "type": "boolean"
    },
    "favoritesCount": {
      "type": "integer"
    },
    "author": {
      "$ref": "#/Profile"
    }
  },
  "required": [
    "slug",
    "title",
    "description",
    "body",
    "tagList",
    "createdAt",
    "updatedAt",
    "favorited",
    "favoritesCount",
    "author"
  ]
}`

var ArticleSingleRespDefinition = `{
` + ArticleDefinition + `,
  "type": "object",
  "properties": {
    "article": {
      "$ref": "#/Article"
    }
  },
  "required": [
    "article"
  ]
}`

var ArticleMultipleRespDefinition = `{
` + ArticleDefinition + `,
  "type": "object",
  "properties": {
    "articles": {
      "type": "array",
      "items": {
        "$ref": "#/Article"
      }
    },
    "articlesCount": {
      "type": "integer"
    }
  },
  "required": [
    "articles",
    "articlesCount"
  ]
}`

var CommentDefinition = `
` + ProfileDefinition + `,
"Comment": {
  "type": "object",
  "properties": {
    "id": {
      "type": "integer"
    },
    "createdAt": {
      "type": "string",
      "format": "date-time"
    },
    "updatedAt": {
      "type": "string",
      "format": "date-time"
    },
    "body": {
      "type": "string"
    },
    "author": {
      "$ref": "#/Profile"
    }
  },
  "required": [
    "id",
    "createdAt",
    "updatedAt",
    "body",
    "author"
  ]
}`

var CommentsSimgleResponse = `{
` + CommentDefinition + `,
  "type": "object",
  "properties": {
    "comment": {
      "$ref": "#/Comment"
    }
  },
  "required": [
    "comment"
  ]
}`

var CommentsMultipleResponse = `{
` + CommentDefinition + `,
  "type": "object",
  "properties": {
    "comments": {
      "type": "array",
      "items": {
        "$ref": "#/Comment"
      }
    }
  },
  "required": [
    "comments"
  ]
}`

var TagsResponse = `{
  "type": "object",
  "properties": {
    "tags": {
      "type": "array",
      "items": {
        "type": "string"
      }
    }
  },
  "required": [
    "tags"
  ]
}`

var ErrorResponse = `{
  "type": "object",
  "properties": {
    "errors": {
      "type": "object",
      "properties": {
        "body": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      },
      "required": [
        "body"
      ]
    }
  },
  "required": [
    "errors"
  ]
}`
