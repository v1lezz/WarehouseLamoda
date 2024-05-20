# For start
make run

OR

docker-compose up -d

# Requests and answers

#### Request
curl -X GET http://localhost:9000/getGood?goodID=1
#### Answer
{"data":null,"error":"good with this id is not found"}

#### Request
curl -X POST -d '{"id":1,"name":"good1","size":1}' http://localhost:9000/createGood
#### Answer
{"data": {"name": "good1","size": 1,"id": 1,"warehouses": null },"error": null}

#### Request
curl -X GET http://localhost:9000/getGood?goodID=1
#### Answer
{"data":{"name":"good1","size":1,"id":1,"warehouses":[]},"error":null}

#### Request
curl -X PUT -d '{"id": 1,"name": "good2","size": 2}' http://localhost:9000/updateGood
#### Answer
{"data": {"name": "good2","size": 2,"id": 1,"warehouses": null },"error": null }

#### Request
curl -X DELETE http://localhost:9000/deleteGood?goodID=1
#### Answer
{"data":null,"error":null}

#### Request
curl -X POST -d '{"id":1,"name":"ws1","is_available":true}' http://localhost:9000/createWarehouse
#### Answer
{"data":{"id":1,"name":"ws1","is_available":true,"goods":null},"error":null}

#### Request
curl -X POST http://localhost:9000/addGoodOnWarehouse?goodID=1&warehouseID=1&count=10
#### Answer
{"data":null,"error":null}

#### Request
curl -X GET http://localhost:9000/getCountGoods?warehouseID=1
#### Answer
{"data":{"count":10},"error":null}

#### Request
curl -X PATCH -d '[{"good_id":1,"warehouse_id":1}]' http://localhost:9000/reserveGood
#### Answer
{"data":{"reserved":[1],"error_reservation":[]},"error":null}

#### Request
curl -X GET http://localhost:9000/getGood?goodID=1
#### Answer
{"data":{"name":"good1","size":1,"id":1,"warehouses":[{"name":"ws1","is_available":false,"count":10,"reserved":1}]},"error":null}

#### Request
curl -X PATCH -d '[{"good_id":1,"warehouse_id":1}]' http://localhost:9000/releaseReservationGood
#### Answer
{"data":{"released":[1],"error_release":[]},"error":null}

#### Request
curl -X GET http://localhost:9000/getGood?goodID=1
#### Answer
{"data":{"name":"good1","size":1,"id":1,"warehouses":[{"name":"ws1","is_available":false,"count":10,"reserved":0}]},"error":null}
