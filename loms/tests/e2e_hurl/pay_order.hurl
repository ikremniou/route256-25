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
[Captures]
order_id: jsonpath "$.orderId"

POST http://localhost:8084/order/pay
{
    "order_id": {{order_id}}
}
HTTP 200

POST http://localhost:8084/order/pay
{
    "order_id": {{order_id}}
}
HTTP 200

POST http://localhost:8084/order/pay
{
    "order_id": 0
}
HTTP 400

POST http://localhost:8084/order/pay
{
    "order_id": 99999999999
}
HTTP 404
