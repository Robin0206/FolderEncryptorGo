Memory-Efficient and parallel FolderEncryptor which can encrypt small and large folders.<br>
Encryption is done using ChaCha20/HMAC. <br>

In addition to the encryption it also generates a table with the old filename, <br>
the new randomly generated filename, the used nonce, the used salt and the <br>
Message Authentication Code. The table also gets encrypted for obfuscation <br>
with ChaCha20-Poly1305.<br>
<br>
If the Poly1305 check fails, the program stops to prevent destroying the files <br>
with the wrong password.<br>
<br>
Key Derivation is done using PKDF2.<br>
<br>
In the future i will probably add:<br>
-Menu for tweaking the Iteration-count for PKDF2<br>
-Obfuscation of folder names (currently only filenames are obfuscated)<br>

Usage on Linux:<br>
<br>
there is a binary in the out folder.<br>
you still can compile your own one using go build.<br>

<br>
add the following line to the end of your .bashrc<br>
alias fenc='/path/to/the/binary/NameOfTheCompiledBinary ${PWD}'<br>
<br>
cd into the folder and type fenc<br>
<br>
Usage on Windows(untested):<br>
<br>
add the folder containing the compiled exe to your path-variable<br>
<br>
use it with \<Name of exe\> \<path-to-folder\><br>
