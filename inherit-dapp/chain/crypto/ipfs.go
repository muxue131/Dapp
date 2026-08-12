package crypto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultIPFSAPI      = "http://localhost:5001/api/v0"
	MaxDownloadSize     = 50 * 1024 * 1024 // 50MB max download size
)

// IPFSClient provides methods to interact with an IPFS node
type IPFSClient struct {
	APIURL    string
	HTTPClient *http.Client
}

// NewIPFSClient creates a new IPFS client
func NewIPFSClient(apiURL string) *IPFSClient {
	if apiURL == "" {
		apiURL = DefaultIPFSAPI
	}
	return &IPFSClient{
		APIURL: apiURL,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// UploadResult represents the result of an IPFS upload
type UploadResult struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

// Upload uploads data to IPFS and returns the CID
func (c *IPFSClient) Upload(data []byte, filename string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := part.Write(data); err != nil {
		return "", fmt.Errorf("failed to write data: %w", err)
	}

	writer.Close()

	url := fmt.Sprintf("%s/add?pin=true", c.APIURL)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload to IPFS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("IPFS upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var result UploadResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode IPFS response: %w", err)
	}

	return result.Hash, nil
}

// UploadFile uploads a file to IPFS
func (c *IPFSClient) UploadFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return c.Upload(data, filepath.Base(filePath))
}

// Download downloads data from IPFS by CID
func (c *IPFSClient) Download(cid string) ([]byte, error) {
	apiURL := fmt.Sprintf("%s/cat?arg=%s", c.APIURL, url.QueryEscape(cid))
	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download from IPFS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IPFS download failed with status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, MaxDownloadSize))
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return data, nil
}

// Pin pins a CID on the IPFS node
func (c *IPFSClient) Pin(cid string) error {
	apiURL := fmt.Sprintf("%s/pin/add?arg=%s", c.APIURL, url.QueryEscape(cid))
	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to pin on IPFS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("IPFS pin failed with status %d", resp.StatusCode)
	}

	return nil
}

// EncryptAndUpload encrypts data with a key and uploads to IPFS
func (c *IPFSClient) EncryptAndUpload(data []byte, key []byte, filename string) (string, *EncryptedData, error) {
	encrypted, err := Encrypt(data, key)
	if err != nil {
		return "", nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	encryptedBytes, err := json.Marshal(encrypted)
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal encrypted data: %w", err)
	}

	cid, err := c.Upload(encryptedBytes, filename)
	if err != nil {
		return "", nil, fmt.Errorf("failed to upload to IPFS: %w", err)
	}

	return cid, encrypted, nil
}

// DownloadAndDecrypt downloads encrypted data from IPFS and decrypts it
func (c *IPFSClient) DownloadAndDecrypt(cid string, key []byte) ([]byte, error) {
	data, err := c.Download(cid)
	if err != nil {
		return nil, fmt.Errorf("failed to download from IPFS: %w", err)
	}

	var encrypted EncryptedData
	if err := json.Unmarshal(data, &encrypted); err != nil {
		return nil, fmt.Errorf("failed to unmarshal encrypted data: %w", err)
	}

	decrypted, err := Decrypt(&encrypted, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return decrypted, nil
}
