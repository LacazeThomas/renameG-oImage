package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/codingsince1985/geo-golang/frenchapigouv"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func scan_recursive(dir_path string, ignore []string) []string {

	folders := []string{}

	// Scan
	filepath.Walk(dir_path, func(path string, f os.FileInfo, err error) error {
		_continue := false
		for _, i := range ignore {
			if strings.Index(path, i) != -1 {
				_continue = true
			}
		}
		if _continue == false {
			f, err = os.Stat(path)
			CheckErr(err)
			f_mode := f.Mode()
			if f_mode.IsDir() {
				folders = append(folders, path)
			}
		}
		return nil
	})
	return folders
}

func main() {

	folders := scan_recursive("./", []string{"./"})

	// Folders
	for _, folder := range folders {
		realFolder := "./" + folder + "/"

		files, err := ioutil.ReadDir(realFolder)
		CheckErr(err)

		for i, file := range files {
			log.Println(i, file.Name())
			if strings.Contains(file.Name(), ".jpg") {

				f, err := os.Open(realFolder + file.Name())
				CheckErr(err)

				exif.RegisterParsers(mknote.All...)

				x, err := exif.Decode(f)
				CheckErr(err)

				lat, long, _ := x.LatLong()

				newName := geocoder(lat, long, f.Name())
				log.Println(file.Name())
				renameImage(file.Name(), newName, realFolder)
			}
		}
	}
}

func renameImage(fname, newName, realFolder string) {
	err := os.Rename(realFolder+fname, realFolder+newName+"_"+fname)
	CheckErr(err)
}

func geocoder(lat float64, long float64, originalName string) string {
	address, _ := frenchapigouv.Geocoder().ReverseGeocode(lat, long)

	returnName := originalName
	if address != nil {
		//fmt.Printf("Detailed address: %#v\n", address)
		returnName = fmt.Sprintf("%s %s %s %g %g", address.City, address.HouseNumber, address.Street, lat, long)
	}
	return returnName
}

//CheckErr print error
func CheckErr(err error) {
	if err != nil {
		log.Print(err)
	}
}
