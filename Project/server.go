package main






func main() {
  ln, err := net.Listen("tcp","127.0.0.1:30000")
  if err != nil {
    fmt.Println("Error")
    os.Exit(1)
  }
  defer ln.close

  


}
