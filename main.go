package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

func main() {

	if len(os.Args) == 1 || os.Args[1] == "" {
		fmt.Println("Please use the application this way: nickflix-upload <movie-file-name> <subtitle-file-name (optional)>")
		return
	}

	// FTP server details
	ftpHost := "192.168.0.46"
	ftpPort := "21"
	ftpUser := "ftpuser"
	ftpPassword := "051205"

	var hasSubtitle bool = len(os.Args) >= 3 && os.Args[2] != ""

	currentDir, _ := os.Getwd()

	movieFileName := os.Args[1]
	noSpacesMovieFileName := strings.ReplaceAll(movieFileName, " ", "-")

	localMovieFilePath := filepath.Join(currentDir, movieFileName)
	remoteMovieFilePath := "/files/movies/" + noSpacesMovieFileName

	fmt.Printf("Starting reading movie file...\n\n")
	movieFileContents, err := os.ReadFile(localMovieFilePath)
	if err != nil {
		fmt.Println("Error reading movie local file:", err)
		return
	}

	// Connect to FTP server
	ftpClient, err := ftp.Dial(ftpHost+":"+ftpPort, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		fmt.Println("Error connecting to FTP server:", err)
		return
	}
	defer ftpClient.Quit()

	// Login to FTP server
	err = ftpClient.Login(ftpUser, ftpPassword)
	if err != nil {
		fmt.Println("Error logging in to FTP server:", err)
		return
	}

	fmt.Printf("Starting Movie upload. This might take a while...\n\n")

	err = ftpClient.Stor(remoteMovieFilePath, bytes.NewReader(movieFileContents))
	if err != nil {
		fmt.Println("Error uploading file to FTP server:", err)
		return
	}

	fmt.Printf("Movie uploaded successfully to FTP server!\n\n")
	fmt.Printf("Link: ../secret-movies/%s\n\n", noSpacesMovieFileName)

	if hasSubtitle {
		subtitleFileName := os.Args[2]
		noSpacesSubtitleFileName := strings.ReplaceAll(subtitleFileName, " ", "-")
		remoteSubtitleFilePath := "/files/movies/subtitles/" + noSpacesSubtitleFileName
		localSubtitleFilePath := filepath.Join(currentDir, subtitleFileName)
		fmt.Printf("Starting reading subtitle file...\n\n")
		subtitleFileContents, err := os.ReadFile(localSubtitleFilePath)
		if err != nil {
			fmt.Println("Error reading subtitle local file:", err)
			return
		}

		fmt.Printf("Starting subtitle upload. This might take a while\n\n")
		err = ftpClient.Stor(remoteSubtitleFilePath, bytes.NewReader(subtitleFileContents))

		if err != nil {
			fmt.Println("Error uploading subtitle file to FTP server:", err)
			return
		}

		fmt.Printf("Subtitle uploaded successfully to FTP server!\n\n")
		fmt.Printf("Link: ../secret-movies/subtitles/%s", noSpacesMovieFileName)
	}

}
