"""PY-9: Certainty Factor calculation.

CF(H,E) = MB(H,E) − MD(H,E)
CFcombine = CFold + CFnew × (1 − CFold)
"""
from app.knowledge_base.models import Rule


def cf_single(rule: Rule) -> float:
    return rule.mb - rule.md


def combine_cf(cf_old: float, cf_new: float) -> float:
    return cf_old + cf_new * (1 - cf_old)


def aggregate_cf_by_disease(matched_rules: list[Rule]) -> dict[int, tuple[str, float]]:
    """Return {disease_id: (disease_nama, combined_cf)} for all matched rules."""
    disease_cf: dict[int, tuple[str, float]] = {}
    for rule in matched_rules:
        cf = cf_single(rule)
        if rule.disease_id not in disease_cf:
            disease_cf[rule.disease_id] = (rule.disease_nama, cf)
        else:
            nama, old_cf = disease_cf[rule.disease_id]
            disease_cf[rule.disease_id] = (nama, combine_cf(old_cf, cf))
    return disease_cf
