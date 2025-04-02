from typing import Dict, List, Optional, Tuple
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import text

from app.receipt.schema import Receipt


async def save_receipt(
    session: AsyncSession,
    hash: str,
    bucket: str,
    key: str,
    text_val: str,
    name: str,
    outing_id: str,
    receipt: Receipt,
):
    try:
        result = await session.execute(
            text("""
                INSERT INTO receipt_images (hash, bucket, key, raw_text, file_name, outing_id)
                VALUES (:hash, :bucket, :key, :raw_text, :file_name, :outing_id)
                RETURNING id
            """),
            {
                "hash": hash,
                "bucket": bucket,
                "key": key,
                "raw_text": text_val,
                "file_name": name,
                "outing_id": outing_id,
            },
        )
        image_id = result.scalar()

        # Insert into receipts
        result = await session.execute(
            text("""
                INSERT INTO receipts (
                    receipt_image_id, restaurant, address, opened,
                    order_number, order_type, table_number, server,
                    subtotal, sales_tax, total, copy
                )
                VALUES (
                    :receipt_image_id, :restaurant, :address, :opened,
                    :order_number, :order_type, :table_number, :server,
                    :subtotal, :sales_tax, :total, :copy
                )
                RETURNING id
            """),
            {
                "receipt_image_id": image_id,
                "restaurant": receipt.restaurant,
                "address": receipt.address,
                "opened": receipt.opened,
                "order_number": receipt.order_number,
                "order_type": receipt.order_type,
                "table_number": receipt.table,
                "server": receipt.server,
                "subtotal": receipt.subtotal,
                "sales_tax": receipt.sales_tax,
                "total": receipt.total,
                "copy": receipt.copy,
            },
        )
        receipt_id = result.scalar()

        # Insert into order_items
        for item in receipt.items:
            await session.execute(
                text("""
                    INSERT INTO order_items (receipt_id, name, price, quantity)
                    VALUES (:receipt_id, :name, :price, :quantity)
                """),
                {
                    "receipt_id": receipt_id,
                    "name": item.name,
                    "price": item.price,
                    "quantity": item.quantity,
                },
            )

        # Insert into other_fees
        for fee in receipt.other_fees:
            await session.execute(
                text("""
                    INSERT INTO other_fees (receipt_id, name, price)
                    VALUES (:receipt_id, :name, :price)
                """),
                {"receipt_id": receipt_id, "name": fee.name, "price": fee.price},
            )
        await session.commit()
    except Exception as e:
        raise RuntimeError(f"Error saving receipt: {e}")


async def get_receipt_by_hash(session: AsyncSession, hash: str) -> Optional[dict]:
    query = text("""
        SELECT 
            ri.*,
            COALESCE(oi.items, '[]') AS items,
            COALESCE(of.fees, '[]') AS fees
        FROM receipt_images ri
        JOIN receipts r ON ri.id = r.receipt_image_id
        LEFT JOIN (
            SELECT receipt_id, json_agg(oi.*) AS items
            FROM order_items oi
            GROUP BY receipt_id
        ) oi ON r.id = oi.receipt_id
        LEFT JOIN (
            SELECT receipt_id, json_agg(of.*) AS fees
            FROM other_fees of
            GROUP BY receipt_id
        ) of ON r.id = of.receipt_id
        WHERE ri.hash = :hash
        LIMIT 1
    """)

    result = await session.execute(query, {"hash": hash})
    row = result.fetchone()
    if row is None:
        return None

    # Convert result to dict
    return dict(row._mapping)


async def get_receipt_by_id(session: AsyncSession, id: str) -> Optional[dict]:
    query = text("""
        SELECT
            ri.raw_text,
            ri.bucket,
            ri.key,
            r.*,
            COALESCE(oi.items, '[]') AS items,
            COALESCE(of.fees, '[]') AS fees
        FROM receipt_images ri
        JOIN receipts r ON ri.id = r.receipt_image_id
        LEFT JOIN (
            SELECT receipt_id, json_agg(oi.*) AS items
            FROM order_items oi
            GROUP BY receipt_id
        ) oi ON r.id = oi.receipt_id
        LEFT JOIN (
            SELECT receipt_id, json_agg(of.*) AS fees
            FROM other_fees of
            GROUP BY receipt_id
        ) of ON r.id = of.receipt_id
        WHERE ri.id = :id
        LIMIT 1
    """)

    result = await session.execute(query, {"id": id})
    row = result.fetchone()
    if row is None:
        return None

    # Convert result to dict
    return dict(row._mapping)


async def get_receipts(session: AsyncSession) -> List[Dict[str, str]]:
    result = await session.execute(
        text("""
        SELECT ri.id, r.restaurant
        FROM receipt_images ri
        JOIN receipts r ON ri.id = r.receipt_image_id
    """)
    )
    rows = result.fetchall()
    return [{"id": str(row[0]), "restaurant": row[1]} for row in rows]
