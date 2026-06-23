"""PY-6: Pydantic input/output models."""
from typing import Optional
from pydantic import BaseModel, Field, field_validator


class UserProfile(BaseModel):
    age: Optional[int] = Field(None, ge=0, le=130)
    gender: Optional[str] = None          # male | female | other
    weight_kg: Optional[float] = None


class DiagnoseRequest(BaseModel):
    symptoms: list[int] = Field(..., min_length=1)   # array of symptom_id
    user_profile: Optional[UserProfile] = None
    kategori: str = Field("pre_travel", pattern=r"^(pre_travel|post_travel)$")
    language: str = Field("id", pattern=r"^(id|en)$")  # NEW: language preference

    @field_validator("symptoms")
    @classmethod
    def symptoms_unique(cls, v: list[int]) -> list[int]:
        return list(dict.fromkeys(v))


class DiagnosisResult(BaseModel):
    disease_id: int
    disease_nama: str
    disease_nama_en: Optional[str] = None  # NEW: English disease name
    confidence_score: float
    risk_level: str
    recommendation: str
    language: str  # NEW: echo back requested language


class DiagnoseResponse(BaseModel):
    diagnosis: Optional[str]
    confidence_score: Optional[float]
    risk_level: Optional[str]
    recommendation: Optional[str]
    language: str  # NEW: echo back requested language
    all_results: list[DiagnosisResult] = []
