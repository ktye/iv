# cmd/lui - a gui front-end
<p align="center" >
  <img width="760" height="300" src="aplui.png"><br/>
</p>

Codename: *ludwig der Kurze*

Lui is a gui front-end to APL\iv

It uses the a [duit](https://github.com/ktye/duit) widget 
on top of [shiny](https://golang.org/x/exp/shiny) and embedds the APL385 Unicode font.
Keystrokes are translated automatically and no special keyboard
driver is needed.

When pressing the ENTER key, the current line is interpreted
and the result is appended to the end of the editor.
Otherwise, it is a normal text edit widget.

Lui displays image values on the top left corner over the text.
The image disappears at the next input event. *TODO*

Aplui builds as a single binary.
On windows, build with: `go build -ldflags -H=windowsgui` to prevent the console.

The slightly modified GNU APL keyboard layout visible on startup is for the learning phase.
It will eventually disappear.