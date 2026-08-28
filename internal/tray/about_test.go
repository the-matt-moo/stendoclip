package tray

import (
	"testing"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

func TestLayoutAboutFitsClient(t *testing.T) {
	for _, client := range []winapi.Rect{
		{Right: 424, Bottom: 401},
		{Right: 228, Bottom: 210},
	} {
		image, text := layoutAbout(client, 100, true)
		if image.Left < 0 || image.Top < 0 || image.Right > client.Right || image.Bottom > client.Bottom {
			t.Fatalf("image %v outside client %v", image, client)
		}
		if text.Left < 0 || text.Top < 0 || text.Right > client.Right || text.Bottom > client.Bottom {
			t.Fatalf("text %v outside client %v", text, client)
		}
		if text.Top < image.Bottom {
			t.Fatalf("text starts before image ends: image=%v text=%v", image, text)
		}
	}
}
