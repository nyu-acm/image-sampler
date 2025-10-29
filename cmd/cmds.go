package cmd

import (
	"fmt"
	"os"

	lib "github.com/nyu-acm/iso_sampler/pkg"
	"github.com/spf13/cobra"
)

var (
	dirLimit       int
	img            string
	exportLocation string
)

func init() {
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
		// Implementation of the sampling logic goes here
		if err := lib.ProcessImage(img, dirLimit, exportLocation); err != nil {
			fmt.Println("Error processing image:", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
