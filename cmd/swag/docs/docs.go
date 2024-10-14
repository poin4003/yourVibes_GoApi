// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://github.com/poin4003/yourVibes_GoApi",
        "contact": {
            "name": "TEAM HKTP",
            "url": "https://github.com/poin4003/yourVibes_GoApi",
            "email": "pchuy4003@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/posts/": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "When user create post",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "post"
                ],
                "summary": "Post create post",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Title of the post",
                        "name": "title",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Content of the post",
                        "name": "content",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Privacy level",
                        "name": "privacy",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Location of the post",
                        "name": "location",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "Media files for the post",
                        "name": "media",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ResponseData"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrResponse"
                        }
                    }
                }
            }
        },
        "/posts/getMany/{userId}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Retrieve multiple posts filtered by various criteria.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "post"
                ],
                "summary": "Get many posts",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID to filter posts",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Filter by post title",
                        "name": "title",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by content",
                        "name": "content",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by location",
                        "name": "location",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Filter by advertisement",
                        "name": "is_advertisement",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by creation time",
                        "name": "created_at",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Which column to sort by",
                        "name": "sort_by",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Order by descending if true",
                        "name": "isDescending",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Limit of posts per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page number for pagination",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ResponseData"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrResponse"
                        }
                    }
                }
            }
        },
        "/posts/{postId}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Retrieve a post by its unique ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "post"
                ],
                "summary": "Get post by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Post ID",
                        "name": "postId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ResponseData"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrResponse"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "When user need to update information of post or update media",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "post"
                ],
                "summary": "update post",
                "parameters": [
                    {
                        "type": "string",
                        "description": "PostId",
                        "name": "postId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Post title",
                        "name": "title",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Post content",
                        "name": "content",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Post privacy",
                        "name": "privacy",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Post location",
                        "name": "location",
                        "in": "formData"
                    },
                    {
                        "type": "integer",
                        "description": "Array of mediaIds you want to delete",
                        "name": "media_ids",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "Array of media you want to upload",
                        "name": "media",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ResponseData"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrResponse"
                        }
                    }
                }
            }
        },
        "/users/login/": {
            "post": {
                "description": "When user login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "User login",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth_dto.LoginCredentials"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ResponseData"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrResponse"
                        }
                    }
                }
            }
        },
        "/users/register/": {
            "post": {
                "description": "When user registration",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "User Registration",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth_dto.RegisterCredentials"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ResponseData"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrResponse"
                        }
                    }
                }
            }
        },
        "/users/verifyemail/": {
            "post": {
                "description": "Before user registration",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "User verify email",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth_dto.VerifyEmailInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ResponseData"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth_dto.LoginCredentials": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "auth_dto.RegisterCredentials": {
            "type": "object",
            "required": [
                "birthday",
                "email",
                "family_name",
                "name",
                "otp",
                "password",
                "phone_number"
            ],
            "properties": {
                "birthday": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "family_name": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "otp": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 8
                },
                "phone_number": {
                    "type": "string"
                }
            }
        },
        "auth_dto.VerifyEmailInput": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "response.ErrResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "$ref": "#/definitions/response.ErrResponseChild"
                }
            }
        },
        "response.ErrResponseChild": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "detail_err": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "response.PagingResponse": {
            "type": "object",
            "properties": {
                "limit": {
                    "type": "integer"
                },
                "page": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "response.ResponseData": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "Status code",
                    "type": "integer"
                },
                "data": {
                    "description": "Data"
                },
                "message": {
                    "description": "Status message",
                    "type": "string"
                },
                "paging": {
                    "description": "Paging (optional)",
                    "allOf": [
                        {
                            "$ref": "#/definitions/response.PagingResponse"
                        }
                    ]
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "Token without 'Bearer ' prefix",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:8080",
	BasePath:         "/v1/2024",
	Schemes:          []string{},
	Title:            "API Documentation YourVibes backend",
	Description:      "This is a sample YourVibes backend server",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
