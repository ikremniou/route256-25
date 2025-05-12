DELETE http://localhost:8080/user/1/cart
HTTP 204

GET http://localhost:8080/user/1/cart
HTTP 404

POST http://localhost:8080/user/1/cart/2956315
{
    "count": 9
}
HTTP 200

GET http://localhost:8080/user/1/cart
HTTP 200
[Asserts]
jsonpath "$.items[0].count" == 9
jsonpath "$.items[0].price" > 0
jsonpath "$.total_price" > 1000
jsonpath "$.items[0].name" isString


POST http://localhost:8080/user/1/cart/2956315
{
    "count": 9
}
HTTP 200

GET http://localhost:8080/user/1/cart
HTTP 200
[Asserts]
jsonpath "$.items[0].count" == 18
jsonpath "$.items[0].price" > 0
jsonpath "$.total_price" > 2000
jsonpath "$.items[0].name" isString

POST http://localhost:8080/user/1/cart/2618151
{
    "count": 100
}
HTTP 200

GET http://localhost:8080/user/1/cart
HTTP 200
[Asserts]
jsonpath "$.items[*].sku" includes 2956315
jsonpath "$.items[*].sku" includes 2618151
