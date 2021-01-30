# Browser data parser

This library parses and converts data from Firefox- and Chromium-based
browsers.

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

#### Brave

No Brave-specific data is currently parsed.

#### Edge (Chromium)

No Edge-specific data is currently parsed.

## License

This project is made available under the
[Mozilla Public License](https://www.mozilla.org/en-US/MPL/2.0/).
