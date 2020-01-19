// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "title": "Service user",
    "license": {
      "name": "MIT"
    },
    "version": "0.1.0"
  },
  "basePath": "/api/v1",
  "paths": {
    "/email/verification": {
      "post": {
        "security": [],
        "operationId": "verificationEmail",
        "parameters": [
          {
            "name": "email",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Email"
            }
          }
        ],
        "responses": {
          "204": {
            "$ref": "#/responses/NoContent"
          },
          "default": {
            "$ref": "#/responses/GenericError"
          }
        }
      }
    },
    "/login": {
      "post": {
        "security": [],
        "description": "Login for user.",
        "operationId": "login",
        "parameters": [
          {
            "name": "args",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/LoginParam"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/User"
            },
            "headers": {
              "Set-Cookie": {
                "type": "string",
                "description": "Session auth."
              }
            }
          },
          "default": {
            "$ref": "#/responses/GenericError"
          }
        }
      }
    },
    "/logout": {
      "post": {
        "description": "Logout for user",
        "operationId": "logout",
        "responses": {
          "204": {
            "$ref": "#/responses/NoContent"
          },
          "default": {
            "$ref": "#/responses/GenericError"
          }
        }
      }
    },
    "/user": {
      "get": {
        "description": "Open user profile.",
        "operationId": "getUser",
        "parameters": [
          {
            "type": "integer",
            "format": "int32",
            "name": "id",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/User"
            }
          },
          "default": {
            "$ref": "#/responses/GenericError"
          }
        }
      },
      "post": {
        "security": [],
        "description": "New user registration. If it is not sent to username, it will be the userID",
        "operationId": "createUser",
        "parameters": [
          {
            "name": "args",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/CreateUserParams"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/User"
            },
            "headers": {
              "Set-Cookie": {
                "type": "string",
                "description": "Session auth."
              }
            }
          },
          "default": {
            "$ref": "#/responses/GenericError"
          }
        }
      },
      "delete": {
        "description": "Deletion of your account.",
        "operationId": "deleteUser",
        "responses": {
          "204": {
            "$ref": "#/responses/NoContent"
          },
          "default": {
            "$ref": "#/responses/GenericError"
          }
        }
      }
    },
    "/user/email": {
      "patch": {
        "description": "Change email.",
        "operationId": "updateEmail",
        "parameters": [
          {
            "name": "email",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Email"
            }
          }
        ],
        "responses": {
          "204": {
            "$ref": "#/responses/NoContent"
          },
          "default": {
            "$ref": "#/responses/GenericError"
          }
        }
      }
    },
    "/user/password": {
      "patch": {
        "description": "Change password.",
        "operationId": "updatePassword",
        "parameters": [
          {
            "name": "args",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/UpdatePassword"
            }
          }
        ],
        "responses": {
          "204": {
            "$ref": "#/responses/NoContent"
          },
          "default": {
            "$ref": "#/responses/GenericError"
          }
        }
      }
    },
    "/user/username": {
      "patch": {
        "description": "Change username.",
        "operationId": "updateUsername",
        "parameters": [
          {
            "name": "username",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Username"
            }
          }
        ],
        "responses": {
          "204": {
            "$ref": "#/responses/NoContent"
          },
          "default": {
            "$ref": "#/responses/GenericError"
          }
        }
      }
    },
    "/username/verification": {
      "post": {
        "security": [],
        "operationId": "verificationUsername",
        "parameters": [
          {
            "name": "username",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Username"
            }
          }
        ],
        "responses": {
          "204": {
            "$ref": "#/responses/NoContent"
          },
          "default": {
            "$ref": "#/responses/GenericError"
          }
        }
      }
    },
    "/users": {
      "get": {
        "description": "User search.",
        "operationId": "getUsers",
        "parameters": [
          {
            "name": "args",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/ListUsersParams"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "properties": {
                "total": {
                  "type": "integer",
                  "format": "int32"
                },
                "users": {
                  "type": "array",
                  "maxItems": 100,
                  "uniqueItems": true,
                  "items": {
                    "$ref": "#/definitions/User"
                  }
                }
              }
            }
          },
          "default": {
            "$ref": "#/responses/GenericError"
          }
        }
      }
    }
  },
  "definitions": {
    "CreateUserParams": {
      "type": "object",
      "required": [
        "email",
        "password",
        "username"
      ],
      "properties": {
        "email": {
          "$ref": "#/definitions/Email"
        },
        "password": {
          "$ref": "#/definitions/Password"
        },
        "username": {
          "$ref": "#/definitions/Username"
        }
      }
    },
    "Email": {
      "type": "string",
      "format": "email",
      "maxLength": 255,
      "minLength": 1
    },
    "Error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "message": {
          "type": "string"
        }
      }
    },
    "ListUsersParams": {
      "type": "object",
      "required": [
        "username",
        "pagination"
      ],
      "properties": {
        "pagination": {
          "$ref": "#/definitions/Pagination"
        },
        "username": {
          "$ref": "#/definitions/Username"
        }
      }
    },
    "LoginParam": {
      "type": "object",
      "required": [
        "email",
        "password"
      ],
      "properties": {
        "email": {
          "$ref": "#/definitions/Email"
        },
        "password": {
          "$ref": "#/definitions/Password"
        }
      }
    },
    "Pagination": {
      "type": "object",
      "required": [
        "limit",
        "offset"
      ],
      "properties": {
        "limit": {
          "description": "Maximum amount of items to return.",
          "type": "integer",
          "format": "int32",
          "maximum": 100,
          "minimum": 1
        },
        "offset": {
          "description": "Amount of items to skip.",
          "type": "integer",
          "format": "int32",
          "maximum": 9900
        }
      }
    },
    "Password": {
      "type": "string",
      "format": "password",
      "maxLength": 100,
      "minLength": 8
    },
    "UpdatePassword": {
      "type": "object",
      "required": [
        "old",
        "new"
      ],
      "properties": {
        "new": {
          "$ref": "#/definitions/Password"
        },
        "old": {
          "$ref": "#/definitions/Password"
        }
      }
    },
    "User": {
      "type": "object",
      "required": [
        "id",
        "username"
      ],
      "properties": {
        "id": {
          "$ref": "#/definitions/UserID"
        },
        "username": {
          "$ref": "#/definitions/Username"
        }
      }
    },
    "UserID": {
      "type": "integer",
      "format": "int32"
    },
    "Username": {
      "type": "string",
      "maxLength": 30,
      "minLength": 1
    }
  },
  "responses": {
    "GenericError": {
      "description": "Generic error response.",
      "schema": {
        "$ref": "#/definitions/Error"
      }
    },
    "NoContent": {
      "description": "The server successfully processed the request and is not returning any content."
    }
  },
  "securityDefinitions": {
    "cookieKey": {
      "description": "Session auth inside cookie.",
      "type": "apiKey",
      "name": "Cookie",
      "in": "header"
    }
  },
  "security": [
    {
      "cookieKey": []
    }
  ]
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "title": "Service user",
    "license": {
      "name": "MIT"
    },
    "version": "0.1.0"
  },
  "basePath": "/api/v1",
  "paths": {
    "/email/verification": {
      "post": {
        "security": [],
        "operationId": "verificationEmail",
        "parameters": [
          {
            "name": "email",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Email"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "The server successfully processed the request and is not returning any content."
          },
          "default": {
            "description": "Generic error response.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/login": {
      "post": {
        "security": [],
        "description": "Login for user.",
        "operationId": "login",
        "parameters": [
          {
            "name": "args",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/LoginParam"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/User"
            },
            "headers": {
              "Set-Cookie": {
                "type": "string",
                "description": "Session auth."
              }
            }
          },
          "default": {
            "description": "Generic error response.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/logout": {
      "post": {
        "description": "Logout for user",
        "operationId": "logout",
        "responses": {
          "204": {
            "description": "The server successfully processed the request and is not returning any content."
          },
          "default": {
            "description": "Generic error response.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/user": {
      "get": {
        "description": "Open user profile.",
        "operationId": "getUser",
        "parameters": [
          {
            "type": "integer",
            "format": "int32",
            "name": "id",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/User"
            }
          },
          "default": {
            "description": "Generic error response.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "post": {
        "security": [],
        "description": "New user registration. If it is not sent to username, it will be the userID",
        "operationId": "createUser",
        "parameters": [
          {
            "name": "args",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/CreateUserParams"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/User"
            },
            "headers": {
              "Set-Cookie": {
                "type": "string",
                "description": "Session auth."
              }
            }
          },
          "default": {
            "description": "Generic error response.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "delete": {
        "description": "Deletion of your account.",
        "operationId": "deleteUser",
        "responses": {
          "204": {
            "description": "The server successfully processed the request and is not returning any content."
          },
          "default": {
            "description": "Generic error response.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/user/email": {
      "patch": {
        "description": "Change email.",
        "operationId": "updateEmail",
        "parameters": [
          {
            "name": "email",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Email"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "The server successfully processed the request and is not returning any content."
          },
          "default": {
            "description": "Generic error response.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/user/password": {
      "patch": {
        "description": "Change password.",
        "operationId": "updatePassword",
        "parameters": [
          {
            "name": "args",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/UpdatePassword"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "The server successfully processed the request and is not returning any content."
          },
          "default": {
            "description": "Generic error response.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/user/username": {
      "patch": {
        "description": "Change username.",
        "operationId": "updateUsername",
        "parameters": [
          {
            "name": "username",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Username"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "The server successfully processed the request and is not returning any content."
          },
          "default": {
            "description": "Generic error response.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/username/verification": {
      "post": {
        "security": [],
        "operationId": "verificationUsername",
        "parameters": [
          {
            "name": "username",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Username"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "The server successfully processed the request and is not returning any content."
          },
          "default": {
            "description": "Generic error response.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/users": {
      "get": {
        "description": "User search.",
        "operationId": "getUsers",
        "parameters": [
          {
            "name": "args",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/ListUsersParams"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "properties": {
                "total": {
                  "type": "integer",
                  "format": "int32",
                  "minimum": 0
                },
                "users": {
                  "type": "array",
                  "maxItems": 100,
                  "uniqueItems": true,
                  "items": {
                    "$ref": "#/definitions/User"
                  }
                }
              }
            }
          },
          "default": {
            "description": "Generic error response.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "CreateUserParams": {
      "type": "object",
      "required": [
        "email",
        "password",
        "username"
      ],
      "properties": {
        "email": {
          "$ref": "#/definitions/Email"
        },
        "password": {
          "$ref": "#/definitions/Password"
        },
        "username": {
          "$ref": "#/definitions/Username"
        }
      }
    },
    "Email": {
      "type": "string",
      "format": "email",
      "maxLength": 255,
      "minLength": 1
    },
    "Error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "message": {
          "type": "string"
        }
      }
    },
    "ListUsersParams": {
      "type": "object",
      "required": [
        "username",
        "pagination"
      ],
      "properties": {
        "pagination": {
          "$ref": "#/definitions/Pagination"
        },
        "username": {
          "$ref": "#/definitions/Username"
        }
      }
    },
    "LoginParam": {
      "type": "object",
      "required": [
        "email",
        "password"
      ],
      "properties": {
        "email": {
          "$ref": "#/definitions/Email"
        },
        "password": {
          "$ref": "#/definitions/Password"
        }
      }
    },
    "Pagination": {
      "type": "object",
      "required": [
        "limit",
        "offset"
      ],
      "properties": {
        "limit": {
          "description": "Maximum amount of items to return.",
          "type": "integer",
          "format": "int32",
          "maximum": 100,
          "minimum": 1
        },
        "offset": {
          "description": "Amount of items to skip.",
          "type": "integer",
          "format": "int32",
          "maximum": 9900,
          "minimum": 0
        }
      }
    },
    "Password": {
      "type": "string",
      "format": "password",
      "maxLength": 100,
      "minLength": 8
    },
    "UpdatePassword": {
      "type": "object",
      "required": [
        "old",
        "new"
      ],
      "properties": {
        "new": {
          "$ref": "#/definitions/Password"
        },
        "old": {
          "$ref": "#/definitions/Password"
        }
      }
    },
    "User": {
      "type": "object",
      "required": [
        "id",
        "username"
      ],
      "properties": {
        "id": {
          "$ref": "#/definitions/UserID"
        },
        "username": {
          "$ref": "#/definitions/Username"
        }
      }
    },
    "UserID": {
      "type": "integer",
      "format": "int32"
    },
    "Username": {
      "type": "string",
      "maxLength": 30,
      "minLength": 1
    }
  },
  "responses": {
    "GenericError": {
      "description": "Generic error response.",
      "schema": {
        "$ref": "#/definitions/Error"
      }
    },
    "NoContent": {
      "description": "The server successfully processed the request and is not returning any content."
    }
  },
  "securityDefinitions": {
    "cookieKey": {
      "description": "Session auth inside cookie.",
      "type": "apiKey",
      "name": "Cookie",
      "in": "header"
    }
  },
  "security": [
    {
      "cookieKey": []
    }
  ]
}`))
}
