# SSinventory

this is a API to handle inventory problems in a small shop.
I built this API using GO Language and Using gin for framework.

## Installation

Please install golang first in your system and then run:

```bash
go get github.com/crmspy/ssinventory
```

## Usage

```bash
go run main.go
```

if you a developer and want it run with autoload when files change, run:
```bash
./fresh
```

## Main Feature
1) Stores actual stock of products
2) To store product that will be stored into the inventory.
3) To store product, quantity, notes of the products going out of inventory
4) Shows a report for owner to help her analyze and make decision. This report is related to total inventory value of Toko owner.
5) Shows a report for owner to help her analyze and make decision. This report is related to omzet / selling / profit.

## REST API
This is rest api documentation that available in SSinventory



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
m_inventory_line
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
