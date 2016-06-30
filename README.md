## Installing

```
$ git clone git@github.com:brettbuddin/ponyexpress.git
$ cd ponyexpress/
$ make
$ ponyexpress
...
```

## Running Tests

```
$ make test
...
```

## Running Benchmarks

```
$ go test ./mailbox -check.b -v
...
```

## Usage Example

```
$ curl http://localhost:3000/mailboxes -X POST
{
	"mailbox":{
		"id":"958ff9d3-152b-4d05-9b97-536e3331e419"
	}
}

$ curl -H "Content-Type: application/json" \
       -X POST -d '{"message":{"sender": "brett@buddin.us", "subject": "Howdy", "body": "Some text"}}' \
       http://localhost:3000/mailboxes/958ff9d3-152b-4d05-9b97-536e3331e419/messages
{  
    "message":{  
        "id":"a1294fc4-c511-402b-9192-c4195a35b7dd",
        "sender":"brett@buddin.us",
        "subject":"Howdy",
        "body":"Some text",
        "received":"2016-06-30T08:17:36.616531006-04:00"
    }
}

$ curl -H "Content-Type: application/json" \
       -X POST -d '{"message":{"sender": "brett@buddin.us", "subject": "OMG", "body": "The house is on fire."}}' \
       http://localhost:3000/mailboxes/958ff9d3-152b-4d05-9b97-536e3331e419/messages
{  
    "message":{  
        "id":"bcb08e10-e3b3-40c1-914c-b3f1d2af313e",
        "sender":"brett@buddin.us",
        "subject":"OMG",
        "body":"The house is on fire.",
        "received":"2016-06-30T08:17:59.169944422-04:00"
    }
}

$ curl http://localhost:3000/mailboxes/958ff9d3-152b-4d05-9b97-536e3331e419/messages
{  
   "messages":[  
      {  
         "id":"bcb08e10-e3b3-40c1-914c-b3f1d2af313e",
         "sender":"brett@buddin.us",
         "subject":"OMG",
         "body":"The house is on fire.",
         "received":"2016-06-30T08:17:59.169944422-04:00"
      },
      {  
         "id":"a1294fc4-c511-402b-9192-c4195a35b7dd",
         "sender":"brett@buddin.us",
         "subject":"Howdy",
         "body":"Some text",
         "received":"2016-06-30T08:17:36.616531006-04:00"
      }
   ],
   "meta":{  
      "results":2,
      "limit":100,
      "since_id":"",
      "last_id":"bcb08e10-e3b3-40c1-914c-b3f1d2af313e"
   }
}

$ curl http://localhost:3000/mailboxes/958ff9d3-152b-4d05-9b97-536e3331e419/messages\?since_id\=a1294fc4-c511-402b-9192-c4195a35b7dd
{  
   "messages":[  
      {  
         "id":"bcb08e10-e3b3-40c1-914c-b3f1d2af313e",
         "sender":"brett@buddin.us",
         "subject":"OMG",
         "body":"The house is on fire.",
         "received":"2016-06-30T08:17:59.169944422-04:00"
      }
   ],
   "meta":{  
      "results":1,
      "limit":100,
      "since_id":"a1294fc4-c511-402b-9192-c4195a35b7dd",
      "last_id":"bcb08e10-e3b3-40c1-914c-b3f1d2af313e"
   }
}

$ curl http://localhost:3000/mailboxes/958ff9d3-152b-4d05-9b97-536e3331e419/messages/a1294fc4-c511-402b-9192-c4195a35b7dd
{  
   "message":{  
      "id":"a1294fc4-c511-402b-9192-c4195a35b7dd",
      "sender":"brett@buddin.us",
      "subject":"Howdy",
      "body":"Some text",
      "received":"2016-06-30T08:17:36.616531006-04:00"
   }
}

$ curl http://localhost:3000/mailboxes/958ff9d3-152b-4d05-9b97-536e3331e419/messages/a1294fc4-c511-402b-9192-c4195a35b7dd -X DELETE
{  
   "message":{  
      "id":"a1294fc4-c511-402b-9192-c4195a35b7dd",
      "sender":"brett@buddin.us",
      "subject":"Howdy",
      "body":"Some text",
      "received":"2016-06-30T08:17:36.616531006-04:00"
   }
}

$ curl http://localhost:3000/mailboxes/958ff9d3-152b-4d05-9b97-536e3331e419/messages/a1294fc4-c511-402b-9192-c4195a35b7dd
{
	"error":"unknown message: a1294fc4-c511-402b-9192-c4195a35b7dd"
}

$ curl http://localhost:3000/mailboxes/958ff9d3-152b-4d05-9b97-536e3331e419 -X DELETE
{
	"mailbox":{
		"id":"958ff9d3-152b-4d05-9b97-536e3331e419"
	}
}

$ curl http://localhost:3000/mailboxes/958ff9d3-152b-4d05-9b97-536e3331e419/messages
{
	"error":"unknown mailbox: 958ff9d3-152b-4d05-9b97-536e3331e419"
}
```
