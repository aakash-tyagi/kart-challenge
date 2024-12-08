package server

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

func (s *Server) validateCoupon(ctx context.Context, coupon string, chunkSize int) error {
	// Validate coupon length
	if len(coupon) < 8 || len(coupon) > 10 {
		return fmt.Errorf("coupon code must be between 8 and 10 characters")
	}

	files := []string{"couponBase1.txt", "couponBase2.txt", "couponBase3.txt"}
	matchedFiles := make(map[string]struct{})
	var mu sync.Mutex

	var wg sync.WaitGroup

	for _, filePath := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return // Exit if context is cancelled
			default:
				s.processFile(ctx, filePath, coupon, chunkSize, matchedFiles, &mu)
			}
		}(filePath)
	}

	wg.Wait()

	if len(matchedFiles) >= 2 {
		return nil // Coupon is valid
	}

	s.Logger.Errorf("coupon code must exist in at least two files")
	return fmt.Errorf("coupon code must exist in at least two files")
}

// Updated processFile to check for context cancellation
func (s *Server) processFile(ctx context.Context, filePath string, coupon string, chunkSize int, matchedFiles map[string]struct{}, mu *sync.Mutex) {
	// Open the file
	f, err := os.Open(filePath)
	if err != nil {
		s.Logger.Errorf("Error opening file %s: %v", filePath, err)
		return
	}
	defer f.Close()

	// Get file size
	stat, err := f.Stat()
	if err != nil {
		s.Logger.Errorf("Error stating file %s: %v", filePath, err)
		return
	}
	fileSize := stat.Size()

	// Process file in chunks
	chunks := int(fileSize) / chunkSize
	if fileSize%int64(chunkSize) != 0 {
		chunks++
	}

	var chunkWG sync.WaitGroup

	// Buffer to hold the last part of the previous chunk
	var lastChunkEnd string

	for i := 0; i < chunks; i++ {
		chunkWG.Add(1)
		go func(chunkStart int64) {
			defer chunkWG.Done()
			select {
			case <-ctx.Done():
				return // Exit if context is cancelled
			default:
				// Read the current chunk
				buffer := make([]byte, chunkSize)      // Adjusted buffer size
				n, err := f.ReadAt(buffer, chunkStart) // Read into the buffer
				if err != nil && err != io.EOF {
					s.Logger.Errorf("Error reading file chunk from %s: %v", filePath, err)
					return
				}

				// Combine the last part of the previous chunk with the current chunk
				combined := lastChunkEnd + string(buffer[:n])

				// Scan the combined chunk line by line
				scanner := bufio.NewScanner(strings.NewReader(combined))
				for scanner.Scan() {
					// Check for context cancellation
					select {
					case <-ctx.Done():
						return // Exit if context is cancelled
					default:
						line := strings.TrimSpace(scanner.Text())
						if len(line) == 10 && line == coupon { // Check length and match
							mu.Lock()
							if _, exists := matchedFiles[filePath]; !exists {
								matchedFiles[filePath] = struct{}{}
							}
							mu.Unlock()
							return // Early exit on match
						}
					}
				}

				// Store the last part of the current chunk for the next iteration
				if n > 0 {
					lastChunkEnd = string(buffer[max(0, n-10):n]) // Store last 10 characters
				} else {
					lastChunkEnd = "" // Reset if no bytes were read
				}
			}
		}(int64(i * chunkSize))
	}

	chunkWG.Wait()
}

// Updated processChunk to handle EOF and last line correctly
func (s *Server) processChunk(ctx context.Context, f *os.File, chunkStart int64, chunkSize int, coupon string, matchedFiles map[string]struct{}, mu *sync.Mutex, filePath string) {
	buffer := make([]byte, chunkSize)
	n, err := f.ReadAt(buffer, chunkStart)

	if err != nil {
		if err == io.EOF {
			// Handle EOF gracefully
			if n == 0 {
				// If no bytes were read, we can return without logging an error
				return
			}
			// If we read some bytes but hit EOF, we can log that we reached the end
			s.Logger.Infof("Reached EOF for file %s after reading %d bytes", filePath, n)
		} else {
			// Log any other errors
			s.Logger.Errorf("Error reading file chunk from %s: %v", filePath, err)
		}
		return
	}

	// Scan the chunk line by line
	scanner := bufio.NewScanner(bytes.NewReader(buffer[:n]))
	for scanner.Scan() {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return // Exit if context is cancelled
		default:
			if strings.TrimSpace(scanner.Text()) == coupon {
				mu.Lock()
				if _, exists := matchedFiles[filePath]; !exists {
					matchedFiles[filePath] = struct{}{}
				}
				if len(matchedFiles) >= 2 {
					mu.Unlock()
					return
				}
				mu.Unlock()
				return
			}
		}
	}

	// Check if there was an error during scanning
	if err := scanner.Err(); err != nil {
		s.Logger.Errorf("Error scanning chunk in file %s: %v", filePath, err)
	}

	// If we reached EOF and there are still bytes left in the buffer, process them
	if n < chunkSize {
		// This means we are at the end of the file, and we need to check the last line
		if n > 0 {
			// Create a new scanner for the remaining bytes
			lastLineScanner := bufio.NewScanner(bytes.NewReader(buffer[:n]))
			for lastLineScanner.Scan() {
				fmt.Println(lastLineScanner.Text())
				if strings.TrimSpace(lastLineScanner.Text()) == coupon {
					mu.Lock()
					if _, exists := matchedFiles[filePath]; !exists {
						matchedFiles[filePath] = struct{}{}
					}
					if len(matchedFiles) >= 2 {
						mu.Unlock()
						return
					}
					mu.Unlock()
					return
				}
			}
			if err := lastLineScanner.Err(); err != nil {
				s.Logger.Errorf("Error scanning last line in file %s: %v", filePath, err)
			}
		}
	}
}
