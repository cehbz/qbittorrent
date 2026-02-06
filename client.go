package qbittorrent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type InfoHash string

// Client is used to interact with the qBittorrent API
type Client struct {
	username string
	password string
	client   *http.Client
	baseURL  string
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
	DLLimit            int64    `json:"dl_limit"`
	DLSpeed            int64    `json:"dlspeed"`
	Downloaded         int64    `json:"downloaded"`
	DownloadedSession  int64    `json:"downloaded_session"`
	ETA                int64    `json:"eta"`
	FirstLastPiecePrio bool     `json:"f_l_piece_prio"`
	ForceStart         bool     `json:"force_start"`
	Hash               InfoHash `json:"hash"`
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
	IsPrivate              bool     `json:"isPrivate"`
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
		t.Tags = strings.Split(aux.RawTags, ",")
	}
	return nil
}

// TrackerInfo represents a tracker info for a torrent
type TrackerInfo struct {
	URL      string `json:"url"`
	Status   int    `json:"status"`
	Tier     int    `json:"tier"`
	NumPeers int    `json:"num_peers"`
	Msg      string `json:"msg"`
}

type Category map[string]interface{} // no idea what this should be, category=CategoryName&savePath=/path/to/dir

// fields might be missing, in which case we need to switch to pointers and allow "omitempty"
// https://github.com/qbittorrent/qBittorrent/blob/master/src/base/json_api.cpp#L101
// MainData is the data returned by the /api/v2/sync/maindata endpoint
type MainData struct {
	Categories        map[string]Category    `json:"categories"`
	CategoriesRemoved []Category             `json:"categories_removed"`
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
// If httpClient is nil, a new client with cookie jar is created.
func NewClient(username, password, baseURL string, httpClient ...*http.Client) (*Client, error) {
	var client *http.Client

	// Use the provided http.Client if given, otherwise create one with cookie jar
	if len(httpClient) > 0 && httpClient[0] != nil {
		client = httpClient[0]
	} else {
		// Create a cookie jar for automatic cookie management
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create cookie jar: %w", err)
		}
		client = &http.Client{
			Jar: jar,
		}
	}

	// Ensure the client has a cookie jar for session management
	if client.Jar == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create cookie jar: %w", err)
		}
		client.Jar = jar
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
			return nil, fmt.Errorf("AuthLogin error: %w", err)
		}
	}

	return qbClient, nil
}

// AuthLoginCtx logs in to the qBittorrent Web API with context support
func (c *Client) AuthLoginCtx(ctx context.Context) error {
	data := url.Values{}
	data.Set("username", c.username)
	data.Set("password", c.password)

	resp, err := c.doPostResponseCtx(ctx, "/api/v2/auth/login", strings.NewReader(data.Encode()), "application/x-www-form-urlencoded")
	if err != nil {
		return wrapAPIError("AuthLogin", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("AuthLogin error (%d): %s", resp.StatusCode, string(respBody))
	}

	// Cookie jar automatically stores the SID cookie from the response
	return nil
}

// AuthLogin logs in to the qBittorrent Web API
// Deprecated: Use AuthLoginCtx for context support
func (c *Client) AuthLogin() error {
	return c.AuthLoginCtx(context.Background())
}

// TorrentsExportCtx retrieves the .torrent file for a given torrent hash with context support
func (c *Client) TorrentsExportCtx(ctx context.Context, hash string) ([]byte, error) {
	params := url.Values{}
	params.Set("hash", hash)

	return c.doPostValuesCtx(ctx, "/api/v2/torrents/export", params)
}

// TorrentsExport retrieves the .torrent file for a given torrent hash
// Deprecated: Use TorrentsExportCtx for context support
func (c *Client) TorrentsExport(hash string) ([]byte, error) {
	return c.TorrentsExportCtx(context.Background(), hash)
}

// TorrentsAddParams holds the parameters for adding a torrent
type TorrentsAddParams struct {
	Torrents    [][]byte // Raw torrent data
	URLs        []string // Magnet links or URLs
	SavePath    string   // Download folder
	Cookie      string   // Cookie sent to download the .torrent file
	Category    string   // Category for the torrent
	Tags        string   // Tags for the torrent, comma separated
	SkipCheck   bool     // Skip hash checking
	Paused      bool     // Add torrents in the paused state
	RootFolder  *bool    // Create the root folder (default: true)
	Rename      string   // Rename torrent
	UpLimit     int      // Set torrent upload speed limit (bytes/second)
	DlLimit     int      // Set torrent download speed limit (bytes/second)
	RatioLimit  float64  // Set torrent share ratio limit
	SeedingTime int      // Set torrent seeding time limit (minutes)
	AutoTMM     bool     // Whether Automatic Torrent Management should be used
	Sequential  bool     // Enable sequential download
	FirstLast   bool     // Prioritize download first last piece
}

// buildTorrentAddForm builds a multipart form for adding torrents (DRY helper)
func buildTorrentAddForm(params *TorrentsAddParams) (*bytes.Buffer, string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add torrent files
	for i, torrentData := range params.Torrents {
		part, err := writer.CreateFormFile("torrents", fmt.Sprintf("torrent%d.torrent", i))
		if err != nil {
			return nil, "", fmt.Errorf("CreateFormFile error: %w", err)
		}
		if _, err := io.Copy(part, bytes.NewReader(torrentData)); err != nil {
			return nil, "", fmt.Errorf("io.Copy error: %w", err)
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
	if params.Paused {
		_ = writer.WriteField("paused", "true")
	}
	if params.RootFolder != nil {
		_ = writer.WriteField("root_folder", fmt.Sprintf("%t", *params.RootFolder))
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
	if params.AutoTMM {
		_ = writer.WriteField("autoTMM", "true")
	}
	if params.Sequential {
		_ = writer.WriteField("sequentialDownload", "true")
	}
	if params.FirstLast {
		_ = writer.WriteField("firstLastPiecePrio", "true")
	}

	writer.Close()
	return &body, writer.FormDataContentType(), nil
}

// TorrentsAddParamsCtx adds a torrent using a more flexible parameter structure with context support
func (c *Client) TorrentsAddParamsCtx(ctx context.Context, params *TorrentsAddParams) error {
	body, contentType, err := buildTorrentAddForm(params)
	if err != nil {
		return err
	}

	_, err = c.doPostCtx(ctx, "/api/v2/torrents/add", body, contentType)
	if err != nil {
		return wrapAPIError("TorrentsAdd", err)
	}
	return nil
}

// TorrentsAddParams adds a torrent using a more flexible parameter structure
// Deprecated: Use TorrentsAddParamsCtx for context support
func (c *Client) TorrentsAddParams(params *TorrentsAddParams) error {
	return c.TorrentsAddParamsCtx(context.Background(), params)
}

// TorrentsAddCtx adds a torrent to qBittorrent via Web API using multipart/form-data with context support
func (c *Client) TorrentsAddCtx(ctx context.Context, torrentFile string, fileData []byte) error {
	params := &TorrentsAddParams{
		Torrents:  [][]byte{fileData},
		SkipCheck: true,
		Paused:    false,
		AutoTMM:   false,
	}
	return c.TorrentsAddParamsCtx(ctx, params)
}

// TorrentsAdd adds a torrent to qBittorrent via Web API using multipart/form-data
// Deprecated: Use TorrentsAddCtx for context support
func (c *Client) TorrentsAdd(torrentFile string, fileData []byte) error {
	return c.TorrentsAddCtx(context.Background(), torrentFile, fileData)
}

// TorrentsDeleteCtx deletes a torrent from qBittorrent by its hash with context support
func (c *Client) TorrentsDeleteCtx(ctx context.Context, infohash string) error {
	data := url.Values{}
	data.Set("hashes", infohash)
	data.Set("deleteFiles", "true")

	_, err := c.doPostValuesCtx(ctx, "/api/v2/torrents/delete", data)
	if err != nil {
		return wrapAPIError("TorrentsDelete", err)
	}
	return nil
}

// TorrentsDelete deletes a torrent from qBittorrent by its hash
// Deprecated: Use TorrentsDeleteCtx for context support
func (c *Client) TorrentsDelete(infohash string) error {
	return c.TorrentsDeleteCtx(context.Background(), infohash)
}

// SetForceStartCtx enables force start for the torrent with context support
func (c *Client) SetForceStartCtx(ctx context.Context, hash string, value bool) error {
	data := url.Values{}
	data.Set("hashes", hash)
	data.Set("value", fmt.Sprintf("%t", value))

	_, err := c.doPostValuesCtx(ctx, "/api/v2/torrents/setForceStart", data)
	if err != nil {
		return wrapAPIError("SetForceStart", err)
	}
	return nil
}

// SetForceStart enables force start for the torrent
// Deprecated: Use SetForceStartCtx for context support
func (c *Client) SetForceStart(hash string, value bool) error {
	return c.SetForceStartCtx(context.Background(), hash, value)
}

// TorrentsDownloadCtx retrieves the torrent file by its hash from the qBittorrent server with context support
func (c *Client) TorrentsDownloadCtx(ctx context.Context, infohash string) ([]byte, error) {
	return c.doGetCtx(ctx, "/api/v2/torrents/file", url.Values{"hashes": {infohash}})
}

// TorrentsDownload retrieves the torrent file by its hash from the qBittorrent server
// Deprecated: Use TorrentsDownloadCtx for context support
func (c *Client) TorrentsDownload(infohash string) ([]byte, error) {
	return c.TorrentsDownloadCtx(context.Background(), infohash)
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

// buildTorrentsInfoQuery builds query parameters for TorrentsInfo (DRY helper)
func buildTorrentsInfoQuery(params *TorrentsInfoParams) url.Values {
	query := url.Values{}
	if params.Filter != "" {
		query.Set("filter", params.Filter)
	}
	if params.Category != "" {
		query.Set("category", params.Category)
	}
	if params.Tag != "" {
		query.Set("tag", params.Tag)
	}
	if params.Sort != "" {
		query.Set("sort", params.Sort)
	}
	if params.Reverse {
		query.Set("reverse", "true")
	}
	if params.Limit > 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset != 0 {
		query.Set("offset", strconv.Itoa(params.Offset))
	}
	if len(params.Hashes) > 0 {
		query.Set("hashes", strings.Join(params.Hashes, "|"))
	}
	return query
}

// TorrentsInfoCtx retrieves a list of all torrents from the qBittorrent server with context support
func (c *Client) TorrentsInfoCtx(ctx context.Context, params ...*TorrentsInfoParams) ([]TorrentInfo, error) {
	var query url.Values
	if len(params) > 0 && params[0] != nil {
		query = buildTorrentsInfoQuery(params[0])
	}

	respData, err := c.doGetCtx(ctx, "/api/v2/torrents/info", query)
	if err != nil {
		return nil, err
	}

	var torrents []TorrentInfo
	if err := json.Unmarshal(respData, &torrents); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return torrents, nil
}

// TorrentsInfo retrieves a list of all torrents from the qBittorrent server
// Deprecated: Use TorrentsInfoCtx for context support
func (c *Client) TorrentsInfo(params ...*TorrentsInfoParams) ([]TorrentInfo, error) {
	return c.TorrentsInfoCtx(context.Background(), params...)
}

// TorrentsTrackersCtx retrieves the tracker info for a given torrent hash with context support
func (c *Client) TorrentsTrackersCtx(ctx context.Context, hash string) ([]TrackerInfo, error) {
	params := url.Values{}
	params.Set("hash", hash)

	respData, err := c.doGetCtx(ctx, "/api/v2/torrents/trackers", params)
	if err != nil {
		return nil, wrapAPIError("TorrentsTrackers", err)
	}

	var trackers []TrackerInfo
	if err := json.Unmarshal(respData, &trackers); err != nil {
		return nil, fmt.Errorf("failed to decode trackers response: %w", err)
	}

	return trackers, nil
}

// TorrentsTrackers retrieves the tracker info for a given torrent hash
// Deprecated: Use TorrentsTrackersCtx for context support
func (c *Client) TorrentsTrackers(hash string) ([]TrackerInfo, error) {
	return c.TorrentsTrackersCtx(context.Background(), hash)
}

// TorrentsPropertiesCtx retrieves the generic properties for a given torrent hash with context support
func (c *Client) TorrentsPropertiesCtx(ctx context.Context, hash string) (*TorrentsProperties, error) {
	params := url.Values{}
	params.Set("hash", hash)

	respData, err := c.doGetCtx(ctx, "/api/v2/torrents/properties", params)
	if err != nil {
		return nil, wrapAPIError("TorrentsProperties", err)
	}
	if len(respData) == 0 {
		return nil, fmt.Errorf("TorrentsProperties error: empty response")
	}

	var props TorrentsProperties
	if err := json.Unmarshal(respData, &props); err != nil {
		return nil, fmt.Errorf("failed to decode properties response: %w", err)
	}

	return &props, nil
}

// TorrentsProperties retrieves the generic properties for a given torrent hash.
// Deprecated: Use TorrentsPropertiesCtx for context support
func (c *Client) TorrentsProperties(hash string) (*TorrentsProperties, error) {
	return c.TorrentsPropertiesCtx(context.Background(), hash)
}

// Helper to perform tag operations (DRY)
func (c *Client) doTagOperationCtx(ctx context.Context, operation, endpoint, hashes, tags string) error {
	data := url.Values{}
	if hashes != "" {
		data.Set("hashes", hashes)
	}
	data.Set("tags", tags)

	_, err := c.doPostValuesCtx(ctx, endpoint, data)
	if err != nil {
		return wrapAPIError(operation, err)
	}
	return nil
}

// TorrentsAddTagsCtx adds tags to the specified torrents with context support
func (c *Client) TorrentsAddTagsCtx(ctx context.Context, hashes, tags string) error {
	return c.doTagOperationCtx(ctx, "AddTags", "/api/v2/torrents/addTags", hashes, tags)
}

// TorrentsAddTags adds tags to the specified torrents
// Deprecated: Use TorrentsAddTagsCtx for context support
func (c *Client) TorrentsAddTags(hashes, tags string) error {
	return c.TorrentsAddTagsCtx(context.Background(), hashes, tags)
}

// TorrentsRemoveTagsCtx removes tags from the specified torrents with context support
func (c *Client) TorrentsRemoveTagsCtx(ctx context.Context, hashes, tags string) error {
	return c.doTagOperationCtx(ctx, "RemoveTags", "/api/v2/torrents/removeTags", hashes, tags)
}

// TorrentsRemoveTags removes tags from the specified torrents
// Deprecated: Use TorrentsRemoveTagsCtx for context support
func (c *Client) TorrentsRemoveTags(hashes, tags string) error {
	return c.TorrentsRemoveTagsCtx(context.Background(), hashes, tags)
}

// TorrentsGetTagsCtx retrieves the tags for the given torrent hashes with context support
func (c *Client) TorrentsGetTagsCtx(ctx context.Context, hashes string) ([]string, error) {
	params := &TorrentsInfoParams{
		Hashes: []string{hashes},
	}

	torrents, err := c.TorrentsInfoCtx(ctx, params)
	if err != nil {
		return nil, wrapAPIError("TorrentsGetTags", err)
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

// TorrentsGetTags retrieves the tags for the given torrent hashes
// Deprecated: Use TorrentsGetTagsCtx for context support
func (c *Client) TorrentsGetTags(hashes string) ([]string, error) {
	return c.TorrentsGetTagsCtx(context.Background(), hashes)
}

// TorrentsGetAllTagsCtx retrieves all tags from qBittorrent with context support
func (c *Client) TorrentsGetAllTagsCtx(ctx context.Context) ([]string, error) {
	respData, err := c.doGetCtx(ctx, "/api/v2/torrents/tags", nil)
	if err != nil {
		return nil, wrapAPIError("GetAllTags", err)
	}

	var tags []string
	if err := json.Unmarshal(respData, &tags); err != nil {
		return nil, fmt.Errorf("failed to decode tags response: %w", err)
	}

	return tags, nil
}

// TorrentsGetAllTags retrieves all tags from qBittorrent
// Deprecated: Use TorrentsGetAllTagsCtx for context support
func (c *Client) TorrentsGetAllTags() ([]string, error) {
	return c.TorrentsGetAllTagsCtx(context.Background())
}

// TorrentsCreateTagsCtx creates new tags in qBittorrent with context support
func (c *Client) TorrentsCreateTagsCtx(ctx context.Context, tags string) error {
	return c.doTagOperationCtx(ctx, "CreateTags", "/api/v2/torrents/createTags", "", tags)
}

// TorrentsCreateTags creates new tags in qBittorrent
// Deprecated: Use TorrentsCreateTagsCtx for context support
func (c *Client) TorrentsCreateTags(tags string) error {
	return c.TorrentsCreateTagsCtx(context.Background(), tags)
}

// TorrentsDeleteTagsCtx deletes tags from qBittorrent with context support
func (c *Client) TorrentsDeleteTagsCtx(ctx context.Context, tags string) error {
	return c.doTagOperationCtx(ctx, "DeleteTags", "/api/v2/torrents/deleteTags", "", tags)
}

// TorrentsDeleteTags deletes tags from qBittorrent
// Deprecated: Use TorrentsDeleteTagsCtx for context support
func (c *Client) TorrentsDeleteTags(tags string) error {
	return c.TorrentsDeleteTagsCtx(context.Background(), tags)
}

// Helper function to wrap API errors consistently
func wrapAPIError(operation string, err error) error {
	return fmt.Errorf("%s error: %w", operation, err)
}

// doPostResponseCtx POSTs to qBittorrent and returns the HTTP response
func (c *Client) doPostResponseCtx(ctx context.Context, endpoint string, body io.Reader, contentType string) (*http.Response, error) {
	return c.doRequestCtx(ctx, "POST", endpoint, body, contentType)
}

// doPostResponse POSTs to qBittorrent and returns the HTTP response
// Deprecated: Use doPostResponseCtx for context support
func (c *Client) doPostResponse(endpoint string, body io.Reader, contentType string) (*http.Response, error) {
	return c.doPostResponseCtx(context.Background(), endpoint, body, contentType)
}

// doPostCtx makes POSTs to qBittorrent and returns the response body
func (c *Client) doPostCtx(ctx context.Context, endpoint string, body io.Reader, contentType string) ([]byte, error) {
	resp, err := c.doPostResponseCtx(ctx, endpoint, body, contentType)
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

// doPost makes POSTs to qBittorrent and returns the response body
// Deprecated: Use doPostCtx for context support
func (c *Client) doPost(endpoint string, body io.Reader, contentType string) ([]byte, error) {
	return c.doPostCtx(context.Background(), endpoint, body, contentType)
}

// doPostValuesCtx POSTs to qBittorrent with url.Values and returns the response body
func (c *Client) doPostValuesCtx(ctx context.Context, endpoint string, data url.Values) ([]byte, error) {
	return c.doPostCtx(ctx, endpoint, strings.NewReader(data.Encode()), "application/x-www-form-urlencoded")
}

// doPostValues POSTs to qBittorrent with url.Values and returns the response body
// Deprecated: Use doPostValuesCtx for context support
func (c *Client) doPostValues(endpoint string, data url.Values) ([]byte, error) {
	return c.doPostValuesCtx(context.Background(), endpoint, data)
}

// doGetCtx is a helper method for making GET requests to the qBittorrent API with query parameters
func (c *Client) doGetCtx(ctx context.Context, endpoint string, query url.Values) ([]byte, error) {
	resp, err := c.doRequestCtx(ctx, "GET", endpoint, nil, "", withQuery(query))
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
		return nil, fmt.Errorf("ReadAll error: %w", err)
	}
	return responseData, nil
}

// doGet is a helper method for making GET requests to the qBittorrent API with query parameters
// Deprecated: Use doGetCtx for context support
func (c *Client) doGet(endpoint string, query url.Values) ([]byte, error) {
	return c.doGetCtx(context.Background(), endpoint, query)
}

// doRequestCtx is a helper function to handle HTTP requests with context and optional query parameters
func (c *Client) doRequestCtx(ctx context.Context, method, endpoint string, body io.Reader, contentType string, opts ...func(*http.Request) error) (*http.Response, error) {
	apiURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	apiURL.Path = strings.TrimSuffix(apiURL.Path, "/") + endpoint

	// Store body in buffer if it's not nil so we can retry the request
	var bodyBuffer []byte
	if body != nil {
		bodyBuffer, err = io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
	}

	makeRequest := func(ctx context.Context) (*http.Request, error) {
		var bodyReader io.Reader
		if bodyBuffer != nil {
			bodyReader = bytes.NewReader(bodyBuffer)
		}
		req, err := http.NewRequestWithContext(ctx, method, apiURL.String(), bodyReader)
		if err != nil {
			return nil, fmt.Errorf("NewRequest error: %w", err)
		}

		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}

		// Cookie jar automatically sends stored cookies with requests

		// Apply any optional request modifiers
		for _, opt := range opts {
			if err := opt(req); err != nil {
				return nil, err
			}
		}
		return req, nil
	}

	// Make initial request
	req, err := makeRequest(ctx)
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

		if err := c.AuthLoginCtx(ctx); err != nil {
			return nil, fmt.Errorf("re-authentication failed: %w", err)
		}

		// Retry the original request with the new SID
		req, err := makeRequest(ctx)
		if err != nil {
			return nil, err
		}

		return c.client.Do(req)
	}

	return resp, nil
}

// doRequest is a helper function to handle HTTP requests with optional query parameters
// Deprecated: Use doRequestCtx for context support
func (c *Client) doRequest(method, endpoint string, body io.Reader, contentType string, opts ...func(*http.Request) error) (*http.Response, error) {
	return c.doRequestCtx(context.Background(), method, endpoint, body, contentType, opts...)
}

// withQuery returns a request modifier that adds query parameters
func withQuery(query url.Values) func(*http.Request) error {
	return func(req *http.Request) error {
		req.URL.RawQuery = query.Encode()
		return nil
	}
}

// SyncMainDataCtx retrieves the main data from qBittorrent for synchronization with context support
func (c *Client) SyncMainDataCtx(ctx context.Context, rid int) (*MainData, error) {
	params := url.Values{}
	params.Set("rid", strconv.Itoa(rid))

	resp, err := c.doGetCtx(ctx, "/api/v2/sync/maindata", params)
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

// SyncMainData retrieves the main data from qBittorrent for synchronization
// Deprecated: Use SyncMainDataCtx for context support
func (c *Client) SyncMainData(rid int) (*MainData, error) {
	return c.SyncMainDataCtx(context.Background(), rid)
}

// SyncTorrentPeersCtx retrieves peer information for a specific torrent for synchronization with context support
func (c *Client) SyncTorrentPeersCtx(ctx context.Context, hash string, rid int) (*TorrentPeers, error) {
	params := url.Values{}
	params.Set("rid", strconv.Itoa(rid))
	params.Set("hash", hash)

	resp, err := c.doGetCtx(ctx, "/api/v2/sync/torrentPeers", params)
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

// SyncTorrentPeers retrieves peer information for a specific torrent for synchronization
// Deprecated: Use SyncTorrentPeersCtx for context support
func (c *Client) SyncTorrentPeers(hash string, rid int) (*TorrentPeers, error) {
	return c.SyncTorrentPeersCtx(context.Background(), hash, rid)
}
