# cmd/aplui gui front end to iv/apl
<p align="center" >
  <img width="1173" height="463" src="aplui.png"><br/>
</p>

Aplui is a gui frontend to APL\iv

Aplui uses the iv/duit widget and embedds APL385 Unicode font.
Keystrokes are translated automatically and no special keyboard
driver is needed.

When pressing the ENTER key, the current line is interpreted
and the result is appended to the end of the editor.
Otherwise, it is a normal text editor.

Aplui displays image values on the top left corner over the
input text. The image disappears at the next input event.

Aplui builds as a single binary.
On windows, build with: `go build -ldflags -H=windowsgui`
