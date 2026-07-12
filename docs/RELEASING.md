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

Le tap public [mattrbdr/homebrew-tap](https://github.com/mattrbdr/homebrew-tap)
est configuré dans `.goreleaser.yaml`. Avant de publier la première release,
créez un fine-grained personal access token GitHub autorisé à écrire le contenu
du dépôt `mattrbdr/homebrew-tap`, puis enregistrez-le comme secret Actions
`HOMEBREW_TAP_GITHUB_TOKEN` dans `mattrbdr/envtogithub`.

Après une release stable, les utilisateurs pourront installer la CLI avec :

```bash
brew tap mattrbdr/tap
brew install etg
```
