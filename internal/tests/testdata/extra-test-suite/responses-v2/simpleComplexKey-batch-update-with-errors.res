HTTP/1.1 200 OK
Content-Length: 433
Content-Type: application/json
X-RestLi-Protocol-Version: 2.0.0

{
  "results" : {
    "(string:b)" : {
      "status" : 204
    }
  },
  "errors" : {
    "(string:a)" : {
      "exceptionClass" : "com.linkedin.restli.server.RestLiServiceException",
      "stackTrace" : "trace",
      "message" : "message",
      "status" : 400
    },
    "(string:c)" : {
      "exceptionClass" : "com.linkedin.restli.server.RestLiServiceException",
      "stackTrace" : "trace",
      "status" : 500
    }
  }
}
