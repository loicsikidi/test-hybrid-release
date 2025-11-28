# TODO - Future Improvements

## Migration vers SLSA Level 3

### Contexte

Actuellement, le projet utilise `actions/attest-build-provenance@v2` pour générer des attestations de provenance SLSA. Cette approche fonctionne bien et est recommandée par GitHub, mais elle ne permet pas d'atteindre **SLSA Level 3** car :

1. Le workflow n'est pas un "trusted reusable workflow" isolé
2. `slsa-verifier` ne peut pas vérifier les attestations natives GitHub
3. Pas d'isolation complète du build process

### Objectif

Migrer vers **SLSA Level 3** en utilisant le framework officiel SLSA :
- **Builder** : `slsa-framework/slsa-github-generator`
- **Verifier** : `slsa-framework/slsa-verifier`

### Avantages de la migration

1. ✅ **SLSA Level 3 compliance** - Niveau de sécurité maximal pour les projets open source
2. ✅ **Vérification standardisée** - `slsa-verifier` peut vérifier les attestations
3. ✅ **Isolation du build** - Workflow réutilisable isolé avec permissions minimales
4. ✅ **Supply chain security renforcée** - Builder approuvé et audité par la communauté
5. ✅ **Interopérabilité** - Standard SLSA largement reconnu dans l'industrie

### Actions requises

#### 1. Migration du workflow Trust Bundle (YYYY-MM-DD)

**Fichier concerné** : `.github/workflows/release-bundle.yml`

**Changements** :
- [ ] Remplacer `actions/attest-build-provenance@v2` par `slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.0.0`
- [ ] Adapter le workflow pour utiliser le pattern "caller workflow" + "reusable workflow"
- [ ] Générer le bundle dans une étape séparée et l'uploader comme artifact
- [ ] Le reusable workflow SLSA génère l'attestation de manière isolée

**Exemple de structure** :
```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Generate bundle
        run: |
          go run ./cmd/main.go generate
          sha256sum tpm-ca-certificates.pem > checksums.txt

      - name: Upload bundle
        uses: actions/upload-artifact@v4
        with:
          name: bundle
          path: |
            tpm-ca-certificates.pem
            checksums.txt

  provenance:
    needs: [build]
    permissions:
      id-token: write
      contents: write
      actions: read
    uses: slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.0.0
    with:
      base64-subjects: "${{ needs.build.outputs.digests }}"
      upload-assets: true
```

**Références** :
- [SLSA Generic Generator](https://github.com/slsa-framework/slsa-github-generator/blob/main/internal/builders/generic/README.md)
- [Example workflow](https://github.com/slsa-framework/slsa-github-generator/blob/main/.github/workflows/e2e.generic.workflow_dispatch.main.default.yml)

#### 2. Migration du workflow Binary/OCI (vX.Y.Z)

**Fichier concerné** : `.github/workflows/release.yml`

**Options** :

##### Option A : SLSA Generic Generator (pour les binaires Go)
- [ ] Utiliser `generator_generic_slsa3.yml` comme pour le bundle
- [ ] GoReleaser génère les binaires
- [ ] SLSA generator atteste les binaires

##### Option B : SLSA Go Builder (spécifique Go)
- [ ] Utiliser `slsa-framework/slsa-github-generator/.github/workflows/builder_go_slsa3.yml@v2.0.0`
- [ ] Build natif par le SLSA builder (remplace GoReleaser pour la compilation)
- [ ] Meilleure intégration mais nécessite plus de refactoring

**Pour les images OCI** :
- [ ] Utiliser `slsa-framework/slsa-github-generator/.github/workflows/generator_container_slsa3.yml@v2.0.0`
- [ ] Attester les images ko directement

**Références** :
- [SLSA Go Builder](https://github.com/slsa-framework/slsa-github-generator/blob/main/internal/builders/go/README.md)
- [SLSA Container Generator](https://github.com/slsa-framework/slsa-github-generator/blob/main/internal/builders/container/README.md)

#### 3. Mise à jour de la documentation

**Fichiers à mettre à jour** :
- [ ] `README.md` - Section verification avec `slsa-verifier`
- [ ] Release notes template - Instructions de vérification mises à jour
- [ ] `.goreleaser.yaml` - Commentaires sur l'intégration SLSA

**Nouvelles instructions de vérification** :
```bash
# Avec slsa-verifier (maintenant supporté)
slsa-verifier verify-artifact tpm-ca-certificates.pem \
  --provenance-path tpm-ca-certificates.pem.intoto.jsonl \
  --source-uri github.com/loicsikidi/test-hybrid-release \
  --source-tag 2025-11-28
```

#### 4. Tests de validation

- [ ] Tester la génération d'attestation SLSA Level 3 sur un tag de test
- [ ] Vérifier avec `slsa-verifier` que la vérification fonctionne
- [ ] Valider que les attestations sont accessibles via GitHub API
- [ ] Comparer la taille des attestations (SLSA vs GitHub native)
- [ ] Mesurer l'impact sur le temps d'exécution du workflow

### Considérations et compromis

#### Avantages actuels de `actions/attest-build-provenance@v2` :
- ✅ **Simplicité** - Configuration minimale
- ✅ **Intégration native GitHub** - Pas de dépendance externe
- ✅ **Vérification `gh` CLI** - Outil officiel GitHub
- ✅ **Performance** - Génération rapide de l'attestation
- ✅ **Maintenance** - Géré par GitHub directement

#### Avantages de `slsa-framework/slsa-github-generator` :
- ✅ **SLSA Level 3** - Compliance maximale
- ✅ **Isolation** - Build process isolé
- ✅ **`slsa-verifier`** - Vérification standardisée
- ✅ **Auditabilité** - Workflow réutilisable audité
- ✅ **Communauté** - Standard industrie reconnu

#### Inconvénients potentiels :
- ⚠️ **Complexité accrue** - Workflow plus complexe avec jobs multiples
- ⚠️ **Temps d'exécution** - Potentiellement plus lent (isolation)
- ⚠️ **Dépendance externe** - Dépend du maintien du slsa-github-generator
- ⚠️ **Courbe d'apprentissage** - Plus difficile à comprendre et maintenir

### Décision

**Statut** : ⏸️ À évaluer

**Questions à se poser** :
1. Le projet nécessite-t-il vraiment SLSA Level 3 ? (Quels sont les besoins de conformité ?)
2. Les utilisateurs utilisent-ils `slsa-verifier` ? (GitHub CLI suffit-il ?)
3. La complexité additionnelle est-elle justifiée par les bénéfices ?

**Recommandation actuelle** :
- ✅ Garder l'approche actuelle pour l'instant (simple, efficace, maintenue par GitHub)
- 📋 Documenter cette TODO pour une évaluation future
- 🔍 Réévaluer si :
  - Des exigences de conformité SLSA Level 3 apparaissent
  - L'industrie standardise sur `slsa-verifier`
  - Le projet grandit et nécessite plus de garanties

### Ressources

- [SLSA Levels](https://slsa.dev/spec/v1.0/levels)
- [SLSA GitHub Generator](https://github.com/slsa-framework/slsa-github-generator)
- [SLSA Verifier](https://github.com/slsa-framework/slsa-verifier)
- [GitHub Attestations](https://docs.github.com/en/actions/security-guides/using-artifact-attestations-to-establish-provenance-for-builds)
- [Comparison: GitHub Native vs SLSA Generator](https://github.com/slsa-framework/slsa-github-generator/blob/main/SPECIFICATIONS.md#comparison-with-github-attestations)

### Notes

Migration évaluée le : 2025-11-28
Décision finale : En attente d'évaluation business/sécurité
