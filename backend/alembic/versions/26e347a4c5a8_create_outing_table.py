"""create outing table

Revision ID: 26e347a4c5a8
Revises: 2ebcc8f8c208
Create Date: 2025-04-01 19:07:34.884248

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
import sqlalchemy.dialects.postgresql as pg


# revision identifiers, used by Alembic.
revision: str = "26e347a4c5a8"
down_revision: Union[str, None] = "2ebcc8f8c208"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    """Upgrade schema."""
    op.create_table(
        "outings",
        sa.Column(
            "id",
            pg.UUID(as_uuid=True),
            primary_key=True,
            server_default=sa.text("gen_random_uuid()"),
        ),
        sa.Column("name", sa.String(255), nullable=False),
        sa.Column(
            "created_at",
            sa.TIMESTAMP(timezone=True),
            server_default=sa.text("NOW()"),
            nullable=False,
        ),
        sa.Column(
            "updated_at",
            sa.TIMESTAMP(timezone=True),
            server_default=sa.text("NOW()"),
            nullable=False,
        ),
        sa.Column("deleted_at", sa.TIMESTAMP(timezone=True), nullable=True),
    )
    op.add_column(
        "receipt_images",
        sa.Column("outing_id", sa.UUID(), sa.ForeignKey("outings.id"), nullable=True),
    )


def downgrade() -> None:
    """Downgrade schema."""
    op.drop_column("receipt_images", "outing_id")
    op.drop_table("outings")
