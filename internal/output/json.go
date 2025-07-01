package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSON outputs data as formatted JSON
func JSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	
	return nil
}