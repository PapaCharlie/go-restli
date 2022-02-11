HTTP/1.1 200 OK
Content-Length: 212
Content-Type: application/json
X-RestLi-Protocol-Version: 2.0.0

{
  "statuses" : { },
  "results" : {
    "($params:(temp:1),string:1)" : {
      "status": 204
    },
    "($params:(temp:42),string:string%3Awith%3Acolons)" : {
      "status": 205
    }
  },
  "errors" : { }
}
