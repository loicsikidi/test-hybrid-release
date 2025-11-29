# Nix Build et VCS Information - Limitation

## Contexte

Après avoir corrigé le problème de GoReleaser (cf. [002-go-build-vcs-info.md](002-go-build-vcs-info.md)), j'ai testé si l'approche `debug.ReadBuildInfo()` fonctionne avec Nix en buildant depuis un tag git.

## Résultat

**Les informations VCS ne sont PAS disponibles dans les binaires buildés avec Nix**, même en utilisant `fetchgit` avec `leaveDotGit = true`.

### Vérification

```bash
$ nix-build -E 'with import <nixpkgs> {}; callPackage ./default.nix {}'
$ go version -m result/bin/cmd
```

Résultat:
```
mod	github.com/loicsikidi/test-hybrid-release	(devel)
build	-trimpath=true
# Pas de vcs.*, vcs.time, vcs.modified
```

Test du binaire:
```bash
$ ./result/bin/cmd version
Revision: unknown
Version: unknown
BuildTime: unknown
Dirty: false
```

## Pourquoi?

### 1. `buildGoModule` utilise `-trimpath`

Par défaut, `buildGoModule` dans nixpkgs utilise `-trimpath` pour rendre les builds reproductibles. Cela supprime les chemins absolus du binaire, mais **empêche également Go de détecter les informations VCS**.

### 2. Environnement de build isolé

Nix build dans un environnement sandboxé où:
- Le répertoire `.git` est nettoyé pendant le processus de build
- Les commandes `git` ne sont pas disponibles
- Le build se fait dans `/build/source`, pas dans un vrai dépôt git

### 3. `fetchgit` vs `fetchFromGitHub`

- **`fetchFromGitHub`**: Télécharge une tarball (pas de `.git`)
- **`fetchgit` avec `leaveDotGit = true`**: Préserve `.git` MAIS...
  - Nix nettoie quand même le répertoire avant le build
  - Le flag `-trimpath` empêche la détection VCS

## Solutions possibles

### Option 1: Injecter manuellement les VCS info via ldflags (COMPROMIS)

C'est exactement ce que nous voulions éviter, mais c'est la seule façon de faire fonctionner avec Nix:

```nix
buildGoModule {
  # ...

  ldflags = [
    "-s"
    "-w"
    "-X github.com/loicsikidi/test-hybrid-release/internal/version.version=${version}"
    "-X github.com/loicsikidi/test-hybrid-release/internal/version.revision=${src.rev}"
    "-X github.com/loicsikidi/test-hybrid-release/internal/version.time=1970-01-01T00:00:00Z"
  ];
}
```

**Problème**: Cela nécessite de modifier le code pour avoir des variables globales, ce qui va à l'encontre de l'approche `debug.ReadBuildInfo()`.

### Option 2: Désactiver `-trimpath` (CASSE LA REPRODUCTIBILITÉ)

```nix
buildGoModule {
  # ...

  # ATTENTION: Casse la reproductibilité des builds Nix!
  allowGoReference = true;
}
```

**Problème**:
- Va à l'encontre de la philosophie Nix
- Ne garantit pas que VCS info sera disponible car `.git` est quand même nettoyé

### Option 3: Accepter la limitation

Accepter que **l'approche `debug.ReadBuildInfo()` ne fonctionne pas avec Nix** et documenter cette limitation.

Les utilisateurs qui veulent builder avec Nix devront:
- Soit accepter `version unknown`
- Soit passer par GoReleaser pour la release officielle (qui fonctionne)
- Soit utiliser des ldflags manuels pour Nix

## Recommandation

**Option 3**: Documenter la limitation et l'accepter.

### Raisons:

1. **GoReleaser fonctionne correctement** après le fix
2. Les releases officielles (via CI/CD) utilisent GoReleaser, pas Nix
3. Nix est principalement utilisé pour:
   - L'environnement de développement (`shell.nix`)
   - Les builds locaux pour test
4. Forcer les VCS info avec Nix nécessiterait:
   - De dupliquer la logique de version
   - D'aller à l'encontre de l'approche `debug.ReadBuildInfo()`
   - De casser la reproductibilité

### Pour les développeurs

Les développeurs qui buildent localement avec Nix verront `version unknown`. C'est acceptable car:
- Ils savent qu'ils sont en dev
- Les releases officielles ont les bonnes versions
- Ils peuvent toujours utiliser `go build ./cmd` directement pour tester avec VCS info

## Références

- [Nix buildGoModule documentation](https://nixos.org/manual/nixpkgs/stable/#sec-language-go)
- [Go issue about -trimpath and VCS stamping](https://github.com/golang/go/issues/51831)
- L'utilisation de `-trimpath` est intentionnelle dans Nix pour la reproductibilité
