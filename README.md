# qBittorrent Go Client Library

A Go client library for interacting with the [qBittorrent](https://www.qbittorrent.org/) Web API.

## Features

- **Authentication**: Log in and log out of the qBittorrent Web API.
- **App**: Get version, preferences, and default save path.
- **Torrent Management**:
  - Add, delete, pause, resume, recheck, and reannounce torrents.
  - Export torrent files.
  - Retrieve torrent information and properties.
  - Rename torrents, set location, force-start, and manage auto-TMM.
  - List and set priority for files within a torrent.
  - Set per-torrent download/upload/share limits.
- **Tracker Information**: Fetch tracker details for specific torrents.
- **Tag Management**: Create, delete, add, and remove tags from torrents.
- **Category Management**: Create, edit, remove, and list categories.
- **Transfer Info**: Get global transfer statistics.
- **Synchronization**: Get main data and peer information for real-time updates.

## Installation

```bash
go get github.com/cehbz/qbittorrent
```

## Usage

### Initializing the Client

```go
import "github.com/cehbz/qbittorrent"

client, err := qbittorrent.NewClient("username", "password", "http://localhost:8080")
if err != nil {
    log.Fatalf("Failed to create client: %v", err)
}
```

- `username`: Your qBittorrent Web UI username. Empty if none.
- `password`: Your qBittorrent Web UI password. Empty if none.
- `baseURL`: The full base URL where qBittorrent is running (e.g., `"http://localhost:8080"`).

### Adding a Torrent

#### Simple Method
```go
torrentData, err := os.ReadFile("path/to/your.torrent")
if err != nil {
    log.Fatalf("Failed to read torrent file: %v", err)
}

err = client.TorrentsAdd(torrentData)
if err != nil {
    log.Fatalf("Failed to add torrent: %v", err)
}
```

#### Advanced Method with Parameters
```go
params := &qbittorrent.TorrentsAddParams{
    Torrents:      [][]byte{torrentData},
    SavePath:      "/downloads/movies",
    Category:      "movies",
    Tags:          "hd,x265",
    SkipCheck:     true,
    Paused:        boolPtr(false),
    AutoTMM:       boolPtr(true),
    ContentLayout: "Original",   // "Original", "Subfolder", or "NoSubfolder"
    StopCondition: "MetadataReceived", // or "FilesChecked"
}

err = client.TorrentsAddParams(params)
if err != nil {
    log.Fatalf("Failed to add torrent: %v", err)
}
```

Note: `Paused`, `AutoTMM`, `RootFolder`, and `AddToTopOfQueue` are `*bool` fields. Use a helper to set them:
```go
func boolPtr(b bool) *bool { return &b }
```

### Deleting Torrents

```go
err := client.TorrentsDelete([]string{"hash1", "hash2"}, true) // true = delete files
if err != nil {
    log.Fatalf("Failed to delete torrent: %v", err)
}
```

### Pause, Resume, Recheck, Reannounce

```go
client.TorrentsPause([]string{"hash1", "hash2"})
client.TorrentsResume([]string{"hash1"})
client.TorrentsRecheck([]string{"hash1"})
client.TorrentsReannounce([]string{"hash1"})
```

### Exporting a Torrent File

```go
data, err := client.TorrentsExport("torrent-hash")
if err != nil {
    log.Fatalf("Failed to export torrent: %v", err)
}
os.WriteFile("exported.torrent", data, 0644)
```

### Retrieving Torrent Information

```go
// Get all torrents
torrents, err := client.TorrentsInfo()

// Get specific torrents with filters
params := &qbittorrent.TorrentsInfoParams{
    Category: "movies",
    Sort:     "progress",
    Reverse:  true,
    Limit:    10,
}
filteredTorrents, err := client.TorrentsInfo(params)
```

### Retrieving Torrent Properties

```go
props, err := client.TorrentsProperties("torrent-hash")
fmt.Printf("Save path: %s, Total size: %d\n", props.SavePath, props.TotalSize)
```

Timestamp fields (`AdditionDate`, `CreationDate`, `CompletionDate`, `LastSeen`) are parsed as `time.Time` values from Unix seconds.

`TorrentInfo` timestamp fields (`AddedOn`, `CompletionOn`, `LastActivity`, `SeenComplete`) are also parsed as `time.Time`. A value of `-1` from the API is converted to the zero `time.Time`.

### Torrent Settings

```go
client.SetForceStart([]string{"hash1"}, true)
client.TorrentsSetLocation([]string{"hash1"}, "/new/path")
client.TorrentsSetCategory([]string{"hash1"}, "movies")
client.TorrentsSetAutoTMM([]string{"hash1"}, true)
client.TorrentsRename("hash1", "New Name")
```

### Per-Torrent Limits

```go
client.TorrentsSetDownloadLimit([]string{"hash1"}, 1048576) // 1 MB/s
client.TorrentsSetUploadLimit([]string{"hash1"}, 524288)    // 512 KB/s
client.TorrentsSetShareLimits([]string{"hash1"}, 2.0, 1440, 720) // ratio, seeding time, inactive time
```

### Torrent Files

```go
files, err := client.TorrentsFiles("hash1")
for _, f := range files {
    fmt.Printf("%s (%d bytes, priority %d)\n", f.Name, f.Size, f.Priority)
}

// Set file priority (7 = maximum)
client.TorrentsFilePrio("hash1", []int{0, 1}, 7)
```

### Tracker Information

```go
trackers, err := client.TorrentsTrackers("torrent-hash")
for _, tracker := range trackers {
    fmt.Printf("URL: %s, Status: %d\n", tracker.URL, tracker.Status)
}
```

### Tag Management

```go
client.TorrentsCreateTags("hd,x265,4k")
client.TorrentsAddTags([]string{"hash1", "hash2"}, "hd,x265")
client.TorrentsRemoveTags([]string{"hash1"}, "old-tag")

tags, _ := client.TorrentsGetTags([]string{"hash1"})
allTags, _ := client.TorrentsGetAllTags()

client.TorrentsDeleteTags("unused-tag")
```

### Category Management

```go
categories, _ := client.TorrentsCategories()
client.TorrentsCreateCategory("movies", "/downloads/movies")
client.TorrentsEditCategory("movies", "/new/path")
client.TorrentsRemoveCategories([]string{"movies", "tv"})
```

### App Info

```go
version, _ := client.AppVersion()         // e.g. "v4.6.0"
savePath, _ := client.AppDefaultSavePath() // e.g. "/downloads"
prefs, _ := client.AppPreferences()        // map[string]any
client.SetAppPreferences(map[string]any{"save_path": "/new/downloads"})
client.AuthLogout()
```

### Transfer Info

```go
info, _ := client.TransferGetInfo()
fmt.Printf("DL: %d B/s, UL: %d B/s, DHT: %d nodes\n",
    info.DLInfoSpeed, info.UPInfoSpeed, info.DHTNodes)
```

### Synchronization

```go
mainData, _ := client.SyncMainData(0)
nextRid := mainData.Rid

peers, _ := client.SyncTorrentPeers("torrent-hash", 0)
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contribution

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## Acknowledgments

- [qBittorrent Web API Documentation](https://github.com/qbittorrent/qBittorrent/wiki#WebUI-API)
