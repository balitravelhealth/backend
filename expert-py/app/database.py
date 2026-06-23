"""PostgreSQL read-only connection for the expert system (PY-2b)."""
import os
import psycopg
from psycopg.rows import dict_row


def get_dsn() -> str:
    return (
        f"host={os.getenv('DB_HOST', 'db')} "
        f"port={os.getenv('DB_PORT', '5432')} "
        f"dbname={os.getenv('POSTGRES_DB', 'balitravelhealth')} "
        f"user={os.getenv('POSTGRES_USER', 'balitravelhealthdb')} "
        f"password={os.getenv('POSTGRES_PASSWORD', '')} "
        "connect_timeout=5"
    )


def get_connection() -> psycopg.Connection:
    return psycopg.connect(get_dsn(), row_factory=dict_row)
