# GitLab Media Storage - Öffentliche URLs

## WICHTIG: Repository MUSS öffentlich sein!

Damite Videos und andere Medien ohne Login erreichbar sind, MUSS das GitLab-Repository auf **public** gestellt werden.

### Repository öffentlich machen:

1. Gehe zu: https://gitlab.com/zukunftsorientierte.energie/biometrics-videos
2. Klicke auf **Settings** → **General**
3. Unter **Visibility, project features, permissions**:
   - **Project visibility**: Public
   - Klicke auf **Save changes**

### Öffentliche URL konstruieren:

```
https://gitlab.com/{USERNAME}/{PROJECT_NAME}/-/uploads/{UPLOAD_HASH}/{FILENAME}
```

**Beispiel:**
```
https://gitlab.com/zukunftsorientierte.energie/biometrics-videos/-/uploads/d23d181f4278365b97454a3c0602d132/video.mp4
```

### Testen ob URL öffentlich ist:

```bash
curl -I "https://gitlab.com/zukunftsorientierte.energie/biometrics-videos/-/uploads/HASH/video.mp4"
```

**Erwartete Antwort:**
- `HTTP/2 200` = ✅ Öffentlich zugänglich
- `HTTP/2 404` = ❌ Repository ist privat oder URL falsch
- `HTTP/2 302` = ✅ Redirect (auch OK)

### In README.md verwenden:

```markdown
[![Video Thumbnail](video-thumbnail.png)](https://gitlab.com/USER/REPO/-/uploads/HASH/video.mp4)

**[Video abspielen ▶](https://gitlab.com/USER/REPO/-/uploads/HASH/video.mp4)**
```

### Häufige Fehler:

| Fehler | Ursache | Lösung |
|--------|---------|--------|
| 404 Not Found | Repository ist privat | Auf Public stellen |
| 404 Not Found | Falscher Upload-Pfad | `/uploads/` nicht `/raw/` verwenden |
| Login-Redirect | Repository nicht public | Settings → Visibility → Public |

### API Check (Optional):

```bash
curl "https://gitlab.com/api/v4/projects/USER%2FREPO" | jq .visibility
# Sollte "public" zurückgeben
```
