# Go Build VCS Information Issue

## Problème

Le package `github.com/imjasonh/version` retourne `unknown` pour toutes les informations de version au lieu de récupérer automatiquement les métadonnées VCS depuis `debug.ReadBuildInfo()`.

### Symptôme

```bash
$ docker run --rm ghcr.io/loicsikidi/test-hybrid-release:v1.0.10 version
Revision: unknown
Version: unknown
BuildTime: unknown
Dirty: false
```

## Cause racine

Le binaire est construit **sans informations VCS** par GoReleaser. L'inspection du binaire avec `go version -m` montre:

```
path	command-line-arguments
dep	github.com/loicsikidi/test-hybrid-release	(devel)
```

Au lieu de:

```
path	github.com/loicsikidi/test-hybrid-release/cmd
mod	github.com/loicsikidi/test-hybrid-release	v1.0.10
build	vcs=git
build	vcs.revision=fefa9af9bb799ae55e5d655d56ebdd5c5864a79d
build	vcs.time=2025-11-29T07:51:29Z
build	vcs.modified=false
```

### Pourquoi?

GoReleaser est configuré pour builder le binaire en spécifiant un **fichier** au lieu d'un **package**:

```yaml
# dist/config.yaml:81
main: ./cmd/main.go  # ❌ INCORRECT
```

Quand Go build un fichier directement avec `go build ./cmd/main.go`, il:
1. Utilise `command-line-arguments` comme path au lieu du module path
2. Ne peut pas détecter le tag/version du module → `(devel)`
3. **N'embed PAS les informations VCS** dans le binaire

## Solution

Changer la configuration GoReleaser pour builder le **package** au lieu du fichier:

```yaml
# dist/config.yaml:81
main: ./cmd  # ✅ CORRECT
```

## Vérification

### Build correct (avec VCS info)

```bash
cd /tmp/test-clone
git checkout v1.0.10
go build -o /tmp/test-correct ./cmd  # ← package path
go version -m /tmp/test-correct
```

Résultat:
```
mod	github.com/loicsikidi/test-hybrid-release	v1.0.10
build	vcs=git
build	vcs.revision=fefa9af9bb799ae55e5d655d56ebdd5c5864a79d
build	vcs.time=2025-11-29T07:51:29Z
build	vcs.modified=false
```

Test du binaire:
```bash
$ /tmp/test-correct version
Revision: fefa9af9bb799ae55e5d655d56ebdd5c5864a79d
Version: v1.0.10
BuildTime: 2025-11-29T07:51:29Z
Dirty: false
```

### Build incorrect (sans VCS info)

```bash
go build -o /tmp/test-wrong ./cmd/main.go  # ← file path
go version -m /tmp/test-wrong
```

Résultat:
```
path	command-line-arguments
dep	github.com/loicsikidi/test-hybrid-release	(devel)
# Pas d'infos vcs.*
```

## Références

- [github.com/imjasonh/version README](https://github.com/imjasonh/version/blob/main/README.md)
- [Go debug.ReadBuildInfo documentation](https://pkg.go.dev/runtime/debug#ReadBuildInfo)
- Go 1.12+ embed automatiquement les VCS info **si le build est fait depuis un package, pas un fichier**

## Note importante

Les `ldflags` de GoReleaser (lignes 86-90) qui essaient d'injecter version/commit/date dans `main.version`, `main.commit`, `main.date` sont **inutiles** quand on utilise l'approche `debug.ReadBuildInfo()`. Ces variables n'existent même pas dans le code actuel.

Avec la fix du `main: ./cmd`, on peut **supprimer complètement** ces ldflags car Go embed automatiquement toutes ces informations.
