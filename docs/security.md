# Security

NetScope is designed for authorized security assessments only.

## Intended use

NetScope is intended for:

- Penetration testers assessing scoped targets
- Red teams evaluating their own infrastructure
- Blue teams auditing their environment
- System administrators and security engineers
- Incident responders

## Not intended for

NetScope does not include and should not be used for:

- Unauthorized access to systems you do not own
- Credential theft or harvesting
- Destructive activity
- Persistence
- Evading detection on systems you do not have permission to assess

## Secure development

The codebase adheres to the following practices:

- All discovery targets are validated before use
- No unsafe command execution
- Safe parsing of untrusted network responses
- No hardcoded credentials or secrets
- Secure defaults
- No unnecessary privilege escalation
- Timeouts and cancellation on all blocking operations

## Scope and authorization

The tool does not enforce scope boundaries. It is the operator's responsibility to ensure they have authorization to assess targets they feed into NetScope.

## Disclosure

Security issues in NetScope itself should be reported responsibly to the maintainers. Do not open public issues for security vulnerabilities.
