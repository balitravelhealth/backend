"""PY-12: Map confidence score to risk level."""


def classify_risk(cf: float) -> str:
    if cf >= 0.8:
        return "Darurat"
    if cf >= 0.6:
        return "Tinggi"
    if cf >= 0.4:
        return "Sedang"
    return "Rendah"
