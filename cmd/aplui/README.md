# cmd/aplui - a gui front-end
<p align="center" >
  <img width="760" height="300" src="aplui.png"><br/>
</p>

Aplui is a gui front-end to APL\iv

It uses the iv/duit widget on top of shiny and embedds the APL385 Unicode font.
Keystrokes are translated automatically and no special keyboard
driver is needed.

When pressing the ENTER key, the current line is interpreted
and the result is appended to the end of the editor.
Otherwise, it is a normal text editor.

Aplui displays image values on the top left corner over the text.
The image disappears at the next input event.

Aplui builds as a single binary.
On windows, build with: `go build -ldflags -H=windowsgui`
