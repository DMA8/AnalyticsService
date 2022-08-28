// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/approved_tasks": {
            "get": {
                "description": "endpoint return count of approved tasks",
                "produces": [
                    "application/json"
                ],
                "summary": "get count of approved tasks",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Counter"
                        }
                    }
                }
            }
        },
        "/declined_tasks": {
            "get": {
                "description": "endpoint return count of declined tasks",
                "produces": [
                    "application/json"
                ],
                "summary": "get count of declined tasks",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Counter"
                        }
                    }
                }
            }
        },
        "/summary_time": {
            "get": {
                "description": "Return task id and summary time of decision in seconds",
                "produces": [
                    "application/json"
                ],
                "summary": "Get summary time for each task",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.SummaryTime"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Counter": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer",
                    "example": 5
                }
            }
        },
        "models.SummaryTime": {
            "type": "object",
            "properties": {
                "duration": {
                    "type": "integer",
                    "example": 1005
                },
                "task_id": {
                    "description": "TODO поменять на реальный пример",
                    "type": "string",
                    "example": "test123"
                }
            }
        }
    },
    "x-extension-openapi": {
        "example": "value on a json format"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3003",
	BasePath:         "/analytics/v1",
	Schemes:          []string{},
	Title:            "Swagger Analytics API",
	Description:      "Analytics server",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
