package main

func main() {
	config := InitConfig()
	application := InitializeApplication(config.ApplicationConfig)
	server := InitializeServer(config.HTTPConfig, application)
	server.Run()
}
