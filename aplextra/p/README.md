# p - plot package

This package is an interface to ktye/plot.

The plot function `p→p` or `⌼` takes an array on the right and returns a Plot or a PlotArray.

A Plot is an interactive object that can be manipulated in a gui, such as zoom, pan, or select.
It's properties can be changed like a dictionary.

In a non-gui application, the default stringer formats a Plot as a [sixel](https://en.wikipedia.org/wiki/Sixel) encoded image,
such that it can be directly dumped into the terminal.

There are basically two methods to build plots:
- directly by passing numeric arrays
- piecewise by adding Line object to a Plot

Plots can be arranged in a grid by reshaping a PlotArray.

## Plot numeric arrays
- right argument R (numeric array) returns a Plot or PlotArray
	- rank 1: one plot with a single line (dataseries)
	- rank 2: one plot with multiple lines (one per major cell)
	- rank 3: multiple plots shown side by side (can be reshaped)
- left argument L (numeric array, x-axis)
	- last axis must conform to last axis of R
	- other axis may conform, otherwise the first index is used
	- smaller ranks are extended with leading ones
- left argument L (dict)
	- fields are translated to Plot properties:
		- Type, Xlabel, Ylabel, [XYZ]unit (all strings)
		- [XYZ]{min,max} (float)
		- field X is used as X axis (float vector)
	- example: `("Title" "Xmax" "Ymax"#("A B C";8;5;))⌼?2 3⍴5`
- monadic
	- default value: `L←⍳ ¯1↑R`
- plot type is inferred from the data type
	- if any value in `¯2↑R` is complex: "polar" (monadic), "ampang" (dyadic)
	- otherwise "xy"

## Build Plots sequentially
- initialize empty Plot: `P←⌼⍳0`
- set plot properties like a dict: `P["Type" "Title"]←("polar";"An example plot";)`
- add Line object:
	- `L←p→l 0 ⋄ L["X" "Y"]←(⍳10;10-⍳10;) ⋄ L+P`
- see [ktye/plot](https://github.com/ktye/plot/blob/master/plot.go) for a description of the Plot and Line objects	

	
## Default style
The default plot style is controlled package wide and can be changed by set functions:
- `p→dark 0` (bool): dark background
- `p→transparent 1` (bool): transparent background
- `p→colors 0xFF0000 0x00FF00 0x0000FF` 
	- (int vector): line colors (cyclic), ints 0xRRGGBB, empty: reset
- `p→fontsizes 12 8: (2 floats) font sizes for labels and axis tic labels
- `p→size 600 300` (2 int vector) default plot image size in pixels
- `p→gui 1` (bool): don't use sixel stringer

## Examples
```
	lui -a test.apl
	lui -i '⌼?3 10⍴10'           # one xy-plot with 3 lines
	lui -i '2 2⍴⌼?3 3 10⍴1J1'    # reshape 3 polar plots to 2x2 with one empty image
	lui -i '400 400⌼⌼?3 10⍴1J1'  # explicitly convert to an image with a size of 400x400
```

## TODO
- animations (frames within the plot object)
- animations send plots over a channel
- plot at a low rate using a monitor channel, CLT example
- raster images, spectrogram
- histogram (bar plot style)
- convert/add plots to png, pptx, html

