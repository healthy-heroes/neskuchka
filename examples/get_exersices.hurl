# получаем список упражнений
GET http://localhost:8000/api/exersices
HTTP 200


# получаем конкретное упражнение
GET http://localhost:8000/api/exersices/1
HTTP 200

GET http://localhost:8000/api/exersices/10
HTTP 200

[Asserts]
jsonpath "$.data" == null
