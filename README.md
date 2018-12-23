# SSinventory

this is a API to handle inventory problems in a small shop.
I built this API using GO Language and Using gin for framework.

## Installation

Please install golang first in your system and then run:

```bash
go get github.com/crmspy/ssinventory
```

## Usage
 go to folder $GOPATH/github.com/crmspy/ssinventory then run this command:

```bash
#install all Install dependencies
go get -d -v

#or
go get ./...

#run application
go run main.go

#to build and run application
go build main.go && ./main
```

api service will automatically run on port localhost:8080
## Main Feature
1) Stores actual stock of products
2) To store product that will be stored into the inventory.
3) To store product, quantity, notes of the products going out of inventory
4) Shows a report for owner to help her analyze and make decision. This report is related to total inventory value of Toko owner.
5) Shows a report for owner to help her analyze and make decision. This report is related to omzet / selling / profit.

## Overview
1) Simple workflow that you must do is:
2) upload your stock OR
3) create new PO(Purchase Order) with "create order" then update stock at inventory using "inventory update stock"
4) Create SO(Sales Order) transaction 
5) Now you can Download all transaction using api in folder Report

## REST API
This is rest api documentation that available in SSinventory, If you wanna see full api you can see on [SSinventory Api Documentation](https://documenter.getpostman.com/view/2625111/Rzn8QMcJ) link

### @Report

Report Report Good Shipment
```curl
curl --location --request GET "localhost:8080/api/v1/inventory/goodshipment"
```

Report Sales Order
```curl
curl --location --request POST "localhost:8080/api/v1/inventory/salesorder" \
  --form "date_start=2018-12-20" \
  --form "date_end=2018-12-25"
```

Report Available Stock
```curl
curl --location --request GET "localhost:8080/api/v1/inventory/availablestock"
```

Report Value Of Product
```curl
curl --location --request GET "localhost:8080/api/v1/inventory/valueofproduct"
```

Report Good Receipt
```curl
curl --location --request GET "localhost:8080/api/v1/inventory/goodreceipt"
```

### @Transaction
Migration Data From .csv, format collumn that you must use is like this

| Product Code | Product Name      | Qty | Price | Inventory Location |
|--------------|-------------------|-----|-------|--------------------|
| P001         | Belgian Chocolate | 10  | 10000 | General            |
| P002         | Mongo Milk        | 5   | 50000 | General            |
|              |                   |     |       |                    |

If you did't set "Inventory Location" that will automatically set to "General". you can use sample data at folder sample/Format Upload CSV.csv

```
curl --location --request POST "localhost:8080/api/v1/inventory/migration" \
  --form "file=@"
```

Create Order
```
curl --location --request POST "localhost:8080/api/v1/order" \
  --form "t_order_id=S0004" \
  --form "description=pembelian sample product" \
  --form "detail={\"detail\":[{\"m_product_id\":\"SSI-D00791015-LL-BWH\",\"qty\":5,\"price\":100000},{\"m_product_id\":\"SSI-D01220307-XL-SAL\",\"qty\":10,\"price\":175000}]}" \
  --form "order_type=S" \
  --form "order_status=P"
```

Update Stock In Inventory
```
curl --location --request POST "localhost:8080/api/v1/inventory/inout" \
  --form "m_inventory_id=General" \
  --form "qty=5" \
  --form "t_order_line_id=64" \
  --form "description=diterima sebagian"
```
### @Master Data
insert new product or inventory location

Insert Product
```
curl --location --request POST "localhost:8080/api/v1/product" \
  --form "m_product_id=UK12" \
  --form "name=this is my product"
```

Add New Inventory Location
```
curl --location --request POST "localhost:8080/api/v1/inventory" \
  --form "m_inventory_id=MAINMain" \
  --form "name=Main Inventory"
```
## Database Desain & Flag
This is database schema in ssinventory

**Database**
```
#m_product
m_product_id
name

#m_pricelist
m_pricelist_id
name
m_currency_id ('IDR','USD')
pricelist_type ('P','S')
is_active ('Y','N')

#m_pricelist_line_id
m_pricelist_line_id
m_product_price
m_product_id

#t_order
t_order_id
order_type ('P' 'S')
order_amount
order_status ('D','W','P','I','C','X')

#t_order_line
t_order_line_id
t_order_id
m_product_id
orderline_qty
orderline_price
orderline_total_amount
orderline_outstanding
orderline_received

#m_inventory
m_inventory_id
name

#m_inventory_line
m_inventory_line_id 
m_inventory_id
m_product_id
qty_count
last_update

#t_inout
t_inout_id #you can fill this as no receipt document
inout_type  ('S' 'P')
t_order_line_id
m_product_id
inout_date
inout_qty
description
```

**Database**

Flag Information in every table
```
#PRICELIST TYPE
S = "Sales Order"
P = "Purchase Order"

#Order Type
S = "Sales Order"
P = "Purchase Order"

#INOUT TYPE
IN = "IN"
OUT = "OUT"

#ORDER STATUS
D = "Draft"
W = "Waiting for payment"
P = "Paid"
I = "In Progress"
C = "Complete"
X = "Canceled"

#IS ACTIVE
Y = "Yes"
N = "No"
```
## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
