POST http://localhost:8084/order/create
{
    "user": 1,
    "items": [
        {
            "sku": 139275865,
            "count": 1
        },
        {
            "sku": 1076963,
            "count": 1
        }
    ]
}
HTTP 200
[Captures]
order_id: jsonpath "$.orderId"

GET http://localhost:8084/order/info?orderId={{order_id}}
HTTP 200
[Asserts]
jsonpath "$.user" == "1"
jsonpath "$.status" == "awaiting payment"
jsonpath "$.items[0].sku" == "1076963"
jsonpath "$.items[1].sku" == "139275865"

GET http://localhost:8084/order/info?orderId=999999999
HTTP 404

GET http://localhost:8084/order/info?orderId=0
HTTP 400
