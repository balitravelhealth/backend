"""PY-1: FastAPI entrypoint for the expert system service."""
from fastapi import FastAPI, HTTPException
from fastapi.responses import JSONResponse

from app.input_layer.schemas import DiagnoseRequest, DiagnoseResponse, DiagnosisResult
from app.logic_engine.engine import run as run_engine
from app.output_layer.risk_classifier import classify_risk
from app.output_layer.recommendation import build_recommendation
from app.knowledge_base.loader import load_rules

app = FastAPI(title="BaliTravelHealth Expert System", version="1.0.0")


# PY-2: Health-check (used by Gateway fail-safe GO-21)
@app.get("/health")
def health():
    return {"status": "ok"}


# PY-2: Extended health-check that verifies DB connectivity
@app.get("/health/db")
def health_db():
    try:
        rules = load_rules()
        return {"status": "ok", "published_rules": len(rules)}
    except Exception as e:
        return JSONResponse(
            status_code=503,
            content={"status": "error", "detail": str(e)},
        )


# PY-14: Main diagnosis endpoint
@app.post("/diagnose", response_model=DiagnoseResponse)
def diagnose(req: DiagnoseRequest):
    """
    Accept symptoms (array of symptom_id) + optional user_profile + language preference.
    Run Forward Chaining + CF engine against published rules.
    Return top diagnosis with confidence_score, risk_level, recommendation in requested language.

    Language: 'id' (Indonesian) or 'en' (English). Defaults to 'id'.
    """
    # PY-7: input already validated by Pydantic; empty symptoms caught by min_length=1
    candidates = run_engine(req.symptoms, req.user_profile, req.kategori)

    if not candidates:
        return DiagnoseResponse(
            diagnosis=None,
            confidence_score=None,
            risk_level=None,
            recommendation=None,
            language=req.language,  # Echo back requested language
            all_results=[],
        )

    # PY-11/12/13: build structured results for all candidates
    all_results = []
    for c in candidates:
        risk = classify_risk(c.cf)

        # Select recommendations based on language preference
        if req.language == "en":
            rek_default = c.rek_default_en
            disease_name = c.disease_nama_en or c.disease_nama  # Fallback to ID if EN not available
        else:
            rek_default = c.rek_default_id
            disease_name = c.disease_nama

        # Build recommendation in requested language
        rec = build_recommendation(risk, rek_default, language=req.language)

        all_results.append(DiagnosisResult(
            disease_id=c.disease_id,
            disease_nama=c.disease_nama,
            disease_nama_en=c.disease_nama_en,
            confidence_score=round(c.cf, 4),
            risk_level=risk,
            recommendation=rec,
            language=req.language,  # Echo back requested language
        ))

    top = all_results[0]
    top_disease_name = top.disease_nama_en if req.language == "en" else top.disease_nama
    top_disease_name = top_disease_name or top.disease_nama  # Fallback if EN not available

    return DiagnoseResponse(
        diagnosis=top_disease_name,
        confidence_score=top.confidence_score,
        risk_level=top.risk_level,
        recommendation=top.recommendation,
        language=req.language,  # Echo back requested language
        all_results=all_results,
    )
