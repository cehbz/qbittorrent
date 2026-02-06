# qBittorrent Go Client Library

A Go client library for interacting with the [qBittorrent](https://www.qbittorrent.org/) Web API.

## Features

- **Authentication**: Log in to the qBittorrent Web API.
- **Context Support**: All methods now have context-aware versions (`*Ctx`) for timeout and cancellation control.
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
client, err := qbittorrent.NewClient("username", "password", "http://localhost:8080", nil)
if err != nil {
    log.Fatalf("Failed to create client: %v", err)
}
```

- `username`: Your qBittorrent Web UI username. Empty if none.
- `password`: Your qBittorrent Web UI password. Empty if none.
- `baseURL`: The full URL where qBittorrent is running (e.g., `"http://127.0.0.1:8080"`).
- `httpClient`: Optional custom `*http.Client`. Pass `nil` to use `http.DefaultClient`.

### Using Context for Timeout and Cancellation

All API methods now have context-aware versions (suffixed with `Ctx`) for better control over request lifecycles:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Use context-aware methods
torrents, err := client.TorrentsInfoCtx(ctx)
if err != nil {
    log.Fatalf("Failed to retrieve torrents: %v", err)
}
```

**Benefits of Context-Aware Methods:**
- **Timeout Control**: Prevent operations from running indefinitely
- **Cancellation**: Cancel long-running operations when needed
- **Request Tracing**: Integrate with distributed tracing systems
- **Graceful Shutdown**: Clean up resources during application shutdown

**Backward Compatibility:**
All original methods (without `Ctx` suffix) remain available and internally use `context.Background()`.

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

#### With Context (Recommended)
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

torrentData, err := os.ReadFile("path/to/your.torrent")
if err != nil {
    log.Fatalf("Failed to read torrent file: %v", err)
}

err = client.TorrentsAddCtx(ctx, "your.torrent", torrentData)
if err != nil {
    log.Fatalf("Failed to add torrent: %v", err)
}
```

#### Advanced Method with Parameters
```go
ctx := context.Background()
params := &qbittorrent.TorrentsAddParams{
    Torrents:  [][]byte{torrentData},
    SavePath:  "/downloads/movies",
    Category:  "movies",
    Tags:      "hd,x265",
    SkipCheck: true,
    Paused:    false,
    AutoTMM:   true,
}

err = client.TorrentsAddParamsCtx(ctx, params)
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
// Get all torrents (with context recommended)
ctx := context.Background()
torrents, err := client.TorrentsInfoCtx(ctx)
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

filteredTorrents, err := client.TorrentsInfoCtx(ctx, params)
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

## API Reference

### Context-Aware Methods

All public API methods now have context-aware versions for production use:

| Original Method | Context-Aware Version | Description |
|----------------|----------------------|-------------|
| `AuthLogin()` | `AuthLoginCtx(ctx)` | Authenticate with qBittorrent |
| `TorrentsInfo()` | `TorrentsInfoCtx(ctx, params...)` | Get torrent information |
| `TorrentsAdd()` | `TorrentsAddCtx(ctx, file, data)` | Add a torrent |
| `TorrentsAddParams()` | `TorrentsAddParamsCtx(ctx, params)` | Add torrent with parameters |
| `TorrentsDelete()` | `TorrentsDeleteCtx(ctx, hash)` | Delete a torrent |
| `TorrentsExport()` | `TorrentsExportCtx(ctx, hash)` | Export torrent file |
| `TorrentsDownload()` | `TorrentsDownloadCtx(ctx, hash)` | Download torrent file |
| `TorrentsProperties()` | `TorrentsPropertiesCtx(ctx, hash)` | Get torrent properties |
| `TorrentsTrackers()` | `TorrentsTrackersCtx(ctx, hash)` | Get torrent trackers |
| `SetForceStart()` | `SetForceStartCtx(ctx, hash, value)` | Set force start |
| `TorrentsAddTags()` | `TorrentsAddTagsCtx(ctx, hashes, tags)` | Add tags to torrents |
| `TorrentsRemoveTags()` | `TorrentsRemoveTagsCtx(ctx, hashes, tags)` | Remove tags from torrents |
| `TorrentsGetTags()` | `TorrentsGetTagsCtx(ctx, hashes)` | Get torrent tags |
| `TorrentsGetAllTags()` | `TorrentsGetAllTagsCtx(ctx)` | Get all available tags |
| `TorrentsCreateTags()` | `TorrentsCreateTagsCtx(ctx, tags)` | Create new tags |
| `TorrentsDeleteTags()` | `TorrentsDeleteTagsCtx(ctx, tags)` | Delete tags |
| `SyncMainData()` | `SyncMainDataCtx(ctx, rid)` | Sync main data |
| `SyncTorrentPeers()` | `SyncTorrentPeersCtx(ctx, hash, rid)` | Sync torrent peers |

### Best Practices

1. **Always use context-aware methods in production** to enable timeouts and cancellation
2. **Set appropriate timeouts** based on operation type:
   - Quick operations (info, properties): 5-10 seconds
   - File uploads/downloads: 30-60 seconds
   - Long-running operations: Use context with cancellation support
3. **Handle context errors appropriately**:
   ```go
   if err == context.DeadlineExceeded {
       // Handle timeout
   } else if err == context.Canceled {
       // Handle cancellation
   }
   ```
4. **Use context for graceful shutdown**:
   ```go
   ctx, cancel := context.WithCancel(context.Background())
   defer cancel()

   // On shutdown signal
   go func() {
       <-shutdownSignal
       cancel() // This will cancel all in-flight requests
   }()
   ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contribution

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## Acknowledgments

- [qBittorrent Web API Documentation](https://github.com/qbittorrent/qBittorrent/wiki#WebUI-API)
