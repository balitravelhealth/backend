"""PY-8: Forward Chaining — match input symptoms against rule premis."""
from app.knowledge_base.models import Rule


def match_rules(symptom_ids: list[int], rules: list[Rule]) -> list[Rule]:
    """Return rules whose entire premis is covered by the input symptoms."""
    symptom_set = set(symptom_ids)
    return [r for r in rules if set(r.premis).issubset(symptom_set)]
