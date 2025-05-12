POST http://localhost:8080/user/4/cart/321
HTTP 400

POST http://localhost:8080/users/4/cart/321
{
    "count": 10
}
HTTP 404

POST http://localhost:8080/user/1-23/cart/321
{
    "count": 10
}
HTTP 400

POST http://localhost:8080/user/4/cart/-321
{
    "count": 10
}
HTTP 400

POST http://localhost:8080/user/4/cart/0
{
    "count": 10
}
HTTP 400

POST http://localhost:8080/user/4/cart/1
{
    "count1": 10
}
HTTP 400

POST http://localhost:8080/user/4/cart/2
{ }
HTTP 400

POST http://localhost:8080/user/4/cart/3
{ 
    "count": 4294967296
}
HTTP 400

POST http://localhost:8080/user/4/cart/3
{ 
    "count": 0
}
HTTP 400

POST http://localhost:8080/user/4/cart/2956315
{
    "count": 1
}
HTTP *
[Asserts]
status < 300

POST http://localhost:8080/user/4/cart/2956315
{
    "count": 10
}
HTTP *
[Asserts]
status < 300

# Sku does not exist in the products
POST http://localhost:8080/user/4/cart/3
{ 
    "count": 1
}
HTTP 412

DELETE http://localhost:8080/user/4/cart
HTTP 204