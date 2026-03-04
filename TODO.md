# qbittorrent library TODO

## Bugs

- [x] `AutoTMM bool` → `*bool` in `TorrentsAddParams` — zero-value `false` is indistinguishable from "not set", so `autoTMM=false` is never sent. Follow existing `RootFolder *bool` pattern.
- [x] `Paused bool` → `*bool` in `TorrentsAddParams` — same issue; can't override an instance that defaults to paused-on-add.
- [x] `TorrentsDelete` hardcodes `deleteFiles=true` (line 482) — no way to remove a torrent while keeping data. Add a `deleteFiles bool` parameter.
- [x] `TorrentsProperties.IsPrivate` JSON tag is `"isPrivate"` but API returns `"is_private"` — field silently never deserializes.
- [x] Tag splitting in `UnmarshalJSON` uses `,` but qBittorrent returns `, ` (comma-space) — parsed tags get leading spaces. Use `strings.Split` + `strings.TrimSpace`, or split on `, `.
- [x] `TorrentsDownload` hits `/api/v2/torrents/file` which doesn't exist in the official API. Either remove or redirect to `/api/v2/torrents/export` (already handled by `TorrentsExport`).
- [x] `TorrentsAdd` wrapper (line 468) has a dead `torrentFile` string parameter that is accepted but never used.

## Missing methods (critical)

All trivial — single POST with `hashes` form value:

- [x] `Pause(hashes)` — `POST /api/v2/torrents/pause`
- [x] `Resume(hashes)` — `POST /api/v2/torrents/resume`
- [x] `SetLocation(hashes, location)` — `POST /api/v2/torrents/setLocation`
- [x] `Recheck(hashes)` — `POST /api/v2/torrents/recheck`
- [x] `Reannounce(hashes)` — `POST /api/v2/torrents/reannounce`
- [x] `SetCategory(hashes, category)` — `POST /api/v2/torrents/setCategory`
- [x] `SetAutoTMM(hashes, enable)` — `POST /api/v2/torrents/setAutoTMM`
- [x] `Rename(hash, name)` — `POST /api/v2/torrents/rename`

## Missing methods (moderate)

- [x] `AuthLogout` — `POST /api/v2/auth/logout`
- [x] `AppVersion` — `GET /api/v2/app/version`
- [x] `AppPreferences` / `SetAppPreferences` — `GET/POST /api/v2/app/preferences`
- [x] `DefaultSavePath` — `GET /api/v2/app/defaultSavePath`
- [x] `TransferInfo` — `GET /api/v2/transfer/info`
- [x] `SetDownloadLimit` / `SetUploadLimit` / `SetShareLimits` — per-torrent limits
- [x] `TorrentsFiles(hash)` — `GET /api/v2/torrents/files` (list files within a torrent)
- [x] `TorrentsFilePrio(hash, ids, priority)` — `POST /api/v2/torrents/filePrio`
- [x] Category CRUD — `categories`, `createCategory`, `editCategory`, `removeCategories`

## API consistency

- [x] Standardize hash parameters: `TorrentsInfo` takes `[]string` and joins internally, but `TorrentsAddTags`/`TorrentsRemoveTags`/`TorrentsDelete`/`SetForceStart` require pre-formatted pipe-separated strings or single hashes. Pick one convention (prefer `[]string` with internal join).
- [x] `Category` type is `map[string]interface{}` — should be a proper struct with `Name` and `SavePath` fields.

## Missing fields

- [x] `TorrentInfo`: missing `download_path`, `infohash_v1`, `infohash_v2`, `popularity`
- [x] `TrackerInfo`: missing `num_seeds`, `num_leeches`, `num_downloaded`
- [x] `TorrentsAddParams`: missing `contentLayout` (replaces deprecated `root_folder` in qBittorrent 4.4+), `stopCondition`, `addToTopOfQueue`, `inactiveSeedingTimeLimit`
