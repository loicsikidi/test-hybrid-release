## Objectif

Mettre en application sur un projet de test, les objectifs documentés dans `/home/lsikidi/workspace/repos/tpm-trust-bundle/.claude/notes/001-trust-model-and-cicd-pipeline.md`.

Je t'invite à te référer à ce document pour prendre connaissance du contexte et des objectifs.

## Milestones

Nous allons avoir deux types de release:

1. Publication seulement du bundle de confiance (Trust Bundle)
   - Pattern: le tag est de type `YYYY-MM-DD` (ex: `2024-06-15`)
   - Il faudra runner la commande `go run ./cmd/main.go generate` pour générer le document.
   - Cosign v3 sera utilisé pour signer `checksums.txt` généré avec une approche keyless (via `goreleaser`).
   - Dans l'attestation de provence générée par `Github` je souhaite avoir:
     - checksums.txt (généré par l'outil `goreleaser`)
       - note: dans ce cas précis, si ce n'est pas trop compliqué je souhaite que `goreleaser` génère le fichier `checksums.txt` sur la base du fichier `tpm-ca-certificates.pem` généré par la commande `go run ./cmd/main.go generate` mais en aucun cas je ne souhaite pas que `goreleaser` génère un binaire ou une image docker.
     - tpm-ca-certificates.pem (généré par `go run ./cmd/main.go generate`)
2. Publication du binaire et de l'image docker
   - Pattern: le tag est de type `vX.Y.Z` (ex: `v1.2.3`)
   - Je m'attends à ce que le binaire et l'image docker soient générés et publiés par`goreleaser`.
   - Cosign v3 sera utilisé pour signer `checksums.txt` généré avec une approche keyless (via `goreleaser`)
   - Dans l'attestation de provenance générée par `Github` je souhaite avoir:
     - checksums.txt (généré par l'outil `goreleaser`)
- le binaire s'appelle `awesomecli`

Resources:
- https://goreleaser.com/customization/attestations/
- https://goreleaser.com/blog/cosign-v3/

## Github

Le repository existe déjà github.com/loicsikidi/test-hybrid-release mais la CI est vide.

Je te laisse le soin de créer les fichiers nécessaires dans le repository pour atteindre les objectifs mentionnés ci-dessus. Tu peux commit et push directement dans le repository. Tu pourras monitorer l'état des pipelines via `gh`.

## Contraintes

- Utiliser `goreleaser` pour la gestion des releases (si possible, sinon je te laisse le soin de me proposer une alternative)
- Utiliser `cosign v3` en mode keyless pour la signature des artefacts (cf. golreleaser)