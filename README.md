# Browser data parser

This library parses and converts data from Firefox- and Chromium-based
browsers.

[Documentation](https://pkg.go.dev/github.com/andrewarchi/browser)

## Browsers

Key:

- R: format can be read
- W: format can be written
- D: format documentation and/or source has been consulted

### Firefox

Firefox files currently parsed:

- `Profiles/{profile}/addons.json` (R)
- `Profiles/{profile}/bookmarkbackups/bookmarks-{date}_{count}_{hash}.{json|jsonlz4}` (R)
- `Profiles/{profile}/containers.json` (R)
- `Profiles/{profile}/extension-preferences.json` (R)
- `Profiles/{profile}/extension-settings.json` (R)
- `Profiles/{profile}/extensions.json` (R)
- `Profiles/{profile}/handlers.json` (R)
- `Profiles/{profile}/times.json` (R)
- `installs.ini` (R)
- `profiles.ini` (R)

#### Tor Browser

No Tor Browser-specific data is currently parsed.

### Chrome

Chrome files currently parsed:

- `{profile}/Bookmarks` (R)
- `First Run` (R)

Google Takeout files currently parsed:

- `Takeout/Chrome/Autofill.json` (R)
- `Takeout/Chrome/Bookmarks.html` (R)
- `Takeout/Chrome/BrowserHistory.json` (R)
- `Takeout/Chrome/Extensions.json` (R)
- `Takeout/Chrome/SearchEngines.json` (R)
- `Takeout/Chrome/SyncSettings.json` (R)

#### History Trends Unlimited (extension)

All History Trends Unlimited formats are parsed:

- `exported_analysis_history_{date}.{tsv|txt}` (RWD)
- `exported_archived_history_{date}.{tsv|txt}` (RWD)
- `history_autobackup_{date}_{full|incremental}.{tsv|txt|zip}` (RWD)

#### Brave

No Brave-specific data is currently parsed.

#### Edge (Chromium)

No Edge-specific data is currently parsed.

## Contributing

The project is designed to be strict and reject input that violates any
assumptions. For example, additional unknown json fields are disallowed,
which contrasts with the permissive default in `encoding/json`. The goal
is to prevent any data loss, so it fails early. Several types in the
structures are still unknown due to a small sample size of data and are
marked with `jsonutil.UnknownObj` or `jsonutil.UnknownType`. If you
encounter an error while parsing valid data, please
[report an issue](https://github.com/andrewarchi/browser/issues).

I am currently seeking information on the
`Takeout/Chrome/Dictionary.csv` file in Google Takeout.

## License

This project is made available under the
[Mozilla Public License](https://www.mozilla.org/en-US/MPL/2.0/).
