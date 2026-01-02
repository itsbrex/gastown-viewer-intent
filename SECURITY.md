# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.x     | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability, please report it by:

1. **Do NOT** open a public issue
2. Email the maintainers directly (see CODEOWNERS or repo description)
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

We will respond within 48 hours and work with you to address the issue.

## Security Considerations

This tool:
- Shells out to the `bd` CLI — ensure `bd` is from a trusted source
- Runs a local HTTP server — bind to localhost only by default
- Does not transmit data externally — all operations are local

For production deployments, consider:
- Running behind a reverse proxy with authentication
- Restricting network access to trusted clients
- Regular updates of dependencies
