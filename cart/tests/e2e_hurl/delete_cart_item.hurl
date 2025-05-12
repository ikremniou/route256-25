DELETE http://localhost:8080/user/3/cart/1625903
HTTP 204

POST http://localhost:8080/user/3/cart/1625903
{
    "count": 9
}
HTTP 200

DELETE http://localhost:8080/user/3/cart/1625903
HTTP 204

POST http://localhost:8080/user/3/cart/1625903
{
    "count": 9
}
HTTP 200

# CleanUP
DELETE http://localhost:8080/user/3/cart
HTTP 204