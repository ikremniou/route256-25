POST http://localhost:8080/user/2/cart/1148162
{
    "count": 9
}
HTTP 200

POST http://localhost:8080/user/2/cart/1076963
{
    "count": 9
}
HTTP 200

DELETE http://localhost:8080/user/2/cart
HTTP 204

POST http://localhost:8080/user/2/cart/1076963
{
    "count": 9
}
HTTP 200

DELETE http://localhost:8080/user/2/cart
HTTP 204