package main

import (
	"net/http"
	"log"
	"encoding/json"
	"runtime"
	"fmt"
	"strings"
	"os"
	"archive/zip"
	"path/filepath"
	"io"
	"io/ioutil"
	"archive/tar"
	"compress/gzip"
	"bytes"
)

func main() {
	log.Print("Checking for newest version...")
	release := getLatestRelease()
	log.Printf("Latest version: %s", release.Version)
	if isUpdateRequired(release.Version) {
		downloadUrl := resolveDownloadUrl(release)
		log.Printf("Download %s...", downloadUrl)
		temporaryFile := downloadToTemporaryLocation(downloadUrl)
		log.Printf("Gooing to extract %s...", temporaryFile)
		binaryFile := extractBinaryFromArchive(release.Version, temporaryFile)
		log.Printf("Binary was successful stored at: %s", binaryFile)
		storeLatestVersionHint(release.Version)
		os.Remove(temporaryFile)
	} else {
		log.Print("No updated required.")
	}
}

func isUpdateRequired(version string) bool {
	versionFilename := resolveTargetVersionFilename()
	content, err := ioutil.ReadFile(versionFilename)
	if err != nil {
		return true
	}
	actual := strings.TrimSpace(string(content))
	return actual != version
}

func storeLatestVersionHint(version string) {
	versionFilename := resolveTargetVersionFilename()
	err := ioutil.WriteFile(versionFilename, []byte(version), 0555)
	if err != nil {
		log.Fatalf("Could not store file version hint. Got: %v", err)
	}
}

func extractBinaryFromArchive(version string, file string) string {
	base := filepath.Base(file)
	if strings.HasPrefix(base, "hugo.zip.") {
		return extractBinaryFileFromZip(version, file)
	}
	if strings.HasPrefix(base, "hugo.gz.") {
		return extractBinaryFromTarGzip(version, file)
	}
	log.Fatalf("Unsupported file type: %s", file)
	return ""
}

func extractBinaryFromTarGzip(version string, file string) string {
	rawReader, err := os.Open(file)
	if err != nil {
		log.Fatalf("Could not read downloaded archive %s. Got: %v", file, err)
	}
	defer rawReader.Close()
	gzipReader, err := gzip.NewReader(rawReader)
	if err != nil {
		log.Fatalf("Could not read downloaded archive %s. Got: %v", file, err)
	}
	defer gzipReader.Close()
	reader := tar.NewReader(gzipReader)
	target := resolveTargetFilename()
	expected := resolveFilenameInArchive(version)
	for {
		candidate, err := reader.Next()
		if err == io.EOF {
			log.Fatalf("Could not find in downloaded file %s the expected file %s.", file, expected)
		}
		if err != nil {
			log.Fatalf("Could not read downloaded archive %s. Got: %v", file, err)
		}
		if candidate.Name == expected {
			extractBinaryFileFromTarPart(reader, target)
			return target
		}
	}
}

func extractBinaryFileFromTarPart(input io.Reader, target string) {
	output, err := os.Create(target)
	if err != nil {
		log.Fatalf("Could not create file %s. Got: %v", target, err)
	}
	defer output.Close()
	io.Copy(output, input)
	if runtime.GOOS != "windows" {
		err = output.Chmod(0755)
		if err != nil {
			log.Fatalf("Could not make %s executable. Got: %v", target, err)
		}
	}
}

func extractBinaryFileFromZip(version string, file string) string {
	reader, err := zip.OpenReader(file)
	if err != nil {
		log.Fatalf("Could not read downloaded archive %s. Got: %v", file, err)
	}
	defer reader.Close()
	target := resolveTargetFilename()
	expected := resolveFilenameInArchive(version)
	for _, candidate := range reader.File {
		if candidate.Name == expected {
			extractBinaryFileFromZipPart(candidate, target)
			return target
		}
	}
	log.Fatalf("Could not find in downloaded file %s the expected file %s.", file, expected)
	return ""
}

func extractBinaryFileFromZipPart(source *zip.File, target string) {
	input, err := source.Open()
	if err != nil {
		log.Fatalf("Could not open file %s in zip file. Got: %v", source.Name, err)
	}
	defer input.Close()
	output, err := os.Create(target)
	if err != nil {
		log.Fatalf("Could not create file %s. Got: %v", target, err)
	}
	defer output.Close()
	io.Copy(output, input)
	if runtime.GOOS != "windows" {
		err = output.Chmod(0755)
		if err != nil {
			log.Fatalf("Could not make %s executable. Got: %v", target, err)
		}
	}
}

func downloadToTemporaryLocation(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Could not download %s. Got: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		io.Copy(buf, resp.Body)
		log.Fatalf("Got illegal response: %s", buf.String())
	}
	target, err := ioutil.TempFile("", "hugo" + filepath.Ext(url) + ".")
	if err != nil {
		log.Fatalf("Could not create temporary file. Got: %v", err)
	}
	defer target.Close()
	io.Copy(target, resp.Body)
	return target.Name()
}

func thisExecutable() string {
	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		return os.Args[0]
	}
	return path
}

func resolveTargetFilename() string {
	dir := filepath.Dir(thisExecutable())
	if runtime.GOOS == "windows" {
		return filepath.Join(dir, "hugo.exe")
	}
	return filepath.Join(dir, "hugo")
}

func resolveTargetVersionFilename() string {
	dir := filepath.Dir(thisExecutable())
	return filepath.Join(dir, "hugo.version")
}

func resolveDownloadUrl(release Release) string {
	filenameSuffix := resolveDownloadFilename(release.Version)
	for _, asset := range release.Assets {
		if strings.HasPrefix(asset.Name, "hugo_") && strings.HasSuffix(asset.Name, filenameSuffix) {
			return asset.BrowserDownloadUrl
		}
	}
	log.Fatalf("Cannot find a suitable download url for latest release: %v", release.Name)
	return ""
}

func resolveDownloadFilename(version string) string {
	return fmt.Sprintf("hugo_%s_%s-%s.%s",
		version,
		resolveOperatingSystemDownloadFilenamePart(),
		resolveArchitectureDownloadFilenamePart(),
		resolveArchitectureDownloadFilenameExtension(),
	)
}

func resolveOperatingSystemDownloadFilenamePart() string {
	switch runtime.GOOS {
	case "windows":
		return "Windows"
	case "linux":
		return "Linux"
	case "darwin":
		return "macOS"
	}
	log.Fatalf("Unsupported operating system: %s", runtime.GOOS)
	return ""
}

func resolveArchitectureDownloadFilenamePart() string {
	switch runtime.GOARCH {
	case "386":
		return "32bit"
	case "amd64":
		return "64bit"
	}
	log.Fatalf("Unsupported architecture: %s", runtime.GOARCH)
	return ""
}

func resolveArchitectureDownloadFilenameExtension() string {
	switch runtime.GOOS {
	case "windows":
		return "zip"
	case "linux":
		return "tar.gz"
	case "darwin":
		return "zip"
	}
	log.Fatalf("Unsupported architecture: %s", runtime.GOARCH)
	return ""
}

func resolveFilenameInArchive(version string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("hugo_%s_windows_%s.exe", version, runtime.GOARCH)
	}
	return fmt.Sprintf("hugo_%s_%s_%s/hugo_%s_%s_%s", version, runtime.GOOS, runtime.GOARCH, version, runtime.GOOS, runtime.GOARCH)
}

func getLatestRelease() Release {
	url := "https://api.github.com/repos/spf13/hugo/releases/latest"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Could not download %s. Got: %v", url, err)
	}
	token := os.Getenv("GITHUB_TOKEN")
	if len(token) > 0 {
		req.SetBasicAuth("", token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Could not execute request to download latest release information from GitHub. Got: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		io.Copy(buf, resp.Body)
		log.Fatalf("Got illegal response: %s", buf.String())
	}
	decoder := json.NewDecoder(resp.Body)
	result := Release{}
	err = decoder.Decode(&result)
	if err != nil {
		log.Fatalf("Could not decode the response from GitHub. Got: %v", err)
	}
	if !strings.HasPrefix(result.Name, "v") || len(result.Name) < 2 {
		log.Fatalf("Unsupported release name: %s", result.Name)
	}
	result.Version = result.Name[1:]
	return result
}

type Release struct {
	Id      int64 `json:"id"`
	Name    string `json:"name"`
	Version string
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Id                 int64 `json:"id"`
	Name               string `json:"name"`
	BrowserDownloadUrl string `json:"browser_download_url"`
}
