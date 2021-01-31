# Browser data parser

This library parses and converts data from Firefox- and Chromium-based
browsers.

[Documentation](https://pkg.go.dev/github.com/andrewarchi/browser)

## Browsers

### Firefox

Firefox files currently parsed:

- `Profiles/{profile}/addons.json`
- `Profiles/{profile}/bookmarkbackups/bookmarks-{date}_{count}_{hash}.json`
- `Profiles/{profile}/bookmarkbackups/bookmarks-{date}_{count}_{hash}.jsonlz4`
- `Profiles/{profile}/containers.json`
- `Profiles/{profile}/extension-preferences.json`
- `Profiles/{profile}/extension-settings.json`
- `Profiles/{profile}/extensions.json`
- `Profiles/{profile}/handlers.json`
- `Profiles/{profile}/times.json`
- `installs.ini`
- `profiles.ini`

#### Tor Browser

No Tor Browser-specific data is currently parsed.

### Chrome

Chrome files currently parsed:

- `{profile}/Bookmarks`
- `First Run`

Google Takeout files currently parsed:

- `Takeout/Chrome/Autofill.json`
- `Takeout/Chrome/Bookmarks.html`
- `Takeout/Chrome/BrowserHistory.json`
- `Takeout/Chrome/Extensions.json`
- `Takeout/Chrome/SearchEngines.json`
- `Takeout/Chrome/SyncSettings.json`

#### History Trends Unlimited (extension)

History Trends Unlimited files currently parsed:

- `exported_analysis_history_{date}_{time}.tsv`
- `exported_analysis_history_{date}_{time}.zip`

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
