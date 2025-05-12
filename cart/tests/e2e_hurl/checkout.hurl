POST http://localhost:8080/checkout/0
HTTP 400

POST http://localhost:8080/user/5/cart/2956315
{
    "count": 1
}
HTTP 200

POST http://localhost:8080/checkout/5
HTTP 200
[Asserts]
jsonpath "$.orderID" > 1

POST http://localhost:8080/user/6/cart/139275865
{
    "count": 65535
}
HTTP 200

POST http://localhost:8080/checkout/6
HTTP 412