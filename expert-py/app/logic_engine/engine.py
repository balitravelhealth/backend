"""PY-11: Combine all logic steps and produce ranked diagnosis candidates."""
from dataclasses import dataclass

from app.knowledge_base.models import Rule
from app.input_layer.schemas import UserProfile
from app.logic_engine.symptom_matcher import match_rules
from app.logic_engine.risk_calculator import aggregate_cf_by_disease
from app.logic_engine.profile_modifier import apply_profile
from app.knowledge_base.loader import load_rules


@dataclass
class Candidate:
    disease_id: int
    disease_nama: str
    disease_nama_en: str | None
    cf: float
    rek_default_id: dict | None = None  # Indonesian recommendations
    rek_default_en: dict | None = None  # English recommendations


def run(symptom_ids: list[int], profile: UserProfile | None, kategori: str) -> list[Candidate]:
    rules = load_rules(kategori)
    matched = match_rules(symptom_ids, rules)
    disease_cf = aggregate_cf_by_disease(matched)
    disease_cf = apply_profile(disease_cf, profile)

    # Map disease_id to full disease info from rules for response building
    disease_info: dict[int, tuple[str, str | None, dict | None, dict | None]] = {}
    for rule in rules:
        if rule.disease_id not in disease_info:
            disease_info[rule.disease_id] = (
                rule.disease_nama,
                rule.disease_nama_en,
                rule.rek_default_id,
                rule.rek_default_en,
            )

    candidates = [
        Candidate(
            disease_id=did,
            disease_nama=nama,
            disease_nama_en=disease_info.get(did, (nama, None, None, None))[1],
            cf=cf,
            rek_default_id=disease_info.get(did, (nama, None, None, None))[2],
            rek_default_en=disease_info.get(did, (nama, None, None, None))[3],
        )
        for did, (nama, cf) in disease_cf.items()
    ]
    candidates.sort(key=lambda c: c.cf, reverse=True)
    return candidates
