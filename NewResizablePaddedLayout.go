/*
**
** NewResizablePaddedLayout
** Provides a resizable PaddedLayout in Fyne.
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
)


type resizablePaddedLayout struct {
  padding fyne.Size
}


func NewResizablePaddedLayout(paddingSize fyne.Size) fyne.Layout {
  return &resizablePaddedLayout{padding: paddingSize}
}

func (c *resizablePaddedLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
  newSize := fyne.NewSize(size.Width-c.padding.Width*2, size.Height-c.padding.Height*2)
  for _, o := range objects {
    o.Resize(newSize)
    o.Move(fyne.NewPos(c.padding.Width, c.padding.Height))
  }
}

func (c *resizablePaddedLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
  minSize := fyne.NewSize(0, 0)
  for _, o := range objects {
    childMin := o.MinSize()
    minSize = minSize.Max(childMin)
  }
  return minSize.Add(fyne.NewSize(c.padding.Width*2, c.padding.Height*2))
}
