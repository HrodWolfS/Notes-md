# Guide de Release

Ce document explique comment créer une nouvelle release de NotesMD.

## Prérequis

- Accès push au dépôt GitHub
- Git configuré localement
- Code mergé sur `main` et testé

## Processus de Release

### 1. Préparer la release

```bash
# Assurez-vous d'être sur main et à jour
git checkout main
git pull origin main

# Vérifier que tout compile
make build

# Vérifier que les tests passent
go test ./...
```

### 2. Créer le tag de version

Suivez le [Semantic Versioning](https://semver.org/) :
- **Major** (v2.0.0) : Changements incompatibles
- **Minor** (v1.1.0) : Nouvelles fonctionnalités compatibles
- **Patch** (v1.0.1) : Corrections de bugs

```bash
# Exemple pour v1.0.0
VERSION="v1.0.0"

# Créer le tag avec un message
git tag -a $VERSION -m "Release $VERSION

## What's New
- Feature 1
- Feature 2
- Bug fix 1

## Installation
\`\`\`bash
go install github.com/hrodwolf/notesmd/cmd/notesmd@$VERSION
\`\`\`
"

# Pousser le tag vers GitHub
git push origin $VERSION
```

### 3. GitHub Actions s'occupe du reste

Une fois le tag poussé :

1. **Workflow CI** vérifie que tout compile
2. **Workflow Release** se lance automatiquement :
   - Compile pour toutes les plateformes (Linux, macOS, Windows, FreeBSD)
   - Crée les archives avec LICENSE et README
   - Génère les checksums
   - Crée une release GitHub avec tous les binaires
   - Génère le changelog automatiquement

### 4. Vérifier la release

1. Allez sur https://github.com/hrodwolf/notesmd/releases
2. Vérifiez que la release est créée avec tous les binaires
3. Testez le téléchargement d'un binaire
4. Vérifiez le changelog

### 5. Annoncer la release (optionnel)

- Mettre à jour le README si nécessaire
- Poster sur les réseaux sociaux
- Mettre à jour la documentation

## Binaires générés automatiquement

Pour chaque release, les binaires suivants sont créés :

- `nmd_[version]_Linux_x86_64.tar.gz`
- `nmd_[version]_Linux_arm64.tar.gz`
- `nmd_[version]_Linux_armv7.tar.gz`
- `nmd_[version]_Darwin_x86_64.tar.gz` (macOS Intel)
- `nmd_[version]_Darwin_arm64.tar.gz` (macOS Apple Silicon)
- `nmd_[version]_Windows_x86_64.zip`
- `nmd_[version]_FreeBSD_x86_64.tar.gz`
- `checksums.txt` - SHA256 de tous les fichiers

## Exemple complet

```bash
# 1. Préparer
git checkout main
git pull origin main
make build

# 2. Créer le tag
git tag -a v1.0.0 -m "Release v1.0.0"

# 3. Pousser
git push origin v1.0.0

# 4. Attendre ~2-3 minutes que GitHub Actions finisse

# 5. Vérifier
open https://github.com/hrodwolf/notesmd/releases
```

## En cas de problème

### Annuler une release

```bash
# Supprimer le tag localement
git tag -d v1.0.0

# Supprimer le tag sur GitHub
git push origin :refs/tags/v1.0.0

# Supprimer la release sur GitHub (via l'interface web)
```

### Re-créer une release

```bash
# Corriger le problème dans le code
git commit -am "fix: problème X"
git push origin main

# Créer un nouveau tag patch
git tag -a v1.0.1 -m "Release v1.0.1 - Fix X"
git push origin v1.0.1
```

## Notes

- Les tags doivent commencer par `v` (ex: `v1.0.0`)
- GoReleaser génère automatiquement le changelog depuis les commits
- Utilisez des messages de commit conventionnels pour un meilleur changelog :
  - `feat:` pour les nouvelles fonctionnalités
  - `fix:` pour les corrections de bugs
  - `docs:` pour la documentation
  - `chore:` pour les tâches de maintenance

## Ressources

- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [GoReleaser Documentation](https://goreleaser.com/)
