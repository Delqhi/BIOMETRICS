# WEBSHOP.md

## Global-Mandate-Alignment (AGENTS-GLOBAL)

- Webshop-Betrieb folgt globalen Payment-, Security- und Incident-Regeln.
- Integrität von Checkout-, Ledger- und Fraud-Kontrollen ist verpflichtend.
- Änderungen benötigen belastbare Nachweise und Mapping-Sync.

Status: ACTIVE  
Version: 1.0 (Universal)  
Stand: Februar 2026

## Zweck
Universeller Leitfaden für Webshop-Strukturen, Entscheidungslogik und sichere Abwicklung.

## Hinweis
Falls kein Shop benötigt wird, Status auf `NOT_APPLICABLE` setzen und begründen.

## Shop-Kernmodell (Template)

| Bereich | Ziel | Kritikalität |
|---|---|---|
| Katalog | Orientierung und Auswahl | P1 |
| Produktdetail | Vertrauen und Entscheidung | P1 |
| Warenkorb | Vorbereitung Checkout | P0 |
| Checkout | Abschluss | P0 |
| Order-Status | Transparenz | P1 |

## Produktmodell (Template)
- Produkt-ID
- Name
- Kategorie
- Preislogik
- Verfügbarkeitsstatus
- Steuerklasse

## Checkout-Prinzipien
1. klare Schritte
2. transparente Kosten
3. sichere Eingaben
4. nachvollziehbare Bestätigung

## Risiken
- Abbruch im Checkout
- unklare Preislogik
- fehlende Fehlertransparenz

## Checkout-Kompatibilität (Check)
1. Preis- und Gebührenlogik ist vor Abschluss transparent
2. Validierungsfehler sind verständlich und behebbar
3. Statuswechsel von Warenkorb zu Abschluss ist nachvollziehbar
4. Rechtliche Hinweise sind im kritischen Pfad erreichbar
5. Optional eingesetzte NLM-Assets erzeugen keine falschen Versprechen

## NLM-Einsatz
- Produkt-/Prozess-Erklärvideos optional
- Infografik für Preis-/Versandlogik möglich
- nur mit NLM-CLI und Qualitätsprüfung

## Verifikation
- End-to-End Bestellung im Testmodus
- Fehlerfallprüfung bei unvollständigen Eingaben
- Konsistenz mit rechtlichen Anforderungen

## Abnahme-Check WEBSHOP
1. Shop-Kernmodell dokumentiert
2. Checkout-Prinzipien vorhanden
3. Risikoabschnitt vorhanden
4. Verifikationslogik enthalten

---
