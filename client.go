package qbittorrent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type InfoHash string

// joinHashes joins a slice of hashes with the pipe separator expected by the qBittorrent API.
func joinHashes(hashes []string) string {
	return strings.Join(hashes, "|")
}

// Client is used to interact with the qBittorrent API
type Client struct {
	username string
	password string
	client   *http.Client
	baseURL  string
	sid      string // store the SID cookie
	mu       sync.RWMutex
}

// TorrentInfo represents the structured information of a torrent from the qBittorrent API
type TorrentInfo struct {
	AddedOn            int64    `json:"added_on"`
	AmountLeft         int64    `json:"amount_left"`
	AutoTMM            bool     `json:"auto_tmm"`
	Availability       float64  `json:"availability"`
	Category           string   `json:"category"`
	Completed          int64    `json:"completed"`
	CompletionOn       int64    `json:"completion_on"`
	ContentPath        string   `json:"content_path"`
	DownloadPath       string   `json:"download_path"`
	DLLimit            int64    `json:"dl_limit"`
	DLSpeed            int64    `json:"dlspeed"`
	Downloaded         int64    `json:"downloaded"`
	DownloadedSession  int64    `json:"downloaded_session"`
	ETA                int64    `json:"eta"`
	FirstLastPiecePrio bool     `json:"f_l_piece_prio"`
	ForceStart         bool     `json:"force_start"`
	Hash               InfoHash `json:"hash"`
	InfoHashV1         InfoHash `json:"infohash_v1"`
	InfoHashV2         InfoHash `json:"infohash_v2"`
	IsPrivate          bool     `json:"isPrivate"`
	LastActivity       int64    `json:"last_activity"`
	MagnetURI          string   `json:"magnet_uri"`
	MaxRatio           float64  `json:"max_ratio"`
	MaxSeedingTime     int64    `json:"max_seeding_time"`
	Name               string   `json:"name"`
	NumComplete        int64    `json:"num_complete"`
	NumIncomplete      int64    `json:"num_incomplete"`
	NumLeechs          int64    `json:"num_leechs"`
	NumSeeds           int64    `json:"num_seeds"`
	Popularity         float64  `json:"popularity"`
	Priority           int64    `json:"priority"`
	Progress           float64  `json:"progress"`
	Ratio              float64  `json:"ratio"`
	RatioLimit         float64  `json:"ratio_limit"`
	SavePath           string   `json:"save_path"`
	SeedingTime        int64    `json:"seeding_time"`
	SeedingTimeLimit   int64    `json:"seeding_time_limit"`
	SeenComplete       int64    `json:"seen_complete"`
	SequentialDownload bool     `json:"seq_dl"`
	Size               int64    `json:"size"`
	State              string   `json:"state"`
	SuperSeeding       bool     `json:"super_seeding"`
	Tags               []string `json:"-"`
	TimeActive         int64    `json:"time_active"`
	TotalSize          int64    `json:"total_size"`
	Tracker            string   `json:"tracker"`
	UpLimit            int64    `json:"up_limit"`
	Uploaded           int64    `json:"uploaded"`
	UploadedSession    int64    `json:"uploaded_session"`
	UpSpeed            int64    `json:"upspeed"`
}

// TorrentsProperties represents generic properties for a torrent.
/*
  "addition_date": 1770257484,
  "comment": "https://redacted.sh/torrents.php?torrentid=664915",
  "completion_date": 1770257488,
  "created_by": "Transmission/2.84 (14306)",
  "creation_date": 1483593698,
  "dl_limit": -1,
  "dl_speed": 0,
  "dl_speed_avg": 100351515,
  "download_path": "",
  "eta": 8640000,
  "has_metadata": true,
  "hash": "5a0ff0482cb309913568bee1db6d68f7e5ef1f6d",
  "infohash_v1": "5a0ff0482cb309913568bee1db6d68f7e5ef1f6d",
  "infohash_v2": "",
  "is_private": true,
  "last_seen": 1770257488,
  "name": "Jesca Hoop - Hunting My Dress (2010 Deluxe Edition) [FLAC]",
  "nb_connections": 0,
  "nb_connections_limit": -1,
  "peers": 0,
  "peers_total": 0,
  "piece_size": 262144,
  "pieces_have": 1531,
  "pieces_num": 1531,
  "popularity": 0,
  "private": true,
  "reannounce": 1434,
  "save_path": "/home/haynes/torrents/qbittorrent/music",
  "seeding_time": 6698,
  "seeds": 0,
  "seeds_total": 11,
  "share_ratio": 0,
  "time_elapsed": 6702,
  "total_downloaded": 401406062,
  "total_downloaded_session": 401406062,
  "total_size": 401085347,
  "total_uploaded": 0,
  "total_uploaded_session": 0,
  "total_wasted": 320715,
  "up_limit": -1,
  "up_speed": 0,
  "up_speed_avg": 0
*/
type TorrentsProperties struct {
	AdditionDate           time.Time
	Comment                string `json:"comment"`
	CompletionDate         time.Time
	CreatedBy              string `json:"created_by"`
	CreationDate           time.Time
	DLLimit                int64    `json:"dl_limit"`
	DLSpeed                int64    `json:"dl_speed"`
	DLSpeedAvg             int64    `json:"dl_speed_avg"`
	DownloadPath           string   `json:"download_path"`
	ETA                    int64    `json:"eta"`
	HasMetadata            bool     `json:"has_metadata"`
	Hash                   InfoHash `json:"hash"`
	InfoHashV1             InfoHash `json:"infohash_v1"`
	InfoHashV2             InfoHash `json:"infohash_v2"`
	IsPrivate              bool     `json:"is_private"`
	LastSeen               time.Time
	Name                   string  `json:"name"`
	NbConnections          int64   `json:"nb_connections"`
	NbConnectionsLimit     int64   `json:"nb_connections_limit"`
	Peers                  int64   `json:"peers"`
	PeersTotal             int64   `json:"peers_total"`
	PiecesHave             int64   `json:"pieces_have"`
	PieceSize              int64   `json:"piece_size"`
	PiecesNum              int64   `json:"pieces_num"`
	Popularity             float64 `json:"popularity"`
	Private                bool    `json:"private"`
	Reannounce             int64   `json:"reannounce"`
	SavePath               string  `json:"save_path"`
	SeedingTime            int64   `json:"seeding_time"`
	Seeds                  int64   `json:"seeds"`
	SeedsTotal             int64   `json:"seeds_total"`
	ShareRatio             float64 `json:"share_ratio"`
	TimeElapsed            int64   `json:"time_elapsed"`
	TotalDownloaded        int64   `json:"total_downloaded"`
	TotalDownloadedSession int64   `json:"total_downloaded_session"`
	TotalSize              int64   `json:"total_size"`
	TotalUploaded          int64   `json:"total_uploaded"`
	TotalUploadedSession   int64   `json:"total_uploaded_session"`
	TotalWasted            int64   `json:"total_wasted"`
	UpLimit                int64   `json:"up_limit"`
	UpSpeed                int64   `json:"up_speed"`
	UpSpeedAvg             int64   `json:"up_speed_avg"`
}

// TODO: Apply alias-based timestamp parsing to other structs.

// UnmarshalJSON custom unmarshaller for TorrentsProperties to handle timestamps.
func (t *TorrentsProperties) UnmarshalJSON(data []byte) error {
	type Alias TorrentsProperties
	aux := &struct {
		AdditionDate   int64 `json:"addition_date"`
		CompletionDate int64 `json:"completion_date"`
		CreationDate   int64 `json:"creation_date"`
		LastSeen       int64 `json:"last_seen"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.AdditionDate = unixTimeOrZero(aux.AdditionDate)
	t.CompletionDate = unixTimeOrZero(aux.CompletionDate)
	t.CreationDate = unixTimeOrZero(aux.CreationDate)
	t.LastSeen = unixTimeOrZero(aux.LastSeen)

	return nil
}

func unixTimeOrZero(value int64) time.Time {
	if value == -1 {
		return time.Time{}
	}
	return time.Unix(value, 0)
}

// UnmarshalJSON custom unmarshaller for TorrentInfo to handle Tags
func (t *TorrentInfo) UnmarshalJSON(data []byte) error {
	type Alias TorrentInfo
	aux := &struct {
		RawTags string `json:"tags"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.RawTags == "" {
		t.Tags = []string{}
	} else {
		parts := strings.Split(aux.RawTags, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		t.Tags = parts
	}
	return nil
}

// TrackerInfo represents a tracker info for a torrent
type TrackerInfo struct {
	URL           string `json:"url"`
	Status        int    `json:"status"`
	Tier          int    `json:"tier"`
	NumPeers      int    `json:"num_peers"`
	NumSeeds      int    `json:"num_seeds"`
	NumLeeches    int    `json:"num_leeches"`
	NumDownloaded int    `json:"num_downloaded"`
	Msg           string `json:"msg"`
}

// Category represents a torrent category with its save path.
type Category struct {
	Name     string `json:"name"`
	SavePath string `json:"savePath"`
}

// fields might be missing, in which case we need to switch to pointers and allow "omitempty"
// https://github.com/qbittorrent/qBittorrent/blob/master/src/base/json_api.cpp#L101
// MainData is the data returned by the /api/v2/sync/maindata endpoint
type MainData struct {
	Categories        map[string]Category    `json:"categories"`
	CategoriesRemoved []string               `json:"categories_removed"`
	FullUpdate        bool                   `json:"full_update"`
	Rid               int                    `json:"rid"`
	ServerState       ServerState            `json:"server_state"`
	Tags              []string               `json:"tags"`
	TagsRemoved       []string               `json:"tags_removed"`
	Torrents          map[string]TorrentInfo `json:"torrents"`
	TorrentsRemoved   []string               `json:"torrents_removed"`
	Trackers          map[string][]InfoHash  `json:"trackers"` // maps trackers to infohashes
}

type ServerState struct {
	AllTimeDL            int64  `json:"alltime_dl"`
	AllTimeUL            int64  `json:"alltime_ul"`
	AverageTimeQueue     int    `json:"average_time_queue"`
	ConnectionStatus     string `json:"connection_status"`
	DHTNodes             int    `json:"dht_nodes"`
	DLInfoData           int64  `json:"dl_info_data"`
	DLInfoSpeed          int    `json:"dl_info_speed"`
	DLRateLimit          int    `json:"dl_rate_limit"`
	FreeSpaceOnDisk      int64  `json:"free_space_on_disk"`
	GlobalRatio          string `json:"global_ratio"`
	QueuedIOJobs         int    `json:"queued_io_jobs"`
	Queueing             bool   `json:"queueing"`
	ReadCacheHits        string `json:"read_cache_hits"`
	ReadCacheOverload    string `json:"read_cache_overload"`
	RefreshInterval      int    `json:"refresh_interval"`
	TotalBuffersSize     int64  `json:"total_buffers_size"`
	TotalPeerConnections int    `json:"total_peer_connections"`
	TotalQueuedSize      int64  `json:"total_queued_size"`
	TotalWastedSession   int64  `json:"total_wasted_session"`
	UpInfoData           int64  `json:"up_info_data"`
	UpInfoSpeed          int    `json:"up_info_speed"`
	UpRateLimit          int    `json:"up_rate_limit"`
	UseAltSpeedLimits    bool   `json:"use_alt_speed_limits"`
	UseSubcategories     bool   `json:"use_subcategories"`
	WriteCacheOverload   string `json:"write_cache_overload"`
}

type TorrentPeer struct {
	Client       string  `json:"client"`
	Connection   string  `json:"connection"`
	Country      string  `json:"country"`
	CountryCode  string  `json:"country_code"`
	DLSpeed      int64   `json:"dl_speed"`
	Downloaded   int64   `json:"downloaded"`
	Files        string  `json:"files"`
	Flags        string  `json:"flags"`
	FlagsDesc    string  `json:"flags_desc"`
	IP           string  `json:"ip"`
	PeerIDClient string  `json:"peer_id_client"`
	Port         int     `json:"port"`
	Progress     float64 `json:"progress"`
	Relevance    float64 `json:"relevance"`
	Uploaded     int64   `json:"uploaded"`
	UPSpeed      int64   `json:"up_speed"`
}

type TorrentPeers struct {
	FullUpdate bool                   `json:"full_update"`
	Peers      map[string]TorrentPeer `json:"peers"`
	// PeersRemoved map[string][]string    `json:"peers_removed"`
	Rid       int  `json:"rid"`
	ShowFlags bool `json:"show_flags"`
}

// NewClient initializes a new qBittorrent client.
// If httpClient is nil, http.DefaultClient is used.
func NewClient(username, password, baseURL string, httpClient ...*http.Client) (*Client, error) {
	// Use the provided http.Client if given, otherwise use http.DefaultClient
	client := http.DefaultClient
	if len(httpClient) > 0 && httpClient[0] != nil {
		client = httpClient[0]
	}

	// Create and return the Client instance
	qbClient := &Client{
		username: username,
		password: password,
		client:   client,
		baseURL:  baseURL,
	}

	// Authenticate if username and password are provided
	if username != "" && password != "" {
		if err := qbClient.AuthLogin(); err != nil {
			return nil, fmt.Errorf("AuthLogin error: %v", err)
		}
	}

	return qbClient, nil
}

// AuthLogin logs in to the qBittorrent Web API
func (c *Client) AuthLogin() error {
	data := url.Values{}
	data.Set("username", c.username)
	data.Set("password", c.password)

	resp, err := c.doPostResponse("/api/v2/auth/login", strings.NewReader(data.Encode()), "application/x-www-form-urlencoded")
	if err != nil {
		return fmt.Errorf("AuthLogin error: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("AuthLogin error (%d): %s", resp.StatusCode, string(respBody))
	}
	defer resp.Body.Close()

	// Extract the SID cookie from the response
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "SID" {
			c.mu.Lock()
			c.sid = cookie.Value
			c.mu.Unlock()
			break
		}
	}

	return nil
}

// TorrentsExport retrieves the .torrent file for a given torrent hash
func (c *Client) TorrentsExport(hash string) ([]byte, error) {
	params := url.Values{}
	params.Set("hash", hash)

	// Use the GET request helper
	return c.doPostValues("/api/v2/torrents/export", params)
}

// TorrentsAddParams holds the parameters for adding a torrent
type TorrentsAddParams struct {
	Torrents    [][]byte // Raw torrent data
	URLs        []string // Magnet links or URLs
	SavePath    string   // Download folder
	Cookie      string   // Cookie sent to download the .torrent file
	Category    string   // Category for the torrent
	Tags        string   // Tags for the torrent, comma separated
	SkipCheck               bool     // Skip hash checking
	Paused                  *bool    // Add torrents in the paused state
	RootFolder              *bool    // Create the root folder (default: true)
	ContentLayout           string   // Content layout: "Original", "Subfolder", "NoSubfolder"
	StopCondition           string   // Stop condition: "MetadataReceived" or "FilesChecked"
	Rename                  string   // Rename torrent
	UpLimit                 int      // Set torrent upload speed limit (bytes/second)
	DlLimit                 int      // Set torrent download speed limit (bytes/second)
	RatioLimit              float64  // Set torrent share ratio limit
	SeedingTime             int      // Set torrent seeding time limit (minutes)
	InactiveSeedingTimeLimit int     // Set inactive seeding time limit (minutes)
	AutoTMM                 *bool    // Whether Automatic Torrent Management should be used
	Sequential              bool     // Enable sequential download
	FirstLast               bool     // Prioritize download first last piece
	AddToTopOfQueue         *bool    // Add torrent to top of queue
}

// TorrentsAddParams adds a torrent using a more flexible parameter structure
func (c *Client) TorrentsAddParams(params *TorrentsAddParams) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add torrent files
	for i, torrentData := range params.Torrents {
		part, err := writer.CreateFormFile("torrents", fmt.Sprintf("torrent%d.torrent", i))
		if err != nil {
			return fmt.Errorf("CreateFormFile error: %v", err)
		}
		if _, err := io.Copy(part, bytes.NewReader(torrentData)); err != nil {
			return fmt.Errorf("io.Copy error: %v", err)
		}
	}

	// Add URLs
	for _, url := range params.URLs {
		_ = writer.WriteField("urls", url)
	}

	// Add other fields
	if params.SavePath != "" {
		_ = writer.WriteField("savepath", params.SavePath)
	}
	if params.Cookie != "" {
		_ = writer.WriteField("cookie", params.Cookie)
	}
	if params.Category != "" {
		_ = writer.WriteField("category", params.Category)
	}
	if params.Tags != "" {
		_ = writer.WriteField("tags", params.Tags)
	}
	if params.SkipCheck {
		_ = writer.WriteField("skip_checking", "true")
	}
	if params.Paused != nil {
		_ = writer.WriteField("paused", fmt.Sprintf("%t", *params.Paused))
	}
	if params.RootFolder != nil {
		_ = writer.WriteField("root_folder", fmt.Sprintf("%t", *params.RootFolder))
	}
	if params.ContentLayout != "" {
		_ = writer.WriteField("contentLayout", params.ContentLayout)
	}
	if params.StopCondition != "" {
		_ = writer.WriteField("stopCondition", params.StopCondition)
	}
	if params.Rename != "" {
		_ = writer.WriteField("rename", params.Rename)
	}
	if params.UpLimit > 0 {
		_ = writer.WriteField("upLimit", strconv.Itoa(params.UpLimit))
	}
	if params.DlLimit > 0 {
		_ = writer.WriteField("dlLimit", strconv.Itoa(params.DlLimit))
	}
	if params.RatioLimit > 0 {
		_ = writer.WriteField("ratioLimit", fmt.Sprintf("%.2f", params.RatioLimit))
	}
	if params.SeedingTime > 0 {
		_ = writer.WriteField("seedingTimeLimit", strconv.Itoa(params.SeedingTime))
	}
	if params.InactiveSeedingTimeLimit > 0 {
		_ = writer.WriteField("inactiveSeedingTimeLimit", strconv.Itoa(params.InactiveSeedingTimeLimit))
	}
	if params.AutoTMM != nil {
		_ = writer.WriteField("autoTMM", fmt.Sprintf("%t", *params.AutoTMM))
	}
	if params.Sequential {
		_ = writer.WriteField("sequentialDownload", "true")
	}
	if params.FirstLast {
		_ = writer.WriteField("firstLastPiecePrio", "true")
	}
	if params.AddToTopOfQueue != nil {
		_ = writer.WriteField("addToTopOfQueue", fmt.Sprintf("%t", *params.AddToTopOfQueue))
	}

	writer.Close()

	_, err := c.doPost("/api/v2/torrents/add", &body, writer.FormDataContentType())
	if err != nil {
		return fmt.Errorf("TorrentsAdd error: %v", err)
	}
	return nil
}

// boolPtr returns a pointer to the given bool value.
func boolPtr(b bool) *bool {
	return &b
}

// TorrentsAdd adds a torrent to qBittorrent via Web API using multipart/form-data
func (c *Client) TorrentsAdd(fileData []byte) error {
	params := &TorrentsAddParams{
		Torrents:  [][]byte{fileData},
		SkipCheck: true,
		Paused:    boolPtr(false),
		AutoTMM:   boolPtr(false),
	}
	return c.TorrentsAddParams(params)
}

// TorrentsDelete deletes torrents from qBittorrent by their hashes
func (c *Client) TorrentsDelete(hashes []string, deleteFiles bool) error {
	data := url.Values{}
	data.Set("hashes", joinHashes(hashes))
	data.Set("deleteFiles", fmt.Sprintf("%t", deleteFiles))

	_, err := c.doPostValues("/api/v2/torrents/delete", data)
	if err != nil {
		return fmt.Errorf("TorrentsDelete error: %v", err)
	}
	return nil
}

// SetForceStart enables force start for the specified torrents
func (c *Client) SetForceStart(hashes []string, value bool) error {
	data := url.Values{}
	data.Set("hashes", joinHashes(hashes))
	data.Set("value", fmt.Sprintf("%t", value))

	_, err := c.doPostValues("/api/v2/torrents/setForceStart", data)
	if err != nil {
		return fmt.Errorf("SetForceStart error: %v", err)
	}
	return nil
}


// TorrentsInfoParams holds the optional parameters for the TorrentsInfo method
type TorrentsInfoParams struct {
	Filter   string
	Category string
	Tag      string
	Sort     string
	Reverse  bool
	Limit    int
	Offset   int
	Hashes   []string
}

// TorrentsInfo retrieves a list of all torrents from the qBittorrent server
func (c *Client) TorrentsInfo(params ...*TorrentsInfoParams) ([]TorrentInfo, error) {
	var query url.Values
	if len(params) > 0 && params[0] != nil {
		query = url.Values{}
		if params[0].Filter != "" {
			query.Set("filter", params[0].Filter)
		}
		if params[0].Category != "" {
			query.Set("category", params[0].Category)
		}
		if params[0].Tag != "" {
			query.Set("tag", params[0].Tag)
		}
		if params[0].Sort != "" {
			query.Set("sort", params[0].Sort)
		}
		if params[0].Reverse {
			query.Set("reverse", "true")
		}
		if params[0].Limit > 0 {
			query.Set("limit", strconv.Itoa(params[0].Limit))
		}
		if params[0].Offset != 0 {
			query.Set("offset", strconv.Itoa(params[0].Offset))
		}
		if len(params[0].Hashes) > 0 {
			query.Set("hashes", strings.Join(params[0].Hashes, "|"))
		}
	}

	respData, err := c.doGet("/api/v2/torrents/info", query)
	if err != nil {
		return nil, err
	}

	var torrents []TorrentInfo
	if err := json.Unmarshal(respData, &torrents); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return torrents, nil
}

// TorrentsTrackers retrieves the tracker info for a given torrent hash
func (c *Client) TorrentsTrackers(hash string) ([]TrackerInfo, error) {
	params := url.Values{}
	params.Set("hash", hash)

	respData, err := c.doGet("/api/v2/torrents/trackers", params)
	if err != nil {
		return nil, fmt.Errorf("TorrentsTrackers error: %v", err)
	}

	var trackers []TrackerInfo
	if err := json.Unmarshal(respData, &trackers); err != nil {
		return nil, fmt.Errorf("failed to decode trackers response: %v", err)
	}

	return trackers, nil
}

// TorrentsProperties retrieves the generic properties for a given torrent hash.
func (c *Client) TorrentsProperties(hash string) (*TorrentsProperties, error) {
	params := url.Values{}
	params.Set("hash", hash)

	respData, err := c.doGet("/api/v2/torrents/properties", params)
	if err != nil {
		return nil, fmt.Errorf("TorrentsProperties error: %v", err)
	}
	if len(respData) == 0 {
		return nil, fmt.Errorf("TorrentsProperties error: empty response")
	}

	var props TorrentsProperties
	if err := json.Unmarshal(respData, &props); err != nil {
		return nil, fmt.Errorf("failed to decode properties response: %v", err)
	}

	return &props, nil
}

// TorrentsAddTags adds tags to the specified torrents
func (c *Client) TorrentsAddTags(hashes []string, tags string) error {
	data := url.Values{}
	data.Set("hashes", joinHashes(hashes))
	data.Set("tags", tags)

	_, err := c.doPostValues("/api/v2/torrents/addTags", data)
	if err != nil {
		return fmt.Errorf("AddTags error: %v", err)
	}
	return nil
}

// TorrentsRemoveTags removes tags from the specified torrents
func (c *Client) TorrentsRemoveTags(hashes []string, tags string) error {
	data := url.Values{}
	data.Set("hashes", joinHashes(hashes))
	data.Set("tags", tags)

	_, err := c.doPostValues("/api/v2/torrents/removeTags", data)
	if err != nil {
		return fmt.Errorf("RemoveTags error: %v", err)
	}
	return nil
}

// TorrentsGetTags retrieves the tags for the given torrent hashes
func (c *Client) TorrentsGetTags(hashes []string) ([]string, error) {
	params := &TorrentsInfoParams{
		Hashes: hashes,
	}

	torrents, err := c.TorrentsInfo(params)
	if err != nil {
		return nil, fmt.Errorf("TorrentsGetTags error: %v", err)
	}

	tagSet := make(map[string]struct{})
	for _, torrent := range torrents {
		for _, tag := range torrent.Tags {
			tagSet[tag] = struct{}{}
		}
	}

	var tags []string
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	return tags, nil
}

// TorrentsGetAllTags retrieves all tags from qBittorrent
func (c *Client) TorrentsGetAllTags() ([]string, error) {
	respData, err := c.doGet("/api/v2/torrents/tags", nil)
	if err != nil {
		return nil, fmt.Errorf("GetAllTags error: %v", err)
	}

	var tags []string
	if err := json.Unmarshal(respData, &tags); err != nil {
		return nil, fmt.Errorf("failed to decode tags response: %v", err)
	}

	return tags, nil
}

// TorrentsCreateTags creates new tags in qBittorrent
func (c *Client) TorrentsCreateTags(tags string) error {
	data := url.Values{}
	data.Set("tags", tags)

	_, err := c.doPostValues("/api/v2/torrents/createTags", data)
	if err != nil {
		return fmt.Errorf("CreateTags error: %v", err)
	}
	return nil
}

// TorrentsDeleteTags deletes tags from qBittorrent
func (c *Client) TorrentsDeleteTags(tags string) error {
	data := url.Values{}
	data.Set("tags", tags)

	_, err := c.doPostValues("/api/v2/torrents/deleteTags", data)
	if err != nil {
		return fmt.Errorf("DeleteTags error: %v", err)
	}
	return nil
}

// doPostResponse POSTs to qBittorrent and returns the HTTP response
func (c *Client) doPostResponse(endpoint string, body io.Reader, contentType string) (*http.Response, error) {
	return c.doRequest("POST", endpoint, body, contentType)
}

// doPost makes POSTs to qBittorrent and returns the response body
func (c *Client) doPost(endpoint string, body io.Reader, contentType string) ([]byte, error) {
	resp, err := c.doPostResponse(endpoint, body, contentType)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("POST error (%d): %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}

// doPostValues POSTs to qBittorrent with url.Values and returns the response body
func (c *Client) doPostValues(endpoint string, data url.Values) ([]byte, error) {
	return c.doPost(endpoint, strings.NewReader(data.Encode()), "application/x-www-form-urlencoded")
}

// doGet is a helper method for making GET requests to the qBittorrent API with query parameters
func (c *Client) doGet(endpoint string, query url.Values) ([]byte, error) {
	resp, err := c.doRequest("GET", endpoint, nil, "", withQuery(query))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response code: %d, response: %s", resp.StatusCode, string(respBody))
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %v", err)
	}
	return responseData, nil
}

// doRequest is a helper function to handle HTTP requests with optional query parameters
func (c *Client) doRequest(method, endpoint string, body io.Reader, contentType string, opts ...func(*http.Request) error) (*http.Response, error) {
	apiURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %v", err)
	}

	apiURL.Path = strings.TrimSuffix(apiURL.Path, "/") + endpoint

	// Store body in buffer if it's not nil so we can retry the request
	var bodyBuffer []byte
	if body != nil {
		bodyBuffer, err = io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %v", err)
		}
	}

	makeRequest := func() (*http.Request, error) {
		var bodyReader io.Reader
		if bodyBuffer != nil {
			bodyReader = bytes.NewReader(bodyBuffer)
		}
		req, err := http.NewRequest(method, apiURL.String(), bodyReader)
		if err != nil {
			return nil, fmt.Errorf("NewRequest error: %v", err)
		}

		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}

		c.mu.RLock()
		if c.sid != "" {
			req.AddCookie(&http.Cookie{Name: "SID", Value: c.sid})
		}
		c.mu.RUnlock()

		// Apply any optional request modifiers
		for _, opt := range opts {
			if err := opt(req); err != nil {
				return nil, err
			}
		}
		return req, nil
	}

	// Make initial request
	req, err := makeRequest()
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// If we get a 403 Forbidden, try to re-authenticate once and retry the request
	if resp.StatusCode == http.StatusForbidden {
		resp.Body.Close() // Close the first response

		if err := c.AuthLogin(); err != nil {
			return nil, fmt.Errorf("re-authentication failed: %v", err)
		}

		// Retry the original request with the new SID
		req, err := makeRequest()
		if err != nil {
			return nil, err
		}

		return c.client.Do(req)
	}

	return resp, nil
}

// withQuery returns a request modifier that adds query parameters
func withQuery(query url.Values) func(*http.Request) error {
	return func(req *http.Request) error {
		req.URL.RawQuery = query.Encode()
		return nil
	}
}

// SyncMainData retrieves the main data from qBittorrent for synchronization
func (c *Client) SyncMainData(rid int) (*MainData, error) {
	params := url.Values{}
	params.Set("rid", strconv.Itoa(rid))

	resp, err := c.doGet("/api/v2/sync/maindata", params)
	if err != nil {
		return nil, err
	}

	var result MainData
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// SyncTorrentPeers retrieves peer information for a specific torrent for synchronization
func (c *Client) SyncTorrentPeers(hash string, rid int) (*TorrentPeers, error) {
	params := url.Values{}
	params.Set("rid", strconv.Itoa(rid))
	params.Set("hash", hash)

	resp, err := c.doGet("/api/v2/sync/torrentPeers", params)
	if err != nil {
		return nil, err
	}

	var result TorrentPeers
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// --- Auth ---

// AuthLogout logs out from the qBittorrent Web API
func (c *Client) AuthLogout() error {
	_, err := c.doPostValues("/api/v2/auth/logout", url.Values{})
	if err != nil {
		return fmt.Errorf("AuthLogout error: %v", err)
	}
	c.mu.Lock()
	c.sid = ""
	c.mu.Unlock()
	return nil
}

// --- App ---

// AppVersion returns the qBittorrent application version
func (c *Client) AppVersion() (string, error) {
	data, err := c.doGet("/api/v2/app/version", nil)
	if err != nil {
		return "", fmt.Errorf("AppVersion error: %v", err)
	}
	return string(data), nil
}

// AppPreferences returns the qBittorrent application preferences
func (c *Client) AppPreferences() (map[string]any, error) {
	data, err := c.doGet("/api/v2/app/preferences", nil)
	if err != nil {
		return nil, fmt.Errorf("AppPreferences error: %v", err)
	}
	var prefs map[string]any
	if err := json.Unmarshal(data, &prefs); err != nil {
		return nil, fmt.Errorf("failed to decode preferences: %v", err)
	}
	return prefs, nil
}

// SetAppPreferences sets qBittorrent application preferences
func (c *Client) SetAppPreferences(prefs map[string]any) error {
	prefsJSON, err := json.Marshal(prefs)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %v", err)
	}
	data := url.Values{}
	data.Set("json", string(prefsJSON))
	_, err = c.doPostValues("/api/v2/app/setPreferences", data)
	if err != nil {
		return fmt.Errorf("SetAppPreferences error: %v", err)
	}
	return nil
}

// AppDefaultSavePath returns the default save path
func (c *Client) AppDefaultSavePath() (string, error) {
	data, err := c.doGet("/api/v2/app/defaultSavePath", nil)
	if err != nil {
		return "", fmt.Errorf("AppDefaultSavePath error: %v", err)
	}
	return string(data), nil
}

// --- Transfer ---

// TransferInfo holds global transfer statistics
type TransferInfo struct {
	DLInfoSpeed      int64  `json:"dl_info_speed"`
	DLInfoData       int64  `json:"dl_info_data"`
	UPInfoSpeed      int64  `json:"up_info_speed"`
	UPInfoData       int64  `json:"up_info_data"`
	DLRateLimit      int64  `json:"dl_rate_limit"`
	UPRateLimit      int64  `json:"up_rate_limit"`
	DHTNodes         int    `json:"dht_nodes"`
	ConnectionStatus string `json:"connection_status"`
}

// TransferGetInfo retrieves global transfer info
func (c *Client) TransferGetInfo() (*TransferInfo, error) {
	data, err := c.doGet("/api/v2/transfer/info", nil)
	if err != nil {
		return nil, fmt.Errorf("TransferInfo error: %v", err)
	}
	var info TransferInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("failed to decode transfer info: %v", err)
	}
	return &info, nil
}

// --- Torrent actions (critical) ---

// TorrentsPause pauses the specified torrents
func (c *Client) TorrentsPause(hashes []string) error {
	data := url.Values{}
	data.Set("hashes", joinHashes(hashes))
	_, err := c.doPostValues("/api/v2/torrents/pause", data)
	if err != nil {
		return fmt.Errorf("TorrentsPause error: %v", err)
	}
	return nil
}

// TorrentsResume resumes the specified torrents
func (c *Client) TorrentsResume(hashes []string) error {
	data := url.Values{}
	data.Set("hashes", joinHashes(hashes))
	_, err := c.doPostValues("/api/v2/torrents/resume", data)
	if err != nil {
		return fmt.Errorf("TorrentsResume error: %v", err)
	}
	return nil
}

// TorrentsSetLocation sets the save location for the specified torrents
func (c *Client) TorrentsSetLocation(hashes []string, location string) error {
	data := url.Values{}
	data.Set("hashes", joinHashes(hashes))
	data.Set("location", location)
	_, err := c.doPostValues("/api/v2/torrents/setLocation", data)
	if err != nil {
		return fmt.Errorf("TorrentsSetLocation error: %v", err)
	}
	return nil
}

// TorrentsRecheck rechecks the specified torrents
func (c *Client) TorrentsRecheck(hashes []string) error {
	data := url.Values{}
	data.Set("hashes", joinHashes(hashes))
	_, err := c.doPostValues("/api/v2/torrents/recheck", data)
	if err != nil {
		return fmt.Errorf("TorrentsRecheck error: %v", err)
	}
	return nil
}

// TorrentsReannounce reannounces the specified torrents
func (c *Client) TorrentsReannounce(hashes []string) error {
	data := url.Values{}
	data.Set("hashes", joinHashes(hashes))
	_, err := c.doPostValues("/api/v2/torrents/reannounce", data)
	if err != nil {
		return fmt.Errorf("TorrentsReannounce error: %v", err)
	}
	return nil
}

// TorrentsSetCategory sets the category for the specified torrents
func (c *Client) TorrentsSetCategory(hashes []string, category string) error {
	data := url.Values{}
	data.Set("hashes", joinHashes(hashes))
	data.Set("category", category)
	_, err := c.doPostValues("/api/v2/torrents/setCategory", data)
	if err != nil {
		return fmt.Errorf("TorrentsSetCategory error: %v", err)
	}
	return nil
}

// TorrentsSetAutoTMM enables or disables Automatic Torrent Management for the specified torrents
func (c *Client) TorrentsSetAutoTMM(hashes []string, enable bool) error {
	data := url.Values{}
	data.Set("hashes", joinHashes(hashes))
	data.Set("enable", fmt.Sprintf("%t", enable))
	_, err := c.doPostValues("/api/v2/torrents/setAutoManagement", data)
	if err != nil {
		return fmt.Errorf("TorrentsSetAutoTMM error: %v", err)
	}
	return nil
}

// TorrentsRename renames a torrent
func (c *Client) TorrentsRename(hash string, name string) error {
	data := url.Values{}
	data.Set("hash", hash)
	data.Set("name", name)
	_, err := c.doPostValues("/api/v2/torrents/rename", data)
	if err != nil {
		return fmt.Errorf("TorrentsRename error: %v", err)
	}
	return nil
}

// --- Per-torrent limits ---

// TorrentsSetDownloadLimit sets the download speed limit for the specified torrents (bytes/second)
func (c *Client) TorrentsSetDownloadLimit(hashes []string, limit int) error {
	data := url.Values{}
	data.Set("hashes", joinHashes(hashes))
	data.Set("limit", strconv.Itoa(limit))
	_, err := c.doPostValues("/api/v2/torrents/setDownloadLimit", data)
	if err != nil {
		return fmt.Errorf("TorrentsSetDownloadLimit error: %v", err)
	}
	return nil
}

// TorrentsSetUploadLimit sets the upload speed limit for the specified torrents (bytes/second)
func (c *Client) TorrentsSetUploadLimit(hashes []string, limit int) error {
	data := url.Values{}
	data.Set("hashes", joinHashes(hashes))
	data.Set("limit", strconv.Itoa(limit))
	_, err := c.doPostValues("/api/v2/torrents/setUploadLimit", data)
	if err != nil {
		return fmt.Errorf("TorrentsSetUploadLimit error: %v", err)
	}
	return nil
}

// TorrentsSetShareLimits sets share limits for the specified torrents
func (c *Client) TorrentsSetShareLimits(hashes []string, ratioLimit float64, seedingTimeLimit, inactiveSeedingTimeLimit int) error {
	data := url.Values{}
	data.Set("hashes", joinHashes(hashes))
	data.Set("ratioLimit", fmt.Sprintf("%.2f", ratioLimit))
	data.Set("seedingTimeLimit", strconv.Itoa(seedingTimeLimit))
	data.Set("inactiveSeedingTimeLimit", strconv.Itoa(inactiveSeedingTimeLimit))
	_, err := c.doPostValues("/api/v2/torrents/setShareLimits", data)
	if err != nil {
		return fmt.Errorf("TorrentsSetShareLimits error: %v", err)
	}
	return nil
}

// --- Torrent files ---

// TorrentFile represents a file within a torrent
type TorrentFile struct {
	Index        int     `json:"index"`
	Name         string  `json:"name"`
	Size         int64   `json:"size"`
	Progress     float64 `json:"progress"`
	Priority     int     `json:"priority"`
	IsSeed       bool    `json:"is_seed"`
	PieceRange   []int   `json:"piece_range"`
	Availability float64 `json:"availability"`
}

// TorrentsFiles retrieves the files for a given torrent hash
func (c *Client) TorrentsFiles(hash string) ([]TorrentFile, error) {
	params := url.Values{}
	params.Set("hash", hash)
	data, err := c.doGet("/api/v2/torrents/files", params)
	if err != nil {
		return nil, fmt.Errorf("TorrentsFiles error: %v", err)
	}
	var files []TorrentFile
	if err := json.Unmarshal(data, &files); err != nil {
		return nil, fmt.Errorf("failed to decode files response: %v", err)
	}
	return files, nil
}

// TorrentsFilePrio sets the priority for specific files within a torrent
func (c *Client) TorrentsFilePrio(hash string, ids []int, priority int) error {
	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = strconv.Itoa(id)
	}
	data := url.Values{}
	data.Set("hash", hash)
	data.Set("id", strings.Join(idStrs, "|"))
	data.Set("priority", strconv.Itoa(priority))
	_, err := c.doPostValues("/api/v2/torrents/filePrio", data)
	if err != nil {
		return fmt.Errorf("TorrentsFilePrio error: %v", err)
	}
	return nil
}

// --- Category CRUD ---

// TorrentsCategories retrieves all categories
func (c *Client) TorrentsCategories() (map[string]Category, error) {
	data, err := c.doGet("/api/v2/torrents/categories", nil)
	if err != nil {
		return nil, fmt.Errorf("TorrentsCategories error: %v", err)
	}
	var categories map[string]Category
	if err := json.Unmarshal(data, &categories); err != nil {
		return nil, fmt.Errorf("failed to decode categories response: %v", err)
	}
	return categories, nil
}

// TorrentsCreateCategory creates a new category
func (c *Client) TorrentsCreateCategory(name, savePath string) error {
	data := url.Values{}
	data.Set("category", name)
	data.Set("savePath", savePath)
	_, err := c.doPostValues("/api/v2/torrents/createCategory", data)
	if err != nil {
		return fmt.Errorf("TorrentsCreateCategory error: %v", err)
	}
	return nil
}

// TorrentsEditCategory edits an existing category
func (c *Client) TorrentsEditCategory(name, savePath string) error {
	data := url.Values{}
	data.Set("category", name)
	data.Set("savePath", savePath)
	_, err := c.doPostValues("/api/v2/torrents/editCategory", data)
	if err != nil {
		return fmt.Errorf("TorrentsEditCategory error: %v", err)
	}
	return nil
}

// TorrentsRemoveCategories removes the specified categories
func (c *Client) TorrentsRemoveCategories(categories []string) error {
	data := url.Values{}
	data.Set("categories", strings.Join(categories, "\n"))
	_, err := c.doPostValues("/api/v2/torrents/removeCategories", data)
	if err != nil {
		return fmt.Errorf("TorrentsRemoveCategories error: %v", err)
	}
	return nil
}
