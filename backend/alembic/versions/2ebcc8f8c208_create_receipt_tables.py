"""create receipt tables

Revision ID: 2ebcc8f8c208
Revises:
Create Date: 2025-04-01 12:26:58.366158

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
import sqlalchemy.dialects.postgresql as pg


# revision identifiers, used by Alembic.
revision: str = "2ebcc8f8c208"
down_revision: Union[str, None] = None
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    """Upgrade schema."""
    op.create_table(
        "receipt_images",
        sa.Column(
            "id",
            pg.UUID(as_uuid=True),
            primary_key=True,
            server_default=sa.text("gen_random_uuid()"),
        ),
        sa.Column("bucket", sa.String(255), nullable=False),
        sa.Column("key", sa.String(255), nullable=False),
        sa.Column("raw_text", sa.Text, nullable=False),
        sa.Column("file_name", sa.String(255), nullable=False),
        sa.Column("hash", sa.Text, nullable=False),
    )

    op.create_table(
        "receipts",
        sa.Column(
            "id",
            pg.UUID(as_uuid=True),
            primary_key=True,
            server_default=sa.text("gen_random_uuid()"),
        ),
        sa.Column(
            "receipt_image_id",
            pg.UUID(as_uuid=True),
            sa.ForeignKey("receipt_images.id", ondelete="CASCADE"),
            nullable=False,
        ),
        sa.Column("restaurant", sa.String(255), nullable=False),
        sa.Column("address", sa.Text),
        sa.Column("opened", sa.TIMESTAMP),
        sa.Column("order_number", sa.String(50)),
        sa.Column("order_type", sa.String(50)),
        sa.Column("table_number", sa.String(50)),
        sa.Column("server", sa.String(255)),
        sa.Column("subtotal", sa.Numeric(10, 2), nullable=False),
        sa.Column("sales_tax", sa.Numeric(10, 2), nullable=False),
        sa.Column("total", sa.Numeric(10, 2), nullable=False),
        sa.Column("payment_method", sa.String(50)),
        sa.Column("payment_amount_paid", sa.Numeric(10, 2)),
        sa.Column("payment_tip", sa.Numeric(10, 2)),
        sa.Column("copy", sa.String(50)),
        sa.Column(
            "created_at", sa.TIMESTAMP, server_default=sa.text("CURRENT_TIMESTAMP")
        ),
    )

    op.create_table(
        "order_items",
        sa.Column(
            "id",
            pg.UUID(as_uuid=True),
            primary_key=True,
            server_default=sa.text("gen_random_uuid()"),
        ),
        sa.Column(
            "receipt_id",
            pg.UUID(as_uuid=True),
            sa.ForeignKey("receipts.id", ondelete="CASCADE"),
            nullable=False,
        ),
        sa.Column("name", sa.String(255), nullable=False),
        sa.Column("price", sa.Numeric(10, 2), nullable=False),
        sa.Column("quantity", sa.Integer, nullable=False),
    )

    op.create_table(
        "other_fees",
        sa.Column(
            "id",
            pg.UUID(as_uuid=True),
            primary_key=True,
            server_default=sa.text("gen_random_uuid()"),
        ),
        sa.Column(
            "receipt_id",
            pg.UUID(as_uuid=True),
            sa.ForeignKey("receipts.id", ondelete="CASCADE"),
            nullable=False,
        ),
        sa.Column("name", sa.String(255), nullable=False),
        sa.Column("price", sa.Numeric(10, 2), nullable=False),
    )


def downgrade() -> None:
    """Downgrade schema."""
    op.drop_table("other_fees")
    op.drop_table("order_items")
    op.drop_table("receipts")
    op.drop_table("receipt_images")
