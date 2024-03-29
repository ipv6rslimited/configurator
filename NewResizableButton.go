/*
**
** NewResizableButton
** Provides a resizable button in Fyne.
**
** Distributed under the COOL License.
**
** Copyright (c) 2024 IPv6.rs <https://ipv6.rs>
** All Rights Reserved
**
*/

package configurator

import (
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/widget"
)


type ResizableButton struct {
  widget.Button
  width         float32
}


func NewResizableButton(label string, width float32, tapped func()) *ResizableButton {
  btn := &ResizableButton{
    width:  width,
  }
  btn.ExtendBaseWidget(btn)
  btn.Text = label
  btn.OnTapped = tapped
  return btn
}

func (b *ResizableButton) MinSize() fyne.Size {
  minSize := b.Button.MinSize()
  minSize.Width = b.width
  return minSize
}
