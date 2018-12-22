#Inventory

this is a API to handle inventory problems in a small shop.
I built this API using GO Language and Using gin for framework.

the features available in this API are :

-Catatan Jumlah Barang
    Stores actual stock of products
-Catatan Barang Masuk
    To store product that will be stored into the inventory.
-Catatan Barang Keluar
    To store product, quantity, notes of the products going out of inventory
-Laporan Nilai Barang
    Shows a report for owner to help her analyze and make decision. This report is related to total inventory value of Toko owner.
-Laporan Penjualan
    Shows a report for owner to help her analyze and make decision. This report is related to omzet / selling / profit.

Database Desain

#data product
m_product
m_product_id
name

#data pricelist
m_pricelist_id
name
m_currency_id 'IDR'
pricelist_type 'S' 'P'
is_active


#data price product
m_pricelist_line_id
m_product_price
m_product_id

#data Sales order / purchase order
t_order_id
order_type 'P' 'S'
order_amount
order_status 'P','I','C'

#data detail order
t_order_line_id
t_order_id
m_product_id
orderline_qty
orderline_price
orderline_total_amount
orderline_outstanding
orderline_received

#inventory location
m_inventory
m_inventory_id
name

#inventory data
m_inventory_line
m_inventory_line_id 
m_inventory_id
m_product_id
qty_count
last_update

#t_inout
t_inout_id //nomor dokument keluar masuk barang
inout_type = 'S' 'P' 'M'
t_order_line_id
m_product_id
inout_date
inout_qty

#list flag information

pricelist_type
S = 'Sales Order'
P = 'Purchase Order'

order_type
S = 'Sales Order'
P = 'Purchase Order'

inout_type
IN = 'IN'
OUT = 'OUT'

order_status
D = "Draft"
W = "Waiting for payment"
P = "Paid"
I = "In Progress"
C = "Complete"
X = "Canceled"

is_active
Y = "Yes"
N = "No"
