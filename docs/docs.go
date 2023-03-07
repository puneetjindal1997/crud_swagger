package docs

import (
	"bytes"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

var doc = `{
	"swagger": "2.0",
  "basePath": "{{.BasePath}}",
  "host": "{{.Host}}",
  "info": {
    "contact": {},
    "description": "{{.Description}}",
    "license": {},
    "title": "{{.Title}}",
    "version": "{{.Version}}"
  },
  "tags": [
	{
		"name": "CRUD",
		"description": "CRUD/user"
	  },
  ],
  "paths": {
	"/data": {
		"post": {
			"tags": [
				"CRUD"
			],
			"description": "create data",
			"summary": "create data",
			"produces": [
				"application/json"
			],
			"consumes": [
				"application/json"
			],
			"parameters": [
	  {
		"in": "body",
		"name": "user",
		"description": "The user to create.",
		"schema": {
		  "type": "object",
		  "required": [
			"email",
			"name"
		  ],
		  "properties": {
			"email": {
			  "type": "string"
			},
			"name": {
				"type": "string"
			}
		  }
		}
	  }
	],
			"responses": {
				"201": {
					"description": "Success",
				},
				"400": {
					"description": "Client error",
				},
				"422": {
					"description": "Client error",
				},
				"500": {
					"description": "Client error",
				}
			}
		}
	},
	"/data1": {
		"get": {
			"tags": [
				"CRUD"
			],
			"description": "create data",
			"summary": "create data",
			"produces": [
				"application/json"
			],
			"consumes": [
				"application/json"
			],
			"parameters": [],
			"responses": {
				"201": {
					"description": "Success",
				},
				"400": {
					"description": "Client error",
				},
				"422": {
					"description": "Client error",
				},
				"500": {
					"description": "Client error",
				}
			}
		}
	},
  }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
