# Publishing to NPM

## Prerequisites

1. NPM account (create at https://www.npmjs.com/signup)
2. Verified email address
3. Two-factor authentication enabled (recommended)

## First Time Setup

### 1. Login to NPM

```bash
npm login
```

Enter your:
- Username
- Password
- Email
- OTP (if 2FA enabled)

### 2. Verify Login

```bash
npm whoami
```

Should display your npm username.

### 3. Check Package Name Availability

```bash
npm search @kainos/financial-agent-cli
```

If the name is taken, update `package.json` with a different name.

## Publishing Steps

### 1. Update Version

Edit `package.json` and bump the version:

```json
{
  "version": "1.0.0"  // Change to 1.0.1, 1.1.0, 2.0.0, etc.
}
```

Or use npm version command:

```bash
npm version patch  # 1.0.0 -> 1.0.1
npm version minor  # 1.0.0 -> 1.1.0
npm version major  # 1.0.0 -> 2.0.0
```

### 2. Build the Package

```bash
pnpm build
```

This compiles TypeScript to JavaScript in the `dist/` folder.

### 3. Test Locally

Test the built package:

```bash
node dist/index.js --help
node dist/index.js price AAPL
```

### 4. Dry Run

See what will be published:

```bash
npm publish --dry-run
```

Review the output to ensure only necessary files are included.

### 5. Publish

For scoped packages (like @kainos/...):

```bash
npm publish --access public
```

For unscoped packages:

```bash
npm publish
```

### 6. Verify Publication

Check on npm:

```bash
npm view @kainos/financial-agent-cli
```

Or visit: https://www.npmjs.com/package/@kainos/financial-agent-cli

## Post-Publishing

### 1. Test Installation

In a new directory:

```bash
npm install -g @kainos/financial-agent-cli
financial-agent --help
```

### 2. Tag Release (Optional)

If using Git:

```bash
git tag v1.0.0
git push origin v1.0.0
```

### 3. Update Documentation

Update README with:
- Installation instructions
- New features
- Breaking changes

## Publishing Updates

### Patch Release (Bug Fixes)

```bash
npm version patch
pnpm build
npm publish --access public
```

### Minor Release (New Features)

```bash
npm version minor
pnpm build
npm publish --access public
```

### Major Release (Breaking Changes)

```bash
npm version major
pnpm build
npm publish --access public
```

## Unpublishing

⚠️ **Warning**: Unpublishing is permanent and can break dependent packages.

Within 72 hours of publishing:

```bash
npm unpublish @kainos/financial-agent-cli@1.0.0
```

After 72 hours, you can only deprecate:

```bash
npm deprecate @kainos/financial-agent-cli@1.0.0 "This version has been deprecated"
```

## Troubleshooting

### Error: Package name already exists

**Solution**: Change the package name in `package.json` or use a scope:
```json
{
  "name": "@your-username/financial-agent-cli"
}
```

### Error: You must verify your email

**Solution**: 
1. Check your email for verification link
2. Or run: `npm profile get`
3. Resend verification: Visit npmjs.com profile settings

### Error: 403 Forbidden

**Solution**:
1. Check you're logged in: `npm whoami`
2. For scoped packages, use: `npm publish --access public`
3. Verify package name isn't taken

### Error: No permission to publish

**Solution**:
1. Login with correct account: `npm login`
2. Check organization membership (for @org/package)
3. Verify 2FA code if enabled

## Best Practices

1. **Semantic Versioning**: Follow semver (major.minor.patch)
2. **Changelog**: Maintain CHANGELOG.md with version history
3. **Testing**: Always test before publishing
4. **Documentation**: Keep README up to date
5. **Git Tags**: Tag releases in Git
6. **CI/CD**: Automate publishing with GitHub Actions

## Automation with GitHub Actions

Create `.github/workflows/publish.yml`:

```yaml
name: Publish to NPM

on:
  release:
    types: [created]

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '20'
          registry-url: 'https://registry.npmjs.org'
      - run: pnpm install
      - run: pnpm build
      - run: npm publish --access public
        env:
          NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}
```

Add `NPM_TOKEN` to GitHub repository secrets.

## Package Statistics

After publishing, track your package:

- Downloads: https://npm-stat.com/charts.html?package=@kainos/financial-agent-cli
- Bundle size: https://bundlephobia.com/package/@kainos/financial-agent-cli
- Package health: https://snyk.io/advisor/npm-package/@kainos/financial-agent-cli

## Support

- NPM Documentation: https://docs.npmjs.com/
- NPM Support: https://www.npmjs.com/support
- Package Publishing Guide: https://docs.npmjs.com/packages-and-modules/contributing-packages-to-the-registry
