package receipt

import "time"

type Receipt struct {
	Restaurant  string         `json:"restaurant" jsonschema_description:"Name of the restaurant"`
	Address     string         `json:"address" jsonschema_description:"Address of the restaurant"`
	Opened      string         `json:"opened" jsonschema_description:"Date and time the order was opened"`
	OrderNumber string         `json:"order_number" jsonschema_description:"Unique order number"`
	OrderType   string         `json:"order_type" jsonschema_description:"Type of the order (e.g., dine-in, takeout)"`
	Table       string         `json:"table" jsonschema_description:"Table number or identifier"`
	Server      string         `json:"server" jsonschema_description:"Name or ID of the server"`
	Items       []OrderItem    `json:"items" jsonschema_description:"List of items ordered"`
	Subtotal    float64        `json:"subtotal" jsonschema_description:"Subtotal before tax"`
	SalesTax    float64        `json:"sales_tax" jsonschema_description:"Sales tax amount"`
	Total       float64        `json:"total" jsonschema_description:"Total amount of the order"`
	Payment     PaymentDetails `json:"payment" jsonschema_description:"Payment information"`
	Copy        string         `json:"copy" jsonschema_description:"Receipt copy type (e.g., customer, merchant)"`
	OtherFees   []OtherFee     `json:"other_fees" jsonschema_description:"List of additional fees applied to the order"`
}

type OrderItem struct {
	Name     string  `json:"name" jsonschema_description:"Name of the ordered item"`
	Price    float64 `json:"price" jsonschema_description:"Price of the ordered item"`
	Quantity int     `json:"quantity" jsonschema_description:"Quantity of the ordered item"`
}

type PaymentDetails struct {
	Method     string  `json:"method" jsonschema_description:"Payment method (e.g., cash, credit card)"`
	AmountPaid float64 `json:"amount_paid" jsonschema_description:"Total amount paid"`
	Tip        float64 `json:"tip" jsonschema_description:"Tip amount given"`
}

type OtherFee struct {
	Name  string  `json:"name" jsonschema_description:"Name of the additional fee"`
	Price float64 `json:"price" jsonschema_description:"Price of the additional fee"`
}

type ParsedReceipt struct {
	Receipt
	Opened time.Time
}
