# Browser data parser

This library parses and converts data from Firefox- and Chromium-based
browsers.

## Firefox

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

## Chrome

- `{profile}/Bookmarks`
- `First Run`

### Google Takeout

- `Takeout/Chrome/Autofill.json`
- `Takeout/Chrome/Bookmarks.html`
- `Takeout/Chrome/BrowserHistory.json`
- `Takeout/Chrome/Dictionary.csv` (TODO)
- `Takeout/Chrome/Extensions.json`
- `Takeout/Chrome/SearchEngines.json`
- `Takeout/Chrome/SyncSettings.json`

## Brave

No Brave-specific data is currently parsed.
