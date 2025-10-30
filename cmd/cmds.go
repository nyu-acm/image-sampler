package cmd

import (
	"fmt"
	"os"

	pkg "github.com/nyu-acm/iso_sampler/pkg"
	"github.com/spf13/cobra"
)

var (
	dir            string
	dirLimit       int
	img            string
	exportLocation string
)

func init() {
	processDirectoryCmd.Flags().StringVarP(&dir, "directory", "d", "", "Path to the directory containing ISO images")
	processDirectoryCmd.Flags().IntVarP(&dirLimit, "limit", "l", 10, "Maximum number of directories to sample from per image")
	processDirectoryCmd.Flags().StringVarP(&exportLocation, "out", "o", "exports", "Location to export sampled files")
	rootCmd.AddCommand(processDirectoryCmd)
	sampleImageCmd.Flags().StringVarP(&img, "image", "i", "", "Path to the ISO image")
	sampleImageCmd.Flags().IntVarP(&dirLimit, "limit", "l", 10, "Maximum number of directories to sample from")
	sampleImageCmd.Flags().StringVarP(&exportLocation, "out", "o", "exports", "Location to export sampled files")
	rootCmd.AddCommand(sampleImageCmd)
}

var rootCmd = &cobra.Command{
	Use:   "iso_sampler",
	Short: "A tool to sample files from ISO images",
	Long:  `iso_sampler is a command-line tool that allows users to sample files from ISO images based on specified criteria.`,
}

var sampleImageCmd = &cobra.Command{
	Use:   "sample-image",
	Short: "Sample files from an ISO image",
	Long:  `Sample files from an ISO image based on specified criteria such as directory limit and export location.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := pkg.ProcessImage(img, dirLimit, exportLocation); err != nil {
			fmt.Println("Error processing image:", err)
		}
	},
}

var processDirectoryCmd = &cobra.Command{
	Use:   "process-directory",
	Short: "Process a directory",
	Long:  `Process a directory of images by sampling files based on specified criteria.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := pkg.ProcessDirectory(dir, dirLimit, exportLocation); err != nil {
			fmt.Println("Error processing directory:", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
