package multiple

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
)

var supportedFormats = []string{"jpg", "jpeg", "png", "webp"}

// Check if supported format
func isSupportedFormat(extension string) bool {
	for _, f := range supportedFormats {
		if f == extension {
			return true
		}
	}

	return false
}

// Get user input
func getInput(prompt string) string {
	fmt.Println(prompt)
	fmt.Print("> ")

	var input string
	fmt.Scanln(&input)

	return strings.TrimSpace(input)
}

// Get Y/N input
func askYesNo(prompt string) bool {
	fmt.Println(prompt)
	for {
		fmt.Print("> ")

		var input string
		fmt.Scanln(&input)
		input = strings.ToLower(strings.TrimSpace(input))

		switch input {
		case "y":
			return true
		case "n":
			return false
		}

		fmt.Println("\nInvalid Input. Please enter Y or N.")
	}
}

// Read image file
func decodeImage(path string) (image.Image, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	extension := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))

	switch extension {
	case "jpg", "jpeg":
		return jpeg.Decode(f)
	case "png":
		return png.Decode(f)
	case "webp":
		return webp.Decode(f)
	default:
		return nil, fmt.Errorf("unsupported format: %s", extension)
	}
}

// Save image file
func encodeImage(img image.Image, path string, format string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	switch format {
	case "jpg", "jpeg":
		// Convert to RGB as JPEG does not support transparency
		rgbImg := image.NewRGBA(img.Bounds())

		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				rgbImg.Set(x, y, img.At(x, y))
			}
		}

		return jpeg.Encode(f, rgbImg, nil)
	case "png":
		return png.Encode(f, img)
	case "webp":
		return webp.Encode(f, img, &webp.Options{Lossless: false, Quality: 80})
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// Convert all images in folder
func RunConvert() {
	inputFolder := getInput("Enter the path of the folder where original images are saved.")

	if _, err := os.Stat(inputFolder); os.IsNotExist(err) {
		fmt.Printf("\nCannot find the folder from: %s\n", inputFolder)

		return
	}

	outputFolder := getInput("\nEnter the output folder where converted images will be saved.")

	if outputFolder == "" || outputFolder == "/" {
		fmt.Println("\nCannot save converted images to the root directory.")

		return
	}

	// Get target format
	var targetFormat string

	fmt.Println("\nEnter the target image format.")
	fmt.Println("Enter 'Formats' to see supported formats.")

	for {
		fmt.Print("> ")

		var input string
		fmt.Scanln(&input)
		input = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(input, ".")))

		if input == "formats" {
			fmt.Println("\n[ Supported Formats ]")
			fmt.Println("[ JPG, JPEG, PNG, WEBP ]")

			continue
		}

		if !isSupportedFormat(input) {
			fmt.Printf("\nInvalid command or Unsupported format: %s\n", input)

			continue
		}

		targetFormat = input

		break
	}

	useSubfolder := askYesNo("\nDo you want to create a subfolder? (Y/N)")

	if useSubfolder {
		outputFolder = filepath.Join(outputFolder, "Converted Images")
	}

	os.MkdirAll(outputFolder, os.ModePerm)
	fmt.Println("\n" + strings.Repeat("-", 30) + "\n")

	// Collect converted original format
	convertedExtensions := map[string]bool{}
	entries, _ := os.ReadDir(inputFolder)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		extension := strings.ToLower(strings.TrimPrefix(filepath.Ext(name), "."))

		if !isSupportedFormat(extension) {
			continue
		}

		// Skip if same format
		if extension == targetFormat || (targetFormat == "jpeg" && extension == "jpg") || (targetFormat == "jpg" && extension == "jpeg") {
			continue
		}

		inputPath := filepath.Join(inputFolder, name)
		baseName := strings.TrimSuffix(name, filepath.Ext(name))

		if !useSubfolder && askYesNo(fmt.Sprintf("\nDo you want to overwrite the original image '%s'? (Y/N)", name)) {
			// Override: Convert to temporary file and replace
			temporaryPath := filepath.Join(outputFolder, baseName+"_tmp."+targetFormat)
			img, err := decodeImage(inputPath)

			if err != nil {
				fmt.Printf("Error reading %s: %v\n", name, err)
				continue
			}

			if err := encodeImage(img, temporaryPath, targetFormat); err != nil {
				fmt.Printf("Error converting %s: %v\n", name, err)
				continue
			}

			os.Remove(inputPath)
			os.Rename(temporaryPath, filepath.Join(outputFolder, baseName+"."+targetFormat))
		} else {
			outputPath := filepath.Join(outputFolder, baseName+"."+targetFormat)
			img, err := decodeImage(inputPath)

			if err != nil {
				fmt.Printf("Error reading %s: %v\n", name, err)
				continue
			}

			if err := encodeImage(img, outputPath, targetFormat); err != nil {
				fmt.Printf("Error converting %s: %v\n", name, err)
				continue
			}
		}

		convertedExtensions[strings.ToUpper(extension)] = true
	}

	if len(convertedExtensions) == 0 {
		fmt.Println("No images were converted.")
		return
	}

	extensions := []string{}

	for e := range convertedExtensions {
		extensions = append(extensions, e)
	}

	fmt.Printf("\nSuccessfully Converted!\n%s --> %s\n", strings.Join(extensions, ", "), strings.ToUpper(targetFormat))
}
