"""TEST-4: Unit tests for CF formula and logic engine components."""
import pytest

from app.logic_engine.risk_calculator import cf_single, combine_cf, aggregate_cf_by_disease
from app.logic_engine.symptom_matcher import match_rules
from app.output_layer.risk_classifier import classify_risk
from app.knowledge_base.models import Rule


def make_rule(rule_id, premis, disease_id, disease_nama, mb, md) -> Rule:
    return Rule(
        rule_id=rule_id,
        nama=f"Rule {rule_id}",
        premis=premis,
        disease_id=disease_id,
        disease_nama=disease_nama,
        bobot_cf=mb - md,
        mb=mb,
        md=md,
        kategori="pre_travel",
    )


# ── CF formula ────────────────────────────────────────────────────────────────

def test_cf_single_basic():
    rule = make_rule(1, [1], 10, "Bali Belly", mb=0.8, md=0.1)
    assert abs(cf_single(rule) - 0.7) < 1e-9


def test_cf_single_zero_md():
    rule = make_rule(1, [1], 10, "Bali Belly", mb=0.6, md=0.0)
    assert abs(cf_single(rule) - 0.6) < 1e-9


def test_combine_cf_basic():
    # CFcombine = 0.6 + 0.5 * (1 - 0.6) = 0.6 + 0.2 = 0.8
    result = combine_cf(0.6, 0.5)
    assert abs(result - 0.8) < 1e-9


def test_combine_cf_zero_old():
    result = combine_cf(0.0, 0.7)
    assert abs(result - 0.7) < 1e-9


def test_combine_cf_saturates():
    # Two high-CF rules should not exceed 1.0
    cf = combine_cf(0.9, 0.9)
    assert cf <= 1.0


# ── Symptom Matcher ───────────────────────────────────────────────────────────

def test_match_rules_all_present():
    rules = [make_rule(1, [1, 2, 3], 10, "D", 0.8, 0.1)]
    matched = match_rules([1, 2, 3, 4], rules)
    assert len(matched) == 1


def test_match_rules_partial_premis_not_matched():
    rules = [make_rule(1, [1, 2, 3], 10, "D", 0.8, 0.1)]
    matched = match_rules([1, 2], rules)       # missing symptom 3
    assert len(matched) == 0


def test_match_rules_empty_input():
    rules = [make_rule(1, [1], 10, "D", 0.8, 0.1)]
    matched = match_rules([], rules)
    assert len(matched) == 0


def test_match_rules_multiple():
    rules = [
        make_rule(1, [1], 10, "Bali Belly", 0.7, 0.0),
        make_rule(2, [1, 2], 11, "DBD", 0.8, 0.1),
        make_rule(3, [3, 4], 12, "Rabies", 0.9, 0.0),
    ]
    matched = match_rules([1, 2], rules)
    assert {r.rule_id for r in matched} == {1, 2}


# ── Aggregate CF ──────────────────────────────────────────────────────────────

def test_aggregate_single_rule():
    rules = [make_rule(1, [1], 10, "Bali Belly", 0.8, 0.1)]
    agg = aggregate_cf_by_disease(rules)
    assert 10 in agg
    assert abs(agg[10][1] - 0.7) < 1e-9


def test_aggregate_two_rules_same_disease():
    rules = [
        make_rule(1, [1], 10, "Bali Belly", 0.7, 0.0),   # CF = 0.7
        make_rule(2, [2], 10, "Bali Belly", 0.6, 0.0),   # CF = 0.6
    ]
    agg = aggregate_cf_by_disease(rules)
    # CFcombine = 0.7 + 0.6*(1-0.7) = 0.7 + 0.18 = 0.88
    assert abs(agg[10][1] - 0.88) < 1e-9


# ── Risk Classifier ───────────────────────────────────────────────────────────

def test_classify_darurat():
    assert classify_risk(0.8) == "Darurat"
    assert classify_risk(1.0) == "Darurat"


def test_classify_tinggi():
    assert classify_risk(0.6) == "Tinggi"
    assert classify_risk(0.79) == "Tinggi"


def test_classify_sedang():
    assert classify_risk(0.4) == "Sedang"
    assert classify_risk(0.59) == "Sedang"


def test_classify_rendah():
    assert classify_risk(0.0) == "Rendah"
    assert classify_risk(0.39) == "Rendah"
