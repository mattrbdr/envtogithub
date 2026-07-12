# etg — env to GitHub

`etg` synchronise les secrets GitHub d'un environnement depuis un fichier versionnable
sans valeurs sensibles : `env.<environment>.to.github`.

## CLI

Copiez l'exemple et renseignez vos valeurs localement :

```bash
cp env.production.to.github.example env.production.to.github
```

Les fichiers `env.*.to.github` sont ignorés par Git : ne versionnez jamais les
valeurs réelles.

Puis lancez :

```bash
etg --dry-run
etg
etg --dir ./config --repo owner/repository
etg --color=always
etg info
etg list production
```

Chaque clé est créée ou mise à jour avec `gh secret set` dans l'environnement
GitHub déduit du nom du fichier. `gh auth login` doit avoir été exécuté avant une
synchronisation réelle. Le mode `--dry-run` ne contacte pas GitHub et n'affiche
jamais les valeurs.

La couleur est activée automatiquement dans un terminal. Utilisez `--color=never`
pour la désactiver ou `--color=always` pour la forcer ; la variable `NO_COLOR`
la désactive également en mode automatique.

### Inspection distante

```bash
# Environnements GitHub disponibles dans le dépôt courant
etg info

# Noms des secrets d'un environnement (sans jamais afficher les valeurs)
etg list production
```

Ajoutez `--repo owner/repository` à ces commandes pour cibler un autre dépôt.

### Développement de la CLI

La CLI est écrite en Go pour une distribution Homebrew simple à terme.

```bash
cd packages/cli
go test ./...
go run ./cmd/etg --dry-run
go build -o etg ./cmd/etg
```

### Releases

Les tags `v*` déclenchent GoReleaser et publient les binaires macOS/Linux dans
une GitHub Release. Voir [la procédure de release](docs/RELEASING.md).

Après la première release, installez etg avec Homebrew :

```bash
brew tap mattrbdr/tap
brew install etg
```

## Site

Le site Astro de présentation vit dans `apps/web`.

```bash
npm install
npm run dev
npm run build
```

Le déploiement production est défini dans
[.github/workflows/deploy-on-production.yaml](.github/workflows/deploy-on-production.yaml).
Renseignez les paramètres non sensibles dans
`.github/deploy/o2switch-production.env`, puis les secrets dans l'environnement
GitHub `production`.
