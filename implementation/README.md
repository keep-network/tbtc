
### Setup

```
npm run setup
```

### Compilation

```
npm run compile
```

## Lint

Linting is currently only enabled for JS test code.

```
# Show issues
npm run js:lint

# Automatically fix issues
npm run js:lint:fix
```

Eslint errors can be disabled using a comment on the previous line. For example, to disable linter errors for the 'no-unused-vars' rule: `// eslint-disable-next-line no-unused-vars`.