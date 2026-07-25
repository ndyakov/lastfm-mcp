# Security policy

## Supported versions

Until the first stable release, security fixes are made on the `master` branch. After releases begin, the latest minor release line will receive security fixes.

## Reporting a vulnerability

Use GitHub's private security advisory form for this repository. Do not open a public issue and do not include live credentials or personal listening data in a report.

Include the affected version, transport, deployment topology, reproduction steps, impact, and suggested mitigation if available. You should receive an acknowledgement within seven days.

## Deployment guidance

- Keep write and authentication tools disabled unless they are required.
- Prefer startup-provided session keys over sending passwords through an MCP tool.
- Keep Streamable HTTP bound to loopback for local use.
- For remote use, set a strong bearer token, terminate TLS at a trusted proxy, restrict network access, and use one server/session boundary per trust domain.
- Treat MCP clients and models as having access to all data returned by enabled tools.
