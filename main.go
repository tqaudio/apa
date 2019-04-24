package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	AAX  = "AAX"
	AUv2 = "AUv2"
	VST3 = "VST3"
)

func main() {
	if err := run(os.Args); err != nil {
		log.New(os.Stderr, "error: ", 0).Fatal(err)
	}
}

func run(args []string) error {
	var (
		pluginName string
		outputPath string

		auv2SourcePaths      []string
		auv2DestinationPaths []string

		vst3SourcePaths      []string
		vst3DestinationPaths []string

		sourcePaths      []string
		destinationPaths []string
	)

	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.StringVar(&pluginName, "plugin", "", "name of the plugin name")
	f.StringVar(&outputPath, "output", "output.zip", "path to output path")

	if err := f.Parse(args[1:]); err != nil {
		return err
	}
	if pluginName == "" {
		return fmt.Errorf("specify --plugin")
	}

	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	vst3RootPath := filepath.Join(workDir, "build", "VST3", "Release", pluginName+".vst3")

	if _, err := os.Stat(vst3RootPath); os.IsNotExist(err) {
		return err
	}
	if err := filepath.Walk(vst3RootPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		vst3DestinationPath := strings.TrimPrefix(path, vst3RootPath)[1:]
		vst3DestinationPath = filepath.Join(pluginName, "VST3", pluginName+".vst3", vst3DestinationPath)

		vst3SourcePaths = append(vst3SourcePaths, path)
		vst3DestinationPaths = append(vst3DestinationPaths, vst3DestinationPath)

		sourcePaths = append(sourcePaths, path)
		destinationPaths = append(destinationPaths, vst3DestinationPath)

		return nil
	}); err != nil {
		return err
	}

	auv2RootPath := filepath.Join(workDir, "build", "VST3", "Release", pluginName+"auv2.component")
	if _, err := os.Stat(auv2RootPath); err == nil {
		if err := filepath.Walk(auv2RootPath, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			auv2DestinationPath := strings.TrimPrefix(path, auv2RootPath)[1:]
			auv2DestinationPath = filepath.Join(pluginName, "Components", pluginName+".component", auv2DestinationPath)

			if base := filepath.Base(auv2DestinationPath); base == "plugin.vst3" {
				for i, vst3DestinationPath := range vst3DestinationPaths {
					vst3SourcePath := vst3SourcePaths[i]
					d := filepath.Join(
						auv2DestinationPath,
						filepath.Join(strings.Split(vst3DestinationPath, string(filepath.Separator))[3:]...),
					)
					fmt.Println("@@@d", d)
					auv2DestinationPaths = append(auv2DestinationPaths, d)
					auv2SourcePaths = append(auv2SourcePaths, vst3SourcePath)

					sourcePaths = append(sourcePaths, vst3SourcePath)
					destinationPaths = append(destinationPaths, d)
				}
				return nil
			}
			auv2SourcePaths = append(auv2SourcePaths, path)
			auv2DestinationPaths = append(auv2DestinationPaths, auv2DestinationPath)

			sourcePaths = append(sourcePaths, path)
			destinationPaths = append(destinationPaths, auv2DestinationPath)

			return nil
		}); err != nil {
			return err
		}
	}
	if err := createZip(outputPath, sourcePaths, destinationPaths); err != nil {
		return err
	}

	return nil
}
