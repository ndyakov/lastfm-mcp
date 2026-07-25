# Improvement ideas

The initial implementation deliberately keeps state and deployment assumptions small. Useful follow-up work includes:

1. Add OAuth 2.1 at the HTTP boundary for hosted multi-user deployments. A static bearer token is suitable for a private service, not a shared public service.
2. Store a distinct Last.fm session per authenticated MCP principal. The initial server has one process-wide Last.fm client and session.
3. Add bounded caching for immutable metadata and short-lived chart/search results, with Last.fm-aware rate limiting, retry budgets, and observability.
4. Add MCP resources and prompts for common workflows such as listening summaries, discovery, and weekly comparisons. Keep the underlying endpoint tools available for composability.
5. Add result-shaping options so clients can request compact summaries instead of full Last.fm payloads and reduce model context usage.
6. Add OpenTelemetry traces and metrics for MCP calls and outbound Last.fm requests without recording credentials or sensitive tool arguments.
7. Add integration tests against a disposable protocol client and an opt-in live Last.fm test account. Unit tests should remain network-independent.
8. Publish an OCI image with provenance, SBOMs, signed checksums, and automated vulnerability scanning.
9. Revisit the six experimental methods periodically and remove or promote them based on observed Last.fm behavior.
10. Add a session-bootstrap CLI command so users never need to expose credentials through MCP authentication tools.
