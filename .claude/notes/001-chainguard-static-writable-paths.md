# Writable Paths in `cgr.dev/chainguard/static:latest`

## Summary

Research conducted to identify writable filesystem paths in the Chainguard static base image.

## Methodology

Created a test Go binary that attempts to write to various common filesystem paths and ran it inside the `cgr.dev/chainguard/static:latest` container.

## Results

Only **one writable path** was found:

| Path | Status | Notes |
|------|--------|-------|
| `/tmp` | ✓ WRITABLE | Standard temporary directory |
| `/` | ✗ READ-ONLY | Root filesystem is read-only |
| `/var/tmp` | ✗ READ-ONLY | Directory does not exist |
| `/home` | ✗ READ-ONLY | Permission denied |
| `/root` | ✗ READ-ONLY | Permission denied |
| `/etc` | ✗ READ-ONLY | Permission denied |
| `/usr` | ✗ READ-ONLY | Permission denied |
| `/opt` | ✗ READ-ONLY | Permission denied |

## Conclusion

The `cgr.dev/chainguard/static:latest` image provides a highly restricted filesystem with **only `/tmp` writable**. This is consistent with the security-focused design of Chainguard images, which follow the principle of least privilege.

## Recommendations

For applications that need to write temporary data:
- Use `/tmp` as the output directory
- Or write to stdout/stderr for containerized environments
- Avoid assuming other paths are writable

## Date

2025-11-29
