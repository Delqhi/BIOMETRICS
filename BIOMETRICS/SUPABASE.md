# SUPABASE.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Supabase-Betrieb folgt globalen RLS-, Backup- und Audit-Anforderungen.
- Schema-/Policy-Änderungen erfordern verifizierte Impact- und Mapping-Updates.
- Datenhaltung und Retention bleiben compliance-fähig dokumentiert.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Leitfaden für Supabase als primäres Backend (DB/Auth/Storage/Edge).

## 1) Grundsatz
Supabase wird als zentrale Daten- und Auth-Schicht verwendet. Alle Strukturen sind projektagnostisch als Vorlage formuliert.

## 2) Domänenmodell (Template)

| Tabelle | Zweck | Primärschlüssel | Beziehungen | Sensitivität |
|---|---|---|---|---|
| users | Nutzerstammdaten | id | profiles, sessions | hoch |
| projects | Projektmetadaten | id | assets, tasks | mittel |
| assets | Content-Artefakte | id | projects | mittel |
| tasks | Aufgaben und Status | id | projects | niedrig |
| audit_logs | Nachvollziehbarkeit | id | users | hoch |

## 3) Spaltenstandard
Für jede Tabelle dokumentieren:
1. Spaltenname
2. Datentyp
3. Nullbarkeit
4. Default
5. Validierungsregel
6. Beispielwert

## 4) RLS-Standard
Für jede Tabelle definieren:
1. SELECT Policy
2. INSERT Policy
3. UPDATE Policy
4. DELETE Policy
5. Rollenbezug

## 5) Auth-Standard
- Rollenmodell: user, dev, admin, agent
- Session-Handling und Ablaufzeiten dokumentieren
- Recovery-Flows dokumentieren

## 6) Storage-Standard
- Bucket-Zwecke dokumentieren
- Zugriff je Rolle definieren
- Upload-/Download-Richtlinien festlegen

## 7) Edge Functions (Template)

| Function | Zweck | Trigger | Eingabe | Ausgabe | Sicherheitsgrenze |
|---|---|---|---|---|---|
| fn_task_sync | Task-Sync | API call | payload | status | role-check |
| fn_asset_validate | Asset-Prüfung | upload | asset_ref | score | policy-check |

## 8) Migrationsstrategie
1. Änderungen versioniert
2. Rollback-Pfad definiert
3. Datenintegritätscheck nach Migration

## 9) NLM-Datenbezug
- NLM-Assets in `assets` strukturieren
- Qualitätsscore pro Asset erfassen
- Verwendungsstatus pflegen (draft/review/approved/retired)

## 10) Backup/Restore (Template)
- Backup-Intervall: {BACKUP_INTERVAL}
- RPO/RTO Ziele: {RPO_RTO}
- Restore-Testintervall: {RESTORE_TEST_INTERVAL}

## 11) Verifikation
- Schema-Review
- RLS-Review
- Auth-Flow Review
- Migrationstest

## Abnahme-Check SUPABASE
1. Tabellenkatalog vollständig
2. RLS pro Tabelle definiert
3. Auth- und Storage-Strategie dokumentiert
4. Migrations-/Backupprozess vorhanden
5. NLM-Asset-Tracking berücksichtigt

---
