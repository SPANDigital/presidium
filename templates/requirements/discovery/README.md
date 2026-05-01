# Discovery

Drop input materials for agents into this folder. Anything an agent should read but a human will not directly publish lives here.

Typical contents:

- Confluence or Notion exports (PDF, HTML, or Markdown).
- Domain models — DBML files, ERD images, schema.yml / OpenAPI specs.
- Functional decomposition spreadsheets (CSV) and the diagrams generated from them.
- Stakeholder interview transcripts, recordings, or notes.
- Screenshots of the existing system, if there is one.

The folder is flat — no enforced subfolder structure — so agents can ingest whatever shape of evidence is available. Large binary files are typically tracked with Git LFS in the consuming product.

Agents read from this folder and write their structured outputs (YAML models, generated articles, session logs) into the product's own `agent-artifacts/` directory, which is not provided by this template.
