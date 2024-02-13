// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/login/google": {
            "post": {
                "description": "Get userID and paseto token by google oauth2 token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Try login with Google OAuth2 token",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.LoginWithGoogleRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.LoginWithGoogleResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/errors.CodedError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errors.CodedError"
                        }
                    }
                }
            }
        },
        "/validate": {
            "get": {
                "description": "Validate user's token and if correct return UserInfo",
                "produces": [
                    "application/json"
                ],
                "summary": "Validate user's token",
                "parameters": [
                    {
                        "type": "string",
                        "example": "Bearer v2.local.ABCDEFG",
                        "description": "Paseto encrypted token",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.UserInfo"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/errors.CodedError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errors.CodedError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.LoginWithGoogleRequest": {
            "type": "object",
            "properties": {
                "google_oauth_token": {
                    "type": "string"
                }
            }
        },
        "auth.LoginWithGoogleResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "auth.UserInfo": {
            "type": "object",
            "properties": {
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "errors.CodedError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "err": {}
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "description": "Enter the token with the ` + "`" + `Bearer: ` + "`" + ` prefix, e.g. \\\"Bearer abcde12345\\\"",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "tcarzverey.ru:8082",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "HSE MLOps Auth server",
	Description:      "Auth service for mlops project.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
