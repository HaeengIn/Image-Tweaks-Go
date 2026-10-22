package multiple

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"
)

// 지원 확장자 확인
func isSupportedImage(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp"
}

// 입력 폴더 경로 입력 받기
func getInputFolderPath(reader *bufio.Reader) string {
	for {
		fmt.Println("Enter the path of folder where original images are saved")
		fmt.Print("> ")
		path, _ := reader.ReadString('\n')
		path = strings.TrimSpace(strings.Trim(path, `"`))

		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			return path
		}
		fmt.Printf("Invalid Input. Cannot find the folder: %s\n", path)
	}
}

// 출력 폴더 경로 입력 받기
func getOutputFolderPath(reader *bufio.Reader) string {
	fmt.Println("\nEnter the path of folder where optimized images will be saved")
	fmt.Print("> ")
	path, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.Trim(path, `"`))
}

// 최적화 품질 입력 받기 (0~100)
func getOptimizationQuality(reader *bufio.Reader) int {
	for {
		fmt.Println("\nEnter the optimization quality (0 ~ 100)")
		fmt.Print("> ")
		var quality int
		_, err := fmt.Fscan(reader, &quality)
		reader.ReadString('\n') // 버퍼 비우기
		if err == nil && quality >= 0 && quality <= 100 {
			return quality
		}
		fmt.Println("Invalid Input. Please enter an integer between 0 and 100.")
	}
}

// 폴더 내 이미지 총 용량 계산
func calculateTotalSize(folderPath string) int64 {
	var total int64
	entries, _ := os.ReadDir(folderPath)
	for _, entry := range entries {
		if !entry.IsDir() && isSupportedImage(entry.Name()) {
			info, err := entry.Info()
			if err == nil {
				total += info.Size()
			}
		}
	}
	return total
}

// 이미지 파일 열기 (webp 포함)
func openImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".webp" {
		return webp.Decode(f)
	}
	img, _, err := image.Decode(f)
	return img, err
}

// 이미지 최적화 후 저장
func optimizeImage(inputPath, outputPath string, quality int) error {
	img, err := openImage(inputPath)
	if err != nil {
		return err
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	ext := strings.ToLower(filepath.Ext(outputPath))
	switch ext {
	case ".jpg", ".jpeg":
		// JPEG: quality 그대로 사용
		return jpeg.Encode(out, img, &jpeg.Options{Quality: quality})
	case ".png":
		// PNG: 무손실 압축 (DefaultCompression 사용)
		encoder := png.Encoder{CompressionLevel: png.DefaultCompression}
		return encoder.Encode(out, img)
	case ".webp":
		// WebP: 인코딩 미지원, JPEG로 대체 저장
		fmt.Printf("  [경고] WebP 인코딩은 지원되지 않아 JPEG로 저장합니다: %s\n", filepath.Base(outputPath))
		newPath := strings.TrimSuffix(outputPath, ext) + ".jpg"
		jpgOut, err := os.Create(newPath)
		if err != nil {
			return err
		}
		defer jpgOut.Close()
		return jpeg.Encode(jpgOut, img, &jpeg.Options{Quality: quality})
	default:
		return fmt.Errorf("지원하지 않는 포맷: %s", ext)
	}
}

// Multiple Optimize 실행
func RunOptimize() {
	reader := bufio.NewReader(os.Stdin)

	inputFolder := getInputFolderPath(reader)
	outputFolder := getOutputFolderPath(reader)
	quality := getOptimizationQuality(reader)

	// 입출력 폴더가 동일한지 확인
	absInput, _ := filepath.Abs(inputFolder)
	absOutput, _ := filepath.Abs(outputFolder)
	sameFolder := absInput == absOutput

	var useSubfolder, deleteOriginals bool
	saveFolder := outputFolder

	if sameFolder {
		useSubfolder = askYesNo("\nDo you want to create a subfolder? (Y/N)")
		if useSubfolder {
			saveFolder = filepath.Join(inputFolder, "Optimized Images")
		} else {
			deleteOriginals = askYesNo("\nDo you want to delete original images? (Y/N)")
			saveFolder = inputFolder
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println()

	// 출력 폴더 생성
	if !sameFolder {
		os.MkdirAll(outputFolder, 0755)
		saveFolder = outputFolder
	} else if useSubfolder {
		os.MkdirAll(saveFolder, 0755)
	}

	originalSize := calculateTotalSize(inputFolder)

	// 이미지 처리
	entries, _ := os.ReadDir(inputFolder)
	for _, entry := range entries {
		if entry.IsDir() || !isSupportedImage(entry.Name()) {
			continue
		}

		inputFilePath := filepath.Join(inputFolder, entry.Name())
		var outputFilePath string

		if sameFolder {
			if useSubfolder {
				outputFilePath = filepath.Join(saveFolder, entry.Name())
			} else if deleteOriginals {
				outputFilePath = filepath.Join(saveFolder, entry.Name())

				img, err := openImage(inputFilePath)
				if err != nil {
					fmt.Printf("열기 실패: %s (%v)\n", entry.Name(), err)
					continue
				}
				os.Remove(inputFilePath)

				out, err := os.Create(outputFilePath)
				if err != nil {
					fmt.Printf("저장 실패: %s (%v)\n", entry.Name(), err)
					continue
				}
				ext := strings.ToLower(filepath.Ext(outputFilePath))
				if ext == ".jpg" || ext == ".jpeg" {
					jpeg.Encode(out, img, &jpeg.Options{Quality: quality})
				} else if ext == ".png" {
					png.Encode(out, img)
				}
				out.Close()
				fmt.Printf("완료: %s\n", entry.Name())
				continue
			} else {
				// 같은 폴더, 삭제 안 함 → _optimized 접미사 추가
				ext := filepath.Ext(entry.Name())
				nameWithout := strings.TrimSuffix(entry.Name(), ext)
				outputFilePath = filepath.Join(saveFolder, nameWithout+"_optimized"+ext)
			}
		} else {
			outputFilePath = filepath.Join(saveFolder, entry.Name())
		}

		err := optimizeImage(inputFilePath, outputFilePath, quality)
		if err != nil {
			fmt.Printf("실패: %s (%v)\n", entry.Name(), err)
		} else {
			fmt.Printf("완료: %s\n", entry.Name())
		}
	}

	optimizedSize := calculateTotalSize(saveFolder)

	fmt.Println("\nSuccessfully optimized images!")
	fmt.Printf("%d bytes --> %d bytes\n", originalSize, optimizedSize)
}
