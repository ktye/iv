# p - plot package

This package is an interface to ktye/plot.

The plot function `p→p` or `⌼` takes an array on the right and returns an [image](../../apl/image.go).

## Arguments
- right argument R (numeric array)
	- rank 1: one plot with a single line
	- rank 2: one plot with multiple lines (one per major cell)
	- rank 3: multiple plots shown side by side
	- rank 4: multiple images send over a channel (animation)
- right argument R (list of dictionaries) *TODO*
	- each line object is given as a dictionary (see `ktye/plot/plot.go:/^type Line/`):
	- field `X`: float vector
	- field `Y`: float vector
	- field `C`: complex vector
	- field `Id`: int
	- field `Style`: data style numeric vector [LineWidth PointSize Color] or shorter
	- plot type is inferred by the numeric type and L
		- real: xy plot
		- complex monadic: polar diagram
		- complex dyadic: ampang (amplitude and angle over L)
- monadic:
	- default value: L←⍳ ¯1↑R
- left argument L (numeric array)
	- rank 1: single x axis must conform to last axis of R
	- rank > 1: individual x-axis conformant to R
- left argument L (dictionary) more control over the plot *TODO*
	- the following string fields can be string vectors, or are extended for multiple plots:
	- field `Type`: "xy"|"polar"|"ampang"
	- fields `Xlabel, Ylabel, Title`: axis labels
	- fields `Xunit, Yunit, Zunit`: added to the labels
	- field: `Limits`
		- numeric array: xmin xmax ymin ymax zmin zmax (missing or 0 for autoscale)
		- "equal" autoscale but equal axis for all plots
	
## Default style
Some plot style is controlled package wide and can be changed by set functions:
- `p→dark 0` (bool): dark background
- `p→transparent 1` (bool): transparent background
- `p→colors 0xFF0000 0x00FF00 0x0000FF` 
	- (int vector): line colors (cyclic), ints 0xRRGGBB, empty: reset
- `p→size 600 300`: (2 ints) width height of output image
- `p→fontsizes 12 8`: (2 ints) font sizes for labels and axis tic labels

## Example
```
	⌼?10⍴10           single line plot (xy)
	⌼?2 10⍴1J1        polar diagram with two datasets, 10 points each
	(⍳10)⌼?2 10⍴1J1    amplitude and angle over x, two lines each with 10 points
	⌼?2 3 10⍴10       two xy plots side by side, 3 lines each with 10 points per line
	⌼?10 2 3 10⍴10    animation with 10 frames, 2 plots, 3 lines 10 points each
```

## Device
The terminal version of cmd/lui uses [sixel](https://en.wikipedia.org/wiki/Sixel) as an output device.
The terminal must support it (current versions of xterm, or mintty on windows).

Example:
```
	lui -i '⌼?3 10⍴10'
```

## TODO
Interactive widget for package u, used by lui.
