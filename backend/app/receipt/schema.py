from typing import List
from pydantic import BaseModel
from datetime import datetime


receipt_schema = {
    "type": "object",
    "properties": {
        "restaurant": {"type": "string", "description": "Name of the restaurant"},
        "address": {"type": "string", "description": "Address of the restaurant"},
        "opened": {
            "type": "string",
            "description": "Date and time the order was opened",
        },
        "order_number": {"type": "string", "description": "Unique order number"},
        "order_type": {
            "type": "string",
            "description": "Type of the order (e.g., dine-in, takeout)",
        },
        "table": {"type": "string", "description": "Table number or identifier"},
        "server": {"type": "string", "description": "Name or ID of the server"},
        "items": {
            "type": "array",
            "description": "List of items ordered",
            "items": {
                "type": "object",
                "properties": {
                    "name": {
                        "type": "string",
                        "description": "Name of the ordered item",
                    },
                    "price": {
                        "type": "number",
                        "description": "Price of the ordered item",
                    },
                    "quantity": {
                        "type": "integer",
                        "description": "Quantity of the ordered item",
                    },
                },
                "required": ["name", "price", "quantity"],
                "additionalProperties": False,
            },
        },
        "subtotal": {"type": "number", "description": "Subtotal before tax"},
        "sales_tax": {"type": "number", "description": "Sales tax amount"},
        "total": {"type": "number", "description": "Total amount of the order"},
        "payment": {
            "type": "object",
            "description": "Payment information",
            "properties": {
                "method": {
                    "type": "string",
                    "description": "Payment method (e.g., cash, credit card)",
                },
                "amount_paid": {"type": "number", "description": "Total amount paid"},
                "tip": {"type": "number", "description": "Tip amount given"},
            },
            "required": ["method", "amount_paid", "tip"],
            "additionalProperties": False,
        },
        "copy": {
            "type": "string",
            "description": "Receipt copy type (e.g., customer, merchant)",
        },
        "other_fees": {
            "type": "array",
            "description": "List of additional fees applied to the order",
            "items": {
                "type": "object",
                "properties": {
                    "name": {
                        "type": "string",
                        "description": "Name of the additional fee",
                    },
                    "price": {
                        "type": "number",
                        "description": "Price of the additional fee",
                    },
                },
                "required": ["name", "price"],
                "additionalProperties": False,
            },
        },
    },
    "required": [
        "restaurant",
        "address",
        "order_number",
        "order_type",
        "table",
        "server",
        "items",
        "subtotal",
        "sales_tax",
        "total",
        "payment",
        "copy",
        "other_fees",
        "opened",
    ],
    "additionalProperties": False,
}


class OrderItem(BaseModel):
    name: str
    price: float
    quantity: int


class PaymentInfo(BaseModel):
    method: str
    amount_paid: float
    tip: float


class OtherFee(BaseModel):
    name: str
    price: float


class Receipt(BaseModel):
    restaurant: str
    address: str
    opened: datetime | None
    order_number: str
    order_type: str
    table: str
    server: str
    items: List[OrderItem]
    subtotal: float
    sales_tax: float
    total: float
    payment: PaymentInfo
    copy: str
    other_fees: List[OtherFee]
