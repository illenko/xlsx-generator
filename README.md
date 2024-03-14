## Service for xlsx files generation.
### Stack:

* Xlsx Generation: https://github.com/tealeg/xlsx
* DI: https://github.com/uber-go/fx
* Web Handlers: https://github.com/gin-gonic/gin

### Commands:

#### Docker build
````shell
docker build -t xsls-generator:latest .
````

#### Docker run
````shell
docker run -d -t -i -p 8080:8080 --name xsls-generator xsls-generator:latest
````

### Request example:
```json
{
  "sheets": [
    {
      "name": "Payments",
      "additionalInfo": {
        "top": [
          {
            "title": "User",
            "value": "test@gmail.com"
          },
          {
            "title": "Period",
            "value": "01.02.2024-29.02.2024"
          }
        ],
        "bottom": [
          {
            "title": "Total amount",
            "value": "210.0"
          }
        ]
      },
      "columns": [
        {
          "id": "payment_id",
          "title": "Payment ID",
          "type": "string",
          "color": {
            "background": "8FC8DC"
          }
        },
        {
          "id": "category",
          "title": "Category",
          "type": "string",
          "color": {
            "font": "1D1E1A",
            "background": "C4DC8F"
          }
        },
        {
          "id": "amount",
          "title": "Amount",
          "type": "number"
        }
      ],
      "data": [
        {
          "payment_id": "1",
          "category": "P2P Transfer",
          "amount": "100.1"
        },
        {
          "payment_id": "2",
          "category": "Taxi",
          "amount": "50.0"
        },
        {
          "payment_id": "3",
          "category": "Food",
          "amount": "60.0"
        }
      ]
    }
  ]
}
```