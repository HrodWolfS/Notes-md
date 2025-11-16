# 🚀 Guide de démarrage rapide

## Installation en 30 secondes

### Option 1 : Script d'installation automatique (recommandé)

```bash
curl -sSL https://raw.githubusercontent.com/hrodwolf/notesmd/main/install.sh | bash
```

### Option 2 : Go Install

```bash
go install github.com/hrodwolf/notesmd/cmd/notesmd@latest
```

### Option 3 : Makefile (pour développeurs)

```bash
git clone https://github.com/hrodwolf/notesmd.git
cd notesmd
make install-user  # Installe dans ~/bin
```

## Vérifier l'installation

```bash
nmd
```

## Premiers pas

### 1. Ouvrir un dossier de notes

```bash
nmd ~/Documents/notes
```

### 2. Créer votre première note

- Appuyez sur `n`
- Tapez le nom : `ma-premiere-note`
- Ajoutez du contenu en Markdown
- Appuyez sur `Ctrl+S` pour sauvegarder

### 3. Naviguer

- `↑` `↓` : Se déplacer dans la liste
- `Enter` : Ouvrir un fichier ou entrer dans un dossier
- `←` : Retour au dossier parent

### 4. Rechercher

- `/` : Recherche fuzzy dans les noms de fichiers
- `F` : Recherche dans la note ouverte (avec highlight ⚡)

### 5. Organiser

- `b` : Ajouter/retirer des bookmarks
- `B` : Voir tous les bookmarks
- `Ctrl+R` : Fichiers récents

## Raccourcis essentiels

| Touche | Action |
|--------|--------|
| `?` | Aide complète |
| `n` | Nouvelle note |
| `e` | Éditer dans $EDITOR |
| `/` | Rechercher |
| `q` | Quitter |

## Configuration

Votre configuration est dans `~/.config/notesmd/config.json`

```json
{
  "editor": "nvim",
  "default_dir": "~/Documents/notes"
}
```

## Besoin d'aide ?

- Appuyez sur `?` dans l'application
- Consultez le [README complet](README.md)
- [Ouvrir une issue](https://github.com/hrodwolf/notesmd/issues)

---

**Astuce** : Ajoutez un alias dans votre shell pour lancer rapidement vos notes :

```bash
# Dans ~/.bashrc ou ~/.zshrc
alias notes='nmd ~/Documents/notes'
```

Maintenant lancez simplement `notes` ! 🎉
