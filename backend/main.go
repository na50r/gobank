package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/na50r/gobank/backend/sse"

	"encoding/csv"
)

var JWT_SECRET string
var CLIENT string
var ICONS string
var RECIPES string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, continuing...")
	}
	JWT_SECRET = os.Getenv("JWT_SECRET")
	CLIENT = os.Getenv("CLIENT")
	ICONS = os.Getenv("ICONS")
	RECIPES = os.Getenv("RECIPES")
}

func getImageFromFilePath(filePath string) (*Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	absPath, _ := filepath.Abs(filePath)
	name := filepath.Base(absPath)
	if !strings.HasSuffix(name, ".png") {
		return nil, nil
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	img := make([]byte, info.Size())
	_, err = f.Read(img)
	if err != nil {
		return nil, err
	}
	image := new(Image)
	image.Image = img
	image.Name = name
	return image, nil
}

func getFilePathsInDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, entry := range entries {
		if !entry.IsDir() {
			fullPath := filepath.Join(dir, entry.Name())
			paths = append(paths, fullPath)
		}
	}
	return paths, nil
}

func readImages() ([]*Image, error) {
	paths, err := getFilePathsInDir(ICONS)
	if err != nil {
		return nil, err
	}
	images := []*Image{}
	for _, path := range paths {
		image, err := getImageFromFilePath(path)
		if err != nil {
			return nil, err
		}
		if image == nil {
			continue
		}
		images = append(images, image)
	}
	return images, nil
}

func setImages(store Storage) error {
	images, err := readImages()
	if err != nil {
		return err
	}
	for _, image := range images {
		if err := store.AddImage(image.Image, image.Name); err != nil {
			return err
		}
	}
	return nil
}

func readCSV(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return records[1:], nil
}

func setElements(store Storage) error {
	records, err := readCSV(RECIPES)
	if err != nil {
		return err
	}
	for _, record := range records {
		element := new(Element)
		element.A = strings.ToLower(record[0])
		element.B = strings.ToLower(record[1])
		element.Result = strings.ToLower(record[2])
		if err := store.AddElement(element); err != nil {
			return err
		}
	}
	return nil
}


func main() {
	seed := flag.Bool("seed", false, "seed images and elements")
	flag.Parse()

	store, err := NewSQLiteStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	//./bin/gobank --seed
	if *seed {
		if err := setImages(store); err != nil {
			log.Fatal(err)
		}
		if err := setElements(store); err != nil {
			log.Fatal(err)
		}
	}

	//Accounts for ports provided by hosting services
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}

	server := NewAPIServer(":"+PORT, store, sse.NewBroker())
	server.Run()
}
