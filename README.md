# Ascii-converter 

This Ascii converter was implemented for the purpose of the take home interview provided by Planet.
The purpose of the application is to receive image uploads in a variety of formats, and respond with a json response containing the result of the image converted into Ascii characters.

### Building the application

To build the application, you are required to have go installed. Building requires running the following command in the root directory of the project 

```go build```

### Running the application
To run the application, you must currently be present in the root directory of the project and run the executable that was built in the previous step. Using the project name as the executable, all that you should need to run is

```./ascii-converter```

The application has a dependency of a config file that is present in the root directory of the project. If successful, you should receive a log message indicating that the server is up and listening on the port that is configured in the configuration file.

### Testing the application

To manually test the application, serveral test images are provided and are located under the test/ directory in the project. This directory contains images of a variety of different formats, as well as one file that is not an image. To test these files, we need to upload them to the server. The server requires the files to be uploaded as part of multipart form data. The endpoint to uploaded to is at "/upload", and the form key to use is "image". To upload the "gnome.png" sample image using default server configuration, you can run the following curl command from the root directory of the project.

```curl -F "image=@test/gnome.png" localhost:8080/upload ```

### Running unit tests

Unit tests were also written for the code responsible for doing all the image and ascii transformations. To run the unit tests, you should only be required to run:

```go test```

To run the test benchmarks, you can run the following command, which adds on to the test command

```go test -bench=.```

### Notes about the project
Currently the project only handles 4 different types of imagery: jpeg, png, gif, and tiff.
You can modify the size of the artwork, and the symbols in the art by changing the values in the application config.

#### Error Cases
Although the project does what it is supposed to and has tests surrounding the critical logic components, it is not without flaws. One of the areas that could be improved upon is input validation. Currently, there is middlewhere implemented that validates that the form data uploaded contains the correct key, but there is no validation done on the actual imagery. The errors are still caught if the data contains bad imagery, but they are returned in a 500 response instead of a 400. This is because the error is caught outside of the validation step. In order to be caught within validation, the image would have to be decoded once in validation, then once more in the processing, and this would result in a large performance decrease. 

There hasn't been any testing done on extremely large images, so there is no insight on if there are cases like this that are not handled.

#### Notes on dependencies

There were only two dependencies that were used for this project and one of them is part of the go project, but outside of the standed libary. These two dependencies are:

```
golang.org/x/image/tiff
github.com/nfnt/resize
```

The tiff library was included to allow tiff drivers to be used for reading tiff files. It allowed reading of tiff files into the standard image.Image interface used by the go standard library.

The "github.com/nfnt/resize" was a little bit more of an odd case. In order to allow some images when converted to ascii to be visible without hundreds of wrapping lines, there needed to be a limit on the number of values in each row. In order to do this, the images had to be resized to a particular size that was visible on a screen. Image resizing is a non-trivial task, since there are many different ways that images can be interpolated resulting in many different varieties of imagery with the same dimensions as a result. 
My method of creating Ascii involved the three following steps:

```Convert to Gray -> Resize -> Convert to Ascii```

As a result I only required functionality that needed to accept a grayscale image, and output a grayscale image. I looked at several different libraries, but most of them either involved converting first to RGB then resizing, or output only RGB images. "github.com/nfnt/resize" however, took in a grayscale image and output a grayscale image, without converting the format internally. Although it is stated on the repository site that it is no longer being maintained, it is a mature project with almost 2000 stars. Usually I would not choose a project with this current state, but I found it to be to best solution for the problem.