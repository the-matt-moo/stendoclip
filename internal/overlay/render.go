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

	// Select custom font or fall back to stock.
	var previousFont winapi.HGDIOBJ
	if c.font != 0 {
		previousFont = winapi.SelectObject(dc, c.font)
	} else {
		previousFont = winapi.SelectObject(dc, winapi.GetStockObject(winapi.DefaultGUIFont))
	}
	defer winapi.SelectObject(dc, previousFont)

	entry := c.stack.Get(c.index)
	footerHeight := int32(30)
	padX := int32(24)
	padY := int32(12)

	// Measure body text height with DT_CALCRECT.
	availableWidth := client.Right - 2*padX
	availableHeight := client.Bottom - 2*padY - footerHeight
	measure := winapi.Rect{Left: 0, Top: 0, Right: availableWidth, Bottom: 10000}
	_ = winapi.DrawText(dc, entry.Text, &measure, winapi.DTWordBreak|winapi.DTNoPrefix|winapi.DTCalcRect)
	textHeight := measure.Bottom
	if textHeight > availableHeight {
		textHeight = availableHeight
	}

	// Center text vertically in the body area.
	bodyTop := padY + (availableHeight-textHeight)/2
	body := winapi.Rect{Left: padX, Top: bodyTop, Right: client.Right - padX, Bottom: bodyTop + textHeight}
	c.report(winapi.DrawText(dc, entry.Text, &body, winapi.DTWordBreak|winapi.DTNoPrefix|winapi.DTEndEllipsis))

	// Footer: index/total, right-aligned at bottom.
	status := fmt.Sprintf("%d / %d", c.index+1, c.stack.Len())
	if entry.Pinned {
		status = "PINNED    " + status
	}
	footer := winapi.Rect{Left: padX, Top: client.Bottom - padY - footerHeight, Right: client.Right - padX, Bottom: client.Bottom - padY}
	c.report(winapi.DrawText(dc, status, &footer, winapi.DTRight|winapi.DTBottom|winapi.DTSingleLine|winapi.DTNoPrefix))
}
