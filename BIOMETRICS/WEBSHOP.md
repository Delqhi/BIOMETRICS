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

## KI-gestützte Produktanalyse mit Qwen

### qwen_vision_analysis für Produktbilder
Automatische Qualitätsprüfung und Kategorisierung von Produktbildern.

| Analyse | Beschreibung | Aktion |
|---------|--------------|--------|
| Bildauflösung | Min. 1200x1200px prüfen | Auto-Reject bei Failure |
| Hintergrund-Qualität | Reinweiß/propur prüfen | Optimierungs-Vorschlag |
| Farb-Konsistenz | Markenfarben erkennen | Branding-Validierung |
| Text-Lesbarkeit | Overlay-Text prüfen | Accessibility-Check |

### qwen_conversation für Produktberatung
Intelligenter Produktberater im Chat.

| Feature | Beschreibung | Integration |
|---------|--------------|-------------|
| Produktvergleich | Ähnliche Produkte vorschlagen | Retrieval-Augmented |
| Größenberatung | Passform-Empfehlungen | Größentabelle + KI |
| Alternativ-Vorschläge | Bei Nichtverfügbarkeit | Cross-Sell |
| Bewertungs-Zusammenfassung | Sentiment-Analyse | Review-Aggregation |

### API-Integration
```typescript
// Produktbild-Analyse
const productAnalysis = await fetch('/api/qwen/vision', {
  method: 'POST',
  body: JSON.stringify({
    image: productImageBase64,
    analysisType: 'product_listing',
    requirements: {
      minResolution: '1200x1200',
      background: 'white',
      maxProducts: 1
    }
  })
});

// Produktberatung
const recommendation = await fetch('/api/qwen/chat', {
  method: 'POST',
  body: JSON.stringify({
    messages: [{ role: 'user', content: query }],
    skill: 'qwen_conversation',
    context: {
      productId: currentProduct,
      userPreferences: userPrefs,
      chatHistory: conversation
    }
  })
});
```

### Qualitäts-Workflow
1. Upload → qwen_vision_analysis
2. Score < 80% → Optimierungs-Prompt
3. Freigabe → Listing aktivieren
4. Monitoring → Bewertungs-Analyse

---
