package deployment

import "strings"

func parseVersionFromImage(image string) string {
    parts := strings.Split(image, ":")
    if len(parts) > 1 {
        return parts[len(parts)-1]
    }
    return "latest"
}
