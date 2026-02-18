# 📝 SETUP QUIZ - Test Your Knowledge!

**Topic:** Complete Setup Process  
**Questions:** 10  
**Passing Score:** 8/10

---

## Question 1
**Wo besorgst du den NVIDIA API Key?**

- [ ] Bei GitHub
- [x] Auf https://build.nvidia.com/
- [ ] In der opencode.json
- [ ] Bei Hugging Face

**Explanation:** Der offizielle NVIDIA API Key wird auf der NVIDIA Developer Platform (build.nvidia.com) erstellt.

---

## Question 2
**Wie beginnt ein NVIDIA API Key?**

- [ ] github-
- [x] nvapi-
- [ ] hugging-
- [ ] open-

**Explanation:** Alle NVIDIA API Keys beginnen mit dem Prefix "nvapi-".

---

## Question 3
**Wo wird der NVIDIA_API_KEY gespeichert?**

- [ ] In der opencode.json
- [x] In ~/.zshrc als Environment Variable
- [ ] Im BIOMETRICS Repo
- [ ] In der oh-my-opencode.json

**Explanation:** Der Key wird als Environment Variable in ~/.zshrc gespeichert, nicht in Config-Dateien.

---

## Question 4
**Welcher Befehl lädt die Shell-Konfiguration neu?**

- [ ] `reload ~/.zshrc`
- [ ] `bash restart`
- [x] `exec zsh`
- [ ] `source reload`

**Explanation:** `exec zsh` startet eine neue Shell-Session und lädt alle Konfigurationen neu.

---

## Question 5
**Wo befindet sich die opencode.json?**

- [ ] /etc/opencode/opencode.json
- [x] ~/.config/opencode/opencode.json
- [ ] ~/opencode.json
- [ ] /usr/local/etc/opencode.json

**Explanation:** Die korrekte Location ist ~/.config/opencode/opencode.json

---

## Question 6
**Welches Modell verwenden die BIOMETRICS Agenten?**

- [ ] gpt-4
- [ ] claude-3
- [x] nvidia-nim/qwen-3.5-397b
- [ ] gemini-pro

**Explanation:** Alle Haupt-Agenten verwenden NVIDIA Qwen 3.5 397B für maximale Performance.

---

## Question 7
**Wie groß ist das Context-Limit von Qwen 3.5 397B?**

- [ ] 32K Tokens
- [ ] 128K Tokens
- [x] 262K Tokens
- [ ] 1M Tokens

**Explanation:** Qwen 3.5 397B hat ein enormes Context-Limit von 262.144 Tokens.

---

## Question 8
**Welcher Befehl installiert OpenCode?**

- [ ] `brew install opencode`
- [x] `npm install -g opencode`
- [ ] `pip install opencode`
- [ ] `cargo install opencode`

**Explanation:** OpenCode wird als npm Package global installiert.

---

## Question 9
**Was macht `opencode auth add nvidia-nim`?**

- [ ] Installiert NVIDIA Treiber
- [x] Authentifiziert den NVIDIA Provider via OAuth
- [ ] Erstellt einen neuen API Key
- [ ] Testet die Verbindung

**Explanation:** Dieser Befehl startet den OAuth-Flow um den NVIDIA Provider zu authentifizieren.

---

## Question 10
**Wann wird die Shell neu geladen?**

- [ ] Nach jedem Befehl
- [x] Nach Änderungen an ~/.zshrc
- [ ] Nie, das ist automatisch
- [ ] Nur nach Neustart des Computers

**Explanation:** Nach Änderungen an ~/.zshrc muss die Shell mit `exec zsh` neu geladen werden.

---

## 📊 SCORING

- **10/10:** 🎉 Perfect! Du bist bereit für das Setup!
- **8-9/10:** ✅ Sehr gut! Lies nochmal die fehlenden Teile nach.
- **6-7/10:** ⚠️ Gut, aber lies COMPLETE-SETUP.md nochmal gründlich.
- **<6/10:** ❌ Bitte lies die komplette Setup-Anleitung bevor du startest!

---

**📚 Study Resources:**
- [COMPLETE-SETUP.md](../setup/COMPLETE-SETUP.md)
- [PROVIDER-SETUP.md](../setup/PROVIDER-SETUP.md)
- [UNIVERSAL-BLUEPRINT.md](../UNIVERSAL-BLUEPRINT.md)
