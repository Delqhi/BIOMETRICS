# SECURITY.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Security-by-Default und Secret-Hygiene sind nicht verhandelbar.
- Controls, Findings und Re-Tests folgen dem globalen Audit-Modell.
- Jede Ausnahme benötigt Compensating Control und Ablaufdatum.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Sicherheitsrahmen für Entwicklung, Betrieb und Content-Workflows.

## 1) Sicherheitsprinzipien
1. Least Privilege
2. Zero Trust Denkweise
3. Defense in Depth
4. Secure by Default
5. Nachweisbarkeit und Auditierbarkeit

## 2) Schutzobjekte
- Quellcode und Konfiguration
- Benutzer- und Betriebsdaten
- Integrationszugänge
- NLM-generierte Artefakte
- Betriebs- und Audit-Protokolle

## 3) Bedrohungsmodell (Template)

| Bedrohung | Angriffsvektor | Auswirkung | Gegenmaßnahme | Priorität |
|---|---|---|---|---|
| unautorisierter Zugriff | schwache Rechtevergabe | Datenabfluss | Rollenmodell + MFA | P0 |
| Secret Leak | Fehlkonfiguration | Account kompromittiert | Secret Hygiene + Rotation | P0 |
| Prompt Injection | ungeprüfte Inputs | fehlerhafte Artefakte | Input-Validierung + Review | P1 |
| Overclaim in Content | ungeprüfte NLM-Ausgabe | Compliance-Risiko | Qualitätsmatrix + Freigabe | P1 |

## 4) Zugriffskontrolle
- Rollen: user, dev, admin, agent
- Trennung von Leserechten und Schreibrechten
- Kritische Aktionen nur mit hoher Rolle

## 5) Secret-Management
1. Keine Secrets in Repo-Dateien
2. ENV-basiertes Secret Handling
3. Rotation bei Verdacht oder Incident
4. Zugriff nur für notwendige Rollen

## 6) API-Sicherheit
- Authentifizierung für produktive Endpunkte
- Autorisierung je Rolle
- Input-Validierung und Fehlerhärtung
- Rate Limits definieren

## 7) NLM-spezifische Sicherheit
1. NLM nur über NLM-CLI ausführen
2. Keine sensiblen Rohdaten ungefiltert an NLM
3. Output vor Übernahme auf Fakten und Compliance prüfen
4. Verworfenes dokumentieren

## 8) Compliance-Hinweise (Template)
- Rechtsraum: {COMPLIANCE_SCOPE}
- Datenklassen: {DATA_CLASSES}
- Aufbewahrung: {RETENTION_POLICY}
- Löschung: {DELETION_POLICY}

## 9) Incident-Response Quickflow
1. Severity klassifizieren (P0/P1/P2)
2. Schaden begrenzen
3. Secrets rotieren falls nötig
4. Ursache analysieren
5. Fix validieren
6. Postmortem dokumentieren

## 10) Security-Checks pro Zyklus
- P0-Risiken offen? nein
- Kritische Endpunkte geschützt? ja
- Secret-Policy eingehalten? ja
- NLM-Output geprüft? ja

## 11) Nachweise
- Security-Checklisten
- Audit-Logeinträge
- Freigabeprotokolle

## Abnahme-Check SECURITY
1. Threat Model vorhanden
2. Secret-Policy klar
3. Rollenmodell dokumentiert
4. NLM-Sicherheitsregeln enthalten
5. Incident-Prozess definiert

---
