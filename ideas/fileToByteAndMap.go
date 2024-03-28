// Map must be declare in the block that it is used
// As so as function exits pointers are deallocated by garbage collector

file, err := os.Open("file.txt")
if err != nil {
	fmt.Println("Error opening file:", err)
	return
}
defer file.Close()
scanner := bufio.NewScanner(file)
var byteArray []byte
wordIndexMap := make(map[string]*int)
var index int
for scanner.Scan() {
	line := scanner.Text()
	words := strings.Fields(line)
	for _, word := range words {
		byteArray = append(byteArray, []byte(word)...)
		wordIndexMap[word] = &index
		index += len([]byte(word)) // Update the index
	}
}

if err := scanner.Err(); err != nil {
	fmt.Println("Error reading file:", err)
}
