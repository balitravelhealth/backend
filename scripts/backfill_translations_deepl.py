#!/usr/bin/env python3
"""
Backfill English translations using DeepL API.

Usage:
    python backfill_translations_deepl.py --api-key YOUR_DEEPL_KEY

Environment variables:
    DATABASE_URL: PostgreSQL connection string
    DEEPL_API_KEY: DeepL API key (or pass --api-key)
"""

from dotenv import load_dotenv

load_dotenv()

import argparse
import json
import os
import sys
import time
from typing import Any, Dict, List, Optional

import psycopg2
import requests
from psycopg2.extras import RealDictCursor

DEEPL_API_URL = "https://api-free.deepl.com/v2/translate"  # Free tier
DEEPL_API_URL_PRO = "https://api.deepl.com/v2/translate"  # Pro tier


def get_deepl_translator(api_key: str, is_pro: bool = False) -> callable:
    """Create a DeepL translator function."""
    url = DEEPL_API_URL_PRO if is_pro else DEEPL_API_URL

    def translate(text: str, target_lang: str = "EN") -> str:
        """Translate text to target language using DeepL."""
        if not text:
            return text

        # Ensure text is a string
        if not isinstance(text, str):
            text = str(text)

        if not text.strip():
            return text

        # Sanitize: remove null bytes and control characters
        text = text.replace("\x00", "").strip()
        if not text:
            return ""

        # Truncate very long texts to avoid API limits
        max_chars = 10000
        if len(text) > max_chars:
            text = text[:max_chars]

        payload = {
            "text": [text],  # v2 API requires array
            "target_lang": target_lang,
        }
        headers = {
            "Authorization": f"DeepL-Auth-Key {api_key}",
            "Content-Type": "application/json",
        }

        try:
            response = requests.post(url, json=payload, headers=headers, timeout=10)

            # Better error handling
            if response.status_code == 400:
                # Check if it's "Value for 'text' not supported" error
                try:
                    error_info = response.json()
                    if "Value for 'text' not supported" in str(
                        error_info.get("message", "")
                    ):
                        print(
                            f"⚠️  Skipping unsupported text: '{text[:80]}...'",
                            file=sys.stderr,
                        )
                        return text  # Return sanitized version
                except:
                    pass

                print(f"DeepL 400 error for text: '{text[:100]}...'", file=sys.stderr)
                print(f"Response: {response.text}", file=sys.stderr)
                return text
            elif response.status_code == 401:
                print(f"DeepL 401 Unauthorized - check your API key", file=sys.stderr)
                return text
            elif response.status_code == 429:
                print(f"DeepL Rate limit hit - backing off", file=sys.stderr)
                time.sleep(2)
                return text

            response.raise_for_status()
            result = response.json()
            return result["translations"][0]["text"]
        except Exception as e:
            print(f"Translation error for '{text[:50]}...': {e}", file=sys.stderr)
            return text  # Return original on error

    return translate


class TranslationBackfiller:
    """Backfill translations in the database."""

    def __init__(self, db_url: str, translator: callable):
        self.db_url = db_url
        self.translator = translator
        self.stats = {
            "symptoms_translated": 0,
            "diseases_translated": 0,
            "health_risks_translated": 0,
            "emergency_guides_translated": 0,
            "emergency_flows_translated": 0,
            "total_requests": 0,
            "errors": 0,
        }

    def connect(self):
        """Connect to database."""
        self.conn = psycopg2.connect(self.db_url)
        self.cursor = self.conn.cursor(cursor_factory=RealDictCursor)

    def close(self):
        """Close database connection."""
        if self.conn:
            self.conn.close()

    def commit(self):
        """Commit changes."""
        self.conn.commit()

    def translate_text(self, text: Optional[str]) -> Optional[str]:
        """Translate a single text, with rate limiting."""
        if not text:
            return None

        self.stats["total_requests"] += 1

        # Rate limiting: DeepL free tier allows 50 requests/minute
        if self.stats["total_requests"] % 50 == 0:
            print(
                f"Rate limiting: pausing 1 second after {self.stats['total_requests']} requests"
            )
            time.sleep(1)

        result = self.translator(text)
        return result if result and result != text else None

    def backfill_expert_symptoms(self):
        """Backfill expert_symptoms translations."""
        print("\n📝 Translating expert_symptoms...")

        self.cursor.execute("""
            SELECT symptom_id, label_id, deskripsi_id, label_en, deskripsi_en
            FROM expert_symptoms
            ORDER BY symptom_id
        """)

        symptoms = self.cursor.fetchall()

        for i, sym in enumerate(symptoms, 1):
            # Translate label if label_en is not set
            if sym["label_en"] != sym["label_id"]:
                # Already has proper translation
                continue

            try:
                label_en = self.translate_text(sym["label_id"])
                deskripsi_en = (
                    self.translate_text(sym["deskripsi_id"])
                    if sym["deskripsi_id"]
                    else None
                )

                if label_en or deskripsi_en:
                    self.cursor.execute(
                        """
                        UPDATE expert_symptoms
                        SET label_en = %s, deskripsi_en = %s
                        WHERE symptom_id = %s
                    """,
                        (label_en or sym["label_id"], deskripsi_en, sym["symptom_id"]),
                    )

                    self.stats["symptoms_translated"] += 1
                    print(
                        f"  [{i}/{len(symptoms)}] {sym['label_id'][:40]} -> {label_en[:40] if label_en else 'N/A'}"
                    )

            except Exception as e:
                self.stats["errors"] += 1
                print(f"  ❌ Error translating symptom {sym['symptom_id']}: {e}")

        self.commit()

    def backfill_expert_diseases(self):
        """Backfill expert_diseases translations."""
        print("\n🏥 Translating expert_diseases...")

        self.cursor.execute("""
            SELECT id, nama_id, deskripsi_id, rekomendasi_default_id, nama_en
            FROM expert_diseases
            ORDER BY id
        """)

        diseases = self.cursor.fetchall()

        for i, dis in enumerate(diseases, 1):
            # Translate nama if nama_en is not set
            if dis["nama_en"] != dis["nama_id"]:
                # Already has proper translation
                continue

            try:
                nama_en = self.translate_text(dis["nama_id"])
                deskripsi_en = (
                    self.translate_text(dis["deskripsi_id"])
                    if dis["deskripsi_id"]
                    else None
                )

                # For JSONB recommendations, translate each risk level
                rekomendasi_en = None
                if dis["rekomendasi_default_id"] and isinstance(
                    dis["rekomendasi_default_id"], dict
                ):
                    rekomendasi_en = {}
                    for risk_level, recommendation in dis[
                        "rekomendasi_default_id"
                    ].items():
                        rekomendasi_en[risk_level] = self.translate_text(recommendation)

                if nama_en or deskripsi_en or rekomendasi_en:
                    rekomendasi_en_json = (
                        json.dumps(rekomendasi_en) if rekomendasi_en else None
                    )

                    self.cursor.execute(
                        """
                        UPDATE expert_diseases
                        SET nama_en = %s, deskripsi_en = %s, rekomendasi_default_en = %s
                        WHERE id = %s
                    """,
                        (
                            nama_en or dis["nama_id"],
                            deskripsi_en,
                            rekomendasi_en_json,
                            dis["id"],
                        ),
                    )

                    self.stats["diseases_translated"] += 1
                    print(
                        f"  [{i}/{len(diseases)}] {dis['nama_id'][:40]} -> {nama_en[:40] if nama_en else 'N/A'}"
                    )

            except Exception as e:
                self.stats["errors"] += 1
                error_msg = str(e)
                if "Value for 'text' not supported" in error_msg:
                    print(
                        f"  ⚠️  Skipped disease '{dis['nama_id'][:40]}' (unsupported format)"
                    )
                else:
                    print(f"  ❌ Error translating disease {dis['id']}: {e}")

        self.commit()

    def backfill_health_risks(self):
        """Backfill health_risks translations."""
        print("\n⚠️  Translating health_risks...")

        self.cursor.execute("""
            SELECT id, nama_risiko_id, saran_pencegahan_id, rekomendasi_vaksinasi_id, nama_risiko_en
            FROM health_risks
            ORDER BY id
        """)

        risks = self.cursor.fetchall()

        for i, risk in enumerate(risks, 1):
            # Check if already translated
            if risk["nama_risiko_en"] != risk["nama_risiko_id"]:
                continue

            try:
                nama_en = self.translate_text(risk["nama_risiko_id"])
                saran_en = (
                    self.translate_text(risk["saran_pencegahan_id"])
                    if risk["saran_pencegahan_id"]
                    else None
                )
                vaksin_en = (
                    self.translate_text(risk["rekomendasi_vaksinasi_id"])
                    if risk["rekomendasi_vaksinasi_id"]
                    else None
                )

                if nama_en or saran_en or vaksin_en:
                    self.cursor.execute(
                        """
                        UPDATE health_risks
                        SET nama_risiko_en = %s, saran_pencegahan_en = %s, rekomendasi_vaksinasi_en = %s
                        WHERE id = %s
                    """,
                        (
                            nama_en or risk["nama_risiko_id"],
                            saran_en,
                            vaksin_en,
                            risk["id"],
                        ),
                    )

                    self.stats["health_risks_translated"] += 1
                    print(
                        f"  [{i}/{len(risks)}] {risk['nama_risiko_id'][:40]} -> {nama_en[:40] if nama_en else 'N/A'}"
                    )

            except Exception as e:
                self.stats["errors"] += 1
                print(f"  ❌ Error translating health risk {risk['id']}: {e}")

        self.commit()

    def backfill_emergency_guides(self):
        """Backfill emergency_guides translations (JSONB)."""
        print("\n🚨 Translating emergency_guides...")

        self.cursor.execute("""
            SELECT id, kategori, langkah, isi_media_id, isi_media_en
            FROM emergency_guides
            ORDER BY id
        """)

        guides = self.cursor.fetchall()

        for i, guide in enumerate(guides, 1):
            # Check if already translated
            if guide["isi_media_en"] != guide["isi_media_id"]:
                continue

            try:
                isi_en = None
                if guide["isi_media_id"] and isinstance(guide["isi_media_id"], dict):
                    isi_en = {}
                    for key, value in guide["isi_media_id"].items():
                        if isinstance(value, str):
                            isi_en[key] = self.translate_text(value)
                        else:
                            isi_en[key] = value

                if isi_en:
                    isi_en_json = json.dumps(isi_en)

                    self.cursor.execute(
                        """
                        UPDATE emergency_guides
                        SET isi_media_en = %s
                        WHERE id = %s
                    """,
                        (isi_en_json, guide["id"]),
                    )

                    self.stats["emergency_guides_translated"] += 1
                    print(
                        f"  [{i}/{len(guides)}] Emergency guide {guide['kategori']} step {guide['langkah']}"
                    )

            except Exception as e:
                self.stats["errors"] += 1
                print(f"  ❌ Error translating emergency guide {guide['id']}: {e}")

        self.commit()

    def backfill_emergency_guide_flows(self):
        """Backfill emergency_guide_flows translations (nested JSONB)."""
        print("\n🔄 Translating emergency_guide_flows...")

        self.cursor.execute("""
            SELECT id, title_id, nodes_id, title_en
            FROM emergency_guide_flows
            ORDER BY id
        """)

        flows = self.cursor.fetchall()

        for i, flow in enumerate(flows, 1):
            # Check if already translated
            if flow["title_en"] != flow["title_id"]:
                continue

            try:
                title_en = self.translate_text(flow["title_id"])

                nodes_en = None
                if flow["nodes_id"] and isinstance(flow["nodes_id"], list):
                    nodes_en = []
                    for node in flow["nodes_id"]:
                        if isinstance(node, dict):
                            node_en = node.copy()
                            # Translate relevant fields
                            if "title" in node:
                                node_en["title"] = self.translate_text(node["title"])
                            if "instruction" in node:
                                node_en["instruction"] = self.translate_text(
                                    node["instruction"]
                                )
                            if "choices" in node and isinstance(node["choices"], list):
                                node_en["choices"] = [
                                    {
                                        **choice,
                                        "label": self.translate_text(
                                            choice.get("label", "")
                                        ),
                                    }
                                    for choice in node["choices"]
                                ]
                            nodes_en.append(node_en)
                        else:
                            nodes_en.append(node)

                if title_en or nodes_en:
                    nodes_en_json = json.dumps(nodes_en) if nodes_en else None

                    self.cursor.execute(
                        """
                        UPDATE emergency_guide_flows
                        SET title_en = %s, nodes_en = %s
                        WHERE id = %s
                    """,
                        (title_en or flow["title_id"], nodes_en_json, flow["id"]),
                    )

                    self.stats["emergency_flows_translated"] += 1
                    print(
                        f"  [{i}/{len(flows)}] {flow['title_id'][:40]} -> {title_en[:40] if title_en else 'N/A'}"
                    )

            except Exception as e:
                self.stats["errors"] += 1
                print(f"  ❌ Error translating emergency guide flow {flow['id']}: {e}")

        self.commit()

    def print_stats(self):
        """Print translation statistics."""
        print("\n" + "=" * 60)
        print("📊 TRANSLATION SUMMARY")
        print("=" * 60)
        print(f"Symptoms translated:          {self.stats['symptoms_translated']}")
        print(f"Diseases translated:          {self.stats['diseases_translated']}")
        print(f"Health risks translated:      {self.stats['health_risks_translated']}")
        print(
            f"Emergency guides translated:  {self.stats['emergency_guides_translated']}"
        )
        print(
            f"Emergency flows translated:   {self.stats['emergency_flows_translated']}"
        )
        print(f"Total API requests:           {self.stats['total_requests']}")
        print(f"Errors encountered:           {self.stats['errors']}")
        print("=" * 60)


def main():
    parser = argparse.ArgumentParser(
        description="Backfill English translations using DeepL API"
    )
    parser.add_argument(
        "--api-key",
        help="DeepL API key (or set DEEPL_API_KEY env var)",
        default=os.getenv("DEEPL_API_KEY"),
    )
    parser.add_argument(
        "--db-url",
        help="Database URL (or set DATABASE_URL env var)",
        default=os.getenv("DATABASE_URL"),
    )
    parser.add_argument(
        "--pro", action="store_true", help="Use DeepL Pro API instead of free tier"
    )
    parser.add_argument(
        "--verbose", action="store_true", help="Print detailed debug information"
    )

    args = parser.parse_args()

    if not args.api_key:
        print("Error: DeepL API key required (--api-key or DEEPL_API_KEY env var)")
        sys.exit(1)

    if not args.db_url:
        print("Error: Database URL required (--db-url or DATABASE_URL env var)")
        sys.exit(1)

    # Validate API key format
    api_key = args.api_key.strip()
    if len(api_key) < 10:
        print("❌ Error: API key appears to be invalid (too short)")
        sys.exit(1)

    print("🚀 Starting translation backfill with DeepL")
    print(f"API tier: {'Pro' if args.pro else 'Free'}")
    print(f"API key: {api_key[:10]}..." if api_key else "No API key")

    translator = get_deepl_translator(args.api_key, is_pro=args.pro)
    backfiller = TranslationBackfiller(args.db_url, translator)

    try:
        backfiller.connect()

        backfiller.backfill_expert_symptoms()
        backfiller.backfill_expert_diseases()
        backfiller.backfill_health_risks()
        backfiller.backfill_emergency_guides()
        backfiller.backfill_emergency_guide_flows()

        backfiller.print_stats()
        print("\n✅ Translation backfill complete!")

    except Exception as e:
        print(f"\n❌ Error: {e}", file=sys.stderr)
        sys.exit(1)

    finally:
        backfiller.close()


if __name__ == "__main__":
    main()
