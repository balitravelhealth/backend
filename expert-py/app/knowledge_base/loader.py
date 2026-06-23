"""PY-3 / PY-18: Load published rules from DB on every request (strategy A — always fresh)."""
from typing import Optional
import json

from app.database import get_connection
from app.knowledge_base.models import Rule


def load_rules(kategori: Optional[str] = None) -> list[Rule]:
    """Return all published rules, optionally filtered by kategori.

    Fetches both Indonesian and English disease names and recommendations.
    """
    query = """
        SELECT
            r.rule_id,
            r.nama,
            r.premis,
            r.disease_id,
            r.bobot_cf,
            r.mb,
            r.md,
            r.kategori,
            d.nama_id                      AS disease_nama,
            d.nama_en                      AS disease_nama_en,
            d.rekomendasi_default_id       AS rek_default_id,
            d.rekomendasi_default_en       AS rek_default_en
        FROM expert_rules r
        JOIN expert_diseases d ON d.id = r.disease_id
        WHERE r.status = 'published'
    """
    params: list = []
    if kategori:
        query += " AND r.kategori = %s"
        params.append(kategori)

    with get_connection() as conn:
        with conn.cursor() as cur:
            cur.execute(query, params)
            rows = cur.fetchall()

    rules = []
    for row in rows:
        premis = row["premis"]
        if isinstance(premis, str):
            premis = json.loads(premis)

        rek_id = row["rek_default_id"]
        if isinstance(rek_id, str):
            rek_id = json.loads(rek_id)

        rek_en = row["rek_default_en"]
        if isinstance(rek_en, str):
            rek_en = json.loads(rek_en)

        rules.append(Rule(
            rule_id=row["rule_id"],
            nama=row["nama"],
            premis=premis,
            disease_id=row["disease_id"],
            disease_nama=row["disease_nama"],
            disease_nama_en=row["disease_nama_en"],
            bobot_cf=float(row["bobot_cf"]),
            mb=float(row["mb"]),
            md=float(row["md"]),
            kategori=row["kategori"],
            rek_default_id=rek_id,
            rek_default_en=rek_en,
        ))
    return rules
