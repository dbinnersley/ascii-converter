package main

import (
	"image"
	"io"
)

//Application The implementation of the actual application to be run
type Application struct {
	config *ApplicationConfig
}

//InitializeApplication Initializes the application with a given config
func InitializeApplication(config *ApplicationConfig) *Application {
	application := &Application{config}
	return application
}

//ProcessData Process the bytes sent in to the application
func (app *Application) ProcessData(reader io.Reader) ([]string, error) {
	image, _, err := image.Decode(reader)

	if err != nil {
		return []string{}, err
	}

	asciiArt := ImageToASCII(image, app.config.OutputWidth, app.config.CharSet)
	lines := asciiArt.CreateLines()
	return lines, nil
}
