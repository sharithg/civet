from typing import Dict, List
from sqlalchemy import text
from sqlalchemy.ext.asyncio import AsyncSession


async def create_new_outing(session: AsyncSession, name: str):
    result = await session.execute(
        text("""
                INSERT INTO outings (name)
                VALUES (:name)
                RETURNING id
            """),
        {"name": name},
    )
    outing_id = result.scalar()
    await session.commit()
    return outing_id


async def get_receipts_for_outing(
    session: AsyncSession, outing_id: str
) -> List[Dict[str, str]]:
    query = text(
        """
        select r.restaurant, count(oi.id) as order_count, r.total, r.id from receipts r
        join order_items oi on r.id = oi.receipt_id
        join receipt_images ri on r.receipt_image_id = ri.id
        where ri.outing_id = :outing_id
        group by r.id
        """,
    )
    result = await session.execute(query, {"outing_id": outing_id})
    rows = result.fetchall()
    return [
        {
            "restaurant": row[0],
            "order_count": row[1],
            "total": row[2],
            "id": str(row[3]),
        }
        for row in rows
    ]


async def get_outings(session: AsyncSession) -> List[Dict[str, str]]:
    query = text(
        """
        select o.id, o.name, o.created_at, count(ri.id) as total_receipts from outings o
        left join receipt_images ri on o.id = ri.outing_id
        group by o.id
        """,
    )
    result = await session.execute(query)
    rows = result.fetchall()
    return [
        {
            "id": str(row[0]),
            "name": row[1],
            "created_at": row[2],
            "total_receipts": row[3],
        }
        for row in rows
    ]
