# TODO

## High Priority
- [ ] **Fix Jira text formatting in issue view** - Parse and render Jira wiki markup (`{color:...}`, `{code}`, `*bold*`, `_italic_`, etc.) properly instead of showing raw markup

## Distribution
- [ ] Set up Homebrew tap for easy installation (`brew install ichoosetoaccept/tap/jira-cli`)

## Community PRs to Evaluate
- [ ] #909 - JSON (Human Readable) Output (+1830 lines)
- [ ] #908 - Worklog CRUD commands (list/edit/delete) (+1198 lines)
- [ ] #877 - OAuth 3LO support (+2869 lines)
- [ ] #873 - Remote links in issue view (+122 lines)
- [ ] #869 - List to select transition status (+150 lines)
- [ ] #868 - Update glamour lib (+18 lines)
- [ ] #844 - Description in plain mode (+61 lines)
- [ ] #902 - Comment with group visibility (+43 lines)

## Already Integrated
- [x] Cookie-based authentication (ours)
- [x] `jira refresh` command (ours)
- [x] #905 - Global non-interactive mode
- [x] #886 - `--unformatted` flag for issue view
- [x] #887 - `jira api` command
- [x] #894 - `jira sprint create`
- [x] #904 - Sprint CSV/delimiter options
- [x] #899 - `--board` param for sprint list
