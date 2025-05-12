POST http://localhost:8084/order/create
{
    "user": 1,
    "items": [
        {
            "sku": 139275865,
            "count": 1
        }
    ]
}
HTTP 200

POST http://localhost:8084/order/create
{
    "user": 2,
    "items": [
        {
            "sku": 999999999,
            "count": 1
        }
    ]
}
HTTP 400

POST http://localhost:8084/order/create
{
    "user": 3,
    "items": [
        {
            "sku": 102,
            "count": 99999999
        }
    ]
}
HTTP 400

GET http://localhost:8084/stock/info?sku=139275865
HTTP 200
[Captures]
stock_count: jsonpath "$.count"


POST http://localhost:8084/order/create
{
    "user": 4,
    "items": [
        {
            "sku": 139275865,
            "count": 1
        }
    ]
}
HTTP 200

GET http://localhost:8084/stock/info?sku=139275865
HTTP 200
[Asserts]
jsonpath "$.count" != {{stock_count}}


POST http://localhost:8084/order/create
{
    "user": 5,
    "items": [
        {
            "sku": 139275865,
            "count": 99999999
        }
    ]
}
HTTP 400
