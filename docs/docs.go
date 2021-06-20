// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/upload/file": {
            "post": {
                "description": "Completes a multipart upload by assembling previously uploaded parts",
                "consumes": [
                    "multipart/form-data",
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Upload File",
                "parameters": [
                    {
                        "description": "upload file ",
                        "name": "file",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "location of the file in s3 bucket default is apk",
                        "name": "location",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.CompleteMultiPartOutput"
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/upload/finish": {
            "post": {
                "description": "Completes a multipart upload by assembling previously uploaded parts",
                "produces": [
                    "application/json"
                ],
                "summary": "Completes a multipart upload upload",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Upload ID that identifies the multipart upload.",
                        "name": "X-Upload-Id",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.CompleteMultiPartOutput"
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/upload/initlize": {
            "post": {
                "description": "This action initiates a multipart upload and returns an upload ID. This upload ID is used to associate all of the parts in the specific multipart upload. You specify this upload ID in each of your subsequent upload part requests",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Initilize Multipart Upload",
                "parameters": [
                    {
                        "description": "file info",
                        "name": "bottles",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.initilizeRquest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.MultiPartUploadResponse"
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/upload/part": {
            "post": {
                "description": "Uploads a part in a multipart upload. In this operation, you provide part data\nin your request. However, you have an option to specify your existing Amazon S3\nobject as a data source for the part you are uploading. To upload a part from an\nexisting object, you use the UploadPartCopy",
                "consumes": [
                    "multipart/form-data",
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Uploads a part in a multipart upload",
                "parameters": [
                    {
                        "description": "file info",
                        "name": "bottles",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.PartInput"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Upload ID that identifies the multipart upload.",
                        "name": "X-Upload-Id",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.PartUploadOutput"
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.initilizeRquest": {
            "type": "object",
            "properties": {
                "fileType": {
                    "type": "string"
                }
            }
        },
        "service.CompleteMultiPartOutput": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "location": {
                    "type": "string"
                }
            }
        },
        "service.MultiPartUploadResponse": {
            "type": "object",
            "properties": {
                "abortDate": {
                    "description": "If the bucket has a lifecycle rule configured with an action to abort incomplete\nmultipart uploads and the prefix in the lifecycle rule matches the object name\nin the request, the response includes this header. The header indicates when the\ninitiated multipart upload becomes eligible for an abort operation. For more\ninformation, see  Aborting Incomplete Multipart Uploads Using a Bucket Lifecycle\nPolicy\n(https://docs.aws.amazon.com/AmazonS3/latest/dev/mpuoverview.html#mpu-abort-incomplete-mpu-lifecycle-config).\nThe response also includes the x-amz-abort-rule-id header that provides the ID\nof the lifecycle configuration rule that defines this action.",
                    "type": "string"
                },
                "bucket": {
                    "type": "string"
                },
                "key": {
                    "description": "Object key for which the multipart upload was initiated.",
                    "type": "string"
                },
                "uploadId": {
                    "description": "ID for the initiated multipart upload.",
                    "type": "string"
                }
            }
        },
        "service.PartInput": {
            "type": "object",
            "properties": {
                "chunk": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "partNumber": {
                    "type": "integer"
                }
            }
        },
        "service.PartUploadOutput": {
            "type": "object",
            "properties": {
                "etag": {
                    "description": "Entity tag returned when the part was uploaded.",
                    "type": "string"
                },
                "partNumber": {
                    "description": "Part number that identifies the part. This is a positive integer between 1 and\n10,000.",
                    "type": "integer"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}