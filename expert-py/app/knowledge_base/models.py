from dataclasses import dataclass, field
from typing import Optional


@dataclass
class Rule:
    rule_id: int
    nama: str
    premis: list[int]           # array of symptom_id
    disease_id: int
    disease_nama: str           # Indonesian disease name
    disease_nama_en: Optional[str] = None  # English disease name
    bobot_cf: float = 0.0
    mb: float = 0.0
    md: float = 0.0
    kategori: str = "pre_travel"  # pre_travel | post_travel
    rek_default_id: Optional[dict] = None  # Indonesian recommendations
    rek_default_en: Optional[dict] = None  # English recommendations
