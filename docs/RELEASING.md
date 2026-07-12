# Publier etg

La release est automatisée par `.github/workflows/release.yaml`.

## Première publication

1. Publiez le dépôt sur GitHub et poussez la branche `main`.
2. Vérifiez que les tests de la CLI passent :

   ```bash
   cd packages/cli
   go test ./...
   ```

3. Créez puis poussez un tag :

   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

GitHub Actions lance GoReleaser, compile `etg` pour macOS et Linux en AMD64 et
ARM64, puis attache les archives et `checksums.txt` à la GitHub Release.

## Homebrew

La release initiale publie les binaires GitHub. Pour une formule Homebrew,
créez ensuite un dépôt de tap (par exemple `homebrew-tap`) et ajoutez la cible
`brews` correspondante dans `.goreleaser.yaml`, avec un token dédié ayant accès
en écriture au tap.
