// Package cmd defines and implements command-line commands and flags
// used by fdio. Commands and flags are implemented using Cobra.
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/retgits/fdio/util"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Performs a backup (upload) or restore (download) of a database file from Amazon S3",
	Run:   runBackup,
}

// Flags
var (
	bucket  string
	region  string
	restore bool
)

// init registers the command and flags
func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().StringVar(&bucket, "bucket", "", "The name of the S3 bucket (required)")
	backupCmd.Flags().StringVar(&region, "region", "us-west-2", "The region of the S3 bucket (defaults to us-west-2)")
	backupCmd.Flags().BoolVar(&restore, "restore", false, "Restore (download) the file from Amazon S3")
	backupCmd.Flags().BoolVar(&overwrite, "overwrite", false, "When restoring a file, should any existing file be overwritten (defaults to false)")
}

// runBackup is the actual execution of the command
func runBackup(cmd *cobra.Command, args []string) {
	// Create a new session without AWS credentials. This means it will use the default credentials
	// stored on the machine
	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	folder := filepath.Dir(dbFile)
	filename := filepath.Base(dbFile)
	if restore {
		log.Println("Performing restore from Amazon S3")
		log.Printf("filename has been set to: %s \nS3 Bucket has been set to: %s \n", filename, bucket)
		// Check if file exists
		_, err := os.Stat(filename)
		if err == nil && !overwrite {
			log.Fatalf("The file %s already exists and no overwrite flag was specified\n", filename)
		} else if err != nil {
			log.Fatal(err.Error())
		}
		// Download file
		err = util.DownloadFile(awsSession, folder, filename, bucket)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		log.Println("Performing backup to Amazon S3")
		log.Printf("filename has been set to: %s \nS3 Bucket has been set to: %s \n", filename, bucket)
		// Check if file exists
		_, err := os.Stat(filename)
		if err != nil {
			log.Fatalf("Error finding file %s: %s\n", filename, err.Error())
		}
		// Upload file
		err = util.UploadFile(awsSession, folder, filename, bucket)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
