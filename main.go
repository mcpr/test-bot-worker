package main

import (
	"fmt"
	"os"
	"net/http"
	"github.com/fatih/color"
	"io"
    "log"
    "archive/zip"
    "path/filepath"
    "strings"
    "gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/urfave/cli"
)

var mcprAPIBase = "https://mcpr.io/api/v1"

func download(filepath, url string) (err error) {
	fmt.Println("\nDownloading...")

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println(err)
		color.Unset()
		return err
	}
	if resp.StatusCode == http.StatusNotFound {
		color.Set(color.FgRed)
		fmt.Println(resp.StatusCode)
		fmt.Println("The file you requested could not be found...")
		color.Unset()
		os.Exit(1)
		
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil

}

func createDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

func moveFile(in, out string) {
	err := os.Rename(in, out)

	if err != nil {
		os.Exit(1)
		return
	}
}

func downloadPlugin(pluginID, version string) {
	var filepath = "plugins/" + pluginID + "-" + version + ".jar"
	
	var url = mcprAPIBase + "/versions/" + pluginID + "/" + version + "/download"
	createDir("plugins")
	err := download(filepath, url)
	if err != nil {
		log.Fatalf("Download error: %s", err)
	}
}

func Unzip(src, dest string) ([]string, error) {
	fmt.Println("Extracting jar file...")
	var filenames []string
	
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()
	
	for _, f := range r.File {
	
		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()
	
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		filenames = append(filenames, fpath)
	
		if f.FileInfo().IsDir() {
	
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
	
		} else {
	
			// Make File
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}
	
			err = os.MkdirAll(fdir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
				return filenames, err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}
			defer f.Close()
	
			_, err = io.Copy(f, rc)
			if err != nil {
				return filenames, err
			}
	
		}
	}
	return filenames, nil
}

type plugin struct {
    Name string `yaml:"name"`
    Version string `yaml:"version"`
}
func (c *plugin) getInfo(filepath string) *plugin {
	
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	
	return c
}

func runTest (pluginName, pluginVersion string) {
	downloadPlugin(pluginName, pluginVersion)

	var pluginFile = pluginName + "-" + pluginVersion
	files, err := Unzip("plugins/" + pluginFile + ".jar", "tmp/" + pluginFile)
	
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Unzipped: " + strings.Join(files, "\n"))

    var p plugin
	p.getInfo("tmp/" + pluginFile + "/plugin.yml")
	
	fmt.Println("\nPlugin Name: " + p.Name)
	fmt.Println("Plugin Version: " + p.Version)
}

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "test-bot-worker"
    app.Usage = "The MCPR-Test-Bot Worker"
	app.Description = "The MCPR-Test-Bot Worker"
	app.Commands = []cli.Command{
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Run tests on plugin - test-bot-worker test [pluginID] [version]",
			Action: func(c *cli.Context) error {
				var pluginID string
				var ver string

				if c.Args().Get(0) != "" {
					pluginID = c.Args().Get(0)
				} else {
					fmt.Println("Plugin ID must be specified.")
					os.Exit(1)
				}
				if c.Args().Get(1) != "" {
					ver = c.Args().Get(1)
				} else {
					fmt.Println("Plugin version must be specified.")
					os.Exit(1)
				}
				
				runTest(pluginID, ver)
				return nil
			},
		},
	}
	app.Run(os.Args)

}