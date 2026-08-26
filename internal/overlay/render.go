package overlay

import (
	"fmt"

	"github.com/mooreceipts/stendoclip/internal/winapi"
)

func (c *Controller) paint() {
	var paint winapi.PaintStruct
	dc := winapi.BeginPaint(c.hwnd, &paint)
	if dc == 0 {
		return
	}
	defer winapi.EndPaint(c.hwnd, &paint)

	client, err := winapi.GetClientRect(c.hwnd)
	if err != nil {
		c.report(err)
		return
	}
	winapi.FillRect(dc, &client, c.brush)
	winapi.SetBkMode(dc, winapi.Transparent)
	winapi.SetTextColor(dc, winapi.RGB(245, 245, 247))
	previousFont := winapi.SelectObject(dc, winapi.GetStockObject(winapi.DefaultGUIFont))
	defer winapi.SelectObject(dc, previousFont)

	entry := c.stack.Get(c.index)
	body := winapi.Rect{Left: 24, Top: 20, Right: client.Right - 24, Bottom: client.Bottom - 48}
	c.report(winapi.DrawText(dc, entry.Text, &body, winapi.DTWordBreak|winapi.DTNoPrefix|winapi.DTEndEllipsis))

	status := fmt.Sprintf("%d / %d", c.index+1, c.stack.Len())
	if entry.Pinned {
		status = "PINNED    " + status
	}
	footer := winapi.Rect{Left: 24, Top: client.Bottom - 42, Right: client.Right - 24, Bottom: client.Bottom - 16}
	c.report(winapi.DrawText(dc, status, &footer, winapi.DTRight|winapi.DTBottom|winapi.DTSingleLine|winapi.DTNoPrefix))
}
