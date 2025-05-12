GET http://localhost:8084/stock/info?sku=1148162
HTTP 200
[Asserts]
jsonpath "$.count" == 180