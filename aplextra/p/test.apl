⍝ lui -a test.apl

⍝ Directly plot numeric arrays xy, polar and ampang
⌼?3 10 ⍴10                ⍝ xy-plot with 3 lines
(.1×⍳10)⌼?3 10⍴10         ⍝ xy-plot with a shares x axis
⌼1a30+?3 100⍴1J1          ⍝ polar plot
(⍳100)⌼?3 100⍴1J1          ⍝ ampang plot (amplitude and angle over x)

⍝ Arrays of plots
⌼?2 3 10⍴1J1              ⍝ two polar diagrams side-by-side
2 2⍴⌼?3 3 10⍴1J1          ⍝ reshape 3 polar plots to 2x2 with one empty image

⍝ Default styles
p→dark 0                   ⍝ white background
p→size 600 400             ⍝ image size
p→fontsizes 18 18          ⍝ label and tic font sizes
p→colors 0xFF0000 0xFF00 0 ⍝ line colors
⌼?2 3 10⍴10
p→dark 1
p→size 800 400 
p→fontsizes 18 12
p→colors ⍳0

⍝ Axis labels
P←⌼?3 10⍴10
P[`Title`Xlabel`Ylabel]←`XY-Plot`X-Axis`Y-axis
P

⍝ Plot array has reference semantics
P←⌼?2 3 10⍴10
P1←P[1]
P2←P[2]
P1[`Title]←"left plot"
P2[`Title]←"right plot"
P

⍝ Directly plot a numeric array, with property dict on the left.
(`Title`Xlabel`Ylabel`Xmax#"Example plot" `x-axis`Y-Axis 3)⌼?2 3 10⍴3


⍝ Build plot from line objects
P←⌼⍳0
P[`Type]←`xy
L1←p→l 0
L1[`X`Y]←(⍳10;10-⍳10;)
L1[`Style;`Line;`Width]←5
L1[`Style;`Marker;`Size]←8
L2←p→l 1
L2[`X`Y]←(⍳10;⍳10;)
P←L1+L2+P
P[`Title`Xlabel`Ylabel`Xunit`Yunit]←`Title`X-axis`Y-axis`km`km/h
P