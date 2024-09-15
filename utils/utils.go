package utils

import (
    "crypto/sha256"
    "encoding/hex"
    "math/rand"
    "time"
)

// HashString returns the SHA256 hash of a given string
func HashString(input string) string {
    hash := sha256.Sum256([]byte(input))
    return hex.EncodeToString(hash[:])
}

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[seededRand.Intn(len(charset))]
    }
    return string(b)
}

// Contains checks if a string slice contains a specific string
func Contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

// RemoveDuplicates removes duplicate elements from a string slice
func RemoveDuplicates(slice []string) []string {
    keys := make(map[string]bool)
    list := []string{}
    for _, entry := range slice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}

// Cleanup performs any necessary cleanup operations
func Cleanup() {
    // Add any cleanup logic here
    // For example, closing database connections, flushing logs, etc.
}
