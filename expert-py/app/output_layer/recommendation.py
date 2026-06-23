"""PY-13: Build recommendation text from risk level and disease default."""
from typing import Optional


_DEFAULT_RECOMMENDATIONS_ID: dict[str, str] = {
    "Darurat": (
        "Segera kunjungi IGD atau hubungi layanan darurat 119. "
        "Jangan tunda penanganan medis."
    ),
    "Tinggi": (
        "Segera konsultasikan ke dokter atau klinik terdekat dalam 24 jam. "
        "Pantau kondisi secara ketat."
    ),
    "Sedang": (
        "Disarankan berkonsultasi dengan tenaga medis. "
        "Istirahat cukup dan minum air yang banyak."
    ),
    "Rendah": (
        "Pantau gejala selama 1–2 hari. Jika memburuk segera temui dokter. "
        "Pertahankan hidrasi dan istirahat."
    ),
}

_DEFAULT_RECOMMENDATIONS_EN: dict[str, str] = {
    "Darurat": (
        "Immediately visit the emergency room or call emergency services 119. "
        "Do not delay medical treatment."
    ),
    "Tinggi": (
        "Consult a doctor or nearest clinic within 24 hours. "
        "Monitor your condition closely."
    ),
    "Sedang": (
        "Consult a healthcare professional. "
        "Rest adequately and drink plenty of water."
    ),
    "Rendah": (
        "Monitor symptoms for 1–2 days. If worsening, see a doctor immediately. "
        "Stay hydrated and rest."
    ),
}


def build_recommendation(
    risk_level: str,
    rek_default: Optional[dict],
    language: str = "id"
) -> str:
    """Build recommendation in requested language.

    Args:
        risk_level: One of "Darurat", "Tinggi", "Sedang", "Rendah"
        rek_default: Disease-specific recommendations dict (if available)
        language: "id" (Indonesian) or "en" (English)

    Returns:
        Recommendation text in requested language
    """
    # Select default recommendations based on language
    defaults = (
        _DEFAULT_RECOMMENDATIONS_EN
        if language == "en"
        else _DEFAULT_RECOMMENDATIONS_ID
    )

    # Use disease-specific recommendation if available
    if rek_default and risk_level in rek_default:
        return str(rek_default[risk_level])

    # Fall back to generic default
    return defaults.get(risk_level, defaults["Rendah"])
