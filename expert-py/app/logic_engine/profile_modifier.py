"""PY-10: Adjust CF scores based on user profile risk factors."""
from app.input_layer.schemas import UserProfile

# Disease IDs that receive a risk boost for elderly patients
_ELDERLY_RISK_DISEASE_KEYWORDS = {"jantung", "stroke", "heat", "hipertensi"}


def apply_profile(
    disease_cf: dict[int, tuple[str, float]],
    profile: UserProfile | None,
) -> dict[int, tuple[str, float]]:
    if not profile:
        return disease_cf

    result = {}
    for disease_id, (nama, cf) in disease_cf.items():
        adjusted = cf
        # Elderly (>= 60) → boost CF for cardiovascular/heat-related diseases
        if profile.age is not None and profile.age >= 60:
            if any(kw in nama.lower() for kw in _ELDERLY_RISK_DISEASE_KEYWORDS):
                adjusted = min(1.0, cf + 0.1)
        result[disease_id] = (nama, adjusted)
    return result
