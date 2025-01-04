# получаем список упражнений
GET http://localhost:8000/api/exercises
HTTP 200


# получаем конкретное упражнение
GET http://localhost:8000/api/exercises/1
HTTP 200

GET http://localhost:8000/api/exercises/10
HTTP 200

[Asserts]
jsonpath "$.data" == null
