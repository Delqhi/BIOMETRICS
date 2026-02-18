# MAPPING-FRONTEND-BACKEND.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Frontend↔Backend-Verträge sind als verbindliche Kontrollpunkte zu pflegen.
- Jede Änderung erfordert zeitnahes Mapping-Update mit Ownership.
- Inkonsistenzen werden als Delivery-Risiko behandelt.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Abgleich zwischen Frontend-Flows und Backend-Funktionen.

## Mapping-Tabelle

| Frontend-Flow | Command | Endpoint | Backend-Service | Status |
|---|---|---|---|---|
| User Registration | CMD.AUTH.REGISTER | POST /api/auth/register | Supabase Auth + Edge Function | active |
| User Login | CMD.AUTH.LOGIN | POST /api/auth/login | Supabase Auth | active |
| View Products | CMD.PRODUCTS.LIST | GET /api/products | Supabase DB | active |
| Create Product | CMD.PRODUCTS.CREATE | POST /api/products | Supabase DB + RLS | active |
| Update Product | CMD.PRODUCTS.UPDATE | PATCH /api/products/:id | Supabase DB + RLS | active |
| Delete Product | CMD.PRODUCTS.DELETE | DELETE /api/products/:id | Supabase DB + RLS | active |
| Create Order | CMD.ORDERS.CREATE | POST /api/orders | Supabase DB + Edge Function | active |
| List Orders | CMD.ORDERS.LIST | GET /api/orders | Supabase DB | active |
| View Order Detail | CMD.ORDERS.GET | GET /api/orders/:id | Supabase DB | active |
| Upload Asset | CMD.ASSETS.UPLOAD | POST /api/assets | Supabase Storage | active |
| List Assets | CMD.ASSETS.LIST | GET /api/assets | Supabase Storage | active |
| Qwen Vision Analysis | CMD.QWEN.VISION | POST /api/qwen/vision | Vercel Edge Function | active |
| Qwen Code Generation | CMD.QWEN.CODE | POST /api/qwen/chat | Vercel Edge Function | active |
| Qwen OCR | CMD.QWEN.OCR | POST /api/qwen/ocr | Vercel Edge Function | active |
| NLM Video Generation | CMD.NLM.GENERATE.VIDEO | POST /api/nlm/video | NLM Integration | active |
| NLM Infographic | CMD.NLM.GENERATE.INFOGRAPHIC | POST /api/nlm/infographic | NLM Integration | active |

## Prüfregeln
1. Kein Flow ohne API-Anbindung
2. Kein Endpoint ohne Frontend-Nutzen in Kernpfaden
3. Fehlerfälle müssen UI-seitig behandelbar sein

## Abnahme-Check
1. Kernflows vollständig gemappt
2. Service-Zuordnung vorhanden
3. Fehlerpfad pro Kernflow dokumentiert

---
