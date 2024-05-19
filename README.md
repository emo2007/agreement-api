# Agreement API mock datafeed

## Build and Run 
``` sh 
make d.build && make up
```

## API 

### GET /order/{order_id}/delivered
**Params:** 
* order_id (int)   

**Example:**
``` bash
curl --location 'http://localhost:8081/order/14/delivered'
```
**Response:**
Content-Type text/plain
```
0
```
1 (true) / 0 (false)