package cmd

import (
	"log"
	"time"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/file"
	"github.com/jackytck/alti-cli/gql"
	"github.com/jackytck/alti-cli/service"
	"github.com/spf13/cobra"
)

var face string

// setProfileFaceCmd represents the command of setting profile image
var setProfileFaceCmd = &cobra.Command{
	Use:   "set-face",
	Short: "Set and upload my profile image",
	Long:  "Set and upload my profile image. Maximum 20MB.",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		defer func() {
			if verbose {
				elapsed := time.Since(start)
				log.Println("Took", elapsed)
			}
		}()

		// pre-checks general
		if err := service.Check(
			nil,
			service.CheckAPIServer(),
			service.CheckIsLogin(),
			service.CheckFile(face),
		); err != nil {
			log.Println(err)
			return
		}

		// check image
		isImg, err := file.IsImageFile(face)
		errors.Must(err)
		if !isImg {
			log.Printf("%q is not an image!", face)
			return
		}
		bytes, err := file.Filesize(face)
		errors.Must(err)
		size := file.BytesToMB(bytes)
		if size > 20 {
			log.Printf("%q (with %.2fMB) is too large, max is 20MB!", face, size)
			return
		}

		// parse image into base64 string
		bs, err := file.GetBase64String(face)
		errors.Must(err)
		imgStr := "data:image/jpeg;base64," + bs

		// upload
		res, err := gql.SetProfileFace(imgStr)
		errors.Must(err)
		if res != "Success" {
			log.Println("Unknown error! Please try again later.")
			return
		}

		log.Println("Profile image is successfully set!")
	},
}

func init() {
	rootCmd.AddCommand(setProfileFaceCmd)
	setProfileFaceCmd.Flags().StringVarP(&face, "face", "f", face, "File path of profile image.")
	errors.Must(setProfileFaceCmd.MarkFlagRequired("face"))
}
