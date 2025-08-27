## OpenTelemetry Lambda Layer Release Procedure (All Languages)

Releases are automated via GitHub Actions and are triggered by pushing a tag with a specific prefix. When a tag is pushed:
- A draft GitHub Release is created automatically with the same tag name
- The combined layer is built for amd64 and arm64 (where applicable)
- Artifacts are attached to the draft Release
- Layers are published publicly to multiple AWS regions

This guide applies to Go, Python, NodeJS, Java, Ruby combined layers, and the Collector layer.

### Tag prefixes and formats

Use the following tag formats to trigger releases. The version should include a leading "v" and only digits and dots. The workflows derive the layer version by stripping everything up to the last slash and removing any non-numeric prefix (e.g., "v").

- Go combined layer: `combined-layer-go/vX.Y.Z`
- Python combined layer: `combined-layer-python/vX.Y.Z`
- NodeJS combined layer: `combined-layer-nodejs/vX.Y.Z`
- Java combined layer: `combined-layer-java/vX.Y.Z`
- Ruby combined layer: `combined-layer-ruby/vX.Y.Z`
- Collector layer: `layer-collector/vX.Y.Z`

Examples:

```bash
# Go
git tag combined-layer-go/v1.2.3
git push origin combined-layer-go/v1.2.3

# Python
git tag combined-layer-python/v1.2.3
git push origin combined-layer-python/v1.2.3

# NodeJS
git tag combined-layer-nodejs/v1.2.3
git push origin combined-layer-nodejs/v1.2.3

# Java
git tag combined-layer-java/v1.2.3
git push origin combined-layer-java/v1.2.3

# Ruby
git tag combined-layer-ruby/v1.2.3
git push origin combined-layer-ruby/v1.2.3

# Collector
git tag layer-collector/v0.75.0
git push origin layer-collector/v0.75.0
```

### What the workflows do

After the tag push:
- A draft GitHub Release is created automatically
- The layer is built per architecture and uploaded as an artifact
- The artifact is attached to the draft Release
- The layer is published publicly across a matrix of AWS regions and compatible runtimes
- For the Collector, the workflow also appends region-agnostic ARN templates to the Release body

Related workflows (for reference):
- `.github/workflows/release-combined-go-lambda-layer.yml`
- `.github/workflows/release-combined-layer-python.yml`
- `.github/workflows/release-combined-layer-nodejs.yml`
- `.github/workflows/release-combined-layer-java.yml`
- `.github/workflows/release-combined-ruby-lambda-layer.yml`
- `.github/workflows/release-layer-collector.yml`
- `.github/workflows/layer-publish.yml` (reusable publisher)

### Releasing step-by-step

1. Decide the next version `vX.Y.Z` for the layer you want to release.
2. Create and push the appropriate tag (see examples above).
3. Monitor the corresponding GitHub Actions workflow until it completes.
4. Review the draft Release that was created automatically.
   - For combined language layers, you can find published ARNs in the workflow logs (each publish step prints the ARN).
   - For the Collector, ARN templates are appended to the Release body automatically.
5. Edit the draft Release notes if needed (changelog, highlights, ARNs) and publish the Release.

### Notes and tips

- The publisher converts the version dots to underscores in the layer name suffix (e.g., `1.2.3` -> `1_2_3`).
- Supported runtimes and AWS regions are controlled by each workflow. Adjust there if needed.
- Releases use OIDC to assume the publishing role. Ensure the required secrets/roles exist in the repo settings.
- If something goes wrong, you can delete the tag and the draft Release and try again.

