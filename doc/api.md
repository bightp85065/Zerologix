# Order API 



## 1. Create

POST /api/place-order

## Request Body

The request body should be a JSON object containing the following fields:

- `product` (integer): Product ID. 
- `action` (integer): Specifies the action to be taken, either 1:"buy" or 2:"sell".
- `quantity` (integer): The quantity of the asset to be bought or sold.
- `price_type` (integer): Specifies the price type, either 1:"market" or 2:"limit".
- `price` (integer): The market price if `price_type` is "market", or the limit price if `price_type` is "limit".

### Example Request Body
```json
{
  "product": 1,
  "action": 1,
  "quantity": 10,
  "price_type": 2,
  "price": 150
}
```

## Response
Upon successful execution, the API returns a JSON response with the order details:

order_id (string): The unique identifier for the placed order.
create (time): Order generation time.
status (integer): The status of the order, e.g., 1:"pending", 2:"failure", 3:"cancelled", 4:"completed".

### Example Response
```json
{
  "order_id": "123456789",
  "create": time,
  "status": "pending"
}
```

## Error Response
If the request is invalid or encounters an error, the API returns an error response with details.

### Example Response
```json
{
  "err_type": 1,
  "msg": "Invalid request. Quantity must be a positive integer."
}
```

## 2. Cancel
## 3. Get List
## 4. Get One
