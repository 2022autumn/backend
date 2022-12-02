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
        "/es/get/": {
            "get": {
                "description": "根据id获取对象，可以是author，work，institution,venue,concept",
                "tags": [
                    "esSearch"
                ],
                "summary": "txc",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":200,\"res\":{obeject}}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "201": {
                        "description": "{\"status\":201,\"msg\":\"es get err\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"status\":400,\"msg\":\"id type error\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/es/getAuthorRelationNet": {
            "get": {
                "description": "根据author的id获取专家关系网络, 目前会返回Top N的关系网，N=10，后续可以讨论修改N的大小或者传参给我\n\n目前接口时延约为1s, 后续考虑把计算出来的结果存入数据库，二次查询时延降低\n\n接口使用示例 1. author_id=A2764814280  2. author_id=A2900471938",
                "tags": [
                    "esSearch"
                ],
                "summary": "hr",
                "parameters": [
                    {
                        "type": "string",
                        "description": "author_id",
                        "name": "author_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"data\":{response.AuthorRelationNet}}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "201": {
                        "description": "{\"msg\":\"Get Author Relation Net Error\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/es/search/advanced": {
            "post": {
                "description": "高级搜索，query是一个map列表， 每个map包含\"content\" \"field\" \"logic\"\nlogic 仅包含[\"and\", \"or\", \"not\"]\nfield 仅包含[\"title\", \"abstract\", \"venue\", \"publisher\", \"author\", \"institution\", \"concept\"]\n对于年份的筛选，在query里面 field是\"publication_date\" logic默认为and， 该map下有\"begin\" \"end\"分别是开始和结束\nsort=0为默认排序（降序） =1为按引用数降序 =2按发表日期由近到远\nasc=0为降序 =1为升序\n{ \"asc\": false,\"conds\": {\"venue\":\"International Journal for Research in Applied Science and Engineering Technology\",\"author\": \"Zenith Nandy\"},\"page\": 1,\"query\": [{\"field\": \"title\",\"content\": \"python\",\"logic\": \"and\"},{\"field\": \"publication_date\",\"begin\": \"2021-12-01\",\"end\":\"2022-06-01\",\"logic\": \"and\"}],\"size\": 8,\"sort\": 0}",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "esSearch"
                ],
                "summary": "txc",
                "parameters": [
                    {
                        "description": "data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/response.AdvancedSearchQ"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/es/search/base": {
            "post": {
                "description": "基本搜索，Cond里面填筛选条件，key仅包含[\"type\", \"author\", \"institution\", \"publisher\", \"venue\", \"publication_year\"]",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "esSearch"
                ],
                "summary": "txc",
                "parameters": [
                    {
                        "description": "搜索条件",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/response.BaseSearchQ"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":200,\"res\":{obeject}}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "201": {
                        "description": "{\"status\":201,\"err\":\"es search err\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/es/search/doi": {
            "post": {
                "description": "使用doi查找work，未测试，请勿使用",
                "tags": [
                    "esSearch"
                ],
                "summary": "txc",
                "parameters": [
                    {
                        "type": "string",
                        "description": "doi",
                        "name": "doi",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/login": {
            "post": {
                "description": "登录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "ccf",
                "parameters": [
                    {
                        "description": "data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/response.LoginQ"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":200,\"success\":true,\"msg\":\"login success\",\"token\": 666}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"status\":400,\"success\":false,\"msg\":\"username doesn't exist\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "{\"status\":401,\"success\":false,\"msg\":\"password doesn't match\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "填入用户名和密码注册",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "ccf",
                "parameters": [
                    {
                        "description": "data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/response.RegisterQ"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":200,\"msg\":\"register success\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"status\":400,\"msg\":\"username exists\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/headshot": {
            "post": {
                "description": "上传用户头像",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "ccf",
                "parameters": [
                    {
                        "type": "file",
                        "description": "新头像",
                        "name": "Headshot",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":200,\"msg\":\"修改成功\",\"data\":{object}}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"status\":400,\"msg\":\"用户ID不存在\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "{\"status\":401,\"msg\":\"头像文件上传失败\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "402": {
                        "description": "{\"status\":402,\"msg\":\"文件保存失败\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "{\"status\":403,\"msg\":\"保存文件路径到数据库中失败\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/info": {
            "get": {
                "description": "查看用户个人信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "ccf",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user_id",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":200,\"msg\":\"get info of user\",\"data\":{object}}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"status\":400,\"msg\":\"userID not exist\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/mod": {
            "post": {
                "description": "编辑用户信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "ccf",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user_id",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/response.ModifyQ"
                        }
                    },
                    {
                        "type": "string",
                        "description": "个性签名",
                        "name": "user_info",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "电话号码",
                        "name": "phone_number",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Email",
                        "name": "email",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":200,\"msg\":\"修改成功\",\"data\":{object}}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"status\":400,\"msg\":\"用户ID不存在\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "{\"status\":401,\"msg\":err.Error()}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/pwd": {
            "post": {
                "description": "编辑用户信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "ccf",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user_id",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "旧密码",
                        "name": "Password_Old",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "新密码",
                        "name": "Password_New",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"status\":200,\"msg\":\"修改成功\",\"data\":{object}}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"status\":400,\"msg\":\"用户ID不存在\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "{\"status\":401,\"msg\":\"原密码输入错误\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "402": {
                        "description": "{\"status\":402,\"msg\":err1.Error()}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "response.AdvancedSearchQ": {
            "type": "object",
            "properties": {
                "asc": {
                    "type": "boolean"
                },
                "conds": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "page": {
                    "type": "integer"
                },
                "query": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "additionalProperties": {
                            "type": "string"
                        }
                    }
                },
                "size": {
                    "type": "integer"
                },
                "sort": {
                    "type": "integer"
                }
            }
        },
        "response.BaseSearchQ": {
            "type": "object",
            "properties": {
                "asc": {
                    "type": "boolean"
                },
                "conds": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "kind": {
                    "type": "string"
                },
                "page": {
                    "type": "integer"
                },
                "queryWord": {
                    "type": "string"
                },
                "size": {
                    "type": "integer"
                },
                "sort": {
                    "type": "integer"
                }
            }
        },
        "response.LoginQ": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "response.ModifyQ": {
            "type": "object",
            "properties": {
                "email": {
                    "description": "邮箱",
                    "type": "string"
                },
                "fields": {
                    "description": "研究领域",
                    "type": "string"
                },
                "interest_tag": {
                    "description": "兴趣词",
                    "type": "string"
                },
                "name": {
                    "description": "真实姓名",
                    "type": "string"
                },
                "phone": {
                    "description": "电话号码",
                    "type": "string"
                },
                "user_info": {
                    "description": "个性签名",
                    "type": "string"
                }
            }
        },
        "response.RegisterQ": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "minLength": 6
                },
                "username": {
                    "type": "string",
                    "maxLength": 100,
                    "minLength": 3
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
