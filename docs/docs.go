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
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/fund": {
            "get": {
                "tags": [
                    "fund_controller"
                ],
                "responses": {}
            }
        },
        "/fund/records": {
            "get": {
                "tags": [
                    "fund_record_controller"
                ],
                "responses": {}
            },
            "post": {
                "tags": [
                    "fund_record_controller"
                ],
                "parameters": [
                    {
                        "description": "FundRecordDTO",
                        "name": "record",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.FundRecordDTO"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/funds/{fund}": {
            "get": {
                "tags": [
                    "fund_controller"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "fundCode",
                        "name": "fund",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "post": {
                "tags": [
                    "fund_controller"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "fundCode",
                        "name": "fund",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "model.FundRecordDTO": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "code": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "quantity": {
                    "type": "number"
                },
                "type": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
