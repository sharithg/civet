import json
from uuid import UUID
from datetime import datetime, date
from decimal import Decimal
from typing import Any


def to_serializable(obj: Any) -> Any:
    if isinstance(obj, UUID):
        return str(obj)
    elif isinstance(obj, (datetime, date)):
        return obj.isoformat()
    elif isinstance(obj, Decimal):
        return float(obj)
    elif isinstance(obj, dict):
        return {k: to_serializable(v) for k, v in obj.items()}
    elif isinstance(obj, (list, tuple, set)):
        return [to_serializable(i) for i in obj]
    else:
        return obj


def row_to_dict(row: Any) -> dict:
    """Converts SQLAlchemy Row or row._mapping to a serializable dict."""
    mapping = row._mapping if hasattr(row, "_mapping") else row
    return {key: to_serializable(value) for key, value in mapping.items()}
