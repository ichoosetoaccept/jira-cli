# JiraCLI (Enhanced Fork)

> **This is a fork of [ankitpokhrel/jira-cli](https://github.com/ankitpokhrel/jira-cli) with additional features for enterprise environments.**

For complete documentation, installation options, and usage examples, please see the [upstream README](https://github.com/ankitpokhrel/jira-cli#readme).

## What's Different in This Fork

This fork adds features specifically designed for enterprise Jira deployments, particularly those behind SSO, reverse proxies, or requiring certificate-based authentication.

### Cookie-Based Authentication

If your Jira is behind SSO, a reverse proxy (like pfportal), or requires client certificate authentication, this fork adds `cookie` auth support:

```sh
# During setup
jira init
# Select "Local" installation, then "cookie" as authentication type

# Sign in to Jira in your browser (authenticate via SSO/certificate as needed)
# Copy the JSESSIONID cookie value from your browser's DevTools
# Paste it when prompted - the cookie is validated and stored in your system keychain

# When your session expires, refresh without re-running full setup:
jira refresh
```

### Improved Jira Wiki Markup Rendering

Enhanced markdown conversion for `jira issue view` using the [j2m](https://github.com/FokkeZB/J2M) library:
- **Fixed table rendering** - properly normalizes `\r\n` to `\n`
- **Fixed escaped markup** - handles `{*}text{*}` → bold correctly
- **Color tag handling** - strips color tags while preserving content

### Integrated Community PRs

This fork includes several community contributions not yet merged upstream:
- **Worklog CRUD** - full create, read, update, delete for time tracking
- **Sprint create** - create new sprints from CLI
- **Remote links** - add web links to issues
- **API command** - make raw authenticated API requests
- And more...

## Installation

### Homebrew (Recommended)

```sh
brew tap ichoosetoaccept/tap
brew install jira-cli
```

### From Source

```sh
git clone https://github.com/ichoosetoaccept/jira-cli.git
cd jira-cli
make deps install
```

## Quick Start

```sh
# Initialize configuration
jira init

# View an issue with comments
jira issue view ISSUE-123 --comments 5

# Assign to yourself
jira issue assign ISSUE-123 $(jira me)

# Change status
jira issue move ISSUE-123 "In Progress"

# Log time
jira issue worklog add ISSUE-123 "2h 30m" --comment "Code review"
```

## Upstream Documentation

For complete documentation on all commands and features, see:
- [Upstream README](https://github.com/ankitpokhrel/jira-cli#readme)
- [Installation Guide](https://github.com/ankitpokhrel/jira-cli/wiki/Installation)
- [FAQs](https://github.com/ankitpokhrel/jira-cli/discussions/categories/faqs)

## Support

- **This fork**: [Issues](https://github.com/ichoosetoaccept/jira-cli/issues)
- **Upstream project**: Please [support the original author](https://opencollective.com/jira-cli#backers) if you find this tool useful

## License

MIT License - see [LICENSE](LICENSE)

---

<sub>Original project by [@ankitpokhrel](https://github.com/ankitpokhrel) • Fork maintained by [@ichoosetoaccept](https://github.com/ichoosetoaccept)</sub>
