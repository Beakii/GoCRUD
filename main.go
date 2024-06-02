package main

func main() {
	server := NewAPIServer(":80")
	server.Run()
}
