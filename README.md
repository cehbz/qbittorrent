# qBittorrent Go Client Library

A Go client library for interacting with the [qBittorrent](https://www.qbittorrent.org/) Web API.

## Features

- **Authentication**: Log in to the qBittorrent Web API.
- **Torrent Management**:
  - Add new torrents (with flexible parameters).
  - Delete existing torrents.
  - Export torrent files.
  - Download torrent files.
  - Retrieve torrent information.
  - Retrieve torrent generic properties.
  - Manage torrent force-start settings.
- **Tracker Information**: Fetch tracker details for specific torrents.
- **Tag Management**: Create, delete, add, and remove tags from torrents.
- **Synchronization**: Get main data and peer information for real-time updates.

## Installation

To install the package, run:

```bash
go get github.com/cehbz/qbittorrent
```

## Usage

### Importing the Package

```go
import (
    "github.com/cehbz/qbittorrent"
)
```

### Initializing the Client

```go
client, err := qbittorrent.NewClient("username", "password", "localhost", "8080")
if err != nil {
    log.Fatalf("Failed to create client: %v", err)
}
```

- `username`: Your qBittorrent Web UI username. Empty if none.
- `password`: Your qBittorrent Web UI password. Empty if none.
- `addr`: The address where qBittorrent is running (e.g., `"127.0.0.1"`).
- `port`: The port number of the qBittorrent Web UI (e.g., `"8080"`).

### Adding a Torrent

#### Simple Method
```go
torrentData, err := os.ReadFile("path/to/your.torrent")
if err != nil {
    log.Fatalf("Failed to read torrent file: %v", err)
}

err = client.TorrentsAdd("your.torrent", torrentData)
if err != nil {
    log.Fatalf("Failed to add torrent: %v", err)
}
```

#### Advanced Method with Parameters
```go
params := &qbittorrent.TorrentsAddParams{
    Torrents:  [][]byte{torrentData},
    SavePath:  "/downloads/movies",
    Category:  "movies",
    Tags:      "hd,x265",
    SkipCheck: true,
    Paused:    false,
    AutoTMM:   true,
}

err = client.TorrentsAddParams(params)
if err != nil {
    log.Fatalf("Failed to add torrent: %v", err)
}
```

### Deleting a Torrent

```go
err := client.TorrentsDelete("torrent-hash")
if err != nil {
    log.Fatalf("Failed to delete torrent: %v", err)
}
```

### Exporting a Torrent File

```go
data, err := client.TorrentsExport("torrent-hash")
if err != nil {
    log.Fatalf("Failed to export torrent: %v", err)
}

err = os.WriteFile("exported.torrent", data, 0644)
if err != nil {
    log.Fatalf("Failed to write exported torrent file: %v", err)
}
```

### Downloading a Torrent File

```go
data, err := client.TorrentsDownload("torrent-hash")
if err != nil {
    log.Fatalf("Failed to download torrent: %v", err)
}

err = os.WriteFile("downloaded.torrent", data, 0644)
if err != nil {
    log.Fatalf("Failed to write downloaded torrent file: %v", err)
}
```

### Retrieving Torrent Information

```go
// Get all torrents
torrents, err := client.TorrentsInfo()
if err != nil {
    log.Fatalf("Failed to retrieve torrents info: %v", err)
}

for _, torrent := range torrents {
    fmt.Printf("Name: %s, Progress: %.2f%%\n", torrent.Name, torrent.Progress*100)
}

// Get specific torrents with filters
params := &qbittorrent.TorrentsInfoParams{
    Category: "movies",
    Sort:     "progress",
    Reverse:  true,
    Limit:    10,
}

filteredTorrents, err := client.TorrentsInfo(params)
if err != nil {
    log.Fatalf("Failed to retrieve filtered torrents: %v", err)
}
```

### Retrieving Torrent Generic Properties

```go
props, err := client.TorrentsProperties("torrent-hash")
if err != nil {
    log.Fatalf("Failed to retrieve torrent properties: %v", err)
}

fmt.Printf("Save path: %s, Total size: %d\n", props.SavePath, props.TotalSize)
```

`TorrentsProperties` timestamp fields (`AdditionDate`, `CreationDate`, `CompletionDate`, `LastSeen`) are returned as `time.Time` values parsed from Unix seconds.

### Managing Torrent Force Start

```go
err := client.SetForceStart("torrent-hash", true)
if err != nil {
    log.Fatalf("Failed to set force start: %v", err)
}
```

### Fetching Tracker Information

```go
trackers, err := client.TorrentsTrackers("torrent-hash")
if err != nil {
    log.Fatalf("Failed to get trackers: %v", err)
}

for _, tracker := range trackers {
    fmt.Printf("URL: %s, Status: %d\n", tracker.URL, tracker.Status)
}
```

### Tag Management

```go
// Create new tags
err := client.TorrentsCreateTags("hd,x265,4k")
if err != nil {
    log.Fatalf("Failed to create tags: %v", err)
}

// Add tags to torrents
err = client.TorrentsAddTags("torrent-hash1|torrent-hash2", "hd,x265")
if err != nil {
    log.Fatalf("Failed to add tags: %v", err)
}

// Remove tags from torrents
err = client.TorrentsRemoveTags("torrent-hash", "old-tag")
if err != nil {
    log.Fatalf("Failed to remove tags: %v", err)
}

// Get tags for specific torrents
tags, err := client.TorrentsGetTags("torrent-hash")
if err != nil {
    log.Fatalf("Failed to get tags: %v", err)
}

// Get all available tags
allTags, err := client.TorrentsGetAllTags()
if err != nil {
    log.Fatalf("Failed to get all tags: %v", err)
}

// Delete tags
err = client.TorrentsDeleteTags("unused-tag")
if err != nil {
    log.Fatalf("Failed to delete tags: %v", err)
}
```

### Synchronization (Real-time Updates)

```go
// Get main data for synchronization
mainData, err := client.SyncMainData(0) // 0 for initial request
if err != nil {
    log.Fatalf("Failed to get main data: %v", err)
}

// Use the returned rid for subsequent requests
nextRid := mainData.Rid

// Get peer information for a specific torrent
peers, err := client.SyncTorrentPeers("torrent-hash", 0)
if err != nil {
    log.Fatalf("Failed to get peer data: %v", err)
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contribution

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## Acknowledgments

- [qBittorrent Web API Documentation](https://github.com/qbittorrent/qBittorrent/wiki#WebUI-API)
